package repositories

import (
	"database/sql"
	"errors"
	"fmt"

	"library/db"
	"library/models"
)

type AuthorRepository struct {
	dbo *db.DBO
}

func NewAuthorRepository(dbo *db.DBO) *AuthorRepository {
	return &AuthorRepository{dbo: dbo}
}

func (r *AuthorRepository) FindAll() ([]*models.Author, error) {
	rows, err := r.dbo.QueryRows(`SELECT id, last_name, first_name FROM authors ORDER BY last_name`)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var authors []*models.Author
	for rows.Next() {
		a := &models.Author{}
		if err := rows.Scan(&a.Id, &a.LastName, &a.FirstName); err != nil {
			return nil, fmt.Errorf("FindAll scan: %w", err)
		}
		authors = append(authors, a)
	}
	return authors, rows.Err()
}

func (r *AuthorRepository) FindByID(id int) (*models.Author, error) {
	a := &models.Author{}
	err := r.dbo.QueryRow(`SELECT id, last_name, first_name FROM authors WHERE id = $1`, id).
		Scan(&a.Id, &a.LastName, &a.FirstName)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("author %d not found", id)
	}
	if err != nil {
		return nil, fmt.Errorf("FindByID: %w", err)
	}
	return a, nil
}

func (r *AuthorRepository) Create(a *models.Author) (int, error) {
	var newID int
	err := r.dbo.ExecReturning(
		`INSERT INTO authors (last_name, first_name) VALUES ($1, $2) RETURNING id`,
		a.LastName, a.FirstName,
	).Scan(&newID)
	if err != nil {
		return 0, fmt.Errorf("Create author: %w", err)
	}
	return newID, nil
}

func (r *AuthorRepository) Update(a *models.Author) error {
	n, err := r.dbo.Exec(
		`UPDATE authors SET last_name = $1, first_name = $2 WHERE id = $3`,
		a.LastName, a.FirstName, a.Id,
	)
	if err != nil {
		return fmt.Errorf("Update author: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("author %d not found", a.Id)
	}
	return nil
}

func (r *AuthorRepository) Delete(id int) error {
	n, err := r.dbo.Exec(`DELETE FROM authors WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("Delete author: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("author %d not found", id)
	}
	return nil
}
