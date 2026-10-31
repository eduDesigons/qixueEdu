<script setup>
import { NButton, NTabPane } from 'naive-ui'
import { computed, onActivated, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import api from '@/api'
import { CommonPage, DetailTabs } from '@/components/index.js'
import CourseApplyPanel from './CourseApplyPanel.vue'
import PathApplyPanel from './PathApplyPanel.vue'

defineOptions({ name: 'MemberCourseApply' })

const TABS = ['course', 'path']
const route = useRoute()
const router = useRouter()
const coursePanel = ref(null)
const pathPanel = ref(null)
const courseChecked = ref([])
const pathChecked = ref([])
const coursePending = ref(0)
const pathPending = ref(0)
const activeTab = ref(TABS.includes(route.query.tab) ? route.query.tab : 'course')

const currentChecked = computed(() =>
  activeTab.value === 'path' ? pathChecked.value : courseChecked.value,
)
const currentPanel = computed(() =>
  activeTab.value === 'path' ? pathPanel.value : coursePanel.value,
)
const batchLoading = computed(() => {
  const loading = currentPanel.value?.batchLoading
  return Boolean(loading?.value ?? loading)
})

function loadPending() {
  api.home.getStats().then((res) => {
    const data = res.data || {}
    coursePending.value = Number(data.apply_pending_count) || 0
    pathPending.value = Number(data.path_apply_pending_count) || 0
  }).catch(() => {})
}

watch(() => route.query.tab, (tab) => {
  if (TABS.includes(tab) && tab !== activeTab.value)
    activeTab.value = tab
})

watch(activeTab, (tab) => {
  if (route.query.tab === tab)
    return
  router.replace({
    query: {
      ...route.query,
      tab,
    },
  })
})

function batchApprove() {
  currentPanel.value?.batchApprove()
}

function batchReject() {
  currentPanel.value?.batchReject()
}

onMounted(loadPending)
onActivated(() => {
  const tab = route.query.tab
  if (TABS.includes(tab))
    activeTab.value = tab
  loadPending()
})
</script>

<template>
  <CommonPage title="申请">
    <DetailTabs v-model:value="activeTab">
      <template #suffix>
        <NButton
          v-if="currentChecked.length"
          :disabled="batchLoading"
          :loading="batchLoading"
          type="primary"
          secondary
          @click="batchApprove"
        >
          批量通过
        </NButton>
        <NButton
          v-if="currentChecked.length"
          :disabled="batchLoading"
          :loading="batchLoading"
          type="warning"
          secondary
          @click="batchReject"
        >
          批量拒绝
        </NButton>
      </template>
      <NTabPane name="course" display-directive="show:lazy">
        <template #tab>
          <span class="apply-tab">
            入课申请
            <em v-if="coursePending">{{ coursePending }}</em>
          </span>
        </template>
        <CourseApplyPanel
          ref="coursePanel"
          @checked="courseChecked = $event"
          @changed="loadPending"
        />
      </NTabPane>
      <NTabPane name="path" display-directive="show:lazy">
        <template #tab>
          <span class="apply-tab">
            入系申请
            <em v-if="pathPending">{{ pathPending }}</em>
          </span>
        </template>
        <PathApplyPanel
          ref="pathPanel"
          @checked="pathChecked = $event"
          @changed="loadPending"
        />
      </NTabPane>
    </DetailTabs>
  </CommonPage>
</template>

<style scoped>
.apply-tab {
  display: inline-flex;
  align-items: center;
  gap: 6px;
}

.apply-tab em {
  min-width: 16px;
  height: 16px;
  padding: 0 5px;
  border-radius: 999px;
  background: rgb(var(--primary-color));
  color: #fff;
  font-size: 11px;
  font-style: normal;
  font-weight: 650;
  line-height: 16px;
  text-align: center;
  font-variant-numeric: tabular-nums;
}
</style>
