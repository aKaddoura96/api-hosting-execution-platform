package repository

import (
	"database/sql"
	"time"

	"github.com/aKaddoura96/api-hosting-execution-platform/backend/shared/models"
	"github.com/google/uuid"
)

type ExecutionRepository struct {
	db *sql.DB
}

func NewExecutionRepository(db *sql.DB) *ExecutionRepository {
	return &ExecutionRepository{db: db}
}

func (r *ExecutionRepository) Create(execution *models.Execution) error {
	execution.ID = uuid.New().String()
	
	query := `
		INSERT INTO executions (id, api_id, user_id, status_code, duration, request_size, response_size, error)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING executed_at
	`
	
	return r.db.QueryRow(
		query, execution.ID, execution.APIID, execution.UserID,
		execution.StatusCode, execution.Duration.Milliseconds(),
		execution.RequestSize, execution.ResponseSize, execution.Error,
	).Scan(&execution.ExecutedAt)
}

func (r *ExecutionRepository) GetByAPIID(apiID string, limit int) ([]*models.Execution, error) {
	query := `
		SELECT id, api_id, user_id, status_code, duration, request_size, response_size, error, executed_at
		FROM executions
		WHERE api_id = $1
		ORDER BY executed_at DESC
		LIMIT $2
	`
	
	rows, err := r.db.Query(query, apiID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var executions []*models.Execution
	for rows.Next() {
		exec := &models.Execution{}
		var durationMs int64
		var userID sql.NullString
		
		err := rows.Scan(
			&exec.ID, &exec.APIID, &userID, &exec.StatusCode, &durationMs,
			&exec.RequestSize, &exec.ResponseSize, &exec.Error, &exec.ExecutedAt,
		)
		if err != nil {
			return nil, err
		}
		
		if userID.Valid {
			exec.UserID = userID.String
		}
		exec.Duration = time.Duration(durationMs) * time.Millisecond
		
		executions = append(executions, exec)
	}
	
	return executions, nil
}

func (r *ExecutionRepository) GetStats(apiID string, since time.Time) (map[string]interface{}, error) {
	query := `
		SELECT 
			COUNT(*) as total_requests,
			AVG(duration) as avg_duration,
			MIN(duration) as min_duration,
			MAX(duration) as max_duration,
			COUNT(CASE WHEN status_code >= 200 AND status_code < 300 THEN 1 END) as success_count,
			COUNT(CASE WHEN status_code >= 400 THEN 1 END) as error_count
		FROM executions
		WHERE api_id = $1 AND executed_at >= $2
	`
	
	var totalRequests, successCount, errorCount int64
	var avgDuration, minDuration, maxDuration sql.NullFloat64
	
	err := r.db.QueryRow(query, apiID, since).Scan(
		&totalRequests, &avgDuration, &minDuration, &maxDuration,
		&successCount, &errorCount,
	)
	if err != nil {
		return nil, err
	}
	
	stats := map[string]interface{}{
		"total_requests": totalRequests,
		"success_count":  successCount,
		"error_count":    errorCount,
		"success_rate":   0.0,
	}
	
	if totalRequests > 0 {
		stats["success_rate"] = float64(successCount) / float64(totalRequests) * 100
	}
	
	if avgDuration.Valid {
		stats["avg_duration_ms"] = avgDuration.Float64
	}
	if minDuration.Valid {
		stats["min_duration_ms"] = minDuration.Float64
	}
	if maxDuration.Valid {
		stats["max_duration_ms"] = maxDuration.Float64
	}
	
	return stats, nil
}
