<script setup>
import { computed, defineComponent } from 'vue'
import { NCheckbox } from 'naive-ui'
import { useRouter } from 'vue-router'
import { getAssetUrl } from '@/utils'

const props = defineProps({
  list: {
    type: Array,
    default: () => [],
  },
  columns: {
    type: Array,
    default: () => [],
  },
  rowKey: {
    type: String,
    default: 'id',
  },
  checkedKeys: {
    type: Array,
    default: () => [],
  },
})

const emit = defineEmits(['update:checkedKeys'])
const router = useRouter()

const SKIP_KEYS = new Set([
  'id',
  'actions',
  'action',
  'thumb',
  'img_url',
  'cover',
  'image',
  'avatar',
  'course-info',
  'user-info',
  'user',
  'desc',
  'content',
])

const METRIC_KEYS = new Set(['credit', 'coins', 'sum', 'sort', 'score', 'amount', 'price'])
const CHIP_KEYS = new Set(['status', 'type', 'field', 'is_show', 'is_lock', 'is_active', 'is_rec', 'enable', 'category', 'area', 'is_pass', 'is_people_check', 'is_read', 'paper_id'])
const TIME_KEYS = new Set(['created_at', 'updated_at', 'publish_time', 'last_login_at'])
const LINK_KEYS = new Set(['course', 'course_name', 'chapter', 'chapter_name', 'path', 'path_name'])

const CellView = defineComponent({
  name: 'MeCrudCell',
  props: {
    vnode: { default: null },
  },
  setup(cellProps) {
    return () => cellProps.vnode ?? null
  },
})

const hasSelection = computed(() => props.columns.some(col => col.type === 'selection'))
const actionCol = computed(() => props.columns.find(col => isActionCol(col)))
const titleKey = computed(() => guessTitleKey(props.list[0]))

const metricCols = computed(() => props.columns.filter(col => isMetricCol(col)).slice(0, 4))
const chipCols = computed(() => props.columns.filter(col => isChipCol(col)).slice(0, 2))
const timeCol = computed(() => props.columns.find(col => TIME_KEYS.has(col.key)))
const fieldCols = computed(() => {
  const cols = props.columns.filter((col) => {
    if (!isVisibleCol(col) || isMetricCol(col) || isChipCol(col) || TIME_KEYS.has(col.key))
      return false
    if (col.key === titleKey.value)
      return false
    if (identityCol.value && col.key === identityCol.value.key)
      return false
    if (isLinkCol(col))
      return false
    return true
  })
  cols.sort((a, b) => {
    if (a.key === 'remark')
      return -1
    if (b.key === 'remark')
      return 1
    return 0
  })
  return cols.slice(0, 3)
})

const identityCol = computed(() => {
  const keys = new Set(['course-info', 'user-info', 'user'])
  return props.columns.find(col => keys.has(col.key))
    || props.columns.find(col => typeof col.render === 'function' && (col.width || 0) >= 240 && !isActionCol(col) && col.type !== 'selection')
})

const linkCols = computed(() => props.columns.filter((col) => {
  if (!col?.title || col.type === 'selection' || isActionCol(col))
    return false
  if (identityCol.value && col.key === identityCol.value.key)
    return false
  if (!isLinkCol(col) || typeof col.render !== 'function')
    return false
  return true
}).slice(0, 3))

const pageSelectedCount = computed(() => props.list.filter(row => isChecked(row)).length)
const pageChecked = computed(() => props.list.length > 0 && props.list.every(row => isChecked(row)))
const pageIndeterminate = computed(() => {
  const n = pageSelectedCount.value
  return n > 0 && n < props.list.length
})

function isActionCol(col) {
  return col.key === 'actions' || col.key === 'action' || col.title === '操作'
}

function isVisibleCol(col) {
  if (!col?.title || col.type === 'selection' || isActionCol(col))
    return false
  const key = col.key
  if (!key || SKIP_KEYS.has(key))
    return false
  return true
}

function isMetricCol(col) {
  if (!isVisibleCol(col))
    return false
  const key = col.key || ''
  if (METRIC_KEYS.has(key) || /_count$/.test(key))
    return true
  return false
}

function isChipCol(col) {
  if (!isVisibleCol(col))
    return false
  const key = col.key || ''
  if (CHIP_KEYS.has(key))
    return true
  return typeof col.title === 'string' && /状态|类型|分类|区域|通过|批改|阅读/.test(col.title)
}

function isLinkCol(col) {
  const key = col?.key || ''
  // 标题列本身不当成关联链接，避免卡片里名称重复
  if (key === 'title' || key === 'name' || key === 'nick_name')
    return false
  if (LINK_KEYS.has(key))
    return true
  return typeof col?.title === 'string' && /^(课程|课节|系列)$/.test(col.title)
}

function guessTitleKey(row) {
  if (!row)
    return 'title'
  if (row.title)
    return 'title'
  if (row.name)
    return 'name'
  if (row.nick_name || row.user?.nick_name)
    return 'nick_name'
  if (row.field)
    return 'field'
  if (row.remark)
    return 'remark'
  if (row.account)
    return 'account'
  if (row.content)
    return 'content'
  return 'title'
}

function titleOf(row) {
  const raw = row.title
    || row.name
    || row.nick_name
    || row.user?.nick_name
    || row.paper?.course?.title
    || row.field
    || row.remark
    || row.account
    || row.content
  if (typeof raw === 'string' && raw.trim())
    return raw.trim()
  return '未命名'
}

function rowId(row) {
  return row?.[props.rowKey]
}

function isChecked(row) {
  return props.checkedKeys.includes(rowId(row))
}

function setCheckedKeys(next) {
  emit('update:checkedKeys', next)
}

function toggleOne(row, checked) {
  const id = rowId(row)
  if (checked) {
    if (props.checkedKeys.includes(id))
      return
    setCheckedKeys([...props.checkedKeys, id])
    return
  }
  setCheckedKeys(props.checkedKeys.filter(item => item !== id))
}

const pageCheck = computed({
  get: () => pageChecked.value,
  set: (checked) => {
    const ids = props.list.map(rowId)
    if (checked)
      setCheckedKeys([...new Set([...props.checkedKeys, ...ids])])
    else
      setCheckedKeys(props.checkedKeys.filter(id => !ids.includes(id)))
  },
})

function pickUrl(...vals) {
  for (const val of vals) {
    if (typeof val === 'string' && val.trim())
      return getAssetUrl(val)
  }
  return ''
}

function coverOf(row) {
  return pickUrl(
    row.thumb,
    row.img_url,
    row.cover,
    row.image,
    row.course?.thumb,
    row.course?.cover,
    row.path?.thumb,
    row.path?.cover,
  )
}

function avatarOf(row) {
  return pickUrl(row.avatar, row.user?.avatar)
}

/** 入课/入系列申请、答卷记录等：不是课程本体，背景不用课程封面 */
function isUserSideCard(row) {
  if (!row || !(row.user_id || row.user))
    return false
  if (Object.prototype.hasOwnProperty.call(row, 'reason'))
    return true
  if (row.paper_id != null || row.paper)
    return true
  return false
}

function photoOf(row) {
  if (isUserSideCard(row))
    return avatarOf(row)
  return coverOf(row) || avatarOf(row)
}

function subtitleOf(row) {
  const title = titleOf(row)
  return [
    row.account,
    row.mobile,
    row.email,
    row.code,
    row.desc,
    row.user?.account,
    row.user?.mobile,
  ]
    .filter(val => typeof val === 'string' && val.trim() && val.trim() !== title)
    .slice(0, 2)
    .join(' · ')
}

function clickOf(col, row) {
  if (!col || typeof col.render !== 'function')
    return null
  try {
    const vnode = col.render(row)
    const click = vnode?.props?.onClick
    return typeof click === 'function' ? click : null
  }
  catch {
    return null
  }
}

function vnodeText(node) {
  if (node == null || node === false)
    return ''
  if (typeof node === 'string' || typeof node === 'number')
    return String(node)
  if (Array.isArray(node))
    return node.map(vnodeText).find(Boolean) || ''
  if (typeof node === 'object') {
    const kids = node.children
    if (typeof kids === 'string' || typeof kids === 'number')
      return String(kids)
    if (Array.isArray(kids))
      return vnodeText(kids)
    if (typeof kids?.default === 'function')
      return vnodeText(kids.default())
  }
  return ''
}

function labelOf(col, row) {
  const raw = typeof col.render === 'function' ? vnodeText(col.render(row)) : renderCell(col, row)
  const text = String(raw || '').trim()
  if (!text || text === '-' || text === '—')
    return ''
  return text
}

function rowOpener(row) {
  return clickOf(identityCol.value, row)
}

function titleOpener(row) {
  const titleCol = props.columns.find(col => col.key === titleKey.value || col.key === 'title' || col.key === 'name')
  return clickOf(titleCol, row) || rowOpener(row)
}

function isOpenable(row) {
  return !!titleOpener(row)
}

function openRow(row) {
  titleOpener(row)?.()
}

function openTitle(row, e) {
  e.stopPropagation()
  titleOpener(row)?.()
}

function openLink(col, row, e) {
  e.stopPropagation()
  clickOf(col, row)?.()
}

function personOf(row) {
  const name = row.user?.nick_name
  const id = row.user_id || row.user?.id
  if (!name || id == null)
    return null
  if (titleOf(row) === String(name).trim())
    return null
  return { name, id }
}

function openPerson(row, e) {
  e.stopPropagation()
  const col = identityCol.value
  if (col && (col.key === 'user' || col.key === 'user-info' || col.key === 'course-info')) {
    clickOf(col, row)?.()
    return
  }
  const id = row.user_id || row.user?.id
  if (id)
    router.push({ name: 'UserDetail', params: { id } })
}

function renderCell(column, row) {
  if (typeof column.render === 'function')
    return column.render(row)
  const val = row?.[column.key]
  if (val == null || val === '')
    return '—'
  if (typeof val === 'object')
    return val.name || val.title || val.nick_name || '—'
  return String(val)
}

function isTodo(row) {
  return Number(row.is_people_check) === 1
    || Number(row.status) === 0
    || row.is_read === 0
}
</script>

<template>
  <div class="gen-wrap">
    <div v-if="hasSelection" class="entity-select-bar">
      <NCheckbox
        v-model:checked="pageCheck"
        :indeterminate="pageIndeterminate"
      >
        本页全选
      </NCheckbox>
      <span v-if="pageSelectedCount" class="entity-select-bar__count">已选 {{ pageSelectedCount }} 条</span>
    </div>

    <div class="entity-grid">
      <article
        v-for="row in list"
        :key="rowId(row)"
        class="entity-card gen-card"
        :class="{
          'is-checked': isChecked(row),
          'has-check': hasSelection,
          'has-photo': !!photoOf(row),
          'is-openable': isOpenable(row),
          'is-todo': isTodo(row),
        }"
        @click="openRow(row)"
      >
        <img
          v-if="photoOf(row)"
          class="photo-bg"
          :src="photoOf(row)"
          alt=""
        >
        <div v-if="photoOf(row)" class="photo-veil" />
        <div v-if="hasSelection" class="gen-card__check" @click.stop>
          <NCheckbox
            :checked="isChecked(row)"
            @update:checked="toggleOne(row, $event)"
          />
        </div>

        <div v-if="chipCols.length" class="gen-card__chips" @click.stop>
          <CellView
            v-for="(col, i) in chipCols"
            :key="`${col.key}-${i}`"
            :vnode="renderCell(col, row)"
          />
        </div>

        <button
          v-if="isOpenable(row)"
          type="button"
          class="entity-card__title is-link"
          :title="titleOf(row)"
          @click="openTitle(row, $event)"
        >
          {{ titleOf(row) }}
        </button>
        <h3 v-else class="entity-card__title" :title="titleOf(row)">
          {{ titleOf(row) }}
        </h3>
        <p v-if="subtitleOf(row)" class="entity-card__meta">
          {{ subtitleOf(row) }}
        </p>

        <div v-if="personOf(row) || linkCols.length" class="entity-card__rels" @click.stop>
          <button
            v-if="personOf(row)"
            type="button"
            class="entity-card__chip"
            @click="openPerson(row, $event)"
          >
            {{ personOf(row).name }}
          </button>
          <button
            v-for="(col, i) in linkCols"
            :key="`${col.key}-${i}`"
            v-show="labelOf(col, row)"
            type="button"
            class="entity-card__chip"
            :title="labelOf(col, row)"
            @click="openLink(col, row, $event)"
          >
            {{ labelOf(col, row) }}
          </button>
        </div>

        <div v-if="metricCols.length" class="gen-card__stats">
          <div v-for="(col, i) in metricCols" :key="`${col.key}-${i}`" class="gen-card__stat">
            <span class="gen-card__stat-val">
              <CellView :vnode="renderCell(col, row)" />
            </span>
            <span class="gen-card__stat-lab">{{ col.title }}</span>
          </div>
        </div>

        <dl v-if="fieldCols.length" class="gen-card__fields" @click.stop>
          <div v-for="(col, i) in fieldCols" :key="`${col.key}-${i}`" class="gen-card__row">
            <dt>{{ col.title }}</dt>
            <dd>
              <CellView :vnode="renderCell(col, row)" />
            </dd>
          </div>
        </dl>

        <p v-if="timeCol" class="gen-card__time">
          <CellView :vnode="renderCell(timeCol, row)" />
        </p>

        <div v-if="actionCol" class="gen-card__actions" @click.stop>
          <CellView :vnode="renderCell(actionCol, row)" />
        </div>
      </article>
    </div>
  </div>
</template>

<style scoped>
.gen-wrap {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.gen-card {
  position: relative;
  height: 100%;
  cursor: default;
}

.gen-card.is-openable {
  cursor: pointer;
}

.gen-card .entity-card__meta {
  margin: 4px 0 0;
  padding: 0;
}

.gen-card__check {
  position: absolute;
  top: 10px;
  right: 10px;
  z-index: 2;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  border-radius: 8px;
  background: rgb(255 255 255 / 88%);
}

.gen-card.is-checked {
  border-color: rgba(var(--primary-color), 0.28);
}

.gen-card__stats {
  display: flex;
  align-items: stretch;
  margin: 8px 10px 2px;
  padding: 10px 0 2px;
  border-top: 1px solid var(--ctl-hairline, rgba(var(--primary-color), 0.06));
}

.gen-card__stat {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
  padding: 0 6px;
  text-align: center;
}

.gen-card__stat + .gen-card__stat {
  border-left: 1px solid var(--ctl-hairline, rgba(var(--primary-color), 0.06));
}

.gen-card__stat-val {
  font-size: 15px;
  font-weight: 650;
  letter-spacing: -0.03em;
  font-variant-numeric: tabular-nums;
  color: #1a2332;
  overflow: hidden;
  max-width: 100%;
}

.gen-card__stat-val :deep(span),
.gen-card__stat-val :deep(.text),
.gen-card__stat-val :deep(.text-gray) {
  font-size: inherit;
  font-weight: inherit;
  color: inherit;
}

.gen-card__stat-lab {
  font-size: 10px;
  letter-spacing: -0.01em;
  color: #8b97a5;
}

.gen-card__chips {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 6px;
  padding: 0 0 6px;
}

.gen-card__chips :deep(.n-space),
.gen-card__chips :deep(.flex) {
  justify-content: flex-start !important;
  flex-wrap: wrap;
  gap: 6px !important;
}

.gen-card__chips :deep(.n-tag),
.gen-card__fields :deep(.n-tag) {
  --n-height: 20px !important;
  --n-font-size: 11px !important;
  --n-padding: 0 7px !important;
  font-weight: 550 !important;
}

.gen-card__fields {
  margin: 0;
  padding: 4px 14px 6px;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.gen-card__row {
  display: grid;
  grid-template-columns: 64px 1fr;
  gap: 8px;
  align-items: start;
  min-width: 0;
}

.gen-card__row dt {
  font-size: 11px;
  color: #8b97a5;
  padding-top: 1px;
}

.gen-card__row dd {
  margin: 0;
  min-width: 0;
  font-size: 12px;
  color: #1a2332;
  overflow: hidden;
}

.gen-card__time {
  margin: 0;
  padding: 6px 14px 4px;
  font-size: 11px;
  color: #8b97a5;
}

.gen-card__actions {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 6px;
  margin-top: auto;
  padding: 8px 0 2px;
  border-top: 1px solid var(--ctl-hairline, rgba(var(--primary-color), 0.06));
  overflow: visible;
  min-width: 0;
  max-width: 100%;
}

.gen-card__actions :deep(.n-space),
.gen-card__actions :deep(.n-button-group),
.gen-card__actions :deep(.row-actions) {
  display: flex !important;
  flex-wrap: wrap !important;
  align-items: center;
  gap: 6px;
  width: 100% !important;
  max-width: 100%;
  min-width: 0;
  white-space: normal;
}

.gen-card__actions :deep(.n-button) {
  --n-height: 24px !important;
  --n-font-size: 12px !important;
  --n-padding: 0 8px !important;
  --n-border-radius: 7px !important;
  font-weight: 550 !important;
  flex-shrink: 0;
}

.gen-card__actions :deep(.n-button__content) {
  flex-wrap: nowrap !important;
  white-space: nowrap;
}

.gen-card__actions :deep(.n-button--secondary) {
  --n-border: 0 solid transparent !important;
  --n-border-hover: 0 solid transparent !important;
  --n-border-pressed: 0 solid transparent !important;
  --n-border-focus: 0 solid transparent !important;
  box-shadow: none !important;
}

.gen-card__actions :deep(.n-button--primary-type.n-button--secondary) {
  --n-color: rgba(var(--primary-color), 0.1) !important;
  --n-color-hover: rgba(var(--primary-color), 0.16) !important;
  --n-text-color: rgb(var(--primary-color)) !important;
  --n-text-color-hover: rgb(var(--primary-color)) !important;
}

.gen-card__actions :deep(.n-button--error-type.n-button--secondary),
.gen-card__actions :deep(.n-button--warning-type.n-button--secondary) {
  --n-color: rgb(196 91 74 / 0.1) !important;
  --n-text-color: #c45b4a !important;
}
</style>
