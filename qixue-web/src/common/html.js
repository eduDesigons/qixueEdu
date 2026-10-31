/** 富文本相对路径补全域名 */
export function rewriteHtml(html, domain) {
  const source = String(html || '').replace(/<script[\s\S]*?<\/script>/gi, '')
  const base = String(domain || '').trim().replace(/\/+$/, '')
  if (!base) return source
  return source.replace(/(src|href)=(["'])(?!https?:|data:|blob:|\/\/)([^"']+)/gi, (_, attr, quote, url) => {
    const path = url.startsWith('/') ? url : `/${url}`
    return `${attr}=${quote}${base}${path}`
  })
}
