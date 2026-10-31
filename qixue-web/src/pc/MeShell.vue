<template>
  <aside class="pc-me-nav">
    <button type="button" class="pc-me-nav__user" @click="go('/pages/member/profile')">
      <img v-if="userInfo && userInfo.avatar" :src="static(userInfo.avatar)" alt="" />
      <b v-else>{{ initial }}</b>
      <strong>{{ (userInfo && userInfo.nick_name) || '未登录' }}</strong>
      <coin-badge v-if="userInfo && userInfo.id" :coins="userInfo.coins" />
    </button>
    <div v-for="group in groups" :key="group.title">
      <h2>{{ group.title }}</h2>
      <button
        type="button"
        v-for="item in group.items"
        :key="item.path"
        :class="{ on: isOn(item.path) }"
        @click="go(item.path)"
      >
        <img v-if="item.icon" :src="item.icon" alt="" />
        <van-icon v-else-if="item.uicon" :name="item.uicon" size="22" :color="isOn(item.path) ? '#4285f4' : '#8a9199'"  />
        <svg v-else viewBox="0 0 24 24" aria-hidden="true">
          <template v-if="item.svg === 'credit'">
            <rect width="24" height="24" rx="7" fill="#FFE8D6" />
            <circle cx="12" cy="10" r="4" fill="#F5A623" />
            <path d="M9.4 13.2 8.2 19.2 12 17.2l3.8 2-1.2-6" fill="#F08A24" />
          </template>
          <template v-else-if="item.svg === 'coin'">
            <rect width="24" height="24" rx="7" fill="#FFF4D6" />
            <circle cx="12" cy="12" r="6.2" fill="#F5B400" />
            <path d="M12 8.4v7.2M9.7 10.3h3.1a1.5 1.5 0 0 1 0 3h-2.2a1.5 1.5 0 0 0 0 3H14.4" fill="none" stroke="#fff" stroke-width="1.4" stroke-linecap="round" />
          </template>
        </svg>
        <span>{{ item.text }}</span>
        <em v-if="item.badge && unread > 0" class="pc-me-nav__badge">{{ unread > 99 ? '99+' : unread }}</em>
      </button>
    </div>
  </aside>
</template>

<script>
import { navigateTo, redirectTo } from '@/common/nav'
import { mapState } from 'vuex'

export default {
  name: 'MeShell',
  props: {
    path: { type: String, default: '' },
    unread: { type: Number, default: 0 },
  },
  computed: {
    ...mapState({ userInfo: 'userInfo' }),
    initial() {
      const name = (this.userInfo && this.userInfo.nick_name) || '学'
      return name.slice(0, 1)
    },
  },
  data() {
    return {
      groups: [
        {
          title: '学习内容',
          items: [
            { text: '我的课程', path: '/pages/member/lesson', icon: '/static/img/me/exam.png' },
            { text: '我的考试', path: '/pages/member/exam', icon: '/static/img/me/test.png' },
            { text: '我的测评', path: '/pages/member/test', icon: '/static/img/me/note.png' },
            { text: '我的收藏', path: '/pages/member/collection', icon: '/static/img/me/collect.png' },
            { text: '我的评论', path: '/pages/member/comment', icon: '/static/img/me/commit.png' },
            { text: '学习轨迹', path: '/pages/member/history', icon: '/static/img/me/history.png' },
            { text: '系列课程', path: '/pages/member/path', icon: '/static/img/me/path.png' },
            { text: '我的证书', path: '/pages/member/cert', icon: '/static/img/me/certificate.png' },
            { text: '我的消息', path: '/pages/member/message', icon: '/static/img/me/message.png', badge: true },
          ],
        },
        {
          title: '我的资产',
          items: [
            { text: '我的学分', path: '/pages/member/credit', svg: 'credit' },
            { text: '我的积分', path: '/pages/member/coin', svg: 'coin' },
          ],
        },
        {
          title: '个人设置',
          items: [
            { text: '设置', path: '/pages/member/set', icon: '/static/img/me/set.png' },
            { text: '常见问题', path: '/pages/member/question', icon: '/static/img/me/question.png' },
            { text: '工具库', path: '/pages/member/tools', icon: '/static/img/me/tools.svg' },
          ],
        },
      ],
    }
  },
  methods: {
    isOn(path) {
      if (this.path === path) return true
      if (path === '/pages/member/set' && (this.path === '/pages/member/fontSize' || this.path === '/pages/member/about')) return true
      if (path === '/pages/member/question' && this.path === '/pages/member/feedback') return true
      return false
    },
    go(path) {
      if (!path || path === this.path) return
      if (this.path === '/pages/index/me') {
        navigateTo(path)
        return
      }
      redirectTo(path)
    },
  },
}
</script>

<style scoped>
.pc-me-nav {
  position: sticky;
  top: 80px;
  max-height: calc(100vh - 96px);
  overflow: auto;
  padding: 22px 12px 14px;
  border-radius: 12px;
  background: #fff;
  box-shadow: 0 1px 2px rgba(26, 26, 26, 0.04);
}
.pc-me-nav__user {
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 100%;
  margin-bottom: 6px;
  padding: 2px 8px 16px;
  border-bottom: 1px solid #f2f4f7;
  cursor: pointer;
}
.pc-me-nav__user img,
.pc-me-nav__user b {
  width: 72px;
  height: 72px;
  border-radius: 50%;
  object-fit: cover;
  background: #e8f0fe;
  color: #4285f4;
  font-size: calc(28px * var(--font-scale, 1));
  display: inline-flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 0 0 4px #f4f7fc;
}
.pc-me-nav__user :deep(.coin-lv) {
  justify-content: center;
  max-width: 168px;
  margin-top: 6px;
}
.pc-me-nav__user :deep(.coin-lv__icon) {
  width: 18px;
  height: 18px;
  border-radius: 0;
  box-shadow: none;
  background: transparent;
}
.pc-me-nav__user strong {
  margin-top: 12px;
  color: #1a1a1a;
  font-size: calc(16px * var(--font-scale, 1));
  font-weight: 650;
}
h2 {
  margin: 16px 12px 6px;
  color: #9aa3af;
  font-size: calc(12px * var(--font-scale, 1));
  font-weight: 500;
  letter-spacing: 0.04em;
}
.pc-me-nav button:not(.pc-me-nav__user) {
  display: flex;
  align-items: center;
  gap: 10px;
  width: 100%;
  height: 42px;
  padding: 0 10px;
  border-radius: 8px;
  color: #404040;
  font-size: calc(14px * var(--font-scale, 1));
  text-align: left;
  cursor: pointer;
}
.pc-me-nav button:not(.pc-me-nav__user) img,
.pc-me-nav button:not(.pc-me-nav__user) svg,
.pc-me-nav button:not(.pc-me-nav__user) :deep(.van-icon) {
  width: 22px;
  height: 22px;
  flex: none;
  object-fit: contain;
  line-height: 22px;
  text-align: center;
}
.pc-me-nav button:not(.pc-me-nav__user) span {
  flex: 1;
  min-width: 0;
}
.pc-me-nav button:not(.pc-me-nav__user):hover { background: #f6f7fb; }
.pc-me-nav button.on {
  background: #eef4ff;
  color: #4285f4;
  font-weight: 600;
}
.pc-me-nav__badge {
  min-width: 18px;
  height: 18px;
  padding: 0 5px;
  border-radius: 99px;
  background: #ff4d4f;
  color: #fff;
  font-size: calc(12px * var(--font-scale, 1));
  font-style: normal;
  font-weight: 500;
  line-height: 18px;
  text-align: center;
}
</style>
