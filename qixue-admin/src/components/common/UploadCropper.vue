<script setup>
import { NIcon, NUpload } from 'naive-ui'
import { ref } from 'vue'
import api from '@/api'
import { getAssetUrl } from '@/utils/common'

const props = defineProps({
  title: { type: String, default: '上传图片' },
  circled: { type: Boolean, default: false },
})

const emit = defineEmits(['success'])

const imageUrl = ref('')
const uploading = ref(false)

function handleBeforeUpload({ file, onFinish }) {
  if (!file || !file.type) {
    $message.error('请选择图片文件')
    return
  }

  uploading.value = true
  api.upload(file.file, 1).then((res) => {
    imageUrl.value = getAssetUrl(res.data.url)
    emit('success', res.data.url)
    $message.success('上传成功')
  }).catch(() => {
    $message.error('上传失败')
  }).finally(() => {
    uploading.value = false
    onFinish()
  })
}
</script>

<template>
  <div class="upload-cropper">
    <slot>
      <div v-if="imageUrl" class="image-preview">
        <img :src="imageUrl" alt="preview">
      </div>
      <div v-else class="upload-placeholder">
        <NIcon size="40">
          <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
            <path fill="currentColor" d="M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2z" />
          </svg>
        </NIcon>
        <div>{{ title }}</div>
      </div>
    </slot>
    <div class="upload-overlay">
      <NUpload
        :custom-request="handleBeforeUpload"
        accept="image/png, image/jpeg, image/gif, image/jpg, image/webp"
        :show-file-list="false"
      >
        <template #trigger>
          <div class="upload-area" />
        </template>
      </NUpload>
    </div>
  </div>
</template>

<style scoped lang="scss">
.upload-cropper {
  position: relative;
  width: 100%;
  height: 100%;

  .image-preview {
    width: 100%;
    height: 100%;
    border-radius: v-bind('props.circled ? "50%" : "8px"');
    overflow: hidden;

    img {
      width: 100%;
      height: 100%;
      object-fit: cover;
    }
  }

  .upload-placeholder {
    width: 100%;
    height: 100%;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    border: 1px solid var(--ctl-hairline, rgba(var(--primary-color), 0.08));
    border-radius: v-bind('props.circled ? "50%" : "12px"');
    color: rgb(var(--primary-color));
    background: var(--ctl-field, rgba(var(--primary-color), 0.04));
    box-shadow: inset 0 1px 0 rgb(255 255 255 / 40%);
    cursor: pointer;
    transition:
      border-color 0.22s ease,
      background 0.22s ease,
      box-shadow 0.28s cubic-bezier(0.32, 0.72, 0, 1),
      transform 0.22s cubic-bezier(0.32, 0.72, 0, 1);

    &:hover {
      border-color: rgba(var(--primary-color), 0.28);
      background: var(--ctl-field-hover, rgba(var(--primary-color), 0.07));
      box-shadow:
        inset 0 1px 0 rgb(255 255 255 / 45%),
        0 0 0 4px var(--ctl-ring, rgba(var(--primary-color), 0.12));
    }

    &:active {
      transform: scale(0.98);
    }
  }

  .upload-overlay {
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    cursor: pointer;
    z-index: 10;

    .upload-area {
      width: 100%;
      height: 100%;
      position: absolute;
      top: 0;
      left: 0;
    }
  }
}
</style>
