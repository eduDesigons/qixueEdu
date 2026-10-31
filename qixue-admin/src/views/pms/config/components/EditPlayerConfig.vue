<script setup>
import { cloneDeep } from 'lodash-es'
import { NSwitch } from 'naive-ui'
import { ref, watch } from 'vue'
import api from '@/api'
import { SettingsEditorShell } from '@/components/index.js'

defineOptions({ name: 'EditPlayerConfig' })

const props = defineProps({
  embedded: { type: Boolean, default: false },
})
const emit = defineEmits(['update', 'close'])

const open = ref(false)
const loading = ref(false)
const form = ref({
  allow_seek: true,
  allow_speed: true,
})

function show() {
  loading.value = true
  api.config.getInfo('playerConfig').then((res) => {
    if (res.data) {
      form.value = cloneDeep(res.data)
    }
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
  api.config.saveInfo(form.value, 'playerConfig').then(() => {
    open.value = false
    emit('update')
    emit('close')
    $message.success('保存成功')
  }).catch(() => {
    $message.error('保存失败')
  }).finally(() => {
    loading.value = false
  })
}
</script>

<template>
  <SettingsEditorShell
    v-model:show="open"
    :embedded="embedded"
    title="播放器配置"
    :loading="loading"
    @save="confirm"
    @close="emit('close')"
  >
    <div class="settings-group">
      <div class="settings-row">
        <div class="settings-row__meta">
          <div class="settings-row__title">
            允许快进快退
          </div>
          <div class="settings-row__desc">
            开启后，用户可以在播放时快进或快退
          </div>
        </div>
        <div class="settings-row__control">
          <NSwitch v-model:value="form.allow_seek" />
        </div>
      </div>

      <div class="settings-row">
        <div class="settings-row__meta">
          <div class="settings-row__title">
            允许倍速播放
          </div>
          <div class="settings-row__desc">
            开启后，用户可以选择播放倍速（0.5x、1x、1.5x、2x等）
          </div>
        </div>
        <div class="settings-row__control">
          <NSwitch v-model:value="form.allow_speed" />
        </div>
      </div>
    </div>
  </SettingsEditorShell>
</template>
