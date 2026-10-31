<template>
  <NModal
    v-model:show="visible"
    preset="card"
    title="发送站内信"
    class="w-860px"
    :mask-closable="false"
    @after-leave="resetForm"
  >
    <NForm
      ref="formRef"
      class="settings-form"
      :model="form"
      :rules="rules"
      label-placement="top"
    >
      <NFormItem v-if="userName" label="学员">
        <span class="msg-user-chip">{{ userName }}</span>
      </NFormItem>
      <NFormItem label="标题" path="title">
        <NInput
          ref="titleInput"
          v-model:value="form.title"
          maxlength="80"
          show-count
          placeholder="请输入标题"
          @keydown.enter.prevent
        />
      </NFormItem>
      <NFormItem label="内容" path="content">
        <WangEditor v-if="visible" v-model:value="form.content" />
      </NFormItem>
    </NForm>
    <template #footer>
      <NSpace justify="end">
        <NButton secondary @click="visible = false">取消</NButton>
        <NButton type="primary" :loading="submitting" @click="handleSubmit">发送</NButton>
      </NSpace>
    </template>
  </NModal>
</template>

<script setup>
import { NButton, NForm, NFormItem, NInput, NModal, NSpace } from 'naive-ui'
import api from '@/api'
import WangEditor from '@/components/common/WangEditor.vue'

defineOptions({ name: 'SendMessageModal' })

const emit = defineEmits(['success'])

const visible = ref(false)
const submitting = ref(false)
const formRef = ref(null)
const titleInput = ref(null)
const userName = ref('')
const extra = ref({
  userId: null,
  userIds: [],
  type: 'notice',
  relatedType: null,
  relatedId: null,
})
const form = ref({
  title: '',
  content: '',
})

function isEmptyHtml(html) {
  return !String(html || '')
    .replace(/<[^>]+>/g, '')
    .replace(/&nbsp;/gi, ' ')
    .replace(/\s+/g, '')
}

function toEditorHtml(text) {
  const value = String(text || '').trim()
  if (!value)
    return ''
  if (/<[a-z][\s\S]*>/i.test(value))
    return value
  return `<p>${value.replace(/\n/g, '<br>')}</p>`
}

const rules = {
  title: [{ required: true, message: '请输入标题', trigger: 'blur' }],
  content: [{
    required: true,
    trigger: ['blur', 'change'],
    validator() {
      if (isEmptyHtml(form.value.content))
        return new Error('请输入内容')
      return true
    },
  }],
}

function open(options = {}) {
  extra.value = {
    userId: options.userId || null,
    userIds: options.userIds || [],
    type: options.type || 'notice',
    relatedType: options.relatedType || null,
    relatedId: options.relatedId || null,
  }
  userName.value = options.userName || ''
  form.value = {
    title: options.title || '',
    content: toEditorHtml(options.content || ''),
  }
  visible.value = true
  nextTick(() => titleInput.value?.focus?.())
}

function resetForm() {
  form.value = { title: '', content: '' }
  userName.value = ''
  extra.value = {
    userId: null,
    userIds: [],
    type: 'notice',
    relatedType: null,
    relatedId: null,
  }
}

function handleSubmit() {
  if (submitting.value)
    return
  formRef.value?.validate((err) => {
    if (err) return
    const payload = {
      title: form.value.title.trim(),
      content: form.value.content.trim(),
      type: extra.value.type,
      related_type: extra.value.relatedType,
      related_id: extra.value.relatedId,
    }
    if (extra.value.userIds?.length) {
      payload.user_ids = extra.value.userIds
    } else {
      payload.user_id = extra.value.userId
    }
    if (!payload.user_id && !payload.user_ids?.length) {
      window.$message?.warning('请选择学员')
      return
    }
    submitting.value = true
    api.message.send(payload).then(() => {
      window.$message?.success('站内信已发送')
      visible.value = false
      emit('success')
    }).catch(() => {
    }).finally(() => {
      submitting.value = false
    })
  })
}

defineExpose({ open })
</script>

<style scoped>
.msg-user-chip {
  display: inline-flex;
  align-items: center;
  padding: 4px 12px;
  border-radius: 999px;
  font-size: 13px;
  font-weight: 560;
  letter-spacing: -0.01em;
  color: rgb(var(--primary-color));
  background: rgba(var(--primary-color), 0.1);
  border: 1px solid rgba(var(--primary-color), 0.14);
  box-shadow: inset 0 1px 0 rgb(255 255 255 / 40%);
}
</style>
