<template>
  <el-drawer
    v-model="drawer"
    :direction="direction"
    size="50%"
    :before-close="handleClose"
  >
  <template #header>
      <div class="drawer-header" style="padding: 20px;font-size: 16px;">
        <span class="title">{{ type ? '添加计划任务' : '修改计划任务' }}</span>
      </div>
    </template>
    <el-form
      ref="ruleFormRef"
      :model="ruleForm"
      :rules="rules"
      label-width="120px"
      style="padding: 0 20px"
    >
      <el-form-item label="任务类型" prop="taskType" required>
        <el-select v-model="ruleForm.taskType" placeholder="请选择任务类型">
          <el-option label="Shell脚本" value="shell" />
          <el-option label="备份应用" value="backup_app" />
          <el-option label="备份网站" value="backup_website" />
          <el-option label="备份数据库" value="backup_database" />
        </el-select>
      </el-form-item>

      <el-form-item label="任务名称" prop="name" required>
        <el-input v-model="ruleForm.name" placeholder="请输入任务名称" />
      </el-form-item>

     

      <el-form-item label="脚本内容" prop="script" v-if="ruleForm.taskType === 'shell'">
          <el-input
            v-model="ruleForm.script"
            type="textarea"
            :rows="6"
            placeholder="请输入Shell脚本内容" 
            class="code-editor"
            :spellcheck="false"
            font-family="monospace"
            style="font-family: Monaco, Menlo, Consolas, 'Courier New', monospace;
                   background-color: #1e1e1e;
                   color: #d4d4d4;
                   border-radius: 4px;
                   padding: 8px;"
            @input="(val) => {
              // 只允许输入代码字符
              ruleForm.script = val.replace(/[^\x00-\x7F\n\r\t]/g, '');
            }"
          />
      </el-form-item>
      <el-form-item label="执行周期" prop="cycle" required>
        <div v-for="(cycle, index) in ruleForm.cycles" :key="index" class="cycle-row">
          <el-select v-model="cycle.type" @change="(val) => handleCycleChange(val, index)" placeholder="请选择执行周期" style="width: 100px;">
            <el-option label="每分钟" value="minute" />
            <el-option label="每小时" value="hour" />
            <el-option label="每天" value="day" />
            <el-option label="每周" value="week" />
            <el-option label="每月" value="month" />
            <el-option label="自定义" value="custom" />
          </el-select>

          <div class="cycle-inputs" :class="cycle.type">
            <template v-if="cycle.type === 'month'">
              <div v-for="(monthTime, mIndex) in cycle.monthTimes" :key="mIndex" class="time-row">
                <div class="time-input-group">
                  <el-input-number v-model="monthTime.day" :min="1" :max="31" placeholder="日" >
                  <template #suffix>
        <span>日</span>
      </template>
    </el-input-number>
                </div>
                <div class="time-input-group">
                  <el-input-number v-model="monthTime.hour" :min="0" :max="23" placeholder="时" >
                    000
                  <template #suffix>
        <span>时</span>
      </template>
    </el-input-number>
                </div>
                <div class="time-input-group">
                  <el-input-number v-model="monthTime.minute" :min="0" :max="59" placeholder="分" >
                    <template #suffix>
                      <span>分</span>
                    </template>
                  </el-input-number>
                </div>
                <div class="time-actions">
                  <el-button type="primary" link @click="addMonthTime(index)" v-if="mIndex === cycle.monthTimes.length - 1">
                    <el-icon><Plus /></el-icon>
                  </el-button>
                  <el-button type="danger" link @click="removeMonthTime(index, mIndex)" v-if="cycle.monthTimes.length > 1">
                    <el-icon><Delete /></el-icon>
                  </el-button>
                </div>
              </div>
            </template>

            <template v-if="cycle.type === 'week'">
              <div v-for="(weekTime, wIndex) in cycle.weekTimes" :key="wIndex" class="time-row">
                <el-select v-model="weekTime.day" placeholder="星期">
                  <el-option v-for="i in 7" :key="i" :label="`星期${i}`" :value="i" />
                </el-select>
                <el-input-number v-model="weekTime.hour" :min="0" :max="23" placeholder="时" style="width: 190px;"/>
                <el-input-number v-model="weekTime.minute" :min="0" :max="59" placeholder="分" />
                <div class="time-actions">
                  <el-button type="primary" link @click="addWeekTime(index)" v-if="wIndex === cycle.weekTimes.length - 1">
                    <el-icon><Plus /></el-icon>
                  </el-button>
                  <el-button type="danger" link @click="removeWeekTime(index, wIndex)" v-if="cycle.weekTimes.length > 1">
                    <el-icon><Delete /></el-icon>
                  </el-button>
                </div>
              </div>
            </template>

            <template v-if="cycle.type === 'day'">
              <el-input-number v-model="cycle.dayHour" :min="0" :max="23" placeholder="时" />
              <el-input-number v-model="cycle.dayMinute" :min="0" :max="59" placeholder="分" />
            </template>

            <template v-if="cycle.type === 'hour'">
              <el-input-number v-model="cycle.hourMinute" :min="0" :max="59" placeholder="分" />
            </template>

            <template v-if="cycle.type === 'custom'">
              <el-input v-model="cycle.customCron" placeholder="请输入Cron表达式" />
            </template>
          </div>

          <div class="cycle-actions">
            <el-button type="primary" link @click="addCycle" v-if="index === ruleForm.cycles.length - 1">
              <el-icon><Plus /></el-icon>
            </el-button>
            <el-button type="danger" link @click="removeCycle(index)" v-if="ruleForm.cycles.length > 1">
              <el-icon><Delete /></el-icon>
            </el-button>
          </div>
        </div>
      </el-form-item>

      <el-form-item label="保存数量" prop="saveCount" v-if="ruleForm.taskType.includes('backup')">
        <el-input-number v-model="ruleForm.saveCount" :min="1" :max="100" />
        <span class="tip-text">设置保留的备份文件数量</span>
      </el-form-item>

    </el-form>

    <template #footer>
      <div style="flex: auto; text-align: right; padding: 0 20px">
        <el-button @click="handleClose">取消</el-button>
        <el-button type="primary" @click="handleSubmit">确定</el-button>
      </div>
    </template>
  </el-drawer>
</template>

<script setup lang="ts">
import { ref, reactive, watch } from 'vue'
import type { FormInstance, FormRules } from 'element-plus'
import { ElMessage } from 'element-plus'

const props = defineProps<{
  modelValue: boolean
  type?: boolean // true为新增,false为修改
  formData?: any
}>()

const emit = defineEmits(['update:modelValue', 'success'])

const drawer = ref(false)
const direction = ref<'rtl' | 'ltr' | 'ttb' | 'btt'>('rtl')
const ruleFormRef = ref<FormInstance>()

const ruleForm = reactive({
  taskType: 'shell',
  name: '',
  cycles: [{
    type: 'day',
    monthTimes: [{
      day: 1,
      hour: 0,
      minute: 0
    }],
    weekTimes: [{
      day: 1,
      hour: 0,
      minute: 0
    }],
    dayHour: 0,
    dayMinute: 0,
    hourMinute: 0,
    hourInterval: 2,
    customCron: ''
  }],
  saveCount: 7
})

const rules = reactive<FormRules>({
  taskType: [{ required: true, message: '请选择任务类型', trigger: 'change' }],
  name: [{ required: true, message: '请输入任务名称', trigger: 'blur' }],
  cycle: [{ required: true, message: '请设置执行周期', trigger: 'blur' }],
  script: [{ required: true, message: '请输入脚本内容', trigger: 'blur' }],
  saveCount: [{ required: true, message: '请设置保存数量', trigger: 'change' }]
})

const handleCycleChange = (type: string | number | boolean | undefined, index: number) => {
  if (typeof type !== 'string') return
  const cycle = ruleForm.cycles[index]
  
  switch(type) {
    case 'minute':
      cycle.customCron = '* * * * *'
      break
    case 'hour':
      cycle.customCron = '0 * * * *'
      break
    case 'day':
      cycle.customCron = '0 0 * * *'
      break
    case 'week':
      cycle.customCron = '0 0 * * 0'
      break
    case 'month':
      cycle.customCron = '0 0 1 * *'
      break
  }
}

const handleClose = () => {
  drawer.value = false
  emit('update:modelValue', false)
}

const handleSubmit = async () => {
  if (!ruleFormRef.value) return
  
  await ruleFormRef.value.validate(async (valid) => {
    if (valid) {
      try {
        // 转换所有周期为 cron 表达式数组
        const cronExpressions = ruleForm.cycles.flatMap(cycle => {
          switch(cycle.type) {
            case 'hour':
              return [`${cycle.hourMinute} */${cycle.hourInterval} * * *`] // 每N小时X分钟
            
            case 'week':
              return cycle.weekTimes.map(time => 
                `${time.minute} ${time.hour} * * ${time.day}`) // 每周X日Y时Z分
            
            case 'day':
              return [`${cycle.dayMinute} ${cycle.dayHour} * * *`] // 每天X时Y分
            
            case 'month':
              return cycle.monthTimes.map(time => 
                `${time.minute} ${time.hour} ${time.day} * *`) // 每月X日Y时Z分
            
            case 'custom':
              return [cycle.customCron]
            
            default:
              return []
          }
        })

        // 打印生成的 cron 表达式数组
        console.log('Cron expressions:', cronExpressions)

        // TODO: 调用API保存数据，传递 cronExpressions
        const submitData = {
          ...ruleForm,
          cronExpressions
        }
        
        ElMessage.success(props.type ? '添加成功' : '修改成功')
        emit('success')
        handleClose()
      } catch (error) {
        ElMessage.error('操作失败')
      }
    }
  })
}

// 监听父组件传值
watch(() => props.modelValue, (val) => {
  drawer.value = val
})

// 监听formData变化,用于修改时回填数据
watch(() => props.formData, (val) => {
  if (val && !props.type) {
    Object.assign(ruleForm, val)
  }
}, { immediate: true })

// 添加新的周期行
const addCycle = () => {
  ruleForm.cycles.push({
    type: 'day',
    monthTimes: [{
      day: 1,
      hour: 0,
      minute: 0
    }],
    weekTimes: [{
      day: 1,
      hour: 0,
      minute: 0
    }],
    dayHour: 0,
    dayMinute: 0,
    hourMinute: 0,
    hourInterval: 2,
    customCron: ''
  })
}

// 删除周期行
const removeCycle = (index: number) => {
  ruleForm.cycles.splice(index, 1)
}

// 添加月时间
const addMonthTime = (cycleIndex: number) => {
  ruleForm.cycles[cycleIndex].monthTimes.push({
    day: 1,
    hour: 0,
    minute: 0
  })
}

// 删除月时间
const removeMonthTime = (cycleIndex: number, timeIndex: number) => {
  ruleForm.cycles[cycleIndex].monthTimes.splice(timeIndex, 1)
}

// 添加周时间
const addWeekTime = (cycleIndex: number) => {
  ruleForm.cycles[cycleIndex].weekTimes.push({
    day: 1,
    hour: 0,
    minute: 0
  })
}

// 删除周时间
const removeWeekTime = (cycleIndex: number, timeIndex: number) => {
  ruleForm.cycles[cycleIndex].weekTimes.splice(timeIndex, 1)
}
</script>

<style scoped lang="less">
.cycle-row {
  display: flex;
  align-items: center;
  margin-bottom: 16px;
  gap: 16px;

  .cycle-inputs {
    display: flex;
    gap: 8px;
    margin-top: 8px;

    .el-input-number,
    .el-select {
      width: 120px;
    }
  }

  .cycle-actions {
    display: flex;
    gap: 8px;
  }
}

.tip-text {
  color: #909399;
  font-size: 12px;
  margin-left: 10px;
}

.time-row {
  display: flex;
  align-items: center;
  margin-bottom: 8px;
  gap: 8px;

  &:last-child {
    margin-bottom: 0;
  }

  .time-actions {
    display: flex;
    gap: 8px;
  }
}
</style>
