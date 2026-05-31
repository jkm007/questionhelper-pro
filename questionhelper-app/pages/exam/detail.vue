<template>
  <view class="exam-detail-page">
    <!-- Header -->
    <view class="header">
      <view class="header-nav">
        <view class="nav-back" @tap="handleBack">
          <text class="nav-back-icon">←</text>
        </view>
        <text class="nav-title">Exam Details</text>
        <view class="nav-placeholder"></view>
      </view>
    </view>

    <!-- Info Card -->
    <view class="info-card">
      <view class="info-header">
        <text class="info-title">{{ exam.title }}</text>
        <view class="status-badge" :class="getStatusClass(exam.status)">
          <text class="status-text">{{ getStatusText(exam.status) }}</text>
        </view>
      </view>

      <text class="info-desc">{{ exam.description }}</text>

      <view class="info-grid">
        <view class="info-item">
          <text class="info-item-icon">📅</text>
          <view class="info-item-content">
            <text class="info-item-label">Time</text>
            <text class="info-item-value">{{ formatDateRange(exam.startTime, exam.endTime) }}</text>
          </view>
        </view>
        <view class="info-item">
          <text class="info-item-icon">⏱</text>
          <view class="info-item-content">
            <text class="info-item-label">Duration</text>
            <text class="info-item-value">{{ exam.duration }} minutes</text>
          </view>
        </view>
        <view class="info-item">
          <text class="info-item-icon">📊</text>
          <view class="info-item-content">
            <text class="info-item-label">Total Score</text>
            <text class="info-item-value">{{ exam.totalScore }} points</text>
          </view>
        </view>
        <view class="info-item">
          <text class="info-item-icon">👥</text>
          <view class="info-item-content">
            <text class="info-item-label">Participants</text>
            <text class="info-item-value">{{ exam.participantCount }} enrolled</text>
          </view>
        </view>
      </view>
    </view>

    <!-- Rules Section -->
    <view class="section" v-if="exam.rules && exam.rules.length > 0">
      <text class="section-title">Exam Rules</text>
      <view class="rules-card">
        <view v-for="(rule, index) in exam.rules" :key="index" class="rule-item">
          <text class="rule-index">{{ index + 1 }}.</text>
          <text class="rule-text">{{ rule }}</text>
        </view>
      </view>
    </view>

    <!-- Participants -->
    <view class="section">
      <view class="section-header">
        <text class="section-title">Participants</text>
        <text class="section-count">{{ exam.participantCount }} enrolled</text>
      </view>
      <view class="participants-list">
        <view v-for="participant in participants" :key="participant.id" class="participant-item">
          <image class="participant-avatar" :src="participant.avatar || '/static/default-avatar.png'" mode="aspectFill" />
          <view class="participant-info">
            <text class="participant-name">{{ participant.name }}</text>
            <text class="participant-time">Enrolled {{ formatTime(participant.enrollTime) }}</text>
          </view>
        </view>
        <view v-if="participants.length === 0" class="empty-participants">
          <text class="empty-text">No participants yet</text>
        </view>
      </view>
    </view>

    <!-- History Attempts -->
    <view class="section" v-if="exam.status === 'completed' && historyAttempts.length > 0">
      <text class="section-title">My Attempts</text>
      <view class="history-list">
        <view
          v-for="attempt in historyAttempts"
          :key="attempt.id"
          class="history-item"
          @tap="goToResult(attempt.id)"
        >
          <view class="history-left">
            <text class="history-attempt">Attempt {{ attempt.attemptNumber }}</text>
            <text class="history-time">{{ formatTime(attempt.submitTime) }}</text>
          </view>
          <view class="history-right">
            <text class="history-score">{{ attempt.score }}/{{ exam.totalScore }}</text>
            <view class="history-status" :class="attempt.passed ? 'status-pass' : 'status-fail'">
              <text class="history-status-text">{{ attempt.passed ? 'Passed' : 'Failed' }}</text>
            </view>
          </view>
        </view>
      </view>
    </view>

    <!-- Bottom Action -->
    <view class="bottom-bar">
      <view
        class="enter-btn"
        :class="{ 'btn-disabled': !canEnter }"
        @tap="handleEnterExam"
      >
        <text class="enter-btn-text">{{ getEnterButtonText() }}</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { getExamDetail, getExamParticipants, getExamHistory } from '@/api/exam';

interface ExamDetail {
  id: string;
  title: string;
  description: string;
  startTime: string;
  endTime: string;
  duration: number;
  totalScore: number;
  participantCount: number;
  status: 'upcoming' | 'in_progress' | 'completed';
  rules?: string[];
}

interface Participant {
  id: string;
  name: string;
  avatar?: string;
  enrollTime: string;
}

interface HistoryAttempt {
  id: string;
  attemptNumber: number;
  score: number;
  passed: boolean;
  submitTime: string;
  duration: number;
}

const examId = ref('');
const exam = ref<ExamDetail>({
  id: '',
  title: '',
  description: '',
  startTime: '',
  endTime: '',
  duration: 0,
  totalScore: 0,
  participantCount: 0,
  status: 'upcoming',
  rules: [],
});

const participants = ref<Participant[]>([]);
const historyAttempts = ref<HistoryAttempt[]>([]);

const canEnter = computed(() => {
  if (exam.value.status === 'in_progress') return true;
  if (exam.value.status === 'upcoming') {
    const now = Date.now();
    const start = new Date(exam.value.startTime).getTime();
    return now >= start;
  }
  return false;
});

onMounted(() => {
  const pages = getCurrentPages();
  const currentPage = pages[pages.length - 1] as any;
  examId.value = currentPage.options?.id || '';
  loadExamDetail();
  loadParticipants();
  loadHistory();
});

async function loadExamDetail() {
  try {
    uni.showLoading({ title: 'Loading...' });
    const res = await getExamDetail(examId.value);
    if (res.code === 0) {
      exam.value = res.data;
    }
  } catch (e) {
    console.error('Failed to load exam detail:', e);
    uni.showToast({ title: 'Failed to load', icon: 'none' });
  } finally {
    uni.hideLoading();
  }
}

async function loadParticipants() {
  try {
    const res = await getExamParticipants(examId.value, { limit: 10 });
    if (res.code === 0) {
      participants.value = res.data;
    }
  } catch (e) {
    console.error('Failed to load participants:', e);
  }
}

async function loadHistory() {
  try {
    const res = await getExamHistory(examId.value);
    if (res.code === 0) {
      historyAttempts.value = res.data;
    }
  } catch (e) {
    console.error('Failed to load history:', e);
  }
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

function formatDateRange(start: string, end: string): string {
  const startDate = new Date(start);
  const endDate = new Date(end);
  const startStr = `${startDate.getFullYear()}/${startDate.getMonth() + 1}/${startDate.getDate()} ${String(startDate.getHours()).padStart(2, '0')}:${String(startDate.getMinutes()).padStart(2, '0')}`;
  const endStr = `${endDate.getMonth() + 1}/${endDate.getDate()} ${String(endDate.getHours()).padStart(2, '0')}:${String(endDate.getMinutes()).padStart(2, '0')}`;
  return `${startStr} - ${endStr}`;
}

function formatTime(time: string): string {
  if (!time) return '';
  const date = new Date(time);
  return `${date.getMonth() + 1}/${date.getDate()} ${String(date.getHours()).padStart(2, '0')}:${String(date.getMinutes()).padStart(2, '0')}`;
}

function getEnterButtonText(): string {
  if (exam.value.status === 'in_progress') return 'Enter Exam';
  if (exam.value.status === 'upcoming') {
    const now = Date.now();
    const start = new Date(exam.value.startTime).getTime();
    if (now < start) return 'Not Started Yet';
    return 'Enter Exam';
  }
  return 'Exam Ended';
}

function handleEnterExam() {
  if (!canEnter.value) {
    uni.showToast({ title: 'Cannot enter exam now', icon: 'none' });
    return;
  }
  uni.navigateTo({ url: `/pages/exam/answer?id=${examId.value}` });
}

function goToResult(attemptId: string) {
  uni.navigateTo({ url: `/pages/exam/result?attemptId=${attemptId}` });
}

function handleBack() {
  uni.navigateBack();
}
</script>

<style scoped>
.exam-detail-page {
  min-height: 100vh;
  background-color: #f5f5f5;
  padding-bottom: 160rpx;
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

.info-card {
  margin: -20rpx 30rpx 0;
  background: #ffffff;
  border-radius: 20rpx;
  padding: 30rpx;
  box-shadow: 0 4rpx 20rpx rgba(0, 0, 0, 0.08);
  position: relative;
  z-index: 1;
}

.info-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 16rpx;
}

.info-title {
  font-size: 36rpx;
  font-weight: bold;
  color: #333;
  flex: 1;
  margin-right: 16rpx;
}

.status-badge {
  padding: 8rpx 20rpx;
  border-radius: 24rpx;
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
  font-size: 24rpx;
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

.info-desc {
  font-size: 28rpx;
  color: #666;
  line-height: 1.6;
  margin-bottom: 24rpx;
  display: block;
}

.info-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 20rpx;
}

.info-item {
  width: calc(50% - 10rpx);
  display: flex;
  align-items: flex-start;
}

.info-item-icon {
  font-size: 32rpx;
  margin-right: 12rpx;
  margin-top: 4rpx;
}

.info-item-content {
  flex: 1;
}

.info-item-label {
  font-size: 22rpx;
  color: #999;
  display: block;
}

.info-item-value {
  font-size: 26rpx;
  color: #333;
  font-weight: 500;
  margin-top: 4rpx;
  display: block;
}

.section {
  margin: 30rpx;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16rpx;
}

.section-title {
  font-size: 32rpx;
  font-weight: 600;
  color: #333;
  display: block;
}

.section-count {
  font-size: 24rpx;
  color: #999;
}

.rules-card {
  background: #ffffff;
  border-radius: 16rpx;
  padding: 24rpx;
  margin-top: 16rpx;
}

.rule-item {
  display: flex;
  margin-bottom: 16rpx;
}

.rule-item:last-child {
  margin-bottom: 0;
}

.rule-index {
  font-size: 26rpx;
  color: #667eea;
  font-weight: 600;
  margin-right: 12rpx;
  min-width: 32rpx;
}

.rule-text {
  font-size: 26rpx;
  color: #666;
  line-height: 1.5;
  flex: 1;
}

.participants-list {
  background: #ffffff;
  border-radius: 16rpx;
  overflow: hidden;
  margin-top: 16rpx;
}

.participant-item {
  display: flex;
  align-items: center;
  padding: 20rpx 24rpx;
  border-bottom: 1rpx solid #f0f0f0;
}

.participant-item:last-child {
  border-bottom: none;
}

.participant-avatar {
  width: 72rpx;
  height: 72rpx;
  border-radius: 50%;
  margin-right: 20rpx;
  background: #f0f0f0;
}

.participant-info {
  flex: 1;
}

.participant-name {
  font-size: 28rpx;
  color: #333;
  font-weight: 500;
  display: block;
}

.participant-time {
  font-size: 22rpx;
  color: #999;
  margin-top: 4rpx;
  display: block;
}

.empty-participants {
  padding: 40rpx;
  text-align: center;
}

.empty-text {
  font-size: 26rpx;
  color: #999;
}

.history-list {
  background: #ffffff;
  border-radius: 16rpx;
  overflow: hidden;
  margin-top: 16rpx;
}

.history-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 24rpx;
  border-bottom: 1rpx solid #f0f0f0;
}

.history-item:last-child {
  border-bottom: none;
}

.history-attempt {
  font-size: 28rpx;
  color: #333;
  font-weight: 500;
  display: block;
}

.history-time {
  font-size: 22rpx;
  color: #999;
  margin-top: 4rpx;
  display: block;
}

.history-right {
  display: flex;
  align-items: center;
  gap: 16rpx;
}

.history-score {
  font-size: 32rpx;
  font-weight: bold;
  color: #667eea;
}

.history-status {
  padding: 6rpx 16rpx;
  border-radius: 20rpx;
}

.status-pass {
  background: #e8f5e9;
}

.status-fail {
  background: #ffebee;
}

.history-status-text {
  font-size: 22rpx;
}

.status-pass .history-status-text {
  color: #388e3c;
}

.status-fail .history-status-text {
  color: #d32f2f;
}

.bottom-bar {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  padding: 20rpx 30rpx;
  padding-bottom: calc(20rpx + env(safe-area-inset-bottom));
  background: #ffffff;
  box-shadow: 0 -2rpx 10rpx rgba(0, 0, 0, 0.05);
}

.enter-btn {
  height: 96rpx;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border-radius: 48rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.btn-disabled {
  background: #cccccc;
}

.enter-btn-text {
  font-size: 32rpx;
  font-weight: 600;
  color: #ffffff;
}
</style>
