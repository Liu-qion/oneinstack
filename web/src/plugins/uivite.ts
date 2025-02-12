import uiVite from 'ui-vite'
import { App } from 'vue'
let iconObj = {} as any
let siconObj = {} as any
export const install = (app: App) => {
  app.use(uiVite)

  VConf.vicon.url = '/static/vicon'

  iconObj = Cookie.get('iconObj') || {}
  siconObj = Cookie.get('siconObj') || {}
  VConf.vicon.fun.getRes = async (url) => {
    return iconObj[url]
  }

  VConf.vicon.fun.setRes = (url, content) => {
    iconObj[url] = content
    Cookie.set('iconObj', iconObj)
  }

  VConf.vsicon.url = '/static/vicon/svg'
  VConf.vsicon.lib = 'common'

  VConf.vsicon.fun.getRes = async (url) => {
    return siconObj[url]
  }

  VConf.vsicon.fun.setRes = (url, content) => {
    siconObj[url] = content
    Cookie.set('siconObj', siconObj)
  }
}

export default {
  install
}
