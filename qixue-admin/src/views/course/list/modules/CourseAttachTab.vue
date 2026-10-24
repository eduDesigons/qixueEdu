<script setup>
import { NEmpty, NSpin } from 'naive-ui'
import api from '@/api/index.js'
import { MediaPlayer } from '@/components/index.js'
import { useCrud } from '@/composables/index.js'
import { getAssetUrl } from '@/utils'

const props = defineProps({
  courseId: {
    type: [String, Number],
    required: true,
  },
})

const list = ref([])
const loading = ref(false)

const { handleDelete } = useCrud({
  name: '课件',
  doDelete: api.course.courseAttachDelete,
  refresh: () => getList(),
})

function formatSize(size) {
  const n = Number(size) || 0
  if (n >= 1_000_000)
    return `${(n / 1_000_000).toFixed(n >= 10_000_000 ? 0 : 1)} MB`
  if (n >= 1000)
    return `${(n / 1000).toFixed(n >= 10_000 ? 0 : 1)} KB`
  return `${n} B`
}

function mediaOf(item) {
  return {
    url: getAssetUrl(item.path),
    name: item.name,
    title: item.name,
    type: 'document',
    size: item.size,
  }
}

async function getList() {
  if (!props.courseId)
    return
  loading.value = true
  try {
    const { data } = await api.course.courseAttach({ courseId: props.courseId })
    list.value = Array.isArray(data) ? data : (data?.list || [])
  }
  catch (error) {
    $message.error(error?.message || '课件加载失败')
    list.value = []
  }
  finally {
    loading.value = false
  }
}

onMounted(getList)

defineExpose({
  getList,
})
</script>

<template>
  <div class="attach-panel h-500">
    <NSpin :show="loading">
      <div v-if="list.length" class="attach-grid">
        <article
          v-for="item in list"
          :key="item.id"
          class="attach-card"
        >
          <div class="attach-card__preview">
            <MediaPlayer
              :media-info="mediaOf(item)"
              :controls="false"
              compact
              deletable
              @delete="handleDelete(item.id)"
            />
          </div>
          <div class="attach-card__body">
            <h3 class="attach-card__title" :title="item.name">
              {{ item.name }}
            </h3>
            <p class="attach-card__meta">
              <span>{{ formatSize(item.size) }}</span>
            </p>
            <p v-if="item.created_at" class="attach-card__time">
              {{ item.created_at }}
            </p>
          </div>
        </article>
      </div>
      <div v-else-if="!loading" class="attach-empty">
        <NEmpty description="暂无课件，点击上方上传" />
      </div>
    </NSpin>
  </div>
</template>

<style scoped>
.attach-panel {
  min-height: 240px;
}

.attach-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(168px, 1fr));
  gap: 14px;
  padding-bottom: 8px;
}

.attach-card {
  display: flex;
  flex-direction: column;
  overflow: hidden;
  border-radius: 16px;
  border: 1px solid var(--ctl-hairline, rgba(var(--primary-color), 0.08));
  background:
    radial-gradient(120% 80% at 50% 0%, rgba(var(--primary-color), 0.06), transparent 62%),
    rgb(255 255 255 / 72%);
  box-shadow: inset 0 1px 0 rgb(255 255 255 / 55%);
}

.attach-card__preview {
  position: relative;
  height: 132px;
  overflow: hidden;
  background:
    radial-gradient(90% 70% at 50% 30%, rgba(var(--primary-color), 0.1), transparent 68%),
    var(--ctl-field, rgba(var(--primary-color), 0.04));
}

.attach-card__preview :deep(.mp) {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
  min-height: 0;
  border-radius: 0;
}

.attach-card__preview :deep(.mp-thumb) {
  width: 100%;
  height: 100%;
  min-height: 0;
  border-radius: 0;
}

.attach-card__preview :deep(.mp-image) {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.attach-card__preview :deep(.mp-doc-icon) {
  width: 64px;
  height: 64px;
}

.attach-card__preview :deep(.mp-chip--type) {
  font-weight: 650;
  letter-spacing: 0.02em;
}

.attach-card__body {
  display: flex;
  flex-direction: column;
  gap: 4px;
  padding: 12px 12px 14px;
}

.attach-card__title {
  margin: 0;
  font-size: 13px;
  font-weight: 650;
  letter-spacing: -0.016em;
  color: #1a2332;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  line-height: 1.35;
}

.attach-card__meta {
  margin: 0;
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: #5b6b7c;
}

.attach-card__meta i {
  font-style: normal;
  opacity: 0.45;
}

.attach-card__time {
  margin: 0;
  font-size: 11px;
  color: #8a989a;
}

.attach-empty {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 220px;
}
</style>
