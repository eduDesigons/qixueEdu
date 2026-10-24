import {
  Icon, Loading, Empty, Badge, Image as VanImage, Popup, Tab, Tabs,
  Swipe, SwipeItem, Collapse, CollapseItem, NavBar, Circle, Search, Picker, Dialog,
} from 'vant'
import 'vant/es/icon/style'
import 'vant/es/loading/style'
import 'vant/es/empty/style'
import 'vant/es/badge/style'
import 'vant/es/image/style'
import 'vant/es/popup/style'
import 'vant/es/tab/style'
import 'vant/es/tabs/style'
import 'vant/es/swipe/style'
import 'vant/es/swipe-item/style'
import 'vant/es/collapse/style'
import 'vant/es/collapse-item/style'
import 'vant/es/nav-bar/style'
import 'vant/es/circle/style'
import 'vant/es/search/style'
import 'vant/es/picker/style'
import 'vant/es/dialog/style'
import Mescroll from './components/Mescroll.vue'

const HOOKS = ['onLoad', 'onShow', 'onReady', 'onHide', 'onUnload', 'onReachBottom', 'onPullDownRefresh', 'onPageScroll']

function callPageHook(vm, name, arg) {
  const raw = vm.$options[name]
  const list = !raw ? [] : Array.isArray(raw) ? raw : [raw]
  list.forEach((fn) => fn.call(vm, arg))
}

/** 路由页组件：__qxPage 标记，或组件 name 与当前路由 name 一致（防 HMR 丢标记） */
function isQxPage(vm) {
  if (vm.$options && vm.$options.__qxPage) return true
  const type = vm.$ && vm.$.type
  if (type && type.__qxPage) return true
  const routeName = vm.$route && vm.$route.name
  const typeName = type && type.name
  return !!(routeName && typeName && routeName === typeName)
}

/** 注册 Vant、列表滚动组件，以及页面 onLoad/onShow 生命周期 */
export function installApp(app) {
  const chain = (parent, child) => {
    const list = (value) => (!value ? [] : Array.isArray(value) ? value : [value])
    return list(parent).concat(list(child))
  }
  HOOKS.forEach((name) => {
    app.config.optionMergeStrategies[name] = chain
  })

  ;[
    Icon, Loading, Empty, Badge, VanImage, Popup, Tab, Tabs,
    Swipe, SwipeItem, Collapse, CollapseItem, NavBar, Circle, Search, Picker, Dialog,
  ].forEach((c) => app.use(c))

  app.component('mescroll-body', Mescroll)
  app.component('mescroll-uni', Mescroll)

  app.mixin({
    methods: {
      $set(target, key, value) {
        if (Array.isArray(target)) target.splice(key, 1, value)
        else if (target && typeof target === 'object') target[key] = value
        return value
      },
      $delete(target, key) {
        if (Array.isArray(target)) target.splice(key, 1)
        else if (target) delete target[key]
      },
      __callPageHook(name, arg) {
        callPageHook(this, name, arg)
      },
      __bindPageScroll() {
        this.__unbindPageScroll()
        let startY = 0
        this.__onScroll = () => {
          const gap = document.documentElement.scrollHeight - window.scrollY - window.innerHeight
          if (gap < 80) {
            if (!this.__qxBottomLock) {
              this.__qxBottomLock = true
              this.__callPageHook('onReachBottom')
              setTimeout(() => { this.__qxBottomLock = false }, 600)
            }
          } else {
            this.__qxBottomLock = false
          }
          this.__callPageHook('onPageScroll', { scrollTop: window.scrollY })
        }
        this.__onTouchStart = (event) => {
          startY = event.touches ? event.touches[0].clientY : 0
        }
        this.__onTouchEnd = (event) => {
          const y = event.changedTouches ? event.changedTouches[0].clientY : 0
          if (window.scrollY <= 0 && y - startY > 70) this.__callPageHook('onPullDownRefresh')
        }
        window.addEventListener('scroll', this.__onScroll, { passive: true })
        window.addEventListener('touchstart', this.__onTouchStart, { passive: true })
        window.addEventListener('touchend', this.__onTouchEnd, { passive: true })
      },
      __unbindPageScroll() {
        if (this.__onScroll) window.removeEventListener('scroll', this.__onScroll)
        if (this.__onTouchStart) window.removeEventListener('touchstart', this.__onTouchStart)
        if (this.__onTouchEnd) window.removeEventListener('touchend', this.__onTouchEnd)
      },
    },
    mounted() {
      if (!isQxPage(this)) return
      this.__qxSkipActivate = true
      const query = this.$route ? { ...this.$route.query, ...this.$route.params } : {}
      this.__callPageHook('onLoad', query)
      this.__callPageHook('onShow')
      this.__callPageHook('onReady')
      this.__bindPageScroll()
    },
    activated() {
      if (!isQxPage(this)) return
      if (this.__qxSkipActivate) {
        this.__qxSkipActivate = false
        return
      }
      this.__callPageHook('onShow')
      this.__bindPageScroll()
    },
    deactivated() {
      if (!isQxPage(this)) return
      this.__callPageHook('onHide')
      this.__unbindPageScroll()
    },
    beforeUnmount() {
      if (!isQxPage(this)) return
      this.__callPageHook('onHide')
      this.__callPageHook('onUnload')
      this.__unbindPageScroll()
    },
  })
}
