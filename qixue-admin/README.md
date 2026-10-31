# 管理后台

qixue 的管理端。用来维护课程、学员、装修、资源和站点配置。界面基于 [Vue Naive Admin](https://github.com/zclzone/vue-naive-admin)，接口是同仓库的 `qixue-api`，不是模板自带的 NestJS 服务。

本地地址：http://127.0.0.1:3200 。路由是 hash，登录页是 http://127.0.0.1:3200/#/login 。

## 启动

先让 `qixue-api` 在 8787 端口运行。然后：

```bash
cp .env.example .env.development   # 若尚无本地配置
npm install
npm run dev
```

用库里已有的超级管理员登录。

构建：

```bash
npm run build
```

构建产物含 PWA（`manifest.webmanifest` + Service Worker）。线上需 **HTTPS**（或 localhost）才能安装到桌面；有新版本时底部会提示「立即更新」，未安装时会提醒「安装到桌面」。管理端不做视频离线缓存，视频始终直连 CDN。

## 环境

开发配置在 `.env.development`：

| 变量 | 作用 |
| --- | --- |
| `VITE_AXIOS_BASE_URL` | 接口根地址。本地是 `http://127.0.0.1:8787` |
| `VITE_USE_HASH` | `true` 时使用 hash 路由 |
| `VITE_PUBLIC_PATH` | 静态资源公共路径，需以 `/` 开头和结尾 |

生产配置在 `.env.production`。当前 `VITE_AXIOS_BASE_URL` 是 `/api`。接口自己的路径是 `/admin`，没有 `/api` 前缀，所以线上要把 `/api` 转到接口根路径。

页面请求走 `src/utils/http`。返回码不是 200 时会自动提示，业务代码里不必再弹一次。

## 菜单

登录后菜单和按钮来自接口的 `administrator_menus`，再按角色授权。超级管理员角色编码是 `SUPER_ADMIN`。新增页面要在「管理 → 菜单」里登记 `path` 和 `component`，`component` 对应 `src/views` 下的文件，例如 `/src/views/pms/area/index.vue`。

不需要权限控制的页面放在 `src/router/basic-routes.js`。

## 目录

| 路径 | 内容 |
| --- | --- |
| `src/views/course` | 课程、分类、课节、系列课、作业 |
| `src/views/member` | 学员、证书、报名、反馈、消息 |
| `src/views/ops` | 学习、报名、积分、学分、证书等运营数据 |
| `src/views/decoration` | 公告、轮播 |
| `src/views/media` | 视频、音频、图片、文档 |
| `src/views/pms` | 管理员、角色、菜单、配置、区域、一言 |
| `src/views/comment` | 课程评论、课节评论 |
| `src/api` | 按业务拆开的接口 |
| `src/store` | 登录态、权限、页签 |

站点名称、页脚、积分发放数额在「配置」里改，由接口下发，不要写死在页面上。

## 版权

本目录的界面代码沿用 Vue Naive Admin 的 MIT 许可证，见 `LICENSE`。
