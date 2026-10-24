<script setup>
import { Add } from '@vicons/ionicons5'
import { NButton, NIcon, NSpin } from 'naive-ui'
import { computed, onMounted, ref } from 'vue'
import api from '@/api'
import { CommonPage } from '@/components/index.js'
import { useRefreshOnActivate } from '@/composables'
import EditArea from './EditArea.vue'

defineOptions({ name: 'AreaMgt' })

const loading = ref(false)
const rows = ref([])
const selectedId = ref(null)
const editor = ref(null)

const parents = computed(() => rows.value.filter(item => !item.parent_id))
const current = computed(() => rows.value.find(item => item.id === selectedId.value) || null)
const children = computed(() => current.value?.children || [])

onMounted(load)
useRefreshOnActivate(load)

function load() {
  loading.value = true
  api.area.list().then((res) => {
    rows.value = res.data || []
    if (!rows.value.some(item => item.id === selectedId.value))
      selectedId.value = rows.value[0]?.id ?? null
  }).catch(() => {}).finally(() => {
    loading.value = false
  })
}

function onSaved(payload) {
  const id = Number(payload?.id) || 0
  const parentId = Number(payload?.parent_id) || 0
  if (!parentId && id)
    selectedId.value = id
  else if (parentId)
    selectedId.value = parentId
  load()
}

function removeRow(row) {
  const blocked = row.children?.length
    ? '请先删除下面的二级区域'
    : (row.user_count ? '该区域已被使用，不能删除' : '')
  if (blocked) {
    $message.warning(blocked)
    return
  }
  $dialog.warning({
    title: '确认删除',
    content: `确定删除区域「${row.name}」吗？`,
    positiveText: '删除',
    negativeText: '取消',
    onPositiveClick: () => api.area.delete({ id: row.id }).then(() => {
      $message.success('删除成功')
      if (row.id === selectedId.value)
        selectedId.value = null
      load()
    }).catch(() => {}),
  })
}
</script>

<template>
  <CommonPage subtitle="一级是大区，二级是具体区域。编码写在一级上，二级跟一级用同一个编码。">
    <div class="area-layout">
      <NSpin :show="loading" class="area-side">
        <div class="area-side__bar">
          <span>一级区域</span>
          <NButton size="tiny" secondary type="primary" @click="editor?.show(null, 0)">
            <template #icon>
              <NIcon :component="Add" />
            </template>
            新增
          </NButton>
        </div>
        <div class="area-side__list">
          <button
            v-for="group in rows"
            :key="group.id"
            type="button"
            class="area-item"
            :class="{ 'is-active': selectedId === group.id, 'is-off': group.is_show === 0 }"
            @click="selectedId = group.id"
          >
            <strong>{{ group.name }}</strong>
            <em>{{ group.code || '—' }}</em>
          </button>
          <div v-if="!loading && !rows.length" class="area-empty">
            还没有区域
          </div>
        </div>
      </NSpin>

      <section class="area-detail">
        <template v-if="current">
          <header class="area-detail__head">
            <div>
              <h2>{{ current.name }}</h2>
              <p>
                编码 {{ current.code || '—' }}
                · {{ current.user_count || 0 }} 名学员
                · {{ current.is_show === 0 ? '隐藏' : '显示' }}
              </p>
            </div>
            <div class="area-detail__tools">
              <NButton size="small" secondary @click="editor?.show(current)">
                编辑
              </NButton>
              <NButton
                v-if="!children.length && !current.user_count"
                size="small"
                secondary
                type="error"
                @click="removeRow(current)"
              >
                删除
              </NButton>
              <NButton size="small" type="primary" @click="editor?.show(null, current.id)">
                <template #icon>
                  <NIcon :component="Add" />
                </template>
                添加二级
              </NButton>
            </div>
          </header>
          <div class="area-kids">
            <article
              v-for="child in children"
              :key="child.id"
              class="area-kid"
              :class="{ 'is-off': child.is_show === 0 }"
            >
              <div class="area-kid__main">
                <h3>{{ child.name }}</h3>
                <p>{{ child.user_count || 0 }} 名学员 · {{ child.is_show === 0 ? '隐藏' : '显示' }}</p>
              </div>
              <div class="area-kid__actions">
                <NButton size="small" secondary @click="editor?.show(child)">
                  编辑
                </NButton>
                <NButton v-if="!child.user_count" size="small" secondary type="error" @click="removeRow(child)">
                  删除
                </NButton>
              </div>
            </article>
            <div v-if="!children.length" class="area-empty">
              还没有二级区域
            </div>
          </div>
        </template>
        <div v-else class="area-empty">
          从左侧选择一个一级区域
        </div>
      </section>
    </div>
    <EditArea ref="editor" :parents="parents" @update="onSaved" />
  </CommonPage>
</template>

<style scoped>
.area-layout {
  display: flex;
  gap: 16px;
  height: 100%;
  min-height: 0;
  overflow: hidden;
}

.area-side {
  width: 240px;
  flex-shrink: 0;
  min-height: 0;
  display: flex;
  flex-direction: column;
  border: 1px solid var(--ctl-hairline, rgba(32, 42, 58, 0.08));
  border-radius: 16px;
  background: var(--ctl-glass-strong, #fff);
  overflow: hidden;
}

.area-side :deep(.n-spin-container),
.area-side :deep(.n-spin-content) {
  height: 100%;
  min-height: 0;
  display: flex;
  flex-direction: column;
}

.area-side__bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  padding: 12px 12px 8px;
  font-weight: 650;
}

.area-side__list {
  flex: 1;
  min-height: 0;
  overflow: auto;
  padding: 0 8px 12px;
}

.area-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  width: 100%;
  margin: 0;
  padding: 10px 12px;
  border: 0;
  border-radius: 10px;
  background: transparent;
  text-align: left;
  cursor: pointer;
  color: inherit;
}

.area-item strong {
  min-width: 0;
  font-weight: 600;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.area-item em {
  flex-shrink: 0;
  font-style: normal;
  font-size: 12px;
  color: var(--ctl-text-soft, #5a6b6d);
}

.area-item.is-active {
  background: rgba(66, 133, 244, 0.1);
}

.area-item.is-off,
.area-kid.is-off {
  opacity: 0.5;
}

.area-detail {
  flex: 1;
  min-width: 0;
  min-height: 0;
  overflow: auto;
  padding: 4px 4px 16px;
}

.area-detail__head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
  margin-bottom: 16px;
}

.area-detail__head h2 {
  margin: 0;
  font-size: 28px;
  font-weight: 650;
  letter-spacing: -0.03em;
}

.area-detail__head p {
  margin: 6px 0 0;
  color: var(--ctl-text-soft, #5a6b6d);
}

.area-detail__tools {
  display: flex;
  flex-wrap: wrap;
  justify-content: flex-end;
  gap: 8px;
  flex-shrink: 0;
}

.area-kids {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.area-kid {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 12px 14px;
  border: 1px solid var(--ctl-hairline, rgba(32, 42, 58, 0.08));
  border-radius: 12px;
  background: var(--ctl-glass-strong, #fff);
}

.area-kid__main {
  min-width: 0;
}

.area-kid h3 {
  margin: 0;
  font-size: 15px;
  font-weight: 650;
}

.area-kid p {
  margin: 4px 0 0;
  color: var(--ctl-text-soft, #5a6b6d);
  font-size: 13px;
}

.area-kid__actions {
  display: flex;
  gap: 8px;
  flex-shrink: 0;
}

.area-empty {
  padding: 48px 0;
  text-align: center;
  color: var(--ctl-text-soft, #5a6b6d);
}
</style>
