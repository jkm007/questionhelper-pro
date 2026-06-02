<template>
  <view class="report-page">
    <!-- Period Selector -->
    <view class="period-tabs">
      <view
        class="period-tab"
        :class="{ active: activePeriod === 'week' }"
        @tap="switchPeriod('week')"
      >
        <text class="period-tab-text">Weekly</text>
      </view>
      <view
        class="period-tab"
        :class="{ active: activePeriod === 'month' }"
        @tap="switchPeriod('month')"
      >
        <text class="period-tab-text">Monthly</text>
      </view>
    </view>

    <!-- Study Time Summary -->
    <view class="summary-card">
      <view class="summary-main">
        <text class="summary-value">{{ formatTime(report.totalStudyMinutes) }}</text>
        <text class="summary-label">Total Study Time</text>
      </view>
      <view class="summary-divider" />
      <view class="summary-sub">
        <view class="summary-sub-item">
          <text class="sub-value">{{ report.totalPractice }}</text>
          <text class="sub-label">Practice Sessions</text>
        </view>
        <view class="summary-sub-item">
          <text class="sub-value">{{ report.totalExams }}</text>
          <text class="sub-label">Exams Taken</text>
        </view>
      </view>
    </view>

    <!-- Study Time Chart -->
    <view class="chart-section">
      <text class="section-title">Study Time Distribution</text>
      <view class="chart-wrap">
        <canvas canvas-id="studyTimeChart" id="studyTimeChart" class="chart-canvas" type="2d" />
      </view>
      <view v-if="!report.dailyStudyTime.length" class="chart-empty">
        <text class="chart-empty-text">No data available</text>
      </view>
    </view>

    <!-- Practice Stats -->
    <view class="stats-card">
      <text class="section-title">Practice Performance</text>
      <view class="stats-grid">
        <view class="stats-item">
          <text class="stats-number">{{ report.practiceAccuracy }}%</text>
          <text class="stats-desc">Accuracy Rate</text>
        </view>
        <view class="stats-item">
          <text class="stats-number">{{ report.practiceTotal }}</text>
          <text class="stats-desc">Questions Practiced</text>
        </view>
        <view class="stats-item">
          <text class="stats-number">{{ report.practiceCorrect }}</text>
          <text class="stats-desc">Correct Answers</text>
        </view>
        <view class="stats-item">
          <text class="stats-number">{{ report.avgPracticeTime }}s</text>
          <text class="stats-desc">Avg Time/Question</text>
        </view>
      </view>
    </view>

    <!-- Exam Scores -->
    <view class="chart-section">
      <text class="section-title">Exam Score Trend</text>
      <view class="chart-wrap">
        <canvas canvas-id="scoreChart" id="scoreChart" class="chart-canvas" type="2d" />
      </view>
      <view v-if="!report.examScores.length" class="chart-empty">
        <text class="chart-empty-text">No exam data</text>
      </view>
    </view>

    <!-- Strengths & Weaknesses -->
    <view class="analysis-card">
      <text class="section-title">Strengths &amp; Weaknesses</text>

      <!-- Strengths -->
      <view v-if="report.strengths.length" class="analysis-section">
        <view class="analysis-header">
          <view class="dot green" />
          <text class="analysis-label">Strengths</text>
        </view>
        <view v-for="item in report.strengths" :key="item.category" class="analysis-item">
          <text class="item-name">{{ item.category }}</text>
          <view class="item-bar-wrap">
            <view class="item-bar green-bar" :style="{ width: item.accuracy + '%' }" />
          </view>
          <text class="item-value">{{ item.accuracy }}%</text>
        </view>
      </view>

      <!-- Weaknesses -->
      <view v-if="report.weaknesses.length" class="analysis-section">
        <view class="analysis-header">
          <view class="dot red" />
          <text class="analysis-label">Areas to Improve</text>
        </view>
        <view v-for="item in report.weaknesses" :key="item.category" class="analysis-item">
          <text class="item-name">{{ item.category }}</text>
          <view class="item-bar-wrap">
            <view class="item-bar red-bar" :style="{ width: item.accuracy + '%' }" />
          </view>
          <text class="item-value">{{ item.accuracy }}%</text>
        </view>
      </view>

      <view v-if="!report.strengths.length && !report.weaknesses.length" class="chart-empty">
        <text class="chart-empty-text">Practice more to see your analysis</text>
      </view>
    </view>

    <!-- Improvement Tips -->
    <view v-if="report.weaknesses.length" class="tips-card">
      <text class="section-title">Improvement Tips</text>
      <view v-for="(tip, index) in improvementTips" :key="index" class="tip-item">
        <text class="tip-number">{{ index + 1 }}</text>
        <text class="tip-text">{{ tip }}</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, nextTick } from 'vue'

interface CategoryAnalysis {
  category: string
  accuracy: number
  total: number
  correct: number
}

interface ExamScore {
  examName: string
  score: number
  date: string
}

interface DailyStudy {
  date: string
  minutes: number
}

interface LearningReport {
  totalStudyMinutes: number
  totalPractice: number
  totalExams: number
  practiceAccuracy: number
  practiceTotal: number
  practiceCorrect: number
  avgPracticeTime: number
  dailyStudyTime: DailyStudy[]
  examScores: ExamScore[]
  strengths: CategoryAnalysis[]
  weaknesses: CategoryAnalysis[]
}

const activePeriod = ref<'week' | 'month'>('week')

const report = ref<LearningReport>({
  totalStudyMinutes: 0,
  totalPractice: 0,
  totalExams: 0,
  practiceAccuracy: 0,
  practiceTotal: 0,
  practiceCorrect: 0,
  avgPracticeTime: 0,
  dailyStudyTime: [],
  examScores: [],
  strengths: [],
  weaknesses: []
})

const improvementTips = computed(() => {
  const tips: string[] = []
  const weaknesses = report.value.weaknesses
  if (weaknesses.length > 0) {
    tips.push(`Focus on "${weaknesses[0].category}" where your accuracy is ${weaknesses[0].accuracy}%`)
  }
  if (weaknesses.length > 1) {
    tips.push(`Schedule extra practice for "${weaknesses[1].category}"`)
  }
  if (report.value.avgPracticeTime > 60) {
    tips.push('Try to improve your answering speed through timed practice')
  }
  if (report.value.practiceAccuracy < 70) {
    tips.push('Review wrong answers carefully before moving to new questions')
  }
  if (tips.length === 0) {
    tips.push('Keep up the good work! Maintain consistent study habits')
  }
  return tips
})

onMounted(() => {
  fetchReport()
})

async function fetchReport() {
  try {
    uni.showLoading({ title: 'Loading...' })
    const period = activePeriod.value === 'week' ? 7 : 30
    const res = await new Promise<any>((resolve) => {
      uni.request({
        url: '/api/v1/statistics/learning-report',
        method: 'GET',
        data: { period },
        header: {
          Authorization: `Bearer ${uni.getStorageSync('token')}`
        },
        success: (r: any) => resolve(r.data),
        fail: () => resolve(null)
      })
    })
    uni.hideLoading()

    if (res?.code === '00000' && res.data) {
      report.value = {
        totalStudyMinutes: res.data.totalStudyMinutes || 0,
        totalPractice: res.data.totalPractice || 0,
        totalExams: res.data.totalExams || 0,
        practiceAccuracy: res.data.practiceAccuracy || 0,
        practiceTotal: res.data.practiceTotal || 0,
        practiceCorrect: res.data.practiceCorrect || 0,
        avgPracticeTime: res.data.avgPracticeTime || 0,
        dailyStudyTime: res.data.dailyStudyTime || [],
        examScores: res.data.examScores || [],
        strengths: res.data.strengths || [],
        weaknesses: res.data.weaknesses || []
      }
      nextTick(() => {
        drawStudyTimeChart()
        drawScoreChart()
      })
    } else {
      loadMockData()
    }
  } catch (e) {
    uni.hideLoading()
    console.error('Failed to load learning report', e)
    loadMockData()
  }
}

function loadMockData() {
  const days = activePeriod.value === 'week' ? 7 : 30
  const dailyStudyTime: DailyStudy[] = []
  const now = new Date()
  for (let i = days - 1; i >= 0; i--) {
    const d = new Date(now)
    d.setDate(d.getDate() - i)
    dailyStudyTime.push({
      date: `${d.getMonth() + 1}/${d.getDate()}`,
      minutes: Math.floor(Math.random() * 120) + 10
    })
  }

  report.value = {
    totalStudyMinutes: 1260,
    totalPractice: 358,
    totalExams: 12,
    practiceAccuracy: 78,
    practiceTotal: 358,
    practiceCorrect: 279,
    avgPracticeTime: 42,
    dailyStudyTime,
    examScores: [
      { examName: 'Math Quiz', score: 85, date: '5/20' },
      { examName: 'English Test', score: 72, date: '5/22' },
      { examName: 'Physics Mid', score: 91, date: '5/25' },
      { examName: 'Chemistry', score: 68, date: '5/28' },
      { examName: 'Biology', score: 88, date: '5/30' }
    ],
    strengths: [
      { category: 'Algebra', accuracy: 92, total: 50, correct: 46 },
      { category: 'Grammar', accuracy: 88, total: 40, correct: 35 }
    ],
    weaknesses: [
      { category: 'Organic Chemistry', accuracy: 45, total: 30, correct: 14 },
      { category: 'Thermodynamics', accuracy: 52, total: 25, correct: 13 }
    ]
  }

  nextTick(() => {
    drawStudyTimeChart()
    drawScoreChart()
  })
}

function switchPeriod(period: 'week' | 'month') {
  if (activePeriod.value === period) return
  activePeriod.value = period
  fetchReport()
}

function formatTime(minutes: number): string {
  if (minutes < 60) return `${minutes}min`
  const h = Math.floor(minutes / 60)
  const m = minutes % 60
  return m > 0 ? `${h}h ${m}min` : `${h}h`
}

function drawStudyTimeChart() {
  const data = report.value.dailyStudyTime
  if (!data.length) return

  const ctx = uni.createCanvasContext('studyTimeChart')
  const width = 650
  const height = 320
  const padding = { top: 20, right: 20, bottom: 50, left: 50 }
  const chartW = width - padding.left - padding.right
  const chartH = height - padding.top - padding.bottom

  const maxVal = Math.max(...data.map((d) => d.minutes), 1)
  const barW = Math.min((chartW / data.length) * 0.6, 30)
  const gap = chartW / data.length

  // Grid lines
  ctx.setStrokeStyle('#f0f0f0')
  ctx.setLineWidth(1)
  for (let i = 0; i <= 4; i++) {
    const y = padding.top + (chartH / 4) * i
    ctx.beginPath()
    ctx.moveTo(padding.left, y)
    ctx.lineTo(width - padding.right, y)
    ctx.stroke()
  }

  // Bars
  data.forEach((point, index) => {
    const x = padding.left + gap * index + (gap - barW) / 2
    const barH = (point.minutes / maxVal) * chartH
    const y = padding.top + chartH - barH

    const gradient = ctx.createLinearGradient(x, y, x, padding.top + chartH)
    gradient.addColorStop(0, '#4a90d9')
    gradient.addColorStop(1, '#a0c4e8')
    ctx.setFillStyle(gradient)
    ctx.fillRect(x, y, barW, barH)

    // Date label
    ctx.setFillStyle('#999')
    ctx.setFontSize(10)
    ctx.setTextAlign('center')
    if (data.length <= 14 || index % 2 === 0) {
      ctx.fillText(point.date, x + barW / 2, height - 10)
    }
  })

  ctx.draw()
}

function drawScoreChart() {
  const data = report.value.examScores
  if (!data.length) return

  const ctx = uni.createCanvasContext('scoreChart')
  const width = 650
  const height = 320
  const padding = { top: 20, right: 20, bottom: 50, left: 50 }
  const chartW = width - padding.left - padding.right
  const chartH = height - padding.top - padding.bottom

  const maxVal = 100
  const stepX = chartW / Math.max(data.length - 1, 1)

  // Grid lines
  ctx.setStrokeStyle('#f0f0f0')
  ctx.setLineWidth(1)
  for (let i = 0; i <= 4; i++) {
    const y = padding.top + (chartH / 4) * i
    ctx.beginPath()
    ctx.moveTo(padding.left, y)
    ctx.lineTo(width - padding.right, y)
    ctx.stroke()
  }

  // Line
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

  // Points
  data.forEach((point, index) => {
    const x = padding.left + stepX * index
    const y = padding.top + chartH - (point.score / maxVal) * chartH
    ctx.beginPath()
    ctx.arc(x, y, 4, 0, Math.PI * 2)
    ctx.setFillStyle('#67c23a')
    ctx.fill()

    // Score label
    ctx.setFillStyle('#333')
    ctx.setFontSize(10)
    ctx.setTextAlign('center')
    ctx.fillText(String(point.score), x, y - 10)

    // Date label
    ctx.setFillStyle('#999')
    ctx.fillText(point.date || point.examName, x, height - 10)
  })

  ctx.draw()
}
</script>

<style scoped>
.report-page {
  min-height: 100vh;
  background-color: #f5f5f5;
  padding: 24rpx;
  padding-bottom: 60rpx;
}

.period-tabs {
  display: flex;
  flex-direction: row;
  background-color: #fff;
  border-radius: 12rpx;
  padding: 6rpx;
  margin-bottom: 24rpx;
}

.period-tab {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 18rpx 0;
  border-radius: 10rpx;
}

.period-tab.active {
  background-color: #4a90d9;
}

.period-tab-text {
  font-size: 28rpx;
  color: #666;
}

.period-tab.active .period-tab-text {
  color: #fff;
  font-weight: 600;
}

.summary-card {
  background-color: #fff;
  border-radius: 16rpx;
  padding: 36rpx 30rpx;
  margin-bottom: 24rpx;
  display: flex;
  flex-direction: row;
  align-items: center;
}

.summary-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.summary-value {
  font-size: 48rpx;
  color: #4a90d9;
  font-weight: 700;
}

.summary-label {
  font-size: 22rpx;
  color: #999;
  margin-top: 8rpx;
}

.summary-divider {
  width: 1rpx;
  height: 80rpx;
  background-color: #eee;
  margin: 0 30rpx;
}

.summary-sub {
  flex: 1;
  display: flex;
  flex-direction: row;
  gap: 24rpx;
}

.summary-sub-item {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.sub-value {
  font-size: 36rpx;
  color: #333;
  font-weight: 600;
}

.sub-label {
  font-size: 20rpx;
  color: #999;
  margin-top: 6rpx;
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

.chart-canvas {
  width: 100%;
  height: 100%;
}

.chart-empty {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 160rpx;
}

.chart-empty-text {
  font-size: 26rpx;
  color: #ccc;
}

.stats-card {
  background-color: #fff;
  border-radius: 16rpx;
  padding: 28rpx;
  margin-bottom: 24rpx;
}

.stats-grid {
  display: flex;
  flex-direction: row;
  flex-wrap: wrap;
  gap: 16rpx;
}

.stats-item {
  width: calc(50% - 8rpx);
  background-color: #f8fafd;
  border-radius: 12rpx;
  padding: 24rpx 20rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.stats-number {
  font-size: 36rpx;
  color: #4a90d9;
  font-weight: 700;
}

.stats-desc {
  font-size: 22rpx;
  color: #999;
  margin-top: 8rpx;
}

.analysis-card {
  background-color: #fff;
  border-radius: 16rpx;
  padding: 28rpx;
  margin-bottom: 24rpx;
}

.analysis-section {
  margin-bottom: 28rpx;
}

.analysis-section:last-child {
  margin-bottom: 0;
}

.analysis-header {
  display: flex;
  flex-direction: row;
  align-items: center;
  margin-bottom: 16rpx;
}

.dot {
  width: 16rpx;
  height: 16rpx;
  border-radius: 8rpx;
  margin-right: 12rpx;
}

.dot.green {
  background-color: #67c23a;
}

.dot.red {
  background-color: #e74c3c;
}

.analysis-label {
  font-size: 28rpx;
  color: #333;
  font-weight: 600;
}

.analysis-item {
  display: flex;
  flex-direction: row;
  align-items: center;
  margin-bottom: 16rpx;
}

.analysis-item:last-child {
  margin-bottom: 0;
}

.item-name {
  width: 200rpx;
  font-size: 26rpx;
  color: #666;
  flex-shrink: 0;
}

.item-bar-wrap {
  flex: 1;
  height: 20rpx;
  background-color: #f0f0f0;
  border-radius: 10rpx;
  overflow: hidden;
  margin: 0 16rpx;
}

.item-bar {
  height: 100%;
  border-radius: 10rpx;
  transition: width 0.3s ease;
}

.green-bar {
  background: linear-gradient(90deg, #67c23a, #95d475);
}

.red-bar {
  background: linear-gradient(90deg, #e74c3c, #f89898);
}

.item-value {
  width: 80rpx;
  font-size: 26rpx;
  color: #333;
  font-weight: 600;
  text-align: right;
  flex-shrink: 0;
}

.tips-card {
  background-color: #fff;
  border-radius: 16rpx;
  padding: 28rpx;
  margin-bottom: 24rpx;
}

.tip-item {
  display: flex;
  flex-direction: row;
  align-items: flex-start;
  margin-bottom: 16rpx;
}

.tip-item:last-child {
  margin-bottom: 0;
}

.tip-number {
  width: 40rpx;
  height: 40rpx;
  background-color: #4a90d9;
  border-radius: 20rpx;
  color: #fff;
  font-size: 22rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  margin-right: 16rpx;
  text-align: center;
  line-height: 40rpx;
}

.tip-text {
  font-size: 26rpx;
  color: #666;
  line-height: 1.6;
  flex: 1;
}
</style>
