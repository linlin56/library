package models

type Bibliothecaire struct {
	Id         int    `json:"id"`
	Nom        string `json:"nom"`
	Prenom     string `json:"prenom"`
	Login      string `json:"login"`
	MotDePasse string `json:"mot_de_passe"`
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
