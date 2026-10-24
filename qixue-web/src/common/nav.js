import { reactive } from 'vue'

const pageStack = []
let router = null
let navAction = ''

export const navState = reactive({
  title: '',
  override: '',
})

export function bindRouter(instance) {
  router = instance
  instance.afterEach((to) => {
    const entry = {
      route: to.path.replace(/^\//, ''),
      options: { ...to.query },
    }
    if (!pageStack.length || navAction === 'init') {
      pageStack.splice(0, pageStack.length, entry)
    } else if (navAction === 'push') {
      pageStack.push(entry)
    } else if (navAction === 'replace') {
      pageStack.pop()
      pageStack.push(entry)
    } else if (navAction === 'switch' || navAction === 'relaunch') {
      pageStack.splice(0, pageStack.length, entry)
    } else if (navAction === 'back') {
      pageStack.pop()
      if (!pageStack.length) pageStack.push(entry)
    } else if (pageStack[pageStack.length - 1]?.route !== entry.route) {
      pageStack.splice(0, pageStack.length, entry)
    }
    navAction = ''
    navState.override = ''
    navState.title = to.meta.title || ''
  })
}

function parseUrl(url) {
  const raw = String(url || '')
  const [pathPart, qs] = raw.split('?')
  const path = pathPart.startsWith('/') ? pathPart : `/${pathPart}`
  const query = {}
  if (qs) {
    new URLSearchParams(qs).forEach((value, key) => {
      query[key] = value
    })
  }
  return { path, query }
}

function go(url, action) {
  if (!router) return
  const { path, query } = parseUrl(url)
  navAction = action
  const method = action === 'push' ? 'push' : 'replace'
  return router[method]({ path, query }).catch(() => {})
}

export function navigateTo(url) {
  return go(typeof url === 'object' ? url.url : url, 'push')
}

export function redirectTo(url) {
  return go(typeof url === 'object' ? url.url : url, 'replace')
}

export function reLaunch(url) {
  return go(typeof url === 'object' ? url.url : url, 'relaunch')
}

export function switchTab(url) {
  return go(typeof url === 'object' ? url.url : url, 'switch')
}

export function navigateBack(delta = 1) {
  if (!router) return
  const n = typeof delta === 'object' ? (delta.delta || 1) : delta
  navAction = 'back'
  router.go(-Math.abs(n) || -1)
}

export function setPageTitle(title) {
  const t = typeof title === 'object' ? title.title : title
  navState.override = t || ''
  if (t) document.title = t
}

export function getCurrentPages() {
  return pageStack.map((item) => ({ ...item }))
}
