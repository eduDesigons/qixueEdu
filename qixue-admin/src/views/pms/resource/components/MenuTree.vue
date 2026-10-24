<!--------------------------------
 - @Author: Ronnie Zhang
 - @LastEditor: Ronnie Zhang
 - @LastEditTime: 2024/04/01 15:51:34
 - @Email: zclzone@outlook.com
 - Copyright © 2023 Ronnie Zhang(大脸怪) | https://isme.top
 --------------------------------->

<template>
  <div class="menu-tree">
    <div class="menu-tree__head">
      <h3 class="menu-tree__title">菜单</h3>
      <NButton size="small" type="primary" secondary @click="handleAdd()">
        <i class="i-material-symbols:add mr-4 text-14" />
        新增
      </NButton>
    </div>
    <n-input
      v-model:value="pattern"
      placeholder="搜索"
      clearable
      class="menu-tree__search"
    />
    <div class="menu-tree__list">
      <n-tree
        :show-irrelevant-nodes="false"
        :pattern="pattern"
        :data="treeData"
        :selected-keys="[currentMenu?.code]"
        :render-prefix="renderPrefix"
        :render-suffix="renderSuffix"
        :on-update:selected-keys="onSelect"
        key-field="code"
        label-field="name"
        block-line
        default-expand-all
      />
    </div>

    <ResAddOrEdit ref="modalRef" :menus="treeData" @refresh="(data) => emit('refresh', data)" />
  </div>
</template>

<script setup>
import { NButton } from 'naive-ui'
import { withModifiers } from 'vue'
import api from '../api'
import ResAddOrEdit from './ResAddOrEdit.vue'

defineProps({
  treeData: {
    type: Array,
    default: () => [],
  },
  currentMenu: {
    type: Object,
    default: () => null,
  },
})
const emit = defineEmits(['refresh', 'update:currentMenu'])

const pattern = ref('')

const modalRef = ref(null)
async function handleAdd(data = {}) {
  modalRef.value?.handleOpen({
    action: 'add',
    title: '新增菜单',
    row: { type: 'MENU', ...data },
    okText: '保存',
  })
}

function onSelect(keys, option, { action, node }) {
  emit('update:currentMenu', action === 'select' ? node : null)
}

function renderPrefix({ option }) {
  return h('i', { class: `${option.icon}?mask text-16` })
}

function renderSuffix({ option }) {
  return [
    h(
      NButton,
      {
        text: true,
        type: 'primary',
        title: '新增下级菜单',
        size: 'tiny',
        onClick: withModifiers(() => handleAdd({ parentId: option.id }), ['stop']),
      },
      { default: () => '新增' },
    ),

    h(
      NButton,
      {
        text: true,
        type: 'error',
        size: 'tiny',
        style: 'margin-left: 12px;',
        onClick: withModifiers(() => handleDelete(option), ['stop']),
      },
      { default: () => '删除' },
    ),
  ]
}

function handleDelete(item) {
  $dialog.confirm({
    content: `确认删除【${item.name}】？`,
    async confirm() {
      try {
        $message.loading('正在删除', { key: 'deleteMenu' })
        await api.deletePermission(item.id)
        $message.success('删除成功', { key: 'deleteMenu' })
        emit('refresh')
        emit('update:currentMenu', null)
      }
      catch (error) {
        console.error(error)
        $message.destroy('deleteMenu')
      }
    },
  })
}
</script>

<style scoped>
.menu-tree {
  display: flex;
  flex-direction: column;
  gap: 10px;
  min-height: 0;
  height: 100%;
}

.menu-tree__head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  padding: 0 4px 0 8px;
}

.menu-tree__title {
  margin: 0;
  font-family: var(--font-display);
  font-size: 13px;
  font-weight: 560;
  letter-spacing: -0.01em;
  color: var(--ctl-text-soft, #5b6b7c);
}

.menu-tree__search {
  --n-height: var(--ctl-h-lg, 28px) !important;
  --n-border-radius: 10px !important;
}

.menu-tree__list {
  flex: 1;
  min-height: 0;
  overflow: auto;
  padding: 4px;
  border-radius: 14px;
  background: var(--ctl-field, rgba(var(--primary-color), 0.04));
  border: 1px solid var(--ctl-hairline, rgba(var(--primary-color), 0.08));
}

.menu-tree :deep(.n-tree-node-content) {
  border-radius: 10px !important;
  min-height: 34px;
  letter-spacing: -0.01em;
  transition:
    background 0.2s ease,
    transform 0.2s cubic-bezier(0.32, 0.72, 0, 1) !important;
}

.menu-tree :deep(.n-tree-node-content:hover) {
  background: rgb(255 255 255 / 55%) !important;
}

.menu-tree :deep(.n-tree-node-content:active) {
  transform: scale(0.985);
}

.menu-tree :deep(.n-tree-node--selected .n-tree-node-content) {
  background: rgb(255 255 255 / 78%) !important;
  box-shadow: none;
  border: 0;
  color: var(--ctl-text, #1a2332) !important;
  font-weight: 600;
}
</style>
