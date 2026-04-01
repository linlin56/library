package controllers

import (
	"bibliotheque/db"
	"bibliotheque/repositories"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ExemplaireResume struct {
	Id        int    `json:"id"`
	CodeBarre string `json:"code_barre"`
}

type ExemplaireDetail struct {
	Id          int    `json:"id"`
	CodeBarre   string `json:"code_barre"`
	EstEmprunte bool   `json:"est_emprunte"`
	DateFin     string `json:"date_fin_emprunt,omitempty"`
}

// GetExemplairesDisponibles : GET /api/ouvrages/:id/exemplaires
// Retourne les exemplaires disponibles (non empruntés) pour un ouvrage donné.
// @Summary      Lister les exemplaires disponibles
// @Description  Retourne les exemplaires non empruntes pour un ouvrage.
// @Tags         Exemplaires
// @Produce      json
// @Param        id   path      int  true  "ID ouvrage"
// @Success      200  {array}   ExemplaireResume
// @Success      204  {string}  string  "No Content"
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /ouvrages/{id}/exemplaires [get]
func GetExemplairesDisponibles(c *gin.Context) {
	idStr := c.Param("id")
	ouvrageId, err := strconv.Atoi(idStr)
	if err != nil || ouvrageId <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"erreur": "Identifiant invalide."})
		return
	}

	repo := repositories.NewExemplaireRepository(db.GlobalDBO)
	exemplaires, err := repo.FindDisponiblesByOuvrageId(ouvrageId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": "Impossible de charger les exemplaires."})
		return
	}

	result := make([]ExemplaireResume, 0, len(exemplaires))
	for _, e := range exemplaires {
		result = append(result, ExemplaireResume{Id: e.Id, CodeBarre: e.CodeBarre})
	}
	if len(result) == 0 {
		c.Status(http.StatusNoContent)
		return
	}
	c.JSON(http.StatusOK, result)
}

// GetTousExemplaires : GET /api/ouvrages/:id/exemplaires/tous
// Retourne tous les exemplaires (disponibles et empruntés) pour un ouvrage donné.
// @Summary      Lister tous les exemplaires
// @Description  Retourne tous les exemplaires d'un ouvrage.
// @Tags         Exemplaires
// @Produce      json
// @Param        id   path      int  true  "ID ouvrage"
// @Success      200  {array}   ExemplaireDetail
// @Success      204  {string}  string  "No Content"
// @Failure      400  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /ouvrages/{id}/exemplaires/tous [get]
func GetTousExemplaires(c *gin.Context) {
	idStr := c.Param("id")
	ouvrageId, err := strconv.Atoi(idStr)
	if err != nil || ouvrageId <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"erreur": "Identifiant invalide."})
		return
	}

	repo := repositories.NewExemplaireRepository(db.GlobalDBO)
	exemplaires, err := repo.FindAllByOuvrageId(ouvrageId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": "Impossible de charger les exemplaires."})
		return
	}

	result := make([]ExemplaireDetail, 0, len(exemplaires))
	for _, e := range exemplaires {
		d := ExemplaireDetail{
			Id:          e.Id,
			CodeBarre:   e.CodeBarre,
			EstEmprunte: e.EstEmprunte,
		}
		if e.EstEmprunte && !e.DateFinEmprunt.IsZero() {
			d.DateFin = e.DateFinEmprunt.Format("02/01/2006")
		}
		result = append(result, d)
	}
	if len(result) == 0 {
		c.Status(http.StatusNoContent)
		return
	}
	c.JSON(http.StatusOK, result)
}

// CreerExemplaireForOuvrage : POST /api/ouvrages/:id/exemplaires
// Crée un nouvel exemplaire pour un ouvrage donné.
// @Summary      Creer un exemplaire
// @Description  Cree un exemplaire pour un ouvrage donne.
// @Tags         Exemplaires
// @Accept       json
// @Produce      json
// @Param        id       path      int                    true  "ID ouvrage"
// @Param        payload  body      CreerExemplaireRequest true  "Donnees exemplaire"
// @Success      201      {object}  IDResponse
// @Failure      400      {object}  ErrorResponse
// @Router       /ouvrages/{id}/exemplaires [post]
func CreerExemplaireForOuvrage(c *gin.Context) {
	idStr := c.Param("id")
	ouvrageId, err := strconv.Atoi(idStr)
	if err != nil || ouvrageId <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"erreur": "Identifiant invalide."})
		return
	}

	var payload struct {
		CodeBarre string `json:"code_barre"`
		Delai     int    `json:"delai_emprunt_jours"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil || payload.CodeBarre == "" {
		c.JSON(http.StatusBadRequest, gin.H{"erreur": "Code-barre requis."})
		return
	}
	if payload.Delai <= 0 {
		payload.Delai = 15
	}

	repo := repositories.NewExemplaireRepository(db.GlobalDBO)
	newID, err := repo.CreateForOuvrage(ouvrageId, payload.CodeBarre, payload.Delai)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erreur": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": newID})
}
