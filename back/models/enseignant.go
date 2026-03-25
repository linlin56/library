package models

type Enseignant struct {
	Utilisateur
	DepartementId int               `json:"departement_id"`
	Departement   *DepartementEcole `json:"departement,omitempty"`
}

func NewEnseignant(nom string, prenom string, login string, motDePasse string, departement *DepartementEcole) *Enseignant {
	return &Enseignant{
		Utilisateur: Utilisateur{
			Nom:        nom,
			Prenom:     prenom,
			Login:      login,
			MotDePasse: motDePasse,
		},
		Departement: departement,
	}
}
