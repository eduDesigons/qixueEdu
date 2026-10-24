/**********************************
 * @Author: Ronnie Zhang
 * @LastEditor: Ronnie Zhang
 * @LastEditTime: 2023/12/05 21:24:53
 * @Email: zclzone@outlook.com
 * Copyright © 2023 Ronnie Zhang(大脸怪) | https://isme.top
 **********************************/

/** 瞬时跳转不闪进度条；慢一点再出现，更接近 Safari 的克制感 */
const SHOW_DELAY = 140

export function createPageLoadingGuard(router) {
  let showTimer = null
  let started = false

  function clearShowTimer() {
    if (showTimer != null) {
      clearTimeout(showTimer)
      showTimer = null
    }
  }

  function settle(failed) {
    clearShowTimer()
    if (!started)
      return
    started = false
    requestAnimationFrame(() => {
      if (failed)
        $loadingBar.error()
      else
        $loadingBar.finish()
    })
  }

  router.beforeEach(() => {
    // 进度条本身只走到 80%，必须 finish 才会到头。
    // 权限守卫会再 replace 一次。若这里把已出现的进度条标记清掉，
    // 最后那次很快结束的跳转不会调用 finish，条子就停在后面。
    if (started)
      return
    clearShowTimer()
    showTimer = setTimeout(() => {
      showTimer = null
      started = true
      $loadingBar.start()
    }, SHOW_DELAY)
  })

  router.afterEach(() => {
    settle(false)
  })

  router.onError(() => {
    settle(true)
  })
}
