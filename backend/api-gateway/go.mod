module github.com/aKaddoura96/api-hosting-execution-platform/backend/api-gateway

go 1.21

require (
	github.com/aKaddoura96/api-hosting-execution-platform/backend/shared v0.0.0
	github.com/gorilla/mux v1.8.1
	github.com/joho/godotenv v1.5.1
	github.com/rs/cors v1.11.0
)

require (
	github.com/golang-jwt/jwt/v5 v5.2.0 // indirect
	github.com/google/uuid v1.5.0 // indirect
	github.com/lib/pq v1.10.9 // indirect
	golang.org/x/crypto v0.17.0 // indirect
)

replace github.com/aKaddoura96/api-hosting-execution-platform/backend/shared => ../shared
