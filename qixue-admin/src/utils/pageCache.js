/**********************************
 * 页面 / 列表状态缓存时效，避免浏览器内存与 session 堆积
 **********************************/

/** 页面 KeepAlive 有效期（离开后超过此时长再进入会重建） */
export const PAGE_CACHE_TTL_MS = 15 * 60 * 1000

/** 列表翻页/筛选 session 缓存有效期 */
export const LIST_STATE_TTL_MS = 30 * 60 * 1000

/** KeepAlive 最多同时保留的页面实例数 */
export const PAGE_CACHE_MAX = 10

export function isCacheExpired(timestamp, ttlMs) {
  const t = Number(timestamp) || 0
  if (!t)
    return true
  return Date.now() - t > ttlMs
}

/** 清理过期的 me-crud session 缓存（进入系统时调用一次即可） */
export function pruneExpiredListStateCache(ttlMs = LIST_STATE_TTL_MS) {
  try {
    const keys = []
    for (let i = 0; i < sessionStorage.length; i++) {
      const key = sessionStorage.key(i)
      if (key && key.startsWith('me-crud:'))
        keys.push(key)
    }
    const now = Date.now()
    keys.forEach((key) => {
      try {
        const raw = sessionStorage.getItem(key)
        if (!raw) {
          sessionStorage.removeItem(key)
          return
        }
        const saved = JSON.parse(raw)
        const savedAt = Number(saved?.savedAt) || 0
        if (!savedAt || now - savedAt > ttlMs)
          sessionStorage.removeItem(key)
      }
      catch {
        sessionStorage.removeItem(key)
      }
    })
  }
  catch {}
}
