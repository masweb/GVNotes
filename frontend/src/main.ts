import { createApp } from 'vue'
import { createPinia } from 'pinia'
import Coreui from '@coreui/vue'
import App from './App.vue'
import '@/css/main.scss'
import 'splitpanes/dist/splitpanes.css'
import './utils/coreui-dark-vue.js'
import './composables/useValidation'

const app = createApp(App)

app.use(createPinia())
app.use(Coreui)

app.mount('#app')
