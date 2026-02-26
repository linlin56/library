package repositories

import (
	"database/sql"
	"errors"
	"fmt"

	"bibliotheque/db"
	"bibliotheque/models"
)

type EnseignantRepository struct {
	dbo *db.DBO
}

func NewEnseignantRepository(dbo *db.DBO) *EnseignantRepository {
	return &EnseignantRepository{dbo: dbo}
}

func (r *EnseignantRepository) FindAll() ([]*models.Enseignant, error) {
	rows, err := r.dbo.QueryRows(`
		SELECT u.id, u.nom, u.prenom, u.numero_telephone, u.solde_caution,
		       u.login, u.mot_de_passe, u.date_de_naissance, u.email,
		       d.id, d.nom
		FROM enseignants en
		JOIN utilisateurs u ON u.id = en.id
		JOIN departements_ecole d ON d.id = en.departement_id
		ORDER BY u.nom`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var enseignants []*models.Enseignant
	for rows.Next() {
		en, err := scanEnseignant(rows)
		if err != nil {
			return nil, fmt.Errorf("FindAll enseignant scan: %w", err)
		}
		enseignants = append(enseignants, en)
	}
	return enseignants, rows.Err()
}

func (r *EnseignantRepository) FindByID(id int) (*models.Enseignant, error) {
	row := r.dbo.QueryRow(`
		SELECT u.id, u.nom, u.prenom, u.numero_telephone, u.solde_caution,
		       u.login, u.mot_de_passe, u.date_de_naissance, u.email,
		       d.id, d.nom
		FROM enseignants en
		JOIN utilisateurs u ON u.id = en.id
		JOIN departements_ecole d ON d.id = en.departement_id
		WHERE en.id = $1`, id)

	en := &models.Enseignant{}
	d := &models.DepartementEcole{}
	err := row.Scan(&en.Id, &en.Nom, &en.Prenom, &en.NumeroTelephone, &en.SoldeCaution,
		&en.Login, &en.MotDePasse, &en.DateDeNaissance, &en.Email,
		&d.Id, &d.Nom)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("enseignant %d introuvable", id)
	}
	if err != nil {
		return nil, fmt.Errorf("FindByID enseignant: %w", err)
	}
	en.Departement = d
	en.DepartementId = d.Id
	return en, nil
}

func (r *EnseignantRepository) FindByDepartement(departementID int) ([]*models.Enseignant, error) {
	rows, err := r.dbo.QueryRows(`
		SELECT u.id, u.nom, u.prenom, u.numero_telephone, u.solde_caution,
		       u.login, u.mot_de_passe, u.date_de_naissance, u.email,
		       d.id, d.nom
		FROM enseignants en
		JOIN utilisateurs u ON u.id = en.id
		JOIN departements_ecole d ON d.id = en.departement_id
		WHERE en.departement_id = $1
		ORDER BY u.nom`, departementID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var enseignants []*models.Enseignant
	for rows.Next() {
		en, err := scanEnseignant(rows)
		if err != nil {
			return nil, fmt.Errorf("FindByDepartement scan: %w", err)
		}
		enseignants = append(enseignants, en)
	}
	return enseignants, rows.Err()
}

func (r *EnseignantRepository) Create(en *models.Enseignant) (int, error) {
	var newID int
	err := r.dbo.WithTx(func(tx *db.TxDBO) error {
		if err := tx.ExecReturning(`
			INSERT INTO utilisateurs (nom, prenom, numero_telephone, solde_caution, login, mot_de_passe, date_de_naissance, email)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`,
			en.Nom, en.Prenom, en.NumeroTelephone, en.SoldeCaution,
			en.Login, en.MotDePasse, en.DateDeNaissance, en.Email,
		).Scan(&newID); err != nil {
			return fmt.Errorf("insert utilisateur: %w", err)
		}
		_, err := tx.Exec(
			`INSERT INTO enseignants (id, departement_id) VALUES ($1, $2)`,
			newID, en.DepartementId,
		)
		return err
	})
	if err != nil {
		return 0, fmt.Errorf("Create enseignant: %w", err)
	}
	return newID, nil
}

func (r *EnseignantRepository) Update(en *models.Enseignant) error {
	return r.dbo.WithTx(func(tx *db.TxDBO) error {
		if _, err := tx.Exec(`
			UPDATE utilisateurs SET nom=$1, prenom=$2, numero_telephone=$3, solde_caution=$4,
			login=$5, mot_de_passe=$6, date_de_naissance=$7, email=$8 WHERE id=$9`,
			en.Nom, en.Prenom, en.NumeroTelephone, en.SoldeCaution,
			en.Login, en.MotDePasse, en.DateDeNaissance, en.Email, en.Id,
		); err != nil {
			return err
		}
		_, err := tx.Exec(`UPDATE enseignants SET departement_id=$1 WHERE id=$2`, en.DepartementId, en.Id)
		return err
	})
}

func (r *EnseignantRepository) Delete(id int) error {
	return r.dbo.WithTx(func(tx *db.TxDBO) error {
		if _, err := tx.Exec(`DELETE FROM enseignants WHERE id = $1`, id); err != nil {
			return err
		}
		_, err := tx.Exec(`DELETE FROM utilisateurs WHERE id = $1`, id)
		return err
	})
}

func scanEnseignant(rows interface {
	Scan(...any) error
}) (*models.Enseignant, error) {
	en := &models.Enseignant{}
	d := &models.DepartementEcole{}
	err := rows.Scan(&en.Id, &en.Nom, &en.Prenom, &en.NumeroTelephone, &en.SoldeCaution,
		&en.Login, &en.MotDePasse, &en.DateDeNaissance, &en.Email,
		&d.Id, &d.Nom)
	en.Departement = d
	en.DepartementId = d.Id
	return en, err
}