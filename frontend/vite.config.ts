import { fileURLToPath, URL } from 'node:url'
import AutoImport from 'unplugin-auto-import/vite'
import Components from 'unplugin-vue-components/vite'
import vue from '@vitejs/plugin-vue'
import { type Plugin, defineConfig } from 'vite'
import vueDevTools from 'vite-plugin-vue-devtools'
import VueI18nPlugin from '@intlify/unplugin-vue-i18n'
import { resolve } from 'node:path'

// https://vite.dev/config/
export default defineConfig({
 plugins: [
  vue(),
  vueDevTools(),
  VueI18nPlugin.vite({
   include: [resolve(__dirname, './locales/**')],
   runtimeOnly: false
  }),
  AutoImport({
   dts: 'src/auto-imports.d.ts',
   imports: ['vue', 'pinia', 'vue-router', 'vue-i18n', 'vee-validate'],
   include: [/\.[tj]sx?$/, /\.vue$/, /\.vue\?vue/],
   dirs: ['src/composables', 'src/plugins', 'src/services', 'src/utils', 'src/types', 'src/stores']
  }) as Plugin,
  Components({
   dirs: ['src/services/**', 'src/components/**', 'src/views/**']
  }) as Plugin
 ],
 css: {
  preprocessorOptions: {
   scss: {
    silenceDeprecations: ['mixed-decls', 'color-functions', 'global-builtin', 'import'],
    quietDeps: true
   }
  }
 },
 optimizeDeps: {
  include: ['@coreui/coreui', '@coreui/vue', '@tabler/icons-vue']
 },
 resolve: {
  alias: {
   '@': fileURLToPath(new URL('./src', import.meta.url))
  }
 },
 server: {
  // Proxy /images/* to the Go image server (dev mode only).
  // Port must match imageServerPort in app.go.
  proxy: {
   '/images': 'http://127.0.0.1:34201'
  }
 }
})
