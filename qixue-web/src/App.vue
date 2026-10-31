<template>
  <DesktopHeader v-if="isDesktop && $route.path !== '/pages/index/login'" :tabs="tabs" :active="$route.path" :unread="unreadMessageCount" />
  <div class="qx-page" :class="[routeClass, { 'has-nav': showNav, 'has-tab': showTab, 'has-crumb': showDesktopCrumb, 'is-desktop': isDesktop, 'is-me-shell': showMeShell }]">
    <MeShell v-if="showMeShell" :path="$route.path" :unread="unreadMessageCount" />
    <div class="qx-page__body">
    <div v-if="showDesktopCrumb" class="qx-desktop-crumb">
      <button type="button" @click="back">返回</button>
      <h1>{{ title }}</h1>
    </div>
    <h1 v-if="showMeShell && $route.path !== '/pages/index/me'" class="pc-me-shell__title">{{ title }}</h1>
    <uni-page-head v-if="showNav">
      <button class="uni-page-head-btn" type="button" @click="back">‹</button>
      <div class="uni-page-head__title">{{ title }}</div>
    </uni-page-head>
    <router-view v-slot="{ Component, route }">
      <transition :name="route.meta && route.meta.tab ? '' : 'qx-page'" :mode="route.meta && route.meta.tab ? undefined : 'out-in'">
        <keep-alive :include="tabNames">
          <component :is="Component" :key="route.meta.tab ? route.path : route.fullPath" />
        </keep-alive>
      </transition>
    </router-view>
    <uni-tabbar v-if="showTab">
      <button
        v-for="item in tabs"
        :key="item.path"
        class="uni-tabbar__item"
        :class="{ on: $route.path === item.path }"
        type="button"
        @click="openTab(item.path)"
      >
        <img class="uni-tabbar__icon" :src="$route.path === item.path ? item.on : item.icon" alt="" />
        <span class="uni-tabbar__label">{{ item.text }}</span>
        <span v-if="item.badge && unreadMessageCount > 0" class="uni-tabbar__badge">{{ unreadMessageCount > 99 ? '99+' : unreadMessageCount }}</span>
      </button>
    </uni-tabbar>
    </div>
  </div>
  <PwaPrompt />
</template>

<script>
import { getStorage } from '@/common/storage'
import { navState, navigateBack, reLaunch, switchTab } from '@/common/nav'
import { defineAsyncComponent } from 'vue'
import { mapState } from 'vuex'
import PwaPrompt from './components/PwaPrompt.vue'
import { tabNames } from './router'
import { isDesktop as desktopOn, bindDesktop } from './pc/desktop'

const tabs = [
  { path: '/pages/index/home', text: '发现', icon: '/static/tabbar/home.png', on: '/static/tabbar/home-on.png' },
  { path: '/pages/index/course', text: '课程', icon: '/static/tabbar/book.png', on: '/static/tabbar/book-on.png' },
  { path: '/pages/index/learn', text: '学习', icon: '/static/tabbar/learn.png', on: '/static/tabbar/learn-on.png' },
  { path: '/pages/index/me', text: '我的', icon: '/static/tabbar/me.png', on: '/static/tabbar/me-on.png', badge: true },
]

export default {
  components: {
    DesktopHeader: defineAsyncComponent(() => import('./pc/DesktopHeader.vue')),
    MeShell: defineAsyncComponent(() => import('./pc/MeShell.vue')),
    PwaPrompt,
  },
  data() {
    return { tabNames, tabs }
  },
  computed: {
    ...mapState({ unreadMessageCount: 'unreadMessageCount' }),
    isDesktop() {
      return desktopOn.value
    },
    showNav() {
      if (this.isDesktop) return false
      return this.$route.meta && this.$route.meta.custom === false
    },
    showTab() {
      return !!(this.$route.meta && this.$route.meta.tab) && !this.isDesktop
    },
    showMeShell() {
      if (!this.isDesktop) return false
      const path = this.$route.path || ''
      if (path === '/pages/member/rankHistory') return false
      return path === '/pages/index/me' || path.startsWith('/pages/member/')
    },
    showDesktopCrumb() {
      if (!this.isDesktop || this.showMeShell || (this.$route.meta && this.$route.meta.tab)) return false
      const path = this.$route.path
      return path !== '/pages/index/login' && path !== '/pages/course/courseDetail'
    },
    title() {
      return navState.override || (this.$route.meta && this.$route.meta.title) || ''
    },
    routeClass() {
      const path = (this.$route.path || '').replace(/^\//, '').replace(/[^a-zA-Z0-9]+/g, '-')
      return path ? `qx-route-${path}` : ''
    },
  },
  mounted() {
    bindDesktop()
    this.initLogin()
  },
  methods: {
    back() {
      navigateBack()
    },
    openTab(path) {
      switchTab(path)
    },
    async initLogin() {
      try {
        this.$store.dispatch('getSysInfo')
        const token = getStorage('token')
        if (token) {
          this.$store.commit('setToken', token)
          this.$store.commit('setUserInfo', getStorage('userInfo'))
          this.$store.dispatch('getUserInfo').catch(() => {})
        } else if (this.$route.path !== '/pages/index/login') {
          reLaunch('/pages/index/login')
        }
      } catch (e) {}
    },
  },
}
</script>
