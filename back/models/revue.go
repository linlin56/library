package models

import "time"

type Revue struct {
	Ouvrage
	Numero       int       `json:"numero,omitempty"`
	DateParution time.Time `json:"date_parution,omitempty"`
}

func NewRevue(id int, caution float64, titre string, exemplaires int, numero int, dateParution time.Time) *Revue {
	return &Revue{
		Ouvrage: Ouvrage{
			Id:          id,
			Caution:     caution,
			Titre:       titre,
			Exemplaires: exemplaires,
		},
		Numero:       numero,
		DateParution: dateParution,
	}
}
