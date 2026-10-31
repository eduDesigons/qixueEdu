<script setup>
import { Add, Trash } from '@vicons/ionicons5'
import { cloneDeep } from 'lodash-es'
import { NButton, NIcon, NInput } from 'naive-ui'
import { ref, watch } from 'vue'
import api from '@/api'
import { SettingsEditorShell } from '@/components/index.js'
import Upload from '@/components/common/Upload.vue'

defineOptions({ name: 'EditToolConfig' })

const props = defineProps({
  embedded: { type: Boolean, default: false },
})
const emit = defineEmits(['update', 'close'])

const open = ref(false)
const loading = ref(false)
const toolList = ref([])

function normalizeList(raw) {
  let list = []
  if (Array.isArray(raw))
    list = raw
  else if (raw && Array.isArray(raw.list))
    list = raw.list
  else if (raw && Array.isArray(raw.data))
    list = raw.data
  return list.map(item => ({
    title: item?.title || '',
    desc: item?.desc || '',
    url: item?.url || '',
    icon: item?.icon || '',
  }))
}

function show() {
  loading.value = true
  api.config.getInfo('toolConfig').then((res) => {
    toolList.value = normalizeList(res.data)
    if (!props.embedded)
      open.value = true
  }).catch(() => {
    toolList.value = []
    if (!props.embedded)
      open.value = true
  }).finally(() => {
    loading.value = false
  })
}

watch(() => props.embedded, (v) => { if (v) show() }, { immediate: true })

defineExpose({ show })

function addTool() {
  toolList.value.push({ title: '', desc: '', url: '', icon: '' })
}

function removeTool(index) {
  toolList.value.splice(index, 1)
}

function confirm() {
  const list = cloneDeep(toolList.value)
  for (let i = 0; i < list.length; i++) {
    const item = list[i]
    const n = i + 1
    item.title = String(item.title || '').trim()
    item.desc = String(item.desc || '').trim()
    item.url = String(item.url || '').trim()
    item.icon = String(item.icon || '').trim()
    if (!item.title) {
      window.$message.warning(`第${n}个工具的标题不能为空`)
      return
    }
    if (!item.desc) {
      window.$message.warning(`第${n}个工具的简介不能为空`)
      return
    }
    if (!/^https?:\/\//i.test(item.url)) {
      window.$message.warning(`第${n}个工具的网址需以 http:// 或 https:// 开头`)
      return
    }
    if (!item.icon) {
      window.$message.warning(`第${n}个工具请上传图标`)
      return
    }
  }
  loading.value = true
  api.config.saveInfo(list, 'toolConfig').then(() => {
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
    title="工具库配置"
    width="720px"
    :loading="loading"
    @save="confirm"
    @close="emit('close')"
  >
    <div class="editor-panel-head">
      <span>工具列表</span>
      <NButton type="primary" size="small" @click="addTool">
        <template #icon>
          <NIcon :component="Add" />
        </template>
        添加工具
      </NButton>
    </div>

    <div v-if="toolList.length === 0" class="empty-table-placeholder">
      暂无工具，请点击「添加工具」
    </div>
    <div v-else class="tool-list">
      <div v-for="(item, index) in toolList" :key="index" class="tool-card">
        <div class="tool-card__head">
          <span>工具 {{ index + 1 }}</span>
          <NButton size="small" type="error" secondary @click="removeTool(index)">
            <template #icon>
              <NIcon :component="Trash" />
            </template>
            删除
          </NButton>
        </div>
        <div class="tool-card__body">
          <Upload
            type="image"
            :src="item.icon || ''"
            :cate-id="1"
            :show-drag-box="true"
            :preview-width="120"
            :preview-height="120"
            preview-fit="cover"
            :max-edge="256"
            @success="(data) => item.icon = data.url"
          />
          <div class="tool-card__fields">
            <NInput v-model:value="item.title" maxlength="40" placeholder="工具标题" />
            <NInput v-model:value="item.desc" maxlength="80" placeholder="简介" />
            <NInput v-model:value="item.url" placeholder="外链网址，https:// 开头" />
          </div>
        </div>
      </div>
    </div>
  </SettingsEditorShell>
</template>

<style scoped lang="scss">
.tool-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
  max-height: 560px;
  overflow: auto;
}

.tool-card {
  padding: 14px;
  border-radius: 14px;
  background: var(--ctl-field, rgba(var(--primary-color), 0.04));
  border: 1px solid var(--ctl-hairline, rgba(var(--primary-color), 0.08));
}

.tool-card__head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;
  font-weight: 600;
}

.tool-card__body {
  display: flex;
  gap: 14px;
  align-items: flex-start;
}

.tool-card__fields {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 8px;
}
</style>
