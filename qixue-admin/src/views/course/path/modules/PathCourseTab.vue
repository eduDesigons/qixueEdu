<script setup>
import { getAssetUrl } from '@/utils'
import { NAvatar, NButton, NCascader, NDataTable, NInput, NTag } from 'naive-ui'
import { h } from 'vue'
import { useRouter } from 'vue-router'
import api from '@/api/index.js'
import { MeCrud, MeModal } from '@/components'
import { useCrud } from '@/composables/index.js'
import { getCourseTypeName } from '@/utils/dict.js'
import { formatDateTime } from '@/utils/index.js'

const props = defineProps({
  pathId: {
    type: [String, Number],
    required: true,
  },
})
const emit = defineEmits(['userChange'])
const router = useRouter()
function handleViewCourse(courseId) {
  router.push(`/course/detail/${courseId}`)
}
const $table = ref(null)
onMounted(() => {
  $table.value?.handleSearch()
})

const {
  handleDelete,
} = useCrud({
  name: '课程',
  doDelete: api.path.pathCourseDelete,
  refresh: () => $table.value?.handleSearch(),
})

const columns = [
  {
    title: '编号',
    key: 'id',
    width: 80,
  },
  {
    title: '课程',
    key: 'course-info',
    width: 300,
    render({ course }) {
      return h(
        'div',
        {
          class: 'flex cursor-pointer',
          onClick: () => handleViewCourse(course.id),
        },
        [
          h(NAvatar, {
            size: 60,
            src: getAssetUrl(course.thumb) || undefined,
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
                course.title,
              ),
              h(
                'div',
                {
                  class: 'text-12 text-secondary mt-4',
                },
                course.desc,
              ),
            ],
          ),
        ],
      )
    },
  },
  {
    title: '类型',
    key: 'type',
    align: 'center',
    width: 120,
    render: ({ course }) => h(NTag, { type: 'success', size: 'small', round: true, bordered: false }, getCourseTypeName(course.type)),
  },
  {
    title: '分类',
    key: 'category',
    align: 'center',
    width: 120,
    render: ({ course }) => h(NTag, { type: 'info', size: 'small', round: true, bordered: false }, course.category?.name ?? ''),
  },

  {
    title: '课节',
    key: 'chapters_count',
    align: 'center',
    width: 100,
    render: ({ course }) => {
      return h('span', { class: 'text font-bold' }, course.chapters_count)
    },
  },
  {
    title: '可得学分',
    key: 'credit',
    align: 'center',
    width: 100,
    render: ({ course }) => {
      if (course.credit) {
        return h('span', { class: 'text-gray font-bold' }, course.credit)
      }
      return h('span', { class: 'text-gray font-bold' }, '无学分')
    },
  },
  {
    title: '限制学分',
    key: 'min_need_credit',
    align: 'center',
    width: 100,
    render: ({ course }) => {
      if (course.min_need_credit) {
        return h('span', { class: 'text-gray font-bold' }, course.min_need_credit)
      }
      return h('span', { class: 'text-gray font-bold' }, '不限制')
    },
  },
  {
    title: '学习人数',
    key: 'user_count',
    align: 'center',
    width: 100,
    render: ({ course }) => course.users_count || 0,
  },
  {
    title: '状态',
    key: 'status',
    align: 'center',
    width: 160,
    render({ course }) {
      return h('div', [
        course.is_rec ? h(NTag, { type: 'primary', size: 'small', round: true, bordered: false }, '推荐') : null,
        h(NTag, {
          type: course.is_show ? 'success' : 'error',
          size: 'small',
          round: true,
          bordered: false,
        }, course.is_show ? '已发布' : '未发布'),
      ])
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
  {
    title: '操作',
    key: 'actions',
    width: 120,
    align: 'center',
    fixed: 'right',
    hideInExcel: true,
    render(row) {
      return [
        h(
          NButton,
          {
            size: 'small',
            type: 'error',
            style: 'margin-left: 12px;',
            onClick: () => handleDelete(row.id),
          },
          {
            default: () => '移除',
            icon: () => h('i', { class: 'i-material-symbols:delete-outline text-14' }),
          },
        ),
      ]
    },
  },
]

const checkAddCourseList = ref([])
const courseRows = ref([])
const searching = ref(false)
const searchKey = ref('')
const categoryId = ref(null)
const courseCategories = ref([])
const $modalRef = ref(null)
const pickPagination = reactive({
  page: 1,
  pageSize: 10,
  itemCount: 0,
})

const pickColumns = [
  { type: 'selection' },
  {
    title: '课程',
    key: 'title',
    render(row) {
      return h('div', { class: 'flex items-center min-w-0' }, [
        h(NAvatar, {
          size: 40,
          src: getAssetUrl(row.thumb) || undefined,
          class: 'flex-shrink-0',
          style: { borderRadius: '8px' },
        }),
        h('span', { class: 'ml-10 truncate' }, row.title),
      ])
    },
  },
  {
    title: '分类',
    key: 'category',
    width: 140,
    render: row => row.category?.name || '-',
  },
  {
    title: '类型',
    key: 'type',
    width: 90,
    render: row => getCourseTypeName(row.type) || '-',
  },
]

function loadCourseCategories() {
  if (courseCategories.value.length)
    return
  api.course.getCourseCategories().then((res) => {
    courseCategories.value = res.data || []
  }).catch(() => {
    courseCategories.value = []
  })
}

function showAddCourseModal() {
  checkAddCourseList.value = []
  searchKey.value = ''
  categoryId.value = null
  courseRows.value = []
  pickPagination.page = 1
  pickPagination.itemCount = 0
  loadCourseCategories()
  $modalRef.value?.open()
  searchCourses(1)
}

function searchCourses(page = 1) {
  searching.value = true
  pickPagination.page = page
  api.path.pathCourseSearch({
    title: searchKey.value,
    cid: categoryId.value || '',
    path_id: props.pathId,
    page,
    size: pickPagination.pageSize,
  }).then((res) => {
    const data = res.data || {}
    courseRows.value = data.list || []
    pickPagination.itemCount = Number(data.total || 0)
    pickPagination.page = Number(data.current || page)
  }).catch((error) => {
    $message.error(error?.message || '搜索课程失败')
    courseRows.value = []
    pickPagination.itemCount = 0
  }).finally(() => {
    searching.value = false
  })
}

async function confirmAddCourse() {
  if (!checkAddCourseList.value.length) {
    $message.warning('请先勾选课程')
    return false
  }
  try {
    await api.path.pathCourseListSave({
      path_id: props.pathId,
      course_list: checkAddCourseList.value.map(id => Number(id)),
    })
  }
  catch {
    return false
  }
  $message.success('课程已添加')
  $table.value?.handleSearch()
  checkAddCourseList.value = []
  emit('userChange')
  return true
}

defineExpose({
  showAddCourseModal,
  refresh: () => $table.value?.handleSearch(false),
})
</script>

<template>
  <div class="h-500">
    <MeCrud
      ref="$table"
      :query-items="{ pathId: props.pathId }"
      :bound-query="{ pathId: pathId }"
      :is-pagination="false"
      :columns="columns"
      :get-data="api.path.pathCourse"
    />
    <MeModal ref="$modalRef" width="860px" title="添加课程" :on-ok="confirmAddCourse">
      <div class="add-course-query">
        <NCascader
          v-model:value="categoryId"
          :options="courseCategories"
          check-strategy="child"
          :show-path="true"
          filterable
          clearable
          value-field="id"
          label-field="name"
          placeholder="课程分类"
          style="width: 240px"
          @update:value="searchCourses(1)"
        />
        <NInput
          v-model:value="searchKey"
          clearable
          placeholder="课程名称"
          style="width: 240px"
          @keyup.enter="searchCourses(1)"
        />
        <NButton type="primary" secondary @click="searchCourses(1)">
          搜索
        </NButton>
        <span class="add-course-picked">已选 {{ checkAddCourseList.length }} 门</span>
      </div>
      <NDataTable
        v-model:checked-row-keys="checkAddCourseList"
        class="mt-12"
        remote
        :columns="pickColumns"
        :data="courseRows"
        :loading="searching"
        :row-key="row => row.id"
        :pagination="pickPagination"
        :max-height="420"
        @update:page="searchCourses"
      />
    </MeModal>
  </div>
</template>

<style scoped>
.add-course-query {
  display: flex;
  align-items: center;
  gap: 10px;
}

.add-course-picked {
  margin-left: auto;
  color: #8a9697;
  font-size: 13px;
}
</style>
