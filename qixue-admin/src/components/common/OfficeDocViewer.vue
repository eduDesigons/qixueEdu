<script setup>
import { computed, defineAsyncComponent, ref, watch } from 'vue'
import PdfLazyViewer from './PdfLazyViewer.vue'
import PptxLazyViewer from './PptxLazyViewer.vue'
import '@vue-office/docx/lib/index.css'
import '@vue-office/excel/lib/index.css'

const VueOfficeDocx = defineAsyncComponent(() => import('@vue-office/docx'))
const VueOfficeExcel = defineAsyncComponent(() => import('@vue-office/excel'))

const props = defineProps({
  src: { type: String, default: '' },
  filename: { type: String, default: '' },
})

const emit = defineEmits(['rendered', 'error', 'loading'])

const loading = ref(false)
const error = ref('')
const loadHint = ref('加载中')

const fileExt = computed(() => {
  const fromName = String(props.filename || '')
  const fromUrl = String(props.src || '')
  const source = fromName.includes('.') ? fromName : fromUrl
  const match = source.split('?')[0].split('#')[0].split('.').pop()
  return match ? match.toLowerCase() : ''
})

const kind = computed(() => {
  const ext = fileExt.value
  if (ext === 'pdf')
    return 'pdf'
  if (['doc', 'docx'].includes(ext))
    return 'docx'
  if (['xls', 'xlsx', 'csv'].includes(ext))
    return 'excel'
  if (['ppt', 'pptx'].includes(ext))
    return 'pptx'
  return ''
})

const supported = computed(() => Boolean(props.src && kind.value))

watch(
  () => [props.src, kind.value],
  () => {
    error.value = ''
    loadHint.value = kind.value === 'pdf' ? '渲染首屏…' : '加载中'
    if (supported.value) {
      loading.value = true
      emit('loading', true)
    }
    else {
      loading.value = false
      emit('loading', false)
    }
  },
  { immediate: true },
)

function onRendered() {
  loading.value = false
  error.value = ''
  loadHint.value = '加载中'
  emit('loading', false)
  emit('rendered')
}

function onError(err) {
  loading.value = false
  error.value = '文档预览失败，请下载后查看'
  loadHint.value = '加载中'
  emit('loading', false)
  emit('error', err)
}

function onViewerLoading(v) {
  loading.value = v
  emit('loading', v)
}

function onProgress(pct) {
  if (kind.value !== 'pptx')
    return
  if (pct > 0 && pct < 100)
    loadHint.value = `下载中 ${pct}%`
  else if (pct >= 100)
    loadHint.value = '渲染首屏…'
}
</script>

<template>
  <div class="office-doc">
    <div class="office-doc__body">
      <div v-if="loading" class="office-doc__state">
        <span class="office-doc__spinner" />
        <p>{{ loadHint }}</p>
      </div>
      <div v-else-if="error || !supported" class="office-doc__state">
        <i class="i-fe:alert-circle" />
        <p>{{ error || '暂不支持该格式在线预览' }}</p>
        <a v-if="src" :href="src" target="_blank" rel="noopener" referrerpolicy="origin">下载文件</a>
      </div>

      <PdfLazyViewer
        v-if="kind === 'pdf' && src"
        :src="src"
        class="office-doc__viewer"
        @loading="onViewerLoading"
        @rendered="onRendered"
        @error="onError"
      />
      <VueOfficeDocx
        v-else-if="kind === 'docx' && src"
        :src="src"
        class="office-doc__viewer"
        @rendered="onRendered"
        @error="onError"
      />
      <VueOfficeExcel
        v-else-if="kind === 'excel' && src"
        :src="src"
        class="office-doc__viewer"
        @rendered="onRendered"
        @error="onError"
      />
      <PptxLazyViewer
        v-else-if="kind === 'pptx' && src"
        :src="src"
        class="office-doc__viewer"
        @loading="onViewerLoading"
        @progress="onProgress"
        @rendered="onRendered"
        @error="onError"
      />
    </div>
  </div>
</template>

<style scoped lang="scss">
.office-doc {
  position: relative;
  display: flex;
  flex-direction: column;
  width: 100%;
  height: 100%;
  min-height: 0;
  overflow: hidden;
  background: #f4f7f7;
}

.office-doc__body {
  position: relative;
  flex: 1;
  min-height: 0;
  overflow: hidden;
  background:
    radial-gradient(90% 60% at 50% 0%, rgba(var(--primary-color), 0.05), transparent 70%),
    linear-gradient(180deg, #f7f9fc 0%, #eef2f8 100%);
}

.office-doc__viewer {
  width: 100%;
  height: 100%;
}

.office-doc__viewer :deep(.docx-wrapper),
.office-doc__viewer :deep(.vue-office-pdf),
.office-doc__viewer :deep(.vue-office-excel),
.office-doc__viewer :deep(.vue-office-pptx) {
  height: 100%;
}

.office-doc__viewer :deep(.docx-wrapper) {
  background: transparent;
  padding: 18px 12px 28px;
}

.office-doc__viewer :deep(.docx-wrapper > section.docx) {
  box-shadow:
    0 1px 0 rgb(255 255 255 / 80%),
    0 12px 32px rgb(15 40 45 / 08%);
  border-radius: 2px;
}

.office-doc__state {
  position: absolute;
  inset: 0;
  z-index: 2;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 10px;
  color: #5b6b7c;
  background: rgb(247 250 251 / 92%);
  backdrop-filter: blur(6px);
}

.office-doc__state a {
  color: rgb(var(--primary-color));
  font-weight: 600;
  text-decoration: none;
}

.office-doc__spinner {
  width: 22px;
  height: 22px;
  border: 2px solid rgba(var(--primary-color), 0.16);
  border-top-color: rgb(var(--primary-color));
  border-radius: 50%;
  animation: office-spin 0.75s linear infinite;
}

@keyframes office-spin {
  to { transform: rotate(360deg); }
}
</style>
