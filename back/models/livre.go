package models

type Livre struct {
	Ouvrage
	AuteurId int     `json:"auteur_id,omitempty"`
	Auteur   *Auteur `json:"auteur,omitempty"`
	Isbn     string  `json:"isbn,omitempty"`
}

func NewLivre(id int, caution float64, titre string, exemplaires int, auteur Auteur, isbn string) *Livre {
	return &Livre{
		Ouvrage: Ouvrage{
			Id:          id,
			Caution:     caution,
			Titre:       titre,
			Exemplaires: exemplaires,
		},
		Auteur: &auteur,
		Isbn:   isbn,
	}
}
