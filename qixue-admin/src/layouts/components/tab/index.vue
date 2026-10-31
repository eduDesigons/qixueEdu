<!--------------------------------
 - @Author: Ronnie Zhang
 - @LastEditor: Ronnie Zhang
 - @LastEditTime: 2023/12/16 18:50:54
 - @Email: zclzone@outlook.com
 - Copyright © 2023 Ronnie Zhang(大脸怪) | https://isme.top
 --------------------------------->

<template>
  <div id="top-tab" class="safari-tabs">
    <n-tabs
      :value="tabStore.activeTab"
      :closable="tabStore.tabs.length > 1"
      type="card"
      @close="(path) => tabStore.removeTab(path)"
    >
      <n-tab
        v-for="item in tabStore.tabs"
        :key="item.path"
        :name="item.path"
        @click="handleItemClick(item.path)"
        @contextmenu.prevent="handleContextMenu($event, item)"
      >
        {{ item.title }}
      </n-tab>
    </n-tabs>

    <ContextMenu
      v-if="contextMenuOption.show"
      v-model:show="contextMenuOption.show"
      :current-path="contextMenuOption.currentPath"
      :x="contextMenuOption.x"
      :y="contextMenuOption.y"
    />
  </div>
</template>

<script setup>
import { useTabStore } from '@/store'
import ContextMenu from './ContextMenu.vue'

const router = useRouter()
const tabStore = useTabStore()

const contextMenuOption = reactive({
  show: false,
  x: 0,
  y: 0,
  currentPath: '',
})

function handleItemClick(path) {
  tabStore.setActiveTab(path)
  router.push(path)
}

function showContextMenu() {
  contextMenuOption.show = true
}
function hideContextMenu() {
  contextMenuOption.show = false
}
function setContextMenu(x, y, currentPath) {
  Object.assign(contextMenuOption, { x, y, currentPath })
}

async function handleContextMenu(e, tagItem) {
  const { clientX, clientY } = e
  hideContextMenu()
  setContextMenu(clientX, clientY, tagItem.path)
  await nextTick()
  showContextMenu()
}
</script>

<style scoped>
#top-tab {
  --apple-ease: cubic-bezier(0.32, 0.72, 0, 1);
  --safari-tab-h: 28px;
  display: flex;
  align-items: flex-end;
  min-width: 0;
  width: 100%;
  height: var(--safari-tab-h);
  max-height: var(--safari-tab-h);
  padding: 0;
}

#top-tab :deep(.n-tabs) {
  width: 100%;
  height: var(--safari-tab-h);
  max-height: var(--safari-tab-h);
  background: transparent !important;
}

#top-tab :deep(.n-tabs.n-tabs--card-type .n-tabs-nav),
#top-tab :deep(.n-tabs-nav) {
  height: var(--safari-tab-h) !important;
  min-height: var(--safari-tab-h) !important;
  padding: 0 !important;
  margin: 0 !important;
  border: 0 !important;
  border-radius: 0 !important;
  background: transparent !important;
  box-shadow: none !important;
  backdrop-filter: none !important;
  -webkit-backdrop-filter: none !important;
  line-height: 1;
}

#top-tab :deep(.n-tabs-nav-scroll-wrapper),
#top-tab :deep(.n-tabs-nav-scroll-content),
#top-tab :deep(.n-tabs-wrapper),
#top-tab :deep(.n-tabs-tab-wrapper) {
  height: 100% !important;
  background: transparent !important;
  align-items: flex-end !important;
}

#top-tab :deep(.n-tabs-tab-wrapper) {
  position: relative;
  margin: 0 !important;
  padding: 0 !important;
}

#top-tab :deep(.n-tabs-tab-wrapper:not(:first-child) .n-tabs-tab:not(.n-tabs-tab--active)::before) {
  content: '';
  position: absolute;
  left: 0;
  top: 50%;
  width: 1px;
  height: 11px;
  transform: translateY(-50%);
  background: rgba(var(--primary-color), 0.14);
  pointer-events: none;
}

#top-tab :deep(.n-tabs-tab-wrapper:has(.n-tabs-tab--active) .n-tabs-tab::before),
#top-tab :deep(.n-tabs-tab-wrapper:has(.n-tabs-tab--active) + .n-tabs-tab-wrapper .n-tabs-tab::before) {
  opacity: 0;
}

#top-tab :deep(.n-tabs.n-tabs--card-type .n-tabs-tab),
#top-tab :deep(.n-tabs-tab) {
  position: relative;
  z-index: 1;
  display: inline-flex !important;
  align-items: center !important;
  gap: 4px;
  height: var(--safari-tab-h) !important;
  min-height: var(--safari-tab-h) !important;
  max-width: 160px;
  margin: 0 !important;
  padding: 0 10px 0 12px !important;
  border: 0 !important;
  border-radius: 10px 10px 0 0 !important;
  background: transparent !important;
  box-shadow: none !important;
  color: var(--ctl-text-soft, #6a7b8d) !important;
  font-size: 13px !important;
  font-weight: 540 !important;
  letter-spacing: -0.016em;
  line-height: 1.2 !important;
  transition:
    background 0.2s var(--apple-ease),
    color 0.18s ease,
    box-shadow 0.2s var(--apple-ease) !important;
}

#top-tab :deep(.n-tabs-tab:hover:not(.n-tabs-tab--active)) {
  color: var(--ctl-text, #1a2332) !important;
  background: rgb(255 255 255 / 32%) !important;
  border: 0 !important;
  box-shadow: none !important;
  transform: none !important;
}

#top-tab :deep(.n-tabs-tab:active) {
  transform: none !important;
}

#top-tab :deep(.n-tabs.n-tabs--card-type .n-tabs-tab.n-tabs-tab--active),
#top-tab :deep(.n-tabs-tab--active) {
  z-index: 2;
  color: var(--ctl-text, #1a2332) !important;
  font-weight: 600 !important;
  background: rgb(255 255 255 / 72%) !important;
  box-shadow:
    inset 0 1px 0 rgb(255 255 255 / 90%),
    1px 0 0 rgba(var(--primary-color), 0.06),
    -1px 0 0 rgba(var(--primary-color), 0.06) !important;
  border: 0 !important;
}

#top-tab :deep(.n-tabs-tab .n-tabs-tab__label) {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

#top-tab :deep(.n-tabs-tab .n-base-close) {
  margin-left: 1px;
  width: 16px;
  height: 16px;
  border-radius: 999px;
  opacity: 0.38;
  font-size: 11px;
  color: #5b6b7c;
  transition:
    background 0.18s ease,
    color 0.18s ease,
    opacity 0.18s ease,
    transform 0.18s var(--apple-ease);
}

#top-tab :deep(.n-tabs-tab:hover .n-base-close),
#top-tab :deep(.n-tabs-tab--active .n-base-close) {
  opacity: 0.72;
}

#top-tab :deep(.n-tabs-tab .n-base-close:hover) {
  opacity: 1;
  background: rgba(var(--primary-color), 0.12);
  color: var(--ctl-text, #1a2332);
}

#top-tab :deep(.n-tabs-tab .n-base-close:active) {
  transform: scale(0.9);
}

#top-tab :deep(.n-tabs-pad),
#top-tab :deep(.n-tabs-tab-pad),
#top-tab :deep(.n-tabs-scroll-padding) {
  border: none !important;
  background: transparent !important;
}
</style>
