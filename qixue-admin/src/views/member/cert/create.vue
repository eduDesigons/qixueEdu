<template>
  <CommonPage
    back
    fill
    show-headline
    :title="isEdit ? '编辑证书' : '新增证书'"
  >
    <template #action>
      <NButton v-if="bgSrc" secondary @click="fileInputRef?.click()">
        <i class="i-fe:image mr-4 text-14" />
        更换背景
      </NButton>
      <NButton type="primary" :loading="submitting" :disabled="!bgSrc" @click="submit">
        <i v-if="!submitting" class="i-fe:check mr-4 text-14" />
        保存
      </NButton>
    </template>

    <div class="cert-studio">
      <div v-if="loading" class="studio-loading">
        <span class="studio-loading__dot" />
        加载中
      </div>

      <input ref="fileInputRef" type="file" accept="image/*" class="sr-file" @change="onSelectFile">
      <input ref="imageInputRef" type="file" accept="image/png,image/jpeg,image/webp" class="sr-file" @change="onImageFile">
      <input ref="replaceImageInputRef" type="file" accept="image/png,image/jpeg,image/webp" class="sr-file" @change="onReplaceImage">

      <div class="studio-grid">
        <section class="studio-col studio-col--edit">
          <header class="studio-bar">
            <div class="studio-bar__left">
              <span class="studio-eyebrow">画布</span>
              <n-input
                v-if="bgSrc"
                v-model:value="certName"
                autofocus
                placeholder="证书名称"
                class="studio-name"
                maxlength="50"
                size="small"
              />
              <div v-if="bgSrc" class="studio-show">
                <span>上架</span>
                <n-switch v-model:value="isShow" :checked-value="1" :unchecked-value="0" size="small" />
              </div>
            </div>
            <div v-if="bgSrc" class="seg" role="group" aria-label="插入">
              <button
                type="button"
                class="seg__btn"
                draggable="true"
                @dragstart="onDragStart($event, { kind: 'text', label: '文案', content: '' })"
                @click="addTextBox({ label: '文案', content: '' })"
              >
                <i class="i-fe:type" />
                文案
              </button>
              <button
                type="button"
                class="seg__btn"
                draggable="true"
                @dragstart="onDragStart($event, { kind: 'image' })"
                @click="imageInputRef?.click()"
              >
                <i class="i-fe:image" />
                图片
              </button>
            </div>
          </header>

          <Transition name="soft-reveal" mode="out-in">
            <div
              v-if="!bgSrc"
              key="upload"
              class="drop-hero"
              :class="{ 'drop-hero--active': isDragOver }"
              @dragover.prevent="isDragOver = true"
              @dragleave="isDragOver = false"
              @drop.prevent="onDropBg"
              @click="fileInputRef?.click()"
            >
              <div class="drop-hero__glyph">
                <i class="i-fe:image" />
              </div>
              <p class="drop-hero__title">添加证书背景</p>
              <p class="drop-hero__desc">点击选择，或将图片拖到此处</p>
            </div>

            <div v-else key="board" class="board-wrap" @mousedown="onBoardBlankDown">
              <div class="board-frame">
                <div
                  ref="canvasRef"
                  class="board board--edit"
                  @drop="handleDrop"
                  @dragover.prevent
                >
                  <img :src="bgSrc" class="board__bg" alt="证书背景">
                  <VueDragResize
                    v-for="item in boxes"
                    :key="item.id"
                    class="board-shell"
                    :x="item.x"
                    :y="item.y"
                    :width="item.w"
                    :height="item.h"
                    :active="selectedId === item.id"
                    :resizable="true"
                    :draggable="true"
                    :parent-limitation="true"
                    :immediate="true"
                    @activated="selectedId = item.id"
                    @dragging="(rect) => updateBoxRect(item, rect)"
                    @drag-end="(rect) => updateBoxRect(item, rect)"
                    @resizing="(rect) => updateBoxRect(item, rect)"
                    @resize-end="(rect) => updateBoxRect(item, rect)"
                  >
                    <div
                      class="node"
                      :class="{
                        'node--image': isImageBox(item),
                        'node--text': isTextBox(item),
                        'node--on': selectedId === item.id,
                      }"
                      :style="boxItemStyle(item, false)"
                      @mousedown.stop="selectedId = item.id"
                    >
                      <template v-if="isImageBox(item)">
                        <img
                          v-if="imageDisplaySrc(item)"
                          :src="imageDisplaySrc(item)"
                          class="node__img"
                          alt="图片"
                          draggable="false"
                        >
                        <span v-else class="node__empty">图片</span>
                      </template>
                      <span
                        v-else
                        class="node__html"
                        v-html="formatVarHtml(boxRawContent(item) || '选中后在右下编辑')"
                      />
                      <button
                        type="button"
                        class="node__x"
                        title="删除"
                        @pointerdown.stop
                        @click.stop="removeBox(item.id)"
                      >
                        <i class="i-fe:x" />
                      </button>
                    </div>
                  </VueDragResize>
                </div>
              </div>
            </div>
          </Transition>
        </section>

        <aside class="studio-col studio-col--side">
          <div class="side-preview">
            <header class="side-section-head">
              <span class="side-section-head__title">预览</span>
              <span v-if="bgSrc" class="side-section-head__meta">{{ SAMPLE.name }} · {{ SAMPLE.course }}</span>
            </header>

            <div class="preview-stage">
              <div v-if="!bgSrc" class="preview-empty">
                <div class="preview-empty__ring">
                  <i class="i-fe:eye" />
                </div>
                <p class="preview-empty__title">暂无预览</p>
                <p class="preview-empty__desc">上传背景后即可查看成稿</p>
              </div>

              <div v-else class="preview-hero">
                <div class="preview-card">
                  <div class="preview-viewport">
                    <div ref="liveBoardRef" class="board board--live">
                      <img :src="bgSrc" class="board__bg" alt="预览">
                      <div
                        v-for="item in boxes"
                        :key="`pv-${item.id}`"
                        class="live-layer"
                        :style="previewLayerStyle(item)"
                      >
                        <img
                          v-if="isImageBox(item) && imageDisplaySrc(item)"
                          :src="imageDisplaySrc(item)"
                          class="live-layer__img"
                          alt=""
                        >
                        <span
                          v-else-if="isTextBox(item)"
                          class="live-layer__text"
                          :style="previewTextStyle(item)"
                        >{{ boxDisplayText(item) }}</span>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <div class="side-editor">
            <header class="side-section-head">
              <span class="side-section-head__title">编辑</span>
              <span v-if="selectedBox" class="side-section-head__meta">
                {{ isImageBox(selectedBox) ? '图片' : '文案' }}
                · {{ selectedBox.w }}×{{ selectedBox.h }}
              </span>
            </header>

            <div v-if="!selectedBox" class="side-editor__empty">
              <p>在左侧选中元素</p>
              <span>在此编辑文案与样式</span>
            </div>

            <div v-else class="editor-panel">
              <div class="edit-toolbar" role="toolbar" aria-label="样式工具栏">
                <template v-if="isTextBox(selectedBox)">
                  <NColorPicker
                    :value="selectedBox.color || '#333333'"
                    :show-alpha="false"
                    :modes="['hex']"
                    size="small"
                    class="edit-toolbar__color"
                    @update:value="onColorChange"
                  />
                  <div class="edit-toolbar__size">
                    <button type="button" class="tool-ico" title="减小字号" @click="onFontSizeChange((selectedBox.fontSize || 18) - 2)">
                      <span class="tool-ico__a tool-ico__a--sm">A</span>
                    </button>
                    <NInputNumber
                      :value="selectedBox.fontSize || 18"
                      :min="10"
                      :max="96"
                      :step="1"
                      size="small"
                      :show-button="false"
                      class="edit-toolbar__size-input"
                      @update:value="onFontSizeChange"
                    />
                    <button type="button" class="tool-ico" title="增大字号" @click="onFontSizeChange((selectedBox.fontSize || 18) + 2)">
                      <span class="tool-ico__a">A</span>
                    </button>
                  </div>
                  <div class="tool-cluster">
                    <button
                      type="button"
                      class="tool-ico"
                      :class="{ 'tool-ico--on': selectedBox.align === 'left' }"
                      title="左对齐"
                      @click="selectedBox.align = 'left'"
                    >
                      <i class="i-fe:align-left" />
                    </button>
                    <button
                      type="button"
                      class="tool-ico"
                      :class="{ 'tool-ico--on': !selectedBox.align || selectedBox.align === 'center' }"
                      title="居中"
                      @click="selectedBox.align = 'center'"
                    >
                      <i class="i-fe:align-center" />
                    </button>
                    <button
                      type="button"
                      class="tool-ico"
                      :class="{ 'tool-ico--on': selectedBox.align === 'right' }"
                      title="右对齐"
                      @click="selectedBox.align = 'right'"
                    >
                      <i class="i-fe:align-right" />
                    </button>
                  </div>
                  <div class="edit-toolbar__divider" />
                </template>

                <div class="tool-cluster">
                  <button type="button" class="tool-ico" title="水平居中" @click="alignBox('h')">
                    <i class="i-material-symbols:align-horizontal-center" />
                  </button>
                  <button type="button" class="tool-ico" title="垂直居中" @click="alignBox('v')">
                    <i class="i-material-symbols:align-vertical-center" />
                  </button>
                </div>

                <button
                  v-if="isImageBox(selectedBox)"
                  type="button"
                  class="tool-text"
                  @click="replaceImageInputRef?.click()"
                >
                  更换图片
                </button>
              </div>

              <template v-if="isTextBox(selectedBox)">
                <div
                  class="editor-preview"
                  v-html="formatVarHtml(selectedBox.content || '')"
                />
                <n-input
                  ref="textInputRef"
                  v-model:value="selectedBox.content"
                  type="textarea"
                  :rows="4"
                  maxlength="160"
                  show-count
                  placeholder="输入文案…"
                  class="inspector__textarea"
                />
                <div class="var-row">
                  <button
                    v-for="v in variables"
                    :key="v.token"
                    type="button"
                    class="var-chip"
                    :data-key="v.key"
                    @click="insertPlaceholder(v.token)"
                  >
                    {{ v.label }}
                  </button>
                </div>
              </template>
            </div>
          </div>
        </aside>
      </div>
    </div>
  </CommonPage>
</template>

<script setup>
import { NButton, NColorPicker, NInputNumber } from 'naive-ui'
import { computed, nextTick, onMounted, onUnmounted, ref } from 'vue'
import VueDragResize from 'vue-drag-resize-next'
import { useRoute, useRouter } from 'vue-router'
import api from '@/api'
import { CommonPage } from '@/components/index.js'
import { useUnsavedGuard } from '@/composables'
import { getAssetUrl } from '@/utils'
import {
  captureDomPreviewFile,
  ensureCertFont,
  interpolateCertContent,
  renderCertPreviewFile,
} from '@/utils/certPreviewCanvas'
import 'vue-drag-resize-next/lib/style.css'

defineOptions({ name: 'CertCreate' })

const route = useRoute()
const router = useRouter()

const editId = computed(() => route.query.id || '')
const isEdit = computed(() => !!editId.value)

const FIELD_KEYS = ['name', 'course', 'time', 'no']

const variables = [
  { key: 'name', token: '{name}', label: '姓名', icon: 'i-fe:user' },
  { key: 'course', token: '{course}', label: '课程', icon: 'i-fe:book-open' },
  { key: 'time', token: '{time}', label: '日期', icon: 'i-fe:calendar' },
  { key: 'no', token: '{no}', label: '编号', icon: 'i-fe:hash' },
]

const VAR_LABEL = Object.fromEntries(variables.map(v => [v.token, v.label]))

const certName = ref('')
const isShow = ref(1)
const liveBoardRef = ref(null)
const boxes = ref([])
const selectedId = ref('')
const CANVAS_W = 800
const CANVAS_H = 500
const SAMPLE = {
  name: '张三',
  course: '示例课程',
  time: `${new Date().getFullYear()}年${new Date().getMonth() + 1}月${new Date().getDate()}日`,
  no: 'QX202608210001',
}

const textInputRef = ref(null)
const selectedBox = computed(() => boxes.value.find(b => b.id === selectedId.value))

function onBoardBlankDown(e) {
  if (e.target === e.currentTarget || e.target?.classList?.contains('board__bg') || e.target?.classList?.contains('board') || e.target?.classList?.contains('board-frame')) {
    selectedId.value = ''
  }
}

function uid(prefix = 'b') {
  return `${prefix}_${Date.now()}_${Math.random().toString(36).slice(2, 7)}`
}

function boxType(box) {
  if (!box)
    return ''
  if (box.type === 'image' || box.type === 'stamp')
    return 'image'
  if (box.type === 'field' || box.type === 'text')
    return box.type
  const key = String(box.key || '')
  if (key.startsWith('text_') || key.startsWith('text-'))
    return 'text'
  if (
    key.startsWith('stamp_') || key.startsWith('stamp-')
    || key.startsWith('image_') || key.startsWith('image-')
  ) {
    return 'image'
  }
  if (FIELD_KEYS.includes(key))
    return 'field'
  return 'text'
}

function isTextBox(box) {
  const t = boxType(box)
  return t === 'text' || t === 'field'
}

function isImageBox(box) {
  return boxType(box) === 'image'
}

function escapeHtml(s) {
  return String(s || '')
    .replaceAll('&', '&amp;')
    .replaceAll('<', '&lt;')
    .replaceAll('>', '&gt;')
    .replaceAll('"', '&quot;')
}

function formatVarHtml(content) {
  const raw = escapeHtml(content || '')
  if (!raw)
    return '<span class="var-placeholder">输入文案…</span>'
  return raw.replace(/\{(name|course|time|no)\}/g, (m) => {
    const key = m.slice(1, -1)
    const label = VAR_LABEL[m] || key
    return `<span class="var-token" data-key="${key}">${label}</span>`
  })
}

function interpolate(content, values = SAMPLE) {
  return interpolateCertContent(content, values)
}

function boxRawContent(box) {
  if (boxType(box) === 'field')
    return `{${box.key}}`
  return box.content || ''
}

function boxDisplayText(box) {
  return interpolate(boxRawContent(box)) || ' '
}

function imageDisplaySrc(box) {
  if (box.previewUrl)
    return box.previewUrl
  if (box.src)
    return getAssetUrl(box.src)
  return ''
}

function boxItemStyle(box, _isPreview = false) {
  if (isImageBox(box))
    return undefined
  const fontSize = Number(box.fontSize) > 0
    ? Number(box.fontSize)
    : Math.max(12, Math.round(box.h * 0.42))
  return {
    color: box.color || '#333333',
    textAlign: box.align || 'center',
    fontSize: `${fontSize}px`,
  }
}

function previewLayerStyle(box) {
  return {
    left: `${(box.x / CANVAS_W) * 100}%`,
    top: `${(box.y / CANVAS_H) * 100}%`,
    width: `${(box.w / CANVAS_W) * 100}%`,
    height: `${(box.h / CANVAS_H) * 100}%`,
  }
}

function previewTextStyle(box) {
  const fontSize = Number(box.fontSize) > 0
    ? Number(box.fontSize)
    : Math.max(12, Math.round(box.h * 0.42))
  return {
    color: box.color || '#333333',
    textAlign: box.align || 'center',
    // board--live 以 1em = 1 设计像素（相对 500 高）
    fontSize: `${fontSize}em`,
  }
}

function clampBox(item) {
  const minW = isImageBox(item) ? 48 : 40
  const minH = isImageBox(item) ? 48 : 24
  item.w = Math.min(Math.max(minW, item.w), CANVAS_W)
  item.h = Math.min(Math.max(minH, item.h), CANVAS_H)
  item.x = Math.min(Math.max(0, item.x), CANVAS_W - item.w)
  item.y = Math.min(Math.max(0, item.y), CANVAS_H - item.h)
}

function alignBox(axis) {
  const box = selectedBox.value
  if (!box)
    return
  if (axis === 'h')
    box.x = Math.round((CANVAS_W - box.w) / 2)
  else if (axis === 'v')
    box.y = Math.round((CANVAS_H - box.h) / 2)
  clampBox(box)
}

function onColorChange(color) {
  if (selectedBox.value && isTextBox(selectedBox.value))
    selectedBox.value.color = color
}

function onFontSizeChange(size) {
  const box = selectedBox.value
  if (!box || !isTextBox(box))
    return
  const n = Number(size)
  if (!Number.isFinite(n))
    return
  box.fontSize = Math.min(96, Math.max(10, Math.round(n)))
  // 框高度随字号略调，保证单行可读
  const needH = Math.max(box.h, Math.round(box.fontSize * 1.6))
  if (needH > box.h) {
    box.h = Math.min(CANVAS_H - box.y, needH)
    clampBox(box)
  }
}

function insertPlaceholder(token) {
  const box = selectedBox.value
  if (!box || !isTextBox(box))
    return
  if (boxType(box) === 'field') {
    box.type = 'text'
    box.content = `{${box.key}}`
    box.key = uid('text')
  }
  const el = textInputRef.value?.textareaElRef || textInputRef.value?.$el?.querySelector?.('textarea')
  const value = box.content || ''
  if (el && typeof el.selectionStart === 'number') {
    const start = el.selectionStart
    const end = el.selectionEnd
    box.content = `${value.slice(0, start)}${token}${value.slice(end)}`
    nextTick(() => {
      el.focus()
      const pos = start + token.length
      el.setSelectionRange(pos, pos)
    })
  }
  else {
    box.content = `${value}${token}`
    nextTick(() => textInputRef.value?.focus?.())
  }
}

function onKeydown(e) {
  const tag = (e.target?.tagName || '').toLowerCase()
  if (tag === 'input' || tag === 'textarea')
    return
  if (!selectedId.value)
    return

  const box = selectedBox.value
  if (!box)
    return

  const step = e.shiftKey ? 10 : 1
  if (e.key === 'ArrowLeft') {
    e.preventDefault()
    box.x -= step
    clampBox(box)
    return
  }
  if (e.key === 'ArrowRight') {
    e.preventDefault()
    box.x += step
    clampBox(box)
    return
  }
  if (e.key === 'ArrowUp') {
    e.preventDefault()
    box.y -= step
    clampBox(box)
    return
  }
  if (e.key === 'ArrowDown') {
    e.preventDefault()
    box.y += step
    clampBox(box)
    return
  }
  if (e.key === 'Delete' || e.key === 'Backspace') {
    e.preventDefault()
    removeBox(selectedId.value)
  }
}

const bgSrc = ref('')
const bgFile = ref(null)
const templateImageUrl = ref('')
const canvasRef = ref(null)
const fileInputRef = ref(null)
const imageInputRef = ref(null)
const replaceImageInputRef = ref(null)
const submitting = ref(false)
const isDragOver = ref(false)
const loading = ref(false)
const pendingImageAt = ref(null)

const { snapshot, allowLeave } = useUnsavedGuard(() => ({
  name: certName.value,
  is_show: isShow.value,
  boxes: boxes.value.map(b => ({
    id: b.id,
    type: boxType(b),
    key: b.key,
    content: boxRawContent(b),
    src: b.src,
    hasStampFile: !!b.srcFile,
    x: b.x,
    y: b.y,
    w: b.w,
    h: b.h,
    color: b.color,
    align: b.align,
    fontSize: b.fontSize,
  })),
  hasFile: !!bgFile.value,
  templateImageUrl: templateImageUrl.value,
}))

/** 加载时把旧 field 升为「仅含变量」的文案，统一编辑模型 */
function normalizeLoadedBox(p, i) {
  const key = p.key || `item_${i}`
  let type = p.type || (
    FIELD_KEYS.includes(key)
      ? 'field'
      : (String(key).startsWith('stamp') ? 'stamp' : 'text')
  )

  if (type === 'field' && FIELD_KEYS.includes(key)) {
    return {
      id: `edit-${i}-${key}`,
      type: 'text',
      key: uid('text'),
      label: variables.find(v => v.key === key)?.label || key,
      content: `{${key}}`,
      x: p.x ?? 0,
      y: p.y ?? 0,
      w: p.w ?? 200,
      h: p.h ?? 40,
      color: p.color || '#333333',
      align: p.align || 'center',
      fontSize: Number(p.fontSize) > 0 ? Number(p.fontSize) : Math.max(14, Math.round((p.h ?? 40) * 0.42)),
    }
  }

  if (type === 'stamp' || type === 'image' || String(key).startsWith('image') || String(key).startsWith('stamp')) {
    return {
      id: `edit-${i}-${key}`,
      type: 'image',
      key,
      label: '图片',
      src: p.src || '',
      x: p.x ?? 0,
      y: p.y ?? 0,
      w: p.w ?? 100,
      h: p.h ?? 100,
    }
  }

  return {
    id: `edit-${i}-${key}`,
    type: 'text',
    key,
    label: '文案',
    content: p.content || '',
    x: p.x ?? 0,
    y: p.y ?? 0,
    w: p.w ?? 200,
    h: p.h ?? 40,
    color: p.color || '#333333',
    align: p.align || 'center',
    fontSize: Number(p.fontSize) > 0 ? Number(p.fontSize) : Math.max(14, Math.round((p.h ?? 40) * 0.42)),
  }
}

onMounted(() => {
  ensureCertFont()
  window.addEventListener('keydown', onKeydown)
  if (isEdit.value && editId.value) {
    loading.value = true
    api.cert.detail(editId.value)
      .then((res) => {
        const cert = res.data
        if (cert) {
          certName.value = cert.name || ''
          isShow.value = Number(cert.is_show) === 0 ? 0 : 1
          templateImageUrl.value = cert.template_image || ''
          bgSrc.value = getAssetUrl(cert.template_image || '')
          boxes.value = (cert.params || []).map(normalizeLoadedBox)
        }
      })
      .finally(() => {
        loading.value = false
        nextTick(snapshot)
      })
  }
  else {
    nextTick(snapshot)
  }
})

onUnmounted(() => {
  window.removeEventListener('keydown', onKeydown)
})

function onDragStart(e, item) {
  e.dataTransfer.setData('application/json', JSON.stringify(item))
  e.dataTransfer.effectAllowed = 'copy'
}

function updateBoxRect(item, rect) {
  if (!rect)
    return
  item.x = Math.round(rect.x)
  item.y = Math.round(rect.y)
  if (rect.width != null)
    item.w = Math.round(rect.width)
  if (rect.height != null)
    item.h = Math.round(rect.height)
  clampBox(item)
}

function dropPoint(e) {
  const rect = canvasRef.value?.getBoundingClientRect()
  if (!rect)
    return null
  return { x: e.clientX - rect.left, y: e.clientY - rect.top }
}

function createTextBox(raw, x, y) {
  const content = raw.content ?? ''
  const isShort = !content || /^\{(name|course|time|no)\}$/.test(content)
  const fontSize = isShort ? 22 : 16
  const box = {
    id: uid('text'),
    type: 'text',
    key: uid('text'),
    label: raw.label || '文案',
    content,
    x: Math.max(0, x - (isShort ? 100 : 280)),
    y: Math.max(0, y - (isShort ? 20 : 32)),
    w: isShort ? 220 : 560,
    h: Math.round(fontSize * 1.8),
    color: '#333333',
    align: 'center',
    fontSize,
  }
  clampBox(box)
  return box
}

function createImageBox(raw, x, y) {
  const box = {
    id: uid('image'),
    type: 'image',
    key: uid('image'),
    label: raw.label || '图片',
    src: raw.src || '',
    previewUrl: raw.previewUrl || '',
    srcFile: raw.srcFile || null,
    x: Math.max(0, x - 50),
    y: Math.max(0, y - 50),
    w: 100,
    h: 100,
  }
  clampBox(box)
  return box
}

function handleDrop(e) {
  // 从桌面拖入图片文件
  const file = e.dataTransfer?.files?.[0]
  if (file?.type?.startsWith('image/')) {
    const pt = dropPoint(e) || { x: CANVAS_W / 2, y: CANVAS_H / 2 }
    readStampFile(file).then((previewUrl) => {
      const box = createImageBox({
        label: '图片',
        previewUrl,
        srcFile: file,
      }, pt.x, pt.y)
      boxes.value.push(box)
      selectedId.value = box.id
    }).catch(() => $message.warning('图片读取失败'))
    return
  }

  let raw
  try {
    raw = JSON.parse(e.dataTransfer.getData('application/json') || '{}')
  }
  catch {
    return
  }
  if (!raw?.kind)
    return
  const pt = dropPoint(e)
  if (!pt)
    return

  let box = null
  if (raw.kind === 'text') {
    box = createTextBox(raw, pt.x, pt.y)
  }
  else if (raw.kind === 'image' || raw.kind === 'stamp') {
    if (raw.src || raw.previewUrl || raw.srcFile) {
      box = createImageBox(raw, pt.x, pt.y)
    }
    else {
      pendingImageAt.value = pt
      imageInputRef.value?.click()
      return
    }
  }

  if (!box)
    return
  boxes.value.push(box)
  selectedId.value = box.id
  if (isTextBox(box))
    nextTick(() => textInputRef.value?.focus?.())
}

function addTextBox(item) {
  if (!bgSrc.value) {
    $message.warning('请先添加证书背景')
    return
  }
  const box = createTextBox(item || { label: '文案', content: '' }, CANVAS_W / 2, 210)
  boxes.value.push(box)
  selectedId.value = box.id
  nextTick(() => textInputRef.value?.focus?.())
}

function removeBox(id) {
  boxes.value = boxes.value.filter(b => b.id !== id)
  if (selectedId.value === id)
    selectedId.value = ''
}

function applyBgFile(file) {
  if (!file)
    return
  if (!file.type.startsWith('image/')) {
    $message.warning('请选择图片文件')
    return
  }
  bgFile.value = file
  const reader = new FileReader()
  reader.onload = (ev) => {
    bgSrc.value = ev.target?.result ?? ''
  }
  reader.readAsDataURL(file)
}

function onSelectFile(e) {
  const file = e.target?.files?.[0]
  if (!file)
    return
  applyBgFile(file)
  e.target.value = ''
}

function onDropBg(e) {
  isDragOver.value = false
  const file = e.dataTransfer?.files?.[0]
  if (!file)
    return
  applyBgFile(file)
}

function readStampFile(file) {
  return new Promise((resolve, reject) => {
    if (!file?.type?.startsWith('image/')) {
      reject(new Error('请选择图片'))
      return
    }
    const reader = new FileReader()
    reader.onload = (ev) => resolve(ev.target?.result || '')
    reader.onerror = () => reject(new Error('读取失败'))
    reader.readAsDataURL(file)
  })
}

async function onImageFile(e) {
  const file = e.target?.files?.[0]
  e.target.value = ''
  if (!file)
    return
  try {
    const previewUrl = await readStampFile(file)
    if (!bgSrc.value) {
      $message.warning('请先添加证书背景')
      return
    }
    const at = pendingImageAt.value || { x: CANVAS_W - 140, y: CANVAS_H - 140 }
    pendingImageAt.value = null
    const box = createImageBox({
      label: '图片',
      previewUrl,
      srcFile: file,
    }, at.x, at.y)
    boxes.value.push(box)
    selectedId.value = box.id
  }
  catch (err) {
    $message.warning(err?.message || '图片读取失败')
  }
}

async function onReplaceImage(e) {
  const file = e.target?.files?.[0]
  e.target.value = ''
  const box = selectedBox.value
  if (!file || !box || !isImageBox(box))
    return
  try {
    box.previewUrl = await readStampFile(file)
    box.srcFile = file
    box.src = ''
  }
  catch (err) {
    $message.warning(err?.message || '图片读取失败')
  }
}

function serializeParams() {
  return boxes.value.map((b) => {
    if (isImageBox(b)) {
      return {
        type: 'image',
        key: b.key,
        src: b.src || '',
        x: Math.round(b.x),
        y: Math.round(b.y),
        w: Math.round(b.w),
        h: Math.round(b.h),
      }
    }
    return {
      type: 'text',
      key: b.key,
      content: boxRawContent(b),
      x: Math.round(b.x),
      y: Math.round(b.y),
      w: Math.round(b.w),
      h: Math.round(b.h),
      color: b.color || '#333333',
      align: b.align || 'center',
      fontSize: Number(b.fontSize) > 0 ? Number(b.fontSize) : 18,
    }
  })
}

async function uploadPendingImages() {
  for (const box of boxes.value) {
    if (!isImageBox(box) || !box.srcFile)
      continue
    const res = await api.upload(box.srcFile, 1)
    const url = res.data?.url || res.data?.path
    if (!url)
      throw new Error('图片上传失败')
    box.src = url
    box.srcFile = null
    box.previewUrl = ''
  }
}

/** 前端 Canvas 合成预览图并上传，失败时回退 DOM 截图 */
async function uploadPreviewImage() {
  await nextTick()
  let file
  try {
    file = await renderCertPreviewFile({
      bgSrc: bgSrc.value,
      boxes: boxes.value,
      isImageBox,
      getText: boxDisplayText,
      getImageSrc: imageDisplaySrc,
    })
  }
  catch {
    file = await captureDomPreviewFile(liveBoardRef.value)
  }
  const res = await api.upload(file, 1)
  const url = res.data?.url || res.data?.path
  if (!url)
    throw new Error('预览图上传失败')
  return url
}

function submit() {
  if (submitting.value)
    return
  const name = (certName.value || '').trim()
  if (!name) {
    $message.warning('请填写证书名称')
    return
  }

  submitting.value = true

  function doSave(templateUrl, previewUrl) {
    return api.cert.save({
      name,
      params: serializeParams(),
      template_image: templateUrl,
      preview_image: previewUrl || undefined,
      is_show: isShow.value,
      id: editId.value || undefined,
    })
  }

  const run = async () => {
    await uploadPendingImages()
    let templateUrl = templateImageUrl.value
    if (bgFile.value) {
      const res = await api.upload(bgFile.value, 1)
      const url = res.data?.url || res.data?.path
      if (!url)
        throw new Error('背景图上传失败')
      templateUrl = url
      templateImageUrl.value = url
      // 背景已上传到 CDN，预览合成改用 CDN 地址（便于 CORS）
      bgSrc.value = getAssetUrl(url)
      bgFile.value = null
    }
    if (!templateUrl)
      throw new Error('请先选择证书背景图')

    await nextTick()
    let previewUrl = ''
    try {
      previewUrl = await uploadPreviewImage()
    }
    catch (err) {
      console.warn('证书预览图生成失败，将仅保存模板', err)
    }
    await doSave(templateUrl, previewUrl)
  }

  run()
    .then(() => {
      if (isEdit.value)
        $message.success('模板已更新（已发学员证书不受影响）')
      else
        $message.success('保存成功')
      allowLeave()
      router.push('/course/cert')
    })
    .catch((err) => {
      $message.error(err?.message || '保存失败')
    })
    .finally(() => {
      submitting.value = false
    })
}
</script>

<style scoped>
.cert-studio {
  position: relative;
  height: 100%;
  min-height: 0;
  font-family: var(--font-sans);
  letter-spacing: -0.01em;
  color: var(--ctl-text, rgba(0, 0, 0, 0.85));
  animation: soft-in 0.36s var(--ctl-ease, cubic-bezier(0.32, 0.72, 0, 1)) both;
}

.sr-file {
  position: absolute;
  width: 0;
  height: 0;
  opacity: 0;
  pointer-events: none;
}

.studio-loading {
  position: absolute;
  inset: 0;
  z-index: 40;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10px;
  background: rgb(252 252 252 / 86%);
  backdrop-filter: saturate(1.2) blur(18px);
  -webkit-backdrop-filter: saturate(1.2) blur(18px);
  font-size: 13px;
  font-weight: 520;
  color: var(--ctl-text-soft, rgba(0, 0, 0, 0.5));
}

.studio-loading__dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: rgb(var(--primary-color, 49, 108, 114));
  animation: pulse-dot 0.9s ease infinite alternate;
}

.studio-grid {
  display: grid;
  grid-template-columns: minmax(0, 1fr) minmax(440px, 520px);
  gap: 14px;
  height: 100%;
  min-height: 0;
}

.studio-col {
  display: flex;
  flex-direction: column;
  min-width: 0;
  min-height: 0;
  border-radius: var(--ctl-radius-lg, 16px);
  background: var(--ctl-glass-strong, #fcfcfc);
  border: 1px solid var(--ctl-glass-border, rgb(0 0 0 / 5%));
  box-shadow: var(--ctl-shadow-sm);
  overflow: hidden;
}

.studio-col--side {
  background: #f5f5f7;
  border-color: rgb(0 0 0 / 4%);
  box-shadow: none;
}

.studio-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  flex-shrink: 0;
  padding: 11px 14px;
  border-bottom: 1px solid rgb(0 0 0 / 5%);
  background: rgb(252 252 252 / 92%);
  backdrop-filter: saturate(1.2) blur(14px);
  -webkit-backdrop-filter: saturate(1.2) blur(14px);
}

.studio-bar__left {
  display: flex;
  align-items: center;
  gap: 10px;
  min-width: 0;
  flex: 1;
}

.studio-eyebrow {
  flex-shrink: 0;
  font-size: 13px;
  font-weight: 600;
  letter-spacing: -0.02em;
  color: var(--ctl-text, rgba(0, 0, 0, 0.85));
}

.studio-name {
  max-width: 240px;
}

.studio-show {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  flex-shrink: 0;
  padding: 4px 10px;
  border-radius: 999px;
  background: var(--ctl-fill, rgb(0 0 0 / 6%));
  font-size: 12px;
  font-weight: 500;
  color: var(--ctl-text-soft, rgba(0, 0, 0, 0.55));
}

.seg {
  display: inline-flex;
  padding: 3px;
  border-radius: 10px;
  background: var(--ctl-fill, rgb(0 0 0 / 6%));
  gap: 2px;
}

.seg__btn {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  padding: 6px 12px;
  border: 0;
  border-radius: 8px;
  background: transparent;
  color: var(--ctl-text, rgba(0, 0, 0, 0.85));
  font-size: 12px;
  font-weight: 560;
  cursor: grab;
  transition:
    background 0.18s ease,
    transform 0.18s var(--ctl-ease, cubic-bezier(0.32, 0.72, 0, 1));
}

.seg__btn:hover {
  background: #fff;
  box-shadow: 0 1px 2px rgb(0 0 0 / 6%);
}

.seg__btn:active {
  cursor: grabbing;
  transform: scale(var(--ctl-press, 0.96));
}

.seg__btn i {
  font-size: 14px;
  opacity: 0.7;
}

.drop-hero {
  flex: 1;
  margin: 16px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 6px;
  border-radius: 16px;
  border: 1.5px dashed rgb(0 0 0 / 12%);
  background: rgb(0 0 0 / 2%);
  cursor: pointer;
  transition:
    border-color 0.22s ease,
    background 0.22s ease,
    box-shadow 0.28s var(--ctl-ease-soft, cubic-bezier(0.22, 1, 0.36, 1)),
    transform 0.22s var(--ctl-ease, cubic-bezier(0.32, 0.72, 0, 1));
}

.drop-hero:hover,
.drop-hero--active {
  border-color: rgba(var(--primary-color, 49, 108, 114), 0.4);
  border-style: solid;
  background: #fff;
  box-shadow: 0 0 0 4px var(--ctl-ring, rgba(var(--primary-color, 49, 108, 114), 0.18));
  transform: translateY(-1px);
}

.drop-hero__glyph {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 52px;
  height: 52px;
  margin-bottom: 8px;
  border-radius: 14px;
  font-size: 22px;
  color: rgb(var(--primary-color, 49, 108, 114));
  background: rgba(var(--primary-color, 49, 108, 114), 0.1);
}

.drop-hero__title {
  margin: 0;
  font-size: 15px;
  font-weight: 600;
  letter-spacing: -0.02em;
}

.drop-hero__desc {
  margin: 0;
  font-size: 12px;
  color: var(--ctl-text-soft, rgba(0, 0, 0, 0.45));
}

.board-wrap {
  flex: 1;
  min-height: 0;
  overflow: auto;
  padding: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  background:
    radial-gradient(circle at 1px 1px, rgb(0 0 0 / 4.5%) 1px, transparent 0) 0 0 / 18px 18px;
}

.board-frame {
  position: relative;
  width: 800px;
  height: 500px;
  flex-shrink: 0;
}

.board {
  position: relative;
  width: 100%;
  height: 100%;
  border-radius: 12px;
  overflow: hidden;
  background: #f2f2f7;
}

.board--edit {
  border: 1px solid rgb(0 0 0 / 8%);
  box-shadow:
    0 1px 2px rgb(0 0 0 / 4%),
    0 12px 32px rgb(0 0 0 / 6%);
}

.board__bg {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
  object-fit: cover;
  pointer-events: none;
  border-radius: 12px;
}

.board-shell {
  z-index: 2;
}

.node {
  position: relative;
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 0;
  border-radius: 8px;
  border: 1px dashed rgb(0 0 0 / 18%);
  background: rgb(255 255 255 / 72%);
  backdrop-filter: blur(8px) saturate(1.1);
  -webkit-backdrop-filter: blur(8px) saturate(1.1);
  overflow: hidden;
  line-height: 1.35;
  word-break: break-word;
  font-family: 'Noto Sans SC', 'PingFang SC', 'Microsoft YaHei', sans-serif;
  font-weight: 500;
  transition:
    border-color 0.2s ease,
    box-shadow 0.2s ease,
    background 0.2s ease;
}

.node--on {
  border-color: rgba(var(--primary-color, 49, 108, 114), 0.55);
  border-style: solid;
  box-shadow: 0 0 0 3px var(--ctl-ring, rgba(var(--primary-color, 49, 108, 114), 0.2));
}

.node--image {
  background: rgb(255 255 255 / 25%);
  padding: 4px;
}

.node__img,
.live-layer__img {
  max-width: 100%;
  max-height: 100%;
  object-fit: contain;
  pointer-events: none;
}

.node__empty {
  font-size: 12px;
  color: var(--ctl-text-soft, rgba(0, 0, 0, 0.4));
}

.node__html {
  width: 100%;
  max-height: 100%;
  overflow: hidden;
  white-space: pre-wrap;
  word-break: break-word;
}

.node__x {
  position: absolute;
  top: -8px;
  right: -8px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 18px;
  height: 18px;
  padding: 0;
  border: 0;
  border-radius: 999px;
  background: #ff3b30;
  color: #fff;
  font-size: 10px;
  cursor: pointer;
  z-index: 3;
  box-shadow: 0 0 0 2px #fff;
}

.side-preview {
  display: flex;
  flex-direction: column;
  flex: 1.05 1 0;
  min-height: 240px;
}

.side-editor {
  flex: 1.15 1 0;
  min-height: 280px;
  overflow: auto;
  padding: 0 16px 16px;
  border-top: 1px solid rgb(0 0 0 / 5%);
  background: #f5f5f7;
}

.side-editor .side-section-head {
  padding: 14px 2px 12px;
  position: sticky;
  top: 0;
  z-index: 1;
  background: linear-gradient(180deg, #f5f5f7 70%, transparent);
}

.side-editor__empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 6px;
  min-height: 160px;
  padding: 28px;
  border-radius: 14px;
  background: rgb(255 255 255 / 78%);
  border: 1px solid rgb(0 0 0 / 4%);
  text-align: center;
}

.side-editor__empty p {
  margin: 0;
  font-size: 14px;
  font-weight: 600;
  letter-spacing: -0.02em;
  color: rgba(0, 0, 0, 0.55);
}

.side-editor__empty span {
  font-size: 12px;
  color: rgba(0, 0, 0, 0.32);
}

.editor-panel {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding-bottom: 4px;
}

.edit-toolbar {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 8px;
  padding: 8px 10px;
  border-radius: 12px;
  background: #fff;
  border: 1px solid rgb(0 0 0 / 5%);
  box-shadow: 0 1px 2px rgb(0 0 0 / 3%);
}

.edit-toolbar__divider {
  width: 1px;
  height: 22px;
  background: rgb(0 0 0 / 8%);
  flex-shrink: 0;
}

.edit-toolbar__color {
  width: 28px;
  flex-shrink: 0;
}

.edit-toolbar__color :deep(.n-color-picker-trigger) {
  width: 28px;
  height: 28px;
  min-height: 28px;
  padding: 3px;
  border-radius: 8px;
}

.edit-toolbar__color :deep(.n-color-picker-trigger__fill) {
  border-radius: 5px;
  overflow: hidden;
}

.edit-toolbar__color :deep(.n-color-picker-trigger__value) {
  display: none;
}

.edit-toolbar__size {
  display: inline-flex;
  align-items: center;
  gap: 2px;
  padding: 2px;
  border-radius: 9px;
  background: var(--ctl-fill, rgb(0 0 0 / 6%));
}

.edit-toolbar__size-input {
  width: 44px;
}

.edit-toolbar__size-input :deep(.n-input) {
  --n-height: 28px !important;
  text-align: center;
}

.edit-toolbar__size-input :deep(.n-input__input-el) {
  text-align: center;
  font-variant-numeric: tabular-nums;
  font-weight: 600;
  font-size: 12px;
}

.tool-ico__a {
  font-size: 13px;
  font-weight: 700;
  line-height: 1;
}

.tool-ico__a--sm {
  font-size: 11px;
}

.editor-preview {
  min-height: 44px;
  padding: 10px 12px;
  border-radius: 12px;
  background: #fff;
  border: 1px solid rgb(0 0 0 / 5%);
  font-size: 14px;
  line-height: 1.55;
  word-break: break-word;
  white-space: pre-wrap;
  color: var(--ctl-text, rgba(0, 0, 0, 0.78));
  box-shadow: 0 1px 2px rgb(0 0 0 / 2%);
}

.inspector__textarea {
  margin-top: 0;
}

.inspector__textarea :deep(.n-input) {
  border-radius: 12px;
  background: #fff;
}

.side-section-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  flex-shrink: 0;
  padding: 14px 16px 0;
}

.side-section-head__title {
  font-size: 14px;
  font-weight: 600;
  letter-spacing: -0.02em;
  color: var(--ctl-text, rgba(0, 0, 0, 0.85));
}

.side-section-head__meta {
  font-size: 12px;
  font-weight: 500;
  color: rgba(0, 0, 0, 0.35);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.tool-cluster {
  display: inline-flex;
  gap: 2px;
  padding: 2px;
  border-radius: 9px;
  background: var(--ctl-fill, rgb(0 0 0 / 6%));
}

.tool-ico {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border: 0;
  border-radius: 7px;
  background: transparent;
  color: var(--ctl-text, rgba(0, 0, 0, 0.75));
  cursor: pointer;
  font-size: 16px;
  transition: background 0.15s ease, transform 0.15s ease;
}

.tool-ico:hover,
.tool-ico--on {
  background: #fff;
  box-shadow: 0 1px 2px rgb(0 0 0 / 6%);
}

.tool-ico:active {
  transform: scale(0.96);
}

.tool-text {
  border: 0;
  border-radius: 8px;
  padding: 6px 12px;
  background: var(--ctl-fill, rgb(0 0 0 / 6%));
  font-size: 12px;
  font-weight: 560;
  color: var(--ctl-text, rgba(0, 0, 0, 0.8));
  cursor: pointer;
  white-space: nowrap;
}

.tool-text:hover {
  background: var(--ctl-fill-hover, rgb(0 0 0 / 10%));
}

.var-row {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.var-chip {
  border: 0;
  border-radius: 7px;
  padding: 6px 12px;
  font-size: 12px;
  font-weight: 600;
  cursor: pointer;
  transition: transform 0.15s ease, filter 0.15s ease;
}

.var-chip:active {
  transform: scale(0.96);
}

.var-chip[data-key='name'] {
  color: #0b6e4f;
  background: rgb(52 199 89 / 18%);
}

.var-chip[data-key='course'] {
  color: #0a5c8a;
  background: rgb(10 132 255 / 16%);
}

.var-chip[data-key='time'] {
  color: #8a5a00;
  background: rgb(255 159 10 / 18%);
}

.var-chip[data-key='no'] {
  color: #5b3cc4;
  background: rgb(94 92 230 / 16%);
}

:deep(.var-token) {
  display: inline;
  padding: 1px 6px;
  margin: 0 1px;
  border-radius: 5px;
  font-size: 0.92em;
  font-weight: 600;
  white-space: nowrap;
}

:deep(.var-token[data-key='name']) {
  color: #0b6e4f;
  background: rgb(52 199 89 / 20%);
}

:deep(.var-token[data-key='course']) {
  color: #0a5c8a;
  background: rgb(10 132 255 / 18%);
}

:deep(.var-token[data-key='time']) {
  color: #8a5a00;
  background: rgb(255 159 10 / 20%);
}

:deep(.var-token[data-key='no']) {
  color: #5b3cc4;
  background: rgb(94 92 230 / 18%);
}

:deep(.var-placeholder) {
  color: rgba(0, 0, 0, 0.32);
  font-weight: 400;
}

.preview-stage {
  flex: 1;
  min-height: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 14px 18px 18px;
  overflow: auto;
}

.preview-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  color: rgba(0, 0, 0, 0.35);
  text-align: center;
}

.preview-empty__title {
  margin: 4px 0 0;
  font-size: 13px;
  font-weight: 600;
  letter-spacing: -0.02em;
  color: rgba(0, 0, 0, 0.55);
}

.preview-empty__desc {
  margin: 0;
  font-size: 12px;
  color: rgba(0, 0, 0, 0.32);
}

.preview-empty__ring {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 48px;
  height: 48px;
  border-radius: 14px;
  background: rgb(255 255 255 / 80%);
  border: 1px solid rgb(0 0 0 / 4%);
  color: rgba(0, 0, 0, 0.28);
  font-size: 18px;
  box-shadow: 0 1px 2px rgb(0 0 0 / 4%);
}

.preview-hero {
  width: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.preview-card {
  width: 100%;
  overflow: hidden;
  border-radius: 14px;
  background: #fff;
  box-shadow:
    0 0 0 0.5px rgb(0 0 0 / 6%),
    0 2px 8px rgb(0 0 0 / 4%),
    0 20px 44px rgb(0 0 0 / 8%);
}

.preview-viewport {
  position: relative;
  width: 100%;
  aspect-ratio: 8 / 5;
  overflow: hidden;
  background: #f2f2f7;
  container-type: size;
}

.preview-viewport .board--live {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
  overflow: hidden;
  border: 0;
  border-radius: 0;
  box-shadow: none;
  /* 1em = 1 设计像素（画布高 500） */
  font-size: calc(100cqh / 500);
}

.live-layer {
  position: absolute;
  z-index: 2;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  pointer-events: none;
}

.live-layer__text {
  width: 100%;
  line-height: 1.35;
  word-break: break-word;
  font-family: 'Noto Sans SC', 'PingFang SC', 'Microsoft YaHei', sans-serif;
  font-weight: 500;
  letter-spacing: 0;
  white-space: pre-wrap;
}

@keyframes soft-in {
  from {
    opacity: 0;
    transform: translate3d(0, 8px, 0);
  }
  to {
    opacity: 1;
    transform: none;
  }
}

@keyframes pulse-dot {
  from {
    opacity: 0.35;
    transform: scale(0.85);
  }
  to {
    opacity: 1;
    transform: scale(1);
  }
}

@media (max-width: 1280px) {
  .studio-grid {
    grid-template-columns: minmax(0, 1fr) minmax(400px, 460px);
  }
}

@media (max-width: 1100px) {
  .studio-grid {
    grid-template-columns: 1fr;
    grid-template-rows: minmax(320px, 1.05fr) minmax(380px, 1fr);
  }

  .side-preview,
  .side-editor {
    flex: 1 1 auto;
    min-height: 0;
  }
}

@media (prefers-reduced-motion: reduce) {
  .cert-studio,
  .studio-loading__dot {
    animation: none;
  }

  .drop-hero:hover,
  .seg__btn:active,
  .tool-ico:active,
  .var-chip:active {
    transform: none;
  }
}
</style>
