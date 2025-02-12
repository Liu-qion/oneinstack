<script setup lang="ts">
import { onMounted, reactive } from 'vue'
import { ConfProps } from './index.vue'
import { Api } from '@/api/Api'

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

onMounted(() => {
  !conf.server.options.length && conf.server.getOptions()
})
</script>

<template>
  <div class="container">
    <div class="search">
      <div class="btn">
        <el-row :gutter="10" class="mb8">
          <el-col :span="1.5">
            <el-button style="height: 100%" type="primary" @click="conf.drawer.open('add')">添加key</el-button>
          </el-col>
          <el-col :span="1.5">
            <div class="btnItem">
              <span>远程服务器</span>
            </div>
          </el-col>
          <el-col :span="1.5">
            <div class="btnItem">
              <span>备份列表</span>
            </div>
          </el-col>
          <el-col :span="1.5">
            <div class="btnItem">
              <span>清空数据库</span>
            </div>
          </el-col>
        </el-row>
      </div>
      <div class="demo-form-inline">
        <el-select v-model="conf.list.params.id" style="width: 200px" @change="conf.server.onChange">
          <el-option v-for="item in conf.server.options" v-bind="item" />
        </el-select>
      </div>
    </div>
    <div class="box2">
      <el-tabs v-if="conf.dbList.data.length" v-model="conf.list.params.r_db" @tab-click="conf.list.getData">
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
      />
    </div>
  </div>
</template>

<style scoped lang="less"></style>
