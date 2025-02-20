<template>
  <div class="task-container">
    <card-tabs :list="conf.list" :activeIndex="conf.activeIndex" :clickActive="conf.clickActive" />
    <component :is="conf.list[conf.activeIndex].component" />
  </div>
  
</template>

<script setup lang="ts">
import { ref, onBeforeUnmount ,reactive } from 'vue';
import CardTabs from '@/components/card-tabs.vue'
import Terminal from '@/views/pages/terminal/components/terminal.vue'

const conf = reactive({
  activeIndex: 0,
  list: [
    {
      name: '终端',
      index: 0,
      component: Terminal
    },
    // ,{
    //   name: '主机',
    //   index: 1
    // },
    // {
    //   name: '快速命令',
    //   index: 2
    // }
  ],
  clickActive: (item: any) => {
    conf.activeIndex = item.index
  }
})

</script>

<style scoped>
#terminal {
  font-family: 'Courier New', Courier, monospace;
  background-color: #000;
  color: #fff;
  padding: 10px;
  height: 400px;
  overflow-y: auto;
}

.output {
  margin-bottom: 10px;
}

.input {
  display: flex;
  align-items: center;
}

.input span {
  margin-right: 5px;
}

.input input {
  background-color: #000;
  color: #fff;
  border: none;
  outline: none;
  flex: 1;
}

.connection {
  margin-top: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.connection label {
  margin-right: 10px;
}

.connection input {
  margin-right: 10px;
}
</style>