/**********************************
 * @FilePath: common.js
 * @Author: Ronnie Zhang
 * @LastEditor: Ronnie Zhang
 * @LastEditTime: 2023/12/04 22:45:46
 * @Email: zclzone@outlook.com
 * Copyright © 2023 Ronnie Zhang(大脸怪) | https://isme.top
 **********************************/

import dayjs from 'dayjs'
import { useSysStore } from '@/store/modules/sys'

/**
 * Join asset base URL with a path (absolute URLs unchanged).
 */
export function getAssetUrl(url) {
  if (!url)
    return ''
  if (url.startsWith('data:'))
    return url
  // Vite 打包后的本地资源
  if (url.startsWith('/assets/'))
    return url

  // 完整 URL：localhost / 127.0.0.1 视为相对路径，改走资源域名
  if (/^https?:\/\//.test(url)) {
    try {
      const parsed = new URL(url)
      const host = parsed.hostname.toLowerCase()
      if (host === 'localhost' || host === '127.0.0.1' || host === '::1')
        url = `${parsed.pathname}${parsed.search}`
      else
        return url
    }
    catch {
      return url
    }
  }

  const domain = useSysStore().assetDomain
  if (!domain)
    return url

  const base = domain.includes('://')
    ? domain.replace(/\/$/, '')
    : `https://${domain.replace(/\/$/, '')}`
  return `${base}${url.startsWith('/') ? url : `/${url}`}`
}

/**
 * Rewrite relative asset URLs in rich HTML for display.
 */
export function formatRichContentHtml(html) {
  if (!html)
    return ''

  const resolveAttr = (match, attr, quote, url) => {
    if (!url || url.startsWith('#') || /^(mailto:|tel:|javascript:)/i.test(url))
      return match
    return `${attr}=${quote}${getAssetUrl(url)}${quote}`
  }

  return html
    .replace(/\b(src|poster)=("|')([^"']*)("|')/gi, resolveAttr)
    .replace(/\b(src|poster)=([^ \t>]+)/gi, (match, attr, url) => {
      if (!url || url.startsWith('"') || url.startsWith('\''))
        return match
      if (url.startsWith('#') || /^(mailto:|tel:|javascript:)/i.test(url))
        return match
      return `${attr}="${getAssetUrl(url)}"`
    })
}

export function formatDateTime(time = undefined, format = 'YYYY-MM-DD HH:mm:ss') {
  return dayjs(time).format(format)
}

export function formatDate(date = undefined, format = 'YYYY-MM-DD') {
  return formatDateTime(date, format)
}

export function formatRelativeTime(time) {
  if (!time)
    return ''
  const d = dayjs(time)
  if (!d.isValid())
    return String(time)
  const now = dayjs()
  if (d.isSame(now, 'day'))
    return d.format('HH:mm')
  if (d.isSame(now.subtract(1, 'day'), 'day'))
    return `昨天 ${d.format('HH:mm')}`
  if (d.isSame(now, 'year'))
    return d.format('M月D日')
  return d.format('YYYY-MM-DD')
}

/**
 * @param {Function} fn
 * @param {number} wait
 * @returns {Function}  节流函数
 *
 */
export function throttle(fn, wait) {
  let context, args
  let previous = 0

  return function (...argArr) {
    const now = +new Date()
    context = this
    args = argArr
    if (now - previous > wait) {
      fn.apply(context, args)
      previous = now
    }
  }
}

/**
 * @param {Function} method
 * @param {number} wait
 * @param {boolean} immediate
 * @return {*} 防抖函数
 */
export function debounce(method, wait, immediate) {
  let timeout
  return function (...args) {
    const context = this
    if (timeout) {
      clearTimeout(timeout)
    }
    // 立即执行需要两个条件，一是immediate为true，二是timeout未被赋值或被置为null
    if (immediate) {
      /**
       * 如果定时器不存在，则立即执行，并设置一个定时器，wait毫秒后将定时器置为null
       * 这样确保立即执行后wait毫秒内不会被再次触发
       */
      const callNow = !timeout
      timeout = setTimeout(() => {
        timeout = null
      }, wait)
      if (callNow) {
        method.apply(context, args)
      }
    }
    else {
      // 如果immediate为false，则函数wait毫秒后执行
      timeout = setTimeout(() => {
        /**
         * args是一个类数组对象，所以使用fn.apply
         * 也可写作method.call(context, ...args)
         */
        method.apply(context, args)
      }, wait)
    }
  }
}

/**
 * @param {number} time 毫秒数
 * @returns 睡一会儿，让子弹暂停一下
 */
export function sleep(time) {
  return new Promise(resolve => setTimeout(resolve, time))
}

/**
 * @param {HTMLElement} el
 * @param {Function} cb
 * @return {ResizeObserver}
 */
export function useResize(el, cb) {
  const observer = new ResizeObserver((entries) => {
    cb(entries[0].contentRect)
  })
  observer.observe(el)
  return observer
}
