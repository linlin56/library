package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"library/db"
	"library/i18n"
)

// LoanPreview holds the information displayed before confirmation.
type LoanPreview struct {
	Title          string  `json:"title"`
	Barcode        string  `json:"barcode"`
	Deposit        float64 `json:"deposit"`
	CurrentBalance float64 `json:"current_balance"`
	NewBalance     float64 `json:"new_balance"`
}

type LoanRepository struct {
	dbo *db.DBO
}

func NewLoanRepository(dbo *db.DBO) *LoanRepository {
	return &LoanRepository{dbo: dbo}
}

// Check runs the 3 business checks and returns a preview or an explicit error.
func (r *LoanRepository) Check(userId int, barcode string) (*LoanPreview, error) {
	// Fetch the copy and its item
	var copyId int
	var itemId int
	var deposit float64
	var title string
	var isBorrowed bool
	err := r.dbo.QueryRow(`
		SELECT c.id, c.item_id, c.is_borrowed, i.deposit, i.title
		FROM copies c
		JOIN items i ON i.id = c.item_id
		WHERE c.barcode = $1 LIMIT 1`, barcode).
		Scan(&copyId, &itemId, &isBorrowed, &deposit, &title)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, i18n.NewError("copies.barcode_not_found")
	}
	if err != nil {
		return nil, fmt.Errorf("error while looking up the copy: %w", err)
	}
	if isBorrowed {
		return nil, i18n.NewError("copies.already_borrowed")
	}

	// 1.1 - Check that there is no overdue loan
	var overdueCount int
	err = r.dbo.QueryRow(`
		SELECT COUNT(*) FROM copies
		WHERE borrower_id = $1 AND is_borrowed = TRUE AND loan_due_date < $2`,
		userId, time.Now()).Scan(&overdueCount)
	if err != nil {
		return nil, fmt.Errorf("error while checking overdue loans: %w", err)
	}
	if overdueCount > 0 {
		return nil, i18n.NewError("loans.overdue", overdueCount)
	}

	// 1.2 - Check that the user is not already borrowing a copy of the same item
	var existing int
	err = r.dbo.QueryRow(`
		SELECT COUNT(*) FROM copies
		WHERE borrower_id = $1 AND item_id = $2 AND is_borrowed = TRUE`,
		userId, itemId).Scan(&existing)
	if err != nil {
		return nil, fmt.Errorf("error while checking existing loans: %w", err)
	}
	if existing > 0 {
		return nil, i18n.NewError("loans.already_borrowing")
	}

	// 1.3 - Check that the deposit balance is sufficient
	var balance float64
	err = r.dbo.QueryRow(`
		SELECT deposit_balance FROM users WHERE id = $1
		UNION ALL
		SELECT deposit_balance FROM students WHERE id = $1
		UNION ALL
		SELECT deposit_balance FROM teachers WHERE id = $1
		LIMIT 1`, userId).Scan(&balance)
	if err != nil {
		return nil, i18n.NewError("users.not_found")
	}
	if balance < deposit {
		return nil, i18n.NewError("loans.insufficient_funds", balance, deposit)
	}

	return &LoanPreview{
		Title:          title,
		Barcode:        barcode,
		Deposit:        deposit,
		CurrentBalance: balance,
		NewBalance:     balance - deposit,
	}, nil
}

// LoanSummary is returned in the list of a user's active loans.
type LoanSummary struct {
	Id        int    `json:"id"`
	Barcode   string `json:"barcode"`
	Title     string `json:"title"`
	StartDate string `json:"start_date"`
	DueDate   string `json:"due_date"`
	Overdue   bool   `json:"overdue"`
}

// GetLoans returns a user's active loans with the item information.
func (r *LoanRepository) GetLoans(userId int) ([]*LoanSummary, error) {
	rows, err := r.dbo.QueryRows(`
		SELECT c.id, c.barcode, i.title,
		       to_char(c.loan_start_date, 'YYYY-MM-DD'),
		       to_char(c.loan_due_date,   'YYYY-MM-DD'),
		       c.loan_due_date < NOW()
		FROM copies c
		JOIN items i ON i.id = c.item_id
		WHERE c.borrower_id = $1 AND c.is_borrowed = TRUE
		ORDER BY c.loan_due_date ASC`, userId)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	loans := make([]*LoanSummary, 0)
	for rows.Next() {
		loan := &LoanSummary{}
		if err := rows.Scan(&loan.Id, &loan.Barcode, &loan.Title,
			&loan.StartDate, &loan.DueDate, &loan.Overdue); err != nil {
			return nil, fmt.Errorf("GetLoans scan: %w", err)
		}
		loans = append(loans, loan)
	}
	return loans, rows.Err()
}

// OverdueLoan is returned in the global list of overdue loans (librarians).
type OverdueLoan struct {
	Id          int    `json:"id"`
	Barcode     string `json:"barcode"`
	Title       string `json:"title"`
	DueDate     string `json:"due_date"`
	LastName    string `json:"last_name"`
	FirstName   string `json:"first_name"`
	PhoneNumber string `json:"phone_number"`
}

// GetOverdueLoans returns every overdue copy with the borrower's information.
func (r *LoanRepository) GetOverdueLoans() ([]*OverdueLoan, error) {
	rows, err := r.dbo.QueryRows(`
		SELECT c.id, c.barcode, i.title,
		       to_char(c.loan_due_date, 'YYYY-MM-DD'),
		       u.last_name, u.first_name, COALESCE(u.phone_number, '')
		FROM copies c
		JOIN items i ON i.id = c.item_id
		JOIN (
			SELECT id, last_name, first_name, phone_number FROM students
			UNION ALL
			SELECT id, last_name, first_name, phone_number FROM teachers
			UNION ALL
			SELECT id, last_name, first_name, phone_number FROM ONLY users
		) u ON u.id = c.borrower_id
		WHERE c.is_borrowed = TRUE AND c.loan_due_date < NOW()
		ORDER BY c.loan_due_date ASC`)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	loans := make([]*OverdueLoan, 0)
	for rows.Next() {
		loan := &OverdueLoan{}
		if err := rows.Scan(&loan.Id, &loan.Barcode, &loan.Title,
			&loan.DueDate, &loan.LastName, &loan.FirstName, &loan.PhoneNumber); err != nil {
			return nil, fmt.Errorf("GetOverdueLoans scan: %w", err)
		}
		loans = append(loans, loan)
	}
	return loans, rows.Err()
}

// Borrow records the loan once the user has confirmed it.
func (r *LoanRepository) Borrow(userId int, barcode string) error {
	// Check again before recording
	preview, err := r.Check(userId, barcode)
	if err != nil {
		return err
	}

	// Fetch the copy id and the loan period
	var copyId, periodDays int
	err = r.dbo.QueryRow(`SELECT id, loan_period_days FROM copies WHERE barcode = $1`, barcode).
		Scan(&copyId, &periodDays)
	if err != nil {
		return fmt.Errorf("copy not found: %w", err)
	}

	now := time.Now()
	due := now.AddDate(0, 0, periodDays)

	// Mark the copy as borrowed
	n, err := r.dbo.Exec(`
		UPDATE copies
		SET is_borrowed = TRUE,
		    loan_start_date = $1,
		    loan_due_date = $2,
		    borrower_id = $3
		WHERE id = $4`, now, due, userId, copyId)
	if err != nil || n == 0 {
		return fmt.Errorf("unable to record the loan")
	}

	// Deduct the deposit - tries the 3 tables (native PostgreSQL inheritance)
	for _, table := range []string{"users", "students", "teachers"} {
		q := fmt.Sprintf(`UPDATE %s SET deposit_balance = deposit_balance - $1 WHERE id = $2`, table)
		n, err = r.dbo.Exec(q, preview.Deposit, userId)
		if err == nil && n > 0 {
			return nil
		}
	}
	return fmt.Errorf("unable to deduct the deposit from the balance")
}

// ReturnCopy marks a copy as returned and refunds the deposit to the borrower.
func (r *LoanRepository) ReturnCopy(copyId int) error {
	return r.dbo.WithTx(func(tx *db.TxDBO) error {
		// Fetch the item deposit and the borrower id
		var borrowerId int
		var deposit float64
		err := tx.QueryRow(`
			SELECT c.borrower_id, i.deposit
			FROM copies c
			JOIN items i ON i.id = c.item_id
			WHERE c.id = $1 AND c.is_borrowed = TRUE`, copyId).
			Scan(&borrowerId, &deposit)
		if errors.Is(err, sql.ErrNoRows) {
			return i18n.NewError("loans.not_found")
		}
		if err != nil {
			return fmt.Errorf("return copy: %w", err)
		}

		// Mark the copy as returned
		if _, err = tx.Exec(`
			UPDATE copies
			SET is_borrowed = FALSE,
			    borrower_id = NULL,
			    loan_start_date = NULL,
			    loan_due_date = NULL
			WHERE id = $1`, copyId); err != nil {
			return fmt.Errorf("return copy (update): %w", err)
		}

		// Refund the deposit in every inherited table
		for _, table := range []string{"users", "students", "teachers"} {
			q := fmt.Sprintf(`UPDATE %s SET deposit_balance = deposit_balance + $1 WHERE id = $2`, table)
			n, execErr := tx.Exec(q, deposit, borrowerId)
			if execErr == nil && n > 0 {
				return nil
			}
		}
		return fmt.Errorf("unable to refund the deposit")
	})
}
