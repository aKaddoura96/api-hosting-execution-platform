package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/aKaddoura96/api-hosting-execution-platform/backend/shared/repository"
)

// APIKeyMiddleware validates API keys for public API invocation
func APIKeyMiddleware(apiKeyRepo *repository.APIKeyRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract API key from header
			apiKey := r.Header.Get("X-API-Key")
			if apiKey == "" {
				// Also check Authorization header with "Bearer" prefix
				authHeader := r.Header.Get("Authorization")
				if strings.HasPrefix(authHeader, "Bearer ") {
					apiKey = strings.TrimPrefix(authHeader, "Bearer ")
				}
			}

			// For now, allow requests without API key (for testing)
			// In production, you'd want to enforce this more strictly
			if apiKey == "" {
				// Allow but add a flag
				ctx := context.WithValue(r.Context(), "has_api_key", false)
				next.ServeHTTP(w, r.WithContext(ctx))
				return
			}

			// Validate API key
			key, err := apiKeyRepo.GetByKey(apiKey)
			if err != nil {
				http.Error(w, "Invalid API key", http.StatusUnauthorized)
				return
			}

			// Check if key is active
			if !key.IsActive {
				http.Error(w, "API key is inactive", http.StatusUnauthorized)
				return
			}

			// Check if key is expired
			if key.ExpiresAt != nil && key.ExpiresAt.Before(time.Now()) {
				http.Error(w, "API key has expired", http.StatusUnauthorized)
				return
			}

			// Add key info to context
			ctx := context.WithValue(r.Context(), "api_key", key)
			ctx = context.WithValue(ctx, "has_api_key", true)
			ctx = context.WithValue(ctx, "api_key_user_id", key.UserID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
