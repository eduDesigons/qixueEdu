# 学员端

qixue 的学员站点。用来登录、看课、做作业、看证书和消息。技术栈是 Vue 3、Vite、Vant。

本地地址：http://127.0.0.1:5180 。接口前缀是 `/webapp`。

## 启动

先让 `qixue-api` 在 8787 端口运行。学员账号在管理后台添加，这里没有公开注册。

```bash
cp .env.example .env.development   # 若尚无本地配置
npm install
npm run dev
```

构建：

```bash
npm run build
```

构建产物含 PWA（`manifest.webmanifest` + Service Worker）。线上需 **HTTPS**（或 localhost）才能安装到桌面；有新版本时底部会提示「立即更新」，未安装时会提醒「安装到桌面」。

在线播放直接请求 CDN，Service Worker 不拦截 MP4。只有断网且已经缓存完整文件时，才改走本地播放。图片在线优先走网络，失败后再用缓存。

## 环境

开发配置在 `.env.development`：

| 变量 | 作用 |
| --- | --- |
| `VITE_BASE_URL` | 学员接口地址。本地是 `http://127.0.0.1:8787/webapp` |
| `VITE_ENV_NAME` | 环境名 |

没有单独的生产环境文件时，`src/config.js` 会把接口地址回退到 `/api/webapp`。上线前应写成实际地址，或让站点把 `/api` 转到接口根路径。

站点名称不写死在页面里，登录后从接口的站点配置读取。

请求由 `src/common/request.js` 发出。它会带上登录令牌，以及 `version-xcx`、`timestmp-xcx`、`random-xcx`、`sign-xcx`。这组签名是接口要求的，缺了会被拒绝。

本机打开不校验区域编码。用区域域名访问时，学员所属一级区域的编码必须和这个域名一致，否则登不进去。

## 目录

| 路径 | 内容 |
| --- | --- |
| `src/pages/index` | 登录、首页、学习、课程、我的 |
| `src/pages/course` | 课程详情、课节、系列课、评论 |
| `src/pages/exam` | 作业 |
| `src/pages/member` | 资料、证书、积分、学分、收藏、消息、反馈 |
| `src/pc` | 宽屏布局 |
| `src/common` | 请求、导航、存储、选图上传等 |
| `src/install.js` | 注册 Vant 与页面生命周期 |
| `src/store` | 登录态 |

窄屏走 `src/pages`，宽屏走 `src/pc`。改一个页面的数据或交互时，两边如果都有入口，要一起看。
