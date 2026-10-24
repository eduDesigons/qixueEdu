<template>
  <div
    ref="trackRef"
    class="view-toggle"
    :class="{ 'is-dragging': dragging, 'is-reduced': reduced }"
    :style="{ '--p': String(progress) }"
    role="group"
    aria-label="视图切换"
    @pointerdown="onPointerDown"
    @pointermove="onPointerMove"
    @pointerup="onPointerUp"
    @pointercancel="onPointerUp"
  >
    <span class="view-toggle__thumb" aria-hidden="true">
      <svg class="view-toggle__glyph" viewBox="0 0 16 16" width="14" height="14">
        <g class="view-toggle__grid">
          <rect x="1" y="1" width="6" height="6" rx="1.3" />
          <rect x="9" y="1" width="6" height="6" rx="1.3" />
          <rect x="1" y="9" width="6" height="6" rx="1.3" />
          <rect x="9" y="9" width="6" height="6" rx="1.3" />
        </g>
        <g class="view-toggle__list">
          <rect x="1" y="2.1" width="14" height="2.15" rx="1" />
          <rect x="1" y="6.92" width="14" height="2.15" rx="1" />
          <rect x="1" y="11.75" width="14" height="2.15" rx="1" />
        </g>
      </svg>
    </span>
    <button
      type="button"
      class="view-toggle__hit"
      title="卡片视图"
      :aria-pressed="appStore.listView === 'card'"
      @click="onHit('card')"
    />
    <button
      type="button"
      class="view-toggle__hit"
      title="列表视图"
      :aria-pressed="appStore.listView === 'table'"
      @click="onHit('table')"
    />
  </div>
</template>

<script setup>
import { onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useAppStore } from '@/store'

const appStore = useAppStore()
const trackRef = ref(null)
const progress = ref(appStore.listView === 'table' ? 1 : 0)
const dragging = ref(false)
const reduced = ref(false)

let startX = 0
let startP = 0
let moved = 0
let suppressClick = false
let motionQuery = null

function clamp(n, min, max) {
  return Math.min(max, Math.max(min, n))
}

function syncReduced() {
  reduced.value = !!motionQuery?.matches
}

function commit(mode) {
  progress.value = mode === 'table' ? 1 : 0
  appStore.setListView(mode)
}

function onPointerDown(event) {
  if (event.button !== 0 || reduced.value)
    return
  const track = trackRef.value
  if (!track)
    return
  dragging.value = true
  startX = event.clientX
  startP = progress.value
  moved = 0
  track.setPointerCapture?.(event.pointerId)
}

function onPointerMove(event) {
  if (!dragging.value)
    return
  const track = trackRef.value
  if (!track)
    return
  const dx = event.clientX - startX
  moved = Math.max(moved, Math.abs(dx))
  const range = Math.max(1, track.clientWidth * 0.5)
  progress.value = clamp(startP + dx / range, 0, 1)
}

function onPointerUp(event) {
  if (!dragging.value)
    return
  dragging.value = false
  suppressClick = true
  if (moved < 5) {
    const track = trackRef.value
    if (!track)
      return
    const rect = track.getBoundingClientRect()
    const x = event.clientX - rect.left
    commit(x < rect.width / 2 ? 'card' : 'table')
    return
  }
  commit(progress.value > 0.5 ? 'table' : 'card')
}

function onHit(mode) {
  if (suppressClick) {
    suppressClick = false
    return
  }
  commit(mode)
}

watch(
  () => appStore.listView,
  (mode) => {
    if (dragging.value)
      return
    progress.value = mode === 'table' ? 1 : 0
  },
)

onMounted(() => {
  motionQuery = window.matchMedia('(prefers-reduced-motion: reduce)')
  syncReduced()
  motionQuery.addEventListener('change', syncReduced)
})

onBeforeUnmount(() => {
  motionQuery?.removeEventListener('change', syncReduced)
})
</script>

<style scoped>
.view-toggle {
  --vt-ease: cubic-bezier(0.22, 1.35, 0.36, 1);
  position: relative;
  display: grid;
  grid-template-columns: 1fr 1fr;
  align-items: center;
  width: 64px;
  height: 28px;
  padding: 2px;
  border-radius: 9px;
  border: 1px solid rgba(var(--primary-color), 0.1);
  background:
    radial-gradient(28px 20px at 80% 20%, rgba(var(--primary-color), 0.08), transparent 70%),
    linear-gradient(180deg, #f4f7fb 0%, #eef2f8 100%);
  box-shadow: inset 0 1px 0 #fff;
  overflow: hidden;
  touch-action: none;
  user-select: none;
}

.view-toggle__thumb {
  position: absolute;
  top: 2px;
  left: 2px;
  z-index: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  width: calc(50% - 2px);
  height: calc(100% - 4px);
  border-radius: 7px;
  background: #fff;
  color: rgb(var(--primary-color));
  box-shadow:
    0 1px 2px rgb(31 42 43 / 8%),
    0 4px 10px rgb(31 42 43 / 6%),
    inset 0 1px 0 #fff;
  transform: translateX(calc(var(--p) * 100%));
  transition: transform 0.35s var(--vt-ease);
  pointer-events: none;
}

.view-toggle.is-dragging .view-toggle__thumb {
  transition: none;
  box-shadow:
    0 2px 8px rgba(var(--primary-color), 0.2),
    0 6px 16px rgb(31 42 43 / 10%),
    inset 0 1px 0 #fff;
}

.view-toggle.is-reduced .view-toggle__thumb,
.view-toggle.is-reduced .view-toggle__grid,
.view-toggle.is-reduced .view-toggle__list {
  transition: none;
}

.view-toggle__glyph {
  display: block;
  overflow: visible;
}

.view-toggle__grid,
.view-toggle__list {
  transform-box: fill-box;
  transform-origin: center;
  fill: currentColor;
}

.view-toggle__grid {
  opacity: calc(1 - var(--p));
  transform: translateX(calc(var(--p) * -1.4px)) scaleX(calc(1 + var(--p) * 0.18)) scaleY(calc(1 - var(--p) * 0.42));
  transition: opacity 0.2s linear, transform 0.2s linear;
}

.view-toggle__list {
  opacity: var(--p);
  transform: translateX(calc((1 - var(--p)) * 1.4px)) scaleX(calc(0.78 + var(--p) * 0.22));
  transition: opacity 0.2s linear, transform 0.2s linear;
}

.view-toggle.is-dragging .view-toggle__grid,
.view-toggle.is-dragging .view-toggle__list {
  transition: none;
}

.view-toggle__hit {
  position: relative;
  z-index: 2;
  width: 100%;
  height: 100%;
  border: 0;
  padding: 0;
  background: transparent;
  cursor: pointer;
  appearance: none;
}

@media (prefers-reduced-motion: reduce) {
  .view-toggle__thumb,
  .view-toggle__grid,
  .view-toggle__list {
    transition: none;
  }
}
</style>
