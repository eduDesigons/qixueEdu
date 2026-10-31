<!--------------------------------
 - @Author: Ronnie Zhang
 - @LastEditor: Ronnie Zhang
 - @LastEditTime: 2023/12/05 21:28:36
 - @Email: zclzone@outlook.com
 - Copyright © 2023 Ronnie Zhang(大脸怪) | https://isme.top
 --------------------------------->

<template>
  <div class="login-page">
    <LayoutAurora
      preset="hero"
      :colors="['#a8c9cd', '#FFFFFF', '#FFFFFF']"
    />

    <div class="login-page__body">
      <div class="login-page__main">
        <div class="login-card">
          <h2 class="login-card__title">
            <img :src="logoSrc" class="login-card__logo" alt="">
            {{ siteName }}
          </h2>

          <n-input
            v-model:value="loginInfo.username"
            :autofocus="!loginInfo.username"
            class="mt-32 h-40 items-center"
            placeholder="请输入用户名"
            :maxlength="20"
            @keydown.enter.prevent="focusPassword"
          >
            <template #prefix>
              <i class="i-fe:user mr-12 opacity-20" />
            </template>
          </n-input>
          <n-input
            ref="passwordInput"
            v-model:value="loginInfo.password"
            :autofocus="!!loginInfo.username"
            class="mt-20 h-40 items-center"
            type="password"
            show-password-on="mousedown"
            placeholder="请输入密码"
            :maxlength="20"
            @keydown.enter="handleLogin()"
          >
            <template #prefix>
              <i class="i-fe:lock mr-12 opacity-20" />
            </template>
          </n-input>

          <n-checkbox
            class="mt-20"
            :checked="isRemember"
            label="记住我"
            :on-update:checked="(val) => (isRemember = val)"
          />

          <div class="mt-20 flex items-center">
            <n-button
              class="h-40 flex-1 rounded-12 text-16"
              type="primary"
              :loading="loading"
              @click="handleLogin()"
            >
              登录
            </n-button>
          </div>
        </div>
      </div>

      <TheFooter class="login-page__footer" />
    </div>

    
    <SlideCaptcha
      v-if="showCaptcha"
      ref="captchaRef"
      @success="onCaptchaSuccess"
    />
  </div>
</template>

<script setup>
import { nextTick } from 'vue'
import { useStorage } from '@vueuse/core'
import { useAuthStore, useSysStore } from '@/store'
import { getAssetUrl, lStorage, throttle } from '@/utils'
import defaultLogo from '@/assets/images/logo.png'
import { LayoutAurora } from '@/layouts/components'
import SlideCaptcha from '@/components/common/SlideCaptcha.vue'
import api from './api'

const authStore = useAuthStore()
const sysStore = useSysStore()
const router = useRouter()
const route = useRoute()
const siteName = computed(() => String(sysStore.sysConfig?.name || '').trim())
const logoSrc = computed(() => {
  const logo = String(sysStore.sysConfig?.logo || '').trim()
  if (logo)
    return getAssetUrl(logo) || defaultLogo
  return defaultLogo
})

const captchaRef = ref(null)
const passwordInput = ref(null)
const showCaptcha = ref(false)
const pendingLogin = ref(false)
const loginInfo = ref({
  username: '',
  password: '',
  captchaVerification: '',
})

const initCaptcha = throttle(() => {
  loginInfo.value.captchaVerification = ''
  showCaptcha.value = true
  nextTick(() => {
    captchaRef.value?.show()
  })
  pendingLogin.value = true
}, 500)

function onCaptchaSuccess(payload) {
  loginInfo.value.captchaVerification = payload?.captchaVerification ?? payload
  if (pendingLogin.value) {
    pendingLogin.value = false
    doLogin()
  }
}

const localLoginInfo = lStorage.get('loginInfo')
if (localLoginInfo) {
  loginInfo.value.username = localLoginInfo.username || ''
  loginInfo.value.password = localLoginInfo.password || ''
}

const isRemember = useStorage('isRemember', true)
const loading = ref(false)

function focusPassword() {
  passwordInput.value?.focus?.()
}

async function doLogin() {
  try {
    loading.value = true
    const { data } = await api.login(loginInfo.value)
    if (isRemember.value) {
      lStorage.set('loginInfo', { username: loginInfo.value.username, password: loginInfo.value.password })
    }
    else {
      lStorage.remove('loginInfo')
    }
    onLoginSuccess(data)
  }
  catch (error) {
    if (error?.code === 10003) {
      initCaptcha()
    }
    $message.destroy('login')
    console.error(error)
  }
  loading.value = false
}

async function handleLogin(isQuick) {
  const { username, password, captchaVerification } = loginInfo.value
  if (!username || !password)
    return $message.warning('请输入用户名和密码')
  if (isQuick) {
    return doLogin()
  }
  if (!captchaVerification) {
    pendingLogin.value = true
    showCaptcha.value = true
    nextTick(() => captchaRef.value?.show())
    return
  }
  doLogin()
}

async function onLoginSuccess(data = {}) {
  authStore.setToken(data)
  $message.loading('登录中...', { key: 'login' })
  try {
    $message.success('登录成功', { key: 'login' })
    if (route.query.redirect) {
      const path = route.query.redirect
      delete route.query.redirect
      router.push({ path, query: route.query })
    }
    else {
      router.push('/')
    }
  }
  catch (error) {
    console.error(error)
    $message.destroy('login')
  }
}
</script>

<style scoped>
.login-page {
  position: relative;
  width: 100%;
  height: 100%;
  overflow: hidden;
  background: #fff;
}

.login-page__body {
  position: relative;
  z-index: 1;
  display: flex;
  flex-direction: column;
  width: 100%;
  height: 100%;
  box-sizing: border-box;
}

.login-page__main {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24px 16px;
  box-sizing: border-box;
}

.login-card {
  width: 100%;
  max-width: 400px;
  padding: 40px 36px 36px;
  border-radius: 16px;
  background: rgb(255 255 255 / 78%);
  border: 1px solid rgb(49 108 114 / 8%);
  box-shadow:
    0 1px 2px rgb(49 108 114 / 4%),
    0 18px 48px rgb(49 108 114 / 8%);
  backdrop-filter: blur(10px);
}

.login-card__title {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12px;
  margin: 0;
  font-size: 22px;
  font-weight: 500;
  color: #316c72;
  letter-spacing: 0.02em;
}

.login-card__logo {
  height: 44px;
  width: auto;
}

.login-page__footer {
  padding: 8px 16px 16px;
  opacity: 0.72;
}
</style>
