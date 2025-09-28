// main.js
import { createApp } from 'vue'
import App from './App.vue'
import { createPinia } from 'pinia'

// 不要导入 naive-ui CSS，Naive UI 会自动注入样式
// import 'naive-ui/dist/naive-ui.css'
// import 'naive-ui/lib/index.css'

const app = createApp(App)
app.use(createPinia())
app.mount('#app')
