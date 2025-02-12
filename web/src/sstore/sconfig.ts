import System from '@/utils/System'
import { reactive } from 'vue'
import { sutil } from './sutil'

export const sconfig = reactive({
  load: () => {
    sconfig.userInfo = Cookie.get('userInfo') || (null as any)
  },
  userInfo: null as any as any,

  /**
   * 成功登录用户
   */
  login(info: any) {
    // 清空token
    Cookie.clear()
    sutil.reset()

    sconfig.userInfo = info
    Cookie.set('userInfo', sconfig.userInfo)
  },
  /**
   * 退出登录
   */
  logout(toLogin = false) {
    // 清空token
    Cookie.clear()
    sutil.reset()
    sconfig.userInfo = null as any

    // 是否跳转到登录页
    if (toLogin)
      setTimeout(() => {
        System.router.push('/login')
      }, 300)
  }
})

export default sconfig
