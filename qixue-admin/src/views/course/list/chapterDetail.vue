<script setup>
import {
  NButton,
  NDropdown,
  NSpin,
  NTabPane,
} from 'naive-ui'
import { computed, onActivated, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import api from '@/api'
import { CommonPage, DetailTabs } from '@/components/index.js'
import { useTabQuery, peekTabQuery } from '@/composables'
import { formatRichContentHtml } from '@/utils'
import { useTabStore } from '@/store/index.js'
import ChapterCommentTab from './modules/ChapterCommentTab.vue'
import ChapterInfo from './modules/ChapterInfo.vue'
import ChapterUserTab from './modules/ChapterUserTab.vue'
import PaperLogsTab from './modules/PaperLogsTab.vue'
import PaperView from './modules/PaperView.vue'

defineOptions({
  name: 'ChapterDetail',
})

const router = useRouter()
const route = useRoute()

const chapterId = ref(route.params.id)
const courseId = ref(route.params.courseId)

const chapterData = ref({})
const detailLoading = ref(true)

function fetchChapterDetail() {
  if (!chapterId.value) {
    $message.error('课节不存在')
    router.replace('/course/list')
    return
  }
  if (!chapterData.value?.id)
    detailLoading.value = true
  api.chapter.getChapterDetail(chapterId.value).then((res) => {
    chapterData.value = res.data || {}
  }).catch((error) => {
    $message.error(error?.message || '加载课节失败')
  }).finally(() => {
    detailLoading.value = false
  })
}

watch(
  () => [route.params.id, route.params.courseId],
  ([id, cid]) => {
    if (!id || String(id) === String(chapterId.value))
      return
    chapterId.value = id
    courseId.value = cid
    chapterData.value = {}
    fetchChapterDetail()
  },
)

onMounted(() => {
  fetchChapterDetail()
})

const $paperView = ref(null)
const $paperLogsTab = ref(null)
const $chapterUserTab = ref(null)
const $chapterCommentTab = ref(null)

onActivated(() => {
  if (chapterData.value?.id)
    fetchChapterDetail()
  $paperView.value?.reload?.()
  $paperLogsTab.value?.refresh?.()
  $chapterUserTab.value?.refresh?.()
  $chapterCommentTab.value?.refresh?.()
})

const currentTabName = ref(peekTabQuery('content'))
useTabQuery(currentTabName)
const pageTitle = computed(() => chapterData.value?.title || '课节详情')

const moreOptions = [
  { label: '删除课节', key: 'delete' },
]

function handleMoreSelect(key) {
  if (key === 'delete')
    handleDelete()
}

function handleEdit() {
  router.push({
    path: `/course/${courseId.value}/chapterEdit/${chapterId.value}`,
    query: {
      type: chapterData.value.type,
    },
  })
}

const tabStore = useTabStore()
function handleDelete() {
  $dialog.warning({
    title: '提示',
    content: '确定删除该课节吗？',
    positiveText: '确定',
    negativeText: '取消',
    async onPositiveClick() {
      try {
        await api.chapter.chapterDelete(chapterId.value)
        $message.success('课节删除成功')
        tabStore.removeTab(route.path)
        router.push('/course/list')
      }
      catch (error) {
        $message.error(error?.message || '课节删除失败')
      }
    },
  })
}

function editPaper() {
  router.push({
    path: chapterData.value.paper_id ? `/course/${courseId.value}/paperEdit/${chapterData.value.paper_id}` : `/course/${courseId.value}/paperAdd`,
    query: {
      chapterId: chapterId.value,
    },
  })
}
</script>

<template>
  <CommonPage back fill :title="pageTitle">
    <NSpin :show="detailLoading">
      <ChapterInfo :info="chapterData" />
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
          <NButton v-else-if="currentTabName === 'papers'" type="primary" @click="editPaper">
            {{ chapterData.paper_id ? '编辑测评' : '新增测评' }}
          </NButton>
      </template>
      <NTabPane name="content" tab="课节内容" display-directive="show:lazy">
        <div class="rich-content-panel">
          <div class="rich-content" v-html="formatRichContentHtml(chapterData.content)" />
        </div>
      </NTabPane>
      <NTabPane name="papers" tab="测评内容" display-directive="show:lazy">
        <PaperView ref="$paperView" :paper-id="chapterData.paper_id" />
      </NTabPane>
      <NTabPane name="evaluations" tab="测评记录" display-directive="show:lazy">
        <PaperLogsTab v-if="chapterData.paper_id" ref="$paperLogsTab" :paper-id="chapterData.paper_id" />
        <div v-else class="empty-state detail-empty">
          暂无测评记录
        </div>
      </NTabPane>
      <NTabPane name="members" tab="学员列表" display-directive="show:lazy">
        <ChapterUserTab
          v-if="chapterId"
          :key="`chapter-user-${chapterId}`"
          ref="$chapterUserTab"
          :chapter-id="chapterId"
        />
      </NTabPane>
      <NTabPane name="comments" tab="课节评论" display-directive="show:lazy">
        <ChapterCommentTab
          v-if="chapterId"
          :key="`chapter-comment-${chapterId}`"
          ref="$chapterCommentTab"
          :chapter-id="chapterId"
          :course-id="chapterData.course_id || courseId"
        />
      </NTabPane>
    </DetailTabs>
  </CommonPage>
</template>
