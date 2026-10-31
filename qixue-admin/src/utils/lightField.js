/**
 * ColorOS 风格光场：鼠标位置驱动块内柔光与圆球视差
 * @see https://www.coloros.com/version/coloros16/
 */

const SELECTOR = [
  '.page-block',
  '.visual-block',
  '.info-block',
  '.entity-card',
  '.cmt-card',
  '.detail-head',
  '.panel',
  '.todo__row',
  '.bento__feat',
  '.bento__tile',
  '.hero',
].join(',')

let activeEl = null
let ticking = false
let pending = null
let reduceMotion = false

function applyLight(el, clientX, clientY) {
  const rect = el.getBoundingClientRect()
  if (!rect.width || !rect.height)
    return
  const nx = Math.min(1, Math.max(0, (clientX - rect.left) / rect.width))
  const ny = Math.min(1, Math.max(0, (clientY - rect.top) / rect.height))
  el.style.setProperty('--lf-x', `${(nx * 100).toFixed(2)}%`)
  el.style.setProperty('--lf-y', `${(ny * 100).toFixed(2)}%`)
  el.style.setProperty('--lf-nx', nx.toFixed(3))
  el.style.setProperty('--lf-ny', ny.toFixed(3))
  el.classList.add('is-lit')
}

function clearLight(el) {
  if (!el)
    return
  el.classList.remove('is-lit')
}

function onPointerMove(e) {
  if (reduceMotion)
    return
  pending = e
  if (ticking)
    return
  ticking = true
  requestAnimationFrame(() => {
    ticking = false
    const ev = pending
    pending = null
    if (!ev)
      return
    const el = ev.target?.closest?.(SELECTOR) || null
    if (activeEl && activeEl !== el)
      clearLight(activeEl)
    activeEl = el
    if (el)
      applyLight(el, ev.clientX, ev.clientY)
  })
}

export function setupLightField() {
  if (typeof window === 'undefined' || window.__qixueLightField)
    return
  window.__qixueLightField = true

  const mq = window.matchMedia('(prefers-reduced-motion: reduce)')
  reduceMotion = mq.matches
  mq.addEventListener?.('change', (ev) => {
    reduceMotion = ev.matches
    if (reduceMotion) {
      clearLight(activeEl)
      activeEl = null
    }
  })

  document.addEventListener('pointermove', onPointerMove, { passive: true })
  document.addEventListener('pointerleave', () => {
    clearLight(activeEl)
    activeEl = null
  }, true)
}
