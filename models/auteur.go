package models

type Auteur struct {
	Id     int
	Nom    string
	Prenom string
}

func NewAuteur(id int, nom string, prenom string) *Auteur {
	return &Auteur{
		Id:     id,
		Nom:    nom,
		Prenom: prenom,
	}
}
