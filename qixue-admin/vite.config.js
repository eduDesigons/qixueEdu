/**********************************
 * @Author: Ronnie Zhang
 * @LastEditor: Ronnie Zhang
 * @LastEditTime: 2023/12/05 21:31:02
 * @Email: zclzone@outlook.com
 * Copyright © 2023 Ronnie Zhang(大脸怪) | https://isme.top
 **********************************/

import path from 'node:path'
import Vue from '@vitejs/plugin-vue'
import VueJsx from '@vitejs/plugin-vue-jsx'
import Unocss from 'unocss/vite'
import AutoImport from 'unplugin-auto-import/vite'
import { NaiveUiResolver } from 'unplugin-vue-components/resolvers'
import Components from 'unplugin-vue-components/vite'
import { defineConfig, loadEnv } from 'vite'
import removeNoMatch from 'vite-plugin-router-warn'
import { pluginIcons, pluginPagePathes } from './build/plugin-isme'

export default defineConfig(async ({ mode }) => {
  const viteEnv = loadEnv(mode, process.cwd())
  const { VITE_PUBLIC_PATH, VITE_PROXY_TARGET } = viteEnv

  const plugins = [
    Vue(),
    VueJsx(),
  ]

  // 只在开发环境动态导入 VueDevTools
  if (mode === 'development') {
    try {
      const VueDevTools = (await import('vite-plugin-vue-devtools')).default
      plugins.push(VueDevTools())
    } catch (e) {
      // 如果导入失败，忽略错误
      console.warn('Failed to load VueDevTools:', e.message)
    }
  }

  plugins.push(
    Unocss(),
    AutoImport({
      imports: ['vue', 'vue-router'],
      dts: false,
    }),
    Components({
      resolvers: [NaiveUiResolver()],
      dts: false,
    }),
    // 自定义插件，用于生成页面文件的path，并添加到虚拟模块
    pluginPagePathes(),
    // 自定义插件，用于生成自定义icon，并添加到虚拟模块
    pluginIcons(),
    // 移除非必要的vue-router动态路由警告: No match found for location with path
    removeNoMatch(),
  )

  return {
    base: VITE_PUBLIC_PATH || '/',
    plugins,
    resolve: {
      alias: {
        '@': path.resolve(process.cwd(), 'src'),
        '~': path.resolve(process.cwd()),
      },
    },
    server: {
      host: '0.0.0.0',
      port: 3200,
      open: false,
      proxy: {
        '/api': {
          target: VITE_PROXY_TARGET,
          changeOrigin: true,
          rewrite: path => path.replace(/^\/api/, ''),
          secure: false,
          timeout: 120000,
          proxyTimeout: 120000,
          configure: (proxy, options) => {
            // 配置此项可在响应头中看到请求的真实地址
            proxy.on('proxyRes', (proxyRes, req) => {
              proxyRes.headers['x-real-url'] = new URL(req.url || '', options.target)?.href || ''
            })
          },
        },
      },
    },
    build: {
      chunkSizeWarningLimit: 1024, // chunk 大小警告的限制（单位kb）
    },
  }
})
