package main

import (
	"log"
	"net/http"
	"os"
	"remissio-auth/internal/auth"
	"remissio-auth/middleware"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	log.Println("[Info] Loading environment variables...")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("[Error] Error loading .env file")
	}

	log.Println("[Info] Successfully loaded environment")

	DB_CONNECTION_STRING := os.Getenv("DB_CONNECTION_STRING")
	if DB_CONNECTION_STRING == "" {
		log.Fatal("[Error] DB_CONNECTION_STRING is not set")
	}

	log.Println("[Info] Establishing database connection...")
	db, err := gorm.Open(postgres.Open(DB_CONNECTION_STRING), &gorm.Config{})
	if err != nil {
		log.Fatal("[Error] Error establishing DB connection: ", err)
	}
	log.Println("[Info] Successfully established database connection")

	log.Println("[Info] Migrating schema 'User'...")
	err = db.AutoMigrate(&auth.User{})
	if err != nil {
		log.Fatal("[Error] Error while migrating schema 'User': ", err)
	}
	log.Println("[Info] Successfully migrated schema 'User'")

	repo := auth.NewRepository(db)
	service := auth.NewService(repo)
	handler := auth.NewHandler(service)

	log.Println("[Info] Initializing new router")
	router := http.NewServeMux()

	log.Println("[Info] Registering routes")
	registerAuthRoutes(router, *handler)

	SERVER_PORT := os.Getenv("APPLICATION_DEFAULT_PORT")
	if SERVER_PORT == "" {
		SERVER_PORT = ":8080"
	}

	log.Println("[Info] Initializing new server")
	server := http.Server{
		Addr:    SERVER_PORT,
		Handler: middleware.Logging(router),
	}

	log.Printf("[Info] Starting server on port: %s", SERVER_PORT)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("[Error] Unable to start http server on port: %s", SERVER_PORT)
	}
}

func registerAuthRoutes(s *http.ServeMux, h auth.Handler) {
	log.Println("[Info] Registering route POST /remissio/api/v1/auth/register")
	s.HandleFunc("POST /remissio/api/v1/auth/register", h.Register)

	log.Println("[Info] Registering route POST /remissio/api/v1/auth/login")
	s.HandleFunc("POST /remissio/api/v1/auth/login", h.Login)

	log.Println("[Info] Registering route POST /remissio/api/v1/auth/logout")
	s.HandleFunc("POST /remissio/api/v1/auth/logout", h.Logout)

	log.Println("[Info] Registering route GET /test")
	s.HandleFunc("GET /test", h.Test)
}
