<script setup lang="ts">
import { computed, reactive, toRaw } from 'vue'
import { Hide, View, CopyDocument, Back } from '@element-plus/icons-vue'
import { Api } from '@/api/Api'
import { ElMessage, FormInstance } from 'element-plus'
import { FormItem } from '@/components/custom-form.vue'
import System from '@/utils/System'

const conf = reactive({
  themeColor: {
    light: ['#F7911C'],
    dark: ['#EAB170']
  },
  list: {
    loading: true,
    data: [],
    params: {
      type: 'mysql'
    },
    columns: [
      { prop: 'addr', label: '数据库地址' },
      { prop: 'port', label: '端口' },
      { prop: 'root', label: '用户名' },
      { prop: 'password', label: '密码' },
      { prop: 'remark', label: '备注' },
      { prop: 'action', label: '操作' }
    ],
    getData: async () => {
      conf.list.loading = true
      const { data } = await Api.getConnlist(conf.list.params)
      conf.list.loading = false
      conf.list.data = data
    },
    syncData: async (id: number) => {
      await Api.syncDatabaseConn({ id })
      ElMessage.success('同步成功！')
    }
  },
  password: {
    show: {} as { [key: number]: boolean },
    copy: async (value: string) => {
      await navigator.clipboard.writeText(value)
      ElMessage.success('复制成功！')
    }
  },
  drawer: {
    show: false,
    title: '添加数据库',
    type: 'add',
    loading: false,
    open: (type: 'add' | 'edit', row?: any) => {
      conf.drawer.title = '添加数据库'
      if (type === 'edit') {
        conf.drawer.title = '编辑数据库'
        const cloneRow = structuredClone(toRaw(row))
        conf.form.data.value = cloneRow
      }
      conf.drawer.type = type
      conf.drawer.show = true
    },
    onClose: () => {
      conf.form.instance?.clearValidate()
      conf.form.instance?.resetFields()
      conf.drawer.show = false
    },
    onConfirm: () => {
      conf.form.instance?.validate(async (valid) => {
        if (!valid) return
        conf.drawer.loading = true
        const api = conf.drawer.type === 'add' ? Api.addDatabaseConn : Api.updateDatabaseConn
        await api(conf.form.data.value)
        conf.drawer.loading = false
        ElMessage.success(conf.drawer.type === 'add' ? '添加成功！' : '编辑成功！')
        conf.list.getData()
        conf.drawer.show = false
      })
    }
  },
  form: {
    instance: null as FormInstance | null,
    data: {
      value: {
        type: 'mysql'
      } as any,
      items: computed<FormItem[]>(() => {
        switch (conf.list.params.type) {
          case 'mysql':
            return [
              {
                label: '数据库地址',
                prop: 'addr',
                type: 'input',
                rules: [{ required: true, message: '请输入数据库地址', trigger: 'blur' }]
              },
              {
                label: '端口',
                prop: 'port',
                type: 'input',
                rules: [{ required: true, message: '请输入端口', trigger: 'blur' }]
              },
              {
                label: '用户名',
                prop: 'root',
                type: 'input',
                rules: [{ required: true, message: '请输入用户名', trigger: 'blur' }]
              },
              {
                label: '密码',
                prop: 'password',
                type: 'password',
                rules: [{ required: true, message: '请输入密码', trigger: 'blur' }]
              },
              {
                label: '备注',
                prop: 'remark',
                type: 'textarea'
              }
            ]
          default:
            return []
        }
      })
    }
  }
})

conf.list.getData()
</script>

<template>
  <div class="database-container">
    <div class="container">
      <div class="tool-bar">
        <div class="btn-group">
          <el-button type="primary" @click="conf.drawer.open('add')">添加远程服务器</el-button>
        </div>
      </div>
      <div class="box2">
        <div class="drawerHeader">
          <div class="back" @click="System.router.back()">
            <el-icon><Back /></el-icon>
            <span>返回</span>
          </div>
          <span class="title">远程服务器</span>
        </div>
        <custom-table :loading="conf.list.loading" :data="conf.list.data" :columns="conf.list.columns">
          <template #empty>
            <div style="margin-top: 40px">
              <span>
                您的远程服务器列表为空，您可以
                <a
                  class="cursor-pointer"
                  style="color: var(--el-color-primary); text-decoration: underline"
                  @click="conf.drawer.open('add')"
                >
                  添加一个远程服务器
                </a>
              </span>
            </div>
          </template>
          <template #password="{ row, index }">
            <div class="flex items-center" style="gap: 8px">
              <span>
                {{ conf.password.show[index] ? row.password : ''.padEnd(row.password.length, '*') }}
              </span>
              <el-icon
                class="hover-opacity cursor-pointer"
                :size="18"
                @click="conf.password.show[index] = !conf.password.show[index]"
              >
                <View v-if="conf.password.show[index]" />
                <Hide v-else />
              </el-icon>
              <el-icon class="hover-opacity cursor-pointer" @click="conf.password.copy(row.password)">
                <CopyDocument />
              </el-icon>
            </div>
          </template>
          <template #action="{ row }">
            <el-button type="primary" link @click="conf.drawer.open('edit', row)">编辑</el-button>
            <el-button type="primary" link @click="conf.list.syncData(row.id)">同步</el-button>
          </template>
        </custom-table>
      </div>
    </div>

    <custom-drawer
      :visible="conf.drawer.show"
      :title="conf.drawer.title"
      :loading="conf.drawer.loading"
      :on-close="conf.drawer.onClose"
      :on-confirm="conf.drawer.onConfirm"
    >
      <template v-if="conf.drawer.type === 'add' || conf.drawer.type === 'edit'">
        <custom-form :data="conf.form.data" :on-init="(el) => (conf.form.instance = el)"></custom-form>
      </template>
    </custom-drawer>
  </div>
</template>

<style scoped lang="less">
@import '@/styles/element.less';
@import '@/styles/common.less';

.drawerHeader {
  padding: 20px 0;
  display: flex;
  flex-direction: row;
  align-items: center;

  .back {
    margin-right: 20px;
    display: flex;
    align-items: center;
    color: var(--font-color-gray);
    cursor: pointer;
    border-right: 0.4px solid var(--font-color-gray);
    padding-right: 10px;

    span {
      margin-left: 5px;
    }
  }

  .title {
    font-size: 16px;
    font-weight: 600;
    color: var(--font-color-black);
  }
}
</style>
