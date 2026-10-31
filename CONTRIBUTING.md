# Contributing

感谢参与 qixue。提交前请先本地跑通相关检查。

## 行为准则

参与本仓库即表示你同意遵守 [CODE_OF_CONDUCT.md](CODE_OF_CONDUCT.md)。

## 仓库结构

| 目录 | 说明 |
| --- | --- |
| `qixue-api/` | Go API（Gin） |
| `qixue-admin/` | 管理后台（Vue 3） |
| `qixue-web/` | 学员端（Vue 3） |
| `docs/` | 架构与开发文档 |

本地环境与命令：[docs/DEVELOPMENT.md](docs/DEVELOPMENT.md)。  
分层约定：[docs/ARCHITECTURE.md](docs/ARCHITECTURE.md)。

## 开发约定

### API（`qixue-api`）

分层：`handler → service → repo → model`。

- **handler**（`internal/admin/http`、`internal/webapp/http`）：绑定参数、鉴权上下文、调用 service、写 JSON。不要写业务规则或直接 `Table(...)`。共用的 middleware、bind、page、response 放 `internal/http`。
- **service**：业务与事务；入参用 `context.Context`，**不要**依赖 `gin.Context`。管理端放 `internal/admin/service/`，学员端放 `internal/webapp/service/`，两端共用放 `internal/service/`。
- **repo**：查询与持久化；优先 GORM model。管理端 `internal/admin/repo`，学员端 `internal/webapp/repo`，共用放 `internal/repo`。课程评论与课节评论分表分模型。
- **model**：表结构与 `TableName`。

```bash
cd qixue-api
make check   # fmt 可选；vet + test
```

### 前端

```bash
cd qixue-admin && npm install && npm run lint && npm run build
cd qixue-web && npm install && npm run build
```

环境变量模板：各目录下的 `.env.example`。

## Issue / Pull Request

1. Bug / 需求优先用 Issue 模板；安全问题走 [SECURITY.md](SECURITY.md)，不要公开贴利用细节。
2. PR 请填写模板：动机、范围、如何验证。
3. 不提交 `.env`、密钥、本机 `deploy/release` 产物、二进制与日志。
4. API 行为变更时同步更新 `qixue-api/README.md`、根 README 或 `docs/`。
5. 保持 diff 聚焦，避免无关格式化大扫荡。
6. 仓库根可用 `make check` 做集成检查（需 Node + Go）。

## 安全问题

请勿在公开 Issue 中贴出可利用细节，见 [SECURITY.md](SECURITY.md)。
