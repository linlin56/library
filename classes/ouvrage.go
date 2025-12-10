package classes

type Ouvrage struct {
	Id          int
	Caution     float64
	Titre       string
	Exemplaires int // TODO créer la classe Exemplaire
}

func NewOuvrage(id int, caution float64, titre string, exemplaires int) *Ouvrage {
	return &Ouvrage{
		Id:          id,
		Caution:     caution,
		Titre:       titre,
		Exemplaires: exemplaires,
	}
}
