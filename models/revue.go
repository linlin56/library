package models

type Revue struct {
	Ouvrage
	Numero int
}

func NewRevue(id int, caution float64, titre string, exemplaires int, numero int) *Revue {
	return &Revue{
		Ouvrage: Ouvrage{
			Id:          id,
			Caution:     caution,
			Titre:       titre,
			Exemplaires: exemplaires,
		},
		Numero: numero,
	}
}
