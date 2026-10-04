package repositories

import (
	"database/sql"
	"errors"
	"fmt"

	"library/db"
	"library/models"
)

type StudentRepository struct {
	dbo *db.DBO
}

func NewStudentRepository(dbo *db.DBO) *StudentRepository {
	return &StudentRepository{dbo: dbo}
}

func (r *StudentRepository) FindAll() ([]*models.Student, error) {
	rows, err := r.dbo.QueryRows(`
		SELECT u.id, u.last_name, u.first_name, u.phone_number, u.deposit_balance,
		       u.login, u.password, u.birth_date, u.email, s.study_year
		FROM students s
		JOIN users u ON u.id = s.id
		ORDER BY u.last_name`)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var students []*models.Student
	for rows.Next() {
		s, err := scanStudent(rows)
		if err != nil {
			return nil, fmt.Errorf("FindAll student scan: %w", err)
		}
		students = append(students, s)
	}
	return students, rows.Err()
}

func (r *StudentRepository) FindByID(id int) (*models.Student, error) {
	row := r.dbo.QueryRow(`
		SELECT u.id, u.last_name, u.first_name, u.phone_number, u.deposit_balance,
		       u.login, u.password, u.birth_date, u.email, s.study_year
		FROM students s
		JOIN users u ON u.id = s.id
		WHERE s.id = $1`, id)

	s := &models.Student{}
	err := row.Scan(&s.Id, &s.LastName, &s.FirstName, &s.PhoneNumber, &s.DepositBalance,
		&s.Login, &s.Password, &s.BirthDate, &s.Email, &s.StudyYear)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("student %d not found", id)
	}
	if err != nil {
		return nil, fmt.Errorf("FindByID student: %w", err)
	}
	return s, nil
}

func (r *StudentRepository) Create(s *models.Student) (int, error) {
	var newID int
	err := r.dbo.ExecReturning(`
		INSERT INTO students (last_name, first_name, phone_number, deposit_balance, login, password, birth_date, email, address_id, study_year)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`,
		s.LastName, s.FirstName, s.PhoneNumber, s.DepositBalance,
		s.Login, s.Password, s.BirthDate, s.Email, s.AddressId, s.StudyYear,
	).Scan(&newID)
	if err != nil {
		return 0, fmt.Errorf("Create student: %w", err)
	}
	return newID, nil
}

func (r *StudentRepository) Update(s *models.Student) error {
	return r.dbo.WithTx(func(tx *db.TxDBO) error {
		if _, err := tx.Exec(`
			UPDATE users SET last_name=$1, first_name=$2, phone_number=$3, deposit_balance=$4,
			login=$5, password=$6, birth_date=$7, email=$8 WHERE id=$9`,
			s.LastName, s.FirstName, s.PhoneNumber, s.DepositBalance,
			s.Login, s.Password, s.BirthDate, s.Email, s.Id,
		); err != nil {
			return err
		}
		_, err := tx.Exec(`UPDATE students SET study_year=$1 WHERE id=$2`, s.StudyYear, s.Id)
		return err
	})
}

func (r *StudentRepository) Delete(id int) error {
	return r.dbo.WithTx(func(tx *db.TxDBO) error {
		if _, err := tx.Exec(`DELETE FROM students WHERE id = $1`, id); err != nil {
			return err
		}
		_, err := tx.Exec(`DELETE FROM users WHERE id = $1`, id)
		return err
	})
}

func scanStudent(rows interface {
	Scan(...any) error
}) (*models.Student, error) {
	s := &models.Student{}
	err := rows.Scan(&s.Id, &s.LastName, &s.FirstName, &s.PhoneNumber, &s.DepositBalance,
		&s.Login, &s.Password, &s.BirthDate, &s.Email, &s.StudyYear)
	return s, err
}
