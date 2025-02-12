<script setup lang="ts">
import { ref, reactive } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Refresh, Setting, ArrowDown } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'

interface RuleForm {
  name: string
  region: string
  count: string
  desc: string
}

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
const category = ref(['传统项目', 'swoole异步项目', 'thinkphp异步项目', '异步项目', '一键部署', '批量创建'])

const formInline = reactive({
  user: '',
  region: '',
  date: ''
})

const onSubmit = () => {
  console.log('submit!')
}

const statusClick = () => {
  ElMessageBox.confirm('计划任务暂停后将无法继续运行，您真的要停用这个计划任务吗？', '设置计划任务状态', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  })
    .then(() => {
      ElMessage({
        type: 'success',
        message: '退出成功'
      })
    })
    .catch(() => {
      ElMessage({
        type: 'info',
        message: '取消退出'
      })
    })
}

const ruleFormRef = ref<FormInstance>()
const ruleForm = reactive<RuleForm>({
  name: 'Hello',
  region: '',
  count: '',
  desc: ''
})

const rules = reactive<FormRules<RuleForm>>({
  name: [
    { required: true, message: 'Please input Activity name', trigger: 'blur' },
    { min: 3, max: 5, message: 'Length should be 3 to 5', trigger: 'blur' }
  ],
  region: [
    {
      required: true,
      message: 'Please select Activity zone',
      trigger: 'change'
    }
  ],
  count: [
    {
      required: true,
      message: 'Please select Activity count',
      trigger: 'change'
    }
  ],
  desc: [{ required: true, message: '请输入域名', trigger: 'blur' }]
})

const submitForm = async (formEl: FormInstance | undefined) => {
  if (!formEl) return
  await formEl.validate((valid, fields) => {
    if (valid) {
      console.log('submit!')
    } else {
      console.log('error submit!', fields)
    }
  })
}
</script>

<template>
  <div class="container">
    <div class="search">
      <div class="btn">
        <el-row :gutter="10" class="mb8">
          <el-col :span="1.5">
            <el-button type="primary" style="height: 100%">添加任务</el-button>
          </el-col>
          <el-col :span="1.5">
            <div class="btnItem">
              <span>执行任务</span>
            </div>
          </el-col>
          <el-col :span="1.5">
            <div class="btnItem">
              <span>启动任务</span>
            </div>
          </el-col>
          <el-col :span="1.5">
            <div class="btnItem">
              <span>停止任务</span>
            </div>
          </el-col>
          <el-col :span="1.5">
            <div class="btnItem">
              <span>设置任务</span>
            </div>
          </el-col>
          <el-col :span="1.5">
            <div class="btnItem">
              <span>删除任务</span>
            </div>
          </el-col>
        </el-row>
      </div>
      <div class="demo-form-inline flex">
        <div class="btnItem mr-2">
          <el-dropdown>
            <span class="el-dropdown-link">
              全部分类
              <el-icon class="el-icon--right"><arrow-down /></el-icon>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item>分类</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
        <search-input placeholder="请输入域名或备注" style="margin-right: 18px" />
        <el-button :icon="Refresh" type="primary" @click="onSubmit" />
        <el-button :icon="Setting" type="primary" @click="onSubmit" />
      </div>
    </div>

    <div class="box2">
      <el-table class="fileTable" :data="tableData" border style="width: 100%">
        <el-table-column type="selection" width="55" />
        <el-table-column prop="date" label="任务名称" width="180"></el-table-column>
        <el-table-column prop="status" label="状态" width="180">
          <template #default="scope">
            <div style="display: flex; flex-direction: row; align-items: center; cursor: pointer">
              <a style="color: #64ffc9; text-decoration: underline" v-if="scope.row.status == 1" @click="statusClick()">
                运行中
              </a>
              <a style="color: #ff8888; text-decoration: underline" v-if="scope.row.status == 2">已停用</a>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="address" label="执行周期" />
        <el-table-column prop="address" label="保存数量" />
        <el-table-column prop="address" label="备份到" />
        <el-table-column prop="address" label="上次执行时间" />
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

<style scoped lang="less"></style>
