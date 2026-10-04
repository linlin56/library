package repositories

import (
	"fmt"

	"library/db"
)

type LocationSummary struct {
	Id           int    `json:"id"`
	BayNumber    int    `json:"bay_number"`
	ShelfNumber  int    `json:"shelf_number"`
	Level        int    `json:"level"`
	CategoryName string `json:"category_name"`
}

type LocationRepository struct {
	dbo *db.DBO
}

func NewLocationRepository(dbo *db.DBO) *LocationRepository {
	return &LocationRepository{dbo: dbo}
}

func (r *LocationRepository) FindAll() ([]*LocationSummary, error) {
	rows, err := r.dbo.QueryRows(`
		SELECT l.id, l.bay_number, l.shelf_number, l.level, COALESCE(c.name, '')
		FROM locations l
		LEFT JOIN categories c ON c.id = l.category_id
		ORDER BY l.bay_number, l.shelf_number, l.level`)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var locations []*LocationSummary
	for rows.Next() {
		l := &LocationSummary{}
		if err := rows.Scan(&l.Id, &l.BayNumber, &l.ShelfNumber, &l.Level, &l.CategoryName); err != nil {
			return nil, fmt.Errorf("FindAll location scan: %w", err)
		}
		locations = append(locations, l)
	}
	return locations, rows.Err()
}
