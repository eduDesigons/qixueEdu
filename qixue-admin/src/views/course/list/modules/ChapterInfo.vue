<script setup>
import {
  NTag,
} from 'naive-ui'
import { MediaPlayer } from '@/components/index.js'

defineProps({
  info: {
    type: Object,
    default: () => ({}),
  },
})
</script>

<template>
  <div class="detail-head info-block">
    <div class="detail-head__media">
      <MediaPlayer :media-info="info.media" />
    </div>
    <div class="detail-head__body">
      <div class="detail-head__row">
        <span class="detail-head__lead">{{ info.title || '未命名课节' }}</span>
        <NTag
          :type="info.limit_day ? 'warning' : 'default'"
          size="small"
          :bordered="false"
        >
          {{ info.limit_day ? `限制 ${info.limit_day} 天` : '不限时间' }}
        </NTag>
      </div>
      <div class="detail-head__metrics">
        <span><em>{{ info.user_chapter_end_count || 0 }}</em>学完</span>
        <span><em>{{ info.users_count || 0 }}</em>学员</span>
        <span><em>{{ info.comments_count || 0 }}</em>评论</span>
      </div>
      <p class="detail-head__meta">
        创建人 {{ info.creator?.nickName || '—' }}
        <i>·</i>
        创建 {{ info.created_at || '—' }}
        <i>·</i>
        更新 {{ info.updated_at || '—' }}
      </p>
    </div>
  </div>
</template>

<style scoped>
.detail-head__media {
  cursor: zoom-in;
}

.detail-head__media :deep(.mp),
.detail-head__media :deep(.mp-thumb) {
  width: 100%;
  height: 100%;
  border-radius: inherit;
}

.detail-head__media :deep(.mp-play) {
  width: 36px;
  height: 36px;
}

.detail-head__media :deep(.mp-play i) {
  font-size: 15px;
}

@media (max-width: 720px) {
  .detail-head {
    flex-direction: column;
  }
}
</style>
