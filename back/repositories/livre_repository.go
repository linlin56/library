package repositories

import (
	"database/sql"
	"errors"
	"fmt"

	"bibliotheque/db"
	"bibliotheque/models"
)

type LivreRepository struct {
	dbo *db.DBO
}

func NewLivreRepository(dbo *db.DBO) *LivreRepository {
	return &LivreRepository{dbo: dbo}
}

// FindAll retourne tous les livres avec leur auteur (JOIN)
func (r *LivreRepository) FindAll() ([]*models.Livre, error) {
	rows, err := r.dbo.QueryRows(`
		SELECT l.id, o.caution, o.titre, o.exemplaires, l.isbn,
		       a.id, a.nom, a.prenom
		FROM livres l
		JOIN ouvrages o ON o.id = l.id
		JOIN auteurs a ON a.id = l.auteur_id
		ORDER BY o.titre`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var livres []*models.Livre
	for rows.Next() {
		l, err := scanLivre(rows)
		if err != nil {
			return nil, fmt.Errorf("FindAll livre scan: %w", err)
		}
		livres = append(livres, l)
	}
	return livres, rows.Err()
}

func (r *LivreRepository) FindByID(id int) (*models.Livre, error) {
	row := r.dbo.QueryRow(`
		SELECT l.id, o.caution, o.titre, o.exemplaires, l.isbn,
		       a.id, a.nom, a.prenom
		FROM livres l
		JOIN ouvrages o ON o.id = l.id
		JOIN auteurs a ON a.id = l.auteur_id
		WHERE l.id = $1`, id)

	l := &models.Livre{}
	a := &models.Auteur{}
	err := row.Scan(&l.Id, &l.Caution, &l.Titre, &l.Exemplaires, &l.Isbn,
		&a.Id, &a.Nom, &a.Prenom)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("livre %d introuvable", id)
	}
	if err != nil {
		return nil, fmt.Errorf("FindByID livre: %w", err)
	}
	l.Auteur = a
	l.AuteurId = a.Id
	return l, nil
}

func (r *LivreRepository) FindByIsbn(isbn string) (*models.Livre, error) {
	row := r.dbo.QueryRow(`
		SELECT l.id, o.caution, o.titre, o.exemplaires, l.isbn,
		       a.id, a.nom, a.prenom
		FROM livres l
		JOIN ouvrages o ON o.id = l.id
		JOIN auteurs a ON a.id = l.auteur_id
		WHERE l.isbn = $1`, isbn)

	l := &models.Livre{}
	a := &models.Auteur{}
	err := row.Scan(&l.Id, &l.Caution, &l.Titre, &l.Exemplaires, &l.Isbn,
		&a.Id, &a.Nom, &a.Prenom)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("ISBN '%s' introuvable", isbn)
	}
	if err != nil {
		return nil, fmt.Errorf("FindByIsbn: %w", err)
	}
	l.Auteur = a
	l.AuteurId = a.Id
	return l, nil
}

func (r *LivreRepository) FindByAuteur(auteurID int) ([]*models.Livre, error) {
	rows, err := r.dbo.QueryRows(`
		SELECT l.id, o.caution, o.titre, o.exemplaires, l.isbn,
		       a.id, a.nom, a.prenom
		FROM livres l
		JOIN ouvrages o ON o.id = l.id
		JOIN auteurs a ON a.id = l.auteur_id
		WHERE l.auteur_id = $1
		ORDER BY o.titre`, auteurID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var livres []*models.Livre
	for rows.Next() {
		l, err := scanLivre(rows)
		if err != nil {
			return nil, fmt.Errorf("FindByAuteur scan: %w", err)
		}
		livres = append(livres, l)
	}
	return livres, rows.Err()
}

// Create insère dans ouvrages ET livres dans une transaction
func (r *LivreRepository) Create(l *models.Livre) (int, error) {
	var newID int
	err := r.dbo.WithTx(func(tx *db.TxDBO) error {
		// 1. Insérer dans la table parente ouvrages
		if err := tx.ExecReturning(
			`INSERT INTO ouvrages (caution, titre, exemplaires) VALUES ($1, $2, $3) RETURNING id`,
			l.Caution, l.Titre, l.Exemplaires,
		).Scan(&newID); err != nil {
			return fmt.Errorf("insert ouvrage: %w", err)
		}

		// 2. Insérer dans la table livres avec le même id
		_, err := tx.Exec(
			`INSERT INTO livres (id, auteur_id, isbn) VALUES ($1, $2, $3)`,
			newID, l.AuteurId, l.Isbn,
		)
		return err
	})
	if err != nil {
		return 0, fmt.Errorf("Create livre: %w", err)
	}
	return newID, nil
}

func (r *LivreRepository) Update(l *models.Livre) error {
	return r.dbo.WithTx(func(tx *db.TxDBO) error {
		if _, err := tx.Exec(
			`UPDATE ouvrages SET caution=$1, titre=$2, exemplaires=$3 WHERE id=$4`,
			l.Caution, l.Titre, l.Exemplaires, l.Id,
		); err != nil {
			return fmt.Errorf("update ouvrage: %w", err)
		}
		_, err := tx.Exec(
			`UPDATE livres SET auteur_id=$1, isbn=$2 WHERE id=$3`,
			l.AuteurId, l.Isbn, l.Id,
		)
		return err
	})
}

func (r *LivreRepository) Delete(id int) error {
	return r.dbo.WithTx(func(tx *db.TxDBO) error {
		if _, err := tx.Exec(`DELETE FROM livres WHERE id = $1`, id); err != nil {
			return err
		}
		_, err := tx.Exec(`DELETE FROM ouvrages WHERE id = $1`, id)
		return err
	})
}

// --- helper ---

func scanLivre(rows interface {
	Scan(...any) error
}) (*models.Livre, error) {
	l := &models.Livre{}
	a := &models.Auteur{}
	err := rows.Scan(&l.Id, &l.Caution, &l.Titre, &l.Exemplaires, &l.Isbn,
		&a.Id, &a.Nom, &a.Prenom)
	l.Auteur = a
	l.AuteurId = a.Id
	return l, err
}