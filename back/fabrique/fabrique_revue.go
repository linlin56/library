package fabrique

import (
	"time"

	"bibliotheque/models"
)

// FabriqueRevue est la fabrique concrète pour les revues.
type FabriqueRevue struct {
	Titre        string
	Caution      float64
	Numero       int
	DateParution time.Time
}

// CreerOuvrage implémente FabriqueOuvrage et retourne un *models.Revue.
func (f *FabriqueRevue) CreerOuvrage() models.IOuvrage {
	return models.NewRevue(0, f.Caution, f.Titre, 0, f.Numero, f.DateParution)
}
