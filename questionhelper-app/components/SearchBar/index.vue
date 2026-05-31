<template>
  <view class="search-bar" @tap="handleTap">
    <view class="search-inner">
      <text class="search-icon">🔍</text>
      <input
        v-if="!readonly"
        class="search-input"
        :value="modelValue"
        :placeholder="placeholder"
        confirm-type="search"
        @input="$emit('update:modelValue', ($event as any).detail.value)"
        @confirm="$emit('search', ($event as any).detail.value)"
      />
      <text v-else class="search-placeholder">{{ placeholder }}</text>
      <text v-if="modelValue" class="clear-btn" @tap.stop="$emit('update:modelValue', '')">✕</text>
    </view>
  </view>
</template>

<script setup lang="ts">
defineProps({
  modelValue: { type: String, default: '' },
  placeholder: { type: String, default: '搜索...' },
  readonly: { type: Boolean, default: false },
})

defineEmits(['update:modelValue', 'search', 'tap'])

function handleTap() {
  // readonly mode emits tap for navigation
}
</script>

<style lang="scss" scoped>
.search-bar {
  padding: 16rpx 24rpx;
}

.search-inner {
  display: flex;
  align-items: center;
  background: #f5f7fa;
  border-radius: 36rpx;
  padding: 0 24rpx;
  height: 72rpx;
}

.search-icon {
  font-size: 28rpx;
  margin-right: 12rpx;
}

.search-input {
  flex: 1;
  font-size: 28rpx;
  height: 72rpx;
}

.search-placeholder {
  flex: 1;
  font-size: 28rpx;
  color: #909399;
}

.clear-btn {
  font-size: 28rpx;
  color: #c0c4cc;
  padding: 8rpx;
}
</style>
