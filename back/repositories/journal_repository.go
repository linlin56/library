package repositories

import (
	"database/sql"
	"errors"
	"fmt"

	"library/db"
	"library/models"
)

type JournalRepository struct {
	dbo *db.DBO
}

func NewJournalRepository(dbo *db.DBO) *JournalRepository {
	return &JournalRepository{dbo: dbo}
}

func (r *JournalRepository) FindAll() ([]*models.Journal, error) {
	rows, err := r.dbo.QueryRows(`
		SELECT j.id, i.deposit, i.title, i.copies, j.issue_number, j.publication_date
		FROM journals j
		JOIN items i ON i.id = j.id
		ORDER BY i.title, j.issue_number`)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var journals []*models.Journal
	for rows.Next() {
		j := &models.Journal{}
		if err := rows.Scan(&j.Id, &j.Deposit, &j.Title, &j.Copies, &j.IssueNumber, &j.PublicationDate); err != nil {
			return nil, fmt.Errorf("FindAll journal scan: %w", err)
		}
		journals = append(journals, j)
	}
	return journals, rows.Err()
}

func (r *JournalRepository) FindByID(id int) (*models.Journal, error) {
	j := &models.Journal{}
	err := r.dbo.QueryRow(`
		SELECT j.id, i.deposit, i.title, i.copies, j.issue_number, j.publication_date
		FROM journals j
		JOIN items i ON i.id = j.id
		WHERE j.id = $1`, id).
		Scan(&j.Id, &j.Deposit, &j.Title, &j.Copies, &j.IssueNumber, &j.PublicationDate)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("journal %d not found", id)
	}
	if err != nil {
		return nil, fmt.Errorf("FindByID journal: %w", err)
	}
	return j, nil
}

func (r *JournalRepository) FindByTitle(title string) ([]*models.Journal, error) {
	rows, err := r.dbo.QueryRows(`
		SELECT j.id, i.deposit, i.title, i.copies, j.issue_number, j.publication_date
		FROM journals j
		JOIN items i ON i.id = j.id
		WHERE i.title ILIKE $1
		ORDER BY i.title, j.issue_number`, "%"+title+"%")
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var journals []*models.Journal
	for rows.Next() {
		j := &models.Journal{}
		if err := rows.Scan(&j.Id, &j.Deposit, &j.Title, &j.Copies, &j.IssueNumber, &j.PublicationDate); err != nil {
			return nil, fmt.Errorf("FindByTitle journal scan: %w", err)
		}
		journals = append(journals, j)
	}
	return journals, rows.Err()
}

func (r *JournalRepository) Create(j *models.Journal) (int, error) {
	var newID int
	err := r.dbo.WithTx(func(tx *db.TxDBO) error {
		if err := tx.ExecReturning(
			`INSERT INTO items (deposit, title, copies) VALUES ($1, $2, $3) RETURNING id`,
			j.Deposit, j.Title, j.Copies,
		).Scan(&newID); err != nil {
			return fmt.Errorf("insert item: %w", err)
		}
		_, err := tx.Exec(
			`INSERT INTO journals (id, issue_number) VALUES ($1, $2)`,
			newID, j.IssueNumber,
		)
		return err
	})
	if err != nil {
		return 0, fmt.Errorf("Create journal: %w", err)
	}
	return newID, nil
}

func (r *JournalRepository) Update(j *models.Journal) error {
	return r.dbo.WithTx(func(tx *db.TxDBO) error {
		if _, err := tx.Exec(
			`UPDATE items SET deposit=$1, title=$2, copies=$3 WHERE id=$4`,
			j.Deposit, j.Title, j.Copies, j.Id,
		); err != nil {
			return err
		}
		_, err := tx.Exec(`UPDATE journals SET issue_number=$1 WHERE id=$2`, j.IssueNumber, j.Id)
		return err
	})
}

func (r *JournalRepository) Delete(id int) error {
	return r.dbo.WithTx(func(tx *db.TxDBO) error {
		if _, err := tx.Exec(`DELETE FROM journals WHERE id = $1`, id); err != nil {
			return err
		}
		_, err := tx.Exec(`DELETE FROM items WHERE id = $1`, id)
		return err
	})
}
