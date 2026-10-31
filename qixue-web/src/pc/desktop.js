import { ref } from 'vue'

const mq = typeof window !== 'undefined' ? window.matchMedia('(min-width: 1024px)') : null

export const isDesktop = ref(!!(mq && mq.matches))

if (mq) {
  document.documentElement.classList.toggle('is-desktop', mq.matches)
}

export function bindDesktop() {
  if (!mq) return
  const update = () => {
    isDesktop.value = mq.matches
    document.documentElement.classList.toggle('is-desktop', mq.matches)
  }
  update()
  mq.addEventListener('change', update)
}
