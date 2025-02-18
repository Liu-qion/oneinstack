<script setup lang="ts">
import { ref, reactive, defineEmits, watch } from 'vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { Api } from '@/api/Api'


interface RuleForm {
  port: string
}

// 将 defineProps 移到最前面
const props = defineProps<{
  modelValue: boolean,
  form: {
    port: string
  }
}>()

const emit = defineEmits(['update:modelValue'])
const dialogVisible = ref(false)
const drawer = ref(false)
const direction = ref('rtl')
const ruleFormRef = ref<FormInstance>()

const ruleForm = reactive<RuleForm>({
  port: props.form.port
})

// 监听父组件传入的form变化
watch(() => props.form.port, (val) => {
  ruleForm.port = val
})

// 端口范围校验规则
const validatePort = (rule: any, value: string, callback: any) => {
  const portRegex = /^([0-9]{1,5}|[0-9]{1,5}-[0-9]{1,5})$/
  if (!value) {
    callback(new Error('请输入端口'))
  } else if (!portRegex.test(value)) {
    callback(new Error('端口格式不正确,请输入1-65535之间的端口或端口范围'))
  } else {
    const ports = value.split('-')
    const port1 = parseInt(ports[0])
    const port2 = ports.length > 1 ? parseInt(ports[1]) : port1
    
    if (port1 < 1 || port1 > 65535 || port2 < 1 || port2 > 65535) {
      callback(new Error('端口范围必须在1-65535之间'))
    } else if (port2 < port1) {
      callback(new Error('结束端口不能小于起始端口'))
    } else {
      callback()
    }
  }
}

const rules = reactive<FormRules>({
  port: [{ validator: validatePort, trigger: 'blur' }]
})

const handleSubmit = async () => {
  if (!ruleFormRef.value) return
  
  await ruleFormRef.value.validate(async (valid) => {
    if (valid) {
      try {
        // TODO: 调用修改端口的API
        const { data } = await Api.updatePort({port:ruleForm.port})
        ElMessage.success('修改成功')
        dialogVisible.value = false
        emit('update:modelValue', false)
      } catch (error) {
        ElMessage.error('修改失败')
      }
    }
  })
}

const handleClose = () => {
  drawer.value = false
  emit('update:modelValue', false)
}

// 监听 modelValue 变化
watch(() => props.modelValue, (val) => {
  drawer.value = val
  dialogVisible.value = val
})

// 监听变化
watch(() => drawer.value, (val) => {
  emit('update:modelValue', val)
})

watch(() => dialogVisible.value, (val) => {
  emit('update:modelValue', val)
})
</script>

<template>
  <el-drawer
    v-model="dialogVisible"
    width="40%"
    >
    <template #header>
      <div class="drawer-header" style="padding: 20px;font-size: 16px;">
        <span class="title">修改面板端口</span>
      </div>
    </template>
  
    <el-form
      ref="ruleFormRef"
      :model="ruleForm"
      :rules="rules"
      style="padding: 0 20px;"
      label-width="120px"
    >
      <el-form-item label="面板端口" prop="port" label-position="top" required>
        <el-input v-model="ruleForm.port" placeholder="请输入端口号" style="margin-top: 8px;" />
      </el-form-item>
    </el-form>
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit">确定</el-button>
      </span>
    </template>
  </el-drawer>
</template>

<style scoped lang="less">
.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}
</style>
