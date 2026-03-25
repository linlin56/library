package controllers

import (
	"bibliotheque/db"
	"bibliotheque/models"
	"bibliotheque/repositories"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// GetAuteurs : GET /api/auteurs
func GetAuteurs(c *gin.Context) {
	repo := repositories.NewAuteurRepository(db.GlobalDBO)
	auteurs, err := repo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": "Impossible de charger les auteurs."})
		return
	}
	c.JSON(http.StatusOK, auteurs)
}

// GetEmplacements : GET /api/emplacements
func GetEmplacements(c *gin.Context) {
	repo := repositories.NewEmplacementRepository(db.GlobalDBO)
	emplacements, err := repo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": "Impossible de charger les emplacements."})
		return
	}
	c.JSON(http.StatusOK, emplacements)
}

type CreerLivreRequest struct {
	Titre        string  `json:"titre"`
	Caution      float64 `json:"caution"`
	Isbn         string  `json:"isbn"`
	AuteurId     int     `json:"auteur_id"`
	AuteurNom    string  `json:"auteur_nom"`
	AuteurPrenom string  `json:"auteur_prenom"`
	EmplacementId int   `json:"emplacement_id"`
}

// CreerLivre : POST /api/livres
func CreerLivre(c *gin.Context) {
	var req CreerLivreRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erreur": "Corps de requete invalide."})
		return
	}

	req.Titre = strings.TrimSpace(req.Titre)
	req.Isbn = strings.TrimSpace(req.Isbn)
	req.AuteurNom = strings.TrimSpace(req.AuteurNom)
	req.AuteurPrenom = strings.TrimSpace(req.AuteurPrenom)

	if req.Titre == "" || req.Isbn == "" {
		c.JSON(http.StatusBadRequest, gin.H{"erreur": "Le titre et l'ISBN sont requis."})
		return
	}
	if req.EmplacementId <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"erreur": "L'emplacement est requis."})
		return
	}
	if req.Caution < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"erreur": "La caution ne peut pas etre negative."})
		return
	}

	auteurId := req.AuteurId
	if auteurId == 0 {
		if req.AuteurNom == "" || req.AuteurPrenom == "" {
			c.JSON(http.StatusBadRequest, gin.H{"erreur": "Auteur requis : selectionnez un auteur existant ou saisissez nom et prenom."})
			return
		}
		auteurRepo := repositories.NewAuteurRepository(db.GlobalDBO)
		id, err := auteurRepo.Create(&models.Auteur{Nom: req.AuteurNom, Prenom: req.AuteurPrenom})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"erreur": "Impossible de creer l'auteur."})
			return
		}
		auteurId = id
	}

	ouvrageRepo := repositories.NewOuvrageRepository(db.GlobalDBO)
	id, err := ouvrageRepo.CreateLivre(req.Titre, req.Caution, req.Isbn, auteurId, req.EmplacementId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": "Impossible de creer le livre."})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id, "message": "Livre créé avec succès."})
}

type CreerRevueRequest struct {
	Titre         string  `json:"titre"`
	Caution       float64 `json:"caution"`
	Numero        int     `json:"numero"`
	DateParution  string  `json:"date_parution"`
	EmplacementId int     `json:"emplacement_id"`
}

// CreerRevue : POST /api/revues
func CreerRevue(c *gin.Context) {
	var req CreerRevueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erreur": "Corps de requete invalide."})
		return
	}

	req.Titre = strings.TrimSpace(req.Titre)
	req.DateParution = strings.TrimSpace(req.DateParution)

	if req.Titre == "" {
		c.JSON(http.StatusBadRequest, gin.H{"erreur": "Le titre est requis."})
		return
	}
	if req.Numero <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"erreur": "Le numéro doit être un entier positif."})
		return
	}
	if req.Caution < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"erreur": "La caution ne peut pas être négative."})
		return
	}
	if req.DateParution == "" {
		c.JSON(http.StatusBadRequest, gin.H{"erreur": "La date de parution est requise."})
		return
	}
	if req.EmplacementId <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"erreur": "L'emplacement est requis."})
		return
	}

	ouvrageRepo := repositories.NewOuvrageRepository(db.GlobalDBO)
	id, err := ouvrageRepo.CreateRevue(req.Titre, req.Caution, req.Numero, req.DateParution, req.EmplacementId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": "Impossible de créer la revue."})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": id, "message": "Revue créée avec succès."})
}
