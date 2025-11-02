package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/aKaddoura96/api-hosting-execution-platform/backend/executor/container"
	"github.com/aKaddoura96/api-hosting-execution-platform/backend/shared/database"
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
	godotenv.Load()

	// Connect to database
	if err := database.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Initialize container manager
	containerMgr, err := container.NewManager()
	if err != nil {
		log.Fatalf("Failed to initialize container manager: %v", err)
	}
	defer containerMgr.Close()

	// Initialize repositories
	apiRepo := repository.NewAPIRepository(database.DB)

	// Setup router
	router := mux.NewRouter()

	// Health check
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "healthy",
			"service": "executor",
		})
	}).Methods("GET")

	// Deploy API
	router.HandleFunc("/deploy", func(w http.ResponseWriter, r *http.Request) {
		handleDeploy(w, r, containerMgr, apiRepo)
	}).Methods("POST")

	// Stop API
	router.HandleFunc("/stop/{api_id}", func(w http.ResponseWriter, r *http.Request) {
		handleStop(w, r, containerMgr, apiRepo)
	}).Methods("POST")

	// Get API status
	router.HandleFunc("/status/{api_id}", func(w http.ResponseWriter, r *http.Request) {
		handleStatus(w, r, containerMgr, apiRepo)
	}).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("Executor service starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func handleDeploy(w http.ResponseWriter, r *http.Request, mgr *container.Manager, apiRepo *repository.APIRepository) {
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

	// Deploy container
	containerID, err := mgr.DeployAPI(api.ID, api.Runtime, api.CodePath)
	if err != nil {
		log.Printf("Deployment failed: %v", err)
		
		// Update API status to failed
		apiRepo.UpdateStatus(api.ID, "failed", "")
		
		http.Error(w, "Deployment failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Update API status
	if err := apiRepo.UpdateStatus(api.ID, "deployed", containerID); err != nil {
		log.Printf("Failed to update API status: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(DeployResponse{
		Status:      "success",
		ContainerID: containerID,
		Message:     "API deployed successfully",
	})
}

func handleStop(w http.ResponseWriter, r *http.Request, mgr *container.Manager, apiRepo *repository.APIRepository) {
	vars := mux.Vars(r)
	apiID := vars["api_id"]

	// Get API details
	api, err := apiRepo.GetByID(apiID)
	if err != nil {
		http.Error(w, "API not found", http.StatusNotFound)
		return
	}

	if api.ContainerID == "" {
		http.Error(w, "No container running for this API", http.StatusBadRequest)
		return
	}

	// Stop container
	if err := mgr.StopAPI(api.ContainerID); err != nil {
		http.Error(w, "Failed to stop container: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Update API status
	if err := apiRepo.UpdateStatus(api.ID, "stopped", ""); err != nil {
		log.Printf("Failed to update API status: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(DeployResponse{
		Status:  "success",
		Message: "API stopped successfully",
	})
}

func handleStatus(w http.ResponseWriter, r *http.Request, mgr *container.Manager, apiRepo *repository.APIRepository) {
	vars := mux.Vars(r)
	apiID := vars["api_id"]

	// Get API details
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

	// If container exists, check its status
	if api.ContainerID != "" {
		status, err := mgr.GetContainerStatus(api.ContainerID)
		if err == nil {
			response.Status = status
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
