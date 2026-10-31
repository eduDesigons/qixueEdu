<!--------------------------------
 - @Author: Ronnie Zhang
 - @LastEditor: Ronnie Zhang
 - @LastEditTime: 2023/12/13 20:54:55
 - @Email: zclzone@outlook.com
 - Copyright © 2023 Ronnie Zhang(大脸怪) | https://isme.top
 --------------------------------->

<template>
  <div class="layout-shell wh-full flex">
    <LayoutAurora />
    <aside
      class="layout-aside flex-col flex-shrink-0 transition-width-300"
      :class="appStore.collapsed ? 'w-64' : 'w-200'"
    >
      <SideBar />
    </aside>

    <article class="layout-main w-0 flex-col flex-1">
      <AppHeader class="layout-header header-bar-h flex-shrink-0" />
      <div class="layout-main__body">
        <slot />
      </div>
    </article>
  </div>
</template>

<script setup>
import { useAppStore } from '@/store'
import { LayoutAurora } from '@/layouts/components'
import AppHeader from './header/index.vue'
import SideBar from './sidebar/index.vue'

const appStore = useAppStore()
</script>

<style>
.collapsed {
  width: 64px;
}

.header-bar-h {
  height: var(--bar-h, 52px);
  max-height: var(--bar-h, 52px);
  flex-shrink: 0;
}

.layout-main {
  display: flex;
  flex-direction: column;
  min-height: 0;
  overflow: hidden;
}

.layout-main__body {
  position: relative;
  flex: 1;
  min-height: 0;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}


.layout-main__body .page-view > main,
.layout-main__body > main {
  height: 100%;
  max-height: 100%;
  min-height: 0;
  display: flex;
  flex-direction: column;
}

.layout-main__body .page-view > main:not(.cus-scroll),
.layout-main__body > main:not(.cus-scroll) {
  overflow: hidden;
}

.layout-main__body .page-view > main.cus-scroll,
.layout-main__body > main.cus-scroll {
  overflow-x: hidden;
  overflow-y: auto;
}

.layout-aside {
  transition:
    width 0.36s cubic-bezier(0.32, 0.72, 0, 1),
    min-width 0.36s cubic-bezier(0.32, 0.72, 0, 1) !important;
}
</style>
