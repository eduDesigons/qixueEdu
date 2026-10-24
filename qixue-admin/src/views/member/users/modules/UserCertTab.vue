<script setup>
import { NButton, NImage, NTag } from 'naive-ui'
import { h } from 'vue'
import api from '@/api'
import { MeCrud } from '@/components'
import { previewImage } from '@/composables'
import { formatDateTime, getAssetUrl } from '@/utils'

const props = defineProps({
  userId: {
    type: [String, Number],
    required: true,
  },
})

const queryItems = ref({ user_id: props.userId })
const $table = ref(null)
onMounted(() => {
  $table.value?.handleSearch()
})

function preview(url) {
  previewImage(url)
}

const columns = [
  {
    title: '证书',
    key: 'title',
    ellipsis: { tooltip: true },
  },
  {
    title: '来源',
    key: 'source',
    width: 200,
    ellipsis: { tooltip: true },
    render: (row) => {
      if (row.path?.title)
        return `系列 · ${row.path.title}`
      if (row.course?.title)
        return `课程 · ${row.course.title}`
      return '-'
    },
  },
  {
    title: '编号',
    key: 'cert_no',
    width: 200,
  },
  {
    title: '预览',
    key: 'image',
    width: 80,
    align: 'center',
    render: ({ image }) => image
      ? h(NImage, { src: getAssetUrl(image), width: 48, height: 32, objectFit: 'cover' })
      : '-',
  },
  {
    title: '状态',
    key: 'is_valid',
    width: 90,
    align: 'center',
    render: ({ is_valid }) => h(NTag, { type: is_valid ? 'success' : 'error', size: 'small', round: true, bordered: false }, { default: () => (is_valid ? '有效' : '无效') }),
  },
  {
    title: '颁发时间',
    key: 'issued_at',
    width: 180,
    render: row => h('span', formatDateTime(row.issued_at || row.created_at)),
  },
  {
    title: '操作',
    key: 'actions',
    width: 180,
    align: 'center',
    render(row) {
      return h('div', { class: 'flex gap-8 justify-center' }, [
        row.image
          ? h(NButton, { size: 'small', type: 'primary', secondary: true, onClick: () => preview(row.image) }, { default: () => '查看' })
          : null,
        h(NButton, {
          size: 'small',
          type: 'info',
          secondary: true,
          onClick: () => reissue(row),
        }, { default: () => row.image ? '重新生成' : '生成证书' }),
      ].filter(Boolean))
    },
  },
]

function reissue(row) {
  const payload = {
    user_id: props.userId,
    force: true,
  }
  if (row.course_id) {
    payload.course_id = row.course_id
  }
  else if (row.path_id) {
    payload.path_id = row.path_id
  }
  else {
    $message.warning('缺少系列或课程信息')
    return
  }
  $dialog.warning({
    title: '重新生成证书',
    content: '将按最新模板重新生成该证书图片。',
    positiveText: '确定',
    negativeText: '取消',
    async onPositiveClick() {
      try {
        await api.cert.issue(payload)
        $message.success('已重新生成')
        $table.value?.handleSearch()
      }
      catch {
        // 拦截器已提示失败原因
      }
    },
  })
}
</script>

<template>
  <div class="h-500">
    <MeCrud
      ref="$table"
      v-model:query-items="queryItems"
      :bound-query="{ user_id: userId }"
      :columns="columns"
      :get-data="api.user.getUserCertList"
    />
  </div>
</template>
