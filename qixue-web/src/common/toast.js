function ensure(id) {
  let el = document.getElementById(id)
  if (!el) {
    el = document.createElement('div')
    el.id = id
    document.body.appendChild(el)
  }
  return el
}

export function showToast(title, duration = 1500) {
  const el = ensure('qx-toast')
  el.textContent = title || ''
  el.style.cssText = 'position:fixed;left:50%;top:50%;transform:translate(-50%,-50%);z-index:100000;max-width:70%;padding:12px 16px;border-radius:8px;background:rgba(17,17,17,.78);color:#fff;font-size:14px;text-align:center;line-height:1.5;pointer-events:none;'
  clearTimeout(showToast._t)
  showToast._t = setTimeout(() => el.remove(), duration)
}

export function showLoading(title = '加载中') {
  const el = ensure('qx-loading')
  el.style.cssText = 'position:fixed;left:50%;top:50%;transform:translate(-50%,-50%);z-index:100001;min-width:96px;max-width:70%;padding:16px 18px;border-radius:8px;background:rgba(17,17,17,.78);color:#fff;font-size:14px;text-align:center;line-height:1.5;'
  el.textContent = title
}

export function hideLoading() {
  const el = document.getElementById('qx-loading')
  if (el) el.remove()
}
