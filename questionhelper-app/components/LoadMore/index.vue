<template>
  <view class="load-more">
    <text v-if="loading" class="load-more-text">加载中...</text>
    <text v-else-if="finished" class="load-more-text finished">{{ finishedText }}</text>
    <text v-else class="load-more-text" @tap="$emit('loadmore')">{{ loadText }}</text>
  </view>
</template>

<script setup lang="ts">
defineProps({
  loading: { type: Boolean, default: false },
  finished: { type: Boolean, default: false },
  finishedText: { type: String, default: '没有更多了' },
  loadText: { type: String, default: '点击加载更多' },
})

defineEmits(['loadmore'])
</script>

<style lang="scss" scoped>
.load-more {
  padding: 24rpx 0;
  text-align: center;
}

.load-more-text {
  font-size: 26rpx;
  color: #909399;

  &.finished {
    color: #c0c4cc;
  }
}
</style>
