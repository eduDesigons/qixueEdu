<script setup>
import {
  NImage,
  NTag,
} from 'naive-ui'
import { computed } from 'vue'
import { MediaPlayer } from '@/components/index.js'
import { getAssetUrl } from '@/utils'

const props = defineProps({
  courseData: {
    type: Object,
    default: () => ({}),
  },
})

const thumbSrc = computed(() => getAssetUrl(props.courseData.thumb) || '')
const mediaInfo = computed(() => ({
  url: props.courseData.thumb,
  type: 'image',
  title: props.courseData.title || '系列封面',
}))

const certPreviewSrc = computed(() => {
  const cert = props.courseData.cert
  if (!cert)
    return ''
  return getAssetUrl(cert.preview_image || cert.template_image) || ''
})

const pending = computed(() => Number(props.courseData.pending_cert_users || 0))
</script>

<template>
  <div class="detail-head info-block">
    <img
      v-if="thumbSrc"
      class="photo-bg"
      :src="thumbSrc"
      alt=""
    >
    <div v-if="thumbSrc" class="photo-veil" />
    <div class="detail-head__media" :class="{ 'has-img': !!thumbSrc }">
      <MediaPlayer
        v-if="thumbSrc"
        class="detail-head__player"
        :media-info="mediaInfo"
      />
    </div>
    <div class="detail-head__body">
      <div class="detail-head__row">
        <span class="detail-head__lead">{{ courseData.title || '未命名系列' }}</span>
        <NTag v-if="courseData.category?.name" size="small" :bordered="false">
          {{ courseData.category.name }}
        </NTag>
        <NTag v-if="courseData.is_rec" type="primary" size="small" :bordered="false">
          推荐
        </NTag>
        <NTag v-if="courseData.is_open" size="small" :bordered="false">
          公开
        </NTag>
        <NTag v-if="courseData.is_show" type="success" size="small" :bordered="false">
          上架
        </NTag>
        <NTag v-else type="warning" size="small" :bordered="false">
          未发布
        </NTag>
        <NTag v-if="courseData.min_need_credit" type="warning" size="small" :bordered="false">
          限制{{ courseData.min_need_credit }}学分
        </NTag>
        <NTag v-if="courseData.cert?.name" type="info" size="small" :bordered="false">
          证书 {{ courseData.cert.name }}
        </NTag>
        <NTag v-if="pending > 0" type="warning" size="small" :bordered="false">
          {{ pending }}人待补发
        </NTag>
      </div>
      <p class="detail-head__sub">
        {{ courseData.desc || '暂无简介' }}
      </p>
      <div class="detail-head__metrics">
        <span><em>{{ courseData.user_path_end_count || 0 }}</em>结业</span>
        <span><em>{{ courseData.users_count || 0 }}</em>学员</span>
        <span><em>{{ courseData.courses_count || 0 }}</em>课程</span>
        <span><em>{{ courseData.applys_count || 0 }}</em>申请</span>
      </div>
      <p class="detail-head__meta">
        创建人 {{ courseData.creator?.nickName || '—' }}
        <i>·</i>
        创建 {{ courseData.created_at || '—' }}
        <i>·</i>
        更新 {{ courseData.updated_at || '—' }}
      </p>
    </div>
    <div v-if="certPreviewSrc" class="detail-head__cert">
      <NImage
        :src="certPreviewSrc"
        object-fit="cover"
        class="detail-head__cert-img"
        :preview-src="certPreviewSrc"
      />
      <span class="detail-head__cert-name">{{ courseData.cert?.name || '结业证书' }}</span>
    </div>
  </div>
</template>

<style scoped>
.detail-head__cert {
  position: relative;
  z-index: 1;
  flex-shrink: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  margin-left: auto;
  width: 148px;
}

.detail-head__cert-img {
  width: 148px;
  height: 92px;
  border-radius: 10px;
  overflow: hidden;
  box-shadow: 0 8px 24px rgba(15, 23, 42, 0.12);
  background: rgba(255, 255, 255, 0.72);
}

.detail-head__cert-img :deep(img) {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.detail-head__cert-name {
  max-width: 148px;
  font-size: 12px;
  line-height: 1.3;
  color: var(--text-secondary, #64748b);
  text-align: center;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

@media (max-width: 760px) {
  .detail-head__cert {
    margin-left: 0;
    width: 100%;
    align-items: flex-start;
  }
}
</style>
