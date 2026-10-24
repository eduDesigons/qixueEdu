import { cleanupOutdatedCaches, createHandlerBoundToURL, precacheAndRoute } from 'workbox-precaching'
import { NavigationRoute, registerRoute } from 'workbox-routing'
import { NetworkFirst } from 'workbox-strategies'
import { ExpirationPlugin } from 'workbox-expiration'
import { CacheableResponsePlugin } from 'workbox-cacheable-response'

precacheAndRoute(self.__WB_MANIFEST)
cleanupOutdatedCaches()

self.addEventListener('message', (event) => {
  if (event.data && event.data.type === 'SKIP_WAITING')
    self.skipWaiting()
})

self.addEventListener('activate', (event) => {
  event.waitUntil(self.clients.claim())
})

const AUTH_QUERY_KEYS = new Set([
  'token',
  'expires',
  'token_path',
  'signature',
  'sig',
  'auth',
  'auth_key',
  'policy',
])

function stableImageKey(input) {
  const url = new URL(input, self.location.href)
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

function isImageRequest(url, request) {
  if (/\.mp4$/i.test(url.pathname))
    return false
  if (request.destination === 'image')
    return true
  return /\.(png|jpe?g|gif|webp|avif|svg|ico)$/i.test(url.pathname)
}

// 管理端只缓存图片，不拦截、不缓存视频
registerRoute(
  ({ url, request }) => request.method === 'GET' && isImageRequest(url, request),
  new NetworkFirst({
    cacheName: 'qx-admin-images',
    networkTimeoutSeconds: 4,
    plugins: [
      {
        cacheKeyWillBeUsed: async ({ request }) => stableImageKey(request.url),
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
