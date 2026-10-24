# Architecture

启学是一个单体仓库：一个 Go API + 两个 Vue 前端。

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
├── server/          # router.go / router_admin.go / router_webapp.go
├── admin/           # 管理端 HTTP
├── webapp/          # 学员端 HTTP
├── service/
│   ├── admin/       # 管理端用例（banner、course、member…）
│   ├── webapp/      # 学员端用例（auth、course、home…）
│   ├── auth/        # 共享：JWT、限流、密码
│   ├── captcha/
│   ├── coin/
│   ├── enroll/
│   ├── catalog/
│   ├── bunny/
│   ├── menutree/    # 菜单树构建
│   └── …            # 其它跨端能力
├── repo/
├── model/
├── http/            # middleware / bind / page / response
└── platform/        # config / db / cache / applog
```

依赖方向只能向下：`handler → service → repo → model`。  
`service` 不得 import `admin` / `webapp`（HTTP）/ `gin`。

说明：`service/admin/*`、`service/webapp/*` 的 **package 名**仍保留历史名（如 `admincourse`、`webauth`），以免与 HTTP 层同名包冲突；**目录**按端划分即可。

## 路由

- 公开：登录验证码、`/install/index/status`、系统信息等
- `/admin/*`：`AdminAuth`（管理员 JWT）
- `/webapp/*`：`WebappAuth`（学员 JWT）

## 约定

- Service 使用 `context.Context`，不依赖 `gin` / `http/response` / `middleware`
- HTTP 绑定与写响应只在 `admin` / `webapp` handler
- 分页等可复用 `http/page`；新代码优先在 handler 完成绑定后再传入 service
