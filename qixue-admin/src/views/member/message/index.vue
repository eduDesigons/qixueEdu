<script setup>
import { NAvatar, NButton, NInput, NModal, NSelect, NTag } from 'naive-ui'
import { h } from 'vue'
import { useRouter } from 'vue-router'
import api from '@/api'
import { CommonPage, MeCrud, MeQueryItem, CityQueryItem } from '@/components'
import { useCrud, useRefreshOnActivate } from '@/composables'
import { formatDateTime, formatRichContentHtml, getAssetUrl } from '@/utils'

defineOptions({ name: 'MemberMessage' })

const router = useRouter()
const $table = ref(null)
const queryItems = ref({
  nick_name: '',
  title: '',
  type: null,
  is_read: null,
  city_id: null,
})

const showView = ref(false)
const viewRow = ref(null)

onMounted(() => {
  $table.value?.handleSearch()
})

useRefreshOnActivate(() => $table.value?.handleSearch())

const { handleDelete } = useCrud({
  name: '站内信',
  doDelete: id => api.message.delete(id),
  refresh: () => $table.value?.handleSearch(),
})

const typeMap = {
  notice: { name: '系统通知', type: 'info' },
  paper: { name: '作业批改', type: 'warning' },
}

const typeOptions = [
  { label: '系统通知', value: 'notice' },
  { label: '作业批改', value: 'paper' },
]

const readOptions = [
  { label: '未读', value: 0 },
  { label: '已读', value: 1 },
]

function plainText(html) {
  return String(html || '')
    .replace(/<[^>]+>/g, '')
    .replace(/&nbsp;/gi, ' ')
    .replace(/&amp;/gi, '&')
    .replace(/&lt;/gi, '<')
    .replace(/&gt;/gi, '>')
    .replace(/\s+/g, ' ')
    .trim()
}

function openView(row) {
  viewRow.value = row
  showView.value = true
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
    title: '标题',
    key: 'title',
    ellipsis: { tooltip: true },
    render(row) {
      return h(
        'span',
        {
          class: 'entity-link',
          onClick: () => openView(row),
        },
        row.title || '-',
      )
    },
  },
  {
    title: '内容摘要',
    key: 'content',
    ellipsis: { tooltip: true },
    render: ({ content }) => plainText(content) || '-',
  },
  {
    title: '类型',
    key: 'type',
    width: 110,
    align: 'center',
    render: ({ type }) => {
      const item = typeMap[type] || { name: type || '通知', type: 'default' }
      return h(NTag, { size: 'small', type: item.type, round: true, bordered: false }, { default: () => item.name })
    },
  },
  {
    title: '阅读',
    key: 'is_read',
    width: 80,
    align: 'center',
    render: ({ is_read }) => h(
      NTag,
      { size: 'small', type: is_read ? 'success' : 'warning', round: true, bordered: false },
      { default: () => (is_read ? '已读' : '未读') },
    ),
  },
  {
    title: '发送人',
    key: 'admin',
    width: 120,
    render: ({ admin, admin_id }) => (admin_id ? (admin?.nickName || '-') : '系统'),
  },
  {
    title: '发送时间',
    key: 'created_at',
    width: 180,
    render: row => formatDateTime(row.created_at),
  },
  {
    title: '操作',
    key: 'actions',
    width: 160,
    align: 'center',
    fixed: 'right',
    hideInExcel: true,
    render(row) {
      return h('div', { class: 'flex gap-8 justify-center' }, [
        h(NButton, {
          size: 'small',
          type: 'primary',
          secondary: true,
          onClick: () => openView(row),
        }, { default: () => '查看' }),
        h(NButton, {
          size: 'small',
          type: 'error',
          secondary: true,
          onClick: () => handleDelete(row.id),
        }, { default: () => '删除' }),
      ])
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
      :get-data="api.message.list"
    >
      <CityQueryItem v-model="queryItems.city_id" :label-width="50" />
      <MeQueryItem label="学员" :label-width="50">
        <NInput
          v-model:value="queryItems.nick_name"
          placeholder="学员昵称"
          clearable
        />
      </MeQueryItem>
      <MeQueryItem label="标题" :label-width="50">
        <NInput
          v-model:value="queryItems.title"
          placeholder="标题关键词"
          clearable
        />
      </MeQueryItem>
      <MeQueryItem label="类型" :label-width="50">
        <NSelect
          v-model:value="queryItems.type"
          :options="typeOptions"
          placeholder="全部"
          clearable
        />
      </MeQueryItem>
      <MeQueryItem label="阅读" :label-width="50">
        <NSelect
          v-model:value="queryItems.is_read"
          :options="readOptions"
          placeholder="全部"
          clearable
        />
      </MeQueryItem>
    </MeCrud>

    <NModal
      v-model:show="showView"
      preset="card"
      title="站内信详情"
      class="w-720px"
      :bordered="false"
      @after-leave="viewRow = null"
    >
      <template v-if="viewRow">
        <h2 class="msg-title">
          {{ viewRow.title }}
        </h2>
        <div class="msg-meta">
          <span>
            学员：
            <button
              v-if="viewRow.user_id"
              type="button"
              class="entity-link"
              @click="router.push({ name: 'UserDetail', params: { id: viewRow.user_id } })"
            >
              {{ viewRow.user?.nick_name || '-' }}
            </button>
            <template v-else>{{ viewRow.user?.nick_name || '-' }}</template>
          </span>
          <span>{{ typeMap[viewRow.type]?.name || '通知' }}</span>
          <span>{{ viewRow.is_read ? '已读' : '未读' }}</span>
          <span>{{ formatDateTime(viewRow.created_at) }}</span>
        </div>
        <div class="msg-content rich-content" v-html="formatRichContentHtml(viewRow.content)" />
      </template>
    </NModal>
  </CommonPage>
</template>

<style scoped>
.msg-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 8px 10px;
  margin: 0 0 16px;
  padding: 0;
  background: transparent;
  border: 0;
  box-shadow: none;
}

.msg-meta span {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 3px 10px;
  border-radius: 999px;
  background: rgb(255 255 255 / 55%);
  border: 1px solid var(--ctl-hairline, rgba(var(--primary-color), 0.08));
  color: var(--ctl-text-soft, #5b6b7c);
  font-size: 12px;
  font-weight: 540;
  letter-spacing: -0.01em;
}

.msg-meta .entity-link {
  background: none;
  border: 0;
  padding: 0;
  font: inherit;
}

.msg-title {
  margin: 0 0 12px;
  font-family: var(--font-display);
  font-size: 22px;
  font-weight: 600;
  letter-spacing: -0.026em;
  line-height: 1.25;
  color: var(--ctl-text, #1a2332);
}

.msg-content {
  word-break: break-word;
  padding: 14px 16px;
  border-radius: 14px;
  background: var(--ctl-field, rgba(var(--primary-color), 0.03));
  border: 1px solid var(--ctl-hairline, rgba(var(--primary-color), 0.08));
  box-shadow: inset 0 1px 0 rgb(255 255 255 / 40%);
}
</style>
