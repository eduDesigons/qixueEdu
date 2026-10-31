<!--------------------------------
 - @Author: Ronnie Zhang
 - @LastEditor: Ronnie Zhang
 - @LastEditTime: 2023/12/04 22:51:48
 - @Email: zclzone@outlook.com
 - Copyright © 2023 Ronnie Zhang(大脸怪) | https://isme.top
 --------------------------------->

<template>
  <div class="query-item flex items-center">
    <label
      v-if="label || label === 0"
      class="query-item__label flex-shrink-0"
    >
      {{ label }}
    </label>
    <div class="query-item__content">
      <slot />
    </div>
  </div>
</template>

<script setup>
defineProps({
  label: {
    type: String,
    default: '',
  },
  /** @deprecated 全局统一宽度，保留 prop 避免调用方报错 */
  labelWidth: {
    type: Number,
    default: 56,
  },
  /** @deprecated 内容区统一宽度 */
  contentWidth: {
    type: Number,
    default: 180,
  },
})
</script>

<style scoped>
.query-item {
  box-sizing: border-box;
  gap: 8px;
  min-height: 34px;
  margin: 2px 4px 2px 0;
  /* 固定占位，避免在 n-space 里被压成 0 宽导致「类型/状态」消失 */
  flex: 0 0 auto;
  width: 248px;
  max-width: 100%;
}

.query-item__label {
  flex: 0 0 56px;
  width: 56px;
  box-sizing: border-box;
  font-size: 12.5px;
  font-weight: 550;
  letter-spacing: -0.012em;
  line-height: 1.2;
  color: var(--ctl-text-soft, #5b6b7c);
  text-align: right;
  white-space: nowrap;
}

.query-item__content {
  flex: 1 1 auto;
  min-width: 160px;
  width: 180px;
}

.query-item__content :deep(.n-input),
.query-item__content :deep(.n-base-selection),
.query-item__content :deep(.n-date-picker),
.query-item__content :deep(.n-cascader),
.query-item__content :deep(.n-select),
.query-item__content :deep(.n-input-number) {
  width: 100%;
}

.query-item__content :deep(.n-input),
.query-item__content :deep(.n-base-selection) {
  --n-height: var(--ctl-h-lg, 28px) !important;
  --n-border-radius: 7px !important;
  --n-font-size: 13px !important;
  --n-color: var(--ctl-field, rgba(var(--primary-color), 0.055)) !important;
  --n-color-focus: var(--ctl-field-focus, rgb(255 255 255 / 86%)) !important;
  --n-border: 1px solid var(--ctl-field-border, rgba(var(--primary-color), 0.12)) !important;
  --n-border-hover: 1px solid var(--ctl-field-border-hover, rgba(var(--primary-color), 0.18)) !important;
  backdrop-filter: none !important;
  -webkit-backdrop-filter: none !important;
  box-shadow: var(--ctl-field-shadow, inset 0 1px 2px rgba(var(--primary-color), 0.07)) !important;
}

.query-item__content :deep(.n-input:focus-within),
.query-item__content :deep(.n-base-selection.n-base-selection--active),
.query-item__content :deep(.n-base-selection.n-base-selection--focus) {
  box-shadow: var(--ctl-field-shadow-focus, 0 0 0 3.5px rgba(var(--primary-color), 0.22)) !important;
}
</style>
