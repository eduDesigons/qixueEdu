<script setup>
import { Add } from '@vicons/ionicons5'
import { NButton, NIcon, NSpin, NTabPane, NTabs } from 'naive-ui'
import { computed, onMounted, ref } from 'vue'
import api from '@/api'
import { CommonPage } from '@/components/index.js'
import { useRefreshOnActivate } from '@/composables'
import { getAssetUrl } from '@/utils'
import EditCategory from './EditCategory.vue'

defineOptions({ name: 'CourseCategory' })

const tab = ref('course')
const loading = ref(false)
const courseRows = ref([])
const pathRows = ref([])
const editor = ref(null)

const parents = computed(() => courseRows.value.filter(item => !item.parent_id))

onMounted(load)
useRefreshOnActivate(load)

function load() {
  loading.value = true
  const type = tab.value
  api.course.categoryList(type).then((res) => {
    const list = res.data || []
    if (type === 'path')
      pathRows.value = list
    else
      courseRows.value = list
  }).catch(() => {}).finally(() => {
    loading.value = false
  })
}

function onTab(value) {
  tab.value = value
  load()
}

function addOne() {
  editor.value?.show(null, 0)
}

function addChild(row) {
  editor.value?.show(null, row.id)
}

function editRow(row) {
  editor.value?.show(row)
}

function courseTotal(row) {
  return (row.children || []).reduce((sum, item) => sum + (item.count || 0), 0) + (row.count || 0)
}

function iconSrc(icon) {
  if (!icon || String(icon).startsWith('/static/'))
    return ''
  return getAssetUrl(icon)
}

function markText(name) {
  return String(name || '类').trim().slice(0, 1)
}

function removeRow(row) {
  const blocked = tab.value === 'course'
    ? (row.children?.length ? '请先删除下面的二级分类' : (row.count ? '该分类下还有课程，不能删除' : ''))
    : (row.count ? '该分类下还有系列，不能删除' : '')
  if (blocked) {
    $message.warning(blocked)
    return
  }
  $dialog.warning({
    title: '确认删除',
    content: `确定删除分类「${row.name}」吗？`,
    positiveText: '删除',
    negativeText: '取消',
    onPositiveClick: () => {
      return api.course.categoryDelete({ type: tab.value, id: row.id }).then(() => {
        $message.success('删除成功')
        load()
      }).catch(() => {})
    },
  })
}
</script>

<template>
  <CommonPage subtitle="一级分类的基础学分是学员等级的入学学分。二级分类的限制学分会在新建课程时带出。">
    <template #action>
      <NButton type="primary" @click="addOne">
        <template #icon>
          <NIcon :component="Add" />
        </template>
        {{ tab === 'path' ? '新增系列分类' : '新增课程分类' }}
      </NButton>
    </template>
    <div class="cate-page">
      <NTabs :value="tab" type="line" @update:value="onTab">
        <NTabPane name="course" tab="课程分类" />
        <NTabPane name="path" tab="系列分类" />
      </NTabs>
      <NSpin :show="loading" class="cate-spin">
        <div v-if="tab === 'course'" class="cate-groups">
          <section
            v-for="group in courseRows"
            :key="group.id"
            class="cate-group"
            :class="{ 'is-off': group.is_show === 0 }"
          >
            <header class="cate-group__head">
              <div class="cate-mark">
                <img v-if="iconSrc(group.icon)" :src="iconSrc(group.icon)" alt="">
                <span v-else>{{ markText(group.name) }}</span>
              </div>
              <div class="cate-group__copy">
                <h2>{{ group.name }}</h2>
                <p>{{ (group.children || []).length }} 个二级 · {{ courseTotal(group) }} 门课 · 基础学分 {{ group.base_credit || 0 }}</p>
              </div>
              <div class="cate-group__tools">
                <span class="cate-state">{{ group.is_show === 0 ? '隐藏' : '显示' }}</span>
                <NButton size="small" secondary @click="editRow(group)">
                  编辑
                </NButton>
                <NButton v-if="!courseTotal(group)" size="small" secondary type="error" @click="removeRow(group)">
                  删除
                </NButton>
              </div>
            </header>
            <div class="cate-kids">
              <article
                v-for="child in group.children || []"
                :key="child.id"
                class="cate-kid"
                :class="{ 'is-off': child.is_show === 0 }"
              >
                <div class="cate-kid__top">
                  <h3>{{ child.name }}</h3>
                  <span class="cate-state">{{ child.is_show === 0 ? '隐藏' : '显示' }}</span>
                </div>
                <p>{{ child.count || 0 }} 门课 · 限制学分 {{ child.min_need_credit || 0 }}</p>
                <div class="cate-kid__actions">
                  <button type="button" @click="editRow(child)">
                    编辑
                  </button>
                  <button v-if="!child.count" type="button" class="is-danger" @click="removeRow(child)">
                    删除
                  </button>
                </div>
              </article>
              <button type="button" class="cate-add" @click="addChild(group)">
                <NIcon :component="Add" />
                添加二级
              </button>
            </div>
          </section>
          <div v-if="!loading && !courseRows.length" class="cate-empty">
            还没有课程分类
          </div>
        </div>

        <div v-else class="cate-board">
          <article
            v-for="item in pathRows"
            :key="item.id"
            class="cate-path"
          >
            <div class="cate-mark">
              <span>{{ markText(item.name) }}</span>
            </div>
            <div class="cate-path__copy">
              <h3>{{ item.name }}</h3>
              <p>{{ item.count || 0 }} 个系列</p>
            </div>
            <div class="cate-kid__actions">
              <button type="button" @click="editRow(item)">
                编辑
              </button>
              <button v-if="!item.count" type="button" class="is-danger" @click="removeRow(item)">
                删除
              </button>
            </div>
          </article>
          <button type="button" class="cate-add cate-add--path" @click="addOne">
            <NIcon :component="Add" />
            新增系列分类
          </button>
        </div>
      </NSpin>
    </div>
    <EditCategory ref="editor" :type="tab" :parents="parents" @update="load" />
  </CommonPage>
</template>

<style scoped>
.cate-page {
  height: 100%;
  min-height: 0;
  display: flex;
  flex-direction: column;
}

.cate-page :deep(.n-tab-pane) {
  display: none;
}

.cate-spin {
  flex: 1;
  min-height: 0;
}

.cate-spin :deep(.n-spin-container),
.cate-spin :deep(.n-spin-content) {
  height: 100%;
}

.cate-spin :deep(.n-spin-content) {
  overflow: auto;
}

.cate-groups,
.cate-board {
  display: flex;
  flex-direction: column;
  gap: 14px;
  padding: 4px 2px 12px;
}

.cate-board {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
  align-content: start;
}

.cate-group {
  padding: 16px 16px 14px;
  border: 1px solid var(--ctl-hairline, rgba(32, 42, 58, 0.08));
  border-radius: 16px;
  background: var(--ctl-glass-strong, #fff);
}

.cate-group.is-off,
.cate-kid.is-off {
  opacity: 0.55;
}

.cate-group__head {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 12px;
}

.cate-mark {
  width: 42px;
  height: 42px;
  flex-shrink: 0;
  display: grid;
  place-items: center;
  border-radius: 12px;
  background: #eef3fb;
  color: #4285f4;
  font-size: 16px;
  font-weight: 650;
  overflow: hidden;
}

.cate-mark img {
  width: 100%;
  height: 100%;
  object-fit: contain;
}

.cate-group__copy {
  min-width: 0;
  flex: 1;
}

.cate-group__copy h2,
.cate-kid h3,
.cate-path h3 {
  margin: 0;
  color: var(--ctl-text, #1a2332);
  font-size: 16px;
  font-weight: 650;
  letter-spacing: -0.02em;
}

.cate-group__copy p,
.cate-kid p,
.cate-path p {
  margin: 4px 0 0;
  color: var(--ctl-text-soft, #6b7280);
  font-size: 12px;
}

.cate-group__tools {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-shrink: 0;
  margin-left: auto;
}

.cate-state {
  color: var(--ctl-text-soft, #6b7280);
  font-size: 12px;
  line-height: 1;
}

.cate-kids {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(210px, 1fr));
  gap: 10px;
  margin-top: 14px;
}

.cate-kid,
.cate-path {
  display: flex;
  flex-direction: column;
  min-height: 108px;
  padding: 12px 14px;
  border: 1px solid var(--ctl-hairline, rgba(32, 42, 58, 0.08));
  border-radius: 12px;
  background: var(--ctl-field, #f6f8fb);
}

.cate-path {
  gap: 8px;
  background: var(--ctl-glass-strong, #fff);
}

.cate-kid__top {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 8px;
}

.cate-kid h3,
.cate-path h3 {
  font-size: 14px;
}

.cate-kid__actions {
  display: flex;
  gap: 12px;
  margin-top: auto;
  padding-top: 8px;
}

.cate-kid__actions button,
.cate-add {
  border: 0;
  background: transparent;
  color: #4285f4;
  font-size: 13px;
  cursor: pointer;
  padding: 0;
}

.cate-kid__actions .is-danger {
  color: #d14343;
}

.cate-add {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
  min-height: 108px;
  border: 1px dashed rgba(66, 133, 244, 0.45);
  border-radius: 12px;
  background: transparent;
}

.cate-add--path {
  min-height: 108px;
}

.cate-empty {
  padding: 48px 16px;
  text-align: center;
  color: var(--ctl-text-soft, #6b7280);
  border: 1px dashed var(--ctl-hairline, rgba(32, 42, 58, 0.12));
  border-radius: 16px;
}
</style>
