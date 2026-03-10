package repositories

import (
	"database/sql"
	"errors"
	"fmt"

	"bibliotheque/db"
	"bibliotheque/models"
)

type EtudiantRepository struct {
	dbo *db.DBO
}

func NewEtudiantRepository(dbo *db.DBO) *EtudiantRepository {
	return &EtudiantRepository{dbo: dbo}
}

func (r *EtudiantRepository) FindAll() ([]*models.Etudiant, error) {
	rows, err := r.dbo.QueryRows(`
		SELECT u.id, u.nom, u.prenom, u.numero_telephone, u.solde_caution,
		       u.login, u.mot_de_passe, u.date_de_naissance, u.email, e.annee_etude
		FROM etudiants e
		JOIN utilisateurs u ON u.id = e.id
		ORDER BY u.nom`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var etudiants []*models.Etudiant
	for rows.Next() {
		et, err := scanEtudiant(rows)
		if err != nil {
			return nil, fmt.Errorf("FindAll etudiant scan: %w", err)
		}
		etudiants = append(etudiants, et)
	}
	return etudiants, rows.Err()
}

func (r *EtudiantRepository) FindByID(id int) (*models.Etudiant, error) {
	row := r.dbo.QueryRow(`
		SELECT u.id, u.nom, u.prenom, u.numero_telephone, u.solde_caution,
		       u.login, u.mot_de_passe, u.date_de_naissance, u.email, e.annee_etude
		FROM etudiants e
		JOIN utilisateurs u ON u.id = e.id
		WHERE e.id = $1`, id)

	et := &models.Etudiant{}
	err := row.Scan(&et.Id, &et.Nom, &et.Prenom, &et.NumeroTelephone, &et.SoldeCaution,
		&et.Login, &et.MotDePasse, &et.DateDeNaissance, &et.Email, &et.AnneeEtude)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("étudiant %d introuvable", id)
	}
	if err != nil {
		return nil, fmt.Errorf("FindByID etudiant: %w", err)
	}
	return et, nil
}

func (r *EtudiantRepository) Create(et *models.Etudiant) (int, error) {
	var newID int
	err := r.dbo.WithTx(func(tx *db.TxDBO) error {
		if err := tx.ExecReturning(`
			INSERT INTO utilisateurs (nom, prenom, numero_telephone, solde_caution, login, mot_de_passe, date_de_naissance, email)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`,
			et.Nom, et.Prenom, et.NumeroTelephone, et.SoldeCaution,
			et.Login, et.MotDePasse, et.DateDeNaissance, et.Email,
		).Scan(&newID); err != nil {
			return fmt.Errorf("insert utilisateur: %w", err)
		}
		_, err := tx.Exec(
			`INSERT INTO etudiants (id, annee_etude) VALUES ($1, $2)`,
			newID, et.AnneeEtude,
		)
		return err
	})
	if err != nil {
		return 0, fmt.Errorf("Create etudiant: %w", err)
	}
	return newID, nil
}

func (r *EtudiantRepository) Update(et *models.Etudiant) error {
	return r.dbo.WithTx(func(tx *db.TxDBO) error {
		if _, err := tx.Exec(`
			UPDATE utilisateurs SET nom=$1, prenom=$2, numero_telephone=$3, solde_caution=$4,
			login=$5, mot_de_passe=$6, date_de_naissance=$7, email=$8 WHERE id=$9`,
			et.Nom, et.Prenom, et.NumeroTelephone, et.SoldeCaution,
			et.Login, et.MotDePasse, et.DateDeNaissance, et.Email, et.Id,
		); err != nil {
			return err
		}
		_, err := tx.Exec(`UPDATE etudiants SET annee_etude=$1 WHERE id=$2`, et.AnneeEtude, et.Id)
		return err
	})
}

func (r *EtudiantRepository) Delete(id int) error {
	return r.dbo.WithTx(func(tx *db.TxDBO) error {
		if _, err := tx.Exec(`DELETE FROM etudiants WHERE id = $1`, id); err != nil {
			return err
		}
		_, err := tx.Exec(`DELETE FROM utilisateurs WHERE id = $1`, id)
		return err
	})
}

func scanEtudiant(rows interface {
	Scan(...any) error
}) (*models.Etudiant, error) {
	et := &models.Etudiant{}
	err := rows.Scan(&et.Id, &et.Nom, &et.Prenom, &et.NumeroTelephone, &et.SoldeCaution,
		&et.Login, &et.MotDePasse, &et.DateDeNaissance, &et.Email, &et.AnneeEtude)
	return et, err
}