package models

type Etudiant struct {
	Utilisateur
	AnneeEtude string
}

func NewEtudiant(nom string, prenom string, login string, motDePasse string, anneeEtude string) *Etudiant {
	return &Etudiant{
		Utilisateur: Utilisateur{
			Nom:        nom,
			Prenom:     prenom,
			Login:      login,
			MotDePasse: motDePasse,
		},
		AnneeEtude: anneeEtude,
	}
}
