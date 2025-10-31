package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/aKaddoura96/api-hosting-execution-platform/backend/api-gateway/handlers"
	"github.com/aKaddoura96/api-hosting-execution-platform/backend/api-gateway/middleware"
	"github.com/aKaddoura96/api-hosting-execution-platform/backend/shared/database"
	"github.com/aKaddoura96/api-hosting-execution-platform/backend/shared/repository"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Service   string    `json:"service"`
}

func main() {
	// Load environment variables
	godotenv.Load()

	// Connect to database
	if err := database.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Initialize repositories
	userRepo := repository.NewUserRepository(database.DB)
	apiRepo := repository.NewAPIRepository(database.DB)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(userRepo)
	apiHandler := handlers.NewAPIHandler(apiRepo)

	// Setup router
	router := mux.NewRouter()

	// Public routes
	router.HandleFunc("/health", healthHandler).Methods("GET")
	router.HandleFunc("/api/v1/auth/signup", authHandler.Signup).Methods("POST")
	router.HandleFunc("/api/v1/auth/login", authHandler.Login).Methods("POST")
	router.HandleFunc("/api/v1/marketplace/apis", apiHandler.GetPublicAPIs).Methods("GET")
	router.HandleFunc("/api/v1/marketplace/apis/{id}", apiHandler.GetAPI).Methods("GET")

	// Protected routes
	protected := router.PathPrefix("/api/v1").Subrouter()
	protected.Use(middleware.AuthMiddleware)

	// Auth routes
	protected.HandleFunc("/auth/me", authHandler.Me).Methods("GET")

	// API management routes
	protected.HandleFunc("/apis", apiHandler.GetMyAPIs).Methods("GET")
	protected.HandleFunc("/apis", apiHandler.CreateAPI).Methods("POST")
	protected.HandleFunc("/apis/{id}", apiHandler.GetAPI).Methods("GET")
	protected.HandleFunc("/apis/{id}", apiHandler.DeleteAPI).Methods("DELETE")
	protected.HandleFunc("/apis/{id}/upload", apiHandler.UploadCode).Methods("POST")

	// CORS configuration
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:3001"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("API Gateway starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, c.Handler(router)))
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
		Service:   "api-gateway",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
