<template>
  <header class="qx-desktop-header">
    <div class="qx-desktop-header__inner">
      <button class="qx-desktop-header__brand" type="button" @click="open('/pages/index/home')">{{ siteName }}</button>
      <nav class="qx-desktop-header__nav">
        <button
          v-for="item in tabs"
          :key="item.path"
          type="button"
          class="qx-desktop-header__link"
          :class="{ 'is-active': tabActive === item.path }"
          @click="open(item.path)"
        >
          <span>{{ item.text }}</span>
          <em v-if="item.badge && unread > 0">{{ unread > 99 ? '99+' : unread }}</em>
        </button>
      </nav>
      <form class="qx-desktop-header__search" @submit.prevent="submit">
        <svg viewBox="0 0 16 16" aria-hidden="true">
          <circle cx="7" cy="7" r="4.25" fill="none" stroke="currentColor" stroke-width="1.5" />
          <path d="M10.2 10.2 13.2 13.2" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" />
        </svg>
        <input v-model="keyword" type="text" placeholder="搜索课程名称" autocomplete="off" />
        <button v-if="keyword" class="qx-desktop-header__clear" type="button" aria-label="清空" @click="clear">×</button>
      </form>
      <button v-if="userName" class="qx-desktop-header__user" type="button" @click="open('/pages/index/me')">
        <img v-if="userAvatar" :src="userAvatar" alt="" />
        <i v-else>{{ userName.slice(0, 1) }}</i>
        <span>{{ userName }}</span>
        <coin-badge v-if="userInfo && userInfo.id" :coins="userInfo.coins" />
      </button>
      <button v-else class="qx-desktop-header__login" type="button" @click="login">登录</button>
    </div>
  </header>
</template>

<script>
import { emit, off, on } from '@/common/bus'
import { navigateTo, switchTab } from '@/common/nav'
import { mapState } from 'vuex'

export default {
  name: 'DesktopHeader',
  props: {
    tabs: { type: Array, default: () => [] },
    active: { type: String, default: '' },
    unread: { type: Number, default: 0 },
  },
  data() {
    return { keyword: '' }
  },
  computed: {
    ...mapState({ userInfo: 'userInfo' }),
    userName() {
      return (this.userInfo && this.userInfo.nick_name) || ''
    },
    userAvatar() {
      const avatar = this.userInfo && this.userInfo.avatar
      return avatar ? this.static(avatar) : ''
    },
    tabActive() {
      const path = this.active || ''
      if (path.startsWith('/pages/course/')) return '/pages/index/course'
      if (path.startsWith('/pages/exam/')) return '/pages/index/learn'
      if (path.startsWith('/pages/member/')) return '/pages/index/me'
      return this.tabs.some((item) => item.path === path) ? path : ''
    },
  },
  mounted() {
    this.onSync = (keyword) => { this.keyword = keyword || '' }
    on('courseSearchSync', this.onSync)
  },
  beforeUnmount() {
    off('courseSearchSync', this.onSync)
  },
  methods: {
    open(path) {
      if (!path || path === this.active) return
      switchTab(path)
    },
    login() {
      navigateTo('/pages/index/login')
    },
    clear() {
      this.keyword = ''
      this.submit()
    },
    submit() {
      const keyword = this.keyword.trim()
      const go = () => emit('coursePageSearch', keyword)
      if (this.active === '/pages/index/course') {
        go()
        return
      }
      switchTab('/pages/index/course')
    },
  },
}
</script>

<style scoped>
.qx-desktop-header {
  position: sticky;
  top: 0;
  z-index: 1000;
  height: 64px;
  background: #fff;
  border-bottom: 1px solid #eef1f6;
  box-shadow: 0 1px 2px rgba(26, 26, 26, 0.04);
}
.qx-desktop-header__inner {
  display: flex;
  align-items: center;
  width: min(1180px, calc(100% - 48px));
  height: 64px;
  margin: 0 auto;
}
.qx-desktop-header__brand {
  flex: none;
  margin-right: 28px;
  color: #1a1a1a;
  font-size: calc(20px * var(--font-scale, 1));
  font-weight: 700;
  letter-spacing: 0.5px;
  cursor: pointer;
}
.qx-desktop-header__nav {
  display: flex;
  align-items: center;
  height: 64px;
  gap: 4px;
}
.qx-desktop-header__link {
  position: relative;
  display: flex;
  align-items: center;
  height: 36px;
  padding: 0 14px;
  border-radius: 8px;
  color: #5c6370;
  font-size: calc(15px * var(--font-scale, 1));
  cursor: pointer;
}
.qx-desktop-header__link:hover {
  color: #1a1a1a;
  background: #f6f8fb;
}
.qx-desktop-header__link.is-active {
  color: #1a1a1a;
  font-weight: 650;
  background: transparent;
}
.qx-desktop-header__link.is-active span {
  position: relative;
}
.qx-desktop-header__link.is-active span::after {
  content: '';
  position: absolute;
  left: 50%;
  bottom: -10px;
  width: 16px;
  height: 3px;
  margin-left: -8px;
  border-radius: 2px;
  background: #4285f4;
}
.qx-desktop-header__link em {
  min-width: 16px;
  height: 16px;
  margin-left: 6px;
  padding: 0 4px;
  border-radius: 8px;
  background: #f5483b;
  color: #fff;
  font-size: calc(11px * var(--font-scale, 1));
  font-style: normal;
  font-weight: 600;
  line-height: 16px;
  text-align: center;
}
.qx-desktop-header__search {
  display: flex;
  align-items: center;
  flex: 0 1 260px;
  width: 260px;
  height: 36px;
  margin-left: auto;
  padding: 0 10px 0 12px;
  border-radius: 18px;
  background: #f4f6f8;
  color: #9ca3af;
  gap: 8px;
}
.qx-desktop-header__search:focus-within {
  background: #fff;
  box-shadow: 0 0 0 1.5px #4285f4;
  color: #4285f4;
}
.qx-desktop-header__search svg {
  flex: none;
  width: 15px;
  height: 15px;
}
.qx-desktop-header__search input {
  flex: 1;
  min-width: 0;
  border: 0;
  outline: none;
  background: transparent;
  color: #1a1a1a;
  font-size: calc(13px * var(--font-scale, 1));
}
.qx-desktop-header__search input::placeholder {
  color: #9ca3af;
}
.qx-desktop-header__clear {
  flex: none;
  width: 16px;
  height: 16px;
  border-radius: 50%;
  background: #d5d9e0;
  color: #fff;
  font-size: calc(12px * var(--font-scale, 1));
  line-height: 16px;
  cursor: pointer;
}
.qx-desktop-header__login {
  flex: none;
  height: 32px;
  margin-left: 16px;
  padding: 0 16px;
  border-radius: 16px;
  background: #4285f4;
  color: #fff;
  font-size: calc(13px * var(--font-scale, 1));
  font-weight: 500;
  cursor: pointer;
}
.qx-desktop-header__login:hover {
  background: #3a78e6;
}
.qx-desktop-header__user {
  display: flex;
  align-items: center;
  flex: none;
  max-width: 240px;
  height: 36px;
  margin-left: 12px;
  padding: 0 10px 0 4px;
  border-radius: 18px;
  color: #1a1a1a;
  font-size: calc(14px * var(--font-scale, 1));
  gap: 8px;
  cursor: pointer;
}
.qx-desktop-header__user:hover {
  background: #f4f6f8;
}
.qx-desktop-header__user img,
.qx-desktop-header__user i {
  flex: none;
  width: 28px;
  height: 28px;
  border-radius: 50%;
  background: #e8f0fe;
  color: #4285f4;
  font-style: normal;
  font-size: calc(13px * var(--font-scale, 1));
  font-weight: 650;
  line-height: 28px;
  text-align: center;
}
.qx-desktop-header__user img {
  object-fit: cover;
  display: block;
}
.qx-desktop-header__user :deep(.coin-lv) {
  margin-left: 0;
}
.qx-desktop-header__user :deep(.coin-lv__icon) {
  width: 16px;
  height: 16px;
  border-radius: 0;
  background: transparent;
  object-fit: contain;
}
.qx-desktop-header__user span {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}
</style>
