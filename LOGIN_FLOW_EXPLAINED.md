# Giải Thích Flow Logic API Login

## Endpoint: `POST /api/auth/login`

### Request
```json
{
  "username": "admin",
  "password": "admin"
}
```

### Response
```json
{
  "access_token": "eyJhbGci...",
  "token_type": "Bearer",
  "user": {
    "id": "...",
    "name": "admin",
    "display_name": "Admin",
    "email": "admin@example.com",
    "is_admin": true
  },
  "message": "Login successful. Use the access_token in Authorization header."
}
```

---

## Flow Logic Chi Tiết

### Bước 1: Handler nhận request
**File**: `internal/handler/auth.go` - `DirectLogin()`

```go
func (h *AuthHandler) DirectLogin(c echo.Context) error {
    // 1. Parse request body
    var req LoginRequest
    c.Bind(&req)  // {"username": "admin", "password": "admin"}
    
    // 2. Validate input
    if req.Username == "" || req.Password == "" {
        return error("username and password are required")
    }
    
    // 3. Gọi Casdoor để login
    token, err := auth.DirectLogin(req.Username, req.Password)
    
    // 4. Verify token để lấy user info
    claims, err := auth.VerifyToken(token)
    
    // 5. Trả về token + user info
    return c.JSON(200, {
        "access_token": token,
        "user": {...}
    })
}
```

**Logic:**
- Nhận username/password từ client
- Validate input
- Forward request đến Casdoor
- Verify token nhận được
- Trả về token + user info

---

### Bước 2: DirectLogin với Casdoor
**File**: `internal/auth/login.go` - `DirectLogin()`

```go
func DirectLogin(username, password string) (string, error) {
    // 1. Lấy config
    cfg := config.GetConfig()
    
    // 2. Tạo request body
    loginReq := {
        "application": "myapp",      // Từ config
        "username": "admin",         // Từ input
        "password": "admin",         // Từ input
        "type": "token",             // Request JWT token
        "organization": "myorg"      // Từ config
    }
    
    // 3. Gửi POST request đến Casdoor
    POST http://localhost:8000/api/login
    Body: JSON(loginReq)
    
    // 4. Parse response
    Response: {
        "status": "ok",
        "data": "eyJhbGci..."  // JWT token
    }
    
    // 5. Trả về token
    return loginResp.Data
}
```

**Logic:**
- Tạo request body với application, username, password, organization
- Gửi POST đến Casdoor `/api/login`
- Casdoor verify username/password
- Casdoor tạo JWT token và trả về
- Trả về JWT token

**Request đến Casdoor:**
```
POST http://localhost:8000/api/login
Content-Type: application/json

{
  "application": "myapp",
  "username": "admin",
  "password": "admin",
  "type": "token",
  "organization": "myorg"
}
```

**Response từ Casdoor:**
```json
{
  "status": "ok",
  "msg": "",
  "data": "eyJhbGciOiJSUzI1NiIs...",  // JWT token
  "data2": "...",
  "data3": false
}
```

---

### Bước 3: Verify Token
**File**: `internal/auth/casdoor.go` - `VerifyToken()`

```go
func VerifyToken(tokenString string) (*CasdoorClaims, error) {
    // 1. Lấy public key từ Casdoor
    publicKey, err := GetPublicKey()
    
    // 2. Parse và verify JWT token
    token, err := jwt.ParseWithClaims(tokenString, &CasdoorClaims{}, ...)
    
    // 3. Extract claims từ token
    claims := token.Claims.(*CasdoorClaims)
    
    // 4. Trả về claims
    return claims
}
```

**Logic:**
- Lấy public key từ Casdoor (để verify signature)
- Parse JWT token
- Verify signature bằng public key
- Extract claims (user info)
- Trả về claims

**Claims trong JWT token:**
```json
{
  "owner": "myorg",
  "name": "admin",
  "displayName": "Admin",
  "email": "admin@example.com",
  "id": "...",
  "isAdmin": true,
  "roles": [...],
  "sub": "...",  // Subject (user ID)
  "exp": 1234567890,  // Expiration time
  ...
}
```

---

### Bước 4: Get Public Key
**File**: `internal/auth/casdoor.go` - `GetPublicKey()`

```go
func GetPublicKey() (*rsa.PublicKey, error) {
    // 1. Check cache
    if certCache != nil {
        return certCache  // Đã có, return luôn
    }
    
    // 2. Fetch JWKS từ Casdoor
    GET http://localhost:8000/.well-known/jwks
    
    // 3. Parse JWKS JSON
    Response: {
        "keys": [{
            "kty": "RSA",
            "kid": "cert-built-in",
            "x5c": ["MIIE3TCCAsWgAwIBAgIDAeJAMA0GCSqGSIb3DQEBCwUAMCgxDjAMBgNVBAoTBWFkbWluMRYwFAYDVQQDEw1jZXJ0LWJ1aWx0LWluMB4XDTI1MTEwODA4NTEyMloXDTQ1MTEwODA4NTEyMlowKDEOMAwGA1UEChMFYWRtaW4xFjAUBgNVBAMTDWNlcnQtYnVpbHQtaW4wggIiMA0GCSqGSIb3DQEBAQUAA4ICDwAwggIKAoICAQDDmGxOQHvM0VK+W7leY5i2dJtiafJ8plgJO4oCMdXslbb0F+g0qtlrF84oJIPhNeZ2EjFq8l2Eu3elTQESAFhnN90d9g8JF2CVv/cLG/25MQekAOWOeaCq4rYLHO+QwFCSxbI1PoN8WJLat1GdYLYw7lytcLDetAoK+Bb6pqrqqfVeuqmxThB8l8UeDyzVgOPbwPZO7obqOA2PkMTnLP7gqxB440VoE4w1pbzre55WamX6tW4j46rvSHggkPQvj9E/yVyjeCteYPfm9NoDLsB1e5Z+0PUoHaRrZxvsHxmWzaDMM/mnaz9uwQd7GgCzdHDZyLjrD0VyGF/aYZ/dRmeMEPuBqxihgHpIKMpYu4n8wxwSCJQSS3acfAyrzOHxJqTUXNMTPatPiuAZfDoh7HQ9yq6pjO4r1Ck4cehWLWRoX0ufYSgiuPtCN4DKC1hxtVf+sIrIGQBgTPoMnTJog1vIb6dstWYqJmMaSnbmnn2A+8qiDxLMkNPF0jZqCUTZZrfMvM3FAI+Koal5zqOg+1YDHDJRHd8491hg7YZwJg2t7VoaNrPWcRJA6sionVTr6kTbw1m31bwkTNZQCmlvbpvIUAoJM8Ag+stjJmin1SFrlm31VZ0p3j8SUvOuvzJ0JnhicynysuisWOs93ggIIPf0qbwEE/O9Zp+9tlT8h5hBuQIDAQABoxAwDjAMBgNVHRMBAf8EAjAAMA0GCSqGSIb3DQEBCwUAA4ICAQABgTteYo+vgjeOHv013RYYnMuQROpWfL7O0zA5ZjGPMMKK23NOA8ummxeR8XA1+8Tj+OOtLNJItrbi0P/WHx64+IyPMV9TrpT0irBqGxnOnxf1KXh/L3PVVxN3sag0sgeSUwGLpYBb2LDJ+o6WylYv53Ima0vJ+gDi4FrFAnL5WyYeqtaEvygaUpF6XkHtJNu1ow8jMd1DjM3e2AWKzVI0rr7jfK/GeWbQN6it/GU6SzopLeciKtSabH2SQAYU+yrlHdOGB6mf9HZvXOVPlt78AnhMO2N0GxvJ2EakUBEtBV7QXGXTdYEwiAfmEVzxCAogv+iJBgOK6tpJ2ZZdMDGCxwbhJrdmzLXPYLgq1YSHPlakQXnlaB2azL5bzuHzpqIpY32apc2mNu7zAH/ioGsHmFRTdAbAtFeyrdL5B7jhVXzNMFv4gn0qFbOIIpPXVQojD1/LeyVjaScKoi4G0PIBdRvIKRadmTQ3y+fthRqcQ+cpx33LW+5UPSK0JQgMlsDeAwvLwjp+yX1b/OHCBCaIs+HXmgHC0XHFSJ5c6V/QFaY3vwcXlkbtxXPlX0ZwksrP83w48UhzekOi4W9LhT/WPGeh/qk3sn3PxmCOiNKZtXXeOY2i0Vsk5HY65TwsJNygp1G8wlhkgdfRS0KskPaMigq3cPTZfdmyUG8gHFsI2g=="]
        }]
    }
    
    // 4. Decode base64 certificate từ x5c[0]
    certDER = base64.Decode(x5c[0])
    
    // 5. Parse certificate
    cert = x509.ParseCertificate(certDER)
    
    // 6. Extract public key
    pubKey = cert.PublicKey.(*rsa.PublicKey)
    
    // 7. Cache và return
    certCache = pubKey
    return pubKey
}
```

**Logic:**
- Check cache trước (nếu đã có thì return luôn)
- Fetch JWKS từ `/.well-known/jwks` (public endpoint, không cần auth)
- Parse JWKS JSON
- Lấy certificate từ `x5c[0]` (base64 encoded)
- Decode base64 → Parse certificate → Extract public key
- Cache public key để dùng lại

---

## Tổng Kết Flow

```
┌─────────┐
│ Client  │
└────┬────┘
     │ 1. POST /api/auth/login
     │    {"username": "admin", "password": "admin"}
     ▼
┌─────────────────┐
│ Backend Handler │
│ DirectLogin()   │
└────┬────────────┘
     │ 2. Validate input
     │
     │ 3. Call auth.DirectLogin()
     ▼
┌─────────────────┐
│ auth.DirectLogin│
└────┬────────────┘
     │ 4. POST http://localhost:8000/api/login
     │    {"application": "myapp", "username": "admin", ...}
     ▼
┌─────────┐
│ Casdoor │
└────┬────┘
     │ 5. Verify username/password
     │ 6. Generate JWT token
     │ 7. Return token
     ▼
┌─────────────────┐
│ Backend Handler │
└────┬────────────┘
     │ 8. Call auth.VerifyToken()
     ▼
┌─────────────────┐
│ auth.VerifyToken│
└────┬────────────┘
     │ 9. Call GetPublicKey()
     ▼
┌─────────────────┐
│ GetPublicKey()  │
└────┬────────────┘
     │ 10. GET /.well-known/jwks
     │ 11. Parse JWKS → Extract public key
     ▼
┌─────────────────┐
│ auth.VerifyToken│
└────┬────────────┘
     │ 12. Verify JWT signature
     │ 13. Extract claims
     │ 14. Return claims
     ▼
┌─────────────────┐
│ Backend Handler │
└────┬────────────┘
     │ 15. Return token + user info
     ▼
┌─────────┐
│ Client │
└────────┘
```

## Các Bước Chi Tiết

### 1. Client gửi request
```bash
POST /api/auth/login
Body: {"username": "admin", "password": "admin"}
```

### 2. Backend validate
- Check username và password không rỗng

### 3. Backend gọi Casdoor
```bash
POST http://localhost:8000/api/login
Body: {
  "application": "myapp",
  "username": "admin",
  "password": "admin",
  "type": "token",
  "organization": "myorg"
}
```

### 4. Casdoor verify
- Tìm user `admin` trong organization `myorg`
- Verify password
- Tạo JWT token với user info

### 5. Casdoor trả về token
```json
{
  "status": "ok",
  "data": "eyJhbGci..."  // JWT token
}
```

### 6. Backend verify token
- Lấy public key từ JWKS endpoint
- Verify JWT signature
- Extract claims (user info)

### 7. Backend trả về cho client
```json
{
  "access_token": "eyJhbGci...",
  "user": {
    "id": "...",
    "name": "admin",
    "is_admin": true
  }
}
```

## Lưu Ý

1. **Public Key Caching**: Public key được cache sau lần fetch đầu tiên
2. **Token Verification**: Token được verify bằng RSA public key từ Casdoor
3. **User Info**: User info được extract từ JWT claims, không cần query Casdoor
4. **Security**: Password không được lưu trong backend, chỉ forward đến Casdoor

