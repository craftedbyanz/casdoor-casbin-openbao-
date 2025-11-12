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
- `POST /api/admin/policies` - Add policy: `{"subject":"group_name","object":"/api/endpoint","action":"read|write"}`
- `DELETE /api/admin/policies` - Remove policy: `{"subject":"group_name","object":"/api/endpoint","action":"read|write"}`
- `GET /api/admin/roles` - List role assignments
- `POST /api/admin/roles` - Assign user to group: `{"user":"username","role":"group_name"}`
- `DELETE /api/admin/roles` - Remove user from group: `{"user":"username","role":"group_name"}`

### Protected APIs (require authorization)
- `GET /api/users/profile` - User profile
- `GET /api/protected` - Protected resource
- `GET /api/users` - List users (admin only)

## 🔧 How it works

1. **Authentication**: Casdoor JWT verification
2. **Authorization**: Casbin policy enforcement
3. **Flow**: `AuthMiddleware` → `AuthzMiddleware` → `Handler`

## 📝 Group-Based Policy Management

### Create Groups
```bash
# Create dashboard-only group
curl -X POST http://localhost:8080/api/admin/policies \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"subject":"dashboard_group","object":"/api/auth/me","action":"read"}'

# Create transaction group
curl -X POST http://localhost:8080/api/admin/policies \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"subject":"transaction_group","object":"/api/transactions","action":"read"}'

curl -X POST http://localhost:8080/api/admin/policies \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"subject":"transaction_group","object":"/api/transactions/my","action":"read"}'

curl -X POST http://localhost:8080/api/admin/policies \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"subject":"transaction_group","object":"/api/transactions","action":"write"}'

# Create order group
curl -X POST http://localhost:8080/api/admin/policies \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"subject":"order_group","object":"/api/orders","action":"read"}'

curl -X POST http://localhost:8080/api/admin/policies \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"subject":"order_group","object":"/api/orders/my","action":"read"}'

curl -X POST http://localhost:8080/api/admin/policies \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"subject":"order_group","object":"/api/orders","action":"write"}'
```

### Assign Users to Groups
```bash
# Assign user to dashboard-only group
curl -X POST http://localhost:8080/api/admin/roles \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"user":"testuser1","role":"dashboard_group"}'

# Assign user to transaction group
curl -X POST http://localhost:8080/api/admin/roles \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"user":"testuser2","role":"transaction_group"}'

# Assign user to order group
curl -X POST http://localhost:8080/api/admin/roles \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"user":"testuser3","role":"order_group"}'

# Assign user to multiple groups
curl -X POST http://localhost:8080/api/admin/roles \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"user":"poweruser","role":"transaction_group"}'

curl -X POST http://localhost:8080/api/admin/roles \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"user":"poweruser","role":"order_group"}'
```

### Remove Policies
```bash
# Remove specific policy
curl -X DELETE http://localhost:8080/api/admin/policies \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"subject":"transaction_group","object":"/api/transactions","action":"write"}'

# Remove user from group
curl -X DELETE http://localhost:8080/api/admin/roles \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"user":"testuser2","role":"transaction_group"}'
```

### Check Current Policies
```bash
# List all policies
curl -X GET http://localhost:8080/api/admin/policies \
  -H "Authorization: Bearer $ADMIN_TOKEN"

# List all role assignments
curl -X GET http://localhost:8080/api/admin/roles \
  -H "Authorization: Bearer $ADMIN_TOKEN"
```

## 🎯 Group Access Matrix

| Group | Dashboard | Transactions | Orders | Users (Admin) |
|-------|-----------|--------------|--------|--------------|
| dashboard_group | ✅ | ❌ | ❌ | ❌ |
| transaction_group | ✅ | ✅ | ❌ | ❌ |
| order_group | ✅ | ❌ | ✅ | ❌ |
| admin | ✅ | ✅ | ✅ | ✅ |