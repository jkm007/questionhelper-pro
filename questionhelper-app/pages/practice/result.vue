<template>
  <view class="result-page">
    <!-- Score Section -->
    <view class="score-section">
      <view class="score-circle" :class="{ 'score-animate': animated }">
        <view class="score-inner">
          <text class="score-value">{{ animatedScore }}</text>
          <text class="score-unit">%</text>
        </view>
        <view class="score-ring">
          <view class="score-ring-fill" :style="{ '--score-percent': result.accuracy }"></view>
        </view>
      </view>
      <text class="score-label">Accuracy</text>
      <text class="score-message">{{ getScoreMessage(result.accuracy) }}</text>
    </view>

    <!-- Stats Section -->
    <view class="stats-section">
      <view class="stat-card">
        <text class="stat-card-value">{{ result.totalCount }}</text>
        <text class="stat-card-label">Total</text>
      </view>
      <view class="stat-card">
        <text class="stat-card-value stat-correct">{{ result.correctCount }}</text>
        <text class="stat-card-label">Correct</text>
      </view>
      <view class="stat-card">
        <text class="stat-card-value stat-wrong">{{ result.wrongCount }}</text>
        <text class="stat-card-label">Wrong</text>
      </view>
      <view class="stat-card">
        <text class="stat-card-value">{{ formatDuration(result.duration) }}</text>
        <text class="stat-card-label">Duration</text>
      </view>
    </view>

    <!-- Question List -->
    <view class="section">
      <text class="section-title">Question Details</text>
      <view class="question-list">
        <view
          v-for="(item, index) in result.questions"
          :key="item.questionId"
          class="question-item"
          :class="{ 'item-correct': item.isCorrect, 'item-wrong': !item.isCorrect }"
          @tap="toggleQuestionDetail(index)"
        >
          <view class="question-item-header">
            <view class="question-item-left">
              <view class="question-status" :class="item.isCorrect ? 'status-correct' : 'status-wrong'">
                <text class="status-icon">{{ item.isCorrect ? '✓' : '✗' }}</text>
              </view>
              <text class="question-item-num">Q{{ index + 1 }}</text>
            </view>
            <text class="question-item-expand">{{ expandedIndex === index ? '▲' : '▼' }}</text>
          </view>

          <text class="question-item-content">{{ item.content }}</text>

          <!-- Expanded Detail -->
          <view v-if="expandedIndex === index" class="question-detail">
            <view v-for="(option, oIndex) in item.options" :key="oIndex" class="detail-option">
              <text class="detail-option-label">{{ optionLabels[oIndex] }}.</text>
              <text
                class="detail-option-text"
                :class="{
                  'option-is-correct': option.isCorrect,
                  'option-is-wrong-selected': item.userAnswer.includes(oIndex) && !option.isCorrect,
                }"
              >
                {{ option.content }}
              </text>
            </view>
            <view v-if="item.userAnswer.length > 0" class="detail-answer">
              <text class="detail-answer-label">Your Answer:</text>
              <text class="detail-answer-value">{{ item.userAnswer.map(i => optionLabels[i]).join(', ') }}</text>
            </view>
            <view v-if="item.explanation" class="detail-explanation">
              <text class="detail-explanation-label">Explanation:</text>
              <text class="detail-explanation-text">{{ item.explanation }}</text>
            </view>
          </view>
        </view>
      </view>
    </view>

    <!-- Bottom Actions -->
    <view class="bottom-bar">
      <view class="action-btn btn-retry" @tap="handleRetry">
        <text class="btn-text">Retry</text>
      </view>
      <view class="action-btn btn-home" @tap="handleBackHome">
        <text class="btn-text">Back to Practice</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { getPracticeResult } from '@/api/practice';

interface Option {
  content: string;
  isCorrect: boolean;
}

interface QuestionResult {
  questionId: string;
  content: string;
  options: Option[];
  userAnswer: number[];
  correctAnswer: number[];
  isCorrect: boolean;
  explanation?: string;
}

interface PracticeResult {
  totalCount: number;
  correctCount: number;
  wrongCount: number;
  accuracy: number;
  duration: number;
  mode: string;
  questions: QuestionResult[];
}

const optionLabels = ['A', 'B', 'C', 'D', 'E', 'F', 'G', 'H'];

const result = ref<PracticeResult>({
  totalCount: 0,
  correctCount: 0,
  wrongCount: 0,
  accuracy: 0,
  duration: 0,
  mode: '',
  questions: [],
});

const animated = ref(false);
const animatedScore = ref(0);
const expandedIndex = ref(-1);
const recordId = ref('');

onMounted(() => {
  const pages = getCurrentPages();
  const currentPage = pages[pages.length - 1] as any;
  recordId.value = currentPage.options?.recordId || '';
  loadResult();
});

async function loadResult() {
  if (!recordId.value) return;
  try {
    uni.showLoading({ title: 'Loading...' });
    const res = await getPracticeResult(recordId.value);
    if (res.code === 0) {
      result.value = res.data;
      startScoreAnimation();
    }
  } catch (e) {
    console.error('Failed to load result:', e);
    uni.showToast({ title: 'Failed to load', icon: 'none' });
  } finally {
    uni.hideLoading();
  }
}

function startScoreAnimation() {
  animated.value = true;
  const target = result.value.accuracy;
  const duration = 1000;
  const stepTime = 20;
  const steps = duration / stepTime;
  const increment = target / steps;
  let current = 0;
  const interval = setInterval(() => {
    current += increment;
    if (current >= target) {
      current = target;
      clearInterval(interval);
    }
    animatedScore.value = Math.round(current);
  }, stepTime);
}

function getScoreMessage(accuracy: number): string {
  if (accuracy >= 90) return 'Excellent! Keep it up!';
  if (accuracy >= 80) return 'Great job!';
  if (accuracy >= 70) return 'Good work!';
  if (accuracy >= 60) return 'Not bad, keep practicing!';
  return 'Keep trying, you will improve!';
}

function formatDuration(seconds: number): string {
  const mins = Math.floor(seconds / 60);
  const secs = seconds % 60;
  return `${mins}m${secs}s`;
}

function toggleQuestionDetail(index: number) {
  expandedIndex.value = expandedIndex.value === index ? -1 : index;
}

function handleRetry() {
  uni.redirectTo({ url: `/pages/practice/session?mode=${result.value.mode}` });
}

function handleBackHome() {
  uni.navigateBack({ delta: 2 });
}
</script>

<style scoped>
.result-page {
  min-height: 100vh;
  background-color: #f5f5f5;
  padding-bottom: 160rpx;
}

.score-section {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 60rpx 40rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.score-circle {
  width: 240rpx;
  height: 240rpx;
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
}

.score-inner {
  display: flex;
  align-items: baseline;
}

.score-value {
  font-size: 72rpx;
  font-weight: bold;
  color: #ffffff;
}

.score-unit {
  font-size: 32rpx;
  color: rgba(255, 255, 255, 0.8);
  margin-left: 4rpx;
}

.score-ring {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  border-radius: 50%;
  border: 8rpx solid rgba(255, 255, 255, 0.2);
}

.score-ring-fill {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  border-radius: 50%;
  border: 8rpx solid transparent;
  border-top-color: #ffffff;
  transform: rotate(calc(var(--score-percent) * 3.6deg));
  transition: transform 1s ease;
}

.score-label {
  font-size: 28rpx;
  color: rgba(255, 255, 255, 0.8);
  margin-top: 16rpx;
}

.score-message {
  font-size: 30rpx;
  color: #ffffff;
  font-weight: 600;
  margin-top: 12rpx;
}

.stats-section {
  display: flex;
  margin: -30rpx 30rpx 0;
  background: #ffffff;
  border-radius: 20rpx;
  padding: 30rpx;
  box-shadow: 0 4rpx 20rpx rgba(0, 0, 0, 0.08);
}

.stat-card {
  flex: 1;
  text-align: center;
}

.stat-card-value {
  font-size: 36rpx;
  font-weight: bold;
  color: #333;
  display: block;
}

.stat-correct {
  color: #4caf50;
}

.stat-wrong {
  color: #f44336;
}

.stat-card-label {
  font-size: 22rpx;
  color: #999;
  margin-top: 8rpx;
  display: block;
}

.section {
  margin: 30rpx;
}

.section-title {
  font-size: 32rpx;
  font-weight: 600;
  color: #333;
  margin-bottom: 20rpx;
  display: block;
}

.question-list {
  background: #ffffff;
  border-radius: 16rpx;
  overflow: hidden;
}

.question-item {
  border-bottom: 1rpx solid #f0f0f0;
  padding: 24rpx;
}

.question-item:last-child {
  border-bottom: none;
}

.question-item-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12rpx;
}

.question-item-left {
  display: flex;
  align-items: center;
}

.question-status {
  width: 44rpx;
  height: 44rpx;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 16rpx;
}

.status-correct {
  background: #e8f5e9;
}

.status-wrong {
  background: #ffebee;
}

.status-icon {
  font-size: 24rpx;
  font-weight: bold;
}

.status-correct .status-icon {
  color: #4caf50;
}

.status-wrong .status-icon {
  color: #f44336;
}

.question-item-num {
  font-size: 28rpx;
  font-weight: 600;
  color: #333;
}

.question-item-expand {
  font-size: 24rpx;
  color: #999;
}

.question-item-content {
  font-size: 28rpx;
  color: #333;
  line-height: 1.5;
  display: block;
}

.question-detail {
  margin-top: 20rpx;
  padding-top: 20rpx;
  border-top: 1rpx solid #f0f0f0;
}

.detail-option {
  display: flex;
  padding: 8rpx 0;
}

.detail-option-label {
  font-size: 26rpx;
  color: #666;
  margin-right: 12rpx;
  min-width: 40rpx;
}

.detail-option-text {
  font-size: 26rpx;
  color: #333;
  flex: 1;
}

.option-is-correct {
  color: #4caf50;
  font-weight: 600;
}

.option-is-wrong-selected {
  color: #f44336;
  text-decoration: line-through;
}

.detail-answer {
  margin-top: 16rpx;
  display: flex;
  align-items: center;
}

.detail-answer-label {
  font-size: 26rpx;
  color: #666;
  margin-right: 12rpx;
}

.detail-answer-value {
  font-size: 26rpx;
  color: #667eea;
  font-weight: 600;
}

.detail-explanation {
  margin-top: 16rpx;
  background: #fff8e1;
  padding: 16rpx;
  border-radius: 12rpx;
}

.detail-explanation-label {
  font-size: 24rpx;
  color: #e65100;
  font-weight: 600;
  display: block;
  margin-bottom: 8rpx;
}

.detail-explanation-text {
  font-size: 24rpx;
  color: #666;
  line-height: 1.6;
  display: block;
}

.bottom-bar {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  display: flex;
  padding: 20rpx 30rpx;
  padding-bottom: calc(20rpx + env(safe-area-inset-bottom));
  background: #ffffff;
  gap: 20rpx;
  box-shadow: 0 -2rpx 10rpx rgba(0, 0, 0, 0.05);
}

.action-btn {
  flex: 1;
  height: 88rpx;
  border-radius: 44rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.btn-retry {
  background: #f0f0f0;
}

.btn-retry .btn-text {
  color: #666;
}

.btn-home {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.btn-home .btn-text {
  color: #ffffff;
  font-weight: 600;
}

.btn-text {
  font-size: 30rpx;
}
</style>
