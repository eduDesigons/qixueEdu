function isNull(val) {
  return val === null
}

function isUndef(val) {
  return typeof val === 'undefined'
}

export function isNullOrUndef(val) {
  return isNull(val) || isUndef(val)
}

/** Whether the path is an external link (http, mailto, or tel). */
export function isExternal(path) {
  return /^https?:|mailto:|tel:/.test(path)
}
