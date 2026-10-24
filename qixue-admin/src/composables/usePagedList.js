export function usePagedList(getData, queryItems, pageSize = 12) {
  const loading = ref(false)
  const list = ref([])
  const total = ref(0)
  const page = ref(1)

  async function fetchList() {
    loading.value = true
    try {
      while (true) {
        const { data } = await getData({
          page: page.value,
          size: pageSize,
          ...queryItems.value,
        })
        list.value = data?.list || []
        total.value = data?.total ?? 0
        if (total.value && !list.value.length && page.value > 1) {
          page.value -= 1
          continue
        }
        break
      }
    }
    catch (error) {
      console.error(error)
      $message.error(error?.message || '列表加载失败')
      list.value = []
      total.value = 0
    }
    finally {
      loading.value = false
    }
  }

  function handleSearch() {
    page.value = 1
    return fetchList()
  }

  function onPageChange(next) {
    page.value = next
    return fetchList()
  }

  return {
    loading,
    list,
    total,
    page,
    pageSize,
    fetchList,
    handleSearch,
    onPageChange,
  }
}
