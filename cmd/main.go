package main

import (
	"log"
	"net/http"
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

	router := http.NewServeMux()

	port := os.Getenv("APPLICATION_DEFAULT_PORT")
	server := http.Server{
		Addr:    port,
		Handler: router,
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal("Failed to start http server")
	}
}

func registerAuthRoutes(s *http.ServeMux) {
	apiBasePath := os.Getenv("API_BASE_PATH")

	s.HandleFunc(("POST " + apiBasePath + "/register"), func(http.ResponseWriter, *http.Request) {})
	s.HandleFunc(("POST " + apiBasePath + "/login"), func(http.ResponseWriter, *http.Request) {})
}
