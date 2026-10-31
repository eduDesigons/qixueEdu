<script setup>
import { nextTick, onBeforeUnmount, ref, watch } from 'vue'

const props = defineProps({
  src: { type: String, default: '' },
})

const emit = defineEmits(['rendered', 'error', 'loading', 'progress'])

const hostRef = ref(null)
const progressText = ref('')

let viewer = null
let abortCtrl = null
let loadSeq = 0

async function fetchBuffer(url, signal, onProgress) {
  const res = await fetch(url, { signal })
  if (!res.ok)
    throw new Error(`HTTP ${res.status}`)

  const total = Number(res.headers.get('content-length') || 0)
  if (!res.body || !total || !res.body.getReader)
    return res.arrayBuffer()

  const reader = res.body.getReader()
  const chunks = []
  let received = 0
  while (true) {
    const { done, value } = await reader.read()
    if (done)
      break
    chunks.push(value)
    received += value.byteLength
    onProgress?.(Math.min(99, Math.round((received / total) * 100)), received, total)
  }

  const out = new Uint8Array(received)
  let offset = 0
  for (const chunk of chunks) {
    out.set(chunk, offset)
    offset += chunk.byteLength
  }
  return out.buffer
}

function destroyViewer() {
  try {
    viewer?.destroy?.()
  }
  catch {
    // ignore
  }
  viewer = null
  if (hostRef.value)
    hostRef.value.innerHTML = ''
}

async function openPptx(src) {
  const seq = ++loadSeq
  abortCtrl?.abort()
  abortCtrl = new AbortController()
  destroyViewer()

  emit('loading', true)
  emit('progress', 0)
  progressText.value = '下载中…'

  try {
    await nextTick()
    const host = hostRef.value
    if (!host || seq !== loadSeq)
      return

    const buffer = await fetchBuffer(src, abortCtrl.signal, (pct) => {
      if (seq !== loadSeq)
        return
      emit('progress', pct)
      progressText.value = `下载中 ${pct}%`
    })
    if (seq !== loadSeq)
      return

    progressText.value = '解析首屏…'
    emit('progress', 100)

    const { PptxViewer, RECOMMENDED_ZIP_LIMITS } = await import('@aiden0z/pptx-renderer')
    if (seq !== loadSeq)
      return

    let firstPaint = false
    const markReady = () => {
      if (firstPaint || seq !== loadSeq)
        return
      firstPaint = true
      progressText.value = ''
      emit('loading', false)
      emit('rendered')
    }

    viewer = await PptxViewer.open(buffer, host, {
      fitMode: 'contain',
      pdfjs: false,
      lazyMedia: true,
      lazySlides: true,
      zipLimits: RECOMMENDED_ZIP_LIMITS,
      scrollContainer: host,
      listOptions: {
        windowed: true,
        initialSlides: 1,
        batchSize: 2,
        overscanViewport: 1.25,
        showSlideLabels: false,
      },
      onSlideRendered: (index) => {
        if (index === 0)
          markReady()
      },
      onRenderComplete: () => markReady(),
      onSlideError: (_index, err) => {
        if (!firstPaint)
          emit('error', err)
      },
    })

    if (seq !== loadSeq)
      return

    if (!firstPaint)
      markReady()
  }
  catch (err) {
    if (err?.name === 'AbortError' || seq !== loadSeq)
      return
    progressText.value = ''
    emit('loading', false)
    emit('error', err)
  }
}

watch(
  () => props.src,
  (src) => {
    if (src)
      openPptx(src)
    else
      destroyViewer()
  },
  { immediate: true },
)

onBeforeUnmount(() => {
  loadSeq += 1
  abortCtrl?.abort()
  destroyViewer()
})
</script>

<template>
  <div class="pptx-lazy">
    <div ref="hostRef" class="pptx-lazy__host" />
    <div v-if="progressText" class="pptx-lazy__hint">
      {{ progressText }}
    </div>
  </div>
</template>

<style scoped lang="scss">
.pptx-lazy {
  position: relative;
  width: 100%;
  height: 100%;
  min-height: 0;
  overflow: hidden;
  background: #111;
}

.pptx-lazy__host {
  width: 100%;
  height: 100%;
  overflow: auto;
  overscroll-behavior: contain;
}

.pptx-lazy__hint {
  position: absolute;
  left: 50%;
  bottom: 18px;
  z-index: 2;
  transform: translateX(-50%);
  padding: 6px 12px;
  border-radius: 999px;
  font-size: 12px;
  color: #eef2f8;
  background: rgb(15 40 45 / 72%);
  pointer-events: none;
  white-space: nowrap;
}
</style>
