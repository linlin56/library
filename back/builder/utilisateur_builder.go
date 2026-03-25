package builder

import (
	"bibliotheque/models"
	"time"
)

// Pattern builder
type UtilisateurBuilder struct {
	u models.Utilisateur
}

func NewUtilisateurBuilder() *UtilisateurBuilder {
	return &UtilisateurBuilder{}
}

func (b *UtilisateurBuilder) WithNom(nom string) *UtilisateurBuilder {
	b.u.Nom = nom
	return b
}

func (b *UtilisateurBuilder) WithPrenom(prenom string) *UtilisateurBuilder {
	b.u.Prenom = prenom
	return b
}

func (b *UtilisateurBuilder) WithLogin(login string) *UtilisateurBuilder {
	b.u.Login = login
	return b
}

func (b *UtilisateurBuilder) WithMotDePasse(mdp string) *UtilisateurBuilder {
	b.u.MotDePasse = mdp
	return b
}

func (b *UtilisateurBuilder) WithEmail(email string) *UtilisateurBuilder {
	b.u.Email = email
	return b
}

func (b *UtilisateurBuilder) WithNumeroTelephone(tel string) *UtilisateurBuilder {
	b.u.NumeroTelephone = tel
	return b
}

func (b *UtilisateurBuilder) WithSoldeCaution(solde float64) *UtilisateurBuilder {
	b.u.SoldeCaution = solde
	return b
}

func (b *UtilisateurBuilder) WithDateDeNaissance(date time.Time) *UtilisateurBuilder {
	b.u.DateDeNaissance = date
	return b
}

func (b *UtilisateurBuilder) Build() *models.Utilisateur {
	return &b.u
}
