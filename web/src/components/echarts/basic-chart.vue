<script setup lang="ts">
import { watch, onMounted, onUnmounted, useTemplateRef } from 'vue'
import * as echarts from 'echarts'
import type { EChartsOption, EChartsType } from 'echarts'

interface Props {
  option: EChartsOption
  onInit?: (instance: EChartsType) => void
}

const props = defineProps<Props>()

const chartBox = useTemplateRef('chartBox')
let chartInstance: EChartsType

watch(
  () => props.option,
  (newVal) => {
    newVal && chartInstance?.setOption(newVal)
  }
)

const observer = new ResizeObserver(() => chartInstance.resize())
onMounted(() => {
  chartInstance = echarts.init(chartBox.value)
  props.onInit?.(chartInstance)
  observer.observe(chartBox.value as HTMLElement)
})

onUnmounted(() => {
  chartInstance.dispose()
  observer.disconnect()
})
</script>

<template>
  <div ref="chartBox" class="chart-box" />
</template>

<style scoped lang="less">
.chart-box {
  width: 100%;
  height: 100%;
}
</style>
