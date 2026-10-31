const listeners = new Map()

export function on(event, handler) {
  const list = listeners.get(event) || []
  list.push(handler)
  listeners.set(event, list)
}

export function off(event, handler) {
  if (!handler) {
    listeners.delete(event)
    return
  }
  const list = listeners.get(event) || []
  listeners.set(event, list.filter((item) => item !== handler))
}

export function emit(event, payload) {
  ;(listeners.get(event) || []).forEach((handler) => handler(payload))
}
