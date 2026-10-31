<script setup>
import {
  NAvatar,
  NTag,
} from 'naive-ui'
import { computed } from 'vue'
import { getAssetUrl } from '@/utils'

const props = defineProps({
  userData: {
    type: Object,
    default: () => ({}),
  },
})

const thumbSrc = computed(() => getAssetUrl(props.userData.avatar) || '')
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
    <NAvatar
      round
      :size="64"
      :src="thumbSrc || undefined"
      class="detail-head__avatar"
    />
    <div class="detail-head__body">
      <div class="detail-head__row">
        <span class="detail-head__lead">{{ userData.nick_name || '未命名学员' }}</span>
        <NTag v-if="userData.area?.name" size="small" :bordered="false">
          {{ userData.area.name }}
        </NTag>
        <NTag v-if="Number(userData.is_lock) === 1" type="error" size="small" :bordered="false">
          已拉黑
        </NTag>
      </div>
      <p class="detail-head__sub">
        {{ userData.mobile || '暂无手机号' }}
      </p>
      <div class="detail-head__metrics">
        <span><em>{{ userData.certs_count || 0 }}</em>证书</span>
        <span><em>{{ userData.courses_count || 0 }}</em>课程</span>
        <span><em>{{ userData.credit || 0 }}</em>学分</span>
        <span><em>{{ userData.coins || 0 }}</em>积分</span>
      </div>
      <p class="detail-head__meta">
        注册 {{ userData.created_at || '—' }}
        <i>·</i>
        最近活跃 {{ userData.updated_at || '—' }}
      </p>
    </div>
  </div>
</template>
