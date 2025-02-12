const manager = {
  ERouter: {
    /**
     * 路由改变
     */
    change: '',

    /**
     * 重新加载当前页面
     */
    reload: '',

    /**
     * 浏览器显示和隐藏通知
     */
    browserShow: ''
  },
  EUser: {
    /**
     * 打开登录弹框
     */
    showLogin: ''
  },
  EPage: {
    /**
     * page组件页面正在滚动
     */
    scroll: '',

    /** page组件页面被滚动到底部50px */
    scrollBottom: '',

    /** 页面高度被改变 */
    changeHeight: ''
  },
  EKey: {
    /**
     * 页面主题
     */
    pageTheme: '',
    /**
     * 开发接口地址
     */
    apiurl: '',
    /**
     * 语言
     */
    language: ''
  }
}

Object.keys(manager).forEach((ikey) => {
  //@ts-ignore
  const item = manager[ikey]
  Object.keys(item).forEach((jkey) => {
    if (ikey == 'EKey') {
      //@ts-ignore
      const _val = manager[ikey][jkey]
      if (!_val) {
        //@ts-ignore
        manager[ikey][jkey] = `${jkey}`
      }
      if (typeof _val === 'object') {
        Object.keys(_val).forEach((kkey) => {
          //@ts-ignore
          if(!manager[ikey][jkey][kkey]) manager[ikey][jkey][kkey] = `${jkey}_${kkey}`
        })
      }
    } else {
      //@ts-ignore
      manager[ikey][jkey] = `${ikey}_${jkey}`
    }
  })
})

export const ERouter = manager.ERouter
export const EUser = manager.EUser
export const EPage = manager.EPage

/**
 * 缓存key
 */
export const EKey = manager.EKey
