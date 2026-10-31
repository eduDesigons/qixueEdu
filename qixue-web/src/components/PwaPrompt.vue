<template>
  <div class="qx-pwa">
    <transition name="qx-pwa-slide">
      <div v-if="needRefresh" class="qx-pwa__banner qx-pwa__banner--update" role="status">
        <div class="qx-pwa__text">
          <strong>发现新版本</strong>
          <span>更新后可获得最新功能与修复</span>
        </div>
        <div class="qx-pwa__actions">
          <button type="button" class="qx-pwa__btn qx-pwa__btn--ghost" @click="laterUpdate">稍后</button>
          <button type="button" class="qx-pwa__btn qx-pwa__btn--primary" :disabled="updating" @click="doUpdate">
            {{ updating ? '更新中…' : '立即更新' }}
          </button>
        </div>
      </div>
    </transition>

    <transition name="qx-pwa-slide">
      <div v-if="showInstall && !needRefresh" class="qx-pwa__banner qx-pwa__banner--install" role="dialog" aria-label="安装到桌面">
        <div class="qx-pwa__text">
          <strong>安装到桌面</strong>
          <span v-if="iosGuide">分享 →「添加到主屏幕」，无浏览器地址栏，更像 App</span>
          <span v-else>独立窗口打开，隐藏地址栏，使用体验更接近 App</span>
        </div>
        <div class="qx-pwa__actions">
          <button type="button" class="qx-pwa__btn qx-pwa__btn--ghost" @click="dismissInstall">稍后</button>
          <button v-if="!iosGuide" type="button" class="qx-pwa__btn qx-pwa__btn--primary" @click="doInstall">安装</button>
          <button v-else type="button" class="qx-pwa__btn qx-pwa__btn--primary" @click="dismissInstall">知道了</button>
        </div>
      </div>
    </transition>
  </div>
</template>

<script>
import {
  applyUpdate,
  dismissInstallPrompt,
  dismissUpdatePrompt,
  promptInstall,
  subscribePwa,
} from '@/common/pwa'

export default {
  name: 'PwaPrompt',
  data() {
    return {
      needRefresh: false,
      showInstall: false,
      iosGuide: false,
      updating: false,
      unsub: null,
    }
  },
  mounted() {
    this.unsub = subscribePwa((s) => {
      this.needRefresh = !!s.needRefresh
      this.showInstall = !!s.showInstall
      this.iosGuide = !!s.iosGuide
    })
  },
  beforeUnmount() {
    if (this.unsub) this.unsub()
  },
  methods: {
    laterUpdate() {
      dismissUpdatePrompt()
    },
    async doUpdate() {
      if (this.updating) return
      this.updating = true
      try {
        await applyUpdate()
      }
      finally {
        this.updating = false
      }
    },
    dismissInstall() {
      dismissInstallPrompt()
    },
    async doInstall() {
      await promptInstall()
    },
  },
}
</script>

<style scoped>
.qx-pwa {
  position: fixed;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 100050;
  pointer-events: none;
  padding: 12px 12px calc(12px + env(safe-area-inset-bottom, 0px));
  display: flex;
  justify-content: center;
}

.qx-pwa__banner {
  pointer-events: auto;
  width: min(520px, 100%);
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 14px;
  border-radius: 14px;
  background: rgba(17, 24, 39, 0.94);
  color: #fff;
  box-shadow: 0 10px 28px rgba(15, 23, 42, 0.28);
  backdrop-filter: blur(10px);
}

.qx-pwa__text {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.qx-pwa__text strong {
  font-size: 14px;
  font-weight: 600;
  line-height: 1.3;
}

.qx-pwa__text span {
  font-size: 12px;
  line-height: 1.4;
  color: rgba(255, 255, 255, 0.72);
}

.qx-pwa__actions {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-shrink: 0;
}

.qx-pwa__btn {
  border: 0;
  border-radius: 999px;
  padding: 7px 12px;
  font-size: 13px;
  line-height: 1.2;
  cursor: pointer;
}

.qx-pwa__btn:disabled {
  opacity: 0.65;
  cursor: default;
}

.qx-pwa__btn--ghost {
  background: transparent;
  color: rgba(255, 255, 255, 0.78);
}

.qx-pwa__btn--primary {
  background: var(--qx-primary, #4286f5);
  color: #fff;
  font-weight: 600;
}

.qx-pwa-slide-enter-active,
.qx-pwa-slide-leave-active {
  transition: transform 0.28s ease, opacity 0.28s ease;
}

.qx-pwa-slide-enter-from,
.qx-pwa-slide-leave-to {
  transform: translateY(120%);
  opacity: 0;
}

@media (min-width: 960px) {
  .qx-pwa {
    padding-bottom: 20px;
  }
}
</style>
