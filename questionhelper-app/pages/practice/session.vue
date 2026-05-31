<template>
  <view class="session-page">
    <!-- Custom Navigation Bar -->
    <view class="nav-bar">
      <view class="nav-left" @tap="handleBack">
        <text class="nav-back-icon">←</text>
      </view>
      <view class="nav-center">
        <text class="nav-title">{{ modeName }}</text>
      </view>
      <view class="nav-right" @tap="toggleBookmark">
        <text class="nav-bookmark">{{ isBookmarked ? '★' : '☆' }}</text>
      </view>
    </view>

    <!-- Timer and Progress -->
    <view class="timer-bar">
      <text class="timer-text">{{ formatTime(timerSeconds) }}</text>
      <view class="progress-bar">
        <view class="progress-fill" :style="{ width: progressPercent + '%' }"></view>
      </view>
      <text class="progress-text">{{ currentIndex + 1 }}/{{ questions.length }}</text>
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
              <text class="question-number">Q{{ qIndex + 1 }}</text>
              <text class="question-type">{{ question.type }}</text>
            </view>
            <text class="question-content">{{ question.content }}</text>

            <!-- Options -->
            <view class="options-list">
              <view
                v-for="(option, oIndex) in question.options"
                :key="oIndex"
                class="option-item"
                :class="{
                  'option-selected': isOptionSelected(qIndex, oIndex),
                  'option-correct': showAnswer && option.isCorrect,
                  'option-wrong': showAnswer && isOptionSelected(qIndex, oIndex) && !option.isCorrect,
                }"
                @tap="selectOption(qIndex, oIndex)"
              >
                <view class="option-index">
                  <text class="option-index-text">{{ optionLabels[oIndex] }}</text>
                </view>
                <text class="option-text">{{ option.content }}</text>
                <view v-if="showAnswer && option.isCorrect" class="option-icon">
                  <text class="icon-correct">✓</text>
                </view>
                <view v-if="showAnswer && isOptionSelected(qIndex, oIndex) && !option.isCorrect" class="option-icon">
                  <text class="icon-wrong">✗</text>
                </view>
              </view>
            </view>

            <!-- Answer Explanation -->
            <view v-if="showAnswer && question.explanation" class="explanation-card">
              <text class="explanation-title">Explanation</text>
              <text class="explanation-content">{{ question.explanation }}</text>
            </view>
          </view>
        </scroll-view>
      </swiper-item>
    </swiper>

    <!-- Bottom Actions -->
    <view class="bottom-bar">
      <view class="action-btn btn-prev" :class="{ 'btn-disabled': currentIndex === 0 }" @tap="goPrev">
        <text class="btn-text">Previous</text>
      </view>
      <view class="action-btn btn-show-answer" @tap="toggleShowAnswer">
        <text class="btn-text">{{ showAnswer ? 'Hide Answer' : 'Show Answer' }}</text>
      </view>
      <view v-if="currentIndex < questions.length - 1" class="action-btn btn-next" @tap="goNext">
        <text class="btn-text">Next</text>
      </view>
      <view v-else class="action-btn btn-submit" @tap="handleSubmit">
        <text class="btn-text">Submit All</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue';
import { getPracticeQuestions, submitPracticeAnswers, toggleQuestionBookmark } from '@/api/practice';

interface Option {
  content: string;
  isCorrect: boolean;
}

interface Question {
  id: string;
  content: string;
  type: string;
  options: Option[];
  explanation?: string;
}

const mode = ref('random');
const modeName = ref('Practice');
const questions = ref<Question[]>([]);
const currentIndex = ref(0);
const answers = ref<Record<number, number[]>>({});
const showAnswer = ref(false);
const isBookmarked = ref(false);
const timerSeconds = ref(0);
let timerInterval: ReturnType<typeof setInterval> | null = null;

const optionLabels = ['A', 'B', 'C', 'D', 'E', 'F', 'G', 'H'];

const progressPercent = computed(() => {
  if (questions.value.length === 0) return 0;
  return Math.round(((currentIndex.value + 1) / questions.value.length) * 100);
});

onMounted(() => {
  const pages = getCurrentPages();
  const currentPage = pages[pages.length - 1] as any;
  mode.value = currentPage.options?.mode || 'random';
  modeName.value = getModeDisplayName(mode.value);
  loadQuestions();
  startTimer();
});

onUnmounted(() => {
  stopTimer();
});

function getModeDisplayName(m: string): string {
  const map: Record<string, string> = {
    random: 'Random Practice',
    sequential: 'Sequential Practice',
    wrong: 'Wrong Questions',
    favorites: 'Favorites',
    mock: 'Mock Exam',
    challenge: 'Challenge',
    timed: 'Timed Practice',
  };
  return map[m] || 'Practice';
}

async function loadQuestions() {
  try {
    uni.showLoading({ title: 'Loading...' });
    const res = await getPracticeQuestions({ mode: mode.value });
    if (res.code === 0) {
      questions.value = res.data;
    }
  } catch (e) {
    console.error('Failed to load questions:', e);
    uni.showToast({ title: 'Failed to load', icon: 'none' });
  } finally {
    uni.hideLoading();
  }
}

function startTimer() {
  timerInterval = setInterval(() => {
    timerSeconds.value++;
  }, 1000);
}

function stopTimer() {
  if (timerInterval) {
    clearInterval(timerInterval);
    timerInterval = null;
  }
}

function formatTime(seconds: number): string {
  const hrs = Math.floor(seconds / 3600);
  const mins = Math.floor((seconds % 3600) / 60);
  const secs = seconds % 60;
  if (hrs > 0) {
    return `${String(hrs).padStart(2, '0')}:${String(mins).padStart(2, '0')}:${String(secs).padStart(2, '0')}`;
  }
  return `${String(mins).padStart(2, '0')}:${String(secs).padStart(2, '0')}`;
}

function isOptionSelected(qIndex: number, oIndex: number): boolean {
  return answers.value[qIndex]?.includes(oIndex) || false;
}

function selectOption(qIndex: number, oIndex: number) {
  if (showAnswer.value) return;
  const question = questions.value[qIndex];
  if (!answers.value[qIndex]) {
    answers.value[qIndex] = [];
  }
  const idx = answers.value[qIndex].indexOf(oIndex);
  if (idx > -1) {
    answers.value[qIndex].splice(idx, 1);
  } else {
    if (question.type === 'single') {
      answers.value[qIndex] = [oIndex];
    } else {
      answers.value[qIndex].push(oIndex);
    }
  }
}

function onSwiperChange(e: any) {
  currentIndex.value = e.detail.current;
  showAnswer.value = false;
  checkBookmark();
}

function goPrev() {
  if (currentIndex.value > 0) {
    currentIndex.value--;
    showAnswer.value = false;
    checkBookmark();
  }
}

function goNext() {
  if (currentIndex.value < questions.value.length - 1) {
    currentIndex.value++;
    showAnswer.value = false;
    checkBookmark();
  }
}

function toggleShowAnswer() {
  showAnswer.value = !showAnswer.value;
}

async function toggleBookmark() {
  const questionId = questions.value[currentIndex.value]?.id;
  if (!questionId) return;
  try {
    const res = await toggleQuestionBookmark(questionId);
    if (res.code === 0) {
      isBookmarked.value = !isBookmarked.value;
      uni.showToast({
        title: isBookmarked.value ? 'Bookmarked' : 'Unbookmarked',
        icon: 'none',
      });
    }
  } catch (e) {
    console.error('Failed to toggle bookmark:', e);
  }
}

function checkBookmark() {
  const question = questions.value[currentIndex.value];
  if (question) {
    isBookmarked.value = (question as any).isBookmarked || false;
  }
}

async function handleSubmit() {
  const unanswered = questions.value.length - Object.keys(answers.value).length;
  if (unanswered > 0) {
    const confirm = await new Promise<boolean>((resolve) => {
      uni.showModal({
        title: 'Confirm Submit',
        content: `You have ${unanswered} unanswered questions. Submit anyway?`,
        success: (res) => resolve(res.confirm),
      });
    });
    if (!confirm) return;
  }

  stopTimer();
  try {
    uni.showLoading({ title: 'Submitting...' });
    const submitData = {
      mode: mode.value,
      duration: timerSeconds.value,
      answers: Object.entries(answers.value).map(([qIndex, optionIndexes]) => ({
        questionId: questions.value[Number(qIndex)].id,
        selectedOptions: optionIndexes,
      })),
    };
    const res = await submitPracticeAnswers(submitData);
    if (res.code === 0) {
      uni.redirectTo({ url: `/pages/practice/result?recordId=${res.data.recordId}` });
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
    title: 'Leave Practice',
    content: 'Your progress will be lost. Are you sure?',
    success: (res) => {
      if (res.confirm) {
        stopTimer();
        uni.navigateBack();
      }
    },
  });
}
</script>

<style scoped>
.session-page {
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

.nav-title {
  font-size: 32rpx;
  font-weight: 600;
  color: #333;
}

.nav-bookmark {
  font-size: 40rpx;
  color: #ff9800;
  text-align: right;
}

.timer-bar {
  display: flex;
  align-items: center;
  padding: 16rpx 30rpx;
  background: #ffffff;
  border-bottom: 1rpx solid #f0f0f0;
}

.timer-text {
  font-size: 28rpx;
  color: #667eea;
  font-weight: 600;
  min-width: 120rpx;
}

.progress-bar {
  flex: 1;
  height: 12rpx;
  background: #e0e0e0;
  border-radius: 6rpx;
  margin: 0 20rpx;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: linear-gradient(90deg, #667eea, #764ba2);
  border-radius: 6rpx;
  transition: width 0.3s ease;
}

.progress-text {
  font-size: 24rpx;
  color: #999;
  min-width: 80rpx;
  text-align: right;
}

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
  align-items: center;
  margin-bottom: 20rpx;
}

.question-number {
  font-size: 28rpx;
  font-weight: bold;
  color: #667eea;
  margin-right: 16rpx;
}

.question-type {
  font-size: 22rpx;
  color: #999;
  background: #f0f0f0;
  padding: 4rpx 16rpx;
  border-radius: 20rpx;
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

.option-correct {
  border-color: #4caf50;
  background: #e8f5e9;
}

.option-wrong {
  border-color: #f44336;
  background: #ffebee;
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

.option-selected .option-index {
  background: #667eea;
}

.option-correct .option-index {
  background: #4caf50;
}

.option-wrong .option-index {
  background: #f44336;
}

.option-index-text {
  font-size: 26rpx;
  color: #666;
  font-weight: 600;
}

.option-selected .option-index-text,
.option-correct .option-index-text,
.option-wrong .option-index-text {
  color: #ffffff;
}

.option-text {
  flex: 1;
  font-size: 28rpx;
  color: #333;
  line-height: 1.5;
}

.option-icon {
  margin-left: 16rpx;
}

.icon-correct {
  font-size: 32rpx;
  color: #4caf50;
  font-weight: bold;
}

.icon-wrong {
  font-size: 32rpx;
  color: #f44336;
  font-weight: bold;
}

.explanation-card {
  margin-top: 30rpx;
  background: #fff8e1;
  border-radius: 16rpx;
  padding: 24rpx;
  border-left: 6rpx solid #ff9800;
}

.explanation-title {
  font-size: 26rpx;
  font-weight: 600;
  color: #e65100;
  display: block;
  margin-bottom: 12rpx;
}

.explanation-content {
  font-size: 26rpx;
  color: #666;
  line-height: 1.6;
  display: block;
}

.bottom-bar {
  display: flex;
  padding: 20rpx 30rpx;
  padding-bottom: calc(20rpx + env(safe-area-inset-bottom));
  background: #ffffff;
  gap: 16rpx;
  box-shadow: 0 -2rpx 10rpx rgba(0, 0, 0, 0.05);
}

.action-btn {
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

.btn-prev .btn-text {
  color: #666;
}

.btn-disabled {
  opacity: 0.5;
}

.btn-show-answer {
  background: #fff3e0;
}

.btn-show-answer .btn-text {
  color: #e65100;
}

.btn-next {
  background: #667eea;
}

.btn-next .btn-text {
  color: #ffffff;
}

.btn-submit {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.btn-submit .btn-text {
  color: #ffffff;
  font-weight: 600;
}

.btn-text {
  font-size: 28rpx;
}
</style>
