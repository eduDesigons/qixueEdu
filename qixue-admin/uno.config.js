/**********************************
 * @Author: Ronnie Zhang
 * @LastEditor: Ronnie Zhang
 * @LastEditTime: 2023/12/05 21:30:57
 * @Email: zclzone@outlook.com
 * Copyright © 2023 Ronnie Zhang(大脸怪) | https://isme.top
 **********************************/

import { FileSystemIconLoader } from '@iconify/utils/lib/loader/node-loaders'
import presetRemToPx from '@unocss/preset-rem-to-px'
import { defineConfig, presetAttributify, presetIcons, presetWind3 } from 'unocss'
import { getIcons } from './build/index.js'

const icons = getIcons()
export default defineConfig({
  presets: [
    presetWind3(),
    presetAttributify(),
    presetIcons({
      warn: true,
      prefix: ['i-'],
      extraProperties: {
        display: 'inline-block',
        width: '1em',
        height: '1em',
      },
      collections: {
        me: FileSystemIconLoader('./src/assets/icons/isme'),
        fe: FileSystemIconLoader('./src/assets/icons/feather'),
      },
    }),
    presetRemToPx({ baseFontSize: 4 }),
  ],
  safelist: icons.map(icon => `${icon} ${icon}?mask`.split(' ')).flat(),
  shortcuts: [
    ['wh-full', 'w-full h-full'],
    ['f-c-c', 'flex justify-center items-center'],
    ['flex-col', 'flex flex-col'],
    ['card-border', 'border border-solid border-[var(--glass-stroke)]'],
    ['auto-bg', 'bg-transparent'],
    ['auto-bg-hover', 'hover:bg-[rgba(var(--primary-color),0.06)] hover:dark:bg-[rgb(255_255_255_/_8%)]'],
    ['auto-bg-highlight', 'bg-[rgba(var(--primary-color),0.06)] dark:bg-[rgb(255_255_255_/_8%)]'],
    ['text-highlight', 'rounded-8 px-8 py-2 auto-bg-highlight'],
    ['text-secondary', 'text-[var(--text-secondary)]'],
  ],
  rules: [
    [
      'card-shadow',
      { 'box-shadow': '0 1px 2px -2px #00000029, 0 3px 6px #0000001f, 0 5px 12px 4px #00000017' },
    ],
  ],
  theme: {
    colors: {
      primary: 'rgba(var(--primary-color))',
      dark: '#18181c',
      light_border: 'var(--glass-stroke)',
      dark_border: 'rgb(255 255 255 / 12%)',
    },
  },
})
