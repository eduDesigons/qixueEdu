<script setup>
import { NFormItem, NInput } from 'naive-ui'
import { ref, watch } from 'vue'
import api from '@/api'
import { SettingsEditorShell } from '@/components/index.js'
import WangEditor from '@/components/common/WangEditor.vue'

defineOptions({ name: 'EditEnrollConfig' })

const props = defineProps({
  embedded: { type: Boolean, default: false },
})
const emit = defineEmits(['update', 'close'])

const open = ref(false)
const loading = ref(false)
const form = ref({
  title: '入学要求',
  content: '',
})

function show() {
  loading.value = true
  api.config.getInfo('enrollConfig').then((res) => {
    const data = res.data && typeof res.data === 'object' ? res.data : {}
    form.value = {
      title: data.title || '入学要求',
      content: typeof data.content === 'string' ? data.content : '',
    }
    if (!props.embedded)
      open.value = true
  }).catch(() => {
    form.value = { title: '入学要求', content: '' }
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
  api.config.saveInfo(form.value, 'enrollConfig').then(() => {
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
    title="入学要求"
    width="900px"
    :loading="loading"
    @save="confirm"
    @close="emit('close')"
  >
    <NFormItem label="标题">
      <NInput v-model:value="form.title" placeholder="入学要求" />
    </NFormItem>
    <WangEditor v-model:value="form.content" />
  </SettingsEditorShell>
</template>
