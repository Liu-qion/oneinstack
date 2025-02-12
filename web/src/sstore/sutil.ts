import i18n from '@/lang'
import { reactive } from 'vue'
import System from '../utils/System'

export const sutil = reactive({
  /**
   * 请求根地址
   */
  baseUrl: '',
  reset: () => {},
  store: {} as any,
  getStore(name: string) {
    return sutil.store[name]
  },
  isNull(val: any) {
    return StrUtil.isNull(val)
  },

  /**
   * rem转px
   * @param rem
   * @returns
   */
  rem2px: (rem: number) => {
    return H5Util.rem2px(rem)
  },

  /**
   * px转rem
   * @param px
   * @returns
   */
  px2rem: (px: number) => {
    return H5Util.px2rem(px)
  },

  /**
   * 返回上一页，如果没进行过跳转，进入首页
   */
  pageBack: () => {
    if (window.history.state.back) {
      System.router.back()
    } else {
      System.router.replace('/')
    }
  },

  // 数据处理
  dataHandling: (amount: number | string, fixed = 4, endFixed?: number) => {
    if (StrUtil.isNull(amount)) amount = '0'
    let unit = '',
      isNegative = false
    amount = parseFloat(amount.toString()) || 0
    if (amount < 0) {
      isNegative = true
      amount = Math.abs(amount)
    }
    if (amount < 10000) {
    } else if (amount < 1000000) {
      // 如果金额在10,000到1,000,000之间，则显示为xxxx.xxW（万）
      amount /= 1000
      unit = 'K'
    } else if (amount < 1000000000) {
      // 如果金额在1,000,000到1,000,000,000之间，则显示为xxxx.xxM（百万）
      amount /= 1000000
      unit = 'M'
    } else {
      // 如果金额超过1,000,000,000，则显示为xxxx.xxB（十亿）
      amount /= 1000000000
      unit = 'B'
    }

    //对数据进行保留n位处理并去掉末位小数
    let [integerPart, decimalPart = ''] = String(amount).split('.')
    decimalPart = decimalPart.substring(0, fixed)
    decimalPart = decimalPart.replace(/0+$/, '')
    //对末位进行0补齐
    if (endFixed) {
      if (decimalPart.length < endFixed) {
        while (decimalPart.length < endFixed) {
          decimalPart += '0'
        }
      }
    }
    amount = decimalPart ? `${integerPart}.${decimalPart}` : integerPart
    if (isNegative) amount = '-' + amount
    return amount + unit
  },

  /** 获取时间区间选择 */
  getDateOptions: (all = true) => {
    const interval = []
    if (all)
      interval.push({
        label: i18n.t('dateOptions.AllDates'),
        value: null,
        key: 'all'
      })
    interval.push({
      label: i18n.t('dateOptions.Today'),
      value: TimeUtil.somedayse(),
      key: 'today'
    })
    interval.push({
      label: i18n.t('dateOptions.Yesterday'),
      value: TimeUtil.somedayse(-1),
      key: 'yesterday'
    })

    const t = new Date()

    const _day = TimeUtil.somedayse
    // 本周起始时间-结束时间
    const today = t.getDay()
    const startTime = _day(1 - today)[0]
    const endTime = _day(7 - today)[1]
    const thisWeek = {
      label: i18n.t('dateOptions.ThisWeek'),
      value: [startTime, endTime],
      key: 'thisWeek'
    }
    interval.push(thisWeek)

    // 上周起始时间-结束时间
    const lastWeekSTime = _day(1 - today - 7)[0]
    const lastWeekETime = _day(7 - today - 7)[1]
    const lastWeek = {
      label: i18n.t('dateOptions.LastWeek'),
      value: [lastWeekSTime, lastWeekETime],
      key: 'lastWeek'
    }
    interval.push(lastWeek)

    // 本月起始时间-结束时间
    const ms = new Date()
    ms.setDate(1)
    const mstartTime = _day(0, ms)
    ms.setMonth(ms.getMonth() + 1)
    ms.setDate(1)
    ms.setDate(ms.getDate() - 1)
    const mendTime = _day(0, ms)
    const thisMonth = {
      label: i18n.t('dateOptions.ThisMonth'),
      value: [mstartTime[0], mendTime[1]],
      key: 'thisMonth'
    }
    interval.push(thisMonth)

    return interval
  },

  /**
   * 获取时间提示
   * @param time 时间戳
   * @returns 时间提示
   */
  getTimeTip: (time: number) => {
    const _thisTime = new Date()
    const _formatTime = new Date(time)
    const now = _thisTime.getTime()
    const diff = now - time

    // //刚刚
    // if (diff < 60 * 1000) return i18n.t('chatRoom.just_now')

    //今天
    if (diff < 24 * 60 * 60 * 1000) {
      const _dd = _thisTime.Format('dd')
      const _ddFormat = _formatTime.Format('dd')
      if (_dd == _ddFormat) return _formatTime.Format('hh:mm')
    }

    //昨天
    if (diff < 48 * 60 * 60 * 1000) {
      const _yd = new Date(TimeUtil.somedayse(-1)[0])
      const _dd = _yd.Format('dd')
      const _ddFormat = _formatTime.Format('dd')
      if (_dd == _ddFormat) return i18n.t('chatRoom.yesterday')
    }

    //今年的显示月日
    const _year = _thisTime.getFullYear()
    const _yearFormat = _formatTime.getFullYear()
    if (_year == _yearFormat) return _formatTime.Format('MM/dd')

    return _formatTime.Format('MM/dd/yyyy')
  },
  getFileType: (name: any) => {
    const idx = name.lastIndexOf('.')
    return name.slice(idx + 1).split('?')[0]
  },

  /**
   * 字节转换
   * @param bytes
   * @returns
   */
  bytesTransform: (bytes: number) => {
    if (bytes < 1024) {
      return { value: bytes, unit: 'B', strValue: `${bytes}B` }
    } else if (bytes < 1024 * 1024) {
      const value = parseFloat((bytes / 1024).toFixed(2))
      return { value, unit: 'KB', strValue: `${value}KB` }
    } else if (bytes < 1024 * 1024 * 1024) {
      const value = parseFloat((bytes / 1024 / 1024).toFixed(2))
      return { value, unit: 'MB', strValue: `${value}MB` }
    } else if (bytes < 1024 * 1024 * 1024 * 1024) {
      const value = parseFloat((bytes / 1024 / 1024 / 1024).toFixed(2))
      return { value, unit: 'GB', strValue: `${value}GB` }
    } else {
      const value = parseFloat((bytes / 1024 / 1024 / 1024 / 1024).toFixed(2))
      return { value, unit: 'TB', strValue: `${value}TB` }
    }
  },

  /**
   * 匹配文件路径
   * @param path 文件路径
   * @returns 匹配结果
   */
  matchFilePath: (path: string) => {
    return /^\/(?:[^/]+\/)*[^/]+$/.test(path)
  },

  /**
   * 获取css变量
   * @param str 变量名称
   * @return 变量值
   */
  getCssVariable(str: string) {
    const root = getComputedStyle(document.documentElement, null)
    return root.getPropertyValue(str)
  }
})

export default sutil
