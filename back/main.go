// @title           Library API
// @version         1.0
// @description     REST API for managing a university library.
// @host            localhost:8080
// @BasePath        /api

package main

import (
	"library/controllers"
	"library/db"
	_ "library/docs"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	db.Init()
	defer func() { _ = db.GlobalDBO.Close() }()

	r := gin.Default()

	corsOrigin := os.Getenv("CORS_ORIGIN")
	if corsOrigin == "" {
		corsOrigin = "http://localhost:4200"
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{corsOrigin},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept-Language"},
		AllowCredentials: true,
	}))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/api")
	{
		api.GET("/home", controllers.Home)
		api.POST("/login", controllers.Login)
		api.GET("/items", controllers.GetItems)
		api.GET("/items/:id/copies", controllers.GetAvailableCopies)
		api.GET("/items/:id/copies/all", controllers.GetAllCopies)
		api.POST("/items/:id/copies", controllers.CreateCopyForItem)
		api.GET("/authors", controllers.GetAuthors)
		api.GET("/locations", controllers.GetLocations)
		api.POST("/books", controllers.CreateBook)
		api.POST("/journals", controllers.CreateJournal)
		api.GET("/departments", controllers.GetDepartments)
		api.POST("/users", controllers.CreateUser)
		api.GET("/users/:id/deposit", controllers.GetDeposit)
		api.PUT("/users/:id/deposit", controllers.UpdateTotalDeposit)
		api.GET("/loans/check", controllers.CheckLoan)
		api.GET("/loans", controllers.ListLoans)
		api.GET("/loans/overdue", controllers.ListOverdueLoans)
		api.POST("/loans", controllers.Borrow)
		api.DELETE("/loans/:id", controllers.ReturnLoan)
		api.GET("/users/search", controllers.SearchUsers)
	}

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
