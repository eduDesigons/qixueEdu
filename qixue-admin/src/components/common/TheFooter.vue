<script setup>
import { computed } from 'vue'
import { useSysStore } from '@/store'

const DEFAULT_FOOTER = 'The night is nearly over; The day is almost here. So let us put aside the deeds of darkness and put on the armor of light.'

const text = computed(() => {
  const raw = useSysStore().sysConfig?.footer
  if (raw === undefined || raw === null)
    return DEFAULT_FOOTER
  return String(raw).trim()
})
</script>

<template>
  <footer v-if="text" class="the-footer f-c-c">
    <p>{{ text }}</p>
  </footer>
</template>

<style scoped>
.the-footer {
  padding: 10px 20px 16px;
}

.the-footer p {
  margin: 0;
  max-width: 720px;
  text-align: center;
  font-size: 11.5px;
  font-weight: 450;
  line-height: 1.55;
  letter-spacing: -0.01em;
  color: var(--ctl-text-soft, rgba(var(--primary-color), 0.42));
  opacity: 0.62;
  white-space: pre-wrap;
}
</style>
