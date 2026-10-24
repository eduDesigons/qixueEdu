<script setup>
import { NButton, NImage, NModal, NSwitch, NTabPane, NTag } from 'naive-ui'
import { h, nextTick, onActivated, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import api from '@/api'
import { CommonPage, DetailTabs, MeCrud, MeQueryItem, CityQueryItem } from '@/components'
import { previewImage, useCrud } from '@/composables'
import { formatDateTime, getAssetUrl } from '@/utils'

defineOptions({ name: 'MemberCert' })

const router = useRouter()
const route = useRoute()
const $table = ref(null)
const $issuedTable = ref(null)
const queryItems = ref({})
const issuedQuery = ref({
  city_id: null,
})
const activeTab = ref(route.query.tab === 'issued' ? 'issued' : 'tpl')

const showViewModal = ref(false)
const viewCert = ref(null)
const viewLoading = ref(false)

onMounted(() => {
  searchActiveTab()
})

onActivated(() => {
  if (route.query.tab === 'issued' || route.query.tab === 'tpl')
    activeTab.value = route.query.tab
  searchActiveTab()
})

watch(() => route.query.tab, (tab) => {
  if (tab === 'issued' || tab === 'tpl')
    activeTab.value = tab
})

function searchActiveTab(tab = activeTab.value) {
  nextTick(() => {
    if (tab === 'issued')
      $issuedTable.value?.handleSearch()
    else
      $table.value?.handleSearch()
  })
}

const { handleDelete } = useCrud({
  name: '证书',
  doDelete: (id) => { return api.cert.delete(id) },
  refresh: () => $table.value?.handleSearch(),
})

function openView(row) {
  showViewModal.value = true
  viewCert.value = null
  viewLoading.value = true
  api.cert.detail(row.id).then((res) => {
    viewCert.value = res.data
  }).catch((error) => {
    $message.error(error?.message || '加载证书失败')
    showViewModal.value = false
  }).finally(() => {
    viewLoading.value = false
  })
}

function closeViewModal() {
  showViewModal.value = false
  viewCert.value = null
}

async function toggleShow(row, value) {
  const next = Number(value) === 1 ? 1 : 0
  const prev = Number(row.is_show) === 1 ? 1 : 0
  if (next === prev)
    return
  row._showLoading = true
  row.is_show = next
  try {
    await api.cert.setShow({ id: row.id, is_show: next })
    $message.success(next ? '已上架' : '已下架')
  }
  catch {
    row.is_show = prev
  }
  finally {
    row._showLoading = false
  }
}

// 占位 key 对应的展示文字（与 create 页左侧一致）
const paramKeyLabels = {
  name: '姓名',
  course: '名称/标题',
  time: '时间',
  no: '证书编号',
}
function getParamLabel(key) {
  return paramKeyLabels[key] || key
}

const SAMPLE_TEXT = {
  name: '张三',
  course: '示例课程 / 系列名称',
  time: '2026年8月21日',
  no: 'QX202608210001',
}
function getParamPreview(p) {
  return SAMPLE_TEXT[p.key] || getParamLabel(p.key)
}

// 预览图缩放比例（原画布 800x500，弹窗内缩小）
const PREVIEW_SCALE = 0.35
const PREVIEW_W = 800 * PREVIEW_SCALE
const PREVIEW_H = 500 * PREVIEW_SCALE

const columns = [
  {
    title: '证书名称',
    key: 'name',
    width: 180,
    ellipsis: { tooltip: true },
  },
  {
    title: '预览图',
    key: 'preview_image',
    width: 120,
    align: 'center',
    render(row) {
      const url = row.preview_image || row.template_image
      if (!url)
        return h('span', { class: 'text-gray-400' }, '-')
      return h(NImage, {
        width: 80,
        height: 50,
        src: getAssetUrl(url),
        style: { objectFit: 'cover', borderRadius: '10px' },
      })
    },
  },
  {
    title: '状态',
    key: 'is_show',
    width: 100,
    align: 'center',
    render(row) {
      return h(NSwitch, {
        value: Number(row.is_show) === 1 ? 1 : 0,
        checkedValue: 1,
        uncheckedValue: 0,
        size: 'small',
        loading: !!row._showLoading,
        onUpdateValue: (v) => {
          toggleShow(row, v)
        },
      })
    },
  },
  {
    title: '占位框数',
    key: 'params',
    width: 100,
    align: 'center',
    render: ({ params }) => h(NTag, { size: 'small', type: 'info', round: true, bordered: false }, { default: () => (Array.isArray(params) ? params.length : 0) }),
  },
  {
    title: '创建时间',
    key: 'created_at',
    width: 180,
    align: 'center',
    render: row => h('span', formatDateTime(row.created_at)),
  },
  {
    title: '操作',
    key: 'actions',
    width: 220,
    align: 'center',
    fixed: 'right',
    hideInExcel: true,
    render(row) {
      return h('div', { class: 'flex gap-8 justify-center' }, [
        h(NButton, {
          size: 'small',
          type: 'info',
          secondary: true,
          onClick: () => openView(row),
        }, { default: () => '查看' }),
        h(NButton, {
          size: 'small',
          type: 'primary',
          secondary: true,
          onClick: () => router.push({ path: '/course/cert/create', query: { id: row.id } }),
        }, { default: () => '编辑' }),
        h(NButton, {
          size: 'small',
          type: 'error',
          secondary: true,
          onClick: () => handleDelete(row.id),
        }, { default: () => '删除' }),
      ])
    },
  },
]

function previewIssued(url) {
  previewImage(url)
}

const issuedColumns = [
  {
    title: '学员',
    key: 'user',
    width: 140,
    render: ({ user }) => user?.nick_name || '-',
  },
  {
    title: '证书',
    key: 'title',
    ellipsis: { tooltip: true },
  },
  {
    title: '来源',
    key: 'source',
    width: 200,
    ellipsis: { tooltip: true },
    render: (row) => {
      if (row.path?.title)
        return `系列 · ${row.path.title}`
      if (row.course?.title)
        return `课程 · ${row.course.title}`
      return '-'
    },
  },
  {
    title: '编号',
    key: 'cert_no',
    width: 200,
  },
  {
    title: '预览',
    key: 'image',
    width: 80,
    align: 'center',
    render: ({ image }) => image
      ? h(NImage, { src: getAssetUrl(image), width: 48, height: 32, objectFit: 'cover' })
      : '-',
  },
  {
    title: '颁发时间',
    key: 'issued_at',
    width: 180,
    render: row => h('span', formatDateTime(row.issued_at || row.created_at)),
  },
  {
    title: '操作',
    key: 'actions',
    width: 160,
    align: 'center',
    render(row) {
      return h('div', { class: 'flex gap-8 justify-center' }, [
        h(NButton, { size: 'small', secondary: true, type: 'primary', onClick: () => previewIssued(row.image) }, { default: () => '查看' }),
        h(NButton, {
          size: 'small',
          type: 'info',
          secondary: true,
          onClick: () => reissue(row),
        }, { default: () => '重新生成' }),
      ])
    },
  },
]

function reissue(row) {
  const payload = {
    user_id: row.user_id,
    force: true,
  }
  if (row.course_id) {
    payload.course_id = row.course_id
  }
  else if (row.path_id) {
    payload.path_id = row.path_id
  }
  else {
    $message.warning('缺少系列或课程信息')
    return
  }
  $dialog.warning({
    title: '重新生成证书',
    content: `将按最新模板为「${row.user?.nick_name || '学员'}」重新生成证书。`,
    positiveText: '确定',
    negativeText: '取消',
    async onPositiveClick() {
      try {
        await api.cert.issue(payload)
        $message.success('已重新生成')
        $issuedTable.value?.handleSearch()
      }
      catch {
        // 拦截器已提示失败原因
      }
    },
  })
}

watch(activeTab, (tab) => {
  if (route.query.tab !== tab) {
    router.replace({
      query: {
        ...route.query,
        tab,
      },
    })
  }
  searchActiveTab(tab)
})
</script>

<template>
  <CommonPage>
    <template #action>
      <NButton v-if="activeTab === 'tpl'" type="primary" @click="router.push('/course/cert/create')">
        <i class="i-material-symbols:add mr-4 text-18" />
        新增证书
      </NButton>
    </template>

    <DetailTabs v-model:value="activeTab">
      <NTabPane name="tpl" tab="证书模板" display-directive="show:lazy">
        <MeCrud
          ref="$table"
          v-model:query-items="queryItems"
          collapse-query
          :columns="columns"
          :get-data="api.cert.list"
        >
          <MeQueryItem label="证书名称" :label-width="80">
            <n-input
              v-model:value="queryItems.name"
              placeholder="请输入证书名称"
              clearable
            />
          </MeQueryItem>
        </MeCrud>
      </NTabPane>
      <NTabPane name="issued" tab="已颁发" display-directive="show:lazy">
        <MeCrud
          ref="$issuedTable"
          v-model:query-items="issuedQuery"
          collapse-query
          :columns="issuedColumns"
          :get-data="api.cert.issuedList"
        >
          <CityQueryItem v-model="issuedQuery.city_id" />
          <MeQueryItem label="关键词" :label-width="80">
            <n-input
              v-model:value="issuedQuery.name"
              placeholder="学员 / 证书名 / 编号"
              clearable
            />
          </MeQueryItem>
        </MeCrud>
      </NTabPane>
    </DetailTabs>

    <!-- 查看证书弹窗：仅显示图片 + 占位框 -->
    <NModal
      v-model:show="showViewModal"
      preset="card"
      class="cert-view-modal"
      :show-icon="false"
      :bordered="false"
      :closable="true"
      @after-leave="closeViewModal"
    >
      <div v-if="viewLoading" class="cert-view-loading">
        加载中...
      </div>
      <div
        v-else-if="viewCert?.preview_image || viewCert?.template_image"
        class="cert-preview"
        :style="{ width: `${PREVIEW_W}px`, height: `${PREVIEW_H}px` }"
      >
        <img
          :src="getAssetUrl(viewCert.preview_image || viewCert.template_image)"
          class="cert-preview__bg"
          alt="证书"
        >
        <template v-if="!viewCert.preview_image">
          <div
            v-for="(p, i) in (viewCert.params || [])"
            :key="i"
            class="cert-preview__box"
            :style="{
              left: `${p.x * PREVIEW_SCALE}px`,
              top: `${p.y * PREVIEW_SCALE}px`,
              width: `${p.w * PREVIEW_SCALE}px`,
              height: `${p.h * PREVIEW_SCALE}px`,
              color: p.color || '#333',
            }"
          >
            {{ getParamPreview(p) }}
          </div>
        </template>
      </div>
    </NModal>
  </CommonPage>
</template>

<style scoped>
.cert-view-modal {
  width: auto;
  max-width: 95vw;
  padding: 0;
}

.cert-view-modal :deep(.n-card) {
  border-radius: 16px !important;
  overflow: hidden;
}

.cert-view-modal :deep(.n-card__content) {
  padding: 12px !important;
}

.cert-view-loading {
  padding: 28px 36px;
  text-align: center;
  color: #8b97a5;
  font-size: 13px;
  letter-spacing: -0.01em;
}

.cert-preview {
  position: relative;
  border-radius: 14px;
  border: 1px solid var(--ctl-hairline, rgba(var(--primary-color), 0.08));
  overflow: hidden;
  flex-shrink: 0;
  box-shadow: inset 0 1px 0 rgb(255 255 255 / 55%);
}

.cert-preview__bg {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
  object-fit: cover;
  pointer-events: none;
}

.cert-preview__box {
  position: absolute;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 10px;
  font-weight: 550;
  letter-spacing: -0.01em;
  color: var(--ctl-text, #1a2332);
  background: var(--ctl-glass-strong, rgb(255 255 255 / 82%));
  border: 1px dashed rgba(var(--primary-color), 0.35);
  border-radius: 6px;
  box-sizing: border-box;
  backdrop-filter: none;
  -webkit-backdrop-filter: none;
}
</style>
