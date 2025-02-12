import { App } from 'vue'
import toolsJavascript from './tools-javascript'
import toolsVue3 from './tools-vue3'
import uivite from './uivite'
import pvue from './pvue'
import chat from './chat'
export const install = (app: App) => {
  app.use(toolsJavascript)
  app.use(toolsVue3)
  app.use(uivite)
  app.use(pvue)
  app.use(chat)
}

export default {
  install
}
