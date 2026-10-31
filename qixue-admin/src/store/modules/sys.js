/**********************************
 * @FilePath: sys.js
 **********************************/

import { defineStore } from 'pinia'
import api from '@/api'
import { applySiteBranding } from '@/utils/siteBranding'
import { request } from '@/utils'
import { useAuthStore } from './auth'

export const useSysStore = defineStore('sys', {
  state: () => ({
    sysConfig: {},
    playerConfig: {},
    coinConfig: {},
    loaded: false,
  }),
  getters: {
    assetDomain: state => state.sysConfig?.assetDomain || '',
  },
  actions: {
    setSysInfo(data = {}) {
      this.sysConfig = data.sysConfig || {}
      this.playerConfig = data.playerConfig || {}
      this.coinConfig = data.coinConfig || {}
      this.loaded = true
      applySiteBranding({ sysConfig: this.sysConfig })
    },
    async fetchSysInfo() {
      const token = useAuthStore().accessToken
      const res = token
        ? await api.common.getSysInfo()
        : await request.get('/webapp/sys/getSysInfo', { needToken: false, needTip: false })
      this.setSysInfo(res.data)
      return res.data
    },
  },
})
