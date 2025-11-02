# API Invocation & Deployment Guide

## 🎯 Overview

This guide explains how the API deployment and execution system works, how to invoke APIs, and troubleshoot common issues.

---

## 🚀 The Complete Flow

### 1. Create an API
**UI:** Dashboard → "+ Create API" button

**What happens:**
- API record created in database with status="pending"
- Unique endpoint generated: `/execute/{userID}/{apiName}`
- No code uploaded yet

### 2. Upload Code
**UI:** Dashboard → Click API card → "Upload Code" section

**What happens:**
- File uploaded to `uploads/{userID}/{apiID}/filename`
- `code_path` field updated in database
- API still has status="pending"

### 3. Deploy API
**UI:** API Detail Page → "🚀 Deploy" button

**Complete Flow:**
```
Frontend (port 3000)
    ↓ POST /api/v1/apis/{id}/deploy
API Gateway (port 8080)
    ↓ POST /deploy with {api_id: "..."}
Executor Service (port 8081)
    ↓ Validates code exists
    ↓ Updates status to "deployed"
    ↓ Returns success
Database
    ↓ status: "pending" → "deployed"
```

**What Actually Happens:**
- **Current Implementation (Simplified):**
  - Just marks the API as "deployed" in database
  - No persistent container created
  - Code executes on-demand when you test it

- **Future Implementation (Production):**
  - Would create a persistent Docker container
  - Container runs 24/7 waiting for requests
  - Has its own endpoint that proxies to the container

### 4. Test/Invoke API

#### **Option A: Test in UI**
**UI:** API Detail Page → "Test API" section → Write code → "▶️ Run Code"

**Flow:**
```
Frontend
    ↓ POST http://localhost:8081/execute
    ↓ Body: { code: "...", runtime: "python" }
Executor Service
    ↓ Creates temporary Docker container
    ↓ Runs code in isolated environment
    ↓ Captures output
    ↓ Destroys container
    ↓ Returns result
Frontend
    ↓ Displays output
```

**Resource Limits:**
- Memory: 256MB
- CPU: 0.5 cores
- Max processes: 50
- Timeout: 30 seconds (default)
- Network: Isolated

#### **Option B: Invoke via cURL**

**Test Execution (any code):**
```bash
curl -X POST http://localhost:8081/execute \
  -H "Content-Type: application/json" \
  -d '{
    "code": "print(\"Hello World\")",
    "runtime": "python",
    "timeout_sec": 10
  }'
```

**Invoke Deployed API (NOT YET IMPLEMENTED):**
```bash
# This would be the future implementation
curl -X POST http://localhost:8080/execute/{userID}/{apiName} \
  -H "Content-Type: application/json" \
  -H "X-API-Key: your-api-key" \
  -d '{"input": {"param": "value"}}'
```

### 5. Stop API
**UI:** API Detail Page → "⏸️ Stop" button (appears when deployed)

**What happens:**
- Updates status from "deployed" → "stopped"
- In future: would stop the persistent container
- API no longer accepts requests

---

## 🔧 Understanding the Services

### **Port 3000: Frontend (Next.js)**
- User interface
- React components
- Talks to API Gateway

### **Port 8080: API Gateway (Go)**
- Authentication (JWT)
- API CRUD operations
- Routes to executor service
- Handles user requests

### **Port 8081: Executor Service (Go)**
- Code execution engine
- Docker-in-Docker
- Sandboxed containers
- Security isolation

### **Port 5432: PostgreSQL**
- Stores users, APIs, keys
- Execution logs
- Analytics data

### **Port 6379: Redis**
- Caching
- Rate limiting
- Session storage

---

## 🐛 Common Issues & Solutions

### Issue 1: CORS Error on "Run Code"
**Error:**
```
Access to fetch at 'http://localhost:8081/execute' from origin 
'http://localhost:3000' has been blocked by CORS policy
```

**Solution:** ✅ FIXED!
Added CORS middleware to executor service that allows:
- Origin: `http://localhost:3000`
- Methods: GET, POST, PUT, DELETE, OPTIONS
- Headers: Content-Type, Authorization

**To test fix:**
```bash
docker compose restart executor
```

### Issue 2: "Stop" Button Fails
**Error:** "Failed to stop: ..."

**Why it happens:**
- Current implementation doesn't create persistent containers
- Nothing to stop

**Solution (Current):**
- This is expected behavior
- Status updates to "stopped" in database
- You can re-deploy later

**Future Implementation:**
- Will actually stop running Docker containers
- Container lifecycle management

### Issue 3: Deployed API Has No Endpoint
**Current State:**
- Deploy button marks status as "deployed"
- But no actual endpoint to invoke the API yet

**Missing Piece:**
- Need to implement `/execute/{userID}/{apiName}` endpoint
- Should read code from `code_path`
- Execute it and return results

**Quick Implementation Idea:**
```go
// In API Gateway
router.HandleFunc("/execute/{userID}/{apiName}", func(w http.ResponseWriter, r *http.Request) {
    // 1. Find API by endpoint
    // 2. Read code from code_path
    // 3. Call executor service with code
    // 4. Return result
}).Methods("POST")
```

---

## 📊 How Different Runtimes Work

### Python APIs
```python
# Your uploaded file: weather.py
def main(input_data):
    city = input_data.get('city', 'Dubai')
    # Call weather API
    return {"temp": 25, "city": city}

# When invoked:
# - Docker runs: python:3.9-alpine
# - Executes your code
# - Returns JSON result
```

### Node.js APIs
```javascript
// Your uploaded file: api.js
module.exports = async function(input) {
    const { name } = input;
    // Your logic here
    return { message: `Hello ${name}` };
};

// When invoked:
// - Docker runs: node:18-alpine
// - Executes your code
// - Returns JSON result
```

### Go APIs
```go
// Your uploaded file: main.go
package main

import "encoding/json"

func main() {
    // Your API logic
    result := map[string]interface{}{
        "status": "success",
    }
    json.NewEncoder(os.Stdout).Encode(result)
}

// When invoked:
// - Docker runs: golang:1.23-alpine
// - Compiles and executes
// - Returns JSON result
```

---

## 🔐 Security Features

### Isolation
- Each execution runs in a fresh Docker container
- Network isolated (no internet access)
- Filesystem isolated
- Process limits enforced

### Resource Limits
```go
resources := container.Resources{
    Memory:     256 * 1024 * 1024, // 256MB
    MemorySwap: 256 * 1024 * 1024,
    CPUPeriod:  100000,
    CPUQuota:   50000, // 0.5 CPU
    PidsLimit:  50,    // Max processes
}
```

### Timeout Protection
- Default: 30 seconds
- Configurable per request
- Container force-killed after timeout

---

## 🚦 API Lifecycle States

```
┌─────────┐    Upload Code    ┌─────────┐
│ pending │ ──────────────────→│ pending │
└─────────┘                    │(w/ code)│
                               └─────────┘
                                    │
                            Click Deploy
                                    ↓
                               ┌─────────┐
                               │deployed │
                               └─────────┘
                                    │
                             Click Stop
                                    ↓
                               ┌─────────┐
                               │ stopped │
                               └─────────┘
                                    │
                           Click Deploy Again
                                    ↓
                               ┌─────────┐
                               │deployed │
                               └─────────┘
```

---

## 🧪 Testing Your API

### 1. Test in UI (Recommended)
1. Go to API detail page
2. Scroll to "Test API" section
3. Write your code or use sample
4. Click "▶️ Run Code"
5. See results instantly

### 2. Test with cURL
```bash
# Python example
curl -X POST http://localhost:8081/execute \
  -H "Content-Type: application/json" \
  -d '{
    "code": "import json; print(json.dumps({\"result\": \"success\"}))",
    "runtime": "python"
  }'

# Node.js example
curl -X POST http://localhost:8081/execute \
  -H "Content-Type: application/json" \
  -d '{
    "code": "console.log(JSON.stringify({result: \"success\"}))",
    "runtime": "nodejs"
  }'
```

### 3. Test with Postman
- Method: POST
- URL: `http://localhost:8081/execute`
- Headers: `Content-Type: application/json`
- Body (raw JSON):
```json
{
  "code": "print('Hello from Postman')",
  "runtime": "python",
  "timeout_sec": 10
}
```

---

## 📈 What's Next?

### Immediate Todo (To Make It Production-Ready):

1. **Implement Actual Invocation Endpoint**
   - Add `/execute/{userID}/{apiName}` in API Gateway
   - Read code from uploaded file
   - Execute via executor service
   - Return results to caller

2. **API Key Authentication**
   - Generate keys for each user
   - Validate on each API call
   - Track usage per key

3. **Rate Limiting**
   - Limit calls per API key
   - Prevent abuse
   - Different tiers (free/paid)

4. **Analytics Tracking**
   - Log every API call
   - Track latency
   - Count errors
   - Monitor usage

5. **Persistent Containers (Optional)**
   - For faster response times
   - Keep containers warm
   - Better for high-traffic APIs

---

## 💡 Quick Reference

### Check Service Status
```bash
docker compose ps
```

### View Logs
```bash
# All services
docker compose logs -f

# Specific service
docker compose logs -f executor
docker compose logs -f api-gateway
docker compose logs -f frontend
```

### Restart After Changes
```bash
# Rebuild and restart
docker compose up --build -d

# Just restart
docker compose restart

# Restart specific service
docker compose restart executor
```

### Test Executor Directly
```bash
curl http://localhost:8081/health
```

---

## 🎓 Summary

**Current State:**
- ✅ Create APIs
- ✅ Upload code
- ✅ Deploy (marks as deployed)
- ✅ Test code in UI (now works with CORS fix!)
- ✅ Stop APIs (updates status)
- ❌ No actual endpoint to invoke deployed APIs yet

**To Invoke a Deployed API:**
1. Upload your code file
2. Click Deploy
3. **Currently:** Test it in the UI test section
4. **Future:** Call the API's unique endpoint with API key

**The deploy button makes your API "live"** but the invocation endpoint still needs to be implemented for external calls.
