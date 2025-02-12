<script setup lang="ts">
import CustomTable, { type ColumnItem } from '@/components/custom-table.vue'
import { Refresh } from '@element-plus/icons-vue'

const column: ColumnItem[] = [
  {
    label: '标题内容',
    prop: '',
    width: 200
  },
  {
    label: '提交时间',
    prop: ''
  },
  {
    label: '类型',
    prop: ''
  }
]
</script>

<template>
  <div>
    <div class="flex justify-between">
      <div class="flex">
        <div class="btnItem mr-2">全部已读</div>
        <div class="btnItem mr-2">全部删除</div>
      </div>
      <div>
        <el-button :icon="Refresh" />
      </div>
    </div>
    <div class="box2" style="padding: 24px 30px">
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
:deep(.el-select),
:deep(.el-input) {
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
