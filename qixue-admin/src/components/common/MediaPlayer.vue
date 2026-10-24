<template>
  <div
    class="mp"
    :class="{
      'is-empty': !canPreview,
      'is-previewable': canPreview && interactive,
      'is-compact': compact,
    }"
    :style="thumbStyle"
    @click="openPreview"
  >
    <div class="mp-thumb" :class="`mp-thumb--${kind}`">
      <template v-if="kind === 'video'">
        <img
          v-if="videoPoster && !thumbBroken"
          class="mp-image mp-image--poster"
          :src="videoPoster"
          alt=""
          @error="thumbBroken = true"
        >
        <span class="mp-play" aria-hidden="true">
          <i class="i-fe:play" />
        </span>
        <span v-if="durationText" class="mp-chip mp-chip--time">{{ durationText }}</span>
        <span class="mp-chip mp-chip--type">视频</span>
      </template>

      <template v-else-if="kind === 'audio'">
        <div class="mp-audio-face">
          <i class="i-fe:music mp-audio-face__ico" />
          <div class="mp-bars" aria-hidden="true">
            <span v-for="n in 5" :key="n" />
          </div>
        </div>
        <span class="mp-chip mp-chip--type">音频</span>
        <span v-if="durationText" class="mp-chip mp-chip--time">{{ durationText }}</span>
      </template>

      <template v-else-if="kind === 'image'">
        <img
          v-if="!thumbBroken"
          class="mp-image"
          :src="displayMedia.url"
          alt=""
          @error="thumbBroken = true"
        >
        <i v-else class="i-fe:image mp-empty-ico" />
      </template>

      <template v-else-if="kind === 'document'">
        <img class="mp-doc-icon" :src="docIcon" alt="">
        <span class="mp-chip mp-chip--type">{{ fileExt ? fileExt.toUpperCase() : '文档' }}</span>
      </template>

      <template v-else-if="kind === 'app'">
        <i class="i-fe:package mp-app-ico" />
        <span class="mp-chip mp-chip--type">{{ fileExt ? fileExt.toUpperCase() : '安装包' }}</span>
      </template>

      <template v-else>
        <i class="i-fe:file mp-empty-ico" />
      </template>

      <span v-if="canPreview && interactive && kind !== 'video' && kind !== 'audio'" class="mp-expand" aria-hidden="true">
        <i class="i-fe:maximize-2" />
      </span>
    </div>

    <NModal
      v-if="interactive"
      v-model:show="visible"
      :mask-closable="true"
      :auto-focus="false"
      class="mp-modal"
      :style="{ width: 'auto', maxWidth: '96vw', padding: 0, background: 'transparent' }"
      @after-leave="teardown"
    >
      <div class="mp-stage" :class="[`mp-stage--${kind}`, { 'is-immersive': immersive }]" @click.stop>
        <header class="mp-stage__head">
          <div class="mp-stage__copy">
            <strong>{{ displayName }}</strong>
            <small>{{ metaLine }}</small>
          </div>
          <div class="mp-stage__actions">
            <button
              v-if="deletable"
              type="button"
              class="mp-text-btn mp-text-btn--danger"
              @click.stop="emitDelete"
            >
              <i class="i-material-symbols:delete-outline text-14" />
              删除
            </button>
            <button
              v-if="canImmersive && !mediaError"
              type="button"
              class="mp-icon-btn"
              :title="immersive ? '退出全屏' : '全屏'"
              @click.stop="toggleImmersive"
            >
              <i :class="immersive ? 'i-fe:minimize' : 'i-fe:maximize'" />
            </button>
            <a
              v-if="displayMedia.url"
              class="mp-icon-btn"
              :href="displayMedia.url"
              target="_blank"
              rel="noopener" referrerpolicy="origin"
              title="新窗口打开"
              @click.stop
            >
              <i class="i-fe:external-link" />
            </a>
            <a
              v-if="displayMedia.url"
              class="mp-icon-btn"
              :href="displayMedia.url"
              target="_blank"
              rel="noopener" referrerpolicy="origin"
              download
              title="下载"
              @click.stop
            >
              <i class="i-fe:download" />
            </a>
            <button type="button" class="mp-icon-btn mp-icon-btn--close" title="关闭" @click="visible = false">
              <i class="i-fe:x" />
            </button>
          </div>
        </header>

        <div class="mp-stage__body">
          <div class="mp-stage__canvas">
          <div v-if="mediaError" class="mp-status">
            <i class="i-fe:alert-circle" />
            <p>资源加载失败</p>
            <a v-if="displayMedia.url" :href="displayMedia.url" target="_blank" rel="noopener" referrerpolicy="origin">打开原文件</a>
          </div>

          <template v-else>
            <div
              v-if="kind === 'video'"
              ref="screenRef"
              class="mp-screen"
              :class="{ 'is-paused': !playing }"
              @mousemove="revealBar"
              @mouseleave="scheduleHideBar"
            >
              <video
                ref="videoRef"
                class="mp-video"
                playsinline
                preload="metadata"
                :poster="videoPoster || undefined"
                :src="displayMedia.url"
                @click="togglePlay"
                @canplay="mediaReady = true"
                @playing="onPlaying"
                @play="playing = true"
                @pause="playing = false"
                @ended="onMediaEnded"
                @timeupdate="onMediaTime"
                @loadedmetadata="onMediaMeta"
                @error="mediaError = true"
              />
              <button
                v-show="!playing"
                type="button"
                class="mp-big-play"
                @click.stop="togglePlay"
              >
                <i class="i-fe:play" />
              </button>
              <div class="mp-bar" :class="{ 'is-show': barShown || !playing }" @click.stop>
                <button type="button" class="mp-bar__btn" @click="togglePlay">
                  <i :class="playing ? 'i-fe:pause' : 'i-fe:play'" />
                </button>
                <span class="mp-bar__time">{{ formatClock(mediaCurrent) }}</span>
                <div
                  class="mp-bar__track"
                  @pointerdown="onSeekPointerDown"
                >
                  <i class="mp-bar__fill" :style="{ width: `${mediaProgress}%` }" />
                  <i class="mp-bar__knob" :style="{ left: `${mediaProgress}%` }" />
                </div>
                <span class="mp-bar__time">{{ formatClock(mediaDuration || props.mediaInfo?.duration || 0) }}</span>
                <button type="button" class="mp-bar__btn" @click="toggleMute">
                  <i :class="muted ? 'i-fe:volume-x' : 'i-fe:volume-2'" />
                </button>
                <button type="button" class="mp-bar__btn" @click="toggleFullscreen">
                  <i class="i-fe:maximize" />
                </button>
              </div>
            </div>

            <div v-else-if="kind === 'audio'" class="mp-audio">
              <div class="mp-audio__art" :class="{ 'is-on': playing }" @click="togglePlay">
                <i class="i-fe:music" />
                <div class="mp-bars mp-bars--lg" aria-hidden="true">
                  <span v-for="n in 7" :key="n" />
                </div>
              </div>
              <div class="mp-audio__meta">
                <b>{{ displayName }}</b>
                <em>{{ sizeText || '音频文件' }}</em>
              </div>
              <div class="mp-bar mp-bar--dock is-show">
                <button type="button" class="mp-bar__btn" @click="togglePlay">
                  <i :class="playing ? 'i-fe:pause' : 'i-fe:play'" />
                </button>
                <span class="mp-bar__time">{{ formatClock(mediaCurrent) }}</span>
                <div class="mp-bar__track" @pointerdown="onSeekPointerDown">
                  <i class="mp-bar__fill" :style="{ width: `${mediaProgress}%` }" />
                  <i class="mp-bar__knob" :style="{ left: `${mediaProgress}%` }" />
                </div>
                <span class="mp-bar__time">{{ formatClock(mediaDuration || props.mediaInfo?.duration || 0) }}</span>
              </div>
              <audio
                ref="audioRef"
                class="mp-audio__el"
                preload="metadata"
                :src="displayMedia.url"
                @play="playing = true"
                @pause="playing = false"
                @ended="onMediaEnded"
                @timeupdate="onMediaTime"
                @loadedmetadata="onMediaMeta"
                @error="mediaError = true"
              />
            </div>

            <img
              v-else-if="kind === 'image'"
              class="mp-photo"
              :src="displayMedia.url"
              :alt="displayName"
              @load="mediaReady = true"
              @error="mediaError = true"
            >

            <OfficeDocViewer
              v-else-if="kind === 'document'"
              class="mp-doc"
              :src="displayMedia.url"
              :filename="docFilename"
              @loading="onDocLoading"
              @error="mediaError = true"
              @rendered="mediaError = false"
            />

            <div v-else-if="kind === 'app'" class="mp-app">
              <i class="i-fe:package" />
              <p>安装包无法在线预览</p>
              <a :href="displayMedia.url" target="_blank" rel="noopener" referrerpolicy="origin">下载文件</a>
            </div>

            <div v-else class="mp-app">
              <p>无法预览该类型</p>
            </div>
          </template>

          <div v-if="showSpinner" class="mp-status mp-status--overlay">
            <span class="mp-spinner" />
            <p>加载中</p>
          </div>
          </div>
        </div>
      </div>
    </NModal>
  </div>
</template>

<script setup>
import { NModal } from 'naive-ui'
import { computed, onBeforeUnmount, ref, watch } from 'vue'
import IconDoc from '@/assets/icons/document/doc.png'
import IconExcel from '@/assets/icons/document/excel.png'
import IconPdf from '@/assets/icons/document/pdf.png'
import IconPpt from '@/assets/icons/document/ppt.png'
import IconWord from '@/assets/icons/document/word.png'
import { getAssetUrl } from '@/utils/common'
import OfficeDocViewer from './OfficeDocViewer.vue'

const props = defineProps({
  mediaInfo: {
    type: Object,
    default: () => ({
      url: '',
      type: '',
      thumb: '',
      name: '',
    }),
  },
  width: { type: [Number, String], default: undefined },
  height: { type: [Number, String], default: undefined },
  controls: { type: Boolean, default: true },
  interactive: { type: Boolean, default: true },
  compact: { type: Boolean, default: false },
  deletable: { type: Boolean, default: false },
})

const emit = defineEmits(['delete'])

const visible = ref(false)
const playing = ref(false)
const muted = ref(false)
const barShown = ref(true)
const mediaCurrent = ref(0)
const mediaDuration = ref(0)
const mediaReady = ref(false)
const mediaError = ref(false)
const thumbBroken = ref(false)
const docLoading = ref(false)
const immersive = ref(false)
const videoRef = ref(null)
const audioRef = ref(null)
const screenRef = ref(null)
let hideBarTimer = 0
let seeking = false

const displayMedia = computed(() => ({
  ...props.mediaInfo,
  url: props.mediaInfo?.url ? getAssetUrl(props.mediaInfo.url) : '',
  thumb: props.mediaInfo?.thumb ? getAssetUrl(props.mediaInfo.thumb) : '',
}))

const canPreview = computed(() => Boolean(displayMedia.value.url))

const fileExt = computed(() => {
  const fromName = String(props.mediaInfo?.title || props.mediaInfo?.name || '')
  const fromUrl = String(props.mediaInfo?.url || displayMedia.value.url || '')
  const source = fromName.includes('.') ? fromName : fromUrl
  const match = source.split('?')[0].split('.').pop()
  return match ? match.toLowerCase() : ''
})

const kind = computed(() => {
  const type = String(props.mediaInfo?.type || '').toLowerCase()
  if (['video', 'audio', 'image', 'document', 'app'].includes(type))
    return type
  const ext = fileExt.value
  if (['mp4', 'webm', 'mov', 'm3u8', 'avi', 'mkv'].includes(ext))
    return 'video'
  if (['mp3', 'wav', 'm4a', 'aac', 'ogg', 'flac'].includes(ext))
    return 'audio'
  if (['png', 'jpg', 'jpeg', 'gif', 'webp', 'bmp', 'svg'].includes(ext))
    return 'image'
  if (['apk', 'ipa'].includes(ext))
    return 'app'
  if (ext)
    return 'document'
  return 'unknown'
})

const canImmersive = computed(() => ['document', 'video', 'image'].includes(kind.value))

const videoPoster = computed(() => displayMedia.value.thumb || '')

const thumbStyle = computed(() => {
  const style = {}
  if (props.width != null && props.width !== '')
    style.width = typeof props.width === 'number' ? `${props.width}px` : props.width
  if (props.height != null && props.height !== '')
    style.height = typeof props.height === 'number' ? `${props.height}px` : props.height
  return style
})

const displayName = computed(() => {
  return props.mediaInfo?.title || props.mediaInfo?.name || '未命名资源'
})

function formatSize(bytes) {
  const n = Number(bytes) || 0
  if (n <= 0)
    return ''
  if (n < 1024)
    return `${n} B`
  if (n < 1024 * 1024)
    return `${(n / 1024).toFixed(1)} KB`
  if (n < 1024 * 1024 * 1024)
    return `${(n / 1024 / 1024).toFixed(1)} MB`
  return `${(n / 1024 / 1024 / 1024).toFixed(2)} GB`
}

function formatDuration(sec) {
  const n = Math.round(Number(sec) || 0)
  if (n <= 0)
    return ''
  const h = Math.floor(n / 3600)
  const m = Math.floor((n % 3600) / 60)
  const s = n % 60
  if (h)
    return `${h}:${String(m).padStart(2, '0')}:${String(s).padStart(2, '0')}`
  return `${m}:${String(s).padStart(2, '0')}`
}

function formatClock(sec) {
  const n = Math.max(0, Math.round(Number(sec) || 0))
  const h = Math.floor(n / 3600)
  const m = Math.floor((n % 3600) / 60)
  const s = n % 60
  if (h)
    return `${h}:${String(m).padStart(2, '0')}:${String(s).padStart(2, '0')}`
  return `${m}:${String(s).padStart(2, '0')}`
}

const sizeText = computed(() => formatSize(props.mediaInfo?.size))
const durationText = computed(() => formatDuration(props.mediaInfo?.duration))

const mediaProgress = computed(() => {
  if (!mediaDuration.value)
    return 0
  return Math.min(100, (mediaCurrent.value / mediaDuration.value) * 100)
})

const showSpinner = computed(() => {
  if (mediaError.value)
    return false
  if (kind.value === 'image')
    return visible.value && !mediaReady.value
  return false
})
const typeLabel = computed(() => {
  return {
    video: '视频',
    audio: '音频',
    image: '图片',
    document: '文档',
    app: '安装包',
  }[kind.value] || '资源'
})
const metaLine = computed(() => {
  return [typeLabel.value, sizeText.value, durationText.value].filter(Boolean).join(' · ')
})

const docIcons = {
  pdf: IconPdf,
  doc: IconWord,
  docx: IconWord,
  xls: IconExcel,
  xlsx: IconExcel,
  ppt: IconPpt,
  pptx: IconPpt,
  default: IconDoc,
}
const docIcon = computed(() => docIcons[fileExt.value] || docIcons.default)

const docFilename = computed(() => {
  const name = props.mediaInfo?.title || props.mediaInfo?.name
  if (name && /\.[a-z0-9]{2,8}$/i.test(String(name)))
    return String(name)
  const u = props.mediaInfo?.url
  if (u) {
    try {
      const seg = new URL(u, typeof location !== 'undefined' ? location.href : undefined).pathname.split('/').filter(Boolean).pop()
      if (seg?.includes('.'))
        return decodeURIComponent(seg)
    }
    catch {
      const seg = String(u).split('?')[0].split('/').pop()
      if (seg?.includes('.'))
        return decodeURIComponent(seg)
    }
  }
  return `document.${fileExt.value || 'pdf'}`
})

function onDocLoading(val) {
  docLoading.value = Boolean(val)
}

function toggleImmersive() {
  immersive.value = !immersive.value
}

function teardown() {
  clearTimeout(hideBarTimer)
  if (document.fullscreenElement)
    document.exitFullscreen?.().catch(() => {})
  videoRef.value?.pause?.()
  audioRef.value?.pause?.()
  playing.value = false
  muted.value = false
  mediaCurrent.value = 0
  mediaReady.value = false
  mediaError.value = false
  docLoading.value = false
  immersive.value = false
}

function emitDelete() {
  visible.value = false
  emit('delete')
}

function resetMediaState() {
  mediaReady.value = false
  mediaError.value = false
  thumbBroken.value = false
  mediaCurrent.value = 0
  mediaDuration.value = 0
  playing.value = false
  muted.value = false
  barShown.value = true
}

function mediaEl() {
  return kind.value === 'audio' ? audioRef.value : videoRef.value
}

function togglePlay() {
  const el = mediaEl()
  if (!el)
    return
  if (el.paused)
    el.play().catch(() => { mediaError.value = true })
  else
    el.pause()
}

function onPlaying() {
  playing.value = true
  mediaReady.value = true
  scheduleHideBar()
}

function onMediaTime() {
  if (seeking)
    return
  mediaCurrent.value = mediaEl()?.currentTime || 0
}

function onMediaMeta() {
  const el = mediaEl()
  mediaDuration.value = el?.duration || 0
  if (el)
    el.muted = muted.value
}

function onMediaEnded() {
  playing.value = false
  mediaCurrent.value = 0
  barShown.value = true
}

function applySeek(clientX, track) {
  const el = mediaEl()
  if (!el || !el.duration)
    return
  const rect = track.getBoundingClientRect()
  const ratio = Math.min(1, Math.max(0, (clientX - rect.left) / rect.width))
  el.currentTime = ratio * el.duration
  mediaCurrent.value = el.currentTime
}

function onSeekPointerDown(e) {
  const track = e.currentTarget
  seeking = true
  track.setPointerCapture?.(e.pointerId)
  applySeek(e.clientX, track)
  const move = (ev) => applySeek(ev.clientX, track)
  const up = () => {
    seeking = false
    window.removeEventListener('pointermove', move)
    window.removeEventListener('pointerup', up)
  }
  window.addEventListener('pointermove', move)
  window.addEventListener('pointerup', up)
}

function toggleMute() {
  const el = mediaEl()
  if (!el)
    return
  el.muted = !el.muted
  muted.value = el.muted
}

function toggleFullscreen() {
  const node = screenRef.value
  if (!node)
    return
  if (document.fullscreenElement)
    document.exitFullscreen?.().catch(() => {})
  else
    node.requestFullscreen?.().catch(() => {})
}

function revealBar() {
  barShown.value = true
  scheduleHideBar()
}

function scheduleHideBar() {
  clearTimeout(hideBarTimer)
  if (!playing.value)
    return
  hideBarTimer = window.setTimeout(() => {
    if (playing.value)
      barShown.value = false
  }, 1600)
}

function openPreview(e) {
  e?.stopPropagation?.()
  if (!canPreview.value || !props.interactive)
    return
  resetMediaState()
  visible.value = true
}

watch(visible, (val) => {
  if (!val)
    teardown()
})

watch(
  () => [props.mediaInfo?.url, props.mediaInfo?.thumb],
  () => {
    thumbBroken.value = false
  },
)

onBeforeUnmount(() => {
  teardown()
})
</script>

<style scoped lang="scss">
.mp {
  display: flex;
  width: 100%;
  height: 100%;
  min-width: 0;
  min-height: 0;
  align-items: stretch;
  justify-content: stretch;
  border-radius: 12px;
}

.mp.is-previewable {
  cursor: pointer;
}

.mp-thumb {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  width: 100%;
  height: 100%;
  min-height: 72px;
  border-radius: 12px;
  background: var(--ctl-field, rgba(var(--primary-color), 0.04));
  transition: filter 0.22s cubic-bezier(0.32, 0.72, 0, 1);
}

.mp.is-previewable:hover .mp-play {
  transform: scale(1.06);
}

.mp.is-previewable:active .mp-thumb {
  filter: brightness(0.96);
}

.mp-thumb--video {
  background: #102224;
}

.mp-thumb--video::after {
  content: '';
  position: absolute;
  inset: 0;
  z-index: 0;
  pointer-events: none;
  background: linear-gradient(180deg, transparent 42%, rgb(8 18 20 / 62%) 100%);
}

.mp-thumb--audio {
  background:
    radial-gradient(circle at 30% 20%, rgba(var(--primary-color), 0.22), transparent 55%),
    linear-gradient(160deg, rgba(var(--primary-color), 0.16) 0%, rgb(16 32 34 / 8%) 100%);
}

.mp-thumb--document,
.mp-thumb--app,
.mp-thumb--unknown {
  flex-direction: column;
  gap: 8px;
}

.mp-image {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
}

.mp-image--poster {
  z-index: 0;
}

.mp-play {
  position: relative;
  z-index: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 44px;
  height: 44px;
  border-radius: 50%;
  color: #fff;
  background: rgb(var(--primary-color));
  box-shadow: 0 8px 20px rgba(var(--primary-color), 0.28);
  transition: transform 0.22s cubic-bezier(0.32, 0.72, 0, 1);
}

.mp-play i {
  margin-left: 2px;
  font-size: 18px;
}

.mp-expand {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  background: rgb(8 16 18 / 32%);
  opacity: 0;
  transition: opacity 0.2s ease;
}

.mp-expand i {
  font-size: 18px;
}

.mp.is-previewable:hover .mp-expand {
  opacity: 1;
}

.mp.is-compact .mp-chip--type {
  display: none;
}

.mp.is-compact .mp-play {
  width: 32px;
  height: 32px;
}

.mp.is-compact .mp-play i {
  font-size: 14px;
}

.mp-chip {
  position: absolute;
  z-index: 1;
  height: 20px;
  padding: 0 7px;
  border-radius: 999px;
  font-size: 11px;
  font-weight: 650;
  letter-spacing: -0.01em;
  line-height: 20px;
  color: #fff;
  background: rgb(8 16 18 / 58%);
}

.mp-chip--type {
  top: 8px;
  left: 8px;
}

.mp-chip--time {
  left: 8px;
  bottom: 8px;
  font-variant-numeric: tabular-nums;
}

.mp-audio-face {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 10px;
  color: rgb(var(--primary-color));
}

.mp-audio-face__ico,
.mp-app-ico,
.mp-empty-ico {
  font-size: 28px;
}

.mp-app-ico,
.mp-empty-ico {
  color: rgb(var(--primary-color));
}

.mp-bars {
  display: flex;
  align-items: flex-end;
  gap: 3px;
  height: 16px;
}

.mp-bars span {
  width: 3px;
  height: 6px;
  border-radius: 99px;
  background: rgb(var(--primary-color));
  opacity: 0.75;
  transform-origin: bottom;
  animation: mp-bar 1.1s ease-in-out infinite;
}

.mp-bars span:nth-child(2) { animation-delay: 0.12s; height: 11px; }
.mp-bars span:nth-child(3) { animation-delay: 0.22s; height: 16px; }
.mp-bars span:nth-child(4) { animation-delay: 0.08s; height: 9px; }
.mp-bars span:nth-child(5) { animation-delay: 0.28s; height: 13px; }
.mp-bars span:nth-child(6) { animation-delay: 0.18s; height: 8px; }
.mp-bars span:nth-child(7) { animation-delay: 0.3s; height: 14px; }

.mp-bars--lg {
  height: 28px;
}

.mp-bars--lg span {
  width: 4px;
}

.mp-audio__art:not(.is-on) .mp-bars span,
.mp-thumb--audio .mp-bars span {
  animation: none;
}

.mp-audio__art.is-on .mp-bars span {
  animation: mp-bar 1.1s ease-in-out infinite;
}

@keyframes mp-bar {
  0%,
  100% { transform: scaleY(0.45); }
  50% { transform: scaleY(1); }
}

.mp-doc-icon {
  width: 42px;
  height: 42px;
  object-fit: contain;
}

.mp-stage {
  width: max-content;
  max-width: min(1280px, 96vw);
  min-width: min(320px, 96vw);
  overflow: hidden;
  border-radius: 22px;
  border: 1px solid rgb(255 255 255 / 55%);
  background:
    linear-gradient(180deg, #ffffff 0%, #f5f8f8 100%);
  color: #1c2b2e;
  box-shadow:
    0 0 0 1px rgb(15 40 45 / 4%),
    0 28px 80px rgb(15 40 45 / 22%),
    0 8px 24px rgb(15 40 45 / 08%);
  animation: mp-stage-in 0.38s cubic-bezier(0.22, 1, 0.36, 1);
}

.mp-stage--video {
  width: min(1080px, 94vw);
  max-width: 96vw;
}

.mp-stage--document {
  width: min(1480px, 96vw);
  max-width: 96vw;
}

.mp-stage.is-immersive {
  width: 100vw !important;
  max-width: 100vw !important;
  min-width: 0;
  height: 100vh;
  max-height: 100vh;
  border-radius: 0;
  border: 0;
  display: flex;
  flex-direction: column;
  box-shadow: none;
  animation: none;
}

.mp-stage.is-immersive .mp-stage__head {
  flex-shrink: 0;
  border-radius: 0;
}

.mp-stage.is-immersive .mp-stage__body {
  flex: 1;
  min-height: 0;
  padding: 0;
  align-items: stretch;
}

.mp-stage.is-immersive .mp-stage__canvas {
  width: 100%;
  height: 100%;
  min-height: 0;
  border-radius: 0;
}

.mp-stage.is-immersive .mp-doc,
.mp-stage.is-immersive .mp-screen,
.mp-stage.is-immersive .mp-photo {
  min-height: 0;
  width: 100%;
  height: 100%;
}

.mp-stage.is-immersive .mp-photo {
  object-fit: contain;
  background: #0f1c20;
}

.mp-stage.is-immersive.mp-stage--video .mp-stage__body,
.mp-stage.is-immersive.mp-stage--image .mp-stage__body {
  background: #0f1c20;
}

.mp-stage.is-immersive.mp-stage--video .mp-stage__canvas,
.mp-stage.is-immersive.mp-stage--image .mp-stage__canvas {
  background: #0f1c20;
  box-shadow: none;
}

.mp-stage--audio,
.mp-stage--app,
.mp-stage--unknown {
  width: min(480px, 94vw);
  min-width: min(320px, 94vw);
}

@keyframes mp-stage-in {
  from {
    opacity: 0;
    transform: translateY(12px) scale(0.98);
  }
  to {
    opacity: 1;
    transform: none;
  }
}

.mp-stage--document,
.mp-stage--app,
.mp-stage--unknown,
.mp-stage--video,
.mp-stage--audio,
.mp-stage--image {
  background: linear-gradient(180deg, #ffffff 0%, #f5f8f8 100%);
  color: #1c2b2e;
}

.mp-stage__head {
  display: flex;
  align-items: center;
  gap: 16px;
  min-height: 64px;
  padding: 14px 16px 14px 22px;
  border-bottom: 1px solid rgb(15 40 45 / 06%);
  background: rgb(255 255 255 / 92%);
  backdrop-filter: blur(16px) saturate(1.2);
  -webkit-backdrop-filter: blur(16px) saturate(1.2);
}

.mp-stage__copy {
  min-width: 0;
  flex: 1;
}

.mp-stage__copy strong {
  display: block;
  overflow: hidden;
  font-family: var(--font-display, inherit);
  font-size: 16px;
  font-weight: 650;
  letter-spacing: -0.03em;
  line-height: 1.25;
  color: #1c2b2e;
  white-space: nowrap;
  text-overflow: ellipsis;
}

.mp-stage__copy small {
  display: block;
  margin-top: 3px;
  font-size: 12px;
  font-weight: 500;
  letter-spacing: -0.01em;
  color: #6b7a7c;
}

.mp-stage__actions {
  display: flex;
  align-items: center;
  gap: 8px;
}

.mp-text-btn {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  height: 34px;
  padding: 0 14px;
  border: 0;
  border-radius: 999px;
  font-size: 13px;
  font-weight: 550;
  letter-spacing: -0.01em;
  cursor: pointer;
  transition:
    transform 0.2s cubic-bezier(0.32, 0.72, 0, 1),
    background 0.2s ease,
    color 0.2s ease;
}

.mp-text-btn--danger {
  color: #c23b32;
  background: rgb(194 59 50 / 10%);
}

.mp-text-btn--danger:hover {
  background: rgb(194 59 50 / 16%);
}

.mp-text-btn:active {
  transform: scale(0.96);
}

.mp-icon-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 34px;
  height: 34px;
  border: 0;
  border-radius: 999px;
  color: #1c2b2e;
  background: rgba(var(--primary-color), 0.08);
  text-decoration: none;
  cursor: pointer;
  transition:
    transform 0.2s cubic-bezier(0.32, 0.72, 0, 1),
    background 0.2s ease,
    color 0.2s ease;
}

.mp-icon-btn--close {
  background: rgb(15 40 45 / 06%);
  color: #3a4a4c;
}

.mp-icon-btn:hover {
  background: rgba(var(--primary-color), 0.14);
}

.mp-icon-btn--close:hover {
  background: rgb(15 40 45 / 10%);
  color: #1c2b2e;
}

.mp-icon-btn:active {
  transform: scale(0.92);
}

.mp-stage__body {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 0;
  padding: 14px 16px 16px;
  background:
    radial-gradient(120% 80% at 50% 0%, rgba(var(--primary-color), 0.06), transparent 62%),
    #e8efef;
}

.mp-stage__canvas {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  width: max-content;
  max-width: 100%;
  min-height: 0;
  border-radius: 16px;
  overflow: hidden;
  background: #f3f7f7;
  box-shadow:
    inset 0 0 0 1px rgb(15 40 45 / 05%),
    0 1px 0 rgb(255 255 255 / 70%);
}

.mp-stage--video .mp-stage__canvas,
.mp-stage--document .mp-stage__canvas,
.mp-stage--image .mp-stage__canvas {
  width: 100%;
}

.mp-stage--audio .mp-stage__canvas,
.mp-stage--app .mp-stage__canvas,
.mp-stage--unknown .mp-stage__canvas {
  width: 100%;
  background: transparent;
  box-shadow: none;
}

.mp-status {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 48px 20px;
  text-align: center;
}

.mp-status a {
  color: rgb(var(--primary-color));
  font-weight: 600;
  text-decoration: none;
}

.mp-status--overlay {
  position: absolute;
  inset: 0;
  z-index: 2;
  border-radius: 16px;
  background: rgb(255 255 255 / 78%);
  backdrop-filter: blur(8px);
  -webkit-backdrop-filter: blur(8px);
  color: #1c2b2e;
}

.mp-spinner {
  width: 22px;
  height: 22px;
  border: 2px solid rgba(var(--primary-color), 0.16);
  border-top-color: rgb(var(--primary-color));
  border-radius: 50%;
  animation: mp-spin 0.75s linear infinite;
}

@keyframes mp-spin {
  to { transform: rotate(360deg); }
}

.mp-screen {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  min-height: 0;
  overflow: hidden;
  border-radius: 0;
  background: #0c1214;
  cursor: pointer;
}

.mp-video {
  display: block;
  width: 100%;
  max-height: min(78vh, 720px);
  background: #0c1214;
}

.mp-screen:fullscreen,
.mp-screen:-webkit-full-screen {
  width: 100%;
  height: 100%;
  background: #000;
}

.mp-screen:fullscreen .mp-video,
.mp-screen:-webkit-full-screen .mp-video {
  max-height: 100%;
  height: 100%;
  object-fit: contain;
}

.mp-big-play {
  position: absolute;
  z-index: 2;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 72px;
  height: 72px;
  border: 0;
  border-radius: 50%;
  color: #fff;
  background: rgb(var(--primary-color) / 92%);
  box-shadow:
    0 12px 32px rgba(var(--primary-color), 0.35),
    inset 0 1px 0 rgb(255 255 255 / 25%);
  cursor: pointer;
  backdrop-filter: blur(8px);
  transition: transform 0.22s cubic-bezier(0.32, 0.72, 0, 1);
}

.mp-big-play i {
  margin-left: 3px;
  font-size: 28px;
}

.mp-big-play:hover {
  transform: scale(1.06);
}

.mp-big-play:active {
  transform: scale(0.94);
}

.mp-bar {
  position: absolute;
  left: 16px;
  right: 16px;
  bottom: 16px;
  z-index: 3;
  display: flex;
  align-items: center;
  gap: 10px;
  height: 48px;
  padding: 0 12px 0 10px;
  border-radius: 14px;
  border: 1px solid rgb(255 255 255 / 10%);
  background: rgb(16 24 26 / 72%);
  color: #eef3f3;
  backdrop-filter: blur(18px) saturate(1.3);
  -webkit-backdrop-filter: blur(18px) saturate(1.3);
  opacity: 0;
  transform: translateY(8px);
  pointer-events: none;
  transition:
    opacity 0.22s cubic-bezier(0.32, 0.72, 0, 1),
    transform 0.22s cubic-bezier(0.32, 0.72, 0, 1);
}

.mp-bar.is-show {
  opacity: 1;
  transform: none;
  pointer-events: auto;
}

.mp-bar--dock {
  position: relative;
  left: auto;
  right: auto;
  bottom: auto;
  width: min(100%, 520px);
  pointer-events: auto;
  background: rgb(28 43 46 / 88%);
}

.mp-bar__btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 34px;
  height: 34px;
  flex-shrink: 0;
  border: 0;
  border-radius: 999px;
  color: inherit;
  background: transparent;
  cursor: pointer;
  transition: background 0.18s ease;
}

.mp-bar__btn:hover {
  background: rgb(255 255 255 / 12%);
}

.mp-bar__btn i {
  font-size: 16px;
}

.mp-bar__btn .i-fe-play {
  margin-left: 1px;
}

.mp-bar__time {
  min-width: 42px;
  font-size: 12px;
  font-variant-numeric: tabular-nums;
  letter-spacing: -0.02em;
  opacity: 0.72;
}

.mp-bar__track {
  position: relative;
  flex: 1;
  height: 4px;
  border-radius: 99px;
  background: rgb(255 255 255 / 16%);
  cursor: pointer;
  transition: height 0.15s ease;
}

.mp-bar__track:hover {
  height: 6px;
}

.mp-bar__fill {
  display: block;
  height: 100%;
  border-radius: inherit;
  background: #fff;
}

.mp-bar__knob {
  position: absolute;
  top: 50%;
  width: 13px;
  height: 13px;
  margin-left: -6.5px;
  border-radius: 50%;
  background: #fff;
  box-shadow: 0 1px 6px rgb(0 0 0 / 35%);
  transform: translateY(-50%);
  pointer-events: none;
}

.mp-photo {
  display: block;
  max-width: min(92vw, 1200px);
  max-height: min(78vh, 800px);
  width: auto;
  height: auto;
  object-fit: contain;
  border-radius: 0;
  background:
    linear-gradient(45deg, rgb(49 108 114 / 5%) 25%, transparent 25%),
    linear-gradient(-45deg, rgb(49 108 114 / 5%) 25%, transparent 25%),
    linear-gradient(45deg, transparent 75%, rgb(49 108 114 / 5%) 75%),
    linear-gradient(-45deg, transparent 75%, rgb(49 108 114 / 5%) 75%);
  background-size: 18px 18px;
  background-position: 0 0, 0 9px, 9px -9px, -9px 0;
  background-color: #f3f7f7;
}

.mp-audio {
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 22px;
  width: min(560px, 100%);
  margin: auto;
  padding: 36px 28px 32px;
  border-radius: 20px;
  background: rgb(255 255 255 / 78%);
  border: 1px solid rgb(255 255 255 / 70%);
  box-shadow: 0 12px 32px rgb(15 40 45 / 06%);
}

.mp-audio__art {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 14px;
  width: 168px;
  height: 168px;
  border-radius: 36px;
  color: #fff;
  cursor: pointer;
  background:
    radial-gradient(circle at 30% 20%, rgba(var(--primary-color), 0.55), transparent 58%),
    linear-gradient(160deg, #2b5d62 0%, #132628 100%);
  box-shadow: 0 16px 40px rgba(var(--primary-color), 0.28);
  transition: transform 0.22s cubic-bezier(0.32, 0.72, 0, 1);
}

.mp-audio__art:hover {
  transform: scale(1.03);
}

.mp-audio__art.is-on {
  animation: mp-art-pulse 2.4s ease-in-out infinite;
}

@keyframes mp-art-pulse {
  0%, 100% { box-shadow: 0 16px 40px rgba(var(--primary-color), 0.28); }
  50% { box-shadow: 0 20px 48px rgba(var(--primary-color), 0.4); }
}

.mp-audio__art i {
  font-size: 40px;
}

.mp-audio__meta {
  text-align: center;
}

.mp-audio__meta b {
  display: block;
  font-family: var(--font-display, inherit);
  font-size: 18px;
  font-weight: 650;
  letter-spacing: -0.028em;
  color: #1c2b2e;
}

.mp-audio__meta em {
  display: block;
  margin-top: 4px;
  font-style: normal;
  font-size: 13px;
  color: #6b7a7c;
}

.mp-audio__el {
  position: absolute;
  width: 0;
  height: 0;
  opacity: 0;
  pointer-events: none;
}

.mp-stage--document .mp-status--overlay {
  background: rgb(255 255 255 / 82%);
  color: #1c2b2e;
}

.mp-stage--document .mp-spinner {
  border-color: rgba(var(--primary-color), 0.18);
  border-top-color: rgb(var(--primary-color));
}

.mp-stage--document .mp-stage__canvas {
  background: #f4f7f7;
  min-height: min(78vh, 760px);
  height: min(78vh, 760px);
  box-shadow: none;
}

.mp-doc {
  display: flex;
  flex-direction: column;
  width: 100%;
  height: 100%;
  min-height: min(78vh, 760px);
  overflow: hidden;
  border-radius: 0;
  border: 0;
  background: transparent;
}

.mp-app {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  width: min(420px, 100%);
  margin: auto;
  padding: 56px 28px;
  border-radius: 20px;
  color: #5a6b6d;
  background: rgb(255 255 255 / 78%);
  border: 1px solid rgb(255 255 255 / 70%);
  box-shadow: 0 12px 32px rgb(15 40 45 / 06%);
}

.mp-app i {
  font-size: 40px;
  color: rgb(var(--primary-color));
}

.mp-app a {
  color: rgb(var(--primary-color));
  text-decoration: none;
  font-weight: 600;
}

.mp-app a:hover {
  text-decoration: underline;
}
</style>

<style>
.mp-modal.n-modal {
  padding: 0 !important;
  background: transparent !important;
  box-shadow: none !important;
  width: auto !important;
  max-width: 96vw !important;
  border: 0 !important;
}

.mp-modal.n-modal:has(.mp-stage.is-immersive) {
  max-width: 100vw !important;
  width: 100vw !important;
  margin: 0 !important;
}

.mp-modal.n-modal .n-card,
.mp-modal .n-card {
  width: auto !important;
  max-width: 96vw !important;
  background: transparent !important;
  border: 0 !important;
  box-shadow: none !important;
  --n-color-modal: transparent !important;
  --n-color-target: transparent !important;
}

.mp-modal.n-modal:has(.mp-stage.is-immersive) .n-card,
.mp-modal:has(.mp-stage.is-immersive) .n-card {
  max-width: 100vw !important;
  width: 100vw !important;
  height: 100vh !important;
  border-radius: 0 !important;
}

.mp-modal.n-modal .n-card > .n-card__content,
.mp-modal .n-card .n-card__content {
  padding: 0 !important;
  width: auto !important;
  max-width: 96vw !important;
}

.mp-modal.n-modal:has(.mp-stage.is-immersive) .n-card > .n-card__content,
.mp-modal:has(.mp-stage.is-immersive) .n-card .n-card__content {
  max-width: 100vw !important;
  width: 100vw !important;
  height: 100vh !important;
}

.n-modal-container:has(.mp-modal) .n-modal-mask {
  background: rgb(15 28 32 / 48%) !important;
  backdrop-filter: blur(10px) saturate(1.1);
  -webkit-backdrop-filter: blur(10px) saturate(1.1);
}

.n-modal-container:has(.mp-stage.is-immersive) {
  display: flex !important;
  align-items: stretch !important;
  justify-content: stretch !important;
}

.n-modal-container:has(.mp-stage.is-immersive) .n-modal-body-wrapper,
.n-modal-container:has(.mp-stage.is-immersive) .n-modal-scroll-content {
  padding: 0 !important;
  max-width: 100vw !important;
  width: 100vw !important;
  height: 100vh !important;
}

.n-modal-container:has(.mp-modal) .n-card {
  animation: none !important;
  transform: none !important;
}
</style>
