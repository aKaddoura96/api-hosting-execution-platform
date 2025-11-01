package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/aKaddoura96/api-hosting-execution-platform/backend/shared/database"
	"github.com/aKaddoura96/api-hosting-execution-platform/backend/shared/models"
	"github.com/aKaddoura96/api-hosting-execution-platform/backend/shared/repository"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	// Connect to database
	if err := database.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Initialize repositories
	execRepo := repository.NewExecutionRepository(database.DB)
	apiRepo := repository.NewAPIRepository(database.DB)

	// Setup router
	router := mux.NewRouter()

	// Health check
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "healthy",
			"service": "analytics",
		})
	}).Methods("GET")

	// Log execution
	router.HandleFunc("/log", func(w http.ResponseWriter, r *http.Request) {
		handleLogExecution(w, r, execRepo)
	}).Methods("POST")

	// Get API stats
	router.HandleFunc("/stats/{api_id}", func(w http.ResponseWriter, r *http.Request) {
		handleGetStats(w, r, execRepo, apiRepo)
	}).Methods("GET")

	// Get execution history
	router.HandleFunc("/history/{api_id}", func(w http.ResponseWriter, r *http.Request) {
		handleGetHistory(w, r, execRepo, apiRepo)
	}).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	log.Printf("Analytics service starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

type LogExecutionRequest struct {
	APIID        string `json:"api_id"`
	UserID       string `json:"user_id,omitempty"`
	StatusCode   int    `json:"status_code"`
	Duration     int64  `json:"duration_ms"`
	RequestSize  int64  `json:"request_size"`
	ResponseSize int64  `json:"response_size"`
	Error        string `json:"error,omitempty"`
}

func handleLogExecution(w http.ResponseWriter, r *http.Request, execRepo *repository.ExecutionRepository) {
	var req LogExecutionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	execution := &models.Execution{
		APIID:        req.APIID,
		UserID:       req.UserID,
		StatusCode:   req.StatusCode,
		Duration:     time.Duration(req.Duration) * time.Millisecond,
		RequestSize:  req.RequestSize,
		ResponseSize: req.ResponseSize,
		Error:        req.Error,
	}

	if err := execRepo.Create(execution); err != nil {
		log.Printf("Failed to log execution: %v", err)
		http.Error(w, "Failed to log execution", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "logged",
		"id":     execution.ID,
	})
}

func handleGetStats(w http.ResponseWriter, r *http.Request, execRepo *repository.ExecutionRepository, apiRepo *repository.APIRepository) {
	vars := mux.Vars(r)
	apiID := vars["api_id"]

	// Check if API exists
	if _, err := apiRepo.GetByID(apiID); err != nil {
		http.Error(w, "API not found", http.StatusNotFound)
		return
	}

	// Get time period (default: last 24 hours)
	hoursParam := r.URL.Query().Get("hours")
	hours := 24
	if hoursParam != "" {
		if h, err := strconv.Atoi(hoursParam); err == nil {
			hours = h
		}
	}

	since := time.Now().Add(-time.Duration(hours) * time.Hour)
	stats, err := execRepo.GetStats(apiID, since)
	if err != nil {
		log.Printf("Failed to get stats: %v", err)
		http.Error(w, "Failed to get stats", http.StatusInternalServerError)
		return
	}

	stats["api_id"] = apiID
	stats["period_hours"] = hours

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

func handleGetHistory(w http.ResponseWriter, r *http.Request, execRepo *repository.ExecutionRepository, apiRepo *repository.APIRepository) {
	vars := mux.Vars(r)
	apiID := vars["api_id"]

	// Check if API exists
	if _, err := apiRepo.GetByID(apiID); err != nil {
		http.Error(w, "API not found", http.StatusNotFound)
		return
	}

	// Get limit (default: 100)
	limitParam := r.URL.Query().Get("limit")
	limit := 100
	if limitParam != "" {
		if l, err := strconv.Atoi(limitParam); err == nil && l > 0 && l <= 1000 {
			limit = l
		}
	}

	executions, err := execRepo.GetByAPIID(apiID, limit)
	if err != nil {
		log.Printf("Failed to get execution history: %v", err)
		http.Error(w, "Failed to get execution history", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"api_id":     apiID,
		"count":      len(executions),
		"executions": executions,
	})
}
