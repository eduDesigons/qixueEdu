<!--------------------------------
 - @Author: Ronnie Zhang
 - @LastEditor: Ronnie Zhang
 - @LastEditTime: 2023/12/16 18:50:18
 - @Email: zclzone@outlook.com
 - Copyright © 2023 Ronnie Zhang(大脸怪) | https://isme.top
 --------------------------------->

<template>
  <div
    id="menu-collapse"
    class="menu-collapse f-c-c cursor-pointer text-15"
    @click="appStore.switchCollapsed"
  >
    <i :class="appStore.collapsed ? 'i-line-md-menu-unfold-left' : 'i-line-md-menu-fold-left'" />
  </div>
</template>

<script setup>
import { useAppStore } from '@/store'

const appStore = useAppStore()
</script>

<style scoped>
.menu-collapse {
  width: 28px;
  height: 28px;
  border-radius: 9px;
  color: var(--ctl-text-soft, #5b6b7c);
  border: 1px solid transparent;
  transition:
    transform 0.24s cubic-bezier(0.32, 0.72, 0, 1),
    background 0.22s ease,
    box-shadow 0.28s cubic-bezier(0.22, 1, 0.36, 1),
    border-color 0.22s ease,
    color 0.22s ease;
}

.menu-collapse:hover {
  color: var(--ctl-text, #1a2332);
  background: var(--ctl-field, rgba(var(--primary-color), 0.06));
  border-color: var(--ctl-field-border, rgba(var(--primary-color), 0.12));
  box-shadow: inset 0 1px 0 rgb(255 255 255 / 55%);
}

.menu-collapse:active {
  transform: scale(0.94);
}
</style>
