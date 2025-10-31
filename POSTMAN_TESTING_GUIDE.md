# Postman Testing Guide

This guide will help you test the API Platform using Postman.

## Prerequisites

1. **Start the Backend**
   ```bash
   # Make sure database is running
   docker-compose up -d
   
   # Run migrations
   .\scripts\run-migrations.ps1
   
   # Start API Gateway
   cd backend/api-gateway
   go run main.go
   ```

2. **Install Postman**
   - Download from https://www.postman.com/downloads/

## Import the Collection

### Method 1: Import Files
1. Open Postman
2. Click **Import** (top left)
3. Drag and drop both files:
   - `API-Platform.postman_collection.json`
   - `API-Platform.postman_environment.json`
4. Click **Import**

### Method 2: Direct Import
1. In Postman, click **Import**
2. Select **File** tab
3. Browse to project root and select the JSON files
4. Click **Open**

## Set the Environment

1. In the top-right corner, click the environment dropdown
2. Select **"API Platform - Local"**
3. Verify base_url is set to `http://localhost:8080`

## Testing Workflow

### 1. Health Check ✅
**Verify API Gateway is running**

- Select: `Health Check`
- Click: **Send**
- Expected: `200 OK` with JSON response
```json
{
    "status": "healthy",
    "timestamp": "2025-10-31T...",
    "service": "api-gateway"
}
```

---

### 2. Create Developer Account 👨‍💻

**Register as a developer**

- Folder: `Authentication`
- Request: `Signup - Developer`
- Click: **Send**
- Expected: `200 OK`

**Response:**
```json
{
    "token": "eyJhbGc...",
    "user": {
        "id": "uuid",
        "email": "developer@example.com",
        "name": "John Developer",
        "role": "developer",
        "created_at": "2025-10-31T..."
    }
}
```

**✨ The token is automatically saved to environment variables!**

---

### 3. Verify Authentication 🔐

**Get current user details**

- Request: `Get Current User`
- Click: **Send**
- Expected: `200 OK` with your user details

---

### 4. Create APIs 🚀

**Create a Python API**

- Folder: `API Management`
- Request: `Create API - Python`
- Click: **Send**
- Expected: `200 OK`

**Response:**
```json
{
    "id": "api-uuid",
    "name": "weather-api",
    "description": "Simple weather API...",
    "version": "v1",
    "runtime": "python",
    "visibility": "public",
    "status": "pending",
    "endpoint": "/execute/abc12345/weather-api",
    "created_at": "2025-10-31T..."
}
```

**✨ The API ID is automatically saved!**

Try creating more APIs:
- `Create API - Node.js` (private)
- `Create API - Paid` (paid tier)

---

### 5. List Your APIs 📋

**View all your APIs**

- Request: `Get My APIs`
- Click: **Send**
- Expected: `200 OK` with array of your APIs

---

### 6. Get Specific API Details 🔍

**View details of one API**

- Request: `Get API by ID`
- Click: **Send**
- Expected: `200 OK` with API details

**Note:** Uses the `{{api_id}}` variable saved from create request

---

### 7. Upload Code 📦

**Upload a code file to your API**

- Request: `Upload Code`
- In **Body** tab, click **Select Files**
- Choose a test file:
  - `.py` for Python
  - `.js` for Node.js
  - `.go` for Go
- Click: **Send**
- Expected: `200 OK`

**Response:**
```json
{
    "message": "Code uploaded successfully",
    "path": "uploads/user-id/api-id/filename.py"
}
```

**Create a test file:**
```python
# test.py
def handler(request):
    return {"message": "Hello from API Platform!"}
```

---

### 8. Browse Marketplace 🛒

**View public APIs (no authentication needed)**

- Folder: `Marketplace`
- Request: `Get Public APIs`
- Click: **Send**
- Expected: `200 OK` with public APIs

**This works without authentication!**

---

### 9. Create Consumer Account 👤

**Register as a consumer**

- Request: `Signup - Consumer`
- Click: **Send**
- Expected: `200 OK`

**✨ Consumer token saved separately as `consumer_token`**

---

### 10. Delete API 🗑️

**Remove an API**

- Request: `Delete API`
- Click: **Send**
- Expected: `204 No Content`

---

## Variables Explained

The collection uses these variables (auto-saved):

| Variable | Purpose | Saved From |
|----------|---------|------------|
| `base_url` | API Gateway URL | Environment |
| `auth_token` | JWT for developer | Login/Signup |
| `consumer_token` | JWT for consumer | Consumer signup |
| `api_id` | Last created API ID | Create API |
| `user_id` | Current user ID | Login/Signup |

## Testing Scenarios

### Scenario 1: Complete Developer Flow
1. ✅ Health Check
2. ✅ Signup - Developer
3. ✅ Get Current User
4. ✅ Create API - Python (public)
5. ✅ Upload Code
6. ✅ Get My APIs
7. ✅ Get Public APIs (marketplace)

### Scenario 2: Consumer Flow
1. ✅ Health Check
2. ✅ Get Public APIs (no auth)
3. ✅ Signup - Consumer
4. ✅ Get Current User
5. ✅ Get Public API Details

### Scenario 3: Multiple APIs
1. ✅ Login - Developer
2. ✅ Create API - Python (public)
3. ✅ Create API - Node.js (private)
4. ✅ Create API - Paid
5. ✅ Get My APIs (see all 3)
6. ✅ Get Public APIs (see only public one)

## Expected Errors

### 401 Unauthorized
**Cause:** Missing or invalid token
**Solution:** Run `Login - Developer` first

### 403 Forbidden
**Cause:** Trying to access/modify someone else's API
**Solution:** Ensure you're the API owner

### 404 Not Found
**Cause:** API ID doesn't exist
**Solution:** Create an API first or check `{{api_id}}` variable

### 400 Bad Request
**Cause:** Invalid request body
**Solution:** Check request body format matches examples

## Tips & Tricks

### 1. Auto-save Tokens
The collection includes test scripts that automatically save tokens and IDs to variables.

### 2. Manual Token Override
To test with a different token:
1. Click environment dropdown
2. Click 👁️ icon
3. Edit `auth_token` value
4. Click **Save**

### 3. Test with Multiple Users
1. Use `Signup - Developer` with different email
2. Token is auto-saved
3. All subsequent requests use new account

### 4. View Environment Variables
- Click environment dropdown → 👁️ icon
- See all saved variables
- Copy token for debugging

### 5. Quick Re-run
Use **Cmd/Ctrl + Enter** to quickly re-send requests

## Troubleshooting

### "Connection refused"
- Backend not running
- Check: `go run main.go` in `backend/api-gateway`

### "Database connection failed"
- PostgreSQL not running
- Check: `docker-compose ps`
- Run: `docker-compose up -d`

### "Invalid token"
- Token expired (24h expiry)
- Re-login: Run `Login - Developer`

### "Environment not set"
- Select environment in top-right dropdown
- Choose: **API Platform - Local**

## Running Full Test Suite

1. Start backend: `go run main.go`
2. Open Postman
3. Select environment: **API Platform - Local**
4. Run requests in order:
   - Health Check
   - Signup - Developer
   - Get Current User
   - Create API - Python
   - Upload Code
   - Get My APIs
   - Get Public APIs
   - Delete API

**All should return 200/204 status codes!**

## Next Steps

After testing with Postman:
1. Test frontend UI at http://localhost:3000
2. Verify data in PostgreSQL:
   ```bash
   psql postgres://apiplatform:dev_password@localhost:5432/apiplatform
   SELECT * FROM users;
   SELECT * FROM apis;
   ```
3. Check uploaded files:
   ```bash
   dir uploads
   ```

---

**Happy Testing! 🚀**

Need help? Check [GETTING_STARTED.md](GETTING_STARTED.md) or [PROJECT_STATUS.md](PROJECT_STATUS.md)
