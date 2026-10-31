import { previewImages } from '@/common/media'

export function upx2px(value) {
  return Number(value) * (Math.min(window.innerWidth || 375, 480) / 750)
}

export function getSystemInfoSync() {
  const width = window.innerWidth || 375
  const height = window.innerHeight || 667
  return {
    windowWidth: width,
    windowHeight: height,
    windowTop: 0,
    screenWidth: width,
    screenHeight: height,
    pixelRatio: window.devicePixelRatio || 1,
    statusBarHeight: 0,
    platform: 'h5',
    safeAreaInsets: { top: 0, right: 0, bottom: 0, left: 0 },
    safeArea: { left: 0, top: 0, right: width, bottom: height, width, height },
  }
}

export function getSystemInfo(options = {}) {
  const info = getSystemInfoSync()
  options.success && options.success(info)
  options.complete && options.complete(info)
  return info
}

export function pageScrollTo({ scrollTop = 0, duration = 0 } = {}) {
  window.scrollTo({ top: scrollTop, behavior: duration ? 'smooth' : 'auto' })
}

export function createSelectorQuery() {
  let root = document
  const query = {
    in(vm) {
      root = vm && vm.$el ? vm.$el : document
      return query
    },
    select(selector) {
      let callback = null
      const chain = {
        boundingClientRect(fn) {
          callback = fn
          return chain
        },
        fields(_opt, fn) {
          callback = fn
          return chain
        },
        exec(done) {
          const el = root.querySelector ? root.querySelector(selector) : null
          const rect = el ? el.getBoundingClientRect() : null
          if (callback) callback(rect)
          if (done) done([rect])
        },
      }
      return chain
    },
  }
  return query
}

function mediaContext(id, tag) {
  const root = typeof id === 'string' ? document.getElementById(id) : null
  const media = root && root.tagName === tag.toUpperCase()
    ? root
    : (root && root.querySelector && root.querySelector(tag)) || document.querySelector(tag)
  return {
    play: () => media && media.play && media.play(),
    pause: () => media && media.pause && media.pause(),
    seek: (time) => {
      if (media) media.currentTime = Number(time) || 0
    },
    playbackRate: (rate) => {
      if (media) media.playbackRate = Number(rate) || 1
    },
    requestFullScreen: () => media && media.requestFullscreen && media.requestFullscreen(),
    exitFullScreen: () => document.exitFullscreen && document.exitFullscreen(),
  }
}

export function createAudioContext(id) {
  return mediaContext(id, 'audio')
}

export function createVideoContext(id) {
  return mediaContext(id, 'video')
}

export function showModal(options = {}) {
  const opt = options || {}
  return new Promise((resolve) => {
    const mask = document.createElement('div')
    mask.style.cssText = 'position:fixed;inset:0;z-index:100002;background:rgba(0,0,0,.45);display:flex;align-items:center;justify-content:center;padding:24px;'
    const box = document.createElement('div')
    box.style.cssText = 'width:min(300px,100%);background:#fff;border-radius:12px;overflow:hidden;text-align:center;'
    const title = document.createElement('div')
    title.style.cssText = 'padding:22px 16px 8px;font-size:16px;font-weight:600;color:#1a1a1a;'
    title.textContent = opt.title || '提示'
    const content = document.createElement('div')
    content.style.cssText = 'padding:0 16px 20px;font-size:14px;color:#666;line-height:1.5;'
    content.textContent = opt.content || ''
    const actions = document.createElement('div')
    actions.style.cssText = 'display:flex;border-top:1px solid #eee;'
    const finish = (confirm) => {
      mask.remove()
      const result = { confirm, cancel: !confirm }
      opt.success && opt.success(result)
      opt.complete && opt.complete(result)
      resolve(result)
    }
    if (opt.showCancel !== false) {
      const cancel = document.createElement('button')
      cancel.textContent = opt.cancelText || '取消'
      cancel.style.cssText = 'flex:1;height:48px;border:0;background:#fff;font-size:16px;color:#666;'
      cancel.onclick = () => finish(false)
      actions.appendChild(cancel)
    }
    const ok = document.createElement('button')
    ok.textContent = opt.confirmText || '确定'
    ok.style.cssText = 'flex:1;height:48px;border:0;border-left:1px solid #eee;background:#fff;font-size:16px;color:#4286F5;'
    ok.onclick = () => finish(true)
    actions.appendChild(ok)
    box.appendChild(title)
    if (opt.content) box.appendChild(content)
    box.appendChild(actions)
    mask.appendChild(box)
    document.body.appendChild(mask)
  })
}

export function setClipboard(data) {
  const text = String(typeof data === 'object' && data && 'data' in data ? data.data : data ?? '')
  const done = () => {
    if (typeof data === 'object' && data && data.success) data.success()
  }
  if (navigator.clipboard && navigator.clipboard.writeText) {
    return navigator.clipboard.writeText(text).then(done).catch(done)
  }
  done()
  return Promise.resolve()
}

export function openLocation({ latitude, longitude, name, address } = {}) {
  const query = [name, address].filter(Boolean).join(' ') || `${latitude},${longitude}`
  window.open(`https://uri.amap.com/marker?position=${longitude},${latitude}&name=${encodeURIComponent(query)}`, '_blank')
}

export function previewImage({ urls = [], current } = {}) {
  previewImages(urls, current)
}

export function timeFrom(timestamp, format) {
  let ts = Number(timestamp)
  if (!ts) return ''
  if (String(Math.floor(ts)).length === 10) ts *= 1000
  const date = new Date(ts)
  const pad = (n) => String(n).padStart(2, '0')
  if (format) {
    return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())}`
  }
  const diff = (Date.now() - ts) / 1000
  if (diff < 60) return '刚刚'
  if (diff < 3600) return `${Math.floor(diff / 60)}分钟前`
  if (diff < 86400) return `${Math.floor(diff / 3600)}小时前`
  if (diff < 2592000) return `${Math.floor(diff / 86400)}天前`
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())}`
}
