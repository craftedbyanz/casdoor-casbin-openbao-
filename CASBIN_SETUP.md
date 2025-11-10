# Casbin Integration Setup

## 🚀 Quick Start

### 1. Start services
```bash
docker-compose up -d
go mod tidy
go run cmd/server/main.go
```

### 2. Get admin token from Casdoor
```bash
ADMIN_TOKEN=$(curl -s -X POST "http://localhost:8000/api/login" \
  -d '{"application":"app-built-in","username":"admin","password":"123456","type":"token"}' \
  -H "Content-Type: application/json" | jq -r .data)
```

### 3. Check policies (auto-initialized on startup)
```bash
# Check policies from Casbin memory:
curl -X GET http://localhost:8080/api/admin/policies \
  -H "Authorization: Bearer $ADMIN_TOKEN"

# Check raw DB table casbin_rule:
curl -X GET http://localhost:8080/api/admin/debug/casbin-rules \
  -H "Authorization: Bearer $ADMIN_TOKEN"

# Manually reinitialize if needed:
curl -X POST http://localhost:8080/api/admin/init \
  -H "Authorization: Bearer $ADMIN_TOKEN"
```

## 📋 API Endpoints

### Admin APIs (require admin token)
- `POST /api/admin/init` - Initialize default policies
- `GET /api/admin/policies` - List all policies
- `POST /api/admin/policies` - Add policy
- `DELETE /api/admin/policies` - Remove policy
- `GET /api/admin/roles` - List role assignments
- `POST /api/admin/roles` - Assign role
- `DELETE /api/admin/roles` - Remove role

### Protected APIs (require authorization)
- `GET /api/users/profile` - User profile
- `GET /api/protected` - Protected resource
- `GET /api/users` - List users (admin only)

## 🔧 How it works

1. **Authentication**: Casdoor JWT verification
2. **Authorization**: Casbin policy enforcement
3. **Flow**: `AuthMiddleware` → `AuthzMiddleware` → `Handler`

## 📝 Policy Examples

```
# Policies (subject, object, action)
admin, /api/users, read
user, /api/users/profile, read

# Role assignments (user, role)
admin, admin
testuser, user
```