<template>
  <n-config-provider
    class="wh-full"
    :locale="zhCN"
    :date-locale="dateZhCN"
    :theme-overrides="naiveThemeOverrides"
  >
    <NMessageProvider>
      <NDialogProvider>
        <NNotificationProvider>
          <router-view v-if="Layout" v-slot="{ Component, route: curRoute }">
            <component :is="Layout">
              <!--
                Transition + KeepAlive：有翻页动画且列表可缓存。
                缓存有数量上限与时效，过期后再进入会重建。
              -->
              <div class="page-view">
                <Transition name="fade-slide" appear>
                  <KeepAlive :include="keepAliveInclude" :max="PAGE_CACHE_MAX">
                    <component
                      :is="Component"
                      v-if="!tabStore.reloading && Component"
                      :key="pageCacheKey(curRoute)"
                    />
                  </KeepAlive>
                </Transition>
              </div>
            </component>
          </router-view>
          <AppImagePreview />
          <PwaPrompt />
        </NNotificationProvider>
      </NDialogProvider>
    </NMessageProvider>
  </n-config-provider>
</template>

<script setup>
import { dateZhCN, NDialogProvider, NMessageProvider, NNotificationProvider, zhCN } from 'naive-ui'
import { AppImagePreview, PwaPrompt } from '@/components'
import { naiveThemeOverrides } from '@/settings'
import { useTabStore } from '@/store'
import { PAGE_CACHE_MAX, pruneExpiredListStateCache } from '@/utils/pageCache'

const layouts = new Map()
function getLayout(name) {
  if (layouts.get(name))
    return layouts.get(name)
  const layout = markRaw(defineAsyncComponent(() => import(`@/layouts/${name}/index.vue`)))
  layouts.set(name, layout)
  return layout
}

const route = useRoute()
const Layout = computed(() => {
  if (!route.matched?.length)
    return null
  const name = route.meta?.layout === 'empty' ? 'empty' : 'normal'
  return getLayout(name)
})

const tabStore = useTabStore()

/** 仅缓存未过期且菜单勾了 KeepAlive 的页 */
const keepAliveInclude = computed(() => {
  const names = tabStore.tabs
    .filter(item => item.keepAlive && item.name && tabStore.isTabCacheFresh(item))
    .map(item => String(item.name))
  if (route.meta?.keepAlive && route.name)
    names.push(String(route.name))
  return [...new Set(names)]
})

function pageCacheKey(curRoute) {
  const tab = tabStore.tabs.find(item => item.path === curRoute.path)
  return `${curRoute.path}:${tab?.refreshAt || 0}`
}

onMounted(() => {
  tabStore.pruneExpiredPageCache()
  pruneExpiredListStateCache()
})
</script>
