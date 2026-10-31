function resolveAssetUrl(src, assetDomain) {
	if (!src)
		return ''
	const s = String(src).trim()
	if (!s)
		return ''
	if (s.indexOf('http') === 0 || s.indexOf('data:') === 0 || s.indexOf('blob:') === 0)
		return s
	const domain = String(assetDomain || '').trim().replace(/\/+$/, '')
	if (!domain)
		return s
	return domain + (s.startsWith('/') ? s : '/' + s)
}

let appliedManifestKey = ''

/** Apply the site name, favicon, and install manifest from sysConfig. */
export function applySiteBranding(sysInfo) {
	try {
		if (typeof document === 'undefined')
			return
		const cfg = (sysInfo && sysInfo.sysConfig) || {}
		const name = String(cfg.name || '').trim()
		if (name)
			document.title = name
		if (!document.head)
			return
		const icon = resolveAssetUrl(cfg.logo, cfg.assetDomain)
		if (icon) {
			let link = document.getElementById('qixue-site-favicon')
			if (!link) {
				link = document.createElement('link')
				link.id = 'qixue-site-favicon'
				link.rel = 'icon'
				document.head.appendChild(link)
			}
			link.type = /\.svg(\?|$)/i.test(icon)
				? 'image/svg+xml'
				: (/\.png(\?|$)/i.test(icon) ? 'image/png' : (/\.jpe?g(\?|$)/i.test(icon) ? 'image/jpeg' : 'image/x-icon'))
			const sep = icon.indexOf('?') >= 0 ? '&' : '?'
			link.href = icon + sep + 'v=' + encodeURIComponent(String(name || icon).slice(0, 32))
		}
		if (!name && !icon)
			return
		const key = name + '\n' + icon
		if (key === appliedManifestKey)
			return
		appliedManifestKey = key
		const shortName = (name || '在线学习').slice(0, 12)
		const title = document.querySelector('meta[name="apple-mobile-web-app-title"]')
		if (title)
			title.setAttribute('content', shortName)
		if (icon)
			setLink('apple-touch-icon', icon)
		const iconType = /\.svg(\?|$)/i.test(icon)
			? 'image/svg+xml'
			: (/\.jpe?g(\?|$)/i.test(icon) ? 'image/jpeg' : (/\.webp(\?|$)/i.test(icon) ? 'image/webp' : 'image/png'))
		const manifest = {
			id: absoluteUrl('/'),
			name: name || '在线学习',
			short_name: shortName,
			description: name || '在线学习',
			theme_color: '#f8f8f8',
			background_color: '#f8f8f8',
			display: 'standalone',
			display_override: ['window-controls-overlay', 'standalone', 'minimal-ui', 'fullscreen'],
			orientation: 'any',
			lang: 'zh-CN',
			dir: 'ltr',
			start_url: absoluteUrl('/?source=pwa'),
			scope: absoluteUrl('/'),
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
		if (prev.indexOf('blob:') === 0)
			URL.revokeObjectURL(prev)
	}
	catch (e) {}
}

function absoluteUrl(path) {
	try {
		return new URL(path, document.baseURI).href
	}
	catch (e) {
		return path
	}
}

function setLink(rel, href) {
	let link = document.querySelector('link[rel="' + rel + '"]')
	if (!link) {
		link = document.createElement('link')
		link.rel = rel
		document.head.appendChild(link)
	}
	link.href = href
	return link
}
