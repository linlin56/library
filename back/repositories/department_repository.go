package repositories

import (
	"database/sql"
	"errors"
	"fmt"

	"library/db"
	"library/models"
)

type DepartmentRepository struct {
	dbo *db.DBO
}

func NewDepartmentRepository(dbo *db.DBO) *DepartmentRepository {
	return &DepartmentRepository{dbo: dbo}
}

func (r *DepartmentRepository) FindAll() ([]*models.Department, error) {
	rows, err := r.dbo.QueryRows(`SELECT id, name FROM departments ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var departments []*models.Department
	for rows.Next() {
		d := &models.Department{}
		if err := rows.Scan(&d.Id, &d.Name); err != nil {
			return nil, fmt.Errorf("FindAll department scan: %w", err)
		}
		departments = append(departments, d)
	}
	return departments, rows.Err()
}

func (r *DepartmentRepository) FindByID(id int) (*models.Department, error) {
	d := &models.Department{}
	err := r.dbo.QueryRow(`SELECT id, name FROM departments WHERE id = $1`, id).
		Scan(&d.Id, &d.Name)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("department %d not found", id)
	}
	if err != nil {
		return nil, fmt.Errorf("FindByID department: %w", err)
	}
	return d, nil
}

func (r *DepartmentRepository) Create(d *models.Department) (int, error) {
	var newID int
	err := r.dbo.ExecReturning(
		`INSERT INTO departments (name) VALUES ($1) RETURNING id`,
		d.Name,
	).Scan(&newID)
	if err != nil {
		return 0, fmt.Errorf("Create department: %w", err)
	}
	return newID, nil
}

func (r *DepartmentRepository) Update(d *models.Department) error {
	n, err := r.dbo.Exec(`UPDATE departments SET name=$1 WHERE id=$2`, d.Name, d.Id)
	if err != nil {
		return fmt.Errorf("Update department: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("department %d not found", d.Id)
	}
	return nil
}

func (r *DepartmentRepository) Delete(id int) error {
	n, err := r.dbo.Exec(`DELETE FROM departments WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("Delete department: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("department %d not found", id)
	}
	return nil
}
