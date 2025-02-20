<template>
  <x-route-event></x-route-event>
  <x-router-view :show="conf.show"></x-router-view>
  <div v-title data-title="我是标题"></div>
</template>

<script setup lang="ts">
import { nextTick, reactive } from 'vue'
import { ERouter } from './enum/Enum'
import sconfig from './sstore/sconfig'
import { Api } from '@/api/Api'

//获取网页标题并修改网页标题
const getSystemInfo = async () => {
  const authorization = localStorage.getItem('userInfo') ? JSON.parse(localStorage.getItem('userInfo') || '')?.token : '';
  try {
    if(authorization){
      const { data: res } = await Api.getSystemInfo()
      document.title = res.title+'' 
    }
  } catch (error) {
  }
}
defineOptions({
  name: 'App'
})

const conf = reactive({
  show: true,
  reload: () => {
    conf.show = false
    nextTick(() => {
      conf.show = true
    })
  }
})

CEvent.on(ERouter.browserShow, (isShow) => {
  if (sconfig.userInfo) {
    if (isShow) {
    } else {
    }
  }
})

CEvent.on(ERouter.reload, conf.reload)
getSystemInfo()
</script>
