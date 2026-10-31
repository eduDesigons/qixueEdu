<template>
  <CommonPage>
    <div class="content-wrapper">
      <aside class="media-nav">
        <div class="media-nav__list">
          <button
            v-for="cate in cateList"
            :key="cate.id"
            type="button"
            class="media-nav__item"
            :class="{ 'is-active': activeMenu === cate.id }"
            @click="changeCate(cate.id)"
          >
            {{ cate.name }}
          </button>
        </div>
      </aside>

      <div class="right">
        <div class="media-toolbar">
          <n-input
            v-model:value="keyword"
            clearable
            placeholder="搜索资源标题"
            class="media-search"
            @keyup.enter="handleSearch"
          >
            <template #suffix>
              <i class="i-fe:search cursor-pointer" @click="handleSearch" />
            </template>
          </n-input>
          <n-button type="primary" secondary @click="handleSearch">
            搜索
          </n-button>
        </div>

        <div class="media-view-host view-host">
          <Transition :name="viewSwapName">
            <div v-if="isCardLayout" key="card" class="list-box cus-scroll">
              <n-spin :show="loading">
                <div v-if="mediaList && mediaList.length" class="media-grid">
                  <div class="media-grid__inner">
                    <article
                      v-for="item in mediaList"
                      :key="item.id"
                      class="media-card"
                    >
                      <div class="media-card__preview">
                        <MediaPlayer :media-info="item" :controls="false" />
                      </div>
                      <div class="media-card__title" :title="item.title">
                        {{ item.title || '未命名资源' }}
                      </div>
                      <div class="media-card__meta">
                        {{ metaLine(item) }}
                      </div>
                    </article>
                  </div>
                </div>
                <div v-else-if="!loading" class="empty-state media-empty">
                  <p>{{ keyword.trim() ? '没有匹配的资源' : '暂无资源' }}</p>
                  <n-button v-if="keyword.trim()" size="small" secondary @click="clearSearch">
                    清空搜索
                  </n-button>
                </div>
              </n-spin>
            </div>
            <div v-else key="table" class="media-table-wrap">
              <NDataTable
                :columns="columns"
                :data="mediaList"
                :loading="loading"
                :bordered="false"
                :bottom-bordered="false"
                :row-key="row => row.id"
                flex-height
                class="media-table"
              />
            </div>
          </Transition>
        </div>

        <div class="media-pager">
          <n-pagination
            v-model:page="pageInfo.page"
            :page-size="pageInfo.size"
            :item-count="total"
            @update:page="getMediaList"
          />
        </div>
      </div>
    </div>
  </CommonPage>
</template>

<script setup>
import { NDataTable, NTag } from 'naive-ui'
import { computed, h, ref, watch } from 'vue'
import { CommonPage, MediaPlayer } from '@/components/index.js'
import { useRefreshOnActivate } from '@/composables'
import { useAppStore } from '@/store'
import api from './api'

const route = useRoute()
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

const typeMap = {
  MediaImage: 'image',
  MediaVideo: 'video',
  MediaAudio: 'audio',
  MediaDocument: 'document',
  MediaFile: 'file',
}
const typeLabelMap = {
  image: '图片',
  video: '视频',
  audio: '音频',
  document: '文档',
  file: '文件',
  app: '安装包',
}
const cateType = typeMap[route.name] ?? 'image'

const activeMenu = ref('')
const keyword = ref('')
const loading = ref(false)
const cateList = ref([{
  id: '',
  name: '全部',
  sort: 0,
}])
function getCateList() {
  api.cateList(cateType).then((res) => {
    cateList.value = [{
      id: '',
      name: '全部',
      sort: 0,
    }, ...res.data]
  })
}
getCateList()

const pageInfo = ref({
  page: 1,
  size: 12,
})
const total = ref(0)
const mediaList = ref([])
let searchTimer = 0
let skipKeywordWatch = true

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

function typeLabel(item) {
  return typeLabelMap[String(item?.type || cateType).toLowerCase()] || '资源'
}

const columns = [
  {
    title: '预览',
    key: 'preview',
    width: 96,
    render(row) {
      return h('div', { class: 'media-table-preview' }, [
        h(MediaPlayer, {
          mediaInfo: row,
          controls: false,
          compact: true,
        }),
      ])
    },
  },
  {
    title: '标题',
    key: 'title',
    ellipsis: { tooltip: true },
    render(row) {
      return row.title || '未命名资源'
    },
  },
  {
    title: '类型',
    key: 'type',
    width: 96,
    render(row) {
      return h(NTag, {
        size: 'small',
        round: true,
        bordered: false,
        type: 'default',
      }, {
        default: () => typeLabel(row),
      })
    },
  },
  {
    title: '大小',
    key: 'size',
    width: 110,
    render(row) {
      return formatSize(row.size)
    },
  },
  {
    title: '时长',
    key: 'duration',
    width: 96,
    render(row) {
      return formatDuration(row.duration) || '—'
    },
  },
  {
    title: '时间',
    key: 'created_at',
    width: 168,
    render(row) {
      return row.created_at || '—'
    },
  },
]

function getMediaList() {
  loading.value = true
  api.mediaList({
    ...pageInfo.value,
    cateId: activeMenu.value,
    type: cateType,
    title: keyword.value,
  }).then((res) => {
    mediaList.value = res.data.list
    total.value = res.data.total
  }).catch((error) => {
    $message.error(error?.message || '加载资源失败')
  }).finally(() => {
    loading.value = false
  })
}

useRefreshOnActivate(() => getMediaList())

getMediaList()
queueMicrotask(() => {
  skipKeywordWatch = false
})

function handleSearch() {
  pageInfo.value.page = 1
  getMediaList()
}

function clearSearch() {
  keyword.value = ''
}

function changeCate(id) {
  activeMenu.value = id
  pageInfo.value.page = 1
  getMediaList()
}

watch(keyword, () => {
  if (skipKeywordWatch)
    return
  window.clearTimeout(searchTimer)
  searchTimer = window.setTimeout(() => {
    handleSearch()
  }, 280)
})
</script>

<style lang="scss" scoped>
.content-wrapper {
  height: 100%;
  display: flex;
  gap: 16px;
  overflow: hidden;
  padding: 0 0 8px;
}

.media-nav {
  width: 188px;
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  min-height: 0;
  padding: 4px 4px 8px 0 !important;
  overflow: hidden;
}

.media-nav__list {
  flex: 1;
  min-height: 0;
  overflow: auto;
  display: flex;
  flex-direction: column;
  gap: 1px;
  padding: 4px;
  border-radius: 14px;
  background: var(--ctl-field, rgba(var(--primary-color), 0.04));
  border: 1px solid var(--ctl-hairline, rgba(var(--primary-color), 0.08));
}

.media-nav__item {
  appearance: none;
  width: 100%;
  min-height: 36px;
  padding: 8px 12px;
  border: 0;
  border-radius: 11px;
  background: transparent;
  text-align: left;
  font-size: 13px;
  font-weight: 540;
  letter-spacing: -0.01em;
  color: var(--ctl-text, #1a2332);
  cursor: pointer;
  transition:
    background 0.2s ease,
    transform 0.2s cubic-bezier(0.32, 0.72, 0, 1);
}

.media-nav__item:hover {
  background: rgb(255 255 255 / 55%);
}

.media-nav__item:active {
  transform: scale(0.985);
}

.media-nav__item.is-active {
  background: rgb(255 255 255 / 78%);
  font-weight: 600;
}

.right {
  height: 100%;
  display: flex;
  flex-direction: column;
  min-width: 0;
  min-height: 0;
  flex: 1;
  padding: 4px 0 8px 4px !important;
  overflow: hidden !important;
}

.media-toolbar {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 12px;
}

.media-search {
  width: 280px;
  max-width: 100%;
  --n-height: var(--ctl-h-lg, 28px) !important;
  --n-border-radius: 10px !important;
}

.media-view-host {
  flex: 1;
  min-height: 0;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.list-box {
  height: 100%;
  min-height: 0;
  overflow: auto !important;
  width: 100%;
  backface-visibility: hidden;
}

.media-table-wrap {
  height: 100%;
  min-height: 0;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  width: 100%;
  backface-visibility: hidden;
}

.media-grid {
  padding-bottom: 8px;
}

.media-grid__inner {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(236px, 1fr));
  gap: 14px;
}

.media-card {
  display: flex;
  flex-direction: column;
  overflow: hidden;
  border-radius: 16px;
  border: 1px solid var(--ctl-hairline, rgba(var(--primary-color), 0.08));
  background: rgb(255 255 255 / 58%);
  box-shadow: none;
  transition:
    background 0.2s ease,
    border-color 0.2s ease;

  &:hover {
    background: rgb(255 255 255 / 82%);
    border-color: var(--ctl-field-border, rgba(var(--primary-color), 0.16));
  }
}

.media-card__preview {
  position: relative;
  display: block;
  aspect-ratio: 3 / 2;
  height: auto;
  overflow: hidden;
  background:
    radial-gradient(120% 80% at 50% 0%, rgba(var(--primary-color), 0.12), transparent 62%),
    var(--ctl-field, rgba(var(--primary-color), 0.04));
  border-bottom: 1px solid var(--ctl-hairline, rgba(var(--primary-color), 0.06));
}

.media-card__preview :deep(.mp) {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
  min-height: 0;
  border-radius: 0;
}

.media-card__preview :deep(.mp-thumb) {
  width: 100%;
  height: 100%;
  min-height: 0;
  border-radius: 0;
}

.media-card__preview :deep(.mp-image) {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.media-card__title {
  padding: 12px 14px 6px;
  min-height: calc(1.35em * 2 + 18px);
  font-size: 14px;
  font-weight: 600;
  letter-spacing: -0.018em;
  line-height: 1.35;
  color: #1a2332;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.media-card__meta {
  margin-top: auto;
  padding: 0 14px 14px;
  font-size: 11px;
  letter-spacing: -0.01em;
  color: #8b97a5;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.media-empty {
  min-height: 240px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 10px;
}

.media-pager {
  display: flex;
  justify-content: flex-end;
  padding-top: 12px;
}

.media-table {
  height: 100%;
  overflow: hidden;
  border: 0;
  background: transparent;
}

.media-table :deep(.n-data-table) {
  --n-merged-border-color: transparent !important;
  --n-th-color: transparent !important;
  --n-td-color: transparent !important;
  background: transparent !important;
}

.media-table :deep(.n-data-table-th) {
  background: transparent !important;
  font-size: 11px !important;
  font-weight: 550 !important;
  letter-spacing: -0.012em;
  color: #8b97a5 !important;
  border-bottom: 1px solid var(--ctl-hairline, rgba(var(--primary-color), 0.08)) !important;
  padding: 6px 12px 8px !important;
}

.media-table :deep(.n-data-table-td) {
  background: transparent !important;
  font-size: 13px !important;
  vertical-align: middle !important;
  border-bottom: 1px solid var(--ctl-hairline, rgba(var(--primary-color), 0.06)) !important;
  padding: 10px 12px !important;
}

.media-table :deep(.media-table-preview) {
  width: 72px;
  height: 48px;
  border-radius: 10px;
  overflow: hidden;
  background: var(--ctl-field, rgba(var(--primary-color), 0.04));
  border: 1px solid var(--ctl-hairline, rgba(var(--primary-color), 0.08));
}

.media-table :deep(.media-table-preview .mp) {
  width: 100%;
  height: 100%;
  min-height: 0;
  border-radius: 0;
}

.media-table :deep(.media-table-preview .mp-thumb) {
  width: 100%;
  height: 100%;
  min-height: 0;
  border-radius: 0;
}

.media-table :deep(.media-table-preview .mp-image) {
  width: 100%;
  height: 100%;
  object-fit: cover;
}
</style>
