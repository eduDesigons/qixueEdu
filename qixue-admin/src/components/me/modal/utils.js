/**********************************
 * @Author: Ronnie Zhang
 * @LastEditor: Ronnie Zhang
 * @LastEditTime: 2024/01/13 17:41:26
 * @Email: zclzone@outlook.com
 * Copyright © 2023 Ronnie Zhang(大脸怪) | https://isme.top
 **********************************/

function getCss(element, key) {
  return element.currentStyle
    ? element.currentStyle[key]
    : window.getComputedStyle(element, null)[key]
}

export function initDrag(bar, box) {
  if (!bar || !box)
    return
  const params = {
    left: 0,
    top: 0,
    currentX: 0,
    currentY: 0,
    flag: false,
  }

  if (getCss(box, 'left') !== 'auto') {
    params.left = getCss(box, 'left')
  }
  if (getCss(box, 'top') !== 'auto') {
    params.top = getCss(box, 'top')
  }

  bar.style.cursor = 'move'
  bar.onmousedown = function (e) {
    params.flag = true
    e.preventDefault()
    params.currentX = e.clientX
    params.currentY = e.clientY
  }
  document.onmouseup = function () {
    params.flag = false
    if (getCss(box, 'left') !== 'auto') {
      params.left = getCss(box, 'left')
    }
    if (getCss(box, 'top') !== 'auto') {
      params.top = getCss(box, 'top')
    }
  }
  document.onmousemove = function (e) {
    if (e.target !== bar && !params.flag)
      return

    e.preventDefault()
    if (params.flag) {
      const nowX = e.clientX
      const nowY = e.clientY
      const disX = nowX - params.currentX
      const disY = nowY - params.currentY

      const left = Number.parseInt(params.left) + disX
      const top = Number.parseInt(params.top) + disY

      box.style.left = `${left}px`
      box.style.top = `${top}px`
    }
  }
}
