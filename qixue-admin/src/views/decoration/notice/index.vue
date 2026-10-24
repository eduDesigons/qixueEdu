<script setup>
import { Add } from '@vicons/ionicons5'
import { NButton, NDataTable, NIcon, NTag } from 'naive-ui'
import { onMounted, ref, watch } from 'vue'
import api from '@/api'
import { CommonPage } from '@/components/index.js'
import MeCrudCards from '@/components/me/crud/Cards.vue'
import { useRefreshOnActivate } from '@/composables'
import { useAppStore } from '@/store'
import EditNotice from './components/EditNotice.vue'

defineOptions({ name: 'DecorationNotice' })
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
const noticeList = ref([])
const loading = ref(false)
const editNoticeRef = ref()

const columns = [
  {
    title: '标题',
    key: 'title',
    width: 200,
  },
  {
    title: '发布时间',
    key: 'publish_time',
    width: 180,
  },
  {
    title: '状态',
    key: 'status',
    width: 100,
    render(row) {
      return h(NTag, {
        type: row.status === 1 ? 'success' : 'default',
        size: 'small',
        round: true,
        bordered: false,
      }, {
        default: () => row.status === 1 ? '已发布' : '未发布',
      })
    },
  },
  {
    title: '操作',
    key: 'actions',
    width: 150,
    render(row) {
      return h('div', { style: 'display: flex; gap: 8px;' }, [
        h(NButton, {
          size: 'small',
          type: 'primary',
          secondary: true,
          onClick: () => editIt(row),
        }, {
          default: () => '编辑',
        }),
        h(NButton, {
          size: 'small',
          type: 'error',
          secondary: true,
          onClick: () => delIt(row),
        }, {
          default: () => '删除',
        }),
      ])
    },
  },
]

onMounted(() => {
  getNoticeList()
})

useRefreshOnActivate(() => getNoticeList())

function getNoticeList() {
  loading.value = true
  api.decoration.getNoticeList({}).then((res) => {
    noticeList.value = res.data?.list || []
  }).finally(() => {
    loading.value = false
  })
}

function addOne() {
  editNoticeRef.value?.show()
}

function editIt(row) {
  editNoticeRef.value?.show(row)
}

function delIt(row) {
  $dialog.warning({
    title: '确认删除',
    content: '确定要删除这条公告吗？',
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: () => {
      api.decoration.deleteNotice(row.id).then(() => {
        $message.success('删除成功')
        getNoticeList()
      })
    },
  })
}
</script>

<template>
  <CommonPage>
    <template #action>
      <NButton type="primary" @click="addOne()">
        <template #icon>
          <NIcon :component="Add" />
        </template>
        新增公告
      </NButton>
    </template>
    <div class="notice-view-host view-host">
      <Transition :name="viewSwapName">
        <div v-if="isCardLayout" key="card" class="notice-cards cus-scroll">
          <n-spin :show="loading">
            <MeCrudCards v-if="noticeList.length" :list="noticeList" :columns="columns" />
            <div v-else-if="!loading" class="notice-empty">
              暂无数据
            </div>
          </n-spin>
        </div>
        <div v-else key="table" class="notice-table-wrap">
          <NDataTable
            :columns="columns"
            :data="noticeList"
            :loading="loading"
            :bordered="false"
            :bottom-bordered="false"
            flex-height
            class="notice-table"
          />
        </div>
      </Transition>
    </div>
    <EditNotice ref="editNoticeRef" @update="getNoticeList()" />
  </CommonPage>
</template>

<style scoped lang="scss">
.notice-view-host {
  height: 100%;
  min-height: 0;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.notice-cards {
  height: 100%;
  min-height: 0;
  overflow: auto;
  width: 100%;
  backface-visibility: hidden;
}

.notice-table-wrap {
  height: 100%;
  min-height: 0;
  display: flex;
  flex-direction: column;
  overflow: hidden;
  width: 100%;
  backface-visibility: hidden;
}

.notice-empty {
  min-height: 200px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 13px;
  font-weight: 500;
  color: var(--ctl-text-soft, #5b6b7c);
}

.notice-table {
  height: 100%;
  overflow: hidden;
  border: 0;
  background: transparent;
}

.notice-table :deep(.n-data-table) {
  --n-merged-border-color: transparent !important;
  --n-th-color: transparent !important;
  --n-td-color: transparent !important;
  background: transparent !important;
}

.notice-table :deep(.n-data-table-th) {
  background: transparent !important;
  font-size: 11px !important;
  font-weight: 550 !important;
  letter-spacing: -0.012em;
  color: #8b97a5 !important;
  border-bottom: 1px solid var(--ctl-hairline, rgba(var(--primary-color), 0.08)) !important;
  padding: 6px 12px 8px !important;
}

.notice-table :deep(.n-data-table-td) {
  background: transparent !important;
  font-size: 13px !important;
  vertical-align: middle !important;
  border-bottom: 1px solid var(--ctl-hairline, rgba(var(--primary-color), 0.06)) !important;
  padding: 8px 12px !important;
}

.notice-table :deep(.n-data-table-tr:not(.n-data-table-tr--summary):hover .n-data-table-td) {
  background: rgba(var(--primary-color), 0.05) !important;
}

.notice-table :deep(.n-tag) {
  --n-height: 20px !important;
  --n-font-size: 11px !important;
  --n-padding: 0 7px !important;
}

.notice-table :deep(.n-data-table-td .n-button) {
  --n-height: 24px !important;
  --n-font-size: 12px !important;
  --n-padding: 0 8px !important;
  --n-border-radius: 7px !important;
}

.notice-table :deep(.n-button--secondary) {
  --n-border: 0 solid transparent !important;
  --n-border-hover: 0 solid transparent !important;
  box-shadow: none !important;
}

.notice-table :deep(.n-button--primary-type.n-button--secondary) {
  --n-color: rgba(var(--primary-color), 0.1) !important;
  --n-text-color: rgb(var(--primary-color)) !important;
}

.notice-table :deep(.n-button--error-type.n-button--secondary) {
  --n-color: rgb(196 91 74 / 0.1) !important;
  --n-text-color: #c45b4a !important;
}
</style>
