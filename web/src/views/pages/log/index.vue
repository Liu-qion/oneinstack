<script setup lang="ts">
import CardTabs from '@/components/card-tabs.vue'
import { reactive } from 'vue'
import Panel from './components/panel.vue'

const conf = reactive({
  list: [
    {
      name: '面板日志',
      index: 0,
      component: Panel
    },
    {
      name: '网站日志',
      index: 1
    },
    {
      name: '日志审计',
      index: 2
    },
    {
      name: 'SSH登录日志',
      index: 3
    },
    {
      name: '软件日志',
      index: 4
    }
  ],
  activeIndex: 0,
  clickActive: (item: any) => {
    conf.activeIndex = item.index
  }
})
</script>

<template>
  <div class="log-container">
    <card-tabs :list="conf.list" :active-index="conf.activeIndex" :click-active="conf.clickActive" />
    <component :is="conf.list[conf.activeIndex].component" />
  </div>
</template>

<style scoped lang="less">

</style>
