<!--------------------------------
 - @Author: Ronnie Zhang
 - @LastEditor: Ronnie Zhang
 - @LastEditTime: 2023/12/16 18:50:35
 - @Email: zclzone@outlook.com
 - Copyright © 2023 Ronnie Zhang(大脸怪) | https://isme.top
 --------------------------------->

<template>
  <n-menu
    ref="menu"
    class="side-menu"
    accordion
    :indent="12"
    :collapsed-icon-size="20"
    :collapsed-width="64"
    :collapsed="appStore.collapsed"
    :options="permissionStore.menus"
    :value="activeKey"
    @update:value="handleMenuSelect"
  />
</template>

<script setup>
import { useAppStore, usePermissionStore } from '@/store'
import { isExternal } from '@/utils'

const router = useRouter()
const route = useRoute()
const appStore = useAppStore()
const permissionStore = usePermissionStore()

const activeKey = computed(() => route.meta?.parentKey || route.name)

const menu = ref(null)
watch(route, async () => {
  await nextTick()
  menu.value?.showOption()
})

function handleMenuSelect(key, item) {
  if (isExternal(item.originPath)) {
    $dialog.confirm({
      type: 'info',
      title: `请选择打开方式`,
      positiveText: '外链打开',
      negativeText: '在本站内嵌打开',
      confirm() {
        window.open(item.originPath)
      },
      cancel: () => {
        router.push(item.path)
      },
    })
  }
  else {
    if (!item.path)
      return
    router.push(item.path)
  }
}
</script>

<style>
.side-menu:not(.n-menu--collapsed) {
  .n-menu-item-content {
    &::before {
      left: 6px;
      right: 6px;
      border-radius: 10px;
    }
    &.n-menu-item-content--selected::before {
      border-left: 3px solid rgb(var(--primary-color));
      background: rgb(49 108 114 / 8%);
    }
  }
}
</style>
