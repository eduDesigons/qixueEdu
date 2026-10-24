<template>
  <span class="tip flex items-center">
    <span class="tip__label">{{ label }}</span>
    <n-popover v-if="content" trigger="hover" placement="top">
      <template #trigger>
        <button type="button" class="tip__btn" tabindex="-1" aria-label="说明">
          <i class="i-material-symbols:info-outline tip__icon" />
        </button>
      </template>
      <span class="tip__content">{{ content }}</span>
    </n-popover>
  </span>
</template>

<script setup>
defineProps({
  label: {
    type: String,
    required: true,
  },
  content: {
    type: String,
    default: '',
  },
})
</script>

<style scoped>
.tip {
  gap: 4px;
}

.tip__btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 18px;
  height: 18px;
  margin: 0;
  padding: 0;
  border: 0;
  border-radius: 999px;
  background: transparent;
  color: var(--ctl-text-soft, #8b97a5);
  cursor: help;
  transition:
    color 0.2s ease,
    background 0.2s ease,
    transform 0.2s var(--ctl-ease, cubic-bezier(0.32, 0.72, 0, 1));
}

.tip__btn:hover {
  color: rgb(var(--primary-color));
  background: var(--ctl-field, rgba(var(--primary-color), 0.08));
}

.tip__btn:active {
  transform: scale(0.92);
}

.tip__icon {
  font-size: 14px;
}

.tip__label {
  font-weight: 550;
  letter-spacing: -0.01em;
}

.tip__content {
  display: inline-block;
  max-width: 260px;
  font-size: 13px;
  line-height: 1.5;
  letter-spacing: -0.01em;
}
</style>
