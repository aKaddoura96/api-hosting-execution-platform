# End-to-End Testing Guide

## 🎯 Complete Testing Workflow

This guide walks you through testing the entire platform from signup to API invocation.

---

## ✅ Prerequisites

Make sure all services are running:
```bash
docker compose ps
```

You should see all services "Up" or "Healthy":
- api-platform-postgres (Healthy)
- api-platform-redis (Healthy)
- api-platform-gateway (Up)
- api-platform-executor (Up)
- api-platform-analytics (Up)
- api-platform-frontend (Up)

If not, start them:
```bash
docker compose up -d
```

---

## 📝 Test Scenario: Weather API

We'll create a Python API that returns weather data, deploy it, and invoke it.

---

### **Step 1: Sign Up & Login**

#### UI Method:
1. Open http://localhost:3000
2. Click "Get Started"
3. Fill the form:
   - Name: `Test User`
   - Email: `test@example.com`
   - Password: `SecurePass123!`
   - Role: `Host and monetize my APIs`
4. Click "Create account"
5. You'll be redirected to Dashboard

#### API Method:
```bash
# Signup
curl -X POST http://localhost:8080/api/v1/auth/signup \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test User",
    "email": "test@example.com",
    "password": "SecurePass123!",
    "role": "developer"
  }'

# Login to get token
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "SecurePass123!"
  }'

# Save the token from response
export TOKEN="your-jwt-token-here"
```

---

### **Step 2: Create API**

#### UI Method:
1. In Dashboard, click "+ Create API"
2. Fill the form:
   - Name: `weather-api`
   - Description: `Returns weather data for cities`
   - Runtime: `Python`
   - Visibility: `Public (Free)`
3. Click "Create"
4. API card appears with status "pending"

#### API Method:
```bash
curl -X POST http://localhost:8080/api/v1/apis \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "weather-api",
    "description": "Returns weather data for cities",
    "version": "v1",
    "runtime": "python",
    "visibility": "public"
  }'

# Save the API ID from response
export API_ID="your-api-id-here"
```

---

### **Step 3: Create Code File**

Create a file `weather.py` with this content:

```python
import json
import sys

# Read input from stdin (passed by executor)
input_data = {}
if len(sys.argv) > 1:
    input_data = json.loads(sys.argv[1])

# Get city from input, default to Dubai
city = input_data.get('city', 'Dubai')

# Mock weather data (in production, you'd call a real API)
weather_data = {
    'Dubai': {'temp': 28, 'condition': 'Sunny', 'humidity': 60},
    'London': {'temp': 15, 'condition': 'Cloudy', 'humidity': 75},
    'Tokyo': {'temp': 20, 'condition': 'Rainy', 'humidity': 85},
    'New York': {'temp': 18, 'condition': 'Clear', 'humidity': 55}
}

# Get weather for requested city
weather = weather_data.get(city, {
    'temp': 25,
    'condition': 'Unknown',
    'humidity': 50
})

# Return result
result = {
    'city': city,
    'temperature': weather['temp'],
    'condition': weather['condition'],
    'humidity': weather['humidity'],
    'unit': 'celsius'
}

# Print as JSON (executor captures stdout)
print(json.dumps(result))
```

---

### **Step 4: Upload Code**

#### UI Method:
1. Click on your "weather-api" card
2. Scroll to "Upload Code" section
3. Drag & drop `weather.py` OR click "Browse Files"
4. Wait for "✅ weather.py uploaded successfully!"
5. API status still shows "pending" (needs deploy)

#### API Method:
```bash
curl -X POST http://localhost:8080/api/v1/apis/$API_ID/upload \
  -H "Authorization: Bearer $TOKEN" \
  -F "code=@weather.py"
```

---

### **Step 5: Deploy API**

#### UI Method:
1. Still on API detail page
2. Click "🚀 Deploy" button (top right)
3. Confirm deployment
4. Wait for "✅ API deployed successfully!"
5. Status badge changes to "deployed" (green)
6. Stop button appears

#### API Method:
```bash
curl -X POST http://localhost:8080/api/v1/apis/$API_ID/deploy \
  -H "Authorization: Bearer $TOKEN"
```

---

### **Step 6: Test in UI**

1. Scroll to "Test API" section
2. The code editor shows sample Python code
3. Replace it with:
```python
print('{"message": "Testing from UI"}')
```
4. Click "▶️ Run Code"
5. See output in the result box below
6. Should show: `{"message": "Testing from UI"}`

---

### **Step 7: Invoke Deployed API**

Now let's invoke the actual deployed API with its endpoint.

#### Get the Endpoint:
- Look at your API card in dashboard
- Copy the endpoint (e.g., `/execute/abc12345/weather-api`)

#### Invoke via cURL:

**Default request (Dubai):**
```bash
curl -X POST http://localhost:8080/execute/abc12345/weather-api \
  -H "Content-Type: application/json" \
  -d '{}'
```

**With input (London):**
```bash
curl -X POST http://localhost:8080/execute/abc12345/weather-api \
  -H "Content-Type: application/json" \
  -d '{
    "input": {
      "city": "London"
    }
  }'
```

**Expected Response:**
```json
{
  "output": "{\"city\": \"London\", \"temperature\": 15, \"condition\": \"Cloudy\", \"humidity\": 75, \"unit\": \"celsius\"}\n",
  "error": "",
  "status_code": 200,
  "duration_ms": 245,
  "exit_code": 0,
  "result": {
    "city": "London",
    "temperature": 15,
    "condition": "Cloudy",
    "humidity": 75,
    "unit": "celsius"
  }
}
```

---

### **Step 8: Test Different Cities**

```bash
# Tokyo
curl -X POST http://localhost:8080/execute/abc12345/weather-api \
  -H "Content-Type: application/json" \
  -d '{"input": {"city": "Tokyo"}}'

# New York
curl -X POST http://localhost:8080/execute/abc12345/weather-api \
  -H "Content-Type: application/json" \
  -d '{"input": {"city": "New York"}}'

# Unknown city (will use defaults)
curl -X POST http://localhost:8080/execute/abc12345/weather-api \
  -H "Content-Type: application/json" \
  -d '{"input": {"city": "Paris"}}'
```

---

### **Step 9: Test Edit Functionality**

#### UI Method:
1. Go back to Dashboard
2. Hover over your API card
3. Click the pencil icon (✏️) in top-right
4. Edit modal appears
5. Change description to: `Weather API with multiple cities`
6. Click "Save Changes"
7. Card updates with new description

---

### **Step 10: Test Stop Functionality**

#### UI Method:
1. Click on API card to open detail page
2. Click "⏸️ Stop" button
3. Confirm
4. Status changes to "stopped"
5. Deploy button appears again

#### Try to Invoke Stopped API:
```bash
curl -X POST http://localhost:8080/execute/abc12345/weather-api \
  -H "Content-Type: application/json" \
  -d '{}'
```

**Expected Response:**
```
API is not deployed. Current status: stopped
```

---

### **Step 11: Re-deploy**

1. Click "🚀 Deploy" again
2. API goes back to "deployed" status
3. Can be invoked again

---

## 🧪 Testing Other Runtimes

### **Node.js API Example**

Create `hello.js`:
```javascript
const input = process.env.INPUT ? JSON.parse(process.env.INPUT) : {};
const name = input.name || 'World';

const result = {
    message: `Hello, ${name}!`,
    timestamp: new Date().toISOString()
};

console.log(JSON.stringify(result));
```

1. Create API with runtime "Node.js"
2. Upload `hello.js`
3. Deploy
4. Invoke:
```bash
curl -X POST http://localhost:8080/execute/abc12345/hello-api \
  -H "Content-Type: application/json" \
  -d '{"input": {"name": "Ahmad"}}'
```

### **Go API Example**

Create `main.go`:
```go
package main

import (
	"encoding/json"
	"fmt"
)

type Result struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Version string `json:"version"`
}

func main() {
	result := Result{
		Status:  "success",
		Message: "Hello from Go!",
		Version: "1.0.0",
	}
	
	data, _ := json.Marshal(result)
	fmt.Println(string(data))
}
```

1. Create API with runtime "Go"
2. Upload `main.go`
3. Deploy
4. Invoke

---

## 🔍 Troubleshooting

### Check Service Logs
```bash
# All services
docker compose logs -f

# Specific service
docker compose logs -f api-gateway
docker compose logs -f executor
docker compose logs -f frontend
```

### API Returns Empty
- Check that code file exists in `uploads/` directory
- Verify code prints to stdout
- Check executor logs for errors

### CORS Errors
- Make sure executor service has CORS enabled
- Restart executor: `docker compose restart executor`

### Database Connection Issues
```bash
# Check postgres
docker compose logs postgres

# Restart if needed
docker compose restart postgres
```

### Services Not Starting
```bash
# Stop all
docker compose down

# Rebuild and start
docker compose up --build -d

# Check status
docker compose ps
```

---

## ✅ Success Criteria

After completing all steps, you should have:

- ✅ User account created and logged in
- ✅ API created via UI
- ✅ Code uploaded successfully
- ✅ API deployed with green status
- ✅ Tested code execution in UI
- ✅ Invoked API via cURL with correct response
- ✅ Edited API details
- ✅ Stopped and re-deployed API
- ✅ All operations work smoothly

---

## 📊 What to Test Next

1. **Multiple APIs**: Create several APIs and test them all
2. **Different Runtimes**: Test Python, Node.js, and Go
3. **Error Handling**: Upload invalid code and see error messages
4. **Marketplace**: Make API public and view in marketplace
5. **Concurrency**: Invoke API multiple times simultaneously
6. **Large Payloads**: Test with bigger input data
7. **Timeout**: Test with long-running code

---

## 🎉 You're All Set!

Your API platform is fully functional with:
- Complete authentication system
- API CRUD operations with edit modal
- Code upload with drag & drop
- Deployment system
- **NEW: Actual API invocation via unique endpoints**
- Beautiful, responsive UI
- Marketplace for public APIs
- Testing interface

**Next steps:** Consider implementing API keys and analytics for production readiness!
