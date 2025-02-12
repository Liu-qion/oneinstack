<script setup lang="ts">
import { Back } from '@element-plus/icons-vue'

interface Props {
  visible: boolean
  title: string
  cancelText?: string
  confirmText?: string
  loading?: boolean
  onClose?: () => void
  onConfirm?: () => void
}

const props = withDefaults(defineProps<Props>(), {
  visible: false,
  title: '',
  loading: false,
  cancelText: '取消',
  confirmText: '确认'
})
</script>

<template>
  <el-drawer :model-value="visible" direction="rtl" size="40%" :show-close="false" @close="onClose">
    <template #header>
      <div class="drawerHeader">
        <div class="back" @click="onClose">
          <el-icon><Back /></el-icon>
          <span>返回</span>
        </div>
        <span class="title">{{ props.title }}</span>
      </div>
    </template>
    <template #default>
      <slot />
    </template>
    <template #footer>
      <div style="flex: auto">
        <el-button
          style="border: 1px solid rgb(var(--border-color-gray)); color: var(--font-color-black)"
          color="transparent"
          @click="onClose"
        >
          {{ cancelText }}
        </el-button>
        <el-button :loading="loading" type="primary" @click="onConfirm">{{ confirmText }}</el-button>
      </div>
    </template>
  </el-drawer>
</template>

<style scoped lang="less">
.drawerHeader {
  background: rgb(var(--category-item-bg-color));
  padding: 20px;
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
