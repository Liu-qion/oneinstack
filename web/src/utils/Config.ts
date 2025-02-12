import i18n from '@/lang'
import plugins from '@/plugins'
import { initRouter } from '@/router'
import { initApp } from '@/sstore'
import { App } from 'vue'
import HttpConfig from './HttpConfig'
import System from './System'
import Components from '@/components'
export default class Config {
  /**
   * 初始化
   */
  static async init(app: App) {
    // 初始化系统参数
    System.init()

    // 初始化插件
    app.use(plugins)

    // 初始化最后打包版本
    window.v = window.version = '#{version}'

    // 初始化i18n
    app.use(i18n as any)
    await i18n.setLang()

    // 初始化http
    HttpConfig.init(System.env)

    // 初始化应用数据
    initApp()

    // 初始化组件
    app.use(Components)

    // 初始化路由
    const router = initRouter()
    app.use(router)
    router.isReady().then(async () => {
      // 初始化完成
      app.mount('#app')
    })
  }
}
