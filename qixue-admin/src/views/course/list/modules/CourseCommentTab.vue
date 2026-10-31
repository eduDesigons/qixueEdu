<script setup>
import { getAssetUrl, formatRelativeTime } from '@/utils'
import { NAvatar, NButton, NInput, NModal, NSpin } from 'naive-ui'
import { h, nextTick, watch } from 'vue'
import { useRouter } from 'vue-router'
import api from '@/api/index.js'
import { MeCrud, MeQueryItem, CityQueryItem, CommentCardGrid } from '@/components'
import { useCrud } from '@/composables/index.js'

const props = defineProps({
  courseId: {
    type: [String, Number],
    required: true,
  },
})

const queryItems = ref({
  course_id: props.courseId,
  content: '',
  city_id: null,
})

const replyModalShow = ref(false)
const replyContent = ref('')
const replyCommentId = ref(0)
const replyLoading = ref(false)
const replyInput = ref(null)
const repliesModalShow = ref(false)
const repliesLoading = ref(false)
const replies = ref([])
const router = useRouter()
const $table = ref(null)
const checkedIds = ref([])
const batchLoading = ref(false)
const crudKey = ref(0)

onMounted(() => {
  $table.value?.handleSearch()
})

const { handleDelete } = useCrud({
  name: '评论',
  doDelete: api.course.courseCommentDelete,
  refresh: () => $table.value?.handleSearch(),
})

function onChecked(rowKeys) {
  checkedIds.value = rowKeys || []
}

function openBatchDelete() {
  if (!checkedIds.value.length) {
    $message.warning('请先勾选要删除的评论')
    return
  }
  $dialog.warning({
    title: '批量删除',
    content: `确定删除选中的 ${checkedIds.value.length} 条评论？其回复也会一并删除。`,
    positiveText: '删除',
    negativeText: '取消',
    onPositiveClick: async () => {
      batchLoading.value = true
      try {
        await api.course.courseCommentBatchDelete(checkedIds.value)
        $message.success('删除成功')
        checkedIds.value = []
        crudKey.value += 1
        await nextTick()
        $table.value?.handleSearch()
      }
      catch (e) {
        $message.error(e?.message || '删除失败')
      }
      finally {
        batchLoading.value = false
      }
    },
  })
}

function openReply(row) {
  replyCommentId.value = row.id
  replyContent.value = ''
  replyModalShow.value = true
}

async function openReplies(row) {
  repliesModalShow.value = true
  repliesLoading.value = true
  replies.value = []
  try {
    const res = await api.course.courseCommentReplies({ comment_id: row.id })
    replies.value = res.data || []
  }
  catch (e) {
    $message.error(e?.message || '加载回复失败')
  }
  finally {
    repliesLoading.value = false
  }
}

function closeReplyModal() {
  replyModalShow.value = false
  replyCommentId.value = 0
  replyContent.value = ''
}

async function submitReply() {
  const content = replyContent.value?.trim()
  if (!content) {
    $message.warning('请输入回复内容')
    return
  }
  replyLoading.value = true
  try {
    await api.course.courseCommentReply({ comment_id: replyCommentId.value, content })
    $message.success('回复成功')
    closeReplyModal()
    $table.value?.handleSearch()
  }
  catch (e) {
    $message.error(e?.message || '回复失败')
  }
  finally {
    replyLoading.value = false
  }
}

watch(replyModalShow, (show) => {
  if (show)
    nextTick(() => replyInput.value?.focus?.())
})

function goUser(item) {
  if (item.user_id === 0)
    return
  router.push({ name: 'UserDetail', params: { id: item.user_id } })
}

defineExpose({
  refresh: () => $table.value?.handleSearch(false),
})

const columns = [
  { type: 'selection', fixed: 'left' },
  {
    title: '编号',
    key: 'id',
    width: 100,
    resizable: true,
  },
  {
    title: '学员',
    key: 'user',
    render({ user, user_id }) {
      return h(
        'div',
        {
          class: 'flex cursor-pointer items-center',
          onClick: () => router.push({ name: 'UserDetail', params: { id: user_id } }),
        },
        [
          h(NAvatar, {
            size: 36,
            src: getAssetUrl(user.avatar) || undefined,
          }),
          h('div', { class: 'text-16 font-bold text-primary ml-10' }, user.nick_name),
        ],
      )
    },
  },
  {
    title: '评论内容',
    key: 'content',
    resizable: true,
  },
  {
    title: '回复数',
    key: 'children_count',
    align: 'center',
    render: ({ children_count }) => h('span', { class: 'font-bold' }, children_count),
  },
  {
    title: '评论时间',
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
    render(row) {
      return [
        row.children_count > 0 && h(
          NButton,
          {
            size: 'small',
            type: 'info',
            secondary: true,
            style: 'margin-right: 8px;',
            onClick: () => openReplies(row),
          },
          { default: () => '查看回复' },
        ),
        h(
          NButton,
          {
            size: 'small',
            type: 'primary',
            secondary: true,
            style: 'margin-right: 8px;',
            onClick: () => openReply(row),
          },
          {
            default: () => '回复',
            icon: () => h('i', { class: 'i-material-symbols:reply text-14' }),
          },
        ),
        h(
          NButton,
          {
            size: 'small',
            type: 'error',
            secondary: true,
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
    <NModal
      v-model:show="repliesModalShow"
      preset="card"
      title="评论回复"
      style="width: 640px"
      class="replies-modal"
      @after-leave="replies = []"
    >
      <NSpin :show="repliesLoading">
        <div v-if="replies.length" class="replies-list">
          <div v-for="item in replies" :key="item.id" class="replies-item">
            <div class="replies-item__head">
              <strong>{{ item.user?.nick_name || '学员' }}</strong>
              <span>{{ formatRelativeTime(item.created_at) }}</span>
            </div>
            <div class="replies-item__body">
              {{ item.content }}
            </div>
          </div>
        </div>
        <div v-else-if="!repliesLoading" class="empty-state replies-empty">
          暂无回复
        </div>
      </NSpin>
    </NModal>
    <NModal
      v-model:show="replyModalShow"
      preset="dialog"
      title="回复评论"
      positive-text="提交"
      negative-text="取消"
      :loading="replyLoading"
      :close-on-positive-click="false"
      @positive-click="submitReply"
      @after-leave="closeReplyModal"
    >
      <NInput
        ref="replyInput"
        v-model:value="replyContent"
        type="textarea"
        placeholder="请输入回复内容"
        :rows="4"
        :autosize="{ minRows: 4, maxRows: 8 }"
      />
    </NModal>
    <div v-if="checkedIds.length" class="mb-12">
      <NButton type="error" secondary :loading="batchLoading" :disabled="batchLoading" @click="openBatchDelete">
        批量删除
      </NButton>
    </div>
    <MeCrud
      :key="crudKey"
      ref="$table"
      v-model:query-items="queryItems"
      :bound-query="{ course_id: courseId }"
      :columns="columns"
      default-expand-all
      :get-data="api.course.courseComment"
      @on-checked="onChecked"
      @onDataChange="checkedIds = []"
    >
      <CityQueryItem v-model="queryItems.city_id" />
      <MeQueryItem label="评论内容" :label-width="80">
        <NInput
          v-model:value="queryItems.content"
          clearable
          placeholder="请输入评论关键字"
        />
      </MeQueryItem>
      <template #card="{ list }">
        <CommentCardGrid
          :list="list"
          v-model:checked-ids="checkedIds"
          @user="goUser"
          @replies="openReplies"
          @reply="openReply"
          @delete="item => handleDelete(item.id)"
        />
      </template>
    </MeCrud>
  </div>
</template>

<style scoped>
.replies-list {
  max-height: 480px;
  overflow: auto;
  display: flex;
  flex-direction: column;
  gap: 1px;
  padding: 4px;
  border-radius: 14px;
  background: var(--ctl-field, rgba(var(--primary-color), 0.04));
  border: 1px solid var(--ctl-hairline, rgba(var(--primary-color), 0.08));
}

.replies-item {
  padding: 12px 14px;
  border-radius: 11px;
  background: rgb(255 255 255 / 55%);
}

.replies-item__head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 8px;
}

.replies-item__head strong {
  font-size: 13px;
  font-weight: 600;
  color: var(--ctl-text, #1a2332);
}

.replies-item__head span {
  font-size: 12px;
  color: var(--ctl-text-soft, #7a8b8d);
}

.replies-item__body {
  font-size: 13px;
  line-height: 1.65;
  color: #3d4f51;
  white-space: pre-wrap;
  word-break: break-word;
}
</style>
