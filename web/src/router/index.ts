import { createRouter, createWebHashHistory } from 'vue-router'

import System from '@/utils/System'
import prefetchRouteData from './prefetch'
import { CRouter } from 'tools-vue3'
import routes from './routes'
import sconfig from '@/sstore/sconfig'

export const initRouter = () => {
  const _routes = CRouter.init({
    modules: import.meta.glob('/src/views/**/*.vue'),
    theme: Cookie.get('pageTheme'),
    excludeReg: ['/com/', '/components/'],
    children: routes,
    init: [
      {
        path: '/',
        redirect: '/home'
      },
      {
        path: '/:matchOthers(.*)*',
        redirect: '/'
      }
    ],
    alias: {
      '/user/login/index': '/login',
      '/user/register/register': '/register'
    },
    pathHook(path) {
      return path.replace(/^\/pages/, '').replace(/\/index$/, '')
    },
  })

  console.log('routes', _routes)
  /**
   * 路由映射
   */
  const _routesMap = {} as any
  ObjectUtil.getObjectByChildren(_routes, (obj) => {
    if (obj.component) _routesMap[obj.path] = obj
  })

  const router = createRouter({
    history: createWebHashHistory(),
    routes: _routes
  })

  const whiteList = ['/login', '/login/scan']
  router.beforeEach(async (to, from, next) => {
    let _name = (to.meta.name as any) || ''
    let _title = System.env.title
    if (_name) _title = _title + ' - ' + _name
    document.title = _title
    document.documentElement.scrollTop = 0
    if (!sconfig.userInfo?.token && !whiteList.includes(to.path)) return next('/login')
    else next()
  })

  router.afterEach((guard) => {
    /**
     * 预加载路由
     */
    const _arr = prefetchRouteData[guard.path]
    if (_arr) {
      _arr.forEach((v: string) => {
        if (_routesMap[v]) {
          _routesMap[v].component()
        }
      })
    }
  })

  System.router = router

  return router
}
