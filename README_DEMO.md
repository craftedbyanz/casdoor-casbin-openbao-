# Casdoor Integration Demo - Backend Golang với Echo

Đây là demo tích hợp Casdoor authentication vào backend Golang sử dụng framework Echo.

## Kiến trúc

```
┌──────────────┐
│   Frontend   │
│ (Web/API)    │
└──────┬───────┘
       │ Login via OAuth2/OIDC (Casdoor)
       ▼
┌──────────────┐
│   Casdoor    │
│ (AuthN IdP)  │
│ ├─ User, Org │
│ ├─ Role      │
│ └─ JWT issue │
└──────┬───────┘
       │ JWT Token
       ▼
┌──────────────┐
│   Backend    │
│ (Golang Echo)│
│ ├─ Verify JWT │
│ ├─ Middleware │
│ └─ Protected  │
│    Routes     │
└──────────────┘
```

## Yêu cầu

- Go 1.21+
- Docker và Docker Compose
- Casdoor đang chạy (qua docker-compose)

## Cài đặt

### 1. Khởi động Casdoor và PostgreSQL

```bash
docker-compose up -d
```

Casdoor sẽ chạy tại: `http://localhost:8000`

### 2. Cấu hình Casdoor Application

1. Truy cập `http://localhost:8000`
2. Đăng nhập với admin/admin (hoặc tài khoản admin mặc định)
3. Vào **Applications** → **Add**
4. Điền thông tin:
   - **Name**: app-built-in (hoặc tên bạn muốn)
   - **Organization**: built-in
   - **Redirect URLs**: `http://localhost:8080/api/auth/callback`
   - **Client ID**: (ghi lại)
   - **Client Secret**: (ghi lại)

### 3. Cấu hình Backend

Tạo file `.env` từ `.env.example`:

```bash
cp .env.example .env
```

Chỉnh sửa `.env`:

```env
SERVER_PORT=8080
SERVER_HOST=localhost

CASDOOR_ENDPOINT=http://localhost:8000
CASDOOR_CLIENT_ID=your_client_id_here
CASDOOR_CLIENT_SECRET=your_client_secret_here
CASDOOR_ORGANIZATION=built-in
CASDOOR_APPLICATION=app-built-in
CASDOOR_REDIRECT_URL=http://localhost:8080/api/auth/callback
```

### 4. Cài đặt dependencies và chạy

```bash
go mod download
go run cmd/server/main.go
```

Server sẽ chạy tại: `http://localhost:8080`

## API Endpoints

### Public Endpoints

#### 1. Health Check
```bash
GET /health
```

#### 2. Get Login URL
```bash
GET /api/auth/login
```

Response:
```json
{
  "login_url": "http://localhost:8000/api/login/oauth/authorize?...",
  "state": "random_state_string",
  "message": "Redirect to this URL to login"
}
```

#### 3. OAuth Callback
```bash
GET /api/auth/callback?code=xxx&state=xxx
```

Response:
```json
{
  "access_token": "eyJhbGci...",
  "token_type": "Bearer",
  "expires_in": 3600,
  "message": "Login successful. Use the access_token in Authorization header."
}
```

### Protected Endpoints (Require Bearer Token)

#### 4. Get Current User Info
```bash
GET /api/auth/me
Authorization: Bearer <token>
```

#### 5. Get User Profile
```bash
GET /api/users/profile
Authorization: Bearer <token>
```

#### 6. Get Protected Resource
```bash
GET /api/protected
Authorization: Bearer <token>
```

#### 7. Get All Users (Admin Only)
```bash
GET /api/users
Authorization: Bearer <token>
```

## Cách sử dụng

### Bước 1: Lấy Login URL

```bash
curl http://localhost:8080/api/auth/login
```

### Bước 2: Đăng nhập qua Browser

Mở `login_url` từ response ở bước 1 trong browser, đăng nhập với tài khoản Casdoor.

Sau khi đăng nhập thành công, browser sẽ redirect về:
```
http://localhost:8080/api/auth/callback?code=xxx&state=xxx
```

### Bước 3: Lấy Access Token

Copy `code` từ URL và gọi:

```bash
curl "http://localhost:8080/api/auth/callback?code=YOUR_CODE&state=YOUR_STATE"
```

Response sẽ chứa `access_token`.

### Bước 4: Sử dụng Token

```bash
# Get user info
curl -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  http://localhost:8080/api/auth/me

# Get profile
curl -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  http://localhost:8080/api/users/profile

# Get protected resource
curl -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  http://localhost:8080/api/protected

# Get users (admin only)
curl -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  http://localhost:8080/api/users
```

## Authentication Flow

1. **Client** → `GET /api/auth/login` → Nhận login URL
2. **Client** → Redirect user đến Casdoor login page
3. **User** → Đăng nhập trên Casdoor
4. **Casdoor** → Redirect về `/api/auth/callback?code=xxx`
5. **Backend** → Exchange code lấy access token từ Casdoor
6. **Backend** → Trả về access token cho client
7. **Client** → Sử dụng token trong header `Authorization: Bearer <token>`
8. **Backend** → Verify JWT token trong middleware
9. **Backend** → Cho phép truy cập protected routes

## Middleware

### AuthMiddleware

Middleware này verify JWT token từ Casdoor và lưu user info vào context:

```go
protectedGroup.Use(auth.AuthMiddleware())
```

### RequireAdmin

Middleware kiểm tra user có phải admin không:

```go
adminGroup.Use(auth.AuthMiddleware())
adminGroup.Use(auth.RequireAdmin())
```

## Cấu trúc Code

```
.
├── cmd/
│   └── server/
│       └── main.go          # Main server
├── config/
│   └── config.go            # Configuration
├── internal/
│   ├── auth/
│   │   ├── casdoor.go       # Casdoor JWT verification
│   │   └── middleware.go    # Auth middleware
│   ├── config/
│   │   └── config.go        # Config initialization
│   └── handler/
│       ├── auth.go          # Auth handlers (login, callback)
│       └── user.go          # User handlers (protected routes)
├── docker-compose.yml       # Casdoor + PostgreSQL
└── go.mod                   # Dependencies
```

## Troubleshooting

### Lỗi: "failed to fetch certificate"

- Kiểm tra Casdoor đang chạy: `curl http://localhost:8000/api/certs`
- Kiểm tra `CASDOOR_ENDPOINT` trong `.env`

### Lỗi: "invalid token"

- Kiểm tra token còn hạn không
- Kiểm tra token đúng format không (Bearer token)
- Kiểm tra certificate từ Casdoor có thể fetch được không

### Lỗi: "token exchange failed"

- Kiểm tra `CASDOOR_CLIENT_ID` và `CASDOOR_CLIENT_SECRET`
- Kiểm tra `CASDOOR_REDIRECT_URL` khớp với cấu hình trong Casdoor
- Kiểm tra code chưa expire (thường chỉ dùng được 1 lần)

## Notes

- JWT token từ Casdoor được verify bằng public key từ `/api/certs`
- Certificate được cache sau lần fetch đầu tiên
- Token có thể chứa thông tin: user ID, name, email, roles, isAdmin
- Protected routes yêu cầu Bearer token trong header

## Tích hợp với Frontend

Frontend có thể:

1. Gọi `GET /api/auth/login` để lấy login URL
2. Redirect user đến login URL
3. Nhận callback với code
4. Gọi `GET /api/auth/callback?code=xxx` để lấy token
5. Lưu token và sử dụng cho các API calls sau

Hoặc có thể implement OAuth flow hoàn toàn ở frontend và chỉ gửi token lên backend để verify.

