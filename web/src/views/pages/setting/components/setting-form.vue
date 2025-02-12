<script setup lang="ts">
type FormItemType =
  | 'input'
  | 'select'
  | 'switch'
  | 'time'
  | 'datetime'
  | 'number'
  | 'password'
  | 'file'
  | 'custom'
  | 'none'

interface ActionBtn {
  type?: 'default' | 'primary'
  text: string
  click: () => void
}

export interface FormItem {
  label: string
  prop: string
  value: string | number
  type: FormItemType
  disabled?: boolean
  change?: (value: any) => void
  action?: ActionBtn | ActionBtn[]
  tip?: string
  options?: any[]
}

interface Props {
  data: FormItem[]
}

const props = defineProps<Props>()
</script>

<template>
  <el-form label-width="180px">
    <el-form-item v-for="item in data" :key="item.prop" :label="item.label" label-position="left">
      <div class="mr-2">
        <el-input
          v-if="item.type === 'input' || item.type === 'password'"
          v-model="item.value"
          :type="item.type === 'password' ? 'password' : 'text'"
          clearable
          :disabled="item.disabled"
        />
        <el-switch v-else-if="item.type === 'switch'" v-model="item.value" :disabled="item.disabled" />
        <el-input v-else-if="item.type === 'file'" v-model="item.value" :disabled="item.disabled">
          <template #append>
            <v-s-icon name="folders" />
          </template>
        </el-input>
      </div>
      <!-- 按钮 -->
      <template v-if="item.action">
        <el-button
          v-if="Array.isArray(item.action)"
          v-for="action in item.action"
          :type="action.type || 'default'"
          class="mr-2"
          @click="action.click"
        >
          {{ action.text }}
        </el-button>
        <el-button v-else :type="item.action?.type || 'default'" class="mr-2" @click="item.action.click">
          {{ item.action.text }}
        </el-button>
      </template>
      <!-- 提示 -->
      <template v-if="item.tip">
        <span v-html="item.tip" class="tips-text"></span>
      </template>
    </el-form-item>
  </el-form>
</template>

<style scoped lang="less">
:deep(.el-form) {
  &-item__content {
    flex-wrap: nowrap;
  }
}

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
    background: transparent !important;
    border: 0.1px solid #7d838980 !important;
    box-shadow: none !important;
    border-left: none !important;
  }
}

.tips-text {
  font-size: 14px;
  color: var(--font-color-gray-light);
}
</style>
