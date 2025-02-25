<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Api } from '@/api/Api'

let tableData = ref([])

const getData = async () => {

  try {
    const { data: res } = await Api.getPlanTaskLog({
      id: 5,
      //   pageSize: pagination.pageSize,
      //   q: searchValue.value
    })
    console.log(res, 'res')
    tableData.value = res.data
  } catch (error) {

  }
}
const formatDate = (dateStr: string) => {
  if (!dateStr) return ''
  const date = new Date(dateStr)
  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  const seconds = String(date.getSeconds()).padStart(2, '0')
  return `${year}年${month}月${day}日 ${hours}:${minutes}:${seconds}`
}
onMounted(() => {
  getData()
})
</script>
<template>
  <div class="task-container">
    <el-card>
      <el-table :data="tableData" border style="width: 100%" empty-text="暂无数据" :row-key="(row) => row.id">
      
          <el-table-column prop="status" label="状态" width="120">
            <template #default="scope">
              <div style="display: flex; flex-direction: row; align-items: center; cursor: pointer">
                <a style="color: #64ffc9;" v-if="scope.row.status ==='success'">成功</a>
                <a style="color: #ff8888;" v-if="scope.row.status === 'failed'">失败</a>
              </div>
            </template>
            </el-table-column>
            <el-table-column prop="start_time" label="执行时间"  >
              <template #default="scope">
                <div style="display: flex; flex-direction: row; align-items: center; cursor: pointer">
                  <span>{{ formatDate(scope.row.start_time) }}</span> -- <span>{{ formatDate(scope.row.end_time) }}</span>
                </div>
              </template>
            </el-table-column>
           
            <el-table-column prop="output" label="输出日志" />
           
</el-table> 
</el-card>
  </div>
</template>
<style scoped lang="less">
。task-container {
  padding: 26px;
}
</style>