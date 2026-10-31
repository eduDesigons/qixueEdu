<script setup>
import { NButton, NSelect, NTag } from 'naive-ui'
import { h } from 'vue'
import { useRouter } from 'vue-router'
import api from '@/api'
import { MeCrud, MeQueryItem, PaperGradeModal, PaperRecordCards, SendMessageModal } from '@/components'
import { formatDateTime } from '@/utils'

const props = defineProps({
  userId: {
    type: [String, Number],
    required: true,
  },
})

const router = useRouter()
const queryItems = ref({ user_id: props.userId, type: 'course', is_pass: null, is_people_check: null })
const passOptions = [
  { label: '通过', value: 1 },
  { label: '未通过', value: 0 },
]
const gradedOptions = [
  { label: '已批改', value: 0 },
  { label: '待批改', value: 1 },
]
const $table = ref(null)
const $gradeModal = ref(null)
const $msgModal = ref(null)
onMounted(() => {
  $table.value?.handleSearch()
})

const columns = [
  {
    title: '课程',
    key: 'course',
    width: 220,
    ellipsis: { tooltip: true },
    render: ({ paper }) => {
      const courseId = paper?.course_id || paper?.course?.id
      const title = paper?.course?.title ?? '-'
      if (!paper || !courseId) return title
      return h(
        'span',
        { class: 'cursor-pointer text-primary hover:underline', onClick: () => router.push(`/course/detail/${courseId}`) },
        title,
      )
    },
  },
  {
    title: '得分',
    key: 'score',
    align: 'center',
    width: 80,
    render: ({ paper, score }) => (paper ? `${score ?? 0}/${paper.pass_score ?? 0}` : '-'),
  },
  {
    title: '是否通过',
    key: 'is_pass',
    align: 'center',
    width: 100,
    render: ({ is_pass }) => h(NTag, { type: is_pass === 1 ? 'success' : 'error', size: 'small', round: true, bordered: false }, is_pass === 1 ? '通过' : '未通过'),
  },
  {
    title: '提交时间',
    key: 'created_at',
    align: 'center',
    width: 180,
    render: row => h('span', formatDateTime(row.created_at)),
  },
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
            secondary: true,
            onClick: () => $gradeModal.value?.open(row.id, isGraded),
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

function openGrade(row, viewOnly = false) {
  $gradeModal.value?.open(row.id, viewOnly)
}

function openSendMessage(row) {
  const courseTitle = row.paper?.course?.title || '课程'
  $msgModal.value?.open({
    userId: row.user_id || props.userId,
    userName: row.user?.nick_name,
    type: 'paper',
    title: '作业批改通知',
    content: `同学，你提交的「${courseTitle}」作业已批改${row.score != null ? `，成绩 ${row.score} 分` : ''}。`,
    relatedType: 'user_paper',
    relatedId: row.id,
  })
}
</script>

<template>
  <div class="h-500">
    <MeCrud
      ref="$table"
      v-model:query-items="queryItems"
      :bound-query="{ user_id: userId }"
      :columns="columns"
      :get-data="api.user.getUserPaperList"
    >
      <MeQueryItem label="是否通过" :label-width="80">
        <NSelect
          v-model:value="queryItems.is_pass"
          :options="passOptions"
          placeholder="全部"
          clearable
        />
      </MeQueryItem>
      <MeQueryItem label="是否批改" :label-width="80">
        <NSelect
          v-model:value="queryItems.is_people_check"
          :options="gradedOptions"
          placeholder="全部"
          clearable
        />
      </MeQueryItem>
      <template #card="{ list }">
        <PaperRecordCards
          :list="list"
          :show-user="false"
          @open="openGrade"
          @message="openSendMessage"
        />
      </template>
    </MeCrud>
    <PaperGradeModal ref="$gradeModal" @success="$table?.handleSearch()" />
    <SendMessageModal ref="$msgModal" />
  </div>
</template>
