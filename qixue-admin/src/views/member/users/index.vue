<script setup>
import { NAvatar, NButton, NCascader, NInput, NSelect, NTag, useMessage } from 'naive-ui'
import { useRouter } from 'vue-router'
import api from '@/api'
import { MeCrud, MeModal, MeQueryItem, CityQueryItem, SendMessageModal } from '@/components'
import { CommonPage } from '@/components/index.js'
import { useCrud, useRefreshOnActivate } from '@/composables'

import { getAssetUrl, formatDateTime, h } from '@/utils'

defineOptions({ name: 'MemberUser' })

const router = useRouter()
const $table = ref(null)
const $msgModal = ref(null)
/** QueryBar筛选参数（可选） */
const queryItems = ref({
  nick_name: '',
  account: '',
  email: '',
  city_id: null,
  area_id: '',
  is_active: '',
  is_lock: '',
})

const message = useMessage()

useRefreshOnActivate(() => $table.value?.handleSearch())

function locked(row) {
  return Number(row?.is_lock) === 1
}

function handleSetLock(row) {
  if (!row?.id)
    return
  const next = locked(row) ? -1 : 1
  $dialog.warning({
    title: '提示',
    content: next === 1 ? '拉黑后该学员无法登录。确定拉黑？' : '确定解除拉黑？',
    positiveText: '确定',
    negativeText: '取消',
    async onPositiveClick() {
      try {
        await api.user.setLock({ id: row.id, is_lock: next })
        message.success(next === 1 ? '已拉黑' : '已解除拉黑')
        $table.value?.handleSearch()
      }
      catch (e) {
        message.error(e?.message || '操作失败')
      }
    },
  })
}

onMounted(() => {
  $table.value?.handleSearch()
})

// 区域选项
const areas = ref([])
function getAreas() {
  api.common.getAreaTree().then((res) => {
    areas.value = res.data || []
  }).catch((error) => {
    message.error(error?.message || '获取区域列表失败')
    areas.value = []
  })
}
getAreas()

const categoryOptions = ref([])
function loadCategoryOptions() {
  api.course.getCourseCategories().then((res) => {
    categoryOptions.value = (res.data || []).map(item => ({
      label: item.name,
      value: item.id,
      credit: Number(item.base_credit) || 0,
    }))
  }).catch(() => {
    categoryOptions.value = []
  })
}
loadCategoryOptions()

// 使用 useCrud 管理新增和编辑
const {
  modalRef,
  modalFormRef,
  modalForm,
  modalAction,
  handleAdd,
  handleEdit,
  handleDelete: _handleDelete,
} = useCrud({
  name: '学员',
  initForm: {
    nick_name: '',
    account: '',
    mobile: '',
    area_id: null,
    category_id: null,
    credit: 0,
    coins: 0,
    password: '',
    email: '',
  },
  doCreate: data => api.user.createUser(data),
  doUpdate: (data) => {
    // 编辑时，如果密码为空则不传递密码字段
    const processedData = { ...data }
    if (!processedData.password) {
      delete processedData.password
    }
    const hit = categoryOptions.value.find(item => item.value === processedData.category_id)
    if (hit) {
      const current = Number(processedData.credit) || 0
      // 当前学分已高于基础学分时不覆盖
      if (current <= hit.credit)
        processedData.credit = hit.credit
    }
    return api.user.updateUser(processedData)
  },
  doDelete: id => Promise.resolve({ id }), // 删除功能暂不实现
  refresh: () => $table.value?.handleSearch(),
})

const selectedCredit = computed(() => {
  const hit = categoryOptions.value.find(item => item.value === modalForm.value?.category_id)
  return hit ? hit.credit : null
})

// 🌟 1. 單獨宣告說明文檔的 ref
const tipModalRef = ref(null)
const enrollDoc = ref({ title: '入学要求', content: '' })
function handleOpenDoc() {
  api.config.getInfo('enrollConfig').then((res) => {
    const data = res.data && typeof res.data === 'object' ? res.data : {}
    enrollDoc.value = {
      title: data.title || '入学要求',
      content: typeof data.content === 'string' ? data.content : '',
    }
  }).catch(() => {
    enrollDoc.value = { title: '入学要求', content: '' }
  }).finally(() => {
    tipModalRef.value?.open()
  })
}
// 编辑时按当前学分回填学员等级（一级分类）
function handleEditWithArea(row) {
  const credit = Number(row?.credit) || 0
  const hit = categoryOptions.value.find(item => item.credit === credit)
  handleEdit({
    ...row,
    category_id: hit ? hit.value : null,
    password: '',
  })
}

const columns = [
  {
    title: '编号',
    key: 'id',
    width: 80,
    align: 'center',
  },
  {
    title: '学员',
    key: 'course-info',
    width: 300,
    render(row) {
      return h(
        'div',
        {
          class: 'flex cursor-pointer',
          onClick: () => handleViewPath(row),
        },
        [
          h(NAvatar, {
            size: 60,
            round: true,
            src: getAssetUrl(row.avatar) || undefined,
            class: 'flex-shrink-0',
          }),
          h(
            'div',
            {
              class: 'ml-12 flex-col justify-center',
            },
            [
              h(
                'div',
                {
                  class: 'text-14 font-bold text-primary',
                },
                row.nick_name,
              ),
              h(
                'div',
                {
                  class: 'text-12 text-secondary mt-1',
                },
                row.mobile,
              ),
              h(
                'div',
                {
                  class: 'text-12 text-secondary mt-1',
                },
                row.account,
              ),
            ],
          ),
        ],
      )
    },
  },
  {
    title: '区域',
    key: 'category',
    align: 'center',
    width: 120,
    render: ({ area }) => {
      if (!area)
        return h(NTag, { type: 'info', size: 'small', round: true, bordered: false }, { default: () => '未设置' })
      return h(NTag, { type: 'info', size: 'small', round: true, bordered: false }, { default: () => `${area.parent_name} ${area.name}` || '未设置' })
    },
  },

  // 课节数
  {
    title: '课程',
    key: 'courses_count',
    align: 'center',
    width: 100,
    render: ({ courses_count }) => {
      return h('span', { class: 'text font-bold' }, courses_count)
    },
  },
  {
    title: '系列课',
    key: 'coins',
    align: 'center',
    width: 100,
    render: ({ paths_count }) => {
      return h('span', { class: 'text font-bold' }, paths_count)
    },
  },
  {
    title: '收藏',
    key: 'coins',
    align: 'center',
    width: 100,
    render: ({ likes_count }) => {
      return h('span', { class: 'text font-bold' }, likes_count)
    },
  },
  {
    title: '证书',
    key: 'coins',
    align: 'center',
    width: 100,
    render: ({ certs_count }) => {
      return h('span', { class: 'text font-bold' }, certs_count)
    },
  },
  // 所需学分
  {
    title: '学分',
    key: 'credit',
    align: 'center',
    width: 100,
    render: ({ credit }) => {
      return h('span', { class: 'text font-bold' }, credit)
    },
  },
  {
    title: '积分',
    key: 'coins',
    align: 'center',
    width: 100,
    render: ({ coins }) => {
      return h('span', { class: 'text font-bold' }, coins)
    },
  },
  {
    title: '状态',
    key: 'status',
    align: 'center',
    width: 160,
    render(row) {
      if (!row)
        return null
      return h('div', { class: 'flex gap-6 justify-center flex-wrap' }, [
        row.is_active === 1 ? h(NTag, { type: 'primary', bordered: false, size: 'small', round: true }, { default: () => '已认证' }) : null,
        locked(row) ? h(NTag, { type: 'error', bordered: false, size: 'small', round: true }, { default: () => '已拉黑' }) : null,
      ].filter(Boolean))
    },
  },
  {
    title: '注册时间',
    key: 'created_at',
    align: 'center',
    width: 180,
    render(row) {
      return h('span', formatDateTime(row.created_at))
    },
  },
  {
    title: '操作',
    key: 'action',
    align: 'center',
    width: 300,
    fixed: 'right',
    render(row) {
      return h('div', { class: 'row-actions' }, [
        h(NButton, {
          size: 'small',
          type: 'primary',
          secondary: true,
          onClick: () => {
            if (row && row.id)
              handleEditWithArea(row)
          },
        }, { default: () => '编辑' }),

        h(NButton, {
          size: 'small',
          type: 'info',
          secondary: true,
          onClick: () => handleSendMessage(row),
        }, { default: () => '站内信' }),

        h(NButton, {
          size: 'small',
          type: locked(row) ? 'warning' : 'error',
          secondary: true,
          onClick: () => handleSetLock(row),
        }, { default: () => locked(row) ? '解除拉黑' : '拉黑' }),
      ])
    },
  },
]

function handleSendMessage(row) {
  if (!row?.id) {
    message.warning('用户信息不完整')
    return
  }
  $msgModal.value?.open({
    userId: row.id,
    userName: row.nick_name,
    type: 'notice',
  })
}

// 查看学员详情
function handleViewPath(row) {
  router.push(`/member/user/${row.id}`)
}

// 表单验证规则
const rules = computed(() => {
  const action = modalAction.value || ''
  return {
    nick_name: [
      { required: true, message: '请输入昵称', trigger: 'blur' },
    ],
    area_id: [
      {
        required: true,
        type: 'any', // 如果確認 id 是數字，也可以寫 type: 'number'
        message: '区域不能为空',
        trigger: ['blur', 'change'],
      },
    ],
    account: action === 'add'
      ? [
          { required: true, message: '请输入账号', trigger: 'blur' },
        ]
      : [],
    category_id: action === 'add'
      ? [{ required: true, type: 'number', message: '请选择学员等级', trigger: ['blur', 'change'] }]
      : [],
    password: action === 'add'
      ? [
          { required: true, message: '请输入密码', trigger: 'blur' },
          { min: 6, message: '密码长度不能少于6位', trigger: 'blur' },
        ]
      : [
          { min: 6, message: '密码长度不能少于6位', trigger: 'blur' },
        ],
  }
})
</script>

<template>
  <CommonPage>
    <template #action>
      <div class="flex items-center shrink-0 whitespace-nowrap gap-12">
        <NButton type="warning" secondary @click="handleOpenDoc">
          <i class="i-material-symbols:description mr-4" />
          入学要求
        </NButton>

        <NButton type="primary" @click="handleAdd">
          <i class="i-material-symbols:add mr-4" />
          新增学员
        </NButton>
      </div>
    </template>
    <!-- 搜索组件   -->
    <MeCrud
      ref="$table"
      v-model:query-items="queryItems"
      collapse-query
      :scroll-x="1000"
      :columns="columns"
      :get-data="api.user.getUserList"
    >
      <MeQueryItem label="姓名" :label-width="60">
        <NInput
          v-model:value="queryItems.nick_name"
          type="text"
          placeholder="学员名称"
          clearable
        />
      </MeQueryItem>
      <MeQueryItem label="账号" :label-width="70">
        <NInput
          v-model:value="queryItems.account"
          type="text"
          placeholder="学员账号"
          clearable
        />
      </MeQueryItem>
      <!-- <MeQueryItem label="邮箱" :label-width="80">
         <NInput
           v-model:value="queryItems.email"
           type="text"
           placeholder="学员邮箱"
           clearable
         />
       </MeQueryItem> -->
      <CityQueryItem v-model="queryItems.city_id" :label-width="60" />
      <MeQueryItem label="区域" :label-width="60">
        <NCascader
          v-model:value="queryItems.area_id"
          :options="areas"
          check-strategy="child"
          :show-path="true"
          :filterable="true"
          clearable
          value-field="id"
          label-field="name"
        />
      </MeQueryItem>
      <MeQueryItem label="状态" :label-width="30">
        <NSelect
          v-model:value="queryItems.is_active"
          clearable
          :options="[
            { label: '已认证', value: 1 },
            { label: '未认证', value: -1 },
          ]"
        />
      </MeQueryItem>
      <MeQueryItem label="拉黑" :label-width="50">
        <NSelect
          v-model:value="queryItems.is_lock"
          clearable
          :options="[
            { label: '正常', value: -1 },
            { label: '已拉黑', value: 1 },
          ]"
        />
      </MeQueryItem>
    </MeCrud>

    <!-- 新增/编辑弹窗 -->
    <MeModal
      ref="modalRef"
    >
      <NForm
        ref="modalFormRef"
        class="settings-form"
        :model="modalForm"
        :rules="rules"
        label-placement="top"
      >
        <NFormItem v-if="modalAction === 'edit'" label="ID" path="id">
          <NInput v-model:value="modalForm.id" disabled />
        </NFormItem>
        <NFormItem v-if="modalAction === 'edit'" label="账号" path="account">
          <NInput v-model:value="modalForm.account" disabled />
        </NFormItem>
        <NFormItem label="昵称" path="nick_name">
          <NInput v-model:value="modalForm.nick_name" placeholder="请输入昵称" />
        </NFormItem>
        <NFormItem label="手机号" path="mobile">
          <NInput v-model:value="modalForm.mobile" placeholder="请输入手机号" />
        </NFormItem>
        <NFormItem label="区域" path="area_id">
          <NCascader
            v-model:value="modalForm.area_id"
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
        <NFormItem label="学员等级" path="category_id">
          <NSelect
            v-model:value="modalForm.category_id"
            :options="categoryOptions"
            placeholder="请选择一级分类"
          />
        </NFormItem>
        <p v-if="selectedCredit !== null" class="user-credit-hint">
          基础学分 {{ selectedCredit }}
          <template v-if="modalAction === 'edit'">
            · 当前 {{ modalForm.credit || 0 }}
            <template v-if="(Number(modalForm.credit) || 0) > selectedCredit">（高于基础，重选不降分）</template>
          </template>
        </p>
        <p v-else-if="modalAction === 'edit'" class="user-credit-hint">
          当前学分 {{ modalForm.credit || 0 }}
        </p>

        <NFormItem label="密码" path="password">
          <NInput
            v-model:value="modalForm.password"
            type="password"
            :placeholder="modalAction === 'add' ? '请输入密码' : '留空则不修改密码'"
            show-password-on="click"
          />
        </NFormItem>
      </NForm>
    </MeModal>

    <MeModal
      ref="tipModalRef"
      :title="enrollDoc.title || '入学要求'"
      width="760px"
    >
      <div v-if="enrollDoc.content" class="doc-content" v-html="enrollDoc.content" />
      <div v-else class="doc-content">请先在「设置 → 入学要求」中填写内容。</div>
    </MeModal>

    <SendMessageModal ref="$msgModal" />
  </CommonPage>
</template>

<style scoped>
.user-credit-hint {
  margin: -8px 0 12px;
  color: var(--ctl-text-soft, #6b7280);
  font-size: 12px;
}

.doc-content {
  font-family: var(--font-sans);
  letter-spacing: -0.01em;
  max-width: 720px;
  max-height: min(68vh, 640px);
  margin: 0 auto;
  overflow-y: auto;
  padding: 4px 6px 16px 2px;
}

.doc-content :deep(h3) {
  font-family: var(--font-display) !important;
  font-size: 22px !important;
  font-weight: 600 !important;
  letter-spacing: -0.024em !important;
  line-height: 1.25 !important;
  margin: 28px 0 12px !important;
  color: var(--ctl-text, #1a2332);
}

.doc-content :deep(h3:first-child) {
  margin-top: 4px !important;
}

.doc-content :deep(h4) {
  font-size: 15px !important;
  font-weight: 600 !important;
  letter-spacing: -0.014em !important;
  margin: 18px 0 8px !important;
  color: var(--ctl-text, #1a2332);
}

.doc-content :deep(p) {
  font-size: 14px !important;
  line-height: 1.7;
  color: var(--ctl-text-soft, #5b6b7c);
}

.doc-content::-webkit-scrollbar {
  width: 6px;
}

.doc-content::-webkit-scrollbar-thumb {
  background-color: rgba(var(--primary-color), 0.22);
  border-radius: 999px;
}

</style>
