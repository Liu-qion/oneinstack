<script setup lang="ts">
import SearchInput from '@/components/search-input.vue'
import { ref, reactive } from 'vue'

const addRulemodal = ref(false)
const exportRulemodal = ref(false)
const keyNum = ref(0)
const value2 = ref('')
const formInline = reactive({
  user: '',
  region: '',
  date: ''
})

const tableData = [
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
]
const handleAdd = () => {
  addRulemodal.value = true
  keyNum.value++
}
const exportClick = () => {
  exportRulemodal.value = true
  keyNum.value++
}
const oncloseRule = (value: boolean) => {
  addRulemodal.value = value
  keyNum.value++
}
const oncloseExport = (value: boolean) => {
  exportRulemodal.value = value
  keyNum.value++
}
const getList = () => {}
const onSubmit = () => {
  console.log('submit!')
}
</script>

<template>
  <div>
    <div class="box1">
      <div class="flex justify-between items-center">
        <div>
          <span>防火墙开关</span>
          <el-switch v-model="value2" class="ml-2" />
        </div>
        <div style="margin-inline: 60px 150px">
          <span>禁ping</span>
          <el-switch v-model="value2" class="ml-2" />
        </div>
        <div class="webLog">
          Web日志：
          <span class="mr-2" style="color: var(--el-color-primary)">/www/wwwlogs</span>
          5.20MB
        </div>
      </div>
      <el-button>清空</el-button>
    </div>
    <div class="search">
      <div class="btn">
        <el-row :gutter="10" class="mb8">
          <el-col :span="1.5">
            <el-button type="primary" style="height: 100%" @click="handleAdd">添加端口规则</el-button>
          </el-col>
          <el-col :span="1.5">
            <div class="btnItem" @click="exportClick">
              <span>导入规则</span>
            </div>
          </el-col>
          <el-col :span="1.5">
            <div class="btnItem">
              <span>导出规则</span>
            </div>
          </el-col>
          <el-col :span="1.5">
            <div class="flex">
              <div class="btnItem">
                <span>所有方向</span>
              </div>
              <div class="btnItem">
                <span>入站</span>
              </div>
              <div class="btnItem">
                <span>出站</span>
              </div>
            </div>
          </el-col>
        </el-row>
      </div>
      <div class="demo-form-inline">
        <search-input placeholder="请输入关键词进行搜索" />
      </div>
    </div>
    <div class="box2">
      <div class="miscellaneous">
        <div>端口规则：2</div>
        <div>端口转发：0</div>
        <div>地区规划：0</div>
        <div>IP规则：0</div>
      </div>
      <el-table :data="tableData" border style="width: 100%">
        <el-table-column type="selection" width="55" />
        <el-table-column prop="date" label="协议" width="180" />
        <el-table-column prop="status" label="端口" width="180">
          <template #default="scope">
            <div style="display: flex; flex-direction: row; align-items: center; cursor: pointer">
              <a style="color: #64ffc9; text-decoration: underline" v-if="scope.row.status == 1">运行中</a>
              <a style="color: #ff8888; text-decoration: underline" v-if="scope.row.status == 2">已停用</a>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="address" label="状态" />
        <el-table-column prop="address" label="策略" />
        <el-table-column prop="address" label="方向" />
        <el-table-column prop="address" label="来源" />
        <el-table-column prop="address" label="备注" />
        <el-table-column prop="address" label="时间" />
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
.miscellaneous {
  width: 100%;
  height: 60px;
  background: rgba(var(--category-item-bg-color), 0.6);
  border-radius: 4px;
  margin-bottom: 12px;
  padding-inline: 26px;
  display: flex;
  align-items: center;
  gap: 98px;
  color: var(--font-color-black);
  font-size: 14px;
}
</style>
