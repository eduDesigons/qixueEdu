/**********************************
 * @Author: Ronnie Zhang
 * @LastEditor: Ronnie Zhang
 * @LastEditTime: 2023/12/05 21:25:00
 * @Email: zclzone@outlook.com
 * Copyright © 2023 Ronnie Zhang(大脸怪) | https://isme.top
 **********************************/

import { useSysStore } from '@/store'

function resolveSiteTitle() {
  const name = useSysStore().sysConfig?.name
  return name && String(name).trim() ? String(name).trim() : ''
}

export function createPageTitleGuard(router) {
  router.afterEach((to) => {
    const siteTitle = resolveSiteTitle()
    const pageTitle = to.meta?.title
    if (pageTitle && siteTitle)
      document.title = `${pageTitle} | ${siteTitle}`
    else
      document.title = pageTitle || siteTitle || ''
  })
}
