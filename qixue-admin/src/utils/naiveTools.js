/**********************************
 * @FilePath: naiveTools.js
 * @Author: Ronnie Zhang
 * @LastEditor: Ronnie Zhang
 * @LastEditTime: 2023/12/04 22:45:20
 * @Email: zclzone@outlook.com
 * Copyright © 2023 Ronnie Zhang(大脸怪) | https://isme.top
 **********************************/

import * as NaiveUI from 'naive-ui'
import { naiveThemeOverrides } from '@/settings'
import { isNullOrUndef } from '@/utils'

export function setupMessage(NMessage) {
  class Message {
    static instance
    constructor() {
      if (Message.instance)
        return Message.instance
      Message.instance = this
      this.message = {}
      this.removeTimer = {}
    }

    removeMessage(key, duration = 5000) {
      this.removeTimer[key] && clearTimeout(this.removeTimer[key])
      this.removeTimer[key] = setTimeout(() => {
        this.message[key]?.destroy()
      }, duration)
    }

    destroy(key, duration = 200) {
      setTimeout(() => {
        this.message[key]?.destroy()
      }, duration)
    }

    showMessage(type, content, option = {}) {
      if (Array.isArray(content)) {
        return content.forEach(msg => NMessage[type](msg, option))
      }

      if (!option.key) {
        return NMessage[type](content, option)
      }

      const currentMessage = this.message[option.key]
      if (currentMessage) {
        currentMessage.type = type
        currentMessage.content = content
      }
      else {
        this.message[option.key] = NMessage[type](content, {
          ...option,
          duration: 0,
          onAfterLeave: () => {
            delete this.message[option.key]
          },
        })
      }
      this.removeMessage(option.key, option.duration)
    }

    loading(content, option) {
      this.showMessage('loading', content, option)
    }

    success(content, option) {
      this.showMessage('success', content, option)
    }

    error(content, option) {
      this.showMessage('error', content, option)
    }

    info(content, option) {
      this.showMessage('info', content, option)
    }

    warning(content, option) {
      this.showMessage('warning', content, option)
    }
  }

  return new Message()
}

export function setupDialog(NDialog) {
  const rawCreate = NDialog.create.bind(NDialog)
  const rawTyped = {
    info: NDialog.info.bind(NDialog),
    success: NDialog.success.bind(NDialog),
    warning: NDialog.warning.bind(NDialog),
    error: NDialog.error.bind(NDialog),
  }

  function resolvePositiveType(option = {}, dialogType = 'warning') {
    if (option.positiveButtonProps?.type)
      return option.positiveButtonProps.type
    if (option.danger || dialogType === 'error')
      return 'error'
    // 确认操作统一主色，避免 warning 弹窗自带的异常黄色按钮
    return 'primary'
  }

  function wrapPositiveClick(option = {}, getDialog) {
    const fn = option.onPositiveClick || option.confirm
    if (!fn)
      return undefined
    return async (e) => {
      const dialog = getDialog?.()
      if (dialog)
        dialog.loading = true
      try {
        return await fn(e)
      }
      finally {
        if (dialog)
          dialog.loading = false
      }
    }
  }

  function normalizeOption(option = {}, dialogType = 'warning', getDialog) {
    const {
      confirm: _confirm,
      cancel,
      danger: _danger,
      onPositiveClick: _onPositiveClick,
      onNegativeClick,
      onMaskClick,
      positiveButtonProps,
      ...rest
    } = option

    const showIcon = Object.prototype.hasOwnProperty.call(option, 'showIcon')
      ? option.showIcon
      : !isNullOrUndef(option.title)

    return {
      showIcon,
      positiveText: '确定',
      negativeText: '取消',
      class: 'apple-dialog',
      ...rest,
      positiveButtonProps: {
        type: resolvePositiveType(option, dialogType),
        ...positiveButtonProps,
      },
      onPositiveClick: wrapPositiveClick(option, getDialog),
      onNegativeClick: onNegativeClick || cancel,
      onMaskClick: onMaskClick || cancel,
    }
  }

  function openDialog(rawFn, option = {}, dialogType = 'warning') {
    const holder = { dialog: undefined }
    holder.dialog = rawFn(normalizeOption(option, dialogType, () => holder.dialog))
    return holder.dialog
  }

  ;['info', 'success', 'warning', 'error'].forEach((type) => {
    NDialog[type] = function (option = {}) {
      return openDialog(rawTyped[type], option, type)
    }
  })

  NDialog.create = function (option = {}) {
    const type = option.type || 'warning'
    return openDialog(rawCreate, option, type)
  }

  NDialog.confirm = function (option = {}) {
    const type = option.type || 'warning'
    return openDialog(rawTyped[type], option, type)
  }

  return NDialog
}

export function setupNaiveDiscreteApi() {
  const configProviderProps = computed(() => ({
    themeOverrides: {
      ...naiveThemeOverrides,
      LoadingBar: {
        height: '1.5px',
        colorLoading: 'rgba(49, 108, 114, 0.92)',
        colorError: '#D4534A',
      },
    },
  }))
  const { message, dialog, notification, loadingBar } = NaiveUI.createDiscreteApi(
    ['message', 'dialog', 'notification', 'loadingBar'],
    {
      configProviderProps,
      loadingBarProviderProps: {
        loadingBarStyle: {
          loading: {
            background: 'linear-gradient(90deg, rgba(49,108,114,0), rgba(49,108,114,0.5) 45%, #316C72)',
            boxShadow: '0 0 10px rgba(49, 108, 114, 0.28)',
          },
          error: {
            background: 'linear-gradient(90deg, rgba(212,83,74,0.2), #D4534A)',
            boxShadow: '0 0 10px rgba(212, 83, 74, 0.28)',
          },
        },
      },
    },
  )

  window.$loadingBar = loadingBar
  window.$notification = notification
  window.$message = setupMessage(message)
  window.$dialog = setupDialog(dialog)
}
