<script setup>
import { Add, Trash } from '@vicons/ionicons5'
import { cloneDeep } from 'lodash-es'
import { NButton, NIcon, NInput, NTable } from 'naive-ui'
import { ref, watch } from 'vue'
import api from '@/api'
import { SettingsEditorShell } from '@/components/index.js'

defineOptions({ name: 'EditFaqConfig' })

const props = defineProps({
  embedded: { type: Boolean, default: false },
})
const emit = defineEmits(['update', 'close'])

const open = ref(false)
const loading = ref(false)
const faqList = ref([])

function show() {
  loading.value = true
  api.config.getInfo('faqConfig').then((res) => {
    // 处理不同的数据结构
    let list = []
    if (res.data) {
      if (Array.isArray(res.data)) {
        list = res.data
      }
      else if (Array.isArray(res.data.list)) {
        list = res.data.list
      }
      else if (Array.isArray(res.data.data)) {
        list = res.data.data
      }
      else if (Array.isArray(res.data.data?.list)) {
        list = res.data.data.list
      }
    }
    // 确保每个项都有 title 和 answer 字段
    faqList.value = list.length > 0
      ? cloneDeep(list).map(item => ({
          title: item.title || '',
          answer: item.answer || '',
        }))
      : []
    if (!props.embedded)
      open.value = true
  }).catch((err) => {
    console.error('加载常见问题配置失败:', err)
    faqList.value = []
    if (!props.embedded)
      open.value = true
  }).finally(() => {
    loading.value = false
  })
}

watch(() => props.embedded, (v) => { if (v) show() }, { immediate: true })

defineExpose({ show })

function addFaq() {
  // 确保 faqList 是数组
  if (!Array.isArray(faqList.value)) {
    faqList.value = []
  }
  faqList.value.push({
    title: '',
    answer: '',
  })
}

function removeFaq(index) {
  faqList.value.splice(index, 1)
}

function confirm() {
  // 验证必填项
  for (let i = 0; i < faqList.value.length; i++) {
    const item = faqList.value[i]
    if (!item.title || !item.title.trim()) {
      window.$message.warning(`第${i + 1}条问题的标题不能为空`)
      return
    }
    if (!item.answer || !item.answer.trim()) {
      window.$message.warning(`第${i + 1}条问题的答案不能为空`)
      return
    }
  }

  loading.value = true
  api.config.saveInfo(faqList.value, 'faqConfig').then(() => {
    open.value = false
    emit('update')
    emit('close')
    window.$message.success('保存成功')
  }).catch(() => {
    window.$message.error('保存失败')
  }).finally(() => {
    loading.value = false
  })
}
</script>

<template>
  <SettingsEditorShell
    v-model:show="open"
    :embedded="embedded"
    title="常见问题配置"
    width="800px"
    :loading="loading"
    @save="confirm"
    @close="emit('close')"
  >
    <div class="editor-panel-head">
      <span>问题列表</span>
      <NButton type="primary" size="small" @click="addFaq">
        <template #icon>
          <NIcon :component="Add" />
        </template>
        添加问题
      </NButton>
    </div>

    <div v-if="faqList.length === 0" class="empty-table-placeholder">
      暂无问题配置，请点击「添加问题」
    </div>
    <div v-else class="faq-list-container">
      <NTable :bordered="false" :single-line="false">
          <thead>
            <tr>
              <th style="width: 80px;">
                序号
              </th>
              <th style="width: 250px;">
                问题标题
              </th>
              <th style="width: 350px;">
                问题答案
              </th>
              <th style="width: 100px;">
                操作
              </th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(item, index) in faqList" :key="index">
              <td style="padding: 12px; font-weight: 500;">
                {{ index + 1 }}
              </td>
              <td style="padding: 12px;">
                <NInput
                  v-model:value="item.title"
                  placeholder="请输入问题标题"
                />
              </td>
              <td style="padding: 12px;">
                <NInput
                  v-model:value="item.answer"
                  type="textarea"
                  :autosize="{ minRows: 2, maxRows: 4 }"
                  placeholder="请输入问题答案"
                />
              </td>
              <td style="padding: 12px;">
                <NButton
                  size="small"
                  type="error"
                  @click="removeFaq(index)"
                >
                  <template #icon>
                    <NIcon :component="Trash" />
                  </template>
                </NButton>
              </td>
            </tr>
          </tbody>
        </NTable>
    </div>
  </SettingsEditorShell>
</template>

<style scoped lang="scss">
.faq-list-container {
  max-height: 500px;
  overflow: auto;
  padding: 4px;
  border-radius: 14px;
  background: var(--ctl-field, rgba(var(--primary-color), 0.04));
  border: 1px solid var(--ctl-hairline, rgba(var(--primary-color), 0.08));
}

.faq-list-container :deep(.n-table) {
  background: transparent;
}
</style>

