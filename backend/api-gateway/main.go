package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

type HealthResponse struct {
	Status    string    `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	Service   string    `json:"service"`
}

func main() {
	// Load environment variables
	godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	router := mux.NewRouter()

	// Health endpoint
	router.HandleFunc("/health", healthHandler).Methods("GET")
	
	// API routes (to be implemented)
	router.HandleFunc("/api/v1/endpoints", listEndpointsHandler).Methods("GET")

	log.Printf("API Gateway starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
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

func listEndpointsHandler(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement endpoint listing
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Endpoint listing coming soon",
	})
}
