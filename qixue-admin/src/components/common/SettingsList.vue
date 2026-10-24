<template>
  <section class="settings-list">
    <div v-if="label || $slots.extra" class="settings-list__head">
      <h3 v-if="label" class="settings-list__label">
        {{ label }}
      </h3>
      <span v-else />
      <div v-if="$slots.extra" class="settings-list__extra">
        <slot name="extra" />
      </div>
    </div>
    <div class="settings-list__group">
      <slot />
    </div>
  </section>
</template>

<script setup>
defineProps({
  label: {
    type: String,
    default: '',
  },
})
</script>

<style scoped>
.settings-list {
  max-width: 760px;
  margin: 0 auto 24px;
}

.settings-list__head {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  gap: 12px;
  margin: 0 0 6px;
  padding: 0 2px;
}

.settings-list__label {
  margin: 0;
  font-family: var(--font-sans);
  font-size: 11px;
  font-weight: 600;
  letter-spacing: -0.016em;
  color: var(--ctl-text-soft, rgba(0, 0, 0, 0.5));
}

.settings-list__extra :deep(.n-button) {
  --n-height: var(--ctl-h-lg, 28px) !important;
  --n-font-size: 13px !important;
  --n-border-radius: 7px !important;
  font-weight: 500 !important;
}

.settings-list__group {
  display: flex;
  flex-direction: column;
  gap: 0;
  padding: 0;
  border: 0;
  border-radius: 0;
  background: transparent;
  overflow: visible;
}
</style>
