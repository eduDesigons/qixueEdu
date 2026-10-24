<script setup>
import { NInput } from 'naive-ui'
import { nextTick, ref } from 'vue'
import api from '@/api'
import { CityQueryItem, CommonPage, MeCrud, MeQueryItem } from '@/components'
import { useRefreshOnActivate } from '@/composables'
import { formatDateTime } from '@/utils'
import OpsBoard from './OpsBoard.vue'
import { sumCell, useOpsUserNav } from './opsCells.js'

defineOptions({ name: 'OpsCredit' })

const { userCell } = useOpsUserNav()
const $table = ref(null)
const queryItems = ref({ nick_name: '', field: '', city_id: null })

useRefreshOnActivate(() => $table.value?.handleSearch())

function showTable() {
  nextTick(() => $table.value?.handleSearch())
}

const columns = [
  { title: '学员', key: 'user', width: 180, render: row => userCell(row) },
  { title: '类型', key: 'field', width: 120, align: 'center' },
  { title: '学分', key: 'sum', width: 100, align: 'center', render: ({ sum }) => sumCell(sum) },
  { title: '备注', key: 'remark', ellipsis: { tooltip: true } },
  {
    title: '时间',
    key: 'created_at',
    width: 170,
    align: 'center',
    render: row => formatDateTime(row.created_at),
  },
]
</script>

<template>
  <CommonPage title="学分">
    <OpsBoard scene="credit" :city-id="queryItems.city_id" @show-table="showTable">
    <MeCrud
      ref="$table"
      v-model:query-items="queryItems"
      collapse-query
      :columns="columns"
      :get-data="api.ops.creditList"
    >
      <CityQueryItem v-model="queryItems.city_id" />
      <MeQueryItem label="学员" :label-width="60">
        <NInput v-model:value="queryItems.nick_name" clearable placeholder="昵称" class="w-160" />
      </MeQueryItem>
      <MeQueryItem label="类型" :label-width="60">
        <NInput v-model:value="queryItems.field" clearable placeholder="field" class="w-140" />
      </MeQueryItem>
    </MeCrud>
    </OpsBoard>
  </CommonPage>
</template>
