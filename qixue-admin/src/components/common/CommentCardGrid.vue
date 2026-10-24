<script setup>
import { computed } from 'vue'
import { NCheckbox } from 'naive-ui'
import { getAssetUrl, formatRelativeTime } from '@/utils'

const props = defineProps({
  list: {
    type: Array,
    default: () => [],
  },
  checkedIds: {
    type: Array,
    default: () => [],
  },
  showCourse: {
    type: Boolean,
    default: false,
  },
  showChapter: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['update:checkedIds', 'user', 'course', 'chapter', 'replies', 'reply', 'delete'])

function isChecked(id) {
  return props.checkedIds.includes(id)
}

const pageSelectedCount = computed(() => props.list.filter(item => isChecked(item.id)).length)
const pageChecked = computed(() => props.list.length > 0 && props.list.every(item => isChecked(item.id)))
const pageIndeterminate = computed(() => {
  const n = pageSelectedCount.value
  return n > 0 && n < props.list.length
})

function setCheckedIds(next) {
  emit('update:checkedIds', next)
}

function toggleOne(id, checked) {
  if (checked) {
    if (isChecked(id))
      return
    setCheckedIds([...props.checkedIds, id])
    return
  }
  setCheckedIds(props.checkedIds.filter(item => item !== id))
}

const pageCheck = computed({
  get: () => pageChecked.value,
  set: (checked) => {
    const on = !!checked
    const ids = props.list.map(item => item.id)
    if (on)
      setCheckedIds([...new Set([...props.checkedIds, ...ids])])
    else
      setCheckedIds(props.checkedIds.filter(id => !ids.includes(id)))
  },
})
</script>

<template>
  <div class="cmt-wrap">
    <div class="entity-select-bar">
      <NCheckbox
        v-model:checked="pageCheck"
        :indeterminate="pageIndeterminate"
      >
        本页全选
      </NCheckbox>
      <span v-if="pageSelectedCount" class="entity-select-bar__count">已选 {{ pageSelectedCount }} 条</span>
    </div>

    <div class="cmt-grid">
      <article
        v-for="item in list"
        :key="item.id"
        class="cmt-card"
        :class="{ 'is-checked': isChecked(item.id), 'has-photo': !!getAssetUrl(item.user?.avatar || '') }"
        @click="toggleOne(item.id, !isChecked(item.id))"
      >
        <img
          v-if="getAssetUrl(item.user?.avatar || '')"
          class="photo-bg"
          :src="getAssetUrl(item.user?.avatar || '')"
          alt=""
        >
        <div v-if="getAssetUrl(item.user?.avatar || '')" class="photo-veil" />
        <div class="cmt-card__check" @click.stop>
          <NCheckbox
            :checked="isChecked(item.id)"
            @update:checked="toggleOne(item.id, $event)"
          />
        </div>

        <header class="cmt-card__head">
          <div class="cmt-card__who">
            <button
              v-if="item.user_id !== 0"
              type="button"
              class="cmt-card__name entity-link"
              @click.stop="emit('user', item)"
            >
              {{ item.user?.nick_name || '-' }}
            </button>
            <span v-else class="cmt-card__name is-admin">{{ item.user?.nick_name || '管理员' }}</span>
            <p class="cmt-card__when">
              {{ formatRelativeTime(item.created_at) }}
            </p>
          </div>
        </header>

        <div
          v-if="item.user_id === 0 || (showChapter && item.chapter?.title) || (showCourse && item.course?.title)"
          class="cmt-card__chips"
        >
          <span v-if="item.user_id === 0" class="entity-card__badge">管理员</span>
          <button
            v-if="showChapter && item.chapter?.title"
            type="button"
            class="entity-card__chip"
            @click.stop="emit('chapter', item)"
          >
            {{ item.chapter.title }}
          </button>
          <button
            v-if="showCourse && item.course?.title"
            type="button"
            class="entity-card__chip"
            @click.stop="emit('course', item)"
          >
            {{ item.course.title }}
          </button>
        </div>

        <p class="cmt-card__body" :title="item.content?.trim() || ''">
          {{ item.content?.trim() || '暂无内容' }}
        </p>

        <footer class="cmt-card__foot" @click.stop>
          <button
            v-if="item.children_count > 0"
            type="button"
            class="cmt-card__link"
            @click="emit('replies', item)"
          >
            {{ item.children_count }} 条回复
          </button>
          <span v-else class="cmt-card__muted">暂无回复</span>
          <span class="cmt-card__spacer" />
          <button type="button" class="cmt-card__link" @click="emit('reply', item)">
            回复
          </button>
          <button type="button" class="cmt-card__link is-danger" @click="emit('delete', item)">
            删除
          </button>
        </footer>
      </article>
    </div>
  </div>
</template>

<style scoped>
.cmt-wrap {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.cmt-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 14px;
  padding-bottom: 4px;
}

.cmt-card {
  position: relative;
  display: flex;
  flex-direction: column;
  min-height: 168px;
  padding: 16px 16px 10px;
  overflow: hidden;
  border-radius: 16px;
  border: 1px solid var(--ctl-hairline, rgba(var(--primary-color), 0.08));
  cursor: pointer;
  text-align: left;
}

.cmt-card.is-checked {
  border-color: rgba(var(--primary-color), 0.28);
}

.cmt-card__check {
  position: absolute;
  top: 12px;
  right: 12px;
  z-index: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  border-radius: 8px;
  background: rgb(255 255 255 / 88%);
}

.cmt-card__head {
  display: flex;
  align-items: center;
  gap: 10px;
  padding-right: 32px;
  min-width: 0;
}

.cmt-card__head :deep(.n-avatar) {
  width: 40px !important;
  height: 40px !important;
}

.cmt-card__who {
  min-width: 0;
  flex: 1;
}

.cmt-card__name {
  display: block;
  max-width: 100%;
  padding: 0;
  border: 0;
  background: transparent;
  font-size: 14px;
  font-weight: 650;
  letter-spacing: -0.018em;
  color: rgb(var(--primary-color));
  text-align: left;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  cursor: pointer;
}

.cmt-card__name.is-admin {
  color: #1a2332;
  cursor: default;
}

.cmt-card__when {
  margin: 2px 0 0;
  font-size: 11px;
  letter-spacing: -0.01em;
  color: #8b97a5;
}

.cmt-card__chips {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-top: 10px;
}

.cmt-card__chip {
  max-width: 100%;
  height: 22px;
  padding: 0 8px;
  border: 0;
  border-radius: 999px;
  background: rgb(255 255 255 / 78%);
  font-size: 11px;
  font-weight: 550;
  letter-spacing: -0.01em;
  color: #5b6b7c;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  cursor: pointer;
}

.cmt-card__body {
  flex: 1;
  margin: 12px 0 0;
  font-size: 14px;
  line-height: 1.6;
  letter-spacing: -0.014em;
  color: #1a2332;
  display: -webkit-box;
  -webkit-line-clamp: 4;
  -webkit-box-orient: vertical;
  overflow: hidden;
  word-break: break-word;
}

.cmt-card__foot {
  display: flex;
  align-items: center;
  gap: 4px;
  margin-top: 12px;
  padding-top: 8px;
  border-top: 1px solid var(--ctl-hairline, rgba(var(--primary-color), 0.06));
}

.cmt-card__spacer {
  flex: 1;
}

.cmt-card__muted {
  font-size: 12px;
  color: #8b97a5;
}

.cmt-card__link {
  appearance: none;
  height: 26px;
  padding: 0 8px;
  border: 0;
  border-radius: 8px;
  background: transparent;
  font-size: 12px;
  font-weight: 550;
  letter-spacing: -0.01em;
  color: rgb(var(--primary-color));
  cursor: pointer;
}

.cmt-card__link:hover {
  background: rgba(var(--primary-color), 0.08);
}

.cmt-card__link.is-danger {
  color: #c45b4a;
}

.cmt-card__link.is-danger:hover {
  background: rgb(196 91 74 / 0.1);
}
</style>
