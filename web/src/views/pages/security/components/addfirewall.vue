<template>
    <el-dialog v-model="drawer2" :direction="direction" :show-close="true" @close="cancelClick" :title="type?'添加防火墙':'编辑防火墙'"   >
      
      <template #default>
        <el-form
    ref="ruleFormRef"
    :model="ruleForm"
    :rules="rules"
    class="demo-ruleForm"
    status-icon
  >
  <el-form-item label="规则方向" prop="direction" >
      <el-select v-model="ruleForm.direction" placeholder="请选择规则方向">
        <el-option label="入" value="in" />
        <el-option label="出" value="out" />
      </el-select>
    </el-form-item>

    <el-form-item label="协议类型" prop="protocol">
      <el-select v-model="ruleForm.protocol"  placeholder="请选择协议类型">
        <el-option label="tcp" value="tcp" />
        <el-option label="udp" value="udp" />
        <el-option label="icmp" value="icmp" />
      </el-select>
    </el-form-item>
<!-- 
    <el-form-item label="来源" prop="ips"  placeholder="请选择来源"> 
        <el-option label="tcp" value="tcp" />
        <el-option label="udp" value="udp" />
        <el-option label="部分ip"   v-model:value="ruleForm.ips"  v-if=""/>
    </el-form-item>  -->

    <el-form-item label="端口" prop="ports">
      <el-input v-model="ruleForm.ports" placeholder="支持多个端口" />
    </el-form-item>

    <el-form-item label="状态" prop="state">
      <el-switch v-model="ruleForm.state" />
    </el-form-item>
    <el-form-item label="备注" prop="ports">
      <el-input v-model="ruleForm.ports" placeholder="可为空" />
    </el-form-item> 
    <el-form-item label="支持添加多个端口,如80,88"  :required="true">
    </el-form-item>
    <el-form-item label="支持添加多个端口范围,如80,88,90-99,100-101"  :required="true">
    </el-form-item>
    
    
  </el-form>
      </template>
      <template #footer>
        <div style="flex: auto">
          <el-button @click="cancelClick">取消</el-button>
          <el-button type="primary" @click="confirmClick">确定</el-button>
        </div>
      </template>
    </el-dialog>
  </template>
  
  <script lang="ts" setup>
    import { reactive, ref } from 'vue'
import type { ComponentSize, FormInstance, FormRules } from 'element-plus'
  import { ElMessageBox } from 'element-plus'
  import type { DrawerProps } from 'element-plus'

  const drawer2 = ref(false)
  const direction = ref<DrawerProps['direction']>('rtl')
  const radio1 = ref('Option 1')
  const ruleForm = ref({
    id:'',
    direction:'in',// 规则方向 "in" 或 "out"
    protocol:'tcp', //协议类型 "tcp", "udp", "icmp"
    ips:'',  // 支持多个 IP 地址或子网
    ports:'',//支持多个端口
    state:'',//状态
    remark:'',//备注
  
  })
  const type = ref(true)
  const ruleFormRef = ref<FormInstance>()
  const rules = ref<FormRules>({})
  
  const handleClose = (done: () => void) => {
    ElMessageBox.confirm('Are you sure you want to close this?')
      .then(() => {
        done()
      })
      .catch(() => {
        // catch error
      })
  }
  const options = ref([])
  const locationOptions = ref([])
  const submitForm = (formEl: FormInstance | undefined) => {
    if (!formEl) return
    formEl.validate((valid) => {
      if (valid) {
        console.log('submit!')
      } else {
        console.log('error submit!')
      }
    })
  }
  const resetForm = (formEl: FormInstance | undefined) => {
    if (!formEl) return
    formEl.resetFields()
  }
  const cancelClick = () => {
    console.log('cancelClick')
    
    drawer2.value = false
    // emit('close', false)
  }

  function confirmClick() {
    drawer2.value = false
    console.log('confirmClick')
    console.log(ruleForm.value)


    // ElMessageBox.confirm(`Are you confirm to chose ${radio1.value} ?`)
    //   .then(() => {
    //     drawer2.value = false
    //   })
    //   .catch(() => {
    //     // catch error
    //   })
  }
  

</script>
  