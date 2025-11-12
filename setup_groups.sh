#!/bin/bash

# Setup Group-Based Policies for ANZ Operations Portal
# Usage: ./setup_groups.sh

echo "🔐 Setting up group-based policies..."

# Get admin token
echo "Getting admin token..."
ADMIN_TOKEN=$(curl -s -X POST "http://localhost:8000/api/login" \
  -d '{"application":"app-built-in","username":"admin","password":"123456","type":"token"}' \
  -H "Content-Type: application/json" | jq -r .data)

if [ "$ADMIN_TOKEN" = "null" ] || [ -z "$ADMIN_TOKEN" ]; then
    echo "❌ Failed to get admin token. Make sure Casdoor is running on port 8000"
    exit 1
fi

echo "✅ Got admin token"

# 1. Create Dashboard-Only Group
echo "📊 Creating dashboard_group..."
curl -s -X POST http://localhost:8080/api/admin/policies \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"subject":"dashboard_group","object":"/api/auth/me","action":"read"}' > /dev/null

# 2. Create Transaction Group
echo "💰 Creating transaction_group..."
curl -s -X POST http://localhost:8080/api/admin/policies \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"subject":"transaction_group","object":"/api/auth/me","action":"read"}' > /dev/null

curl -s -X POST http://localhost:8080/api/admin/policies \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"subject":"transaction_group","object":"/api/transactions/my","action":"read"}' > /dev/null

curl -s -X POST http://localhost:8080/api/admin/policies \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"subject":"transaction_group","object":"/api/transactions","action":"write"}' > /dev/null

# 3. Create Order Group
echo "🛒 Creating order_group..."
curl -s -X POST http://localhost:8080/api/admin/policies \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"subject":"order_group","object":"/api/auth/me","action":"read"}' > /dev/null

curl -s -X POST http://localhost:8080/api/admin/policies \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"subject":"order_group","object":"/api/orders/my","action":"read"}' > /dev/null

curl -s -X POST http://localhost:8080/api/admin/policies \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"subject":"order_group","object":"/api/orders","action":"write"}' > /dev/null

# 4. Create Full Access Group
echo "👑 Creating full_access_group..."
curl -s -X POST http://localhost:8080/api/admin/policies \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"subject":"full_access_group","object":"/api/auth/me","action":"read"}' > /dev/null

curl -s -X POST http://localhost:8080/api/admin/policies \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"subject":"full_access_group","object":"/api/transactions","action":"read"}' > /dev/null

curl -s -X POST http://localhost:8080/api/admin/policies \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"subject":"full_access_group","object":"/api/transactions/my","action":"read"}' > /dev/null

curl -s -X POST http://localhost:8080/api/admin/policies \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"subject":"full_access_group","object":"/api/transactions","action":"write"}' > /dev/null

curl -s -X POST http://localhost:8080/api/admin/policies \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"subject":"full_access_group","object":"/api/orders","action":"read"}' > /dev/null

curl -s -X POST http://localhost:8080/api/admin/policies \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"subject":"full_access_group","object":"/api/orders/my","action":"read"}' > /dev/null

curl -s -X POST http://localhost:8080/api/admin/policies \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"subject":"full_access_group","object":"/api/orders","action":"write"}' > /dev/null

echo "✅ Groups created successfully!"

# Show current policies
echo ""
echo "📋 Current policies:"
curl -s -X GET http://localhost:8080/api/admin/policies \
  -H "Authorization: Bearer $ADMIN_TOKEN" | jq .

echo ""
echo "🎯 Example: Assign users to groups"
echo "# Dashboard only:"
echo "curl -X POST http://localhost:8080/api/admin/roles -H \"Authorization: Bearer $ADMIN_TOKEN\" -H \"Content-Type: application/json\" -d '{\"user\":\"user1\",\"role\":\"dashboard_group\"}'"
echo ""
echo "# Transaction access:"
echo "curl -X POST http://localhost:8080/api/admin/roles -H \"Authorization: Bearer $ADMIN_TOKEN\" -H \"Content-Type: application/json\" -d '{\"user\":\"user2\",\"role\":\"transaction_group\"}'"
echo ""
echo "# Order access:"
echo "curl -X POST http://localhost:8080/api/admin/roles -H \"Authorization: Bearer $ADMIN_TOKEN\" -H \"Content-Type: application/json\" -d '{\"user\":\"user3\",\"role\":\"order_group\"}'"
echo ""
echo "# Full access:"
echo "curl -X POST http://localhost:8080/api/admin/roles -H \"Authorization: Bearer $ADMIN_TOKEN\" -H \"Content-Type: application/json\" -d '{\"user\":\"poweruser\",\"role\":\"full_access_group\"}'"