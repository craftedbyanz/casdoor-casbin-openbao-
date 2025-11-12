RBAC Model (Role-Based Access Control):

Subject (sub): User/Role thực hiện hành động

Object (obj): Resource được truy cập (API endpoint, data)

Action (act): Hành động (read, write, delete)

Policy: Quy tắc phân quyền p, alice, data1, read

Role: Gán role cho user g, alice, admin

### Flow hoạt động

User login → Casdoor trả JWT với user info

Request API → AuthMiddleware verify JWT

CasbinMiddleware check permission: enforce(user, resource, action)

Nếu có quyền → Handler xử lý

Nếu không → 403 Forbidden
