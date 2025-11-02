# API Platform - User Guide

## 🎯 Overview
This is a MENA-first API hosting and marketplace platform where developers can deploy, test, and monetize their APIs. The platform supports Python, Node.js, and Go runtimes with Docker-based sandboxed execution.

---

## 🚀 Getting Started

### Access the Platform
Open your browser and go to: **http://localhost:3000**

---

## 📍 Navigation Flow

### 1️⃣ **Home Page** (`/`)
**First Stop for Everyone**

#### When Logged Out:
- **Hero Section**: Platform introduction with call-to-action buttons
- **Navigation Options**:
  - "Sign in" → Login page
  - "Get Started" → Signup page
  - "Explore Marketplace" → Browse public APIs
- **Features**: View platform capabilities (Instant Deploy, Analytics, Billing, Security)

#### When Logged In:
- **Navigation Changes To**:
  - "Dashboard" → Your API management
  - "Hi, [Your Name]" → Welcome message
  - "Logout" → Sign out
- **Hero CTAs Change To**:
  - "Go to Dashboard" → Access your APIs
  - "Browse APIs" → Marketplace

---

### 2️⃣ **Sign Up** (`/signup`)
**Create Your Account**

**Required Fields**:
- Full Name
- Email
- Password (min 8 characters)
- Role Selection:
  - "Host and monetize my APIs" → Developer
  - "Use APIs from marketplace" → Consumer

**After Signup**: Automatically redirected to Dashboard

---

### 3️⃣ **Login** (`/login`)
**Access Your Account**

**Required Fields**:
- Email
- Password

**After Login**: Redirected to Dashboard

**Test Credentials** (if you want to test with existing account):
- Email: `demo@apiplatform.com`
- Password: `SecurePass123!`

---

### 4️⃣ **Dashboard** (`/dashboard`)
**Your API Management Hub** (Must be logged in)

#### Top Navigation:
- **Logo** → Click to return home
- **Marketplace** → Browse public APIs
- **Hi, [Name]** → Shows your name
- **Logout** → Sign out

#### Main Content:

##### When You Have No APIs:
- Empty state with rocket emoji 🚀
- Message: "No APIs yet"
- Big blue button: "Create Your First API"

##### When You Have APIs:
- **Grid of API Cards** (clickable)
- Each card shows:
  - **API Name** (bold title)
  - **Description**
  - **Status Badge** (color-coded):
    - 🟢 Green = Deployed
    - 🟡 Yellow = Pending
    - 🔴 Red = Failed/Stopped
  - **Runtime**: Python/Node.js/Go
  - **Version**: v1, v2, etc.
  - **Visibility**: Private/Public/Paid
  - **Endpoint**: Your API URL

##### Create New API:
Click **"+ Create API"** button (top right)

**Modal Form Opens**:
- **Name**: Your API name (required)
- **Description**: What your API does
- **Runtime**: Choose Python, Node.js, or Go
- **Visibility**: 
  - Private (only you)
  - Public (free for all)
  - Paid (monetized)

**After Creation**: 
- API appears in your dashboard with "pending" status
- Click on the API card to open detail page

---

### 5️⃣ **API Detail Page** (`/dashboard/api/[id]`)
**Upload, Test, and Deploy Your API** (Must be logged in)

This page has **3 main sections**:

---

#### 📤 **Section 1: Upload Code**

**Drag & Drop Zone** with file icon 📁

**Two Ways to Upload**:
1. **Drag & Drop**: Drag a code file directly into the dashed box
2. **Browse Files**: Click blue "Browse Files" button

**Supported File Types**:
- `.py` (Python)
- `.js` (Node.js)
- `.ts` (TypeScript)
- `.go` (Go)

**Upload Process**:
1. Select or drop your file
2. "Uploading..." message appears
3. Success: ✅ "[filename] uploaded successfully!" (green box)
4. Error: ❌ "Failed: [error message]" (red box)
5. API status automatically updates

**File Size Limit**: 10 MB max

---

#### 🧪 **Section 2: Test API**

**Live Code Testing Interface**

**Code Editor**:
- Large textarea with your code
- Pre-filled with sample code based on runtime:
  - Python: `print("Hello from API!")`
  - Node.js: `console.log("Hello from API!");`
  - Go: `package main...`

**How to Test**:
1. Write or paste your code in the textarea
2. Click **"▶️ Run Code"** button (blue)
3. Wait for execution (button shows "Running...")
4. Results appear below

**Result Display**:
- **Success** (gray box):
  - Shows output
  - Duration in milliseconds
  - Exit code
- **Error** (red box):
  - Shows error message
  - Stack trace if available

**Live Testing**: Tests run against the actual executor service (port 8081)

---

#### 🚀 **Section 3: Deploy**

**Header Section** (top of page):
- API name and description
- Status badges (runtime, version, status)
- **"🚀 Deploy" button** (green, top right)

**Deploy Process** (coming soon):
- Click "🚀 Deploy"
- Confirmation dialog
- API becomes live and accessible via endpoint

---

## 🎨 Visual Flow Diagram

```
┌──────────────────────────────────────────────────────────────┐
│                      HOME PAGE (/)                            │
│  [Logo: AP] API Platform           [Sign in] [Get Started]   │
├──────────────────────────────────────────────────────────────┤
│                                                                │
│         🚀 Deploy & Monetize Your APIs Instantly              │
│                                                                │
│     [Start Building Free] [Explore Marketplace]               │
│                                                                │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐          │
│  │👨‍💻 Developers│  │🛒 Consumers │  │🌍 MENA-First│          │
│  └─────────────┘  └─────────────┘  └─────────────┘          │
└──────────────────────────────────────────────────────────────┘
                    │                     │
        ┌───────────┴──────────┐         │
        ↓                      ↓         ↓
┌──────────────┐      ┌──────────────┐
│  SIGNUP (/)  │      │  LOGIN (/)   │
│              │      │              │
│ Full Name    │      │ Email        │
│ Email        │      │ Password     │
│ Password     │      │              │
│ Role         │      │ [Sign in]    │
│              │      └──────────────┘
│ [Create]     │              │
└──────────────┘              │
        │                     │
        └──────────┬──────────┘
                   ↓
┌──────────────────────────────────────────────────────────────┐
│                   DASHBOARD (/dashboard)                      │
│  [Logo] API Platform    [Marketplace] Hi, John [Logout]      │
├──────────────────────────────────────────────────────────────┤
│  My APIs                              [+ Create API]          │
│                                                                │
│  ┌──────────────────────────────┐                            │
│  │ 📦 Weather API         🟡pending                          │
│  │ Get real-time weather data                                │
│  │ Runtime: python • v1 • private                            │
│  │ /execute/abc123/weather                                   │
│  └──────────────────────────────┘  ← Click to open           │
│                                                                │
│  ┌──────────────────────────────┐                            │
│  │ 💱 Currency API        🟢deployed                         │
│  │ Convert currencies                                        │
│  │ Runtime: nodejs • v1 • public                             │
│  │ /execute/abc123/currency                                  │
│  └──────────────────────────────┘                            │
└──────────────────────────────────────────────────────────────┘
                   │ (click on API card)
                   ↓
┌──────────────────────────────────────────────────────────────┐
│            API DETAIL (/dashboard/api/xyz123)                 │
│  [← Back]                                      [Logout]       │
├──────────────────────────────────────────────────────────────┤
│  Weather API                              [🚀 Deploy]         │
│  Get real-time weather data                                   │
│  [python] [v1] [pending]                                     │
├──────────────────────────────────────────────────────────────┤
│  📤 UPLOAD CODE                                               │
│  ┌────────────────────────────────────────────────────┐     │
│  │              📁                                     │     │
│  │   Drag & drop your code file here                  │     │
│  │                  or                                 │     │
│  │            [Browse Files]                           │     │
│  │   Supported: .py, .js, .go, .ts (max 10MB)        │     │
│  └────────────────────────────────────────────────────┘     │
│  ✅ weather.py uploaded successfully!                        │
├──────────────────────────────────────────────────────────────┤
│  🧪 TEST API                                                  │
│  ┌────────────────────────────────────────────────────┐     │
│  │ import requests                                     │     │
│  │                                                      │     │
│  │ def get_weather(city):                              │     │
│  │     # your code here...                             │     │
│  │                                                      │     │
│  └────────────────────────────────────────────────────┘     │
│  [▶️ Run Code]                                               │
│                                                                │
│  Result:                                                      │
│  ┌────────────────────────────────────────────────────┐     │
│  │ Output: {"temp": 22, "city": "Dubai"}              │     │
│  │ Duration: 313ms | Exit Code: 0                     │     │
│  └────────────────────────────────────────────────────┘     │
└──────────────────────────────────────────────────────────────┘
```

---

## 🎯 Typical User Journey

### For Developers (Hosting APIs):

1. **Home** → Click "Get Started"
2. **Signup** → Fill form, select "Host and monetize my APIs"
3. **Dashboard** → Automatically redirected
4. **Create API** → Click "+ Create API"
5. **Fill Modal**:
   - Name: "Weather API"
   - Runtime: Python
   - Visibility: Public
6. **API Created** → Click on the new API card
7. **Upload Code** → Drag your `weather.py` file
8. **Test Code** → Write test code, click "Run Code"
9. **Deploy** → Click "🚀 Deploy" when ready
10. **API Live!** → Share your endpoint

### For Consumers (Using APIs):

1. **Home** → Click "Explore Marketplace"
2. **Browse APIs** → See available public APIs
3. **Signup/Login** → Create account
4. **Get API Key** → From dashboard
5. **Integrate** → Use the provided endpoint
6. **Monitor Usage** → Track in dashboard

---

## 📱 Responsive Design

The platform is **fully responsive**:

- **Mobile** (< 640px): 
  - Stacked layouts
  - Full-width buttons
  - Hamburger menu (simplified nav)
  
- **Tablet** (640-1024px):
  - 2-column grids
  - Side-by-side buttons
  
- **Desktop** (> 1024px):
  - 3-column grids
  - Full navigation bar
  - Optimal spacing

**Test on any device** - the UI adapts automatically!

---

## 🔐 Authentication States

The platform is **context-aware**:

- **Logged Out**: Shows signup/login prompts
- **Logged In**: Shows dashboard access and user name
- **Protected Routes**: Dashboard pages require login
- **Auto-redirect**: Login required pages redirect to `/login`

---

## ✅ Visual Design Highlights

### Colors & States:
- **Blue/Indigo Gradient**: Primary actions (Create, Deploy, Login)
- **Green**: Success, deployed status
- **Yellow**: Pending, warning states
- **Red**: Errors, delete actions, logout
- **Gray**: Neutral, secondary info

### Interactive Elements:
- **Hover Effects**: Cards lift with shadow, borders change color
- **Loading States**: Spinners, disabled buttons with opacity
- **Drag States**: Drop zone highlights blue when dragging
- **Status Badges**: Color-coded with rounded pills

### Typography:
- **Bold Gradients**: Hero headlines
- **Clear Hierarchy**: H1 → H2 → H3 sizing
- **Monospace**: Code snippets and endpoints
- **Sans-serif**: Clean, modern UI text

---

## 🐛 Testing the Flow

**Quick Test Path**:
1. Start services: `docker compose up`
2. Open: http://localhost:3000
3. Click "Get Started"
4. Sign up with test data
5. Create an API named "Test API"
6. Upload a simple Python file: `print("Hello")`
7. Test it in the code editor
8. See the output!

---

## 🔧 Technical Stack

**Frontend** (Port 3000):
- Next.js 14 (React)
- TypeScript
- Tailwind CSS
- Context API for auth

**Backend** (Port 8080):
- Go API Gateway
- JWT authentication
- PostgreSQL database

**Executor** (Port 8081):
- Go service
- Docker-based sandboxing
- Multi-runtime support

**All running in Docker Compose** for easy deployment!

---

## 💡 Tips

1. **Upload Before Testing**: Upload your code file first, then use the test interface
2. **Status Colors Matter**: Green = good to go, Yellow = needs code, Red = check errors
3. **Click Cards**: API cards are fully clickable - click anywhere on the card
4. **Drag Files**: Drag & drop is faster than browsing
5. **Auto-Reload**: After upload, the page refreshes API data automatically

---

## 🎉 You're Ready!

Visit **http://localhost:3000** and start building your API marketplace!

**Need Help?** Check the code, everything is commented and organized.
