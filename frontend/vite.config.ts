import { fileURLToPath, URL } from 'node:url'
import AutoImport from 'unplugin-auto-import/vite'
import Components from 'unplugin-vue-components/vite'
import vue from '@vitejs/plugin-vue'
import { type Plugin, defineConfig, loadEnv } from 'vite'

// https://vite.dev/config/
export default defineConfig({
 plugins: [
  vue(),
  AutoImport({
   dts: 'src/auto-imports.d.ts',
   imports: ['vue', 'pinia', 'vue-router', 'vue-i18n', 'vee-validate'],
   include: [/\.[tj]sx?$/, /\.vue$/, /\.vue\?vue/],
   dirs: ['src/composables/**', 'src/plugins/**', 'src/services/**', 'src/utils/**', 'src/types/**', 'src/stores/**']
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
 }
})
