package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/aKaddoura96/api-hosting-execution-platform/backend/shared/repository"
	"github.com/gorilla/mux"
)

type DeployHandler struct {
	apiRepo      *repository.APIRepository
	executorURL string
}

func NewDeployHandler(apiRepo *repository.APIRepository) *DeployHandler {
	executorURL := os.Getenv("EXECUTOR_URL")
	if executorURL == "" {
		executorURL = "http://localhost:8081"
	}

	return &DeployHandler{
		apiRepo:     apiRepo,
		executorURL: executorURL,
	}
}

type DeployRequest struct {
	APIID string `json:"api_id"`
}

type DeployResponse struct {
	Status      string `json:"status"`
	ContainerID string `json:"container_id,omitempty"`
	Message     string `json:"message"`
}

func (h *DeployHandler) DeployAPI(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	apiID := vars["id"]
	userID := r.Context().Value("user_id").(string)

	// Get API details
	api, err := h.apiRepo.GetByID(apiID)
	if err != nil {
		http.Error(w, "API not found", http.StatusNotFound)
		return
	}

	// Check ownership
	if api.UserID != userID {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	// Check if code is uploaded
	if api.CodePath == "" {
		http.Error(w, "Please upload code before deploying", http.StatusBadRequest)
		return
	}

	// Call executor service
	deployReq := DeployRequest{APIID: apiID}
	reqBody, _ := json.Marshal(deployReq)

	resp, err := http.Post(
		h.executorURL+"/deploy",
		"application/json",
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		http.Error(w, "Failed to communicate with executor service", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Deployment failed", resp.StatusCode)
		return
	}

	var deployResp DeployResponse
	if err := json.NewDecoder(resp.Body).Decode(&deployResp); err != nil {
		http.Error(w, "Failed to parse response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(deployResp)
}

func (h *DeployHandler) StopAPI(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	apiID := vars["id"]
	userID := r.Context().Value("user_id").(string)

	// Get API details
	api, err := h.apiRepo.GetByID(apiID)
	if err != nil {
		http.Error(w, "API not found", http.StatusNotFound)
		return
	}

	// Check ownership
	if api.UserID != userID {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	// Call executor service
	url := fmt.Sprintf("%s/stop/%s", h.executorURL, apiID)
	req, _ := http.NewRequest("POST", url, nil)
	
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to communicate with executor service", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Stop failed", resp.StatusCode)
		return
	}

	var stopResp DeployResponse
	if err := json.NewDecoder(resp.Body).Decode(&stopResp); err != nil {
		http.Error(w, "Failed to parse response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stopResp)
}

func (h *DeployHandler) GetAPIStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	apiID := vars["id"]

	// Call executor service
	url := fmt.Sprintf("%s/status/%s", h.executorURL, apiID)
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, "Failed to communicate with executor service", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Failed to get status", resp.StatusCode)
		return
	}

	var statusResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&statusResp); err != nil {
		http.Error(w, "Failed to parse response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(statusResp)
}
