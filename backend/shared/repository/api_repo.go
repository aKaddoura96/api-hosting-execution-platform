package repository

import (
	"database/sql"
	"fmt"

	"github.com/aKaddoura96/api-hosting-execution-platform/backend/shared/models"
	"github.com/google/uuid"
)

type APIRepository struct {
	db *sql.DB
}

func NewAPIRepository(db *sql.DB) *APIRepository {
	return &APIRepository{db: db}
}

func (r *APIRepository) Create(api *models.API) error {
	api.ID = uuid.New().String()
	
	query := `
		INSERT INTO apis (id, user_id, name, description, version, runtime, visibility, status, endpoint, code_path)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING created_at, updated_at
	`
	
	return r.db.QueryRow(
		query, api.ID, api.UserID, api.Name, api.Description, api.Version,
		api.Runtime, api.Visibility, api.Status, api.Endpoint, api.CodePath,
	).Scan(&api.CreatedAt, &api.UpdatedAt)
}

func (r *APIRepository) GetByID(id string) (*models.API, error) {
	api := &models.API{}
	
	query := `
		SELECT id, user_id, name, description, version, runtime, visibility, status,
		       endpoint, COALESCE(code_path, ''), COALESCE(container_id, ''), created_at, updated_at
		FROM apis WHERE id = $1
	`
	
	err := r.db.QueryRow(query, id).Scan(
		&api.ID, &api.UserID, &api.Name, &api.Description, &api.Version,
		&api.Runtime, &api.Visibility, &api.Status, &api.Endpoint, &api.CodePath,
		&api.ContainerID, &api.CreatedAt, &api.UpdatedAt,
	)
	
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("API not found")
	}
	
	return api, err
}

func (r *APIRepository) GetByUserID(userID string) ([]*models.API, error) {
	query := `
		SELECT id, user_id, name, description, version, runtime, visibility, status,
		       endpoint, COALESCE(code_path, ''), COALESCE(container_id, ''), created_at, updated_at
		FROM apis WHERE user_id = $1
		ORDER BY created_at DESC
	`
	
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var apis []*models.API
	for rows.Next() {
		api := &models.API{}
		err := rows.Scan(
			&api.ID, &api.UserID, &api.Name, &api.Description, &api.Version,
			&api.Runtime, &api.Visibility, &api.Status, &api.Endpoint, &api.CodePath,
			&api.ContainerID, &api.CreatedAt, &api.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		apis = append(apis, api)
	}
	
	return apis, nil
}

func (r *APIRepository) GetPublicAPIs() ([]*models.API, error) {
	query := `
		SELECT id, user_id, name, description, version, runtime, visibility, status,
		       endpoint, COALESCE(code_path, ''), COALESCE(container_id, ''), created_at, updated_at
		FROM apis WHERE visibility = 'public' AND status = 'deployed'
		ORDER BY created_at DESC
	`
	
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var apis []*models.API
	for rows.Next() {
		api := &models.API{}
		err := rows.Scan(
			&api.ID, &api.UserID, &api.Name, &api.Description, &api.Version,
			&api.Runtime, &api.Visibility, &api.Status, &api.Endpoint, &api.CodePath,
			&api.ContainerID, &api.CreatedAt, &api.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		apis = append(apis, api)
	}
	
	return apis, nil
}

func (r *APIRepository) UpdateStatus(id, status, containerID string) error {
	query := `
		UPDATE apis 
		SET status = $1, container_id = NULLIF($2, ''), updated_at = CURRENT_TIMESTAMP
		WHERE id = $3
	`
	
	_, err := r.db.Exec(query, status, containerID, id)
	return err
}

func (r *APIRepository) UpdateCodePath(id, codePath string) error {
	query := `
		UPDATE apis 
		SET code_path = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`
	
	_, err := r.db.Exec(query, codePath, id)
	return err
}

func (r *APIRepository) Update(api *models.API) error {
	query := `
		UPDATE apis 
		SET name = $1, description = $2, visibility = $3, endpoint = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $5
	`
	_, err := r.db.Exec(query, api.Name, api.Description, api.Visibility, api.Endpoint, api.ID)
	return err
}

func (r *APIRepository) Delete(id string) error {
	query := `DELETE FROM apis WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
