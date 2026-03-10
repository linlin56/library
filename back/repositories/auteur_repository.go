package repositories

import (
	"database/sql"
	"errors"
	"fmt"

	"bibliotheque/db"
	"bibliotheque/models"
)

type AuteurRepository struct {
	dbo *db.DBO
}

func NewAuteurRepository(dbo *db.DBO) *AuteurRepository {
	return &AuteurRepository{dbo: dbo}
}

func (r *AuteurRepository) FindAll() ([]*models.Auteur, error) {
	rows, err := r.dbo.QueryRows(`SELECT id, nom, prenom FROM auteurs ORDER BY nom`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var auteurs []*models.Auteur
	for rows.Next() {
		a := &models.Auteur{}
		if err := rows.Scan(&a.Id, &a.Nom, &a.Prenom); err != nil {
			return nil, fmt.Errorf("FindAll scan: %w", err)
		}
		auteurs = append(auteurs, a)
	}
	return auteurs, rows.Err()
}

func (r *AuteurRepository) FindByID(id int) (*models.Auteur, error) {
	a := &models.Auteur{}
	err := r.dbo.QueryRow(`SELECT id, nom, prenom FROM auteurs WHERE id = $1`, id).
		Scan(&a.Id, &a.Nom, &a.Prenom)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("auteur %d introuvable", id)
	}
	if err != nil {
		return nil, fmt.Errorf("FindByID: %w", err)
	}
	return a, nil
}

func (r *AuteurRepository) Create(a *models.Auteur) (int, error) {
	var newID int
	err := r.dbo.ExecReturning(
		`INSERT INTO auteurs (nom, prenom) VALUES ($1, $2) RETURNING id`,
		a.Nom, a.Prenom,
	).Scan(&newID)
	if err != nil {
		return 0, fmt.Errorf("Create auteur: %w", err)
	}
	return newID, nil
}

func (r *AuteurRepository) Update(a *models.Auteur) error {
	n, err := r.dbo.Exec(
		`UPDATE auteurs SET nom = $1, prenom = $2 WHERE id = $3`,
		a.Nom, a.Prenom, a.Id,
	)
	if err != nil {
		return fmt.Errorf("Update auteur: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("auteur %d introuvable", a.Id)
	}
	return nil
}

func (r *AuteurRepository) Delete(id int) error {
	n, err := r.dbo.Exec(`DELETE FROM auteurs WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("Delete auteur: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("auteur %d introuvable", id)
	}
	return nil
}