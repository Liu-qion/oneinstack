import { App } from 'vue'
import routeEvent from './routeEvent/index.vue'
import routerView from './routerView/index.vue'
import { uscom } from './uscom'
export const install = (app: App) => {
  uscom.init()
  app.component('xRouteEvent', routeEvent)
  app.component('xRouterView', routerView)
}

export default {
  install
}
