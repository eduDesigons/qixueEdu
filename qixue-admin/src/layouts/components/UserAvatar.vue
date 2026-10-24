<!--------------------------------
 - @Author: Ronnie Zhang
 - @LastEditor: Ronnie Zhang
 - @LastEditTime: 2023/12/16 18:50:42
 - @Email: zclzone@outlook.com
 - Copyright © 2023 Ronnie Zhang(大脸怪) | https://isme.top
 --------------------------------->

<template>
  <n-dropdown :options="options" @select="handleSelect">
    <div id="user-dropdown" class="user-chip flex cursor-pointer items-center">
      <n-avatar round :size="26" :src="$assetUrl(userStore.avatar)" />
      <div v-if="userStore.userInfo" class="user-chip__meta ml-8 flex-col flex-shrink-0 items-start">
        <span class="user-chip__name">{{ userStore.nickName ?? userStore.username }}</span>
      </div>
      <i class="i-fe:chevron-down user-chip__caret ml-2" />
    </div>
  </n-dropdown>

  <RoleSelect ref="roleSelectRef" />
</template>

<script setup>
import api from '@/api'
import { RoleSelect } from '@/layouts/components'
import { useAuthStore, usePermissionStore, useUserStore } from '@/store'

const router = useRouter()
const userStore = useUserStore()
const authStore = useAuthStore()
const permissionStore = usePermissionStore()

const options = reactive([
  {
    label: '个人资料',
    key: 'profile',
    icon: () => h('i', { class: 'i-material-symbols:person-outline text-14' }),
    show: computed(() => permissionStore.accessRoutes?.some(item => item.path === '/profile')),
  },
  {
    label: '切换角色',
    key: 'toggleRole',
    icon: () => h('i', { class: 'i-basil:exchange-solid text-14' }),
    show: computed(() => userStore.roles.length > 1),
  },
  {
    label: '退出登录',
    key: 'logout',
    icon: () => h('i', { class: 'i-mdi:exit-to-app text-14' }),
  },
])

const roleSelectRef = ref(null)
function handleSelect(key) {
  switch (key) {
    case 'profile':
      router.push('/profile')
      break
    case 'toggleRole':
      roleSelectRef.value?.open({
        onOk() {
          location.reload()
        },
      })
      break
    case 'logout':
      $dialog.confirm({
        title: '提示',
        type: 'info',
        content: '确认退出？',
        async confirm() {
          try {
            await api.logout()
          }
          catch (error) {
            console.error(error)
          }
          authStore.logout()
          $message.success('已退出登录')
        },
      })
      break
  }
}
</script>

<style scoped>
.user-chip {
  gap: 0;
  padding: 3px 8px 3px 3px;
  border-radius: 999px;
  border: 1px solid transparent;
  transition:
    background 0.22s ease,
    border-color 0.22s ease,
    box-shadow 0.28s cubic-bezier(0.22, 1, 0.36, 1),
    transform 0.22s cubic-bezier(0.32, 0.72, 0, 1);
}

.user-chip:hover {
  background: var(--ctl-field, rgba(var(--primary-color), 0.06));
  border-color: var(--ctl-field-border, rgba(var(--primary-color), 0.12));
  box-shadow: inset 0 1px 0 rgb(255 255 255 / 55%);
}

.user-chip:active {
  transform: scale(0.97);
}

.user-chip__name {
  font-size: 12px;
  font-weight: 560;
  letter-spacing: -0.012em;
  color: var(--text-primary, #1a2332);
  line-height: 1.2;
  max-width: 84px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.user-chip__caret {
  font-size: 11px;
  color: var(--ctl-text-soft, #8b97a5);
  opacity: 0.8;
  transition: transform 0.22s cubic-bezier(0.32, 0.72, 0, 1);
}

.user-chip:hover .user-chip__caret {
  transform: translateY(1px);
}
</style>
