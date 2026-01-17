import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import i18n from './i18n'

// 导入全局样式
import './assets/main.css'

// 导入Font Awesome
import '@fortawesome/fontawesome-free/css/all.min.css'

const app = createApp(App)

app.use(createPinia())
app.use(router)
app.use(i18n)

// 在应用挂载后同步i18n store和vue-i18n的状态
app.mount('#app')
