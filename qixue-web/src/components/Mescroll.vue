<template>
  <div class="mescroll-body" :style="boxStyle" @scroll.passive="onScroll" @touchstart="onTouchStart" @touchend="onTouchEnd">
    <slot />
    <div v-if="showEmpty" class="mescroll-empty">
      <img v-if="emptyIcon" :src="emptyIcon" alt="" />
      <div>{{ emptyTip }}</div>
    </div>
    <div v-else-if="noMore && total > 0" class="mescroll-tip">没有更多了</div>
    <div v-else-if="loading" class="mescroll-tip">加载中</div>
  </div>
</template>

<script>
export default {
  name: 'mescroll-body',
  props: {
    down: { type: Object, default: () => ({}) },
    up: { type: Object, default: () => ({}) },
    height: [Number, String],
  },
  emits: ['init', 'down', 'up', 'emptyclick'],
  data() {
    return {
      loading: false,
      noMore: false,
      showEmpty: false,
      total: 0,
      pageNum: 0,
      startY: 0,
    }
  },
  computed: {
    pageSize() {
      return (this.up && this.up.page && this.up.page.size) || 10
    },
    emptyTip() {
      return (this.up && this.up.empty && this.up.empty.tip) || '暂无数据'
    },
    emptyIcon() {
      return (this.up && this.up.empty && this.up.empty.icon) || ''
    },
    boxStyle() {
      if (this.height === undefined || this.height === '') return null
      const height = typeof this.height === 'number' ? `${this.height}px` : this.height
      return { height, overflow: 'auto' }
    },
  },
  mounted() {
    const mescroll = this.createMescroll()
    this.mescroll = mescroll
    this.$emit('init', mescroll)
    if (!this.up || this.up.auto !== false) {
      this.$nextTick(() => mescroll.triggerUpScroll())
    }
  },
  methods: {
    createMescroll() {
      const vm = this
      return {
        optUp: { use: true },
        optDown: { use: !(vm.down && vm.down.use === false) },
        resetUpScroll() {
          vm.pageNum = 0
          vm.noMore = false
          vm.showEmpty = false
          this.triggerUpScroll()
        },
        triggerUpScroll() {
          if (vm.loading) return
          vm.pageNum += 1
          vm.loading = true
          vm.$emit('up', { num: vm.pageNum, size: vm.pageSize })
        },
        triggerDownScroll() {
          vm.$emit('down')
        },
        endSuccess(size = 0, hasNext) {
          vm.loading = false
          const count = Number(size) || 0
          vm.total = vm.pageNum === 1 ? count : vm.total + count
          const more = typeof hasNext === 'boolean' ? hasNext : count >= vm.pageSize
          vm.noMore = !more
          vm.showEmpty = vm.pageNum === 1 && count === 0
        },
        endErr() {
          vm.loading = false
          if (vm.pageNum > 0) vm.pageNum -= 1
        },
        onReachBottom() {
          if (!vm.loading && !vm.noMore) this.triggerUpScroll()
        },
        onPullDownRefresh() {
          vm.$emit('down')
        },
        onPageScroll() {},
      }
    },
    onScroll(event) {
      if (!this.height) return
      const el = event.target
      if (el.scrollTop + el.clientHeight >= el.scrollHeight - 40 && this.mescroll) {
        this.mescroll.onReachBottom()
      }
    },
    onTouchStart(event) {
      this.startY = event.changedTouches ? event.changedTouches[0].clientY : 0
    },
    onTouchEnd(event) {
      const y = event.changedTouches ? event.changedTouches[0].clientY : 0
      const top = this.height ? event.currentTarget.scrollTop : window.scrollY
      if (top <= 0 && y - this.startY > 70) this.$emit('down')
    },
  },
}
</script>
