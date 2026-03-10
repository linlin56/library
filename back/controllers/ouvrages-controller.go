package controllers

import (
	"bibliotheque/db"
	"bibliotheque/repositories"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetOuvrages(c *gin.Context) {
	titre := c.Query("titre")
	repo := repositories.NewOuvrageRepository(db.GlobalDBO)

	if titre != "" {
		ouvrages, err := repo.FindByTitre(titre)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"erreur": err.Error()})
			return
		}
		c.JSON(http.StatusOK, ouvrages)
		return
	}

	ouvrages, err := repo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"erreur": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ouvrages)
}
