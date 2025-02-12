<script setup lang="ts">
import { reactive, useTemplateRef } from 'vue'
import LoginContainer from './components/login-container.vue'
import System from '@/utils/System'
import { Api } from '@/api/Api'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import sconfig from '@/sstore/sconfig'

const conf = reactive({
  form: {
    username: '',
    password: '',
    remember: false
  },
  loading: false,
  formRef: useTemplateRef<FormInstance>('formRef'),
  rules: {
    username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
    password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
  } as FormRules,
  handleLogin: () => {
    conf.formRef?.validate(async valid => {
      if (!valid) return
      conf.loading = true
      const { data: res } = await Api.login({
        ...conf.form,
        final: () => conf.loading = false
      })
      ElMessage.success('登录成功')
      sconfig.login(res)
      setTimeout(() => {
        System.router.push('/')
      }, 500)
    })
  },
  handleToAppLogin: () => {
    System.router.push('/login/scan')
  }
})
</script>

<template>
  <login-container :loading="conf.loading" currentActive="login">
    <template #default="{ className }">
      <div class="logo-box flex justify-center fit-width">
        <v-s-icon class="login-content-right__main-logo" name="logo" width="175" height="144" />
      </div>
      <el-form ref="formRef" :model="conf.form" :rules="conf.rules" class="fit-width">
        <el-form-item :class="className.formItemGap" prop="username">
          <el-input v-model="conf.form.username" placeholder="输入您的账户/手机号/邮箱号" clearable>
            <template #prefix>
              <v-s-icon name="user" size="30" />
            </template>
          </el-input>
        </el-form-item>
        <el-form-item prop="password">
          <el-input v-model="conf.form.password" placeholder="输入您的密码" type="password" clearable>
            <template #prefix>
              <v-s-icon name="password" size="30" />
            </template>
          </el-input>
        </el-form-item>
        <div class="flex justify-between">
          <el-checkbox v-model="conf.form.remember" label="记住密码" />
        </div>
        <button :class="className.loginBtn" @click.prevent="conf.handleLogin">登录</button>
      </el-form>
      <div :class="className.other">
        <el-divider class="login-content-right__main-other-title">
          <el-link @click="conf.handleToAppLogin">App扫码登录</el-link>
        </el-divider>
      </div>
    </template>
  </login-container>
</template>

<style scoped lang="less">
@primary-color: #f7911c;
@font-gray: #a2a2a2;
@border-gray: #e3e3e3;
@bg-color: #f5f5f5;

.login-content-right__main {
  &-logo {
    margin-bottom: 116px;
  }

  &-other {
    &-btn {
      margin-top: 14px;
      display: flex;
      justify-content: space-between;
      align-items: center;
      gap: 14px;
      font-size: 16px;
      color: #2b2b2b;
      cursor: pointer;

      .vsicon {
        border-radius: 50%;
        transition: box-shadow 0.3s;
      }

      &:hover {
        color: @primary-color;

        .vsicon {
          box-shadow: 0 0 10px @primary-color;
        }
      }
    }
  }
}
</style>
