package classes

type Livre struct {
	Ouvrage
	Auteur Auteur
	Isbn   string
}

func NewLivre(id int, caution float64, titre string, exemplaires int, auteur Auteur, isbn string) *Livre {
	return &Livre{
		Ouvrage: Ouvrage{
			Id:          id,
			Caution:     caution,
			Titre:       titre,
			Exemplaires: exemplaires,
		},
		Auteur: auteur,
		Isbn:   isbn,
	}
}
