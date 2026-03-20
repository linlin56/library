package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"bibliotheque/db"
)

// PreviewEmprunt contient les informations affichees avant confirmation.
type PreviewEmprunt struct {
Titre       string  `json:"titre"`
CodeBarre   string  `json:"code_barre"`
Caution     float64 `json:"caution"`
SoldeActuel float64 `json:"solde_actuel"`
NouveauSolde float64 `json:"nouveau_solde"`
}

type EmpruntRepository struct {
dbo *db.DBO
}

func NewEmpruntRepository(dbo *db.DBO) *EmpruntRepository {
return &EmpruntRepository{dbo: dbo}
}

// Verifier effectue les 3 verifications metier et renvoie un apercu ou une erreur explicite.
func (r *EmpruntRepository) Verifier(utilisateurId int, codeBarre string) (*PreviewEmprunt, error) {
// Recuperer l'exemplaire et l'ouvrage associe
var exId int
var ouvrageId int
var caution float64
var titre string
var estEmprunte bool
err := r.dbo.QueryRow(`
		SELECT e.id, e.ouvrage_id, e.est_emprunte, o.caution, o.titre
		FROM exemplaires e
		JOIN ouvrages o ON o.id = e.ouvrage_id
		WHERE e.code_barre = $1 LIMIT 1`, codeBarre).
Scan(&exId, &ouvrageId, &estEmprunte, &caution, &titre)
if errors.Is(err, sql.ErrNoRows) {
return nil, fmt.Errorf("Code barre introuvable.")
}
if err != nil {
return nil, fmt.Errorf("Erreur lors de la recherche de l'exemplaire: %w", err)
}
if estEmprunte {
return nil, fmt.Errorf("Cet exemplaire est deja emprunte.")
}

	// 1.1 - Verifier aucun emprunt en retard
	var nbRetard int
err = r.dbo.QueryRow(`
SELECT COUNT(*) FROM exemplaires
WHERE emprunteur_id = $1 AND est_emprunte = TRUE AND date_fin_emprunt < $2`,
utilisateurId, time.Now()).Scan(&nbRetard)
if err != nil {
return nil, fmt.Errorf("Erreur lors de la verification des emprunts en retard: %w", err)
}
if nbRetard > 0 {
return nil, fmt.Errorf("Vous avez %d emprunt(s) en retard. Veuillez les retourner avant d'emprunter.", nbRetard)
}

// 1.2 - Verifier qu'il n'emprunte pas deja un exemplaire du meme ouvrage
var existant int
err = r.dbo.QueryRow(`
		SELECT COUNT(*) FROM exemplaires
		WHERE emprunteur_id = $1 AND ouvrage_id = $2 AND est_emprunte = TRUE`,
utilisateurId, ouvrageId).Scan(&existant)
if err != nil {
return nil, fmt.Errorf("Erreur lors de la verification des emprunts existants: %w", err)
}
if existant > 0 {
return nil, fmt.Errorf("Vous empruntez deja un exemplaire de cet ouvrage.")
}

// 1.3 - Verifier le solde caution suffisant
var solde float64
err = r.dbo.QueryRow(`
		SELECT solde_caution FROM utilisateurs WHERE id = $1
		UNION ALL
		SELECT solde_caution FROM etudiants WHERE id = $1
		UNION ALL
		SELECT solde_caution FROM enseignants WHERE id = $1
		LIMIT 1`, utilisateurId).Scan(&solde)
if err != nil {
return nil, fmt.Errorf("Utilisateur introuvable.")
}
if solde < caution {
return nil, fmt.Errorf("Solde insuffisant (%.2f EUR disponible, %.2f EUR requis).", solde, caution)
}

return &PreviewEmprunt{
Titre:        titre,
CodeBarre:    codeBarre,
Caution:      caution,
SoldeActuel:  solde,
	NouveauSolde: solde - caution,
}, nil
}
// EmpruntItem est retourné dans la liste des emprunts actifs d'un utilisateur.
type EmpruntItem struct {
	Id        int    `json:"id"`
	CodeBarre string `json:"code_barre"`
	Titre     string `json:"titre"`
	DateDebut string `json:"date_debut"`
	DateFin   string `json:"date_fin"`
	EnRetard  bool   `json:"en_retard"`
}

// GetEmprunts retourne les emprunts actifs d'un utilisateur avec les informations de l'ouvrage.
func (r *EmpruntRepository) GetEmprunts(utilisateurId int) ([]*EmpruntItem, error) {
	rows, err := r.dbo.QueryRows(`
		SELECT e.id, e.code_barre, o.titre,
		       to_char(e.date_debut_emprunt, 'YYYY-MM-DD'),
		       to_char(e.date_fin_emprunt,   'YYYY-MM-DD'),
		       e.date_fin_emprunt < NOW()
		FROM exemplaires e
		JOIN ouvrages o ON o.id = e.ouvrage_id
		WHERE e.emprunteur_id = $1 AND e.est_emprunte = TRUE
		ORDER BY e.date_fin_emprunt ASC`, utilisateurId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*EmpruntItem, 0)
	for rows.Next() {
		item := &EmpruntItem{}
		if err := rows.Scan(&item.Id, &item.CodeBarre, &item.Titre,
			&item.DateDebut, &item.DateFin, &item.EnRetard); err != nil {
			return nil, fmt.Errorf("GetEmprunts scan: %w", err)
		}
		items = append(items, item)
	}
	return items, rows.Err()
}
// Emprunter enregistre l'emprunt apres validation de l'utilisateur.
func (r *EmpruntRepository) Emprunter(utilisateurId int, codeBarre string) error {
	// Re-verifier avant d'enregistrer
preview, err := r.Verifier(utilisateurId, codeBarre)
if err != nil {
return err
}

// Recuperer l'id de l'exemplaire et le delai
var exId, delaiJours int
err = r.dbo.QueryRow(`SELECT id, delai_emprunt_jours FROM exemplaires WHERE code_barre = $1`, codeBarre).
Scan(&exId, &delaiJours)
if err != nil {
return fmt.Errorf("Exemplaire introuvable: %w", err)
}

now := time.Now()
fin := now.AddDate(0, 0, delaiJours)

// Marquer l'exemplaire comme emprunte
n, err := r.dbo.Exec(`
		UPDATE exemplaires
		SET est_emprunte = TRUE,
		    date_debut_emprunt = $1,
		    date_fin_emprunt = $2,
		    emprunteur_id = $3
		WHERE id = $4`, now, fin, utilisateurId, exId)
if err != nil || n == 0 {
return fmt.Errorf("Impossible d'enregistrer l'emprunt.")
}

// Deduire la caution - essaie les 3 tables (heritage natif PostgreSQL)
	for _, table := range []string{"utilisateurs", "etudiants", "enseignants"} {
		q := fmt.Sprintf(`UPDATE %s SET solde_caution = solde_caution - $1 WHERE id = $2`, table)
		n, err = r.dbo.Exec(q, preview.Caution, utilisateurId)
		if err == nil && n > 0 {
			return nil
		}
	}
	return fmt.Errorf("Impossible de deduire la caution du solde.")
}
