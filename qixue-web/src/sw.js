import { cleanupOutdatedCaches, createHandlerBoundToURL, precacheAndRoute } from 'workbox-precaching'
import { NavigationRoute, registerRoute } from 'workbox-routing'
import { NetworkFirst } from 'workbox-strategies'
import { ExpirationPlugin } from 'workbox-expiration'
import { CacheableResponsePlugin } from 'workbox-cacheable-response'
import { createPartialResponse } from 'workbox-range-requests'
import { OFFLINE_VIDEO_PATH, isCompleteVideoMark, isMp4Url, mediaCacheNames, stableMediaKey } from './common/mediaCacheCore.js'

const NAMES = mediaCacheNames('qx-web')

precacheAndRoute(self.__WB_MANIFEST)
cleanupOutdatedCaches()

self.addEventListener('message', (event) => {
  if (event.data && event.data.type === 'SKIP_WAITING')
    self.skipWaiting()
})

self.addEventListener('activate', (event) => {
  event.waitUntil(self.clients.claim())
})

function isImageRequest(url, request) {
  if (isMp4Url(url.href))
    return false
  if (request.destination === 'image')
    return true
  return /\.(png|jpe?g|gif|webp|avif|svg|ico)$/i.test(url.pathname)
}

// 只处理离线播放用的站内地址。CDN 上的 MP4 不拦截，避免影响在线拖动和 Range。
registerRoute(
  ({ url, request }) => request.method === 'GET'
    && url.origin === self.location.origin
    && url.pathname === OFFLINE_VIDEO_PATH,
  async ({ request, url }) => {
    const key = url.searchParams.get('key') || ''
    const cache = await caches.open(NAMES.videos)
    const cached = key ? await cache.match(key) : null
    if (!cached || cached.status !== 200 || !isCompleteVideoMark(cached.headers))
      return new Response('', { status: 404 })
    if (request.headers.has('range'))
      return createPartialResponse(request, cached)
    return cached
  },
  'GET',
)

registerRoute(
  ({ url, request }) => request.method === 'GET' && isImageRequest(url, request),
  new NetworkFirst({
    cacheName: NAMES.images,
    networkTimeoutSeconds: 4,
    plugins: [
      {
        cacheKeyWillBeUsed: async ({ request }) => stableMediaKey(request.url),
      },
      new CacheableResponsePlugin({ statuses: [0, 200] }),
      new ExpirationPlugin({
        maxEntries: 300,
        maxAgeSeconds: 30 * 24 * 60 * 60,
        purgeOnQuotaError: true,
      }),
    ],
  }),
  'GET',
)

registerRoute(new NavigationRoute(createHandlerBoundToURL('index.html'), {
  denylist: [/^\/api/, /^\/webapp/, /^\/admin/],
}))
