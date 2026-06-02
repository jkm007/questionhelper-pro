<template>
  <view class="ranking-page">
    <!-- Top 3 Podium -->
    <view v-if="topThree.length" class="podium">
      <view v-for="(item, i) in topThree" :key="item.rank" class="podium-item" :class="'rank-' + (i + 1)">
        <text class="podium-medal">{{ medals[i] }}</text>
        <text class="podium-name">{{ item.name }}</text>
        <text class="podium-score">{{ item.score }}分</text>
        <text class="podium-time">{{ item.time }}</text>
      </view>
    </view>

    <!-- My Rank -->
    <view v-if="myRank" class="my-rank-card">
      <text class="my-rank-label">我的排名</text>
      <view class="my-rank-info">
        <text class="my-rank-pos">第 {{ myRank.rank }} 名</text>
        <text class="my-rank-score">{{ myRank.score }}分</text>
      </view>
    </view>

    <!-- Full Ranking List -->
    <view class="section">
      <text class="section-title">完整排名</text>
      <view v-for="(item, i) in rankings" :key="i" class="rank-item" :class="{ highlight: item.isMe }">
        <text class="rank-pos">{{ item.rank }}</text>
        <view class="rank-info">
          <text class="rank-name">{{ item.name }}{{ item.isMe ? '（我）' : '' }}</text>
          <text class="rank-time">{{ item.time }}</text>
        </view>
        <text class="rank-score">{{ item.score }}分</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'

const medals = ['🥇', '🥈', '🥉']

interface RankItem {
  rank: number
  name: string
  score: number
  time: string
  isMe: boolean
}

const rankings = ref<RankItem[]>([])
const examId = ref('')

const topThree = computed(() => rankings.value.filter((r) => r.rank <= 3).sort((a, b) => a.rank - b.rank))
const myRank = computed(() => rankings.value.find((r) => r.isMe))

onMounted(() => {
  const pages = getCurrentPages()
  const page = pages[pages.length - 1] as any
  examId.value = page.options?.examId || ''
  fetchRankings()
})

async function fetchRankings() {
  try {
    // TODO: replace with actual API call using examId
    rankings.value = [
      { rank: 1, name: '王小明', score: 98, time: '45:12', isMe: false },
      { rank: 2, name: '刘思琪', score: 95, time: '38:45', isMe: false },
      { rank: 3, name: '陈浩然', score: 93, time: '52:30', isMe: false },
      { rank: 4, name: '我', score: 88, time: '40:15', isMe: true },
      { rank: 5, name: '赵丽华', score: 85, time: '55:00', isMe: false },
      { rank: 6, name: '孙志强', score: 82, time: '48:20', isMe: false }
    ]
  } catch (e) {
    console.error('Failed to load rankings', e)
  }
}
</script>

<style scoped>
.ranking-page {
  min-height: 100vh;
  background-color: #f5f5f5;
  padding: 20rpx 24rpx;
  padding-bottom: 80rpx;
}

.podium {
  display: flex;
  flex-direction: row;
  align-items: flex-end;
  justify-content: center;
  gap: 20rpx;
  margin-bottom: 24rpx;
  padding: 30rpx 0 0;
}

.podium-item {
  flex: 1;
  background-color: #fff;
  border-radius: 16rpx;
  padding: 24rpx 16rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
}

.podium-item.rank-1 {
  padding-bottom: 40rpx;
  border: 2rpx solid #ffd700;
}

.podium-item.rank-2 {
  border: 2rpx solid #c0c0c0;
}

.podium-item.rank-3 {
  border: 2rpx solid #cd7f32;
}

.podium-medal {
  font-size: 52rpx;
  margin-bottom: 10rpx;
}

.podium-name {
  font-size: 26rpx;
  font-weight: 600;
  color: #333;
  margin-bottom: 6rpx;
}

.podium-score {
  font-size: 28rpx;
  font-weight: 700;
  color: #4a90d9;
}

.podium-time {
  font-size: 20rpx;
  color: #bbb;
  margin-top: 4rpx;
}

.my-rank-card {
  background: linear-gradient(135deg, #4a90d9, #6db3f2);
  border-radius: 16rpx;
  padding: 28rpx 30rpx;
  margin-bottom: 24rpx;
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: space-between;
}

.my-rank-label {
  font-size: 28rpx;
  color: rgba(255, 255, 255, 0.85);
}

.my-rank-info {
  display: flex;
  flex-direction: row;
  align-items: center;
  gap: 20rpx;
}

.my-rank-pos {
  font-size: 32rpx;
  font-weight: 700;
  color: #fff;
}

.my-rank-score {
  font-size: 28rpx;
  color: rgba(255, 255, 255, 0.9);
}

.section {
  background-color: #fff;
  border-radius: 16rpx;
  padding: 28rpx 30rpx;
}

.section-title {
  font-size: 30rpx;
  font-weight: 600;
  color: #333;
  margin-bottom: 20rpx;
}

.rank-item {
  display: flex;
  flex-direction: row;
  align-items: center;
  padding: 22rpx 0;
  border-bottom: 1rpx solid #f5f5f5;
}

.rank-item:last-child {
  border-bottom: none;
}

.rank-item.highlight {
  background-color: #f0f7ff;
  margin: 0 -30rpx;
  padding: 22rpx 30rpx;
  border-radius: 12rpx;
}

.rank-pos {
  width: 60rpx;
  font-size: 30rpx;
  font-weight: 700;
  color: #333;
  text-align: center;
}

.rank-info {
  flex: 1;
  margin-left: 16rpx;
}

.rank-name {
  font-size: 28rpx;
  color: #333;
}

.rank-time {
  font-size: 22rpx;
  color: #bbb;
}

.rank-score {
  font-size: 28rpx;
  font-weight: 600;
  color: #4a90d9;
}
</style>
