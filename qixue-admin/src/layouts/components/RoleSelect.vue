<!--------------------------------
 - @Author: Ronnie Zhang
 - @LastEditor: Ronnie Zhang
 - @LastEditTime: 2023/12/12 09:03:32
 - @Email: zclzone@outlook.com
 - Copyright © 2023 Ronnie Zhang(大脸怪) | https://isme.top
 --------------------------------->

<template>
  <MeModal ref="modalRef" title="请选择角色" width="380px" class="role-select-modal">
    <n-radio-group v-model:value="roleCode" class="role-select cus-scroll-y max-h-420 w-full">
      <n-space vertical :size="8" class="w-full">
        <n-radio-button
          v-for="role in roles"
          :key="role.id"
          class="role-select__option"
          :class="{ 'is-active': role.code === roleCode }"
          :value="role.code"
        >
          {{ role.name }}
        </n-radio-button>
      </n-space>
    </n-radio-group>

    <template #footer>
      <div class="flex gap-12">
        <n-button class="flex-1" size="large" secondary @click="logout()">
          退出登录
        </n-button>
        <n-button
          :loading="okLoading"
          class="flex-1"
          type="primary"
          size="large"
          :disabled="userStore.currentRole?.code === roleCode"
          @click="setCurrentRole"
        >
          确认
        </n-button>
      </div>
    </template>
  </MeModal>
</template>

<script setup>
import api from '@/api'
import { MeModal } from '@/components'
import { useModal } from '@/composables'
import { useAuthStore, useUserStore } from '@/store'

const userStore = useUserStore()
const authStore = useAuthStore()

const roles = ref(userStore.roles || [])
const roleCode = ref(userStore.currentRole?.code ?? roles.value[0]?.code ?? '')

const [modalRef, okLoading] = useModal()
function open(options) {
  modalRef.value?.open({
    ...options,
  })
}

async function setCurrentRole() {
  try {
    okLoading.value = true
    const { data } = await api.switchCurrentRole(roleCode.value)
    await authStore.switchCurrentRole(data)
    okLoading.value = false
    $message.success('切换成功')
    modalRef.value?.handleOk()
  }
  catch (error) {
    console.error(error)
    okLoading.value = false
    return false
  }
}

async function logout() {
  await api.logout()
  authStore.logout()
  modalRef.value?.close()
  $message.success('已退出登录')
}

defineExpose({
  open,
})
</script>

<style scoped>
.role-select {
  padding: 8px 0 4px;
}

.role-select :deep(.n-radio-group) {
  width: 100%;
}

.role-select__option {
  width: 100% !important;
  height: 44px !important;
  border-radius: 14px !important;
  border: 1px solid var(--ctl-field-border, rgba(var(--primary-color), 0.12)) !important;
  background: var(--ctl-field, rgba(var(--primary-color), 0.05)) !important;
  color: #1a2332 !important;
  font-size: 15px !important;
  font-weight: 540 !important;
  letter-spacing: -0.012em;
  line-height: 44px !important;
  text-align: center;
  box-shadow: inset 0 1px 0 rgb(255 255 255 / 50%);
  transition:
    transform 0.22s cubic-bezier(0.32, 0.72, 0, 1),
    background 0.22s ease,
    box-shadow 0.28s cubic-bezier(0.22, 1, 0.36, 1),
    border-color 0.22s ease,
    color 0.22s ease !important;
}

.role-select__option:hover {
  background: var(--ctl-field-hover, rgba(var(--primary-color), 0.08)) !important;
}

.role-select__option:active {
  transform: scale(0.985);
}

.role-select__option.is-active {
  background: rgba(var(--primary-color), 0.12) !important;
  border-color: rgba(var(--primary-color), 0.28) !important;
  color: rgb(var(--primary-color)) !important;
  font-weight: 600 !important;
  box-shadow: inset 0 1px 0 rgb(255 255 255 / 55%) !important;
}
</style>
