/**
 * 开发环境log处理
 */
export const initLog = () => {
  const _log = console.log
  console.log = (...args: any[]) => {
    _log(...args)
  }

  const _warn = console.warn
  console.warn = (...args: any[]) => {
    for (let i = 0; i < args.length; i++) {
      const str = args[i]
      if (typeof str === 'string') {
        const arr = ['Use of eval', 'minification', '::v-deep']
        for (let i = 0; i < arr.length; i++) {
          if (str.indexOf(arr[i]) !== -1) return
        }
      }
    }
    _warn(...args)
  }
}
