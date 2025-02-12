<script setup lang="ts">
import panelSetting from './panel-setting.vue'
import securitySetting from './security-setting.vue'
</script>

<template>
  <div class="all-container">
    <panel-setting :isCard="false" />
    <security-setting :isCard="false" />
  </div>
</template>

<style scoped lang="less">
.all-container {
  background: rgb(var(--bg-card-color));
  border-radius: 5px;
  padding: 26px 30px;
}
</style>
