package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/aKaddoura96/api-hosting-execution-platform/backend/shared/models"
	"github.com/aKaddoura96/api-hosting-execution-platform/backend/shared/repository"
	"github.com/gorilla/mux"
)

type ExecuteHandler struct {
	apiRepo     *repository.APIRepository
	executorURL string
}

func NewExecuteHandler(apiRepo *repository.APIRepository) *ExecuteHandler {
	executorURL := os.Getenv("EXECUTOR_URL")
	if executorURL == "" {
		executorURL = "http://localhost:8081"
	}

	return &ExecuteHandler{
		apiRepo:     apiRepo,
		executorURL: executorURL,
	}
}

type ExecuteRequest struct {
	Input      map[string]interface{} `json:"input"`
	TimeoutSec int                    `json:"timeout_sec,omitempty"`
}

type ExecuteResponse struct {
	Output     string                 `json:"output"`
	Error      string                 `json:"error,omitempty"`
	StatusCode int                    `json:"status_code"`
	DurationMS int                    `json:"duration_ms"`
	ExitCode   int                    `json:"exit_code"`
	Result     map[string]interface{} `json:"result,omitempty"`
}

// ExecuteAPI handles requests to invoke a deployed API
func (h *ExecuteHandler) ExecuteAPI(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	endpoint := r.URL.Path // Full path like /execute/abc12345/my-api

	// Find API by endpoint
	apis, err := h.apiRepo.GetPublicAPIs() // Get all public APIs
	if err != nil {
		http.Error(w, "Failed to lookup API", http.StatusInternalServerError)
		return
	}

	// Try to find matching API by endpoint
	var targetAPI *models.API
	for _, api := range apis {
		if api.Endpoint == endpoint {
			targetAPI = api
			break
		}
	}

	// If not found in public, try all APIs (for testing)
	if targetAPI == nil {
		userID := vars["user_id"]
		if userID != "" {
			allAPIs, _ := h.apiRepo.GetByUserID(userID)
			for _, api := range allAPIs {
				if api.Endpoint == endpoint {
					targetAPI = api
					break
				}
			}
		}
	}

	if targetAPI == nil {
		http.Error(w, "API not found", http.StatusNotFound)
		return
	}

	// Check if API is deployed
	if targetAPI.Status != "deployed" {
		http.Error(w, fmt.Sprintf("API is not deployed. Current status: %s", targetAPI.Status), http.StatusBadRequest)
		return
	}

	// Check if code exists
	if targetAPI.CodePath == "" {
		http.Error(w, "No code uploaded for this API", http.StatusBadRequest)
		return
	}

	// Read the uploaded code
	codeBytes, err := os.ReadFile(targetAPI.CodePath)
	if err != nil {
		http.Error(w, "Failed to read API code", http.StatusInternalServerError)
		return
	}

	// Parse input from request
	var execReq ExecuteRequest
	if r.Body != nil {
		json.NewDecoder(r.Body).Decode(&execReq)
	}

	// Prepare executor request
	executorReq := map[string]interface{}{
		"code":    string(codeBytes),
		"runtime": targetAPI.Runtime,
		"input":   execReq.Input,
	}

	if execReq.TimeoutSec > 0 {
		executorReq["timeout_sec"] = execReq.TimeoutSec
	}

	// Call executor service
	reqBody, _ := json.Marshal(executorReq)
	resp, err := http.Post(
		h.executorURL+"/execute",
		"application/json",
		bytes.NewBuffer(reqBody),
	)
	if err != nil {
		http.Error(w, "Failed to execute API code", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read response
	respBody, _ := io.ReadAll(resp.Body)

	// Return result
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	w.Write(respBody)
}

// GetAPIByEndpoint is a helper to find API by its endpoint
func (h *ExecuteHandler) GetAPIByEndpoint(endpoint string) (*models.API, error) {
	// This is a simple implementation - in production you'd want to optimize this
	// by adding an index on the endpoint column or using a cache

	// Try public APIs first
	apis, err := h.apiRepo.GetPublicAPIs()
	if err != nil {
		return nil, err
	}

	for _, api := range apis {
		if api.Endpoint == endpoint {
			return api, nil
		}
	}

	return nil, fmt.Errorf("API not found")
}
