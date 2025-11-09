# Luồng Tích Hợp Casdoor với Hệ Thống Hiện Tại

## Tình huống hiện tại
- Hệ thống có 1 account admin/admin fix cứng trong code
- Chưa có quản lý user
- Muốn dùng Casdoor làm auth server

## Các Flow Tích Hợp

### Option 1: OAuth2/OIDC Flow (Flow hiện tại - Phức tạp hơn)

**Luồng:**
```
1. User → Backend: GET /api/auth/login
2. Backend → User: Redirect URL đến Casdoor
3. User → Casdoor: Đăng nhập
4. Casdoor → Backend: Redirect về /api/auth/callback?code=xxx
5. Backend → Casdoor: Exchange code lấy access_token
6. Backend → User: Trả về access_token
7. User → Backend: Gửi request với Bearer token
8. Backend: Verify JWT token từ Casdoor
```

**Ưu điểm:**
- Chuẩn OAuth2/OIDC
- Phù hợp cho web app với frontend riêng
- User quản lý session trên Casdoor

**Nhược điểm:**
- Phức tạp hơn
- Cần redirect flow
- Không phù hợp cho API-only backend

---

### Option 2: Direct Login API (Đơn giản hơn - Khuyến nghị)

**Luồng:**
```
1. User → Backend: POST /api/auth/login {username, password}
2. Backend → Casdoor: POST /api/login {application, username, password, type: "token"}
3. Casdoor → Backend: Trả về JWT token
4. Backend → User: Trả về JWT token
5. User → Backend: Gửi request với Bearer token
6. Backend: Verify JWT token từ Casdoor
```

**Ưu điểm:**
- Đơn giản, giống flow hiện tại (username/password)
- Không cần redirect
- Phù hợp cho API-only backend
- Dễ migrate từ code fix cứng

**Nhược điểm:**
- Backend phải xử lý password (nhưng chỉ forward đến Casdoor)
- Không có OAuth flow

---

### Option 3: Hybrid (Kết hợp cả 2)

**Luồng:**
- Web app: Dùng OAuth2 flow
- API/Mobile: Dùng Direct Login API

---

## Khuyến nghị: Option 2 (Direct Login API)

Vì hệ thống hiện tại:
- Đã có login với username/password
- Chỉ cần thay thế logic verify admin/admin
- Không cần OAuth flow phức tạp

## Implementation Plan

### Bước 1: Tạo Casdoor Application
1. Vào Casdoor: http://localhost:8000
2. Tạo Application:
   - Name: `your-app-name`
   - Organization: `built-in` (hoặc tạo org mới)
   - Enable Password: `true`
   - Enable Code: `true` (nếu cần OAuth sau)
3. Copy Client ID và Client Secret

### Bước 2: Migrate Users
1. Tạo user admin trong Casdoor:
   - Username: `admin`
   - Password: `admin` (hoặc password mới)
   - Organization: `built-in`
   - Is Admin: `true`

### Bước 3: Update Backend Code

**Thay thế:**
```go
// Code cũ
if username == "admin" && password == "admin" {
    // Login success
}
```

**Bằng:**
```go
// Code mới
token, err := casdoorLogin(username, password)
if err != nil {
    // Login failed
}
// Verify token và lấy user info
```

### Bước 4: Update Middleware

**Thay thế:**
```go
// Code cũ
if token == "admin-secret-token" {
    // Allow access
}
```

**Bằng:**
```go
// Code mới
claims, err := VerifyCasdoorToken(token)
if err != nil {
    // Unauthorized
}
// Use claims.UserID, claims.IsAdmin, etc.
```

## Code Structure

```
internal/
├── auth/
│   ├── casdoor.go          # Casdoor JWT verification (đã có)
│   ├── middleware.go       # Auth middleware (đã có)
│   └── login.go            # Direct login API (mới)
├── handler/
│   ├── auth.go             # Auth handlers
│   └── user.go             # User handlers
```

## API Endpoints

### Direct Login (Khuyến nghị)
```
POST /api/auth/login
Body: {
  "username": "admin",
  "password": "admin"
}
Response: {
  "access_token": "eyJhbGci...",
  "user": {
    "id": "...",
    "name": "admin",
    "is_admin": true
  }
}
```

### OAuth Flow (Nếu cần)
```
GET /api/auth/login          # Get OAuth login URL
GET /api/auth/callback       # OAuth callback
```

### Protected Routes
```
GET /api/protected
Authorization: Bearer <token>
```

## Migration Steps

1. **Setup Casdoor** (đã có)
2. **Tạo Application trong Casdoor**
3. **Tạo user admin trong Casdoor**
4. **Update login handler** - Thay thế logic verify admin/admin
5. **Update auth middleware** - Thay thế logic verify token
6. **Test và deploy**

## So sánh với Flow hiện tại

| Aspect | Flow hiện tại (OAuth) | Direct Login (Khuyến nghị) |
|--------|----------------------|---------------------------|
| Complexity | Cao | Thấp |
| User Experience | Redirect flow | Direct API call |
| Phù hợp cho | Web app với frontend | API-only backend |
| Migration effort | Cao | Thấp |
| Security | Cao (OAuth2) | Cao (JWT) |

## Kết luận

**Khuyến nghị: Dùng Direct Login API (Option 2)**

Vì:
- Đơn giản hơn, dễ migrate
- Phù hợp với hệ thống hiện tại (username/password)
- Vẫn secure với JWT verification
- Có thể thêm OAuth flow sau nếu cần

