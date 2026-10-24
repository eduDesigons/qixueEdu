<script setup>
import { NForm, NFormItem, NInput, NInputNumber, NModal, NSelect, NSwitch } from 'naive-ui'
import { computed, ref } from 'vue'
import api from '@/api'
import Upload from '@/components/common/Upload.vue'

const props = defineProps({
  type: { type: String, default: 'course' },
  parents: { type: Array, default: () => [] },
})
const emit = defineEmits(['update'])

const open = ref(false)
const loading = ref(false)
const parentLocked = ref(false)
const form = ref(emptyForm())

const parentOptions = computed(() => {
  const stayChild = form.value.id && form.value.parent_id && form.value.count > 0
  const options = stayChild ? [] : [{ label: '作为一级分类', value: 0 }]
  for (const item of props.parents) {
    if (item.id === form.value.id)
      continue
    options.push({ label: item.name, value: item.id })
  }
  return options
})

function emptyForm(parentId = 0) {
  return {
    id: null,
    name: '',
    parent_id: parentId,
    sort: 0,
    is_show: 1,
    icon: '',
    base_credit: 0,
    min_need_credit: 0,
    count: 0,
  }
}

function show(row = null, parentId = 0) {
  parentLocked.value = !row && parentId > 0
  if (row) {
    form.value = {
      id: row.id,
      name: row.name || '',
      parent_id: Number(row.parent_id) || 0,
      sort: Number(row.sort) || 0,
      is_show: Number(row.is_show) === 0 ? 0 : 1,
      icon: row.icon || '',
      base_credit: Number(row.base_credit) || 0,
      min_need_credit: Number(row.min_need_credit) || 0,
      count: Number(row.count) || 0,
    }
    const hasChildren = Array.isArray(row.children) && row.children.length > 0
    parentLocked.value = Number(row.parent_id) === 0 && (hasChildren || Number(row.count) > 0)
  }
  else {
    form.value = emptyForm(parentId)
  }
  open.value = true
}

defineExpose({ show })

const title = computed(() => {
  const kind = props.type === 'path' ? '系列分类' : '课程分类'
  return form.value.id ? `编辑${kind}` : `新增${kind}`
})

function confirm() {
  const name = String(form.value.name || '').trim()
  if (!name) {
    $message.warning('请填写分类名称')
    return
  }
  loading.value = true
  api.course.categorySave({
    type: props.type,
    id: form.value.id,
    name,
    parent_id: props.type === 'path' ? 0 : Number(form.value.parent_id) || 0,
    sort: Number(form.value.sort) || 0,
    is_show: Number(form.value.is_show) === 0 ? 0 : 1,
    icon: form.value.icon || '',
    base_credit: Math.max(0, Math.floor(Number(form.value.base_credit) || 0)),
    min_need_credit: Math.max(0, Math.floor(Number(form.value.min_need_credit) || 0)),
  }).then(() => {
    open.value = false
    emit('update')
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
      <NFormItem v-if="type === 'course'" label="上级">
        <NSelect
          v-model:value="form.parent_id"
          :options="parentOptions"
          :disabled="parentLocked"
        />
      </NFormItem>
      <NFormItem label="名称">
        <NInput v-model:value="form.name" maxlength="32" placeholder="请输入分类名称" />
      </NFormItem>
      <NFormItem label="排序">
        <NInputNumber v-model:value="form.sort" :min="0" class="w-full" />
      </NFormItem>
      <NFormItem v-if="type === 'course' && !form.parent_id" label="基础学分">
        <NInputNumber v-model:value="form.base_credit" :min="0" :precision="0" class="w-full" />
      </NFormItem>
      <NFormItem v-if="type === 'course' && form.parent_id" label="限制学分">
        <NInputNumber v-model:value="form.min_need_credit" :min="0" :precision="0" class="w-full" />
      </NFormItem>
      <NFormItem v-if="type === 'course'" label="显示">
        <NSwitch v-model:value="form.is_show" :checked-value="1" :unchecked-value="0" />
      </NFormItem>
      <NFormItem v-if="type === 'course' && !form.parent_id" label="首页图标">
        <Upload
          type="image"
          :src="form.icon"
          :show-drag-box="true"
          :preview-width="72"
          :preview-height="72"
          preview-fit="contain"
          @success="(data) => { form.icon = data?.url || '' }"
        />
      </NFormItem>
    </NForm>
    <template #footer>
      <div class="cate-modal__footer">
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
.cate-modal__footer {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
}

.w-full {
  width: 100%;
}
</style>
