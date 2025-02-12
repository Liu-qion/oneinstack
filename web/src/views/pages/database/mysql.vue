<script setup lang="ts">
import sapp from '@/sstore/sapp'
import System from '@/utils/System'
import { CircleClose, ArrowDown, Setting, Hide, View, CopyDocument } from '@element-plus/icons-vue'
import type { ConfProps } from './index.vue'
import { onMounted } from 'vue';

const { conf } = defineProps<ConfProps>()

onMounted(() => {
  conf.list.getData()
})
</script>

<template>
  <div v-if="conf.showTips" class="tip">
    <div class="flex items-center fit-width">
      <v-s-icon name="warning" size="22" :color="conf.themeColor[sapp.theme]" />
      <span class="ellipsis" style="margin-left: 32px; max-width: 85%">
        请在添加数据库后，务必到[
        <span style="color: var(--el-color-primary)">计划任务</span>
        ]添加定时备份任务，以确保您的数据安全。温馨提示：通过第三方或者MySQL命令行创建的数据库需要点击"从服务器获取"才能在计划任务中备份
      </span>
    </div>
    <el-icon class="cursor-pointer" size="26" color="#A2A2A2" @click="conf.showTips = false"><CircleClose /></el-icon>
  </div>
  <div class="container">
    <div class="search">
      <div class="btn">
        <el-row :gutter="10" class="mb8">
          <el-col :span="1.5">
            <el-button style="height: 100%" type="primary" @click="conf.drawer.open('add')">添加数据库</el-button>
          </el-col>
          <el-col :span="1.5">
            <div class="btnItem">
              <span>root密码</span>
            </div>
          </el-col>
          <el-col :span="1.5">
            <div class="btnItem" @click="System.router.push('/database/remote')">
              <span>远程数据库</span>
            </div>
          </el-col>
          <el-col :span="1.5">
            <div class="btnItem">
              <span>phpMyAdmin</span>
            </div>
          </el-col>
          <el-col :span="1.5">
            <div class="btnItem">
              <el-dropdown>
                <span class="el-dropdown-link">
                  高级设置
                  <el-icon class="el-icon--right"><arrow-down /></el-icon>
                </span>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item>修改默认页面</el-dropdown-item>
                    <el-dropdown-item>默认站点</el-dropdown-item>
                    <el-dropdown-item>PHP命令行版本</el-dropdown-item>
                    <el-dropdown-item>HTTPS防窜站</el-dropdown-item>
                    <el-dropdown-item>TLS设置</el-dropdown-item>
                    <el-dropdown-item>全局设置</el-dropdown-item>
                    <el-dropdown-item>关联数据库</el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </div>
          </el-col>
          <el-col :span="1.5">
            <div class="btnItem" style="cursor: pointer">
              <span style="font-size: 16px; margin-right: 8px">nignx</span>
              <v-s-icon name="triangle" size="16" style="transform: rotate(-90deg)" />
            </div>
          </el-col>
          <el-col :span="1.5">
            <div class="btnItem">
              <el-dropdown>
                <span class="el-dropdown-link">
                  全部
                  <el-icon class="el-icon--right"><arrow-down /></el-icon>
                </span>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item>分类1</el-dropdown-item>
                    <el-dropdown-item>分类2</el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </div>
          </el-col>
        </el-row>
      </div>
      <div class="demo-form-inline">
        <search-input
          v-model="conf.list.params.name"
          placeholder="请输入数据库名称"
          style="margin-right: 18px"
          @search="conf.list.getData"
        />
        <el-button :icon="Setting" type="primary" />
      </div>
    </div>
    <div class="box2">
      <custom-table
        v-model:page="conf.list.params.page"
        :loading="conf.list.loading"
        :data="conf.list.data"
        :auto-pagination="false"
        :total="conf.list.total"
        :page-size="conf.list.params.pageSize"
        :columns="conf.list.columns"
        @update:page="conf.list.getData"
      >
        <template #empty>
          <div style="margin-top: 40px">
            <span>
              您的数据库列表为空，您可以
              <a
                class="cursor-pointer"
                style="color: var(--el-color-primary); text-decoration: underline"
                @click="conf.drawer.open('add')"
              >
                添加一个数据库
              </a>
            </span>
          </div>
        </template>
        <template #password="{ row, index }">
          <div class="flex items-center" style="gap: 8px">
            <span>
              {{ conf.password.show[index] ? row.password : ''.padEnd(row.password.length, '*') }}
            </span>
            <el-icon
              class="hover-opacity cursor-pointer"
              :size="18"
              @click="conf.password.show[index] = !conf.password.show[index]"
            >
              <View v-if="conf.password.show[index]" />
              <Hide v-else />
            </el-icon>
            <el-icon class="hover-opacity cursor-pointer" @click="conf.password.copy(row.password)">
              <CopyDocument />
            </el-icon>
          </div>
        </template>
      </custom-table>
    </div>
  </div>
</template>

<style scoped lang="less">
.database-container {
  .tip {
    width: 100%;
    height: 58px;
    background: rgb(var(--bg-card-color));
    margin-top: 20px;
    padding: 18px 28px;
    display: flex;
    align-items: center;
    justify-content: space-between;

    span {
      font-size: 16px;
      color: var(--font-color-gray-light);
    }
  }
}
</style>
