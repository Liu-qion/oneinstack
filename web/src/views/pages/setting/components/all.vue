<script setup lang="ts">
import panelSetting from './panel-setting.vue'
import securitySetting from './security-setting.vue'
import { Api } from '@/api/Api'
import { onMounted, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'

let allinfo = ref('')
const getSystemInfo = async () => {
  try {
    const { data: res } = await Api.getSystemInfo()
    allinfo.value = res
    // if (res) {
    //   // 更新settingData的值
    //   conf.settingData.forEach(item => {
    //     console.log(item,'item')
    //     item.value = res[item.prop] ? res[item.prop]  :res.user[item.prop]
    //   })
    // }
  } catch (error) {
    ElMessage.error('获取系统信息失败')
  }
}

onMounted(() => {
  getSystemInfo()
})
</script>

<template>
  <div class="all-container" v-if="allinfo">
    <panel-setting v-if="allinfo" :isCard="false" :allinfo="allinfo" />
    <security-setting v-if="allinfo" :isCard="false" :allinfo="allinfo" />
  </div>
</template>

<style scoped lang="less">
.all-container {
  background: rgb(var(--bg-card-color));
  border-radius: 5px;
  padding: 26px 30px;
}
</style>
