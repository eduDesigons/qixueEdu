<script setup>
import { NInput, NSelect } from 'naive-ui'
import { nextTick, ref } from 'vue'
import api from '@/api'
import { CityQueryItem, CommonPage, MeCrud, MeQueryItem } from '@/components'
import { useRefreshOnActivate } from '@/composables'
import { formatDateTime } from '@/utils'
import OpsBoard from './OpsBoard.vue'
import { endOptions, endTag, useOpsUserNav } from './opsCells.js'

defineOptions({ name: 'OpsPathEnroll' })

const { userCell } = useOpsUserNav()
const $table = ref(null)
const queryItems = ref({ nick_name: '', title: '', is_end: null, city_id: null })

useRefreshOnActivate(() => $table.value?.handleSearch())

function showTable() {
  nextTick(() => $table.value?.handleSearch())
}

const columns = [
  { title: '学员', key: 'user', width: 180, render: row => userCell(row) },
  {
    title: '系列',
    key: 'path',
    ellipsis: { tooltip: true },
    render: ({ path }) => path?.title || '-',
  },
  { title: '状态', key: 'is_end', width: 100, align: 'center', render: ({ is_end }) => endTag(is_end) },
  {
    title: '加入时间',
    key: 'created_at',
    width: 170,
    align: 'center',
    render: row => formatDateTime(row.created_at),
  },
]
</script>

<template>
  <CommonPage title="入系">
    <OpsBoard scene="path" :city-id="queryItems.city_id" @show-table="showTable">
    <MeCrud
      ref="$table"
      v-model:query-items="queryItems"
      collapse-query
      :columns="columns"
      :get-data="api.ops.pathEnrollList"
    >
      <CityQueryItem v-model="queryItems.city_id" />
      <MeQueryItem label="学员" :label-width="60">
        <NInput v-model:value="queryItems.nick_name" clearable placeholder="昵称" class="w-160" />
      </MeQueryItem>
      <MeQueryItem label="系列" :label-width="60">
        <NInput v-model:value="queryItems.title" clearable placeholder="系列名" class="w-160" />
      </MeQueryItem>
      <MeQueryItem label="状态" :label-width="60">
        <NSelect
          v-model:value="queryItems.is_end"
          :options="endOptions"
          clearable
          class="w-120"
          @update:value="$table?.handleSearch(false)"
        />
      </MeQueryItem>
    </MeCrud>
    </OpsBoard>
  </CommonPage>
</template>
