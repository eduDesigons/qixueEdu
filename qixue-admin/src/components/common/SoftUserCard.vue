<script setup>
import { NAvatar } from 'naive-ui'
import { computed } from 'vue'
import { getAssetUrl } from '@/utils'

defineOptions({ name: 'SoftUserCard' })

const props = defineProps({
  avatar: { type: String, default: '' },
  name: { type: String, default: '未命名' },
  subtitle: { type: String, default: '' },
  status: { type: String, default: '' },
  statusTone: {
    type: String,
    default: 'idle', // done | doing | info | warn | error | idle | ok
  },
  time: { type: String, default: '' },
  metrics: {
    type: Array,
    default: () => [],
  },
  desc: { type: String, default: '' },
  clickable: { type: Boolean, default: true },
  deletable: { type: Boolean, default: false },
})

const emit = defineEmits(['click', 'delete'])

const avatarUrl = computed(() => getAssetUrl(props.avatar) || undefined)

function onCardClick() {
  if (props.clickable)
    emit('click')
}

function onDelete(e) {
  e.stopPropagation()
  emit('delete')
}
</script>

<template>
  <article
    class="soft-user-card entity-card"
    :class="{ 'is-clickable': clickable }"
    @click="onCardClick"
  >
    <div class="soft-user-card__top">
      <div class="soft-user-card__who">
        <NAvatar
          v-if="avatarUrl"
          round
          :size="36"
          :src="avatarUrl"
          class="soft-user-card__avatar"
        />
        <div class="soft-user-card__name-wrap">
          <strong class="soft-user-card__name" :title="name">
            {{ name || '未命名' }}
          </strong>
          <span v-if="subtitle" class="soft-user-card__sub">
            {{ subtitle }}
          </span>
        </div>
      </div>
      <span
        v-if="status"
        class="soft-user-card__status"
        :class="`is-${statusTone}`"
      >
        {{ status }}
      </span>
      <slot v-else name="status" />
    </div>

    <p v-if="desc" class="soft-user-card__desc" :title="desc">
      {{ desc }}
    </p>
    <slot name="body" />

    <div v-if="metrics.length" class="soft-user-card__stats">
      <div
        v-for="(m, i) in metrics"
        :key="`${m.label}-${i}`"
        class="soft-user-card__stat"
      >
        <strong>{{ m.value }}</strong>
        <span>{{ m.label }}</span>
      </div>
    </div>
    <slot name="metrics" />

    <div class="soft-user-card__foot">
      <time class="soft-user-card__time">
        {{ time || '—' }}
      </time>
      <div class="soft-user-card__actions" @click.stop>
        <slot name="actions" />
        <button
          v-if="deletable"
          type="button"
          class="soft-user-card__del"
          title="删除"
          @click="onDelete"
        >
          <i class="i-material-symbols:delete-outline text-14" />
          删除
        </button>
      </div>
    </div>
  </article>
</template>

<style scoped>
.soft-user-card {
  position: relative;
  display: flex;
  flex-direction: column;
  min-height: 148px;
  padding: 14px;
}

.soft-user-card.is-clickable {
  cursor: pointer;
}

.soft-user-card__top {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 10px;
}

.soft-user-card__who {
  display: flex;
  align-items: center;
  gap: 10px;
  min-width: 0;
  flex: 1;
}

.soft-user-card__avatar {
  flex-shrink: 0;
  box-shadow: 0 0 0 2px rgb(255 255 255 / 80%);
}

.soft-user-card__name-wrap {
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.soft-user-card__name {
  display: block;
  max-width: 100%;
  font-size: 14px;
  font-weight: 650;
  letter-spacing: -0.018em;
  line-height: 1.3;
  color: #1a2332;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.soft-user-card__sub {
  font-size: 11px;
  font-weight: 520;
  color: #8b97a5;
}

.soft-user-card__status {
  flex-shrink: 0;
  display: inline-flex;
  align-items: center;
  height: 24px;
  padding: 0 10px;
  border-radius: 999px;
  font-size: 11px;
  font-weight: 650;
  letter-spacing: -0.01em;
  white-space: nowrap;
  color: #5b6b7c;
  background: rgba(var(--primary-color), 0.08);
  border: 1px solid rgba(var(--primary-color), 0.1);
}

.soft-user-card__status.is-done,
.soft-user-card__status.is-ok {
  color: #2d7a45;
  background: rgb(52 199 89 / 16%);
  border-color: rgb(52 199 89 / 18%);
}

.soft-user-card__status.is-doing,
.soft-user-card__status.is-error {
  color: #a14b3a;
  background: rgb(255 59 48 / 10%);
  border-color: rgb(255 59 48 / 12%);
}

.soft-user-card__status.is-warn {
  color: #8a5a00;
  background: rgb(255 159 10 / 16%);
  border-color: rgb(255 159 10 / 18%);
}

.soft-user-card__status.is-info {
  color: #2f6f8f;
  background: rgb(90 170 210 / 14%);
  border-color: rgb(90 170 210 / 18%);
}

.soft-user-card__status.is-idle {
  color: #8b97a5;
  background: rgb(138 150 151 / 10%);
  border-color: rgb(138 150 151 / 14%);
}

.soft-user-card__desc {
  margin: 12px 0 0;
  padding: 0;
  font-size: 12px;
  font-weight: 500;
  line-height: 1.45;
  color: #5b6b7c;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.soft-user-card__stats {
  display: flex;
  align-items: stretch;
  margin-top: 14px;
  padding: 10px 0;
  border-top: 1px solid var(--ctl-hairline, rgba(var(--primary-color), 0.06));
  border-bottom: 1px solid var(--ctl-hairline, rgba(var(--primary-color), 0.06));
}

.soft-user-card__stat {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
  text-align: center;
}

.soft-user-card__stat + .soft-user-card__stat {
  border-left: 1px solid var(--ctl-hairline, rgba(var(--primary-color), 0.06));
}

.soft-user-card__stat strong {
  font-size: 16px;
  font-weight: 700;
  letter-spacing: -0.03em;
  color: #1a2332;
  font-variant-numeric: tabular-nums;
  line-height: 1.2;
}

.soft-user-card__stat span {
  font-size: 11px;
  font-weight: 500;
  color: #8b97a5;
}

.soft-user-card__foot {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  margin-top: auto;
  padding-top: 10px;
  min-height: 34px;
}

.soft-user-card__time {
  font-size: 11px;
  font-weight: 520;
  letter-spacing: -0.01em;
  color: #8b97a5;
  font-variant-numeric: tabular-nums;
}

.soft-user-card__actions {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  opacity: 0;
  transform: translateY(3px) scale(0.96);
  pointer-events: none;
  transition:
    opacity 0.2s cubic-bezier(0.32, 0.72, 0, 1),
    transform 0.2s cubic-bezier(0.32, 0.72, 0, 1);
}

.soft-user-card:hover .soft-user-card__actions,
.soft-user-card:focus-within .soft-user-card__actions {
  opacity: 1;
  transform: none;
  pointer-events: auto;
}

.soft-user-card__del {
  appearance: none;
  display: inline-flex;
  align-items: center;
  gap: 4px;
  height: 26px;
  padding: 0 10px;
  border: 0;
  border-radius: 8px;
  font-size: 12px;
  font-weight: 560;
  color: #c4342d;
  background: rgb(255 59 48 / 8%);
  cursor: pointer;
  transition: background 0.16s ease;
}

.soft-user-card__del:hover {
  background: rgb(255 59 48 / 14%);
}

.soft-user-card__actions :deep(.n-button) {
  --n-height: 26px !important;
  --n-font-size: 12px !important;
  --n-padding: 0 10px !important;
  --n-border-radius: 8px !important;
}

@media (prefers-reduced-motion: reduce) {
  .soft-user-card__actions {
    transition: none;
  }
}
</style>
