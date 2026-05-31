<template>
  <view class="exam-history-page">
    <!-- Filter Tabs -->
    <view class="filter-bar">
      <view
        v-for="tab in statusTabs"
        :key="tab.value"
        class="filter-tab"
        :class="{ active: currentStatus === tab.value }"
        @tap="switchStatus(tab.value)"
      >
        <text class="tab-text" :class="{ active: currentStatus === tab.value }">{{ tab.label }}</text>
      </view>
    </view>

    <!-- Exam List -->
    <scroll-view class="exam-list" scroll-y @scrolltolower="loadMore">
      <view v-if="examList.length === 0 && !loading" class="empty-state">
        <image class="empty-icon" src="/static/icon-empty.png" mode="aspectFit" />
        <text class="empty-text">No exam records</text>
      </view>

      <view
        v-for="exam in examList"
        :key="exam.id"
        class="exam-card"
        @tap="goToResult(exam)"
      >
        <view class="exam-header">
          <text class="exam-title">{{ exam.title }}</text>
          <view class="status-badge" :class="exam.status">
            <text class="status-text">{{ statusMap[exam.status] || exam.status }}</text>
          </view>
        </view>
        <view class="exam-info">
          <view class="info-item">
            <image class="info-icon" src="/static/icon-score.png" mode="aspectFit" />
            <text class="info-label">Score:</text>
            <text class="info-value highlight">{{ exam.score ?? '--' }}</text>
            <text class="info-unit">/ {{ exam.totalScore }}</text>
          </view>
          <view class="info-item">
            <image class="info-icon" src="/static/icon-clock.png" mode="aspectFit" />
            <text class="info-label">Duration:</text>
            <text class="info-value">{{ formatDuration(exam.duration) }}</text>
          </view>
        </view>
        <view class="exam-footer">
          <text class="exam-date">{{ exam.submitTime || exam.createTime }}</text>
          <view class="detail-btn">
            <text class="detail-text">View Details</text>
            <image class="detail-arrow" src="/static/icon-arrow-right.png" mode="aspectFit" />
          </view>
        </view>
      </view>

      <view v-if="loading" class="loading-wrap">
        <text class="loading-text">Loading...</text>
      </view>
      <view v-if="noMore && examList.length > 0" class="no-more-wrap">
        <text class="no-more-text">No more records</text>
      </view>
    </scroll-view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getExamHistory } from '@/api/exam'

interface ExamRecord {
  id: string
  title: string
  score: number | null
  totalScore: number
  status: 'completed' | 'in_progress' | 'timeout' | 'graded'
  duration: number
  submitTime: string
  createTime: string
}

const statusTabs = [
  { label: 'All', value: 'all' },
  { label: 'Completed', value: 'completed' },
  { label: 'In Progress', value: 'in_progress' },
  { label: 'Graded', value: 'graded' }
]

const statusMap: Record<string, string> = {
  completed: 'Completed',
  in_progress: 'In Progress',
  timeout: 'Timed Out',
  graded: 'Graded'
}

const currentStatus = ref('all')
const examList = ref<ExamRecord[]>([])
const loading = ref(false)
const noMore = ref(false)
const page = ref(1)
const pageSize = 20

onMounted(() => {
  fetchHistory()
})

async function fetchHistory() {
  if (loading.value) return
  loading.value = true
  try {
    const params: any = {
      page: page.value,
      pageSize
    }
    if (currentStatus.value !== 'all') {
      params.status = currentStatus.value
    }
    const res = await getExamHistory(params)
    const list = res.data?.list || []
    if (page.value === 1) {
      examList.value = list
    } else {
      examList.value.push(...list)
    }
    if (list.length < pageSize) {
      noMore.value = true
    }
  } catch (e) {
    console.error('Failed to load exam history', e)
  } finally {
    loading.value = false
  }
}

function switchStatus(status: string) {
  if (currentStatus.value === status) return
  currentStatus.value = status
  page.value = 1
  noMore.value = false
  examList.value = []
  fetchHistory()
}

function loadMore() {
  if (noMore.value || loading.value) return
  page.value++
  fetchHistory()
}

function formatDuration(seconds: number): string {
  if (!seconds) return '--'
  const mins = Math.floor(seconds / 60)
  const secs = seconds % 60
  if (mins > 60) {
    const hrs = Math.floor(mins / 60)
    const remainMins = mins % 60
    return `${hrs}h ${remainMins}m`
  }
  return `${mins}m ${secs}s`
}

function goToResult(exam: ExamRecord) {
  uni.navigateTo({
    url: `/pages/exam/result?id=${exam.id}`
  })
}
</script>

<style scoped>
.exam-history-page {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background-color: #f5f5f5;
}

.filter-bar {
  display: flex;
  flex-direction: row;
  background-color: #fff;
  padding: 20rpx 24rpx;
  gap: 16rpx;
  border-bottom: 1rpx solid #f0f0f0;
}

.filter-tab {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 14rpx 0;
  border-radius: 30rpx;
  background-color: #f5f5f5;
}

.filter-tab.active {
  background-color: #4a90d9;
}

.tab-text {
  font-size: 26rpx;
  color: #666;
}

.tab-text.active {
  color: #fff;
  font-weight: 600;
}

.exam-list {
  flex: 1;
  padding: 20rpx 24rpx;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 200rpx 0;
}

.empty-icon {
  width: 160rpx;
  height: 160rpx;
  margin-bottom: 20rpx;
  opacity: 0.6;
}

.empty-text {
  font-size: 28rpx;
  color: #999;
}

.exam-card {
  background-color: #fff;
  border-radius: 16rpx;
  padding: 28rpx;
  margin-bottom: 20rpx;
}

.exam-header {
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: space-between;
}

.exam-title {
  font-size: 30rpx;
  color: #333;
  font-weight: 600;
  flex: 1;
  margin-right: 16rpx;
}

.status-badge {
  padding: 6rpx 16rpx;
  border-radius: 8rpx;
  flex-shrink: 0;
}

.status-badge.completed {
  background-color: #e6f7e6;
}

.status-badge.in_progress {
  background-color: #fff3e0;
}

.status-badge.timeout {
  background-color: #ffeaea;
}

.status-badge.graded {
  background-color: #e6f0ff;
}

.status-text {
  font-size: 22rpx;
  color: #666;
}

.exam-info {
  display: flex;
  flex-direction: row;
  gap: 40rpx;
  margin-top: 20rpx;
}

.info-item {
  display: flex;
  flex-direction: row;
  align-items: center;
}

.info-icon {
  width: 32rpx;
  height: 32rpx;
  margin-right: 8rpx;
}

.info-label {
  font-size: 26rpx;
  color: #999;
}

.info-value {
  font-size: 26rpx;
  color: #333;
  margin-left: 6rpx;
}

.info-value.highlight {
  color: #4a90d9;
  font-weight: 600;
  font-size: 28rpx;
}

.info-unit {
  font-size: 24rpx;
  color: #999;
}

.exam-footer {
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: space-between;
  margin-top: 20rpx;
  padding-top: 16rpx;
  border-top: 1rpx solid #f5f5f5;
}

.exam-date {
  font-size: 24rpx;
  color: #bbb;
}

.detail-btn {
  display: flex;
  flex-direction: row;
  align-items: center;
}

.detail-text {
  font-size: 26rpx;
  color: #4a90d9;
}

.detail-arrow {
  width: 24rpx;
  height: 24rpx;
  margin-left: 6rpx;
}

.loading-wrap,
.no-more-wrap {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 30rpx 0;
}

.loading-text,
.no-more-text {
  font-size: 24rpx;
  color: #999;
}
</style>
