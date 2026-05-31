<template>
  <view class="statistics-page">
    <!-- Summary Cards -->
    <view class="summary-row">
      <view class="summary-card">
        <text class="summary-value">{{ summary.totalPractice }}</text>
        <text class="summary-label">Total Practice</text>
      </view>
      <view class="summary-card">
        <text class="summary-value">{{ summary.totalExams }}</text>
        <text class="summary-label">Total Exams</text>
      </view>
      <view class="summary-card">
        <text class="summary-value">{{ summary.avgAccuracy }}%</text>
        <text class="summary-label">Avg Accuracy</text>
      </view>
    </view>

    <!-- Practice Trend Chart -->
    <view class="chart-section">
      <text class="section-title">Practice Trend (Recent 30 Days)</text>
      <view class="chart-wrap">
        <canvas
          canvas-id="practiceTrend"
          id="practiceTrend"
          class="chart-canvas"
          type="2d"
        />
      </view>
      <view v-if="!practiceTrendData.length" class="chart-empty">
        <text class="chart-empty-text">No data available</text>
      </view>
    </view>

    <!-- Exam Score Trend Chart -->
    <view class="chart-section">
      <text class="section-title">Exam Score Trend</text>
      <view class="chart-wrap">
        <canvas
          canvas-id="scoreTrend"
          id="scoreTrend"
          class="chart-canvas"
          type="2d"
        />
      </view>
      <view v-if="!scoreTrendData.length" class="chart-empty">
        <text class="chart-empty-text">No data available</text>
      </view>
    </view>

    <!-- Category Accuracy Distribution -->
    <view class="chart-section">
      <text class="section-title">Category Accuracy Distribution</text>
      <view class="category-list">
        <view v-for="item in categoryData" :key="item.name" class="category-item">
          <view class="category-header">
            <text class="category-name">{{ item.name }}</text>
            <text class="category-accuracy">{{ item.accuracy }}%</text>
          </view>
          <view class="progress-bar">
            <view class="progress-fill" :style="{ width: item.accuracy + '%' }" />
          </view>
          <text class="category-count">{{ item.correct }}/{{ item.total }} correct</text>
        </view>
      </view>
      <view v-if="!categoryData.length" class="chart-empty">
        <text class="chart-empty-text">No data available</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted, nextTick } from 'vue'
import { getPersonalStatistics, getPracticeTrend, getScoreTrend, getCategoryAccuracy } from '@/api/statistics'

interface Summary {
  totalPractice: number
  totalExams: number
  avgAccuracy: number
}

interface TrendPoint {
  date: string
  count?: number
  score?: number
}

interface CategoryItem {
  name: string
  total: number
  correct: number
  accuracy: number
}

const summary = ref<Summary>({
  totalPractice: 0,
  totalExams: 0,
  avgAccuracy: 0
})

const practiceTrendData = ref<TrendPoint[]>([])
const scoreTrendData = ref<TrendPoint[]>([])
const categoryData = ref<CategoryItem[]>([])

onMounted(async () => {
  await fetchAllData()
  nextTick(() => {
    drawPracticeTrend()
    drawScoreTrend()
  })
})

async function fetchAllData() {
  try {
    const [summaryRes, practiceRes, scoreRes, categoryRes] = await Promise.all([
      getPersonalStatistics(),
      getPracticeTrend({ days: 30 }),
      getScoreTrend({ limit: 10 }),
      getCategoryAccuracy()
    ])
    if (summaryRes.data) {
      summary.value = summaryRes.data
    }
    practiceTrendData.value = practiceRes.data?.list || []
    scoreTrendData.value = scoreRes.data?.list || []
    categoryData.value = categoryRes.data?.list || []
  } catch (e) {
    console.error('Failed to load statistics', e)
  }
}

function drawPracticeTrend() {
  if (!practiceTrendData.value.length) return
  const ctx = uni.createCanvasContext('practiceTrend')
  const data = practiceTrendData.value
  const width = 680
  const height = 360
  const padding = { top: 30, right: 30, bottom: 50, left: 50 }
  const chartW = width - padding.left - padding.right
  const chartH = height - padding.top - padding.bottom

  const maxVal = Math.max(...data.map((d) => d.count || 0), 1)
  const stepX = chartW / Math.max(data.length - 1, 1)

  // Draw grid lines
  ctx.setStrokeStyle('#f0f0f0')
  ctx.setLineWidth(1)
  for (let i = 0; i <= 4; i++) {
    const y = padding.top + (chartH / 4) * i
    ctx.beginPath()
    ctx.moveTo(padding.left, y)
    ctx.lineTo(width - padding.right, y)
    ctx.stroke()
  }

  // Draw line
  ctx.setStrokeStyle('#4a90d9')
  ctx.setLineWidth(2)
  ctx.beginPath()
  data.forEach((point, index) => {
    const x = padding.left + stepX * index
    const y = padding.top + chartH - ((point.count || 0) / maxVal) * chartH
    if (index === 0) {
      ctx.moveTo(x, y)
    } else {
      ctx.lineTo(x, y)
    }
  })
  ctx.stroke()

  // Draw points
  data.forEach((point, index) => {
    const x = padding.left + stepX * index
    const y = padding.top + chartH - ((point.count || 0) / maxVal) * chartH
    ctx.beginPath()
    ctx.arc(x, y, 4, 0, Math.PI * 2)
    ctx.setFillStyle('#4a90d9')
    ctx.fill()
  })

  ctx.draw()
}

function drawScoreTrend() {
  if (!scoreTrendData.value.length) return
  const ctx = uni.createCanvasContext('scoreTrend')
  const data = scoreTrendData.value
  const width = 680
  const height = 360
  const padding = { top: 30, right: 30, bottom: 50, left: 50 }
  const chartW = width - padding.left - padding.right
  const chartH = height - padding.top - padding.bottom

  const maxVal = 100
  const stepX = chartW / Math.max(data.length - 1, 1)

  // Draw grid lines
  ctx.setStrokeStyle('#f0f0f0')
  ctx.setLineWidth(1)
  for (let i = 0; i <= 4; i++) {
    const y = padding.top + (chartH / 4) * i
    ctx.beginPath()
    ctx.moveTo(padding.left, y)
    ctx.lineTo(width - padding.right, y)
    ctx.stroke()
  }

  // Draw bars
  const barWidth = Math.min(stepX * 0.5, 40)
  data.forEach((point, index) => {
    const x = padding.left + stepX * index - barWidth / 2
    const barH = ((point.score || 0) / maxVal) * chartH
    const y = padding.top + chartH - barH
    ctx.setFillStyle('#67c23a')
    ctx.fillRect(x, y, barWidth, barH)
  })

  ctx.draw()
}
</script>

<style scoped>
.statistics-page {
  min-height: 100vh;
  background-color: #f5f5f5;
  padding: 24rpx;
  padding-bottom: 40rpx;
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
  padding: 24rpx 16rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.summary-value {
  font-size: 40rpx;
  color: #4a90d9;
  font-weight: 700;
}

.summary-label {
  font-size: 22rpx;
  color: #999;
  margin-top: 8rpx;
}

.chart-section {
  background-color: #fff;
  border-radius: 16rpx;
  padding: 28rpx;
  margin-bottom: 24rpx;
}

.section-title {
  font-size: 30rpx;
  color: #333;
  font-weight: 600;
  margin-bottom: 20rpx;
}

.chart-wrap {
  width: 100%;
  height: 360rpx;
  position: relative;
}

.chart-canvas {
  width: 100%;
  height: 100%;
}

.chart-empty {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 200rpx;
}

.chart-empty-text {
  font-size: 26rpx;
  color: #ccc;
}

.category-list {
  padding-top: 8rpx;
}

.category-item {
  margin-bottom: 24rpx;
}

.category-item:last-child {
  margin-bottom: 0;
}

.category-header {
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10rpx;
}

.category-name {
  font-size: 26rpx;
  color: #333;
}

.category-accuracy {
  font-size: 26rpx;
  color: #4a90d9;
  font-weight: 600;
}

.progress-bar {
  width: 100%;
  height: 16rpx;
  background-color: #f0f0f0;
  border-radius: 8rpx;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: linear-gradient(90deg, #4a90d9, #67c23a);
  border-radius: 8rpx;
  transition: width 0.3s ease;
}

.category-count {
  font-size: 22rpx;
  color: #bbb;
  margin-top: 6rpx;
}
</style>
