package repositories

import (
	"database/sql"
	"errors"
	"fmt"

	"bibliotheque/db"
	"bibliotheque/models"
)

type OuvrageRepository struct {
	dbo *db.DBO
}

func NewOuvrageRepository(dbo *db.DBO) *OuvrageRepository {
	return &OuvrageRepository{dbo: dbo}
}

func (r *OuvrageRepository) FindAll() ([]*models.Ouvrage, error) {
	rows, err := r.dbo.QueryRows(`SELECT id, caution, titre, exemplaires FROM ouvrages ORDER BY titre`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ouvrages []*models.Ouvrage
	for rows.Next() {
		o := &models.Ouvrage{}
		if err := rows.Scan(&o.Id, &o.Caution, &o.Titre, &o.Exemplaires); err != nil {
			return nil, fmt.Errorf("FindAll ouvrage scan: %w", err)
		}
		ouvrages = append(ouvrages, o)
	}
	return ouvrages, rows.Err()
}

func (r *OuvrageRepository) FindByID(id int) (*models.Ouvrage, error) {
	o := &models.Ouvrage{}
	err := r.dbo.QueryRow(`SELECT id, caution, titre, exemplaires FROM ouvrages WHERE id = $1`, id).
		Scan(&o.Id, &o.Caution, &o.Titre, &o.Exemplaires)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("ouvrage %d introuvable", id)
	}
	if err != nil {
		return nil, fmt.Errorf("FindByID ouvrage: %w", err)
	}
	return o, nil
}

func (r *OuvrageRepository) FindByTitre(titre string) ([]*models.Ouvrage, error) {
	rows, err := r.dbo.QueryRows(
		`SELECT id, caution, titre, exemplaires FROM ouvrages WHERE titre ILIKE $1 ORDER BY titre`,
		"%"+titre+"%",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ouvrages []*models.Ouvrage
	for rows.Next() {
		o := &models.Ouvrage{}
		if err := rows.Scan(&o.Id, &o.Caution, &o.Titre, &o.Exemplaires); err != nil {
			return nil, fmt.Errorf("FindByTitre scan: %w", err)
		}
		ouvrages = append(ouvrages, o)
	}
	return ouvrages, rows.Err()
}

func (r *OuvrageRepository) Create(o *models.Ouvrage) (int, error) {
	var newID int
	err := r.dbo.ExecReturning(
		`INSERT INTO ouvrages (caution, titre, exemplaires) VALUES ($1, $2, $3) RETURNING id`,
		o.Caution, o.Titre, o.Exemplaires,
	).Scan(&newID)
	if err != nil {
		return 0, fmt.Errorf("Create ouvrage: %w", err)
	}
	return newID, nil
}

func (r *OuvrageRepository) Update(o *models.Ouvrage) error {
	n, err := r.dbo.Exec(
		`UPDATE ouvrages SET caution=$1, titre=$2, exemplaires=$3 WHERE id=$4`,
		o.Caution, o.Titre, o.Exemplaires, o.Id,
	)
	if err != nil {
		return fmt.Errorf("Update ouvrage: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("ouvrage %d introuvable", o.Id)
	}
	return nil
}

func (r *OuvrageRepository) Delete(id int) error {
	n, err := r.dbo.Exec(`DELETE FROM ouvrages WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("Delete ouvrage: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("ouvrage %d introuvable", id)
	}
	return nil
}