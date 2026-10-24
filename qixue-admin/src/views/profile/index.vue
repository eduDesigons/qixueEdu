<!--------------------------------
 - @Author: Ronnie Zhang
 - @LastEditor: Ronnie Zhang
 - @LastEditTime: 2023/12/05 21:30:11
 - @Email: zclzone@outlook.com
 - Copyright © 2023 Ronnie Zhang(大脸怪) | https://isme.top
 --------------------------------->

<template>
  <AppPage show-footer>
    <div class="profile">
      <header class="page-mast">
        <h1 class="page-mast__title">个人资料</h1>
        <p class="page-mast__sub">账号信息与安全设置</p>
      </header>
      <section class="profile-hero">
        <n-avatar round :size="88" :src="$assetUrl(userStore.avatar)" class="profile-hero__avatar" />
        <div class="profile-hero__meta">
          <h2>{{ userStore.nickName || userStore.username }}</h2>
          <p>@{{ userStore.username }}</p>
          <div class="profile-hero__actions">
            <n-button secondary @click="openAvatarModal()">
              更改头像
            </n-button>
            <n-button type="primary" secondary @click="pwdModalRef.open()">
              <i class="i-fe:edit mr-4" />
              修改密码
            </n-button>
          </div>
        </div>
      </section>

      <SettingsList label="个人资料">
        <template #extra>
          <n-button type="primary" secondary size="small" @click="openProfileModal()">
            <i class="i-fe:edit mr-4" />
            编辑
          </n-button>
        </template>
        <SettingsListRow title="昵称" :chevron="false">
          {{ userStore.nickName || '—' }}
        </SettingsListRow>
        <SettingsListRow title="性别" :chevron="false">
          {{ genders.find((item) => item.value === userStore.userInfo?.gender)?.label ?? '未知' }}
        </SettingsListRow>
        <SettingsListRow title="地址" :chevron="false">
          {{ userStore.userInfo?.address || '—' }}
        </SettingsListRow>
        <SettingsListRow title="邮箱" :chevron="false">
          {{ userStore.userInfo?.email || '—' }}
        </SettingsListRow>
      </SettingsList>
    </div>

    <MeModal ref="avatarModalRef" width="480px" title="更改头像" :on-ok="handleAvatarSave">
      <div class="profile-avatar-upload">
        <Upload
          type="image"
          :src="newAvatar"
          :show-drag-box="true"
          :allow-library="true"
          :preview-width="240"
          :preview-height="240"
          preview-fit="cover"
          @success="onAvatarUpload"
        />
      </div>
      <p class="profile-hint">
        支持上传图片，或从资源库选择。
      </p>
    </MeModal>

    <MeModal ref="pwdModalRef" title="修改密码" width="440px" :on-ok="handlePwdSave">
      <n-form
        ref="pwdFormRef"
        class="settings-form"
        :model="pwdForm"
        label-placement="top"
        require-mark-placement="right"
      >
        <n-form-item label="原密码" path="oldPassword" :rule="required">
          <n-input v-model:value="pwdForm.oldPassword" type="password" placeholder="请输入原密码" show-password-on="mousedown" />
        </n-form-item>
        <n-form-item label="新密码" path="newPassword" :rule="required">
          <n-input v-model:value="pwdForm.newPassword" type="password" placeholder="请输入新密码" show-password-on="mousedown" />
        </n-form-item>
      </n-form>
    </MeModal>

    <MeModal ref="profileModalRef" title="修改资料" width="440px" :on-ok="handleProfileSave">
      <n-form ref="profileFormRef" class="settings-form" :model="profileForm" label-placement="top">
        <n-form-item label="昵称" path="nickName">
          <n-input v-model:value="profileForm.nickName" placeholder="请输入昵称" />
        </n-form-item>
        <n-form-item label="性别" path="gender">
          <n-select
            v-model:value="profileForm.gender"
            :options="genders"
            placeholder="请选择性别"
          />
        </n-form-item>
        <n-form-item label="地址" path="address">
          <n-input v-model:value="profileForm.address" placeholder="请输入地址" />
        </n-form-item>
        <n-form-item label="邮箱" path="email">
          <n-input v-model:value="profileForm.email" placeholder="请输入邮箱" />
        </n-form-item>
      </n-form>
    </MeModal>
  </AppPage>
</template>

<script setup>
import { MeModal, SettingsList, SettingsListRow, Upload } from '@/components'
import { useForm, useModal } from '@/composables'
import { useUserStore } from '@/store'
import { getUserInfo } from '@/store/helper'
import api from './api'

const userStore = useUserStore()
const required = {
  required: true,
  message: '此为必填项',
  trigger: ['blur', 'change'],
}

async function refreshUser() {
  const next = await getUserInfo()
  userStore.setUser({
    ...userStore.userInfo,
    ...next,
  })
}

const [avatarModalRef] = useModal()
const newAvatar = ref(userStore.avatar || '')

function openAvatarModal() {
  newAvatar.value = userStore.avatar || ''
  avatarModalRef.value?.open()
}

function onAvatarUpload(res) {
  newAvatar.value = res?.url || ''
}

async function handleAvatarSave() {
  if (!newAvatar.value) {
    $message.error('请先上传头像')
    return false
  }
  await api.updateProfile({ avatar: newAvatar.value })
  await refreshUser()
  $message.success('头像修改成功')
}

const [pwdModalRef] = useModal()
const [pwdFormRef, pwdForm, pwdValidation] = useForm({
  oldPassword: '',
  newPassword: '',
})
async function handlePwdSave() {
  await pwdValidation()
  await api.changePassword(pwdForm.value)
  pwdForm.value.oldPassword = ''
  pwdForm.value.newPassword = ''
  $message.success('密码修改成功')
}

const [profileModalRef] = useModal()
const [profileFormRef, profileForm, profileValidation] = useForm({
  nickName: userStore.nickName,
  gender: userStore.userInfo?.gender,
  address: userStore.userInfo?.address,
  email: userStore.userInfo?.email,
})
const genders = [
  { label: '男', value: 1 },
  { label: '女', value: 2 },
]

function openProfileModal() {
  profileForm.value = {
    nickName: userStore.nickName,
    gender: userStore.userInfo?.gender,
    address: userStore.userInfo?.address,
    email: userStore.userInfo?.email,
  }
  profileModalRef.value?.open()
}

async function handleProfileSave() {
  await profileValidation()
  await api.updateProfile({
    nickName: profileForm.value.nickName,
    gender: profileForm.value.gender,
    address: profileForm.value.address,
    email: profileForm.value.email,
  })
  await refreshUser()
  $message.success('资料修改成功')
}
</script>

<style scoped>
.profile {
  max-width: 760px;
  margin: 0 auto;
  padding: 8px 4px 24px;
}

.profile-hero {
  display: flex;
  align-items: center;
  gap: 20px;
  margin-bottom: 20px;
  padding: 8px 16px 20px;
}

.profile-hero__avatar {
  box-shadow: 0 0 0 3px rgb(255 255 255 / 80%);
  flex-shrink: 0;
}

.profile-hero__meta h2 {
  margin: 0;
  font-family: var(--font-display);
  font-size: var(--title-lg, 28px);
  font-weight: 600;
  letter-spacing: -0.028em;
  color: var(--text-primary, #1a2332);
}

.profile-hero__meta p {
  margin: 6px 0 0;
  font-size: 13px;
  letter-spacing: -0.01em;
  color: var(--text-secondary, #5b6b7c);
}

.profile-hero__actions {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 16px;
}

.profile-hero__actions :deep(.n-button) {
  --n-height: var(--ctl-h-lg, 28px) !important;
  --n-font-size: 13px !important;
  --n-border-radius: 7px !important;
  font-weight: 550 !important;
  letter-spacing: -0.012em;
}

.profile-avatar-upload {
  display: flex;
  justify-content: center;
  padding: 8px 0 4px;
}

.profile-hint {
  margin: 10px 0 0;
  font-size: 12px;
  line-height: 1.5;
  color: #8b97a5;
  letter-spacing: -0.01em;
  text-align: center;
}

@media (max-width: 640px) {
  .profile-hero {
    flex-direction: column;
    align-items: flex-start;
  }
}
</style>
