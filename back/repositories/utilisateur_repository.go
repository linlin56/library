package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"bibliotheque/db"
	"bibliotheque/models"
)

type UtilisateurRepository struct {
	dbo *db.DBO
}

func NewUtilisateurRepository(dbo *db.DBO) *UtilisateurRepository {
	return &UtilisateurRepository{dbo: dbo}
}

func (r *UtilisateurRepository) FindAll() ([]*models.Utilisateur, error) {
	rows, err := r.dbo.QueryRows(`
		SELECT id, nom, prenom, numero_telephone, solde_caution, login, mot_de_passe, date_de_naissance, email
		FROM utilisateurs ORDER BY nom`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var utilisateurs []*models.Utilisateur
	for rows.Next() {
		u, err := scanUtilisateur(rows)
		if err != nil {
			return nil, err
		}
		utilisateurs = append(utilisateurs, u)
	}
	return utilisateurs, rows.Err()
}

func (r *UtilisateurRepository) FindByID(id int) (*models.Utilisateur, error) {
	row := r.dbo.QueryRow(`
		SELECT id, nom, prenom, numero_telephone, solde_caution, login, mot_de_passe, date_de_naissance, email
		FROM utilisateurs WHERE id = $1`, id)

	u := &models.Utilisateur{}
	err := row.Scan(&u.Id, &u.Nom, &u.Prenom, &u.NumeroTelephone,
		&u.SoldeCaution, &u.Login, &u.MotDePasse, &u.DateDeNaissance, &u.Email)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("utilisateur %d introuvable", id)
	}
	if err != nil {
		return nil, fmt.Errorf("FindByID utilisateur: %w", err)
	}
	return u, nil
}

func (r *UtilisateurRepository) FindByLogin(login string) (*models.Utilisateur, error) {
	// UNION ALL sur toutes les tables héritées.
	// "FROM ONLY utilisateurs" exclut les enfants, on les ajoute explicitement
	// pour éviter que les colonnes supplémentaires (annee_etude, departement_id…)
	// ne fassent échouer le Scan.
	row := r.dbo.QueryRow(`
		SELECT id, nom, prenom, numero_telephone, solde_caution, login, mot_de_passe, date_de_naissance, email
		FROM ONLY utilisateurs WHERE login = $1
		UNION ALL
		SELECT id, nom, prenom, numero_telephone, solde_caution, login, mot_de_passe, date_de_naissance, email
		FROM etudiants WHERE login = $1
		UNION ALL
		SELECT id, nom, prenom, numero_telephone, solde_caution, login, mot_de_passe, date_de_naissance, email
		FROM enseignants WHERE login = $1
		LIMIT 1`, login)

	u := &models.Utilisateur{}
	err := row.Scan(&u.Id, &u.Nom, &u.Prenom, &u.NumeroTelephone,
		&u.SoldeCaution, &u.Login, &u.MotDePasse, &u.DateDeNaissance, &u.Email)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("login '%s' introuvable", login)
	}
	if err != nil {
		return nil, fmt.Errorf("FindByLogin: %w", err)
	}
	return u, nil
}

func (r *UtilisateurRepository) Create(u *models.Utilisateur) (int, error) {
	var newID int
	err := r.dbo.ExecReturning(`
		INSERT INTO utilisateurs (nom, prenom, numero_telephone, solde_caution, login, mot_de_passe, date_de_naissance, email)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`,
		u.Nom, u.Prenom, u.NumeroTelephone, u.SoldeCaution,
		u.Login, u.MotDePasse, u.DateDeNaissance, u.Email,
	).Scan(&newID)
	if err != nil {
		return 0, fmt.Errorf("Create utilisateur: %w", err)
	}
	return newID, nil
}

func (r *UtilisateurRepository) Update(u *models.Utilisateur) error {
	n, err := r.dbo.Exec(`
		UPDATE utilisateurs SET nom=$1, prenom=$2, numero_telephone=$3, solde_caution=$4,
		login=$5, mot_de_passe=$6, date_de_naissance=$7, email=$8 WHERE id=$9`,
		u.Nom, u.Prenom, u.NumeroTelephone, u.SoldeCaution,
		u.Login, u.MotDePasse, u.DateDeNaissance, u.Email, u.Id,
	)
	if err != nil {
		return fmt.Errorf("Update utilisateur: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("utilisateur %d introuvable", u.Id)
	}
	return nil
}

func (r *UtilisateurRepository) Delete(id int) error {
	n, err := r.dbo.Exec(`DELETE FROM utilisateurs WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("Delete utilisateur: %w", err)
	}
	if n == 0 {
		return fmt.Errorf("utilisateur %d introuvable", id)
	}
	return nil
}

// FindEmpruntsActifs retourne les exemplaires actuellement empruntés par un utilisateur
func (r *UtilisateurRepository) FindEmpruntsActifs(utilisateurID int) ([]*models.Exemplaire, error) {
	rows, err := r.dbo.QueryRows(`
		SELECT id, est_emprunte, date_debut_emprunt, date_fin_emprunt, delai_emprunt_jours, code_barre, emprunteur_id
		FROM exemplaires WHERE emprunteur_id = $1 AND est_emprunte = true`, utilisateurID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exemplaires []*models.Exemplaire
	for rows.Next() {
		e := &models.Exemplaire{}
		if err := rows.Scan(&e.Id, &e.EstEmprunte, &e.DateDebutEmprunt, &e.DateFinEmprunt,
			&e.DelaiEmpruntJours, &e.CodeBarre, &e.EmprunteurId); err != nil {
			return nil, fmt.Errorf("FindEmpruntsActifs scan: %w", err)
		}
		exemplaires = append(exemplaires, e)
	}
	return exemplaires, rows.Err()
}

// UtilisateurResume est utilisé dans la recherche d'utilisateurs par les bibliothécaires.
type UtilisateurResume struct {
	Id              int    `json:"id"`
	Nom             string `json:"nom"`
	Prenom          string `json:"prenom"`
	NumeroTelephone string `json:"numero_telephone"`
	Role            string `json:"role"`
}

// RechercherUtilisateurs cherche des utilisateurs par nom, prénom, code postal et/ou numéro de téléphone.
// Au moins un critère non vide est attendu (validé dans le controller).
func (r *UtilisateurRepository) RechercherUtilisateurs(nom, prenom, codePostal, numeroTelephone string) ([]*UtilisateurResume, error) {
	const q = `
		SELECT u.id, u.nom, u.prenom, COALESCE(u.numero_telephone, ''), 'etudiant'
		FROM etudiants u
		LEFT JOIN adresses a ON a.id = u.adresse_id
		WHERE ($1 = '' OR u.nom ILIKE '%' || $1 || '%')
		  AND ($2 = '' OR u.prenom ILIKE '%' || $2 || '%')
		  AND ($3 = '' OR a.code_postal = $3)
		  AND ($4 = '' OR COALESCE(u.numero_telephone, '') ILIKE '%' || $4 || '%')
		UNION ALL
		SELECT u.id, u.nom, u.prenom, COALESCE(u.numero_telephone, ''), 'enseignant'
		FROM enseignants u
		LEFT JOIN adresses a ON a.id = u.adresse_id
		WHERE ($1 = '' OR u.nom ILIKE '%' || $1 || '%')
		  AND ($2 = '' OR u.prenom ILIKE '%' || $2 || '%')
		  AND ($3 = '' OR a.code_postal = $3)
		  AND ($4 = '' OR COALESCE(u.numero_telephone, '') ILIKE '%' || $4 || '%')
		UNION ALL
		SELECT u.id, u.nom, u.prenom, COALESCE(u.numero_telephone, ''), 'utilisateur'
		FROM ONLY utilisateurs u
		LEFT JOIN adresses a ON a.id = u.adresse_id
		WHERE ($1 = '' OR u.nom ILIKE '%' || $1 || '%')
		  AND ($2 = '' OR u.prenom ILIKE '%' || $2 || '%')
		  AND ($3 = '' OR a.code_postal = $3)
		  AND ($4 = '' OR COALESCE(u.numero_telephone, '') ILIKE '%' || $4 || '%')
		ORDER BY nom, prenom`

	rows, err := r.dbo.QueryRows(q, nom, prenom, codePostal, numeroTelephone)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]*UtilisateurResume, 0)
	for rows.Next() {
		u := &UtilisateurResume{}
		if err := rows.Scan(&u.Id, &u.Nom, &u.Prenom, &u.NumeroTelephone, &u.Role); err != nil {
			return nil, fmt.Errorf("RechercherUtilisateurs scan: %w", err)
		}
		results = append(results, u)
	}
	return results, rows.Err()
}

func (r *UtilisateurRepository) LoginExists(login string) (bool, error) {
	var count int
	err := r.dbo.QueryRow(`
		SELECT COUNT(*) FROM (
			SELECT id FROM ONLY utilisateurs WHERE login = $1
			UNION ALL
			SELECT id FROM etudiants WHERE login = $1
			UNION ALL
			SELECT id FROM enseignants WHERE login = $1
			UNION ALL
			SELECT id FROM bibliothecaires WHERE login = $1
		) AS t`, login).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("LoginExists: %w", err)
	}
	return count > 0, nil
}

// --- helper interne ---

func scanUtilisateur(rows interface {
	Scan(...any) error
}) (*models.Utilisateur, error) {
	u := &models.Utilisateur{}
	var ddn time.Time
	err := rows.Scan(&u.Id, &u.Nom, &u.Prenom, &u.NumeroTelephone,
		&u.SoldeCaution, &u.Login, &u.MotDePasse, &ddn, &u.Email)
	u.DateDeNaissance = ddn
	return u, err
}