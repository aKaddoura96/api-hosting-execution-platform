package main

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/aKaddoura96/api-hosting-execution-platform/backend/executor/runtime"
	"github.com/aKaddoura96/api-hosting-execution-platform/backend/shared/database"
	"github.com/aKaddoura96/api-hosting-execution-platform/backend/shared/logger"
	"github.com/aKaddoura96/api-hosting-execution-platform/backend/shared/repository"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
	// Initialize logger
	log := logger.NewLogger("executor")
	logger.SetDefaultLogger(log)

	godotenv.Load()
	log.Info("Starting Executor Service", map[string]interface{}{
		"version": "2.0.0",
		"mode":    "production",
	})

	// Connect to database
	log.Info("Connecting to database...")
	if err := database.Connect(); err != nil {
		log.Fatal("Failed to connect to database", map[string]interface{}{"error": err.Error()})
	}
	defer database.Close()
	log.Info("Database connected successfully")

	// Initialize executor
	executor, err := runtime.NewExecutor()
	if err != nil {
		log.Fatal("Failed to initialize executor", map[string]interface{}{"error": err.Error()})
	}
	defer executor.Close()
	log.Info("Executor initialized successfully")

	// Initialize repositories
	apiRepo := repository.NewAPIRepository(database.DB)

	// Setup router
	router := mux.NewRouter()

	// Health check
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "healthy",
			"service": "executor",
			"version": "2.0.0",
		})
	}).Methods("GET")

	// Execute code endpoint
	router.HandleFunc("/execute", func(w http.ResponseWriter, r *http.Request) {
		handleExecute(w, r, executor)
	}).Methods("POST")

	// Deploy API endpoint (for uploaded code)
	router.HandleFunc("/deploy", func(w http.ResponseWriter, r *http.Request) {
		handleDeploy(w, r, apiRepo)
	}).Methods("POST")

	// Stop API endpoint
	router.HandleFunc("/stop/{api_id}", func(w http.ResponseWriter, r *http.Request) {
		handleStop(w, r, apiRepo)
	}).Methods("POST")

	// Get API status
	router.HandleFunc("/status/{api_id}", func(w http.ResponseWriter, r *http.Request) {
		handleStatus(w, r, apiRepo)
	}).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Info("Executor service ready", map[string]interface{}{
		"port":    port,
		"address": "http://localhost:" + port,
	})

	router.Use(logger.HTTPLoggingMiddleware(log))

	// CORS configuration
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:3001"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	if err := http.ListenAndServe(":"+port, c.Handler(router)); err != nil {
		log.Fatal("Server failed", map[string]interface{}{"error": err.Error()})
	}
}

type ExecuteRequest struct {
	Code       string                 `json:"code"`
	Runtime    string                 `json:"runtime"`
	Input      map[string]interface{} `json:"input,omitempty"`
	TimeoutSec int                    `json:"timeout_sec,omitempty"`
}

func handleExecute(w http.ResponseWriter, r *http.Request, executor *runtime.Executor) {
	var req ExecuteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Code == "" {
		http.Error(w, "Code is required", http.StatusBadRequest)
		return
	}

	if req.Runtime == "" {
		http.Error(w, "Runtime is required", http.StatusBadRequest)
		return
	}

	// Execute code
	execReq := &runtime.ExecutionRequest{
		Code:       req.Code,
		Runtime:    req.Runtime,
		Input:      req.Input,
		TimeoutSec: req.TimeoutSec,
	}

	result, err := executor.Execute(execReq)
	if err != nil {
		logger.Error("Execution failed", map[string]interface{}{
			"error":   err.Error(),
			"runtime": req.Runtime,
		})
		http.Error(w, "Execution failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

type DeployRequest struct {
	APIID string `json:"api_id"`
}

type DeployResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func handleDeploy(w http.ResponseWriter, r *http.Request, apiRepo *repository.APIRepository) {
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

	// Check if code is uploaded
	if api.CodePath == "" {
		http.Error(w, "No code uploaded for this API", http.StatusBadRequest)
		return
	}

	// Update status to deployed
	// In a real implementation, we would deploy the code to a persistent container
	// For now, we'll mark it as deployed and it can be executed via /execute endpoint
	if err := apiRepo.UpdateStatus(api.ID, "deployed", ""); err != nil {
		logger.Error("Failed to update API status", map[string]interface{}{"error": err.Error()})
		http.Error(w, "Failed to deploy API", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(DeployResponse{
		Status:  "success",
		Message: "API deployed successfully and ready for execution",
	})
}

func handleStop(w http.ResponseWriter, r *http.Request, apiRepo *repository.APIRepository) {
	vars := mux.Vars(r)
	apiID := vars["api_id"]

	// Update status to stopped
	if err := apiRepo.UpdateStatus(apiID, "stopped", ""); err != nil {
		http.Error(w, "Failed to stop API", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(DeployResponse{
		Status:  "success",
		Message: "API stopped successfully",
	})
}

type StatusResponse struct {
	APIID  string `json:"api_id"`
	Status string `json:"status"`
}

func handleStatus(w http.ResponseWriter, r *http.Request, apiRepo *repository.APIRepository) {
	vars := mux.Vars(r)
	apiID := vars["api_id"]

	// Get API details
	api, err := apiRepo.GetByID(apiID)
	if err != nil {
		http.Error(w, "API not found", http.StatusNotFound)
		return
	}

	response := StatusResponse{
		APIID:  api.ID,
		Status: api.Status,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
