<script setup lang="ts">
import CustomDialog from '@/components/custom-dialog.vue'

interface Props {
  show: boolean
  memo: {
    id: number
    content: string
  }
  close: () => void
  update: () => void
}

defineProps<Props>()
</script>

<template>
  <custom-dialog :show="show" title="备忘录" width="624px" :show-close="false" @update:show="close">
    <el-input v-model="memo.content" type="textarea" placeholder="请输入备忘录内容" />
    <template #footer>
      <el-button type="primary" @click="update">保存</el-button>
      <el-button @click="close">取消</el-button>
    </template>
  </custom-dialog>
</template>

<style scoped lang="less"></style>
