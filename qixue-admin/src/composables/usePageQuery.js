const tabMemory = new Map()

function memoryKey(route) {
  return route.path
}

export function peekTabQuery(defaultTab = 'content') {
  const route = useRoute()
  const fromQuery = route.query.tab
  if (typeof fromQuery === 'string' && fromQuery)
    return fromQuery
  return tabMemory.get(memoryKey(route)) || defaultTab
}

export function useTabQuery(currentTabName, defaultTab = 'content') {
  const route = useRoute()
  const router = useRouter()

  function remember(tab) {
    if (tab)
      tabMemory.set(memoryKey(route), tab)
  }

  function syncFromRouteOrMemory() {
    const tab = (typeof route.query.tab === 'string' && route.query.tab)
      || tabMemory.get(memoryKey(route))
      || defaultTab
    if (tab && currentTabName.value !== tab)
      currentTabName.value = tab
    remember(tab)
    if (tab && route.query.tab !== tab) {
      router.replace({
        query: {
          ...route.query,
          tab,
        },
      })
    }
  }

  watch(() => route.query.tab, (tab) => {
    if (typeof tab === 'string' && tab && tab !== currentTabName.value) {
      currentTabName.value = tab
      remember(tab)
    }
  })

  watch(currentTabName, (tab) => {
    if (!tab)
      return
    remember(tab)
    if (route.query.tab === tab)
      return
    router.replace({
      query: {
        ...route.query,
        tab,
      },
    })
  })

  onMounted(syncFromRouteOrMemory)
  onActivated(syncFromRouteOrMemory)
}

export function useUnsavedGuard(getValue) {
  let baseline = ''
  let armed = false
  let bypass = false

  function snapshot() {
    try {
      baseline = JSON.stringify(getValue() ?? null)
      armed = true
    }
    catch {
      armed = false
    }
  }

  function isDirty() {
    if (!armed || bypass)
      return false
    try {
      return JSON.stringify(getValue() ?? null) !== baseline
    }
    catch {
      return false
    }
  }

  function allowLeave() {
    bypass = true
  }

  function confirmLeave() {
    return new Promise((resolve) => {
      if (!isDirty()) {
        resolve(true)
        return
      }
      let settled = false
      const done = (ok) => {
        if (settled)
          return
        settled = true
        resolve(ok)
      }
      $dialog.warning({
        title: '未保存的修改',
        content: '离开后当前修改会丢失，确定离开？',
        positiveText: '离开',
        negativeText: '留下',
        onPositiveClick: () => done(true),
        onNegativeClick: () => done(false),
        onClose: () => done(false),
      })
    })
  }

  onBeforeRouteLeave(async (to, from, next) => {
    if (!isDirty()) {
      next()
      return
    }
    const ok = await confirmLeave()
    next(ok)
  })

  function onBeforeUnload(e) {
    if (!isDirty())
      return
    e.preventDefault()
    e.returnValue = ''
  }

  onMounted(() => window.addEventListener('beforeunload', onBeforeUnload))
  onActivated(() => window.addEventListener('beforeunload', onBeforeUnload))
  onDeactivated(() => window.removeEventListener('beforeunload', onBeforeUnload))
  onBeforeUnmount(() => window.removeEventListener('beforeunload', onBeforeUnload))

  return { snapshot, allowLeave, isDirty }
}
