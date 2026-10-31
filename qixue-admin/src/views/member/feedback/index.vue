<script setup>
import { NAvatar, NButton, NImage, NInput, NModal, NSelect, NTag } from 'naive-ui'
import { h } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import api from '@/api'
import { CommonPage, MeCrud, MeQueryItem, CityQueryItem } from '@/components'
import { useRefreshOnActivate } from '@/composables'
import { formatDateTime, getAssetUrl } from '@/utils'

defineOptions({ name: 'MemberFeedback' })

const route = useRoute()
const router = useRouter()
const $table = ref(null)
const queryItems = ref({
  nick_name: '',
  status: null,
  city_id: null,
})
const showReply = ref(false)
const replyRow = ref(null)
const replyText = ref('')
const replying = ref(false)

function applyStatusQuery() {
  const raw = route.query.status
  if (raw === undefined || raw === null || raw === '')
    return
  const status = Number(raw)
  if (status === 0 || status === 1)
    queryItems.value.status = status
}

onMounted(() => {
  applyStatusQuery()
  $table.value?.handleSearch()
})

useRefreshOnActivate(() => {
  applyStatusQuery()
  $table.value?.handleSearch()
})

const statusOptions = [
  { label: '待回复', value: 0 },
  { label: '已回复', value: 1 },
]

function imageList(row) {
  const images = row?.images
  return Array.isArray(images) ? images.filter(Boolean) : []
}

function openReply(row) {
  replyRow.value = row
  replyText.value = row.reply || ''
  showReply.value = true
}

function submitReply() {
  if (!replyRow.value || replying.value)
    return
  const reply = replyText.value.trim()
  if (!reply) {
    window.$message?.warning('请填写回复')
    return
  }
  replying.value = true
  api.feedback.reply({ id: replyRow.value.id, reply }).then(() => {
    window.$message?.success('回复成功')
    showReply.value = false
    $table.value?.handleSearch()
  }).finally(() => {
    replying.value = false
  })
}

const columns = [
  {
    title: '学员',
    key: 'user',
    width: 180,
    render({ user, user_id }) {
      if (!user)
        return '-'
      return h(
        'div',
        {
          class: 'flex items-center cursor-pointer',
          onClick: () => router.push(`/member/user/${user_id}`),
        },
        [
          h(NAvatar, {
            size: 32,
            src: getAssetUrl(user.avatar) || undefined,
            round: true,
            class: 'flex-shrink-0',
          }),
          h('span', { class: 'ml-8 text-primary' }, user.nick_name || '-'),
        ],
      )
    },
  },
  {
    title: '内容',
    key: 'feedback_text',
    minWidth: 240,
    ellipsis: { tooltip: true },
    render(row) {
      return h('span', { class: 'entity-link', onClick: () => openReply(row) }, row.content || '-')
    },
  },
  {
    title: '图片',
    key: 'images',
    width: 160,
    render(row) {
      const list = imageList(row)
      if (!list.length)
        return '-'
      return h('div', { class: 'flex gap-6' }, list.map(url => h(NImage, {
        width: 36,
        height: 36,
        src: getAssetUrl(url),
        objectFit: 'cover',
        style: 'border-radius: 6px',
      })))
    },
  },
  {
    title: '状态',
    key: 'status',
    width: 90,
    align: 'center',
    render: ({ status }) => h(
      NTag,
      { size: 'small', type: Number(status) === 1 ? 'success' : 'warning', round: true, bordered: false },
      { default: () => (Number(status) === 1 ? '已回复' : '待回复') },
    ),
  },
  {
    title: '提交时间',
    key: 'created_at',
    width: 170,
    render: row => formatDateTime(row.created_at),
  },
  {
    title: '操作',
    key: 'actions',
    width: 100,
    align: 'center',
    fixed: 'right',
    render(row) {
      return h(NButton, {
        size: 'small',
        type: 'primary',
        secondary: true,
        onClick: () => openReply(row),
      }, { default: () => (Number(row.status) === 1 ? '查看' : '回复') })
    },
  },
]
</script>

<template>
  <CommonPage>
    <MeCrud
      ref="$table"
      v-model:query-items="queryItems"
      collapse-query
      :columns="columns"
      :get-data="api.feedback.list"
    >
      <CityQueryItem v-model="queryItems.city_id" :label-width="50" />
      <MeQueryItem label="学员" :label-width="50">
        <NInput v-model:value="queryItems.nick_name" placeholder="学员昵称" clearable />
      </MeQueryItem>
      <MeQueryItem label="状态" :label-width="50">
        <NSelect
          v-model:value="queryItems.status"
          :options="statusOptions"
          placeholder="全部"
          clearable
          style="width: 120px"
        />
      </MeQueryItem>
    </MeCrud>

    <NModal v-model:show="showReply" preset="card" title="意见建议" style="width: 560px">
      <div v-if="replyRow" class="feedback-reply">
        <div class="feedback-reply__user">{{ replyRow.user?.nick_name || '学员' }}</div>
        <div class="feedback-reply__content">{{ replyRow.content }}</div>
        <div v-if="imageList(replyRow).length" class="feedback-reply__images">
          <NImage
            v-for="(url, i) in imageList(replyRow)"
            :key="i"
            :src="getAssetUrl(url)"
            width="72"
            height="72"
            object-fit="cover"
          />
        </div>
        <NInput
          v-model:value="replyText"
          type="textarea"
          placeholder="填写回复，学员会在站内信收到"
          :autosize="{ minRows: 4, maxRows: 8 }"
          maxlength="1000"
          show-count
        />
      </div>
      <template #footer>
        <div class="flex justify-end gap-12">
          <NButton @click="showReply = false">取消</NButton>
          <NButton type="primary" :loading="replying" @click="submitReply">回复</NButton>
        </div>
      </template>
    </NModal>
  </CommonPage>
</template>

<style scoped>
.feedback-reply__user {
  margin-bottom: 8px;
  font-weight: 600;
}

.feedback-reply__content {
  margin-bottom: 12px;
  line-height: 1.6;
  white-space: pre-wrap;
}

.feedback-reply__images {
  display: flex;
  gap: 8px;
  margin-bottom: 16px;
}

.feedback-reply__images :deep(img) {
  border-radius: 8px;
}
</style>
