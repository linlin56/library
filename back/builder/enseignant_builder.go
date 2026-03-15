package builder

import (
	"bibliotheque/models"
	"time"
)

// EnseignantBuilder étend UtilisateurBuilder pour construire un Enseignant.
type EnseignantBuilder struct {
	UtilisateurBuilder
	departement *models.DepartementEcole
}

func NewEnseignantBuilder() *EnseignantBuilder {
	return &EnseignantBuilder{}
}

func (b *EnseignantBuilder) WithNom(nom string) *EnseignantBuilder {
	b.u.Nom = nom
	return b
}

func (b *EnseignantBuilder) WithPrenom(prenom string) *EnseignantBuilder {
	b.u.Prenom = prenom
	return b
}

func (b *EnseignantBuilder) WithLogin(login string) *EnseignantBuilder {
	b.u.Login = login
	return b
}

func (b *EnseignantBuilder) WithMotDePasse(mdp string) *EnseignantBuilder {
	b.u.MotDePasse = mdp
	return b
}

func (b *EnseignantBuilder) WithEmail(email string) *EnseignantBuilder {
	b.u.Email = email
	return b
}

func (b *EnseignantBuilder) WithNumeroTelephone(tel string) *EnseignantBuilder {
	b.u.NumeroTelephone = tel
	return b
}

func (b *EnseignantBuilder) WithSoldeCaution(solde float64) *EnseignantBuilder {
	b.u.SoldeCaution = solde
	return b
}

func (b *EnseignantBuilder) WithDateDeNaissance(date time.Time) *EnseignantBuilder {
	b.u.DateDeNaissance = date
	return b
}

func (b *EnseignantBuilder) WithDepartement(dept *models.DepartementEcole) *EnseignantBuilder {
	b.departement = dept
	return b
}

func (b *EnseignantBuilder) Build() *models.Enseignant {
	return &models.Enseignant{
		Utilisateur: b.u,
		Departement: b.departement,
	}
}
