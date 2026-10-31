# 本地开发

## 前置

| 依赖 | 说明 |
| --- | --- |
| Go | 版本见 `qixue-api/go.mod` |
| Node.js | 20+（管理端 / 学员端） |
| MySQL | 8.x |
| Redis | 任意近期稳定版 |

## 一键检查（仓库根）

```bash
make check      # API vet+test + 两端 build
make build      # 仅构建三端产物
make api-check  # 仅 API
```

Windows 若无 GNU make，可进入子目录执行各自 README 中的命令。

## API（`qixue-api`）

```bash
cp .env.example .env   # 填写 DB / Redis / JWT 等
# 空库：install/schema.sql + install/menus.sql，再写入超级管理员
go run ./cmd/server    # 默认 :8787
make check             # vet + test
make build             # → bin/server
```

关键变量见 `.env.example`。不要提交 `.env`。

分层：`handler → service → repo → model`，详见 [ARCHITECTURE.md](ARCHITECTURE.md)。

## 管理端（`qixue-admin`）

```bash
cp .env.example .env.development   # 若尚无本地文件
npm install
npm run dev    # :3200
npm run lint   # eslint
npm run build
```

`VITE_AXIOS_BASE_URL` 指向 API 根（本地一般为 `http://127.0.0.1:8787`）。

## 学员端（`qixue-web`）

```bash
cp .env.example .env.development   # 若尚无本地文件
npm install
npm run dev    # :5180
npm run build
```

`VITE_BASE_URL` 须包含 `/webapp`（本地一般为 `http://127.0.0.1:8787/webapp`）。

学员账号由管理后台添加，无公开注册。

## 提交前

1. API 改动：`cd qixue-api && make check`
2. 前端改动：对应目录 `npm run build`（管理端建议再跑 `npm run lint`）
3. 不提交密钥、`deploy/release`、本机二进制与日志
4. 行为或接口变更时更新相关 README / docs
