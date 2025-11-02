package main

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/aKaddoura96/api-hosting-execution-platform/backend/api-gateway/handlers"
	"github.com/aKaddoura96/api-hosting-execution-platform/backend/api-gateway/middleware"
	"github.com/aKaddoura96/api-hosting-execution-platform/backend/shared/database"
	"github.com/aKaddoura96/api-hosting-execution-platform/backend/shared/logger"
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
	// Initialize logger
	log := logger.NewLogger("api-gateway")
	logger.SetDefaultLogger(log)

	// Load environment variables
	godotenv.Load()
	log.Info("Starting API Gateway", map[string]interface{}{
		"version": "1.0.0",
	})

	// Connect to database
	log.Info("Connecting to database...")
	if err := database.Connect(); err != nil {
		log.Fatal("Failed to connect to database", map[string]interface{}{
			"error": err.Error(),
		})
	}
	defer database.Close()
	log.Info("Database connected successfully")

	// Initialize repositories
	log.Info("Initializing repositories")
	userRepo := repository.NewUserRepository(database.DB)
	apiRepo := repository.NewAPIRepository(database.DB)

	// Initialize handlers
	log.Info("Initializing handlers")
	authHandler := handlers.NewAuthHandler(userRepo)
	apiHandler := handlers.NewAPIHandler(apiRepo)
	deployHandler := handlers.NewDeployHandler(apiRepo)

	// Setup router
	log.Info("Setting up routes")
	router := mux.NewRouter()

	// Add logging middleware
	router.Use(logger.HTTPLoggingMiddleware(log))

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
	
	// Deployment routes
	protected.HandleFunc("/apis/{id}/deploy", deployHandler.DeployAPI).Methods("POST")
	protected.HandleFunc("/apis/{id}/stop", deployHandler.StopAPI).Methods("POST")
	protected.HandleFunc("/apis/{id}/status", deployHandler.GetAPIStatus).Methods("GET")

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

	log.Info("API Gateway ready", map[string]interface{}{
		"port": port,
		"address": "http://localhost:" + port,
	})
	
	if err := http.ListenAndServe(":"+port, c.Handler(router)); err != nil {
		log.Fatal("Server failed", map[string]interface{}{
			"error": err.Error(),
		})
	}
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
