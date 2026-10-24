<!--------------------------------
 - @Author: Ronnie Zhang
 - @LastEditor: Ronnie Zhang
 - @LastEditTime: 2023/12/04 22:51:42
 - @Email: zclzone@outlook.com
 - Copyright © 2023 Ronnie Zhang(大脸怪) | https://isme.top
 --------------------------------->

<template>
  <div ref="wrapRef" class="query-wrap h-full flex flex-col overflow-hidden">
    <NCollapseTransition :show="!collapseQuery || showQuery">
      <form v-if="hasQuerySlot" class="query-panel" @submit.prevent="handleSearchSubmit()">
        <div class="query-panel__fields">
          <n-space :wrap="true" :size="[16, 8]">
            <slot />
          </n-space>
        </div>
        <div class="query-panel__actions">
          <n-button secondary @click="handleReset">
            重置
          </n-button>
          <n-button type="primary" attr-type="submit">
            搜索
          </n-button>
        </div>
      </form>
    </NCollapseTransition>

    <div class="me-view-host view-host">
      <Transition :name="viewSwapName">
        <div v-if="isCardLayout" key="card" class="me-card-layout">
          <div class="me-card-scroll cus-scroll">
            <n-spin :show="loading" :delay="120">
              <slot v-if="tableData.length && hasCardSlot" name="card" :list="tableData" />
              <MeCrudCards
                v-else-if="tableData.length"
                :list="tableData"
                :columns="tableColumns"
                :row-key="rowKey"
                :checked-keys="checkedKeys"
                @update:checked-keys="onChecked"
              />
              <div v-else-if="!loading" class="me-empty">
                <i class="i-fe:inbox me-empty__icon" />
                <p>{{ hasActiveFilters ? '没有符合筛选的结果' : '暂无数据' }}</p>
                <n-button v-if="hasActiveFilters" size="small" secondary @click="handleReset">
                  清空筛选
                </n-button>
              </div>
            </n-spin>
          </div>
          <div v-if="isPagination" class="me-card-pager">
            <n-pagination
              size="small"
              :page="pagination.page"
              :page-size="pagination.pageSize"
              :item-count="pagination.itemCount"
              :page-sizes="pagination.pageSizes"
              show-size-picker
              @update:page="onPageChange"
              @update:page-size="onPageSizeChange"
            >
              <template #prefix="{ itemCount }">
                共 {{ itemCount }} 条数据
              </template>
            </n-pagination>
          </div>
        </div>
        <div v-else key="table" class="me-table-layout">
          <NDataTable
            :remote="remote"
            :loading="loading"
            :scroll-x="scrollX"
            :columns="tableColumns"
            :data="tableData"
            :row-key="(row) => row[rowKey]"
            :checked-row-keys="hasCardSlot ? undefined : checkedKeys"
            :pagination="isPagination ? pagination : false"
            :bordered="false"
            :bottom-bordered="false"
            flex-height
            class="flex-1 me-table"
            @update:checked-row-keys="onChecked"
            @update:page="onPageChange"
            @update:page-size="onPageSizeChange"
          >
            <template #empty>
              <div class="me-empty">
                <i class="i-fe:inbox me-empty__icon" />
                <p>{{ hasActiveFilters ? '没有符合筛选的结果' : '暂无数据' }}</p>
                <n-button v-if="hasActiveFilters" size="small" secondary @click="handleReset">
                  清空筛选
                </n-button>
              </div>
            </template>
          </NDataTable>
        </div>
      </Transition>
    </div>

    <Teleport v-if="canTeleport" :to="filterHostEl">
      <n-button
        size="small"
        :type="showQuery || hasActiveFilters ? 'primary' : 'default'"
        :secondary="!(showQuery || hasActiveFilters)"
        @click="showQuery = !showQuery"
      >
        <i class="i-fe:filter mr-4" />
        筛选
        <span v-if="hasActiveFilters" class="query-host-badge">{{ activeFilterCount }}</span>
        <i class="i-fe:chevron-down ml-4 query-toggle-icon" :class="{ 'is-open': showQuery }" />
      </n-button>
    </Teleport>
    <n-button
      v-else-if="hasQuerySlot && collapseQuery && isInActiveTab"
      size="small"
      class="query-local-trigger"
      :type="showQuery || hasActiveFilters ? 'primary' : 'default'"
      :secondary="!(showQuery || hasActiveFilters)"
      @click="showQuery = !showQuery"
    >
      <i class="i-fe:filter mr-4" />
      筛选
      <span v-if="hasActiveFilters" class="query-host-badge">{{ activeFilterCount }}</span>
    </n-button>
  </div>
</template>

<script setup>
import { NCollapseTransition, NDataTable } from 'naive-ui'
import { utils, writeFile } from 'xlsx'
import { useAppStore } from '@/store'
import { LIST_STATE_TTL_MS } from '@/utils/pageCache'
import MeCrudCards from './Cards.vue'

const props = defineProps({
  /**
   * @remote true: 后端分页  false： 前端分页
   */
  remote: {
    type: Boolean,
    default: true,
  },
  /**
   * @isPagination 是否分页
   */
  isPagination: {
    type: Boolean,
    default: true,
  },
  scrollX: {
    type: Number,
    default: 1200,
  },
  rowKey: {
    type: String,
    default: 'id',
  },
  columns: {
    type: Array,
    required: true,
  },
  /** queryBar中的参数 */
  queryItems: {
    type: Object,
    default() {
      return {}
    },
  },
  /**
   * ! 约定接口入参出参
   * 分页模式需约定分页接口入参
   *    @pageSize 分页参数：一页展示多少条，默认10
   *    @pageNo   分页参数：页码，默认1
   * 需约定接口出参
   *    @pageData 分页模式必须,非分页模式如果没有pageData则取上一层data
   *    @total    分页模式必须，非分页模式如果没有total则取上一层data.length
   */
  getData: {
    type: Function,
    required: true,
  },
  /** 分页缓存 key，同一页多个表格时请传入不同值 */
  cacheKey: {
    type: String,
    default: '',
  },
  /**
   * 强制绑定的上下文参数（如 course_id / chapter_id / path_id / user_id）
   * 每次请求都会合并进入参，且清空筛选时不会被洗掉
   */
  boundQuery: {
    type: Object,
    default: null,
  },
  /** 为 true 时默认收起筛选，通过标题旁按钮展开 */
  collapseQuery: {
    type: Boolean,
    default: true,
  },
})

const emit = defineEmits(['update:queryItems', 'onChecked', 'onDataChange'])
const route = useRoute()
const loading = ref(false)
const initQuery = { ...props.queryItems }
const tableData = ref([])
const pagination = reactive({
  page: 1,
  pageSize: 10,
  showSizePicker: true,
  pageSizes: [10, 20, 50],
  size: 'small',
  prefix({ itemCount }) {
    return `共 ${itemCount} 条数据`
  },
})

function normalizeBoundQuery(raw = props.boundQuery) {
  if (!raw || typeof raw !== 'object')
    return {}
  const next = {}
  Object.keys(raw).forEach((key) => {
    const val = raw[key]
    if (val === undefined)
      return
    next[key] = val
  })
  return next
}

function hasBoundQueryValue(bound = normalizeBoundQuery()) {
  return Object.keys(bound).some((key) => {
    const val = bound[key]
    return val !== null && val !== undefined && val !== ''
  })
}

function syncBoundQueryIntoItems(force = false) {
  const bound = normalizeBoundQuery()
  if (!Object.keys(bound).length)
    return false
  const current = props.queryItems || {}
  let changed = false
  const next = { ...current }
  Object.keys(bound).forEach((key) => {
    if (bound[key] === undefined)
      return
    if (force || String(current[key] ?? '') !== String(bound[key] ?? '')) {
      next[key] = bound[key]
      changed = true
    }
  })
  if (changed)
    emit('update:queryItems', next)
  return changed
}

function isActionColumn(col) {
  if (!col || col.type === 'selection')
    return false
  return col.key === 'actions' || col.key === 'action' || col.title === '操作'
}

function hasRowActionsClass(vnode) {
  if (!vnode || typeof vnode !== 'object' || Array.isArray(vnode))
    return false
  const cls = vnode.props?.class
  if (typeof cls === 'string')
    return cls.split(/\s+/).includes('row-actions')
  if (Array.isArray(cls))
    return cls.some(item => String(item).split(/\s+/).includes('row-actions'))
  if (cls && typeof cls === 'object')
    return Boolean(cls['row-actions'])
  return false
}

/** 操作列统一包一层，保证列表视图按钮横排对齐 */
const tableColumns = computed(() => {
  return (props.columns || []).map((col) => {
    if (!isActionColumn(col) || typeof col.render !== 'function')
      return col
    const rawRender = col.render
    return {
      ...col,
      render: (...args) => {
        const result = rawRender(...args)
        if (result == null || result === false)
          return result
        if (hasRowActionsClass(result))
          return result
        const children = Array.isArray(result) ? result.filter(Boolean) : [result]
        return h('div', { class: 'row-actions' }, children)
      },
    }
  })
})

function getEnclosingTabPaneName(instance) {
  let parent = instance?.parent
  while (parent) {
    const type = parent.type
    if (type?.__TAB_PANE__ || type?.name === 'TabPane')
      return parent.props?.name
    parent = parent.parent
  }
  return null
}

/** 详情 Tab 内嵌表不写 sessionStorage，避免换课节时筛选项串台 */
const isNestedInTab = !!getEnclosingTabPaneName(getCurrentInstance())

function getCacheKey() {
  const colKeys = (props.columns || []).map(item => item.key).filter(Boolean).join(',')
  const query = props.queryItems || {}
  const ctx = {}
  Object.keys(query).forEach((key) => {
    if (/(_id|Id)$/.test(key) || ['type', 'userId', 'pathId', 'paperId', 'courseId'].includes(key))
      ctx[key] = query[key]
  })
  return `me-crud:${props.cacheKey || route.path}:${JSON.stringify(ctx)}:${colKeys}`
}

function saveState() {
  if (!props.isPagination || isNestedInTab)
    return
  try {
    sessionStorage.setItem(getCacheKey(), JSON.stringify({
      page: pagination.page,
      pageSize: pagination.pageSize,
      query: { ...props.queryItems },
      savedAt: Date.now(),
    }))
  }
  catch {}
}

function restoreState() {
  if (!props.isPagination || isNestedInTab)
    return
  try {
    const raw = sessionStorage.getItem(getCacheKey())
    if (!raw)
      return
    const saved = JSON.parse(raw)
    const savedAt = Number(saved?.savedAt) || 0
    if (!savedAt || Date.now() - savedAt > LIST_STATE_TTL_MS) {
      sessionStorage.removeItem(getCacheKey())
      return
    }
    if (saved.pageSize)
      pagination.pageSize = saved.pageSize
    if (saved.page)
      pagination.page = saved.page
    if (saved.query && typeof saved.query === 'object') {
      const next = { ...props.queryItems, ...saved.query }
      Object.keys(initQuery).forEach((key) => {
        if (/(_id|Id)$/.test(key) || ['type', 'userId', 'pathId', 'paperId', 'courseId'].includes(key))
          next[key] = props.queryItems[key] ?? initQuery[key]
      })
      emit('update:queryItems', next)
    }
  }
  catch {}
}

onBeforeMount(() => {
  syncBoundQueryIntoItems(true)
  restoreState()
  // restore 可能盖掉绑定 id，再同步一次
  syncBoundQueryIntoItems(true)
})

onDeactivated(() => {
  saveState()
})

const skipFirstActivated = ref(true)
onActivated(async () => {
  if (skipFirstActivated.value) {
    skipFirstActivated.value = false
    return
  }
  syncBoundQueryIntoItems()
  // KeepAlive 已保留内存态；仅保存当前页，不强制重拉（避免闪一下）
  saveState()
})

watch(
  () => props.boundQuery,
  (next, prev) => {
    if (!next || typeof next !== 'object')
      return
    const bound = normalizeBoundQuery(next)
    const prevBound = normalizeBoundQuery(prev)
    const nextKey = JSON.stringify(bound)
    const prevKey = JSON.stringify(prevBound)
    if (nextKey === prevKey)
      return

    const current = props.queryItems || {}
    const cleaned = { ...current }
    Object.keys(cleaned).forEach((key) => {
      if (Object.prototype.hasOwnProperty.call(bound, key) || isContextQueryKey(key))
        return
      cleaned[key] = Object.prototype.hasOwnProperty.call(initQuery, key) ? initQuery[key] : null
    })
    Object.keys(bound).forEach((key) => {
      cleaned[key] = bound[key]
    })
    emit('update:queryItems', cleaned)
    pagination.page = 1
    nextTick(() => handleQuery())
  },
  { deep: true },
)

const wrapRef = ref(null)
const isCrudVisible = ref(true)
const slots = useSlots()
const filterHost = inject('pageFilterHost', null)
const detailActiveTab = inject('detailActiveTab', null)
const enclosingTabName = getEnclosingTabPaneName(getCurrentInstance())
const showQuery = ref(!props.collapseQuery)

const hasQuerySlot = computed(() => !!slots.default)
const hasCardSlot = computed(() => !!slots.card)
const appStore = useAppStore()
const isCardLayout = computed(() => appStore.listView === 'card')
const viewSwapName = ref('view-to-card')
watch(
  () => appStore.listView,
  (next, prev) => {
    if (!prev || prev === next)
      return
    viewSwapName.value = next === 'table' ? 'view-to-table' : 'view-to-card'
  },
)
const checkedKeys = ref([])
const filterHostEl = computed(() => filterHost?.value || null)
const isInActiveTab = computed(() => {
  if (enclosingTabName == null || !detailActiveTab)
    return true
  return detailActiveTab.value === enclosingTabName
})
const canTeleport = computed(() =>
  props.collapseQuery
  && hasQuerySlot.value
  && !!filterHostEl.value
  && isCrudVisible.value
  && isInActiveTab.value,
)

function isEmptyQueryValue(value) {
  return value === null || value === undefined || value === ''
}

function isContextQueryKey(key) {
  return /(_id|Id)$/.test(key) || ['type', 'userId', 'pathId', 'paperId', 'courseId'].includes(key)
}

const activeFilterCount = computed(() => {
  const current = props.queryItems || {}
  const keys = new Set([...Object.keys(initQuery), ...Object.keys(current)])
  let count = 0
  for (const key of keys) {
    if (isContextQueryKey(key))
      continue
    const a = initQuery[key]
    const b = current[key]
    if (isEmptyQueryValue(a) && isEmptyQueryValue(b))
      continue
    if (a !== b)
      count++
  }
  return count
})

const hasActiveFilters = computed(() => activeFilterCount.value > 0)

let visibilityObserver = null
onMounted(() => {
  const el = wrapRef.value
  if (!el || typeof IntersectionObserver === 'undefined')
    return
  visibilityObserver = new IntersectionObserver((entries) => {
    const entry = entries[0]
    isCrudVisible.value = !!(entry && entry.isIntersecting)
  }, { threshold: 0.05 })
  visibilityObserver.observe(el)
})
onBeforeUnmount(() => {
  visibilityObserver?.disconnect()
  visibilityObserver = null
})

async function handleQuery() {
  try {
    const bound = normalizeBoundQuery()
    if (Object.keys(bound).length) {
      syncBoundQueryIntoItems()
      if (!hasBoundQueryValue(bound)) {
        $message.error('列表参数缺失')
        return
      }
    }
    loading.value = true
    let paginationParams = {}
    // 如果非分页模式或者使用前端分页,则无需传分页参数
    if (props.isPagination && props.remote) {
      paginationParams = { page: pagination.page, size: pagination.pageSize }
    }
    const { data } = await props.getData({
      ...props.queryItems,
      ...bound,
      ...paginationParams,
    })
    tableData.value = data?.list || data
    pagination.itemCount = data.total ?? data.length
    if (pagination.itemCount && !tableData.value.length && pagination.page > 1) {
      // 如果当前页数据为空，且总条数不为0，则返回上一页数据
      onPageChange(pagination.page - 1)
      return
    }
    saveState()
    resetChecked()
  }
  catch (error) {
    console.error(error)
    $message.error(error?.message || '列表加载失败')
  }
  finally {
    emit('onDataChange', tableData.value)
    loading.value = false
  }
}

function handleSearch(keepCurrentPage = true) {
  if (keepCurrentPage || !props.remote) {
    handleQuery()
  }
  else {
    onPageChange(1)
  }
}

function handleSearchSubmit() {
  pagination.page = 1
  handleQuery()
}

async function handleReset() {
  const current = props.queryItems || {}
  const bound = normalizeBoundQuery()
  const next = { ...initQuery }
  // 先清空非上下文字段，再强制保留当前 *_id 等上下文，避免 chapter_id / course_id 被洗成 null
  for (const key of Object.keys(current)) {
    if (isContextQueryKey(key) || Object.prototype.hasOwnProperty.call(bound, key))
      next[key] = bound[key] ?? current[key] ?? initQuery[key] ?? null
    else if (!(key in initQuery))
      next[key] = null
  }
  for (const key of Object.keys(initQuery)) {
    if (isContextQueryKey(key) || Object.prototype.hasOwnProperty.call(bound, key))
      next[key] = bound[key] ?? current[key] ?? initQuery[key] ?? null
    else
      next[key] = initQuery[key]
  }
  Object.keys(bound).forEach((key) => {
    next[key] = bound[key]
  })
  emit('update:queryItems', next)
  await nextTick()
  pagination.page = 1
  try {
    sessionStorage.removeItem(getCacheKey())
  }
  catch {}
  handleQuery()
}
function onPageChange(currentPage) {
  pagination.page = currentPage
  if (props.remote) {
    handleQuery()
  }
  else {
    saveState()
  }
}
function onPageSizeChange(pageSize) {
  pagination.pageSize = pageSize
  pagination.page = 1
  if (props.remote) {
    handleQuery()
  }
  else {
    saveState()
  }
}
pagination.onUpdatePageSize = onPageSizeChange
function onChecked(rowKeys) {
  checkedKeys.value = rowKeys || []
  if (props.columns.some(item => item.type === 'selection'))
    emit('onChecked', checkedKeys.value)
}

function resetChecked() {
  if (!checkedKeys.value.length)
    return
  onChecked([])
}
function handleExport(columns = props.columns, data = tableData.value) {
  if (!data?.length)
    return $message.warning('没有数据')
  const columnsData = columns.filter(item => !!item.title && !item.hideInExcel)
  const thKeys = columnsData.map(item => item.key)
  const thData = columnsData.map(item => item.title)
  const trData = data.map(item => thKeys.map(key => item[key]))
  const sheet = utils.aoa_to_sheet([thData, ...trData])
  const workBook = utils.book_new()
  utils.book_append_sheet(workBook, sheet, '数据报表')
  writeFile(workBook, '数据报表.xlsx')
}

defineExpose({
  handleSearch,
  handleReset,
  handleExport,
  tableData,
})
</script>

<style scoped>
.query-wrap {
  position: relative;
}

.query-panel {
  position: relative;
  overflow: visible;
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 14px;
  margin: 0 0 12px;
  padding: 2px 0 12px;
  border: 0;
  border-bottom: 1px solid var(--ctl-hairline, rgba(var(--primary-color), 0.08));
  border-radius: 0;
  background: transparent;
  box-shadow: none;
  transition: none;
}

.query-panel:focus-within {
  background: transparent;
  box-shadow: none;
}

.query-panel__fields,
.query-panel__actions {
  position: relative;
  z-index: 1;
}

.query-panel__fields {
  min-width: 0;
  flex: 1;
}

.query-panel__fields :deep(.n-space) {
  gap: 8px 12px !important;
  align-items: center !important;
}

.query-panel__fields :deep(.n-space > *) {
  flex-shrink: 0 !important;
  max-width: 100%;
}

.query-panel__actions {
  display: flex;
  flex-shrink: 0;
  align-items: center;
  gap: 8px;
  min-height: 34px;
}

.query-panel__actions :deep(.n-button) {
  --n-height: var(--ctl-h-lg, 28px) !important;
  --n-font-size: 13px !important;
  --n-border-radius: 10px !important;
  font-weight: 550 !important;
  letter-spacing: -0.012em;
  transition: transform 0.2s cubic-bezier(0.32, 0.72, 0, 1) !important;
}

.query-panel__actions :deep(.n-button:not(.n-button--disabled):active) {
  transform: scale(0.97);
}

.query-host-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 16px;
  height: 16px;
  margin-left: 6px;
  padding: 0 4px;
  border-radius: 999px;
  background: rgb(var(--primary-color));
  color: #fff;
  font-size: 11px;
  font-weight: 600;
  line-height: 16px;
}

.me-table {
  border-radius: 0;
  overflow: hidden;
  border: 0;
  background: transparent;
}

.me-table :deep(.n-data-table) {
  --n-border-radius: 0 !important;
  --n-merged-border-color: transparent !important;
  --n-th-color: transparent !important;
  --n-td-color: transparent !important;
  --n-td-color-hover: rgba(var(--primary-color), 0.06) !important;
  background: transparent !important;
}

.me-table :deep(.n-data-table-wrapper),
.me-table :deep(.n-data-table-base-table),
.me-table :deep(.n-data-table-base-table-header),
.me-table :deep(.n-data-table-base-table-body) {
  background: transparent !important;
}

.me-table :deep(.n-data-table-th) {
  font-size: 11px !important;
  letter-spacing: -0.012em;
  font-weight: 550 !important;
  color: #8b97a5 !important;
  background: transparent !important;
  border-bottom: 1px solid var(--ctl-hairline, rgba(var(--primary-color), 0.08)) !important;
  padding: 6px 12px 8px !important;
}

.me-table :deep(.n-data-table-td) {
  font-size: 13px !important;
  letter-spacing: -0.014em;
  line-height: 1.35;
  vertical-align: middle !important;
  background: transparent !important;
  border-bottom: 1px solid var(--ctl-hairline, rgba(var(--primary-color), 0.07)) !important;
  padding: 8px 12px !important;
}

.me-table :deep(.n-data-table-th--selection),
.me-table :deep(.n-data-table-td--selection) {
  padding-left: 6px !important;
  padding-right: 4px !important;
  border-bottom-color: transparent !important;
}

.me-table :deep(.n-data-table-tr:last-child .n-data-table-td) {
  border-bottom-color: transparent !important;
}

.me-table :deep(.n-data-table-td .flex:not(.flex-col)) {
  align-items: center;
}

.me-table :deep(.n-data-table-td .flex-col) {
  align-items: flex-start;
}

.me-table :deep(.n-data-table-tr:not(.n-data-table-tr--summary):hover .n-data-table-td) {
  background: rgba(var(--primary-color), 0.05) !important;
}

.me-table :deep(.n-data-table-tr--checked .n-data-table-td),
.me-table :deep(.n-data-table-tr.n-data-table-tr--checked:hover .n-data-table-td) {
  background: rgba(var(--primary-color), 0.08) !important;
}

.me-table :deep(.n-avatar) {
  width: 40px !important;
  height: 40px !important;
  --n-merged-size: 40px !important;
}

.me-table :deep(.n-avatar:not(.n-avatar--round)) {
  border-radius: 10px !important;
}

.me-table :deep(.n-checkbox) {
  --n-size: 16px !important;
}

.me-table :deep(.n-data-table-td .text-14) {
  font-size: 13px !important;
  font-weight: 600 !important;
  letter-spacing: -0.016em;
}

.me-table :deep(.n-data-table-td .text-12) {
  font-size: 11px !important;
  letter-spacing: -0.01em;
}

.me-table :deep(.n-data-table-td .ml-12) {
  margin-left: 10px !important;
}

.me-table :deep(.n-data-table-td .mt-4) {
  margin-top: 2px !important;
}

.me-table :deep(.n-data-table-td .line-clamp-2) {
  -webkit-line-clamp: 1;
}

.me-table :deep(.n-tag) {
  --n-height: 20px !important;
  --n-font-size: 11px !important;
  --n-padding: 0 7px !important;
  font-weight: 550 !important;
}

.me-table :deep(.n-data-table-td .n-button) {
  --n-height: 24px !important;
  --n-font-size: 12px !important;
  --n-padding: 0 8px !important;
  --n-border-radius: 7px !important;
  font-weight: 550 !important;
  flex-shrink: 0;
}

.me-table :deep(.n-button--secondary) {
  --n-border: 0 solid transparent !important;
  --n-border-hover: 0 solid transparent !important;
  --n-border-pressed: 0 solid transparent !important;
  --n-border-focus: 0 solid transparent !important;
  box-shadow: none !important;
}

.me-table :deep(.n-button--primary-type.n-button--secondary) {
  --n-color: rgba(var(--primary-color), 0.1) !important;
  --n-color-hover: rgba(var(--primary-color), 0.16) !important;
  --n-color-pressed: rgba(var(--primary-color), 0.2) !important;
  --n-text-color: rgb(var(--primary-color)) !important;
  --n-text-color-hover: rgb(var(--primary-color)) !important;
}

.me-table :deep(.n-button--info-type.n-button--secondary) {
  --n-color: rgb(56 132 156 / 0.1) !important;
  --n-color-hover: rgb(56 132 156 / 0.16) !important;
}

.me-table :deep(.n-button--error-type.n-button--secondary) {
  --n-color: rgb(196 91 74 / 0.1) !important;
  --n-color-hover: rgb(196 91 74 / 0.16) !important;
  --n-text-color: #c45b4a !important;
  --n-text-color-hover: #c45b4a !important;
}

.me-table :deep(.n-button--warning-type.n-button--secondary) {
  --n-color: rgb(196 140 60 / 0.12) !important;
  --n-color-hover: rgb(196 140 60 / 0.18) !important;
  --n-text-color: #a67c2a !important;
  --n-text-color-hover: #a67c2a !important;
}

.me-table :deep(.n-data-table-td:last-child) {
  white-space: nowrap;
}

.me-table :deep(.n-data-table-td:last-child .n-data-table-td__ellipsis),
.me-table :deep(.n-data-table-td:last-child > div),
.me-table :deep(.row-actions) {
  display: inline-flex !important;
  flex-wrap: nowrap !important;
  align-items: center !important;
  justify-content: center;
  gap: 6px !important;
  vertical-align: middle;
  white-space: nowrap;
  max-width: 100%;
}

.gen-card :deep(.row-actions) {
  display: flex !important;
  flex-wrap: wrap !important;
  align-items: center !important;
  justify-content: flex-start;
  gap: 6px !important;
  max-width: 100%;
  min-width: 0;
  white-space: normal;
}

.me-table :deep(.n-data-table-td:last-child .n-button),
.me-table :deep(.n-data-table-td:last-child .n-space),
.me-table :deep(.n-data-table-td:last-child .n-space-item),
.me-table :deep(.row-actions .n-button),
.me-table :deep(.row-actions .n-space),
.me-table :deep(.row-actions .n-space-item),
.gen-card :deep(.row-actions .n-button) {
  flex-shrink: 0;
  margin: 0 !important;
  margin-left: 0 !important;
  margin-right: 0 !important;
}

.me-table :deep(.n-data-table-td:last-child .n-space),
.me-table :deep(.row-actions .n-space) {
  display: inline-flex !important;
  flex-wrap: nowrap !important;
  align-items: center !important;
  gap: 6px !important;
}

.me-table :deep(.n-data-table-td:last-child .n-button + .n-button),
.me-table :deep(.row-actions .n-button + .n-button) {
  margin-left: 0 !important;
}

.me-table :deep(.n-data-table__pagination) {
  padding: 8px 2px 0 !important;
  margin: 0 !important;
  border-top: 1px solid var(--ctl-hairline, rgba(var(--primary-color), 0.06));
  background: transparent !important;
}

.me-table :deep(.n-pagination),
.me-card-pager :deep(.n-pagination) {
  padding: 0;
  gap: 2px !important;
  --n-item-size: 26px !important;
  --n-item-font-size: 12px !important;
}

.me-table :deep(.n-pagination-item),
.me-card-pager :deep(.n-pagination-item) {
  min-width: 26px !important;
  height: 26px !important;
  border: 0 !important;
  border-radius: 8px !important;
  background: transparent !important;
  font-weight: 550 !important;
}

.me-table :deep(.n-pagination-item:not(.n-pagination-item--disabled):hover),
.me-card-pager :deep(.n-pagination-item:not(.n-pagination-item--disabled):hover) {
  background: rgba(var(--primary-color), 0.08) !important;
}

.me-table :deep(.n-pagination-item--active),
.me-card-pager :deep(.n-pagination-item--active) {
  background: rgba(var(--primary-color), 0.14) !important;
  color: rgb(var(--primary-color)) !important;
  font-weight: 650 !important;
}

.me-table :deep(.n-pagination-prefix),
.me-card-pager :deep(.n-pagination-prefix) {
  font-size: 12px;
  color: #8b97a5;
}

.me-empty {
  min-height: 220px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 10px;
  padding: 48px 16px;
  font-size: 13px;
  font-weight: 500;
  letter-spacing: -0.012em;
  color: var(--ctl-text-soft, #5b6b7c);
}

.me-empty__icon {
  width: 28px;
  height: 28px;
  opacity: 0.35;
  color: rgb(var(--primary-color));
}

.me-empty p {
  margin: 0;
}

.query-host-group {
  display: flex;
  align-items: center;
  gap: 8px;
}

.me-view-host {
  display: flex;
  flex-direction: column;
  flex: 1;
  min-height: 0;
  overflow: hidden;
}

.me-card-layout {
  display: flex;
  flex-direction: column;
  flex: 1;
  min-height: 0;
  width: 100%;
  backface-visibility: hidden;
}

.me-table-layout {
  display: flex;
  flex-direction: column;
  flex: 1;
  min-height: 0;
  width: 100%;
  overflow: hidden;
  backface-visibility: hidden;
}

.me-card-scroll {
  flex: 1;
  min-height: 0;
  overflow: auto;
  padding: 2px 2px 8px;
}

.me-card-pager {
  display: flex;
  justify-content: flex-end;
  padding: 8px 2px 0;
  border-top: 1px solid var(--ctl-hairline, rgba(var(--primary-color), 0.06));
}

.query-local-trigger {
  position: absolute;
  top: 0;
  right: 0;
  z-index: 2;
}

.query-toggle-icon {
  transition: transform 0.24s cubic-bezier(0.32, 0.72, 0, 1);
}

.query-toggle-icon.is-open {
  transform: rotate(180deg);
}
</style>
