// @title           Bibliothèque API
// @version         1.0
// @description     API REST de gestion de bibliothèque universitaire.
// @host            localhost:8080
// @BasePath        /api

package main

import (
	"bibliotheque/controllers"
	"bibliotheque/db"
	_ "bibliotheque/docs"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api")
	{
		api.GET("/accueil", controllers.Accueil)
		api.POST("/connexion", controllers.Connexion)
		api.GET("/ouvrages", controllers.GetOuvrages)
		api.GET("/ouvrages/:id/exemplaires", controllers.GetExemplairesDisponibles)
		api.GET("/ouvrages/:id/exemplaires/tous", controllers.GetTousExemplaires)
		api.POST("/ouvrages/:id/exemplaires", controllers.CreerExemplaireForOuvrage)
		api.GET("/auteurs", controllers.GetAuteurs)
		api.GET("/emplacements", controllers.GetEmplacements)
		api.POST("/livres", controllers.CreerLivre)
		api.POST("/revues", controllers.CreerRevue)
		api.GET("/departements", controllers.GetDepartements)
		api.POST("/utilisateurs", controllers.CreerUtilisateur)
		api.GET("/utilisateurs/:id/caution", controllers.GetCaution)
		api.PUT("/utilisateurs/:id/caution", controllers.UpdateCautionTotale)
		api.GET("/emprunts/verifier", controllers.VerifierEmprunt)
		api.GET("/emprunts", controllers.ListerEmprunts)
		api.GET("/emprunts/retard", controllers.ListerEmpruntsEnRetard)
		api.POST("/emprunts", controllers.Emprunter)
		api.GET("/utilisateurs/rechercher", controllers.RechercherUtilisateurs)
	}

	r.Run(":8080")
}
