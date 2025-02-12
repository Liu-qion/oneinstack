<script setup lang="ts">
interface Props {
  placeholder?: string
}

interface Emits {
  (event: 'search', value: string): void
}

const emit = defineEmits<Emits>()

withDefaults(defineProps<Props>(), {
  placeholder: 'Please Input'
})

const searchValue = defineModel('modelValue', {
  type: String,
  default: ''
})

const handleSearch = () => {
  emit('search', searchValue.value)
}
</script>

<template>
  <el-input
    v-model="searchValue"
    :placeholder="placeholder"
    clearable
    @keydown.enter="handleSearch"
    @clear="handleSearch"
  >
    <template #suffix>
      <v-s-icon name="search" size="16" class="cursor-pointer" @click="handleSearch" />
    </template>
  </el-input>
</template>

<style scoped lang="less">
.el-input {
  --el-input-border-radius: 18px;
  --el-input-width: 262px;
  --el-input-height: 36px;
  --el-input-border-color: rgb(var(--border-input-color));
  --el-input-text-color: var(--font-color-black);
  font-size: 14px;

  :deep(.el-input__wrapper) {
    background-color: rgba(var(--bg-input-color), 0.5);
    padding-inline: 20px;
  }
}
</style>
