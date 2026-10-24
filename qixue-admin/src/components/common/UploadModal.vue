<script setup>
import { ref } from 'vue'

const props = defineProps({
  size: { type: Number },
  title: { type: String, default: '裁剪头像' },
  circled: { type: Boolean, default: true },
  height: { type: Number, default: 600 },
  width: { type: Number, default: 600 },
  autoUpload: { type: Boolean, default: true }, // 是否自动上传
  options: { type: Object, default: () => ({}) },
})
const emit = defineEmits(['success'])
const src = ref('')

function handleOk(obj) {
  if (props.autoUpload) {
    // 实现自动上传逻辑

  }
  else {
    src.value = ''
    emit('success', obj.b64)
  }
}

function handleBeforeUpload({ file, onFinish }) {
  if (props.size && file.size > 1024 * 1024 * props.size) {
    src.value = ''
    $message.warning(`文件太大了，不能超过${props.size}M`)
    return
  }
  const reader = new FileReader()
  reader.readAsDataURL(file.file)
  src.value = ''
  reader.onload = function (e) {
    src.value = (e.target?.result) ?? ''
  }
  onFinish()
}
</script>

<template>
  <div class="box">
    <slot>
      <n-button type="primary">
        上传
      </n-button>
    </slot>
    <div class="upload">
      <n-upload :custom-request="handleBeforeUpload" accept="image/png, image/jpeg, image/gif, image/jpg" :show-file-list="false">
        <div class="upload" />
      </n-upload>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.box {
  position: relative;
  .upload {
    position: absolute;
    left: 0;
    right: 0;
    top: 0;
    bottom: 0;
    cursor: pointer;
  }
}
</style>
