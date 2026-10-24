<template>
  <div class="qx-pwa">
    <Transition name="qx-pwa-slide">
      <div v-if="needRefresh" class="qx-pwa__banner" role="status">
        <div class="qx-pwa__text">
          <strong>发现新版本</strong>
          <span>更新后可获得最新功能与修复</span>
        </div>
        <div class="qx-pwa__actions">
          <n-button quaternary size="small" @click="laterUpdate">
            稍后
          </n-button>
          <n-button type="primary" size="small" :loading="updating" @click="doUpdate">
            立即更新
          </n-button>
        </div>
      </div>
    </Transition>

    <Transition name="qx-pwa-slide">
      <div v-if="showInstall && !needRefresh" class="qx-pwa__banner" role="dialog" aria-label="安装到桌面">
        <div class="qx-pwa__text">
          <strong>安装到桌面</strong>
          <span v-if="iosGuide">分享 →「添加到主屏幕」，无浏览器地址栏，更像 App</span>
          <span v-else>独立窗口打开，隐藏地址栏与标签栏，更像桌面应用</span>
        </div>
        <div class="qx-pwa__actions">
          <n-button quaternary size="small" @click="dismissInstall">
            稍后
          </n-button>
          <n-button v-if="!iosGuide" type="primary" size="small" @click="doInstall">
            安装
          </n-button>
          <n-button v-else type="primary" size="small" @click="dismissInstall">
            知道了
          </n-button>
        </div>
      </div>
    </Transition>
  </div>
</template>

<script setup>
import {
  applyUpdate,
  dismissInstallPrompt,
  dismissUpdatePrompt,
  promptInstall,
  subscribePwa,
} from '@/utils/pwa'

const needRefresh = ref(false)
const showInstall = ref(false)
const iosGuide = ref(false)
const updating = ref(false)

let unsub = null

onMounted(() => {
  unsub = subscribePwa((s) => {
    needRefresh.value = !!s.needRefresh
    showInstall.value = !!s.showInstall
    iosGuide.value = !!s.iosGuide
  })
})

onBeforeUnmount(() => {
  if (unsub)
    unsub()
})

function laterUpdate() {
  dismissUpdatePrompt()
}

async function doUpdate() {
  if (updating.value)
    return
  updating.value = true
  try {
    await applyUpdate()
  }
  finally {
    updating.value = false
  }
}

function dismissInstall() {
  dismissInstallPrompt()
}

async function doInstall() {
  await promptInstall()
}
</script>

<style scoped>
.qx-pwa {
  position: fixed;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 10050;
  pointer-events: none;
  padding: 16px;
  display: flex;
  justify-content: center;
}

.qx-pwa__banner {
  pointer-events: auto;
  width: min(480px, 100%);
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 16px;
  border-radius: 14px;
  background: rgba(28, 43, 46, 0.94);
  color: #fff;
  box-shadow: 0 12px 32px rgba(28, 43, 46, 0.24);
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
  gap: 4px;
  flex-shrink: 0;
}

.qx-pwa__actions :deep(.n-button) {
  --n-text-color: rgba(255, 255, 255, 0.82) !important;
}

.qx-pwa__actions :deep(.n-button--primary-type) {
  --n-text-color: #fff !important;
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
</style>
