import { App } from 'vue'
import { InertiaScroller } from './pvue-com/InertiaScroller'

export const install = (app: App) => {
  /** 支持pc端滚动指令，对需要pc端滚动的元素添加此指令 v-scroll, 传入空对象，支持toTop,toLeft,showDom方法 */
  app.directive('scroll', {
    mounted(el: HTMLDivElement, bindinfo) {
      const value = bindinfo.value
      if (value && typeof value === 'object') {
        //滚动到顶部指定位置
        value.toTop = (top: number, useAni: boolean = true) => {
          if (useAni) {
            el.scrollTo({
              top,
              behavior: 'smooth'
            })
          } else {
            el.scrollTop = top
          }
        }
        //滚动到左侧指定位置
        value.toLeft = (left: number, useAni: boolean = true) => {
          if (useAni) {
            el.scrollTo({
              left,
              behavior: 'smooth'
            })
          } else {
            el.scrollLeft = left
          }
        }

        //滚动到最右侧
        value.toRight = (useAni: boolean = true) => {
          if (useAni) {
            el.scrollTo({
              left: el.scrollWidth,
              behavior: 'smooth'
            })
          } else {
            el.scrollLeft = el.scrollWidth
          }
        }

        //滚动到最下侧
        value.toBottom = (useAni: boolean = true) => {
          if (useAni) {
            el.scrollTo({
              top: el.scrollHeight,
              behavior: 'smooth'
            })
          } else {
            el.scrollTop = el.scrollHeight
          }
        }

        //将dom完全显示出来
        value.showDom = (dom: HTMLElement | string, useAni: boolean = true) => {
          const domEl = typeof dom === 'string' ? document.querySelector(dom) : dom
          if (!domEl) return
          const rect = domEl.getBoundingClientRect()
          value.toTop(el.scrollTop + rect.top - el.clientHeight / 2, useAni)
          value.toLeft(el.scrollLeft + rect.left - el.clientWidth / 2, useAni)
        }
      }
      new InertiaScroller(el)
    }
  })
}

export default {
  install
}
