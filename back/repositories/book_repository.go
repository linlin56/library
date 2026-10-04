package repositories

import (
	"database/sql"
	"errors"
	"fmt"

	"library/db"
	"library/models"
)

type BookRepository struct {
	dbo *db.DBO
}

func NewBookRepository(dbo *db.DBO) *BookRepository {
	return &BookRepository{dbo: dbo}
}

// FindAll returns every book with its author (JOIN)
func (r *BookRepository) FindAll() ([]*models.Book, error) {
	rows, err := r.dbo.QueryRows(`
		SELECT b.id, i.deposit, i.title, i.copies, b.isbn,
		       a.id, a.last_name, a.first_name
		FROM books b
		JOIN items i ON i.id = b.id
		JOIN authors a ON a.id = b.author_id
		ORDER BY i.title`)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var books []*models.Book
	for rows.Next() {
		b, err := scanBook(rows)
		if err != nil {
			return nil, fmt.Errorf("FindAll book scan: %w", err)
		}
		books = append(books, b)
	}
	return books, rows.Err()
}

func (r *BookRepository) FindByID(id int) (*models.Book, error) {
	row := r.dbo.QueryRow(`
		SELECT b.id, i.deposit, i.title, i.copies, b.isbn,
		       a.id, a.last_name, a.first_name
		FROM books b
		JOIN items i ON i.id = b.id
		JOIN authors a ON a.id = b.author_id
		WHERE b.id = $1`, id)

	b := &models.Book{}
	a := &models.Author{}
	err := row.Scan(&b.Id, &b.Deposit, &b.Title, &b.Copies, &b.Isbn,
		&a.Id, &a.LastName, &a.FirstName)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("book %d not found", id)
	}
	if err != nil {
		return nil, fmt.Errorf("FindByID book: %w", err)
	}
	b.Author = a
	b.AuthorId = a.Id
	return b, nil
}

func (r *BookRepository) FindByIsbn(isbn string) (*models.Book, error) {
	row := r.dbo.QueryRow(`
		SELECT b.id, i.deposit, i.title, i.copies, b.isbn,
		       a.id, a.last_name, a.first_name
		FROM books b
		JOIN items i ON i.id = b.id
		JOIN authors a ON a.id = b.author_id
		WHERE b.isbn = $1`, isbn)

	b := &models.Book{}
	a := &models.Author{}
	err := row.Scan(&b.Id, &b.Deposit, &b.Title, &b.Copies, &b.Isbn,
		&a.Id, &a.LastName, &a.FirstName)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("ISBN '%s' not found", isbn)
	}
	if err != nil {
		return nil, fmt.Errorf("FindByIsbn: %w", err)
	}
	b.Author = a
	b.AuthorId = a.Id
	return b, nil
}

func (r *BookRepository) FindByAuthor(authorID int) ([]*models.Book, error) {
	rows, err := r.dbo.QueryRows(`
		SELECT b.id, i.deposit, i.title, i.copies, b.isbn,
		       a.id, a.last_name, a.first_name
		FROM books b
		JOIN items i ON i.id = b.id
		JOIN authors a ON a.id = b.author_id
		WHERE b.author_id = $1
		ORDER BY i.title`, authorID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var books []*models.Book
	for rows.Next() {
		b, err := scanBook(rows)
		if err != nil {
			return nil, fmt.Errorf("FindByAuthor scan: %w", err)
		}
		books = append(books, b)
	}
	return books, rows.Err()
}

// Create inserts into items AND books inside a transaction
func (r *BookRepository) Create(b *models.Book) (int, error) {
	var newID int
	err := r.dbo.WithTx(func(tx *db.TxDBO) error {
		// 1. Insert into the parent table items
		if err := tx.ExecReturning(
			`INSERT INTO items (deposit, title, copies) VALUES ($1, $2, $3) RETURNING id`,
			b.Deposit, b.Title, b.Copies,
		).Scan(&newID); err != nil {
			return fmt.Errorf("insert item: %w", err)
		}

		// 2. Insert into the books table with the same id
		_, err := tx.Exec(
			`INSERT INTO books (id, author_id, isbn) VALUES ($1, $2, $3)`,
			newID, b.AuthorId, b.Isbn,
		)
		return err
	})
	if err != nil {
		return 0, fmt.Errorf("Create book: %w", err)
	}
	return newID, nil
}

func (r *BookRepository) Update(b *models.Book) error {
	return r.dbo.WithTx(func(tx *db.TxDBO) error {
		if _, err := tx.Exec(
			`UPDATE items SET deposit=$1, title=$2, copies=$3 WHERE id=$4`,
			b.Deposit, b.Title, b.Copies, b.Id,
		); err != nil {
			return fmt.Errorf("update item: %w", err)
		}
		_, err := tx.Exec(
			`UPDATE books SET author_id=$1, isbn=$2 WHERE id=$3`,
			b.AuthorId, b.Isbn, b.Id,
		)
		return err
	})
}

func (r *BookRepository) Delete(id int) error {
	return r.dbo.WithTx(func(tx *db.TxDBO) error {
		if _, err := tx.Exec(`DELETE FROM books WHERE id = $1`, id); err != nil {
			return err
		}
		_, err := tx.Exec(`DELETE FROM items WHERE id = $1`, id)
		return err
	})
}

// --- helper ---

func scanBook(rows interface {
	Scan(...any) error
}) (*models.Book, error) {
	b := &models.Book{}
	a := &models.Author{}
	err := rows.Scan(&b.Id, &b.Deposit, &b.Title, &b.Copies, &b.Isbn,
		&a.Id, &a.LastName, &a.FirstName)
	b.Author = a
	b.AuthorId = a.Id
	return b, err
}
