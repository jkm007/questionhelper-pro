<template>
  <view class="score-chart">
    <view v-if="title" class="chart-header">
      <text class="chart-title">{{ title }}</text>
    </view>
    <view class="chart-body">
      <view class="score-display">
        <text class="score-value" :class="scoreClass">{{ score }}</text>
        <text v-if="total" class="score-total">/{{ total }}</text>
      </view>
      <view v-if="showBar" class="score-bar">
        <view class="bar-track">
          <view class="bar-fill" :style="{ width: percent + '%' }" :class="barClass"></view>
        </view>
        <text class="bar-percent">{{ percent }}%</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  score: number
  total?: number
  title?: string
  showBar?: boolean
}>()

const percent = computed(() => {
  if (!props.total || props.total === 0) return 0
  return Math.round((props.score / props.total) * 100)
})

const scoreClass = computed(() => {
  if (percent.value >= 80) return 'score-high'
  if (percent.value >= 60) return 'score-mid'
  return 'score-low'
})

const barClass = computed(() => {
  if (percent.value >= 80) return 'bar-high'
  if (percent.value >= 60) return 'bar-mid'
  return 'bar-low'
})
</script>

<style lang="scss" scoped>
.score-chart {
  background: #ffffff;
  border-radius: 16rpx;
  padding: 24rpx;
}

.chart-header {
  margin-bottom: 16rpx;
}

.chart-title {
  font-size: 28rpx;
  font-weight: 600;
  color: #303133;
}

.chart-body {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16rpx;
}

.score-display {
  display: flex;
  align-items: baseline;
}

.score-value {
  font-size: 56rpx;
  font-weight: bold;

  &.score-high { color: #67c23a; }
  &.score-mid { color: #e6a23c; }
  &.score-low { color: #f56c6c; }
}

.score-total {
  font-size: 28rpx;
  color: #909399;
  margin-left: 4rpx;
}

.score-bar {
  width: 100%;
  display: flex;
  align-items: center;
  gap: 12rpx;
}

.bar-track {
  flex: 1;
  height: 12rpx;
  background: #f0f2f5;
  border-radius: 6rpx;
  overflow: hidden;
}

.bar-fill {
  height: 100%;
  border-radius: 6rpx;
  transition: width 0.5s ease;

  &.bar-high { background: #67c23a; }
  &.bar-mid { background: #e6a23c; }
  &.bar-low { background: #f56c6c; }
}

.bar-percent {
  font-size: 24rpx;
  color: #909399;
  min-width: 60rpx;
  text-align: right;
}
</style>
