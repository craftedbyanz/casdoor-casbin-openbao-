# 🔐 Policy Management Commands

## Quick Setup
```bash
# Make scripts executable
chmod +x setup_groups.sh test_groups.sh

# Setup all groups
./setup_groups.sh

# Test specific user/group
./test_groups.sh testuser dashboard_group
```

## Manual Commands

### 1. Get Admin Token
```bash
ADMIN_TOKEN=$(curl -s -X POST "http://localhost:8000/api/login" \
  -d '{"application":"app-built-in","username":"admin","password":"123456","type":"token"}' \
  -H "Content-Type: application/json" | jq -r .data)
```

### 2. Add Policies (Create Groups)
```bash
# Dashboard-only group
curl -X POST http://localhost:8080/api/admin/policies \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"subject":"dashboard_group","object":"/api/auth/me","action":"read"}'

# Transaction group
curl -X POST http://localhost:8080/api/admin/policies \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"subject":"transaction_group","object":"/api/transactions/my","action":"read"}'

# Order group  
curl -X POST http://localhost:8080/api/admin/policies \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"subject":"order_group","object":"/api/orders/my","action":"read"}'
```

### 3. Assign Users to Groups
```bash
# Dashboard only
curl -X POST http://localhost:8080/api/admin/roles \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"user":"user1","role":"dashboard_group"}'

# Transaction access
curl -X POST http://localhost:8080/api/admin/roles \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"user":"user2","role":"transaction_group"}'

# Order access
curl -X POST http://localhost:8080/api/admin/roles \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"user":"user3","role":"order_group"}'
```

### 4. Remove Policies
```bash
# Remove specific policy
curl -X DELETE http://localhost:8080/api/admin/policies \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"subject":"transaction_group","object":"/api/transactions/my","action":"read"}'

# Remove user from group
curl -X DELETE http://localhost:8080/api/admin/roles \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"user":"user2","role":"transaction_group"}'
```

### 5. Check Current State
```bash
# List all policies
curl -X GET http://localhost:8080/api/admin/policies \
  -H "Authorization: Bearer $ADMIN_TOKEN"

# List all role assignments
curl -X GET http://localhost:8080/api/admin/roles \
  -H "Authorization: Bearer $ADMIN_TOKEN"

# Check raw database
curl -X GET http://localhost:8080/api/admin/debug/casbin-rules \
  -H "Authorization: Bearer $ADMIN_TOKEN"
```

## 🎯 Group Access Matrix

| Group | Dashboard | Transactions | Orders | Admin |
|-------|-----------|--------------|--------|-------|
| `dashboard_group` | ✅ | ❌ | ❌ | ❌ |
| `transaction_group` | ✅ | ✅ | ❌ | ❌ |
| `order_group` | ✅ | ❌ | ✅ | ❌ |
| `full_access_group` | ✅ | ✅ | ✅ | ❌ |
| `admin` | ✅ | ✅ | ✅ | ✅ |

## 🧪 Test Access
```bash
# Get user token
USER_TOKEN=$(curl -s -X POST "http://localhost:8000/api/login" \
  -d '{"application":"app-built-in","username":"testuser","password":"123456","type":"token"}' \
  -H "Content-Type: application/json" | jq -r .data)

# Test endpoints
curl -X GET http://localhost:8080/api/auth/me -H "Authorization: Bearer $USER_TOKEN"
curl -X GET http://localhost:8080/api/transactions/my -H "Authorization: Bearer $USER_TOKEN"
curl -X GET http://localhost:8080/api/orders/my -H "Authorization: Bearer $USER_TOKEN"
```