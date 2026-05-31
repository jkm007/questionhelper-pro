<template>
  <view class="exam-list-page">
    <!-- Header -->
    <view class="header">
      <text class="header-title">Exams</text>
      <text class="header-subtitle">Manage your examinations</text>
    </view>

    <!-- Tabs -->
    <view class="tab-bar">
      <view
        v-for="tab in tabs"
        :key="tab.key"
        class="tab-item"
        :class="{ 'tab-active': activeTab === tab.key }"
        @tap="switchTab(tab.key)"
      >
        <text class="tab-text" :class="{ 'tab-text-active': activeTab === tab.key }">{{ tab.name }}</text>
        <view v-if="activeTab === tab.key" class="tab-indicator"></view>
      </view>
    </view>

    <!-- Exam List -->
    <scroll-view
      scroll-y
      class="exam-scroll"
      refresher-enabled
      :refresher-triggered="isRefreshing"
      @refresherrefresh="onRefresh"
      @scrolltolower="onLoadMore"
    >
      <view v-if="examList.length === 0 && !loading" class="empty-state">
        <text class="empty-icon">📝</text>
        <text class="empty-text">No exams found</text>
        <text class="empty-hint">{{ getEmptyHint() }}</text>
      </view>

      <view v-else class="exam-cards">
        <view
          v-for="exam in examList"
          :key="exam.id"
          class="exam-card"
          @tap="goToDetail(exam.id)"
        >
          <view class="card-header">
            <text class="card-title">{{ exam.title }}</text>
            <view class="status-badge" :class="getStatusClass(exam.status)">
              <text class="status-text">{{ getStatusText(exam.status) }}</text>
            </view>
          </view>

          <view class="card-info">
            <view class="info-row">
              <text class="info-icon">📅</text>
              <text class="info-text">{{ formatDateRange(exam.startTime, exam.endTime) }}</text>
            </view>
            <view class="info-row">
              <text class="info-icon">⏱</text>
              <text class="info-text">Duration: {{ exam.duration }} min</text>
            </view>
            <view class="info-row">
              <text class="info-icon">👥</text>
              <text class="info-text">{{ exam.participantCount }} participants</text>
            </view>
          </view>

          <view v-if="exam.status === 'upcoming'" class="card-footer">
            <text class="footer-text">Starts in {{ getTimeUntil(exam.startTime) }}</text>
          </view>
          <view v-else-if="exam.status === 'in_progress'" class="card-footer footer-active">
            <text class="footer-text-active">In Progress</text>
            <text class="footer-arrow">></text>
          </view>
          <view v-else-if="exam.status === 'completed'" class="card-footer">
            <view v-if="exam.myScore !== undefined" class="score-info">
              <text class="score-label">Your Score:</text>
              <text class="score-value">{{ exam.myScore }}/{{ exam.totalScore }}</text>
            </view>
            <text v-else class="footer-text">View Details</text>
          </view>
        </view>
      </view>

      <!-- Load More -->
      <view v-if="hasMore && examList.length > 0" class="load-more">
        <text class="load-more-text">{{ loading ? 'Loading...' : 'Load More' }}</text>
      </view>
      <view v-if="!hasMore && examList.length > 0" class="no-more">
        <text class="no-more-text">No more exams</text>
      </view>
    </scroll-view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { getExamList } from '@/api/exam';

interface Exam {
  id: string;
  title: string;
  description: string;
  startTime: string;
  endTime: string;
  duration: number;
  participantCount: number;
  status: 'upcoming' | 'in_progress' | 'completed';
  totalScore: number;
  myScore?: number;
}

const tabs = [
  { key: 'upcoming', name: 'Upcoming' },
  { key: 'in_progress', name: 'In Progress' },
  { key: 'completed', name: 'Completed' },
];

const activeTab = ref('upcoming');
const examList = ref<Exam[]>([]);
const loading = ref(false);
const isRefreshing = ref(false);
const hasMore = ref(true);
const page = ref(1);
const pageSize = 10;

onMounted(() => {
  loadExamList();
});

async function loadExamList(isRefresh = false) {
  if (loading.value) return;
  loading.value = true;

  try {
    if (isRefresh) {
      page.value = 1;
      hasMore.value = true;
    }

    const res = await getExamList({
      status: activeTab.value,
      page: page.value,
      pageSize,
    });

    if (res.code === 0) {
      const newExams = res.data.list;
      if (isRefresh) {
        examList.value = newExams;
      } else {
        examList.value.push(...newExams);
      }
      hasMore.value = newExams.length >= pageSize;
      page.value++;
    }
  } catch (e) {
    console.error('Failed to load exams:', e);
    uni.showToast({ title: 'Failed to load', icon: 'none' });
  } finally {
    loading.value = false;
    isRefreshing.value = false;
  }
}

function switchTab(tab: string) {
  if (activeTab.value === tab) return;
  activeTab.value = tab;
  examList.value = [];
  loadExamList(true);
}

function onRefresh() {
  isRefreshing.value = true;
  loadExamList(true);
}

function onLoadMore() {
  if (!hasMore.value || loading.value) return;
  loadExamList();
}

function getStatusClass(status: string): string {
  const map: Record<string, string> = {
    upcoming: 'badge-upcoming',
    in_progress: 'badge-active',
    completed: 'badge-completed',
  };
  return map[status] || '';
}

function getStatusText(status: string): string {
  const map: Record<string, string> = {
    upcoming: 'Upcoming',
    in_progress: 'In Progress',
    completed: 'Completed',
  };
  return map[status] || status;
}

function getEmptyHint(): string {
  const map: Record<string, string> = {
    upcoming: 'No upcoming exams',
    in_progress: 'No exams in progress',
    completed: 'No completed exams',
  };
  return map[activeTab.value] || '';
}

function formatDateRange(start: string, end: string): string {
  const startDate = new Date(start);
  const endDate = new Date(end);
  const startStr = `${startDate.getMonth() + 1}/${startDate.getDate()} ${String(startDate.getHours()).padStart(2, '0')}:${String(startDate.getMinutes()).padStart(2, '0')}`;
  const endStr = `${endDate.getMonth() + 1}/${endDate.getDate()} ${String(endDate.getHours()).padStart(2, '0')}:${String(endDate.getMinutes()).padStart(2, '0')}`;
  return `${startStr} - ${endStr}`;
}

function getTimeUntil(time: string): string {
  const target = new Date(time).getTime();
  const now = Date.now();
  const diff = target - now;
  if (diff <= 0) return 'soon';
  const days = Math.floor(diff / (1000 * 60 * 60 * 24));
  const hours = Math.floor((diff % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60));
  if (days > 0) return `${days}d ${hours}h`;
  const minutes = Math.floor((diff % (1000 * 60 * 60)) / (1000 * 60));
  return `${hours}h ${minutes}m`;
}

function goToDetail(id: string) {
  uni.navigateTo({ url: `/pages/exam/detail?id=${id}` });
}
</script>

<style scoped>
.exam-list-page {
  min-height: 100vh;
  background-color: #f5f5f5;
  display: flex;
  flex-direction: column;
}

.header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 60rpx 40rpx 40rpx;
}

.header-title {
  font-size: 44rpx;
  font-weight: bold;
  color: #ffffff;
  display: block;
}

.header-subtitle {
  font-size: 26rpx;
  color: rgba(255, 255, 255, 0.8);
  margin-top: 8rpx;
  display: block;
}

.tab-bar {
  display: flex;
  background: #ffffff;
  padding: 0 20rpx;
  box-shadow: 0 2rpx 10rpx rgba(0, 0, 0, 0.05);
}

.tab-item {
  flex: 1;
  height: 88rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  position: relative;
}

.tab-text {
  font-size: 28rpx;
  color: #666;
}

.tab-text-active {
  color: #667eea;
  font-weight: 600;
}

.tab-indicator {
  position: absolute;
  bottom: 0;
  width: 60rpx;
  height: 6rpx;
  background: #667eea;
  border-radius: 3rpx;
}

.exam-scroll {
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

.exam-cards {
  padding: 20rpx 30rpx;
  display: flex;
  flex-direction: column;
  gap: 20rpx;
}

.exam-card {
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
  padding: 6rpx 16rpx;
  border-radius: 20rpx;
  flex-shrink: 0;
}

.badge-upcoming {
  background: #e3f2fd;
}

.badge-active {
  background: #e8f5e9;
}

.badge-completed {
  background: #f5f5f5;
}

.status-text {
  font-size: 22rpx;
}

.badge-upcoming .status-text {
  color: #1976d2;
}

.badge-active .status-text {
  color: #388e3c;
}

.badge-completed .status-text {
  color: #666;
}

.card-info {
  display: flex;
  flex-direction: column;
  gap: 12rpx;
}

.info-row {
  display: flex;
  align-items: center;
}

.info-icon {
  font-size: 28rpx;
  margin-right: 12rpx;
}

.info-text {
  font-size: 26rpx;
  color: #666;
}

.card-footer {
  margin-top: 20rpx;
  padding-top: 20rpx;
  border-top: 1rpx solid #f0f0f0;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.footer-text {
  font-size: 24rpx;
  color: #999;
}

.footer-active {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.footer-text-active {
  font-size: 26rpx;
  color: #388e3c;
  font-weight: 600;
}

.footer-arrow {
  font-size: 28rpx;
  color: #999;
}

.score-info {
  display: flex;
  align-items: center;
}

.score-label {
  font-size: 24rpx;
  color: #999;
  margin-right: 8rpx;
}

.score-value {
  font-size: 28rpx;
  color: #667eea;
  font-weight: 600;
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
