<template>
  <view class="exam-result-page">
    <!-- Score Display -->
    <view class="score-section" :class="result.passed ? 'section-pass' : 'section-fail'">
      <view class="score-circle" :class="{ 'score-animate': animated }">
        <view class="score-inner">
          <text class="score-value">{{ animatedScore }}</text>
          <text class="score-unit">pts</text>
        </view>
      </view>
      <view class="pass-status">
        <text class="pass-icon">{{ result.passed ? '✓' : '✗' }}</text>
        <text class="pass-text">{{ result.passed ? 'PASSED' : 'FAILED' }}</text>
      </view>
      <text class="score-total">Total: {{ result.totalScore }} points</text>
    </view>

    <!-- Stats Section -->
    <view class="stats-section">
      <view class="stat-card">
        <text class="stat-card-value">{{ result.score }}</text>
        <text class="stat-card-label">Score</text>
      </view>
      <view class="stat-card">
        <text class="stat-card-value">#{{ result.rank }}</text>
        <text class="stat-card-label">Rank</text>
      </view>
      <view class="stat-card">
        <text class="stat-card-value">{{ result.correctRate }}%</text>
        <text class="stat-card-label">Correct Rate</text>
      </view>
      <view class="stat-card">
        <text class="stat-card-value">{{ formatDuration(result.duration) }}</text>
        <text class="stat-card-label">Duration</text>
      </view>
    </view>

    <!-- Question Breakdown -->
    <view class="section">
      <text class="section-title">Question Breakdown</text>
      <view class="breakdown-list">
        <view
          v-for="(item, index) in result.questions"
          :key="item.questionId"
          class="breakdown-item"
          :class="{ 'item-correct': item.isCorrect, 'item-wrong': !item.isCorrect }"
          @tap="toggleDetail(index)"
        >
          <view class="breakdown-header">
            <view class="breakdown-left">
              <view class="breakdown-status" :class="item.isCorrect ? 'status-correct' : 'status-wrong'">
                <text class="status-icon">{{ item.isCorrect ? '✓' : '✗' }}</text>
              </view>
              <view class="breakdown-info">
                <text class="breakdown-num">Q{{ index + 1 }}</text>
                <text class="breakdown-score">{{ item.score }}/{{ item.totalScore }} pts</text>
              </view>
            </view>
            <text class="breakdown-expand">{{ expandedIndex === index ? '▲' : '▼' }}</text>
          </view>

          <text class="breakdown-content">{{ item.content }}</text>

          <!-- Detail -->
          <view v-if="expandedIndex === index" class="breakdown-detail">
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
            <view class="detail-answers">
              <view class="detail-row">
                <text class="detail-label">Your Answer:</text>
                <text class="detail-value" :class="item.isCorrect ? 'value-correct' : 'value-wrong'">
                  {{ item.userAnswer.length > 0 ? item.userAnswer.map(i => optionLabels[i]).join(', ') : 'Not answered' }}
                </text>
              </view>
              <view class="detail-row">
                <text class="detail-label">Correct Answer:</text>
                <text class="detail-value value-correct">
                  {{ item.correctAnswer.map(i => optionLabels[i]).join(', ') }}
                </text>
              </view>
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
      <view class="action-btn btn-answers" @tap="handleViewAnswers">
        <text class="btn-text">View Answers</text>
      </view>
      <view class="action-btn btn-back" @tap="handleBack">
        <text class="btn-text">Back to Exam</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { getExamResult } from '@/api/exam';

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
  score: number;
  totalScore: number;
  explanation?: string;
}

interface ExamResult {
  score: number;
  totalScore: number;
  passed: boolean;
  rank: number;
  correctRate: number;
  duration: number;
  questions: QuestionResult[];
}

const optionLabels = ['A', 'B', 'C', 'D', 'E', 'F', 'G', 'H'];

const attemptId = ref('');
const result = ref<ExamResult>({
  score: 0,
  totalScore: 0,
  passed: false,
  rank: 0,
  correctRate: 0,
  duration: 0,
  questions: [],
});

const animated = ref(false);
const animatedScore = ref(0);
const expandedIndex = ref(-1);

onMounted(() => {
  const pages = getCurrentPages();
  const currentPage = pages[pages.length - 1] as any;
  attemptId.value = currentPage.options?.attemptId || '';
  loadResult();
});

async function loadResult() {
  if (!attemptId.value) return;
  try {
    uni.showLoading({ title: 'Loading...' });
    const res = await getExamResult(attemptId.value);
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
  const target = result.value.score;
  const duration = 1500;
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

function formatDuration(seconds: number): string {
  const mins = Math.floor(seconds / 60);
  const secs = seconds % 60;
  return `${mins}m${secs}s`;
}

function toggleDetail(index: number) {
  expandedIndex.value = expandedIndex.value === index ? -1 : index;
}

function handleViewAnswers() {
  uni.navigateTo({ url: `/pages/exam/answer?id=${result.value.examId}&review=true&attemptId=${attemptId.value}` });
}

function handleBack() {
  uni.navigateBack();
}
</script>

<style scoped>
.exam-result-page {
  min-height: 100vh;
  background-color: #f5f5f5;
  padding-bottom: 160rpx;
}

.score-section {
  padding: 60rpx 40rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.section-pass {
  background: linear-gradient(135deg, #4caf50 0%, #81c784 100%);
}

.section-fail {
  background: linear-gradient(135deg, #f44336 0%, #e57373 100%);
}

.score-circle {
  width: 240rpx;
  height: 240rpx;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.2);
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 24rpx;
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
  font-size: 28rpx;
  color: rgba(255, 255, 255, 0.8);
  margin-left: 8rpx;
}

.pass-status {
  display: flex;
  align-items: center;
  margin-bottom: 12rpx;
}

.pass-icon {
  font-size: 40rpx;
  color: #ffffff;
  margin-right: 12rpx;
  font-weight: bold;
}

.pass-text {
  font-size: 36rpx;
  font-weight: bold;
  color: #ffffff;
}

.score-total {
  font-size: 28rpx;
  color: rgba(255, 255, 255, 0.8);
}

.stats-section {
  display: flex;
  margin: -30rpx 30rpx 0;
  background: #ffffff;
  border-radius: 20rpx;
  padding: 30rpx;
  box-shadow: 0 4rpx 20rpx rgba(0, 0, 0, 0.08);
  position: relative;
  z-index: 1;
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

.breakdown-list {
  background: #ffffff;
  border-radius: 16rpx;
  overflow: hidden;
}

.breakdown-item {
  border-bottom: 1rpx solid #f0f0f0;
  padding: 24rpx;
}

.breakdown-item:last-child {
  border-bottom: none;
}

.breakdown-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12rpx;
}

.breakdown-left {
  display: flex;
  align-items: center;
}

.breakdown-status {
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

.breakdown-info {
  display: flex;
  flex-direction: column;
}

.breakdown-num {
  font-size: 28rpx;
  font-weight: 600;
  color: #333;
}

.breakdown-score {
  font-size: 22rpx;
  color: #999;
  margin-top: 4rpx;
}

.breakdown-expand {
  font-size: 24rpx;
  color: #999;
}

.breakdown-content {
  font-size: 28rpx;
  color: #333;
  line-height: 1.5;
  display: block;
}

.breakdown-detail {
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

.detail-answers {
  margin-top: 16rpx;
  padding: 16rpx;
  background: #f5f5f5;
  border-radius: 12rpx;
}

.detail-row {
  display: flex;
  margin-bottom: 8rpx;
}

.detail-row:last-child {
  margin-bottom: 0;
}

.detail-label {
  font-size: 26rpx;
  color: #666;
  margin-right: 12rpx;
  min-width: 160rpx;
}

.detail-value {
  font-size: 26rpx;
  font-weight: 600;
}

.value-correct {
  color: #4caf50;
}

.value-wrong {
  color: #f44336;
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

.btn-answers {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.btn-answers .btn-text {
  color: #ffffff;
  font-weight: 600;
}

.btn-back {
  background: #f0f0f0;
}

.btn-back .btn-text {
  color: #666;
}

.btn-text {
  font-size: 30rpx;
}
</style>
