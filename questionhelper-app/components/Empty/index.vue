<template>
  <view class="empty-container">
    <image class="empty-image" src="/static/images/empty/empty.png" mode="aspectFit"></image>
    <text class="empty-text">{{ text }}</text>
    <button v-if="showButton" class="empty-btn" @tap="$emit('action')">
      {{ buttonText }}
    </button>
  </view>
</template>

<script setup lang="ts">
withDefaults(defineProps<{
  text?: string
  showButton?: boolean
  buttonText?: string
}>(), {
  text: '暂无数据',
  showButton: false,
  buttonText: '去看看'
})

defineEmits(['action'])
</script>

<style lang="scss" scoped>
.empty-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 100rpx 0;
}

.empty-image {
  width: 300rpx;
  height: 300rpx;
  margin-bottom: 30rpx;
}

.empty-text {
  font-size: 28rpx;
  color: #909399;
  margin-bottom: 30rpx;
}

.empty-btn {
  background-color: #4A90D9;
  color: #ffffff;
  font-size: 28rpx;
  padding: 16rpx 60rpx;
  border-radius: 40rpx;
  border: none;
}
</style>
