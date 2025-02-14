<script setup lang="ts">
import CardTabs from '@/components/card-tabs.vue'
import { computed, reactive } from 'vue'
import { Api } from '@/api/Api'
import { ElMessage, FormInstance } from 'element-plus'
import { FormItem } from '@/components/custom-form.vue'
import System from '@/utils/System'
import { ColumnItem } from '@/components/custom-table.vue'

export interface ConfProps {
  conf: typeof conf
}

const conf = reactive({
  themeColor: {
    light: ['#F7911C'],
    dark: ['#EAB170']
  },
  showTips: true,
  tabs: {
    activeIndex: 0,
    list: [
      {
        name: 'MySQL',
        index: 0,
        value: 'mysql'
      },
      // {
      //   name: 'SQLServer',
      //   index: 1
      // },
      // {
      //   name: 'MongoDB',
      //   index: 2
      // },
      {
        name: 'Redis',
        index: 3,
        value: 'redis'
      },
      // {
      //   name: 'PgSQL',
      //   index: 4
      // }
    ],
    clickActive: (item: any) => {
      if (conf.tabs.activeIndex === item.index) return
      conf.list.params = {
        page: 1,
        pageSize: 10,
        type: 'mysql'
      }
      conf.tabs.activeIndex = item.index
      conf.list.params.type = item.value
      conf.form.data.value.type = item.value
      conf.list.loading = false
      conf.list.total = 0
      conf.list.data = []
      System.router.push(`/database/${item.value}`)
    }
  },
  list: {
    loading: true,
    data: [],
    total: 0,
    params: {
      page: 1,
      pageSize: 10,
      type: 'mysql'
    } as any,
    columns: computed<ColumnItem[]>(() => {
      switch (conf.list.params.type) {
        case 'mysql':
          return [
            { prop: 'name', label: '数据库名' },
            { prop: 'user', label: '用户名' },
            { prop: 'password', label: '密码' },
            { prop: 'capacity', label: '容量', placeholder: '未配置' },
            { prop: 'p_addr', label: '数据库位置' },
            { prop: 'action', label: '操作' }
          ]
        case 'redis':
          return [
            { prop: 'key', label: '键' },
            { prop: 'type', label: '类型' },
            { prop: 'length', label: '长度' },
            { prop: 'expiration', label: '有效期' }
          ]
        default:
          return []
      }
    }),
    getData: async () => {
      conf.list.loading = true
      const api = conf.list.params.type === 'redis' ? Api.getRedisList : Api.getDatabaseList
      const { data: res } = await api(conf.list.params)
      if (conf.list.params.type !== 'redis')
        res.data.forEach((_: any, index: number) => (conf.password.show[index] = false))
      conf.list.loading = false
      conf.list.total = res.total
      conf.list.data = res[conf.list.params.type === 'redis' ? 'keys' : 'data']
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
    open: (type: 'add') => {
      conf.drawer.type = type
      conf.drawer.show = true
      conf.drawer.title = '添加数据库'
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
        await Api.addDatabaseLib(conf.form.data.value)
        conf.drawer.loading = false
        ElMessage.success('添加成功！')
        conf.list.getData()
        conf.drawer.show = false
      })
    }
  },
  form: {
    instance: null as FormInstance | null,
    data: {
      value: {
        encoding: 'utf8',
        type: 'mysql'
      } as any,
      items: computed<FormItem[]>(() => {
        switch (conf.list.params.type) {
          case 'mysql':
            return [
              {
                label: '数据库名',
                prop: 'name',
                type: 'custom',
                rules: [{ required: true, message: '请输入数据库名', trigger: 'blur' }]
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
                label: '添加至',
                prop: 'id',
                asyncOptions: async () => {
                  const { data } = await Api.getConnlist(conf.list.params)
                  conf.form.data.value.id = data[0].id
                  return data.map((item: any) => ({
                    label: item.remark ? `${item.remark}(${item.addr})` : item.addr,
                    value: item.id
                  }))
                },
                type: 'select',
                rules: [{ required: true, message: '请选择', trigger: 'change' }]
              }
            ]
          default:
            return []
        }
      })
    }
  }
})

const routeName = (System.getRouterPath() as string).match(/(?<=\/database\/)\w*/)?.[0]
conf.tabs.activeIndex = conf.tabs.list.find((item) => item.value === routeName)!.index
conf.list.params.type = routeName
</script>

<template>
  <div class="database-container">
    <card-tabs :list="conf.tabs.list" :active-index="conf.tabs.activeIndex" :click-active="conf.tabs.clickActive" />
    <router-view :conf="conf" />

    <custom-drawer
      :visible="conf.drawer.show"
      :title="conf.drawer.title"
      :loading="conf.drawer.loading"
      :on-close="conf.drawer.onClose"
      :on-confirm="conf.drawer.onConfirm"
    >
      <template v-if="conf.drawer.type === 'add'">
        <custom-form :data="conf.form.data" :on-init="(el) => (conf.form.instance = el)">
          <template #name="{ row }">
            <el-input v-model="conf.form.data.value.name" :placeholder="row.placeholder">
              <template #append>
                <el-select v-model="conf.form.data.value.encoding" style="width: 120px">
                  <el-option label="utf8" value="utf8" />
                  <el-option label="utf8mb4" value="utf8mb4" />
                  <el-option label="gbk" value="gbk" />
                  <el-option label="big5" value="big5" />
                </el-select>
              </template>
            </el-input>
          </template>
        </custom-form>
      </template>
    </custom-drawer>
  </div>
</template>

<style scoped lang="less">
@import '@/styles/element.less';
@import '@/styles/common.less';

.database-container {
  .tip {
    width: 100%;
    height: 58px;
    background: rgb(var(--bg-card-color));
    margin-top: 20px;
    padding: 18px 28px;
    display: flex;
    align-items: center;
    justify-content: space-between;

    span {
      font-size: 16px;
      color: var(--font-color-gray-light);
    }
  }
}
</style>
