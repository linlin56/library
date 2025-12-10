package db

import (
	"fmt"
	"log"
	"os"

	"bibliotheque/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {
	// Récupère les variables d'environnement depuis le fichier db.env
	godotenv.Load("db/db.env")

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC", host, user, password, dbname, port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	DB = db

	if err := DB.AutoMigrate(getModels()...); err != nil {
		log.Fatalf("auto migrate failed: %v", err)
	}
}

// Retourne la liste des modèles à migrer pour créer toutes les tables SQL
func getModels() []interface{} {
	return []interface{}{
		&models.Auteur{},
		&models.Bibliothecaire{},
		&models.DepartementEcole{},
		&models.Enseignant{},
		&models.Etudiant{},
		&models.Exemplaire{},
		&models.Livre{},
		&models.Ouvrage{},
		&models.Revue{},
		&models.Utilisateur{},
	}
}
