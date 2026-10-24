<script setup>
import { cloneDeep } from 'lodash-es'
import { NForm, NFormItem, NInput } from 'naive-ui'
import { ref, watch } from 'vue'
import api from '@/api'
import Upload from '@/components/common/Upload.vue'
import { SettingsEditorShell } from '@/components/index.js'
import { useSysStore } from '@/store'

const DEFAULT_FOOTER = 'The night is nearly over; The day is almost here. So let us put aside the deeds of darkness and put on the armor of light.'

defineOptions({ name: 'EditSiteConfig' })

const props = defineProps({
  embedded: { type: Boolean, default: false },
})
const emit = defineEmits(['update', 'close'])

const open = ref(false)
const loading = ref(false)
const form = ref({
  name: '',
  logo: '',
  footer: DEFAULT_FOOTER,
})

function show() {
  loading.value = true
  api.config.getInfo('sysConfig').then((res) => {
    if (res.data) {
      const data = cloneDeep(res.data)
      delete data.app_version
      delete data.app_download_url
      delete data.assetDomain
      if (data.footer == null)
        data.footer = DEFAULT_FOOTER
      form.value = data
    }
    if (!props.embedded)
      open.value = true
  }).finally(() => {
    loading.value = false
  })
}

watch(() => props.embedded, (v) => { if (v) show() }, { immediate: true })

defineExpose({ show })

function handleLogoSuccess(data) {
  form.value.logo = data.url || ''
}

function confirm() {
  loading.value = true
  api.config.saveInfo(form.value, 'sysConfig').then(async () => {
    open.value = false
    try {
      await useSysStore().fetchSysInfo()
    }
    catch {}
    emit('update')
    emit('close')
    $message.success('保存成功')
  }).catch(() => {
    $message.error('保存失败')
  }).finally(() => {
    loading.value = false
  })
}
</script>

<template>
  <SettingsEditorShell
    v-model:show="open"
    :embedded="embedded"
    title="网站信息配置"
    :loading="loading"
    @save="confirm"
    @close="emit('close')"
  >
    <NForm class="settings-form" label-placement="top">
      <NFormItem label="网站名称">
        <NInput
          v-model:value="form.name"
          placeholder="请输入网站名称"
        />
      </NFormItem>

      <NFormItem label="页脚">
        <NInput
          v-model:value="form.footer"
          type="textarea"
          placeholder="登录页和页面底部显示的文字，留空则不显示"
          :autosize="{ minRows: 2, maxRows: 4 }"
          :maxlength="300"
          show-count
        />
      </NFormItem>

      <NFormItem label="Logo">
        <Upload
          type="image"
          :src="form.logo"
          :cate-id="1"
          :show-drag-box="true"
          :preview-width="140"
          :preview-height="140"
          preview-fit="cover"
          @success="handleLogoSuccess"
        />
      </NFormItem>
    </NForm>
  </SettingsEditorShell>
</template>
