package repository

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"

	"github.com/aKaddoura96/api-hosting-execution-platform/backend/shared/models"
	"github.com/google/uuid"
)

type APIKeyRepository struct {
	db *sql.DB
}

func NewAPIKeyRepository(db *sql.DB) *APIKeyRepository {
	return &APIKeyRepository{db: db}
}

func (r *APIKeyRepository) Create(apiKey *models.APIKey) error {
	apiKey.ID = uuid.New().String()
	
	// Generate a random API key
	if apiKey.Key == "" {
		key, err := generateAPIKey()
		if err != nil {
			return err
		}
		apiKey.Key = key
	}
	
	query := `
		INSERT INTO api_keys (id, user_id, api_id, key, name, is_active, expires_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING created_at
	`
	
	return r.db.QueryRow(
		query, apiKey.ID, apiKey.UserID, apiKey.APIID, apiKey.Key,
		apiKey.Name, apiKey.IsActive, apiKey.ExpiresAt,
	).Scan(&apiKey.CreatedAt)
}

func (r *APIKeyRepository) GetByKey(key string) (*models.APIKey, error) {
	apiKey := &models.APIKey{}
	
	query := `
		SELECT id, user_id, api_id, key, name, is_active, expires_at, created_at
		FROM api_keys WHERE key = $1
	`
	
	var apiID sql.NullString
	err := r.db.QueryRow(query, key).Scan(
		&apiKey.ID, &apiKey.UserID, &apiID, &apiKey.Key,
		&apiKey.Name, &apiKey.IsActive, &apiKey.ExpiresAt, &apiKey.CreatedAt,
	)
	
	if apiID.Valid {
		apiKey.APIID = apiID.String
	}
	
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("API key not found")
	}
	
	return apiKey, err
}

func (r *APIKeyRepository) GetByUserID(userID string) ([]*models.APIKey, error) {
	query := `
		SELECT id, user_id, api_id, key, name, is_active, expires_at, created_at
		FROM api_keys WHERE user_id = $1
		ORDER BY created_at DESC
	`
	
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var apiKeys []*models.APIKey
	for rows.Next() {
		apiKey := &models.APIKey{}
		var apiID sql.NullString
		
		err := rows.Scan(
			&apiKey.ID, &apiKey.UserID, &apiID, &apiKey.Key,
			&apiKey.Name, &apiKey.IsActive, &apiKey.ExpiresAt, &apiKey.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		
		if apiID.Valid {
			apiKey.APIID = apiID.String
		}
		
		apiKeys = append(apiKeys, apiKey)
	}
	
	return apiKeys, nil
}

func (r *APIKeyRepository) Deactivate(id string) error {
	query := `UPDATE api_keys SET is_active = false WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func generateAPIKey() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return "apk_" + hex.EncodeToString(bytes), nil
}
