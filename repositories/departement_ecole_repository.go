package repositories

import (
	"database/sql"
	"errors"
	"fmt"

	"bibliotheque/db"
	"bibliotheque/models"
)

type DepartementEcoleRepository struct {
	dbo *db.DBO
}

func NewDepartementEcoleRepository(dbo *db.DBO) *DepartementEcoleRepository {
	return &DepartementEcoleRepository{dbo: dbo}
}

func (r *DepartementEcoleRepository) FindAll() ([]*models.DepartementEcole, error) {
	rows, err := r.dbo.QueryRows(`SELECT id, nom FROM departements_ecole ORDER BY nom`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var departements []*models.DepartementEcole
	for rows.Next() {
		d := &models.DepartementEcole{}
		if err := rows.Scan(&d.Id, &d.Nom); err != nil {
			return nil, fmt.Errorf("FindAll departement scan: %w", err)
		}
		departements = append(departements, d)
	}
	return departements, rows.Err()
}

func (r *DepartementEcoleRepository) FindByID(id int) (*models.DepartementEcole, error) {
	d := &models.DepartementEcole{}
	err := r.dbo.QueryRow(`SELECT id, nom FROM departements_ecole WHERE id = $1`, id).
		Scan(&d.Id, &d.Nom)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("département %d introuvable", id)
	}
	if err != nil {
		return nil, fmt.Errorf("FindByID departement: %w", err)
	}
	return d, nil
}

func (r *DepartementEcoleRepository) Create(d *models.DepartementEcole) (int, error) {
	var newID int
	err := r.dbo.ExecReturning(
		`INSERT INTO departements_ecole (nom) VALUES ($1) RETURNING id`,
		d.Nom,
	).Scan(&newID)
	if err != nil {
		return 0, fmt.Errorf("Create departement: %w", err)
	}
	return newID, nil
}

func (r *DepartementEcoleRepository) Update(d *models.DepartementEcole) error {
	n, err := r.dbo.Exec(`UPDATE departements_ecole SET nom=$1 WHERE id=$2`, d.Nom, d.Id)
	if err != nil {
		return fmt.Errorf("Update departement: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("département %d introuvable", d.Id)
	}
	return nil
}

func (r *DepartementEcoleRepository) Delete(id int) error {
	n, err := r.dbo.Exec(`DELETE FROM departements_ecole WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("Delete departement: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("département %d introuvable", id)
	}
	return nil
}