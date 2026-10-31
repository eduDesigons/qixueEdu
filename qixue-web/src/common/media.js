/**
 * 原生选图 / 预览，不依赖 uni。
 */

export function pickImages({ count = 1, accept = 'image/*' } = {}) {
  return new Promise((resolve, reject) => {
    const input = document.createElement('input')
    input.type = 'file'
    input.accept = accept
    if (count > 1) input.multiple = true
    input.style.display = 'none'
    const cleanup = () => {
      input.remove()
    }
    input.onchange = () => {
      const files = [...(input.files || [])].slice(0, count)
      cleanup()
      if (!files.length) {
        reject(new Error('cancel'))
        return
      }
      resolve(files.map((file) => ({
        file,
        url: URL.createObjectURL(file),
        name: file.name || 'image.jpg',
      })))
    }
    input.oncancel = () => {
      cleanup()
      reject(new Error('cancel'))
    }
    document.body.appendChild(input)
    input.click()
  })
}

export function pickFiles({ count = 1, accept = '' } = {}) {
  return new Promise((resolve, reject) => {
    const input = document.createElement('input')
    input.type = 'file'
    if (accept) input.accept = accept
    if (count > 1) input.multiple = true
    input.style.display = 'none'
    const cleanup = () => input.remove()
    input.onchange = () => {
      const files = [...(input.files || [])].slice(0, count)
      cleanup()
      if (!files.length) {
        reject(new Error('cancel'))
        return
      }
      resolve(files.map((file) => ({
        file,
        url: URL.createObjectURL(file),
        name: file.name || 'file',
      })))
    }
    input.oncancel = () => {
      cleanup()
      reject(new Error('cancel'))
    }
    document.body.appendChild(input)
    input.click()
  })
}

export function previewImages(urls, current) {
  const list = (urls || []).filter(Boolean)
  if (!list.length) return
  let index = typeof current === 'number' ? current : list.indexOf(current)
  if (index < 0) index = 0

  const mask = document.createElement('div')
  mask.style.cssText = 'position:fixed;inset:0;z-index:100010;background:rgba(0,0,0,.92);display:flex;align-items:center;justify-content:center;'
  const img = document.createElement('img')
  img.style.cssText = 'max-width:100%;max-height:100%;object-fit:contain;'
  const show = () => { img.src = list[index] }
  show()
  mask.onclick = () => mask.remove()
  mask.appendChild(img)
  if (list.length > 1) {
    mask.onwheel = (e) => {
      e.preventDefault()
      index = (index + (e.deltaY > 0 ? 1 : -1) + list.length) % list.length
      show()
    }
  }
  document.body.appendChild(mask)
}

/** 把 File / Blob / blob: / data: 统一成可上传的 Blob */
export async function toUploadBlob(input) {
  if (input instanceof Blob) return input
  if (typeof input === 'string') {
    const res = await fetch(input)
    return res.blob()
  }
  throw new Error('invalid upload source')
}
