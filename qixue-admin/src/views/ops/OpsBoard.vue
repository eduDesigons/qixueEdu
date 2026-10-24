<script setup>
import { BarChart, LineChart, PieChart } from 'echarts/charts'
import { GridComponent, LegendComponent, TooltipComponent } from 'echarts/components'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { NSpin, NTabPane, NTabs } from 'naive-ui'
import { computed, onMounted, ref, watch } from 'vue'
import VChart from 'vue-echarts'
import api from '@/api'
import { useRefreshOnActivate } from '@/composables'

defineOptions({ name: 'OpsBoard' })

const emit = defineEmits(['show-table'])

const props = defineProps({
  scene: { type: String, required: true },
  cityId: { type: [Number, String], default: null },
})

use([CanvasRenderer, LineChart, BarChart, PieChart, GridComponent, TooltipComponent, LegendComponent])

const loading = ref(true)
const board = ref(null)
const view = ref('chart')

function onTab(value) {
  view.value = value
  if (value === 'table')
    emit('show-table')
}

const tooltip = {
  backgroundColor: 'rgba(31, 42, 43, 0.92)',
  borderWidth: 0,
  extraCssText: 'border-radius: 12px; padding: 8px 12px; box-shadow: 0 10px 24px rgba(0,0,0,.18);',
  textStyle: { color: '#fff', fontSize: 12 },
}

function themeRgb() {
  if (typeof document === 'undefined')
    return '49, 108, 114'
  return getComputedStyle(document.documentElement).getPropertyValue('--primary-color').trim() || '49, 108, 114'
}

function themeColor(alpha) {
  const rgb = themeRgb()
  return alpha == null ? `rgb(${rgb})` : `rgba(${rgb}, ${alpha})`
}

function load() {
  loading.value = true
  api.ops.chart({
    scene: props.scene,
    city_id: props.cityId || '',
  }).then((res) => {
    board.value = res.data || null
  }).catch(() => {
    board.value = null
  }).finally(() => {
    loading.value = false
  })
}

onMounted(load)
useRefreshOnActivate(load)
watch(() => props.cityId, load)

const slices = computed(() => board.value?.slices || [])
const rankItems = computed(() => board.value?.rank?.items || [])

const trendOption = computed(() => {
  const trend = board.value?.trend
  if (!trend?.dates?.length)
    return null
  const primary = themeColor()
  return {
    tooltip: { ...tooltip, trigger: 'axis' },
    grid: { left: 4, right: 12, bottom: 4, top: 16, containLabel: true },
    xAxis: {
      type: 'category',
      boundaryGap: false,
      data: trend.dates,
      axisLine: { show: false },
      axisTick: { show: false },
      axisLabel: { color: '#8a9697', fontSize: 11 },
    },
    yAxis: {
      type: 'value',
      minInterval: 1,
      splitLine: { lineStyle: { color: themeColor(0.08), type: 'dashed' } },
      axisLabel: { color: '#8a9697', fontSize: 11 },
    },
    series: [{
      name: trend.name,
      type: 'line',
      smooth: true,
      showSymbol: false,
      lineStyle: { width: 2.4, color: primary },
      itemStyle: { color: primary },
      areaStyle: {
        color: {
          type: 'linear',
          x: 0,
          y: 0,
          x2: 0,
          y2: 1,
          colorStops: [
            { offset: 0, color: themeColor(0.22) },
            { offset: 1, color: themeColor(0.02) },
          ],
        },
      },
      data: trend.values,
    }],
  }
})

const pieOption = computed(() => {
  if (!slices.value.length)
    return null
  return {
    tooltip: { ...tooltip, trigger: 'item' },
    legend: {
      bottom: 0,
      itemWidth: 8,
      itemHeight: 8,
      textStyle: { color: '#6b7879', fontSize: 11 },
    },
    color: [themeColor(), '#c9972a', '#5d6d82', '#8b9bb4', '#7eb8bd', '#d17a45', '#6f8f72', '#b08968'],
    series: [{
      type: 'pie',
      radius: ['46%', '68%'],
      center: ['50%', '44%'],
      label: { show: false },
      data: slices.value,
    }],
  }
})

const rankOption = computed(() => {
  const items = [...rankItems.value].reverse()
  if (!items.length)
    return null
  return {
    tooltip: { ...tooltip, trigger: 'axis' },
    grid: { left: 4, right: 24, bottom: 4, top: 8, containLabel: true },
    xAxis: {
      type: 'value',
      minInterval: 1,
      splitLine: { lineStyle: { color: themeColor(0.08), type: 'dashed' } },
      axisLabel: { color: '#8a9697', fontSize: 11 },
    },
    yAxis: {
      type: 'category',
      data: items.map(item => item.name),
      axisLine: { show: false },
      axisTick: { show: false },
      axisLabel: { color: '#5c6b6c', fontSize: 12, width: 120, overflow: 'truncate' },
    },
    series: [{
      type: 'bar',
      barMaxWidth: 14,
      itemStyle: { color: themeColor(0.85), borderRadius: [0, 6, 6, 0] },
      data: items.map(item => item.value),
    }],
  }
})

const rankHeight = computed(() => `${Math.max(168, rankItems.value.length * 28 + 36)}px`)
</script>

<template>
  <div class="ops-page">
    <NTabs :value="view" type="line" class="ops-tabs" @update:value="onTab">
      <NTabPane name="chart" tab="图表" />
      <NTabPane name="table" tab="明细" />
    </NTabs>
    <NSpin v-if="view === 'chart'" :show="loading" class="ops-spin">
      <section class="ops-board">
        <div class="ops-stats">
          <article v-for="item in board?.stats || []" :key="item.label" class="ops-stat">
            <span>{{ item.label }}</span>
            <strong>{{ Number(item.value || 0).toLocaleString() }}</strong>
          </article>
        </div>
        <div class="ops-charts">
          <article class="ops-card" :class="{ 'is-solo': !pieOption }">
            <header>
              <h3>{{ board?.trend?.name || '趋势' }}</h3>
              <p>近 14 天</p>
            </header>
            <VChart v-if="trendOption" class="ops-chart" :option="trendOption" autoresize />
            <div v-else class="ops-chart ops-chart--empty">暂无趋势</div>
          </article>
          <article v-if="pieOption" class="ops-card">
            <header>
              <h3>构成</h3>
            </header>
            <VChart class="ops-chart" :option="pieOption" autoresize />
          </article>
          <article v-if="rankOption" class="ops-card ops-card--rank">
            <header>
              <h3>{{ board?.rank?.name || '排行' }}</h3>
            </header>
            <VChart class="ops-chart ops-chart--rank" :option="rankOption" :style="{ height: rankHeight }" autoresize />
          </article>
        </div>
      </section>
    </NSpin>
    <div v-else class="ops-page__table">
      <slot />
    </div>
  </div>
</template>

<style scoped>
.ops-page {
  height: 100%;
  min-height: 0;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.ops-page :deep(.n-tab-pane) {
  display: none;
}

.ops-tabs {
  flex: 0 0 auto;
}

.ops-spin {
  flex: 1;
  min-height: 0;
}

.ops-spin :deep(.n-spin-content) {
  height: 100%;
  overflow: auto;
}

.ops-board {
  display: flex;
  flex-direction: column;
  gap: 12px;
  padding: 4px 2px 12px;
}

.ops-stats {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 10px;
}

.ops-stat,
.ops-card {
  border: 1px solid var(--ctl-hairline, rgba(32, 42, 58, 0.08));
  border-radius: 16px;
  background: var(--ctl-glass-strong, #fff);
}

.ops-stat {
  padding: 12px 14px;
}

.ops-stat span,
.ops-card header p {
  color: var(--ctl-text-soft, #6b7280);
  font-size: 12px;
}

.ops-stat strong {
  display: block;
  margin-top: 4px;
  color: var(--ctl-text, #1a2332);
  font-size: 22px;
  font-weight: 650;
  letter-spacing: -0.03em;
}

.ops-charts {
  display: grid;
  grid-template-columns: minmax(0, 1.7fr) minmax(240px, 0.9fr);
  gap: 12px;
}

.ops-card {
  min-width: 0;
  padding: 12px 14px 8px;
}

.ops-card.is-solo,
.ops-card--rank {
  grid-column: 1 / -1;
}

.ops-card header {
  display: flex;
  align-items: baseline;
  justify-content: space-between;
  gap: 8px;
}

.ops-card h3,
.ops-card header p {
  margin: 0;
}

.ops-card h3 {
  color: var(--ctl-text, #1a2332);
  font-size: 14px;
  font-weight: 650;
}

.ops-chart {
  width: 100%;
  height: 280px;
}

.ops-chart--rank {
  height: auto;
  min-height: 168px;
}

.ops-chart--empty {
  display: flex;
  align-items: center;
  justify-content: center;
  color: #8a9697;
  font-size: 13px;
}

.ops-page__table {
  flex: 1;
  min-height: 0;
}

@media (max-width: 960px) {
  .ops-stats,
  .ops-charts {
    grid-template-columns: 1fr;
  }
}
</style>
