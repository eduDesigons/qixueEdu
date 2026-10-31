/** 媒体缓存纯逻辑：页面与 Service Worker 共用，不依赖 DOM。 */

export const AUTH_QUERY_KEYS = new Set([
  'token',
  'expires',
  'token_path',
  'signature',
  'sig',
  'auth',
  'auth_key',
  'policy',
])

export const OFFLINE_VIDEO_PATH = '/qx-offline-video'
export const OFFLINE_AUDIO_PATH = '/qx-offline-audio'

export function offlinePlaybackUrl(stableKey, kind = 'video') {
  const path = kind === 'audio' ? OFFLINE_AUDIO_PATH : OFFLINE_VIDEO_PATH
  return `${path}?key=${encodeURIComponent(stableKey)}`
}

export function mediaCacheNames(prefix) {
  return {
    images: `${prefix}-images`,
    videos: `${prefix}-videos`,
    audio: `${prefix}-audio`,
    content: `${prefix}-content`,
    indexKey: `${prefix}-video-index`,
    audioIndexKey: `${prefix}-audio-index`,
  }
}

export function isMp4Url(input) {
  try {
    const url = new URL(input, 'https://local.invalid')
    return /\.mp4$/i.test(url.pathname)
  }
  catch {
    return /\.mp4(\?|#|$)/i.test(String(input || ''))
  }
}

export function isAudioUrl(input) {
  try {
    const url = new URL(input, 'https://local.invalid')
    return /\.(mp3|m4a|aac|wav|ogg)$/i.test(url.pathname)
  }
  catch {
    return /\.(mp3|m4a|aac|wav|ogg)(\?|#|$)/i.test(String(input || ''))
  }
}

const CONTENT_ROOTS = ['/home/', '/course/', '/path/', '/sys/', '/rank/']

/** 课程、首页、系列、系统说明、榜单的 GET。学习进度和账号接口不放进这份缓存。 */
export function isContentApiUrl(input) {
  try {
    const url = new URL(input, 'https://local.invalid')
    const at = url.pathname.indexOf('/webapp/')
    if (at < 0)
      return false
    const rest = url.pathname.slice(at + '/webapp'.length)
    return CONTENT_ROOTS.some((prefix) => rest.startsWith(prefix))
  }
  catch {
    return false
  }
}
export function stableMediaKey(input) {
  const url = new URL(input, 'https://local.invalid')
  const drop = []
  url.searchParams.forEach((_, key) => {
    const lower = key.toLowerCase()
    if (AUTH_QUERY_KEYS.has(lower) || lower.startsWith('x-amz-'))
      drop.push(key)
  })
  drop.forEach((key) => url.searchParams.delete(key))
  url.hash = ''
  return url.href
}

export function parseContentRangeTotal(header) {
  if (!header)
    return null
  const match = /\/(\d+)\s*$/.exec(String(header))
  if (!match)
    return null
  const total = Number(match[1])
  return Number.isFinite(total) ? total : null
}

export function isCompleteVideoMark(headers) {
  if (!headers || typeof headers.get !== 'function')
    return false
  return headers.get('x-qx-video-complete') === '1'
}

/** 只有带完整标记、HTTP 200、且长度大于 0 的响应才算缓存完成。 */
export function isCompleteMediaResponse(response) {
  if (!response || response.status !== 200 || !isCompleteVideoMark(response.headers))
    return false
  const len = Number(response.headers.get('content-length') || 0)
  return Number.isFinite(len) && len > 0
}

/**
 * 只有完整的 HTTP 200 才能写入视频缓存。
 * 206 片段、空包、长度不一致一律拒绝。
 */
export function assertDownloadComplete({ status, received, contentLength, contentRangeTotal }) {
  if (status !== 200) {
    return {
      ok: false,
      reason: status === 206
        ? '响应是 206 片段，不能当作完整视频缓存'
        : `完整视频需要 HTTP 200，实际 ${status}`,
    }
  }
  if (!Number.isFinite(received) || received < 1) {
    return { ok: false, reason: '响应体为空，未写入缓存' }
  }
  if (contentLength != null && contentLength !== received) {
    return {
      ok: false,
      reason: `下载 ${received} 字节，与 Content-Length ${contentLength} 不一致`,
    }
  }
  if (contentRangeTotal != null && contentRangeTotal !== received) {
    return {
      ok: false,
      reason: `下载 ${received} 字节，与 Content-Range 总长度 ${contentRangeTotal} 不一致`,
    }
  }
  return { ok: true }
}
