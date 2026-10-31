<script setup>
import { NAvatar, NSelect, NTag } from 'naive-ui'
import { useRouter } from 'vue-router'
import api from '@/api'
import { MeCrud, MeQueryItem } from '@/components'
import { CommonPage } from '@/components/index.js'
import { useCrud, useRefreshOnActivate } from '@/composables'
import { getAssetUrl, formatDateTime, h } from '@/utils'

defineOptions({ name: 'CourseCollect' })

const router = useRouter()
const $table = ref(null)
const queryItems = ref({})

onMounted(() => {
  $table.value?.handleSearch()
})

useRefreshOnActivate(() => $table.value?.handleSearch())

const pathCategories = ref([])
function getCategories() {
  api.path.getPathCategories().then((res) => {
    pathCategories.value = res.data || []
  }).catch((error) => {
    $message.error(error?.message || '获取系列分类失败')
    pathCategories.value = []
  })
}
getCategories()

useCrud({
  name: '系列',
  doDelete: api.path.pathDelete,
  refresh: () => $table.value?.handleSearch(),
})

const columns = [
  {
    title: '系列',
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
    title: '分类',
    key: 'category',
    align: 'center',
    width: 120,
    render: ({ category }) => h(NTag, { type: 'info', size: 'small', round: true, bordered: false }, { default: () => category?.name || '—' }),
  },
  {
    title: '课程',
    key: 'courses_count',
    align: 'center',
    width: 100,
    render: ({ courses_count }) => h('span', { class: 'text font-bold' }, courses_count),
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

function handleViewPath(row) {
  router.push(`/course/path/${row.id}`)
}

function handleAddPath() {
  router.push('/course/path/add')
}
</script>

<template>
  <CommonPage>
    <template #action>
      <NButton type="primary" @click="handleAddPath()">
        <i class="i-material-symbols:add mr-4 text-18" />
        添加系列课
      </NButton>
    </template>

    <MeCrud
      ref="$table"
      v-model:query-items="queryItems"
      collapse-query
      :scroll-x="1000"
      :columns="columns"
      :get-data="api.path.pathList"
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
        <NSelect
          v-model:value="queryItems.category_id"
          :options="pathCategories"
          filterable
          clearable
          value-field="id"
          label-field="name"
        />
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
            @click="handleViewPath(item)"
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
                <span class="entity-card__badge">{{ item.courses_count || 0 }} 门课</span>
                <span v-if="item.min_need_credit" class="entity-card__badge">限 {{ item.min_need_credit }} 学分</span>
                <span class="entity-card__badge">{{ item.users_count || 0 }} 人学</span>
                <span v-if="item.cert?.name" class="entity-card__badge">{{ item.cert.name }}</span>
              </div>
            </div>
          </article>
        </div>
      </template>
    </MeCrud>
  </CommonPage>
</template>
