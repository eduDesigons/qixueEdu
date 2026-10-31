<script setup>
import { getAssetUrl } from '@/utils'
import { NAvatar, NButton, NInput, NTag } from 'naive-ui'
import { h } from 'vue'
import { useRouter } from 'vue-router'
import api from '@/api/index.js'
import { MeCrud, MeQueryItem, CityQueryItem, PaperGradeModal, PaperRecordCards } from '@/components'

const props = defineProps({
  paperId: {
    type: [String, Number],
    required: true,
  },
})

const queryItems = ref({
  paper_id: props.paperId,
  nick_name: '',
  city_id: null,
})

watch(() => props.paperId, (id) => {
  queryItems.value.paper_id = id
  $table.value?.handleSearch(false)
})

const router = useRouter()
const $table = ref(null)
const $gradeModal = ref(null)
onMounted(() => {
  $table.value?.handleSearch()
})

defineExpose({
  refresh: () => $table.value?.handleSearch(false),
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
          h(NAvatar, {
            size: 40,
            src: getAssetUrl(user?.avatar) || undefined,
            class: 'flex-shrink-0',
          }),
          h(
            'div',
            { class: 'ml-12' },
            [
              h('div', { class: 'font-medium' }, user?.nick_name || '-'),
              h(NTag, { type: 'info', size: 'small', round: true, bordered: false }, `${user?.credit ?? 0}学分`),
            ],
          ),
        ],
      )
    },
  },
  {
    title: '课程',
    key: 'course',
    width: 150,
    ellipsis: { tooltip: true },
    render: ({ paper }) => {
      const courseId = paper?.course_id || paper?.course?.id
      const title = paper?.course?.title ?? '-'
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
    key: 'chapter',
    width: 120,
    ellipsis: { tooltip: true },
    render: ({ paper }) => {
      const courseId = paper?.course_id || paper?.course?.id
      const chapterId = paper?.chapter_id || paper?.chapter?.id
      const title = paper?.chapter?.title ?? '-'
      if (!courseId || !chapterId) return title
      return h(
        'span',
        { class: 'cursor-pointer text-primary hover:underline', onClick: () => router.push(`/course/${courseId}/chapterDetail/${chapterId}`) },
        title,
      )
    },
  },
  {
    title: '成绩',
    key: 'score',
    align: 'center',
    render: row => h(
      'span',
      {
        class: 'text font-bold text-primary',
      },
      row.score,
    ),
  },
  {
    title: '是否通过',
    key: 'is_pass',
    align: 'center',
    width: 100,
    render: ({ is_pass }) => h(NTag, { type: is_pass ? 'success' : 'error', size: 'small', round: true, bordered: false }, is_pass ? '通过' : '未通过'),
  },
  {
    title: '是否批改',
    key: 'is_people_check',
    align: 'center',
    width: 100,
    render: ({ is_people_check }) => h(NTag, { type: is_people_check ? 'warning' : 'info', size: 'small', round: true, bordered: false }, is_people_check ? '待批改' : '已批改'),
  },
  {
    title: '提交次数',
    key: 'submit_count',
    align: 'center',
    width: 90,
  },
  {
    title: '提交时间',
    key: 'created_at',
    align: 'center',
    width: 170,
  },
  {
    title: '操作',
    key: 'actions',
    width: 100,
    align: 'center',
    fixed: 'right',
    render(row) {
      const isGraded = !row.is_people_check
      return h(
        NButton,
        {
          size: 'small',
          type: isGraded ? 'default' : 'primary',
          onClick: () => $gradeModal.value?.open(row.id, isGraded),
        },
        { default: () => isGraded ? '查看' : '批改' },
      )
    },
  },
]

function openGrade(row, viewOnly = false) {
  $gradeModal.value?.open(row.id, viewOnly)
}
</script>

<template>
  <div class="h-500">
    <MeCrud
      ref="$table"
      v-model:query-items="queryItems"
      :bound-query="{ paper_id: paperId }"
      :columns="columns"
      :get-data="api.paper.paperLogs"
    >
      <CityQueryItem v-model="queryItems.city_id" :label-width="50" />
      <MeQueryItem label="学员" :label-width="50">
        <NInput
          v-model:value="queryItems.nick_name"
          clearable
          placeholder="搜索学员昵称"
          class="w-180"
          @keyup.enter="$table?.handleSearch()"
        />
      </MeQueryItem>
      <template #card="{ list }">
        <PaperRecordCards
          :list="list"
          show-chapter
          :show-message="false"
          @open="openGrade"
        />
      </template>
    </MeCrud>
    <PaperGradeModal ref="$gradeModal" @success="$table?.handleSearch()" />
  </div>
</template>
