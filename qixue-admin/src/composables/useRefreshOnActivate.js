export function useRefreshOnActivate(refresh) {
  let skip = true
  onActivated(() => {
    if (skip) {
      skip = false
      return
    }
    refresh?.()
  })
}
