<template>
  <view class="page">
    <!-- Summary Cards -->
    <view class="summary-row">
      <view class="summary-card">
        <text class="summary-number">{{ summary.totalWrong }}</text>
        <text class="summary-label">总错题数</text>
      </view>
      <view class="summary-card">
        <text class="summary-number mastered-num">{{ summary.masteredCount }}</text>
        <text class="summary-label">已掌握</text>
      </view>
      <view class="summary-card">
        <text class="summary-number rate-num">{{ summary.masteryRate }}%</text>
        <text class="summary-label">掌握率</text>
      </view>
    </view>

    <!-- Category Distribution -->
    <view class="chart-section">
      <text class="chart-title">错题分类分布</text>
      <view class="chart-container">
        <view v-for="(item, idx) in categoryData" :key="idx" class="bar-item">
          <view class="bar-label-wrap">
            <text class="bar-label">{{ item.category }}</text>
            <text class="bar-value">{{ item.count }}题</text>
          </view>
          <view class="bar-bg">
            <view
              class="bar-fill"
              :style="{ width: (item.count / maxCategoryCount * 100) + '%' }"
              :class="'bar-color-' + (idx % 5)"
            ></view>
          </view>
        </view>
      </view>
      <view v-if="categoryData.length === 0" class="chart-empty">
        <text class="chart-empty-text">暂无数据</text>
      </view>
    </view>

    <!-- Knowledge Weak Spots -->
    <view class="chart-section">
      <text class="chart-title">薄弱知识点 TOP10</text>
      <view class="weak-list">
        <view v-for="(item, idx) in weakPoints" :key="idx" class="weak-item">
          <view class="weak-rank" :class="'rank-' + idx">
            <text class="rank-num">{{ idx + 1 }}</text>
          </view>
          <view class="weak-content">
            <text class="weak-name">{{ item.name }}</text>
            <view class="weak-bar-bg">
              <view
                class="weak-bar-fill"
                :style="{ width: (item.wrongRate) + '%' }"
              ></view>
            </view>
          </view>
          <text class="weak-rate">{{ item.wrongRate }}%</text>
        </view>
      </view>
      <view v-if="weakPoints.length === 0" class="chart-empty">
        <text class="chart-empty-text">暂无数据</text>
      </view>
    </view>

    <!-- Mastery Trend -->
    <view class="chart-section">
      <text class="chart-title">掌握度趋势 (近30天)</text>
      <view class="trend-chart">
        <view class="trend-y-axis">
          <text class="y-label">100%</text>
          <text class="y-label">75%</text>
          <text class="y-label">50%</text>
          <text class="y-label">25%</text>
          <text class="y-label">0%</text>
        </view>
        <scroll-view class="trend-scroll" scroll-x>
          <view class="trend-bars" :style="{ width: trendData.length * 48 + 'rpx' }">
            <view v-for="(item, idx) in trendData" :key="idx" class="trend-bar-wrap">
              <view class="trend-bar-bg">
                <view
                  class="trend-bar-fill"
                  :style="{ height: item.rate + '%' }"
                ></view>
              </view>
              <text class="trend-date">{{ item.dateLabel }}</text>
            </view>
          </view>
        </scroll-view>
      </view>
      <view v-if="trendData.length === 0" class="chart-empty">
        <text class="chart-empty-text">暂无数据</text>
      </view>
    </view>

    <!-- Recommendations -->
    <view class="recommend-section">
      <text class="chart-title">复习建议</text>
      <view v-for="(rec, idx) in recommendations" :key="idx" class="rec-card">
        <view class="rec-icon-wrap" :class="'rec-type-' + rec.type">
          <text class="rec-icon">{{ rec.icon }}</text>
        </view>
        <view class="rec-content">
          <text class="rec-title">{{ rec.title }}</text>
          <text class="rec-desc">{{ rec.description }}</text>
        </view>
        <text class="rec-action" @tap="onRecAction(rec)">&#xe61e;</text>
      </view>
      <view v-if="recommendations.length === 0" class="chart-empty">
        <text class="chart-empty-text">暂无建议</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { getWrongAnalysis, getMasteryTrend, getRecommendations } from '@/api/wrong'

interface CategoryItem {
  category: string
  count: number
}

interface WeakPoint {
  name: string
  wrongRate: number
}

interface TrendItem {
  dateLabel: string
  rate: number
}

interface Recommendation {
  type: string
  icon: string
  title: string
  description: string
  actionUrl: string
}

const summary = ref({
  totalWrong: 0,
  masteredCount: 0,
  masteryRate: 0
})

const categoryData = ref<CategoryItem[]>([])
const weakPoints = ref<WeakPoint[]>([])
const trendData = ref<TrendItem[]>([])
const recommendations = ref<Recommendation[]>([])

const maxCategoryCount = computed(() => {
  if (categoryData.value.length === 0) return 1
  return Math.max(...categoryData.value.map(item => item.count), 1)
})

async function loadData() {
  uni.showLoading({ title: '加载中' })
  try {
    const [analysisRes, trendRes, recRes] = await Promise.all([
      getWrongAnalysis(),
      getMasteryTrend(),
      getRecommendations()
    ])

    if (analysisRes.data) {
      summary.value = {
        totalWrong: analysisRes.data.totalWrong || 0,
        masteredCount: analysisRes.data.masteredCount || 0,
        masteryRate: analysisRes.data.masteryRate || 0
      }
      categoryData.value = analysisRes.data.categoryDistribution || []
      weakPoints.value = (analysisRes.data.weakPoints || []).slice(0, 10)
    }

    trendData.value = trendRes.data || []
    recommendations.value = recRes.data || []
  } catch (e) {
    uni.showToast({ title: '加载失败', icon: 'none' })
  } finally {
    uni.hideLoading()
  }
}

function onRecAction(rec: Recommendation) {
  if (rec.actionUrl) {
    uni.navigateTo({ url: rec.actionUrl })
  }
}

onMounted(() => {
  loadData()
})
</script>

<style scoped>
.page {
  min-height: 100vh;
  background-color: #f5f6fa;
  padding-bottom: 40rpx;
}

.summary-row {
  display: flex;
  gap: 16rpx;
  padding: 24rpx;
}

.summary-card {
  flex: 1;
  background-color: #fff;
  border-radius: 16rpx;
  padding: 28rpx 20rpx;
  text-align: center;
  box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.04);
}

.summary-number {
  font-size: 44rpx;
  font-weight: 700;
  color: #333;
  display: block;
  margin-bottom: 8rpx;
}

.mastered-num {
  color: #52c41a;
}

.rate-num {
  color: #1677ff;
}

.summary-label {
  font-size: 24rpx;
  color: #999;
}

.chart-section {
  background-color: #fff;
  margin: 0 24rpx 20rpx;
  border-radius: 16rpx;
  padding: 28rpx;
  box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.04);
}

.chart-title {
  font-size: 30rpx;
  font-weight: 600;
  color: #333;
  margin-bottom: 24rpx;
  display: block;
}

.chart-container {
  display: flex;
  flex-direction: column;
  gap: 20rpx;
}

.bar-item {
  display: flex;
  flex-direction: column;
  gap: 8rpx;
}

.bar-label-wrap {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.bar-label {
  font-size: 26rpx;
  color: #555;
}

.bar-value {
  font-size: 24rpx;
  color: #999;
}

.bar-bg {
  height: 24rpx;
  background-color: #f0f1f5;
  border-radius: 12rpx;
  overflow: hidden;
}

.bar-fill {
  height: 100%;
  border-radius: 12rpx;
  transition: width 0.5s ease;
}

.bar-color-0 { background-color: #1677ff; }
.bar-color-1 { background-color: #52c41a; }
.bar-color-2 { background-color: #faad14; }
.bar-color-3 { background-color: #722ed1; }
.bar-color-4 { background-color: #ff4d4f; }

.chart-empty {
  padding: 48rpx 0;
  text-align: center;
}

.chart-empty-text {
  font-size: 26rpx;
  color: #ccc;
}

.weak-list {
  display: flex;
  flex-direction: column;
  gap: 20rpx;
}

.weak-item {
  display: flex;
  align-items: center;
  gap: 16rpx;
}

.weak-rank {
  width: 44rpx;
  height: 44rpx;
  border-radius: 50%;
  background-color: #f0f1f5;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.weak-rank.rank-0 { background-color: #ff4d4f; }
.weak-rank.rank-1 { background-color: #fa8c16; }
.weak-rank.rank-2 { background-color: #faad14; }

.rank-num {
  font-size: 24rpx;
  font-weight: 600;
  color: #666;
}

.rank-0 .rank-num,
.rank-1 .rank-num,
.rank-2 .rank-num {
  color: #fff;
}

.weak-content {
  flex: 1;
  min-width: 0;
}

.weak-name {
  font-size: 26rpx;
  color: #333;
  margin-bottom: 8rpx;
  display: block;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.weak-bar-bg {
  height: 12rpx;
  background-color: #f0f1f5;
  border-radius: 6rpx;
  overflow: hidden;
}

.weak-bar-fill {
  height: 100%;
  background-color: #ff4d4f;
  border-radius: 6rpx;
  transition: width 0.5s ease;
}

.weak-rate {
  font-size: 24rpx;
  color: #ff4d4f;
  flex-shrink: 0;
  min-width: 72rpx;
  text-align: right;
}

.trend-chart {
  display: flex;
  gap: 12rpx;
  height: 320rpx;
}

.trend-y-axis {
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  flex-shrink: 0;
  width: 72rpx;
}

.y-label {
  font-size: 20rpx;
  color: #ccc;
  text-align: right;
}

.trend-scroll {
  flex: 1;
  overflow: hidden;
}

.trend-bars {
  display: flex;
  align-items: flex-end;
  height: 280rpx;
  gap: 4rpx;
  padding-bottom: 40rpx;
  position: relative;
}

.trend-bar-wrap {
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 40rpx;
  height: 100%;
  justify-content: flex-end;
}

.trend-bar-bg {
  width: 20rpx;
  height: 240rpx;
  background-color: #f0f1f5;
  border-radius: 10rpx 10rpx 0 0;
  position: relative;
  overflow: hidden;
  display: flex;
  align-items: flex-end;
}

.trend-bar-fill {
  width: 100%;
  background-color: #1677ff;
  border-radius: 10rpx 10rpx 0 0;
  transition: height 0.5s ease;
}

.trend-date {
  font-size: 18rpx;
  color: #ccc;
  margin-top: 8rpx;
  white-space: nowrap;
}

.recommend-section {
  background-color: #fff;
  margin: 0 24rpx 20rpx;
  border-radius: 16rpx;
  padding: 28rpx;
  box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.04);
}

.rec-card {
  display: flex;
  align-items: center;
  gap: 20rpx;
  padding: 20rpx 0;
  border-bottom: 1rpx solid #f5f5f5;
}

.rec-card:last-child {
  border-bottom: none;
}

.rec-icon-wrap {
  width: 72rpx;
  height: 72rpx;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.rec-type-urgent { background-color: #fff1f0; }
.rec-type-review { background-color: #e8f3ff; }
.rec-type-practice { background-color: #f6ffed; }

.rec-icon {
  font-size: 36rpx;
}

.rec-content {
  flex: 1;
  min-width: 0;
}

.rec-title {
  font-size: 28rpx;
  font-weight: 600;
  color: #333;
  display: block;
  margin-bottom: 6rpx;
}

.rec-desc {
  font-size: 24rpx;
  color: #999;
  line-height: 1.5;
}

.rec-action {
  font-size: 32rpx;
  color: #ccc;
  padding: 8rpx;
  flex-shrink: 0;
}
</style>
