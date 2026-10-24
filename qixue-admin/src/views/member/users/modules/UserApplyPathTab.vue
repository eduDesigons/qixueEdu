<script setup>
import { getAssetUrl } from '@/utils'
import { NAvatar, NButton, NTag } from 'naive-ui'
import { h } from 'vue'
import { useRouter } from 'vue-router'
import api from '@/api/index.js'
import { MeCrud } from '@/components'
import { useCrud } from '@/composables/index.js'

const props = defineProps({
  pathId: {
    type: [String, Number],
    required: true,
  },
})

const emit = defineEmits(['userChange'])

const queryItems = ref({
  path_id: props.pathId,
})

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
  0: {
    name: '待审核',
    type: 'info',
  },
  1: {
    name: '通过',
    type: 'success',
  },
  2: {
    name: '拒绝',
    type: 'error',
  },
}

function changeStatus(row, status) {
  if (row.status === status) {
    return
  }
  api.path.userPathApplyStatusChange({ id: row.id, status }).then(() => {
    $table.value?.handleSearch()
    if (status === 1) {
      emit('userChange')
    }
  })
}

// 试卷列表数据
const columns = [
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
            src: getAssetUrl(user.avatar) || undefined,
            class: 'flex-shrink-0',
          }),
          h(
            'div',
            {
              class: 'ml-12 flex-col justify-center',
            },
            [
              h(
                'div',
                {
                  class: 'text-16 font-bold text-primary',
                },
                user.nick_name,
              ),
              h(
                NTag,
                {
                  type: 'info',
                  size: 'small',
                  bordered: false,
                },
                `${user.credit}学分`,
              ),
            ],
          ),
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
    key: 'is_end',
    align: 'center',
    render: ({ status }) => h(NTag, { type: statusMap[status].type, size: 'small', round: true, bordered: false }, statusMap[status].name),
  },
  {
    title: '申请时间',
    key: 'created_at',
    align: 'center',
  },
  {
    title: '操作',
    key: 'actions',
    width: 120,
    align: 'center',
    fixed: 'right',
    hideInExcel: true,
    render(row) {
      return [
        row.status === 0 && h(
          NButton,
          {
            size: 'small',
            type: 'error',
            style: 'margin-left: 12px;',
            onClick: () => changeStatus(row, 1),
          },
          {
            default: () => '通过',
            icon: () => h('i', { class: 'i-material-symbols:delete-outline text-14' }),
          },
        ),
        row.status === 0 && h(
          NButton,
          {
            size: 'small',
            type: 'error',
            style: 'margin-left: 12px;',
            onClick: () => changeStatus(row, 2),
          },
          {
            default: () => '拒绝',
            icon: () => h('i', { class: 'i-material-symbols:delete-outline text-14' }),
          },
        ),
        h(
          NButton,
          {
            size: 'small',
            type: 'error',
            style: 'margin-left: 12px;',
            onClick: () => handleDelete(row.id),
          },
          {
            default: () => '删除',
            icon: () => h('i', { class: 'i-material-symbols:delete-outline text-14' }),
          },
        ),
      ]
    },
  },
]
</script>

<template>
  <div class="h-500">
    <MeCrud
      ref="$table"
      v-model:query-items="queryItems"
      :bound-query="{ path_id: pathId }"
      :columns="columns"
      :get-data="api.path.userPathApplyList"
    />
  </div>
</template>
