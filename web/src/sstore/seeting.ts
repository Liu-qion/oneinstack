import { EKey } from '@/enum/Enum'
import { reactive } from 'vue'

export const sapp = reactive({
  load: () => {
    // sapp.theme = Cookie.get(EKey.pageTheme) || 'light'
  },
  theme: 'light' as 'light' | 'dark',
  setTheme(theme: 'light' | 'dark') {
    sapp.theme = theme
    Cookie.set(EKey.pageTheme, theme)
    document.documentElement.className = theme
  }
})

export default sapp
