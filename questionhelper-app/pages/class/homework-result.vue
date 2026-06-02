<template>
  <view class="result-page">
    <!-- Score Header -->
    <view class="score-card">
      <view class="score-circle" :class="{ passed: result.score >= result.passingScore }">
        <text class="score-value">{{ result.score }}</text>
        <text class="score-unit">分</text>
      </view>
      <text class="score-status">{{ result.score >= result.passingScore ? '通过' : '未通过' }}</text>
      <text class="score-total">满分 {{ result.totalScore }} 分，及格 {{ result.passingScore }} 分</text>
    </view>

    <!-- Teacher Feedback -->
    <view v-if="result.feedback" class="section">
      <text class="section-title">教师评语</text>
      <text class="feedback-text">{{ result.feedback }}</text>
    </view>

    <!-- Correct Answers -->
    <view v-if="result.answersReleased" class="section">
      <text class="section-title">参考答案</text>
      <view v-for="(item, i) in result.answers" :key="i" class="answer-item">
        <view class="answer-header">
          <text class="answer-index">{{ i + 1 }}.</text>
          <text class="answer-status" :class="item.correct ? 'correct' : 'wrong'">
            {{ item.correct ? '正确' : '错误' }}
          </text>
        </view>
        <text class="answer-question">{{ item.question }}</text>
        <text class="answer-correct">正确答案：{{ item.correctAnswer }}</text>
        <text v-if="!item.correct" class="answer-yours">你的答案：{{ item.yourAnswer }}</text>
      </view>
    </view>

    <view v-else class="section">
      <text class="section-title">参考答案</text>
      <view class="empty-hint">
        <text class="empty-text">答案尚未发布，请等待教师批改完成</text>
      </view>
    </view>

    <!-- Back Button -->
    <view class="back-btn" @tap="goBack">
      <text class="back-text">返回班级</text>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

interface AnswerItem {
  question: string
  correctAnswer: string
  yourAnswer: string
  correct: boolean
}

interface HomeworkResult {
  score: number
  totalScore: number
  passingScore: number
  feedback: string
  answersReleased: boolean
  answers: AnswerItem[]
}

const classId = ref('')
const homeworkId = ref('')

const result = ref<HomeworkResult>({
  score: 0,
  totalScore: 100,
  passingScore: 60,
  feedback: '',
  answersReleased: false,
  answers: []
})

onMounted(() => {
  const pages = getCurrentPages()
  const page = pages[pages.length - 1] as any
  classId.value = page.options?.classId || ''
  homeworkId.value = page.options?.homeworkId || ''
  fetchResult()
})

async function fetchResult() {
  try {
    // TODO: replace with actual API call using classId, homeworkId
    result.value = {
      score: 78,
      totalScore: 100,
      passingScore: 60,
      feedback: '整体完成不错，但选择题部分还需加强，建议复习第三章内容。',
      answersReleased: true,
      answers: [
        { question: '什么是二叉树的前序遍历？', correctAnswer: '根-左-右', yourAnswer: '根-左-右', correct: true },
        { question: 'TCP三次握手的第二步是什么？', correctAnswer: 'SYN+ACK', yourAnswer: 'ACK', correct: false }
      ]
    }
  } catch (e) {
    console.error('Failed to load result', e)
  }
}

function goBack() {
  uni.navigateBack()
}
</script>

<style scoped>
.result-page {
  min-height: 100vh;
  background-color: #f5f5f5;
  padding: 20rpx 24rpx;
  padding-bottom: 80rpx;
}

.score-card {
  background-color: #fff;
  border-radius: 16rpx;
  padding: 48rpx 30rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-bottom: 24rpx;
}

.score-circle {
  width: 180rpx;
  height: 180rpx;
  border-radius: 90rpx;
  border: 8rpx solid #e74c3c;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  margin-bottom: 20rpx;
}

.score-circle.passed {
  border-color: #27ae60;
}

.score-value {
  font-size: 56rpx;
  font-weight: 700;
  color: #333;
}

.score-unit {
  font-size: 22rpx;
  color: #999;
}

.score-status {
  font-size: 30rpx;
  font-weight: 600;
  color: #e74c3c;
  margin-bottom: 8rpx;
}

.score-circle.passed + .score-status {
  color: #27ae60;
}

.score-total {
  font-size: 24rpx;
  color: #999;
}

.section {
  background-color: #fff;
  border-radius: 16rpx;
  padding: 28rpx 30rpx;
  margin-bottom: 24rpx;
}

.section-title {
  font-size: 30rpx;
  font-weight: 600;
  color: #333;
  margin-bottom: 20rpx;
}

.feedback-text {
  font-size: 28rpx;
  color: #555;
  line-height: 1.6;
}

.answer-item {
  padding: 20rpx 0;
  border-bottom: 1rpx solid #f5f5f5;
}

.answer-item:last-child {
  border-bottom: none;
}

.answer-header {
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 8rpx;
}

.answer-index {
  font-size: 28rpx;
  font-weight: 600;
  color: #333;
}

.answer-status {
  font-size: 24rpx;
  font-weight: 600;
}

.answer-status.correct {
  color: #27ae60;
}

.answer-status.wrong {
  color: #e74c3c;
}

.answer-question {
  font-size: 26rpx;
  color: #555;
  margin-bottom: 8rpx;
}

.answer-correct {
  font-size: 24rpx;
  color: #27ae60;
}

.answer-yours {
  font-size: 24rpx;
  color: #e74c3c;
  margin-top: 4rpx;
}

.empty-hint {
  padding: 40rpx 0;
  display: flex;
  align-items: center;
  justify-content: center;
}

.empty-text {
  font-size: 26rpx;
  color: #ccc;
}

.back-btn {
  background-color: #4a90d9;
  border-radius: 44rpx;
  padding: 26rpx 0;
  display: flex;
  align-items: center;
  justify-content: center;
}

.back-text {
  font-size: 30rpx;
  color: #fff;
  font-weight: 600;
}
</style>
