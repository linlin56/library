package test

import (
	"bibliotheque/models"
	"testing"
	"time"
)

// Ces tests vérifient que les constructeurs ont le comportement attendu
// ===================== Auteur =====================

func TestNewAuteur_ChampsInitialises(t *testing.T) {
	a := models.NewAuteur(1, "Hugo", "Victor")
	if a.Id != 1 {
		t.Errorf("Id: attendu 1, obtenu %d", a.Id)
	}
	if a.Nom != "Hugo" {
		t.Errorf("Nom: attendu Hugo, obtenu %s", a.Nom)
	}
	if a.Prenom != "Victor" {
		t.Errorf("Prenom: attendu Victor, obtenu %s", a.Prenom)
	}
}

func TestNewAuteur_NomVide(t *testing.T) {
	a := models.NewAuteur(0, "", "")
	if a.Nom != "" {
		t.Errorf("Nom vide attendu, obtenu %s", a.Nom)
	}
}

// ===================== Ouvrage =====================

func TestNewOuvrage_ChampsInitialises(t *testing.T) {
	o := models.NewOuvrage(10, 3.50, "Les Miserables", 5)
	if o.Id != 10 {
		t.Errorf("Id: attendu 10, obtenu %d", o.Id)
	}
	if o.Caution != 3.50 {
		t.Errorf("Caution: attendu 3.50, obtenu %.2f", o.Caution)
	}
	if o.Titre != "Les Miserables" {
		t.Errorf("Titre incorrect: %s", o.Titre)
	}
	if o.Exemplaires != 5 {
		t.Errorf("Exemplaires: attendu 5, obtenu %d", o.Exemplaires)
	}
}

func TestOuvrage_GetId(t *testing.T) {
	o := models.NewOuvrage(42, 0, "", 0)
	if o.GetId() != 42 {
		t.Errorf("GetId: attendu 42, obtenu %d", o.GetId())
	}
}

func TestOuvrage_GetTitre(t *testing.T) {
	o := models.NewOuvrage(0, 0, "Titre Test", 0)
	if o.GetTitre() != "Titre Test" {
		t.Errorf("GetTitre: attendu Titre Test, obtenu %s", o.GetTitre())
	}
}

func TestOuvrage_GetCaution(t *testing.T) {
	o := models.NewOuvrage(0, 7.25, "", 0)
	if o.GetCaution() != 7.25 {
		t.Errorf("GetCaution: attendu 7.25, obtenu %.2f", o.GetCaution())
	}
}

func TestOuvrage_SatisfaitInterfaceIOuvrage(t *testing.T) {
	o := models.NewOuvrage(1, 2.0, "Test", 1)
	var _ models.IOuvrage = o
}

// ===================== Livre =====================

func TestNewLivre_ChampsInitialises(t *testing.T) {
	auteur := *models.NewAuteur(1, "Tolkien", "J.R.R.")
	l := models.NewLivre(3, 4.0, "Le Seigneur des Anneaux", 2, auteur, "978-0618260584")
	if l.Id != 3 {
		t.Errorf("Id: attendu 3, obtenu %d", l.Id)
	}
	if l.Caution != 4.0 {
		t.Errorf("Caution: attendu 4.0, obtenu %.2f", l.Caution)
	}
	if l.Titre != "Le Seigneur des Anneaux" {
		t.Errorf("Titre incorrect: %s", l.Titre)
	}
	if l.Exemplaires != 2 {
		t.Errorf("Exemplaires: attendu 2, obtenu %d", l.Exemplaires)
	}
	if l.Isbn != "978-0618260584" {
		t.Errorf("Isbn incorrect: %s", l.Isbn)
	}
}

func TestNewLivre_AuteurNonNil(t *testing.T) {
	auteur := *models.NewAuteur(5, "Zola", "Emile")
	l := models.NewLivre(1, 2.0, "Germinal", 3, auteur, "978-2070413430")
	if l.Auteur == nil {
		t.Fatal("Auteur ne doit pas etre nil")
	}
	if l.Auteur.Nom != "Zola" {
		t.Errorf("Auteur.Nom: attendu Zola, obtenu %s", l.Auteur.Nom)
	}
}

func TestLivre_SatisfaitInterfaceIOuvrage(t *testing.T) {
	auteur := *models.NewAuteur(1, "Hugo", "Victor")
	l := models.NewLivre(1, 2.0, "Notre-Dame de Paris", 1, auteur, "978-0140443530")
	var _ models.IOuvrage = l
}

func TestLivre_GetTitrePropageDepuisOuvrage(t *testing.T) {
	auteur := *models.NewAuteur(1, "Camus", "Albert")
	l := models.NewLivre(1, 1.0, "L Etranger", 1, auteur, "978-2070360024")
	if l.GetTitre() != "L Etranger" {
		t.Errorf("GetTitre incorrect: %s", l.GetTitre())
	}
}

// ===================== Revue =====================

func TestNewRevue_ChampsInitialises(t *testing.T) {
	date, _ := time.Parse("2006-01-02", "2024-03-01")
	r := models.NewRevue(7, 1.50, "Science et Vie", 3, 42, date)
	if r.Id != 7 {
		t.Errorf("Id: attendu 7, obtenu %d", r.Id)
	}
	if r.Caution != 1.50 {
		t.Errorf("Caution: attendu 1.50, obtenu %.2f", r.Caution)
	}
	if r.Titre != "Science et Vie" {
		t.Errorf("Titre incorrect: %s", r.Titre)
	}
	if r.Numero != 42 {
		t.Errorf("Numero: attendu 42, obtenu %d", r.Numero)
	}
	if !r.DateParution.Equal(date) {
		t.Errorf("DateParution incorrecte: attendu %v, obtenu %v", date, r.DateParution)
	}
}

func TestRevue_SatisfaitInterfaceIOuvrage(t *testing.T) {
	r := models.NewRevue(1, 1.0, "Revue Test", 0, 1, time.Now())
	var _ models.IOuvrage = r
}

func TestRevue_GetCautionCorrecte(t *testing.T) {
	r := models.NewRevue(1, 2.75, "Test", 0, 1, time.Now())
	if r.GetCaution() != 2.75 {
		t.Errorf("GetCaution: attendu 2.75, obtenu %.2f", r.GetCaution())
	}
}

// ===================== Utilisateur =====================

func TestNewUtilisateur_ChampsInitialises(t *testing.T) {
	ddn, _ := time.Parse("2006-01-02", "1990-04-15")
	u := models.NewUtilisateur(1, "Martin", "Jean", "0612345678", 20.0, "jmartin", "mdp", ddn, "jean@exemple.fr")
	if u.Id != 1 {
		t.Errorf("Id: attendu 1, obtenu %d", u.Id)
	}
	if u.Nom != "Martin" {
		t.Errorf("Nom: attendu Martin, obtenu %s", u.Nom)
	}
	if u.SoldeCaution != 20.0 {
		t.Errorf("SoldeCaution: attendu 20.0, obtenu %.2f", u.SoldeCaution)
	}
	if u.Login != "jmartin" {
		t.Errorf("Login: attendu jmartin, obtenu %s", u.Login)
	}
	if u.Email != "jean@exemple.fr" {
		t.Errorf("Email incorrect: %s", u.Email)
	}
}

func TestNewUtilisateur_AdresseIdNilParDefaut(t *testing.T) {
	u := models.NewUtilisateur(0, "", "", "", 0, "", "", time.Time{}, "")
	if u.AdresseId != nil {
		t.Error("AdresseId doit etre nil par defaut")
	}
}

// ===================== Etudiant =====================

func TestNewEtudiant_ChampsInitialises(t *testing.T) {
	e := models.NewEtudiant("Dupont", "Marie", "mdupont", "mdp", "L3")
	if e.Nom != "Dupont" {
		t.Errorf("Nom: attendu Dupont, obtenu %s", e.Nom)
	}
	if e.Prenom != "Marie" {
		t.Errorf("Prenom: attendu Marie, obtenu %s", e.Prenom)
	}
	if e.Login != "mdupont" {
		t.Errorf("Login: attendu mdupont, obtenu %s", e.Login)
	}
	if e.AnneeEtude != "L3" {
		t.Errorf("AnneeEtude: attendu L3, obtenu %s", e.AnneeEtude)
	}
}

func TestEtudiant_HeritageUtilisateur(t *testing.T) {
	e := models.NewEtudiant("Test", "User", "tuser", "mdp", "M1")
	e.SoldeCaution = 15.0
	if e.SoldeCaution != 15.0 {
		t.Errorf("SoldeCaution herite: attendu 15.0, obtenu %.2f", e.SoldeCaution)
	}
}

// ===================== Enseignant =====================

func TestNewEnseignant_ChampsInitialises(t *testing.T) {
	dept := models.NewDepartementEcole(1, "Informatique")
	ens := models.NewEnseignant("Roche", "Pierre", "proche", "mdp", dept)
	if ens.Nom != "Roche" {
		t.Errorf("Nom: attendu Roche, obtenu %s", ens.Nom)
	}
	if ens.Prenom != "Pierre" {
		t.Errorf("Prenom: attendu Pierre, obtenu %s", ens.Prenom)
	}
	if ens.Login != "proche" {
		t.Errorf("Login: attendu proche, obtenu %s", ens.Login)
	}
}

func TestNewEnseignant_DepartementNonNil(t *testing.T) {
	dept := models.NewDepartementEcole(2, "Physique")
	ens := models.NewEnseignant("Bernard", "Claire", "cbernard", "mdp", dept)
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

func TestNewEnseignant_SansDepartement(t *testing.T) {
	ens := models.NewEnseignant("Test", "User", "tuser", "mdp", nil)
	if ens.Departement != nil {
		t.Error("Departement doit etre nil si nil passe en parametre")
	}
}

// ===================== Bibliothecaire =====================

func TestNewBibliothecaire_ChampsInitialises(t *testing.T) {
	b := models.NewBibliothecaire(1, "Petit", "Valerie", "vpetit", "mdp")
	if b.Id != 1 {
		t.Errorf("Id: attendu 1, obtenu %d", b.Id)
	}
	if b.Nom != "Petit" {
		t.Errorf("Nom: attendu Petit, obtenu %s", b.Nom)
	}
	if b.Login != "vpetit" {
		t.Errorf("Login: attendu vpetit, obtenu %s", b.Login)
	}
	if b.MotDePasse != "mdp" {
		t.Errorf("MotDePasse: attendu mdp, obtenu %s", b.MotDePasse)
	}
}

// ===================== DepartementEcole =====================

func TestNewDepartementEcole_ChampsInitialises(t *testing.T) {
	d := models.NewDepartementEcole(3, "Mathematiques")
	if d.Id != 3 {
		t.Errorf("Id: attendu 3, obtenu %d", d.Id)
	}
	if d.Nom != "Mathematiques" {
		t.Errorf("Nom: attendu Mathematiques, obtenu %s", d.Nom)
	}
}

func TestNewDepartementEcole_IdZeroValide(t *testing.T) {
	d := models.NewDepartementEcole(0, "Sans ID")
	if d.Id != 0 {
		t.Errorf("Id: attendu 0, obtenu %d", d.Id)
	}
}

// ===================== Exemplaire =====================

func TestNewExemplaire_ChampsInitialises(t *testing.T) {
	ex := models.NewExemplaire("EX-0001", 15)
	if ex.CodeBarre != "EX-0001" {
		t.Errorf("CodeBarre: attendu EX-0001, obtenu %s", ex.CodeBarre)
	}
	if ex.DelaiEmpruntJours != 15 {
		t.Errorf("DelaiEmpruntJours: attendu 15, obtenu %d", ex.DelaiEmpruntJours)
	}
}

func TestNewExemplaire_NonEmprunterParDefaut(t *testing.T) {
	ex := models.NewExemplaire("EX-0042", 7)
	if ex.EstEmprunte {
		t.Error("EstEmprunte doit etre false par defaut")
	}
}

func TestNewExemplaire_DateVides(t *testing.T) {
	ex := models.NewExemplaire("EX-0010", 15)
	if !ex.DateDebutEmprunt.IsZero() {
		t.Error("DateDebutEmprunt doit etre zero par defaut")
	}
	if !ex.DateFinEmprunt.IsZero() {
		t.Error("DateFinEmprunt doit etre zero par defaut")
	}
}
