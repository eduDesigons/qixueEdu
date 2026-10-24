/** PWA：Service Worker 更新 + 安装到桌面 */

const INSTALL_DISMISS_KEY = 'qx-admin-pwa-install-dismissed'
const INSTALL_COOLDOWN_MS = 7 * 24 * 60 * 60 * 1000

const state = {
  needRefresh: false,
  offlineReady: false,
  canInstall: false,
  iosGuide: false,
  showInstall: false,
  installed: false,
  /** @type {null | ((reloadPage?: boolean) => Promise<void>)} */
  updateSW: null,
  /** @type {null | Event} */
  deferredPrompt: null,
}

const listeners = new Set()
let installReady = false
let updateDismissed = false

function emit() {
  listeners.forEach((fn) => {
    try {
      fn(getPwaState())
    }
    catch {}
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
  if (typeof window === 'undefined')
    return false
  if (window.matchMedia) {
    if (window.matchMedia('(display-mode: standalone)').matches)
      return true
    if (window.matchMedia('(display-mode: fullscreen)').matches)
      return true
    if (window.matchMedia('(display-mode: minimal-ui)').matches)
      return true
    if (window.matchMedia('(display-mode: window-controls-overlay)').matches)
      return true
  }
  if (window.navigator && window.navigator.standalone)
    return true
  return false
}

/** 给 html 打独立窗口标记，驱动 App 壳样式 */
export function applyAppShell() {
  if (typeof document === 'undefined')
    return
  const root = document.documentElement
  const on = isStandalone()
  root.classList.toggle('qx-pwa-standalone', on)
  root.classList.toggle('qx-pwa-browser', !on)
}

function isIosSafari() {
  if (typeof navigator === 'undefined')
    return false
  const ua = navigator.userAgent || ''
  const iOS = /iPad|iPhone|iPod/.test(ua) || (navigator.platform === 'MacIntel' && navigator.maxTouchPoints > 1)
  const webkit = /WebKit/.test(ua)
  const isCriOS = /CriOS|FxiOS|EdgiOS/.test(ua)
  return iOS && webkit && !isCriOS
}

function isInstallDismissed() {
  try {
    const raw = localStorage.getItem(INSTALL_DISMISS_KEY)
    if (!raw)
      return false
    const at = Number(raw)
    if (!Number.isFinite(at))
      return false
    return Date.now() - at < INSTALL_COOLDOWN_MS
  }
  catch {
    return false
  }
}

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
  catch {}
  state.showInstall = false
  emit()
}

export async function promptInstall() {
  const evt = state.deferredPrompt
  if (!evt || typeof evt.prompt !== 'function') {
    if (state.iosGuide)
      return { outcome: 'ios-guide' }
    return { outcome: 'unavailable' }
  }
  evt.prompt()
  let outcome = 'dismissed'
  try {
    const choice = await evt.userChoice
    outcome = (choice && choice.outcome) || 'dismissed'
  }
  catch {}
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
  if (typeof state.updateSW === 'function')
    await state.updateSW(true)
  else
    window.location.reload()
}

function bindInstallEvents() {
  if (typeof window === 'undefined')
    return

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
      catch {}
    })
  }
  if (!state.installed && isIosSafari())
    state.iosGuide = true
  refreshInstallVisibility()
  setTimeout(() => {
    installReady = true
    refreshInstallVisibility()
  }, 3500)

  if (!import.meta.env.PROD)
    return

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
      onRegisteredSW(_swUrl, registration) {
        if (!registration)
          return
        setInterval(() => {
          try {
            registration.update()
          }
          catch {}
        }, 60 * 60 * 1000)
        document.addEventListener('visibilitychange', () => {
          if (document.visibilityState === 'visible') {
            try {
              registration.update()
            }
            catch {}
          }
        })
      },
    })
  }
  catch (e) {
    console.warn('[pwa] register failed', e)
  }
}
