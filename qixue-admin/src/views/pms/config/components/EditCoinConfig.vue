<script setup>
import { NInputNumber } from 'naive-ui'
import { ref, watch } from 'vue'
import api from '@/api'
import { SettingsEditorShell } from '@/components/index.js'
import { useSysStore } from '@/store'

defineOptions({ name: 'EditCoinConfig' })

const props = defineProps({
  embedded: { type: Boolean, default: false },
})
const emit = defineEmits(['update', 'close'])

const GRANT_DEFAULTS = {
  daily: 5,
  watch: 10,
  read: 5,
  comment: 5,
  commentDailyCap: 3,
  commentMinLen: 4,
  test: 10,
  course: 20,
}

const open = ref(false)
const loading = ref(false)
const dayCoin = ref(5)
const grant = ref({ ...GRANT_DEFAULTS })

function readAmount(data, key) {
  if (data[key] === undefined || data[key] === null || data[key] === '')
    return GRANT_DEFAULTS[key]
  return Math.max(0, Number(data[key]) || 0)
}

function show() {
  loading.value = true
  api.config.getInfo('coinConfig').then((res) => {
    const data = res.data && !Array.isArray(res.data) ? res.data : {}
    dayCoin.value = Math.max(1, Number(data.dayCoin) || 5)
    grant.value = Object.fromEntries(
      Object.keys(GRANT_DEFAULTS).map(key => [key, readAmount(data, key)]),
    )
    if (!props.embedded)
      open.value = true
  }).catch((err) => {
    console.error('加载积分配置失败:', err)
    dayCoin.value = 5
    grant.value = { ...GRANT_DEFAULTS }
    if (!props.embedded)
      open.value = true
  }).finally(() => {
    loading.value = false
  })
}

watch(() => props.embedded, (v) => { if (v) show() }, { immediate: true })

defineExpose({ show })

function confirm() {
  const value = Math.max(1, Number(dayCoin.value) || 5)
  const amounts = {}
  for (const key of Object.keys(GRANT_DEFAULTS))
    amounts[key] = Math.max(0, Number(grant.value[key]) || 0)
  loading.value = true
  api.config.saveInfo({
    style: 'qq',
    dayCoin: value,
    ...amounts,
  }, 'coinConfig').then(async () => {
    open.value = false
    try {
      await useSysStore().fetchSysInfo()
    }
    catch {}
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
    title="积分"
    width="640px"
    :loading="loading"
    @save="confirm"
    @close="emit('close')"
  >
    <div class="editor-panel-head">
      <span>发放</span>
    </div>
    <div class="grant-list">
      <label>每日登录<NInputNumber v-model:value="grant.daily" :min="0" :max="100000" /></label>
      <label>看完音视频<NInputNumber v-model:value="grant.watch" :min="0" :max="100000" /></label>
      <label>读完图文<NInputNumber v-model:value="grant.read" :min="0" :max="100000" /></label>
      <label>发表评论<NInputNumber v-model:value="grant.comment" :min="0" :max="100000" /></label>
      <label>评论每天最多<NInputNumber v-model:value="grant.commentDailyCap" :min="0" :max="100" /></label>
      <label>评论最少字数<NInputNumber v-model:value="grant.commentMinLen" :min="0" :max="500" /></label>
      <label>通过测评<NInputNumber v-model:value="grant.test" :min="0" :max="100000" /></label>
      <label>课程结业<NInputNumber v-model:value="grant.course" :min="0" :max="100000" /></label>
    </div>
    <p class="qq-note grant-note">
      填 0 则该项不发放。评论按自然日计次，短于最少字数的不计。
    </p>

    <div class="editor-panel-head">
      <span>升级规则</span>
    </div>
    <div class="qq-step">
      <span>1 个活跃日对应积分</span>
      <NInputNumber v-model:value="dayCoin" :min="1" :max="100000" />
    </div>
    <p class="qq-note">
      按 QQ 等级换算。达到 N 级需要的活跃日是 N×(N+4)，从 N 级升到下一级需要 2×N+5 个活跃日。星星、月亮、太阳、皇冠仍按 1、4、16、64 级叠在昵称旁。
    </p>
  </SettingsEditorShell>
</template>

<style scoped lang="scss">
.grant-list {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px 20px;
  margin: 8px 0 12px;
}

.grant-list label {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.grant-list :deep(.n-input-number) {
  width: 120px;
  flex-shrink: 0;
}

.grant-note {
  margin-bottom: 22px;
}

.qq-step {
  display: flex;
  align-items: center;
  gap: 12px;
  margin: 8px 0 12px;
}

.qq-note {
  margin: 0;
  color: var(--text-secondary, #5b6b7c);
  line-height: 1.6;
}
</style>
