<script setup>
import { NTag } from 'naive-ui'
import { h } from 'vue'
import api from '@/api'
import { MeCrud } from '@/components'
import { formatDateTime } from '@/utils'

const props = defineProps({
  userId: {
    type: [String, Number],
    required: true,
  },
})

const queryItems = ref({ user_id: props.userId })
const $table = ref(null)
onMounted(() => {
  $table.value?.handleSearch()
})

const columns = [
  { title: '类型', key: 'field', align: 'center', width: 120 },
  { title: '学分', key: 'sum', align: 'center', width: 100, render: ({ sum }) => h('span', sum > 0 ? `+${sum}` : String(sum)) },
  { title: '备注', key: 'remark', align: 'left', ellipsis: { tooltip: true } },
  {
    title: '时间',
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
      :get-data="api.user.getUserCreditRecordList"
    />
  </div>
</template>
