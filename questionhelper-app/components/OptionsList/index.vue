<template>
  <view class="options-list">
    <view
      v-for="(option, index) in options"
      :key="index"
      class="option-item"
      :class="{
        'option-selected': isSelected(index),
        'option-correct': showAnswer && option.isCorrect,
        'option-wrong': showAnswer && isSelected(index) && !option.isCorrect,
      }"
      @tap="handleSelect(index)"
    >
      <view class="option-index">
        <text class="option-index-text">{{ labels[index] }}</text>
      </view>
      <text class="option-text">{{ option.content }}</text>
      <view v-if="showAnswer && option.isCorrect" class="option-icon">
        <text class="icon-correct">✓</text>
      </view>
      <view v-if="showAnswer && isSelected(index) && !option.isCorrect" class="option-icon">
        <text class="icon-wrong">✗</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
const props = defineProps<{
  options: { content: string; isCorrect?: boolean }[]
  selected: number[]
  showAnswer?: boolean
  disabled?: boolean
}>()

const emit = defineEmits(['select'])

const labels = ['A', 'B', 'C', 'D', 'E', 'F', 'G', 'H']

function isSelected(index: number): boolean {
  return props.selected.includes(index)
}

function handleSelect(index: number) {
  if (props.disabled || props.showAnswer) return
  emit('select', index)
}
</script>

<style lang="scss" scoped>
.options-list {
  display: flex;
  flex-direction: column;
  gap: 16rpx;
}

.option-item {
  display: flex;
  align-items: center;
  background: #ffffff;
  border-radius: 16rpx;
  padding: 24rpx;
  border: 2rpx solid #e4e7ed;
  transition: all 0.2s;

  &.option-selected {
    border-color: #409eff;
    background: #ecf5ff;
  }

  &.option-correct {
    border-color: #67c23a;
    background: #f0f9eb;
  }

  &.option-wrong {
    border-color: #f56c6c;
    background: #fef0f0;
  }
}

.option-index {
  width: 48rpx;
  height: 48rpx;
  border-radius: 50%;
  background: #f0f2f5;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 20rpx;
  flex-shrink: 0;

  .option-selected & { background: #409eff; }
  .option-correct & { background: #67c23a; }
  .option-wrong & { background: #f56c6c; }
}

.option-index-text {
  font-size: 24rpx;
  color: #606266;
  font-weight: 600;

  .option-selected &,
  .option-correct &,
  .option-wrong & { color: #ffffff; }
}

.option-text {
  flex: 1;
  font-size: 28rpx;
  color: #303133;
  line-height: 1.5;
}

.option-icon {
  margin-left: 16rpx;
}

.icon-correct {
  font-size: 32rpx;
  color: #67c23a;
  font-weight: bold;
}

.icon-wrong {
  font-size: 32rpx;
  color: #f56c6c;
  font-weight: bold;
}
</style>
