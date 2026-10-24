<script setup>
import { NForm, NFormItem, NInput, NInputNumber, NModal, NSelect, NSwitch } from 'naive-ui'
import { computed, ref } from 'vue'
import api from '@/api'

const props = defineProps({
  parents: { type: Array, default: () => [] },
})
const emit = defineEmits(['update'])

const open = ref(false)
const loading = ref(false)
const parentLocked = ref(false)
const form = ref(emptyForm())

const parentOptions = computed(() => {
  const options = [{ label: '作为一级区域', value: 0 }]
  for (const item of props.parents) {
    if (item.id === form.value.id)
      continue
    options.push({ label: item.name, value: item.id })
  }
  return options
})

function emptyForm(parentId = 0, code = '') {
  return {
    id: null,
    name: '',
    parent_id: parentId,
    code,
    sort: 0,
    is_show: 1,
  }
}

function show(row = null, parentId = 0) {
  parentLocked.value = !row && parentId > 0
  if (row) {
    const hasChildren = Array.isArray(row.children) && row.children.length > 0
    parentLocked.value = Number(row.parent_id) === 0 && hasChildren
    form.value = {
      id: row.id,
      name: row.name || '',
      parent_id: Number(row.parent_id) || 0,
      code: row.code || '',
      sort: Number(row.sort) || 0,
      is_show: Number(row.is_show) === 0 ? 0 : 1,
    }
  }
  else {
    const parent = props.parents.find(item => item.id === parentId)
    form.value = emptyForm(parentId, parent?.code || '')
  }
  open.value = true
}

defineExpose({ show })

const title = computed(() => (form.value.id ? '编辑区域' : '新增区域'))

function confirm() {
  const name = String(form.value.name || '').trim()
  const parentId = Number(form.value.parent_id) || 0
  if (!name) {
    $message.warning('请填写名称')
    return
  }
  let code = ''
  if (!parentId) {
    code = String(form.value.code || '').trim().toLowerCase()
    if (!/^[a-z0-9]{1,30}$/.test(code)) {
      $message.warning('编码只能是字母和数字')
      return
    }
  }
  loading.value = true
  api.area.save({
    id: form.value.id,
    name,
    parent_id: parentId,
    code,
    sort: Number(form.value.sort) || 0,
    is_show: Number(form.value.is_show) === 0 ? 0 : 1,
  }).then((res) => {
    open.value = false
    emit('update', { id: Number(res?.data?.id) || form.value.id, parent_id: parentId })
    $message.success('保存成功')
  }).catch(() => {}).finally(() => {
    loading.value = false
  })
}
</script>

<template>
  <NModal
    v-model:show="open"
    :title="title"
    preset="card"
    style="width: 520px;"
    :mask-closable="false"
  >
    <NForm class="settings-form" label-placement="top">
      <NFormItem label="上级">
        <NSelect
          v-model:value="form.parent_id"
          :options="parentOptions"
          :disabled="parentLocked"
        />
      </NFormItem>
      <NFormItem label="名称">
        <NInput v-model:value="form.name" maxlength="32" placeholder="请输入区域名称" />
      </NFormItem>
      <NFormItem v-if="!form.parent_id" label="编码">
        <NInput v-model:value="form.code" maxlength="30" placeholder="字母或数字，用于登录域名" />
      </NFormItem>
      <NFormItem label="排序">
        <NInputNumber v-model:value="form.sort" :min="0" class="w-full" />
      </NFormItem>
      <NFormItem label="显示">
        <NSwitch v-model:value="form.is_show" :checked-value="1" :unchecked-value="0" />
      </NFormItem>
    </NForm>
    <template #footer>
      <div class="area-modal__footer">
        <n-button @click="open = false">
          取消
        </n-button>
        <n-button type="primary" :loading="loading" @click="confirm">
          保存
        </n-button>
      </div>
    </template>
  </NModal>
</template>

<style scoped>
.area-modal__footer {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
}

.w-full {
  width: 100%;
}
</style>
