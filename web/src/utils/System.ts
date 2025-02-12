import sutil from '@/sstore/sutil'
import { Router } from 'vue-router'
import { globalType } from '../../build/env/globalVar'
import { ElLoading, ElMessage } from 'element-plus'
export default class System {

  /**
   * 环境变量
   */
  static env: globalType = {} as any

  /**
   * 初始化系统参数
   */
  static init() {

    // 初始化环境变量
    System.env = JSON.parse('#{global}')
  }

  /**
   * 路由
   */
  static router: Router

  /**
   * 获取路由路径
   */
  static getRouterPath = (): any => {
    return this.router.currentRoute.value.path
  }

  /**
   * 获取路由参数,自动获取query、params中的参数， 哪个有就返回哪个
   */
  static getRouterParams = (): any => {
    const _query = this.router.currentRoute.value.query || {}
    const _params = this.router.currentRoute.value.params || {}
    if (Object.keys(_query).length) {
      return _query
    } else if (Object.keys(_params).length) {
      return _params
    }
    return {}
  }

  /**
   * 失败提示
   * @param content
   * @param options
   */
  static er = (content: string, options: any = {}) => {
    const optionObj = Object.assign(
      {
        message: content,
        type: 'error'
      },
      options
    )
    ElMessage(optionObj)
  }

  /**
   * 成功提示
   * @param content
   * @param options
   */
  static sr = (content: string, options: any = {}) => {
    const optionObj = Object.assign(
      {
        message: content,
        type: 'success'
      },
      options
    )
    ElMessage(optionObj)
  }

  private static loadingObj: {
    setText: (text: string) => void
    close: () => void
    visible: import('vue').Ref<boolean>
    fullscreen: import('vue').Ref<boolean>
    lock: import('vue').Ref<boolean>
    closed?: import('vue').Ref<(() => void) | undefined> | undefined
  } = null as any

  private static loadingTimer: any = null
  /**
   * 加载中
   * @param show
   * @param text
   */
  static loading = (show: boolean = true, text: string = '加载中') => {
    Timer.un(this.loadingTimer)
    this.loadingTimer = Timer.once(() => {
      if (show) {
        this.loadingObj = ElLoading.service({
          lock: true,
          text,
          background: 'rgba(0, 0, 0, 0.3)'
        })
      } else {
        if (this.loadingObj) {
          this.loadingObj.close()
          this.loadingObj = null as any
        }
      }
    }, 50)
  }

  /**
   * 打开新地址
   * @param url
   * @param isCenter
   */
  static openUrl = (url: string, isCenter: boolean = true) => {
    if (isCenter) {
      let left = screen.width / 2 - 1000 / 2
      let top = screen.height / 2 - 800 / 2 - 30
      window.open(
        url,
        '_blank',
        'toolbar=no, location=yes, directories=no, status=yes, menubar=no, scrollbars=yes, resizable=no, copyhistory=no, width=' +
          1000 +
          ', height=' +
          800 +
          ', top=' +
          top +
          ', left=' +
          left
      )
    } else {
      window.open(
        url,
        'DescriptiveWindowName' + StrUtil.getId(),
        'resizable,scrollbars=yes,status=1,width=1024, height=600, top=0, left=0'
      )
    }
  }

  static moduleObj = {
    'lottie': async () => {
      return await import('@lottiefiles/lottie-player')
    },
    'html2canvas': async () => {
      return await import('html2canvas')
    }
  }

  /**
   * 加载模块
   * @param module
   */
  static loadModule = async (module: keyof typeof System.moduleObj): Promise<any> => {
    const _module = await System.moduleObj[module]()
    return _module.default
  }

  /**
   * 截图返回图片url数据
   * @param param
   * @returns
   */
  static getImgPic = (param: { id: string; windowWidth?: number; windowHeight?: number }): Promise<string> => {
    return new Promise(async (resolve, reject) => {
      let htmlDom: any = document.getElementById(param.id)
      const html2canvas = await System.loadModule('html2canvas')
      try {
        html2canvas(htmlDom, {
          logging: false,
          allowTaint: true,
          scale: window.devicePixelRatio,
          width: param.windowWidth,
          height: param.windowHeight,
          windowWidth: param.windowWidth,
          windowHeight: param.windowHeight,
          useCORS: true,
          backgroundColor: '#ffffff'
        }).then(function (canvas: any) {
          let imgUrl = canvas.toDataURL('image/png')
          resolve(imgUrl)
        })
      } catch (error) {
        reject(error)
      }
    })
  }

  /**
   *  下载文件
   */
  static async download(
    param:
      | string
      | {
          /**
           * 下载地址或者文件内容
           */
          content: string
          /**
           * 文件名，包含后缀，默认自动生成时间戳加文件类型
           */
          name?: string
          /**
           * 是否显示提示
           */
          toast?: boolean
        }
  ): Promise<any> {
    if (typeof param === 'string') {
      param = { content: param }
    }
    let { content, name, toast = true } = param
    if (!name) {
      const getDataType = (ctx: string) => {
        const res = ctx.match(/^data:(.*);base64,(.*)$/)
        if (res) return res[1].split('/').pop()
        return 'png'
      }
      name = Date.now() + '.' + (content.startsWith('http') ? sutil.getFileType(content) : getDataType(content))
    }
    if (content.startsWith('http')) {
      const res = await fetch(content)
      const blob = await res.blob()
      content = URL.createObjectURL(blob)
    }
    const tempLink = document.createElement('a')
    tempLink.style.display = 'none'
    tempLink.href = content
    tempLink.setAttribute('download', name)
    tempLink.setAttribute('target', '_blank')
    document.body.appendChild(tempLink)
    tempLink.click()
    document.body.removeChild(tempLink)
    window.URL.revokeObjectURL(content)

    if (toast) {
    }
  }

}
