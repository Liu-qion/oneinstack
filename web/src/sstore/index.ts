import { sapp } from './sapp'
import { sconfig } from './sconfig'
import { sutil } from './sutil'

sutil.store = {
  sconfig,
  sapp,
  sutil
}

export const reset = () => {
  Object.keys(sutil.store).forEach((key) => {
    if (key !== 'sutil') sutil.store[key].reset?.()
  })
}

export const load = () => {
  Object.keys(sutil.store).forEach((key) => {
    if (key !== 'sutil') sutil.store[key].load?.()
  })
}

//初始化应用
export const initApp = async () => {
  load()

  //设置icon
  let link = document.createElement('link')
  link.rel = 'icon'
  link.href = '/favicon.ico'
  document.getElementsByTagName('head')[0].appendChild(link)
}

sutil.reset = reset

export default () => {}
