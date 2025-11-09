# Giải Thích: Tại Sao Cần Public Key Từ Casdoor?

## Câu Hỏi

Tại sao phải call lên Casdoor để lấy public key? Public key có phải đi theo application không?

## Trả Lời Ngắn Gọn

**Public key là của Casdoor (Identity Provider), không phải của application.**

- Casdoor **sign** JWT token bằng **private key** của Casdoor
- Backend **verify** JWT token bằng **public key** của Casdoor
- Public key là **chung cho tất cả applications** trong Casdoor instance
- Không phải mỗi application có public key riêng

---

## Giải Thích Chi Tiết

### 1. JWT Signing Process

```
┌─────────┐
│ Casdoor │
└────┬────┘
     │
     │ 1. User login thành công
     │
     │ 2. Tạo JWT token với claims:
     │    {
     │      "sub": "user-id",
     │      "name": "admin",
     │      "isAdmin": true,
     │      ...
     │    }
     │
     │ 3. Sign JWT với Casdoor PRIVATE KEY
     │    Signature = RSA_Sign(token_data, casdoor_private_key)
     │
     │ 4. Trả về JWT token:
     │    eyJhbGci...header...payload...signature
     │
     ▼
┌─────────┐
│ Backend │
└─────────┘
```

**Quan trọng:**
- JWT được **sign** bởi **Casdoor private key**
- Private key chỉ có ở Casdoor (bí mật)
- Public key được dùng để **verify** signature

---

### 2. JWT Verification Process

```
┌─────────┐
│ Backend │
└────┬────┘
     │
     │ 1. Nhận JWT token từ client
     │    eyJhbGci...header...payload...signature
     │
     │ 2. Lấy public key từ Casdoor
     │    GET /.well-known/jwks
     │
     │ 3. Verify signature:
     │    RSA_Verify(token_data, signature, casdoor_public_key)
     │
     │ 4. Nếu signature hợp lệ:
     │    - Token được sign bởi Casdoor
     │    - Token không bị giả mạo
     │    - Extract claims từ token
     │
     ▼
┌─────────┐
│ Success │
└─────────┘
```

**Quan trọng:**
- Public key dùng để **verify** signature
- Nếu verify thành công → Token đúng từ Casdoor
- Nếu verify thất bại → Token giả mạo hoặc không từ Casdoor

---

### 3. Tại Sao Public Key Là Của Casdoor, Không Phải Application?

#### Public Key = Casdoor Instance Level

```
┌─────────────────────────────────┐
│      Casdoor Instance           │
│  (http://localhost:8000)        │
│                                  │
│  ┌──────────────────────────┐   │
│  │  Private Key (Secret)    │   │
│  │  - Dùng để SIGN tokens   │   │
│  │  - Chỉ có ở Casdoor      │   │
│  └──────────────────────────┘   │
│                                  │
│  ┌──────────────────────────┐   │
│  │  Public Key (Public)     │   │
│  │  - Dùng để VERIFY tokens │   │
│  │  - Công khai, ai cũng    │   │
│  │    có thể lấy được       │   │
│  └──────────────────────────┘   │
│                                  │
│  ┌──────────┐  ┌──────────┐    │
│  │ App 1    │  │ App 2    │    │
│  │ (myapp)  │  │ (app2)   │    │
│  └──────────┘  └──────────┘    │
│                                  │
│  Tất cả apps dùng CHUNG          │
│  public key này                  │
└─────────────────────────────────┘
```

**Lý do:**
- Casdoor là **Identity Provider (IdP)**
- Tất cả applications trong Casdoor instance dùng **chung 1 cặp key**
- Private key: Sign tokens cho tất cả apps
- Public key: Verify tokens từ tất cả apps

---

### 4. So Sánh Với Application

#### Application Level (Không đúng)

```
❌ Mỗi application có public key riêng
App 1 → Public Key 1
App 2 → Public Key 2
App 3 → Public Key 3
```

**Vấn đề:**
- Mỗi app phải quản lý public key riêng
- Phức tạp hơn
- Không cần thiết

#### Casdoor Instance Level (Đúng)

```
✅ Tất cả applications dùng chung public key
App 1 ──┐
App 2 ──┼──→ Public Key (chung)
App 3 ──┘
```

**Lợi ích:**
- Đơn giản hơn
- Dễ quản lý
- Public key là của Casdoor instance

---

### 5. Tại Sao Phải Fetch Public Key?

#### Option 1: Fetch từ Casdoor (Hiện tại)

```go
func GetPublicKey() (*rsa.PublicKey, error) {
    // Fetch từ /.well-known/jwks
    GET http://localhost:8000/.well-known/jwks
    
    // Parse và extract public key
    // Cache để dùng lại
}
```

**Ưu điểm:**
- Tự động update nếu Casdoor đổi key
- Không cần config thủ công
- Luôn có public key mới nhất

**Nhược điểm:**
- Cần network call (nhưng có cache)
- Phụ thuộc vào Casdoor endpoint

#### Option 2: Config Static Public Key

```go
// Trong .env
CASDOOR_CERTIFICATE="-----BEGIN CERTIFICATE-----..."

// Trong code
if cfg.Casdoor.Certificate != "" {
    // Dùng certificate từ config
}
```

**Ưu điểm:**
- Không cần network call
- Hoạt động offline

**Nhược điểm:**
- Phải update thủ công nếu Casdoor đổi key
- Dễ quên update

---

### 6. Flow Thực Tế

```
┌─────────┐
│ Client  │
└────┬────┘
     │ 1. Login với username/password
     ▼
┌─────────┐
│ Backend │
└────┬────┘
     │ 2. Forward đến Casdoor
     ▼
┌─────────┐
│ Casdoor │
└────┬────┘
     │ 3. Verify username/password
     │ 4. Sign JWT với PRIVATE KEY
     │    Signature = Sign(token, casdoor_private_key)
     │ 5. Trả về JWT token
     ▼
┌─────────┐
│ Backend │
└────┬────┘
     │ 6. Nhận JWT token
     │ 7. Lấy PUBLIC KEY từ Casdoor
     │    GET /.well-known/jwks
     │ 8. Verify signature
     │    Verify(token, signature, casdoor_public_key)
     │ 9. Extract claims
     ▼
┌─────────┐
│ Client │
└────────┘
```

---

### 7. Tại Sao Không Cần Public Key Theo Application?

#### JWT Token Structure

```
JWT Token = Header.Payload.Signature

Header: {
  "alg": "RS256",
  "kid": "cert-built-in",  // Key ID (của Casdoor)
  "typ": "JWT"
}

Payload: {
  "sub": "user-id",
  "name": "admin",
  "azp": "myapp",  // Application ID
  "iss": "http://localhost:8000",  // Casdoor issuer
  ...
}

Signature: RSA_Sign(header.payload, casdoor_private_key)
```

**Quan trọng:**
- Signature được sign bởi **Casdoor private key**
- Application ID (`azp`) chỉ là **claim trong token**, không ảnh hưởng đến signature
- Public key để verify signature → Phải là của Casdoor

---

### 8. Ví Dụ Thực Tế

#### Scenario: 2 Applications

```
Application 1: "myapp"
Application 2: "app2"

Cả 2 apps đều dùng Casdoor instance: http://localhost:8000
```

**Khi user login với App 1:**
1. Casdoor sign JWT với **Casdoor private key**
2. Token có `azp: "myapp"` (application ID)
3. Backend verify với **Casdoor public key** (chung)
4. Verify thành công → Token hợp lệ

**Khi user login với App 2:**
1. Casdoor sign JWT với **Casdoor private key** (cùng key)
2. Token có `azp: "app2"` (application ID khác)
3. Backend verify với **Casdoor public key** (cùng key)
4. Verify thành công → Token hợp lệ

**Kết luận:**
- Cùng public key để verify
- Application ID chỉ để biết token từ app nào
- Không cần public key riêng cho mỗi app

---

### 9. Có Thể Config Static Public Key Không?

**Có**, nếu muốn:

```env
# .env
CASDOOR_CERTIFICATE="-----BEGIN CERTIFICATE-----
MIIE3TCCAsWgAwIBAgIDAeJAMA0GCSqGSIb3DQEBCwUAMCgxDjAMBgNVBAoTBWFk
...
-----END CERTIFICATE-----"
```

Code sẽ ưu tiên dùng certificate từ config nếu có.

---

## Tóm Tắt

1. **Public key là của Casdoor instance**, không phải của application
2. **Tất cả applications** trong Casdoor instance dùng **chung public key**
3. **JWT được sign** bởi Casdoor private key → Cần Casdoor public key để verify
4. **Application ID** (`azp`) chỉ là claim trong token, không ảnh hưởng đến signature
5. **Fetch public key** từ Casdoor để đảm bảo luôn có key mới nhất
6. **Có thể cache** public key để tránh fetch mỗi lần
7. **Có thể config static** public key nếu muốn

---

## Code Flow

```go
// 1. Verify token
func VerifyToken(tokenString string) (*CasdoorClaims, error) {
    // 2. Lấy public key (của Casdoor instance)
    publicKey, err := GetPublicKey()
    
    // 3. Verify JWT signature với public key
    token, err := jwt.ParseWithClaims(tokenString, &CasdoorClaims{}, 
        func(token *jwt.Token) (interface{}, error) {
            return publicKey, nil  // Dùng Casdoor public key
        })
    
    // 4. Extract claims
    claims := token.Claims.(*CasdoorClaims)
    // claims.Azp = "myapp"  // Application ID (chỉ là claim)
    
    return claims
}
```

**Quan trọng:**
- Public key = Casdoor instance level
- Application ID = Claim trong token (không ảnh hưởng đến verification)

