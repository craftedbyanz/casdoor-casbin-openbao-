# Quick Start Guide

## Bước 1: Khởi động Casdoor

```bash
docker-compose up -d
```

Chờ vài giây để Casdoor khởi động, sau đó truy cập: http://localhost:8000

## Bước 2: Cấu hình Casdoor Application

1. Đăng nhập vào Casdoor:
   - URL: http://localhost:8000
   - Username: `admin`
   - Password: `123` (hoặc password mặc định của bạn)

2. Tạo Application:
   - Vào **Applications** → **Add**
   - Điền thông tin:
     - **Name**: `app-built-in`
     - **Organization**: `built-in`
     - **Redirect URLs**: `http://localhost:8080/api/auth/callback`
     - **Grant Types**: ✅ **Enable "Authorization Code"** (QUAN TRỌNG!)
     - **Client ID**: Copy và lưu lại
     - **Client Secret**: Copy và lưu lại
   
   **⚠️ LƯU Ý QUAN TRỌNG**: Phải enable **Authorization Code** trong Grant Types, nếu không sẽ gặp lỗi "Grant_type: is not supported in this application"

## Bước 3: Cấu hình Backend

Tạo file `.env`:

```bash
cat > .env << EOF
SERVER_PORT=8080
SERVER_HOST=localhost

CASDOOR_ENDPOINT=http://localhost:8000
CASDOOR_CLIENT_ID=your_client_id_from_step_2
CASDOOR_CLIENT_SECRET=your_client_secret_from_step_2
CASDOOR_ORGANIZATION=built-in
CASDOOR_APPLICATION=app-built-in
CASDOOR_REDIRECT_URL=http://localhost:8080/api/auth/callback
EOF
```

Thay `your_client_id_from_step_2` và `your_client_secret_from_step_2` bằng giá trị thực từ Casdoor.

## Bước 4: Cài đặt và chạy Backend

```bash
# Cài đặt dependencies
go mod download

# Chạy server
go run cmd/server/main.go
```

Hoặc sử dụng Makefile:

```bash
make setup  # Lần đầu tiên
make run    # Chạy server
```

## Bước 5: Test API

### 1. Lấy Login URL

```bash
curl http://localhost:8080/api/auth/login
```

Response sẽ có `login_url`. Mở URL này trong browser.

### 2. Đăng nhập

- Browser sẽ redirect đến Casdoor login page
- Đăng nhập với tài khoản Casdoor
- Sau khi đăng nhập, browser sẽ redirect về callback URL với `code`

### 3. Lấy Access Token

Copy `code` từ URL callback và gọi:

```bash
curl "http://localhost:8080/api/auth/callback?code=YOUR_CODE&state=YOUR_STATE"
```

Response sẽ chứa `access_token`.

### 4. Sử dụng Token

```bash
# Set token
TOKEN="your_access_token_here"

# Get user info
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/auth/me

# Get profile
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/users/profile

# Get protected resource
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:8080/api/protected
```

## Troubleshooting

### Lỗi: "failed to fetch certificate"
- Kiểm tra Casdoor đang chạy: `curl http://localhost:8000/api/certs`
- Kiểm tra `CASDOOR_ENDPOINT` trong `.env`

### Lỗi: "token exchange failed"
- Kiểm tra `CASDOOR_CLIENT_ID` và `CASDOOR_CLIENT_SECRET` đúng chưa
- Kiểm tra `CASDOOR_REDIRECT_URL` khớp với cấu hình trong Casdoor
- Code chỉ dùng được 1 lần, nếu đã dùng rồi thì phải login lại

### Lỗi: "invalid token"
- Kiểm tra token chưa hết hạn
- Kiểm tra format: `Authorization: Bearer <token>`
- Token phải từ Casdoor và chưa bị revoke

## API Endpoints

### Public
- `GET /health` - Health check
- `GET /api/auth/login` - Get OAuth login URL
- `GET /api/auth/callback?code=xxx` - OAuth callback

### Protected (Require Bearer Token)
- `GET /api/auth/me` - Get current user info
- `GET /api/users/profile` - Get user profile
- `GET /api/protected` - Get protected resource
- `GET /api/users` - Get all users (admin only)

## Next Steps

- Xem `README_DEMO.md` để biết chi tiết về kiến trúc và cách tích hợp
- Tích hợp với frontend
- Thêm Casbin authorization
- Thêm OpenBao/Vault cho secrets management

