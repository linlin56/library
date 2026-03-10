package db

import (
	"log"
	"strconv"

	"github.com/joho/godotenv"
)

// DBO global accessible depuis tout le projet
var GlobalDBO *DBO

// Init charge le fichier db.env et initialise le DBO global
func Init() {
	if err := godotenv.Load("db/db.env"); err != nil {
		log.Println("⚠️  Fichier db/db.env non trouvé, utilisation des variables d'environnement système")
	}

	port, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
	if err != nil {
		log.Fatalf("DB_PORT invalide: %v", err)
	}

	cfg := Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     port,
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", ""),
		DBName:   getEnv("DB_NAME", "bibliotheque"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	dbo, err := NewDBO(cfg)
	if err != nil {
		log.Fatalf("❌ Impossible de se connecter à la base de données: %v", err)
	}

	GlobalDBO = dbo

	if err := GlobalDBO.SeedIfEmpty("db/scripts/insert.sql"); err != nil {
		log.Fatalf("❌ Erreur peuplement base de données: %v", err)
	}
}