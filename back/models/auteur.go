package models

// Note (valable pour tous les modèles)
// Le `json` sert à définir en camel_case quand l'API envoie des JSON pour le front
type Auteur struct {
	Id     int    `json:"id"`
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
