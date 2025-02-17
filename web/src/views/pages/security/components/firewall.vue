<script setup lang="ts">
import SearchInput from '@/components/search-input.vue'
import { ref, reactive } from 'vue'
import { Api } from '@/api/Api'
import { onMounted } from 'vue'
import Addfirewall from './addfirewall.vue'
import { ElMessage } from 'element-plus'



const addRulemodal = ref(false)
const exportRulemodal = ref(false)
const keyNum = ref(0)
const value1 = ref(false)
const value2 = ref(false)



const tableData = ref([
  // {
  //   date: 'www.baidu.com',
  //   status: 1,
  //   address: 'No. 189, Grove St, Los Angeles'
  // },
  // {
  //   date: '2016-05-02',
  //   status: 2,
  //   address: 'No. 189, Grove St, Los Angeles'
  // },
  // {
  //   date: '2016-05-04',
  //   status: 1,
  //   address: 'No. 189, Grove St, Los Angeles'
  // },
  // {
  //   date: '2016-05-01',
  //   status: 1,
  //   address: 'No. 189, Grove St, Los Angeles'
  // }
])
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
  getData()
}
const oncloseExport = (value: boolean) => {
  exportRulemodal.value = value
  keyNum.value++
}
const handleAllDirection = (value: string) => {
  
}
const handleOpenPing =async (value: string) => {
  if(value == 'ping'){
    const { data: res } = await Api.openPing('')
    if (res.code !=0){
      ElMessage.error(res.msg)
    }
  }else{
    const { data1: res } = await Api.stopFirewall('')
    if (res.code !=0){
      ElMessage.error(res.msg)
    }
  }

}

const getList = () => {}
const onSubmit = () => {
  console.log('submit!')
}
const getData = async () => {
  const { data: res } = await Api.getFirewallRule('')
  tableData.value = res.data
  console.log(res,'res')
}
const getFirewallInfo = async () => {
  const { data: res } = await Api.getFirewallInfo({})

  value1.value =res.info.enabled ? res.info.enabled : false
  value2.value =res.info.pingBlocked ? res.info.pingBlocked : false
}
onMounted(() => {
  getData()
  getFirewallInfo()
})
</script>

<template>
  <div>
    <div class="box1">
      <div class="flex justify-between items-center">
        <div>
          <span>防火墙开关</span>
          <el-switch v-model="value1" class="ml-2"  @change="handleOpenPing('wall')"/>
        </div>
        <div style="margin-inline: 60px 150px">
          <span>禁ping</span>
          <el-switch v-model="value2" class="ml-2"  @change="handleOpenPing('ping')"/>
        </div>
        <div class="webLog">
          Web日志：
          <span class="mr-2" style="color: var(--el-color-primary)">/www/wwwlogs</span>
          5.20MB
        </div>
      </div>
      <el-button>清空</el-button>
    </div>
    <div class="tool-bar">
      <el-space class="btn-group">
        <el-button type="primary" @click="handleAdd">添加端口规则</el-button>
        <!-- <el-button type="primary" @click="exportClick">导入规则</el-button>
        <el-button type="primary" @click="exportClick">导出规则</el-button> -->
        <el-button-group>
          <el-button type="primary" @click="handleAllDirection('')">所有方向</el-button>
          <el-button type="primary"  @click="handleAllDirection('in')" >入站</el-button>
          <el-button type="primary"  @click="handleAllDirection('out')">出站</el-button>
        </el-button-group>
      </el-space>
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

  <Addfirewall v-model="addRulemodal" @close="oncloseRule" style="padding: 20px;" :key="keyNum" v-if="addRulemodal" />
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
