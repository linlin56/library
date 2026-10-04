package db

import (
	"log"
	"strconv"

	"github.com/joho/godotenv"
)

/**
 * Notes: to manage the database connection, I chose to use a DBO.
 * There is an ORM for Go called GORM, but it does not support PostgreSQL inheritance.
 * I started by designing the database and I find inheritance very useful,
 * so I preferred to use a DBO and repositories.
 */

// Global DBO, available from the whole project
// It is a singleton
var GlobalDBO *DBO

// Init loads the db.env file and initializes the global DBO
func Init() {
	// godotenv loads a .env file
	if err := godotenv.Load("db/db.env"); err != nil {
		log.Println("db/db.env file not found, using system environment variables")
	}

	port, err := strconv.Atoi(getEnv("DB_PORT", "5432"))
	if err != nil {
		log.Fatalf("Invalid DB_PORT: %v", err)
	}

	cfg := Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     port,
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", ""),
		DBName:   getEnv("DB_NAME", "library"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	dbo, err := NewDBO(cfg)
	if err != nil {
		log.Fatalf("Unable to connect to the database: %v", err)
	}

	GlobalDBO = dbo

	if err := GlobalDBO.SeedIfEmpty("db/scripts/insert.sql"); err != nil {
		log.Fatalf("Error while seeding the database: %v", err)
	}
}