<script setup lang="ts">
import { Api } from '@/api/Api'
import { reactive } from 'vue'
import type Node from 'element-plus/es/components/tree/src/model/node'
import type { TreeOptionProps } from 'element-plus/es/components/tree/src/tree.type'

interface Tree {
  name: string
  leaf?: boolean
  isDir?: boolean
}

interface Emits {
  (e: 'select', path: string): void
}

const props: TreeOptionProps = {
  label: 'name'
}

const emit = defineEmits<Emits>()

const state = reactive({
  path: '/',
  fileList: [],
  loading: false
})

const getFileList = async () => {
  state.loading = true
  const { data: res } = await Api.getFileList({
    path: state.path,
    final: () => (state.loading = false)
  })
  state.fileList = res.files?.filter((item: any) => item.isDir) ?? []
}

const loadNode = async (node: Node, resolve: (data: Tree[]) => void) => {
  let currentNode = node,
    paths = []
  while (currentNode.parent) {
    paths.unshift(currentNode.data.name)
    currentNode = currentNode.parent
  }
  state.path = `/${paths.join('/')}`
  emit('select', state.path)
  await getFileList()
  return resolve(state.fileList)
}
</script>

<template>
  <el-tree v-loading="state.loading" :props="props" :load="loadNode" lazy>
    <template #default="{ data }">
      <div class="flex items-center" style="gap: 20px">
        <v-s-icon name="folder" />
        <span>{{ data.name }}</span>
      </div>
    </template>
  </el-tree>
</template>

<style scoped lang="less"></style>
