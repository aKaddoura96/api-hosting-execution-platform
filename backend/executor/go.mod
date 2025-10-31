module github.com/aKaddoura96/api-hosting-execution-platform/backend/executor

go 1.21

require (
	github.com/aKaddoura96/api-hosting-execution-platform/backend/shared v0.0.0
	github.com/docker/docker v24.0.7+incompatible
	github.com/docker/go-connections v0.4.0
	github.com/gorilla/mux v1.8.1
	github.com/joho/godotenv v1.5.1
)

replace github.com/aKaddoura96/api-hosting-execution-platform/backend/shared => ../shared
