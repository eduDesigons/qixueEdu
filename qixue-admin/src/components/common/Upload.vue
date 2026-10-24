<script setup>
import { Add, Close, CloudUploadOutline } from '@vicons/ionicons5'
import { NIcon } from 'naive-ui'
import { computed, onMounted, ref, useSlots, watch } from 'vue'
import { MediaPlayer } from '@/components/index.js'
import MediaPickerModal from './MediaPickerModal.vue'
import api from '@/api'

const props = defineProps({
  type: { type: String, default: 'image' },
  src: { type: String, default: '' },
  showPreview: { type: Boolean, default: true },
  thumb: { type: String, default: '' },
  cateId: { type: Number, default: 0 },
  showDragBox: { type: Boolean, default: false },
  multiple: { type: Boolean, default: false },
  allowLibrary: { type: Boolean, default: true },
  /** 预览区宽高，数字按 px；不传则保持原有自适应 */
  previewWidth: { type: [Number, String], default: undefined },
  previewHeight: { type: [Number, String], default: undefined },
  /** 预览图 object-fit，头像场景可用 cover */
  previewFit: { type: String, default: 'contain' },
  /** 封面最长边。传入后上传时压到这个尺寸，0 表示不额外缩小 */
  maxEdge: { type: Number, default: 0 },
})
const emit = defineEmits(['success'])
const slots = useSlots()
const libraryShow = ref(false)
const showLibraryEntry = computed(() => props.allowLibrary && (props.showDragBox || !!slots.default))

function toCssSize(val) {
  if (val == null || val === '')
    return undefined
  return typeof val === 'number' ? `${val}px` : String(val)
}

const hasPreviewSize = computed(() => props.previewWidth != null || props.previewHeight != null)
const previewBoxStyle = computed(() => {
  const style = {}
  const w = toCssSize(props.previewWidth)
  const h = toCssSize(props.previewHeight)
  if (w) {
    style.width = w
    style.maxWidth = w
  }
  if (h) {
    style.height = h
    style.minHeight = h
  }
  return style
})

const info = ref({
  url: '',
  title: '',
  thumb: null,
  type: 'image',
})
watch(() => props.src, (val) => {
  if (val) {
    info.value = {
      ...info.value,
      url: props.src,
      type: props.type,
      thumb: props.thumb,
    }
  }
})

watch(() => props.type, (newType) => {
  // 更新 info 中的 type
  info.value = {
    ...info.value,
    type: newType,
  }
})

const accept = computed(() => {
  switch (props.type) {
    case 'image':
      return 'image/png, image/jpeg, image/gif, image/jpg, image/webp'
    case 'video':
      return 'video/mp4, video/mpeg, video/quicktime'
    case 'audio':
      return 'audio/mpeg, audio/wav, audio/mp3'
    case 'document':
      return 'application/pdf, application/msword, application/vnd.openxmlformats-officedocument.wordprocessingml.document, application/vnd.ms-excel, application/vnd.openxmlformats-officedocument.spreadsheetml.sheet, application/vnd.ms-powerpoint, application/vnd.openxmlformats-officedocument.presentationml.presentation'
    case 'app':
      return '.apk,.ipa,application/vnd.android.package-archive,application/octet-stream'
    default:
      return 'image/png, image/jpeg, image/gif, image/jpg'
  }
})

const uploadProgress = ref(0)
const uploadLoaded = ref(0)
const uploadTotal = ref(0)
const uploadingCount = ref(0)
const uploadPhase = ref('idle') // prepare | upload | finalize | done
const uploadFileName = ref('')
const uploadThumbUrl = ref('')
const uploadIndeterminate = ref(false)
const showUploadCard = ref(false)
const isUploading = computed(() => uploadingCount.value > 0 || showUploadCard.value)
const UPLOAD_MSG_KEY = 'upload'
const RING_R = 34
const RING_C = 2 * Math.PI * RING_R

function formatUploadSize(bytes) {
  const n = Number(bytes) || 0
  if (n < 1024)
    return `${n} B`
  if (n < 1024 * 1024)
    return `${(n / 1024).toFixed(1)} KB`
  if (n < 1024 * 1024 * 1024)
    return `${(n / 1024 / 1024).toFixed(1)} MB`
  return `${(n / 1024 / 1024 / 1024).toFixed(2)} GB`
}

const progressPercent = computed(() =>
  Math.min(100, Math.max(0, Number(uploadProgress.value) || 0)),
)

const ringDashoffset = computed(() => {
  const p = uploadIndeterminate.value ? 0.22 : progressPercent.value / 100
  return RING_C * (1 - p)
})

const phaseLabel = computed(() => {
  switch (uploadPhase.value) {
    case 'prepare':
      return '正在准备'
    case 'upload':
      return uploadIndeterminate.value ? '正在上传' : '正在上传到云端'
    case 'finalize':
      return '即将完成'
    case 'done':
      return '上传完成'
    default:
      return '上传中'
  }
})

const progressSizeText = computed(() => {
  if (uploadTotal.value > 0)
    return `${formatUploadSize(uploadLoaded.value)} / ${formatUploadSize(uploadTotal.value)}`
  if (uploadTotal.value === 0 && uploadFileName.value)
    return uploadFileName.value
  return ''
})

const progressHint = computed(() => {
  if (uploadPhase.value === 'prepare')
    return '解析视频封面与时长'
  if (uploadPhase.value === 'finalize')
    return '写入资源库'
  if (uploadPhase.value === 'done')
    return '可以继续编辑了'
  if (uploadTotal.value > 0 && progressPercent.value > 0 && progressPercent.value < 100) {
    const left = Math.max(0, uploadTotal.value - uploadLoaded.value)
    return left > 0 ? `剩余 ${formatUploadSize(left)}` : '收尾中…'
  }
  return uploadFileName.value || '请稍候'
})

function clearUploadThumb() {
  if (uploadThumbUrl.value) {
    URL.revokeObjectURL(uploadThumbUrl.value)
    uploadThumbUrl.value = ''
  }
}

function resetUploadProgress() {
  uploadProgress.value = 0
  uploadLoaded.value = 0
  uploadTotal.value = 0
  uploadPhase.value = 'idle'
  uploadFileName.value = ''
  uploadIndeterminate.value = false
  showUploadCard.value = false
  clearUploadThumb()
}

function setUploadProgress(percent, loaded, total) {
  uploadIndeterminate.value = false
  uploadProgress.value = Math.min(100, Math.max(0, Math.round(percent)))
  if (typeof loaded === 'number')
    uploadLoaded.value = loaded
  if (typeof total === 'number' && total > 0)
    uploadTotal.value = total
}
function fileTypeCategory(file) {
  const mime = file.type || ''
  const ext = (file.name.split('.').pop() || '').toLowerCase()
  if (mime.startsWith('image/'))
    return 'image'
  if (mime.startsWith('video/'))
    return 'video'
  if (mime.startsWith('audio/'))
    return 'audio'
  if (ext === 'apk' || ext === 'ipa')
    return 'app'
  return 'document'
}

function getVideoMeta(file) {
  return new Promise((resolve) => {
    const video = document.createElement('video')
    video.preload = 'metadata'
    const objectUrl = URL.createObjectURL(file)
    let resolved = false

    const done = (duration = 0, coverBlob = null) => {
      if (resolved)
        return
      resolved = true
      URL.revokeObjectURL(objectUrl)
      resolve({ duration, coverBlob })
    }

    video.onloadedmetadata = () => {
      const duration = Number.isFinite(video.duration) ? Math.floor(video.duration) : 0
      video.currentTime = 1
      video.onseeked = () => {
        try {
          const canvas = document.createElement('canvas')
          canvas.width = video.videoWidth || 320
          canvas.height = video.videoHeight || 180
          const ctx = canvas.getContext('2d')
          if (!ctx) {
            done(duration, null)
            return
          }
          ctx.drawImage(video, 0, 0, canvas.width, canvas.height)
          canvas.toBlob((blob) => {
            done(duration, blob || null)
          }, 'image/jpeg', 0.8)
        }
        catch {
          done(duration, null)
        }
      }
    }

    video.onerror = () => done(0, null)
    video.src = objectUrl
    setTimeout(() => done(Number.isFinite(video.duration) ? Math.floor(video.duration) : 0, null), 8000)
  })
}

async function saveUploadedToDb(data) {
  const payload = {
    ...data,
    cate_id: props.cateId || 0,
  }
  const res = await api.saveUploaded(payload)
  return res.data
}

async function uploadFile(file) {
  try {
    let result
    uploadFileName.value = file.name || ''
    if (fileTypeCategory(file) === 'video') {
      const fileSize = Number(file.size) || 0
      uploadPhase.value = 'prepare'
      setUploadProgress(2, 0, fileSize)
      const { duration, coverBlob } = await getVideoMeta(file)
      if (coverBlob) {
        clearUploadThumb()
        uploadThumbUrl.value = URL.createObjectURL(coverBlob)
      }
      const thumbFile = coverBlob
        ? new File([coverBlob], 'thumb.jpg', { type: 'image/jpeg' })
        : null
      setUploadProgress(6, 0, fileSize)
      const initRes = await api.initVideoDirectUpload({
        fileName: file.name || 'video.mp4',
        fileSize,
        mimeType: file.type || 'video/mp4',
        duration: duration || 0,
        thumbFile,
      })
      const token = initRes.data
      uploadPhase.value = 'upload'
      setUploadProgress(10, 0, fileSize)
      await api.uploadVideoToS3(token, file, (percent, loaded, total) => {
        // 准备 0–10%，S3 直传 10–98%，落库 98–100%
        const mapped = 10 + Math.round((Math.min(100, percent) / 100) * 88)
        setUploadProgress(mapped, loaded, total || fileSize)
      })
      uploadPhase.value = 'finalize'
      setUploadProgress(98, fileSize, fileSize)
      const saved = await saveUploadedToDb({
        url: token.objectKey,
        title: token.title,
        type: 'video',
        vid: token.videoId,
        size: file.size,
        duration: token.duration ?? duration ?? 0,
        thumb: token.thumbUrl || '',
      })
      result = {
        ...saved,
        thumb: saved.thumb || null,
        type: saved.type || props.type || info.value.type,
      }
      uploadPhase.value = 'done'
      setUploadProgress(100, fileSize, fileSize)
    }
    else {
      uploadPhase.value = 'upload'
      uploadIndeterminate.value = true
      uploadProgress.value = 0
      uploadTotal.value = Number(file.size) || 0
      uploadLoaded.value = 0
      const uploadRes = await api.upload(file, props.cateId, props.maxEdge)
      const saved = uploadRes.data
      result = {
        ...saved,
        thumb: saved.thumb || null,
        type: saved.type || props.type || info.value.type,
      }
      uploadPhase.value = 'done'
      uploadIndeterminate.value = false
      setUploadProgress(100, uploadTotal.value, uploadTotal.value)
    }

    if (!props.multiple) {
      info.value = result
    }

    emit('success', result)
    if (!props.multiple)
      $message.success('上传成功')

    return result
  }
  catch (error) {
    resetUploadProgress()
    console.error('上传失败:', error)
    $message.error(`上传失败: ${error.message || '未知错误'}`)
    throw error
  }
}

onMounted(() => {
  info.value = {
    ...info.value,
    url: props.src,
    type: props.type,
  }
})

async function handleBeforeUpload({ file, onFinish }) {
  if (!file || !file.type) {
    $message.error('请选择文件')
    onFinish()
    return
  }

  const isFirst = uploadingCount.value === 0
  uploadingCount.value++
  if (isFirst) {
    resetUploadProgress()
    showUploadCard.value = true
  }

  try {
    await uploadFile(file.file)
  }
  catch {
    // 错误已在 uploadFile 中处理
  }
  finally {
    uploadingCount.value--
    onFinish()
    if (uploadingCount.value === 0) {
      if (props.multiple)
        $message.success('上传完成', { key: UPLOAD_MSG_KEY })
      // 完成态稍留片刻再收起卡片
      if (uploadPhase.value === 'done') {
        setTimeout(() => {
          if (uploadingCount.value === 0)
            resetUploadProgress()
        }, 780)
      }
      else {
        resetUploadProgress()
      }
    }
  }
}

function applyPicked(item) {
  const result = {
    id: item.id,
    url: item.url,
    title: item.title || '未命名资源',
    thumb: item.thumb || null,
    type: item.type || props.type,
    size: item.size,
    duration: item.duration,
  }
  if (!props.multiple)
    info.value = result
  emit('success', result)
}

function handleLibrarySelect(payload) {
  const items = Array.isArray(payload) ? payload : [payload]
  items.forEach(applyPicked)
  $message.success(items.length > 1 ? `已选用 ${items.length} 个资源` : '已选用资源库文件')
}

function handleRemove() {
  $dialog.confirm({
    content: '确定删除吗？',
    async confirm() {
      info.value = {
        url: '',
        title: '',
        thumb: null,
        type: props.type,
      }
      emit('success', { url: '' })
    },
  })
}
</script>

<template>
  <div
    class="upload-box"
    :class="{
      'upload-box--library': showLibraryEntry && (!info.url || multiple),
      'upload-box--library-overlay': showLibraryEntry && showDragBox && (!info.url || multiple),
      'upload-box--sized': hasPreviewSize,
    }"
    :style="hasPreviewSize ? previewBoxStyle : undefined"
  >
    <div
      v-if="showPreview && info.url"
      class="upload-preview-host h-full w-full"
      :class="{
        'mr-10': !hasPreviewSize,
        'h-300 w-auto': type === 'video' && !hasPreviewSize,
        'upload-preview-host--sized': hasPreviewSize,
        'upload-preview-host--library': showLibraryEntry && !isUploading,
      }"
    >
      <div class="upload-preview relative h-full w-full">
        <MediaPlayer
          :media-info="info"
          :width="previewWidth"
          :height="previewHeight"
          :style="previewFit !== 'contain' ? { '--mp-image-fit': previewFit } : undefined"
        />
        <div class="close-icon-box" @click.stop="handleRemove">
          <NIcon v-if="info.url" class="close-icon" :component="Close" size="16" />
        </div>
      </div>
      <button
        v-if="showLibraryEntry && !isUploading"
        type="button"
        class="upload-library-btn upload-library-btn--below-preview"
        @click.stop="libraryShow = true"
      >
        更换资源
      </button>
    </div>
    <div v-else class="box">
      <slot>
        <Transition name="upload-progress">
          <div
            v-if="isUploading"
            key="upload-progress"
            class="up-card"
            :class="{
              'up-card--done': uploadPhase === 'done',
              'up-card--has-thumb': !!uploadThumbUrl,
            }"
            role="status"
            aria-live="polite"
          >
            <div class="up-card__visual">
              <div class="up-card__thumb" aria-hidden="true">
                <img v-if="uploadThumbUrl" :src="uploadThumbUrl" alt="">
                <div v-else class="up-card__thumb-fallback">
                  <NIcon :component="CloudUploadOutline" :size="24" />
                </div>
              </div>
              <svg
                class="up-card__ring"
                :class="{ 'is-spinning': uploadIndeterminate }"
                viewBox="0 0 80 80"
                aria-hidden="true"
              >
                <circle class="up-card__ring-rail" cx="40" cy="40" :r="RING_R" />
                <circle
                  class="up-card__ring-fill"
                  :class="{ 'is-indeterminate': uploadIndeterminate }"
                  cx="40"
                  cy="40"
                  :r="RING_R"
                  :stroke-dasharray="RING_C"
                  :stroke-dashoffset="ringDashoffset"
                />
              </svg>
              <div class="up-card__pct">
                <div class="up-card__pct-inner">
                  <template v-if="uploadIndeterminate">
                    <span class="up-card__pct-dots" />
                  </template>
                  <template v-else>
                    <span class="up-card__pct-num">{{ progressPercent }}</span>
                    <span class="up-card__pct-unit">%</span>
                  </template>
                </div>
              </div>
            </div>

            <div class="up-card__copy">
              <div class="up-card__title">
                {{ phaseLabel }}
              </div>
              <div class="up-card__hint">
                {{ progressHint }}
              </div>
              <div class="up-card__track" aria-hidden="true">
                <div
                  class="up-card__fill"
                  :class="{ 'is-indeterminate': uploadIndeterminate }"
                  :style="uploadIndeterminate ? undefined : { width: `${progressPercent}%` }"
                />
              </div>
              <div v-if="progressSizeText" class="up-card__meta">
                <span class="up-card__size">{{ progressSizeText }}</span>
                <span v-if="uploadPhase !== 'idle'" class="up-card__steps">
                  <i :class="{ 'is-on': ['prepare', 'upload', 'finalize', 'done'].includes(uploadPhase) }" />
                  <i :class="{ 'is-on': ['upload', 'finalize', 'done'].includes(uploadPhase) }" />
                  <i :class="{ 'is-on': ['finalize', 'done'].includes(uploadPhase) }" />
                </span>
              </div>
            </div>
          </div>
        </Transition>
        <template v-if="!isUploading">
          <template v-if="!info.url">
            <!-- 拖拽上传区域 -->
            <n-upload
              v-if="showDragBox"
              class="upload-drag-wrap"
              :class="{ 'upload-drag-wrap--sized': hasPreviewSize }"
              :style="hasPreviewSize ? previewBoxStyle : undefined"
              :custom-request="handleBeforeUpload"
              :accept="accept"
              :multiple="multiple"
              :show-file-list="false"
            >
              <n-upload-dragger>
                <div class="upload-area-content">
                  <NIcon :component="Add" size="40" />
                  <div class="upload-area-content-text">
                    点击或拖拽上传{{ type === 'image' ? '图片' : type === 'video' ? '视频' : type === 'audio' ? '音频' : type === 'app' ? '安装包(APK/IPA)' : '文件' }}
                  </div>
                </div>
              </n-upload-dragger>
            </n-upload>
            <!-- 普通按钮上传 -->
            <n-button v-else type="primary">
              上传
            </n-button>
          </template>
        </template>
      </slot>
      <!-- 普通上传覆盖层（用于点击上传，不显示拖拽区域时） -->
      <div v-if="!isUploading && !showDragBox && (!info.url || multiple)" class="upload">
        <n-upload :custom-request="handleBeforeUpload" :accept="accept" :multiple="multiple" :show-file-list="false">
          <div class="upload" />
        </n-upload>
      </div>
    </div>
    <button
      v-if="showLibraryEntry && !isUploading && (!info.url || multiple)"
      type="button"
      class="upload-library-btn"
      :class="{ 'upload-library-btn--ghost': showDragBox }"
      @click.stop="libraryShow = true"
    >
      从资源库选择
    </button>
    <MediaPickerModal
      v-model:show="libraryShow"
      :type="type"
      :multiple="multiple"
      @select="handleLibrarySelect"
    />
  </div>
</template>

<style lang="scss" scoped>
.upload-preview-host {
  position: relative;
}

.upload-preview {
  overflow: hidden;
  border-radius: 12px;
}

.upload-preview :deep(.mp-image) {
  object-fit: contain;
}

.upload-box--library {
  position: relative;
  display: inline-flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 0;
}

.upload-box--library-overlay {
  position: relative;
  display: block;
}

.upload-box--library-overlay .upload-library-btn {
  position: absolute;
  left: 50%;
  bottom: 10px;
  z-index: 20;
  transform: translateX(-50%) translateY(4px);
  opacity: 0;
  pointer-events: none;
}

.upload-box--library-overlay:hover .upload-library-btn,
.upload-box--library-overlay:focus-within .upload-library-btn,
.upload-box--library-overlay .upload-library-btn:focus-visible {
  opacity: 1;
  pointer-events: auto;
  transform: translateX(-50%) translateY(0);
}

.upload-box--library-overlay .upload-library-btn:active {
  transform: translateX(-50%) translateY(0) scale(0.97);
}

.upload-box--library .box {
  flex: 1 1 auto;
  min-height: 0;
  height: auto;
}

.upload-box--library .box :deep(.n-upload),
.upload-box--library .box :deep(.n-upload-trigger) {
  width: 100%;
  height: 100%;
}

.upload-library-btn {
  appearance: none;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  box-sizing: border-box;
  height: 28px;
  padding: 0 12px;
  border: 1px solid rgba(var(--primary-color), 0.14);
  border-radius: 10px;
  background: rgb(255 255 255 / 92%);
  color: rgb(var(--primary-color));
  font-size: 12px;
  font-weight: 550;
  letter-spacing: -0.01em;
  line-height: 1;
  cursor: pointer;
  z-index: 20;
  white-space: nowrap;
  box-shadow: 0 6px 16px rgb(15 40 45 / 10%);
  transition:
    opacity 0.18s ease,
    transform 0.2s cubic-bezier(0.32, 0.72, 0, 1),
    background 0.2s ease,
    border-color 0.2s ease,
    box-shadow 0.2s ease;
}

.upload-library-btn:hover {
  background: #fff;
  border-color: rgba(var(--primary-color), 0.28);
  box-shadow: 0 8px 18px rgb(15 40 45 / 12%);
}

/* 空态：悬停上传区后底部居中露出「从资源库选择」 */
.upload-box--library > .upload-library-btn {
  position: absolute;
  left: 50%;
  bottom: 10px;
  top: auto;
  z-index: 30;
  margin: 0;
  opacity: 0;
  pointer-events: none;
  transform: translateX(-50%) translateY(4px);
}

.upload-box--library:hover > .upload-library-btn,
.upload-box--library:focus-within > .upload-library-btn,
.upload-box--library > .upload-library-btn:focus-visible {
  opacity: 1;
  pointer-events: auto;
  transform: translateX(-50%) translateY(0);
}

.upload-box--library > .upload-library-btn:active {
  transform: translateX(-50%) translateY(0) scale(0.97);
}

.upload-library-btn--ghost {
  background: rgb(255 255 255 / 88%);
  color: rgb(var(--primary-color));
}

.upload-library-btn--ghost:hover {
  background: #fff;
  color: rgb(var(--primary-color));
}

/* 更换资源：预览底部居中浮层，样式与「从资源库选择」一致 */
.upload-library-btn--below-preview {
  position: absolute;
  left: 50%;
  bottom: 10px;
  top: auto;
  z-index: 30;
  margin: 0;
  opacity: 0;
  pointer-events: none;
  transform: translateX(-50%) translateY(4px);
}

.upload-preview-host--library:hover > .upload-library-btn--below-preview,
.upload-preview-host--library:focus-within > .upload-library-btn--below-preview,
.upload-library-btn--below-preview:focus-visible {
  opacity: 1;
  pointer-events: auto;
  transform: translateX(-50%) translateY(0);
}

.upload-library-btn--below-preview:active {
  transform: translateX(-50%) translateY(0) scale(0.97);
}

.close-icon-box {
  position: absolute;
  right: -8px;
  top: -8px;
  z-index: 10;
  background: rgb(255 255 255 / 88%);
  border-radius: 999px;
  padding: 0;
  border: 1px solid rgba(var(--primary-color), 0.22);
  color: rgb(var(--primary-color));
  cursor: pointer;
  width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: inset 0 1px 0 rgb(255 255 255 / 90%);
  backdrop-filter: none;
  -webkit-backdrop-filter: none;
  transition:
    transform 0.2s cubic-bezier(0.32, 0.72, 0, 1),
    background 0.2s ease,
    color 0.2s ease;

  &:hover {
    background: rgb(var(--primary-color));
    color: #fff;
    border-color: rgb(var(--primary-color));
  }

  &:active {
    transform: scale(0.92);
  }
}

.box {
  position: relative;
  width: 100%;
  height: 100%;

  .upload {
    position: absolute;
    left: 0;
    right: 0;
    top: 0;
    bottom: 0;
    cursor: pointer;
    z-index: 10;
  }
}

// n-upload-dragger 组件会自动处理拖拽样式，这里保留原有样式作为备用
.upload-area {
  width: 100%;
  height: 100%;
  border: 1px solid var(--ctl-hairline, rgba(var(--primary-color), 0.08));
  border-radius: 16px;
  background: var(--ctl-field, rgba(var(--primary-color), 0.04));
  box-shadow: inset 0 1px 0 rgb(255 255 255 / 40%);
  transition:
    border-color 0.22s ease,
    background 0.22s ease,
    box-shadow 0.28s cubic-bezier(0.32, 0.72, 0, 1);
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;

  &:hover {
    border-color: rgba(var(--primary-color), 0.28);
    background: var(--ctl-field-hover, rgba(var(--primary-color), 0.07));
    box-shadow:
      inset 0 1px 0 rgb(255 255 255 / 45%),
      0 0 0 4px var(--ctl-ring, rgba(var(--primary-color), 0.12));
  }
}

:deep(.n-upload-dragger) {
  display: flex !important;
  align-items: center !important;
  justify-content: center !important;
  width: 100%;
  height: 100%;
  box-sizing: border-box !important;
  border: 1px solid var(--ctl-hairline, rgba(var(--primary-color), 0.08)) !important;
  border-radius: 16px !important;
  background: var(--ctl-field, rgba(var(--primary-color), 0.04)) !important;
  box-shadow: inset 0 1px 0 rgb(255 255 255 / 40%);
  backdrop-filter: none !important;
  -webkit-backdrop-filter: none !important;
  transition:
    border-color 0.22s ease,
    background 0.22s ease,
    box-shadow 0.28s cubic-bezier(0.32, 0.72, 0, 1) !important;

  &:hover {
    border-color: rgba(var(--primary-color), 0.28) !important;
    background: var(--ctl-field-hover, rgba(var(--primary-color), 0.07)) !important;
    box-shadow:
      inset 0 1px 0 rgb(255 255 255 / 45%),
      0 0 0 4px var(--ctl-ring, rgba(var(--primary-color), 0.12)) !important;
  }
}

.upload-area-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 12px;
  width: 100%;
  height: 100%;
  min-height: 120px;
  box-sizing: border-box;
  padding: 16px;
  color: rgb(var(--primary-color));
  pointer-events: none;
  text-align: center;
}

.upload-area-content-text {
  font-size: 13px;
  font-weight: 540;
  letter-spacing: -0.01em;
  color: var(--text-secondary, #5b6b7c);
  text-align: center;
}

/* 上传进度：Apple 玻璃卡片 */
.upload-progress-enter-active,
.upload-progress-leave-active {
  transition:
    opacity 0.4s cubic-bezier(0.32, 0.72, 0, 1),
    transform 0.4s cubic-bezier(0.32, 0.72, 0, 1);
}

.upload-progress-enter-from,
.upload-progress-leave-to {
  opacity: 0;
  transform: translateY(12px) scale(0.96);
}

.upload-progress-enter-to,
.upload-progress-leave-from {
  opacity: 1;
  transform: translateY(0) scale(1);
}

.up-card {
  --up-ease: cubic-bezier(0.32, 0.72, 0, 1);
  display: flex;
  align-items: center;
  gap: 18px;
  width: 100%;
  min-height: 132px;
  padding: 18px 20px;
  box-sizing: border-box;
  border-radius: 18px;
  border: 1px solid rgb(255 255 255 / 72%);
  background:
    linear-gradient(165deg, rgb(255 255 255 / 78%) 0%, rgb(248 251 251 / 88%) 100%);
  box-shadow:
    inset 0 1px 0 rgb(255 255 255 / 90%),
    0 1px 2px rgb(15 40 45 / 4%),
    0 18px 40px rgb(15 40 45 / 8%);
  backdrop-filter: blur(22px) saturate(1.35);
  -webkit-backdrop-filter: blur(22px) saturate(1.35);
  overflow: hidden;
}

.up-card::before {
  content: '';
  position: absolute;
  inset: 0 auto auto 0;
  width: 140px;
  height: 140px;
  margin: -48px 0 0 -36px;
  border-radius: 50%;
  background: radial-gradient(circle, rgba(var(--primary-color), 0.18), transparent 68%);
  pointer-events: none;
}

.up-card {
  position: relative;
}

.up-card--done {
  border-color: rgba(var(--primary-color), 0.22);
  box-shadow:
    inset 0 1px 0 rgb(255 255 255 / 92%),
    0 0 0 4px rgba(var(--primary-color), 0.08),
    0 18px 40px rgb(15 40 45 / 8%);
}

.up-card__visual {
  position: relative;
  flex-shrink: 0;
  width: 80px;
  height: 80px;
}

.up-card__thumb {
  position: absolute;
  inset: 10px;
  border-radius: 50%;
  overflow: hidden;
  background:
    radial-gradient(120% 80% at 50% 0%, rgba(var(--primary-color), 0.16), transparent 70%),
    rgba(var(--primary-color), 0.06);
  box-shadow: inset 0 0 0 1px rgb(255 255 255 / 45%);
}

.up-card__thumb img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
}

.up-card__thumb-fallback {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100%;
  color: rgb(var(--primary-color));
  font-size: 22px;
  opacity: 0.85;
}

.up-card__ring {
  position: absolute;
  inset: 0;
  width: 80px;
  height: 80px;
  transform: rotate(-90deg);
  transform-origin: 50% 50%;
}

.up-card__ring.is-spinning {
  animation: up-ring-spin 1.15s linear infinite;
}

.up-card__ring-rail,
.up-card__ring-fill {
  fill: none;
  stroke-width: 4.5;
  stroke-linecap: round;
}

.up-card__ring-rail {
  stroke: rgba(var(--primary-color), 0.12);
}

.up-card__ring-fill {
  stroke: rgb(var(--primary-color));
  transition: stroke-dashoffset 0.28s var(--up-ease);
  filter: drop-shadow(0 0 6px rgba(var(--primary-color), 0.28));
}

.up-card__ring-fill.is-indeterminate {
  stroke-dashoffset: 160 !important;
  transition: none;
}

@keyframes up-ring-spin {
  to {
    transform: rotate(270deg);
  }
}

.up-card__pct {
  position: absolute;
  inset: 0;
  display: grid;
  place-items: center;
  z-index: 2;
  pointer-events: none;
  font-variant-numeric: tabular-nums;
  color: var(--ctl-text, #1a2332);
}

.up-card__pct-inner {
  display: inline-flex;
  align-items: baseline;
  justify-content: center;
  gap: 1px;
  line-height: 1;
  /* 数字+百分号作为整体，光学垂直居中微调 */
  transform: translateY(0.5px);
}

.up-card--has-thumb .up-card__pct {
  color: #fff;
  text-shadow: 0 1px 8px rgb(0 0 0 / 35%);
}

.up-card--has-thumb .up-card__thumb::after {
  content: '';
  position: absolute;
  inset: 0;
  background: rgb(12 24 28 / 28%);
  z-index: 1;
}

.up-card__pct-num {
  font-size: 18px;
  font-weight: 650;
  letter-spacing: -0.03em;
  line-height: 1;
}

.up-card__pct-unit {
  font-size: 10px;
  font-weight: 600;
  opacity: 0.72;
  letter-spacing: -0.01em;
  line-height: 1;
}

.up-card__pct-dots {
  width: 18px;
  height: 4px;
  border-radius: 999px;
  background:
    linear-gradient(90deg,
      transparent 0%,
      currentColor 30%,
      currentColor 70%,
      transparent 100%);
  background-size: 200% 100%;
  opacity: 0.55;
  animation: up-dots-slide 1.2s var(--up-ease) infinite;
}

@keyframes up-dots-slide {
  0% { background-position: 100% 0; }
  100% { background-position: -100% 0; }
}

.up-card__copy {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.up-card__title {
  font-size: 15px;
  font-weight: 620;
  letter-spacing: -0.022em;
  color: var(--ctl-text, #1a2332);
  line-height: 1.25;
}

.up-card__hint {
  font-size: 12px;
  font-weight: 500;
  letter-spacing: -0.01em;
  color: var(--ctl-text-soft, #5b6b7c);
  line-height: 1.35;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.up-card__track {
  position: relative;
  width: 100%;
  height: 6px;
  margin-top: 6px;
  border-radius: 999px;
  background: rgba(var(--primary-color), 0.1);
  overflow: hidden;
  box-shadow: inset 0 1px 1px rgb(15 40 45 / 6%);
}

.up-card__fill {
  height: 100%;
  width: 0;
  border-radius: inherit;
  background:
    linear-gradient(90deg,
      rgba(var(--primary-color), 0.75) 0%,
      rgb(var(--primary-color)) 55%,
      #5a9ea4 100%);
  box-shadow: 0 0 12px rgba(var(--primary-color), 0.28);
  transition: width 0.28s var(--up-ease);
  position: relative;
}

.up-card__fill::after {
  content: '';
  position: absolute;
  inset: 0;
  background: linear-gradient(
    105deg,
    transparent 35%,
    rgb(255 255 255 / 45%) 50%,
    transparent 65%
  );
  background-size: 220% 100%;
  animation: up-fill-shine 1.8s var(--up-ease) infinite;
}

.up-card__fill.is-indeterminate {
  width: 38% !important;
  animation: up-fill-indeterminate 1.35s var(--up-ease) infinite;
}

.up-card__fill.is-indeterminate::after {
  animation: none;
}

@keyframes up-fill-shine {
  0% { background-position: 120% 0; }
  100% { background-position: -120% 0; }
}

@keyframes up-fill-indeterminate {
  0% { transform: translateX(-120%); }
  100% { transform: translateX(320%); }
}

.up-card__meta {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  margin-top: 2px;
}

.up-card__size {
  font-size: 11px;
  font-weight: 550;
  letter-spacing: -0.012em;
  font-variant-numeric: tabular-nums;
  color: var(--ctl-text-soft, #5b6b7c);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  min-width: 0;
}

.up-card__steps {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  flex-shrink: 0;
}

.up-card__steps i {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: rgba(var(--primary-color), 0.18);
  transition:
    background 0.28s var(--up-ease),
    transform 0.28s var(--up-ease),
    box-shadow 0.28s var(--up-ease);
}

.up-card__steps i.is-on {
  background: rgb(var(--primary-color));
  transform: scale(1.15);
  box-shadow: 0 0 0 3px rgba(var(--primary-color), 0.14);
}

@media (prefers-reduced-motion: reduce) {
  .up-card__ring-fill,
  .up-card__fill,
  .up-card__fill::after,
  .up-card__pct-dots,
  .up-card__ring.is-spinning,
  .up-card__fill.is-indeterminate {
    animation: none !important;
    transition: none !important;
  }
}

.upload-box--sized {
  display: inline-flex;
  flex-direction: column;
  align-items: stretch;
  max-width: 100%;
}

.upload-box--sized .box {
  width: 100%;
  height: 100%;
  min-height: 0;
}

.upload-box--sized .up-card {
  flex-direction: column;
  align-items: center;
  justify-content: center;
  text-align: center;
  min-height: 100%;
  height: 100%;
  padding: 12px 14px;
  gap: 10px;
}

.upload-box--sized .up-card__visual {
  width: 72px;
  height: 72px;
}

.upload-box--sized .up-card__ring {
  width: 72px;
  height: 72px;
}

.upload-box--sized .up-card__copy {
  width: 100%;
  align-items: center;
}

.upload-box--sized .up-card__hint,
.upload-box--sized .up-card__size {
  max-width: 100%;
}

.upload-box--sized .up-card__meta {
  width: 100%;
  justify-content: center;
}

.upload-box--sized .up-card__steps {
  display: none;
}

.upload-box--sized .up-card__pct-num {
  font-size: 15px;
}

.upload-box--sized .up-card__title {
  font-size: 13px;
}

.upload-box--sized .up-card__hint {
  font-size: 11px;
}

.upload-preview-host--sized {
  margin-right: 0 !important;
}

.upload-preview-host--sized .upload-preview,
.upload-drag-wrap--sized,
.upload-drag-wrap--sized :deep(.n-upload-trigger),
.upload-drag-wrap--sized :deep(.n-upload-dragger) {
  width: 100%;
  height: 100%;
}

.upload-drag-wrap--sized :deep(.n-upload-trigger) {
  display: flex !important;
}

.upload-preview-host--sized :deep(.mp),
.upload-preview-host--sized :deep(.mp-thumb),
.upload-preview-host--sized :deep(.mp-image) {
  width: 100%;
  height: 100%;
  min-height: 100%;
  object-fit: var(--mp-image-fit, contain);
}

.upload-drag-wrap--sized :deep(.upload-area-content) {
  min-height: 0;
  padding: 12px;
}
</style>
