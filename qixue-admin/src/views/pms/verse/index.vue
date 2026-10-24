<script setup>
import { Add } from '@vicons/ionicons5'
import { NButton, NForm, NFormItem, NIcon, NInput, NModal } from 'naive-ui'
import { h, onMounted, ref } from 'vue'
import api from '@/api'
import { CommonPage, MeCrud, MeQueryItem } from '@/components'
import { useRefreshOnActivate } from '@/composables'

defineOptions({ name: 'Verse' })

const $table = ref(null)
const queryItems = ref({ keyword: '' })
const open = ref(false)
const saving = ref(false)
const form = ref({ id: null, content: '' })

onMounted(() => $table.value?.handleSearch())
useRefreshOnActivate(() => $table.value?.handleSearch())

function showEditor(row) {
  form.value = {
    id: row?.id || null,
    content: row?.content || '',
  }
  open.value = true
}

function save() {
  const content = String(form.value.content || '').trim()
  if (!content) {
    $message.warning('请填写一言')
    return
  }
  saving.value = true
  api.verse.save({ id: form.value.id, content }).then(() => {
    $message.success('保存成功')
    open.value = false
    $table.value?.handleSearch()
  }).catch(() => {}).finally(() => {
    saving.value = false
  })
}

function removeRow(row) {
  $dialog.warning({
    title: '确认删除',
    content: '确定删除这条一言吗？',
    positiveText: '删除',
    negativeText: '取消',
    onPositiveClick: () => api.verse.delete({ id: row.id }).then(() => {
      $message.success('删除成功')
      $table.value?.handleSearch()
    }).catch(() => {}),
  })
}

const columns = [
  {
    title: '一言',
    key: 'content',
    ellipsis: { tooltip: true },
  },
  {
    title: '操作',
    key: 'actions',
    width: 140,
    align: 'center',
    render: row => h('div', { style: 'display: flex; gap: 8px; justify-content: center;' }, [
      h(NButton, { size: 'small', secondary: true, onClick: () => showEditor(row) }, { default: () => '编辑' }),
      h(NButton, { size: 'small', secondary: true, type: 'error', onClick: () => removeRow(row) }, { default: () => '删除' }),
    ]),
  },
]
</script>

<template>
  <CommonPage title="一言" subtitle="首页每次随机展示一条。">
    <template #action>
      <NButton type="primary" @click="showEditor()">
        <template #icon>
          <NIcon :component="Add" />
        </template>
        新增一言
      </NButton>
    </template>
    <MeCrud
      ref="$table"
      v-model:query-items="queryItems"
      collapse-query
      :columns="columns"
      :get-data="api.verse.list"
    >
      <MeQueryItem label="内容" :label-width="50">
        <NInput v-model:value="queryItems.keyword" clearable placeholder="搜索一言" class="w-220" />
      </MeQueryItem>
    </MeCrud>
    <NModal
      v-model:show="open"
      :title="form.id ? '编辑一言' : '新增一言'"
      preset="card"
      style="width: 560px;"
      :mask-closable="false"
    >
      <NForm label-placement="top">
        <NFormItem label="内容">
          <NInput
            v-model:value="form.content"
            type="textarea"
            maxlength="500"
            show-count
            :autosize="{ minRows: 4, maxRows: 8 }"
            placeholder="请输入一言"
          />
        </NFormItem>
      </NForm>
      <template #footer>
        <div class="verse-modal__footer">
          <NButton @click="open = false">
            取消
          </NButton>
          <NButton type="primary" :loading="saving" @click="save">
            保存
          </NButton>
        </div>
      </template>
    </NModal>
  </CommonPage>
</template>

<style scoped>
.verse-modal__footer {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
}
</style>
