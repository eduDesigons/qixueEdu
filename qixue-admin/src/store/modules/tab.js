/**********************************
 * @Author: Ronnie Zhang
 * @LastEditor: Ronnie Zhang
 * @LastEditTime: 2023/12/05 21:25:52
 * @Email: zclzone@outlook.com
 * Copyright © 2023 Ronnie Zhang(大脸怪) | https://isme.top
 **********************************/

import { defineStore } from 'pinia'
import { PAGE_CACHE_TTL_MS, isCacheExpired } from '@/utils/pageCache'
import { useRouterStore } from './router'

export const useTabStore = defineStore('tab', {
  state: () => ({
    tabs: [],
    activeTab: '',
    reloading: false,
    /** 访问顺序，关闭标签时回到上一页 */
    visitStack: [],
  }),
  getters: {
    activeIndex() {
      return this.tabs.findIndex(item => item.path === this.activeTab)
    },
  },
  actions: {
    pushVisit(path) {
      if (!path)
        return
      const stack = this.visitStack
      if (stack[stack.length - 1] === path)
        return
      stack.push(path)
      if (stack.length > 40)
        this.visitStack = stack.slice(-30)
    },
    touchTab(path) {
      const tab = this.tabs.find(item => item.path === path)
      if (!tab)
        return
      tab.lastActiveAt = Date.now()
    },
    isTabCacheFresh(tab) {
      if (!tab?.keepAlive)
        return false
      if (tab.path === this.activeTab)
        return true
      const last = tab.lastActiveAt || tab.refreshAt || 0
      return !isCacheExpired(last, PAGE_CACHE_TTL_MS)
    },
    /** 过期的 KeepAlive 页 bump refreshAt，下次进入重建，释放旧实例 */
    pruneExpiredPageCache() {
      const now = Date.now()
      this.tabs.forEach((tab) => {
        if (!tab.keepAlive || tab.path === this.activeTab)
          return
        const last = tab.lastActiveAt || tab.refreshAt || 0
        if (!last || now - last <= PAGE_CACHE_TTL_MS)
          return
        tab.refreshAt = now
        tab.lastActiveAt = 0
      })
    },
    async setActiveTab(path) {
      await nextTick()
      this.pruneExpiredPageCache()
      this.activeTab = path
      this.touchTab(path)
      this.pushVisit(path)
    },
    setTabs(tabs) {
      this.tabs = tabs
      const alive = new Set(tabs.map(item => item.path))
      this.visitStack = this.visitStack.filter(path => alive.has(path))
    },
    /** 只关标签，不跳转（详情页返回配合 history.back） */
    closeTabSilent(path) {
      if (!path)
        return
      this.tabs = this.tabs.filter(tab => tab.path !== path)
      this.visitStack = this.visitStack.filter(item => item !== path)
      if (this.activeTab === path)
        this.activeTab = this.visitStack[this.visitStack.length - 1] || this.tabs[this.tabs.length - 1]?.path || ''
    },
    addTab(tab = {}) {
      const findIndex = this.tabs.findIndex(item => item.path === tab.path)
      const now = Date.now()
      if (findIndex !== -1) {
        const prev = this.tabs[findIndex]
        const expired = prev.keepAlive && isCacheExpired(prev.lastActiveAt || prev.refreshAt, PAGE_CACHE_TTL_MS)
        this.tabs.splice(findIndex, 1, {
          ...prev,
          ...tab,
          refreshAt: expired ? now : prev.refreshAt,
          lastActiveAt: now,
        })
      }
      else {
        this.setTabs([...this.tabs, { ...tab, refreshAt: now, lastActiveAt: now }])
      }
      this.setActiveTab(tab.path)
    },
    async reloadTab(path) {
      const findItem = this.tabs.find(item => item.path === path)
      if (!findItem)
        return
      findItem.refreshAt = Date.now()
      findItem.lastActiveAt = Date.now()
      $loadingBar.start()
      this.reloading = true
      await nextTick()
      this.reloading = false
      setTimeout(() => {
        document.documentElement.scrollTo({ left: 0, top: 0 })
        $loadingBar.finish()
      }, 100)
    },
    resolvePreviousTab(excludePath) {
      const remaining = this.tabs.filter(tab => tab.path !== excludePath)
      const alive = new Set(remaining.map(item => item.path))

      for (let i = this.visitStack.length - 1; i >= 0; i--) {
        const path = this.visitStack[i]
        if (path !== excludePath && alive.has(path))
          return path
      }

      const detailMatch = excludePath.match(/^(.*)\/(?:detail|edit)(?:\/|$)/)
      if (detailMatch) {
        const base = detailMatch[1]
        const listPath = remaining.find(tab =>
          tab.path === `${base}/list` || tab.path === base || tab.path.startsWith(`${base}/`),
        )
        if (listPath)
          return listPath.path
        const courseList = remaining.find(tab => tab.path === '/course/list')
        if (courseList && excludePath.startsWith('/course/'))
          return courseList.path
        const pathList = remaining.find(tab => tab.path === '/course/path')
        if (pathList && excludePath.startsWith('/course/path'))
          return pathList.path
      }

      const chapterMatch = excludePath.match(/^(\/course\/\d+)\/chapter/)
      if (chapterMatch) {
        const courseDetail = remaining.find(tab => tab.path.startsWith('/course/detail/'))
        if (courseDetail)
          return courseDetail.path
      }

      return remaining[remaining.length - 1]?.path || ''
    },
    async removeTab(path) {
      const nextPath = path === this.activeTab ? this.resolvePreviousTab(path) : ''
      const remaining = this.tabs.filter(tab => tab.path !== path)
      this.visitStack = this.visitStack.filter(item => item !== path)

      if (path === this.activeTab && nextPath) {
        const failure = await useRouterStore().router?.push(nextPath)
        if (failure)
          return
      }
      this.setTabs(remaining)
    },
    removeOther(curPath = this.activeTab) {
      this.setTabs(this.tabs.filter(tab => tab.path === curPath))
      this.visitStack = this.visitStack.filter(path => path === curPath)
      if (curPath !== this.activeTab)
        useRouterStore().router?.push(curPath)
    },
    removeLeft(curPath) {
      const curIndex = this.tabs.findIndex(item => item.path === curPath)
      const filterTabs = this.tabs.filter((item, index) => index >= curIndex)
      this.setTabs(filterTabs)
      if (!filterTabs.find(item => item.path === this.activeTab))
        useRouterStore().router?.push(filterTabs[filterTabs.length - 1].path)
    },
    removeRight(curPath) {
      const curIndex = this.tabs.findIndex(item => item.path === curPath)
      const filterTabs = this.tabs.filter((item, index) => index <= curIndex)
      this.setTabs(filterTabs)
      if (!filterTabs.find(item => item.path === this.activeTab))
        useRouterStore().router?.push(filterTabs[filterTabs.length - 1].path)
    },
    resetTabs() {
      this.$reset()
    },
  },
  persist: {
    pick: ['tabs', 'visitStack', 'activeTab'],
    storage: sessionStorage,
  },
})
