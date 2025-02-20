<script setup lang="ts">
import { reactive ,ref} from 'vue'
import { ChildEmits, ChildProps } from '../index.vue'
import CustomDrawer from '@/components/custom-drawer.vue'
import CustomForm, { type FormItem, type Props as FormProps } from '@/components/custom-form.vue'
import { FormInstance } from 'element-plus'
import { Api } from '@/api/Api'
import { Scope } from 'tools-vue3'
import { ArrowDown } from '@element-plus/icons-vue'

const emit = defineEmits<ChildEmits>()

withDefaults(defineProps<ChildProps>(), {
  list: () => []
})
const formRef = ref<FormInstance | null>(null)
const drawer = reactive({
  show: false,
  title: '安装',
  onClose: () => {
    formRef.value?.clearValidate()
    formRef.value?.resetFields()
    drawer.show = false
  },
  onConfirm: () => {
    formRef.value?.validate((valid:any) => {
      if (!valid) return
      handleInstall()
      drawer.show = false
    })
  }
})

// let formRef: FormInstance
const installForm = reactive<FormProps['data']>({
  value: {
    key: '',
    version: ''
  },
  items: []
})

const handleSelectVersion = (v: string, item: any) => {
  installForm.value.key = item.key
  installForm.value.version = v
}

const handleClickInstall = (item: any) => {
  console.log('installForm.value.key', installForm.value.key,installForm.value.key!='',item.versions[0],item)
  
  if (installForm.value.version!='') {
    if(item && item.versions && item.versions.includes(installForm.value.version)) {
      handleSelectVersion( installForm.value.version, item)
    }else{
      handleSelectVersion(item.versions[0], item)
    }
  }else{
    handleSelectVersion(item.versions[0], item)
  }
  
 
  if (!item.params) return handleInstall()
  const config = JSON.parse(item.params)
  installForm.items = (config as []).map<FormItem>((item: any) => ({
    label: item.name,
    type: 'input',
    prop: item.key,
    rules: [{ required: item.required === 'true', message: `请输入${item.name}`, trigger: 'blur' }]
  }))
  drawer.show = true
}

const handleInstall = async () => {
  const { data: res } = await Api.installSoft(installForm.value)
  handleCheckInstallLog(res.installName)
}

const dialog = reactive({
  show: false,
  content: '',
  onClose: () => {
    timer && timer.clear()
    emit('refresh')
  }
})

const timer = Scope.Timer()
const handleCheckInstallLog = async (value: string) => {
  dialog.show = true
  timer.on(
    async () => {
      const { data: res2 } = await Api.getInstallLog({ fn: value })
      dialog.content = res2.logs
    },
    2000,
    true
  )
}
</script>

<template>
  <div >
    <div class="title">应用</div>
    <div class="list">
      <template v-if="list.length">
        <div v-for="item in list" class="item">
          <div style="padding: 28px 26px">
            <div class="sundry">
              <div class="icon">
                <img :src="item.icon" alt="" />
              </div>
              <div class="content">
                <div>
                  <span class="menuTitle">{{ item.name }}</span>
                  <span v-if="item.tags" class="remark">（{{ item.tags }}）</span>
                  <span
                    v-if="item.status !== 0"
                    class="status ml-2"
                    :class="{ error: item.status === 3, success: item.status === 2, installing: item.status === 1 }"
                    @click="handleCheckInstallLog(item.log)"
                  >
                    {{ item.status === 1 ? '安装中' : item.status === 2 ? '已安装' : '安装失败' }}
                  </span>
                </div>
                <div class="tip">{{ item.describe }}</div>
              </div>
            </div>
            <div class="xian" />
            <div class="below">
              <div>
                <el-dropdown :disabled="!(item.versions.length > 1)">
                  <span class="flex items-center" style="color: var(--font-color-gray-light)">
                    <span>安装版本：</span>
                    {{ installForm.value.key === item.key ? installForm.value.version : item.versions[0] }}
                    <el-icon v-if="item.versions.length > 1" class="el-icon--right">
                      <arrow-down />
                    </el-icon>
                  </span>
                  <template #dropdown>
                    <el-dropdown-menu>
                      <el-dropdown-item v-for="(v, i) in item.versions" :key="i" @click="handleSelectVersion(v, item)">
                        {{ v }}
                      </el-dropdown-item>
                    </el-dropdown-menu>
                  </template>
                </el-dropdown>
              </div>
              <div class="btn" :class="{ installed: item.installed }" @click="handleClickInstall(item)">
                {{ item.installed ? '已安装' : '安装' }}
              </div>
            </div>
          </div>
        </div>
      </template>
      <div v-else class="no-data">
        <img src="/static/images/empty.webp" alt="" />
        <span>暂无应用</span>
      </div>
    </div>
    <custom-drawer
      :visible="drawer.show"
      :title="drawer.title"
      :on-close="drawer.onClose"
      :on-confirm="drawer.onConfirm"
    >
      <custom-form :data="installForm" :on-init="(el) => (formRef = el)" />
    </custom-drawer>
    <custom-dialog v-model:show="dialog.show" title="安装日志" :on-close="dialog.onClose">
      <el-input v-model="dialog.content" type="textarea" readonly style="min-height: 200px" />
    </custom-dialog>
  </div>
</template>

<style scoped lang="less">
.title {
  font-weight: 500;
  font-size: 18px;
  color: var(--font-color-black);
}

.list {
  width: 100%;
  display: flex;
  flex-direction: row;
  flex-wrap: wrap;
  margin-top: 20px;

  .item {
    width: 32%;
    height: 220px;
    background: rgb(var(--bg-hover-color));
    border-radius: 8px;
    margin-left: 2%;
    margin-bottom: 22px;
    border: 2px solid transparent;

    &:hover {
      border-color: var(--el-color-primary);
      cursor: pointer;

      .menuTitle {
        color: var(--el-color-primary);
      }
    }

    .menuTitle {
      font-size: 22px;
      color: var(--font-color-black);
    }

    .remark {
      color: var(--el-color-primary);
      font-size: 14px;
    }

    .sundry {
      margin-bottom: 20px;
      display: flex;
      flex-direction: row;
      align-items: center;

      .icon {
        width: 86px;
        height: 86px;
        background: #ffffff;
        border-radius: 8px;
        overflow: hidden;

        img {
          width: 100%;
          height: 100%;
        }
      }

      .content {
        margin-left: 28px;

        .tip {
          color: var(--font-color-gray-light);
          font-size: 16px;
          margin-top: 18px;
        }

        .status {
          display: inline-block;
          padding: 2px 8px;
          font-size: 12px;
          border-radius: 20px;
          border: 1px solid;

          &.error {
            border-color: rgb(var(--error-color));
            color: rgb(var(--error-color));
          }

          &.success {
            border-color: rgb(var(--success-color));
            color: rgb(var(--success-color));
          }

          &.installing {
            border-color: rgb(var(--primary-color));
            color: rgb(var(--primary-color));
          }
        }
      }
    }

    .xian {
      border-bottom: 1px solid rgba(186, 186, 186, 0.32);
    }

    .below {
      display: flex;
      flex-direction: row;
      justify-content: space-between;
      margin-top: 20px;
      color: var(--font-color-gray-light);

      .btn {
        width: 69px;
        height: 34px;
        background: transparent;
        border: 1px solid var(--el-color-primary);
        border-radius: 4px;
        display: flex;
        flex-direction: row;
        align-items: center;
        justify-content: center;
        font-size: 14px;
        cursor: pointer;
        color: var(--el-color-primary);

        &:hover,
        &.installed {
          background: var(--el-color-primary);
          color: var(--font-color-white);
        }
      }
    }
  }

  .item:nth-of-type(3n-2) {
    margin-left: 0;
  }

  .no-data {
    width: 100%;
    display: flex;
    flex-flow: column nowrap;
    justify-content: center;
    align-items: center;
    color: var(--font-color-gray-light);
  }
}

:deep(.el-textarea__inner) {
  min-height: 200px !important;
}
</style>
