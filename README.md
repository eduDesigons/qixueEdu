# 启学（qixue）

开源在线学习系统：学员端上课，管理后台维护课程与学员，Go API 同时服务两端。

[架构说明](docs/ARCHITECTURE.md) · [参与贡献](CONTRIBUTING.md) · [安全披露](SECURITY.md) · [License (MIT)](LICENSE)

## 目录

```
qixue/
├── qixue-api/      # Go API（Gin）
├── qixue-admin/    # 管理后台（Vue 3 + Vite + Naive UI）
├── qixue-web/      # 学员端（Vue 3 + Vite + Vant）
├── docs/           # 架构等文档
└── .github/        # CI
```

| 服务 | 本地地址 | 接口前缀 |
| --- | --- | --- |
| `qixue-api` | http://127.0.0.1:8787 | `/admin`、`/webapp` |
| `qixue-admin` | http://127.0.0.1:3200 | → API `/admin` |
| `qixue-web` | http://127.0.0.1:5180 | → API `/webapp` |

## 环境

- Go（版本见 `qixue-api/go.mod`）
- MySQL 8
- Redis
- Node.js 20+

## 快速开始

1. 复制 `qixue-api/.env.example` → `qixue-api/.env`，填写数据库、Redis、JWT、存储等。**不要提交 `.env`。**
2. 空库导入 `qixue-api/install/schema.sql` 与 `menus.sql`，再写入超级管理员（或从已有环境迁移）。
3. 证书字体在 `qixue-api/storage/fonts/`（缺字体会跳过文字绘制）。

### 启动 API

```bash
cd qixue-api
go run ./cmd/server
# 或：make check && make build && ./bin/server
```

默认 `:8787`（`GO_API_ADDR`）。排行榜快照 ticker：`RANK_SNAPSHOT_INTERVAL_SEC=0` 可关。

### 启动前端

```bash
cd qixue-admin && npm install && npm run dev   # :3200
cd qixue-web && npm install && npm run dev     # :5180
```

- 管理端接口：`qixue-admin/.env.development` → `VITE_AXIOS_BASE_URL`
- 学员端接口：`qixue-web/.env.development` → `VITE_BASE_URL`（含 `/webapp`）

学员账号由后台添加，无公开注册。本机访问不校验区域编码；区域域名访问时学员区域须匹配。

## 功能概览

- **学员端**：课程、系列课、课节、作业、证书、积分、学分、评论、排行、消息、反馈
- **管理端**：分类与课程、系列课、作业、证书、学员、报名、评论、公告/轮播、资源、区域、站点配置、积分规则、一言、管理员/角色/菜单

## 存储

`STORAGE_DRIVER=local` → `qixue-api/public`；`bunny`（默认）→ 对象存储（见 `.env.example`）。

## 生产与发版

```bash
cd qixue-api && make build
cd qixue-admin && npm run build
cd qixue-web && npm run build
```

发版读 `qixue-api/.env` 的 `STATIC_PUBLISH`（`bunny` | `local`）：

```bash
cd qixue-api
go run ./cmd/publish admin web api
```

`api` 产物默认在 `deploy/release/qixue-api/`（该目录不入库）。服务器日常只覆盖 `server`，保留线上 `.env`。细节见 [`qixue-api/README.md`](qixue-api/README.md)。

## 第三方来源

管理后台基于 [Vue Naive Admin](https://github.com/zclzone/vue-naive-admin)，其许可证见 `qixue-admin/LICENSE`。

## License

本仓库代码以 [MIT License](LICENSE) 发布（`qixue-admin` 上游组件另见其目录内 LICENSE）。
