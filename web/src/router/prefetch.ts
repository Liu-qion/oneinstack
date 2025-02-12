/**
 * 预加载路由数据
 */
const prefetchRouteData = {
} as any

Object.keys(prefetchRouteData).forEach((key) => {
  const _arr = [...prefetchRouteData[key]]
  _arr.forEach((route: string) => {
    prefetchRouteData[route] = [key, ..._arr]
    prefetchRouteData[route] = prefetchRouteData[route].filter((r: string) => r !== route)
  })
})

export default prefetchRouteData
