package controllers

import (
	"bibliotheque/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

// L'écran d'accueil affiche des stats sur la bibliothèque
// Ça permet de voir que la BDD a été remplie
type AccueilStats struct {
	NbLivres                int `json:"nb_livres"`
	NbRevues                int `json:"nb_revues"`
	NbExemplairesDisponibles int `json:"nb_exemplaires_disponibles"`
	NbUtilisateurs          int `json:"nb_utilisateurs"`
}

func Accueil(c *gin.Context) {
	var stats AccueilStats
	db.GlobalDBO.QueryRow(`SELECT COUNT(*) FROM livres`).Scan(&stats.NbLivres)
	db.GlobalDBO.QueryRow(`SELECT COUNT(*) FROM revues`).Scan(&stats.NbRevues)
	db.GlobalDBO.QueryRow(`SELECT COUNT(*) FROM exemplaires WHERE est_emprunte = FALSE`).Scan(&stats.NbExemplairesDisponibles)
	db.GlobalDBO.QueryRow(`SELECT COUNT(*) FROM utilisateurs`).Scan(&stats.NbUtilisateurs)

	c.JSON(http.StatusOK, stats)
}
