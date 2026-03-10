package main

import (
	"bibliotheque/controllers"
	"bibliotheque/db"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	db.Init()
	defer db.GlobalDBO.Close()

	r := gin.Default()

	r.Use(cors.New(cors.Config{
AllowOrigins:     []string{"http://localhost:4200"},
AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
AllowCredentials: true,
}))

	api := r.Group("/api")
	{
		api.GET("/accueil", controllers.Accueil)
		api.POST("/connexion", controllers.Connexion)
		api.GET("/ouvrages", controllers.GetOuvrages)
	}

	r.Run(":8080")
}
