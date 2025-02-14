<script setup lang="ts">
import { reactive } from 'vue'
import { ConfProps } from './index.vue'
import { Api } from '@/api/Api'
import { WarningFilled } from '@element-plus/icons-vue'
import System from '@/utils/System'

const { conf: parentConf } = defineProps<ConfProps>()

const conf = reactive({
  ...parentConf,
  server: {
    options: [] as { label: string; value: number }[],
    getOptions: async () => {
      const { data } = await Api.getConnlist(conf.list.params)
      conf.server.options = data.map((item: any) => ({
        label: item.remark ? `${item.remark}(${item.addr})` : item.addr,
        value: item.id
      }))
    },
    onChange: (value: number) => {
      conf.list.params.id = value
      conf.dbList.params.id = value
      conf.dbList.getData()
      handleTabClick({ paneName: conf.dbList.data[0] })
    }
  },
  dbList: {
    params: {
      page: 1,
      pageSize: 999,
      id: 0,
      type: 'redis'
    },
    data: [] as any[],
    getData: async () => {
      const { data: res } = await Api.getDatabaseList(conf.list.params)
      conf.dbList.data = res.data
    }
  }
})

conf.list.loading = false
!conf.server.options.length && conf.server.getOptions()

const handleTabClick = ({ paneName }: { paneName: string | number | undefined }) => {
  conf.list.params.r_db = paneName
  conf.list.getData()
}
</script>

<template>
  <div class="container">
    <div class="tool-bar">
      <el-space class="btn-group">
        <!-- <el-button type="primary" @click="conf.drawer.open('add')">添加key</el-button> -->
        <el-button type="primary" @click="System.router.push('/database/remote')">远程服务器</el-button>
        <!-- <el-button type="primary">备份列表</el-button>
        <el-button type="primary">清空数据库</el-button> -->
      </el-space>
      <div class="demo-form-inline flex" style="gap: 16px">
        <span class="flex items-center" style="color: var(--el-color-primary); gap: 8px">
          <el-icon :size="18"><WarningFilled /></el-icon>
          当前所有操作都关联至
        </span>
        <el-select
          v-model="conf.list.params.id"
          placeholder="请选择一个服务器"
          style="width: 200px"
          @change="conf.server.onChange"
        >
          <el-option v-for="item in conf.server.options" v-bind="item" />
        </el-select>
      </div>
    </div>
    <div class="box2">
      <el-tabs v-if="conf.dbList.data.length" v-model="conf.list.params.r_db" @tab-click="handleTabClick">
        <el-tab-pane v-for="item in conf.dbList.data" :key="item.index" :label="`DB${item.name}`" :name="item.name" />
      </el-tabs>
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
          <div class="no-data">
            <img src="/static/images/empty.webp" alt="" />
            <span>暂无数据</span>
          </div>
        </template>
      </custom-table>
    </div>
  </div>
</template>

<style scoped lang="less">
.no-data {
  width: 100%;
  display: flex;
  flex-flow: column nowrap;
  justify-content: center;
  align-items: center;
  color: var(--font-color-gray-light);
}
</style>
