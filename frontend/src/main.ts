import { createApp } from 'vue'
import { createPinia } from 'pinia'
import Coreui from '@coreui/vue'
import { configure } from 'vee-validate'
import App from './App.vue'
import '@/css/main.scss'
import 'splitpanes/dist/splitpanes.css'
import './utils/coreui-dark-vue.js'

// Validate only on blur/change, not on every keystroke
configure({ validateOnInput: false, validateOnBlur: true })

const app = createApp(App)

app.use(createPinia())
app.use(Coreui)

app.mount('#app')
