package controllers

import (
	"bibliotheque/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

// L'écran d'accueil affiche des stats sur la bibliothèque
// Ça permet de voir que la BDD a été remplie
type AccueilStats struct {
	NbLivres                 int `json:"nb_livres"`
	NbRevues                 int `json:"nb_revues"`
	NbExemplairesDisponibles int `json:"nb_exemplaires_disponibles"`
	NbUtilisateurs           int `json:"nb_utilisateurs"`
}

// Accueil godoc
// @Summary      Statistiques d'accueil
// @Description  Retourne les statistiques globales de la bibliothèque.
// @Tags         Accueil
// @Produce      json
// @Success      200  {object}  AccueilStats
// @Router       /accueil [get]
func Accueil(c *gin.Context) {
	var stats AccueilStats
	_ = db.GlobalDBO.QueryRow(`SELECT COUNT(*) FROM livres`).Scan(&stats.NbLivres)
	_ = db.GlobalDBO.QueryRow(`SELECT COUNT(*) FROM revues`).Scan(&stats.NbRevues)
	_ = db.GlobalDBO.QueryRow(`SELECT COUNT(*) FROM exemplaires WHERE est_emprunte = FALSE`).Scan(&stats.NbExemplairesDisponibles)
	_ = db.GlobalDBO.QueryRow(`SELECT COUNT(*) FROM utilisateurs`).Scan(&stats.NbUtilisateurs)

	c.JSON(http.StatusOK, stats)
}
