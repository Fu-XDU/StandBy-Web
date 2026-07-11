import './assets/main.css'

import { createApp } from 'vue'
import App from './App.vue'
import { initDeviceIdFromUrl } from '@/composables/useRemoteDevice'

import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import 'element-plus/theme-chalk/dark/css-vars.css'

initDeviceIdFromUrl()

const app = createApp(App)
app.use(ElementPlus)
app.mount('#app')
