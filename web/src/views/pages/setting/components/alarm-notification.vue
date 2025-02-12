<script setup lang="ts">
import CustomTable, { type ColumnItem } from '@/components/custom-table.vue'

const column: ColumnItem[] = [
  {
    label: '标题',
    prop: '',
    width: 200
  },
  {
    label: '状态',
    prop: ''
  },
  {
    label: '告警方式',
    prop: ''
  },
  {
    label: '告警条件',
    prop: ''
  },
  {
    label: '操作',
    prop: ''
  }
]
</script>

<template>
  <div>
    <div class="flex">
      <div class="btnItem mr-2">告警列表</div>
      <div class="btnItem mr-2">告警设置</div>
      <div class="btnItem">告警日志</div>
    </div>
    <div class="box2" style="padding: 24px 30px">
      <el-form label-width="auto">
        <el-form-item label="任务类型">
          <el-select placeholder="网站证书(SSL)到期"></el-select>
        </el-form-item>
        <el-form-item label="网站">
          <el-select placeholder="所有网站"></el-select>
        </el-form-item>
        <el-form-item label="剩余天数">
          <el-input>
            <template #append>
              <span>天</span>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item label="发送次数">
          <el-input>
            <template #append>
              <span>次</span>
            </template>
          </el-input>
        </el-form-item>
        <el-form-item label="告警方式">
          <el-checkbox-group>
            <el-checkbox>邮箱</el-checkbox>
            <el-checkbox>钉钉</el-checkbox>
            <el-checkbox>飞书</el-checkbox>
          </el-checkbox-group>
        </el-form-item>
      </el-form>
      <div class="flex justify-between" style="margin: 36px 0 20px">
        <el-button type="primary">添加任务</el-button>
        <div>
          <el-button>启动</el-button>
          <el-button>停止</el-button>
          <el-button>停止</el-button>
        </div>
      </div>
      <custom-table :data="[]" :columns="column" selection>
        <template #empty>
          <div class="flex column justify-between items-center" style="margin-top: 26px">
            <img src="/static/images/empty.webp" alt="" />
            <span style="color: var(--font-color-gray-light)">暂无任务，点击添加任务开始操作吧</span>
          </div>
        </template>
      </custom-table>
    </div>
  </div>
</template>

<style scoped lang="less">
:deep(.el-select), :deep(.el-input) {
  width: 297px;

  &.is-disabled {
    .el-input__wrapper,
    .el-input-group__append {
      background: rgb(var(--category-item-bg-color)) !important;
    }
  }

  &__wrapper {
    background: transparent !important;
    box-shadow: 0 0 0 0.4px #7d8389 inset !important;
  }

  &-group__append {
    background: rgb(var(--category-item-bg-color)) !important;
    border: 0.1px solid #7d838980 !important;
    box-shadow: none !important;
    border-left: none !important;
  }
}
</style>
