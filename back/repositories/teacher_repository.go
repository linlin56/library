package repositories

import (
	"database/sql"
	"errors"
	"fmt"

	"library/db"
	"library/models"
)

type TeacherRepository struct {
	dbo *db.DBO
}

func NewTeacherRepository(dbo *db.DBO) *TeacherRepository {
	return &TeacherRepository{dbo: dbo}
}

func (r *TeacherRepository) FindAll() ([]*models.Teacher, error) {
	rows, err := r.dbo.QueryRows(`
		SELECT u.id, u.last_name, u.first_name, u.phone_number, u.deposit_balance,
		       u.login, u.password, u.birth_date, u.email,
		       d.id, d.name
		FROM teachers t
		JOIN users u ON u.id = t.id
		JOIN departments d ON d.id = t.department_id
		ORDER BY u.last_name`)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var teachers []*models.Teacher
	for rows.Next() {
		t, err := scanTeacher(rows)
		if err != nil {
			return nil, fmt.Errorf("FindAll teacher scan: %w", err)
		}
		teachers = append(teachers, t)
	}
	return teachers, rows.Err()
}

func (r *TeacherRepository) FindByID(id int) (*models.Teacher, error) {
	row := r.dbo.QueryRow(`
		SELECT u.id, u.last_name, u.first_name, u.phone_number, u.deposit_balance,
		       u.login, u.password, u.birth_date, u.email,
		       d.id, d.name
		FROM teachers t
		JOIN users u ON u.id = t.id
		JOIN departments d ON d.id = t.department_id
		WHERE t.id = $1`, id)

	t := &models.Teacher{}
	d := &models.Department{}
	err := row.Scan(&t.Id, &t.LastName, &t.FirstName, &t.PhoneNumber, &t.DepositBalance,
		&t.Login, &t.Password, &t.BirthDate, &t.Email,
		&d.Id, &d.Name)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("teacher %d not found", id)
	}
	if err != nil {
		return nil, fmt.Errorf("FindByID teacher: %w", err)
	}
	t.Department = d
	t.DepartmentId = d.Id
	return t, nil
}

func (r *TeacherRepository) FindByDepartment(departmentID int) ([]*models.Teacher, error) {
	rows, err := r.dbo.QueryRows(`
		SELECT u.id, u.last_name, u.first_name, u.phone_number, u.deposit_balance,
		       u.login, u.password, u.birth_date, u.email,
		       d.id, d.name
		FROM teachers t
		JOIN users u ON u.id = t.id
		JOIN departments d ON d.id = t.department_id
		WHERE t.department_id = $1
		ORDER BY u.last_name`, departmentID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var teachers []*models.Teacher
	for rows.Next() {
		t, err := scanTeacher(rows)
		if err != nil {
			return nil, fmt.Errorf("FindByDepartment scan: %w", err)
		}
		teachers = append(teachers, t)
	}
	return teachers, rows.Err()
}

func (r *TeacherRepository) Create(t *models.Teacher) (int, error) {
	var newID int
	err := r.dbo.ExecReturning(`
		INSERT INTO teachers (last_name, first_name, phone_number, deposit_balance, login, password, birth_date, email, address_id, department_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`,
		t.LastName, t.FirstName, t.PhoneNumber, t.DepositBalance,
		t.Login, t.Password, t.BirthDate, t.Email, t.AddressId, t.DepartmentId,
	).Scan(&newID)
	if err != nil {
		return 0, fmt.Errorf("Create teacher: %w", err)
	}
	return newID, nil
}

func (r *TeacherRepository) Update(t *models.Teacher) error {
	return r.dbo.WithTx(func(tx *db.TxDBO) error {
		if _, err := tx.Exec(`
			UPDATE users SET last_name=$1, first_name=$2, phone_number=$3, deposit_balance=$4,
			login=$5, password=$6, birth_date=$7, email=$8 WHERE id=$9`,
			t.LastName, t.FirstName, t.PhoneNumber, t.DepositBalance,
			t.Login, t.Password, t.BirthDate, t.Email, t.Id,
		); err != nil {
			return err
		}
		_, err := tx.Exec(`UPDATE teachers SET department_id=$1 WHERE id=$2`, t.DepartmentId, t.Id)
		return err
	})
}

func (r *TeacherRepository) Delete(id int) error {
	return r.dbo.WithTx(func(tx *db.TxDBO) error {
		if _, err := tx.Exec(`DELETE FROM teachers WHERE id = $1`, id); err != nil {
			return err
		}
		_, err := tx.Exec(`DELETE FROM users WHERE id = $1`, id)
		return err
	})
}

func scanTeacher(rows interface {
	Scan(...any) error
}) (*models.Teacher, error) {
	t := &models.Teacher{}
	d := &models.Department{}
	err := rows.Scan(&t.Id, &t.LastName, &t.FirstName, &t.PhoneNumber, &t.DepositBalance,
		&t.Login, &t.Password, &t.BirthDate, &t.Email,
		&d.Id, &d.Name)
	t.Department = d
	t.DepartmentId = d.Id
	return t, err
}
