<script setup lang="ts">
import sapp from '@/sstore/sapp'
import { toRefs } from 'vue'

const { theme } = toRefs(sapp)
sapp.setTheme(theme.value)

const colors = {
  light: {
    default: ['#c8c8c8'],
    select: ['#ffffff']
  },
  dark: {
    default: ['#b7b7b7'],
    select: ['#212e3d']
  }
}

const toggleTheme = (event: MouseEvent) => {
  // @ts-ignore
  if (!document.startViewTransition) return sapp.setTheme(theme.value === 'light' ? 'dark' : 'light')

  const x = event.clientX
  const y = event.clientY
  const endRadius = Math.hypot(Math.max(x, innerWidth - x), Math.max(y, innerHeight - y))

  let isDark: boolean

  // @ts-ignore
  const transition = document.startViewTransition(() => {
    isDark = theme.value === 'dark'
    sapp.setTheme(theme.value === 'light' ? 'dark' : 'light')
  })

  transition.ready.then(() => {
    const clipPath = [`circle(0px at ${x}px ${y}px)`, `circle(${endRadius}px at ${x}px ${y}px)`]
    document.documentElement.animate(
      {
        clipPath: isDark ? clipPath.reverse() : clipPath
      },
      {
        duration: 500,
        easing: 'ease-in',
        pseudoElement: isDark ? '::view-transition-old(root)' : '::view-transition-new(root)'
      }
    )
  })
}
</script>

<template>
  <div class="theme-switch" @click="toggleTheme">
    <div :class="{ active: theme === 'light' }" class="theme-switch__core">
      <v-s-icon name="day" :color="colors[theme][theme === 'light' ? 'select' : 'default']" size="16" />
    </div>
    <div :class="{ active: theme === 'dark' }" class="theme-switch__core">
      <v-s-icon name="night" :color="colors[theme][theme === 'dark' ? 'select' : 'default']" size="16" />
    </div>
  </div>
</template>

<style scoped lang="less">
.theme-switch {
  width: 60px;
  height: 30px;
  background-color: rgb(var(--bg-card-color));
  border-radius: 15px;
  display: flex;
  justify-content: center;
  align-items: center;
  cursor: pointer;

  &__core {
    width: 26px;
    height: 26px;
    border-radius: 50%;
    display: flex;
    justify-content: center;
    align-items: center;
    transition: all 0.3s;

    &.active {
      background-color: var(--el-color-primary);
    }

    img {
      width: 16px;
      height: 16px;
    }
  }
}
</style>
<style>
::view-transition-old(root),
::view-transition-new(root) {
  animation: none;
  mix-blend-mode: normal;
}

/* 进入dark模式和退出dark模式时，两个图像的位置顺序正好相反 */
.dark::view-transition-old(root) {
  z-index: 1;
}
.dark::view-transition-new(root) {
  z-index: 999;
}

::view-transition-old(root) {
  z-index: 999;
}
::view-transition-new(root) {
  z-index: 1;
}
</style>
