package controllers

import (
	"bibliotheque/db"
	"bibliotheque/repositories"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetCaution : GET /api/utilisateurs/:id/caution
// Retourne le solde_caution et la caution_totale d'un utilisateur.
// @Summary      Consulter la caution d'un utilisateur
// @Description  Retourne le solde de caution disponible et la caution totale.
// @Tags         Utilisateurs
// @Produce      json
// @Param        id   path      int  true  "ID utilisateur"
// @Success      200  {object}  repositories.CautionInfo
// @Failure      400  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Router       /utilisateurs/{id}/caution [get]
func GetCaution(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"erreur": "Identifiant invalide."})
		return
	}

	repo := repositories.NewUtilisateurRepository(db.GlobalDBO)
	info, err := repo.GetCaution(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"erreur": err.Error()})
		return
	}

	c.JSON(http.StatusOK, info)
}

// UpdateCautionTotale : PUT /api/utilisateurs/:id/caution
// Met à jour la caution_totale d'un utilisateur (réservé à la bibliothécaire).
// Recalcule solde_caution en conservant le montant emprunté.
// Retourne 400 si la nouvelle valeur est inférieure au montant emprunté.
// @Summary      Mettre a jour la caution totale
// @Description  Met a jour la caution totale d'un utilisateur et recalcule son solde.
// @Tags         Utilisateurs
// @Accept       json
// @Produce      json
// @Param        id       path      int                   true  "ID utilisateur"
// @Param        payload  body      UpdateCautionRequest  true  "Nouvelle caution"
// @Success      200      {object}  repositories.CautionInfo
// @Failure      400      {object}  ErrorResponse
// @Router       /utilisateurs/{id}/caution [put]
func UpdateCautionTotale(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"erreur": "Identifiant invalide."})
		return
	}

	var payload struct {
		CautionTotale float64 `json:"caution_totale"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil || payload.CautionTotale < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"erreur": "Valeur de caution invalide."})
		return
	}

	repo := repositories.NewUtilisateurRepository(db.GlobalDBO)
	if err := repo.UpdateCautionTotale(id, payload.CautionTotale); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erreur": err.Error()})
		return
	}

	// Retourner la nouvelle info de caution
	info, err := repo.GetCaution(id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "Caution mise à jour."})
		return
	}
	c.JSON(http.StatusOK, info)
}
