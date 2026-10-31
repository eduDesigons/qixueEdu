<script setup>
import { getAssetUrl } from '@/utils'
import { NAvatar, NTag } from 'naive-ui'
import { h } from 'vue'
import { useRouter } from 'vue-router'
import api from '@/api/index.js'
import { MeCrud, MeQueryItem } from '@/components'
import { formatDateTime } from '@/utils/index.js'

const props = defineProps({
  userId: {
    type: [String, Number],
    required: true,
  },
})

const router = useRouter()
const queryItems = ref({
  userId: props.userId,
  title: '',
})

const $table = ref(null)
onMounted(() => {
  $table.value?.handleSearch()
})

function handleViewPath(row) {
  router.push({
    name: 'CourseDetail',
    params: {
      id: row.course_id,
    },
  })
}

const columns = [
  {
    title: '课程',
    key: 'course-info',
    width: 300,
    render(row) {
      return h(
        'div',
        {
          class: 'flex cursor-pointer',
          onClick: () => handleViewPath(row),
        },
        [
          h(NAvatar, {
            size: 60,
            src: getAssetUrl(row.course.thumb) || undefined,
            class: 'flex-shrink-0',
          }),
          h(
            'div',
            {
              class: 'ml-12 flex-col justify-center',
            },
            [
              h(
                'div',
                {
                  class: 'text-16 font-bold text-primary',
                },
                row.course.title,
              ),
              h(
                'div',
                {
                  class: 'text-12 text-secondary mt-4',
                },
                row.course.desc,
              ),
            ],
          ),
        ],
      )
    },
  },
  {
    title: '分类',
    key: 'category',
    align: 'center',
    width: 120,
    render: ({ course }) => h(NTag, { type: 'info', size: 'small', round: true, bordered: false }, course.category?.name ?? ''),
  },

  {
    title: '总课节',
    key: 'courses_count',
    align: 'center',
    width: 100,
    render: ({ course }) => {
      return h('span', { class: 'text font-bold' }, course.chapters_count)
    },
  },
  {
    title: '已学完',
    key: 'courses_count',
    align: 'center',
    width: 100,
    render: ({ has_end_chapter }) => {
      return h('span', { class: 'text font-bold' }, has_end_chapter)
    },
  },
  {
    title: '状态',
    key: 'status',
    align: 'center',
    width: 160,
    render(row) {
      return h(NTag, {
        type: row.is_show ? 'success' : 'error',
        size: 'small',
        round: true,
        bordered: false,
      }, row.is_end ? '已学完' : '未学完')
    },
  },
  {
    title: '添加时间',
    key: 'created_at',
    align: 'center',
    width: 180,
    render(row) {
      return h('span', formatDateTime(row.created_at))
    },
  },
]
</script>

<template>
  <div class="h-500">
    <MeCrud
      ref="$table"
      v-model:query-items="queryItems"
      :bound-query="{ userId }"
      :columns="columns"
      :get-data="api.user.getUserCourseList"
    >
      <MeQueryItem label="课程" :label-width="50">
        <n-input v-model:value="queryItems.title" placeholder="课程名称" clearable />
      </MeQueryItem>
    </MeCrud>
  </div>
</template>
