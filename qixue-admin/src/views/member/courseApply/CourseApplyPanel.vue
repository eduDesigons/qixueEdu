<script setup>
import { getAssetUrl } from '@/utils'
import { NAvatar, NButton, NDatePicker, NInput, NSelect, NTag } from 'naive-ui'
import { computed, h, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import api from '@/api'
import { MeCrud, MeQueryItem, CityQueryItem } from '@/components'

const emit = defineEmits(['checked', 'changed'])

const router = useRouter()
const $table = ref(null)
const checkedIds = ref([])
const batchLoading = ref(false)
const crudKey = ref(0)
const queryItems = ref({
  course_title: '',
  nick_name: '',
  city_id: null,
  status: 0,
  start_time: '',
  end_time: '',
})

const timeRangeValue = computed(() => {
  const { start_time, end_time } = queryItems.value
  if (!start_time || !end_time)
    return null
  return [new Date(start_time).getTime(), new Date(end_time).getTime()]
})

function formatLocalDay(date) {
  const d = new Date(date)
  const y = d.getFullYear()
  const m = String(d.getMonth() + 1).padStart(2, '0')
  const day = String(d.getDate()).padStart(2, '0')
  return `${y}-${m}-${day}`
}

function onTimeRangeUpdate(val) {
  if (!val || val.length !== 2) {
    queryItems.value.start_time = ''
    queryItems.value.end_time = ''
    return
  }
  queryItems.value.start_time = `${formatLocalDay(val[0])} 00:00:00`
  queryItems.value.end_time = `${formatLocalDay(val[1])} 23:59:59`
}

onMounted(() => {
  $table.value?.handleSearch()
})

const loadingId = ref(null)
const loadingAction = ref(null)

const statusMap = {
  0: { name: '待审核', type: 'info' },
  1: { name: '通过', type: 'success' },
  2: { name: '拒绝', type: 'error' },
}

const statusOptions = [
  { label: '待审核', value: 0 },
  { label: '通过', value: 1 },
  { label: '拒绝', value: 2 },
]

function onChecked(rowKeys) {
  checkedIds.value = rowKeys || []
  emit('checked', checkedIds.value)
}

async function changeStatus(row, status) {
  if (row.status === status)
    return
  loadingId.value = row.id
  loadingAction.value = status === 1 ? 'approve' : 'reject'
  try {
    await api.course.courseUserApplyStatusChange({ id: row.id, status })
    $message.success(status === 1 ? '已通过' : '已拒绝')
    $table.value?.handleSearch()
    emit('changed')
  }
  catch (error) {
    $message.error(error?.message || '操作失败')
  }
  finally {
    loadingId.value = null
    loadingAction.value = null
  }
}

async function handleBatchChangeStatus(status) {
  if (!checkedIds.value.length)
    return $message.warning('请先勾选需要操作的申请')
  batchLoading.value = true
  try {
    await api.course.courseUserApplyStatusChange({ ids: checkedIds.value, status })
    $message.success(status === 1 ? '批量通过成功' : '批量拒绝成功')
    checkedIds.value = []
    emit('checked', [])
    emit('changed')
    crudKey.value += 1
    await nextTick()
    $table.value?.handleSearch()
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
    closable: false,
    maskClosable: false,
    confirm: async () => {
      await handleBatchChangeStatus(status)
    },
    cancel: () => {
      checkedIds.value = []
      emit('checked', [])
      crudKey.value += 1
      nextTick(() => $table.value?.handleSearch())
    },
  })
}

async function handleDeleteWithLoading(row) {
  loadingId.value = row.id
  loadingAction.value = 'delete'
  try {
    await api.course.courseUserApplyDelete(row.id)
    $message.success('删除成功')
    $table.value?.handleSearch()
    emit('changed')
  }
  finally {
    loadingId.value = null
    loadingAction.value = null
  }
}

const columns = [
  { type: 'selection', fixed: 'left', disabled: row => row?.status !== 0 },
  {
    title: '编号',
    key: 'id',
    width: 80,
    align: 'center',
    resizable: true,
  },
  {
    title: '学员',
    key: 'user',
    width: 220,
    render({ user, user_id }) {
      return h(
        'div',
        {
          class: 'flex cursor-pointer items-center',
          onClick: () => router.push({ name: 'UserDetail', params: { id: user_id } }),
        },
        [
          h(NAvatar, { size: 48, round: true, src: getAssetUrl(user?.avatar) || undefined, class: 'flex-shrink-0' }),
          h('div', { class: 'ml-12 flex flex-col justify-center gap-4' }, [
            h('div', { class: 'text-14 font-bold text-primary' }, user?.nick_name ?? '-'),
            h(NTag, { type: 'info', size: 'small', bordered: false, round: true }, { default: () => `${user?.credit ?? 0}学分` }),
          ]),
        ],
      )
    },
  },
  {
    title: '课程',
    key: 'course',
    width: 200,
    ellipsis: { tooltip: true },
    render({ course, course_id }) {
      if (!course)
        return '-'
      return h(
        'div',
        {
          class: 'cursor-pointer text-primary',
          onClick: () => router.push(`/course/detail/${course_id}`),
        },
        course.title,
      )
    },
  },
  {
    title: '申请原因',
    key: 'reason',
    resizable: true,
    ellipsis: { tooltip: true },
    render: ({ reason }) => reason?.trim() || h('span', { class: 'text-gray-400' }, '未填写'),
  },
  {
    title: '状态',
    key: 'status',
    width: 100,
    align: 'center',
    render: ({ status }) => h(NTag, { type: statusMap[status].type, size: 'small', round: true, bordered: false }, { default: () => statusMap[status].name }),
  },
  {
    title: '申请时间',
    key: 'created_at',
    width: 160,
    align: 'center',
  },
  {
    title: '操作',
    key: 'actions',
    width: 200,
    align: 'center',
    fixed: 'right',
    hideInExcel: true,
    render(row) {
      const isRowLoading = loadingId.value === row.id
      const action = loadingAction.value
      return h('div', { class: 'flex items-center justify-center gap-8' }, [
        row.status === 0 && h(
          NButton,
          {
            size: 'small',
            type: 'primary',
            secondary: true,
            loading: isRowLoading && action === 'approve',
            disabled: isRowLoading,
            onClick: () => changeStatus(row, 1),
          },
          { default: () => '通过' },
        ),
        row.status === 0 && h(
          NButton,
          {
            size: 'small',
            type: 'warning',
            secondary: true,
            loading: isRowLoading && action === 'reject',
            disabled: isRowLoading,
            onClick: () => changeStatus(row, 2),
          },
          { default: () => '拒绝' },
        ),
        row.status === 0 && h(
          NButton,
          {
            size: 'small',
            type: 'error',
            secondary: true,
            loading: isRowLoading && action === 'delete',
            disabled: isRowLoading,
            onClick: () => {
              $dialog.warning({
                content: '确定删除该申请？',
                title: '提示',
                positiveText: '确定',
                negativeText: '取消',
                onPositiveClick: () => handleDeleteWithLoading(row),
              })
            },
          },
          { default: () => '删除', icon: () => h('i', { class: 'i-material-symbols:delete-outline text-14' }) },
        ),
      ])
    },
  },
]

defineExpose({
  batchApprove: () => openBatchConfirm(1),
  batchReject: () => openBatchConfirm(2),
  batchLoading,
})
</script>

<template>
  <MeCrud
    :key="crudKey"
    ref="$table"
    v-model:query-items="queryItems"
    collapse-query
    cache-key="member-course-apply"
    :scroll-x="1100"
    :columns="columns"
    :get-data="api.course.courseUserApplyList"
    @on-checked="onChecked"
    @onDataChange="onChecked([])"
  >
    <CityQueryItem v-model="queryItems.city_id" />
    <MeQueryItem label="课程名称" :label-width="80">
      <NInput
        v-model:value="queryItems.course_title"
        clearable
        placeholder="请输入课程名称"
      />
    </MeQueryItem>
    <MeQueryItem label="学员名称" :label-width="80">
      <NInput
        v-model:value="queryItems.nick_name"
        clearable
        placeholder="请输入学员名称"
      />
    </MeQueryItem>
    <MeQueryItem label="状态" :label-width="80">
      <NSelect
        v-model:value="queryItems.status"
        clearable
        placeholder="全部状态"
        :options="statusOptions"
      />
    </MeQueryItem>
    <MeQueryItem label="申请时间" :label-width="80">
      <NDatePicker
        :value="timeRangeValue"
        type="daterange"
        clearable
        placeholder="请选择申请时间"
        style="width: 100%"
        @update:value="onTimeRangeUpdate"
      />
    </MeQueryItem>
  </MeCrud>
</template>
