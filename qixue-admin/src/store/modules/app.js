/**********************************
 * @Author: Ronnie Zhang
 * @LastEditor: Ronnie Zhang
 * @LastEditTime: 2023/12/05 21:25:31
 * @Email: zclzone@outlook.com
 * Copyright © 2023 Ronnie Zhang(大脸怪) | https://isme.top
 **********************************/

import { defineStore } from 'pinia'

export const useAppStore = defineStore('app', {
  state: () => ({
    collapsed: false,
    listView: 'card',
  }),
  actions: {
    switchCollapsed() {
      this.collapsed = !this.collapsed
    },
    setCollapsed(b) {
      this.collapsed = b
    },
    toggleListView() {
      this.listView = this.listView === 'card' ? 'table' : 'card'
    },
    setListView(mode) {
      if (mode !== 'card' && mode !== 'table')
        return
      this.listView = mode
    },
  },
  persist: [
    {
      pick: ['collapsed'],
      storage: sessionStorage,
    },
    {
      pick: ['listView'],
      storage: localStorage,
    },
  ],
})
