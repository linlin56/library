package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"library/db"
	"library/i18n"
	"library/models"
)

type CopyRepository struct {
	dbo *db.DBO
}

func NewCopyRepository(dbo *db.DBO) *CopyRepository {
	return &CopyRepository{dbo: dbo}
}

func (r *CopyRepository) FindAll() ([]*models.Copy, error) {
	rows, err := r.dbo.QueryRows(`
		SELECT id, is_borrowed, loan_start_date, loan_due_date,
		       loan_period_days, barcode, borrower_id
		FROM copies ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var copies []*models.Copy
	for rows.Next() {
		c, err := scanCopy(rows)
		if err != nil {
			return nil, fmt.Errorf("FindAll copy scan: %w", err)
		}
		copies = append(copies, c)
	}
	return copies, rows.Err()
}

func (r *CopyRepository) FindByID(id int) (*models.Copy, error) {
	row := r.dbo.QueryRow(`
		SELECT id, is_borrowed, loan_start_date, loan_due_date,
		       loan_period_days, barcode, borrower_id
		FROM copies WHERE id = $1`, id)

	c := &models.Copy{}
	var borrowerID sql.NullInt64
	err := row.Scan(&c.Id, &c.IsBorrowed, &c.LoanStartDate, &c.LoanDueDate,
		&c.LoanPeriodDays, &c.Barcode, &borrowerID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("copy %d not found", id)
	}
	if err != nil {
		return nil, fmt.Errorf("FindByID copy: %w", err)
	}
	if borrowerID.Valid {
		c.BorrowerId = int(borrowerID.Int64)
	}
	return c, nil
}

func (r *CopyRepository) FindByBarcode(barcode string) (*models.Copy, error) {
	row := r.dbo.QueryRow(`
		SELECT id, is_borrowed, loan_start_date, loan_due_date,
		       loan_period_days, barcode, borrower_id
		FROM copies WHERE barcode = $1`, barcode)

	c := &models.Copy{}
	var borrowerID sql.NullInt64
	err := row.Scan(&c.Id, &c.IsBorrowed, &c.LoanStartDate, &c.LoanDueDate,
		&c.LoanPeriodDays, &c.Barcode, &borrowerID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("barcode '%s' not found", barcode)
	}
	if err != nil {
		return nil, fmt.Errorf("FindByBarcode: %w", err)
	}
	if borrowerID.Valid {
		c.BorrowerId = int(borrowerID.Int64)
	}
	return c, nil
}

func (r *CopyRepository) FindAvailable() ([]*models.Copy, error) {
	rows, err := r.dbo.QueryRows(`
		SELECT id, is_borrowed, loan_start_date, loan_due_date,
		       loan_period_days, barcode, borrower_id
		FROM copies WHERE is_borrowed = false ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var copies []*models.Copy
	for rows.Next() {
		c, err := scanCopy(rows)
		if err != nil {
			return nil, fmt.Errorf("FindAvailable scan: %w", err)
		}
		copies = append(copies, c)
	}
	return copies, rows.Err()
}

func (r *CopyRepository) FindAllByItemId(itemId int) ([]*models.Copy, error) {
	rows, err := r.dbo.QueryRows(`
		SELECT id, is_borrowed, barcode, loan_due_date
		FROM copies
		WHERE item_id = $1
		ORDER BY barcode`, itemId)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var result []*models.Copy
	for rows.Next() {
		c := &models.Copy{}
		var dueDate sql.NullTime
		if err := rows.Scan(&c.Id, &c.IsBorrowed, &c.Barcode, &dueDate); err != nil {
			return nil, fmt.Errorf("FindAllByItemId scan: %w", err)
		}
		if dueDate.Valid {
			c.LoanDueDate = dueDate.Time
		}
		result = append(result, c)
	}
	return result, rows.Err()
}

func (r *CopyRepository) CreateForItem(itemId int, barcode string, periodDays int) (int, error) {
	// Check that the barcode is unique (across all items)
	var exists bool
	if err := r.dbo.QueryRow(`SELECT EXISTS(SELECT 1 FROM copies WHERE barcode = $1)`, barcode).Scan(&exists); err != nil {
		return 0, fmt.Errorf("CreateForItem (check): %w", err)
	}
	if exists {
		return 0, i18n.NewError("copies.barcode_exists")
	}

	var newID int
	err := r.dbo.ExecReturning(
		`INSERT INTO copies (item_id, barcode, loan_period_days, is_borrowed)
		 VALUES ($1, $2, $3, false) RETURNING id`,
		itemId, barcode, periodDays,
	).Scan(&newID)
	if err != nil {
		return 0, fmt.Errorf("CreateForItem: %w", err)
	}
	return newID, nil
}

func (r *CopyRepository) FindAvailableByItemId(itemId int) ([]*models.Copy, error) {
	rows, err := r.dbo.QueryRows(`
		SELECT id, is_borrowed, barcode, loan_period_days
		FROM copies
		WHERE item_id = $1 AND is_borrowed = false
		ORDER BY barcode`, itemId)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var result []*models.Copy
	for rows.Next() {
		c := &models.Copy{}
		if err := rows.Scan(&c.Id, &c.IsBorrowed, &c.Barcode, &c.LoanPeriodDays); err != nil {
			return nil, fmt.Errorf("FindAvailableByItemId scan: %w", err)
		}
		result = append(result, c)
	}
	return result, rows.Err()
}

func (r *CopyRepository) Create(c *models.Copy) (int, error) {
	var newID int
	err := r.dbo.ExecReturning(
		`INSERT INTO copies (is_borrowed, loan_start_date, loan_due_date, loan_period_days, barcode, borrower_id)
		 VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		c.IsBorrowed, c.LoanStartDate, c.LoanDueDate,
		c.LoanPeriodDays, c.Barcode, nullableInt(c.BorrowerId),
	).Scan(&newID)
	if err != nil {
		return 0, fmt.Errorf("Create copy: %w", err)
	}
	return newID, nil
}

func (r *CopyRepository) Update(c *models.Copy) error {
	n, err := r.dbo.Exec(`
		UPDATE copies SET is_borrowed=$1, loan_start_date=$2, loan_due_date=$3,
		loan_period_days=$4, barcode=$5, borrower_id=$6 WHERE id=$7`,
		c.IsBorrowed, c.LoanStartDate, c.LoanDueDate,
		c.LoanPeriodDays, c.Barcode, nullableInt(c.BorrowerId), c.Id,
	)
	if err != nil {
		return fmt.Errorf("Update copy: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("copy %d not found", c.Id)
	}
	return nil
}

// Borrow marks a copy as borrowed by a user
func (r *CopyRepository) Borrow(copyID int, borrowerID int, periodDays int) error {
	start := time.Now()
	due := start.AddDate(0, 0, periodDays)
	n, err := r.dbo.Exec(`
		UPDATE copies SET is_borrowed=true, loan_start_date=$1, loan_due_date=$2,
		borrower_id=$3, loan_period_days=$4 WHERE id=$5 AND is_borrowed=false`,
		start, due, borrowerID, periodDays, copyID,
	)
	if err != nil {
		return fmt.Errorf("Borrow: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("copy %d not found or already borrowed", copyID)
	}
	return nil
}

// Return marks a copy as returned
func (r *CopyRepository) Return(copyID int) error {
	n, err := r.dbo.Exec(`
		UPDATE copies SET is_borrowed=false, borrower_id=NULL,
		loan_start_date=NULL, loan_due_date=NULL WHERE id=$1`,
		copyID,
	)
	if err != nil {
		return fmt.Errorf("Return: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("copy %d not found", copyID)
	}
	return nil
}

func (r *CopyRepository) Delete(id int) error {
	n, err := r.dbo.Exec(`DELETE FROM copies WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("Delete copy: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("copy %d not found", id)
	}
	return nil
}

// --- helpers ---

func scanCopy(rows interface {
	Scan(...any) error
}) (*models.Copy, error) {
	c := &models.Copy{}
	var borrowerID sql.NullInt64
	err := rows.Scan(&c.Id, &c.IsBorrowed, &c.LoanStartDate, &c.LoanDueDate,
		&c.LoanPeriodDays, &c.Barcode, &borrowerID)
	if borrowerID.Valid {
		c.BorrowerId = int(borrowerID.Int64)
	}
	return c, err
}

// nullableInt converts 0 to NULL for optional foreign keys
func nullableInt(id int) sql.NullInt64 {
	if id == 0 {
		return sql.NullInt64{Valid: false}
	}
	return sql.NullInt64{Int64: int64(id), Valid: true}
}
