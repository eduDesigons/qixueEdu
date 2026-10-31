<!--------------------------------
 - @Author: Ronnie Zhang
 - @LastEditor: Ronnie Zhang
 - @LastEditTime: 2023/12/04 22:51:21
 - @Email: zclzone@outlook.com
 - Copyright © 2023 Ronnie Zhang(大脸怪) | https://isme.top
 --------------------------------->

<template>
  <main class="h-full flex-col flex-1 overflow-hidden min-h-0">
    <header
      v-if="showHeader && (back || $slots.header)"
      class="page-nav"
      :class="{ 'page-nav--ghost': back && !$slots.header }"
    >
      <slot v-if="$slots.header" name="header" />
      <template v-else>
        <div class="page-nav__leading">
          <button
            v-if="back"
            type="button"
            class="page-back"
            aria-label="返回"
            @click="handleBack"
          >
            <i class="i-material-symbols:chevron-left page-back__icon" />
            <span class="page-back__label">返回</span>
          </button>
          <slot name="title-prefix" />

          <h2 v-if="showNavTitle" class="page-title">
            {{ navTitle }}
          </h2>
          <slot name="title-suffix" />
        </div>
        <div class="page-nav__action">
          <slot name="action" />
        </div>
      </template>
    </header>

    <AppCard
      class="cus-scroll page-body page-surface h-0 min-h-0 flex-1"
      :class="{ 'page-body--fill': !back || fill, 'page-body--detail': back }"
      :bordered="false"
    >
      <header v-if="headlineVisible && !$slots.header" class="page-mast">
        <div class="page-mast__row">
          <div class="page-mast__copy">
            <h1 class="page-mast__title">
              {{ headlineText }}
            </h1>
            <p v-if="subtitle" class="page-mast__sub">
              {{ subtitle }}
            </p>
          </div>
          <div v-if="!back && showHeader" class="page-mast__toolbar">
            <div
              :ref="setFilterHost"
              class="page-filter-host"
            />
            <slot name="title-prefix" />
            <slot name="title-suffix" />
            <div class="page-nav__action">
              <slot name="action" />
            </div>
          </div>
        </div>
      </header>
      <div class="page-body__main">
        <slot />
      </div>
    </AppCard>

    <slot name="footer">
      <AppCard v-if="showFooter" class="flex-shrink-0 py-12">
        <TheFooter />
      </AppCard>
    </slot>
  </main>
</template>

<script setup>
import { AppCard } from '@/components/index.js'
import { useTabStore } from '@/store'

const props = defineProps({
  back: {
    type: Boolean,
    default: false,
  },
  showFooter: {
    type: Boolean,
    default: false,
  },
  showHeader: {
    type: Boolean,
    default: true,
  },
  title: {
    type: String,
    default: undefined,
  },
  headline: {
    type: String,
    default: undefined,
  },
  subtitle: {
    type: String,
    default: undefined,
  },
  showHeadline: {
    type: Boolean,
    default: undefined,
  },
  fill: {
    type: Boolean,
    default: false,
  },
})
const route = useRoute()
const router = useRouter()
const tabStore = useTabStore()
const filterHostRef = ref(null)
provide('pageFilterHost', filterHostRef)

function setFilterHost(el) {
  filterHostRef.value = el
}

const navTitle = computed(() => props.title ?? route.meta?.title)
const headlineText = computed(() => props.headline ?? props.title ?? route.meta?.title)
const showNavTitle = computed(() => props.back && !headlineVisible.value)
const headlineVisible = computed(() => {
  if (typeof props.showHeadline === 'boolean')
    return props.showHeadline
  return !props.back
})

function handleBack() {
  const currentPath = route.path
  // 详情返回：关掉当前标签后走浏览器历史，才能回到列表并保留翻页
  tabStore.closeTabSilent(currentPath)
  if (window.history.length > 1) {
    router.back()
    return
  }
  const fallback = tabStore.resolvePreviousTab(currentPath)
    || tabStore.tabs[tabStore.tabs.length - 1]?.path
  if (fallback)
    router.push(fallback)
  else
    router.push('/')
}
</script>

<style scoped>
.page-body {
  padding: 16px var(--content-inset-x, 20px) 24px;
  min-height: 0 !important;
  overflow-x: hidden !important;
  overflow-y: auto !important;
  -webkit-overflow-scrolling: touch;
}

.page-body--fill {
  display: flex !important;
  flex-direction: column;
  overflow: hidden !important;
}

.page-body--fill .page-mast {
  flex-shrink: 0;
}

.page-body--detail {
  padding-top: 8px;
}

.page-body__main {
  min-height: 0;
}

.page-body--fill .page-body__main {
  flex: 1 1 auto;
  min-height: 0;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.page-body--detail.page-body--fill .page-body__main,
.page-body--detail .page-body__main {
  gap: 36px;
}

.page-body--fill .page-body__main > * {
  flex: 1 1 auto;
  min-height: 0;
}

.page-body--fill .page-body__main > .n-spin-container,
.page-body--fill .page-body__main > .n-spin {
  flex: 0 0 auto;
  height: auto;
  min-height: 0;
}
</style>
