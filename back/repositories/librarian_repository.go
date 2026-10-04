package repositories

import (
	"database/sql"
	"errors"
	"fmt"

	"library/db"
	"library/models"
)

type LibrarianRepository struct {
	dbo *db.DBO
}

func NewLibrarianRepository(dbo *db.DBO) *LibrarianRepository {
	return &LibrarianRepository{dbo: dbo}
}

func (r *LibrarianRepository) FindAll() ([]*models.Librarian, error) {
	rows, err := r.dbo.QueryRows(`SELECT id, last_name, first_name, login, password FROM librarians ORDER BY last_name`)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var librarians []*models.Librarian
	for rows.Next() {
		l := &models.Librarian{}
		if err := rows.Scan(&l.Id, &l.LastName, &l.FirstName, &l.Login, &l.Password); err != nil {
			return nil, fmt.Errorf("FindAll librarian scan: %w", err)
		}
		librarians = append(librarians, l)
	}
	return librarians, rows.Err()
}

func (r *LibrarianRepository) FindByID(id int) (*models.Librarian, error) {
	l := &models.Librarian{}
	err := r.dbo.QueryRow(`SELECT id, last_name, first_name, login, password FROM librarians WHERE id = $1`, id).
		Scan(&l.Id, &l.LastName, &l.FirstName, &l.Login, &l.Password)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("librarian %d not found", id)
	}
	if err != nil {
		return nil, fmt.Errorf("FindByID librarian: %w", err)
	}
	return l, nil
}

func (r *LibrarianRepository) FindByLogin(login string) (*models.Librarian, error) {
	l := &models.Librarian{}
	err := r.dbo.QueryRow(`SELECT id, last_name, first_name, login, password FROM librarians WHERE login = $1`, login).
		Scan(&l.Id, &l.LastName, &l.FirstName, &l.Login, &l.Password)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("login '%s' not found", login)
	}
	if err != nil {
		return nil, fmt.Errorf("FindByLogin librarian: %w", err)
	}
	return l, nil
}

func (r *LibrarianRepository) Create(l *models.Librarian) (int, error) {
	var newID int
	err := r.dbo.ExecReturning(
		`INSERT INTO librarians (last_name, first_name, login, password) VALUES ($1, $2, $3, $4) RETURNING id`,
		l.LastName, l.FirstName, l.Login, l.Password,
	).Scan(&newID)
	if err != nil {
		return 0, fmt.Errorf("Create librarian: %w", err)
	}
	return newID, nil
}

func (r *LibrarianRepository) Update(l *models.Librarian) error {
	n, err := r.dbo.Exec(
		`UPDATE librarians SET last_name=$1, first_name=$2, login=$3, password=$4 WHERE id=$5`,
		l.LastName, l.FirstName, l.Login, l.Password, l.Id,
	)
	if err != nil {
		return fmt.Errorf("Update librarian: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("librarian %d not found", l.Id)
	}
	return nil
}

func (r *LibrarianRepository) Delete(id int) error {
	n, err := r.dbo.Exec(`DELETE FROM librarians WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("Delete librarian: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("librarian %d not found", id)
	}
	return nil
}
