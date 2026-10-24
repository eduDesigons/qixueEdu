import {
  assertDownloadComplete,
  isCompleteVideoMark,
  isMp4Url,
  mediaCacheNames,
  offlinePlaybackUrl,
  parseContentRangeTotal,
  stableMediaKey,
} from './mediaCacheCore.js'

export const NAMES = mediaCacheNames('qx-web')
const MAX_VIDEOS = 8
const MAX_VIDEO_BYTES = 1200 * 1024 * 1024

function readIndex() {
  try {
    const raw = localStorage.getItem(NAMES.indexKey)
    const data = raw ? JSON.parse(raw) : []
    return Array.isArray(data) ? data : []
  }
  catch {
    return []
  }
}

function writeIndex(items) {
  localStorage.setItem(NAMES.indexKey, JSON.stringify(items))
}

export function formatBytes(n) {
  const v = Number(n) || 0
  if (v < 1024)
    return `${v} B`
  if (v < 1024 * 1024)
    return `${(v / 1024).toFixed(1)} KB`
  if (v < 1024 * 1024 * 1024)
    return `${(v / 1024 / 1024).toFixed(1)} MB`
  return `${(v / 1024 / 1024 / 1024).toFixed(2)} GB`
}

export async function getVideoCacheStatus(url) {
  if (!url || typeof caches === 'undefined')
    return { cached: false, bytes: 0 }
  const key = stableMediaKey(url)
  const cache = await caches.open(NAMES.videos)
  const hit = await cache.match(key)
  if (!hit || hit.status !== 200 || !isCompleteVideoMark(hit.headers)) {
    writeIndex(readIndex().filter((item) => item.key !== key))
    return { cached: false, bytes: 0, key }
  }
  const bytes = Number(hit.headers.get('content-length') || 0)
  return {
    cached: true,
    bytes,
    key,
    cachedAt: hit.headers.get('x-qx-cached-at') || '',
  }
}

/** 检查 CDN 是否真的允许跨域读取，以及 Range 是否返回 206。不写入缓存。 */
export async function probeVideo(url) {
  let res
  try {
    res = await fetch(url, {
      method: 'GET',
      mode: 'cors',
      credentials: 'omit',
      cache: 'no-store',
      headers: { Range: 'bytes=0-1' },
    })
  }
  catch {
    return {
      ok: false,
      cors: false,
      range: false,
      status: 0,
      total: null,
      error: '跨域读取失败。CDN 需对当前站点返回 Access-Control-Allow-Origin，并允许 GET 与 Range。签名参数已保留在原始地址中。',
    }
  }
  const contentRange = res.headers.get('content-range')
  const total = parseContentRangeTotal(contentRange)
  const range = res.status === 206 && total != null
  try { await res.body?.cancel() }
  catch { /* ignore */ }
  return {
    ok: res.status === 200 || res.status === 206,
    cors: true,
    range,
    status: res.status,
    total,
    contentType: res.headers.get('content-type') || '',
    error: (res.status === 200 || res.status === 206)
      ? ''
      : `探测失败（HTTP ${res.status}）。请确认签名 URL 仍有效。`,
  }
}

async function storageRoom() {
  if (!navigator.storage?.estimate)
    return { usage: 0, quota: 0 }
  const { usage = 0, quota = 0 } = await navigator.storage.estimate()
  return { usage, quota }
}

async function evictOldest(exceptKey) {
  const cache = await caches.open(NAMES.videos)
  const items = readIndex()
    .filter((item) => item.key !== exceptKey)
    .sort((a, b) => String(a.cachedAt).localeCompare(String(b.cachedAt)))
  const victim = items[0]
  if (!victim)
    return false
  await cache.delete(victim.key)
  writeIndex(readIndex().filter((item) => item.key !== victim.key))
  return true
}

async function ensureSpace(bytes, exceptKey) {
  const index = readIndex().filter((item) => item.key !== exceptKey)
  const used = index.reduce((sum, item) => sum + (Number(item.bytes) || 0), 0)
  if (index.length >= MAX_VIDEOS || used + bytes > MAX_VIDEO_BYTES) {
    const removed = await evictOldest(exceptKey)
    if (!removed) {
      throw new Error(`视频缓存已达上限（${MAX_VIDEOS} 个或 ${formatBytes(MAX_VIDEO_BYTES)}），请先清除旧视频`)
    }
    return ensureSpace(bytes, exceptKey)
  }
  const { usage, quota } = await storageRoom()
  if (quota && usage + bytes > quota * 0.9) {
    const removed = await evictOldest(exceptKey)
    if (!removed)
      throw new Error('浏览器存储空间不足，无法缓存该视频')
    return ensureSpace(bytes, exceptKey)
  }
}

export async function cacheFullVideo(url, { title = '', signal, onProgress } = {}) {
  if (!isMp4Url(url))
    throw new Error('仅支持缓存 MP4 地址')
  const probe = await probeVideo(url)
  if (!probe.cors || !probe.ok)
    throw new Error(probe.error || '无法确认 CDN 跨域与 Range 响应')
  const key = stableMediaKey(url)
  if (probe.total)
    await ensureSpace(probe.total, key)

  const res = await fetch(url, {
    method: 'GET',
    mode: 'cors',
    credentials: 'omit',
    cache: 'no-store',
    signal,
  })
  if (res.status !== 200) {
    throw new Error(res.status === 206
      ? '完整下载收到 206 片段，已中止，避免把片段当成整段视频'
      : `完整下载需要 HTTP 200，实际 ${res.status}`)
  }
  const contentLengthHeader = res.headers.get('content-length')
  const contentLength = contentLengthHeader ? Number(contentLengthHeader) : null
  const reader = res.body?.getReader()
  if (!reader)
    throw new Error('浏览器无法读取视频响应体')
  const chunks = []
  let received = 0
  while (true) {
    const { done, value } = await reader.read()
    if (done)
      break
    chunks.push(value)
    received += value.byteLength
    onProgress?.({
      received,
      total: contentLength || probe.total || 0,
    })
  }
  const check = assertDownloadComplete({
    status: 200,
    received,
    contentLength: Number.isFinite(contentLength) ? contentLength : null,
    contentRangeTotal: probe.total,
  })
  if (!check.ok)
    throw new Error(check.reason)

  await ensureSpace(received, key)
  const type = (res.headers.get('content-type') || 'video/mp4').split(';')[0]
  const blob = new Blob(chunks, { type })
  const headers = new Headers({
    'Content-Type': type,
    'Content-Length': String(blob.size),
    'Accept-Ranges': 'bytes',
    'X-Qx-Video-Complete': '1',
    'X-Qx-Cached-At': new Date().toISOString(),
  })
  const cache = await caches.open(NAMES.videos)
  await cache.put(key, new Response(blob, { status: 200, statusText: 'OK', headers }))
  const items = readIndex().filter((item) => item.key !== key)
  items.push({
    key,
    title: title || '',
    bytes: blob.size,
    cachedAt: headers.get('x-qx-cached-at'),
    range: probe.range,
  })
  writeIndex(items)
  return {
    bytes: blob.size,
    range: probe.range,
    warning: probe.range ? '' : 'CDN 未对 Range 返回 206。联网播放仍走原始请求；离线拖动由本地完整文件切段提供。',
  }
}

/** 仅在离线时使用。联网播放继续走原始地址，不经过这段逻辑。 */
export async function offlinePlaybackSrc(url) {
  const status = await getVideoCacheStatus(url)
  if (!status.cached)
    return ''
  if (navigator.serviceWorker?.controller)
    return offlinePlaybackUrl(status.key)
  const cache = await caches.open(NAMES.videos)
  const hit = await cache.match(status.key)
  if (!hit)
    return ''
  const blob = await hit.blob()
  return URL.createObjectURL(blob)
}

export async function removeCachedVideo(url) {
  const key = stableMediaKey(url)
  if (typeof caches !== 'undefined') {
    const cache = await caches.open(NAMES.videos)
    await cache.delete(key)
  }
  writeIndex(readIndex().filter((item) => item.key !== key))
}

export async function clearCachedVideos() {
  if (typeof caches !== 'undefined')
    await caches.delete(NAMES.videos)
  writeIndex([])
}
