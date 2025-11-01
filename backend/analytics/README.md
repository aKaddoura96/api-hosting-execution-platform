# Analytics Service

Collects and analyzes API execution metrics, providing usage statistics and performance insights.

## Features

- **Execution Logging** - Track every API request
- **Performance Metrics** - Duration, success rate, error rate
- **Usage Statistics** - Request counts, trends
- **Execution History** - Detailed logs of recent requests

## API Endpoints

### Log Execution
```bash
POST /log
Body: {
  "api_id": "uuid",
  "user_id": "uuid",
  "status_code": 200,
  "duration_ms": 45,
  "request_size": 1024,
  "response_size": 2048,
  "error": ""
}
```

### Get API Statistics
```bash
GET /stats/{api_id}?hours=24
Response: {
  "api_id": "...",
  "period_hours": 24,
  "total_requests": 150,
  "success_count": 145,
  "error_count": 5,
  "success_rate": 96.67,
  "avg_duration_ms": 45.5,
  "min_duration_ms": 12.3,
  "max_duration_ms": 234.6
}
```

### Get Execution History
```bash
GET /history/{api_id}?limit=100
Response: {
  "api_id": "...",
  "count": 100,
  "executions": [...]
}
```

### Health Check
```bash
GET /health
Response: {"status": "healthy", "service": "analytics"}
```

## Running Locally

```bash
cd backend/analytics
cp .env.example .env
go mod download
go run main.go
```

## Environment Variables

```env
PORT=8082
DATABASE_URL=postgres://apiplatform:dev_password@localhost:5432/apiplatform?sslmode=disable
```

## Metrics Collected

- **Request Count** - Total API invocations
- **Success Rate** - Percentage of successful requests
- **Error Rate** - Percentage of failed requests
- **Duration** - Min, max, average response time
- **Payload Size** - Request and response sizes
- **Timestamps** - When requests occurred

## Integration

Other services call the analytics service to log executions:

```go
// Log an API execution
analytics.LogExecution(apiID, userID, statusCode, duration, reqSize, respSize, error)
```

## Dashboards (Future)

- Real-time metrics
- Historical trends
- Per-user analytics
- Geographic distribution
- Popular APIs

## TODO

- [ ] Real-time metrics streaming
- [ ] Time-series data aggregation
- [ ] Alerting on anomalies
- [ ] Export to monitoring systems
- [ ] Usage-based billing integration
