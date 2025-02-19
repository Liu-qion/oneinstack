import 'tools-css/index.css'
import 'tools-javascript'
import 'tools-websocket'
import { createApp } from 'vue'
import { createPinia } from 'pinia';
import App from './App.vue'
import './styles/index.less'
import Config from './utils/Config'

const pinia = createPinia();
const app = createApp(App)
//#ifvar-dev
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import { createTerminal } from 'vue-web-terminal'
// import 'vue-web-terminal/lib/theme/dark.css'
app.use(ElementPlus)
// 全局注册 Terminal 组件
app.use(createTerminal());
app.use(pinia);
//#endvar
Config.init(app)

