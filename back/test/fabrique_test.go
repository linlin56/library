package test

import (
	"bibliotheque/fabrique"
	"bibliotheque/models"
	"testing"
	"time"
)

// ===================== FabriqueLivre =====================

func TestFabriqueLivre_CreerOuvrage_RetourneLivre(t *testing.T) {
	auteur := models.NewAuteur(1, "Hugo", "Victor")
	f := &fabrique.FabriqueLivre{
		Titre:   "Notre-Dame de Paris",
		Caution: 3.0,
		Isbn:    "978-0140443530",
		Auteur:  auteur,
	}
	ouvrage := f.CreerOuvrage()
	livre, ok := ouvrage.(*models.Livre)
	if !ok {
		t.Fatal("CreerOuvrage doit retourner un *models.Livre")
	}
	if livre.Titre != "Notre-Dame de Paris" {
		t.Errorf("Titre incorrect: %s", livre.Titre)
	}
	if livre.Isbn != "978-0140443530" {
		t.Errorf("Isbn incorrect: %s", livre.Isbn)
	}
	if livre.Caution != 3.0 {
		t.Errorf("Caution: attendu 3.0, obtenu %.2f", livre.Caution)
	}
}

func TestFabriqueLivre_CreerOuvrage_AuteurPropague(t *testing.T) {
	auteur := models.NewAuteur(5, "Zola", "Emile")
	f := &fabrique.FabriqueLivre{
		Titre:   "Germinal",
		Caution: 2.0,
		Isbn:    "978-2070413430",
		Auteur:  auteur,
	}
	livre := f.CreerOuvrage().(*models.Livre)
	if livre.Auteur == nil {
		t.Fatal("Auteur ne doit pas etre nil dans le livre cree")
	}
	if livre.Auteur.Nom != "Zola" {
		t.Errorf("Auteur.Nom: attendu Zola, obtenu %s", livre.Auteur.Nom)
	}
}

func TestFabriqueLivre_CreerOuvrage_IdInitialise0(t *testing.T) {
	auteur := models.NewAuteur(1, "Test", "Auteur")
	f := &fabrique.FabriqueLivre{Titre: "Test", Caution: 1.0, Isbn: "000", Auteur: auteur}
	livre := f.CreerOuvrage().(*models.Livre)
	if livre.Id != 0 {
		t.Errorf("Id doit etre 0 (genere par la DB), obtenu %d", livre.Id)
	}
}

func TestFabriqueLivre_SatisfaitFabriqueOuvrage(t *testing.T) {
	auteur := models.NewAuteur(1, "Test", "Auteur")
	var _ fabrique.FabriqueOuvrage = &fabrique.FabriqueLivre{Titre: "Test", Auteur: auteur}
}

func TestFabriqueLivre_CreerOuvrage_SatisfaitIOuvrage(t *testing.T) {
	auteur := models.NewAuteur(1, "Test", "Auteur")
	f := &fabrique.FabriqueLivre{Titre: "Test", Caution: 1.0, Isbn: "000", Auteur: auteur}
	var _ models.IOuvrage = f.CreerOuvrage()
}

// ===================== FabriqueRevue =====================

func TestFabriqueRevue_CreerOuvrage_RetourneRevue(t *testing.T) {
	date, _ := time.Parse("2006-01-02", "2025-01-01")
	f := &fabrique.FabriqueRevue{
		Titre:        "Science et Vie",
		Caution:      1.50,
		Numero:       99,
		DateParution: date,
	}
	ouvrage := f.CreerOuvrage()
	revue, ok := ouvrage.(*models.Revue)
	if !ok {
		t.Fatal("CreerOuvrage doit retourner un *models.Revue")
	}
	if revue.Titre != "Science et Vie" {
		t.Errorf("Titre incorrect: %s", revue.Titre)
	}
	if revue.Numero != 99 {
		t.Errorf("Numero: attendu 99, obtenu %d", revue.Numero)
	}
	if revue.Caution != 1.50 {
		t.Errorf("Caution: attendu 1.50, obtenu %.2f", revue.Caution)
	}
}

func TestFabriqueRevue_CreerOuvrage_DateParutionPropagee(t *testing.T) {
	date, _ := time.Parse("2006-01-02", "2024-06-15")
	f := &fabrique.FabriqueRevue{Titre: "Test", Caution: 1.0, Numero: 1, DateParution: date}
	revue := f.CreerOuvrage().(*models.Revue)
	if !revue.DateParution.Equal(date) {
		t.Errorf("DateParution incorrecte")
	}
}

func TestFabriqueRevue_CreerOuvrage_IdInitialise0(t *testing.T) {
	f := &fabrique.FabriqueRevue{Titre: "Test", Caution: 1.0, Numero: 1, DateParution: time.Now()}
	revue := f.CreerOuvrage().(*models.Revue)
	if revue.Id != 0 {
		t.Errorf("Id doit etre 0 (genere par la DB), obtenu %d", revue.Id)
	}
}

func TestFabriqueRevue_SatisfaitFabriqueOuvrage(t *testing.T) {
	var _ fabrique.FabriqueOuvrage = &fabrique.FabriqueRevue{Titre: "Test", DateParution: time.Now()}
}

func TestFabriqueRevue_CreerOuvrage_SatisfaitIOuvrage(t *testing.T) {
	f := &fabrique.FabriqueRevue{Titre: "Test", Caution: 1.0, Numero: 1, DateParution: time.Now()}
	var _ models.IOuvrage = f.CreerOuvrage()
}

// ===================== Polymorphisme abstrait =====================

func TestFabrique_PolymorphismeAbstrait(t *testing.T) {
	auteur := models.NewAuteur(1, "Hugo", "Victor")
	fabriques := []fabrique.FabriqueOuvrage{
		&fabrique.FabriqueLivre{
			Titre:   "Notre-Dame de Paris",
			Caution: 2.0,
			Isbn:    "978-0140443530",
			Auteur:  auteur,
		},
		&fabrique.FabriqueRevue{
			Titre:        "Science et Vie",
			Caution:      1.0,
			Numero:       5,
			DateParution: time.Now(),
		},
	}
	for i, f := range fabriques {
		o := f.CreerOuvrage()
		if o.GetTitre() == "" {
			t.Errorf("fabrique[%d]: GetTitre ne doit pas etre vide", i)
		}
		if o.GetCaution() <= 0 {
			t.Errorf("fabrique[%d]: GetCaution doit etre positif, obtenu %.2f", i, o.GetCaution())
		}
	}
}
