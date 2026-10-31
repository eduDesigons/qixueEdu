<script setup>
import { NButton, NInput, NTag } from 'naive-ui'
import { h, nextTick, ref } from 'vue'
import api from '@/api'
import { CityQueryItem, CommonPage, MeCrud, MeQueryItem } from '@/components'
import { previewImage, useRefreshOnActivate } from '@/composables'
import { formatDateTime } from '@/utils'
import OpsBoard from './OpsBoard.vue'
import { useOpsUserNav } from './opsCells.js'

defineOptions({ name: 'OpsCert' })

const { userCell } = useOpsUserNav()
const $table = ref(null)
const queryItems = ref({ name: '', city_id: null })

useRefreshOnActivate(() => $table.value?.handleSearch())

function showTable() {
  nextTick(() => $table.value?.handleSearch())
}

const columns = [
  { title: '学员', key: 'user', width: 160, render: row => userCell(row) },
  {
    title: '证书',
    key: 'title',
    ellipsis: { tooltip: true },
    render: row => row.title || row.cert?.name || '-',
  },
  {
    title: '来源',
    key: 'source',
    width: 180,
    ellipsis: { tooltip: true },
    render: (row) => {
      if (row.course_id)
        return `课程 · ${row.course?.title || row.course_id}`
      if (row.path_id)
        return `系列 · ${row.path?.title || row.path_id}`
      return '-'
    },
  },
  {
    title: '状态',
    key: 'is_valid',
    width: 90,
    align: 'center',
    render: ({ is_valid }) => h(
      NTag,
      { size: 'small', round: true, bordered: false, type: Number(is_valid) === 1 ? 'success' : 'default' },
      { default: () => (Number(is_valid) === 1 ? '有效' : '已撤回') },
    ),
  },
  {
    title: '颁发时间',
    key: 'issued_at',
    width: 170,
    align: 'center',
    render: row => formatDateTime(row.issued_at || row.created_at),
  },
  {
    title: '操作',
    key: 'actions',
    width: 90,
    align: 'center',
    render: (row) => {
      if (!row.image)
        return h('span', { class: 'text-secondary' }, '生成中')
      return h(NButton, {
        size: 'small',
        secondary: true,
        onClick: () => previewImage(row.image),
      }, { default: () => '查看' })
    },
  },
]
</script>

<template>
  <CommonPage title="证书">
    <OpsBoard scene="cert" :city-id="queryItems.city_id" @show-table="showTable">
    <MeCrud
      ref="$table"
      v-model:query-items="queryItems"
      collapse-query
      :columns="columns"
      :get-data="api.cert.issuedList"
    >
      <CityQueryItem v-model="queryItems.city_id" />
      <MeQueryItem label="关键词" :label-width="60">
        <NInput v-model:value="queryItems.name" clearable placeholder="学员 / 证书 / 编号" class="w-200" />
      </MeQueryItem>
    </MeCrud>
    </OpsBoard>
  </CommonPage>
</template>
