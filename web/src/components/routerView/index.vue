<template>
  <router-view v-slot="{ Component }" v-if="show">
    <keep-alive :include="[...usrouterview.include, ...usrouterview.includeDefault]">
      <component :is="Component" />
    </keep-alive>
  </router-view>
</template>
<script setup lang="ts">
import usrouterview from './usrouterview'

defineProps({
  show: { default: true }
})
</script>
