<!--------------------------------
 - @Author: Ronnie Zhang
 - @LastEditor: Ronnie Zhang
 - @LastEditTime: 2023/12/16 18:52:48
 - @Email: zclzone@outlook.com
 - Copyright © 2023 Ronnie Zhang(大脸怪) | https://isme.top
 --------------------------------->

<template>
  <AppCard class="layout-header header-bar">
    <div class="header-bar__row">
      <div class="header-bar__leading">
        <MenuCollapse />
      </div>

      <div class="header-bar__tabs">
        <AppTab />
      </div>

      <div class="header-bar__actions">
        <LayoutViewToggle />
        <span class="header-bar__divider" aria-hidden="true" />
        <UserAvatar />
      </div>
    </div>
  </AppCard>
</template>

<script setup>
import { AppTab, LayoutViewToggle, MenuCollapse, UserAvatar } from '@/layouts/components'
</script>

<style scoped>
.header-bar {
  display: flex;
  flex-direction: column;
  justify-content: flex-end;
  height: var(--bar-h, 52px);
  max-height: var(--bar-h, 52px);
  min-height: var(--bar-h, 52px);
  flex-shrink: 0;
  padding: 0 20px !important;
  overflow: hidden;
}

.header-bar__row {
  display: flex;
  align-items: flex-end;
  gap: 8px;
  width: 100%;
  min-width: 0;
  height: 28px;
}

.header-bar__leading,
.header-bar__actions {
  display: flex;
  flex-shrink: 0;
  align-items: center;
  height: 28px;
  margin-bottom: 0;
}

.header-bar__leading {
  padding-left: 2px;
}

.header-bar__tabs {
  flex: 1;
  min-width: 0;
  height: 28px;
  display: flex;
  align-items: flex-end;
}

.header-bar__tabs :deep(#top-tab) {
  width: 100%;
  height: 28px;
  padding: 0;
}

.header-bar__actions {
  gap: 1px;
  padding-right: 2px;
  padding-bottom: 1px;
}

.header-bar__divider {
  width: 1px;
  height: 12px;
  margin: 0 6px;
  border-radius: 999px;
  background: var(--ctl-hairline, rgba(var(--primary-color), 0.14));
}

.header-bar__leading :deep(.menu-collapse),
.header-bar__actions :deep(.header-icon-btn) {
  width: 22px;
  height: 22px;
  font-size: 14px;
  border-radius: 6px;
}
</style>
