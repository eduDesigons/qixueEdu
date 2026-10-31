<script setup>
import { NButton, NEmpty, NInput, NModal, NPagination, NSpin } from 'naive-ui'
import { computed, nextTick, ref, watch } from 'vue'
import api from '@/api'
import { getAssetUrl } from '@/utils/common'
import MediaPlayer from './MediaPlayer.vue'

const props = defineProps({
  type: { type: String, default: 'image' },
  types: { type: Array, default: () => [] },
  multiple: { type: Boolean, default: false },
})

const emit = defineEmits(['select'])
const visible = defineModel('show', { type: Boolean, default: false })

const TYPE_LABEL = {
  image: '图片',
  video: '视频',
  audio: '音频',
  document: '文件',
  app: '安装包',
  file: '文件',
}

const loading = ref(false)
const keyword = ref('')
const cateId = ref('')
const cateList = ref([{ id: '', name: '全部' }])
const list = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = 12
const selectedMap = ref({})

const typeOptions = computed(() => {
  const list = (props.types || []).filter(Boolean)
  return list.length ? list : [props.type]
})
const activeType = ref(props.type)
const searchInput = ref(null)

const modalTitle = computed(() => `从资源库选择${TYPE_LABEL[activeType.value] || '资源'}`)
const selectedItems = computed(() => Object.values(selectedMap.value))
const selectedCount = computed(() => selectedItems.value.length)
const emptyText = computed(() => keyword.value.trim() ? '没有匹配的资源' : '暂无资源')

function isSelected(id) {
  return Boolean(selectedMap.value[id])
}

function formatSize(bytes) {
  const n = Number(bytes) || 0
  if (n < 1024)
    return `${n} B`
  if (n < 1024 * 1024)
    return `${(n / 1024).toFixed(1)} KB`
  return `${(n / 1024 / 1024).toFixed(1)} MB`
}

function formatDuration(sec) {
  const n = Math.round(Number(sec) || 0)
  if (n <= 0)
    return ''
  const m = Math.floor(n / 60)
  const s = n % 60
  return `${m}:${String(s).padStart(2, '0')}`
}

function metaLine(item) {
  const parts = [formatSize(item.size)]
  const duration = formatDuration(item.duration)
  if (duration)
    parts.push(duration)
  return parts.join(' · ')
}

/** 规范化给 MediaPlayer：补全 CDN，保证 type / 扩展名可识别 */
function mediaOf(item) {
  return {
    ...item,
    url: item?.url ? getAssetUrl(item.url) : '',
    thumb: item?.thumb ? getAssetUrl(item.thumb) : '',
    title: item?.title || item?.name || '',
    type: item?.type || activeType.value,
  }
}

function clearSearch() {
  keyword.value = ''
}

function fetchCates() {
  api.mediaCateList(activeType.value).then((res) => {
    cateList.value = [{ id: '', name: '全部' }, ...(res.data || [])]
  }).catch(() => {
    cateList.value = [{ id: '', name: '全部' }]
  })
}

function fetchList() {
  loading.value = true
  api.mediaList({
    page: page.value,
    size: pageSize,
    type: activeType.value,
    cateId: cateId.value,
    title: keyword.value,
  }).then((res) => {
    list.value = res.data?.list || []
    total.value = res.data?.total || 0
  }).catch((error) => {
    $message.error(error?.message || '加载资源失败')
    list.value = []
    total.value = 0
  }).finally(() => {
    loading.value = false
  })
}

function handleSearch() {
  page.value = 1
  fetchList()
}

let searchTimer = 0
let skipKeywordWatch = false
watch(keyword, () => {
  if (!visible.value || skipKeywordWatch)
    return
  window.clearTimeout(searchTimer)
  searchTimer = window.setTimeout(() => {
    handleSearch()
  }, 280)
})

function changeCate(id) {
  cateId.value = id
  page.value = 1
  fetchList()
}

function changeType(nextType) {
  if (activeType.value === nextType)
    return
  activeType.value = nextType
  cateId.value = ''
  page.value = 1
  selectedMap.value = {}
  fetchCates()
  fetchList()
}

function toggleSelect(item) {
  if (props.multiple) {
    const next = { ...selectedMap.value }
    if (next[item.id])
      delete next[item.id]
    else
      next[item.id] = item
    selectedMap.value = next
    return
  }
  emit('select', item)
  visible.value = false
}

function confirmMultiple() {
  if (!selectedCount.value) {
    $message.warning('请选择资源')
    return
  }
  emit('select', selectedItems.value)
  visible.value = false
}

watch(visible, (show) => {
  window.clearTimeout(searchTimer)
  if (!show)
    return
  skipKeywordWatch = true
  activeType.value = typeOptions.value[0] || props.type
  keyword.value = ''
  cateId.value = ''
  page.value = 1
  selectedMap.value = {}
  fetchCates()
  fetchList()
  queueMicrotask(() => {
    skipKeywordWatch = false
  })
  nextTick(() => searchInput.value?.focus?.())
})
</script>

<template>
  <NModal
    v-model:show="visible"
    preset="card"
    :title="modalTitle"
    class="media-picker-modal"
    :style="{ width: '800px' }"
    :segmented="{ content: true, footer: multiple }"
  >
    <div class="picker">
      <div class="picker__toolbar">
        <NInput
          ref="searchInput"
          v-model:value="keyword"
          clearable
          placeholder="搜索资源标题"
          class="picker__search"
          @keyup.enter="handleSearch"
        >
          <template #suffix>
            <i class="i-fe:search cursor-pointer" @click="handleSearch" />
          </template>
        </NInput>
        <NButton type="primary" secondary @click="handleSearch">
          搜索
        </NButton>
      </div>

      <div v-if="typeOptions.length > 1" class="picker__types">
        <button
          v-for="item in typeOptions"
          :key="item"
          type="button"
          class="picker__type"
          :class="{ 'is-active': activeType === item }"
          @click="changeType(item)"
        >
          {{ TYPE_LABEL[item] || item }}
        </button>
      </div>

      <div v-if="cateList.length > 1" class="picker__cates">
        <button
          v-for="cate in cateList"
          :key="cate.id === '' ? 'all' : cate.id"
          type="button"
          class="picker__cate"
          :class="{ 'is-active': cateId === cate.id }"
          @click="changeCate(cate.id)"
        >
          {{ cate.name }}
        </button>
      </div>

      <NSpin :show="loading">
        <div v-if="list.length" class="picker__grid">
          <article
            v-for="item in list"
            :key="item.id"
            class="picker__card"
            :class="{ 'is-selected': isSelected(item.id) }"
          >
            
            <div class="picker__thumb" @click.stop>
              <MediaPlayer
                :media-info="mediaOf(item)"
                compact
              />
            </div>
            <div class="picker__meta" @click="toggleSelect(item)">
              <div class="picker__name" :title="item.title">
                {{ item.title || '未命名资源' }}
              </div>
              <div class="picker__size">
                {{ metaLine(item) }}
              </div>
              <button type="button" class="picker__use" @click.stop="toggleSelect(item)">
                {{ multiple && isSelected(item.id) ? '取消' : '选择' }}
              </button>
            </div>
            <span v-if="isSelected(item.id)" class="picker__check">已选</span>
          </article>
        </div>
        <NEmpty v-else-if="!loading" :description="emptyText" class="picker__empty">
          <template v-if="keyword.trim()" #extra>
            <NButton size="small" secondary @click="clearSearch">
              清空搜索
            </NButton>
          </template>
        </NEmpty>
      </NSpin>

      <div v-if="total > pageSize" class="picker__pager">
        <NPagination
          v-model:page="page"
          :item-count="total"
          :page-size="pageSize"
          @update:page="fetchList"
        />
      </div>
    </div>

    <template v-if="multiple" #footer>
      <div class="picker__footer">
        <span class="picker__count">已选 {{ selectedCount }} 项</span>
        <NButton @click="visible = false">
          取消
        </NButton>
        <NButton type="primary" :disabled="!selectedCount" @click="confirmMultiple">
          使用所选
        </NButton>
      </div>
    </template>
  </NModal>
</template>

<style scoped>
.picker {
  display: flex;
  flex-direction: column;
  min-height: 420px;
}

.picker__toolbar {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 12px;
}

.picker__search {
  width: 280px;
  max-width: 100%;
}

.picker__types,
.picker__cates {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-bottom: 12px;
}

.picker__type {
  appearance: none;
  height: 30px;
  padding: 0 12px;
  border: 0;
  border-radius: 8px;
  background: var(--ctl-field, rgba(var(--primary-color), 0.06));
  color: var(--ctl-text, #1a2332);
  font-size: 12px;
  font-weight: 600;
  cursor: pointer;
}

.picker__type.is-active {
  background: rgb(var(--primary-color));
  color: #fff;
}

.picker__cate {
  appearance: none;
  height: 28px;
  padding: 0 10px;
  border: 0;
  border-radius: 999px;
  background: var(--ctl-field, rgba(var(--primary-color), 0.06));
  color: var(--ctl-text, #1a2332);
  font-size: 12px;
  font-weight: 550;
  cursor: pointer;
}

.picker__cate.is-active {
  background: rgb(var(--primary-color));
  color: #fff;
}

.picker__grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(148px, 1fr));
  gap: 10px;
}

.picker__card {
  display: flex;
  flex-direction: column;
  overflow: hidden;
  border: 1px solid var(--ctl-hairline, rgba(var(--primary-color), 0.08));
  border-radius: 12px;
  background: rgb(255 255 255 / 55%);
  position: relative;
  transition:
    border-color 0.2s ease,
    background 0.2s ease;
}

.picker__card:hover {
  background: rgb(255 255 255 / 78%);
}

.picker__card.is-selected {
  border-color: rgb(var(--primary-color));
  background: rgba(var(--primary-color), 0.08);
}

.picker__thumb {
  position: relative;
  height: 110px;
  overflow: hidden;
  background:
    radial-gradient(90% 70% at 50% 30%, rgba(var(--primary-color), 0.1), transparent 68%),
    var(--ctl-field, rgba(var(--primary-color), 0.04));
}

.picker__thumb :deep(.mp) {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
  min-height: 0;
  border-radius: 0;
}

.picker__thumb :deep(.mp-thumb) {
  width: 100%;
  height: 100%;
  min-height: 0;
  border-radius: 0;
}

.picker__thumb :deep(.mp-image) {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.picker__thumb :deep(.mp-doc-icon) {
  width: 52px;
  height: 52px;
}

.picker__meta {
  display: flex;
  flex-direction: column;
  gap: 2px;
  padding: 8px 10px 10px;
  cursor: pointer;
}

.picker__name {
  font-size: 12px;
  font-weight: 600;
  color: #1a2332;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.picker__size {
  font-size: 11px;
  color: #8b97a5;
}

.picker__use {
  appearance: none;
  align-self: flex-start;
  margin-top: 6px;
  height: 24px;
  padding: 0 10px;
  border: 0;
  border-radius: 8px;
  background: rgba(var(--primary-color), 0.1);
  color: rgb(var(--primary-color));
  font-size: 11px;
  font-weight: 650;
  cursor: pointer;
}

.picker__use:hover {
  background: rgba(var(--primary-color), 0.16);
}

.picker__check {
  position: absolute;
  top: 8px;
  right: 8px;
  z-index: 2;
  height: 20px;
  padding: 0 6px;
  border-radius: 999px;
  background: rgb(var(--primary-color));
  color: #fff;
  font-size: 11px;
  font-weight: 650;
  line-height: 20px;
  pointer-events: none;
}

.picker__empty {
  min-height: 240px;
}

.picker__pager {
  display: flex;
  justify-content: flex-end;
  padding-top: 12px;
}

.picker__footer {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  gap: 8px;
  width: 100%;
}

.picker__count {
  margin-right: auto;
  font-size: 12px;
  color: #5b6b7c;
}
</style>
