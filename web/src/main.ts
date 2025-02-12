import 'tools-css/index.css'
import 'tools-javascript'
import 'tools-websocket'
import { createApp } from 'vue'
import App from './App.vue'
import './styles/index.less'
import Config from './utils/Config'
const app = createApp(App)
//#ifvar-dev
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
app.use(ElementPlus)
//#endvar
Config.init(app)

