package repositories

import (
	"database/sql"
	"errors"
	"fmt"

	"library/db"
	"library/models"
)

// SearchFilters holds the advanced search criteria.
type SearchFilters struct {
	Title       string
	Author      string
	Isbn        string
	Barcode     string
	IssueNumber string
	Available   bool
}

// ItemLocation describes where an item is shelved.
type ItemLocation struct {
	BayNumber   int `json:"bay_number"`
	ShelfNumber int `json:"shelf_number"`
	Level       int `json:"level"`
}

// ItemResult is the enriched view returned by Search.
type ItemResult struct {
	Id              int           `json:"id"`
	Title           string        `json:"title"`
	Deposit         float64       `json:"deposit"`
	Type            string        `json:"type"`
	Isbn            string        `json:"isbn,omitempty"`
	Author          string        `json:"author,omitempty"`
	IssueNumber     *int          `json:"issue_number,omitempty"`
	PublicationDate string        `json:"publication_date,omitempty"`
	AvailableCopies int           `json:"available_copies"`
	Location        *ItemLocation `json:"location,omitempty"`
}

const searchSQL = `
SELECT
    b.id, b.deposit, b.title,
    b.isbn,
    COALESCE(a.first_name || ' ' || a.last_name, '') AS author,
    NULL::int AS issue_number,
    NULL::text AS publication_date,
    'book' AS type,
    (SELECT COUNT(*) FROM copies c WHERE c.item_id = b.id AND NOT c.is_borrowed) AS available_copies,
    l.bay_number, l.shelf_number, l.level
FROM books b
LEFT JOIN authors a ON a.id = b.author_id
LEFT JOIN locations l ON l.id = b.location_id
WHERE ($1 = '' OR b.title ILIKE '%' || $1 || '%')
  AND ($2 = '' OR a.last_name ILIKE '%' || $2 || '%' OR a.first_name ILIKE '%' || $2 || '%')
  AND ($3 = '' OR b.isbn ILIKE '%' || $3 || '%')
  AND ($4 = '' OR EXISTS(SELECT 1 FROM copies c WHERE c.item_id = b.id AND c.barcode = $4))
  AND (NOT $5 OR EXISTS(SELECT 1 FROM copies c WHERE c.item_id = b.id AND NOT c.is_borrowed))
  AND $6 = ''
UNION ALL
SELECT
    j.id, j.deposit, j.title,
    NULL::varchar AS isbn,
    NULL::varchar AS author,
    j.issue_number,
    TO_CHAR(j.publication_date, 'YYYY-MM-DD') AS publication_date,
    'journal' AS type,
    (SELECT COUNT(*) FROM copies c WHERE c.item_id = j.id AND NOT c.is_borrowed) AS available_copies,
    l.bay_number, l.shelf_number, l.level
FROM journals j
LEFT JOIN locations l ON l.id = j.location_id
WHERE ($1 = '' OR j.title ILIKE '%' || $1 || '%')
  AND ($4 = '' OR EXISTS(SELECT 1 FROM copies c WHERE c.item_id = j.id AND c.barcode = $4))
  AND ($6 = '' OR j.issue_number::text = $6)
  AND (NOT $5 OR EXISTS(SELECT 1 FROM copies c WHERE c.item_id = j.id AND NOT c.is_borrowed))
  AND $2 = ''
  AND $3 = ''
ORDER BY title
`

type ItemRepository struct {
	dbo *db.DBO
}

func NewItemRepository(dbo *db.DBO) *ItemRepository {
	return &ItemRepository{dbo: dbo}
}

// Search runs an advanced search with the given criteria.
func (r *ItemRepository) Search(f SearchFilters) ([]ItemResult, error) {
	rows, err := r.dbo.QueryRows(searchSQL,
		f.Title, f.Author, f.Isbn, f.Barcode, f.Available, f.IssueNumber,
	)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	results := make([]ItemResult, 0)
	for rows.Next() {
		var res ItemResult
		var isbn, author, publicationDate sql.NullString
		var issueNumber, bay, shelf, level sql.NullInt32
		if err := rows.Scan(&res.Id, &res.Deposit, &res.Title, &isbn, &author, &issueNumber,
			&publicationDate, &res.Type, &res.AvailableCopies, &bay, &shelf, &level); err != nil {
			return nil, fmt.Errorf("Search scan: %w", err)
		}
		if isbn.Valid {
			res.Isbn = isbn.String
		}
		if author.Valid {
			res.Author = author.String
		}
		if issueNumber.Valid {
			n := int(issueNumber.Int32)
			res.IssueNumber = &n
		}
		if publicationDate.Valid {
			res.PublicationDate = publicationDate.String
		}
		if bay.Valid {
			res.Location = &ItemLocation{
				BayNumber:   int(bay.Int32),
				ShelfNumber: int(shelf.Int32),
				Level:       int(level.Int32),
			}
		}
		results = append(results, res)
	}
	return results, rows.Err()
}

func (r *ItemRepository) FindAll() ([]*models.Item, error) {
	rows, err := r.dbo.QueryRows(`SELECT id, deposit, title FROM items ORDER BY title`)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	items := make([]*models.Item, 0)
	for rows.Next() {
		i := &models.Item{}
		if err := rows.Scan(&i.Id, &i.Deposit, &i.Title); err != nil {
			return nil, fmt.Errorf("FindAll item scan: %w", err)
		}
		items = append(items, i)
	}
	return items, rows.Err()
}

func (r *ItemRepository) FindByID(id int) (*models.Item, error) {
	i := &models.Item{}
	err := r.dbo.QueryRow(`SELECT id, deposit, title FROM items WHERE id = $1`, id).
		Scan(&i.Id, &i.Deposit, &i.Title)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("item %d not found", id)
	}
	if err != nil {
		return nil, fmt.Errorf("FindByID item: %w", err)
	}
	return i, nil
}

func (r *ItemRepository) Create(i *models.Item) (int, error) {
	var newID int
	err := r.dbo.ExecReturning(`INSERT INTO items (deposit, title) VALUES ($1, $2) RETURNING id`, i.Deposit, i.Title).Scan(&newID)
	if err != nil {
		return 0, fmt.Errorf("Create item: %w", err)
	}
	return newID, nil
}

func (r *ItemRepository) Update(i *models.Item) error {
	n, err := r.dbo.Exec(
		`UPDATE items SET deposit=$1, title=$2 WHERE id=$3`,
		i.Deposit, i.Title, i.Id,
	)
	if err != nil {
		return fmt.Errorf("Update item: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("item %d not found", i.Id)
	}
	return nil
}

func (r *ItemRepository) Delete(id int) error {
	n, err := r.dbo.Exec(`DELETE FROM items WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("Delete item: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("item %d not found", id)
	}
	return nil
}

func (r *ItemRepository) CreateBook(title string, deposit float64, isbn string, authorId int, locationId int) (int, error) {
	var newID int
	err := r.dbo.ExecReturning(
		`INSERT INTO books (title, deposit, isbn, author_id, location_id) VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		title, deposit, isbn, authorId, locationId,
	).Scan(&newID)
	if err != nil {
		return 0, fmt.Errorf("CreateBook: %w", err)
	}
	return newID, nil
}

func (r *ItemRepository) CreateJournal(title string, deposit float64, issueNumber int, publicationDate string, locationId int) (int, error) {
	var newID int
	err := r.dbo.ExecReturning(
		`INSERT INTO journals (title, deposit, issue_number, publication_date, location_id) VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		title, deposit, issueNumber, publicationDate, locationId,
	).Scan(&newID)
	if err != nil {
		return 0, fmt.Errorf("CreateJournal: %w", err)
	}
	return newID, nil
}
