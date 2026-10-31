<script setup>
import {
  NButton,
  NDropdown,
  NSpin,
  NTabPane,
} from 'naive-ui'
import { computed, onActivated, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import api from '@/api'
import { CommonPage, DetailTabs } from '@/components/index.js'
import { useTabQuery, peekTabQuery } from '@/composables'
import { formatRichContentHtml } from '@/utils'
import PathCertTab from './modules/PathCertTab.vue'
import PathCourseTab from './modules/PathCourseTab.vue'
import PathInfo from './modules/PathInfo.vue'
import PathUserApplyTab from './modules/PathUserApplyTab.vue'
import PathUserTab from './modules/PathUserTab.vue'

defineOptions({
  name: 'CourseDetail',
})

const router = useRouter()
const route = useRoute()

const currentTabName = ref(peekTabQuery('content'))
useTabQuery(currentTabName)
const pageTitle = computed(() => pathData.value?.title || '系列详情')

const pathId = ref(route.params.id)

const pathData = ref({})
const detailLoading = ref(true)

function fetchPathDetail() {
  if (!pathId.value) {
    $message.error('系列不存在')
    router.replace('/course/path')
    return
  }
  if (!pathData.value?.id)
    detailLoading.value = true
  api.path.getPathDetail(pathId.value).then((res) => {
    pathData.value = res.data || {}
  }).catch((error) => {
    $message.error(error?.message || '加载系列失败')
  }).finally(() => {
    detailLoading.value = false
  })
}

onMounted(() => {
  fetchPathDetail()
})

onActivated(() => {
  if (pathData.value?.id)
    fetchPathDetail()
  $pathCourseTab.value?.refresh?.()
  $pathUserTab.value?.refresh?.()
  $pathCertTab.value?.refresh?.()
})

function handleEdit() {
  router.push({
    path: `/course/path/edit/${pathId.value}`,
  })
}

function handleDelete() {
  $dialog.warning({
    title: '提示',
    content: '确定删除该系列吗？',
    positiveText: '确定',
    negativeText: '取消',
    async onPositiveClick() {
      try {
        await api.path.pathDelete(pathId.value)
        $message.success('系列删除成功')
        router.push('/course/path')
      }
      catch (error) {
        $message.error(error?.message || '系列删除失败')
      }
    },
  })
}
const $pathCourseTab = ref(null)
function handleAddCourse() {
  $pathCourseTab.value?.showAddCourseModal()
}

const $pathUserTab = ref(null)
const $pathCertTab = ref(null)
function handleAddUser() {
  $pathUserTab.value?.showAddUserModal()
}

function handleCertBackfill() {
  $pathCertTab.value?.askBackfill?.()
}

const moreOptions = [
  { label: '删除系列', key: 'delete' },
]

function handleMoreSelect(key) {
  if (key === 'delete')
    handleDelete()
}
</script>

<template>
  <CommonPage back fill :title="pageTitle">
    <NSpin :show="detailLoading">
      <PathInfo :course-data="pathData" @refresh="fetchPathDetail" />
    </NSpin>

    <DetailTabs v-model:value="currentTabName">
      <template #suffix>
          <template v-if="currentTabName === 'content'">
            <NButton type="primary" @click="handleEdit">
              编辑
            </NButton>
            <NDropdown trigger="click" :options="moreOptions" @select="handleMoreSelect">
              <NButton secondary class="page-more-btn" aria-label="更多">
                <i class="i-fe:more-horizontal text-14" />
              </NButton>
            </NDropdown>
          </template>
          <NButton v-else-if="currentTabName === 'chapters'" type="primary" @click="handleAddCourse">
            添加课程
          </NButton>
          <NButton v-else-if="currentTabName === 'courseUser'" type="primary" @click="handleAddUser">
            添加学员
          </NButton>
          <NButton
            v-else-if="currentTabName === 'courseCert'"
            type="primary"
            @click="handleCertBackfill"
          >
            检测补发
            <template v-if="Number(pathData.pending_cert_users) > 0">
              （{{ pathData.pending_cert_users }}）
            </template>
          </NButton>
      </template>
      <NTabPane name="content" tab="介绍内容" display-directive="show:lazy">
        <div class="rich-content-panel">
          <div class="rich-content" v-html="formatRichContentHtml(pathData.content)" />
        </div>
      </NTabPane>
      <NTabPane name="chapters" tab="课程列表" display-directive="show:lazy">
        <PathCourseTab :key="`path-course-${pathId}`" ref="$pathCourseTab" :path-id="pathId" @user-change="fetchPathDetail()" />
      </NTabPane>
      <NTabPane name="courseUser" tab="学员列表" display-directive="show:lazy">
        <PathUserTab
          :key="`path-user-${pathId}`"
          ref="$pathUserTab"
          :path-id="pathId"
          :has-cert="!!pathData.cert_id"
          @user-change="fetchPathDetail()"
        />
      </NTabPane>
      <NTabPane
        v-if="pathData.cert_id"
        name="courseCert"
        tab="证书颁发"
        display-directive="show:lazy"
      >
        <PathCertTab
          :key="`path-cert-${pathId}`"
          ref="$pathCertTab"
          :path-id="pathId"
          :pending-count="pathData.pending_cert_users"
          @change="fetchPathDetail()"
        />
      </NTabPane>
      <NTabPane name="applyCourse" tab="入系申请" display-directive="show:lazy">
        <PathUserApplyTab :key="`path-apply-${pathId}`" :path-id="pathId" @user-change="fetchPathDetail()" />
      </NTabPane>
    </DetailTabs>
  </CommonPage>
</template>
