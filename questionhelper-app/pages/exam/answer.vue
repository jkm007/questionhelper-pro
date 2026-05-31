<template>
  <view class="exam-answer-page">
    <!-- Custom Navigation Bar -->
    <view class="nav-bar">
      <view class="nav-left" @tap="handleBack">
        <text class="nav-back-icon">←</text>
      </view>
      <view class="nav-center">
        <text class="nav-timer">{{ formatCountdown(remainingSeconds) }}</text>
        <text class="nav-timer-label">Remaining</text>
      </view>
      <view class="nav-right" @tap="toggleQuestionGrid">
        <text class="nav-grid-icon">☰</text>
      </view>
    </view>

    <!-- Question Number Grid Drawer -->
    <view v-if="showGrid" class="grid-overlay" @tap="toggleQuestionGrid">
      <view class="grid-drawer" @tap.stop>
        <view class="grid-header">
          <text class="grid-title">Question Navigator</text>
          <text class="grid-close" @tap="toggleQuestionGrid">✕</text>
        </view>
        <scroll-view scroll-y class="grid-scroll">
          <view class="grid-body">
            <view
              v-for="(q, index) in questions"
              :key="q.id"
              class="grid-item"
              :class="{
                'grid-current': index === currentIndex,
                'grid-answered': answers[q.id] !== undefined,
                'grid-flagged': flaggedQuestions.has(q.id),
              }"
              @tap="jumpToQuestion(index)"
            >
              <text class="grid-item-text">{{ index + 1 }}</text>
            </view>
          </view>
        </scroll-view>
        <view class="grid-legend">
          <view class="legend-item">
            <view class="legend-dot legend-current"></view>
            <text class="legend-text">Current</text>
          </view>
          <view class="legend-item">
            <view class="legend-dot legend-answered"></view>
            <text class="legend-text">Answered</text>
          </view>
          <view class="legend-item">
            <view class="legend-dot legend-flagged"></view>
            <text class="legend-text">Flagged</text>
          </view>
        </view>
      </view>
    </view>

    <!-- Question Content -->
    <swiper
      class="question-swiper"
      :current="currentIndex"
      @change="onSwiperChange"
      :duration="200"
    >
      <swiper-item v-for="(question, qIndex) in questions" :key="question.id">
        <scroll-view scroll-y class="question-scroll">
          <view class="question-card">
            <view class="question-header">
              <view class="question-meta">
                <text class="question-number">Q{{ qIndex + 1 }}</text>
                <text class="question-score">{{ question.score }} pts</text>
              </view>
              <view class="flag-btn" @tap="toggleFlag(question.id)">
                <text class="flag-icon" :class="{ 'flag-active': flaggedQuestions.has(question.id) }">
                  {{ flaggedQuestions.has(question.id) ? '🚩' : '⚑' }}
                </text>
              </view>
            </view>

            <text class="question-content">{{ question.content }}</text>

            <!-- Options -->
            <view class="options-list">
              <view
                v-for="(option, oIndex) in question.options"
                :key="oIndex"
                class="option-item"
                :class="{ 'option-selected': isOptionSelected(question.id, oIndex) }"
                @tap="selectOption(question.id, oIndex)"
              >
                <view class="option-index" :class="{ 'index-selected': isOptionSelected(question.id, oIndex) }">
                  <text class="option-index-text">{{ optionLabels[oIndex] }}</text>
                </view>
                <text class="option-text">{{ option.content }}</text>
              </view>
            </view>
          </view>
        </scroll-view>
      </swiper-item>
    </swiper>

    <!-- Answer Sheet at Bottom -->
    <view class="answer-sheet">
      <view class="sheet-progress">
        <text class="sheet-progress-text">{{ answeredCount }}/{{ questions.length }} answered</text>
        <view class="sheet-progress-bar">
          <view class="sheet-progress-fill" :style="{ width: answerProgress + '%' }"></view>
        </view>
      </view>
      <view class="sheet-actions">
        <view class="sheet-btn btn-prev" :class="{ 'btn-disabled': currentIndex === 0 }" @tap="goPrev">
          <text class="sheet-btn-text">Prev</text>
        </view>
        <view class="sheet-btn btn-submit" @tap="handleSubmit">
          <text class="sheet-btn-text">Submit</text>
        </view>
        <view v-if="currentIndex < questions.length - 1" class="sheet-btn btn-next" @tap="goNext">
          <text class="sheet-btn-text">Next</text>
        </view>
        <view v-else class="sheet-btn btn-next" @tap="handleSubmit">
          <text class="sheet-btn-text">Submit</text>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue';
import { getExamQuestions, submitExamAnswers, saveExamProgress } from '@/api/exam';

interface Option {
  content: string;
}

interface ExamQuestion {
  id: string;
  content: string;
  type: string;
  options: Option[];
  score: number;
}

const optionLabels = ['A', 'B', 'C', 'D', 'E', 'F', 'G', 'H'];

const examId = ref('');
const questions = ref<ExamQuestion[]>([]);
const currentIndex = ref(0);
const answers = ref<Record<string, number[]>>({});
const flaggedQuestions = ref(new Set<string>());
const remainingSeconds = ref(0);
const showGrid = ref(false);

let countdownInterval: ReturnType<typeof setInterval> | null = null;
let autoSaveInterval: ReturnType<typeof setInterval> | null = null;

const answeredCount = computed(() => Object.keys(answers.value).length);
const answerProgress = computed(() => {
  if (questions.value.length === 0) return 0;
  return Math.round((answeredCount.value / questions.value.length) * 100);
});

onMounted(() => {
  const pages = getCurrentPages();
  const currentPage = pages[pages.length - 1] as any;
  examId.value = currentPage.options?.id || '';
  loadQuestions();
  startCountdown();
  startAutoSave();
});

onUnmounted(() => {
  stopCountdown();
  stopAutoSave();
});

async function loadQuestions() {
  try {
    uni.showLoading({ title: 'Loading...' });
    const res = await getExamQuestions(examId.value);
    if (res.code === 0) {
      questions.value = res.data.questions;
      remainingSeconds.value = res.data.remainingSeconds;
      if (res.data.savedAnswers) {
        answers.value = res.data.savedAnswers;
      }
    }
  } catch (e) {
    console.error('Failed to load questions:', e);
    uni.showToast({ title: 'Failed to load', icon: 'none' });
  } finally {
    uni.hideLoading();
  }
}

function startCountdown() {
  countdownInterval = setInterval(() => {
    if (remainingSeconds.value > 0) {
      remainingSeconds.value--;
    } else {
      stopCountdown();
      handleTimeUp();
    }
  }, 1000);
}

function stopCountdown() {
  if (countdownInterval) {
    clearInterval(countdownInterval);
    countdownInterval = null;
  }
}

function startAutoSave() {
  autoSaveInterval = setInterval(() => {
    saveProgress();
  }, 30000);
}

function stopAutoSave() {
  if (autoSaveInterval) {
    clearInterval(autoSaveInterval);
    autoSaveInterval = null;
  }
}

function formatCountdown(seconds: number): string {
  const hrs = Math.floor(seconds / 3600);
  const mins = Math.floor((seconds % 3600) / 60);
  const secs = seconds % 60;
  if (hrs > 0) {
    return `${String(hrs).padStart(2, '0')}:${String(mins).padStart(2, '0')}:${String(secs).padStart(2, '0')}`;
  }
  return `${String(mins).padStart(2, '0')}:${String(secs).padStart(2, '0')}`;
}

function isOptionSelected(questionId: string, oIndex: number): boolean {
  return answers.value[questionId]?.includes(oIndex) || false;
}

function selectOption(questionId: string, oIndex: number) {
  const question = questions.value.find((q) => q.id === questionId);
  if (!question) return;

  if (!answers.value[questionId]) {
    answers.value[questionId] = [];
  }

  const idx = answers.value[questionId].indexOf(oIndex);
  if (idx > -1) {
    answers.value[questionId].splice(idx, 1);
  } else {
    if (question.type === 'single') {
      answers.value[questionId] = [oIndex];
    } else {
      answers.value[questionId].push(oIndex);
    }
  }
}

function onSwiperChange(e: any) {
  currentIndex.value = e.detail.current;
}

function goPrev() {
  if (currentIndex.value > 0) {
    currentIndex.value--;
  }
}

function goNext() {
  if (currentIndex.value < questions.value.length - 1) {
    currentIndex.value++;
  }
}

function jumpToQuestion(index: number) {
  currentIndex.value = index;
  showGrid.value = false;
}

function toggleQuestionGrid() {
  showGrid.value = !showGrid.value;
}

function toggleFlag(questionId: string) {
  if (flaggedQuestions.value.has(questionId)) {
    flaggedQuestions.value.delete(questionId);
  } else {
    flaggedQuestions.value.add(questionId);
  }
}

async function saveProgress() {
  try {
    await saveExamProgress({
      examId: examId.value,
      answers: answers.value,
      currentIndex: currentIndex.value,
    });
  } catch (e) {
    console.error('Failed to save progress:', e);
  }
}

async function handleTimeUp() {
  uni.showToast({ title: 'Time is up!', icon: 'none' });
  await submitExam();
}

async function handleSubmit() {
  const unanswered = questions.value.length - answeredCount.value;
  if (unanswered > 0) {
    const confirm = await new Promise<boolean>((resolve) => {
      uni.showModal({
        title: 'Submit Exam',
        content: `You have ${unanswered} unanswered questions. Are you sure you want to submit?`,
        success: (res) => resolve(res.confirm),
      });
    });
    if (!confirm) return;
  }
  await submitExam();
}

async function submitExam() {
  stopCountdown();
  stopAutoSave();
  try {
    uni.showLoading({ title: 'Submitting...' });
    const submitData = {
      examId: examId.value,
      answers: Object.entries(answers.value).map(([questionId, selectedOptions]) => ({
        questionId,
        selectedOptions,
      })),
    };
    const res = await submitExamAnswers(submitData);
    if (res.code === 0) {
      uni.redirectTo({ url: `/pages/exam/result?attemptId=${res.data.attemptId}` });
    }
  } catch (e) {
    console.error('Failed to submit:', e);
    uni.showToast({ title: 'Submit failed', icon: 'none' });
  } finally {
    uni.hideLoading();
  }
}

function handleBack() {
  uni.showModal({
    title: 'Leave Exam',
    content: 'Your progress has been saved. Are you sure you want to leave?',
    success: async (res) => {
      if (res.confirm) {
        await saveProgress();
        stopCountdown();
        stopAutoSave();
        uni.navigateBack();
      }
    },
  });
}
</script>

<style scoped>
.exam-answer-page {
  min-height: 100vh;
  background-color: #f5f5f5;
  display: flex;
  flex-direction: column;
}

.nav-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 88rpx;
  padding: 0 30rpx;
  background: #ffffff;
  padding-top: var(--status-bar-height, 44rpx);
  border-bottom: 1rpx solid #f0f0f0;
}

.nav-left,
.nav-right {
  width: 80rpx;
  height: 88rpx;
  display: flex;
  align-items: center;
}

.nav-back-icon {
  font-size: 36rpx;
  color: #333;
}

.nav-center {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.nav-timer {
  font-size: 36rpx;
  font-weight: bold;
  color: #f44336;
}

.nav-timer-label {
  font-size: 20rpx;
  color: #999;
}

.nav-grid-icon {
  font-size: 40rpx;
  color: #333;
  text-align: right;
}

/* Grid Drawer */
.grid-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  z-index: 100;
  display: flex;
  justify-content: flex-end;
}

.grid-drawer {
  width: 60%;
  height: 100%;
  background: #ffffff;
  display: flex;
  flex-direction: column;
}

.grid-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 30rpx;
  border-bottom: 1rpx solid #f0f0f0;
}

.grid-title {
  font-size: 30rpx;
  font-weight: 600;
  color: #333;
}

.grid-close {
  font-size: 36rpx;
  color: #999;
}

.grid-scroll {
  flex: 1;
  height: 0;
}

.grid-body {
  display: flex;
  flex-wrap: wrap;
  padding: 20rpx;
  gap: 16rpx;
}

.grid-item {
  width: calc(20% - 13rpx);
  height: 72rpx;
  border-radius: 12rpx;
  background: #f5f5f5;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 2rpx solid #e0e0e0;
}

.grid-current {
  border-color: #667eea;
  background: #e8eaf6;
}

.grid-answered {
  background: #e8f5e9;
  border-color: #4caf50;
}

.grid-flagged {
  background: #fff3e0;
  border-color: #ff9800;
}

.grid-item-text {
  font-size: 26rpx;
  color: #333;
  font-weight: 500;
}

.grid-legend {
  display: flex;
  justify-content: space-around;
  padding: 20rpx 30rpx;
  border-top: 1rpx solid #f0f0f0;
}

.legend-item {
  display: flex;
  align-items: center;
}

.legend-dot {
  width: 24rpx;
  height: 24rpx;
  border-radius: 6rpx;
  margin-right: 8rpx;
}

.legend-current {
  background: #e8eaf6;
  border: 2rpx solid #667eea;
}

.legend-answered {
  background: #e8f5e9;
  border: 2rpx solid #4caf50;
}

.legend-flagged {
  background: #fff3e0;
  border: 2rpx solid #ff9800;
}

.legend-text {
  font-size: 22rpx;
  color: #666;
}

/* Question Content */
.question-swiper {
  flex: 1;
  height: 0;
}

.question-scroll {
  height: 100%;
}

.question-card {
  padding: 30rpx;
}

.question-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24rpx;
}

.question-meta {
  display: flex;
  align-items: center;
  gap: 16rpx;
}

.question-number {
  font-size: 32rpx;
  font-weight: bold;
  color: #667eea;
}

.question-score {
  font-size: 24rpx;
  color: #999;
  background: #f0f0f0;
  padding: 4rpx 16rpx;
  border-radius: 20rpx;
}

.flag-btn {
  padding: 8rpx;
}

.flag-icon {
  font-size: 36rpx;
  color: #ccc;
}

.flag-active {
  color: #ff9800;
}

.question-content {
  font-size: 30rpx;
  color: #333;
  line-height: 1.6;
  margin-bottom: 30rpx;
  display: block;
}

.options-list {
  display: flex;
  flex-direction: column;
  gap: 16rpx;
}

.option-item {
  display: flex;
  align-items: center;
  background: #ffffff;
  border-radius: 16rpx;
  padding: 24rpx;
  border: 2rpx solid #e8e8e8;
  transition: all 0.2s ease;
}

.option-selected {
  border-color: #667eea;
  background: #f0f2ff;
}

.option-index {
  width: 52rpx;
  height: 52rpx;
  border-radius: 50%;
  background: #f0f0f0;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 20rpx;
  flex-shrink: 0;
}

.index-selected {
  background: #667eea;
}

.option-index-text {
  font-size: 26rpx;
  color: #666;
  font-weight: 600;
}

.index-selected .option-index-text {
  color: #ffffff;
}

.option-text {
  flex: 1;
  font-size: 28rpx;
  color: #333;
  line-height: 1.5;
}

/* Answer Sheet */
.answer-sheet {
  background: #ffffff;
  padding: 20rpx 30rpx;
  padding-bottom: calc(20rpx + env(safe-area-inset-bottom));
  box-shadow: 0 -2rpx 10rpx rgba(0, 0, 0, 0.05);
}

.sheet-progress {
  display: flex;
  align-items: center;
  margin-bottom: 16rpx;
}

.sheet-progress-text {
  font-size: 24rpx;
  color: #666;
  margin-right: 16rpx;
  min-width: 160rpx;
}

.sheet-progress-bar {
  flex: 1;
  height: 12rpx;
  background: #e0e0e0;
  border-radius: 6rpx;
  overflow: hidden;
}

.sheet-progress-fill {
  height: 100%;
  background: linear-gradient(90deg, #667eea, #764ba2);
  border-radius: 6rpx;
  transition: width 0.3s ease;
}

.sheet-actions {
  display: flex;
  gap: 16rpx;
}

.sheet-btn {
  flex: 1;
  height: 80rpx;
  border-radius: 40rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.btn-prev {
  background: #f0f0f0;
}

.btn-prev .sheet-btn-text {
  color: #666;
}

.btn-disabled {
  opacity: 0.5;
}

.btn-submit {
  background: #fff3e0;
}

.btn-submit .sheet-btn-text {
  color: #e65100;
  font-weight: 600;
}

.btn-next {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.btn-next .sheet-btn-text {
  color: #ffffff;
  font-weight: 600;
}

.sheet-btn-text {
  font-size: 28rpx;
}
</style>
