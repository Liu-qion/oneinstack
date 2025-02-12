<script setup lang="ts">
import { reactive, useTemplateRef } from 'vue'
import LoginContainer from './components/login-container.vue'
import System from '@/utils/System'
import { Api } from '@/api/Api'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import sconfig from '@/sstore/sconfig'

const conf = reactive({
  loading: false
})
</script>

<template>
  <login-container :loading="conf.loading" currentActive="login">
    <div class="back-btn flex items-center" @click="System.router.back()">
      <v-s-icon name="arrow-down" :color="['#6f6f6f']" size="24" style="transform: rotate(90deg); margin-right: 16px" />
      <span>返回</span>
    </div>
    <div class="code-box">
      <div class="code"></div>
      <div class="code-box-title">
        <span>用App扫描当前二维码进行登录</span>
      </div>
    </div>
  </login-container>
</template>

<style scoped lang="less">
.back-btn {
  position: absolute;
  top: 20px;
  left: 20px;
  cursor: pointer;
  font-weight: 500;
  font-size: 24px;
  color: #1c2a38;

  &:hover {
    opacity: 0.8;
  }
}

.code-box {
  width: 328px;
  height: 316px;
  background: #ffffff;
  border-radius: 14px;
  position: relative;
  display: flex;
  justify-content: center;
  align-items: center;

  .code {
    width: 300px;
    height: 300px;
  }

  &-title {
    width: 100%;
    position: absolute;
    bottom: -42px;
    left: 50%;
    transform: translateX(-50%);
    font-weight: 500;
    font-size: 18px;
    color: #1c2a38;
    text-align: center;
  }
}
</style>
