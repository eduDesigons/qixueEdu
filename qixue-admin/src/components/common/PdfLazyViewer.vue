<script setup>
/**
 * PDF 预览：沿用 @vue-office/pdf（pdf.js + 虚拟列表），不必换插件。
 * 优化点：按容器宽度渲染；disableAutoFetch + 更大 rangeChunk，首屏按 Range 拉取。
 */
import { computed, defineAsyncComponent, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'

const VueOfficePdf = defineAsyncComponent(() => import('@vue-office/pdf/lib/v3/vue-office-pdf.mjs'))

const props = defineProps({
  src: { type: String, default: '' },
})

const emit = defineEmits(['rendered', 'error', 'loading'])

const bodyRef = ref(null)
const ready = ref(false)
const mountWidth = ref(0)

let ro = null

const pdfOptions = computed(() => ({
  width: Math.max(320, mountWidth.value - 16),
  gap: 12,
  disableAutoFetch: true,
  disableStream: false,
  rangeChunkSize: 262144,
}))

function measure() {
  const el = bodyRef.value
  if (!el)
    return
  const w = Math.floor(el.getBoundingClientRect().width)
  if (w <= 40)
    return
  if (!ready.value) {
    mountWidth.value = w
    ready.value = true
    return
  }
  // 全屏切换等大幅变宽时再重挂，避免滚动时频繁重建
  if (Math.abs(w - mountWidth.value) >= 48)
    mountWidth.value = w
}

watch(
  () => props.src,
  () => {
    if (props.src) {
      emit('loading', true)
      ready.value = false
      mountWidth.value = 0
      nextTick(() => {
        measure()
        if (!ready.value)
          requestAnimationFrame(measure)
      })
    }
  },
  { immediate: true },
)

onMounted(() => {
  const el = bodyRef.value
  if (!el || typeof ResizeObserver === 'undefined')
    return
  ro = new ResizeObserver(() => measure())
  ro.observe(el)
  measure()
})

onBeforeUnmount(() => {
  ro?.disconnect()
  ro = null
})

function onRendered() {
  emit('loading', false)
  emit('rendered')
}

function onError(err) {
  emit('loading', false)
  emit('error', err)
}
</script>

<template>
  <div ref="bodyRef" class="pdf-lazy">
    <VueOfficePdf
      v-if="src && ready && mountWidth"
      :key="`${src}:${mountWidth}`"
      :src="src"
      class="pdf-lazy__viewer"
      :options="pdfOptions"
      @rendered="onRendered"
      @error="onError"
    />
  </div>
</template>

<style scoped lang="scss">
.pdf-lazy {
  width: 100%;
  height: 100%;
  min-height: 0;
  overflow: hidden;
}

.pdf-lazy__viewer {
  width: 100%;
  height: 100%;
}

.pdf-lazy__viewer :deep(.vue-office-pdf) {
  height: 100%;
}
</style>
