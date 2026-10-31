<script setup>
import {
  NButton,
  NForm,
  NFormItem,
  NInput,
  NCascader,
  NModal,
  NSpin,
  NTabPane,
  NDropdown,
} from 'naive-ui'
import { computed, onActivated, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import api from '@/api'
import { CommonPage, DetailTabs, SendMessageModal } from '@/components/index.js'
import { useTabQuery, peekTabQuery } from '@/composables'
import UserCourseChapterTab from '@/views/member/users/modules/UserCourseChapterTab.vue'
import PathUserTab from './modules/PathUserTab.vue'
import UserCourseTab from './modules/UserCourseTab.vue'
import UserInfo from './modules/UserInfo.vue'
import UserPathApplyByUserTab from './modules/UserPathApplyByUserTab.vue'
import UserCourseApplyByUserTab from './modules/UserCourseApplyByUserTab.vue'
import UserCreditRecordTab from './modules/UserCreditRecordTab.vue'
import UserCoinsRecordTab from './modules/UserCoinsRecordTab.vue'
import UserTestRecordTab from './modules/UserTestRecordTab.vue'
import UserExamRecordTab from './modules/UserExamRecordTab.vue'
import UserCertTab from './modules/UserCertTab.vue'

defineOptions({
  name: 'UserDetail',
})

const router = useRouter()
const route = useRoute()
const currentTabName = ref(peekTabQuery('pathList'))
useTabQuery(currentTabName, 'pathList')
const userId = ref(route.params.id || route.query.id || '')
const userInfo = ref({})
const detailLoading = ref(true)
const $msgModal = ref(null)

const editModalShow = ref(false)
const editFormRef = ref(null)
const editForm = ref({
  id: '',
  nick_name: '',
  mobile: '',
  area_id: null,
  email: '',
  password: '',
})
const areas = ref([])

const editRules = {
  nick_name: [{ required: true, message: '请输入昵称', trigger: 'blur' }],
  password: [
    {
      validator(_rule, value) {
        if (!value)
          return true
        if (value.length < 8 || value.length > 64)
          return new Error('密码长度需要在8-64个字符之间')
        if (!/[a-zA-Z]/.test(value) || !/\d/.test(value))
          return new Error('密码必须同时包含字母和数字')
        return true
      },
      trigger: 'blur',
    },
  ],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    {
      validator(_rule, value) {
        const email = String(value || '').trim()
        if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email))
          return new Error('请输入正确的邮箱')
        return true
      },
      trigger: 'blur',
    },
  ],
}

function fetchDetail() {
  if (!userId.value) return
  if (!userInfo.value?.id)
    detailLoading.value = true
  api.user.getUserDetail(userId.value).then((res) => {
    userInfo.value = res.data || {}
  }).catch((error) => {
    $message.error(error?.message || '加载学员失败')
  }).finally(() => {
    detailLoading.value = false
  })
}

function openEdit() {
  editForm.value = {
    id: userInfo.value.id,
    nick_name: userInfo.value.nick_name ?? '',
    mobile: userInfo.value.mobile ?? '',
    area_id: userInfo.value.area_id ?? null,
    email: userInfo.value.email ?? '',
    password: '',
  }
  editModalShow.value = true
}

function handleEditSubmit() {
  editFormRef.value?.validate((err) => {
    if (err) return
    const data = { ...editForm.value }
    if (!data.password) delete data.password
    api.user.updateUser(data).then(() => {
      $message.success('保存成功')
      editModalShow.value = false
      fetchDetail()
    }).catch((e) => {
      $message.error(e?.message || '保存失败')
    })
  })
}

function openSendMessage() {
  if (!userInfo.value?.id) {
    $message.warning('学员信息未加载完成')
    return
  }
  $msgModal.value?.open({
    userId: userInfo.value.id,
    userName: userInfo.value.nick_name,
    type: 'notice',
  })
}

function handleSetLock() {
  if (!userInfo.value?.id) {
    $message.warning('学员信息未加载完成')
    return
  }
  const locked = Number(userInfo.value.is_lock) === 1
  const next = locked ? -1 : 1
  $dialog.warning({
    title: '提示',
    content: next === 1 ? '拉黑后该学员无法登录。确定拉黑？' : '确定解除拉黑？',
    positiveText: '确定',
    negativeText: '取消',
    async onPositiveClick() {
      try {
        await api.user.setLock({ id: userInfo.value.id, is_lock: next })
        $message.success(next === 1 ? '已拉黑' : '已解除拉黑')
        fetchDetail()
      }
      catch (error) {
        $message.error(error?.message || '操作失败')
      }
    },
  })
}

function handleDelete() {
  $dialog.warning({
    content: '确定删除该学员？删除后不可恢复。',
    title: '提示',
    positiveText: '确定',
    negativeText: '取消',
    async onPositiveClick() {
      try {
        await api.user.userDelete(userId.value)
        $message.success('学员删除成功')
        router.push('/member/users')
      } catch (error) {
        $message.error(error?.message || '删除失败')
      }
    },
  })
}

onMounted(() => {
  fetchDetail()
  api.common.getAreaTree().then((res) => {
    areas.value = res.data || []
  }).catch((error) => {
    $message.error(error?.message || '获取区域列表失败')
    areas.value = []
  })
})

onActivated(() => {
  if (userInfo.value?.id)
    fetchDetail()
})

watch(() => route.params.id || route.query.id, (id) => {
  if (id) {
    userId.value = id
    fetchDetail()
  }
}, { immediate: false })

const pageTitle = computed(() => userInfo.value?.nick_name || '学员详情')
const moreOptions = computed(() => [
  { label: Number(userInfo.value.is_lock) === 1 ? '解除拉黑' : '拉黑', key: 'lock' },
  { label: '删除学员', key: 'delete' },
])

function handleMoreSelect(key) {
  if (key === 'lock')
    handleSetLock()
  if (key === 'delete')
    handleDelete()
}
</script>

<template>
  <CommonPage back fill :title="pageTitle">
    <NSpin :show="detailLoading">
      <UserInfo :user-data="userInfo" />
    </NSpin>

    <DetailTabs v-model:value="currentTabName">
      <template #suffix>
          <NButton secondary @click="openSendMessage">
            站内信
          </NButton>
          <NButton type="primary" @click="openEdit">
            编辑
          </NButton>
          <NDropdown trigger="click" :options="moreOptions" @select="handleMoreSelect">
            <NButton secondary class="page-more-btn" aria-label="更多">
              <i class="i-fe:more-horizontal text-14" />
            </NButton>
          </NDropdown>
      </template>
      <NTabPane name="pathList" tab="系列列表" display-directive="show:lazy">
        <PathUserTab :key="`user-path-${userId}`" :user-id="userId" />
      </NTabPane>
        <NTabPane name="courseList" tab="课程列表" display-directive="show:lazy">
          <UserCourseTab :key="`user-course-${userId}`" :user-id="userId" />
        </NTabPane>
        <NTabPane name="chapterList" tab="课节列表" display-directive="show:lazy">
          <UserCourseChapterTab :key="`user-chapter-${userId}`" :user-id="userId" />
        </NTabPane>
        <NTabPane name="pathApply" tab="入系申请" display-directive="show:lazy">
          <UserPathApplyByUserTab :key="`user-path-apply-${userId}`" :user-id="userId" />
        </NTabPane>
        <NTabPane name="courseApply" tab="入课申请" display-directive="show:lazy">
          <UserCourseApplyByUserTab :key="`user-course-apply-${userId}`" :user-id="userId" />
        </NTabPane>
        <NTabPane name="creditRecord" tab="学分记录" display-directive="show:lazy">
          <UserCreditRecordTab :key="`user-credit-${userId}`" :user-id="userId" />
        </NTabPane>
        <NTabPane name="coinsRecord" tab="积分记录" display-directive="show:lazy">
          <UserCoinsRecordTab :key="`user-coins-${userId}`" :user-id="userId" />
        </NTabPane>
        <NTabPane name="testRecord" tab="测试记录" display-directive="show:lazy">
          <UserTestRecordTab :key="`user-test-${userId}`" :user-id="userId" />
        </NTabPane>
        <NTabPane name="examRecord" tab="考试记录" display-directive="show:lazy">
          <UserExamRecordTab :key="`user-exam-${userId}`" :user-id="userId" />
        </NTabPane>
        <NTabPane name="certList" tab="证书" display-directive="show:lazy">
          <UserCertTab :key="`user-cert-${userId}`" :user-id="userId" />
        </NTabPane>
    </DetailTabs>

    <NModal
      v-model:show="editModalShow"
      title="编辑学员"
      preset="card"
      style="width: 560px"
      :mask-closable="false"
      :positive-button-props="{ text: '确定' }"
      :negative-button-props="{ text: '取消' }"
      :close-on-positive-click="false"
      @positive-click="handleEditSubmit"
    >
      <NForm
        ref="editFormRef"
        class="settings-form"
        :model="editForm"
        :rules="editRules"
        label-placement="top"
      >
        <NFormItem label="ID" path="id">
          <NInput v-model:value="editForm.id" disabled />
        </NFormItem>
        <NFormItem label="昵称" path="nick_name">
          <NInput v-model:value="editForm.nick_name" placeholder="请输入昵称" />
        </NFormItem>
        <NFormItem label="手机号" path="mobile">
          <NInput v-model:value="editForm.mobile" placeholder="请输入手机号" />
        </NFormItem>
        <NFormItem label="区域" path="area_id">
          <NCascader
            v-model:value="editForm.area_id"
            :options="areas"
            check-strategy="child"
            :show-path="true"
            :filterable="true"
            clearable
            placeholder="请选择区域"
            value-field="id"
            label-field="name"
          />
        </NFormItem>
        <NFormItem label="邮箱" path="email">
          <NInput v-model:value="editForm.email" placeholder="邮箱" type="email" />
        </NFormItem>
        <NFormItem label="密码" path="password">
          <NInput
            v-model:value="editForm.password"
            type="password"
            placeholder="留空则不修改，8-64位且含字母和数字"
            show-password-on="click"
          />
        </NFormItem>
      </NForm>
    </NModal>
    <SendMessageModal ref="$msgModal" />
  </CommonPage>
</template>

