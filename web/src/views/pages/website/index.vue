<script setup lang="ts">
import { computed, reactive, toRaw } from 'vue'
import SearchInput from '@/components/search-input.vue'
import { Refresh, Setting, ArrowDown, CaretBottom } from '@element-plus/icons-vue'
import CardTabs from '@/components/card-tabs.vue'
import CustomTable from '@/components/custom-table.vue'
import { Api } from '@/api/Api'
import type { FormItem } from '@/components/custom-form.vue'
import { FormInstance } from 'element-plus'

const conf = reactive({
  tabs: {
    activeIndex: 0,
    list: [
      {
        name: 'PHP项目',
        index: 0,
        value: 'php'
      },
      // {
      //   name: 'JAVA项目',
      //   index: 1
      // },
      // {
      //   name: 'Node项目',
      //   index: 2
      // },
      // {
      //   name: 'Go项目',
      //   index: 3
      // },
      // {
      //   name: 'Python项目',
      //   index: 4
      // },
      {
        name: '反向代理',
        index: 5,
        value: 'proxy'
      },
      {
        name: 'HTML项目',
        index: 6,
        value: 'static'
      }
      // {
      //   name: '其他项目',
      //   index: 7
      // }
    ],
    clickActive: (item: any) => {
      conf.tabs.activeIndex = item.index
      conf.website.params.type = item.value
      conf.website.getData()
    }
  },
  website: {
    data: [],
    total: 0,
    columns: [
      { prop: 'name', label: '网站名', width: 200 },
      { prop: 'domain', label: '域名', width: 200 },
      { prop: 'root_dir', label: '根目录' },
      { prop: 'remark', label: '备注', width: 200 },
      { prop: 'action', label: '操作' }
    ],
    params: {
      type: 'php',
      page: 1,
      pageSize: 10
    } as any,
    loading: true,
    getData: async () => {
      conf.website.loading = true
      const { data: res } = await Api.getWebsiteList(conf.website.params)
      conf.website.loading = false
      conf.website.total = res.total
      conf.website.data = res.data
    },
    handleAdd: () => {
      conf.drawer.open('add')
      conf.form.data.value.type = conf.website.params.type
    }
  },
  drawer: {
    show: false,
    title: '创建网站',
    type: 'add',
    loading: false,
    open: (type: 'add' | 'edit', row?: any) => {
      conf.drawer.title = '创建网站'
      conf.drawer.type = type
      if (type === 'edit') {
        conf.drawer.title = '设置网站'
        const cloneRow = structuredClone(toRaw(row))
        const domain = cloneRow.domain?.split(',')
        conf.form.data.value = cloneRow
        conf.form.data.value.hostDomain = domain[0].trim()
        domain.shift()
        conf.form.data.value.otherDomain = domain.join('\n')
      }
      conf.drawer.show = true
    },
    onConfirm: () => {
      conf.form.instance?.validate(async (valid) => {
        if (!valid) return
      		let otherDomain =''
      		if(conf.form.data.value.otherDomain){
      			otherDomain = conf.form.data.value.otherDomain?.split('\n')
      		}else{
      			otherDomain=''
      		}
      		
        
        conf.form.data.value.domain = otherDomain!=''
          ? `${conf.form.data.value.hostDomain.trim()},${otherDomain}`
          : conf.form.data.value.hostDomain
        console.log(conf.form.data.value)
        conf.drawer.loading = true
        const api = conf.drawer.type === 'add' ? Api.addWebsite : Api.updateWebsite
        await api(conf.form.data.value)
        conf.drawer.loading = false
        conf.drawer.show = false
        conf.website.getData()
      })
    },
    onClose: () => {
      conf.form.data.value = {}
      conf.form.instance?.resetFields()
      conf.form.instance?.clearValidate()
      conf.drawer.show = false
    }
  },
  form: {
    instance: null as FormInstance | null,
    data: {
      value: {} as any,
      items: computed<FormItem[]>(() => {
        switch (conf.website.params.type) {
          case 'php':
          case 'static':
            return [
              {
                label: '主域名',
                type: 'input',
                placeholder: '支持域名:端口',
                prop: 'hostDomain',
                rules: [
                  { required: true, message: '请输入主域名', trigger: 'blur' },
                  {
                    pattern: /^([0-9a-zA-Z-]{1,}\.)+([a-zA-Z]{2,})(\:\d{2,})?$/,
                    message: '域名格式错误',
                    trigger: 'blur'
                  }
                ]
              },
              {
                label: '其他域名',
                type: 'textarea',
                placeholder: '一行一个域名，支持*和IP地址，支持"域名:端口"',
                prop: 'otherDomain'
              },
              {
                label: '根目录',
                type: 'input',
                prop: 'root_dir',
                rules: [
                  { required: true, message: '请选择根目录', trigger: 'blur' },
                  { pattern: /^\/(?:[^/]+\/)*[^/]+$/, message: '路径格式错误' }
                ]
              },
              {
                label: '备注',
                type: 'textarea',
                prop: 'remark'
              }
            ]
          case 'proxy':
            conf.form.data.value.pact = 'http'
            conf.form.data.value.tar_url = '$http_host'
            return [
              {
                label: '主域名',
                type: 'input',
                placeholder: '支持域名:端口',
                prop: 'hostDomain',
                rules: [
                  { required: true, message: '请输入主域名', trigger: 'blur' },
                  { pattern: /^([0-9a-zA-Z-]{1,}\.)+([a-zA-Z]{2,})$/, message: '域名格式错误', trigger: 'blur' }
                ]
              },
              {
                label: '其他域名',
                type: 'textarea',
                placeholder: '一行一个域名，支持*和IP地址，支持"域名:端口"',
                prop: 'otherDomain'
              },
              {
                label: '代理地址',
                type: 'custom',
                placeholder: '例：127.0.0.1:8080',
                prop: 'send_url',
                rules: [{ required: true, message: '请输入代理地址', trigger: 'blur' }]
              },
              {
                label: '备注',
                type: 'textarea',
                prop: 'remark'
              }
            ]
          default:
            return []
        }
      })
    }
  },
  dialog: {
    show: false,
    title: '网站删除确认',
    type: 'delete',
    row: {} as any,
    open: (type: 'delete', row?: any) => {
      conf.dialog.type = type
      conf.dialog.row = row
      switch (type) {
        case 'delete':
          conf.dialog.title = '网站删除确认'
          break
      }
      conf.dialog.show = true
    },
    close: () => {
      conf.dialog.show = false
    },
    confirm: async () => {
      await Api.delWebsite({ id: conf.dialog.row.id })
      conf.website.getData()
      conf.dialog.show = false
    }
  }
})

conf.website.getData()
</script>

<template>
  <div class="website-container">
    <card-tabs
      :list="conf.tabs.list"
      :active-index="conf.tabs.activeIndex"
      :click-active="conf.tabs.clickActive"
    />
    <div class="tool-bar">
      <el-space class="btn-group" :size="14" style="width: 100%;" >
        <el-button type="primary" @click="conf.website.handleAdd" >添加站点</el-button>
		
        <!-- <el-dropdown>
            <el-button type="primary">
              <span class="el-dropdown-link">
                高级设置
                <el-icon class="el-icon--right"><arrow-down /></el-icon>
              </span>
            </el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item>修改默认页面</el-dropdown-item>
                <el-dropdown-item>默认站点</el-dropdown-item>
                <el-dropdown-item>PHP命令行版本</el-dropdown-item>
                <el-dropdown-item>HTTPS防窜站</el-dropdown-item>
                <el-dropdown-item>TLS设置</el-dropdown-item>
                <el-dropdown-item>全局设置</el-dropdown-item>
                <el-dropdown-item>关联数据库</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
          <el-button type="primary">
            <span>漏洞扫描（0）</span>
          </el-button>
          <el-button type="primary">
            <span style="font-size: 14px; margin-right: 8px">nignx</span>
            <el-icon><CaretBottom /></el-icon>
          </el-button>
          <el-dropdown>
            <el-button type="primary">
              <span class="el-dropdown-link">
                全部分类
                <el-icon class="el-icon--right"><arrow-down /></el-icon>
              </span>
            </el-button>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item>修改默认页面</el-dropdown-item>
                <el-dropdown-item>默认站点</el-dropdown-item>
                <el-dropdown-item>PHP命令行版本</el-dropdown-item>
                <el-dropdown-item>HTTPS防窜站</el-dropdown-item>
                <el-dropdown-item>TLS设置</el-dropdown-item>
                <el-dropdown-item>全局设置</el-dropdown-item>
                <el-dropdown-item>关联数据库</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown> -->
		
      </el-space>
      <div class="demo-form-inline">
       <el-space class="btn-group" :size="14" >
       	<search-input
       	  v-model="conf.website.params.name"
       	  placeholder="请输入域名"
       	  style="margin-right: 18px"
       	  @search="conf.website.getData()"
       	/>
       	<el-button :icon="Refresh" type="primary" @click="conf.website.getData()" />
       	<!-- <el-button :icon="Setting" type="primary" /> -->
       </el-space>
      </div>
    </div>
    <div class="box2">
      <custom-table
        v-model:page="conf.website.params.page"
        :loading="conf.website.loading"
        :data="conf.website.data"
        :columns="conf.website.columns"
        :auto-pagination="false"
        :total="conf.website.total"
        :page-size="conf.website.params.pageSize"
        @update:page="conf.website.getData"
		
      >
        <template #action="{ row }">
          <el-button type="primary" link @click="conf.drawer.open('edit', row)">设置</el-button>
          <el-button type="danger" link @click="conf.dialog.open('delete', row)">删除</el-button>
        </template>
      </custom-table>
    </div>

    <custom-drawer
      :visible="conf.drawer.show"
      :title="conf.drawer.title"
      :loading="conf.drawer.loading"
      :on-close="conf.drawer.onClose"
      :on-confirm="conf.drawer.onConfirm"
    >
      <custom-form v-if="conf.drawer.show" :data="conf.form.data" :on-init="(el) => (conf.form.instance = el)">
        <template #send_url="{ row }">
          <el-input v-model="conf.form.data.value.send_url" :placeholder="row.placeholder">
            <template #prepend>
              <el-select v-model="conf.form.data.value.pact" style="width: 80px">
                <el-option label="http" value="http://" />
                <el-option label="https" value="https://" />
              </el-select>
            </template>
          </el-input>
        </template>
      </custom-form>
    </custom-drawer>

    <custom-dialog v-model="conf.dialog.show" :title="conf.dialog.title">
      <template v-if="conf.dialog.type === 'delete'">
        <el-alert title="确定删除所选网站？" type="warning" show-icon :closable="false" />
      </template>
      <template #footer>
        <el-button @click="conf.dialog.close">取消</el-button>
        <el-button type="primary" @click="conf.dialog.confirm">确认</el-button>
      </template>
    </custom-dialog>
  </div>
</template>

<style scoped lang="less">
	

</style>
