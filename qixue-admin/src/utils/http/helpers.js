/**********************************
 * @FilePath: helpers.js
 * @Author: Ronnie Zhang
 * @LastEditor: Ronnie Zhang
 * @LastEditTime: 2023/12/04 22:46:22
 * @Email: zclzone@outlook.com
 * Copyright © 2023 Ronnie Zhang(大脸怪) | https://isme.top
 **********************************/

import { useAuthStore } from '@/store'

let isConfirming = false

function showReloginDialog(content) {
  if (isConfirming)
    return
  isConfirming = true
  $dialog.confirm({
    title: '提示',
    type: 'info',
    content,
    positiveText: '重新登录',
    negativeText: null,
    closable: false,
    maskClosable: false,
    closeOnEsc: false,
    confirm() {
      useAuthStore().logout()
      window.$message?.success('已退出登录')
      isConfirming = false
    },
  })
}

export function resolveResError(code, message, needTip = true) {
  switch (code) {
    case 401:
      if (!needTip)
        return
      showReloginDialog('登录已过期，请重新登录')
      return false
    case 11007:
    case 11008:
      if (!needTip)
        return
      showReloginDialog(`${message}，请重新登录`)
      return false
    case 403:
      message = '请求被拒绝'
      break
    case 404:
      message = '请求资源或接口不存在'
      break
    case 500:
      message = '服务器发生异常'
      break
    default:
      message = message ?? `【${code}】: 未知异常!`
      break
  }
  needTip && window.$message?.error(message)
  return message
}
