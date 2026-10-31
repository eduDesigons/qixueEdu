<script setup>
import { Add } from '@vicons/ionicons5'
import { NButton, NIcon, NImage, NTag } from 'naive-ui'
import { computed, onMounted, ref } from 'vue'
import { VueDraggable } from 'vue-draggable-plus'
import api from '@/api'
import { AppPage } from '@/components/index.js'
import { useRefreshOnActivate } from '@/composables'
import EditBanner from './components/EditBanner.vue'
import { getAssetUrl, formatRichContentHtml } from '@/utils/common'

defineOptions({ name: 'DecorationBanner' })

const currentIndex = ref(0)
const bannerList = ref([])

function changeCurrentIndex(index) {
  currentIndex.value = index
}

onMounted(() => {
  getBannerList()
})

useRefreshOnActivate(() => getBannerList())

function getBannerList() {
  api.decoration.getBannerList().then((res) => {
    bannerList.value = res.data || []
  })
}

function onSortEnd() {
  const ids = bannerList.value.map(item => item.id)
  api.decoration.updateBannerSort(ids).then(() => {
    $message.success('排序已保存')
  }).catch(() => {
    getBannerList()
  })
}

const info = computed(() => {
  if (bannerList.value.length > 0) {
    if (currentIndex.value > bannerList.value.length - 1) {
      return null
    }
    return bannerList.value[currentIndex.value]
  }
  else {
    return null
  }
})

function delIt() {
  if (info.value != null && info.value.id) {
    $dialog.warning({
      title: '确认删除',
      content: '确定要删除这个轮播图吗？',
      positiveText: '确定',
      negativeText: '取消',
      onPositiveClick: () => {
        api.decoration.delBanner(info.value.id).then(() => {
          bannerList.value.splice(currentIndex.value, 1)
          if (currentIndex.value !== 0 && currentIndex.value > bannerList.value.length - 1) {
            currentIndex.value -= 1
          }
          $message.success('删除成功')
        })
      },
    })
  }
}

const editBannerRef = ref()
function editIt() {
  editBannerRef.value?.show(info.value)
}
function addOne() {
  editBannerRef.value?.show()
}
</script>

<template>
  <AppPage full>
    <div class="banner-wrapper">
    <header class="page-mast">
      <h1 class="page-mast__title">轮播图</h1>
      <p class="page-mast__sub">拖拽左侧缩略图可调整顺序</p>
    </header>
    <div class="banner-container">
      <div class="banner-list">
        <div class="banner-item add" @click="addOne()">
          <NIcon :component="Add" :size="30" />
        </div>
        <VueDraggable
          v-model="bannerList"
          tag="div"
          class="banner-draggable"
          :animation="150"
          ghost-class="banner-ghost"
          chosen-class="banner-chosen"
          @end="onSortEnd"
        >
          <div
            v-for="(item, index) in bannerList"
            :key="item.id"
            class="banner-item"
            :class="{ active: currentIndex === index, disabled: item.status === 2 }"
            @click="changeCurrentIndex(index)"
          >
            <NImage :src="getAssetUrl(item.img_url)" :alt="item.name" object-fit="contain" :preview-disabled="true" />
          </div>
        </VueDraggable>
      </div>
      <div class="banner-detail">
        <div v-if="info" class="banner-detail__body">
          <div class="banner-detail__head">
            <div>
              <h2>{{ info.name || '未命名' }}</h2>
              <div class="banner-detail__tags">
                <NTag :type="info.status === 1 ? 'success' : 'default'" size="small" round :bordered="false">
                  {{ info.status === 1 ? '已启用' : '已停用' }}
                </NTag>
                <NTag type="info" size="small" round :bordered="false">
                  {{ ['', '仅展示', '链接', '内容'][info.type] }}
                </NTag>
              </div>
            </div>
            <div class="banner-detail__actions">
              <NButton type="primary" @click="editIt()">
                编辑
              </NButton>
              <NButton secondary @click="delIt()">
                删除
              </NButton>
            </div>
          </div>
          <div class="banner-preview">
            <NImage :src="getAssetUrl(info.img_url)" object-fit="contain" />
          </div>
          <div v-if="info.type === 2" class="banner-link">
            <span>链接</span>
            <a :href="info.link_url" target="_blank">{{ info.link_url }}</a>
          </div>
          <div v-if="info.type === 3 && info.content" class="banner-content">
            <div class="editor-panel-head">
              内容预览
            </div>
            <div class="banner-content__well rich-content" v-html="formatRichContentHtml(info.content)" />
          </div>
        </div>
        <div v-else class="empty-state">
          请选择或添加轮播图
        </div>
      </div>
    </div>

    <EditBanner ref="editBannerRef" @update="getBannerList()" />
    </div>
  </AppPage>
</template>

<style scoped lang="scss">
.banner-wrapper {
  height: 100%;
  min-height: 0;
  display: flex;
  flex-direction: column;
}

.banner-wrapper .page-mast {
  padding: 4px 0 12px;
}

.banner-container {
  display: flex;
  overflow: hidden;
  flex: 1;
  min-height: 0;
  gap: 16px;
}

.banner-list {
  padding: 8px;
  display: flex;
  flex-direction: column;
  gap: 8px;
  height: 100%;
  overflow-y: auto;
  width: 220px;
  flex-shrink: 0;
  border-radius: 16px;
  background: var(--ctl-field, rgba(var(--primary-color), 0.035));
  border: 1px solid var(--ctl-hairline, rgba(var(--primary-color), 0.08));
  box-shadow: inset 0 1px 0 rgb(255 255 255 / 40%);
}

.banner-draggable {
  display: flex;
  flex-direction: column;
  gap: 8px;
  flex: 1;
  min-height: 0;
}

.banner-list .banner-item {
  width: 196px;
  height: 84px;
  border-radius: 14px;
  overflow: hidden;
  cursor: pointer;
  border: 1px solid var(--ctl-hairline, rgba(var(--primary-color), 0.08));
  box-shadow: none;
  background: var(--ctl-field, rgba(var(--primary-color), 0.04));
  transition:
    transform 0.24s cubic-bezier(0.32, 0.72, 0, 1),
    border-color 0.22s ease,
    box-shadow 0.28s cubic-bezier(0.32, 0.72, 0, 1);

  &:hover {
    border-color: var(--ctl-field-border, rgba(var(--primary-color), 0.16));
  }

  &:active {
    transform: scale(0.98);
  }

  :deep(.n-image),
  :deep(img) {
    display: block;
    width: 100%;
    height: 100%;
    object-fit: contain;
  }

  &.add {
    height: 72px;
    display: flex;
    align-items: center;
    justify-content: center;
    border: 1px solid var(--ctl-hairline, rgba(var(--primary-color), 0.08));
    color: rgb(var(--primary-color));
    background: rgb(255 255 255 / 55%);
    border-radius: 14px;
    box-shadow: none;
    font-weight: 550;
    letter-spacing: -0.01em;
  }

  &.active {
    border-color: rgba(var(--primary-color), 0.45);
    box-shadow: 0 0 0 3px rgba(var(--primary-color), 0.12);
  }

  &.disabled {
    filter: grayscale(100%);
    opacity: 0.6;
  }
}

.banner-ghost {
  opacity: 0.5;
  background: rgba(var(--primary-color), 0.1);
}

.banner-chosen {
  opacity: 0.95;
  box-shadow: 0 0 0 3px rgba(var(--primary-color), 0.12);
}

.banner-detail {
  flex: 1;
  padding: 4px 4px 12px;
  overflow-y: auto;
  min-width: 0;
}

.banner-detail__head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 16px;
  margin-bottom: 16px;
}

.banner-detail__head h2 {
  margin: 0;
  font-family: var(--font-display);
  font-size: var(--title-lg, 32px);
  font-weight: 600;
  letter-spacing: -0.028em;
  line-height: 1.15;
  color: var(--ctl-text, #1a2332);
}

.banner-detail__tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-top: 10px;
}

.banner-detail__actions {
  display: flex;
  gap: 8px;
  flex-shrink: 0;
  padding-top: 6px;
}

.banner-preview {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  max-width: 720px;
  min-height: 200px;
  aspect-ratio: 16 / 7;
  border-radius: 16px;
  overflow: hidden;
  background: var(--ctl-field, rgba(var(--primary-color), 0.04));
  border: 1px solid var(--ctl-hairline, rgba(var(--primary-color), 0.08));
}

.banner-preview :deep(.n-image) {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100%;
}

.banner-preview :deep(img) {
  display: block;
  max-width: 100%;
  max-height: 100%;
  width: auto;
  height: auto;
  object-fit: contain;
}

.banner-link {
  display: flex;
  align-items: center;
  gap: 12px;
  max-width: 720px;
  margin-top: 16px;
  padding: 12px 14px;
  border-radius: 14px;
  background: var(--ctl-field, rgba(var(--primary-color), 0.04));
  border: 1px solid var(--ctl-hairline, rgba(var(--primary-color), 0.08));
  font-size: 13px;
  letter-spacing: -0.01em;
}

.banner-link span {
  flex-shrink: 0;
  font-weight: 550;
  color: var(--ctl-text-soft, #5b6b7c);
}

.banner-link a {
  min-width: 0;
  color: rgb(var(--primary-color));
  word-break: break-all;
  text-decoration: none;
}

.banner-content {
  max-width: 720px;
  margin-top: 20px;
}

.banner-content__well {
  max-height: 400px;
  overflow-y: auto;
  padding: 16px;
  border-radius: 14px;
  background: var(--ctl-field, rgba(var(--primary-color), 0.04));
  border: 1px solid var(--ctl-hairline, rgba(var(--primary-color), 0.08));
}

.empty-state {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 240px;
  color: var(--ctl-text-soft, #8b97a5);
  font-size: 13px;
  font-weight: 500;
  letter-spacing: -0.012em;
  border-radius: 16px;
  background: var(--ctl-field, rgba(var(--primary-color), 0.04));
  border: 1px solid var(--ctl-hairline, rgba(var(--primary-color), 0.08));
}
</style>
