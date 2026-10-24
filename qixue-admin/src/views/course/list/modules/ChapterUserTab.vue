<script setup>
import { NAvatar, NButton, NTag } from 'naive-ui'
import { computed, h, nextTick, onMounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import api from '@/api/index.js'
import { CityQueryItem, MeCrud, MeQueryItem, SoftUserCard } from '@/components'
import { getAssetUrl } from '@/utils'

const props = defineProps({
  chapterId: {
    type: [String, Number],
    required: true,
  },
})

function emptyQuery(chapterId) {
  return {
    chapter_id: chapterId,
    nick_name: null,
    mobile: null,
    city_id: null,
  }
}

const queryItems = ref(emptyQuery(props.chapterId))

const router = useRouter()
const unlockingIds = ref(new Set())
const $table = ref(null)

const resolvedChapterId = computed(() => {
  const fromProp = props.chapterId
  const fromQuery = queryItems.value?.chapter_id
  if (fromProp != null && fromProp !== '')
    return fromProp
  if (fromQuery != null && fromQuery !== '')
    return fromQuery
  return null
})

/** 请求时强制带上 chapter_id，避免筛选重置 / 异步时丢参 */
function fetchChapterUsers(params = {}) {
  const chapterId = resolvedChapterId.value
  if (chapterId == null || chapterId === '')
    return Promise.reject(new Error('课节参数缺失'))
  return api.chapter.chapterUser({
    ...params,
    chapter_id: chapterId,
  })
}

function ensureChapterIdInQuery() {
  const chapterId = props.chapterId
  if (chapterId == null || chapterId === '')
    return false
  if (String(queryItems.value.chapter_id) !== String(chapterId))
    queryItems.value = { ...queryItems.value, chapter_id: chapterId }
  return true
}

function searchUsers(resetPage = false) {
  if (!ensureChapterIdInQuery())
    return
  nextTick(() => {
    if (resetPage)
      $table.value?.handleSearch(false)
    else
      $table.value?.handleSearch()
  })
}

onMounted(() => {
  searchUsers()
})

watch(
  () => props.chapterId,
  (id, prev) => {
    if (id == null || id === '')
      return
    if (prev != null && prev !== '' && String(id) === String(prev)) {
      ensureChapterIdInQuery()
      return
    }
    queryItems.value = emptyQuery(id)
    searchUsers(true)
  },
)

defineExpose({
  refresh: () => searchUsers(true),
})

function openUser(row) {
  router.push({ name: 'UserDetail', params: { id: row.user_id } })
}

function isUnlocking(id) {
  return unlockingIds.value.has(id)
}

async function doUnlock(row) {
  if (!row?.id || Number(row.un_lock) === 1 || isUnlocking(row.id))
    return
  unlockingIds.value = new Set([...unlockingIds.value, row.id])
  try {
    await api.chapter.unlockChapterUser(row.id)
    $message.success('解锁成功')
    $table.value?.handleSearch(false)
  }
  catch (error) {
    $message.error(error?.message || '解锁失败')
  }
  finally {
    const next = new Set(unlockingIds.value)
    next.delete(row.id)
    unlockingIds.value = next
  }
}

function confirmUnlock(row, e) {
  e?.stopPropagation?.()
  if (!row || Number(row.un_lock) === 1)
    return
  const name = row.user?.nick_name || '该学员'
  $dialog.warning({
    title: '帮学员解锁',
    content: `确定为「${name}」解锁本课节吗？`,
    positiveText: '解锁',
    negativeText: '取消',
    onPositiveClick: () => doUnlock(row),
  })
}

const columns = [
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
            size: 40,
            round: true,
            src: getAssetUrl(user?.avatar) || undefined,
            class: 'flex-shrink-0',
          }),
          h('div', { class: 'text-16 font-bold text-primary ml-10' }, user?.nick_name || '未命名'),
        ],
      )
    },
  },
  {
    title: '解锁',
    key: 'un_lock',
    align: 'center',
    width: 100,
    render: ({ un_lock }) => h(
      NTag,
      {
        type: un_lock ? 'success' : 'warning',
        size: 'small',
        round: true,
        bordered: false,
      },
      un_lock ? '已解锁' : '未解锁',
    ),
  },
  {
    title: '状态',
    key: 'is_end',
    align: 'center',
    width: 100,
    render: ({ is_end }) => h(
      NTag,
      {
        type: is_end ? 'success' : 'error',
        size: 'small',
        round: true,
        bordered: false,
      },
      is_end ? '已学完' : '未学完',
    ),
  },
  {
    title: '加入时间',
    key: 'created_at',
    align: 'center',
  },
  {
    title: '操作',
    key: 'actions',
    width: 108,
    align: 'center',
    fixed: 'right',
    hideInExcel: true,
    render(row) {
      if (Number(row.un_lock) === 1) {
        return h('span', { class: 'chapter-user-unlocked' }, '已解锁')
      }
      return h('div', { class: 'row-actions' }, [
        h(
          NButton,
          {
            size: 'small',
            type: 'primary',
            secondary: true,
            loading: isUnlocking(row.id),
            onClick: () => confirmUnlock(row),
          },
          {
            default: () => '解锁',
            icon: () => h('i', { class: 'i-fe:unlock text-14' }),
          },
        ),
      ])
    },
  },
]
</script>

<template>
  <div class="chapter-user-tab">
    <MeCrud
      ref="$table"
      v-model:query-items="queryItems"
      :bound-query="{ chapter_id: chapterId }"
      :columns="columns"
      :get-data="fetchChapterUsers"
    >
      <CityQueryItem v-model="queryItems.city_id" :label-width="50" />
      <MeQueryItem label="用户名" :label-width="50">
        <n-input
          v-model:value="queryItems.nick_name"
          type="text"
          placeholder="请输入用户名"
          clearable
        />
      </MeQueryItem>
      <MeQueryItem label="电话" :label-width="50">
        <n-input
          v-model:value="queryItems.mobile"
          type="text"
          placeholder="请输入电话"
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
            :subtitle="Number(item.un_lock) === 1 ? '已解锁' : '未解锁'"
            :status="item.is_end ? '已学完' : '未学完'"
            :status-tone="item.is_end ? 'done' : 'doing'"
            :time="item.created_at"
            @click="openUser(item)"
          >
            <template #body>
              <div class="chapter-user-lock" :class="Number(item.un_lock) === 1 ? 'is-on' : 'is-off'">
                <i :class="Number(item.un_lock) === 1 ? 'i-fe:unlock' : 'i-fe:lock'" />
                <span>{{ Number(item.un_lock) === 1 ? '课节已解锁' : '课节未解锁' }}</span>
              </div>
            </template>
            <template v-if="Number(item.un_lock) !== 1" #actions>
              <NButton
                size="tiny"
                type="primary"
                secondary
                :loading="isUnlocking(item.id)"
                @click="confirmUnlock(item, $event)"
              >
                <template #icon>
                  <i class="i-fe:unlock text-14" />
                </template>
                帮TA解锁
              </NButton>
            </template>
          </SoftUserCard>
        </div>
      </template>
    </MeCrud>
  </div>
</template>

<style scoped>
.chapter-user-tab {
  height: 500px;
  min-height: 0;
}

.chapter-user-tab :deep(.chapter-user-unlocked) {
  font-size: 12px;
  font-weight: 520;
  color: #8b97a5;
}

.chapter-user-lock {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  margin-top: 12px;
  height: 28px;
  padding: 0 10px;
  border-radius: 9px;
  font-size: 12px;
  font-weight: 560;
  letter-spacing: -0.01em;
}

.chapter-user-lock i {
  font-size: 14px;
}

.chapter-user-lock.is-on {
  color: #2d7a45;
  background: rgb(52 199 89 / 12%);
  border: 1px solid rgb(52 199 89 / 16%);
}

.chapter-user-lock.is-off {
  color: #8a5a00;
  background: rgb(255 159 10 / 12%);
  border: 1px solid rgb(255 159 10 / 16%);
}
</style>
