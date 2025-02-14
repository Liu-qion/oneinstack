<script setup lang="ts">
import SearchInput from '@/components/search-input.vue'
import { reactive } from 'vue'

const conf = reactive({
  activeName: 0,
  tabs: [
    {
      name: '操作日志',
      index: 0
    },
    {
      name: '运行日志',
      index: 1
    },
    {
      name: '计划任务日志',
      index: 2
    }
  ],
  drawer: false,
  tableData: [
    {
      date: 'www.baidu.com',
      status: 1,
      address: 'No. 189, Grove St, Los Angeles'
    },
    {
      date: '2016-05-02',
      status: 2,
      address: 'No. 189, Grove St, Los Angeles'
    },
    {
      date: '2016-05-04',
      status: 1,
      address: 'No. 189, Grove St, Los Angeles'
    },
    {
      date: '2016-05-01',
      status: 1,
      address: 'No. 189, Grove St, Los Angeles'
    }
  ],
  handleAdd: () => {
    conf.drawer = true
  },
  handleClick: () => {
    console.log('click')
  }
})
</script>

<template>
  <div>
    <div class="box2">
      <div class="category">
        <el-tabs v-model="conf.activeName" @tab-click="conf.handleClick">
          <el-tab-pane v-for="item in conf.tabs" :key="item.index" :label="item.name" :name="item.index" />
        </el-tabs>
      </div>
      <div class="tool-bar">
        <el-space :size="14" class="btn-group">
          <el-button type="primary" @click="conf.handleAdd">刷新日志</el-button>
          <el-button type="primary" @click="conf.handleAdd">清空日志</el-button>
          <el-button type="primary" @click="conf.handleAdd">IP操作统计</el-button>
        </el-space>
        <div class="demo-form-inline">
          <SearchInput placeholder="请输入搜索关键字" />
        </div>
      </div>

      <el-table :data="conf.tableData" border style="width: 100%">
        <el-table-column type="selection" width="55" />
        <el-table-column prop="date" label="网站名" width="180" />
        <el-table-column prop="status" label="状态" width="180">
          <template #default="scope">
            <div style="display: flex; flex-direction: row; align-items: center; cursor: pointer">
              <a style="color: #64ffc9; text-decoration: underline" v-if="scope.row.status == 1">运行中</a>
              <a style="color: #ff8888; text-decoration: underline" v-if="scope.row.status == 2">已停用</a>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="address" label="备份" />
        <el-table-column prop="address" label="跟目录" />
        <el-table-column prop="address" label="日流量" />
        <el-table-column prop="address" label="到期时间" />
        <el-table-column prop="address" label="备注" />
        <el-table-column prop="address" label="PHP" />
        <el-table-column prop="address" label="SSL证书" />
        <el-table-column prop="address" label="操作">
          <template #default>
            <el-button link type="primary" size="small">设置</el-button>
            <el-button link type="primary" size="small">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="pagination">
        <el-pagination background layout="prev, pager, next" :total="1000" />
      </div>
    </div>
  </div>
</template>

<style scoped lang="less">
.category {
  height: 60px;
  background: rgba(var(--category-item-bg-color), 0.6);
  border-radius: 4px;
}

:deep(.el-tabs__nav-wrap) {
  &::after {
    background: transparent;
  }
}
</style>
