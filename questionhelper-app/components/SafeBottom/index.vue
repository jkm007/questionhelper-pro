<template>
  <view class="safe-bottom" :style="{ height: safeHeight + 'px' }"></view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

const safeHeight = ref(0)

onMounted(() => {
  const sysInfo = uni.getSystemInfoSync()
  safeHeight.value = sysInfo.safeAreaInsets?.bottom || 0
})
</script>

<style scoped>
.safe-bottom {
  width: 100%;
}
</style>
