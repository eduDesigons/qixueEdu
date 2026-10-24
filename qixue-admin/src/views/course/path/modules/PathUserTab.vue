<script setup>
import { NAvatar, NButton, NTag } from 'naive-ui'
import { h, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import api from '@/api/index.js'
import { CityQueryItem, MeCrud, MeModal, MeQueryItem, SoftUserCard } from '@/components'
import { previewImage, useCrud } from '@/composables/index.js'
import { getAssetUrl } from '@/utils'

const props = defineProps({
  pathId: {
    type: [String, Number],
    required: true,
  },
  /** 系列是否绑定了结业证书模板 */
  hasCert: {
    type: Boolean,
    default: false,
  },
})

const emit = defineEmits(['userChange'])

const queryItems = ref({
  path_id: props.pathId,
  nick_name: null,
  mobile: null,
  city_id: null,
})

const router = useRouter()

const $table = ref(null)
onMounted(() => {
  $table.value?.handleSearch()
})

const {
  handleDelete,
} = useCrud({
  name: '学员',
  doDelete: api.path.userPathDelete,
  refresh: () => $table.value?.handleSearch(),
})

function openUser(row) {
  router.push({ name: 'UserDetail', params: { id: row.user_id } })
}

function openCert(row, e) {
  e?.stopPropagation?.()
  previewImage(row.user_cert?.image)
}

function removeUser(row, e) {
  e?.stopPropagation?.()
  handleDelete(row.id, {
    title: '移出学员',
    content: '将移出本系列及学习记录；已颁发的结业证书仍会保留，可在证书管理或学员详情中查看。',
  })
}

function certTag(user_cert) {
  if (!user_cert)
    return h('span', { class: 'text-secondary' }, '-')
  if (!user_cert.image)
    return h(NTag, { type: 'warning', size: 'small', round: true, bordered: false }, '待生成')
  return h(NTag, { type: 'success', size: 'small', round: true, bordered: false }, '已颁发')
}

function certMetric(user_cert) {
  if (!user_cert)
    return '未颁发'
  if (!user_cert.image)
    return '待生成'
  return '已颁发'
}

function issueCert(row, e) {
  e?.stopPropagation?.()
  if (!props.hasCert) {
    $message.warning('该系列未绑定证书模板')
    return
  }
  if (!row.is_end) {
    $message.warning('学员尚未结业')
    return
  }
  const name = row.user?.nick_name || '该学员'
  const isRegen = !!(row.user_cert)
  $dialog.warning({
    title: isRegen ? '重新生成证书' : '补发结业证书',
    content: isRegen
      ? `将按最新模板为「${name}」重新生成证书图片。`
      : `确认为「${name}」补发结业证书？`,
    positiveText: '确定',
    negativeText: '取消',
    async onPositiveClick() {
      try {
        await api.cert.issue({
          user_id: row.user_id,
          path_id: props.pathId,
          force: isRegen,
        })
        $message.success(isRegen ? '已重新生成' : '已补发')
        $table.value?.handleSearch()
        emit('userChange')
      }
      catch {
        // 拦截器已提示
      }
    },
  })
}

const columns = [
  {
    title: '学员',
    key: 'user',
    render({ user, user_id }) {
      return h(
        'div',
        {
          class: 'flex cursor-pointer',
          onClick: () => router.push({ name: 'UserDetail', params: { id: user_id } }),
        },
        [
          h(NAvatar, {
            size: 60,
            src: getAssetUrl(user?.avatar) || undefined,
            class: 'flex-shrink-0',
          }),
          h(
            'div',
            { class: 'ml-12 flex-col justify-center' },
            [
              h('div', { class: 'text-16 font-bold text-primary' }, user?.nick_name || '未命名'),
              h(NTag, { type: 'info', size: 'small', bordered: false }, `${user?.credit ?? 0}学分`),
            ],
          ),
        ],
      )
    },
  },
  {
    title: '状态',
    key: 'is_end',
    align: 'center',
    width: 100,
    render: ({ is_end }) => h(NTag, { type: is_end ? 'success' : 'error', size: 'small', round: true, bordered: false }, is_end ? '已学完' : '未学完'),
  },
  {
    title: '证书',
    key: 'user_cert',
    align: 'center',
    width: 110,
    render: ({ user_cert }) => certTag(user_cert),
  },
  {
    title: '加入时间',
    key: 'created_at',
    align: 'center',
  },
  {
    title: '操作',
    key: 'actions',
    width: 220,
    align: 'center',
    fixed: 'right',
    hideInExcel: true,
    render(row) {
      const buttons = []
      if (row.user_cert?.image) {
        buttons.push(h(
          NButton,
          {
            size: 'small',
            type: 'info',
            secondary: true,
            onClick: () => openCert(row),
          },
          { default: () => '查看证书' },
        ))
      }
      else if (props.hasCert && row.is_end) {
        buttons.push(h(
          NButton,
          {
            size: 'small',
            type: 'primary',
            secondary: true,
            onClick: () => issueCert(row),
          },
          { default: () => row.user_cert ? '重新生成' : '补发' },
        ))
      }
      buttons.push(h(
        NButton,
        {
          size: 'small',
          type: 'error',
          secondary: true,
          onClick: () => removeUser(row),
        },
        {
          default: () => '删除',
          icon: () => h('i', { class: 'i-material-symbols:delete-outline text-14' }),
        },
      ))
      return h('div', { class: 'row-actions', style: { display: 'flex', gap: '8px', justifyContent: 'center', flexWrap: 'wrap' } }, buttons)
    },
  },
]

const checkAddUserList = ref([])
const addUserSearchList = ref([])
const $modalRef = ref(null)
function showAddUserModal() {
  addUserSearchList.value = []
  checkAddUserList.value = []
  searchKey.value = ''
  $modalRef.value?.open()
}

const searchKey = ref('')
function filterAddUser() {
  if (searchKey.value.length) {
    api.path.userPathSearch({ nick_name: searchKey.value }).then((res) => {
      addUserSearchList.value = res.data.map(item => ({
        ...item,
        value: item.id,
        label: item.nick_name,
      }))
    }).catch((error) => {
      $message.error(error?.message || '搜索学员失败')
      addUserSearchList.value = []
    })
  }
  else {
    addUserSearchList.value = []
    searchKey.value = ''
  }
}
async function confirmAddUser() {
  if (!checkAddUserList.value.length) {
    $message.warning('请先勾选')
    return false
  }
  try {
    await api.path.userPathListSave({
      path_id: props.pathId,
      user_list: checkAddUserList.value,
    })
  }
  catch {
    return false
  }
  $table.value?.handleSearch()
  addUserSearchList.value = []
  checkAddUserList.value = []
  searchKey.value = ''
  emit('userChange')
  return true
}

function renderLabel({ option }) {
  return h(
    'div',
    { style: { display: 'flex', margin: '6px 0' } },
    {
      default: () => [
        h(NAvatar, {
          round: true,
          src: getAssetUrl(option.avatar) || undefined,
          size: 'small',
          fallbackSrc: 'https://07akioni.oss-cn-beijing.aliyuncs.com/07akioni.jpeg',
        }),
        h('div', { style: { display: 'flex', marginLeft: '6px', alignSelf: 'center' } }, { default: () => option.label }),
      ],
    },
  )
}

defineExpose({
  showAddUserModal,
  refresh: () => $table.value?.handleSearch(false),
})
</script>

<template>
  <div class="h-500">
    <MeCrud
      ref="$table"
      v-model:query-items="queryItems"
      :bound-query="{ path_id: pathId }"
      :columns="columns"
      :get-data="api.path.pathUser"
    >
      <CityQueryItem v-model="queryItems.city_id" :label-width="50" />
      <MeQueryItem label="用户名" :label-width="50">
        <n-input
          v-model:value="queryItems.nick_name"
          type="text"
          placeholder="请输入用户名"
          clearable
        />
      </MeQueryItem>
      <MeQueryItem label="电话" :label-width="50">
        <n-input
          v-model:value="queryItems.mobile"
          type="text"
          placeholder="请输入电话"
          clearable
        />
      </MeQueryItem>

      <template #card="{ list }">
        <div class="entity-grid">
          <SoftUserCard
            v-for="item in list"
            :key="item.id"
            :avatar="item.user?.avatar"
            :name="item.user?.nick_name"
            :subtitle="`${item.user?.credit ?? 0} 学分`"
            :status="item.is_end ? '已学完' : '未学完'"
            :status-tone="item.is_end ? 'done' : 'doing'"
            :time="item.created_at"
            :metrics="[{ value: certMetric(item.user_cert), label: '证书' }]"
            deletable
            @click="openUser(item)"
            @delete="removeUser(item)"
          >
            <template #actions>
              <NButton
                v-if="item.user_cert?.image"
                size="tiny"
                type="info"
                secondary
                @click="openCert(item, $event)"
              >
                查看证书
              </NButton>
              <NButton
                v-else-if="hasCert && item.is_end"
                size="tiny"
                type="primary"
                secondary
                @click="issueCert(item, $event)"
              >
                {{ item.user_cert ? '重新生成' : '补发' }}
              </NButton>
            </template>
          </SoftUserCard>
        </div>
      </template>
    </MeCrud>
    <MeModal ref="$modalRef" width="520px" title="添加学员" :on-ok="confirmAddUser">
      <n-input-group>
        <n-input v-model:value.trim="searchKey" />
        <NButton type="primary" secondary @click="filterAddUser()">
          {{ searchKey.length ? '搜索' : '重置' }}
        </NButton>
      </n-input-group>
      <div class="mt-10">
        <n-transfer
          v-model:value="checkAddUserList"
          virtual-scroll
          :options="addUserSearchList"
          :render-target-label="renderLabel"
          :render-source-label="renderLabel"
        />
      </div>
    </MeModal>
  </div>
</template>
