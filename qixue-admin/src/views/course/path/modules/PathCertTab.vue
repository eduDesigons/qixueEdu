<script setup>
/**
 * 系列详情「证书颁发」页签 —— 复用课程侧组件逻辑（pathId）
 */
import CourseCertTab from '@/views/course/list/modules/CourseCertTab.vue'

defineOptions({ name: 'PathCertTab' })

defineProps({
  pathId: {
    type: [String, Number],
    required: true,
  },
  pendingCount: {
    type: [Number, String],
    default: 0,
  },
})

const emit = defineEmits(['change'])
const $inner = ref(null)

defineExpose({
  refresh: () => $inner.value?.refresh?.(),
  askBackfill: () => $inner.value?.askBackfill?.(),
})
</script>

<template>
  <CourseCertTab
    ref="$inner"
    :path-id="pathId"
    :pending-count="pendingCount"
    @change="emit('change')"
  />
</template>
