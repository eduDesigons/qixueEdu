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
├── server/     # 路由装配：router.go / router_admin.go / router_webapp.go
├── admin/      # 管理端 HTTP
├── webapp/     # 学员端 HTTP
├── service/    # 用例（无 gin）
├── repo/       # 数据访问
├── model/      # 表实体
├── http/       # middleware / bind / page / response
├── menutree/   # 菜单树（跨层共享，避免循环依赖）
└── platform/   # config / db / cache / applog
```

依赖方向只能向下：`handler → service → repo → model`。  
`service` 不得 import `admin` / `webapp` / `gin`。

## 路由

- 公开：登录验证码、`/install/index/status`、系统信息等
- `/admin/*`：`AdminAuth`（管理员 JWT）
- `/webapp/*`：`WebappAuth`（学员 JWT）

## 渐进整改

部分较老的 `service/admin*` 仍直接接收 `gin.Context`（历史遗留）。新代码与改动请按 `adminbanner` / `adminverse` / `adminfeedback` 的方式：service 返回数据或 `error`，由 handler 写响应。
