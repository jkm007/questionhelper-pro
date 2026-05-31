<template>
  <view class="practice-page">
    <!-- Header -->
    <view class="header">
      <text class="header-title">Practice</text>
      <text class="header-subtitle">Improve your skills</text>
    </view>

    <!-- Today Stats -->
    <view class="stats-card">
      <text class="stats-card-title">Today's Progress</text>
      <view class="stats-row">
        <view class="stat-item">
          <text class="stat-value">{{ todayStats.count }}</text>
          <text class="stat-label">Questions</text>
        </view>
        <view class="stat-item">
          <text class="stat-value">{{ todayStats.accuracy }}%</text>
          <text class="stat-label">Accuracy</text>
        </view>
        <view class="stat-item">
          <text class="stat-value">{{ formatDuration(todayStats.duration) }}</text>
          <text class="stat-label">Duration</text>
        </view>
      </view>
    </view>

    <!-- Quick Start -->
    <view class="quick-start" @tap="handleQuickStart">
      <view class="quick-start-content">
        <text class="quick-start-icon">▶</text>
        <view class="quick-start-text">
          <text class="quick-start-title">Quick Start</text>
          <text class="quick-start-desc">Start practice with recommended settings</text>
        </view>
      </view>
      <text class="quick-start-arrow">></text>
    </view>

    <!-- Practice Modes -->
    <view class="section">
      <text class="section-title">Practice Modes</text>
      <view class="mode-grid">
        <view
          v-for="mode in practiceModes"
          :key="mode.key"
          class="mode-card"
          :style="{ backgroundColor: mode.color }"
          @tap="handleModeSelect(mode.key)"
        >
          <text class="mode-icon">{{ mode.icon }}</text>
          <text class="mode-name">{{ mode.name }}</text>
          <text class="mode-desc">{{ mode.description }}</text>
        </view>
      </view>
    </view>

    <!-- Recent Records -->
    <view class="section">
      <view class="section-header">
        <text class="section-title">Recent Practice</text>
        <text class="section-more" @tap="goToAllRecords">View All</text>
      </view>
      <view v-if="recentRecords.length === 0" class="empty-state">
        <text class="empty-text">No recent practice records</text>
      </view>
      <view v-else class="record-list">
        <view
          v-for="record in recentRecords"
          :key="record.id"
          class="record-item"
          @tap="goToRecordDetail(record.id)"
        >
          <view class="record-left">
            <text class="record-mode">{{ getModeName(record.mode) }}</text>
            <text class="record-time">{{ formatTime(record.createTime) }}</text>
          </view>
          <view class="record-right">
            <text class="record-accuracy">{{ record.accuracy }}%</text>
            <text class="record-count">{{ record.correctCount }}/{{ record.totalCount }}</text>
          </view>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { getRecentRecords, getTodayStats } from '@/api/practice';

interface TodayStats {
  count: number;
  accuracy: number;
  duration: number;
}

interface PracticeRecord {
  id: string;
  mode: string;
  accuracy: number;
  correctCount: number;
  totalCount: number;
  createTime: string;
}

const practiceModes = [
  { key: 'random', name: 'Random', icon: '🔀', description: 'Random questions', color: '#E3F2FD' },
  { key: 'sequential', name: 'Sequential', icon: '📋', description: 'In order', color: '#E8F5E9' },
  { key: 'wrong', name: 'Wrong Questions', icon: '❌', description: 'Review mistakes', color: '#FFF3E0' },
  { key: 'favorites', name: 'Favorites', icon: '⭐', description: 'Saved questions', color: '#FCE4EC' },
  { key: 'mock', name: 'Mock Exam', icon: '📝', description: 'Simulate exam', color: '#E8EAF6' },
  { key: 'challenge', name: 'Challenge', icon: '🏆', description: 'Challenge mode', color: '#F3E5F5' },
  { key: 'timed', name: 'Timed', icon: '⏱', description: 'Time limited', color: '#E0F7FA' },
];

const todayStats = ref<TodayStats>({
  count: 0,
  accuracy: 0,
  duration: 0,
});

const recentRecords = ref<PracticeRecord[]>([]);

onMounted(async () => {
  await loadTodayStats();
  await loadRecentRecords();
});

async function loadTodayStats() {
  try {
    const res = await getTodayStats();
    if (res.code === 0) {
      todayStats.value = res.data;
    }
  } catch (e) {
    console.error('Failed to load today stats:', e);
  }
}

async function loadRecentRecords() {
  try {
    const res = await getRecentRecords({ limit: 5 });
    if (res.code === 0) {
      recentRecords.value = res.data;
    }
  } catch (e) {
    console.error('Failed to load recent records:', e);
  }
}

function formatDuration(seconds: number): string {
  const mins = Math.floor(seconds / 60);
  if (mins < 60) return `${mins}m`;
  const hours = Math.floor(mins / 60);
  const remainMins = mins % 60;
  return `${hours}h${remainMins}m`;
}

function formatTime(time: string): string {
  if (!time) return '';
  const date = new Date(time);
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const day = String(date.getDate()).padStart(2, '0');
  const hour = String(date.getHours()).padStart(2, '0');
  const minute = String(date.getMinutes()).padStart(2, '0');
  return `${month}-${day} ${hour}:${minute}`;
}

function getModeName(mode: string): string {
  const found = practiceModes.find((m) => m.key === mode);
  return found ? found.name : mode;
}

function handleQuickStart() {
  uni.navigateTo({ url: '/pages/practice/session?mode=random' });
}

function handleModeSelect(mode: string) {
  uni.navigateTo({ url: `/pages/practice/session?mode=${mode}` });
}

function goToAllRecords() {
  uni.navigateTo({ url: '/pages/practice/records' });
}

function goToRecordDetail(id: string) {
  uni.navigateTo({ url: `/pages/practice/result?recordId=${id}` });
}
</script>

<style scoped>
.practice-page {
  min-height: 100vh;
  background-color: #f5f5f5;
  padding-bottom: 40rpx;
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

.stats-card {
  margin: -30rpx 30rpx 0;
  background: #ffffff;
  border-radius: 20rpx;
  padding: 30rpx;
  box-shadow: 0 4rpx 20rpx rgba(0, 0, 0, 0.08);
}

.stats-card-title {
  font-size: 28rpx;
  color: #333;
  font-weight: 600;
  margin-bottom: 24rpx;
  display: block;
}

.stats-row {
  display: flex;
  justify-content: space-around;
}

.stat-item {
  text-align: center;
}

.stat-value {
  font-size: 40rpx;
  font-weight: bold;
  color: #667eea;
  display: block;
}

.stat-label {
  font-size: 22rpx;
  color: #999;
  margin-top: 8rpx;
  display: block;
}

.quick-start {
  margin: 30rpx;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 20rpx;
  padding: 30rpx;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.quick-start-content {
  display: flex;
  align-items: center;
}

.quick-start-icon {
  font-size: 48rpx;
  color: #ffffff;
  margin-right: 24rpx;
}

.quick-start-title {
  font-size: 30rpx;
  font-weight: bold;
  color: #ffffff;
  display: block;
}

.quick-start-desc {
  font-size: 22rpx;
  color: rgba(255, 255, 255, 0.8);
  margin-top: 4rpx;
  display: block;
}

.quick-start-arrow {
  font-size: 32rpx;
  color: rgba(255, 255, 255, 0.8);
}

.section {
  margin: 30rpx;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20rpx;
}

.section-title {
  font-size: 32rpx;
  font-weight: 600;
  color: #333;
  display: block;
}

.section-more {
  font-size: 24rpx;
  color: #667eea;
}

.mode-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 20rpx;
}

.mode-card {
  width: calc(50% - 10rpx);
  border-radius: 16rpx;
  padding: 24rpx;
  box-sizing: border-box;
}

.mode-icon {
  font-size: 48rpx;
  display: block;
  margin-bottom: 12rpx;
}

.mode-name {
  font-size: 28rpx;
  font-weight: 600;
  color: #333;
  display: block;
}

.mode-desc {
  font-size: 22rpx;
  color: #666;
  margin-top: 8rpx;
  display: block;
}

.empty-state {
  text-align: center;
  padding: 60rpx 0;
}

.empty-text {
  font-size: 26rpx;
  color: #999;
}

.record-list {
  background: #ffffff;
  border-radius: 16rpx;
  overflow: hidden;
}

.record-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 24rpx 30rpx;
  border-bottom: 1rpx solid #f0f0f0;
}

.record-item:last-child {
  border-bottom: none;
}

.record-mode {
  font-size: 28rpx;
  color: #333;
  font-weight: 500;
  display: block;
}

.record-time {
  font-size: 22rpx;
  color: #999;
  margin-top: 8rpx;
  display: block;
}

.record-right {
  text-align: right;
}

.record-accuracy {
  font-size: 32rpx;
  font-weight: bold;
  color: #667eea;
  display: block;
}

.record-count {
  font-size: 22rpx;
  color: #999;
  margin-top: 4rpx;
  display: block;
}
</style>
