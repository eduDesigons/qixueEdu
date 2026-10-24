<!--------------------------------
 - @Author: Ronnie Zhang
 - @LastEditor: Ronnie Zhang
 - @LastEditTime: 2023/12/04 22:51:12
 - @Email: zclzone@outlook.com
 - Copyright © 2023 Ronnie Zhang(大脸怪) | https://isme.top
 --------------------------------->

<template>
  <main class="app-page h-full min-h-0 flex-col flex-1">
    <div
      class="app-page__body"
      :class="{ 'app-page__body--fill': full }"
    >
      <slot />
    </div>
    <slot name="footer">
      <TheFooter v-if="showFooter" class="app-page__footer" />
    </slot>
    <n-back-top :bottom="20" />
  </main>
</template>

<script setup>
defineProps({
  full: {
    type: Boolean,
    default: false,
  },
  showFooter: {
    type: Boolean,
    default: false,
  },
})
</script>

<style scoped>
.app-page {
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.app-page__body {
  flex: 1 1 auto;
  min-height: 0;
  overflow-x: hidden;
  overflow-y: auto;
  margin: 12px 16px 0;
  display: flex;
  flex-direction: column;
}

.app-page__body--fill {
  overflow: hidden;
}

.app-page__footer {
  flex-shrink: 0;
  margin: 0;
}
</style>
