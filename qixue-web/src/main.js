import { getStorage } from '@/common/storage'
import { createApp } from 'vue'
import { bindRouter } from '@/common/nav'
import { applyAppShell, setupPwa } from '@/common/pwa'
import { installApp } from './install'
import { applyFontScale } from './common/fontScale'
import { http } from './common/request'
import mixin from './common/mixin'
import App from './App.vue'
import CoinBadge from './components/CoinBadge.vue'
import router from './router'
import store from './store'
import config from './config'
import './styles/base.css'
import './styles/global.scss'
import './styles/desktop.css'
import './styles/app-shell.css'

applyFontScale()
bindRouter(router)
applyAppShell()
setupPwa()

const savedToken = getStorage('token')
if (savedToken) {
  store.commit('setToken', savedToken)
  const savedUser = getStorage('userInfo')
  if (savedUser && typeof savedUser === 'object') store.commit('setUserInfo', savedUser)
}
try {
  const savedSys = getStorage('sysInfo')
  if (savedSys && savedSys.sysConfig) store.commit('setSysInfo', savedSys)
} catch (e) {}

const app = createApp(App)
installApp(app)
app.component('CoinBadge', CoinBadge)
app.mixin(mixin)
app.use(store)
app.use(router)
app.config.globalProperties.$config = config
app.config.globalProperties.$http = http
router.isReady().then(() => app.mount('#app'))
