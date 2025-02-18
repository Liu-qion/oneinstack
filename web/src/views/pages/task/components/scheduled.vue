<script setup lang="ts">
import { ref, reactive } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Refresh, Setting, ArrowDown } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'
import { Api } from '@/api/Api'
import { onMounted } from 'vue'
import AddTask from './add-task.vue'

interface RuleForm {
  name: string
  region: string
  count: string
  desc: string
}

const tableData = ref([])
const getData = async () => {
  try {
    const { data: res } = await Api.getPlanTaskList({
      page: pagination.currentPage,
      pageSize: pagination.pageSize,
      q: searchValue.value
    })
    console.log(res,'res')
    if (res ) {  // 确保请求成功
      tableData.value = res.data || []  // 更新表格数据
      pagination.total = res.total || 0  // 更新总数
    } else {
      ElMessage.error(res?.message || '获取数据失败')
    }
  } catch (error) {
    ElMessage.error('获取数据失败')
    tableData.value = []
    pagination.total = 0
  }
}
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

const pagination = reactive({
  currentPage: 1,
  pageSize: 10,
  total: 0
})

const filterDirection = ref('')
const searchValue = ref('')

const addTaskVisible = ref(false)
const addTask = () => {
  addTaskVisible.value = true
}

onMounted(() => {
  getData()
})
</script>

<template>
  <div class="container">
    <div class="tool-bar">
      <el-space class="btn-group">
        <el-button type="primary" @click="addTask">添加任务</el-button>
        <el-button type="primary">执行任务</el-button>
        <el-button type="primary">启动任务</el-button>
        <el-button type="primary">停止任务</el-button>
        <el-button type="primary">设置任务</el-button>
        <el-button type="primary">删除任务</el-button>
      </el-space>
      <div class="demo-form-inline flex">
        <el-dropdown>
          <el-button type="primary" class="mr-2">
            <span class="el-dropdown-link">
              全部分类
              <el-icon class="el-icon--right"><arrow-down /></el-icon>
            </span>
          </el-button>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item>分类</el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
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
    <AddTask v-model="addTaskVisible" :type="true" />
  </div>
</template>

<style scoped lang="less"></style>
