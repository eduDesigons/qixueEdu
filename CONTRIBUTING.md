# Contributing

感谢参与启学（qixue）。提交前请先本地跑通相关检查。

## 仓库结构

| 目录 | 说明 |
| --- | --- |
| `qixue-api/` | Go API（Gin） |
| `qixue-admin/` | 管理后台（Vue 3） |
| `qixue-web/` | 学员端（Vue 3） |

## 开发约定

### API（`qixue-api`）

分层：`handler → service → repo → model`。

- **handler**（`internal/admin`、`internal/webapp`）：绑定参数、鉴权上下文、调用 service、写 JSON。不要写业务规则或直接 `Table(...)`。
- **service**：业务与事务；入参用 `context.Context`，**不要**依赖 `gin.Context`。
- **repo**：查询与持久化。
- **model**：表结构与 `TableName`。

参考已对齐的包：`service/adminbanner`、`service/adminverse`、`service/adminfeedback`。

```bash
cd qixue-api
make check   # fmt + vet + test
```

### 前端

```bash
cd qixue-admin && npm install && npm run lint
cd qixue-web && npm install && npm run build
```

## Pull Request

1. 说明动机与改动范围。
2. 不提交 `.env`、密钥、本机 `deploy/release` 产物。
3. API 行为变更时同步更新 `qixue-api/README.md` 或根 README。
4. 保持 diff 聚焦，避免无关格式化大扫荡。

## 安全问题

请勿在公开 issue 中贴出可利用细节，见 [SECURITY.md](SECURITY.md)。
