module github.com/aKaddoura96/api-hosting-execution-platform/backend/analytics

go 1.21

require (
	github.com/aKaddoura96/api-hosting-execution-platform/backend/shared v0.0.0
	github.com/gorilla/mux v1.8.1
	github.com/joho/godotenv v1.5.1
)

replace github.com/aKaddoura96/api-hosting-execution-platform/backend/shared => ../shared
