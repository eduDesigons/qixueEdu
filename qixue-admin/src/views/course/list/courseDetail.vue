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
import { CommonPage, DetailTabs, Upload } from '@/components/index.js'
import { useTabQuery, peekTabQuery } from '@/composables'
import { formatRichContentHtml } from '@/utils'
import CourseCommentTab from '@/views/course/list/modules/CourseCommentTab.vue'
import PaperView from '@/views/course/list/modules/PaperView.vue'
import ChaptersTab from './modules/ChaptersTab.vue'
import CourseAttachTab from './modules/CourseAttachTab.vue'
import CourseCertTab from './modules/CourseCertTab.vue'
import CourseInfo from './modules/CourseInfo.vue'
import CourseUserApplyTab from './modules/CourseUserApplyTab.vue'
import CourseUserTab from './modules/CourseUserTab.vue'
import PaperLogsTab from './modules/PaperLogsTab.vue'

defineOptions({
  name: 'CourseDetail',
})

const router = useRouter()
const route = useRoute()

const currentTabName = ref(peekTabQuery('content'))
useTabQuery(currentTabName)

const courseId = ref(route.params.id)

// 课件上传落到资源库「文档」分类（避免写死不存在的 cate-id）
const attachCateId = ref(0)
async function resolveAttachCateId() {
  try {
    const res = await api.mediaCateList('document')
    const list = res.data || []
    const preferred = list.find(item => Number(item.is_init) === 1) || list[0]
    attachCateId.value = Number(preferred?.id) || 0
  }
  catch {
    attachCateId.value = 0
  }
}

const paperIdForLogs = computed(() => route.query.paperId || courseData.value?.paper_id)
const pageTitle = computed(() => courseData.value?.title || '课程详情')

const courseData = ref({})
const detailLoading = ref(true)

function fetchCourseDetail() {
  if (!courseId.value) {
    $message.error('课程不存在')
    router.replace('/course/list')
    return
  }
  if (!courseData.value?.id)
    detailLoading.value = true
  api.course.getCourseDetail(courseId.value).then((res) => {
    courseData.value = res.data || {}
  }).catch((error) => {
    $message.error(error?.message || '加载课程失败')
  }).finally(() => {
    detailLoading.value = false
  })
}

onMounted(() => {
  resolveAttachCateId()
  fetchCourseDetail()
})

const $chaptersTab = ref(null)
const $paperView = ref(null)
const $paperLogsTab = ref(null)
const $courseCommentTab = ref(null)

onActivated(() => {
  if (courseData.value?.id)
    fetchCourseDetail()
  $chaptersTab.value?.refresh?.()
  $courseUserTab.value?.refresh?.()
  $courseCertTab.value?.refresh?.()
  $courseAttachTab.value?.getList?.()
  $paperView.value?.reload?.()
  $paperLogsTab.value?.refresh?.()
  $courseCommentTab.value?.refresh?.()
})

function handleEdit() {
  router.push({
    path: `/course/list/editCourse/${courseId.value}`,
  })
}

function handleDelete() {
  $dialog.warning({
    title: '提示',
    content: '确定删除该课程吗？',
    positiveText: '确定',
    negativeText: '取消',
    async onPositiveClick() {
      try {
        await api.course.courseDelete(courseId.value)
        window.$message.success('课程删除成功')
        router.push('/course/list')
      }
      catch (error) {
        $message.error(error?.message || '课程删除失败')
      }
    },
  })
}
function handleAddChapter() {
  router.push({
    path: `/course/${courseId.value}/chapterAdd`,
    query: {
      type: courseData.value?.type,
    },
  })
}

function editPaper() {
  router.push({
    path: courseData.value.paper_id ? `/course/${courseId.value}/paperEdit/${courseData.value.paper_id}` : `/course/${courseId.value}/paperAdd`,
  })
}

const $courseAttachTab = ref(null)
function handleUploadCourseAttach(res) {
  api.course.courseAttachSave({
    course_id: courseId.value,
    name: res.title || res.name,
    path: res.url,
    size: res.size,
  }).then(() => {
    $courseAttachTab.value?.getList()
  }).catch((error) => {
    $message.error(error?.message || '课件保存失败')
  })
}

const $courseUserTab = ref(null)
const $courseCertTab = ref(null)
function handleAddUser() {
  $courseUserTab.value?.showAddUserModal()
}

function handleCertBackfill() {
  $courseCertTab.value?.askBackfill?.()
}

const moreOptions = [
  { label: '删除课程', key: 'delete' },
]

function handleMoreSelect(key) {
  if (key === 'delete')
    handleDelete()
}
</script>

<template>
  <CommonPage back fill :title="pageTitle">
    <NSpin :show="detailLoading">
      <CourseInfo :course-data="courseData" @refresh="fetchCourseDetail" />
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
          <NButton v-else-if="currentTabName === 'chapters'" type="primary" @click="handleAddChapter">
            添加课节
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
            <template v-if="Number(courseData.pending_cert_users) > 0">
              （{{ courseData.pending_cert_users }}）
            </template>
          </NButton>
          <Upload v-else-if="currentTabName === 'courseAttach'" type="document" :cate-id="attachCateId" :show-preview="false" multiple @success="handleUploadCourseAttach">
            <NButton type="primary">
              上传课件
            </NButton>
          </Upload>
          <NButton v-else-if="currentTabName === 'paperView'" type="primary" @click="editPaper">
            {{ courseData.paper_id ? '编辑考试' : '新增考试' }}
          </NButton>
      </template>
      <NTabPane name="content" tab="介绍内容" display-directive="show:lazy">
        <div class="rich-content-panel">
          <div class="rich-content" v-html="formatRichContentHtml(courseData.content)" />
        </div>
      </NTabPane>
      <NTabPane name="chapters" tab="课节列表" display-directive="show:lazy">
        <ChaptersTab :key="`chapters-${courseId}`" ref="$chaptersTab" :course-id="courseId" />
      </NTabPane>
      <NTabPane name="paperView" tab="试卷内容" display-directive="show:lazy">
        <PaperView ref="$paperView" :paper-id="courseData.paper_id" />
      </NTabPane>
      <NTabPane name="paperLogs" tab="考试记录" display-directive="show:lazy">
        <PaperLogsTab v-if="paperIdForLogs" :key="`paper-logs-${paperIdForLogs}`" ref="$paperLogsTab" :paper-id="paperIdForLogs" />
        <div v-else class="empty-state detail-empty">暂无试卷记录</div>
      </NTabPane>
      <NTabPane name="courseAttach" tab="课件列表" display-directive="show:lazy">
        <CourseAttachTab :key="`attach-${courseId}`" ref="$courseAttachTab" :course-id="courseId" />
      </NTabPane>
      <NTabPane name="courseUser" tab="学员列表" display-directive="show:lazy">
        <CourseUserTab
          :key="`course-user-${courseId}`"
          ref="$courseUserTab"
          :course-id="courseId"
          :has-cert="!!courseData.cert_id"
          @user-change="fetchCourseDetail()"
        />
      </NTabPane>
      <NTabPane
        v-if="courseData.cert_id"
        name="courseCert"
        tab="证书颁发"
        display-directive="show:lazy"
      >
        <CourseCertTab
          :key="`course-cert-${courseId}`"
          ref="$courseCertTab"
          :course-id="courseId"
          :pending-count="courseData.pending_cert_users"
          @change="fetchCourseDetail()"
        />
      </NTabPane>
      <NTabPane name="applyCourse" tab="入课申请" display-directive="show:lazy">
        <CourseUserApplyTab :key="`course-apply-${courseId}`" :course-id="courseId" @user-change="fetchCourseDetail()" />
      </NTabPane>
      <NTabPane name="comments" tab="课程评论" display-directive="show:lazy">
        <CourseCommentTab :key="`course-comment-${courseId}`" ref="$courseCommentTab" :course-id="courseId" @user-change="fetchCourseDetail()" />
      </NTabPane>
    </DetailTabs>
  </CommonPage>
</template>
