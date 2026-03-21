package models

type Auteur struct {
	Id     int    `gorm:"primaryKey;autoIncrement" json:"id"`
	Nom    string `json:"nom"`
	Prenom string `json:"prenom"`
}

func NewAuteur(id int, nom string, prenom string) *Auteur {
	return &Auteur{
		Id:     id,
		Nom:    nom,
		Prenom: prenom,
	}
}
