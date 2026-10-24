import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
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
