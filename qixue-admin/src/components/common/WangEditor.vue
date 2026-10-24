<script setup>
import { Editor, Toolbar } from '@wangeditor/editor-for-vue'
// 引入 css
import { onBeforeUnmount, ref, shallowRef, watch } from 'vue'
import api from '@/api'
import MediaPickerModal from './MediaPickerModal.vue'
import { bindMediaLibrary, ensureMediaLibraryMenu, MEDIA_LIBRARY_MENU_KEY, unbindMediaLibrary } from './wangMediaLibraryMenu'
import { getAssetUrl } from '@/utils/common'
import '@wangeditor/editor/dist/css/style.css'

ensureMediaLibraryMenu()

const props = defineProps({
  value: { type: String, default: '' },
  readOnly: { type: Boolean, default: false },
})

const emit = defineEmits(['update:value'])

// 编辑器实例，必须用 shallowRef
const editorRef = shallowRef()

watch(() => props.readOnly, (isReadOnly) => {
  if (isReadOnly) {
    editorRef.value.disable()
  }
  else {
    editorRef.value.enable()
  }
})

// 内容 HTML
const valueHtml = ref(props.value)
watch(() => valueHtml.value, (vl) => {
  emit('update:value', vl)
})
// 监听回显
watch(() => props.value, (vl) => {
  valueHtml.value = vl
})

const mode = 'default'
const toolbarConfig = {
  toolbarKeys: [
    'blockquote',
    'header1',
    'header2',
    'header3',
    '|',
    'bold',
    'underline',
    'italic',
    'through',
    'color',
    'bgColor',
    'clearStyle',
    '|',
    'bulletedList',
    'numberedList',
    'todo',
    'justifyLeft',
    'justifyRight',
    'justifyCenter',
    '|',
    'insertLink',
    'uploadImage',
    MEDIA_LIBRARY_MENU_KEY,
    'uploadVideo',
    'insertTable',
    'codeBlock',
    '|',
    'undo',
    'redo',
    '|',
    'fullScreen',
  ],
}

const editorConfig = {
  placeholder: '请输入内容...',
  MENU_CONF: {
    insertImage: {
      parseImageSrc: (src) => getAssetUrl(src) || src,
    },
    uploadImage: {
      base64LimitSize: 5 * 1024, // 5kb
      async customUpload(file, insertFn) {
        api.upload(file, 4).then((res) => {
          insertFn(res.data.url, '', '')
        })
      },
    },
    uploadVideo: {
      async customUpload(file, insertFn) {
        api.upload(file, 14).then((res) => {
          insertFn(res.data.url, '')
        })
      },
    },
  },
}

// 组件销毁时，也及时销毁编辑器
onBeforeUnmount(() => {
  const editor = editorRef.value
  if (editor == null)
    return
  unbindMediaLibrary(editor)
  editor.destroy()
})

function handleCreated(editor) {
  editorRef.value = editor // 记录 editor 实例，重要！
  bindMediaLibrary(editor, () => {
    libraryShow.value = true
  })
  if (props.readOnly) {
    editorRef.value.disable()
  }
}

const libraryShow = ref(false)
function insertFromLibrary(item) {
  const editor = editorRef.value
  if (!editor || editor.isDisabled?.())
    return
  const src = getAssetUrl(item.url) || item.url
  if (!src)
    return
  const poster = item.thumb ? (getAssetUrl(item.thumb) || item.thumb) : ''
  try {
    editor.focus()
    if (item.type === 'video') {
      editor.insertNode({
        type: 'video',
        src,
        poster,
        children: [{ text: '' }],
      })
    }
    else {
      editor.insertNode({
        type: 'image',
        src,
        alt: item.title || '',
        href: '',
        style: {},
        children: [{ text: '' }],
      })
    }
  }
  catch {
    if (item.type === 'video')
      editor.dangerouslyInsertHtml(`<video src="${src}" controls></video>`)
    else
      editor.dangerouslyInsertHtml(`<img src="${src}" alt="">`)
  }
  emit('update:value', editor.getHtml?.() || valueHtml.value)
}
</script>

<template>
  <div class="wang-editor">
    <Toolbar
      v-if="!props.readOnly"
      class="wang-editor__toolbar"
      :editor="editorRef"
      :default-config="toolbarConfig"
      mode="simple"
    />
    <Editor
      v-model="valueHtml"
      class="wang-editor__body"
      :default-config="editorConfig"
      :mode="mode"
      @on-created="handleCreated"
    />
    <MediaPickerModal
      v-model:show="libraryShow"
      type="image"
      :types="['image', 'video']"
      @select="insertFromLibrary"
    />
  </div>
</template>

<style scoped>
.wang-editor {
  --apple-ease: cubic-bezier(0.32, 0.72, 0, 1);
  height: 500px;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  border: 1px solid var(--ctl-hairline, rgba(var(--primary-color), 0.08));
  border-radius: 16px;
  background: var(--ctl-field, rgba(var(--primary-color), 0.04));
}

.wang-editor__toolbar {
  flex-shrink: 0;
  border-bottom: 1px solid var(--ctl-hairline, rgba(var(--primary-color), 0.08)) !important;
  background: rgb(255 255 255 / 42%) !important;
  padding: 4px 6px;
}

.wang-editor__toolbar :deep(.w-e-toolbar) {
  background: transparent !important;
  border: 0 !important;
}

.wang-editor__toolbar :deep(.w-e-bar-item button) {
  border-radius: 8px !important;
  transition:
    background 0.18s ease,
    transform 0.18s var(--apple-ease) !important;
}

.wang-editor__toolbar :deep(.w-e-bar-item button:hover) {
  background: var(--ctl-field-hover, rgba(var(--primary-color), 0.08)) !important;
}

.wang-editor__toolbar :deep(.w-e-bar-item button:active) {
  transform: scale(0.94);
}

.wang-editor__body {
  flex: 1;
  min-height: 0 !important;
  height: auto !important;
  overflow-y: hidden;
  background: transparent !important;
}

.wang-editor__body :deep(.w-e-text-container),
.wang-editor__body :deep(.w-e-scroll) {
  height: 100% !important;
  background: transparent !important;
}

.wang-editor__body :deep([data-slate-editor]) {
  min-height: 100%;
  padding: 16px 18px 48px !important;
  box-sizing: border-box;
  font-family: var(--font-sans);
  letter-spacing: -0.01em;
  line-height: 1.7;
  color: var(--ctl-text, #1a2332);
}

.wang-editor__body :deep(.w-e-text-placeholder) {
  color: var(--ctl-text-soft, #8b97a5) !important;
  font-style: normal !important;
  letter-spacing: -0.01em;
}

.w-e-full-screen-container {
  z-index: 999;
}
</style>
