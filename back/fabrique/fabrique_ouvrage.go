package fabrique

import "bibliotheque/models"

// FabriqueOuvrage est la fabrique abstraite.
type FabriqueOuvrage interface {
	CreerOuvrage() models.IOuvrage
}
