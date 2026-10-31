<template>
  <CommonPage title="测评记录管理">
    <MeCrud
      ref="$table"
      :query-items="queryItems"
      collapse-query
      :columns="columns"
      :get-data="api.paper.assessmentList"
    >
      <CityQueryItem v-model="queryItems.city_id" :label-width="60" />
      <MeQueryItem label="学员" :label-width="60">
        <NInput
          v-model:value="queryItems.nick_name"
          clearable
          placeholder="搜索学员昵称"
          class="w-180"
          @keyup.enter="$table?.handleSearch(false)"
        />
      </MeQueryItem>
      <MeQueryItem label="课程" :label-width="60">
        <NSelect
          v-model:value="queryItems.course_id"
          :options="courseOptions"
          :loading="courseLoading"
          filterable
          remote
          clearable
          placeholder="搜索课程名称"
          class="w-180"
          @search="searchCourses"
          @update:value="onCourseChange"
        />
      </MeQueryItem>
      <MeQueryItem v-if="queryItems.course_id" label="课节" :label-width="60">
        <NSelect
          v-model:value="queryItems.chapter_id"
          :options="chapterOptions"
          :loading="chapterLoading"
          filterable
          clearable
          placeholder="选择课节"
          class="w-180"
          @update:value="$table?.handleSearch(false)"
        />
      </MeQueryItem>
      <MeQueryItem label="是否通过" :label-width="80">
        <NSelect
          v-model:value="queryItems.is_pass"
          :options="passOptions"
          placeholder="是否通过"
          clearable
          class="w-120"
          @update:value="$table?.handleSearch(false)"
        />
      </MeQueryItem>
      <MeQueryItem label="是否批改" :label-width="80">
        <NSelect
          v-model:value="queryItems.is_people_check"
          :options="gradedOptions"
          placeholder="是否批改"
          clearable
          class="w-120"
          @update:value="$table?.handleSearch(false)"
        />
      </MeQueryItem>
      <template #card="{ list }">
        <PaperRecordCards
          :list="list"
          show-chapter
          @open="openGrade"
          @message="openSendMessage"
        />
      </template>
    </MeCrud>
    <PaperGradeModal ref="$gradeModal" @success="onGradeSuccess" />
    <SendMessageModal ref="$msgModal" />
  </CommonPage>
</template>

<script setup>
import { getAssetUrl } from '@/utils'
import { NAvatar, NButton, NInput, NSelect, NTag } from 'naive-ui'
import { MeCrud, MeQueryItem, CityQueryItem } from '@/components'
import { CommonPage, PaperGradeModal, PaperRecordCards, SendMessageModal } from '@/components/index.js'
import { useRefreshOnActivate } from '@/composables'
import api from '@/api'
import { h, nextTick, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'

defineOptions({ name: 'PaperAssessment' })

const router = useRouter()
const $table = ref(null)
const $gradeModal = ref(null)
const $msgModal = ref(null)
const queryItems = ref({
  nick_name: '',
  course_id: null,
  chapter_id: null,
  city_id: null,
  is_pass: null,
  is_people_check: null,
})

const courseOptions = ref([])
const courseLoading = ref(false)
const chapterOptions = ref([])
const chapterLoading = ref(false)
let courseSearchTimer = null
const passOptions = [
  { label: '通过', value: 1 },
  { label: '未通过', value: 0 },
]
const gradedOptions = [
  { label: '已批改', value: 0 },
  { label: '待批改', value: 1 },
]

function mapCourseOptions(list) {
  return (list || []).map(c => ({ label: c.title, value: c.id }))
}

function searchCourses(query) {
  if (courseSearchTimer)
    clearTimeout(courseSearchTimer)
  courseSearchTimer = setTimeout(() => {
    courseLoading.value = true
    api.course.courseList({ title: query || '', size: 20, page: 1 }).then((res) => {
      const selectedId = queryItems.value.course_id
      const selected = courseOptions.value.find(o => o.value === selectedId)
      courseOptions.value = mapCourseOptions(res.data?.list)
      if (selected && !courseOptions.value.some(o => o.value === selected.value))
        courseOptions.value.unshift(selected)
    }).catch((error) => {
      $message.error(error?.message || '搜索课程失败')
      courseOptions.value = []
    }).finally(() => {
      courseLoading.value = false
    })
  }, 300)
}

function loadChapters(courseId) {
  if (!courseId) {
    chapterOptions.value = []
    return
  }
  chapterLoading.value = true
  api.chapter.chapterList({ courseId }).then((res) => {
    chapterOptions.value = (res.data || []).map(item => ({ label: item.title, value: item.id }))
  }).catch((error) => {
    $message.error(error?.message || '加载课节失败')
    chapterOptions.value = []
  }).finally(() => {
    chapterLoading.value = false
  })
}

function onCourseChange() {
  queryItems.value.chapter_id = null
  loadChapters(queryItems.value.course_id)
  $table.value?.handleSearch(false)
}

onMounted(() => {
  searchCourses('')
  $table.value?.handleSearch()
})

useRefreshOnActivate(() => $table.value?.handleSearch())

const columns = [
  {
    title: '学员',
    key: 'user',
    width: 180,
    render({ user, user_id }) {
      return h(
        'div',
        {
          class: 'flex cursor-pointer items-center',
          onClick: () => router.push({ name: 'UserDetail', params: { id: user_id } }),
        },
        [
          h(NAvatar, { size: 40, src: getAssetUrl(user?.avatar) || undefined, class: 'flex-shrink-0' }),
          h('div', { class: 'ml-12' }, [
            h('div', { class: 'font-medium' }, user?.nick_name || '-'),
            h('div', { class: 'text-12 text-secondary' }, `${user?.credit || 0}学分`),
          ]),
        ],
      )
    },
  },
  {
    title: '课程',
    key: 'course_name',
    width: 150,
    ellipsis: { tooltip: true },
    render: ({ paper }) => {
      const courseId = paper?.course_id || paper?.course?.id
      const title = paper?.course?.title || '-'
      if (!courseId) return title
      return h(
        'span',
        { class: 'cursor-pointer text-primary hover:underline', onClick: () => router.push(`/course/detail/${courseId}`) },
        title,
      )
    },
  },
  {
    title: '课节',
    key: 'chapter_name',
    width: 150,
    ellipsis: { tooltip: true },
    render: ({ paper }) => {
      const courseId = paper?.course_id || paper?.course?.id
      const chapterId = paper?.chapter_id || paper?.chapter?.id
      const title = paper?.chapter?.title || '-'
      if (!courseId || !chapterId) return title
      return h(
        'span',
        { class: 'cursor-pointer text-primary hover:underline', onClick: () => router.push(`/course/${courseId}/chapterDetail/${chapterId}`) },
        title,
      )
    },
  },
  { title: '成绩', key: 'score', width: 80, align: 'center' },
  {
    title: '是否通过',
    key: 'is_pass',
    width: 100,
    align: 'center',
    render: ({ is_pass }) => h(NTag, { type: is_pass ? 'success' : 'error', size: 'small', round: true, bordered: false }, is_pass ? '通过' : '未通过'),
  },
  {
    title: '是否批改',
    key: 'is_people_check',
    width: 100,
    align: 'center',
    render: ({ is_people_check }) => h(NTag, { type: is_people_check ? 'warning' : 'info', size: 'small', round: true, bordered: false }, is_people_check ? '待批改' : '已批改'),
  },
  { title: '提交次数', key: 'submit_count', width: 90, align: 'center' },
  { title: '提交时间', key: 'created_at', width: 170 },
  {
    title: '操作',
    key: 'actions',
    width: 180,
    align: 'center',
    fixed: 'right',
    render(row) {
      const isGraded = !row.is_people_check
      return h('div', { style: { display: 'flex', gap: '8px', justifyContent: 'center', flexWrap: 'wrap' } }, [
        h(
          NButton,
          {
            size: 'small',
            type: isGraded ? 'default' : 'primary',
            onClick: () => openGrade(row, isGraded),
          },
          { default: () => isGraded ? '查看' : '批改' },
        ),
        h(
          NButton,
          {
            size: 'small',
            type: 'info',
            secondary: true,
            onClick: () => openSendMessage(row),
          },
          { default: () => '站内信' },
        ),
      ])
    },
  },
]

function pendingQueue(exceptId) {
  return ($table.value?.tableData || []).filter(row =>
    Number(row?.is_people_check) === 1 && row.id !== exceptId,
  )
}

function openGrade(row, viewOnly = false) {
  $gradeModal.value?.open(row.id, viewOnly, {
    hasNext: !viewOnly && pendingQueue(row.id).length > 0,
  })
}

function onGradeSuccess(payload = {}) {
  const nextRow = payload.next ? pendingQueue(payload.id)[0] : null
  $table.value?.handleSearch()
  if (nextRow)
    nextTick(() => openGrade(nextRow, false))
}

function openSendMessage(row) {
  const courseTitle = row.paper?.course?.title || '课程'
  const chapterTitle = row.paper?.chapter?.title
  const target = chapterTitle ? `「${courseTitle}」课节「${chapterTitle}」` : `「${courseTitle}」`
  $msgModal.value?.open({
    userId: row.user_id,
    userName: row.user?.nick_name,
    type: 'paper',
    title: '作业批改通知',
    content: `同学，你提交的${target}作业已批改${row.score != null ? `，成绩 ${row.score} 分` : ''}。`,
    relatedType: 'user_paper',
    relatedId: row.id,
  })
}
</script>
