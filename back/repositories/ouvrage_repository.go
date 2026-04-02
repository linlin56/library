package repositories

import (
	"database/sql"
	"errors"
	"fmt"

	"bibliotheque/db"
	"bibliotheque/models"
)

// FiltresRecherche contient les criteres de la recherche avancée.
type FiltresRecherche struct {
Titre      string
Auteur     string
Isbn       string
CodeBarre  string
CodeRevue  string
Disponible bool
}

// OuvrageResultat est la vue enrichie retournee par RechercheAvancee.
type OuvrageResultat struct {
Id                    int     `json:"id"`
Titre                 string  `json:"titre"`
Caution               float64 `json:"caution"`
Type                  string  `json:"type"`
Isbn                  string  `json:"isbn,omitempty"`
Auteur                string  `json:"auteur,omitempty"`
Numero                *int    `json:"numero,omitempty"`
DateParution          string  `json:"date_parution,omitempty"`
ExemplairesDisponibles int    `json:"exemplaires_disponibles"`
Emplacement           string  `json:"emplacement,omitempty"`
}

const rechercheSQL = `
SELECT
    l.id, l.caution, l.titre,
    l.isbn,
    COALESCE(a.prenom || ' ' || a.nom, '') AS auteur,
    NULL::int AS numero,
    NULL::text AS date_parution,
    'livre' AS type,
    (SELECT COUNT(*) FROM exemplaires e WHERE e.ouvrage_id = l.id AND NOT e.est_emprunte) AS exemplaires_disponibles,
    COALESCE('travée ' || em.numero_travee || ', étagère ' || em.numero_etagere || ', niveau ' || em.niveau, '') AS emplacement
FROM livres l
LEFT JOIN auteurs a ON a.id = l.auteur_id
LEFT JOIN emplacements em ON em.id = l.emplacement_id
WHERE ($1 = '' OR l.titre ILIKE '%' || $1 || '%')
  AND ($2 = '' OR a.nom ILIKE '%' || $2 || '%' OR a.prenom ILIKE '%' || $2 || '%')
  AND ($3 = '' OR l.isbn ILIKE '%' || $3 || '%')
  AND ($4 = '' OR EXISTS(SELECT 1 FROM exemplaires e WHERE e.ouvrage_id = l.id AND e.code_barre = $4))
  AND (NOT $5 OR EXISTS(SELECT 1 FROM exemplaires e WHERE e.ouvrage_id = l.id AND NOT e.est_emprunte))
  AND $6 = ''
UNION ALL
SELECT
    r.id, r.caution, r.titre,
    NULL::varchar AS isbn,
    NULL::varchar AS auteur,
    r.numero,
    TO_CHAR(r.date_parution, 'YYYY-MM-DD') AS date_parution,
    'revue' AS type,
    (SELECT COUNT(*) FROM exemplaires e WHERE e.ouvrage_id = r.id AND NOT e.est_emprunte) AS exemplaires_disponibles,
    COALESCE('travée ' || em.numero_travee || ', étagère ' || em.numero_etagere || ', niveau ' || em.niveau, '') AS emplacement
FROM revues r
LEFT JOIN emplacements em ON em.id = r.emplacement_id
WHERE ($1 = '' OR r.titre ILIKE '%' || $1 || '%')
  AND ($4 = '' OR EXISTS(SELECT 1 FROM exemplaires e WHERE e.ouvrage_id = r.id AND e.code_barre = $4))
  AND ($6 = '' OR r.numero::text = $6)
  AND (NOT $5 OR EXISTS(SELECT 1 FROM exemplaires e WHERE e.ouvrage_id = r.id AND NOT e.est_emprunte))
  AND $2 = ''
  AND $3 = ''
ORDER BY titre
`

type OuvrageRepository struct {
dbo *db.DBO
}

func NewOuvrageRepository(dbo *db.DBO) *OuvrageRepository {
return &OuvrageRepository{dbo: dbo}
}

// Rechercher execute une recherche avancée avec les criteres fournis.
func (r *OuvrageRepository) Rechercher(f FiltresRecherche) ([]OuvrageResultat, error) {
rows, err := r.dbo.QueryRows(rechercheSQL,
f.Titre, f.Auteur, f.Isbn, f.CodeBarre, f.Disponible, f.CodeRevue,
)
if err != nil {
return nil, err
}
defer func() { _ = rows.Close() }()

resultats := make([]OuvrageResultat, 0)
for rows.Next() {
var res OuvrageResultat
var isbn, auteur, dateParution, emplacement sql.NullString
var numero sql.NullInt32
if err := rows.Scan(&res.Id, &res.Caution, &res.Titre, &isbn, &auteur, &numero, &dateParution, &res.Type, &res.ExemplairesDisponibles, &emplacement); err != nil {
return nil, fmt.Errorf("Rechercher scan: %w", err)
}
if isbn.Valid {
res.Isbn = isbn.String
}
if auteur.Valid {
res.Auteur = auteur.String
}
if numero.Valid {
n := int(numero.Int32)
res.Numero = &n
}
if dateParution.Valid {
res.DateParution = dateParution.String
}
if emplacement.Valid {
res.Emplacement = emplacement.String
}
resultats = append(resultats, res)
}
return resultats, rows.Err()
}

func (r *OuvrageRepository) FindAll() ([]*models.Ouvrage, error) {
rows, err := r.dbo.QueryRows(`SELECT id, caution, titre FROM ouvrages ORDER BY titre`)
if err != nil {
return nil, err
}
defer func() { _ = rows.Close() }()

ouvrages := make([]*models.Ouvrage, 0)
for rows.Next() {
o := &models.Ouvrage{}
if err := rows.Scan(&o.Id, &o.Caution, &o.Titre); err != nil {
return nil, fmt.Errorf("FindAll ouvrage scan: %w", err)
}
ouvrages = append(ouvrages, o)
}
return ouvrages, rows.Err()
}

func (r *OuvrageRepository) FindByID(id int) (*models.Ouvrage, error) {
o := &models.Ouvrage{}
err := r.dbo.QueryRow(`SELECT id, caution, titre FROM ouvrages WHERE id = $1`, id).
Scan(&o.Id, &o.Caution, &o.Titre)
if errors.Is(err, sql.ErrNoRows) {
return nil, fmt.Errorf("ouvrage %d introuvable", id)
}
if err != nil {
return nil, fmt.Errorf("FindByID ouvrage: %w", err)
}
return o, nil
}

func (r *OuvrageRepository) Create(o *models.Ouvrage) (int, error) {
var newID int
err := r.dbo.ExecReturning(`INSERT INTO ouvrages (caution, titre) VALUES ($1, $2) RETURNING id`, o.Caution, o.Titre).Scan(&newID)
if err != nil {
return 0, fmt.Errorf("Create ouvrage: %w", err)
}
return newID, nil
}

func (r *OuvrageRepository) Update(o *models.Ouvrage) error {
n, err := r.dbo.Exec(
`UPDATE ouvrages SET caution=$1, titre=$2 WHERE id=$3`,
o.Caution, o.Titre, o.Id,
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

func (r *OuvrageRepository) CreateLivre(titre string, caution float64, isbn string, auteurId int, emplacementId int) (int, error) {
	var newID int
	err := r.dbo.ExecReturning(
		`INSERT INTO livres (titre, caution, isbn, auteur_id, emplacement_id) VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		titre, caution, isbn, auteurId, emplacementId,
	).Scan(&newID)
	if err != nil {
		return 0, fmt.Errorf("CreateLivre: %w", err)
	}
	return newID, nil
}

func (r *OuvrageRepository) CreateRevue(titre string, caution float64, numero int, dateParution string, emplacementId int) (int, error) {
	var newID int
	err := r.dbo.ExecReturning(
		`INSERT INTO revues (titre, caution, numero, date_parution, emplacement_id) VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		titre, caution, numero, dateParution, emplacementId,
	).Scan(&newID)
	if err != nil {
		return 0, fmt.Errorf("CreateRevue: %w", err)
	}
	return newID, nil
}
