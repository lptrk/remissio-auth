package main

import (
	"log"
	"os"
	"remissio-auth/internal/auth"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading environment variables")
	}

	dbConnStr := os.Getenv("DB_CONNECTION_STRING")
	db, err := gorm.Open(postgres.Open(dbConnStr), &gorm.Config{})
	if err != nil {
		log.Fatal("Error establishing database connection")
	}

	err = db.AutoMigrate(&auth.User{})
	if err != nil {
		log.Fatal("Error migrating schema 'User'")
	}
}
