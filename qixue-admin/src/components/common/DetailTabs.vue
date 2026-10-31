<script setup>
import { NTabs } from 'naive-ui'

defineOptions({ name: 'DetailTabs' })

const current = defineModel('value', { type: String })

const filterHostRef = ref(null)
provide('pageFilterHost', filterHostRef)
provide('detailActiveTab', current)
</script>

<template>
  <NTabs
    v-model:value="current"
    type="line"
    class="detail-tabs"
  >
    <template #suffix>
      <div class="detail-tabs__actions">
        <div
          ref="filterHostRef"
          class="page-filter-host page-filter-host--tabs"
        />
        <slot name="suffix" />
      </div>
    </template>
    <slot />
  </NTabs>
</template>
