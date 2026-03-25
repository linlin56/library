package models

import "time"

type Utilisateur struct {
	Id              int           `json:"id"`
	Nom             string        `json:"nom"`
	Prenom          string        `json:"prenom"`
	NumeroTelephone string        `json:"numero_telephone"`
	SoldeCaution    float64       `json:"solde_caution"`
	Login           string        `json:"login"`
	MotDePasse      string        `json:"mot_de_passe"`
	DateDeNaissance time.Time     `json:"date_de_naissance"`
	Email           string        `json:"email"`
	AdresseId       *int          `json:"adresse_id,omitempty"`
	Emprunts        []*Exemplaire `json:"emprunts,omitempty"`
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
