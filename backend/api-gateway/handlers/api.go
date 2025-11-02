package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/aKaddoura96/api-hosting-execution-platform/backend/shared/models"
	"github.com/aKaddoura96/api-hosting-execution-platform/backend/shared/repository"
	"github.com/gorilla/mux"
)

type APIHandler struct {
	apiRepo *repository.APIRepository
}

func NewAPIHandler(apiRepo *repository.APIRepository) *APIHandler {
	return &APIHandler{apiRepo: apiRepo}
}

type CreateAPIRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Version     string `json:"version"`
	Runtime     string `json:"runtime"`
	Visibility  string `json:"visibility"`
}

func (h *APIHandler) CreateAPI(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	var req CreateAPIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if req.Name == "" || req.Runtime == "" {
		http.Error(w, "Name and runtime are required", http.StatusBadRequest)
		return
	}

	if req.Version == "" {
		req.Version = "v1"
	}

	if req.Visibility == "" {
		req.Visibility = "private"
	}

	// Generate endpoint URL
	endpoint := fmt.Sprintf("/execute/%s/%s", userID[:8], req.Name)

	// Create API record
	api := &models.API{
		UserID:      userID,
		Name:        req.Name,
		Description: req.Description,
		Version:     req.Version,
		Runtime:     req.Runtime,
		Visibility:  req.Visibility,
		Status:      "pending",
		Endpoint:    endpoint,
		CodePath:    "", // Will be set on upload
	}

	if err := h.apiRepo.Create(api); err != nil {
		http.Error(w, "Failed to create API: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(api)
}

func (h *APIHandler) UploadCode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	apiID := vars["id"]
	userID := r.Context().Value("user_id").(string)

	// Get API
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

	// Parse multipart form
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB max
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("code")
	if err != nil {
		http.Error(w, "Code file is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Create upload directory
	uploadDir := filepath.Join("uploads", userID, apiID)
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		http.Error(w, "Failed to create upload directory", http.StatusInternalServerError)
		return
	}

	// Save file
	codePath := filepath.Join(uploadDir, header.Filename)
	dst, err := os.Create(codePath)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	// Update API record with code path
	if err := h.apiRepo.UpdateCodePath(apiID, codePath); err != nil {
		http.Error(w, "Failed to update code path", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Code uploaded successfully",
		"path":    codePath,
	})
}

func (h *APIHandler) GetMyAPIs(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	apis, err := h.apiRepo.GetByUserID(userID)
	if err != nil {
		http.Error(w, "Failed to get APIs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(apis)
}

func (h *APIHandler) GetPublicAPIs(w http.ResponseWriter, r *http.Request) {
	apis, err := h.apiRepo.GetPublicAPIs()
	if err != nil {
		http.Error(w, "Failed to get APIs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(apis)
}

func (h *APIHandler) GetAPI(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	apiID := vars["id"]

	api, err := h.apiRepo.GetByID(apiID)
	if err != nil {
		http.Error(w, "API not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(api)
}

func (h *APIHandler) DeleteAPI(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	apiID := vars["id"]
	userID := r.Context().Value("user_id").(string)

	// Get API
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

	// TODO: Stop container if running

	// Delete API
	if err := h.apiRepo.Delete(apiID); err != nil {
		http.Error(w, "Failed to delete API", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
