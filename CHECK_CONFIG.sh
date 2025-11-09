#!/bin/bash

echo "=== Kiểm tra cấu hình Casdoor Integration ==="
echo ""

# Check if .env exists
if [ ! -f .env ]; then
    echo "❌ File .env không tồn tại!"
    echo "Tạo file .env với nội dung:"
    echo "CASDOOR_ENDPOINT=http://localhost:8000"
    echo "CASDOOR_CLIENT_ID=your_client_id"
    echo "CASDOOR_CLIENT_SECRET=your_client_secret"
    echo "CASDOOR_REDIRECT_URL=http://localhost:8080/api/auth/callback"
    exit 1
fi

echo "✅ File .env tồn tại"
echo ""

# Check if Casdoor is running
echo "Kiểm tra Casdoor đang chạy..."
if curl -s http://localhost:8000/api/get-account > /dev/null 2>&1; then
    echo "✅ Casdoor đang chạy tại http://localhost:8000"
else
    echo "❌ Casdoor không chạy. Hãy chạy: docker-compose up -d"
    exit 1
fi
echo ""

# Check OpenID Configuration
echo "Kiểm tra OpenID Configuration..."
AUTH_ENDPOINT=$(curl -s http://localhost:8000/.well-known/openid-configuration | grep -o '"authorization_endpoint":"[^"]*"' | cut -d'"' -f4)
TOKEN_ENDPOINT=$(curl -s http://localhost:8000/.well-known/openid-configuration | grep -o '"token_endpoint":"[^"]*"' | cut -d'"' -f4)

if [ -n "$AUTH_ENDPOINT" ]; then
    echo "✅ Authorization endpoint: $AUTH_ENDPOINT"
else
    echo "❌ Không tìm thấy authorization endpoint"
fi

if [ -n "$TOKEN_ENDPOINT" ]; then
    echo "✅ Token endpoint: $TOKEN_ENDPOINT"
else
    echo "❌ Không tìm thấy token endpoint"
fi
echo ""

# Check .env variables
echo "Kiểm tra biến môi trường trong .env..."
source .env 2>/dev/null || true

if [ -z "$CASDOOR_CLIENT_ID" ] || [ "$CASDOOR_CLIENT_ID" = "your_client_id" ] || [ "$CASDOOR_CLIENT_ID" = "app-123" ]; then
    echo "❌ CASDOOR_CLIENT_ID chưa được cấu hình đúng"
    echo "   Hãy lấy Client ID từ Casdoor: http://localhost:8000"
    echo "   Applications → Chọn application → Copy Client ID"
else
    echo "✅ CASDOOR_CLIENT_ID: $CASDOOR_CLIENT_ID"
fi

if [ -z "$CASDOOR_CLIENT_SECRET" ] || [ "$CASDOOR_CLIENT_SECRET" = "your_client_secret" ] || [ "$CASDOOR_CLIENT_SECRET" = "secret-456" ]; then
    echo "❌ CASDOOR_CLIENT_SECRET chưa được cấu hình đúng"
    echo "   Hãy lấy Client Secret từ Casdoor: http://localhost:8000"
    echo "   Applications → Chọn application → Copy Client Secret"
else
    echo "✅ CASDOOR_CLIENT_SECRET: [HIDDEN]"
fi

if [ -z "$CASDOOR_ENDPOINT" ]; then
    echo "❌ CASDOOR_ENDPOINT chưa được cấu hình"
else
    echo "✅ CASDOOR_ENDPOINT: $CASDOOR_ENDPOINT"
fi

if [ -z "$CASDOOR_REDIRECT_URL" ]; then
    echo "❌ CASDOOR_REDIRECT_URL chưa được cấu hình"
else
    echo "✅ CASDOOR_REDIRECT_URL: $CASDOOR_REDIRECT_URL"
fi
echo ""

# Test login endpoint
echo "Kiểm tra login endpoint..."
LOGIN_RESPONSE=$(curl -s http://localhost:8080/api/auth/login 2>/dev/null)
if [ -z "$LOGIN_RESPONSE" ]; then
    echo "❌ Server chưa chạy. Hãy chạy: go run cmd/server/main.go"
else
    if echo "$LOGIN_RESPONSE" | grep -q "client_id"; then
        CLIENT_ID_IN_URL=$(echo "$LOGIN_RESPONSE" | grep -o 'client_id=[^&]*' | cut -d'=' -f2)
        if [ -z "$CLIENT_ID_IN_URL" ] || [ "$CLIENT_ID_IN_URL" = "" ]; then
            echo "❌ Client ID rỗng trong login URL"
            echo "   Response: $LOGIN_RESPONSE"
        else
            echo "✅ Login URL có client_id: $CLIENT_ID_IN_URL"
        fi
    else
        echo "⚠️  Không tìm thấy client_id trong response"
        echo "   Response: $LOGIN_RESPONSE"
    fi
fi
echo ""

echo "=== Kết thúc kiểm tra ==="
echo ""
echo "Nếu có lỗi, xem TROUBLESHOOTING.md để biết cách sửa"

