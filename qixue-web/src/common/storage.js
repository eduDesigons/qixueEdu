const PREFIX = 'qx:'

function keyOf(key) {
  return PREFIX + key
}

export function getStorage(key) {
  const raw = localStorage.getItem(keyOf(key))
  if (raw == null) return ''
  try {
    return JSON.parse(raw)
  } catch (e) {
    return raw
  }
}

export function setStorage(key, data) {
  localStorage.setItem(keyOf(key), JSON.stringify(data))
}

export function removeStorage(key) {
  localStorage.removeItem(keyOf(key))
}
