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

// GetExemplairesDisponibles : GET /api/ouvrages/:id/exemplaires
// Retourne les exemplaires disponibles (non empruntés) pour un ouvrage donné.
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
	c.JSON(http.StatusOK, result)
}
