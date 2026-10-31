<script setup>
import { cloneDeep } from 'lodash-es'
import { NButton, NDatePicker, NForm, NFormItem, NInput, NModal, NSwitch } from 'naive-ui'
import { ref } from 'vue'
import api from '@/api'
import WangEditor from '@/components/common/WangEditor.vue'

const emit = defineEmits(['update'])

const open = ref(false)
const showLoading = ref(false)

const form = ref({
  id: null,
  title: '',
  publish_time: null,
  content: '',
  status: 1,
})

function show(data = null) {
  if (data) {
    form.value = cloneDeep(data)
    if (form.value.publish_time) {
      form.value.publish_time = new Date(form.value.publish_time)
    }
  }
  else {
    form.value = {
      id: null,
      title: '',
      publish_time: new Date(),
      content: '',
      status: 1,
    }
  }
  open.value = true
}

defineExpose({ show })

function confirm() {
  if (!form.value.title) {
    $message.warning('请输入标题')
    return
  }
  if (!form.value.publish_time) {
    $message.warning('请选择发布时间')
    return
  }

  showLoading.value = true
  const submitData = {
    ...form.value,
    publish_time: form.value.publish_time ? new Date(form.value.publish_time).getTime() / 1000 : null,
  }
  api.decoration.saveNotice(submitData).then(() => {
    open.value = false
    emit('update')
    $message.success('保存成功')
  }).catch(() => {
    $message.error('保存失败')
  }).finally(() => {
    showLoading.value = false
  })
}
</script>

<template>
  <NModal
    v-model:show="open"
    :title="form.id ? '编辑公告' : '新增公告'"
    preset="card"
    style="width: 900px;"
    :mask-closable="false"
  >
    <NForm class="settings-form edit-notice" :model="form" label-placement="top">
      <NFormItem label="标题">
        <NInput v-model:value="form.title" placeholder="请输入标题" />
      </NFormItem>

      <NFormItem label="发布时间">
        <NDatePicker
          v-model:value="form.publish_time"
          type="datetime"
          placeholder="请选择发布时间"
          style="width: 100%"
        />
      </NFormItem>

      <div class="form-switch-group">
        <NFormItem label="是否发布">
          <NSwitch
            v-model:value="form.status"
            :checked-value="1"
            :unchecked-value="0"
          >
            <template #checked>
              已发布
            </template>
            <template #unchecked>
              未发布
            </template>
          </NSwitch>
        </NFormItem>
      </div>

      <NFormItem label="内容" class="edit-notice__content">
        <WangEditor v-model:value="form.content" />
      </NFormItem>
    </NForm>

    <template #footer>
      <div class="edit-notice__footer">
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
.edit-notice {
  max-width: none;
}

.edit-notice__content :deep(.n-form-item-blank) {
  width: 100%;
}

.edit-notice__footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}
</style>
