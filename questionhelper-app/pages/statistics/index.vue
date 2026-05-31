<template>
  <view class="statistics-page">
    <!-- Summary Cards -->
    <view class="summary-grid">
      <view class="summary-card blue">
        <text class="summary-value">{{ summary.practiceCount }}</text>
        <text class="summary-label">Practice Sessions</text>
      </view>
      <view class="summary-card green">
        <text class="summary-value">{{ summary.examCount }}</text>
        <text class="summary-label">Exams Taken</text>
      </view>
      <view class="summary-card orange">
        <text class="summary-value">{{ summary.accuracy }}%</text>
        <text class="summary-label">Accuracy</text>
      </view>
      <view class="summary-card purple">
        <text class="summary-value">{{ formatDuration(summary.studyDuration) }}</text>
        <text class="summary-label">Study Duration</text>
      </view>
    </view>

    <!-- Weekly Practice Trend -->
    <view class="chart-section">
      <text class="section-title">Weekly Practice Trend</text>
      <view class="chart-wrap">
        <canvas
          canvas-id="weeklyPractice"
          id="weeklyPractice"
          class="chart-canvas"
          type="2d"
        />
      </view>
      <view v-if="!weeklyPracticeData.length" class="chart-empty">
        <text class="chart-empty-text">No practice data this week</text>
      </view>
    </view>

    <!-- Category Accuracy Radar -->
    <view class="chart-section">
      <text class="section-title">Category Accuracy</text>
      <view class="chart-wrap radar-wrap">
        <canvas
          canvas-id="categoryRadar"
          id="categoryRadar"
          class="chart-canvas"
          type="2d"
        />
      </view>
      <view class="category-legend">
        <view v-for="item in categoryData" :key="item.name" class="legend-item">
          <view class="legend-dot" :style="{ backgroundColor: item.color }" />
          <text class="legend-name">{{ item.name }}</text>
          <text class="legend-value">{{ item.accuracy }}%</text>
        </view>
      </view>
    </view>

    <!-- Exam Score Trend -->
    <view class="chart-section">
      <text class="section-title">Exam Score Trend</text>
      <view class="chart-wrap">
        <canvas
          canvas-id="examScoreTrend"
          id="examScoreTrend"
          class="chart-canvas"
          type="2d"
        />
      </view>
      <view v-if="!examScoreData.length" class="chart-empty">
        <text class="chart-empty-text">No exam data available</text>
      </view>
    </view>

    <!-- Rankings -->
    <view class="ranking-section">
      <text class="section-title">Rankings</text>
      <view v-if="rankings.length === 0" class="chart-empty">
        <text class="chart-empty-text">No ranking data</text>
      </view>
      <view v-for="(rank, index) in rankings" :key="rank.userId" class="rank-item">
        <view class="rank-position" :class="{ top: index < 3 }">
          <text class="rank-num">{{ index + 1 }}</text>
        </view>
        <image class="rank-avatar" :src="rank.avatar || '/static/default-avatar.png'" mode="aspectFill" />
        <view class="rank-info">
          <text class="rank-name">{{ rank.nickname }}</text>
          <text class="rank-stats">{{ rank.practiceCount }} practices | {{ rank.accuracy }}% accuracy</text>
        </view>
        <view v-if="rank.isCurrentUser" class="you-badge">
          <text class="you-text">You</text>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted, nextTick } from 'vue'
import {
  getStatisticsSummary,
  getWeeklyPractice,
  getCategoryAccuracy,
  getExamScoreTrend,
  getRankings
} from '@/api/statistics'

interface Summary {
  practiceCount: number
  examCount: number
  accuracy: number
  studyDuration: number
}

interface WeekDay {
  day: string
  count: number
}

interface CategoryItem {
  name: string
  accuracy: number
  color: string
}

interface ScorePoint {
  title: string
  score: number
}

interface RankItem {
  userId: string
  nickname: string
  avatar: string
  practiceCount: number
  accuracy: number
  isCurrentUser: boolean
}

const summary = ref<Summary>({
  practiceCount: 0,
  examCount: 0,
  accuracy: 0,
  studyDuration: 0
})

const weeklyPracticeData = ref<WeekDay[]>([])
const categoryData = ref<CategoryItem[]>([])
const examScoreData = ref<ScorePoint[]>([])
const rankings = ref<RankItem[]>([])

const categoryColors = ['#4a90d9', '#67c23a', '#e6a23c', '#e74c3c', '#9b59b6', '#1abc9c']

onMounted(async () => {
  await fetchAllData()
  nextTick(() => {
    drawWeeklyPractice()
    drawCategoryRadar()
    drawExamScoreTrend()
  })
})

async function fetchAllData() {
  try {
    const [summaryRes, weeklyRes, categoryRes, scoreRes, rankingRes] = await Promise.all([
      getStatisticsSummary(),
      getWeeklyPractice(),
      getCategoryAccuracy(),
      getExamScoreTrend({ limit: 10 }),
      getRankings({ limit: 10 })
    ])
    if (summaryRes.data) {
      summary.value = summaryRes.data
    }
    weeklyPracticeData.value = weeklyRes.data?.list || []
    const rawCategories = categoryRes.data?.list || []
    categoryData.value = rawCategories.map((item: any, index: number) => ({
      ...item,
      color: categoryColors[index % categoryColors.length]
    }))
    examScoreData.value = scoreRes.data?.list || []
    rankings.value = rankingRes.data?.list || []
  } catch (e) {
    console.error('Failed to load statistics', e)
  }
}

function formatDuration(minutes: number): string {
  if (!minutes) return '0h'
  if (minutes < 60) return minutes + 'min'
  const hrs = Math.floor(minutes / 60)
  const mins = minutes % 60
  return mins > 0 ? `${hrs}h ${mins}m` : `${hrs}h`
}

function drawWeeklyPractice() {
  if (!weeklyPracticeData.value.length) return
  const ctx = uni.createCanvasContext('weeklyPractice')
  const data = weeklyPracticeData.value
  const width = 680
  const height = 320
  const padding = { top: 20, right: 20, bottom: 50, left: 20 }
  const chartW = width - padding.left - padding.right
  const chartH = height - padding.top - padding.bottom

  const maxVal = Math.max(...data.map((d) => d.count), 1)
  const barWidth = Math.min(chartW / data.length * 0.6, 60)
  const gap = chartW / data.length

  // Draw bars
  data.forEach((item, index) => {
    const x = padding.left + gap * index + (gap - barWidth) / 2
    const barH = (item.count / maxVal) * chartH
    const y = padding.top + chartH - barH

    // Bar gradient
    ctx.setFillStyle('#4a90d9')
    ctx.fillRect(x, y, barWidth, barH)

    // Day label
    ctx.setFillStyle('#999')
    ctx.setFontSize(11)
    ctx.setTextAlign('center')
    ctx.fillText(item.day, x + barWidth / 2, height - padding.bottom + 20)

    // Count label
    if (item.count > 0) {
      ctx.setFillStyle('#4a90d9')
      ctx.fillText(String(item.count), x + barWidth / 2, y - 8)
    }
  })

  ctx.draw()
}

function drawCategoryRadar() {
  if (!categoryData.value.length) return
  const ctx = uni.createCanvasContext('categoryRadar')
  const data = categoryData.value
  const centerX = 340
  const centerY = 200
  const radius = 140
  const sides = data.length

  if (sides < 3) return

  const angleStep = (Math.PI * 2) / sides

  // Draw grid rings
  for (let ring = 1; ring <= 4; ring++) {
    const r = (radius / 4) * ring
    ctx.setStrokeStyle('#f0f0f0')
    ctx.setLineWidth(1)
    ctx.beginPath()
    for (let i = 0; i <= sides; i++) {
      const angle = angleStep * i - Math.PI / 2
      const x = centerX + r * Math.cos(angle)
      const y = centerY + r * Math.sin(angle)
      if (i === 0) {
        ctx.moveTo(x, y)
      } else {
        ctx.lineTo(x, y)
      }
    }
    ctx.stroke()
  }

  // Draw axis lines
  for (let i = 0; i < sides; i++) {
    const angle = angleStep * i - Math.PI / 2
    ctx.setStrokeStyle('#e0e0e0')
    ctx.beginPath()
    ctx.moveTo(centerX, centerY)
    ctx.lineTo(centerX + radius * Math.cos(angle), centerY + radius * Math.sin(angle))
    ctx.stroke()
  }

  // Draw data polygon
  ctx.beginPath()
  data.forEach((item, index) => {
    const angle = angleStep * index - Math.PI / 2
    const r = (item.accuracy / 100) * radius
    const x = centerX + r * Math.cos(angle)
    const y = centerY + r * Math.sin(angle)
    if (index === 0) {
      ctx.moveTo(x, y)
    } else {
      ctx.lineTo(x, y)
    }
  })
  ctx.closePath()
  ctx.setFillStyle('rgba(74, 144, 217, 0.2)')
  ctx.fill()
  ctx.setStrokeStyle('#4a90d9')
  ctx.setLineWidth(2)
  ctx.stroke()

  // Draw data points
  data.forEach((item, index) => {
    const angle = angleStep * index - Math.PI / 2
    const r = (item.accuracy / 100) * radius
    const x = centerX + r * Math.cos(angle)
    const y = centerY + r * Math.sin(angle)
    ctx.beginPath()
    ctx.arc(x, y, 4, 0, Math.PI * 2)
    ctx.setFillStyle('#4a90d9')
    ctx.fill()
  })

  // Draw labels
  ctx.setFillStyle('#666')
  ctx.setFontSize(11)
  ctx.setTextAlign('center')
  data.forEach((item, index) => {
    const angle = angleStep * index - Math.PI / 2
    const labelR = radius + 25
    const x = centerX + labelR * Math.cos(angle)
    const y = centerY + labelR * Math.sin(angle)
    ctx.fillText(item.name, x, y + 4)
  })

  ctx.draw()
}

function drawExamScoreTrend() {
  if (!examScoreData.value.length) return
  const ctx = uni.createCanvasContext('examScoreTrend')
  const data = examScoreData.value
  const width = 680
  const height = 320
  const padding = { top: 30, right: 30, bottom: 50, left: 50 }
  const chartW = width - padding.left - padding.right
  const chartH = height - padding.top - padding.bottom

  const maxVal = 100
  const stepX = chartW / Math.max(data.length - 1, 1)

  // Grid lines
  ctx.setStrokeStyle('#f0f0f0')
  ctx.setLineWidth(1)
  for (let i = 0; i <= 5; i++) {
    const y = padding.top + (chartH / 5) * i
    ctx.beginPath()
    ctx.moveTo(padding.left, y)
    ctx.lineTo(width - padding.right, y)
    ctx.stroke()
  }

  // Draw area fill
  ctx.beginPath()
  ctx.moveTo(padding.left, padding.top + chartH)
  data.forEach((point, index) => {
    const x = padding.left + stepX * index
    const y = padding.top + chartH - (point.score / maxVal) * chartH
    ctx.lineTo(x, y)
  })
  ctx.lineTo(padding.left + stepX * (data.length - 1), padding.top + chartH)
  ctx.closePath()
  ctx.setFillStyle('rgba(103, 194, 58, 0.15)')
  ctx.fill()

  // Draw line
  ctx.setStrokeStyle('#67c23a')
  ctx.setLineWidth(2)
  ctx.beginPath()
  data.forEach((point, index) => {
    const x = padding.left + stepX * index
    const y = padding.top + chartH - (point.score / maxVal) * chartH
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
    const y = padding.top + chartH - (point.score / maxVal) * chartH
    ctx.beginPath()
    ctx.arc(x, y, 5, 0, Math.PI * 2)
    ctx.setFillStyle('#67c23a')
    ctx.fill()
    ctx.beginPath()
    ctx.arc(x, y, 3, 0, Math.PI * 2)
    ctx.setFillStyle('#fff')
    ctx.fill()
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

.summary-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 16rpx;
  margin-bottom: 24rpx;
}

.summary-card {
  width: calc(50% - 8rpx);
  border-radius: 16rpx;
  padding: 28rpx 24rpx;
  display: flex;
  flex-direction: column;
}

.summary-card.blue {
  background: linear-gradient(135deg, #4a90d9, #357abd);
}

.summary-card.green {
  background: linear-gradient(135deg, #67c23a, #529b2e);
}

.summary-card.orange {
  background: linear-gradient(135deg, #e6a23c, #cf9236);
}

.summary-card.purple {
  background: linear-gradient(135deg, #9b59b6, #8e44ad);
}

.summary-value {
  font-size: 44rpx;
  color: #fff;
  font-weight: 700;
}

.summary-label {
  font-size: 24rpx;
  color: rgba(255, 255, 255, 0.8);
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
  height: 320rpx;
}

.radar-wrap {
  height: 400rpx;
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

.category-legend {
  display: flex;
  flex-wrap: wrap;
  gap: 16rpx;
  margin-top: 20rpx;
}

.legend-item {
  display: flex;
  flex-direction: row;
  align-items: center;
  width: calc(50% - 8rpx);
}

.legend-dot {
  width: 16rpx;
  height: 16rpx;
  border-radius: 8rpx;
  margin-right: 10rpx;
}

.legend-name {
  font-size: 24rpx;
  color: #666;
  flex: 1;
}

.legend-value {
  font-size: 24rpx;
  color: #333;
  font-weight: 600;
}

.ranking-section {
  background-color: #fff;
  border-radius: 16rpx;
  padding: 28rpx;
  margin-bottom: 24rpx;
}

.rank-item {
  display: flex;
  flex-direction: row;
  align-items: center;
  padding: 20rpx 0;
  border-bottom: 1rpx solid #f5f5f5;
}

.rank-item:last-child {
  border-bottom: none;
}

.rank-position {
  width: 48rpx;
  height: 48rpx;
  border-radius: 24rpx;
  background-color: #f5f5f5;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 16rpx;
}

.rank-position.top {
  background-color: #ff9800;
}

.rank-num {
  font-size: 26rpx;
  color: #666;
  font-weight: 600;
}

.rank-position.top .rank-num {
  color: #fff;
}

.rank-avatar {
  width: 64rpx;
  height: 64rpx;
  border-radius: 32rpx;
  margin-right: 16rpx;
}

.rank-info {
  flex: 1;
}

.rank-name {
  font-size: 28rpx;
  color: #333;
  font-weight: 500;
}

.rank-stats {
  font-size: 22rpx;
  color: #999;
  margin-top: 4rpx;
}

.you-badge {
  padding: 6rpx 16rpx;
  background-color: #e6f0ff;
  border-radius: 12rpx;
}

.you-text {
  font-size: 22rpx;
  color: #4a90d9;
}
</style>
