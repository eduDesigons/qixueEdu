<script setup>
import { NButton, NForm, NFormItem, NInput, NInputNumber } from 'naive-ui'
import { useRoute, useRouter } from 'vue-router'
import api from '@/api/index.js'
import { CommonPage, Tip, Upload, WangEditor } from '@/components/index.js'
import { useForm, useUnsavedGuard } from '@/composables'
import { useTabStore } from '@/store'

defineOptions({
  name: 'ChapterAdd',
})

const router = useRouter()
const route = useRoute()

const type = route.query.type ? route.query.type * 1 : 1
const [formRef, form, validation] = useForm({
  id: null,
  title: '',
  course_id: route.params.courseId,
  media_id: null,
  media_url: '',
  media_thumb: '',
  type,
  desc: '',
  content: '',
  limit_day: null,
})

const tabStore = useTabStore()
const saving = ref(false)
const { snapshot, allowLeave } = useUnsavedGuard(() => form.value)

if (route.params.id) {
  api.chapter.getChapterDetail(route.params.id).then((res) => {
    form.value = res.data
    nextTick(snapshot)
  }).catch((error) => {
    $message.error(error?.message || '加载课节失败')
  })
}
else {
  onMounted(() => nextTick(snapshot))
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
  saving.value = true
  try {
    await api.chapter.chapterSave(form.value)
    $message.success(route.params.id ? '课节更新成功' : '课节创建成功')
    allowLeave()
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

// 返回课程列表
function handleBack() {
  tabStore.removeTab(route.path)
  router.back()
}
</script>

<template>
  <CommonPage
    back
    show-headline
    :title="route.params.id ? '编辑课节' : '新增课节'"
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
            message: '请输入课程标题',
            trigger: ['input', 'blur'],
          },
          desc: {
            required: true,
            message: '请输入简短介绍',
            trigger: ['input', 'blur'],
          },
          media_id: {
            required: type === 1 || type === 2,
            type: 'number',
            message: '请上传资源',
            trigger: ['input', 'blur'],
          },
        }"
      >
        <NFormItem label="课节标题" path="title">
          <NInput
            v-model:value="form.title"
            autofocus
            placeholder="请输入课节标题"
          />
        </NFormItem>

        <NFormItem v-if="form.type === 1 || form.type === 2" label="资源" path="media_id">
          <Upload
            :src="form.media_url"
            :thumb="form.media_thumb || ''"
            :type="form.type === 1 ? 'video' : 'audio'"
            :show-drag-box="true"
            :cate-id="form.type === 1 ? 16 : 8"
            :preview-width="169"
            :preview-height="169"
            preview-fit="cover"
            @success="(data) => { form.media_id = data.id; form.media_url = data.url; form.media_thumb = data.thumb }"
          />
        </NFormItem>

        <NFormItem label="限制天数" path="limit_day">
          <template #label>
            <Tip label="限制天数" content="不填则不限制，超过限制天数，解锁下节会扣除积分" />
          </template>
          <NInputNumber
            v-model:value="form.limit_day"
            :min="0"
            placeholder="请输入限制天数"
          />
        </NFormItem>

        <NFormItem label="内容" path="content">
          <WangEditor v-model:value="form.content" />
        </NFormItem>
      </NForm>
  </CommonPage>
</template>
