/**********************************
 * @Description: 入口文件
 * @FilePath: main.js
 * @Author: Ronnie Zhang
 * @LastEditor: Ronnie Zhang
 * @LastEditTime: 2023/12/04 22:41:32
 * @Email: zclzone@outlook.com
 * Copyright © 2023 Ronnie Zhang(大脸怪) | https://isme.top
 **********************************/

import { createApp } from 'vue'
import api from './api'
import App from './App.vue'
import 'babel-polyfill'

import { setupDirectives } from './directives'
import { setupRouter } from './router'
import { setupStore } from './store'
import { getAssetUrl, setupNaiveDiscreteApi } from './utils'
import { setupLightField } from './utils/lightField'
import { applyPrimaryTheme } from './utils/theme'
import '@/styles/reset.css'
import '@/styles/global.css'
import 'uno.css'

async function bootstrap() {
  // 强制亮色：清除历史暗黑模式 class / vueuse 缓存
  document.documentElement.classList.remove('dark')
  document.documentElement.removeAttribute('data-theme')
  try {
    localStorage.removeItem('vueuse-color-scheme')
  }
  catch {}

  applyPrimaryTheme()

  const app = createApp(App)
  app.config.globalProperties.$assetUrl = getAssetUrl
  setupStore(app)
  setupDirectives(app)
  await setupRouter(app)
  app.mount('#app')
  setupNaiveDiscreteApi()
  setupLightField()
}

bootstrap()
