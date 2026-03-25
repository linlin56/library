package models

// IOuvrage est l'interface produit de la fabrique abstraite.
// Livre et Revue la satisfont automatiquement par embedding de Ouvrage.
type IOuvrage interface {
	GetId() int
	GetTitre() string
	GetCaution() float64
}

type Ouvrage struct {
	Id          int     `json:"id"`
	Caution     float64 `json:"caution"`
	Titre       string  `json:"titre"`
	Exemplaires int     `json:"exemplaires,omitempty"`
}

func (o *Ouvrage) GetId() int          { return o.Id }
func (o *Ouvrage) GetTitre() string    { return o.Titre }
func (o *Ouvrage) GetCaution() float64 { return o.Caution }

func NewOuvrage(id int, caution float64, titre string, exemplaires int) *Ouvrage {
	return &Ouvrage{
		Id:          id,
		Caution:     caution,
		Titre:       titre,
		Exemplaires: exemplaires,
	}
}
