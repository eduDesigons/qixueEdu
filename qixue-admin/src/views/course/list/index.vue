<script setup>
import { NAvatar, NCascader, NTag } from 'naive-ui'
import { useRouter } from 'vue-router'
import api from '@/api'
import { MeCrud, MeQueryItem } from '@/components'
import { CommonPage } from '@/components/index.js'
import { useCrud, useRefreshOnActivate } from '@/composables'
import { getAssetUrl, formatDateTime, h } from '@/utils'
import { courseTypes, getCourseTypeName } from '@/utils/dict'

defineOptions({ name: 'CourseList' })

const router = useRouter()
const $table = ref(null)
const queryItems = ref({})

onMounted(() => {
  $table.value?.handleSearch()
})

useRefreshOnActivate(() => $table.value?.handleSearch())

const courseCategories = ref([])
function getCourseCategories() {
  api.course.getCourseCategories().then((res) => {
    courseCategories.value = res.data || []
  }).catch((error) => {
    $message.error(error?.message || '获取课程分类失败')
    courseCategories.value = []
  })
}
getCourseCategories()

useCrud({
  name: '课程',
  doDelete: api.course.courseDel,
  refresh: () => $table.value?.handleSearch(),
})

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
          onClick: () => handleViewCourse(row),
        },
        [
          h(NAvatar, {
            size: 60,
            src: getAssetUrl(row.thumb) || undefined,
            class: 'flex-shrink-0',
            style: { borderRadius: '12px' },
          }),
          h(
            'div',
            {
              class: 'ml-12 flex-col justify-center min-w-0',
            },
            [
              h('div', { class: 'text-14 font-bold text-primary truncate' }, row.title),
              h('div', { class: 'text-12 text-secondary mt-4 line-clamp-2' }, row.desc || '暂无简介'),
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
    render: ({ type }) => h(NTag, { type: 'success', size: 'small', round: true, bordered: false }, { default: () => getCourseTypeName(type) || '课程' }),
  },
  {
    title: '分类',
    key: 'category',
    align: 'center',
    width: 120,
    render: ({ category }) => h(NTag, { type: 'info', size: 'small', round: true, bordered: false }, { default: () => category?.name || '—' }),
  },
  {
    title: '课节',
    key: 'chapters_count',
    align: 'center',
    width: 100,
    render: ({ chapters_count }) => h('span', { class: 'text font-bold' }, chapters_count),
  },
  {
    title: '可得学分',
    key: 'credit',
    align: 'center',
    width: 100,
    render: ({ credit }) => credit
      ? h('span', { class: 'text-gray font-bold' }, credit)
      : h('span', { class: 'text-secondary' }, '无学分'),
  },
  {
    title: '限制学分',
    key: 'min_need_credit',
    align: 'center',
    width: 100,
    render: ({ min_need_credit }) => min_need_credit
      ? h('span', { class: 'text-gray font-bold' }, min_need_credit)
      : h('span', { class: 'text-secondary' }, '不限制'),
  },
  {
    title: '结业证书',
    key: 'cert',
    align: 'center',
    width: 140,
    ellipsis: { tooltip: true },
    render: ({ cert }) => cert?.name || h('span', { class: 'text-secondary' }, '未绑定'),
  },
  {
    title: '学习人数',
    key: 'users_count',
    align: 'center',
    width: 100,
    render: ({ users_count }) => users_count || 0,
  },
  {
    title: '状态',
    key: 'status',
    align: 'center',
    width: 180,
    render(row) {
      return h('div', { class: 'flex gap-6 justify-center flex-wrap' }, [
        row.is_rec ? h(NTag, { type: 'primary', bordered: false, size: 'small', round: true }, { default: () => '推荐' }) : null,
        row.is_open ? h(NTag, { type: 'info', bordered: false, size: 'small', round: true }, { default: () => '公开' }) : null,
        h(NTag, {
          type: row.is_show ? 'success' : 'error',
          bordered: false,
          size: 'small',
          round: true,
        }, { default: () => (row.is_show ? '已发布' : '未发布') }),
      ].filter(Boolean))
    },
  },
  {
    title: '创建时间',
    key: 'created_at',
    align: 'center',
    width: 180,
    render(row) {
      return h('span', formatDateTime(row.created_at))
    },
  },
]

function handleViewCourse(row) {
  router.push(`/course/detail/${row.id}`)
}

function handleAddCourse() {
  router.push('/course/list/courseAdd')
}
</script>

<template>
  <CommonPage>
    <template #action>
      <NButton type="primary" @click="handleAddCourse()">
        <i class="i-material-symbols:add mr-4 text-18" />
        新增课程
      </NButton>
    </template>

    <MeCrud
      ref="$table"
      v-model:query-items="queryItems"
      collapse-query
      :scroll-x="1000"
      :columns="columns"
      :get-data="api.course.courseList"
    >
      <MeQueryItem label="名称" :label-width="80">
        <n-input
          v-model:value="queryItems.title"
          type="text"
          placeholder="请输入课程名称"
          clearable
        />
      </MeQueryItem>
      <MeQueryItem label="分类" :label-width="80">
        <NCascader
          v-model:value="queryItems.category_id"
          :options="courseCategories"
          check-strategy="child"
          :show-path="true"
          :filterable="true"
          value-field="id"
          label-field="name"
        />
      </MeQueryItem>
      <MeQueryItem label="类型" :label-width="50">
        <n-select v-model:value="queryItems.type" clearable :options="courseTypes" />
      </MeQueryItem>
      <MeQueryItem label="状态" :label-width="50">
        <n-select
          v-model:value="queryItems.is_show"
          clearable
          :options="[
            { label: '已发布', value: 1 },
            { label: '未发布', value: 0 },
          ]"
        />
      </MeQueryItem>

      <template #card="{ list }">
        <div class="entity-grid">
          <article
            v-for="item in list"
            :key="item.id"
            class="entity-card entity-card--cover"
            :class="{ 'is-draft': !item.is_show }"
            @click="handleViewCourse(item)"
          >
            <div class="entity-card__cover">
              <img
                v-if="item.thumb"
                class="entity-card__cover-img"
                :src="getAssetUrl(item.thumb)"
                alt=""
              >
              <div v-else class="entity-card__cover-empty" />
              <div class="entity-card__badges entity-card__badges--on-cover">
                <span v-if="getCourseTypeName(item.type)" class="entity-card__badge">
                  {{ getCourseTypeName(item.type) }}
                </span>
                <span v-if="item.category?.name" class="entity-card__badge">
                  {{ item.category.name }}
                </span>
                <span v-if="item.is_rec" class="entity-card__badge is-rec">推荐</span>
                <span v-if="!item.is_show" class="entity-card__badge is-off">未发布</span>
              </div>
            </div>
            <div class="entity-card__body">
              <h3 class="entity-card__title" :title="item.title">
                {{ item.title }}
              </h3>
              <div class="entity-card__rels">
                <span class="entity-card__badge">{{ item.chapters_count || 0 }} 课节</span>
                <span v-if="item.credit" class="entity-card__badge">{{ item.credit }} 学分</span>
                <span class="entity-card__badge">{{ item.users_count || 0 }} 人学</span>
              </div>
            </div>
          </article>
        </div>
      </template>
    </MeCrud>
  </CommonPage>
</template>
