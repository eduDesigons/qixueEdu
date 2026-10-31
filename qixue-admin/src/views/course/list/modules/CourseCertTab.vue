<script setup>
import { NAvatar, NButton, NImage, NTag } from 'naive-ui'
import { computed, h, nextTick, onMounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import api from '@/api'
import { CityQueryItem, MeCrud, MeQueryItem } from '@/components'
import { previewImage } from '@/composables'
import { formatDateTime, getAssetUrl } from '@/utils'

defineOptions({ name: 'EntityCertTab' })

const props = defineProps({
  /** 课程 id（与 pathId 二选一） */
  courseId: {
    type: [String, Number],
    default: null,
  },
  /** 系列 id */
  pathId: {
    type: [String, Number],
    default: null,
  },
  /** 待补发人数（来自详情接口） */
  pendingCount: {
    type: [Number, String],
    default: 0,
  },
})

const emit = defineEmits(['change'])

const router = useRouter()
const $table = ref(null)
const backfilling = ref(false)

const isCourse = computed(() => props.courseId != null && props.courseId !== '')
const scopeId = computed(() => (isCourse.value ? props.courseId : props.pathId))
const pending = computed(() => Number(props.pendingCount || 0))

const queryItems = ref({
  nick_name: null,
  city_id: null,
})

function fetchIssued(params = {}) {
  const payload = {
    ...params,
    name: queryItems.value.nick_name || params.name || '',
  }
  if (isCourse.value)
    payload.course_id = props.courseId
  else
    payload.path_id = props.pathId
  return api.cert.issuedList(payload)
}

function search(resetPage = false) {
  nextTick(() => {
    if (resetPage)
      $table.value?.handleSearch(false)
    else
      $table.value?.handleSearch()
  })
}

onMounted(() => search())

watch(
  () => [props.courseId, props.pathId],
  () => search(true),
)

function openUser(row) {
  if (row.user_id)
    router.push({ name: 'UserDetail', params: { id: row.user_id } })
}

function preview(row) {
  previewImage(row?.image)
}

function askBackfill() {
  if (!scopeId.value)
    return
  const n = pending.value
  $dialog.warning({
    title: '检测并补发',
    content: n > 0
      ? `当前有 ${n} 名已结业学员尚未获得有效证书，是否立即补发？${n > 500 ? '（单次最多 500 人）' : ''}`
      : '当前没有待补发学员。仍要再检查一遍吗？',
    positiveText: n > 0 ? '补发' : '检查',
    negativeText: '取消',
    async onPositiveClick() {
      backfilling.value = true
      try {
        const payload = isCourse.value
          ? { course_id: props.courseId }
          : { path_id: props.pathId }
        const res = await api.cert.backfill(payload)
        const count = res?.data?.count
        if (count > 0)
          $message.success(`已为 ${count} 名学员补发证书`)
        else
          $message.success(res?.message || '暂无需要补发的学员')
        search()
        emit('change')
      }
      catch {
      }
      finally {
        backfilling.value = false
      }
    },
  })
}

function reissue(row) {
  const name = row.user?.nick_name || '该学员'
  $dialog.warning({
    title: '重发证书',
    content: `将按最新模板为「${name}」重新生成证书图片并恢复为有效。`,
    positiveText: '确定',
    negativeText: '取消',
    async onPositiveClick() {
      try {
        const payload = {
          user_id: row.user_id,
          force: true,
        }
        if (isCourse.value)
          payload.course_id = props.courseId
        else
          payload.path_id = props.pathId
        await api.cert.issue(payload)
        $message.success('已重发')
        search()
        emit('change')
      }
      catch {
      }
    },
  })
}

function revoke(row) {
  const name = row.user?.nick_name || '该学员'
  $dialog.warning({
    title: '撤回证书',
    content: `撤回后「${name}」在学员端将看不到该证书，记录仍保留，可随时重发。`,
    positiveText: '撤回',
    negativeText: '取消',
    async onPositiveClick() {
      try {
        await api.cert.revoke({ id: row.id })
        $message.success('已撤回')
        search()
        emit('change')
      }
      catch {
      }
    },
  })
}

const columns = [
  {
    title: '学员',
    key: 'user',
    width: 180,
    render({ user, user_id }) {
      return h(
        'div',
        {
          class: 'flex cursor-pointer items-center',
          onClick: () => openUser({ user_id }),
        },
        [
          h(NAvatar, {
            size: 40,
            src: getAssetUrl(user?.avatar) || undefined,
            class: 'flex-shrink-0',
          }),
          h('div', { class: 'ml-12 font-medium' }, user?.nick_name || '-'),
        ],
      )
    },
  },
  {
    title: '证书',
    key: 'title',
    ellipsis: { tooltip: true },
  },
  {
    title: '编号',
    key: 'cert_no',
    width: 180,
    ellipsis: { tooltip: true },
  },
  {
    title: '预览',
    key: 'image',
    width: 80,
    align: 'center',
    render: ({ image }) => image
      ? h(NImage, { src: getAssetUrl(image), width: 48, height: 32, objectFit: 'cover' })
      : '-',
  },
  {
    title: '状态',
    key: 'is_valid',
    width: 90,
    align: 'center',
    render: ({ is_valid }) => h(
      NTag,
      {
        type: Number(is_valid) === 1 ? 'success' : 'warning',
        size: 'small',
        round: true,
        bordered: false,
      },
      { default: () => (Number(is_valid) === 1 ? '有效' : '已撤回') },
    ),
  },
  {
    title: '颁发时间',
    key: 'issued_at',
    width: 170,
    render: row => h('span', formatDateTime(row.issued_at || row.created_at)),
  },
  {
    title: '操作',
    key: 'actions',
    width: 220,
    align: 'center',
    fixed: 'right',
    hideInExcel: true,
    render(row) {
      const buttons = [
        h(NButton, {
          size: 'small',
          secondary: true,
          type: 'info',
          disabled: !row.image,
          onClick: () => preview(row),
        }, { default: () => '查看' }),
        h(NButton, {
          size: 'small',
          secondary: true,
          type: 'primary',
          onClick: () => reissue(row),
        }, { default: () => '重发' }),
      ]
      if (Number(row.is_valid) === 1) {
        buttons.push(h(NButton, {
          size: 'small',
          secondary: true,
          type: 'error',
          onClick: () => revoke(row),
        }, { default: () => '撤回' }))
      }
      return h('div', { class: 'flex gap-8 justify-center flex-wrap' }, buttons)
    },
  },
]

defineExpose({
  refresh: () => search(true),
  askBackfill,
  backfilling,
})
</script>

<template>
  <div class="h-500">
    <MeCrud
      ref="$table"
      v-model:query-items="queryItems"
      :columns="columns"
      :get-data="fetchIssued"
    >
      <CityQueryItem v-model="queryItems.city_id" :label-width="50" />
      <MeQueryItem label="学员" :label-width="50">
        <n-input
          v-model:value="queryItems.nick_name"
          type="text"
          placeholder="昵称 / 编号"
          clearable
        />
      </MeQueryItem>
    </MeCrud>
  </div>
</template>
