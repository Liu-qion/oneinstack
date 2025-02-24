<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Api } from '@/api/Api'

const getData = async () => {
  multipleSelection.value = []
  try {
    const { data: res } = await Api.getPlanTaskLog({
      id: 5,
    //   pageSize: pagination.pageSize,
    //   q: searchValue.value
    })
    console.log(res, 'res')
    if (res) {
      tableData.value = res.data || []
      pagination.total = res.total || 0
    } else {
      ElMessage.error(res?.message || '获取数据失败')
    }
  } catch (error) {
    ElMessage.error('获取数据失败')
    tableData.value = []
    pagination.total = 0
  }
}
onMounted(() => {
  getData()
})
</script>

<template>
  <div class="task-container">
    <el-card>
        <!-- <el-table :data="tableData" border style="width: 100%" empty-text="暂无数据" :row-key="(row) => row.id">
        <el-table-column prop="direction" label="方向" width="100">
          <template #default="scope">
            <div style="display: flex; flex-direction: row; align-items: center; cursor: pointer">
              <el-tag type="primary"  v-if="scope.row.direction == 'in'">输入</el-tag>
              <el-tag type="warning"  v-if="scope.row.direction == 'out'">输出</el-tag>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="protocol" label="协议" width="120" />
        <el-table-column prop="state" label="状态" width="120">
          <template #default="scope">
            <div style="display: flex; flex-direction: row; align-items: center; cursor: pointer">
              <a style="color: #64ffc9;" v-if="scope.row.status ==='success'">成功</a>
              <a style="color: #ff8888;" v-if="scope.row.state === 'failed'">失败</a>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="ports" label="开始时间"  width="180" />
        <el-table-column prop="output" label="结束时间" />
        <el-table-column prop="ips" label="来源" />
      
        <el-table-column prop="action" label="操作">
          <template #default="scope">
            <el-button link type="primary" size="small" @click="handleSet(scope.row)">设置</el-button>  
            <el-button link type="primary" size="small" @click="handleDelete(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table> -->
    </el-card>
    
  </div>
</template>

<style scoped lang="less">
。task-container{
    padding: 26px;
}
</style>