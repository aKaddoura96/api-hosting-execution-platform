package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/aKaddoura96/api-hosting-execution-platform/backend/shared/database"
	"github.com/aKaddoura96/api-hosting-execution-platform/backend/shared/logger"
	"github.com/aKaddoura96/api-hosting-execution-platform/backend/shared/repository"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type DeployRequest struct {
	APIID string `json:"api_id"`
}

type DeployResponse struct {
	Status      string `json:"status"`
	ContainerID string `json:"container_id,omitempty"`
	Message     string `json:"message"`
}

type StatusResponse struct {
	APIID       string `json:"api_id"`
	Status      string `json:"status"`
	ContainerID string `json:"container_id,omitempty"`
}

func main() {
	// Initialize logger
	log := logger.NewLogger("executor")
	logger.SetDefaultLogger(log)
	
	godotenv.Load()
	log.Info("Starting Executor Service (Stub Mode)", map[string]interface{}{
		"version": "1.0.0",
		"mode": "stub",
	})

	// Connect to database
	log.Info("Connecting to database...")
	if err := database.Connect(); err != nil {
		log.Fatal("Failed to connect to database", map[string]interface{}{"error": err.Error()})
	}
	defer database.Close()
	log.Info("Database connected successfully")

	// Initialize repositories
	apiRepo := repository.NewAPIRepository(database.DB)

	// Setup router
	router := mux.NewRouter()

	// Health check
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "healthy",
			"service": "executor",
			"note":    "Docker integration pending - using stub mode",
		})
	}).Methods("GET")

	// Deploy API (stub)
	router.HandleFunc("/deploy", func(w http.ResponseWriter, r *http.Request) {
		var req DeployRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		// Get API details
		api, err := apiRepo.GetByID(req.APIID)
		if err != nil {
			http.Error(w, "API not found", http.StatusNotFound)
			return
		}

		// Simulate deployment
		stubContainerID := "stub-container-" + api.ID
		apiRepo.UpdateStatus(api.ID, "deployed", stubContainerID)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(DeployResponse{
			Status:      "success",
			ContainerID: stubContainerID,
			Message:     "API deployed successfully (stub mode)",
		})
	}).Methods("POST")

	// Stop API (stub)
	router.HandleFunc("/stop/{api_id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		apiID := vars["api_id"]

		apiRepo.UpdateStatus(apiID, "stopped", "")

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(DeployResponse{
			Status:  "success",
			Message: "API stopped successfully (stub mode)",
		})
	}).Methods("POST")

	// Get API status (stub)
	router.HandleFunc("/status/{api_id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		apiID := vars["api_id"]

		api, err := apiRepo.GetByID(apiID)
		if err != nil {
			http.Error(w, "API not found", http.StatusNotFound)
			return
		}

		response := StatusResponse{
			APIID:       api.ID,
			Status:      api.Status,
			ContainerID: api.ContainerID,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Info("Executor service ready", map[string]interface{}{
		"port": port,
		"address": "http://localhost:" + port,
		"mode": "stub",
	})
	
	router.Use(logger.HTTPLoggingMiddleware(log))
	
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatal("Server failed", map[string]interface{}{"error": err.Error()})
	}
}
