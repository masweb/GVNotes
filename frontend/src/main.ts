import { createApp } from 'vue'
import { createPinia } from 'pinia'
import Coreui from '@coreui/vue'
import { i18n } from '@/i18n/i18n'

import App from './App.vue'
import '@/css/main.scss'
import 'splitpanes/dist/splitpanes.css'
import './utils/coreui-dark-vue.js'
import './composables/useValidation'

const pinia = createPinia()
const app = createApp(App)

app.use(pinia)
app.use(i18n)
app.use(Coreui)

app.mount('#app')
