<script setup>
import { NAvatar, NButton, NInput, NSelect, NTag } from 'naive-ui'
import { h, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import api from '@/api/index.js'
import { CityQueryItem, MeCrud, MeQueryItem, SoftUserCard } from '@/components'
import { useCrud } from '@/composables/index.js'
import { getAssetUrl } from '@/utils'

const props = defineProps({
  pathId: {
    type: [String, Number],
    required: true,
  },
})

const emit = defineEmits(['userChange'])

const queryItems = ref({
  path_id: props.pathId,
  nick_name: '',
  status: null,
  city_id: null,
})

const statusOptions = [
  { label: '待审核', value: 0 },
  { label: '通过', value: 1 },
  { label: '拒绝', value: 2 },
]

const router = useRouter()

const $table = ref(null)
onMounted(() => {
  $table.value?.handleSearch()
})

const {
  handleDelete,
} = useCrud({
  name: '学员',
  doDelete: api.path.userPathApplyDelete,
  refresh: () => $table.value?.handleSearch(),
})

const statusMap = {
  0: { name: '待审核', type: 'info', tone: 'info' },
  1: { name: '通过', type: 'success', tone: 'done' },
  2: { name: '拒绝', type: 'error', tone: 'error' },
}

const checkedIds = ref([])
const batchLoading = ref(false)

function onChecked(rowKeys) {
  checkedIds.value = rowKeys || []
}

function openUser(row) {
  router.push({ name: 'UserDetail', params: { id: row.user_id } })
}

function changeStatus(row, status) {
  if (row.status === status)
    return
  api.path.userPathApplyStatusChange({ id: row.id, status }).then(() => {
    $table.value?.handleSearch()
    if (status === 1)
      emit('userChange')
  }).catch((error) => {
    $message.error(error?.message || '操作失败')
  })
}

async function handleBatchChangeStatus(status) {
  if (!checkedIds.value.length) {
    $message.warning('请先勾选需要操作的申请')
    return
  }
  batchLoading.value = true
  try {
    await api.path.userPathApplyStatusChange({ ids: checkedIds.value, status })
    $message.success(status === 1 ? '批量通过成功' : '批量拒绝成功')
    checkedIds.value = []
    $table.value?.handleSearch()
    if (status === 1)
      emit('userChange')
  }
  catch (error) {
    $message.error(error?.message || '批量操作失败')
  }
  finally {
    batchLoading.value = false
  }
}

function openBatchConfirm(status) {
  if (!checkedIds.value.length) {
    $message.warning('请先勾选需要操作的申请')
    return
  }
  const isPass = status === 1
  $dialog.confirm({
    title: '批量审核',
    content: `确定将选中的申请全部${isPass ? '通过' : '拒绝'}？`,
    positiveText: isPass ? '通过' : '拒绝',
    negativeText: '取消',
    type: 'warning',
    confirm: () => handleBatchChangeStatus(status),
  })
}

function renderActions(row) {
  const buttons = []
  if (row.status === 0) {
    buttons.push(
      h(NButton, {
        size: 'small',
        type: 'success',
        secondary: true,
        onClick: () => changeStatus(row, 1),
      }, {
        default: () => '通过',
        icon: () => h('i', { class: 'i-material-symbols:check-circle-outline text-14' }),
      }),
      h(NButton, {
        size: 'small',
        type: 'warning',
        secondary: true,
        onClick: () => changeStatus(row, 2),
      }, {
        default: () => '拒绝',
        icon: () => h('i', { class: 'i-material-symbols:close text-14' }),
      }),
    )
  }
  buttons.push(
    h(NButton, {
      size: 'small',
      type: 'error',
      secondary: true,
      onClick: () => handleDelete(row.id),
    }, {
      default: () => '删除',
      icon: () => h('i', { class: 'i-material-symbols:delete-outline text-14' }),
    }),
  )
  return h('div', { class: 'row-actions', style: { display: 'flex', gap: '8px', justifyContent: 'center', flexWrap: 'wrap' } }, buttons)
}

const columns = [
  { type: 'selection', fixed: 'left', disabled: row => row?.status !== 0 },
  {
    title: '学员',
    key: 'user',
    render({ user, user_id }) {
      return h(
        'div',
        {
          class: 'flex cursor-pointer',
          onClick: () => router.push({ name: 'UserDetail', params: { id: user_id } }),
        },
        [
          h(NAvatar, {
            size: 60,
            src: getAssetUrl(user?.avatar) || undefined,
            class: 'flex-shrink-0',
          }),
          h('div', { class: 'ml-12 flex-col justify-center' }, [
            h('div', { class: 'text-16 font-bold text-primary' }, user?.nick_name || '未命名'),
            h(NTag, { type: 'info', size: 'small', bordered: false }, `${user?.credit ?? 0}学分`),
          ]),
        ],
      )
    },
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
    render: ({ status }) => h(NTag, { type: statusMap[status]?.type || 'default', size: 'small', round: true, bordered: false }, { default: () => statusMap[status]?.name || '—' }),
  },
  {
    title: '申请时间',
    key: 'created_at',
    align: 'center',
  },
  {
    title: '操作',
    key: 'actions',
    width: 220,
    align: 'center',
    fixed: 'right',
    hideInExcel: true,
    render: renderActions,
  },
]
</script>

<template>
  <div class="h-500">
    <div v-if="checkedIds.length" class="mb-12 flex items-center">
      <NButton type="primary" :loading="batchLoading" :disabled="batchLoading" @click="openBatchConfirm(1)">
        批量通过
      </NButton>
      <NButton class="ml-12" type="warning" :loading="batchLoading" :disabled="batchLoading" @click="openBatchConfirm(2)">
        批量拒绝
      </NButton>
    </div>
    <MeCrud
      ref="$table"
      v-model:query-items="queryItems"
      :bound-query="{ path_id: pathId }"
      :columns="columns"
      :get-data="api.path.userPathApplyList"
      @on-checked="onChecked"
      @onDataChange="checkedIds = []"
    >
      <CityQueryItem v-model="queryItems.city_id" :label-width="50" />
      <MeQueryItem label="学员" :label-width="50">
        <NInput
          v-model:value="queryItems.nick_name"
          placeholder="学员昵称"
          clearable
        />
      </MeQueryItem>
      <MeQueryItem label="状态" :label-width="50">
        <NSelect
          v-model:value="queryItems.status"
          :options="statusOptions"
          placeholder="全部"
          clearable
        />
      </MeQueryItem>

      <template #card="{ list }">
        <div class="entity-grid">
          <SoftUserCard
            v-for="item in list"
            :key="item.id"
            :avatar="item.user?.avatar"
            :name="item.user?.nick_name"
            :subtitle="`${item.user?.credit ?? 0} 学分`"
            :status="statusMap[item.status]?.name || '—'"
            :status-tone="statusMap[item.status]?.tone || 'idle'"
            :desc="item.reason?.trim() || '未填写申请原因'"
            :time="item.created_at"
            deletable
            @click="openUser(item)"
            @delete="handleDelete(item.id)"
          >
            <template v-if="item.status === 0" #actions>
              <NButton size="tiny" type="success" secondary @click="changeStatus(item, 1)">
                通过
              </NButton>
              <NButton size="tiny" type="warning" secondary @click="changeStatus(item, 2)">
                拒绝
              </NButton>
            </template>
          </SoftUserCard>
        </div>
      </template>
    </MeCrud>
  </div>
</template>
