/** 用网站信息里的名称和 Logo 替换安装清单。 */

let appliedManifestKey = ''

export function applySiteBranding(sysInfo) {
  try {
    if (typeof document === 'undefined' || !document.head)
      return
    const cfg = sysInfo?.sysConfig || {}
    const name = String(cfg.name || '').trim()
    const rawLogo = String(cfg.logo || '').trim()
    let icon = ''
    if (rawLogo.startsWith('data:') || rawLogo.startsWith('blob:') || /^https?:\/\//i.test(rawLogo))
      icon = rawLogo
    else if (rawLogo) {
      const domain = String(cfg.assetDomain || '').trim().replace(/\/+$/, '')
      if (!domain)
        icon = absoluteUrl(rawLogo)
      else {
        const base = domain.includes('://') ? domain : `https://${domain}`
        icon = base + (rawLogo.startsWith('/') ? rawLogo : `/${rawLogo}`)
      }
    }
    if (!name && !icon)
      return
    const key = `${name}\n${icon}`
    if (key === appliedManifestKey)
      return
    appliedManifestKey = key
    const shortName = (name || '管理后台').slice(0, 12)
    const title = document.querySelector('meta[name="apple-mobile-web-app-title"]')
    if (title)
      title.setAttribute('content', shortName)
    if (icon)
      setLink('apple-touch-icon', icon)
    const iconType = /\.svg(\?|$)/i.test(icon)
      ? 'image/svg+xml'
      : (/\.jpe?g(\?|$)/i.test(icon) ? 'image/jpeg' : (/\.webp(\?|$)/i.test(icon) ? 'image/webp' : 'image/png'))
    const manifest = {
      id: absoluteUrl('./'),
      name: name || '管理后台',
      short_name: shortName,
      description: name || '管理后台',
      theme_color: '#316C72',
      background_color: '#f4f7f7',
      display: 'standalone',
      display_override: ['window-controls-overlay', 'standalone', 'minimal-ui'],
      orientation: 'any',
      lang: 'zh-CN',
      dir: 'ltr',
      start_url: absoluteUrl('./?source=pwa'),
      scope: absoluteUrl('./'),
      icons: icon
        ? [
            { src: icon, sizes: '192x192', type: iconType, purpose: 'any' },
            { src: icon, sizes: '512x512', type: iconType, purpose: 'any' },
          ]
        : [
            { src: absoluteUrl('pwa/icon-192.png'), sizes: '192x192', type: 'image/png', purpose: 'any' },
            { src: absoluteUrl('pwa/icon-512.png'), sizes: '512x512', type: 'image/png', purpose: 'any' },
            { src: absoluteUrl('pwa/icon-maskable-512.png'), sizes: '512x512', type: 'image/png', purpose: 'maskable' },
          ],
    }
    const blob = new Blob([JSON.stringify(manifest)], { type: 'application/manifest+json' })
    const url = URL.createObjectURL(blob)
    const link = setLink('manifest', url)
    const prev = link.dataset.dynamicManifest || ''
    link.dataset.dynamicManifest = url
    if (prev.startsWith('blob:'))
      URL.revokeObjectURL(prev)
  }
  catch {}
}

function absoluteUrl(path) {
  try {
    return new URL(path, document.baseURI).href
  }
  catch {
    return path
  }
}

function setLink(rel, href) {
  let link = document.querySelector(`link[rel="${rel}"]`)
  if (!link) {
    link = document.createElement('link')
    link.rel = rel
    document.head.appendChild(link)
  }
  link.href = href
  return link
}
