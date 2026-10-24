<script setup>
import { cloneDeep } from 'lodash-es'
import { NButton, NForm, NFormItem, NInput, NModal, NRadioButton, NRadioGroup, NSwitch } from 'naive-ui'
import { ref } from 'vue'
import api from '@/api'
import Upload from '@/components/common/Upload.vue'
import WangEditor from '@/components/common/WangEditor.vue'

const emit = defineEmits(['update'])

const open = ref(false)
const showLoading = ref(false)

const form = ref({
  id: null,
  name: '',
  img_url: '',
  sort: 1,
  link_url: '',
  type: 1,
  status: 1,
  content: '',
})

function show(data = null) {
  if (data) {
    form.value = cloneDeep(data)
  }
  else {
    form.value = {
      id: null,
      name: '',
      img_url: '',
      sort: 1,
      link_url: '',
      type: 1,
      status: 1,
      content: '',
    }
  }
  open.value = true
}

defineExpose({ show })

function confirm() {
  showLoading.value = true
  api.decoration.saveBannerInfo(form.value).then(() => {
    open.value = false
    emit('update')
    $message.success('保存成功')
  }).catch(() => {
    $message.error('保存失败')
  }).finally(() => {
    showLoading.value = false
  })
}

function handleUploadSuccess(data) {
  if (data && data.url) {
    form.value.img_url = data.url
    $message.success('上传成功')
  }
}
</script>

<template>
  <NModal
    v-model:show="open"
    :title="form.id ? '编辑轮播' : '新增轮播'"
    preset="card"
    style="width: 800px;"
    :mask-closable="false"
  >
    <div class="edit-banner">
      <div class="image-section">
        <Upload
          type="image"
          :src="form.img_url"
          :cate-id="5"
          :show-drag-box="true"
          :preview-width="392"
          :preview-height="169"
          preview-fit="contain"
          :max-edge="1200"
          @success="handleUploadSuccess"
        />
      </div>

      <NForm class="settings-form form-section" :model="form" label-placement="top">
        <NFormItem label="标题">
          <NInput v-model:value="form.name" placeholder="请输入标题" />
        </NFormItem>

        <div class="form-switch-group">
          <NFormItem label="状态">
            <NSwitch
              v-model:value="form.status"
              :checked-value="1"
              :unchecked-value="2"
            >
              <template #checked>
                已启用
              </template>
              <template #unchecked>
                已停用
              </template>
            </NSwitch>
          </NFormItem>
        </div>

        <NFormItem label="类型">
          <NRadioGroup v-model:value="form.type" size="small">
            <NRadioButton :value="1">
              仅展示
            </NRadioButton>
            <NRadioButton :value="2">
              链接
            </NRadioButton>
            <NRadioButton :value="3">
              内容
            </NRadioButton>
          </NRadioGroup>
        </NFormItem>

        <NFormItem v-if="form.type === 2" label="链接">
          <NInput v-model:value="form.link_url" placeholder="https://xxx.com" />
        </NFormItem>
      </NForm>
    </div>

    <div v-if="form.type === 3" class="content-section">
      <WangEditor v-model:value="form.content" />
    </div>

    <template #footer>
      <div class="edit-banner__footer">
        <NButton secondary @click="open = false">
          取消
        </NButton>
        <NButton type="primary" :loading="showLoading" @click="confirm">
          确定
        </NButton>
      </div>
    </template>
  </NModal>
</template>

<style scoped lang="scss">
.edit-banner {
  display: flex;
  gap: 20px;
  align-items: flex-start;
}

.image-section {
  flex: 0 0 auto;
  width: 392px;
  height: 169px;
}

.form-section {
  flex: 1;
  min-width: 0;
  max-width: none;
}

.content-section {
  margin-top: 18px;
}

.edit-banner__footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style>
