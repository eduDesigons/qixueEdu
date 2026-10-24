/** 固定主色 CSS 变量（不再支持运行时换色） */
import { defaultPrimaryRgb } from '@/settings'

export function applyPrimaryTheme() {
  const rgb = defaultPrimaryRgb
  const root = document.documentElement
  root.style.setProperty('--primary-color', rgb)
  document.body.style.setProperty('--primary-color', rgb)
  root.style.setProperty('--ctl-fill', `rgba(${rgb}, 0.06)`)
  root.style.setProperty('--ctl-fill-hover', `rgba(${rgb}, 0.11)`)
  root.style.setProperty('--ctl-fill-active', `rgba(${rgb}, 0.15)`)
  root.style.setProperty('--ctl-ring', `rgba(${rgb}, 0.22)`)
  root.style.setProperty('--ctl-ring-strong', `rgba(${rgb}, 0.32)`)
  root.style.setProperty('--ctl-hairline', `rgba(${rgb}, 0.1)`)
  root.style.setProperty('--glass-stroke', `rgba(${rgb}, 0.09)`)
}
