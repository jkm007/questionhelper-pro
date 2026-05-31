<template>
  <view class="exam-card" @tap="$emit('tap')">
    <view class="card-header">
      <text class="card-title">{{ exam.title }}</text>
      <view class="status-badge" :class="statusClass">
        <text class="status-text">{{ statusText }}</text>
      </view>
    </view>
    <view class="card-info">
      <view class="info-row">
        <text class="info-icon">📅</text>
        <text class="info-text">{{ formatDate(exam.startTime) }} - {{ formatDate(exam.endTime) }}</text>
      </view>
      <view class="info-row">
        <text class="info-icon">⏱</text>
        <text class="info-text">{{ exam.duration }}分钟</text>
      </view>
      <view v-if="exam.participantCount !== undefined" class="info-row">
        <text class="info-icon">👥</text>
        <text class="info-text">{{ exam.participantCount }}人参加</text>
      </view>
    </view>
    <view v-if="exam.myScore !== undefined" class="card-score">
      <text class="score-label">得分</text>
      <text class="score-value">{{ exam.myScore }}/{{ exam.totalScore }}</text>
    </view>
  </view>
</template>

<script setup lang="ts">
import { computed } from 'vue'

interface Exam {
  id: number | string
  title: string
  startTime: string
  endTime: string
  duration: number
  participantCount?: number
  status?: string
  totalScore?: number
  myScore?: number
}

const props = defineProps<{
  exam: Exam
}>()

defineEmits(['tap'])

const statusClass = computed(() => {
  const map: Record<string, string> = {
    upcoming: 'badge-upcoming',
    in_progress: 'badge-active',
    ongoing: 'badge-active',
    completed: 'badge-completed',
  }
  return map[props.exam.status || ''] || 'badge-upcoming'
})

const statusText = computed(() => {
  const map: Record<string, string> = {
    upcoming: '未开始',
    in_progress: '进行中',
    ongoing: '进行中',
    completed: '已结束',
  }
  return map[props.exam.status || ''] || '未开始'
})

function formatDate(dateStr: string): string {
  if (!dateStr) return ''
  const d = new Date(dateStr)
  return `${d.getMonth() + 1}/${d.getDate()} ${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`
}
</script>

<style lang="scss" scoped>
.exam-card {
  background: #ffffff;
  border-radius: 16rpx;
  padding: 24rpx;
  box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.06);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 16rpx;
}

.card-title {
  font-size: 30rpx;
  font-weight: 600;
  color: #303133;
  flex: 1;
  margin-right: 16rpx;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.status-badge {
  padding: 4rpx 16rpx;
  border-radius: 20rpx;
  flex-shrink: 0;
}

.badge-upcoming {
  background: #e3f2fd;
  .status-text { color: #1976d2; }
}

.badge-active {
  background: #e8f5e9;
  .status-text { color: #388e3c; }
}

.badge-completed {
  background: #f5f5f5;
  .status-text { color: #909399; }
}

.status-text {
  font-size: 22rpx;
}

.card-info {
  display: flex;
  flex-direction: column;
  gap: 8rpx;
}

.info-row {
  display: flex;
  align-items: center;
}

.info-icon {
  font-size: 24rpx;
  margin-right: 8rpx;
}

.info-text {
  font-size: 24rpx;
  color: #909399;
}

.card-score {
  margin-top: 16rpx;
  padding-top: 16rpx;
  border-top: 1rpx solid #f0f0f0;
  display: flex;
  align-items: center;
}

.score-label {
  font-size: 24rpx;
  color: #909399;
  margin-right: 12rpx;
}

.score-value {
  font-size: 28rpx;
  color: #409eff;
  font-weight: 600;
}
</style>
