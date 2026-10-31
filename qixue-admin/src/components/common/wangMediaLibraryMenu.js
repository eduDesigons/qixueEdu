import { Boot } from '@wangeditor/editor'

export const MEDIA_LIBRARY_MENU_KEY = 'mediaLibrary'

const openers = new WeakMap()

class MediaLibraryMenu {
  constructor() {
    this.title = '资源库'
    this.tag = 'button'
    this.iconSvg = '<svg viewBox="0 0 1024 1024"><path d="M64 128h360l96 96h440v672H64V128z"/></svg>'
  }

  getValue() {
    return ''
  }

  isActive() {
    return false
  }

  isDisabled(editor) {
    return !editor || !!editor.isDisabled?.()
  }

  exec(editor) {
    openers.get(editor)?.()
  }
}

let registered = false

export function ensureMediaLibraryMenu() {
  if (registered)
    return
  registered = true
  try {
    Boot.registerMenu({
      key: MEDIA_LIBRARY_MENU_KEY,
      factory() {
        return new MediaLibraryMenu()
      },
    })
  }
  catch {
  }
}

export function bindMediaLibrary(editor, open) {
  if (editor)
    openers.set(editor, open)
}

export function unbindMediaLibrary(editor) {
  if (editor)
    openers.delete(editor)
}
