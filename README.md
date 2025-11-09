```
┌──────────────┐
│   Frontend   │
│ (Web/API)    │
└──────┬───────┘
       │ Login via OIDC (Casdoor)
       ▼
┌──────────────┐
│   Casdoor    │
│ (AuthN IdP)  │
│ ├─ User, Org │
│ ├─ Role      │
│ └─ JWT issue │
└──────┬───────┘
       │ JWT
       ▼
┌──────────────┐
│   Backend    │
│ (Golang)     │
│ ├─ Verify JWT │
│ ├─ Casbin     │
│ │  (AuthZ)   │
│ └─ Use Vault │
│    (Secrets) │
└──────┬───────┘
       │ Policy Enforcement
       ▼
┌──────────────┐
│ PostgreSQL   │
│ (data store) │
└──────────────┘
```

### Casdoor — Authentication (OIDC Provider)

1. Run casdoor

```
docker run -d --name casdoor -p 8000:8000 casbin/casdoor
```

2. Tạp Application

```
http://localhost:8000, login admin → Applications → Add

Điền:
    - Redirect URL: http://localhost:8080/api/callback
    - Client ID, Secret: (ghi lại để backend dùng)
```

3. Lấy public key
   Casdoor có endpoint để verify JWT

```
GET /api/certs
```

> Cần public key này trong backend.

### OpenBao (Vault) — Secrets Management

Dùng để lưu:

    - CASDOOR_CLIENT_ID
    - CASDOOR_CLIENT_SECRET
    - JWT_PUBLIC_KEY
    - DB password
    - Casbin adapter credentials…

1. Run OpenBao (giống Vault)

```
docker run -d --name openbao -p 8200:8200 openbao/openbao
```

2. Lưu secrets

```
vault kv put secret/casdoor client_id="xxx" client_secret="yyy"
vault kv put secret/jwt public_key="-----BEGIN PUBLIC KEY-----..."
vault kv put secret/db dsn="postgres://user:pass@localhost/db"
```

### Casbin — Authorization Engine

Casbin local trong backend, dùng DB adapter (vd: PostgreSQL).

Cài đặt adapter:

```
go get github.com/casbin/xorm-adapter/v3
go get github.com/casbin/casbin/v2
```

Cấu hình model (file rbac_model.conf):

```
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
```

```
[Client] → gửi request kèm Bearer token → [Echo Backend]
          → AuthMiddleware:
               - Casdoor SDK xác thực JWT (AuthN)
               - Casbin kiểm tra policy (AuthZ)
          → Controller xử lý logic
```

```
{
  "status": "ok",
  "msg": "",
  "sub": "",
  "name": "",
  "data": "eyJhbGciOiJSUzI1NiIsImtpZCI6ImNlcnQtYnVpbHQtaW4iLCJ0eXAiOiJKV1QifQ.eyJvd25lciI6ImJ1aWx0LWluIiwibmFtZSI6ImFkbWluIiwiY3JlYXRlZFRpbWUiOiIyMDI1LTExLTA4VDA4OjUxOjIxWiIsInVwZGF0ZWRUaW1lIjoiMjAyNS0xMS0wOFQwOToxMDozM1oiLCJkZWxldGVkVGltZSI6IiIsImlkIjoiNmJjNTEwMzgtYzUxMy00NDg0LTg2NWYtOGEwZjk2MDQ2MDNlIiwidHlwZSI6Im5vcm1hbC11c2VyIiwicGFzc3dvcmQiOiIiLCJwYXNzd29yZFNhbHQiOiIiLCJwYXNzd29yZFR5cGUiOiJwbGFpbiIsImRpc3BsYXlOYW1lIjoiQWRtaW4iLCJmaXJzdE5hbWUiOiIiLCJsYXN0TmFtZSI6IiIsImF2YXRhciI6Imh0dHBzOi8vY2RuLmNhc2Jpbi5vcmcvaW1nL2Nhc2Jpbi5zdmciLCJhdmF0YXJUeXBlIjoiIiwicGVybWFuZW50QXZhdGFyIjoiIiwiZW1haWwiOiJhZG1pbkBleGFtcGxlLmNvbSIsImVtYWlsVmVyaWZpZWQiOmZhbHNlLCJwaG9uZSI6IjEyMzQ1Njc4OTEwIiwiY291bnRyeUNvZGUiOiJVUyIsInJlZ2lvbiI6IiIsImxvY2F0aW9uIjoiIiwiYWRkcmVzcyI6W10sImFmZmlsaWF0aW9uIjoiRXhhbXBsZSBJbmMuIiwidGl0bGUiOiIiLCJpZENhcmRUeXBlIjoiIiwiaWRDYXJkIjoiIiwiaG9tZXBhZ2UiOiIiLCJiaW8iOiIiLCJsYW5ndWFnZSI6IiIsImdlbmRlciI6IiIsImJpcnRoZGF5IjoiIiwiZWR1Y2F0aW9uIjoiIiwic2NvcmUiOjIwMDAsImthcm1hIjowLCJyYW5raW5nIjoxLCJpc0RlZmF1bHRBdmF0YXIiOmZhbHNlLCJpc09ubGluZSI6ZmFsc2UsImlzQWRtaW4iOnRydWUsImlzRm9yYmlkZGVuIjpmYWxzZSwiaXNEZWxldGVkIjpmYWxzZSwic2lnbnVwQXBwbGljYXRpb24iOiJhcHAtYnVpbHQtaW4iLCJoYXNoIjoiIiwicHJlSGFzaCI6IiIsInJlZ2lzdGVyVHlwZSI6IkFkZCBVc2VyIiwicmVnaXN0ZXJTb3VyY2UiOiJidWlsdC1pbi9hZG1pbiIsImFjY2Vzc0tleSI6IiIsImFjY2Vzc1NlY3JldCI6IiIsImdpdGh1YiI6IiIsImdvb2dsZSI6IiIsInFxIjoiIiwid2VjaGF0IjoiIiwiZmFjZWJvb2siOiIiLCJkaW5ndGFsayI6IiIsIndlaWJvIjoiIiwiZ2l0ZWUiOiIiLCJsaW5rZWRpbiI6IiIsIndlY29tIjoiIiwibGFyayI6IiIsImdpdGxhYiI6IiIsImNyZWF0ZWRJcCI6IjEyNy4wLjAuMSIsImxhc3RTaWduaW5UaW1lIjoiIiwibGFzdFNpZ25pbklwIjoiIiwicHJlZmVycmVkTWZhVHlwZSI6IiIsInJlY292ZXJ5Q29kZXMiOm51bGwsInRvdHBTZWNyZXQiOiIiLCJtZmFQaG9uZUVuYWJsZWQiOmZhbHNlLCJtZmFFbWFpbEVuYWJsZWQiOmZhbHNlLCJsZGFwIjoiIiwicHJvcGVydGllcyI6e30sInJvbGVzIjpbXSwicGVybWlzc2lvbnMiOltdLCJncm91cHMiOltdLCJsYXN0U2lnbmluV3JvbmdUaW1lIjoiIiwic2lnbmluV3JvbmdUaW1lcyI6MCwibWFuYWdlZEFjY291bnRzIjpudWxsLCJ0b2tlblR5cGUiOiJhY2Nlc3MtdG9rZW4iLCJ0YWciOiJzdGFmZiIsImF6cCI6ImVhNTI1YzE5ZjZmNzVjMmY4NDE5IiwiaXNzIjoiaHR0cDovL2xvY2FsaG9zdDo4MDAwIiwic3ViIjoiNmJjNTEwMzgtYzUxMy00NDg0LTg2NWYtOGEwZjk2MDQ2MDNlIiwiYXVkIjpbImVhNTI1YzE5ZjZmNzVjMmY4NDE5Il0sImV4cCI6MTc2MzE5NzkwOSwibmJmIjoxNzYyNTkzMTA5LCJpYXQiOjE3NjI1OTMxMDksImp0aSI6ImFkbWluLzBiZWZiYTZmLThhMDgtNDhlMC04ZjlhLWE0MmY2ZmRmYmM4NCJ9.HfRJ6DpnuK1oVadkkU6AdS6QlQNz6OdIpQ-D_YYyGdxQOV6G17zx2qv8IxOBzY4xnXEmFzoWPLv5UXab2FVsCEWeBV5l0TXBl1XZTtdY8Al0kpnflMGGEV62flKSSi9DozjfX2LpCbFFqCk5iQdvmFy8MJWmofOQSk7RJOmfr19UvwjW0QUF0mWxDIePOu8swz07CV-8fsYNMgFYuVq5BN_IsIdGwN6tn4n7aklksGSYJZbFTo7FMKLcu4gQRnjtnULvaSVRJFfNMGDzSXwk0byrjmYAfB_8F6demTZf3VwMFX74RjrxC2G98ahn9-vm9sRY7F9_uqOGoiA1D9A-5DZ14BcNpxkXjwhOkA0RVfZ1NwhiCc8qKnmT802rX2AQ_im33RGNI-Bg82Lm2MDFBMkwp7VEsKEJD81H_ujRqZEwVpj38mx2ZFds0PDkGCFtnSkVuIQp3LwKqhzrbroZrzjRs-e4SGrGauXGmGwOz-_vtExfIwf5cpwbWoTpshneIo9tWdmScJ_rDZTcwhXQ2Bx9jZj5xranhY7mV8rBD3rzytfFEPmddQR4CK96X9J6_12TUwy46zrGdj5XmOHhe-OVWz2NbOA7SX2krHAGYkOG5VEHaXK74VCrXqUDmwpIebsIAieqRs2ycyhFdZfeII5Q3CbuwIW1B8RUWejqtUE",
  "data2": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJvd25lciI6ImJ1aWx0LWluIiwibmFtZSI6ImFkbWluIiwiY3JlYXRlZFRpbWUiOiIyMDI1LTExLTA4VDA4OjUxOjIxWiIsInVwZGF0ZWRUaW1lIjoiMjAyNS0xMS0wOFQwOToxMDozM1oiLCJkZWxldGVkVGltZSI6IiIsImlkIjoiNmJjNTEwMzgtYzUxMy00NDg0LTg2NWYtOGEwZjk2MDQ2MDNlIiwidHlwZSI6Im5vcm1hbC11c2VyIiwicGFzc3dvcmQiOiIiLCJwYXNzd29yZFNhbHQiOiIiLCJwYXNzd29yZFR5cGUiOiJwbGFpbiIsImRpc3BsYXlOYW1lIjoiQWRtaW4iLCJmaXJzdE5hbWUiOiIiLCJsYXN0TmFtZSI6IiIsImF2YXRhciI6Imh0dHBzOi8vY2RuLmNhc2Jpbi5vcmcvaW1nL2Nhc2Jpbi5zdmciLCJhdmF0YXJUeXBlIjoiIiwicGVybWFuZW50QXZhdGFyIjoiIiwiZW1haWwiOiJhZG1pbkBleGFtcGxlLmNvbSIsImVtYWlsVmVyaWZpZWQiOmZhbHNlLCJwaG9uZSI6IjEyMzQ1Njc4OTEwIiwiY291bnRyeUNvZGUiOiJVUyIsInJlZ2lvbiI6IiIsImxvY2F0aW9uIjoiIiwiYWRkcmVzcyI6W10sImFmZmlsaWF0aW9uIjoiRXhhbXBsZSBJbmMuIiwidGl0bGUiOiIiLCJpZENhcmRUeXBlIjoiIiwiaWRDYXJkIjoiIiwiaG9tZXBhZ2UiOiIiLCJiaW8iOiIiLCJsYW5ndWFnZSI6IiIsImdlbmRlciI6IiIsImJpcnRoZGF5IjoiIiwiZWR1Y2F0aW9uIjoiIiwic2NvcmUiOjIwMDAsImthcm1hIjowLCJyYW5raW5nIjoxLCJpc0RlZmF1bHRBdmF0YXIiOmZhbHNlLCJpc09ubGluZSI6ZmFsc2UsImlzQWRtaW4iOnRydWUsImlzRm9yYmlkZGVuIjpmYWxzZSwiaXNEZWxldGVkIjpmYWxzZSwic2lnbnVwQXBwbGljYXRpb24iOiJhcHAtYnVpbHQtaW4iLCJoYXNoIjoiIiwicHJlSGFzaCI6IiIsInJlZ2lzdGVyVHlwZSI6IkFkZCBVc2VyIiwicmVnaXN0ZXJTb3VyY2UiOiJidWlsdC1pbi9hZG1pbiIsImFjY2Vzc0tleSI6IiIsImFjY2Vzc1NlY3JldCI6IiIsImdpdGh1YiI6IiIsImdvb2dsZSI6IiIsInFxIjoiIiwid2VjaGF0IjoiIiwiZmFjZWJvb2siOiIiLCJkaW5ndGFsayI6IiIsIndlaWJvIjoiIiwiZ2l0ZWUiOiIiLCJsaW5rZWRpbiI6IiIsIndlY29tIjoiIiwibGFyayI6IiIsImdpdGxhYiI6IiIsImNyZWF0ZWRJcCI6IjEyNy4wLjAuMSIsImxhc3RTaWduaW5UaW1lIjoiIiwibGFzdFNpZ25pbklwIjoiIiwicHJlZmVycmVkTWZhVHlwZSI6IiIsInJlY292ZXJ5Q29kZXMiOm51bGwsInRvdHBTZWNyZXQiOiIiLCJtZmFQaG9uZUVuYWJsZWQiOmZhbHNlLCJtZmFFbWFpbEVuYWJsZWQiOmZhbHNlLCJsZGFwIjoiIiwicHJvcGVydGllcyI6e30sInJvbGVzIjpbXSwicGVybWlzc2lvbnMiOltdLCJncm91cHMiOltdLCJsYXN0U2lnbmluV3JvbmdUaW1lIjoiIiwic2lnbmluV3JvbmdUaW1lcyI6MCwibWFuYWdlZEFjY291bnRzIjpudWxsLCJ0b2tlblR5cGUiOiJyZWZyZXNoLXRva2VuIiwidGFnIjoic3RhZmYiLCJhenAiOiJlYTUyNWMxOWY2Zjc1YzJmODQxOSIsImlzcyI6Imh0dHA6Ly9sb2NhbGhvc3Q6ODAwMCIsInN1YiI6IjZiYzUxMDM4LWM1MTMtNDQ4NC04NjVmLThhMGY5NjA0NjAzZSIsImF1ZCI6WyJlYTUyNWMxOWY2Zjc1YzJmODQxOSJdLCJleHAiOjE3NjMxOTc5MDksIm5iZiI6MTc2MjU5MzEwOSwiaWF0IjoxNzYyNTkzMTA5LCJqdGkiOiJhZG1pbi8wYmVmYmE2Zi04YTA4LTQ4ZTAtOGY5YS1hNDJmNmZkZmJjODQifQ.oZN1LT48szGZDBJBuValBrSwML6nDd1caMRnTBvOpDy2lobhMa4NS81MaCQ-bQJrMvyQ0s6lJMYV_S4rYsk-Jts1c-qZR6kk7kvpQx1I18esvNcMLbhCTBa6pYtg17ZXRDJfpMcQb7Wp5EeLz2N4QCnpUzi5mj9WipqFxuUBzJEx0u3eHki3P94dsdsSBNWJvIQMaJRos_DrE4DViBrWpXBoF_Cr5Iy7R3CPC_I-_pjANnKYFoAibcp4mBeqa6CqfSvXPt7YMlPOI2ui4fRzbEZtpygmzVp792F4VgKrQruJh8tiZZXrjDDcjAI4m4evfKIiAt3Ewgj5GSzW2oPJ03xYFen54GnYDpDWYnelBj5cJUZ7-CrwqB2JGCc2Y_eNCEqpbD9mtD_806dZQ7FhKs5woT7_WizyQ52CJUcsQSk4WQbD34FBc-jfgqrFvgWCm-1Fe3UNXsbedkmOYf9ONQcBH3UHU3pNUc_WFSnskLIGI0GRPvjLrCGiLkLjQ8S8--2DM3c2mcee6PxkNz_cL-83X5HC9jqGe_-rt8LSAAUBkSv2q9GxpbDJj_FSIMIX6-i_bgfL9fbhpRN8VQ39qlvhBvY2I3VpAm9Va-CPKFGvnnYslUZLhtE6NlfkNDu9597FDGyMNXTUffZmsxKZ6OYX__YY9NAG3i6Cw8Vy5Gg",
  "data3": false
}
```

```
# 1. Đăng nhập admin → lấy token
ADMIN_TOKEN=$(curl -s -X POST "http://localhost:8000/api/login" \
  -d '{"application":"app-built-in","username":"admin","password":"123456","type":"token"}' \
  -H "Content-Type: application/json" | jq -r .data)

# 2. Tạo user mới
curl -X POST "http://localhost:8000/api/add-user" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d '{"owner":"org_anz","name":"testuser","email":"testuser@anz.com","password":"pass123"}' \
  -H "Content-Type: application/json"

# 3. Login user → lấy code từ browser → đổi lấy token
# (mở link, copy code, chạy lệnh)

# 4. Gọi API với token user
curl -X GET "http://localhost:8000/api/get-account" -H "Authorization: Bearer $USER_TOKEN"
```


### Cert
```
```
