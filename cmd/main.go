package main

import (
	"log"
	"os"

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
	db, err = gorm.Open(postgres.Open(dbConnStr), &gorm.Config{})
}
