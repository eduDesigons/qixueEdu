<script setup>
import { NDatePicker } from 'naive-ui'
import { ref, watch } from 'vue'
import api from '@/api'
import { SettingsEditorShell } from '@/components/index.js'

defineOptions({ name: 'EditMemberRule' })

const props = defineProps({
  embedded: { type: Boolean, default: false },
})
const emit = defineEmits(['update', 'close'])

const open = ref(false)
const loading = ref(false)
const range = ref(null)

function toStamp(value) {
  if (!value)
    return null
  const ts = Date.parse(String(value).replace(/-/g, '/'))
  return Number.isNaN(ts) ? null : ts
}

function formatStamp(ts) {
  const date = new Date(ts)
  const pad = n => String(n).padStart(2, '0')
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())} ${pad(date.getHours())}:${pad(date.getMinutes())}:${pad(date.getSeconds())}`
}

function show() {
  loading.value = true
  api.config.getInfo('memberRule').then((res) => {
    const data = res.data && typeof res.data === 'object' ? res.data : {}
    const start = toStamp(data.register_start)
    const end = toStamp(data.register_end)
    range.value = start && end ? [start, end] : null
    if (!props.embedded)
      open.value = true
  }).catch(() => {
    range.value = null
    if (!props.embedded)
      open.value = true
  }).finally(() => {
    loading.value = false
  })
}

watch(() => props.embedded, (v) => { if (v) show() }, { immediate: true })
defineExpose({ show })

function confirm() {
  let registerStart = ''
  let registerEnd = ''
  if (Array.isArray(range.value) && range.value.length === 2 && range.value[0] && range.value[1]) {
    if (range.value[0] > range.value[1]) {
      $message.warning('开始时间不能晚于结束时间')
      return
    }
    registerStart = formatStamp(range.value[0])
    registerEnd = formatStamp(range.value[1])
  }
  loading.value = true
  api.config.saveInfo({
    register_start: registerStart,
    register_end: registerEnd,
  }, 'memberRule').then(() => {
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
    title="学员规则"
    :loading="loading"
    @save="confirm"
    @close="emit('close')"
  >
    <div class="rule-range">
      <div class="rule-range__title">
        注册时间范围
      </div>
      <div class="rule-range__desc">
        只限制后台添加学员。留空表示不限制。
      </div>
      <NDatePicker
        v-model:value="range"
        type="datetimerange"
        clearable
        style="width: 100%"
        start-placeholder="开始时间"
        end-placeholder="结束时间"
      />
    </div>
  </SettingsEditorShell>
</template>

<style scoped>
.rule-range {
  margin-top: 4px;
}

.rule-range__title {
  font-size: 15px;
  font-weight: 600;
  color: var(--ctl-text, #1c2b2e);
}

.rule-range__desc {
  margin: 4px 0 10px;
  color: var(--text-secondary, #5b6b7c);
  line-height: 1.6;
}
</style>
