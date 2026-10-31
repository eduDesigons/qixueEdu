const MescrollMixin = {
  data() {
    return { mescroll: null }
  },
  onPullDownRefresh() {
    this.mescroll && this.mescroll.onPullDownRefresh && this.mescroll.onPullDownRefresh()
  },
  onPageScroll(e) {
    this.mescroll && this.mescroll.onPageScroll && this.mescroll.onPageScroll(e)
  },
  onReachBottom() {
    this.mescroll && this.mescroll.onReachBottom && this.mescroll.onReachBottom()
  },
  methods: {
    mescrollInit(mescroll) {
      this.mescroll = mescroll
    },
    downCallback() {
      if (this.mescroll && this.mescroll.optUp && this.mescroll.optUp.use) this.mescroll.resetUpScroll()
      else if (this.mescroll) this.mescroll.endSuccess()
    },
    upCallback() {
      this.mescroll && this.mescroll.endErr()
    },
  },
}

export default MescrollMixin
