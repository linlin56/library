package repositories

import (
	"fmt"

	"bibliotheque/db"
)

type EmplacementResume struct {
	Id             int    `json:"id"`
	NumeroTravee   int    `json:"numero_travee"`
	NumeroEtagere  int    `json:"numero_etagere"`
	Niveau         int    `json:"niveau"`
	CategorieNom   string `json:"categorie_nom"`
}

type EmplacementRepository struct {
	dbo *db.DBO
}

func NewEmplacementRepository(dbo *db.DBO) *EmplacementRepository {
	return &EmplacementRepository{dbo: dbo}
}

func (r *EmplacementRepository) FindAll() ([]*EmplacementResume, error) {
	rows, err := r.dbo.QueryRows(`
		SELECT e.id, e.numero_travee, e.numero_etagere, e.niveau, COALESCE(c.nom, '')
		FROM emplacements e
		LEFT JOIN categories c ON c.id = e.categorie_id
		ORDER BY e.numero_travee, e.numero_etagere, e.niveau`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var emplacements []*EmplacementResume
	for rows.Next() {
		em := &EmplacementResume{}
		if err := rows.Scan(&em.Id, &em.NumeroTravee, &em.NumeroEtagere, &em.Niveau, &em.CategorieNom); err != nil {
			return nil, fmt.Errorf("FindAll emplacement scan: %w", err)
		}
		emplacements = append(emplacements, em)
	}
	return emplacements, rows.Err()
}
