package repositories

import (
	"database/sql"
	"errors"
	"fmt"

	"bibliotheque/db"
	"bibliotheque/models"
)

type RevueRepository struct {
	dbo *db.DBO
}

func NewRevueRepository(dbo *db.DBO) *RevueRepository {
	return &RevueRepository{dbo: dbo}
}

func (r *RevueRepository) FindAll() ([]*models.Revue, error) {
	rows, err := r.dbo.QueryRows(`
		SELECT rv.id, o.caution, o.titre, o.exemplaires, rv.numero, rv.date_parution
		FROM revues rv
		JOIN ouvrages o ON o.id = rv.id
		ORDER BY o.titre, rv.numero`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var revues []*models.Revue
	for rows.Next() {
		rv := &models.Revue{}
		if err := rows.Scan(&rv.Id, &rv.Caution, &rv.Titre, &rv.Exemplaires, &rv.Numero, &rv.DateParution); err != nil {
			return nil, fmt.Errorf("FindAll revue scan: %w", err)
		}
		revues = append(revues, rv)
	}
	return revues, rows.Err()
}

func (r *RevueRepository) FindByID(id int) (*models.Revue, error) {
	rv := &models.Revue{}
	err := r.dbo.QueryRow(`
		SELECT rv.id, o.caution, o.titre, o.exemplaires, rv.numero, rv.date_parution
		FROM revues rv
		JOIN ouvrages o ON o.id = rv.id
		WHERE rv.id = $1`, id).
		Scan(&rv.Id, &rv.Caution, &rv.Titre, &rv.Exemplaires, &rv.Numero, &rv.DateParution)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("revue %d introuvable", id)
	}
	if err != nil {
		return nil, fmt.Errorf("FindByID revue: %w", err)
	}
	return rv, nil
}

func (r *RevueRepository) FindByTitre(titre string) ([]*models.Revue, error) {
	rows, err := r.dbo.QueryRows(`
		SELECT rv.id, o.caution, o.titre, o.exemplaires, rv.numero, rv.date_parution
		FROM revues rv
		JOIN ouvrages o ON o.id = rv.id
		WHERE o.titre ILIKE $1
		ORDER BY o.titre, rv.numero`, "%"+titre+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var revues []*models.Revue
	for rows.Next() {
		rv := &models.Revue{}
		if err := rows.Scan(&rv.Id, &rv.Caution, &rv.Titre, &rv.Exemplaires, &rv.Numero, &rv.DateParution); err != nil {
			return nil, fmt.Errorf("FindByTitre revue scan: %w", err)
		}
		revues = append(revues, rv)
	}
	return revues, rows.Err()
}

func (r *RevueRepository) Create(rv *models.Revue) (int, error) {
	var newID int
	err := r.dbo.WithTx(func(tx *db.TxDBO) error {
		if err := tx.ExecReturning(
			`INSERT INTO ouvrages (caution, titre, exemplaires) VALUES ($1, $2, $3) RETURNING id`,
			rv.Caution, rv.Titre, rv.Exemplaires,
		).Scan(&newID); err != nil {
			return fmt.Errorf("insert ouvrage: %w", err)
		}
		_, err := tx.Exec(
			`INSERT INTO revues (id, numero) VALUES ($1, $2)`,
			newID, rv.Numero,
		)
		return err
	})
	if err != nil {
		return 0, fmt.Errorf("Create revue: %w", err)
	}
	return newID, nil
}

func (r *RevueRepository) Update(rv *models.Revue) error {
	return r.dbo.WithTx(func(tx *db.TxDBO) error {
		if _, err := tx.Exec(
			`UPDATE ouvrages SET caution=$1, titre=$2, exemplaires=$3 WHERE id=$4`,
			rv.Caution, rv.Titre, rv.Exemplaires, rv.Id,
		); err != nil {
			return err
		}
		_, err := tx.Exec(`UPDATE revues SET numero=$1 WHERE id=$2`, rv.Numero, rv.Id)
		return err
	})
}

func (r *RevueRepository) Delete(id int) error {
	return r.dbo.WithTx(func(tx *db.TxDBO) error {
		if _, err := tx.Exec(`DELETE FROM revues WHERE id = $1`, id); err != nil {
			return err
		}
		_, err := tx.Exec(`DELETE FROM ouvrages WHERE id = $1`, id)
		return err
	})
}