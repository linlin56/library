package fabrique

import "bibliotheque/models"

// FabriqueLivre est la fabrique concrète pour les livres.
type FabriqueLivre struct {
	Titre   string
	Caution float64
	Isbn    string
	Auteur  *models.Auteur
}

// CreerOuvrage implémente FabriqueOuvrage et retourne un *models.Livre.
func (f *FabriqueLivre) CreerOuvrage() models.IOuvrage {
	return models.NewLivre(0, f.Caution, f.Titre, 0, *f.Auteur, f.Isbn)
}
