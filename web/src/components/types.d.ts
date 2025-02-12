import '@vue/runtime-core'
    export {}
    declare module '@vue/runtime-core' {
        export interface GlobalComponents {

        xEcharts:typeof import('./echarts/index.vue')['default']
        xRouteEvent:typeof import('./routeEvent/index.vue')['default']
        xRouterView:typeof import('./routerView/index.vue')['default']

        }
    }
    