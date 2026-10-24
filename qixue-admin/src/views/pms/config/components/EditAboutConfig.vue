<script setup>
import { ref, watch } from 'vue'
import api from '@/api'
import { SettingsEditorShell } from '@/components/index.js'
import WangEditor from '@/components/common/WangEditor.vue'

defineOptions({ name: 'EditAboutConfig' })

const props = defineProps({
  embedded: { type: Boolean, default: false },
})
const emit = defineEmits(['update', 'close'])

const open = ref(false)
const loading = ref(false)
const form = ref({
  content: '',
})

function show() {
  loading.value = true
  api.config.getInfo('aboutConfig').then((res) => {
    // 处理不同的数据结构
    if (res.data) {
      if (typeof res.data === 'object' && res.data !== null && typeof res.data.content === 'string') {
        form.value.content = res.data.content
      }
      else if (typeof res.data === 'string') {

        form.value.content = res.data
      }
      else if (typeof res.data.data === 'object' && res.data.data !== null && typeof res.data.data.content === 'string') {
        form.value.content = res.data.data.content
      }
      else {
        form.value.content = ''
      }
    }
    else {
      form.value.content = ''
    }
    if (!props.embedded)
      open.value = true
  }).catch((err) => {
    console.error('加载关于我们配置失败:', err)
    form.value.content = ''
    if (!props.embedded)
      open.value = true
  }).finally(() => {
    loading.value = false
  })
}

watch(() => props.embedded, (v) => { if (v) show() }, { immediate: true })

defineExpose({ show })

function confirm() {
  loading.value = true
  // 保存为对象格式 {content: ''}
  api.config.saveInfo(form.value, 'aboutConfig').then(() => {
    open.value = false
    emit('update')
    emit('close')
    window.$message.success('保存成功')
  }).catch(() => {
    window.$message.error('保存失败')
  }).finally(() => {
    loading.value = false
  })
}
</script>

<template>
  <SettingsEditorShell
    v-model:show="open"
    :embedded="embedded"
    title="关于我们配置"
    width="900px"
    :loading="loading"
    @save="confirm"
    @close="emit('close')"
  >
      <WangEditor v-model:value="form.content" />
  </SettingsEditorShell>
</template>

