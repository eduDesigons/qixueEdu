/** PWA：Service Worker 更新 + 安装到桌面 */

const INSTALL_DISMISS_KEY = 'qx-pwa-install-dismissed'
const INSTALL_COOLDOWN_MS = 7 * 24 * 60 * 60 * 1000

const state = {
  needRefresh: false,
  offlineReady: false,
  canInstall: false,
  /** iOS Safari 无 beforeinstallprompt，改为展示引导 */
  iosGuide: false,
  showInstall: false,
  installed: false,
  /** @type {null | ((reloadPage?: boolean) => Promise<void>)} */
  updateSW: null,
  /** @type {null | Event} */
  deferredPrompt: null,
}

/** 本次会话内暂缓更新提示 */
let updateDismissed = false

const listeners = new Set()

function emit() {
  listeners.forEach((fn) => {
    try {
      fn(getPwaState())
    }
    catch (e) {}
  })
}

export function getPwaState() {
  return {
    ...state,
    needRefresh: state.needRefresh && !updateDismissed,
  }
}

export function dismissUpdatePrompt() {
  updateDismissed = true
  emit()
}

export function subscribePwa(fn) {
  listeners.add(fn)
  fn(getPwaState())
  return () => listeners.delete(fn)
}

function isStandalone() {
  if (typeof window === 'undefined') return false
  if (window.matchMedia) {
    if (window.matchMedia('(display-mode: standalone)').matches) return true
    if (window.matchMedia('(display-mode: fullscreen)').matches) return true
    if (window.matchMedia('(display-mode: minimal-ui)').matches) return true
    if (window.matchMedia('(display-mode: window-controls-overlay)').matches) return true
  }
  if (window.navigator && window.navigator.standalone) return true
  return false
}

/** 给 html 打独立窗口标记，驱动 App 壳样式（安全区 / 无浏览器栏体验） */
export function applyAppShell() {
  if (typeof document === 'undefined') return
  const root = document.documentElement
  const on = isStandalone()
  root.classList.toggle('qx-pwa-standalone', on)
  root.classList.toggle('qx-pwa-browser', !on)
  // iOS 独立窗口：同步状态栏与主题色观感
  let meta = document.querySelector('meta[name="theme-color"]')
  if (!meta) {
    meta = document.createElement('meta')
    meta.setAttribute('name', 'theme-color')
    document.head.appendChild(meta)
  }
  if (on && !meta.dataset.qxAppShell) {
    meta.dataset.qxAppShell = meta.getAttribute('content') || ''
  }
  // 独立窗口：状态栏透明（与 apple black-translucent + viewport-fit=cover 配合）
  if (on) {
    document.querySelectorAll('meta[name="theme-color"]').forEach((el) => {
      el.setAttribute('content', 'transparent')
    })
  } else if (meta.dataset.qxAppShell) {
    meta.setAttribute('content', meta.dataset.qxAppShell || '#4286F5')
  }
}

function isIosSafari() {
  if (typeof navigator === 'undefined') return false
  const ua = navigator.userAgent || ''
  const iOS = /iPad|iPhone|iPod/.test(ua) || (navigator.platform === 'MacIntel' && navigator.maxTouchPoints > 1)
  const webkit = /WebKit/.test(ua)
  const isCriOS = /CriOS|FxiOS|EdgiOS/.test(ua)
  return iOS && webkit && !isCriOS
}

function isInstallDismissed() {
  try {
    const raw = localStorage.getItem(INSTALL_DISMISS_KEY)
    if (!raw) return false
    const at = Number(raw)
    if (!Number.isFinite(at)) return false
    return Date.now() - at < INSTALL_COOLDOWN_MS
  }
  catch (e) {
    return false
  }
}

let installReady = false

function refreshInstallVisibility() {
  state.installed = isStandalone()
  const eligible = state.canInstall || state.iosGuide
  state.showInstall = installReady && !state.installed && eligible && !isInstallDismissed()
  emit()
}

export function dismissInstallPrompt() {
  try {
    localStorage.setItem(INSTALL_DISMISS_KEY, String(Date.now()))
  }
  catch (e) {}
  state.showInstall = false
  emit()
}

export async function promptInstall() {
  const evt = state.deferredPrompt
  if (!evt || typeof evt.prompt !== 'function') {
    if (state.iosGuide) return { outcome: 'ios-guide' }
    return { outcome: 'unavailable' }
  }
  evt.prompt()
  let outcome = 'dismissed'
  try {
    const choice = await evt.userChoice
    outcome = (choice && choice.outcome) || 'dismissed'
  }
  catch (e) {}
  state.deferredPrompt = null
  state.canInstall = false
  if (outcome === 'accepted') {
    state.installed = true
    state.showInstall = false
  }
  else {
    dismissInstallPrompt()
  }
  emit()
  return { outcome }
}

export async function applyUpdate() {
  if (typeof state.updateSW === 'function') {
    await state.updateSW(true)
  }
  else {
    window.location.reload()
  }
}

function bindInstallEvents() {
  if (typeof window === 'undefined') return

  window.addEventListener('beforeinstallprompt', (e) => {
    e.preventDefault()
    state.deferredPrompt = e
    state.canInstall = true
    refreshInstallVisibility()
  })

  window.addEventListener('appinstalled', () => {
    state.deferredPrompt = null
    state.canInstall = false
    state.installed = true
    state.showInstall = false
    emit()
  })

  if (window.matchMedia) {
    window.matchMedia('(display-mode: standalone)').addEventListener('change', refreshInstallVisibility)
  }
}

/** 注册 SW；生产构建由 vite-plugin-pwa 注入 virtual:pwa-register */
export async function setupPwa() {
  bindInstallEvents()
  state.installed = isStandalone()
  applyAppShell()
  if (window.matchMedia) {
    ;['standalone', 'fullscreen', 'minimal-ui', 'window-controls-overlay'].forEach((mode) => {
      try {
        window.matchMedia(`(display-mode: ${mode})`).addEventListener('change', () => {
          state.installed = isStandalone()
          applyAppShell()
          refreshInstallVisibility()
        })
      }
      catch (e) {}
    })
  }
  if (!state.installed && isIosSafari()) {
    state.iosGuide = true
  }
  refreshInstallVisibility()
  // 稍晚再弹出安装提醒，避免盖住首屏登录/导航
  setTimeout(() => {
    installReady = true
    refreshInstallVisibility()
  }, 3500)

  if (!import.meta.env.PROD) return

  try {
    const { registerSW } = await import('virtual:pwa-register')
    state.updateSW = registerSW({
      immediate: true,
      onNeedRefresh() {
        state.needRefresh = true
        emit()
      },
      onOfflineReady() {
        state.offlineReady = true
        emit()
      },
      onRegisteredSW(swUrl, registration) {
        if (!registration) return
        // 定期检查更新（打开页后每小时）
        setInterval(() => {
          try {
            registration.update()
          }
          catch (e) {}
        }, 60 * 60 * 1000)
        // 从后台切回前台时也检查
        document.addEventListener('visibilitychange', () => {
          if (document.visibilityState === 'visible') {
            try {
              registration.update()
            }
            catch (e) {}
          }
        })
      },
    })
  }
  catch (e) {
    console.warn('[pwa] register failed', e)
  }
}
