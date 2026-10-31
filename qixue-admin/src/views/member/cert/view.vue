<script setup>
import { NButton, NImage } from 'naive-ui'
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import api from '@/api'
import { CommonPage, SettingsList, SettingsListRow } from '@/components'
import { getAssetUrl } from '@/utils/common'

defineOptions({ name: 'CertView' })

const route = useRoute()
const router = useRouter()
const cert = ref(null)
const loading = ref(true)

const certId = computed(() => route.query.id || route.params.id)
const pageTitle = computed(() => cert.value?.name || '证书详情')

const paramKeyLabels = {
  name: '姓名',
  course: '名称/标题',
  time: '时间',
  no: '证书编号',
}

onMounted(() => {
  if (!certId.value) {
    $message.warning('缺少证书 id')
    return
  }
  api.cert.detail(certId.value).then((res) => {
    cert.value = res.data
  }).catch((error) => {
    $message.error(error?.message || '加载证书失败')
  }).finally(() => {
    loading.value = false
  })
})

function goEdit() {
  router.push({ path: '/course/cert/create', query: { id: certId.value } })
}

function paramLabel(key) {
  return paramKeyLabels[key] || key
}

function paramSize(p) {
  return `x ${p.x} · y ${p.y} · ${p.w}×${p.h}`
}
</script>

<template>
  <CommonPage back show-headline :title="pageTitle">
    <template #action>
      <NButton type="primary" @click="goEdit">
        编辑
      </NButton>
    </template>

    <div v-if="loading" class="cert-view__status">
      加载中...
    </div>
    <template v-else-if="cert">
      <SettingsList label="基本信息">
        <SettingsListRow title="证书名称" :chevron="false">
          {{ cert.name }}
        </SettingsListRow>
      </SettingsList>

      <section class="cert-view__media">
        <h3 class="editor-panel-head">
          预览图
        </h3>
        <div class="cert-view__img-well">
          <NImage
            v-if="cert.preview_image || cert.template_image"
            :src="getAssetUrl(cert.preview_image || cert.template_image)"
            object-fit="contain"
            class="cert-view__img"
          />
          <div v-else class="cert-view__status cert-view__status--inset">
            暂无预览图
          </div>
        </div>
      </section>

      <SettingsList v-if="cert.params?.length" label="占位框">
        <SettingsListRow
          v-for="(p, i) in cert.params"
          :key="i"
          :title="paramLabel(p.key)"
          :desc="p.key"
          :chevron="false"
        >
          {{ paramSize(p) }}
        </SettingsListRow>
      </SettingsList>
    </template>
    <div v-else class="cert-view__status">
      未找到证书
    </div>
  </CommonPage>
</template>

<style scoped>
.cert-view__media {
  max-width: 760px;
  margin: 0 auto 20px;
}

.cert-view__img-well {
  padding: 8px;
  border-radius: 14px;
  background: var(--ctl-field, rgba(var(--primary-color), 0.04));
  border: 1px solid var(--ctl-hairline, rgba(var(--primary-color), 0.08));
}

.cert-view__img {
  display: block;
  width: 100%;
  max-width: 100%;
  border-radius: 11px;
  overflow: hidden;
}

.cert-view__img :deep(.n-image img) {
  border-radius: 11px;
}

.cert-view__status {
  max-width: 760px;
  margin: 0 auto;
  padding: 48px 20px;
  text-align: center;
  letter-spacing: -0.012em;
  font-weight: 500;
  color: var(--ctl-text-soft, #5b6b7c);
  border-radius: 14px;
  background: var(--ctl-field, rgba(var(--primary-color), 0.04));
  border: 1px solid var(--ctl-hairline, rgba(var(--primary-color), 0.08));
  box-shadow: inset 0 1px 0 rgb(255 255 255 / 40%);
}

.cert-view__status--inset {
  margin: 0;
  padding: 36px 16px;
  background: transparent;
  border: 0;
  box-shadow: none;
}
</style>
