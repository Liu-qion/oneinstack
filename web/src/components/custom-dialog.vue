<script setup lang="ts">
import { CircleClose } from '@element-plus/icons-vue'

interface Props {
  title?: string
  width?: string
  showClose?: boolean
  onClose?: () => void
}

withDefaults(defineProps<Props>(), {
  title: 'title',
  width: '624px',
  showClose: false
})

const show = defineModel<boolean>('show')
</script>

<template>
  <el-dialog
    v-model="show"
    :title="title"
    :width="width"
    :show-close="showClose"
    class="custom-dialog"
    @close="onClose && onClose()"
  >
    <template #header="{ close, titleId }">
      <div class="custom-dialog__header">
        <h4 :id="titleId">{{ title }}</h4>
        <el-icon size="24" color="#A2A2A2" class="cursor-pointer" @click="close"><CircleClose /></el-icon>
      </div>
    </template>
    <template #default>
      <div class="custom-dialog__body">
        <slot />
      </div>
    </template>
    <template #footer>
      <div class="custom-dialog__footer">
        <slot name="footer" />
      </div>
    </template>
  </el-dialog>
</template>

<style scoped lang="less">
:global(.el-dialog.custom-dialog) {
  --el-dialog-bg-color: rgb(var(--bg-card-color));
  --el-dialog-padding-primary: 0;

  &__header {
    overflow: hidden;
  }
}

.custom-dialog {
  &__header {
    padding: 20px 28px;
    background: rgb(var(--category-item-bg-color));
    display: flex;
    align-items: center;
    justify-content: space-between;
    color: var(--font-color-black);
    font-size: 18px;
  }

  &__body {
    padding: 16px 28px 0;
  }

  &__footer {
    padding: 28px;
  }
}
</style>
