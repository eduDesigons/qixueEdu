<template>
  <CommonPage title="考试记录管理">
    <MeCrud
      ref="$table"
      :query-items="queryItems"
      collapse-query
      :columns="columns"
      :get-data="api.paper.examList"
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
import api from '@/api'
import { h, nextTick, onActivated, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'

defineOptions({ name: 'PaperExam' })

const router = useRouter()
const route = useRoute()
const $table = ref(null)
const $gradeModal = ref(null)
const $msgModal = ref(null)

function examQueryFromRoute() {
  if (route.query.pass === '1' || route.query.pass === '0') {
    return {
      is_pass: Number(route.query.pass),
      is_people_check: null,
    }
  }
  return {
    is_pass: null,
    is_people_check: null,
  }
}

const queryItems = ref({
  nick_name: '',
  course_id: null,
  city_id: null,
  ...examQueryFromRoute(),
})

const courseOptions = ref([])
const courseLoading = ref(false)
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

onMounted(() => {
  searchCourses('')
  $table.value?.handleSearch()
})

let skipActivateSearch = true
onActivated(() => {
  const next = examQueryFromRoute()
  const changed = queryItems.value.is_pass !== next.is_pass
    || queryItems.value.is_people_check !== next.is_people_check
  if (changed) {
    queryItems.value.is_pass = next.is_pass
    queryItems.value.is_people_check = next.is_people_check
  }
  if (skipActivateSearch) {
    skipActivateSearch = false
    return
  }
  $table.value?.handleSearch()
})

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
    width: 200,
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
  $msgModal.value?.open({
    userId: row.user_id,
    userName: row.user?.nick_name,
    type: 'paper',
    title: '作业批改通知',
    content: `同学，你提交的「${courseTitle}」作业已批改${row.score != null ? `，成绩 ${row.score} 分` : ''}。`,
    relatedType: 'user_paper',
    relatedId: row.id,
  })
}
</script>
