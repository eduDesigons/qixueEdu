<!--------------------------------
 - @Author: Ronnie Zhang
 - @LastEditor: Ronnie Zhang
 - @LastEditTime: 2023/12/16 18:50:28
 - @Email: zclzone@outlook.com
 - Copyright © 2023 Ronnie Zhang(大脸怪) | https://isme.top
 --------------------------------->

<template>
  <router-link class="layout-aside__brand brand-link f-c-c" to="/" :title="siteName">
    <div class="brand-logo">
      <img :src="logoSrc" :alt="siteName">
    </div>
    <h2
      v-show="!appStore.collapsed"
      class="brand-title ml-8 max-w-112 flex-shrink-0"
    >
      {{ siteName }}
    </h2>
  </router-link>
</template>

<script setup>
import { computed } from 'vue'
import defaultLogo from '@/assets/images/logo.png'
import { useAppStore, useSysStore } from '@/store'
import { getAssetUrl } from '@/utils'

const appStore = useAppStore()
const sysStore = useSysStore()

const siteName = computed(() => String(sysStore.sysConfig?.name || '').trim())

const logoSrc = computed(() => {
  const logo = String(sysStore.sysConfig?.logo || '').trim()
  if (logo)
    return getAssetUrl(logo) || logo
  return defaultLogo
})
</script>

<style scoped>
.brand-link {
  height: var(--bar-h, 52px);
  min-height: var(--bar-h, 52px);
  max-height: var(--bar-h, 52px);
  flex-shrink: 0;
  text-decoration: none;
  border-radius: 10px;
  margin: 0 6px;
  transition:
    background 0.22s ease,
    transform 0.22s cubic-bezier(0.32, 0.72, 0, 1);
}

.brand-link:hover {
  background: var(--ctl-field, rgba(var(--primary-color), 0.06));
}

.brand-link:active {
  transform: scale(0.98);
}

.brand-logo {
  width: 32px;
  height: 32px;
  flex-shrink: 0;
  background: transparent;
}

.brand-logo img {
  display: block;
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.brand-title {
  margin: 0;
  font-family: var(--font-display);
  font-size: 15px;
  font-weight: 600;
  letter-spacing: -0.028em;
  color: rgb(var(--primary-color));
  line-height: 1.2;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.brand-link:hover .brand-title {
  filter: brightness(1.06);
}
</style>
