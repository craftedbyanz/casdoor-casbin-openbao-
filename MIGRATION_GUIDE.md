# Migration Guide: Từ Admin/Admin Fix Cứng sang Casdoor

## Tình huống hiện tại

Hệ thống hiện tại có:
- Login với username/password fix cứng: `admin/admin`
- Không có quản lý user
- Cần migrate sang Casdoor

## Luồng Tích Hợp - Khuyến nghị: Direct Login API

### Flow mới (Direct Login)

```
1. Client → Backend: POST /api/auth/login
   Body: {"username": "admin", "password": "admin"}

2. Backend → Casdoor: POST /api/login
   Body: {application, username, password, type: "token"}

3. Casdoor → Backend: Trả về JWT token

4. Backend → Client: Trả về access_token + user info

5. Client → Backend: Gửi request với Bearer token
   Header: Authorization: Bearer <token>

6. Backend: Verify JWT token từ Casdoor
   - Lấy user info từ token
   - Check is_admin, roles, etc.
```

### So sánh với OAuth Flow

| Aspect | Direct Login (Khuyến nghị) | OAuth Flow |
|--------|---------------------------|------------|
| **Complexity** | Thấp | Cao |
| **User Experience** | Direct API call | Redirect flow |
| **Phù hợp cho** | API-only backend | Web app với frontend |
| **Migration effort** | Thấp (giống flow cũ) | Cao |
| **Security** | Cao (JWT) | Cao (OAuth2) |

## Migration Steps

### Bước 1: Setup Casdoor

```bash
docker-compose up -d
```

### Bước 2: Tạo Application trong Casdoor

1. Vào http://localhost:8000
2. Đăng nhập admin
3. Vào **Applications** → **Add**
4. Điền:
   - **Name**: `your-app-name`
   - **Organization**: `built-in`
   - **Enable Password**: `true`
   - **Enable Code**: `true` (nếu cần OAuth sau)
5. Save và copy **Client ID** và **Client Secret**

### Bước 3: Tạo User trong Casdoor

1. Vào **Users** → **Add**
2. Điền:
   - **Name**: `admin`
   - **Password**: `admin` (hoặc password mới)
   - **Organization**: `built-in`
   - **Is Admin**: `true`
3. Save

### Bước 4: Cấu hình Backend

Tạo file `.env`:

```env
CASDOOR_ENDPOINT=http://localhost:8000
CASDOOR_CLIENT_ID=your_client_id
CASDOOR_CLIENT_SECRET=your_client_secret
CASDOOR_ORGANIZATION=built-in
CASDOOR_APPLICATION=your-app-name
```

### Bước 5: Update Code

#### Trước (Code cũ):

```go
// Login handler
func Login(c echo.Context) error {
    var req struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }
    c.Bind(&req)
    
    // Fix cứng
    if req.Username == "admin" && req.Password == "admin" {
        return c.JSON(200, map[string]interface{}{
            "token": "admin-secret-token",
            "user": map[string]interface{}{
                "username": "admin",
                "is_admin": true,
            },
        })
    }
    
    return c.JSON(401, map[string]interface{}{
        "error": "Invalid credentials",
    })
}

// Middleware
func AuthMiddleware() echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            token := c.Request().Header.Get("Authorization")
            if token != "admin-secret-token" {
                return echo.NewHTTPError(401, "Unauthorized")
            }
            c.Set("user", map[string]interface{}{
                "username": "admin",
                "is_admin": true,
            })
            return next(c)
        }
    }
}
```

#### Sau (Code mới):

```go
// Login handler - Dùng Direct Login API
func Login(c echo.Context) error {
    var req struct {
        Username string `json:"username"`
        Password string `json:"password"`
    }
    c.Bind(&req)
    
    // Login với Casdoor
    token, err := auth.DirectLogin(req.Username, req.Password)
    if err != nil {
        return echo.NewHTTPError(401, "Login failed: "+err.Error())
    }
    
    // Verify token để lấy user info
    claims, err := auth.VerifyToken(token)
    if err != nil {
        return echo.NewHTTPError(500, "Failed to verify token")
    }
    
    return c.JSON(200, map[string]interface{}{
        "access_token": token,
        "user": map[string]interface{}{
            "id": claims.GetUserID(),
            "name": claims.Name,
            "is_admin": claims.IsAdmin,
        },
    })
}

// Middleware - Dùng JWT verification
func AuthMiddleware() echo.MiddlewareFunc {
    return auth.AuthMiddleware() // Đã có sẵn
}
```

### Bước 6: Test

```bash
# 1. Login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"admin"}'

# Response:
# {
#   "access_token": "eyJhbGci...",
#   "user": {
#     "id": "...",
#     "name": "admin",
#     "is_admin": true
#   }
# }

# 2. Sử dụng token
TOKEN="eyJhbGci..."

curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/protected
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
  "token_type": "Bearer",
  "user": {
    "id": "...",
    "name": "admin",
    "is_admin": true
  }
}
```

### OAuth Flow (Nếu cần)

```
GET /api/auth/oauth/login
Response: {
  "login_url": "http://localhost:8000/login/oauth/authorize?..."
}

GET /api/auth/oauth/callback?code=xxx
Response: {
  "access_token": "eyJhbGci..."
}
```

### Protected Routes

```
GET /api/protected
Authorization: Bearer <token>
```

## Lợi ích

1. **Quản lý user tập trung**: Tất cả users trong Casdoor
2. **Bảo mật tốt hơn**: JWT token với signature verification
3. **Mở rộng dễ dàng**: Thêm users, roles, permissions
4. **Tích hợp đơn giản**: API giống flow cũ (username/password)
5. **Có thể thêm OAuth sau**: Nếu cần web app với frontend

## Checklist Migration

- [ ] Setup Casdoor (docker-compose up -d)
- [ ] Tạo Application trong Casdoor
- [ ] Tạo user admin trong Casdoor
- [ ] Cấu hình `.env` với Client ID/Secret
- [ ] Update login handler → Dùng `auth.DirectLogin()`
- [ ] Update middleware → Dùng `auth.AuthMiddleware()`
- [ ] Test login với username/password
- [ ] Test protected routes với Bearer token
- [ ] Migrate users từ hệ thống cũ (nếu có)
- [ ] Deploy và monitor

## Troubleshooting

### Login failed
- Kiểm tra Application có tồn tại trong Casdoor
- Kiểm tra username/password đúng
- Kiểm tra Application có enable Password login

### Token verification failed
- Kiểm tra Casdoor đang chạy
- Kiểm tra certificate có thể fetch được không
- Kiểm tra token chưa hết hạn

### User not found
- Kiểm tra user có tồn tại trong Casdoor
- Kiểm tra Organization name đúng

