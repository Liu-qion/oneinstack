import { Scope } from 'tools-vue3'
import { App } from 'vue'

export const install = (app: App) => {
  // 初始化事件
  Scope.CEventBean = CEventBean
  // 初始化定时器
  Scope.TimerBean = TimerBean
}

export default {
  install
}
