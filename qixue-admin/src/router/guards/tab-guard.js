/**********************************
 * @Author: Ronnie Zhang
 * @LastEditor: Ronnie Zhang
 * @LastEditTime: 2023/12/05 21:25:17
 * @Email: zclzone@outlook.com
 * Copyright © 2023 Ronnie Zhang(大脸怪) | https://isme.top
 **********************************/

import { useTabStore } from '@/store'

export const EXCLUDE_TAB = ['/404', '/403', '/login']

export function createTabGuard(router) {
  router.afterEach((to) => {
    if (EXCLUDE_TAB.includes(to.path))
      return
    const tabStore = useTabStore()
    const { name, path } = to
    const title = to.meta?.title
    const icon = to.meta?.icon
    // 跟随菜单配置；未配置时默认不缓存（避免新增表单脏数据）
    const keepAlive = !!to.meta?.keepAlive
    tabStore.addTab({ name, path, title, icon, keepAlive })
  })
}
