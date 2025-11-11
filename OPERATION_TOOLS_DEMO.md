# 🛠️ Operation Tools Demo Guide

## 🎯 Overview

Complete operation management system with authentication, authorization, and resource ownership control.

## 🔐 Authentication Methods

### **Case 1: Direct Login**
- **Method**: Username/Password
- **Users**: admin/123456, hihi/password
- **Use case**: Internal systems

### **Case 2: Microsoft SSO**
- **Method**: Azure AD integration
- **Use case**: Enterprise environments
- **Flow**: Frontend → Casdoor → Microsoft → Backend

## 📊 Operation Modules

### **🏦 Transaction Management**
- **Purpose**: Financial transaction tracking
- **Features**: Create, view, ownership control
- **Data**: Amount, type, status, description

### **🛒 Order Management**
- **Purpose**: E-commerce order processing
- **Features**: Create, view, status updates
- **Data**: Product, quantity, price, status

## 🔒 Authorization Matrix

### **Admin User Permissions:**
```
✅ View all transactions
✅ View all orders
✅ Update order status
✅ View all users
✅ Access admin endpoints
✅ View own resources
✅ Create transactions/orders
```

### **Regular User Permissions:**
```
✅ View own transactions only
✅ View own orders only
✅ Create transactions/orders
✅ View own profile
❌ View all transactions (403)
❌ View all orders (403)
❌ Update order status (403)
❌ View all users (403)
❌ Access admin endpoints (403)
```

## 🧪 Test Scenarios

### **1. Authentication Testing:**
```bash
# Direct login
POST /api/auth/login
Body: {"username": "admin", "password": "123456"}

# Microsoft SSO
GET /api/auth/microsoft/login
→ Redirect to Microsoft login
```

### **2. Authorization Testing:**

**Admin User (should work):**
```bash
GET /api/transactions        # All transactions
GET /api/orders              # All orders
GET /api/users               # All users
PUT /api/orders/ord_001/status # Update order
```

**Regular User (mixed results):**
```bash
GET /api/transactions/my     # ✅ Own transactions
GET /api/orders/my           # ✅ Own orders
POST /api/transactions       # ✅ Create transaction
GET /api/transactions        # ❌ 403 Forbidden
GET /api/orders              # ❌ 403 Forbidden
```

### **3. Ownership Testing:**
```bash
# User A can access their transaction
GET /api/transactions/txn_001  # ✅ if owned by user A

# User A cannot access other's transaction
GET /api/transactions/txn_002  # ❌ 403 if owned by user B

# Admin can access any transaction
GET /api/transactions/txn_002  # ✅ Admin bypass ownership
```

## 🎮 Demo Flow

### **Step 1: Login**
1. Choose authentication method
2. Get JWT token
3. View user information

### **Step 2: Test Basic Endpoints**
1. Test protected resources
2. Test user profile
3. Test admin endpoints (if admin)

### **Step 3: Test Transaction Module**
1. View own transactions
2. Try to view all transactions
3. Create new transaction
4. Test ownership control

### **Step 4: Test Order Module**
1. View own orders
2. Try to view all orders
3. Create new order
4. Try to update order status (admin only)

## 📋 Expected Results

### **Admin User Results:**
```json
{
  "transactions": [...],     // ✅ All transactions
  "orders": [...],          // ✅ All orders
  "users": [...],           // ✅ All users
  "message": "Success"
}
```

### **Regular User Results:**
```json
// Own resources - Success
{
  "transactions": [...],     // ✅ Only user's transactions
  "orders": [...],          // ✅ Only user's orders
  "message": "Your data retrieved"
}

// Admin resources - Forbidden
{
  "message": "access denied" // ❌ 403 Forbidden
}
```

## 🔧 Technical Implementation

### **Authentication Flow:**
```
JWT Token → AuthMiddleware → Extract user info → Context
```

### **Authorization Flow:**
```
Context → CasbinMiddleware → Check policies → Allow/Deny
```

### **Ownership Flow:**
```
Handler → Check resource ownership → Allow own resources only
```

### **Policy Examples:**
```
# Admin policies
p, admin, /api/transactions, read
p, admin, /api/orders, read

# User policies  
p, user, /api/transactions/my, read
p, user, /api/orders/my, read

# Role assignments
g, admin, admin
g, hihi, user
```

## 🚀 Quick Start

```bash
# 1. Start services
docker-compose up -d
go run cmd/server/main.go

# 2. Open demo
http://localhost:8080

# 3. Login and test all modules
```

## 🎯 Success Criteria

- ✅ Both authentication methods work
- ✅ Admin can access all resources
- ✅ Users can only access own resources
- ✅ Ownership control works correctly
- ✅ Casbin policies are enforced
- ✅ 403 errors for unauthorized access

**Complete operation tools with full RBAC implementation! 🎉**