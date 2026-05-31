<template>
  <view class="rich-text-wrap">
    <rich-text :nodes="content" />
  </view>
</template>

<script setup lang="ts">
defineProps<{
  content: string
}>()
</script>

<style lang="scss" scoped>
.rich-text-wrap {
  font-size: 28rpx;
  color: #303133;
  line-height: 1.8;
  word-break: break-all;
}
</style>
