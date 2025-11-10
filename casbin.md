RBAC Model (Role-Based Access Control):

Subject (sub): User/Role thực hiện hành động

Object (obj): Resource được truy cập (API endpoint, data)

Action (act): Hành động (read, write, delete)

Policy: Quy tắc phân quyền p, alice, data1, read

Role: Gán role cho user g, alice, admin

```
internal/
├── casbin/
│   ├── enforcer.go     # Khởi tạo Casbin enforcer
│   ├── middleware.go   # Authorization middleware
│   ├── policy.go       # Quản lý policies
│   └── service.go      # Business logic
├── database/
│   └── postgres.go     # Database connection
└── handler/
    └── admin.go        # Admin APIs cho policy management
```

### Flow hoạt động

User login → Casdoor trả JWT với user info

Request API → AuthMiddleware verify JWT

CasbinMiddleware check permission: enforce(user, resource, action)

Nếu có quyền → Handler xử lý

Nếu không → 403 Forbidden
