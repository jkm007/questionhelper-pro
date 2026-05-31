<template>
  <view class="answer-sheet">
    <view class="sheet-header">
      <text class="sheet-title">答题卡</text>
      <text class="sheet-close" @tap="$emit('close')">✕</text>
    </view>
    <view class="sheet-stats">
      <view class="stat-item">
        <text class="stat-value answered">{{ answeredCount }}</text>
        <text class="stat-label">已答</text>
      </view>
      <view class="stat-item">
        <text class="stat-value unanswered">{{ total - answeredCount }}</text>
        <text class="stat-label">未答</text>
      </view>
      <view class="stat-item">
        <text class="stat-value flagged">{{ flaggedCount }}</text>
        <text class="stat-label">标记</text>
      </view>
    </view>
    <scroll-view scroll-y class="sheet-body">
      <view class="grid">
        <view
          v-for="(item, index) in questions"
          :key="item.id"
          class="grid-item"
          :class="{
            'grid-current': index === currentIndex,
            'grid-answered': answeredSet.has(item.id),
            'grid-flagged': flaggedSet.has(item.id),
          }"
          @tap="$emit('select', index)"
        >
          <text class="grid-text">{{ index + 1 }}</text>
        </view>
      </view>
    </scroll-view>
    <view class="sheet-footer">
      <view class="legend-item">
        <view class="legend-dot current"></view>
        <text class="legend-text">当前</text>
      </view>
      <view class="legend-item">
        <view class="legend-dot answered"></view>
        <text class="legend-text">已答</text>
      </view>
      <view class="legend-item">
        <view class="legend-dot flagged"></view>
        <text class="legend-text">标记</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  questions: { id: number | string }[]
  currentIndex: number
  answeredSet: Set<number | string>
  flaggedSet: Set<number | string>
}>()

defineEmits(['close', 'select'])

const total = computed(() => props.questions.length)
const answeredCount = computed(() => props.answeredSet.size)
const flaggedCount = computed(() => props.flaggedSet.size)
</script>

<style lang="scss" scoped>
.answer-sheet {
  background: #ffffff;
  border-radius: 24rpx 24rpx 0 0;
  padding-bottom: env(safe-area-inset-bottom);
}

.sheet-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 24rpx 30rpx;
  border-bottom: 1rpx solid #f0f0f0;
}

.sheet-title {
  font-size: 32rpx;
  font-weight: 600;
  color: #303133;
}

.sheet-close {
  font-size: 36rpx;
  color: #909399;
  padding: 8rpx;
}

.sheet-stats {
  display: flex;
  justify-content: space-around;
  padding: 20rpx 30rpx;
  border-bottom: 1rpx solid #f0f0f0;
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.stat-value {
  font-size: 36rpx;
  font-weight: bold;

  &.answered { color: #409eff; }
  &.unanswered { color: #909399; }
  &.flagged { color: #e6a23c; }
}

.stat-label {
  font-size: 22rpx;
  color: #909399;
  margin-top: 4rpx;
}

.sheet-body {
  max-height: 500rpx;
  padding: 20rpx;
}

.grid {
  display: flex;
  flex-wrap: wrap;
  gap: 16rpx;
}

.grid-item {
  width: calc(12.5% - 14rpx);
  aspect-ratio: 1;
  border-radius: 12rpx;
  background: #f5f7fa;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 2rpx solid #e4e7ed;

  &.grid-current {
    border-color: #409eff;
    background: #ecf5ff;
  }

  &.grid-answered {
    background: #e8f5e9;
    border-color: #67c23a;
  }

  &.grid-flagged {
    background: #fdf6ec;
    border-color: #e6a23c;
  }
}

.grid-text {
  font-size: 24rpx;
  color: #606266;
  font-weight: 500;
}

.sheet-footer {
  display: flex;
  justify-content: center;
  gap: 40rpx;
  padding: 16rpx 30rpx;
  border-top: 1rpx solid #f0f0f0;
}

.legend-item {
  display: flex;
  align-items: center;
  gap: 8rpx;
}

.legend-dot {
  width: 20rpx;
  height: 20rpx;
  border-radius: 6rpx;

  &.current {
    background: #ecf5ff;
    border: 2rpx solid #409eff;
  }

  &.answered {
    background: #e8f5e9;
    border: 2rpx solid #67c23a;
  }

  &.flagged {
    background: #fdf6ec;
    border: 2rpx solid #e6a23c;
  }
}

.legend-text {
  font-size: 22rpx;
  color: #909399;
}
</style>
