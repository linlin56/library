package main

import (
	"bibliotheque/controllers"
	"bibliotheque/db"
	"html/template"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialise le DBO global depuis db/db.env
	db.Init()
	defer db.GlobalDBO.Close()

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