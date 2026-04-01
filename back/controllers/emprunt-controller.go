package controllers

import (
	"bibliotheque/db"
	"bibliotheque/repositories"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type EmpruntRequest struct {
	UtilisateurId int    `json:"utilisateur_id" binding:"required"`
	CodeBarre     string `json:"code_barre" binding:"required"`
}

// VerifierEmprunt : GET /api/emprunts/verifier?utilisateur_id=X&code_barre=Y
// @Summary      Verifier un emprunt
// @Description  Verifie la possibilite d'emprunt et retourne un apercu.
// @Tags         Emprunts
// @Produce      json
// @Param        utilisateur_id  query     int     true  "ID utilisateur"
// @Param        code_barre      query     string  true  "Code barre exemplaire"
// @Success      200             {object}  repositories.PreviewEmprunt
// @Failure      400             {object}  ErrorResponse
// @Failure      422             {object}  ErrorResponse
// @Router       /emprunts/verifier [get]
func VerifierEmprunt(c *gin.Context) {
	utilisateurIdStr := c.Query("utilisateur_id")
	codeBarre := c.Query("code_barre")

	utilisateurId, err := strconv.Atoi(utilisateurIdStr)
	if err != nil || utilisateurId <= 0 || codeBarre == "" {
		c.JSON(http.StatusBadRequest, gin.H{"erreur": "Parametres manquants."})
		return
	}

	repo := repositories.NewEmpruntRepository(db.GlobalDBO)
	preview, err := repo.Verifier(utilisateurId, codeBarre)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"erreur": err.Error()})
		return
	}
	c.JSON(http.StatusOK, preview)
}

// Emprunter : POST /api/emprunts
// @Summary      Enregistrer un emprunt
// @Description  Enregistre l'emprunt d'un exemplaire pour un utilisateur.
// @Tags         Emprunts
// @Accept       json
// @Produce      json
// @Param        payload  body      EmpruntRequest  true  "Demande d'emprunt"
// @Success      201      {object}  MessageResponse
// @Failure      400      {object}  ErrorResponse
// @Failure      422      {object}  ErrorResponse
// @Router       /emprunts [post]
func Emprunter(c *gin.Context) {
	var req EmpruntRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erreur": "Parametres manquants."})
		return
	}

	repo := repositories.NewEmpruntRepository(db.GlobalDBO)
	if err := repo.Emprunter(req.UtilisateurId, req.CodeBarre); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"erreur": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Emprunt enregistre avec succes."})
}
