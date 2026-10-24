<template>
  <n-modal v-model:show="visible" :mask-closable="false">
    <div class="slide-captcha">
      <div v-if="loading" class="slide-captcha__loading">
        验证码加载中
      </div>
      <Slide
        v-else-if="data"
        :data="data"
        :config="config"
        :events="events"
      />
    </div>
  </n-modal>
</template>

<script setup>
import { Slide } from 'go-captcha-vue'
import 'go-captcha-vue/dist/style.css'
import api from '@/views/login/api'

const emit = defineEmits(['success'])

const visible = ref(false)
const loading = ref(false)
const data = ref(null)
const token = ref('')

const config = {
  width: 300,
  height: 220,
  title: '请完成安全验证',
}

const events = {
  refresh: () => load(),
  close: () => {
    visible.value = false
  },
  confirm: (point, reset) => {
    api.checkCaptcha({
      token: token.value,
      x: Math.round(point.x),
      y: Math.round(point.y),
    }).then((res) => {
      const captchaVerification = res?.data?.captchaVerification
      if (!captchaVerification) {
        reset()
        return
      }
      emit('success', { captchaVerification })
      visible.value = false
    }).catch(() => {
      reset()
      load()
    })
  },
}

async function load() {
  loading.value = true
  data.value = null
  try {
    const res = await api.getCaptcha()
    const body = res?.data
    if (!body?.token || !body?.image || !body?.thumb)
      throw new Error('验证码数据格式异常')
    token.value = body.token
    data.value = {
      image: body.image,
      thumb: body.thumb,
      thumbX: body.thumbX,
      thumbY: body.thumbY,
      thumbWidth: body.thumbWidth,
      thumbHeight: body.thumbHeight,
    }
  }
  catch {
    $message.error('验证码加载失败，请重试')
    visible.value = false
  }
  finally {
    loading.value = false
  }
}

function show() {
  visible.value = true
  load()
}

defineExpose({ show, refresh: load })
</script>

<style scoped>
.slide-captcha {
  width: 300px;
  /* 覆盖 go-captcha 默认蓝，跟后台主题色一致 */
  --go-captcha-theme-active-color: rgb(var(--primary-color));
  --go-captcha-theme-btn-bg-color: rgb(var(--primary-color));
  --go-captcha-theme-btn-border-color: rgb(var(--primary-color));
  --go-captcha-theme-btn-color: #fff;
  --go-captcha-theme-icon-color: rgb(var(--primary-color));
  --go-captcha-theme-drag-bg-color: rgb(var(--primary-color));
  --go-captcha-theme-drag-icon-color: #fff;
  --go-captcha-theme-drag-bar-color: rgba(var(--primary-color), 0.14);
  --go-captcha-theme-border-color: rgba(var(--primary-color), 0.22);
  --go-captcha-theme-loading-icon-color: rgb(var(--primary-color));
}

.slide-captcha__loading {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 220px;
  color: #607d82;
  background: #fff;
  border-radius: 12px;
}
</style>
