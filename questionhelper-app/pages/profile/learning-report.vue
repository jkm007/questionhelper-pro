<template>
  <view class="report-page">
    <!-- Period Toggle -->
    <view class="period-toggle">
      <view
        v-for="p in periods"
        :key="p.value"
        class="period-tab"
        :class="{ active: period === p.value }"
        @tap="period = p.value"
      >
        <text class="period-text">{{ p.label }}</text>
      </view>
    </view>

    <!-- Summary Cards -->
    <view class="summary-row">
      <view class="summary-card">
        <text class="summary-value">{{ report.totalQuestions }}</text>
        <text class="summary-label">做题总数</text>
      </view>
      <view class="summary-card">
        <text class="summary-value">{{ report.accuracy }}%</text>
        <text class="summary-label">正确率</text>
      </view>
      <view class="summary-card">
        <text class="summary-value">{{ report.studyTime }}h</text>
        <text class="summary-label">学习时长</text>
      </view>
    </view>

    <!-- Exam Trend -->
    <view class="section">
      <text class="section-title">考试成绩趋势</text>
      <view v-if="report.examScores.length" class="trend-list">
        <view v-for="(item, i) in report.examScores" :key="i" class="trend-item">
          <text class="trend-name">{{ item.name }}</text>
          <text class="trend-score" :class="{ passed: item.passed }">{{ item.score }}分</text>
        </view>
      </view>
      <view v-else class="empty-hint">
        <text class="empty-text">暂无考试记录</text>
      </view>
    </view>

    <!-- Strengths & Weaknesses -->
    <view class="section">
      <text class="section-title">知识掌握情况</text>
      <view v-for="(cat, i) in report.categories" :key="i" class="category-item">
        <view class="category-header">
          <text class="category-name">{{ cat.name }}</text>
          <text class="category-rate" :class="{ weak: cat.accuracy < 60 }">{{ cat.accuracy }}%</text>
        </view>
        <view class="progress-bar">
          <view class="progress-fill" :style="{ width: cat.accuracy + '%' }" :class="{ weak: cat.accuracy < 60 }" />
        </view>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

const periods = [
  { label: '本周', value: 'week' },
  { label: '本月', value: 'month' }
]

const period = ref('week')

interface ExamScore { name: string; score: number; passed: boolean }
interface Category { name: string; accuracy: number }
interface Report {
  totalQuestions: number
  accuracy: number
  studyTime: number
  examScores: ExamScore[]
  categories: Category[]
}

const report = ref<Report>({
  totalQuestions: 0,
  accuracy: 0,
  studyTime: 0,
  examScores: [],
  categories: []
})

onMounted(() => { fetchReport() })

async function fetchReport() {
  try {
    // TODO: replace with actual API call
    report.value = {
      totalQuestions: 326,
      accuracy: 78,
      studyTime: 12.5,
      examScores: [
        { name: '期中考试', score: 85, passed: true },
        { name: '单元测试', score: 72, passed: true },
        { name: '随堂测验', score: 55, passed: false }
      ],
      categories: [
        { name: '数据结构', accuracy: 88 },
        { name: '算法设计', accuracy: 65 },
        { name: '操作系统', accuracy: 45 },
        { name: '计算机网络', accuracy: 72 }
      ]
    }
  } catch (e) {
    console.error('Failed to load report', e)
  }
}
</script>

<style scoped>
.report-page {
  min-height: 100vh;
  background-color: #f5f5f5;
  padding: 20rpx 24rpx;
  padding-bottom: 80rpx;
}

.period-toggle {
  display: flex;
  flex-direction: row;
  background-color: #fff;
  border-radius: 16rpx;
  overflow: hidden;
  margin-bottom: 24rpx;
}

.period-tab {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 24rpx 0;
}

.period-tab.active {
  background-color: #4a90d9;
}

.period-text {
  font-size: 28rpx;
  color: #666;
}

.period-tab.active .period-text {
  color: #fff;
}

.summary-row {
  display: flex;
  flex-direction: row;
  gap: 16rpx;
  margin-bottom: 24rpx;
}

.summary-card {
  flex: 1;
  background-color: #fff;
  border-radius: 16rpx;
  padding: 28rpx 16rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.summary-value {
  font-size: 40rpx;
  font-weight: 700;
  color: #4a90d9;
  margin-bottom: 8rpx;
}

.summary-label {
  font-size: 24rpx;
  color: #999;
}

.section {
  background-color: #fff;
  border-radius: 16rpx;
  padding: 28rpx 30rpx;
  margin-bottom: 24rpx;
}

.section-title {
  font-size: 30rpx;
  font-weight: 600;
  color: #333;
  margin-bottom: 20rpx;
}

.trend-list {
  display: flex;
  flex-direction: column;
  gap: 16rpx;
}

.trend-item {
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  align-items: center;
  padding: 16rpx 0;
  border-bottom: 1rpx solid #f5f5f5;
}

.trend-name {
  font-size: 28rpx;
  color: #333;
}

.trend-score {
  font-size: 28rpx;
  font-weight: 600;
  color: #e74c3c;
}

.trend-score.passed {
  color: #27ae60;
}

.empty-hint {
  padding: 40rpx 0;
  display: flex;
  align-items: center;
  justify-content: center;
}

.empty-text {
  font-size: 26rpx;
  color: #ccc;
}

.category-item {
  margin-bottom: 20rpx;
}

.category-header {
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  margin-bottom: 10rpx;
}

.category-name {
  font-size: 26rpx;
  color: #333;
}

.category-rate {
  font-size: 26rpx;
  color: #4a90d9;
  font-weight: 600;
}

.category-rate.weak {
  color: #e74c3c;
}

.progress-bar {
  height: 14rpx;
  background-color: #f0f0f0;
  border-radius: 7rpx;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background-color: #4a90d9;
  border-radius: 7rpx;
  transition: width 0.3s;
}

.progress-fill.weak {
  background-color: #e74c3c;
}
</style>
