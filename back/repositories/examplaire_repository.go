package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"bibliotheque/db"
	"bibliotheque/models"
)

type ExemplaireRepository struct {
	dbo *db.DBO
}

func NewExemplaireRepository(dbo *db.DBO) *ExemplaireRepository {
	return &ExemplaireRepository{dbo: dbo}
}

func (r *ExemplaireRepository) FindAll() ([]*models.Exemplaire, error) {
	rows, err := r.dbo.QueryRows(`
		SELECT id, est_emprunte, date_debut_emprunt, date_fin_emprunt,
		       delai_emprunt_jours, code_barre, emprunteur_id
		FROM exemplaires ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exemplaires []*models.Exemplaire
	for rows.Next() {
		e, err := scanExemplaire(rows)
		if err != nil {
			return nil, fmt.Errorf("FindAll exemplaire scan: %w", err)
		}
		exemplaires = append(exemplaires, e)
	}
	return exemplaires, rows.Err()
}

func (r *ExemplaireRepository) FindByID(id int) (*models.Exemplaire, error) {
	row := r.dbo.QueryRow(`
		SELECT id, est_emprunte, date_debut_emprunt, date_fin_emprunt,
		       delai_emprunt_jours, code_barre, emprunteur_id
		FROM exemplaires WHERE id = $1`, id)

	e := &models.Exemplaire{}
	var emprunteurID sql.NullInt64
	err := row.Scan(&e.Id, &e.EstEmprunte, &e.DateDebutEmprunt, &e.DateFinEmprunt,
		&e.DelaiEmpruntJours, &e.CodeBarre, &emprunteurID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("exemplaire %d introuvable", id)
	}
	if err != nil {
		return nil, fmt.Errorf("FindByID exemplaire: %w", err)
	}
	if emprunteurID.Valid {
		e.EmprunteurId = int(emprunteurID.Int64)
	}
	return e, nil
}

func (r *ExemplaireRepository) FindByCodeBarre(codeBarre string) (*models.Exemplaire, error) {
	row := r.dbo.QueryRow(`
		SELECT id, est_emprunte, date_debut_emprunt, date_fin_emprunt,
		       delai_emprunt_jours, code_barre, emprunteur_id
		FROM exemplaires WHERE code_barre = $1`, codeBarre)

	e := &models.Exemplaire{}
	var emprunteurID sql.NullInt64
	err := row.Scan(&e.Id, &e.EstEmprunte, &e.DateDebutEmprunt, &e.DateFinEmprunt,
		&e.DelaiEmpruntJours, &e.CodeBarre, &emprunteurID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("code barre '%s' introuvable", codeBarre)
	}
	if err != nil {
		return nil, fmt.Errorf("FindByCodeBarre: %w", err)
	}
	if emprunteurID.Valid {
		e.EmprunteurId = int(emprunteurID.Int64)
	}
	return e, nil
}

func (r *ExemplaireRepository) FindDisponibles() ([]*models.Exemplaire, error) {
	rows, err := r.dbo.QueryRows(`
		SELECT id, est_emprunte, date_debut_emprunt, date_fin_emprunt,
		       delai_emprunt_jours, code_barre, emprunteur_id
		FROM exemplaires WHERE est_emprunte = false ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exemplaires []*models.Exemplaire
	for rows.Next() {
		e, err := scanExemplaire(rows)
		if err != nil {
			return nil, fmt.Errorf("FindDisponibles scan: %w", err)
		}
		exemplaires = append(exemplaires, e)
	}
	return exemplaires, rows.Err()
}

func (r *ExemplaireRepository) FindDisponiblesByOuvrageId(ouvrageId int) ([]*models.Exemplaire, error) {
	rows, err := r.dbo.QueryRows(`
		SELECT id, est_emprunte, code_barre, delai_emprunt_jours
		FROM exemplaires
		WHERE ouvrage_id = $1 AND est_emprunte = false
		ORDER BY code_barre`, ouvrageId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []*models.Exemplaire
	for rows.Next() {
		e := &models.Exemplaire{}
		if err := rows.Scan(&e.Id, &e.EstEmprunte, &e.CodeBarre, &e.DelaiEmpruntJours); err != nil {
			return nil, fmt.Errorf("FindDisponiblesByOuvrageId scan: %w", err)
		}
		result = append(result, e)
	}
	return result, rows.Err()
}

func (r *ExemplaireRepository) Create(e *models.Exemplaire) (int, error) {
	var newID int
	err := r.dbo.ExecReturning(
		`INSERT INTO exemplaires (est_emprunte, date_debut_emprunt, date_fin_emprunt, delai_emprunt_jours, code_barre, emprunteur_id)
		 VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`,
		e.EstEmprunte, e.DateDebutEmprunt, e.DateFinEmprunt,
		e.DelaiEmpruntJours, e.CodeBarre, nullableInt(e.EmprunteurId),
	).Scan(&newID)
	if err != nil {
		return 0, fmt.Errorf("Create exemplaire: %w", err)
	}
	return newID, nil
}

func (r *ExemplaireRepository) Update(e *models.Exemplaire) error {
	n, err := r.dbo.Exec(`
		UPDATE exemplaires SET est_emprunte=$1, date_debut_emprunt=$2, date_fin_emprunt=$3,
		delai_emprunt_jours=$4, code_barre=$5, emprunteur_id=$6 WHERE id=$7`,
		e.EstEmprunte, e.DateDebutEmprunt, e.DateFinEmprunt,
		e.DelaiEmpruntJours, e.CodeBarre, nullableInt(e.EmprunteurId), e.Id,
	)
	if err != nil {
		return fmt.Errorf("Update exemplaire: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("exemplaire %d introuvable", e.Id)
	}
	return nil
}

// Emprunter marque un exemplaire comme emprunté par un utilisateur
func (r *ExemplaireRepository) Emprunter(exemplaireID int, emprunteurID int, delaiJours int) error {
	debut := time.Now()
	fin := debut.AddDate(0, 0, delaiJours)
	n, err := r.dbo.Exec(`
		UPDATE exemplaires SET est_emprunte=true, date_debut_emprunt=$1, date_fin_emprunt=$2,
		emprunteur_id=$3, delai_emprunt_jours=$4 WHERE id=$5 AND est_emprunte=false`,
		debut, fin, emprunteurID, delaiJours, exemplaireID,
	)
	if err != nil {
		return fmt.Errorf("Emprunter: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("exemplaire %d introuvable ou déjà emprunté", exemplaireID)
	}
	return nil
}

// Retourner marque un exemplaire comme rendu
func (r *ExemplaireRepository) Retourner(exemplaireID int) error {
	n, err := r.dbo.Exec(`
		UPDATE exemplaires SET est_emprunte=false, emprunteur_id=NULL,
		date_debut_emprunt=NULL, date_fin_emprunt=NULL WHERE id=$1`,
		exemplaireID,
	)
	if err != nil {
		return fmt.Errorf("Retourner: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("exemplaire %d introuvable", exemplaireID)
	}
	return nil
}

func (r *ExemplaireRepository) Delete(id int) error {
	n, err := r.dbo.Exec(`DELETE FROM exemplaires WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("Delete exemplaire: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("exemplaire %d introuvable", id)
	}
	return nil
}

// --- helpers ---

func scanExemplaire(rows interface {
	Scan(...any) error
}) (*models.Exemplaire, error) {
	e := &models.Exemplaire{}
	var emprunteurID sql.NullInt64
	err := rows.Scan(&e.Id, &e.EstEmprunte, &e.DateDebutEmprunt, &e.DateFinEmprunt,
		&e.DelaiEmpruntJours, &e.CodeBarre, &emprunteurID)
	if emprunteurID.Valid {
		e.EmprunteurId = int(emprunteurID.Int64)
	}
	return e, err
}

// nullableInt convertit 0 en NULL pour les clés étrangères optionnelles
func nullableInt(id int) sql.NullInt64 {
	if id == 0 {
		return sql.NullInt64{Valid: false}
	}
	return sql.NullInt64{Int64: int64(id), Valid: true}
}