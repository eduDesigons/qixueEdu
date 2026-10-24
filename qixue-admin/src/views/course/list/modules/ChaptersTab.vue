<script setup>
import { ArrowDown, ArrowUp } from '@vicons/ionicons5'
import { NIcon, NInput, NTag } from 'naive-ui'
import { h, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import api from '@/api/index.js'
import { MeCrud, MeQueryItem } from '@/components'

const props = defineProps({
  courseId: {
    type: [String, Number],
    required: true,
  },
})

const router = useRouter()
const $table = ref(null)
const tableData = ref([])
const queryItems = ref({
  courseId: props.courseId,
  title: '',
})

function handleViewChapter(row) {
  router.push(`/course/${props.courseId}/chapterDetail/${row.id}`)
}

async function handleSort(row, direction) {
  try {
    await api.chapter.chapterSort(row.id, direction)
    window.$message.success(direction === 'up' ? '上移成功' : '下移成功')
    $table.value?.handleSearch()
  }
  catch (error) {
    console.error('排序失败:', error)
    window.$message.error(error?.message || '操作失败，请重试')
  }
}

function handleDataChange(data) {
  tableData.value = data || []
}

function sortBtn(row, direction, disabled) {
  return h('button', {
    type: 'button',
    class: [
      'chapter-sort-btn',
      direction === 'up' ? 'is-up' : 'is-down',
      disabled ? 'is-disabled' : '',
    ],
    disabled,
    title: direction === 'up' ? '上移' : '下移',
    onClick: (e) => {
      e.stopPropagation()
      if (!disabled)
        handleSort(row, direction)
    },
  }, [
    h(NIcon, { size: 14 }, {
      default: () => h(direction === 'up' ? ArrowUp : ArrowDown),
    }),
  ])
}

function sortIndexLabel(value) {
  const n = Number(value)
  if (!Number.isFinite(n) || n <= 0)
    return '—'
  return String(n).padStart(2, '0')
}

const columns = [
  {
    title: '序号',
    key: 'sort',
    width: 72,
    align: 'center',
    render: row => h('span', { class: 'chapter-idx' }, row.sort ?? '—'),
  },
  {
    title: '名称',
    key: 'title',
    resizable: true,
    render: row => h('span', {
      class: 'entity-link chapter-title',
      onClick: () => handleViewChapter(row),
    }, row.title),
  },
  {
    title: '学完',
    key: 'user_chapter_end_count',
    align: 'center',
    width: 80,
    render: ({ user_chapter_end_count }) => h('span', { class: 'chapter-metric' }, user_chapter_end_count || 0),
  },
  {
    title: '学员',
    key: 'users_count',
    align: 'center',
    width: 80,
    render: ({ users_count }) => h('span', { class: 'chapter-metric' }, users_count || 0),
  },
  {
    title: '评论',
    key: 'comments_count',
    align: 'center',
    width: 80,
    render: ({ comments_count }) => h('span', { class: 'chapter-metric' }, comments_count || 0),
  },
  {
    title: '时限',
    key: 'limit_day',
    align: 'center',
    width: 96,
    render: ({ limit_day }) => limit_day
      ? h(NTag, { size: 'small', bordered: false, type: 'warning', round: true }, { default: () => `${limit_day} 天` })
      : h('span', { class: 'chapter-muted' }, '不限'),
  },
  {
    title: '测评',
    key: 'paper_id',
    align: 'center',
    width: 96,
    render: ({ paper_id }) => paper_id
      ? h(NTag, { size: 'small', bordered: false, type: 'success', round: true }, { default: () => '有测评' })
      : h('span', { class: 'chapter-muted' }, '无测评'),
  },
  {
    title: '排序',
    key: 'action',
    align: 'center',
    width: 84,
    fixed: 'right',
    render: (row, index) => {
      const isFirst = index === 0
      const isLast = index === tableData.value.length - 1
      return h('div', { class: 'row-actions chapter-sort-actions' }, [
        sortBtn(row, 'up', isFirst),
        sortBtn(row, 'down', isLast),
      ])
    },
  },
]

onMounted(() => {
  $table.value?.handleSearch()
})

defineExpose({
  refresh: () => $table.value?.handleSearch(false),
})
</script>

<template>
  <div class="chapters-tab">
    <MeCrud
      ref="$table"
      v-model:query-items="queryItems"
      :is-pagination="false"
      :columns="columns"
      :bound-query="{ courseId }"
      :get-data="api.chapter.chapterList"
      @on-data-change="handleDataChange"
    >
      <MeQueryItem label="课节" :label-width="50">
        <NInput
          v-model:value="queryItems.title"
          clearable
          placeholder="搜索课节名称"
          class="w-200"
          @keyup.enter="$table?.handleSearch()"
        />
      </MeQueryItem>

      <template #card="{ list }">
        <div class="entity-grid chapter-card-grid">
          <article
            v-for="(item, index) in list"
            :key="item.id"
            class="entity-card chapter-card"
            @click="handleViewChapter(item)"
          >
            <span class="chapter-card__idx" :title="`第 ${item.sort ?? index + 1} 节`">
              {{ sortIndexLabel(item.sort ?? index + 1) }}
            </span>
            <span v-if="item.paper_id" class="chapter-card__exam">
              有测评
            </span>

            <div class="chapter-card__body">
              <h3 class="chapter-card__title" :title="item.title">
                {{ item.title || '未命名课节' }}
              </h3>

              <div class="chapter-card__stats">
                <div class="chapter-card__stat">
                  <strong>{{ item.user_chapter_end_count || 0 }}</strong>
                  <span>学完</span>
                </div>
                <div class="chapter-card__stat">
                  <strong>{{ item.users_count || 0 }}</strong>
                  <span>学员</span>
                </div>
                <div class="chapter-card__stat">
                  <strong>{{ item.comments_count || 0 }}</strong>
                  <span>评论</span>
                </div>
              </div>

              <div class="chapter-card__foot">
                <span class="chapter-card__limit">
                  {{ item.limit_day ? `时限 ${item.limit_day} 天` : '时限不限' }}
                </span>
                <div
                  class="chapter-sort-actions"
                  @click.stop
                >
                  <button
                    type="button"
                    class="chapter-sort-btn"
                    :class="{ 'is-disabled': index === 0 }"
                    :disabled="index === 0"
                    title="上移"
                    @click="index !== 0 && handleSort(item, 'up')"
                  >
                    <NIcon :size="14">
                      <ArrowUp />
                    </NIcon>
                  </button>
                  <button
                    type="button"
                    class="chapter-sort-btn"
                    :class="{ 'is-disabled': index === list.length - 1 }"
                    :disabled="index === list.length - 1"
                    title="下移"
                    @click="index !== list.length - 1 && handleSort(item, 'down')"
                  >
                    <NIcon :size="14">
                      <ArrowDown />
                    </NIcon>
                  </button>
                </div>
              </div>
            </div>
          </article>
        </div>
      </template>
    </MeCrud>
  </div>
</template>

<style scoped lang="scss">
.chapters-tab {
  height: 500px;
  min-height: 0;
}

.chapters-tab :deep(.chapter-idx) {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 28px;
  height: 24px;
  padding: 0 8px;
  border-radius: 8px;
  font-size: 12px;
  font-weight: 600;
  letter-spacing: -0.02em;
  color: #5b6b7c;
  background: rgba(var(--primary-color), 0.06);
}

.chapters-tab :deep(.chapter-title) {
  font-weight: 560;
  letter-spacing: -0.014em;
}

.chapters-tab :deep(.chapter-metric) {
  font-size: 13px;
  font-weight: 620;
  letter-spacing: -0.02em;
  color: #1a2332;
  font-variant-numeric: tabular-nums;
}

.chapters-tab :deep(.chapter-muted) {
  font-size: 12px;
  font-weight: 500;
  color: #8b97a5;
}

.chapter-card {
  position: relative;
  display: flex;
  flex-direction: column;
  min-height: 168px;
  padding: 0;
  overflow: hidden;
  cursor: pointer;
}

.chapter-card__idx {
  position: absolute;
  top: 12px;
  left: 12px;
  z-index: 2;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 36px;
  height: 28px;
  padding: 0 10px;
  border-radius: 10px;
  font-size: 13px;
  font-weight: 700;
  font-variant-numeric: tabular-nums;
  letter-spacing: -0.03em;
  color: #fff;
  background:
    linear-gradient(
      145deg,
      rgba(var(--primary-color), 0.95) 0%,
      rgba(var(--primary-color), 0.72) 100%
    );
  box-shadow:
    0 1px 2px rgba(var(--primary-color), 0.18),
    0 6px 14px rgba(var(--primary-color), 0.16),
    inset 0 1px 0 rgb(255 255 255 / 28%);
}

.chapter-card__exam {
  position: absolute;
  top: 12px;
  right: 12px;
  z-index: 2;
  display: inline-flex;
  align-items: center;
  height: 24px;
  padding: 0 10px;
  border-radius: 999px;
  font-size: 11px;
  font-weight: 650;
  letter-spacing: -0.01em;
  color: #2d7a45;
  background: rgb(52 199 89 / 16%);
  border: 1px solid rgb(52 199 89 / 18%);
  backdrop-filter: blur(8px);
}

.chapter-card__body {
  position: relative;
  z-index: 1;
  display: flex;
  flex-direction: column;
  flex: 1;
  min-height: 0;
  padding: 52px 14px 12px;
}

.chapter-card__title {
  margin: 0;
  min-height: calc(1.35em * 2);
  font-size: 15px;
  font-weight: 650;
  letter-spacing: -0.02em;
  line-height: 1.35;
  color: #1a2332;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.chapter-card__stats {
  display: flex;
  align-items: stretch;
  margin-top: 14px;
  padding: 10px 0;
  border-top: 1px solid var(--ctl-hairline, rgba(var(--primary-color), 0.06));
  border-bottom: 1px solid var(--ctl-hairline, rgba(var(--primary-color), 0.06));
}

.chapter-card__stat {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
  text-align: center;
}

.chapter-card__stat + .chapter-card__stat {
  border-left: 1px solid var(--ctl-hairline, rgba(var(--primary-color), 0.06));
}

.chapter-card__stat strong {
  font-size: 15px;
  font-weight: 700;
  letter-spacing: -0.03em;
  color: #1a2332;
  font-variant-numeric: tabular-nums;
  line-height: 1.2;
}

.chapter-card__stat span {
  font-size: 11px;
  font-weight: 500;
  color: #8b97a5;
}

.chapter-card__foot {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  margin-top: auto;
  padding-top: 10px;
  min-height: 34px;
}

.chapter-card__limit {
  font-size: 11px;
  font-weight: 540;
  letter-spacing: -0.01em;
  color: #8b97a5;
}

.chapters-tab :deep(.chapter-sort-actions) {
  display: inline-flex !important;
  align-items: center;
  justify-content: center;
  gap: 2px !important;
  opacity: 1;
  transform: none;
  pointer-events: auto;
}

.chapter-card .chapter-sort-actions {
  display: inline-flex !important;
  align-items: center;
  justify-content: center;
  gap: 2px !important;
  opacity: 0;
  transform: translateY(3px) scale(0.96);
  transition:
    opacity 0.2s cubic-bezier(0.32, 0.72, 0, 1),
    transform 0.2s cubic-bezier(0.32, 0.72, 0, 1);
  pointer-events: none;
}

.chapter-card:hover .chapter-sort-actions,
.chapter-card:focus-within .chapter-sort-actions {
  opacity: 1;
  transform: none;
  pointer-events: auto;
}

.chapters-tab :deep(.chapter-sort-btn),
.chapter-sort-btn {
  appearance: none;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 26px;
  height: 26px;
  margin: 0;
  padding: 0;
  border: 0;
  border-radius: 8px;
  color: #5b6b7c;
  background: transparent;
  cursor: pointer;
  transition:
    color 0.16s ease,
    background 0.16s ease,
    transform 0.16s cubic-bezier(0.32, 0.72, 0, 1);
}

.chapters-tab :deep(.chapter-sort-btn:hover:not(.is-disabled)),
.chapter-sort-btn:hover:not(.is-disabled) {
  color: rgb(var(--primary-color));
  background: rgba(var(--primary-color), 0.1);
}

.chapters-tab :deep(.chapter-sort-btn:active:not(.is-disabled)),
.chapter-sort-btn:active:not(.is-disabled) {
  transform: scale(0.92);
}

.chapters-tab :deep(.chapter-sort-btn.is-disabled),
.chapter-sort-btn.is-disabled {
  color: #c5cecf;
  cursor: not-allowed;
}

@media (prefers-reduced-motion: reduce) {
  .chapters-tab :deep(.chapter-sort-actions),
  .chapters-tab :deep(.chapter-sort-btn),
  .chapter-sort-actions,
  .chapter-sort-btn {
    transition: none;
  }
}
</style>
