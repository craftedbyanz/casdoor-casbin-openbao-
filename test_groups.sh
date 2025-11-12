#!/bin/bash

# Test Group-Based Access Control
# Usage: ./test_groups.sh [username] [group_name]

USERNAME=${1:-"testuser"}
GROUP=${2:-"dashboard_group"}

echo "🧪 Testing group access for user: $USERNAME in group: $GROUP"

# Get admin token
ADMIN_TOKEN=$(curl -s -X POST "http://localhost:8000/api/login" \
  -d '{"application":"app-built-in","username":"admin","password":"123456","type":"token"}' \
  -H "Content-Type: application/json" | jq -r .data)

if [ "$ADMIN_TOKEN" = "null" ] || [ -z "$ADMIN_TOKEN" ]; then
    echo "❌ Failed to get admin token"
    exit 1
fi

# 1. Assign user to group
echo "👤 Assigning $USERNAME to $GROUP..."
curl -s -X POST http://localhost:8080/api/admin/roles \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"user\":\"$USERNAME\",\"role\":\"$GROUP\"}" > /dev/null

# 2. Get user token (assuming user exists in Casdoor)
echo "🔑 Getting user token..."
USER_TOKEN=$(curl -s -X POST "http://localhost:8000/api/login" \
  -d "{\"application\":\"app-built-in\",\"username\":\"$USERNAME\",\"password\":\"123456\",\"type\":\"token\"}" \
  -H "Content-Type: application/json" | jq -r .data)

if [ "$USER_TOKEN" = "null" ] || [ -z "$USER_TOKEN" ]; then
    echo "❌ Failed to get user token. User may not exist in Casdoor."
    echo "💡 Create user first or use existing user"
    exit 1
fi

echo "✅ Got user token"

# 3. Test access
echo ""
echo "🧪 Testing access with group: $GROUP"
echo "----------------------------------------"

# Test dashboard access (should work for all groups)
echo "📊 Testing dashboard access (/api/auth/me):"
RESPONSE=$(curl -s -w "%{http_code}" -X GET http://localhost:8080/api/auth/me \
  -H "Authorization: Bearer $USER_TOKEN")
HTTP_CODE="${RESPONSE: -3}"
if [ "$HTTP_CODE" = "200" ]; then
    echo "✅ Dashboard: ALLOWED"
else
    echo "❌ Dashboard: DENIED ($HTTP_CODE)"
fi

# Test transaction access
echo "💰 Testing transaction access (/api/transactions/my):"
RESPONSE=$(curl -s -w "%{http_code}" -X GET http://localhost:8080/api/transactions/my \
  -H "Authorization: Bearer $USER_TOKEN")
HTTP_CODE="${RESPONSE: -3}"
if [ "$HTTP_CODE" = "200" ]; then
    echo "✅ Transactions: ALLOWED"
else
    echo "❌ Transactions: DENIED ($HTTP_CODE)"
fi

# Test order access
echo "🛒 Testing order access (/api/orders/my):"
RESPONSE=$(curl -s -w "%{http_code}" -X GET http://localhost:8080/api/orders/my \
  -H "Authorization: Bearer $USER_TOKEN")
HTTP_CODE="${RESPONSE: -3}"
if [ "$HTTP_CODE" = "200" ]; then
    echo "✅ Orders: ALLOWED"
else
    echo "❌ Orders: DENIED ($HTTP_CODE)"
fi

# Test admin access
echo "👥 Testing admin access (/api/users):"
RESPONSE=$(curl -s -w "%{http_code}" -X GET http://localhost:8080/api/users \
  -H "Authorization: Bearer $USER_TOKEN")
HTTP_CODE="${RESPONSE: -3}"
if [ "$HTTP_CODE" = "200" ]; then
    echo "✅ Admin Users: ALLOWED"
else
    echo "❌ Admin Users: DENIED ($HTTP_CODE)"
fi

echo ""
echo "📋 Current role assignments for $USERNAME:"
curl -s -X GET http://localhost:8080/api/admin/roles \
  -H "Authorization: Bearer $ADMIN_TOKEN" | jq ".roles[] | select(.user == \"$USERNAME\")"

echo ""
echo "🎯 Expected access for $GROUP:"
case $GROUP in
    "dashboard_group")
        echo "✅ Dashboard | ❌ Transactions | ❌ Orders | ❌ Admin"
        ;;
    "transaction_group")
        echo "✅ Dashboard | ✅ Transactions | ❌ Orders | ❌ Admin"
        ;;
    "order_group")
        echo "✅ Dashboard | ❌ Transactions | ✅ Orders | ❌ Admin"
        ;;
    "full_access_group")
        echo "✅ Dashboard | ✅ Transactions | ✅ Orders | ❌ Admin"
        ;;
    "admin")
        echo "✅ Dashboard | ✅ Transactions | ✅ Orders | ✅ Admin"
        ;;
    *)
        echo "Unknown group: $GROUP"
        ;;
esac