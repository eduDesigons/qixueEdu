<script setup>
import { NInput, NSelect, NTag } from 'naive-ui'
import { h, nextTick, ref } from 'vue'
import api from '@/api'
import { CityQueryItem, CommonPage, MeCrud, MeQueryItem } from '@/components'
import { useRefreshOnActivate } from '@/composables'
import { formatDateTime } from '@/utils'
import OpsBoard from './OpsBoard.vue'
import { useOpsUserNav } from './opsCells.js'

defineOptions({ name: 'OpsLearn' })

const { userCell } = useOpsUserNav()
const $table = ref(null)
const queryItems = ref({ nick_name: '', title: '', is_end: 1, city_id: null })

useRefreshOnActivate(() => $table.value?.handleSearch())

function showTable() {
  nextTick(() => $table.value?.handleSearch())
}

const statusOptions = [
  { label: '全部', value: null },
  { label: '已学完', value: 1 },
  { label: '未学完', value: 0 },
]

const columns = [
  { title: '学员', key: 'user', width: 180, render: row => userCell(row) },
  {
    title: '课程',
    key: 'course',
    width: 160,
    ellipsis: { tooltip: true },
    render: ({ course }) => course?.title || '-',
  },
  {
    title: '课节',
    key: 'chapter',
    ellipsis: { tooltip: true },
    render: ({ chapter }) => chapter?.title || '-',
  },
  {
    title: '状态',
    key: 'is_end',
    width: 100,
    align: 'center',
    render: ({ is_end }) => h(
      NTag,
      { size: 'small', round: true, bordered: false, type: Number(is_end) === 1 ? 'success' : 'warning' },
      { default: () => (Number(is_end) === 1 ? '已学完' : '未学完') },
    ),
  },
  {
    title: '更新时间',
    key: 'updated_at',
    width: 170,
    align: 'center',
    render: row => formatDateTime(row.updated_at),
  },
]
</script>

<template>
  <CommonPage title="学习">
    <OpsBoard scene="learn" :city-id="queryItems.city_id" @show-table="showTable">
    <MeCrud
      ref="$table"
      v-model:query-items="queryItems"
      collapse-query
      :columns="columns"
      :get-data="api.ops.learnList"
    >
      <CityQueryItem v-model="queryItems.city_id" />
      <MeQueryItem label="学员" :label-width="60">
        <NInput v-model:value="queryItems.nick_name" clearable placeholder="昵称" class="w-160" />
      </MeQueryItem>
      <MeQueryItem label="关键词" :label-width="60">
        <NInput v-model:value="queryItems.title" clearable placeholder="课程 / 课节" class="w-160" />
      </MeQueryItem>
      <MeQueryItem label="状态" :label-width="60">
        <NSelect
          v-model:value="queryItems.is_end"
          :options="statusOptions"
          clearable
          class="w-120"
          @update:value="$table?.handleSearch(false)"
        />
      </MeQueryItem>
    </MeCrud>
    </OpsBoard>
  </CommonPage>
</template>
