import { cleanupOutdatedCaches, createHandlerBoundToURL, precacheAndRoute } from 'workbox-precaching'
import { NavigationRoute, registerRoute } from 'workbox-routing'
import { CacheFirst, NetworkFirst } from 'workbox-strategies'
import { ExpirationPlugin } from 'workbox-expiration'
import { CacheableResponsePlugin } from 'workbox-cacheable-response'
import { createPartialResponse } from 'workbox-range-requests'
import { OFFLINE_AUDIO_PATH, OFFLINE_VIDEO_PATH, isAudioUrl, isCompleteMediaResponse, isContentApiUrl, isMp4Url, mediaCacheNames, stableMediaKey } from './common/mediaCacheCore.js'

const NAMES = mediaCacheNames('qx-web')

precacheAndRoute(self.__WB_MANIFEST)
cleanupOutdatedCaches()

self.addEventListener('message', (event) => {
  if (event.data && event.data.type === 'SKIP_WAITING')
    self.skipWaiting()
})

async function purgeIncompleteKind(kind) {
  const cache = await caches.open(kind === 'audio' ? NAMES.audio : NAMES.videos)
  const reqs = await cache.keys()
  await Promise.all(reqs.map(async (req) => {
    const cached = await cache.match(req)
    if (!isCompleteMediaResponse(cached))
      await cache.delete(req)
  }))
}

self.addEventListener('activate', (event) => {
  event.waitUntil(Promise.all([
    self.clients.claim(),
    purgeIncompleteKind('video'),
    purgeIncompleteKind('audio'),
  ]))
})

function authTag(header) {
  let hash = 2166136261
  for (let i = 0; i < header.length; i++) {
    hash ^= header.charCodeAt(i)
    hash = Math.imul(hash, 16777619)
  }
  return (hash >>> 0).toString(16)
}

function isImageRequest(url, request) {
  if (isMp4Url(url.href) || isAudioUrl(url.href))
    return false
  if (request.destination === 'image')
    return true
  return /\.(png|jpe?g|gif|webp|avif|svg|ico)$/i.test(url.pathname)
}

async function completeMedia(kind, key) {
  if (!key)
    return null
  const cache = await caches.open(kind === 'audio' ? NAMES.audio : NAMES.videos)
  const cached = await cache.match(key)
  if (isCompleteMediaResponse(cached))
    return cached
  if (cached)
    await cache.delete(key)
  return null
}

// 完整缓存命中时，按原始地址直接回本地文件，并支持拖动切段。没有完整文件则仍走网络。
registerRoute(
  ({ url, request }) => request.method === 'GET' && (isMp4Url(url.href) || isAudioUrl(url.href)),
  async ({ request, url }) => {
    const kind = isAudioUrl(url.href) ? 'audio' : 'video'
    const cached = await completeMedia(kind, stableMediaKey(url.href))
    if (!cached)
      return fetch(request)
    if (request.headers.has('range'))
      return createPartialResponse(request, cached)
    return cached
  },
  'GET',
)

registerRoute(
  ({ url, request }) => request.method === 'GET'
    && url.origin === self.location.origin
    && (url.pathname === OFFLINE_VIDEO_PATH || url.pathname === OFFLINE_AUDIO_PATH),
  async ({ request, url }) => {
    const kind = url.pathname === OFFLINE_AUDIO_PATH ? 'audio' : 'video'
    const cached = await completeMedia(kind, url.searchParams.get('key') || '')
    if (!cached)
      return new Response('', { status: 404 })
    if (request.headers.has('range'))
      return createPartialResponse(request, cached)
    return cached
  },
  'GET',
)

// 完整图片已在本地时直接用这份，有网也不再请求。没有则拉网络并写入，下次走本地。
registerRoute(
  ({ url, request }) => request.method === 'GET' && isImageRequest(url, request),
  new CacheFirst({
    cacheName: NAMES.images,
    plugins: [
      {
        cacheKeyWillBeUsed: async ({ request }) => stableMediaKey(request.url),
      },
      new CacheableResponsePlugin({ statuses: [0, 200] }),
      new ExpirationPlugin({
        purgeOnQuotaError: true,
      }),
    ],
  }),
  'GET',
)

registerRoute(
  ({ url, request }) => request.method === 'GET' && isContentApiUrl(url.href),
  new NetworkFirst({
    cacheName: NAMES.content,
    networkTimeoutSeconds: 3,
    plugins: [
      {
        cacheKeyWillBeUsed: async ({ request }) => {
          const url = new URL(request.url)
          url.searchParams.set('qx_auth', authTag(request.headers.get('authorization') || ''))
          return url.href
        },
      },
      new CacheableResponsePlugin({ statuses: [200] }),
      new ExpirationPlugin({
        maxEntries: 120,
        maxAgeSeconds: 3 * 24 * 60 * 60,
        purgeOnQuotaError: true,
      }),
    ],
  }),
  'GET',
)

registerRoute(new NavigationRoute(createHandlerBoundToURL('index.html'), {
  denylist: [/^\/api/, /^\/webapp/, /^\/admin/],
}))
