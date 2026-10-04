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

type UserRepository struct {
	dbo *db.DBO
}

func NewUserRepository(dbo *db.DBO) *UserRepository {
	return &UserRepository{dbo: dbo}
}

func (r *UserRepository) FindAll() ([]*models.User, error) {
	rows, err := r.dbo.QueryRows(`
		SELECT id, last_name, first_name, phone_number, deposit_balance, login, password, birth_date, email
		FROM users ORDER BY last_name`)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var users []*models.User
	for rows.Next() {
		u, err := scanUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

func (r *UserRepository) FindByID(id int) (*models.User, error) {
	row := r.dbo.QueryRow(`
		SELECT id, last_name, first_name, phone_number, deposit_balance, login, password, birth_date, email
		FROM users WHERE id = $1`, id)

	u := &models.User{}
	err := row.Scan(&u.Id, &u.LastName, &u.FirstName, &u.PhoneNumber,
		&u.DepositBalance, &u.Login, &u.Password, &u.BirthDate, &u.Email)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("user %d not found", id)
	}
	if err != nil {
		return nil, fmt.Errorf("FindByID user: %w", err)
	}
	return u, nil
}

func (r *UserRepository) FindByLogin(login string) (*models.User, error) {
	// UNION ALL over every inherited table.
	// "FROM ONLY users" excludes the children, which are added explicitly
	// so that the extra columns (study_year, department_id…)
	// do not make the Scan fail.
	row := r.dbo.QueryRow(`
		SELECT id, last_name, first_name, phone_number, deposit_balance, login, password, birth_date, email
		FROM ONLY users WHERE login = $1
		UNION ALL
		SELECT id, last_name, first_name, phone_number, deposit_balance, login, password, birth_date, email
		FROM students WHERE login = $1
		UNION ALL
		SELECT id, last_name, first_name, phone_number, deposit_balance, login, password, birth_date, email
		FROM teachers WHERE login = $1
		LIMIT 1`, login)

	u := &models.User{}
	err := row.Scan(&u.Id, &u.LastName, &u.FirstName, &u.PhoneNumber,
		&u.DepositBalance, &u.Login, &u.Password, &u.BirthDate, &u.Email)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("login '%s' not found", login)
	}
	if err != nil {
		return nil, fmt.Errorf("FindByLogin: %w", err)
	}
	return u, nil
}

func (r *UserRepository) Create(u *models.User) (int, error) {
	var newID int
	err := r.dbo.ExecReturning(`
		INSERT INTO users (last_name, first_name, phone_number, deposit_balance, login, password, birth_date, email)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`,
		u.LastName, u.FirstName, u.PhoneNumber, u.DepositBalance,
		u.Login, u.Password, u.BirthDate, u.Email,
	).Scan(&newID)
	if err != nil {
		return 0, fmt.Errorf("Create user: %w", err)
	}
	return newID, nil
}

func (r *UserRepository) Update(u *models.User) error {
	n, err := r.dbo.Exec(`
		UPDATE users SET last_name=$1, first_name=$2, phone_number=$3, deposit_balance=$4,
		login=$5, password=$6, birth_date=$7, email=$8 WHERE id=$9`,
		u.LastName, u.FirstName, u.PhoneNumber, u.DepositBalance,
		u.Login, u.Password, u.BirthDate, u.Email, u.Id,
	)
	if err != nil {
		return fmt.Errorf("Update user: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("user %d not found", u.Id)
	}
	return nil
}

func (r *UserRepository) Delete(id int) error {
	n, err := r.dbo.Exec(`DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("Delete user: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("user %d not found", id)
	}
	return nil
}

// FindActiveLoans returns the copies currently borrowed by a user
func (r *UserRepository) FindActiveLoans(userID int) ([]*models.Copy, error) {
	rows, err := r.dbo.QueryRows(`
		SELECT id, is_borrowed, loan_start_date, loan_due_date, loan_period_days, barcode, borrower_id
		FROM copies WHERE borrower_id = $1 AND is_borrowed = true`, userID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var copies []*models.Copy
	for rows.Next() {
		c := &models.Copy{}
		if err := rows.Scan(&c.Id, &c.IsBorrowed, &c.LoanStartDate, &c.LoanDueDate,
			&c.LoanPeriodDays, &c.Barcode, &c.BorrowerId); err != nil {
			return nil, fmt.Errorf("FindActiveLoans scan: %w", err)
		}
		copies = append(copies, c)
	}
	return copies, rows.Err()
}

// UserSummary is used when librarians search for users.
type UserSummary struct {
	Id          int    `json:"id"`
	LastName    string `json:"last_name"`
	FirstName   string `json:"first_name"`
	PhoneNumber string `json:"phone_number"`
	Role        string `json:"role"`
}

// SearchUsers looks for users by last name, first name, postal code and/or phone number.
// At least one non-empty criterion is expected (validated in the controller).
func (r *UserRepository) SearchUsers(lastName, firstName, postalCode, phoneNumber string) ([]*UserSummary, error) {
	const q = `
		SELECT u.id, u.last_name, u.first_name, COALESCE(u.phone_number, ''), 'student'
		FROM students u
		LEFT JOIN addresses a ON a.id = u.address_id
		WHERE ($1 = '' OR u.last_name ILIKE '%' || $1 || '%')
		  AND ($2 = '' OR u.first_name ILIKE '%' || $2 || '%')
		  AND ($3 = '' OR a.postal_code = $3)
		  AND ($4 = '' OR COALESCE(u.phone_number, '') ILIKE '%' || $4 || '%')
		UNION ALL
		SELECT u.id, u.last_name, u.first_name, COALESCE(u.phone_number, ''), 'teacher'
		FROM teachers u
		LEFT JOIN addresses a ON a.id = u.address_id
		WHERE ($1 = '' OR u.last_name ILIKE '%' || $1 || '%')
		  AND ($2 = '' OR u.first_name ILIKE '%' || $2 || '%')
		  AND ($3 = '' OR a.postal_code = $3)
		  AND ($4 = '' OR COALESCE(u.phone_number, '') ILIKE '%' || $4 || '%')
		UNION ALL
		SELECT u.id, u.last_name, u.first_name, COALESCE(u.phone_number, ''), 'user'
		FROM ONLY users u
		LEFT JOIN addresses a ON a.id = u.address_id
		WHERE ($1 = '' OR u.last_name ILIKE '%' || $1 || '%')
		  AND ($2 = '' OR u.first_name ILIKE '%' || $2 || '%')
		  AND ($3 = '' OR a.postal_code = $3)
		  AND ($4 = '' OR COALESCE(u.phone_number, '') ILIKE '%' || $4 || '%')
		ORDER BY last_name, first_name`

	rows, err := r.dbo.QueryRows(q, lastName, firstName, postalCode, phoneNumber)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	results := make([]*UserSummary, 0)
	for rows.Next() {
		u := &UserSummary{}
		if err := rows.Scan(&u.Id, &u.LastName, &u.FirstName, &u.PhoneNumber, &u.Role); err != nil {
			return nil, fmt.Errorf("SearchUsers scan: %w", err)
		}
		results = append(results, u)
	}
	return results, rows.Err()
}

// DepositInfo holds a user's current balance and total deposit.
type DepositInfo struct {
	DepositBalance float64 `json:"deposit_balance"`
	TotalDeposit   float64 `json:"total_deposit"`
}

// GetDeposit returns a user's deposit_balance and total_deposit (all inherited tables).
func (r *UserRepository) GetDeposit(id int) (*DepositInfo, error) {
	row := r.dbo.QueryRow(`
		SELECT deposit_balance, total_deposit FROM ONLY users WHERE id = $1
		UNION ALL
		SELECT deposit_balance, total_deposit FROM students WHERE id = $1
		UNION ALL
		SELECT deposit_balance, total_deposit FROM teachers WHERE id = $1
		LIMIT 1`, id)

	info := &DepositInfo{}
	err := row.Scan(&info.DepositBalance, &info.TotalDeposit)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, i18n.NewError("users.not_found")
	}
	if err != nil {
		return nil, fmt.Errorf("GetDeposit: %w", err)
	}
	return info, nil
}

// UpdateTotalDeposit updates total_deposit AND recomputes deposit_balance
// while keeping the amount currently borrowed (total_deposit - deposit_balance).
// Returns an error if the new value is lower than the borrowed amount.
func (r *UserRepository) UpdateTotalDeposit(id int, newTotalDeposit float64) error {
	current, err := r.GetDeposit(id)
	if err != nil {
		return err
	}

	borrowedAmount := current.TotalDeposit - current.DepositBalance
	if newTotalDeposit < borrowedAmount {
		return i18n.NewError("deposit.below_borrowed", borrowedAmount)
	}

	newBalance := newTotalDeposit - borrowedAmount

	tables := []string{"students", "teachers", "users"}
	for _, table := range tables {
		q := fmt.Sprintf(`UPDATE %s SET total_deposit = $1, deposit_balance = $2 WHERE id = $3`, table)
		n, err := r.dbo.Exec(q, newTotalDeposit, newBalance, id)
		if err != nil {
			return fmt.Errorf("UpdateTotalDeposit(%s): %w", table, err)
		}
		if n > 0 {
			return nil
		}
	}
	return i18n.NewError("users.not_found")
}

func (r *UserRepository) LoginExists(login string) (bool, error) {
	var count int
	err := r.dbo.QueryRow(`
		SELECT COUNT(*) FROM (
			SELECT id FROM ONLY users WHERE login = $1
			UNION ALL
			SELECT id FROM students WHERE login = $1
			UNION ALL
			SELECT id FROM teachers WHERE login = $1
			UNION ALL
			SELECT id FROM librarians WHERE login = $1
		) AS t`, login).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("LoginExists: %w", err)
	}
	return count > 0, nil
}

// --- internal helper ---

func scanUser(rows interface {
	Scan(...any) error
}) (*models.User, error) {
	u := &models.User{}
	var birthDate time.Time
	err := rows.Scan(&u.Id, &u.LastName, &u.FirstName, &u.PhoneNumber,
		&u.DepositBalance, &u.Login, &u.Password, &birthDate, &u.Email)
	u.BirthDate = birthDate
	return u, err
}
