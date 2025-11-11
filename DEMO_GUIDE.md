# 🎯 Casdoor Authentication Demo Guide

## 🚀 Quick Start

### 1. Start Services
```bash
# Start PostgreSQL + Casdoor
docker-compose up -d

# Start backend server
go run cmd/server/main.go
```

### 2. Open Demo
```
http://localhost:8080
```

## 📋 Demo Features

### **Case 1: Direct Login** 
- **Method**: Username/Password
- **Flow**: Frontend → Backend → Casdoor API
- **Default**: admin / 123456

### **Case 2: OAuth/OIDC Flow**
- **Method**: Standard OAuth 2.0
- **Flow**: Frontend → Casdoor Login Page → Backend Callback
- **Use case**: Web applications

### **Case 3: Microsoft SSO**
- **Method**: Azure AD integration
- **Flow**: Frontend → Casdoor → Microsoft → Backend Callback  
- **Use case**: Enterprise SSO

## 🔧 Test Scenarios

### **Authentication Tests:**
1. **Direct Login**: Test username/password
2. **OAuth Flow**: Test redirect-based login
3. **Microsoft SSO**: Test enterprise login

### **Authorization Tests:**
- **Protected Resource**: Basic authorization
- **User Profile**: User-specific data
- **Users List**: Admin-only endpoint
- **Secrets**: Certificate verification

### **Expected Results:**
- ✅ **Admin user**: Access all endpoints
- ❌ **Regular user**: Limited access (403 on admin endpoints)
- 🔐 **Authorization**: Casbin policies enforced

## 🎯 Demo Flow

### **1. Authentication Phase:**
```
User → Choose login method → Get JWT token → Show user info
```

### **2. Authorization Phase:**
```
JWT token → Test protected endpoints → Show results (200/403)
```

### **3. Policy Enforcement:**
```
Request → AuthMiddleware (JWT) → CasbinMiddleware (Policies) → Handler
```

## 📊 Expected Outputs

### **Successful Login:**
```json
{
  "id": "user-id",
  "name": "admin", 
  "email": "admin@example.com",
  "is_admin": true,
  "access_token": "eyJhbGciOiJSUzI1NiIs..."
}
```

### **Protected Resource (Success):**
```json
{
  "message": "This is a protected resource",
  "user": "admin",
  "data": "Sensitive data..."
}
```

### **Admin Endpoint (Forbidden):**
```json
{
  "message": "access denied"
}
```

## 🔍 Troubleshooting

### **Common Issues:**

**1. Microsoft SSO fails:**
- Check Azure Redirect URI: `http://localhost:8000/callback`
- Verify Casdoor provider config
- Check API permissions in Azure

**2. Authorization fails:**
- Check Casbin policies in DB
- Verify user roles assignment
- Test with admin user first

**3. CORS errors:**
- Server allows all origins for demo
- Check browser console for errors

### **Debug Commands:**
```bash
# Check Casdoor
curl http://localhost:8000/api/get-account

# Check backend health  
curl http://localhost:8080/health

# Check policies
curl http://localhost:8080/api/admin/policies \
  -H "Authorization: Bearer TOKEN"
```

## 🎉 Success Criteria

### **Demo is working when:**
- ✅ All 3 login methods work
- ✅ User info displays correctly
- ✅ Admin can access all endpoints
- ✅ Regular user gets 403 on admin endpoints
- ✅ JWT tokens are valid
- ✅ Casbin policies are enforced

## 📝 Notes

- **Frontend**: Pure HTML/JS (no framework)
- **Backend**: Go Echo with Casdoor + Casbin
- **Database**: PostgreSQL for Casbin policies
- **Authentication**: Casdoor (supports multiple providers)
- **Authorization**: Casbin RBAC model

**Ready to demo! 🚀**