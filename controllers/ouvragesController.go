package controllers

import (
	"bibliotheque/models"
	"bibliotheque/utils"
	"html/template"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RechercheOuvrages(c *gin.Context) {

	// Exemple statique — remplacé plus tard par un service + GORM
	auteurTest := models.NewAuteur(1, "Saint-Exupéry", "Antoine de")
	livreTest := models.NewLivre(1, 10.0, "Le Petit Prince", 3, *auteurTest, "978-0156013987")

	ctx := map[string]interface{}{
		"Ouvrages": []*models.Livre{livreTest},
	}

	contentHTML, err := utils.RenderContentTemplate("ouvrages.html", ctx)
	if err != nil {
		log.Printf("Erreur à la lecture du HTML : %v", err)
		c.String(http.StatusInternalServerError, "Impossible de charger la page.")
		return
	}

	data := PageData{
		Title:   "Résultat de la recherche - Bibliothèque",
		Content: template.HTML(contentHTML),
	}

	c.HTML(http.StatusOK, "default", data)
}
