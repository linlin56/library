package main

import (
	"bibliotheque/controllers"
	"html/template"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()

	r.SetFuncMap(template.FuncMap{
		"safe": func(s string) template.HTML {
			return template.HTML(s)
		},
	})

	r.LoadHTMLGlob("views/**/*.html")
	r.Static("/static", "./static")

	// Routes
	r.GET("/", controllers.Accueil)
	r.GET("/recherche", controllers.Recherche)
	r.POST("/ouvrages", controllers.RechercheOuvrages)

	r.Run()
}
