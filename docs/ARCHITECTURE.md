# Architecture

qixue 是一个单体仓库：一个 Go API + 两个 Vue 前端。

```
浏览器 / App
   │
   ├─ qixue-admin  ──HTTP──►  /admin/*
   └─ qixue-web    ──HTTP──►  /webapp/*
                              │
                         qixue-api (Gin)
                              │
              ┌───────────────┼───────────────┐
              ▼               ▼               ▼
           MySQL            Redis         对象存储
```

## API 分层

```
internal/
├── admin/
│   ├── http/              # 管理端 handler（course、member、pms…）
│   ├── service/
│   └── repo/
├── webapp/
│   ├── http/              # 学员端 handler（auth、course、learn…）
│   ├── service/
│   └── repo/
├── service/               # 两端共用：auth、captcha、coin、enroll、cert、catalog、bunny、menutree …
├── repo/                  # 两端共用：user、课程报名、证书、积分、分类
├── model/
├── http/                  # 共用：middleware / bind / page / response
├── platform/              # config / db / cache / applog
└── server/
```

依赖方向只能向下：`handler → service → repo → model`。  
共用的 `internal/service`、`internal/repo` 不得 import `internal/admin` 或 `internal/webapp`。  
端内 service 不得依赖 `gin`，也不得 import 本端的 HTTP handler。

端内 service 的 **package 名**仍保留历史名（如 `admincourse`、`webauth`），以免和同端 HTTP 包重名。

## 路由

- 公开：登录验证码、系统信息等
- `/admin/*`：`AdminAuth`（管理员 JWT）
- `/webapp/*`：`WebappAuth`（学员 JWT）

## 约定

- Service 使用 `context.Context`，不依赖 `gin` / `http/response` / `middleware`
- HTTP 绑定与写响应只在 `admin` / `webapp` handler
- 分页等可复用 `http/page`；新代码优先在 handler 完成绑定后再传入 service
