import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import { VitePWA } from 'vite-plugin-pwa'
import { NodeTypes } from '@vue/compiler-core'
import { fileURLToPath, URL } from 'node:url'

function markClickable(node) {
  if (node.type !== NodeTypes.ELEMENT || node.tag === 'template' || node.tag === 'slot') return
  const clickable = node.props.some((prop) => {
    if (prop.type !== NodeTypes.DIRECTIVE || prop.name !== 'on' || !prop.arg) return false
    return prop.arg.type === NodeTypes.SIMPLE_EXPRESSION && prop.arg.content === 'click'
  })
  if (!clickable) return
  const klass = node.props.find((prop) => prop.type === NodeTypes.ATTRIBUTE && prop.name === 'class' && prop.value)
  if (klass) {
    const names = klass.value.content.split(/\s+/)
    if (!names.includes('qx-clickable')) klass.value.content = `${klass.value.content} qx-clickable`
    return
  }
  node.props.push({
    type: NodeTypes.ATTRIBUTE,
    name: 'class',
    value: { type: NodeTypes.TEXT, content: 'qx-clickable', loc: node.loc },
    loc: node.loc,
  })
}

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, fileURLToPath(new URL('.', import.meta.url)))
  return {
  base: env.VITE_PUBLIC_PATH || '/',
  plugins: [
    vue({
      template: {
        compilerOptions: {
          isCustomElement: (tag) => tag === 'uni-page-head' || tag === 'uni-tabbar',
          nodeTransforms: [markClickable],
        },
      },
    }),
    VitePWA({
      strategies: 'injectManifest',
      srcDir: 'src',
      filename: 'sw.js',
      registerType: 'prompt',
      includeAssets: ['favicon.png', 'pwa/*.png'],
      injectManifest: {
        globPatterns: ['**/*.{js,css,html,ico,png,svg,webp,woff2,woff,ttf}'],
        maximumFileSizeToCacheInBytes: 3 * 1024 * 1024,
      },
      manifest: {
        id: '/',
        name: '启学',
        short_name: '启学',
        description: '启学学员端：课程学习、作业与证书',
        // 闪屏底与页面一致；状态栏透明由 index.html theme-color=transparent 控制
        theme_color: '#f8f8f8',
        background_color: '#f8f8f8',
        display: 'standalone',
        // 优先独立窗口（无浏览器地址栏/工具栏），桌面 Chromium 可叠层窗口控件
        display_override: ['window-controls-overlay', 'standalone', 'minimal-ui', 'fullscreen'],
        orientation: 'any',
        lang: 'zh-CN',
        dir: 'ltr',
        start_url: '/?source=pwa',
        scope: '/',
        prefer_related_applications: false,
        categories: ['education', 'productivity'],
        launch_handler: {
          client_mode: ['navigate-existing', 'auto'],
        },
        handle_links: 'preferred',
        icons: [
          { src: 'pwa/icon-192.png', sizes: '192x192', type: 'image/png', purpose: 'any' },
          { src: 'pwa/icon-512.png', sizes: '512x512', type: 'image/png', purpose: 'any' },
          { src: 'pwa/icon-maskable-512.png', sizes: '512x512', type: 'image/png', purpose: 'maskable' },
        ],
      },
      devOptions: {
        enabled: false,
      },
    }),
  ],
  resolve: {
    alias: [
      { find: '@/uni_modules/mescroll-uni/components/mescroll-uni/mescroll-mixins.js', replacement: fileURLToPath(new URL('./src/components/mescroll-mixins.js', import.meta.url)) },
      { find: '@/uni_modules/mescroll-uni/components/mescroll-uni/mixins/mescroll-more-item.js', replacement: fileURLToPath(new URL('./src/components/mescroll-more-item.js', import.meta.url)) },
      { find: '@', replacement: fileURLToPath(new URL('./src', import.meta.url)) },
    ],
  },
  css: {
    preprocessorOptions: {
      scss: {
        additionalData: `@import "@/styles/vars.scss";\n`,
      },
    },
  },
  server: {
    port: 5180,
    strictPort: true,
  },
  build: {
    assetsInlineLimit: 8 * 1024,
    cssCodeSplit: false,
    chunkSizeWarningLimit: 2500,
    rollupOptions: {
      output: {
        inlineDynamicImports: true,
      },
    },
  },
  }
})
