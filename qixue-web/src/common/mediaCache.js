import {
  assertDownloadComplete,
  isAudioUrl,
  isCompleteMediaResponse,
  isMp4Url,
  mediaCacheNames,
  offlinePlaybackUrl,
  parseContentRangeTotal,
  stableMediaKey,
} from './mediaCacheCore.js'

export const NAMES = mediaCacheNames('qx-web')
const CONTENT_STORAGE_PREFIXES = [
  'qixue_home',
  'qixue_course',
  'qixue_learn',
  'qixue_rank',
  'qixue_me_',
  'qixue_path',
  'qixue_my_',
  'chapter_progress_',
]
const autoJobs = new Map()

function cacheNameOf(kind) {
  return kind === 'audio' ? NAMES.audio : NAMES.videos
}

function indexStorageKey(kind) {
  return kind === 'audio' ? NAMES.audioIndexKey : NAMES.indexKey
}

function readIndex(kind = 'video') {
  try {
    const raw = localStorage.getItem(indexStorageKey(kind))
    const data = raw ? JSON.parse(raw) : []
    return Array.isArray(data) ? data : []
  }
  catch {
    return []
  }
}

function writeIndex(items, kind = 'video') {
  localStorage.setItem(indexStorageKey(kind), JSON.stringify(items))
}

function dropIndexKey(kind, key) {
  const items = readIndex(kind)
  if (!items.some((item) => item.key === key))
    return
  writeIndex(items.filter((item) => item.key !== key), kind)
}

async function storedMedia(kind, key) {
  if (!key || typeof caches === 'undefined') {
    if (key)
      dropIndexKey(kind, key)
    return null
  }
  const cache = await caches.open(cacheNameOf(kind))
  const hit = await cache.match(key)
  if (isCompleteMediaResponse(hit))
    return hit
  if (hit)
    await cache.delete(key)
  dropIndexKey(kind, key)
  return null
}

export async function purgeIncompleteMedia() {
  if (typeof caches === 'undefined')
    return
  let names = []
  try {
    names = await caches.keys()
  }
  catch {
    return
  }
  for (const kind of ['video', 'audio']) {
    const name = cacheNameOf(kind)
    if (!names.includes(name)) {
      if (readIndex(kind).length)
        writeIndex([], kind)
      continue
    }
    try {
      const cache = await caches.open(name)
      const reqs = await cache.keys()
      const keep = new Set()
      for (const req of reqs) {
        const hit = await cache.match(req)
        if (isCompleteMediaResponse(hit)) {
          keep.add(req.url)
          continue
        }
        await cache.delete(req)
      }
      const next = readIndex(kind).filter((item) => keep.has(item.key))
      if (next.length !== readIndex(kind).length)
        writeIndex(next, kind)
    }
    catch { /* ignore */ }
  }
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
  if (!url)
    return { cached: false, bytes: 0 }
  const key = stableMediaKey(url)
  const hit = await storedMedia('video', key)
  if (!hit)
    return { cached: false, bytes: 0, key }
  return {
    cached: true,
    bytes: Number(hit.headers.get('content-length') || 0),
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

export async function cacheFullMedia(url, { title = '', signal, onProgress, kind = 'video' } = {}) {
  const accept = kind === 'audio' ? isAudioUrl : isMp4Url
  if (!accept(url))
    throw new Error(kind === 'audio' ? '仅支持缓存常见音频地址' : '仅支持缓存 MP4 地址')
  const probe = await probeVideo(url)
  if (!probe.cors || !probe.ok)
    throw new Error(probe.error || '无法确认 CDN 跨域与 Range 响应')
  const key = stableMediaKey(url)
  await storedMedia(kind, key)
  try {
    const res = await fetch(url, {
      method: 'GET',
      mode: 'cors',
      credentials: 'omit',
      cache: 'no-store',
      signal,
    })
    if (res.status !== 200) {
      throw new Error(res.status === 206
        ? '完整下载收到 206 片段，已中止，避免把片段当成整段文件'
        : `完整下载需要 HTTP 200，实际 ${res.status}`)
    }
    const contentLengthHeader = res.headers.get('content-length')
    const contentLength = contentLengthHeader ? Number(contentLengthHeader) : null
    const reader = res.body?.getReader()
    if (!reader)
      throw new Error('浏览器无法读取媒体响应体')
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

    const fallbackType = kind === 'audio' ? 'audio/mpeg' : 'video/mp4'
    const type = (res.headers.get('content-type') || fallbackType).split(';')[0]
    const blob = new Blob(chunks, { type })
    const headers = new Headers({
      'Content-Type': type,
      'Content-Length': String(blob.size),
      'Accept-Ranges': 'bytes',
      'X-Qx-Video-Complete': '1',
      'X-Qx-Cached-At': new Date().toISOString(),
    })
    const cache = await caches.open(cacheNameOf(kind))
    await cache.put(key, new Response(blob, { status: 200, statusText: 'OK', headers }))
    const items = readIndex(kind).filter((item) => item.key !== key)
    items.push({
      key,
      title: title || '',
      bytes: blob.size,
      cachedAt: headers.get('x-qx-cached-at'),
      range: probe.range,
    })
    writeIndex(items, kind)
    return {
      bytes: blob.size,
      range: probe.range,
      warning: probe.range ? '' : 'CDN 未对 Range 返回 206。完整文件已在本机，之后播放和拖动都用这份。',
    }
  }
  catch (err) {
    try { await storedMedia(kind, key) }
    catch { /* ignore */ }
    throw err
  }
}

export function cacheFullVideo(url, opts = {}) {
  return cacheFullMedia(url, { ...opts, kind: 'video' })
}

export function autoCacheMedia(url, options = {}) {
  if (!url)
    return Promise.resolve()
  const kind = options.kind || (isMp4Url(url) ? 'video' : 'audio')
  const key = `${kind}:${stableMediaKey(url)}`
  if (autoJobs.has(key))
    return autoJobs.get(key)
  const job = (async () => {
    const status = kind === 'audio'
      ? await getAudioCacheStatus(url)
      : await getVideoCacheStatus(url)
    if (status.cached)
      return
    await cacheFullMedia(url, { ...options, kind })
  })().catch(() => {}).finally(() => {
    autoJobs.delete(key)
  })
  autoJobs.set(key, job)
  return job
}

/** 本地已有完整文件时返回可播放地址（有网也用这份）；没有则返回空字符串，调用方继续用原地址。 */
export async function localPlaybackSrc(url, kind = 'video') {
  const status = kind === 'audio' ? await getAudioCacheStatus(url) : await getVideoCacheStatus(url)
  if (!status.cached)
    return ''
  if (kind === 'video' && navigator.serviceWorker?.controller)
    return offlinePlaybackUrl(status.key, 'video')
  const cache = await caches.open(kind === 'audio' ? NAMES.audio : NAMES.videos)
  const hit = await cache.match(status.key)
  if (!hit)
    return ''
  const blob = await hit.blob()
  return URL.createObjectURL(blob)
}

/** 视频入口：完整本地文件优先，不区分是否联网。 */
export async function offlinePlaybackSrc(url) {
  return localPlaybackSrc(url, 'video')
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

export async function getAudioCacheStatus(url) {
  if (!url)
    return { cached: false, bytes: 0 }
  const key = stableMediaKey(url)
  const hit = await storedMedia('audio', key)
  if (!hit)
    return { cached: false, bytes: 0, key }
  return { cached: true, bytes: Number(hit.headers.get('content-length') || 0), key }
}

function isContentStorageKey(fullKey) {
  if (!fullKey.startsWith('qx:'))
    return false
  const key = fullKey.slice(3)
  return CONTENT_STORAGE_PREFIXES.some((prefix) => key.startsWith(prefix))
}

function contentStorageBytes() {
  let bytes = 0
  try {
    for (let i = 0; i < localStorage.length; i++) {
      const key = localStorage.key(i)
      if (!isContentStorageKey(key))
        continue
      bytes += (key.length + (localStorage.getItem(key) || '').length) * 2
    }
  }
  catch { /* ignore */ }
  return bytes
}

function clearContentStorage() {
  const keys = []
  try {
    for (let i = 0; i < localStorage.length; i++) {
      const key = localStorage.key(i)
      if (isContentStorageKey(key))
        keys.push(key)
    }
    keys.forEach((key) => localStorage.removeItem(key))
  }
  catch { /* ignore */ }
}

async function measureCache(name) {
  if (typeof caches === 'undefined')
    return 0
  const names = await caches.keys()
  if (!names.includes(name))
    return 0
  const cache = await caches.open(name)
  const reqs = await cache.keys()
  let bytes = 0
  for (const req of reqs) {
    const res = await cache.match(req)
    if (!res)
      continue
    const len = Number(res.headers.get('content-length') || 0)
    if (len > 0) {
      bytes += len
      continue
    }
    try {
      const blob = await res.clone().blob()
      bytes += blob.size
    }
    catch { /* ignore */ }
  }
  return bytes
}

function isAppCache(name) {
  return name.includes('precache') || name.startsWith('workbox-')
}

export async function listCacheGroups() {
  await purgeIncompleteMedia()
  const names = typeof caches === 'undefined' ? [] : await caches.keys()
  const known = new Set([NAMES.videos, NAMES.audio, NAMES.images, NAMES.content])
  const videos = readIndex('video')
  const audios = readIndex('audio')
  const videoBytes = videos.reduce((sum, item) => sum + (Number(item.bytes) || 0), 0)
  const audioBytes = audios.reduce((sum, item) => sum + (Number(item.bytes) || 0), 0)
  const imageBytes = await measureCache(NAMES.images)
  const contentBytes = (await measureCache(NAMES.content)) + contentStorageBytes()
  const appNames = names.filter(isAppCache)
  let appBytes = 0
  for (const name of appNames)
    appBytes += await measureCache(name)
  const otherNames = names.filter((name) => !known.has(name) && !isAppCache(name))
  let otherBytes = 0
  for (const name of otherNames)
    otherBytes += await measureCache(name)
  const checked = (bytes) => bytes > 0
  return [
    {
      id: 'videos',
      name: '视频',
      bytes: videoBytes,
      checked: checked(videoBytes),
      expandable: true,
      items: videos.map((item) => ({
        key: item.key,
        title: item.title || '未命名视频',
        bytes: Number(item.bytes) || 0,
        checked: true,
      })),
    },
    {
      id: 'audio',
      name: '音频',
      bytes: audioBytes,
      checked: checked(audioBytes),
      expandable: true,
      items: audios.map((item) => ({
        key: item.key,
        title: item.title || '未命名音频',
        bytes: Number(item.bytes) || 0,
        checked: true,
      })),
    },
    { id: 'images', name: '图片', bytes: imageBytes, checked: checked(imageBytes), items: [], cacheNames: [NAMES.images] },
    { id: 'content', name: '浏览内容', bytes: contentBytes, checked: checked(contentBytes), items: [], cacheNames: [NAMES.content] },
    { id: 'other', name: '其他', bytes: otherBytes, checked: checked(otherBytes), items: [], cacheNames: otherNames },
    { id: 'app', name: '应用文件', bytes: appBytes, checked: false, items: [], cacheNames: appNames },
  ]
}

export async function cacheUsageBytes() {
  const groups = await listCacheGroups()
  return groups.filter((group) => group.id !== 'app').reduce((sum, group) => sum + group.bytes, 0)
}

async function deleteCaches(names) {
  if (typeof caches === 'undefined')
    return
  for (const name of names || []) {
    if (name)
      await caches.delete(name)
  }
}

async function removeIndexed(kind, key) {
  if (typeof caches !== 'undefined') {
    const cache = await caches.open(cacheNameOf(kind))
    await cache.delete(key)
  }
  writeIndex(readIndex(kind).filter((item) => item.key !== key), kind)
}

export async function clearCacheGroups(groups) {
  for (const group of groups) {
    if (group.id === 'videos' || group.id === 'audio') {
      const kind = group.id === 'audio' ? 'audio' : 'video'
      const items = group.items || []
      const picked = items.filter((item) => item.checked)
      if (!picked.length)
        continue
      if (picked.length === items.length) {
        await deleteCaches([cacheNameOf(kind)])
        writeIndex([], kind)
        continue
      }
      for (const item of picked)
        await removeIndexed(kind, item.key)
      continue
    }
    if (!group.checked)
      continue
    if (group.id === 'content')
      clearContentStorage()
    await deleteCaches(group.cacheNames)
  }
}
