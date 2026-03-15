package builder

import (
	"bibliotheque/models"
	"time"
)

// EtudiantBuilder étend UtilisateurBuilder pour construire un Etudiant.
type EtudiantBuilder struct {
	UtilisateurBuilder
	anneeEtude string
}

func NewEtudiantBuilder() *EtudiantBuilder {
	return &EtudiantBuilder{}
}

// Les méthodes With* du UtilisateurBuilder retournent *UtilisateurBuilder,
// on redéfinit chacune pour conserver le type *EtudiantBuilder.

func (b *EtudiantBuilder) WithNom(nom string) *EtudiantBuilder {
	b.u.Nom = nom
	return b
}

func (b *EtudiantBuilder) WithPrenom(prenom string) *EtudiantBuilder {
	b.u.Prenom = prenom
	return b
}

func (b *EtudiantBuilder) WithLogin(login string) *EtudiantBuilder {
	b.u.Login = login
	return b
}

func (b *EtudiantBuilder) WithMotDePasse(mdp string) *EtudiantBuilder {
	b.u.MotDePasse = mdp
	return b
}

func (b *EtudiantBuilder) WithEmail(email string) *EtudiantBuilder {
	b.u.Email = email
	return b
}

func (b *EtudiantBuilder) WithNumeroTelephone(tel string) *EtudiantBuilder {
	b.u.NumeroTelephone = tel
	return b
}

func (b *EtudiantBuilder) WithSoldeCaution(solde float64) *EtudiantBuilder {
	b.u.SoldeCaution = solde
	return b
}

func (b *EtudiantBuilder) WithDateDeNaissance(date time.Time) *EtudiantBuilder {
	b.u.DateDeNaissance = date
	return b
}

func (b *EtudiantBuilder) WithAnneeEtude(annee string) *EtudiantBuilder {
	b.anneeEtude = annee
	return b
}

func (b *EtudiantBuilder) Build() *models.Etudiant {
	return &models.Etudiant{
		Utilisateur: b.u,
		AnneeEtude:  b.anneeEtude,
	}
}
