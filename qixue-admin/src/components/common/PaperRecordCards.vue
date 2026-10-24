<script setup>
import { NButton } from 'naive-ui'
import { useRouter } from 'vue-router'
import SoftUserCard from './SoftUserCard.vue'
import { getAssetUrl } from '@/utils'

const props = defineProps({
  list: {
    type: Array,
    default: () => [],
  },
  showChapter: {
    type: Boolean,
    default: false,
  },
  showUser: {
    type: Boolean,
    default: true,
  },
  showMessage: {
    type: Boolean,
    default: true,
  },
})

const emit = defineEmits(['open', 'message'])
const router = useRouter()

function photoOf(row) {
  const src = row.user?.avatar || row.avatar
  return src ? getAssetUrl(src) : ''
}

function userName(row) {
  return row.user?.nick_name || row.nick_name || ''
}

function courseTitle(row) {
  return row.paper?.course?.title || ''
}

function chapterTitle(row) {
  return row.paper?.chapter?.title || ''
}

function titleOf(row) {
  if (props.showUser)
    return userName(row) || '未命名学员'
  return courseTitle(row) || chapterTitle(row) || '未命名'
}

function subtitleOf(row) {
  if (props.showUser) {
    const parts = []
    if (courseTitle(row))
      parts.push(courseTitle(row))
    if (props.showChapter && chapterTitle(row))
      parts.push(chapterTitle(row))
    return parts.join(' · ')
  }
  if (props.showChapter && chapterTitle(row))
    return chapterTitle(row)
  return row.is_pass ? '通过' : '未通过'
}

function courseIdOf(row) {
  return row.paper?.course_id || row.paper?.course?.id
}

function chapterIdOf(row) {
  return row.paper?.chapter_id || row.paper?.chapter?.id
}

function isGraded(row) {
  return !row.is_people_check
}

function statusOf(row) {
  return row.is_people_check ? '待批改' : '已批改'
}

function statusToneOf(row) {
  return row.is_people_check ? 'warn' : 'done'
}

function metricsOf(row) {
  const list = [
    { value: row.score ?? '—', label: '成绩' },
    { value: row.is_pass ? '通过' : '未通过', label: '结果' },
  ]
  if (row.submit_count != null)
    list.push({ value: row.submit_count, label: '提交次数' })
  return list
}

function open(row, viewOnly = isGraded(row)) {
  emit('open', row, viewOnly)
}

function goCourse(row, e) {
  e?.stopPropagation?.()
  const id = courseIdOf(row)
  if (id)
    router.push(`/course/detail/${id}`)
}

function goChapter(row, e) {
  e?.stopPropagation?.()
  const courseId = courseIdOf(row)
  const chapterId = chapterIdOf(row)
  if (courseId && chapterId)
    router.push(`/course/${courseId}/chapterDetail/${chapterId}`)
}

function onCardClick(row) {
  open(row)
}
</script>

<template>
  <div class="entity-grid">
    <SoftUserCard
      v-for="item in list"
      :key="item.id"
      class="paper-record-card"
      :class="{ 'is-todo': !!item.is_people_check }"
      :avatar="showUser ? photoOf(item) : ''"
      :name="titleOf(item)"
      :subtitle="subtitleOf(item)"
      :status="statusOf(item)"
      :status-tone="statusToneOf(item)"
      :time="item.created_at"
      :metrics="metricsOf(item)"
      @click="onCardClick(item)"
    >
      <template #body>
        <div
          v-if="(showUser && courseTitle(item)) || (showChapter && chapterTitle(item))"
          class="paper-record-card__rels"
          @click.stop
        >
          <button
            v-if="showUser && courseTitle(item)"
            type="button"
            class="paper-record-card__chip"
            :title="courseTitle(item)"
            @click="goCourse(item, $event)"
          >
            {{ courseTitle(item) }}
          </button>
          <button
            v-if="showChapter && chapterTitle(item)"
            type="button"
            class="paper-record-card__chip"
            :title="chapterTitle(item)"
            @click="goChapter(item, $event)"
          >
            {{ chapterTitle(item) }}
          </button>
        </div>
      </template>

      <template #actions>
        <NButton
          size="tiny"
          :type="isGraded(item) ? 'default' : 'primary'"
          :secondary="isGraded(item)"
          @click="open(item)"
        >
          {{ isGraded(item) ? '查看' : '批改' }}
        </NButton>
        <NButton
          v-if="showMessage"
          size="tiny"
          type="info"
          secondary
          @click="emit('message', item)"
        >
          站内信
        </NButton>
      </template>
    </SoftUserCard>
  </div>
</template>

<style scoped>
.paper-record-card__rels {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-top: 10px;
}

.paper-record-card__chip {
  appearance: none;
  max-width: 100%;
  height: 22px;
  padding: 0 8px;
  border: 0;
  border-radius: 999px;
  font-size: 11px;
  font-weight: 560;
  letter-spacing: -0.01em;
  color: #5b6b7c;
  background: rgba(var(--primary-color), 0.08);
  cursor: pointer;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  transition: background 0.16s ease, color 0.16s ease;
}

.paper-record-card__chip:hover {
  color: rgb(var(--primary-color));
  background: rgba(var(--primary-color), 0.14);
}
</style>
