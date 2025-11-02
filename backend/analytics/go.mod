module github.com/aKaddoura96/api-hosting-execution-platform/backend/analytics

go 1.21

require (
	github.com/aKaddoura96/api-hosting-execution-platform/backend/shared v0.0.0
	github.com/gorilla/mux v1.8.1
	github.com/joho/godotenv v1.5.1
)

require (
	github.com/google/uuid v1.5.0 // indirect
	github.com/lib/pq v1.10.9 // indirect
	golang.org/x/crypto v0.17.0 // indirect
)

replace github.com/aKaddoura96/api-hosting-execution-platform/backend/shared => ../shared
