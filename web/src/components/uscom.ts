import { reactive } from 'vue'

export const uscom = reactive({
  init: () => {
    const theme = Cookie.get('pageTheme')
    const _uscom = uscom as any
    Object.keys(uscom).forEach((key) => {
      if (key !== 'init' && !key.endsWith('Theme')) {
        const themeVal = _uscom[`${key}Theme`]?.[theme]
        if (themeVal) {
          _uscom[key] = ObjectUtil.deepMerge(_uscom[key], themeVal)
        }
      }
    })
  },
  noData: {
    theme: {
      'no-method': {
        src: '/static/img/wallet/no-method.png'
      }
    } as any
  },
  noDataTheme: {
    'blue': {
      theme: {
        'no-method': {
          src: '/static/theme/blue/no-method.webp'
        }
      }
    }
  }
})
