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

/** H5 浏览器标题：只用 sysConfig.name，不拼接页面名 */
export function applySiteDocumentTitle(sysInfo) {
	try {
		if (typeof document === 'undefined')
			return
		const site = String((sysInfo && sysInfo.sysConfig && sysInfo.sysConfig.name) || '').trim()
		if (site)
			document.title = site
	}
	catch (e) {}
}

/** H5 标签页图标：用 sysConfig.logo */
export function applySiteFavicon(sysInfo) {
	try {
		if (typeof document === 'undefined' || !document.head)
			return
		const cfg = (sysInfo && sysInfo.sysConfig) || {}
		const href = resolveAssetUrl(cfg.logo, cfg.assetDomain)
		if (!href)
			return
		let link = document.getElementById('qixue-site-favicon')
		if (!link) {
			link = document.createElement('link')
			link.id = 'qixue-site-favicon'
			link.rel = 'icon'
			document.head.appendChild(link)
		}
		const type = /\.svg(\?|$)/i.test(href)
			? 'image/svg+xml'
			: (/\.png(\?|$)/i.test(href) ? 'image/png' : (/\.jpe?g(\?|$)/i.test(href) ? 'image/jpeg' : 'image/x-icon'))
		link.type = type
		// 加版本避免换 logo 后浏览器仍用旧缓存
		const sep = href.indexOf('?') >= 0 ? '&' : '?'
		link.href = href + sep + 'v=' + encodeURIComponent(String(cfg.name || href).slice(0, 32))
	}
	catch (e) {}
}

/** 同步站点名 + favicon */
export function applySiteBranding(sysInfo) {
	applySiteDocumentTitle(sysInfo)
	applySiteFavicon(sysInfo)
}
