<script setup>
import { NAlert, NButton, NDataTable, NEmpty, NSelect, NSpin, NTag } from 'naive-ui'
import { computed, h, ref } from 'vue'
import api from '@/api'
import { CommonPage } from '@/components'
import { useOpsUserNav } from './opsCells.js'

defineOptions({ name: 'OpsInspect' })

const { userCell } = useOpsUserNav()
const kind = ref(null)
const scanning = ref(false)
const repairing = ref(false)
const scanned = ref(false)
const truncated = ref(false)
const summary = ref([])
const issues = ref([])
const lastRepair = ref(null)
const steps = ref([])
const stepText = ref('')
const busy = computed(() => scanning.value || repairing.value)

const kindOptions = computed(() => {
  const options = [{ label: '全部类型', value: null }]
  for (const item of summary.value) {
    options.push({ label: `${item.label}（${item.count}）`, value: item.kind })
  }
  return options
})

const filtered = computed(() => {
  if (!kind.value)
    return issues.value
  return issues.value.filter(row => row.kind === kind.value)
})

const total = computed(() => issues.value.length)

const pagination = {
  pageSize: 20,
  showSizePicker: true,
  pageSizes: [20, 50, 100],
  prefix: ({ itemCount }) => `共 ${itemCount} 条`,
}

function kindTag(row) {
  const map = {
    course_missing_chapters: 'warning',
    course_orphan_chapters: 'error',
    course_id_mismatch: 'warning',
    course_missing_enroll: 'warning',
    course_progress_mismatch: 'warning',
    course_ended_incomplete: 'error',
    course_missing_credit: 'warning',
    course_missing_coins: 'warning',
    course_missing_cert: 'warning',
    chapter_flag_mismatch: 'warning',
    path_missing_courses: 'warning',
    path_progress_mismatch: 'warning',
    path_ended_incomplete: 'error',
    path_should_end: 'warning',
    path_missing_cert: 'warning',
  }
  return h(
    NTag,
    { size: 'small', round: true, bordered: false, type: map[row.kind] || 'default' },
    { default: () => row.label || row.kind },
  )
}

function targetText(row) {
  const parts = []
  if (row.title)
    parts.push(row.title)
  if (row.course_id && !row.chapter_id)
    parts.push(`课程 #${row.course_id}`)
  if (row.path_id)
    parts.push(`系列 #${row.path_id}`)
  if (row.chapter_id)
    parts.push(`课节 #${row.chapter_id}`)
  return parts.filter(Boolean).join(' · ') || '-'
}

const columns = [
  { title: '学员', key: 'user', width: 180, render: row => userCell(row) },
  { title: '异常', key: 'kind', width: 150, render: kindTag },
  {
    title: '对象',
    key: 'title',
    ellipsis: { tooltip: true },
    render: row => targetText(row),
  },
  {
    title: '说明',
    key: 'detail',
    ellipsis: { tooltip: true },
    render: row => row.detail || '-',
  },
]

function rebuildSummary() {
  const counts = {}
  const labels = {}
  for (const row of issues.value) {
    counts[row.kind] = (counts[row.kind] || 0) + 1
    labels[row.kind] = row.label || labels[row.kind]
  }
  const order = steps.value.length ? steps.value : Object.keys(counts).map(code => ({ kind: code, label: labels[code] || code }))
  summary.value = order
    .filter(step => counts[step.kind])
    .map(step => ({
      kind: step.kind,
      label: labels[step.kind] || step.label || step.kind,
      count: counts[step.kind],
    }))
  truncated.value = Object.values(counts).some(n => n >= 300)
}

function replaceKindIssues(kindCode, rows, isTruncated) {
  issues.value = issues.value.filter(row => row.kind !== kindCode).concat(rows || [])
  rebuildSummary()
  if (isTruncated)
    truncated.value = true
}

async function loadSteps() {
  if (steps.value.length)
    return steps.value
  const res = await api.ops.inspectKinds()
  steps.value = res.data || []
  return steps.value
}

async function scan() {
  if (busy.value)
    return
  scanning.value = true
  lastRepair.value = null
  truncated.value = false
  issues.value = []
  summary.value = []
  kind.value = null
  try {
    const list = await loadSteps()
    for (let i = 0; i < list.length; i++) {
      const step = list[i]
      stepText.value = `正在检测：${step.label}（${i + 1}/${list.length}）`
      const res = await api.ops.inspect({ kind: step.kind })
      const data = res.data || {}
      replaceKindIssues(step.kind, data.issues || [], !!data.truncated)
    }
    scanned.value = true
    window.$message?.success(issues.value.length ? `检出 ${issues.value.length} 条异常` : '未发现异常')
  }
  catch {
    scanned.value = false
    issues.value = []
    summary.value = []
  }
  finally {
    stepText.value = ''
    scanning.value = false
  }
}

async function repairKind(step) {
  let rounds = 0
  let repaired = 0
  let failed = 0
  let remaining = 0
  while (rounds < 80) {
    rounds++
    stepText.value = `正在修复：${step.label}`
    const res = await api.ops.repair({ kind: step.kind })
    const data = res.data || {}
    repaired += Number(data.repaired) || 0
    failed += Number(data.failed) || 0
    remaining = Number(data.remaining) || 0
    if (!remaining)
      break
    if (!(Number(data.repaired) || Number(data.failed)))
      break
  }
  const again = await api.ops.inspect({ kind: step.kind })
  const next = again.data || {}
  replaceKindIssues(step.kind, next.issues || [], !!next.truncated)
  return { repaired, failed, remaining: Number(next.total) || remaining }
}

function repair() {
  if (!issues.value.length || busy.value)
    return
  window.$dialog?.warning({
    title: '开始修复',
    content: `将按类型分步修复当前 ${issues.value.length} 条异常，每类单独请求，避免超时。`,
    positiveText: '开始修复',
    negativeText: '取消',
    onPositiveClick: async () => {
      repairing.value = true
      try {
        const list = await loadSteps()
        const todo = list.filter(step => issues.value.some(row => row.kind === step.kind))
        let repaired = 0
        let failed = 0
        for (const step of todo) {
          const result = await repairKind(step)
          repaired += result.repaired
          failed += result.failed
        }
        lastRepair.value = { repaired, failed, remaining: issues.value.length }
        scanned.value = true
        if (issues.value.length)
          window.$message?.warning(`已处理 ${repaired} 项，仍有 ${issues.value.length} 条待处理`)
        else
          window.$message?.success(`修复完成，已处理 ${repaired} 项`)
      }
      finally {
        stepText.value = ''
        repairing.value = false
      }
    },
  })
}
</script>

<template>
  <CommonPage title="数据巡检" subtitle="核对入课课节、结业学分/积分/证书，以及系列完成数据，缺失的可一键补齐。">
    <template #action>
      <NButton :loading="scanning" :disabled="repairing" type="primary" @click="scan">
        开始检测
      </NButton>
      <NButton
        :loading="repairing"
        :disabled="scanning || !issues.length"
        secondary
        type="warning"
        @click="repair"
      >
        开始修复
      </NButton>
    </template>

    <NSpin :show="busy" class="inspect-spin">
      <template v-if="stepText" #description>{{ stepText }}</template>
      <div class="inspect">
        <NAlert v-if="truncated" type="warning" class="inspect-alert">
          同类异常较多，列表每类最多展示 300 条。修复后可再检测一次。
        </NAlert>
        <NAlert v-if="lastRepair" type="info" class="inspect-alert">
          上次修复处理 {{ lastRepair.repaired || 0 }} 项
          <template v-if="lastRepair.failed">，失败 {{ lastRepair.failed }} 项</template>
          ，剩余 {{ lastRepair.remaining || 0 }} 条。
        </NAlert>

        <div v-if="summary.length" class="inspect-stats">
          <article class="inspect-stat">
            <span>异常合计</span>
            <strong>{{ total }}</strong>
          </article>
          <article v-for="item in summary" :key="item.kind" class="inspect-stat">
            <span>{{ item.label }}</span>
            <strong>{{ item.count }}</strong>
          </article>
        </div>

        <div v-if="scanned && issues.length" class="inspect-toolbar">
          <NSelect
            v-model:value="kind"
            :options="kindOptions"
            class="w-220"
            placeholder="异常类型"
          />
        </div>

        <div v-if="scanned && !issues.length" class="inspect-empty">
          <NEmpty description="未发现课程、课节或系列数据异常" />
        </div>
        <div v-else-if="!scanned" class="inspect-empty">
          <NEmpty description="点击「开始检测」扫描学员学习数据" />
        </div>
        <NDataTable
          v-else
          :columns="columns"
          :data="filtered"
          :row-key="row => row.key"
          :pagination="pagination"
          :scroll-x="900"
          :bordered="false"
          :bottom-bordered="false"
        />
      </div>
    </NSpin>
  </CommonPage>
</template>

<style scoped>
.inspect-spin {
  min-height: 360px;
}

.inspect-spin :deep(.n-spin-content) {
  min-height: 360px;
}

.inspect {
  display: flex;
  flex-direction: column;
  gap: 12px;
  min-height: 360px;
}

.inspect-alert {
  border-radius: 12px;
}

.inspect-stats {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(140px, 1fr));
  gap: 10px;
}

.inspect-stat {
  border: 1px solid var(--ctl-hairline, rgba(32, 42, 58, 0.08));
  border-radius: 16px;
  background: var(--ctl-glass-strong, #fff);
  padding: 12px 14px;
}

.inspect-stat span {
  color: var(--ctl-text-soft, #6b7280);
  font-size: 12px;
}

.inspect-stat strong {
  display: block;
  margin-top: 4px;
  color: var(--ctl-text, #1a2332);
  font-size: 22px;
  font-weight: 650;
  letter-spacing: -0.03em;
}

.inspect-toolbar {
  display: flex;
  justify-content: flex-end;
}

.inspect-empty {
  flex: 1;
  min-height: 280px;
  display: flex;
  align-items: center;
  justify-content: center;
}
</style>
