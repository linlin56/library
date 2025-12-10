package models

type Enseignant struct {
	Utilisateur
	DepartementId int `gorm:"foreignKey"`
	Departement   *DepartementEcole
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
