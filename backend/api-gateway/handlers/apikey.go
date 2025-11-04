package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/aKaddoura96/api-hosting-execution-platform/backend/shared/models"
	"github.com/aKaddoura96/api-hosting-execution-platform/backend/shared/repository"
	"github.com/gorilla/mux"
)

type APIKeyHandler struct {
	apiKeyRepo *repository.APIKeyRepository
}

func NewAPIKeyHandler(apiKeyRepo *repository.APIKeyRepository) *APIKeyHandler {
	return &APIKeyHandler{apiKeyRepo: apiKeyRepo}
}

type CreateAPIKeyRequest struct {
	Name   string `json:"name"`
	APIID  string `json:"api_id,omitempty"`
}

// CreateAPIKey generates a new API key for the user
func (h *APIKeyHandler) CreateAPIKey(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	var req CreateAPIKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate name
	if req.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	// Create API key
	apiKey := &models.APIKey{
		UserID:   userID,
		APIID:    req.APIID,
		Name:     req.Name,
		IsActive: true,
	}

	if err := h.apiKeyRepo.Create(apiKey); err != nil {
		http.Error(w, "Failed to create API key", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(apiKey)
}

// GetMyAPIKeys returns all API keys for the authenticated user
func (h *APIKeyHandler) GetMyAPIKeys(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(string)

	apiKeys, err := h.apiKeyRepo.GetByUserID(userID)
	if err != nil {
		http.Error(w, "Failed to get API keys", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(apiKeys)
}

// DeactivateAPIKey deactivates an API key
func (h *APIKeyHandler) DeactivateAPIKey(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	keyID := vars["id"]
	userID := r.Context().Value("user_id").(string)

	// Get all user keys to verify ownership
	apiKeys, err := h.apiKeyRepo.GetByUserID(userID)
	if err != nil {
		http.Error(w, "Failed to verify ownership", http.StatusInternalServerError)
		return
	}

	// Check if key belongs to user
	found := false
	for _, key := range apiKeys {
		if key.ID == keyID {
			found = true
			break
		}
	}

	if !found {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	// Deactivate the key
	if err := h.apiKeyRepo.Deactivate(keyID); err != nil {
		http.Error(w, "Failed to deactivate API key", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "API key deactivated successfully",
	})
}
