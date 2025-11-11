# Microsoft SSO Demo

## 🚀 Quick Test Microsoft SSO

### 1. Start Server
```bash
go run cmd/server/main.go
```

### 2. Test Microsoft SSO Flow

#### Step 1: Get Microsoft Login URL
```bash
curl -X GET http://localhost:8080/api/auth/microsoft/login
```

**Response:**
```json
{
  "login_url": "http://localhost:8000/login/oauth/authorize?application=myapp&client_id=ea525c19f6f75c2f8419&organization=myorg&provider=microsoft-provider&redirect_uri=http%3A%2F%2Flocalhost%3A8080%2Fapi%2Fauth%2Fcallback&response_type=code&scope=read&state=xyz123",
  "state": "xyz123",
  "message": "Redirect to this URL to login with Microsoft",
  "provider": "microsoft",
  "flow": "SSO via Casdoor → Microsoft → Casdoor → Backend"
}
```

#### Step 2: Manual Test (Browser)
1. Copy `login_url` from response
2. Open in browser
3. Should redirect to Microsoft login
4. After login → redirect back to `/api/auth/callback` with code

#### Step 3: Complete Flow
```bash
# After browser redirect, you'll get callback with code
# Example: http://localhost:8080/api/auth/callback?code=abc123&state=xyz123

# The callback endpoint will exchange code for JWT token
curl -X GET "http://localhost:8080/api/auth/callback?code=YOUR_CODE&state=xyz123"
```

**Response:**
```json
{
  "access_token": "eyJhbGciOiJSUzI1NiIs...",
  "token_type": "Bearer",
  "expires_in": 3600,
  "message": "Login successful. Use the access_token in Authorization header."
}
```

#### Step 4: Use JWT Token
```bash
# Use the access_token from step 3
TOKEN="eyJhbGciOiJSUzI1NiIs..."

# Test protected endpoint
curl -X GET http://localhost:8080/api/auth/me \
  -H "Authorization: Bearer $TOKEN"
```

## 🔄 Complete SSO Flow

```
1. Frontend → GET /api/auth/microsoft/login
2. Frontend → Redirect user to login_url
3. User → Login with Microsoft account
4. Microsoft → Redirect to Casdoor with auth code
5. Casdoor → Process Microsoft user info
6. Casdoor → Redirect to /api/auth/callback?code=xxx
7. Backend → Exchange code for Casdoor JWT
8. Backend → Return JWT to frontend
9. Frontend → Use JWT for API calls
```

## 🎯 Expected User Flow

### New Microsoft User:
1. Login with Microsoft → Casdoor creates new user
2. User gets JWT with Microsoft email/name
3. Can access protected endpoints based on Casbin policies

### Existing User:
1. Login with Microsoft → Casdoor finds existing user
2. User gets JWT with existing permissions
3. Same authorization rules apply

## 🔧 Troubleshooting

### Common Issues:

1. **"failed to generate Microsoft login URL"**
   - Check CASDOOR_CLIENT_ID in .env
   - Verify Casdoor application is configured

2. **Microsoft login fails**
   - Check Azure App Registration redirect URI
   - Verify microsoft-provider is enabled in Casdoor

3. **Token exchange fails**
   - Check CASDOOR_CLIENT_SECRET
   - Verify callback URL matches Azure/Casdoor config

### Debug Commands:
```bash
# Check server health
curl http://localhost:8080/health

# Check Casdoor connection
curl http://localhost:8000/api/get-global-providers

# Check current config
curl http://localhost:8080/api/auth/oauth/login
```

## 📋 Prerequisites Checklist

- ✅ Azure App Registration created
- ✅ Casdoor microsoft-provider configured  
- ✅ Redirect URIs match in all systems
- ✅ Environment variables set
- ✅ Server running on correct port

Ready to test Microsoft SSO! 🎉