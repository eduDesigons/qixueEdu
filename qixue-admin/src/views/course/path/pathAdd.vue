<script setup>
import { NButton, NForm, NFormItem, NImage, NInput, NInputNumber, NSelect, NSwitch } from 'naive-ui'
import { useRoute, useRouter } from 'vue-router'
import api from '@/api/index.js'
import { CommonPage, Tip, Upload, WangEditor } from '@/components/index.js'
import { useForm, useUnsavedGuard } from '@/composables'
import { useTabStore } from '@/store'
import { getAssetUrl } from '@/utils'

defineOptions({
  name: 'AddPath',
})

const router = useRouter()
const route = useRoute()

// 系列分类选项
const categories = ref([])
function getCategories() {
  api.path.getPathCategories().then((res) => {
    categories.value = res.data || []
  }).catch((error) => {
    $message.error(error?.message || '获取系列分类失败')
    categories.value = []
  })
}
getCategories()

const [formRef, form, validation] = useForm({
  id: null,
  title: '',
  thumb: '',
  desc: '',
  content: '',
  category_id: null,
  is_rec: 0,
  is_open: 0,
  min_need_credit: 0,
  is_show: 0,
  cert_id: null,
})

const certOptions = ref([])
const certList = ref([])

function rebuildCertOptions() {
  const currentId = form.value.cert_id
  certOptions.value = certList.value.map(item => ({
    label: Number(item.is_show) === 1 ? item.name : `${item.name}（已下架）`,
    value: item.id,
    disabled: Number(item.is_show) !== 1 && item.id !== currentId,
  }))
}

api.cert.options().then((res) => {
  certList.value = res.data || []
  rebuildCertOptions()
}).catch((error) => {
  $message.error(error?.message || '获取证书列表失败')
  certOptions.value = []
  certList.value = []
})

const courseCount = ref(0)
const canPublish = computed(() => courseCount.value > 0)
const selectedCert = computed(() =>
  certList.value.find(item => item.id === form.value.cert_id) || null,
)

watch(() => form.value.cert_id, () => rebuildCertOptions())

function onShowChange(v) {
  if (v === 1 && !canPublish.value) {
    $message.warning(route.params.id ? '请先添加课程后再发布' : '请先保存系列并添加课程后再发布')
    return
  }
  form.value.is_show = v
}

const tabStore = useTabStore()
const saving = ref(false)
const { snapshot, allowLeave } = useUnsavedGuard(() => form.value)
const initialCertId = ref(null)

if (route.params.id) {
  api.path.getPathDetail(route.params.id).then((res) => {
    form.value = {
      ...res.data,
      cert_id: res.data.cert_id || null,
    }
    initialCertId.value = res.data.cert_id || null
    courseCount.value = Number(res.data.courses_count || 0)
    if (form.value.is_show && courseCount.value <= 0)
      form.value.is_show = 0
    nextTick(snapshot)
  }).catch((error) => {
    $message.error(error?.message || '加载系列失败')
  })
}
else {
  onMounted(() => nextTick(snapshot))
}

function confirmUnbindCert() {
  return new Promise((resolve) => {
    $dialog.warning({
      title: '取消证书绑定',
      content: '取消绑定后，新结业学员将不再自动发证；已颁发的证书不会自动删除。',
      positiveText: '仍要取消',
      negativeText: '返回修改',
      onPositiveClick: () => resolve(true),
      onNegativeClick: () => resolve(false),
      onClose: () => resolve(false),
    })
  })
}

// 提交表单
async function handleSubmit() {
  if (saving.value)
    return
  try {
    await validation()
  }
  catch {
    return
  }

  const hadCert = !!initialCertId.value
  const willUnbind = hadCert && !form.value.cert_id
  if (willUnbind) {
    const ok = await confirmUnbindCert()
    if (!ok)
      return
  }

  if (form.value.is_show && !canPublish.value) {
    $message.warning(route.params.id ? '请先添加课程后再发布' : '请先保存系列并添加课程后再发布')
    return
  }

  saving.value = true
  try {
    const res = await api.path.pathSave(form.value)
    const payload = res?.data || {}
    $message.success(route.params.id ? '系列更新成功' : '系列创建成功')
    allowLeave()

    if (payload.cert_changed && form.value.cert_id) {
      const entityId = payload.id || form.value.id || route.params.id
      const pending = Number(payload.pending_cert_users || 0)
      const issued = Number(payload.issued_cert_users || 0)
      const swapped = !!payload.cert_swapped
      if (pending > 0 || (swapped && issued > 0)) {
        await askCertBackfill({
          pathId: entityId,
          pending,
          issued,
          swapped,
        })
      }
      else {
        $message.info('证书已绑定，当前无待补发学员')
      }
    }

    const currentPath = route.path
    const hadOthers = tabStore.tabs.some(tab => tab.path !== currentPath)
    await tabStore.removeTab(currentPath)
    if (!hadOthers) {
      if (window.history.length > 1)
        router.back()
      else
        router.push('/')
    }
  }
  catch {
    // 拦截器已提示失败原因
  }
  finally {
    saving.value = false
  }
}

function askCertBackfill({ pathId, pending, issued = 0, swapped = false }) {
  return new Promise((resolve) => {
    const lines = []
    if (pending > 0)
      lines.push(`当前有 ${pending} 名已完课学员尚未获得证书。`)
    if (swapped && issued > 0)
      lines.push(`另有 ${issued} 名学员仍持有旧版证书，不会自动更新，可在学员列表或「已颁发」中重新生成。`)
    if (pending > 500)
      lines.push('单次最多补发 500 人，剩余学员可再次操作补发。')
    lines.push(pending > 0 ? '是否立即为未获证学员补发？' : '是否继续？')

    $dialog.warning({
      title: '补发结业证书',
      content: lines.join('\n'),
      positiveText: pending > 0 ? '补发' : '知道了',
      negativeText: pending > 0 ? '暂不补发' : '关闭',
      async onPositiveClick() {
        if (pending <= 0)
          return
        try {
          const res = await api.cert.backfill({ path_id: pathId })
          const count = res?.data?.count
          if (count > 0)
            $message.success(`已为 ${count} 名学员补发证书`)
          else
            $message.success(res?.message || '暂无需要补发的学员')
        }
        catch {
          // 拦截器已提示
        }
      },
      onAfterLeave: resolve,
    })
  })
}

// 返回系列列表
function handleBack() {
  tabStore.removeTab(route.path)
  router.back()
}
</script>

<template>
  <CommonPage
    back
    show-headline
    :title="route.params.id ? '编辑系列' : '新增系列'"
  >
    <template #action>
      <NButton type="primary" :loading="saving" :disabled="saving" @click="handleSubmit">
        保存
      </NButton>
    </template>
    <NForm
      ref="formRef"
      class="settings-form entity-form"
      :model="form"
      label-placement="top"
      :rules="{
          title: {
            required: true,
            message: '请输入系列标题',
            trigger: ['input', 'blur'],
          },
          thumb: {
            required: true,
            type: 'string',
            message: '请上传封面图片',
            trigger: ['input', 'blur'],
          },
          desc: {
            required: true,
            message: '请输入简短介绍',
            trigger: ['input', 'blur'],
          },
          category_id: {
            required: true,
            type: 'number',
            message: '请选择系列分类',
            trigger: ['blur', 'change'],
          },
          min_need_credit: {
            required: true,
            type: 'number',
            message: '请输入最低学分',
            trigger: ['input', 'blur'],
          },
        }"
      >
        <NFormItem label="系列分类" path="category_id">
          <NSelect
            v-model:value="form.category_id"
            placeholder="请选择系列分类"
            :options="categories"
            :show-path="true"
            :filterable="true"
            value-field="id"
            label-field="name"
          />
        </NFormItem>
        <NFormItem label="系列标题" path="title">
          <NInput
            v-model:value="form.title"
            autofocus
            placeholder="请输入系列标题"
          />
        </NFormItem>

        <NFormItem label="封面" path="thumb">
          <Upload
            :src="form.thumb"
            type="image"
            :show-drag-box="true"
            :preview-width="169"
            :preview-height="169"
            preview-fit="cover"
            :max-edge="640"
            @success="(data) => form.thumb = data.url"
          />
        </NFormItem>
        <div class="form-switch-group">
          <NFormItem label="推荐" path="is_rec">
            <template #label>
              <Tip label="推荐" content="推荐后，会出现在首页推荐系列内" />
            </template>
            <NSwitch
              v-model:value="form.is_rec"
              :checked-value="1"
              :unchecked-value="0"
            >
              <template #checked>
                推荐
              </template>
              <template #unchecked>
                关闭
              </template>
            </NSwitch>
          </NFormItem>

          <NFormItem label="是否公开" path="is_open">
            <template #label>
              <Tip label="是否公开" content="公开后，未认证用户可以在系列列表中看到该系列" />
            </template>
            <NSwitch
              v-model:value="form.is_open"
              :checked-value="1"
              :unchecked-value="0"
            >
              <template #checked>
                公开
              </template>
              <template #unchecked>
                私有
              </template>
            </NSwitch>
          </NFormItem>
          <NFormItem label="是否发布" path="is_show">
            <template #label>
              <Tip label="是否发布" content="至少添加 1 门课程后才能发布。新建系列请先保存，再去详情页添加课程。" />
            </template>
            <NSwitch
              :value="form.is_show"
              :checked-value="1"
              :unchecked-value="0"
              @update:value="onShowChange"
            >
              <template #checked>
                已发布
              </template>
              <template #unchecked>
                未发布
              </template>
            </NSwitch>
          </NFormItem>
        </div>

        <div class="form-pair">
          <NFormItem label="限制学分" path="min_need_credit">
            <template #label>
              <Tip label="限制学分" content="学员学分满足多少才可加入本系列" />
            </template>
            <NInputNumber
              v-model:value="form.min_need_credit"
              :min="0"
              placeholder="请输入限制学分"
            />
          </NFormItem>
        </div>

        <NFormItem label="结业证书" path="cert_id">
          <template #label>
            <Tip label="结业证书" content="学员学完本系列全部课程后自动颁发。单门课程也可单独绑定证书。首次绑定或更换模板后保存时，可选择为已完课学员补发。" />
          </template>
          <div class="cert-bind">
            <NSelect
              v-model:value="form.cert_id"
              placeholder="不颁发证书"
              :options="certOptions"
              clearable
              filterable
            />
            <Transition name="soft-reveal">
              <div v-if="selectedCert?.preview_image || selectedCert?.template_image" class="cert-bind__preview">
                <NImage
                  :src="getAssetUrl(selectedCert.preview_image || selectedCert.template_image)"
                  width="200"
                  object-fit="cover"
                  style="border-radius: 10px"
                />
                <span class="cert-bind__name">{{ selectedCert.name }}</span>
              </div>
            </Transition>
          </div>
        </NFormItem>
        <NFormItem label="简短介绍" path="desc">
          <NInput
            v-model:value="form.desc"
            type="textarea"
            placeholder="请输入简短介绍"
            :autosize="{ minRows: 3, maxRows: 6 }"
          />
        </NFormItem>

        <NFormItem label="详情介绍" path="content">
          <WangEditor v-model:value="form.content" />
        </NFormItem>
      </NForm>
  </CommonPage>
</template>

<style scoped>
.cert-bind {
  display: flex;
  flex-direction: column;
  gap: 12px;
  width: 100%;
}

.cert-bind__preview {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 8px;
  padding: 10px;
  border-radius: 14px;
  background: var(--ctl-field, rgba(var(--primary-color), 0.04));
  border: 1px solid var(--ctl-hairline, rgba(var(--primary-color), 0.08));
}

.cert-bind__preview :deep(.n-image) {
  transition: transform 0.28s cubic-bezier(0.32, 0.72, 0, 1);
}

.cert-bind__preview:hover :deep(.n-image) {
  transform: scale(1.02);
}

.cert-bind__name {
  font-size: 13px;
  color: var(--ctl-text-soft, #5b6b7c);
}
</style>
