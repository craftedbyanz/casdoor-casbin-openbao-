# Troubleshooting Guide

## Vấn đề: 404 Not Found khi truy cập login URL

### Nguyên nhân

1. **Client ID chưa được cấu hình đúng** - URL có `client_id=` rỗng
2. **Endpoint OAuth không đúng** - Đã được fix thành `/login/oauth/authorize`

### Giải pháp

#### Bước 1: Kiểm tra file .env

Đảm bảo file `.env` có các biến sau:

```env
CASDOOR_ENDPOINT=http://localhost:8000
CASDOOR_CLIENT_ID=your_real_client_id_here
CASDOOR_CLIENT_SECRET=your_real_client_secret_here
CASDOOR_ORGANIZATION=built-in
CASDOOR_APPLICATION=app-built-in
CASDOOR_REDIRECT_URL=http://localhost:8080/api/auth/callback
```

#### Bước 2: Lấy Client ID và Client Secret từ Casdoor

1. Truy cập Casdoor: http://localhost:8000
2. Đăng nhập với tài khoản admin
3. Vào **Applications** → Tìm application của bạn (hoặc tạo mới)
4. Copy **Client ID** và **Client Secret**
5. Cập nhật vào file `.env`

#### Bước 3: Tạo Application mới (nếu chưa có)

1. Vào **Applications** → **Add**
2. Điền thông tin:
   - **Name**: `app-built-in` (hoặc tên bạn muốn)
   - **Organization**: `built-in`
   - **Redirect URLs**: `http://localhost:8080/api/auth/callback`
   - **Enable Password**: `true` (nếu muốn dùng password login)
3. Save và copy **Client ID** và **Client Secret**

#### Bước 4: Kiểm tra cấu hình

Sau khi cập nhật `.env`, restart server:

```bash
# Stop server (Ctrl+C)
# Start lại
go run cmd/server/main.go
```

#### Bước 5: Test lại

```bash
curl http://localhost:8080/api/auth/login
```

Response phải có `login_url` với `client_id` không rỗng:

```json
{
  "login_url": "http://localhost:8000/login/oauth/authorize?client_id=YOUR_CLIENT_ID&...",
  "state": "...",
  "config": {
    "endpoint": "http://localhost:8000",
    "client_id": true,
    "redirect_url": "http://localhost:8080/api/auth/callback"
  }
}
```

### Kiểm tra Endpoint OAuth

Từ OpenID Configuration của Casdoor:

```bash
curl http://localhost:8000/.well-known/openid-configuration
```

Bạn sẽ thấy:
- `authorization_endpoint`: `http://localhost:8000/login/oauth/authorize`
- `token_endpoint`: `http://localhost:8000/api/login/oauth/access_token`

### Common Issues

#### Issue 1: Client ID rỗng trong URL

**Nguyên nhân**: File `.env` không được load hoặc `CASDOOR_CLIENT_ID` không được set

**Giải pháp**:
1. Đảm bảo file `.env` ở thư mục root của project
2. Kiểm tra `CASDOOR_CLIENT_ID` có giá trị trong `.env`
3. Restart server sau khi cập nhật `.env`

#### Issue 2: 404 Not Found khi truy cập login URL

**Nguyên nhân**: 
- Client ID không tồn tại trong Casdoor
- Application chưa được tạo
- Redirect URL không khớp

**Giải pháp**:
1. Kiểm tra Application có tồn tại trong Casdoor không
2. Đảm bảo Redirect URL trong Application khớp với `CASDOOR_REDIRECT_URL` trong `.env`
3. Kiểm tra Organization name đúng

#### Issue 3: Token exchange failed

**Nguyên nhân**:
- Client Secret sai
- Code đã hết hạn hoặc đã được dùng
- Redirect URI không khớp

**Giải pháp**:
1. Kiểm tra `CASDOOR_CLIENT_SECRET` đúng chưa
2. Code chỉ dùng được 1 lần, nếu đã dùng thì phải login lại
3. Đảm bảo Redirect URI trong callback khớp với cấu hình

### Debug Steps

1. **Kiểm tra config được load đúng không**:
   ```bash
   curl http://localhost:8080/api/auth/login
   ```
   Xem field `config` trong response để kiểm tra.

2. **Kiểm tra Casdoor đang chạy**:
   ```bash
   curl http://localhost:8000/api/get-account
   ```

3. **Kiểm tra OpenID Configuration**:
   ```bash
   curl http://localhost:8000/.well-known/openid-configuration
   ```

4. **Kiểm tra Application trong Casdoor**:
   - Vào http://localhost:8000
   - Đăng nhập admin
   - Vào Applications → Kiểm tra application có tồn tại không

### Logs

Nếu vẫn gặp vấn đề, kiểm tra logs:

1. **Server logs**: Xem console output khi chạy `go run cmd/server/main.go`
2. **Casdoor logs**: 
   ```bash
   docker-compose logs casdoor
   ```

### Liên hệ

Nếu vẫn không giải quyết được, vui lòng cung cấp:
1. Response từ `curl http://localhost:8080/api/auth/login`
2. Nội dung file `.env` (ẩn Client Secret)
3. Logs từ server và Casdoor

