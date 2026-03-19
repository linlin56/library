package models

import "time"

type Utilisateur struct {
	Id              int `gorm:"primaryKey;autoIncrement"`
	Nom             string
	Prenom          string
	NumeroTelephone string
	SoldeCaution    float64
	Login           string
	MotDePasse      string
	DateDeNaissance time.Time
	Email           string
	AdresseId       *int
	Emprunts        []*Exemplaire `gorm:"foreignKey:EmprunteurId"`
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
