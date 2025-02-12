import { App } from 'vue'
export const install = (app: App) => {
  window.uni = {
    getSystemInfoSync: () => {
      return {
        uniPlatform: 'web'
      }
    }
  }
}

export default {
  install
}
