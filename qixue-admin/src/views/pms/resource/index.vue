<!--------------------------------
 - @Author: Ronnie Zhang
 - @LastEditor: Ronnie Zhang
 - @LastEditTime: 2023/12/05 21:28:53
 - @Email: zclzone@outlook.com
 - Copyright © 2023 Ronnie Zhang(大脸怪) | https://isme.top
 --------------------------------->

<template>
  <CommonPage>
    <div class="resource-layout flex">
      <n-spin size="small" :show="treeLoading" class="resource-tree cus-scroll">
        <MenuTree
          v-model:current-menu="currentMenu"
          class="w-full shrink-0"
          :tree-data="treeData"
          @refresh="initData"
        />
      </n-spin>

      <div class="resource-detail cus-scroll w-0 flex-1">
        <template v-if="currentMenu">
          <div class="resource-detail__head">
            <h2 class="resource-detail__title">
              {{ currentMenu.name }}
            </h2>
            <NButton size="small" type="primary" secondary @click="handleEdit(currentMenu)">
              <i class="i-material-symbols:edit-outline mr-4 text-14" />
              编辑
            </NButton>
          </div>
          <SettingsList label="菜单信息">
            <SettingsListRow title="编码" :chevron="false">
              {{ currentMenu.code }}
            </SettingsListRow>
            <SettingsListRow title="名称" :chevron="false">
              {{ currentMenu.name }}
            </SettingsListRow>
            <SettingsListRow title="路由地址" :chevron="false">
              {{ currentMenu.path ?? '--' }}
            </SettingsListRow>
            <SettingsListRow title="组件路径" :chevron="false">
              {{ currentMenu.component ?? '--' }}
            </SettingsListRow>
            <SettingsListRow title="菜单图标" :chevron="false">
              <span v-if="currentMenu.icon" class="flex items-center">
                <i :class="`${currentMenu.icon}?mask text-22 mr-8`" />
                <span class="opacity-50">{{ currentMenu.icon }}</span>
              </span>
              <span v-else>无</span>
            </SettingsListRow>
            <SettingsListRow title="layout" :chevron="false">
              {{ currentMenu.layout || '跟随系统' }}
            </SettingsListRow>
            <SettingsListRow title="是否显示" :chevron="false">
              {{ currentMenu.show ? '是' : '否' }}
            </SettingsListRow>
            <SettingsListRow title="是否启用" :chevron="false">
              {{ currentMenu.enable ? '是' : '否' }}
            </SettingsListRow>
            <SettingsListRow title="KeepAlive" :chevron="false">
              {{ currentMenu.keepAlive ? '是' : '否' }}
            </SettingsListRow>
            <SettingsListRow title="排序" :chevron="false">
              {{ currentMenu.order ?? '--' }}
            </SettingsListRow>
          </SettingsList>

          <div class="editor-panel-head resource-btns-head">
            <span>按钮</span>
            <NButton size="small" type="primary" secondary @click="handleAddBtn">
              <i class="i-fe:plus mr-4 text-14" />
              新增
            </NButton>
          </div>

          <MeCrud
            ref="$table"
            :columns="btnsColumns"
            :scroll-x="-1"
            :get-data="api.getButtons"
            :query-items="{ parentId: currentMenu.id }"
          />
        </template>
        <div v-else class="resource-empty empty-state">
          请选择菜单查看详情
        </div>
      </div>
    </div>
    <ResAddOrEdit ref="modalRef" :menus="treeData" @refresh="initData" />
  </CommonPage>
</template>

<script setup>
import { NButton, NSwitch } from 'naive-ui'
import { MeCrud, SettingsList, SettingsListRow } from '@/components'
import { useRefreshOnActivate } from '@/composables'
import api from './api'
import MenuTree from './components/MenuTree.vue'
import ResAddOrEdit from './components/ResAddOrEdit.vue'

defineOptions({ name: 'Resource_Mgt' })

const treeData = ref([])
const treeLoading = ref(false)
const $table = ref(null)
const currentMenu = ref(null)
async function initData(data) {
  if (data?.type === 'BUTTON') {
    $table.value.handleSearch()
    return
  }
  treeLoading.value = true
  const res = await api.getMenuTree()
  treeData.value = res?.data || []
  treeLoading.value = false

  if (data)
    currentMenu.value = data
}
initData()

useRefreshOnActivate(() => initData())

const modalRef = ref(null)
function handleEdit(item = {}) {
  modalRef.value?.handleOpen({
    action: 'edit',
    title: `编辑菜单 - ${item.name}`,
    row: item,
    okText: '保存',
  })
}

const btnsColumns = [
  { title: '名称', key: 'name' },
  { title: '编码', key: 'code' },
  {
    title: '状态',
    key: 'enable',
    render: row =>
      h(
        NSwitch,
        {
          size: 'small',
          rubberBand: false,
          value: row.enable,
          checkedValue: 1,
          uncheckedValue: 0,
          loading: !!row.enableLoading,
          onUpdateValue: () => handleEnable(row),
        },
        {
          checked: () => '启用',
          unchecked: () => '停用',
        },
      ),
  },
  {
    title: '操作',
    key: 'actions',
    width: 320,
    align: 'right',
    fixed: 'right',
    render(row) {
      return [
        h(
          NButton,
          {
            size: 'small',
            type: 'primary',
            secondary: true,
            style: 'margin-left: 8px;',
            onClick: () => handleEditBtn(row),
          },
          {
            default: () => '编辑',
            icon: () => h('i', { class: 'i-material-symbols:edit-outline text-14' }),
          },
        ),

        h(
          NButton,
          {
            size: 'small',
            type: 'error',
            secondary: true,
            style: 'margin-left: 8px;',
            onClick: () => handleDeleteBtn(row.id),
          },
          {
            default: () => '删除',
            icon: () => h('i', { class: 'i-material-symbols:delete-outline text-14' }),
          },
        ),
      ]
    },
  },
]

watch(
  () => currentMenu.value,
  async (v) => {
    await nextTick()
    if (v)
      $table.value.handleSearch()
  },
)

function handleAddBtn() {
  modalRef.value?.handleOpen({
    action: 'add',
    title: '新增按钮',
    row: { type: 'BUTTON', parentId: currentMenu.value.id },
    okText: '保存',
  })
}

function handleEditBtn(row) {
  modalRef.value?.handleOpen({
    action: 'edit',
    title: `编辑按钮 - ${row.name}`,
    row,
    okText: '保存',
  })
}

function handleDeleteBtn(id) {
  const d = $dialog.warning({
    content: '确定删除？',
    title: '提示',
    positiveText: '确定',
    negativeText: '取消',
    async onPositiveClick() {
      try {
        d.loading = true
        await api.deletePermission(id)
        $message.success('删除成功')
        $table.value.handleSearch()
        d.loading = false
      }
      catch (error) {
        console.error(error)
        d.loading = false
      }
    },
  })
}

async function handleEnable(item) {
  try {
    item.enableLoading = true
    await api.savePermission({
      id: item.id,
      enable: !item.enable,
    })
    $message.success('操作成功')
    $table.value?.handleSearch()
    item.enableLoading = false
  }
  catch (error) {
    console.error(error)
    item.enableLoading = false
  }
}
</script>

<style scoped>
.resource-layout {
  gap: 16px;
  height: 100%;
  min-height: 0;
  overflow: hidden;
}

.resource-tree {
  width: 280px;
  flex-shrink: 0;
  padding: 8px 8px 12px !important;
  overflow: hidden !important;
  display: flex;
  flex-direction: column;
}

.resource-tree :deep(.n-spin-content) {
  height: 100%;
  min-height: 0;
}

.resource-detail {
  padding: 8px 8px 16px 4px !important;
  overflow: auto !important;
  min-height: 0;
}

.resource-detail__head {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 16px;
}

.resource-detail__title {
  margin: 0;
  font-family: var(--font-display);
  font-size: var(--title-lg, 32px);
  font-weight: 600;
  letter-spacing: -0.028em;
  line-height: 1.15;
  color: var(--ctl-text, #1a2332);
}

.resource-detail :deep(.settings-list) {
  max-width: none;
  margin: 0 0 20px;
}

.resource-btns-head {
  padding-left: 16px;
  margin-top: 4px;
}

.resource-empty {
  min-height: 280px;
}
</style>
