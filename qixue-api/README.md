# qixue-api

qixue HTTP API（[Gin](https://github.com/gin-gonic/gin)）。管理端前缀 `/admin`，学员端前缀 `/webapp`。默认监听 `:8787`。

完整项目说明见仓库根目录 [README](../README.md)，分层约定见 [docs/ARCHITECTURE.md](../docs/ARCHITECTURE.md)。

## Layout

```
internal/
├── admin/
│   ├── http/         # 管理端 handler
│   ├── service/
│   └── repo/
├── webapp/
│   ├── http/         # 学员端 handler
│   ├── service/
│   └── repo/
├── service/          # 共用：auth / captcha / coin / enroll / cert …
├── repo/             # 共用：user、报名、证书、积分、分类
├── model/
├── http/             # 共用：middleware / bind / page / response
├── platform/
└── server/
```

| Layer | May | Must not |
| --- | --- | --- |
| handler | bind、鉴权上下文、调 service、JSON | `Table` / 业务规则 |
| service | 规则、事务、调 repo。端内在 `admin/service`、`webapp/service` | 依赖 gin |
| repo | CRUD / 查询。端内在 `admin/repo`、`webapp/repo`，共用放 `repo` | 业务策略 |
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
make tidy    # go mod tidy
```

仓库根目录也可执行 `make api-check` / `make check`（见 [docs/DEVELOPMENT.md](../docs/DEVELOPMENT.md)）。

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
