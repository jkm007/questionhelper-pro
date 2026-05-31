<template>
  <view class="exam-history-page">
    <!-- Header -->
    <view class="header">
      <view class="header-nav">
        <view class="nav-back" @tap="handleBack">
          <text class="nav-back-icon">←</text>
        </view>
        <text class="nav-title">Exam History</text>
        <view class="nav-placeholder"></view>
      </view>
    </view>

    <!-- History List -->
    <scroll-view
      scroll-y
      class="history-scroll"
      refresher-enabled
      :refresher-triggered="isRefreshing"
      @refresherrefresh="onRefresh"
      @scrolltolower="onLoadMore"
    >
      <view v-if="historyList.length === 0 && !loading" class="empty-state">
        <text class="empty-icon">📋</text>
        <text class="empty-text">No exam history</text>
        <text class="empty-hint">Your completed exams will appear here</text>
      </view>

      <view v-else class="history-cards">
        <view
          v-for="item in historyList"
          :key="item.id"
          class="history-card"
          @tap="goToResult(item.id)"
        >
          <view class="card-header">
            <text class="card-title">{{ item.examTitle }}</text>
            <view class="status-badge" :class="item.passed ? 'badge-pass' : 'badge-fail'">
              <text class="status-text">{{ item.passed ? 'Passed' : 'Failed' }}</text>
            </view>
          </view>

          <view class="card-body">
            <view class="score-display">
              <text class="score-value">{{ item.score }}</text>
              <text class="score-divider">/</text>
              <text class="score-total">{{ item.totalScore }}</text>
            </view>
            <view class="score-bar">
              <view
                class="score-bar-fill"
                :class="item.passed ? 'fill-pass' : 'fill-fail'"
                :style="{ width: (item.score / item.totalScore * 100) + '%' }"
              ></view>
            </view>
          </view>

          <view class="card-footer">
            <view class="footer-item">
              <text class="footer-icon">📅</text>
              <text class="footer-text">{{ formatTime(item.submitTime) }}</text>
            </view>
            <view class="footer-item">
              <text class="footer-icon">⏱</text>
              <text class="footer-text">{{ formatDuration(item.duration) }}</text>
            </view>
            <view class="footer-item">
              <text class="footer-icon">📊</text>
              <text class="footer-text">Attempt {{ item.attemptNumber }}</text>
            </view>
          </view>
        </view>
      </view>

      <!-- Load More -->
      <view v-if="hasMore && historyList.length > 0" class="load-more">
        <text class="load-more-text">{{ loading ? 'Loading...' : 'Load More' }}</text>
      </view>
      <view v-if="!hasMore && historyList.length > 0" class="no-more">
        <text class="no-more-text">No more history</text>
      </view>
    </scroll-view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { getExamHistoryList } from '@/api/exam';

interface ExamHistory {
  id: string;
  examId: string;
  examTitle: string;
  score: number;
  totalScore: number;
  passed: boolean;
  submitTime: string;
  duration: number;
  attemptNumber: number;
}

const historyList = ref<ExamHistory[]>([]);
const loading = ref(false);
const isRefreshing = ref(false);
const hasMore = ref(true);
const page = ref(1);
const pageSize = 10;

onMounted(() => {
  loadHistoryList();
});

async function loadHistoryList(isRefresh = false) {
  if (loading.value) return;
  loading.value = true;

  try {
    if (isRefresh) {
      page.value = 1;
      hasMore.value = true;
    }

    const res = await getExamHistoryList({
      page: page.value,
      pageSize,
    });

    if (res.code === 0) {
      const newList = res.data.list;
      if (isRefresh) {
        historyList.value = newList;
      } else {
        historyList.value.push(...newList);
      }
      hasMore.value = newList.length >= pageSize;
      page.value++;
    }
  } catch (e) {
    console.error('Failed to load history:', e);
    uni.showToast({ title: 'Failed to load', icon: 'none' });
  } finally {
    loading.value = false;
    isRefreshing.value = false;
  }
}

function onRefresh() {
  isRefreshing.value = true;
  loadHistoryList(true);
}

function onLoadMore() {
  if (!hasMore.value || loading.value) return;
  loadHistoryList();
}

function formatTime(time: string): string {
  if (!time) return '';
  const date = new Date(time);
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const day = String(date.getDate()).padStart(2, '0');
  const hour = String(date.getHours()).padStart(2, '0');
  const minute = String(date.getMinutes()).padStart(2, '0');
  return `${year}-${month}-${day} ${hour}:${minute}`;
}

function formatDuration(seconds: number): string {
  const mins = Math.floor(seconds / 60);
  const secs = seconds % 60;
  return `${mins}m${secs}s`;
}

function goToResult(attemptId: string) {
  uni.navigateTo({ url: `/pages/exam/result?attemptId=${attemptId}` });
}

function handleBack() {
  uni.navigateBack();
}
</script>

<style scoped>
.exam-history-page {
  min-height: 100vh;
  background-color: #f5f5f5;
  display: flex;
  flex-direction: column;
}

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

.history-scroll {
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

.history-cards {
  padding: 20rpx 30rpx;
  display: flex;
  flex-direction: column;
  gap: 20rpx;
}

.history-card {
  background: #ffffff;
  border-radius: 20rpx;
  padding: 30rpx;
  box-shadow: 0 4rpx 16rpx rgba(0, 0, 0, 0.06);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 20rpx;
}

.card-title {
  font-size: 32rpx;
  font-weight: 600;
  color: #333;
  flex: 1;
  margin-right: 16rpx;
}

.status-badge {
  padding: 8rpx 20rpx;
  border-radius: 24rpx;
  flex-shrink: 0;
}

.badge-pass {
  background: #e8f5e9;
}

.badge-fail {
  background: #ffebee;
}

.status-text {
  font-size: 24rpx;
  font-weight: 600;
}

.badge-pass .status-text {
  color: #388e3c;
}

.badge-fail .status-text {
  color: #d32f2f;
}

.card-body {
  margin-bottom: 20rpx;
}

.score-display {
  display: flex;
  align-items: baseline;
  margin-bottom: 12rpx;
}

.score-value {
  font-size: 48rpx;
  font-weight: bold;
  color: #667eea;
}

.score-divider {
  font-size: 32rpx;
  color: #999;
  margin: 0 8rpx;
}

.score-total {
  font-size: 32rpx;
  color: #999;
}

.score-bar {
  height: 12rpx;
  background: #e0e0e0;
  border-radius: 6rpx;
  overflow: hidden;
}

.score-bar-fill {
  height: 100%;
  border-radius: 6rpx;
  transition: width 0.5s ease;
}

.fill-pass {
  background: linear-gradient(90deg, #4caf50, #81c784);
}

.fill-fail {
  background: linear-gradient(90deg, #f44336, #e57373);
}

.card-footer {
  display: flex;
  justify-content: space-between;
  padding-top: 20rpx;
  border-top: 1rpx solid #f0f0f0;
}

.footer-item {
  display: flex;
  align-items: center;
}

.footer-icon {
  font-size: 24rpx;
  margin-right: 8rpx;
}

.footer-text {
  font-size: 24rpx;
  color: #999;
}

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
