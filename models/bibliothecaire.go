package models

type Bibliothecaire struct {
	Id         int
	Nom        string
	Prenom     string
	Login      string
	MotDePasse string
}

func NewBibliothecaire(id int, nom string, prenom string, login string, motDePasse string) *Bibliothecaire {
	return &Bibliothecaire{
		Id:         id,
		Nom:        nom,
		Prenom:     prenom,
		Login:      login,
		MotDePasse: motDePasse,
	}
}
