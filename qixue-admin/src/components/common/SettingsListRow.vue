<template>
  <button
    type="button"
    class="settings-list-row"
    :class="{ 'is-static': !chevron }"
    @click="$emit('click')"
  >
    <span v-if="$slots.icon" class="settings-list-row__ico">
      <slot name="icon" />
    </span>
    <span class="settings-list-row__meta">
      <strong>{{ title }}</strong>
      <small v-if="desc">{{ desc }}</small>
    </span>
    <span v-if="$slots.default" class="settings-list-row__value">
      <slot />
    </span>
    <i v-if="chevron" class="i-fe:chevron-right settings-list-row__chevron" />
  </button>
</template>

<script setup>
defineProps({
  title: {
    type: String,
    required: true,
  },
  desc: {
    type: String,
    default: '',
  },
  chevron: {
    type: Boolean,
    default: true,
  },
})
defineEmits(['click'])
</script>

<style scoped>
.settings-list-row {
  display: flex;
  align-items: center;
  gap: 12px;
  min-height: 40px;
  padding: 8px 2px;
  border: 0;
  border-radius: 0;
  border-bottom: 1px solid var(--ctl-hairline, rgba(60, 60, 67, 0.18));
  background: transparent;
  text-align: left;
  color: inherit;
  cursor: pointer;
  appearance: none;
  width: 100%;
  transition: background 0.18s ease;
}

.settings-list-row:first-child {
  border-top: 1px solid var(--ctl-hairline, rgba(60, 60, 67, 0.18));
}

.settings-list-row:hover {
  background: rgb(0 0 0 / 4%);
}

.settings-list-row:hover .settings-list-row__chevron {
  opacity: 0.62;
  transform: translateX(2px);
}

.settings-list-row:active {
  background: rgb(0 0 0 / 6%);
}

.settings-list-row.is-static {
  cursor: default;
}

.settings-list-row.is-static:hover,
.settings-list-row.is-static:active {
  background: transparent;
}

.settings-list-row.is-static:hover .settings-list-row__chevron {
  opacity: 0.35;
  transform: none;
}

.settings-list-row__ico {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  border-radius: 6px;
  color: rgb(var(--primary-color));
  background: rgba(var(--primary-color), 0.1);
  flex-shrink: 0;
}

.settings-list-row__meta {
  display: flex;
  flex: 1;
  min-width: 0;
  flex-direction: column;
}

.settings-list-row__meta strong {
  font-size: 13px;
  font-weight: 500;
  letter-spacing: -0.016em;
  color: var(--ctl-text, rgba(0, 0, 0, 0.85));
}

.settings-list-row__meta small {
  margin-top: 2px;
  font-size: 11px;
  letter-spacing: -0.016em;
  color: var(--ctl-text-soft, rgba(0, 0, 0, 0.5));
}

.settings-list-row__value {
  font-size: 13px;
  font-weight: 400;
  letter-spacing: -0.016em;
  color: var(--ctl-text-soft, rgba(0, 0, 0, 0.5));
  text-align: right;
  word-break: break-all;
}

.settings-list-row__chevron {
  flex-shrink: 0;
  font-size: 14px;
  opacity: 0.35;
  transition:
    opacity 0.2s ease,
    transform 0.22s cubic-bezier(0.32, 0.72, 0, 1);
}
</style>
