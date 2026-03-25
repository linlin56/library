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
		api.GET("/ouvrages/:id/exemplaires", controllers.GetExemplairesDisponibles)
		api.GET("/auteurs", controllers.GetAuteurs)
		api.GET("/emplacements", controllers.GetEmplacements)
		api.POST("/livres", controllers.CreerLivre)
		api.POST("/revues", controllers.CreerRevue)
		api.GET("/departements", controllers.GetDepartements)
		api.POST("/utilisateurs", controllers.CreerUtilisateur)
		api.GET("/emprunts/verifier", controllers.VerifierEmprunt)
		api.GET("/emprunts", controllers.ListerEmprunts)
		api.GET("/emprunts/retard", controllers.ListerEmpruntsEnRetard)
		api.POST("/emprunts", controllers.Emprunter)
		api.GET("/utilisateurs/rechercher", controllers.RechercherUtilisateurs)
	}

	r.Run(":8080")
}
