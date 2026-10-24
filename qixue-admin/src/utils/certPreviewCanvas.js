/** Certificate template preview via Canvas. */

export const CERT_CANVAS_W = 800
export const CERT_CANVAS_H = 500
export const CERT_TEXT_PAD = 0
export const CERT_FONT_FAMILY = '"Noto Sans SC", "PingFang SC", "Microsoft YaHei", sans-serif'
export const CERT_FONT_WEIGHT = '500'

let fontReady = null

export function ensureCertFont() {
  if (fontReady)
    return fontReady
  fontReady = (async () => {
    if (!document?.fonts)
      return false
    try {
      await document.fonts.load(`${CERT_FONT_WEIGHT} 16px "Noto Sans SC"`)
      return document.fonts.check(`${CERT_FONT_WEIGHT} 16px "Noto Sans SC"`)
    }
    catch {
      return false
    }
  })()
  return fontReady
}

export function normalizeCertSpaces(text) {
  return String(text ?? '')
    .replace(/\u3000/g, '   ')
    .replace(/[\u00A0\u2000-\u200A\u202F\u205F]/g, ' ')
}

export function interpolateCertContent(content, values = {}) {
  const v = values || {}
  return normalizeCertSpaces(content).replace(
    /[\{\uFF5B]\s*(name|course|time|no)\s*[\}\uFF5D]/gi,
    (_, key) => String(v[String(key).toLowerCase()] ?? ''),
  )
}

function loadImage(src) {
  return new Promise((resolve, reject) => {
    if (!src) {
      reject(new Error('empty image src'))
      return
    }
    const img = new Image()
    if (!String(src).startsWith('blob:') && !String(src).startsWith('data:'))
      img.crossOrigin = 'anonymous'
    img.onload = () => resolve(img)
    img.onerror = () => reject(new Error(`图片加载失败: ${src}`))
    img.src = src
  })
}

function drawCover(ctx, img, tw, th) {
  const ir = img.width / Math.max(1, img.height)
  const tr = tw / th
  let dw
  let dh
  let dx
  let dy
  if (ir > tr) {
    dh = th
    dw = th * ir
    dx = (tw - dw) / 2
    dy = 0
  }
  else {
    dw = tw
    dh = tw / ir
    dx = 0
    dy = (th - dh) / 2
  }
  ctx.drawImage(img, dx, dy, dw, dh)
}

function drawContain(ctx, img, x, y, w, h) {
  const ir = img.width / Math.max(1, img.height)
  const br = w / Math.max(1, h)
  let dw
  let dh
  let dx
  let dy
  if (ir > br) {
    dw = w
    dh = w / ir
    dx = x
    dy = y + (h - dh) / 2
  }
  else {
    dh = h
    dw = h * ir
    dx = x + (w - dw) / 2
    dy = y
  }
  ctx.drawImage(img, dx, dy, dw, dh)
}

function pullLeadingIndent(ctx, text) {
  const m = String(text).match(/^\s+/)
  if (!m)
    return { body: text, indent: 0 }
  const ws = m[0]
  const body = text.slice(ws.length)
  if (!body)
    return { body: '', indent: 0 }
  const marker = '国'
  const indent = Math.max(0, Math.round(ctx.measureText(ws + marker).width - ctx.measureText(marker).width))
  return { body, indent }
}

function wrapText(ctx, text, maxWidth, firstLineIndent = 0) {
  const lines = []
  const paragraphs = String(text ?? '').split('\n')
  let isFirstLine = true
  for (const para of paragraphs) {
    if (!para) {
      lines.push('')
      isFirstLine = false
      continue
    }
    let line = ''
    for (const ch of para) {
      const test = line + ch
      const limit = (isFirstLine && firstLineIndent > 0)
        ? Math.max(8, maxWidth - firstLineIndent)
        : maxWidth
      if (line && ctx.measureText(test).width > limit) {
        lines.push(line)
        line = ch
        isFirstLine = false
      }
      else {
        line = test
      }
    }
    if (line) {
      lines.push(line)
      isFirstLine = false
    }
  }
  return lines.length ? lines : ['']
}

function drawTextInBox(ctx, text, box) {
  const fontSize = Number(box.fontSize) > 0
    ? Number(box.fontSize)
    : Math.max(12, Math.round(box.h * 0.42))
  const align = box.align || 'center'
  const pad = CERT_TEXT_PAD
  ctx.fillStyle = box.color || '#333333'
  ctx.font = `${CERT_FONT_WEIGHT} ${fontSize}px ${CERT_FONT_FAMILY}`
  ctx.textBaseline = 'middle'

  const normalized = normalizeCertSpaces(text)
  let body = normalized
  let indentPx = 0
  if (align === 'left') {
    const pulled = pullLeadingIndent(ctx, normalized)
    body = pulled.body
    indentPx = pulled.indent
  }

  const maxW = Math.max(8, box.w - pad * 2)
  const lines = wrapText(ctx, body, maxW, indentPx).slice(0, 6)
  const lineH = fontSize * 1.35
  const totalH = lines.length * lineH
  let y = box.y + (box.h - totalH) / 2 + lineH / 2
  for (let i = 0; i < lines.length; i++) {
    const line = lines[i]
    const lineIndent = (align === 'left' && i === 0) ? indentPx : 0
    let x
    if (align === 'left') {
      ctx.textAlign = 'left'
      x = box.x + pad + lineIndent
    }
    else if (align === 'right') {
      ctx.textAlign = 'right'
      x = box.x + box.w - pad
    }
    else {
      ctx.textAlign = 'center'
      x = box.x + box.w / 2
    }
    ctx.fillText(line, x, y)
    y += lineH
  }
}

/**
 * @param {{ bgSrc: string, boxes: Array, getText: (box) => string, getImageSrc: (box) => string, isImageBox: (box) => boolean }} opts
 * @returns {Promise<File>}
 */
export async function renderCertPreviewFile(opts) {
  const {
    bgSrc,
    boxes = [],
    getText,
    getImageSrc,
    isImageBox,
    width = CERT_CANVAS_W,
    height = CERT_CANVAS_H,
  } = opts

  await ensureCertFont()

  const canvas = document.createElement('canvas')
  canvas.width = width
  canvas.height = height
  const ctx = canvas.getContext('2d')
  if (!ctx)
    throw new Error('Canvas 不可用')

  ctx.fillStyle = '#ffffff'
  ctx.fillRect(0, 0, width, height)

  const bg = await loadImage(bgSrc)
  drawCover(ctx, bg, width, height)

  for (const box of boxes) {
    if (isImageBox(box)) {
      const src = getImageSrc(box)
      if (!src)
        continue
      try {
        const stamp = await loadImage(src)
        drawContain(ctx, stamp, box.x, box.y, box.w, box.h)
      }
      catch {
        // 单张印章失败不阻断整图
      }
      continue
    }
    const text = getText(box)
    if (text == null || text === '')
      continue
    drawTextInBox(ctx, text, box)
  }

  const blob = await new Promise((resolve, reject) => {
    canvas.toBlob((b) => {
      if (b)
        resolve(b)
      else
        reject(new Error('预览图导出失败'))
    }, 'image/png', 0.92)
  })

  return new File([blob], `cert-preview-${Date.now()}.png`, { type: 'image/png' })
}

/**
 * DOM 回退：把右侧预览板导出为 PNG（html2canvas）
 */
export async function captureDomPreviewFile(el) {
  if (!el)
    throw new Error('预览节点不存在')
  await ensureCertFont()
  const { default: html2canvas } = await import('html2canvas')
  const canvas = await html2canvas(el, {
    useCORS: true,
    allowTaint: false,
    backgroundColor: '#ffffff',
    scale: Math.min(2, window.devicePixelRatio || 1),
    logging: false,
  })
  const blob = await new Promise((resolve, reject) => {
    canvas.toBlob((b) => {
      if (b)
        resolve(b)
      else
        reject(new Error('预览截图失败'))
    }, 'image/png', 0.92)
  })
  return new File([blob], `cert-preview-${Date.now()}.png`, { type: 'image/png' })
}
