<script setup>
import { NSelect, NTag } from 'naive-ui'
import { h } from 'vue'
import api from '@/api'
import { MeCrud, MeQueryItem } from '@/components'
import { formatDateTime } from '@/utils'

const props = defineProps({
  userId: {
    type: [String, Number],
    required: true,
  },
})

const statusMap = {
  0: { name: '待审核', type: 'info' },
  1: { name: '通过', type: 'success' },
  2: { name: '拒绝', type: 'error' },
}

const queryItems = ref({ user_id: props.userId, status: null })
const statusOptions = [
  { label: '待审核', value: 0 },
  { label: '通过', value: 1 },
  { label: '拒绝', value: 2 },
]
const $table = ref(null)
onMounted(() => {
  $table.value?.handleSearch()
})

const columns = [
  {
    title: '课程',
    key: 'course',
    width: 200,
    render: ({ course }) => course?.title ?? '-',
  },
  {
    title: '申请原因',
    key: 'reason',
    align: 'center',
    ellipsis: { tooltip: true },
    render: ({ reason }) => reason?.trim() || h('span', { class: 'text-gray-400' }, '未填写'),
  },
  {
    title: '状态',
    key: 'status',
    align: 'center',
    width: 100,
    render: ({ status }) => h(NTag, { type: statusMap[status]?.type ?? 'default', size: 'small', round: true, bordered: false }, statusMap[status]?.name ?? status),
  },
  {
    title: '申请时间',
    key: 'created_at',
    align: 'center',
    width: 180,
    render: row => h('span', formatDateTime(row.created_at)),
  },
]
</script>

<template>
  <div class="h-500">
    <MeCrud
      ref="$table"
      v-model:query-items="queryItems"
      :bound-query="{ user_id: userId }"
      :columns="columns"
      :get-data="api.course.courseUserApplyList"
    >
      <MeQueryItem label="状态" :label-width="50">
        <NSelect
          v-model:value="queryItems.status"
          :options="statusOptions"
          placeholder="全部"
          clearable
        />
      </MeQueryItem>
    </MeCrud>
  </div>
</template>
