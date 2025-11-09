# Fix: "Grant_type: is not supported in this application"

## Vấn đề

Khi truy cập login URL, bạn gặp lỗi:
```
There was a problem signing you in..
Grant_type: is not supported in this application
```

## Nguyên nhân

Application trong Casdoor chưa được cấu hình để hỗ trợ OAuth2 **Authorization Code Grant**.

## Giải pháp

### Bước 1: Vào Casdoor Application Settings

1. Truy cập: http://localhost:8000
2. Đăng nhập với tài khoản admin
3. Vào **Applications** → Chọn application của bạn (hoặc tạo mới)

### Bước 2: Cấu hình Grant Types

Trong trang cấu hình Application, tìm phần **Grant Types** hoặc **OAuth Settings**:

1. **Enable Grant Types**:
   - ✅ **Authorization Code** (BẮT BUỘC - phải bật)
   - ✅ **Refresh Token** (khuyến nghị)
   - ❌ **Implicit** (không cần cho flow này)
   - ❌ **Password** (không cần cho OAuth flow)
   - ❌ **Client Credentials** (không cần)

2. **Các cấu hình khác cần kiểm tra**:
   - **Redirect URLs**: Phải có `http://localhost:8080/api/auth/callback`
   - **Enable Password**: Có thể bật nếu muốn dùng password login
   - **Enable Sign Up**: Có thể bật nếu muốn cho phép đăng ký

### Bước 3: Lưu cấu hình

1. Click **Save** để lưu thay đổi
2. Đảm bảo Application đã được **Enabled**

### Bước 4: Kiểm tra lại

Sau khi cấu hình xong:

1. **Restart server** (nếu đang chạy):
   ```bash
   # Stop server (Ctrl+C)
   go run cmd/server/main.go
   ```

2. **Test lại login URL**:
   ```bash
   curl http://localhost:8080/api/auth/login
   ```

3. **Mở login_url trong browser** - Bây giờ sẽ không còn lỗi "Grant_type not supported"

## Cấu hình Application mẫu

Nếu tạo Application mới, cấu hình như sau:

### Basic Information
- **Name**: `app-built-in` (hoặc tên bạn muốn)
- **Organization**: `built-in`
- **Display Name**: `My Application`
- **Enable Password**: `true` (nếu muốn dùng password login)

### OAuth Settings
- **Redirect URLs**: 
  ```
  http://localhost:8080/api/auth/callback
  ```
- **Grant Types**:
  - ✅ **Authorization Code** (BẮT BUỘC)
  - ✅ **Refresh Token** (khuyến nghị)

### Advanced Settings
- **Enable Sign Up**: `true` (tùy chọn)
- **Enable Password**: `true` (nếu muốn dùng password login)
- **Enable Code**: `true` (cho OAuth flow)

## Kiểm tra Application đã được cấu hình đúng

Sau khi cấu hình, bạn có thể kiểm tra bằng cách:

1. **Xem Application details**:
   - Vào Applications → Click vào application
   - Kiểm tra "Grant Types" có "Authorization Code" không

2. **Test OAuth flow**:
   ```bash
   # Lấy login URL
   curl http://localhost:8080/api/auth/login
   
   # Mở login_url trong browser
   # Nếu không còn lỗi "Grant_type not supported" thì đã OK
   ```

## Troubleshooting

### Vẫn gặp lỗi sau khi cấu hình

1. **Kiểm tra Application có được Enable không**:
   - Vào Applications → Đảm bảo application có status "Enabled"

2. **Kiểm tra Organization**:
   - Đảm bảo Organization name trong `.env` khớp với Organization của Application

3. **Kiểm tra Client ID**:
   - Đảm bảo `CASDOOR_CLIENT_ID` trong `.env` khớp với Client ID của Application

4. **Clear browser cache**:
   - Thử mở login URL trong incognito/private window

5. **Kiểm tra Casdoor logs**:
   ```bash
   docker-compose logs casdoor | tail -50
   ```

### Lỗi khác

Nếu gặp lỗi khác, xem `TROUBLESHOOTING.md` để biết thêm chi tiết.

