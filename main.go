package main

import (
	"bibliotheque/classes"
	"bibliotheque/utils"
	"html/template"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PageData struct {
	Title   string
	Content string
}

func main() {
	// Le gestionnaire de routes Gin
	r := gin.Default()

	// Set custom template function pour les templates chargés par Gin
	r.SetFuncMap(template.FuncMap{
		"safe": func(s string) template.HTML {
			return template.HTML(s)
		},
	})

	// Load HTML templates from the templates directory (template "defaut" doit utiliser {{.Content | safe}})
	r.LoadHTMLGlob("templates/**/*.html")

	r.Static("/static", "./static")

	// Accueil
	r.GET("/", func(c *gin.Context) {
		contentHTML, err := utils.RenderContentTemplate("accueil.html", nil)
		if err != nil {
			log.Printf("Erreur à la lecture du HTML : %v", err)
			c.String(http.StatusInternalServerError, "Impossible de charger la page.")
			return
		}

		data := PageData{
			Title:   "Accueil - Bibliothèque",
			Content: contentHTML,
		}
		c.HTML(http.StatusOK, "defaut", data)
	})

	// Recherche d'ouvrages
	r.GET("/recherche", func(c *gin.Context) {
		contentHTML, err := utils.RenderContentTemplate("recherche.html", nil)
		if err != nil {
			log.Printf("Erreur à la lecture du HTML : %v", err)
			c.String(http.StatusInternalServerError, "Impossible de charger la page.")
			return
		}
		data := PageData{
			Title:   "Recherche - Bibliothèque",
			Content: contentHTML,
		}
		c.HTML(http.StatusOK, "defaut", data)
	})

	// Résultats d'ouvrages (POST)
	r.POST("/ouvrages", func(c *gin.Context) {
		auteurTest := classes.NewAuteur(1, "Saint-Exupéry", "Antoine de")
		livreTest := classes.NewLivre(1, 10.0, "Le Petit Prince", 3, *auteurTest, "978-0156013987")

		ctx := map[string]interface{}{
			"Ouvrages": []*classes.Livre{livreTest},
		}
		contentHTML, err := utils.RenderContentTemplate("ouvrages.html", ctx)
		if err != nil {
			log.Printf("Erreur à la lecture du HTML : %v", err)
			c.String(http.StatusInternalServerError, "Impossible de charger la page.")
			return
		}
		data := PageData{
			Title:   "Résultat de la recherche - Bibliothèque",
			Content: contentHTML,
		}
		c.HTML(http.StatusOK, "defaut", data)
	})

	r.Run()
}
