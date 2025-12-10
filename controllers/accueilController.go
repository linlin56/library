package controllers

import (
	"bibliotheque/utils"
	"html/template"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PageData struct {
	Title   string
	Content template.HTML
}

func Accueil(c *gin.Context) {
	contentHTML, err := utils.RenderContentTemplate("accueil.html", nil)
	if err != nil {
		log.Printf("Erreur à la lecture du HTML : %v", err)
		c.String(http.StatusInternalServerError, "Impossible de charger la page.")
		return
	}

	data := PageData{
		Title:   "Accueil - Bibliothèque",
		Content: template.HTML(contentHTML),
	}

	c.HTML(http.StatusOK, "default", data)
}
