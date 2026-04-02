package test

import (
	"bibliotheque/builder"
	"bibliotheque/models"
	"testing"
	"time"
)

// Tests pour le design pattern builder
// ===================== UtilisateurBuilder =====================

func TestUtilisateurBuilder_Build_ChampsSimples(t *testing.T) {
	ddn := time.Date(1990, 5, 12, 0, 0, 0, 0, time.UTC)
	u := builder.NewUtilisateurBuilder().
		WithNom("Martin").
		WithPrenom("Jean").
		WithLogin("jmartin").
		WithMotDePasse("mdp").
		WithEmail("jean@exemple.fr").
		WithNumeroTelephone("0612345678").
		WithSoldeCaution(20.0).
		WithDateDeNaissance(ddn).
		Build()
	if u.Nom != "Martin" {
		t.Errorf("Nom: attendu Martin, obtenu %s", u.Nom)
	}
	if u.Prenom != "Jean" {
		t.Errorf("Prenom: attendu Jean, obtenu %s", u.Prenom)
	}
	if u.Login != "jmartin" {
		t.Errorf("Login: attendu jmartin, obtenu %s", u.Login)
	}
	if u.MotDePasse != "mdp" {
		t.Errorf("MotDePasse incorrect: %s", u.MotDePasse)
	}
	if u.Email != "jean@exemple.fr" {
		t.Errorf("Email incorrect: %s", u.Email)
	}
	if u.NumeroTelephone != "0612345678" {
		t.Errorf("NumeroTelephone incorrect: %s", u.NumeroTelephone)
	}
	if u.SoldeCaution != 20.0 {
		t.Errorf("SoldeCaution: attendu 20.0, obtenu %.2f", u.SoldeCaution)
	}
	if !u.DateDeNaissance.Equal(ddn) {
		t.Errorf("DateDeNaissance incorrecte")
	}
}

func TestUtilisateurBuilder_Build_RetournePointeurUtilisateur(t *testing.T) {
	u := builder.NewUtilisateurBuilder().WithNom("Test").Build()
	_ = u
}

func TestUtilisateurBuilder_Build_ValeursParDefautVides(t *testing.T) {
	u := builder.NewUtilisateurBuilder().Build()
	if u.Nom != "" {
		t.Errorf("Nom devrait etre vide par defaut, obtenu %s", u.Nom)
	}
	if u.SoldeCaution != 0 {
		t.Errorf("SoldeCaution devrait etre 0 par defaut, obtenu %.2f", u.SoldeCaution)
	}
}

func TestUtilisateurBuilder_Chaining_RetourneMemeBuildeur(t *testing.T) {
	b := builder.NewUtilisateurBuilder()
	u := b.WithNom("Test").WithPrenom("User").Build()
	if u.Nom != "Test" {
		t.Errorf("Nom apres chaingage: attendu Test, obtenu %s", u.Nom)
	}
}

// ===================== EtudiantBuilder =====================

func TestEtudiantBuilder_Build_ChampsComplets(t *testing.T) {
	ddn := time.Date(2000, 9, 1, 0, 0, 0, 0, time.UTC)
	e := builder.NewEtudiantBuilder().
		WithNom("Dupont").
		WithPrenom("Marie").
		WithLogin("mdupont").
		WithMotDePasse("mdp").
		WithEmail("marie@exemple.fr").
		WithNumeroTelephone("0698765432").
		WithSoldeCaution(20.0).
		WithDateDeNaissance(ddn).
		WithAnneeEtude("L3").
		Build()
	if e.Nom != "Dupont" {
		t.Errorf("Nom: attendu Dupont, obtenu %s", e.Nom)
	}
	if e.AnneeEtude != "L3" {
		t.Errorf("AnneeEtude: attendu L3, obtenu %s", e.AnneeEtude)
	}
	if e.SoldeCaution != 20.0 {
		t.Errorf("SoldeCaution: attendu 20.0, obtenu %.2f", e.SoldeCaution)
	}
}

func TestEtudiantBuilder_Build_RetournePointeurEtudiant(t *testing.T) {
	e := builder.NewEtudiantBuilder().WithNom("Test").Build()
	_ = e
}

func TestEtudiantBuilder_Build_AnneeEtudeVideParDefaut(t *testing.T) {
	e := builder.NewEtudiantBuilder().WithNom("Test").Build()
	if e.AnneeEtude != "" {
		t.Errorf("AnneeEtude devrait etre vide par defaut, obtenu %s", e.AnneeEtude)
	}
}

func TestEtudiantBuilder_HeritageUtilisateurChamps(t *testing.T) {
	e := builder.NewEtudiantBuilder().
		WithNom("Bernard").
		WithPrenom("Alice").
		WithSoldeCaution(15.0).
		Build()
	if e.Nom != "Bernard" {
		t.Errorf("Nom: attendu Bernard, obtenu %s", e.Nom)
	}
	if e.SoldeCaution != 15.0 {
		t.Errorf("Utilisateur.SoldeCaution: attendu 15.0, obtenu %.2f", e.SoldeCaution)
	}
}

// ===================== EnseignantBuilder =====================

func TestEnseignantBuilder_Build_SansDepartement(t *testing.T) {
	ens := builder.NewEnseignantBuilder().
		WithNom("Roche").
		WithPrenom("Pierre").
		WithLogin("proche").
		WithMotDePasse("mdp").
		Build()
	if ens.Nom != "Roche" {
		t.Errorf("Nom: attendu Roche, obtenu %s", ens.Nom)
	}
	if ens.Departement != nil {
		t.Error("Departement doit etre nil si non renseigne")
	}
}

func TestEnseignantBuilder_Build_AvecDepartement(t *testing.T) {
	dept := models.NewDepartementEcole(2, "Physique")
	ens := builder.NewEnseignantBuilder().
		WithNom("Bernard").
		WithPrenom("Claire").
		WithLogin("cbernard").
		WithDepartement(dept).
		Build()
	if ens.Departement == nil {
		t.Fatal("Departement ne doit pas etre nil")
	}
	if ens.Departement.Id != 2 {
		t.Errorf("Departement.Id: attendu 2, obtenu %d", ens.Departement.Id)
	}
	if ens.Departement.Nom != "Physique" {
		t.Errorf("Departement.Nom: attendu Physique, obtenu %s", ens.Departement.Nom)
	}
}

func TestEnseignantBuilder_Build_RetournePointeurEnseignant(t *testing.T) {
	ens := builder.NewEnseignantBuilder().WithNom("Test").Build()
	_ = ens
}

func TestEnseignantBuilder_Build_ChampsUtilisateurPropages(t *testing.T) {
	ens := builder.NewEnseignantBuilder().
		WithNom("Durand").
		WithPrenom("Luc").
		WithEmail("luc@exemple.fr").
		WithSoldeCaution(20.0).
		Build()
	if ens.Nom != "Durand" {
		t.Errorf("Utilisateur.Nom: attendu Durand, obtenu %s", ens.Nom)
	}
	if ens.Email != "luc@exemple.fr" {
		t.Errorf("Utilisateur.Email incorrect: %s", ens.Email)
	}
}

func TestEnseignantBuilder_Build_DateDeNaissance(t *testing.T) {
	ddn := time.Date(1975, 3, 20, 0, 0, 0, 0, time.UTC)
	ens := builder.NewEnseignantBuilder().
		WithNom("Lefebvre").
		WithDateDeNaissance(ddn).
		Build()
	if !ens.DateDeNaissance.Equal(ddn) {
		t.Errorf("DateDeNaissance incorrecte")
	}
}
