# qixue-api

启学 HTTP API（[Gin](https://github.com/gin-gonic/gin)）。管理端前缀 `/admin`，学员端前缀 `/webapp`。默认监听 `:8787`。

完整项目说明见仓库根目录 [README](../README.md)，分层约定见 [docs/ARCHITECTURE.md](../docs/ARCHITECTURE.md)。

## Layout

```
internal/
├── model/      # GORM entities
├── repo/       # data access
├── service/    # use-cases（勿依赖 gin）
├── admin/      # 管理端 HTTP → service
├── webapp/     # 学员端 HTTP → service
├── menutree/   # 菜单树（跨层共享）
├── http/       # middleware / bind / page / response
├── platform/   # config / db / cache / applog
└── server/     # router.go + router_admin.go + router_webapp.go
```

| Layer | May | Must not |
| --- | --- | --- |
| handler | bind、鉴权上下文、调 service、JSON | `Table` / 业务规则 |
| service | 规则、事务、调 repo | 依赖 gin |
| repo | CRUD / 查询 | 业务策略 |
| model | 列 + `TableName` | 胖领域方法 |

依赖在 `server` 中手工装配，无 DI 框架。

## Run

```bash
cp .env.example .env
go run ./cmd/server
```

```bash
make check   # go vet + go test
make build   # → bin/server
```

字体：`storage/fonts/`。空库导入 `install/schema.sql` 与 `install/menus.sql`。

## Publish

```bash
go run ./cmd/publish admin
go run ./cmd/publish web
go run ./cmd/publish api
```

或 `publish.bat`。`STATIC_PUBLISH`：`bunny` | `local`。

`api` 默认 `linux/amd64`（`GO_API_GOOS` / `GO_API_GOARCH`）。产物：`deploy/release/qixue-api/`。

线上默认目录 `/opt/qixue-api`：只替换 `server`，保留服务器 `.env`。

日志：`logs/qixue-api-YYYY-MM-DD.log`（`LOG_KEEP_DAYS`）。

## Rank snapshots

进程内 ticker 写 `rank_snapshots`。`RANK_SNAPSHOT_INTERVAL_SEC=0` 关闭。
