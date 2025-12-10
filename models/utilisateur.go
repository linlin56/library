package models

import "time"

type Utilisateur struct {
	Id              int
	Nom             string
	Prenom          string
	NumeroTelephone string
	SoldeCaution    float64
	Login           string
	MotDePasse      string
	DateDeNaissance time.Time
	Email           string
	Emprunts        []*Exemplaire
}

func NewUtilisateur(id int, nom string, prenom string, numeroTelephone string, soldeCaution float64, login string, motDePasse string, dateDeNaissance time.Time, email string) *Utilisateur {
	return &Utilisateur{
		Id:              id,
		Nom:             nom,
		Prenom:          prenom,
		NumeroTelephone: numeroTelephone,
		SoldeCaution:    soldeCaution,
		Login:           login,
		MotDePasse:      motDePasse,
		DateDeNaissance: dateDeNaissance,
		Email:           email,
	}
}
