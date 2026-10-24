<template>
  <div v-if="embedded" class="settings-editor">
    <slot />
    <footer class="settings-editor__bar">
      <n-button @click="$emit('close')">
        取消
      </n-button>
      <n-button type="primary" :loading="loading" @click="$emit('save')">
        保存
      </n-button>
    </footer>
  </div>
  <NModal
    v-else
    v-model:show="open"
    preset="dialog"
    :title="title"
    :style="width ? { width } : undefined"
    positive-text="保存"
    negative-text="取消"
    :loading="loading"
    @positive-click="$emit('save')"
  >
    <slot />
  </NModal>
</template>

<script setup>
import { NModal } from 'naive-ui'

defineProps({
  embedded: {
    type: Boolean,
    default: false,
  },
  title: {
    type: String,
    default: '',
  },
  width: {
    type: [String, Number],
    default: undefined,
  },
  loading: {
    type: Boolean,
    default: false,
  },
})
defineEmits(['save', 'close'])
const open = defineModel('show', { type: Boolean, default: false })
</script>

<style scoped>
.settings-editor {
  display: flex;
  flex-direction: column;
  gap: 20px;
  min-height: 0;
}

.settings-editor__bar {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  position: sticky;
  bottom: 0;
  padding: 12px 0 4px;
  background: linear-gradient(180deg, transparent, rgb(244 248 248 / 88%) 28%);
}

.settings-editor__bar :deep(.n-button) {
  --n-height: var(--ctl-h-lg, 28px) !important;
  --n-border-radius: 7px !important;
  min-width: 72px;
  transition: transform 0.2s cubic-bezier(0.32, 0.72, 0, 1) !important;
}

.settings-editor__bar :deep(.n-button:not(.n-button--disabled):active) {
  transform: scale(0.97);
}
</style>
