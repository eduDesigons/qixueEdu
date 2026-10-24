<!--------------------------------
 - @Author: Ronnie Zhang
 - @LastEditor: Ronnie Zhang
 - @LastEditTime: 2023/12/05 21:29:56
 - @Email: zclzone@outlook.com
 - Copyright © 2023 Ronnie Zhang(大脸怪) | https://isme.top
 --------------------------------->

<template>
  <CommonPage>
    <template #action>
      <NButton v-permission="'AddUser'" type="primary" @click="handleAdd()">
        <i class="i-material-symbols:add mr-4 text-18" />
        创建新用户
      </NButton>
    </template>

    <MeCrud
      ref="$table"
      v-model:query-items="queryItems"
      collapse-query
      :scroll-x="1000"
      :columns="columns"
      :get-data="api.read"
    >
      <CityQueryItem v-model="queryItems.city_id" :label-width="50" />
      <MeQueryItem label="用户名" :label-width="50">
        <n-input
          v-model:value="queryItems.username"
          type="text"
          placeholder="请输入用户名"
          clearable
        />
      </MeQueryItem>

      <MeQueryItem label="状态" :label-width="50">
        <n-select
          v-model:value="queryItems.enable"
          clearable
          :options="[
            { label: '启用', value: 1 },
            { label: '停用', value: 0 },
          ]"
        />
      </MeQueryItem>
    </MeCrud>

    <MeModal ref="modalRef" width="520px">
      <n-form
        ref="modalFormRef"
        class="settings-form"
        label-placement="top"
        :model="modalForm"
        :disabled="modalAction === 'view'"
      >
        <n-form-item
          label="用户名"
          path="username"
          :rule="{
            required: true,
            message: '请输入用户名',
            trigger: ['input', 'blur'],
          }"
        >
          <n-input v-model:value="modalForm.username" :disabled="modalAction !== 'add'" />
        </n-form-item>
        <n-form-item
          v-if="['add', 'reset'].includes(modalAction)"
          :label="modalAction === 'reset' ? '重置密码' : '初始密码'"
          path="password"
          :rule="{
            required: true,
            message: '请输入密码',
            trigger: ['input', 'blur'],
          }"
        >
          <n-input v-model:value="modalForm.password" type="password" show-password-on="mousedown" />
        </n-form-item>

        <n-form-item v-if="['add', 'setRole'].includes(modalAction)" label="角色" path="roleIds">
          <n-select
            v-model:value="modalForm.roleIds"
            :options="roles"
            label-field="name"
            value-field="id"
            clearable
            filterable
            multiple
          />
        </n-form-item>
        <n-form-item v-if="canAssignCity && ['add', 'setCity'].includes(modalAction)" label="大区" path="city_id">
          <n-select
            v-model:value="modalForm.city_id"
            :options="cityOptions"
            clearable
            filterable
            placeholder="不选则不限"
          />
        </n-form-item>
        <n-form-item v-if="modalAction === 'add'" label="状态" path="enable">
          <NSwitch v-model:value="modalForm.enable" :unchecked-value="0" :checked-value="1">
            <template #checked>
              启用
            </template>
            <template #unchecked>
              停用
            </template>
          </NSwitch>
        </n-form-item>
      </n-form>
      <n-alert v-if="modalAction === 'add'" type="warning" closable>
        详细信息需由用户本人补充修改
      </n-alert>
    </MeModal>
  </CommonPage>
</template>

<script setup>
import { NAvatar, NButton, NSwitch, NTag } from 'naive-ui'
import { MeCrud, MeModal, MeQueryItem, CityQueryItem } from '@/components'
import { useCrud, useRefreshOnActivate } from '@/composables'
import { useUserStore } from '@/store'
import { getAssetUrl, formatDateTime } from '@/utils'
import api from './api'
import commonApi from '@/api/common'

defineOptions({ name: 'UserMgt' })

const $table = ref(null)
const userStore = useUserStore()
const canAssignCity = computed(() => !Number(userStore.userInfo?.city_id || 0))
/** QueryBar筛选参数（可选） */
const queryItems = ref({
  username: '',
  enable: null,
  city_id: null,
})

onMounted(() => {
  $table.value?.handleSearch()
})

useRefreshOnActivate(() => $table.value?.handleSearch())

const roles = ref([])
api.getAllRoles().then(({ data = [] }) => (roles.value = data))

const cityOptions = ref([])
commonApi.getAreaTree().then(({ data = [] }) => {
  cityOptions.value = (data || []).map(item => ({
    label: item.name,
    value: Number(item.id),
  }))
})

const {
  modalRef,
  modalFormRef,
  modalForm,
  modalAction,
  handleAdd,
  handleDelete,
  handleOpen,
  handleSave,
} = useCrud({
  name: '用户',
  initForm: { enable: 1, city_id: null },
  doCreate: api.create,
  doDelete: api.delete,
  doUpdate: api.update,
  refresh: () => $table.value?.handleSearch(),
})

const columns = [
  {
    title: '头像',
    key: 'avatar',
    width: 80,
    render: ({ avatar }) =>
      h(NAvatar, {
        size: 'medium',
        round: true,
        src: getAssetUrl(avatar) || undefined,
      }),
  },
  { title: '用户名', key: 'username', width: 150, ellipsis: { tooltip: true } },
  {
    title: '角色',
    key: 'roles',
    width: 200,
    ellipsis: { tooltip: true },
    render: ({ roles }) => {
      if (roles?.length) {
        return roles.map((item, index) =>
          h(
            NTag,
            {
              type: 'success',
              size: 'small',
              round: true,
              bordered: false,
              style: index > 0 ? 'margin-left: 8px;' : '',
            },
            { default: () => item.name },
          ),
        )
      }
      return h('span', { class: 'text-gray-400' }, '暂无角色')
    },
  },
  {
    title: '大区',
    key: 'city_id',
    width: 120,
    render: ({ city, city_id }) => city?.name || (city_id ? String(city_id) : h('span', { class: 'text-gray-400' }, '不限')),
  },
  {
    title: '创建时间',
    key: 'createDate',
    width: 180,
    render(row) {
      return h('span', formatDateTime(row.createTime))
    },
  },
  {
    title: '状态',
    key: 'enable',
    width: 120,
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
        width: 520,
    align: 'right',
    fixed: 'right',
    hideInExcel: true,
    render(row) {
      return [
        // withPermission(
        //   h(NButton, {
        //     size: 'small',
        //     type: 'primary',
        //     secondary: true,
        //   }, {
        //     default: () => '超管专属',
        //     icon: () => h('i', { class: 'i-carbon:user-role text-14' }),
        //   }),
        //   'SuperAdmin',
        // ),
        h(
          NButton,
          {
            size: 'small',
            type: 'primary',
            class: 'ml-12px',
            secondary: true,
            onClick: () => handleOpenRolesSet(row),
          },
          {
            default: () => '分配角色',
            icon: () => h('i', { class: 'i-carbon:user-role text-14' }),
          },
        ),
        canAssignCity.value && h(
          NButton,
          {
            size: 'small',
            type: 'primary',
            class: 'ml-12px',
            secondary: true,
            onClick: () => handleOpenCitySet(row),
          },
          {
            default: () => '分配大区',
            icon: () => h('i', { class: 'i-carbon:location text-14' }),
          },
        ),
        h(
          NButton,
          {
            size: 'small',
            type: 'primary',
            secondary: true,
            style: 'margin-left: 12px;',
            onClick: () => handleOpen({ action: 'reset', title: '重置密码', row, onOk: onSave }),
          },
          {
            default: () => '重置密码',
            icon: () => h('i', { class: 'i-radix-icons:reset text-14' }),
          },
        ),

        h(
          NButton,
          {
            size: 'small',
            type: 'error',
            secondary: true,
            style: 'margin-left: 12px;',
            onClick: () => handleDelete(row.id),
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

async function handleEnable(row) {
  row.enableLoading = true
  try {
    await api.update({
      id: row.id,
      username: row.username,
      enable: row.enable ? 0 : 1,
    })
    row.enableLoading = false
    $message.success('操作成功')
    $table.value?.handleSearch()
  }
  catch (error) {
    console.error(error)
    row.enableLoading = false
  }
}

function handleOpenRolesSet(row) {
  const roleIds = row.roles.map(item => item.id)
  handleOpen({
    action: 'setRole',
    title: '分配角色',
    row: { id: row.id, username: row.username, roleIds },
    onOk: onSave,
  })
}

function handleOpenCitySet(row) {
  handleOpen({
    action: 'setCity',
    title: '分配大区',
    row: { id: row.id, username: row.username, city_id: row.city_id == null ? null : Number(row.city_id) },
    onOk: onSave,
  })
}

function onSave() {
  if (modalAction.value === 'setRole') {
    return handleSave({
      api: () => api.update(modalForm.value),
      cb: () => $message.success('分配成功'),
    })
  }
  else if (modalAction.value === 'setCity') {
    return handleSave({
      api: () => api.update({
        id: modalForm.value.id,
        username: modalForm.value.username,
        city_id: modalForm.value.city_id ?? '',
      }),
      cb: () => $message.success('大区已更新'),
    })
  }
  else if (modalAction.value === 'reset') {
    return handleSave({
      api: () => api.resetPwd({
        id: modalForm.value.id,
        password: modalForm.value.password,
      }),
      cb: () => $message.success('密码重置成功'),
    })
  }
  handleSave()
}
</script>
