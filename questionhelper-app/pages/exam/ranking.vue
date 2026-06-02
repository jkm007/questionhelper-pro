<template>
  <view class="exam-ranking-page">
    <!-- Header -->
    <view class="header">
      <view class="header-nav">
        <view class="nav-back" @tap="handleBack">
          <text class="nav-back-icon">←</text>
        </view>
        <text class="nav-title">成绩排名</text>
        <view class="nav-placeholder"></view>
      </view>
    </view>

    <!-- Top 3 Podium -->
    <view v-if="rankingList.length >= 3" class="podium-section">
      <view class="podium-row">
        <!-- 2nd Place -->
        <view class="podium-item podium-silver" @tap="onPodiumTap(rankingList[1])">
          <image
            class="podium-avatar"
            :src="rankingList[1].avatar || '/static/default-avatar.png'"
            mode="aspectFill"
          />
          <view class="podium-medal medal-silver">
            <text class="medal-text">2</text>
          </view>
          <text class="podium-name">{{ rankingList[1].nickname }}</text>
          <text class="podium-score">{{ rankingList[1].score }}分</text>
        </view>

        <!-- 1st Place -->
        <view class="podium-item podium-gold" @tap="onPodiumTap(rankingList[0])">
          <view class="podium-crown">
            <text class="crown-icon">👑</text>
          </view>
          <image
            class="podium-avatar podium-avatar-lg"
            :src="rankingList[0].avatar || '/static/default-avatar.png'"
            mode="aspectFill"
          />
          <view class="podium-medal medal-gold">
            <text class="medal-text">1</text>
          </view>
          <text class="podium-name">{{ rankingList[0].nickname }}</text>
          <text class="podium-score">{{ rankingList[0].score }}分</text>
        </view>

        <!-- 3rd Place -->
        <view class="podium-item podium-bronze" @tap="onPodiumTap(rankingList[2])">
          <image
            class="podium-avatar"
            :src="rankingList[2].avatar || '/static/default-avatar.png'"
            mode="aspectFill"
          />
          <view class="podium-medal medal-bronze">
            <text class="medal-text">3</text>
          </view>
          <text class="podium-name">{{ rankingList[2].nickname }}</text>
          <text class="podium-score">{{ rankingList[2].score }}分</text>
        </view>
      </view>
    </view>

    <!-- My Rank Card -->
    <view v-if="myRank" class="my-rank-card">
      <view class="my-rank-left">
        <view class="my-rank-badge">
          <text class="my-rank-num">{{ myRank.rank }}</text>
        </view>
        <image
          class="my-rank-avatar"
          :src="myRank.avatar || '/static/default-avatar.png'"
          mode="aspectFill"
        />
        <view class="my-rank-info">
          <text class="my-rank-name">{{ myRank.nickname }}</text>
          <text class="my-rank-detail">正确率 {{ myRank.correctRate }}%</text>
        </view>
      </view>
      <view class="my-rank-right">
        <text class="my-rank-score">{{ myRank.score }}</text>
        <text class="my-rank-unit">分</text>
      </view>
    </view>

    <!-- Ranking List -->
    <scroll-view
      scroll-y
      class="ranking-scroll"
      refresher-enabled
      :refresher-triggered="isRefreshing"
      @refresherrefresh="onRefresh"
      @scrolltolower="onLoadMore"
    >
      <view v-if="rankingList.length === 0 && !loading" class="empty-state">
        <text class="empty-icon">📊</text>
        <text class="empty-text">暂无排名数据</text>
        <text class="empty-hint">考试结束后将自动生成排名</text>
      </view>

      <view v-else class="ranking-list">
        <view
          v-for="(item, index) in displayList"
          :key="item.userId"
          class="ranking-item"
          :class="{ 'ranking-item-me': item.isMe }"
        >
          <view class="ranking-left">
            <view
              class="ranking-num"
              :class="{
                'num-gold': index + 4 === 1,
                'num-silver': index + 4 === 2,
                'num-bronze': index + 4 === 3,
              }"
            >
              <text class="ranking-num-text">{{ index + 4 }}</text>
            </view>
            <image
              class="ranking-avatar"
              :src="item.avatar || '/static/default-avatar.png'"
              mode="aspectFill"
            />
            <view class="ranking-info">
              <view class="ranking-name-row">
                <text class="ranking-name">{{ item.nickname }}</text>
                <view v-if="item.isMe" class="me-tag">
                  <text class="me-tag-text">我</text>
                </view>
              </view>
              <text class="ranking-detail">
                用时 {{ formatDuration(item.duration) }} · 正确率 {{ item.correctRate }}%
              </text>
            </view>
          </view>
          <view class="ranking-right">
            <text class="ranking-score">{{ item.score }}</text>
            <text class="ranking-score-unit">分</text>
          </view>
        </view>
      </view>

      <!-- Load More -->
      <view v-if="hasMore && rankingList.length > 0" class="load-more">
        <text class="load-more-text">{{ loading ? '加载中...' : '加载更多' }}</text>
      </view>
      <view v-if="!hasMore && rankingList.length > 0" class="no-more">
        <text class="no-more-text">没有更多了</text>
      </view>
    </scroll-view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { request } from '@/api/request'

interface RankingItem {
  userId: string
  nickname: string
  avatar: string
  score: number
  totalScore: number
  rank: number
  duration: number
  correctRate: number
  isMe: boolean
}

const examId = ref('')
const rankingList = ref<RankingItem[]>([])
const loading = ref(false)
const isRefreshing = ref(false)
const hasMore = ref(true)
const page = ref(1)
const pageSize = 20

const myRank = computed(() => {
  return rankingList.value.find((item) => item.isMe) || null
})

// Items beyond the top 3
const displayList = computed(() => {
  return rankingList.value.slice(3)
})

onMounted(() => {
  const pages = getCurrentPages()
  const currentPage = pages[pages.length - 1] as any
  examId.value = currentPage.options?.examId || ''

  if (!examId.value) {
    uni.showToast({ title: '参数错误', icon: 'none' })
    setTimeout(() => uni.navigateBack(), 1500)
    return
  }

  loadRanking()
})

async function loadRanking(isRefresh = false) {
  if (loading.value) return
  loading.value = true

  try {
    if (isRefresh) {
      page.value = 1
      hasMore.value = true
    }

    const res = await request({
      url: `/exam/${examId.value}/ranking`,
      data: { page: page.value, pageSize },
    })

    if (res.code === '00000') {
      const newList = res.data.list || res.data
      if (isRefresh) {
        rankingList.value = newList
      } else {
        rankingList.value.push(...newList)
      }
      hasMore.value = newList.length >= pageSize
      page.value++
    }
  } catch (e) {
    console.error('Failed to load ranking:', e)
    uni.showToast({ title: '加载失败', icon: 'none' })
  } finally {
    loading.value = false
    isRefreshing.value = false
  }
}

function onRefresh() {
  isRefreshing.value = true
  loadRanking(true)
}

function onLoadMore() {
  if (!hasMore.value || loading.value) return
  loadRanking()
}

function formatDuration(seconds: number): string {
  const mins = Math.floor(seconds / 60)
  const secs = seconds % 60
  if (mins > 0) {
    return `${mins}分${secs}秒`
  }
  return `${secs}秒`
}

function onPodiumTap(item: RankingItem) {
  // Could navigate to user profile in the future
}

function handleBack() {
  uni.navigateBack()
}
</script>

<style scoped>
.exam-ranking-page {
  min-height: 100vh;
  background-color: #f5f5f5;
  display: flex;
  flex-direction: column;
}

/* ---- Header ---- */

.header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding-top: var(--status-bar-height, 44rpx);
}

.header-nav {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 88rpx;
  padding: 0 30rpx;
}

.nav-back {
  width: 60rpx;
  height: 60rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.nav-back-icon {
  font-size: 36rpx;
  color: #ffffff;
}

.nav-title {
  font-size: 32rpx;
  font-weight: 600;
  color: #ffffff;
}

.nav-placeholder {
  width: 60rpx;
}

/* ---- Podium (Top 3) ---- */

.podium-section {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 20rpx 30rpx 60rpx;
}

.podium-row {
  display: flex;
  align-items: flex-end;
  justify-content: center;
  gap: 20rpx;
}

.podium-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  flex: 1;
}

.podium-gold {
  order: 1;
}

.podium-silver {
  order: 0;
}

.podium-bronze {
  order: 2;
}

.podium-crown {
  margin-bottom: 8rpx;
}

.crown-icon {
  font-size: 40rpx;
}

.podium-avatar {
  width: 88rpx;
  height: 88rpx;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.3);
  margin-bottom: 12rpx;
  border: 4rpx solid rgba(255, 255, 255, 0.6);
}

.podium-avatar-lg {
  width: 108rpx;
  height: 108rpx;
  border-width: 6rpx;
}

.podium-medal {
  width: 44rpx;
  height: 44rpx;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 12rpx;
}

.medal-gold {
  background: linear-gradient(135deg, #ffd700, #ffb300);
  box-shadow: 0 4rpx 12rpx rgba(255, 215, 0, 0.5);
}

.medal-silver {
  background: linear-gradient(135deg, #e0e0e0, #bdbdbd);
  box-shadow: 0 4rpx 12rpx rgba(189, 189, 189, 0.5);
}

.medal-bronze {
  background: linear-gradient(135deg, #cd7f32, #b8690e);
  box-shadow: 0 4rpx 12rpx rgba(205, 127, 50, 0.5);
}

.medal-text {
  font-size: 24rpx;
  font-weight: bold;
  color: #ffffff;
}

.podium-name {
  font-size: 24rpx;
  color: rgba(255, 255, 255, 0.9);
  max-width: 160rpx;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  margin-bottom: 4rpx;
}

.podium-score {
  font-size: 28rpx;
  font-weight: bold;
  color: #ffffff;
}

/* ---- My Rank Card ---- */

.my-rank-card {
  margin: -30rpx 30rpx 20rpx;
  background: #ffffff;
  border-radius: 20rpx;
  padding: 24rpx 30rpx;
  box-shadow: 0 4rpx 20rpx rgba(0, 0, 0, 0.08);
  display: flex;
  align-items: center;
  justify-content: space-between;
  position: relative;
  z-index: 1;
  border: 2rpx solid #667eea;
}

.my-rank-left {
  display: flex;
  align-items: center;
}

.my-rank-badge {
  width: 52rpx;
  height: 52rpx;
  border-radius: 50%;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 16rpx;
}

.my-rank-num {
  font-size: 26rpx;
  font-weight: bold;
  color: #ffffff;
}

.my-rank-avatar {
  width: 72rpx;
  height: 72rpx;
  border-radius: 50%;
  margin-right: 16rpx;
  background: #f0f0f0;
}

.my-rank-info {
  display: flex;
  flex-direction: column;
}

.my-rank-name {
  font-size: 28rpx;
  font-weight: 600;
  color: #333;
}

.my-rank-detail {
  font-size: 22rpx;
  color: #999;
  margin-top: 4rpx;
}

.my-rank-right {
  display: flex;
  align-items: baseline;
}

.my-rank-score {
  font-size: 40rpx;
  font-weight: bold;
  color: #667eea;
}

.my-rank-unit {
  font-size: 24rpx;
  color: #999;
  margin-left: 4rpx;
}

/* ---- Ranking List ---- */

.ranking-scroll {
  flex: 1;
  height: 0;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 120rpx 40rpx;
}

.empty-icon {
  font-size: 80rpx;
  margin-bottom: 24rpx;
}

.empty-text {
  font-size: 30rpx;
  color: #333;
  font-weight: 600;
  margin-bottom: 12rpx;
}

.empty-hint {
  font-size: 26rpx;
  color: #999;
}

.ranking-list {
  padding: 0 30rpx 20rpx;
  display: flex;
  flex-direction: column;
  gap: 16rpx;
}

.ranking-item {
  background: #ffffff;
  border-radius: 16rpx;
  padding: 24rpx;
  display: flex;
  align-items: center;
  justify-content: space-between;
  box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.04);
}

.ranking-item-me {
  border: 2rpx solid #667eea;
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.05) 0%, rgba(118, 75, 162, 0.05) 100%);
}

.ranking-left {
  display: flex;
  align-items: center;
  flex: 1;
  min-width: 0;
}

.ranking-num {
  width: 48rpx;
  height: 48rpx;
  border-radius: 50%;
  background: #f0f0f0;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 16rpx;
  flex-shrink: 0;
}

.ranking-num-text {
  font-size: 24rpx;
  font-weight: 600;
  color: #666;
}

.ranking-avatar {
  width: 68rpx;
  height: 68rpx;
  border-radius: 50%;
  margin-right: 16rpx;
  background: #f0f0f0;
  flex-shrink: 0;
}

.ranking-info {
  flex: 1;
  min-width: 0;
}

.ranking-name-row {
  display: flex;
  align-items: center;
}

.ranking-name {
  font-size: 28rpx;
  font-weight: 500;
  color: #333;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.me-tag {
  margin-left: 8rpx;
  padding: 2rpx 10rpx;
  background: #667eea;
  border-radius: 8rpx;
}

.me-tag-text {
  font-size: 20rpx;
  color: #ffffff;
}

.ranking-detail {
  font-size: 22rpx;
  color: #999;
  margin-top: 4rpx;
  display: block;
}

.ranking-right {
  display: flex;
  align-items: baseline;
  margin-left: 16rpx;
  flex-shrink: 0;
}

.ranking-score {
  font-size: 36rpx;
  font-weight: bold;
  color: #333;
}

.ranking-score-unit {
  font-size: 22rpx;
  color: #999;
  margin-left: 4rpx;
}

/* ---- Load More ---- */

.load-more {
  text-align: center;
  padding: 30rpx;
}

.load-more-text {
  font-size: 26rpx;
  color: #999;
}

.no-more {
  text-align: center;
  padding: 30rpx;
}

.no-more-text {
  font-size: 26rpx;
  color: #ccc;
}
</style>
