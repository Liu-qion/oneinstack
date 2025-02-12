import { CRouterParamChildren } from 'tools-vue3/dist/utils/CRouter'

/**
 * 路由配置-仅对存在子路由的页面统一配置，以下为示例使用时按需求修改
 */
export default [
  {
    path: '/layout',
    component: '/layout/index',
    childrenReg: '^/pages',
    children: [
      {
        path: '/database',
        component: '/pages/database/index',
        childrenReg: '^/pages/database/',
        redirect: '/database/mysql'
      }
    ]
  }
] as Partial<CRouterParamChildren>[]
