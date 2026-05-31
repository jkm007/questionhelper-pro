<template>
  <view class="timer" :class="{ 'timer-urgent': isUrgent }">
    <text class="timer-text">{{ formatted }}</text>
  </view>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  seconds: number
  urgentThreshold?: number
}>()

const isUrgent = computed(() => props.seconds <= (props.urgentThreshold || 300))

const formatted = computed(() => {
  const hrs = Math.floor(props.seconds / 3600)
  const mins = Math.floor((props.seconds % 3600) / 60)
  const secs = props.seconds % 60
  if (hrs > 0) {
    return `${pad(hrs)}:${pad(mins)}:${pad(secs)}`
  }
  return `${pad(mins)}:${pad(secs)}`
})

function pad(n: number): string {
  return String(n).padStart(2, '0')
}
</script>

<style lang="scss" scoped>
.timer {
  display: inline-flex;
  align-items: center;
  padding: 4rpx 16rpx;
  border-radius: 8rpx;
  background: #f0f2f5;
}

.timer-text {
  font-size: 28rpx;
  font-weight: 600;
  color: #409eff;
  font-variant-numeric: tabular-nums;
}

.timer-urgent {
  background: #fef0f0;

  .timer-text {
    color: #f56c6c;
  }
}
</style>
