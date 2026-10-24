/**********************************
 * @Author: Ronnie Zhang
 * @LastEditor: Ronnie Zhang
 * @LastEditTime: 2023/12/05 21:25:07
 * @Email: zclzone@outlook.com
 * Copyright © 2023 Ronnie Zhang(大脸怪) | https://isme.top
 **********************************/

import { useAuthStore, usePermissionStore, useSysStore, useUserStore } from '@/store'
import { getPermissions, getUserInfo } from '@/store/helper'

const WHITE_LIST = ['/login', '/404']

export function createPermissionGuard(router) {
  router.beforeEach(async (to) => {
    const authStore = useAuthStore()
    const sysStore = useSysStore()
    const token = authStore.accessToken

    if (!sysStore.loaded) {
      try {
        await sysStore.fetchSysInfo()
      }
      catch {}
    }

    /** 没有token */
    if (!token) {
      if (WHITE_LIST.includes(to.path))
        return true
      return { path: 'login', query: { ...to.query, redirect: to.path } }
    }

    // 有token的情况
    if (to.path === '/login')
      return { path: '/' }
    if (WHITE_LIST.includes(to.path))
      return true

    const userStore = useUserStore()
    const permissionStore = usePermissionStore()
    if (!userStore.userInfo) {
      const tasks = [getUserInfo(), getPermissions()]
      if (!sysStore.loaded) {
        tasks.push(sysStore.fetchSysInfo())
      }
      const [user, permissions] = await Promise.all(tasks)
      userStore.setUser(user)
      permissionStore.setPermissions(permissions)
      const routeComponents = import.meta.glob('@/views/**/*.vue')
      permissionStore.accessRoutes.forEach((route) => {
        route.component = routeComponents[route.component] || undefined
        !router.hasRoute(route.name) && router.addRoute(route)
      })
      return { ...to, replace: true }
    }

    const routes = router.getRoutes()
    if (routes.find(route => route.name === to.name))
      return true
    // 否则就是无权限
    console.log('routes', routes)
    console.log('to', to)
    return { name: '403', query: { path: to.fullPath }, state: { from: 'permission-guard' } }
  })
}
