package repositories

import (
	"database/sql"
	"errors"
	"fmt"

	"bibliotheque/db"
	"bibliotheque/models"
)

type BibliothécaireRepository struct {
	dbo *db.DBO
}

func NewBibliothécaireRepository(dbo *db.DBO) *BibliothécaireRepository {
	return &BibliothécaireRepository{dbo: dbo}
}

func (r *BibliothécaireRepository) FindAll() ([]*models.Bibliothecaire, error) {
	rows, err := r.dbo.QueryRows(`SELECT id, nom, prenom, login, mot_de_passe FROM bibliothecaires ORDER BY nom`)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var bibliothecaires []*models.Bibliothecaire
	for rows.Next() {
		b := &models.Bibliothecaire{}
		if err := rows.Scan(&b.Id, &b.Nom, &b.Prenom, &b.Login, &b.MotDePasse); err != nil {
			return nil, fmt.Errorf("FindAll bibliothecaire scan: %w", err)
		}
		bibliothecaires = append(bibliothecaires, b)
	}
	return bibliothecaires, rows.Err()
}

func (r *BibliothécaireRepository) FindByID(id int) (*models.Bibliothecaire, error) {
	b := &models.Bibliothecaire{}
	err := r.dbo.QueryRow(`SELECT id, nom, prenom, login, mot_de_passe FROM bibliothecaires WHERE id = $1`, id).
		Scan(&b.Id, &b.Nom, &b.Prenom, &b.Login, &b.MotDePasse)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("bibliothécaire %d introuvable", id)
	}
	if err != nil {
		return nil, fmt.Errorf("FindByID bibliothecaire: %w", err)
	}
	return b, nil
}

func (r *BibliothécaireRepository) FindByLogin(login string) (*models.Bibliothecaire, error) {
	b := &models.Bibliothecaire{}
	err := r.dbo.QueryRow(`SELECT id, nom, prenom, login, mot_de_passe FROM bibliothecaires WHERE login = $1`, login).
		Scan(&b.Id, &b.Nom, &b.Prenom, &b.Login, &b.MotDePasse)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("login '%s' introuvable", login)
	}
	if err != nil {
		return nil, fmt.Errorf("FindByLogin bibliothecaire: %w", err)
	}
	return b, nil
}

func (r *BibliothécaireRepository) Create(b *models.Bibliothecaire) (int, error) {
	var newID int
	err := r.dbo.ExecReturning(
		`INSERT INTO bibliothecaires (nom, prenom, login, mot_de_passe) VALUES ($1, $2, $3, $4) RETURNING id`,
		b.Nom, b.Prenom, b.Login, b.MotDePasse,
	).Scan(&newID)
	if err != nil {
		return 0, fmt.Errorf("Create bibliothecaire: %w", err)
	}
	return newID, nil
}

func (r *BibliothécaireRepository) Update(b *models.Bibliothecaire) error {
	n, err := r.dbo.Exec(
		`UPDATE bibliothecaires SET nom=$1, prenom=$2, login=$3, mot_de_passe=$4 WHERE id=$5`,
		b.Nom, b.Prenom, b.Login, b.MotDePasse, b.Id,
	)
	if err != nil {
		return fmt.Errorf("Update bibliothecaire: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("bibliothécaire %d introuvable", b.Id)
	}
	return nil
}

func (r *BibliothécaireRepository) Delete(id int) error {
	n, err := r.dbo.Exec(`DELETE FROM bibliothecaires WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("Delete bibliothecaire: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("bibliothécaire %d introuvable", id)
	}
	return nil
}