<template>
  <view class="page">
    <!-- Loading -->
    <view v-if="loading" class="loading-wrap">
      <text class="loading-text">加载中...</text>
    </view>

    <template v-else-if="result">
      <!-- Score Header -->
      <view :class="['score-section', hasScore ? (result.passed ? 'section-pass' : 'section-fail') : 'section-pending']">
        <view class="score-circle">
          <template v-if="hasScore">
            <text class="score-value">{{ result.score }}</text>
            <text class="score-unit">分</text>
          </template>
          <template v-else>
            <text class="score-pending-text">待批改</text>
          </template>
        </view>
        <view v-if="hasScore" class="score-info">
          <text :class="['pass-text', result.passed ? 'text-pass' : 'text-fail']">
            {{ result.passed ? '已通过' : '未通过' }}
          </text>
          <text class="score-total">满分 {{ result.totalScore }} 分</text>
        </view>
        <text v-else class="score-hint">作业已提交，等待老师批改</text>
      </view>

      <!-- Stats Row -->
      <view v-if="hasScore" class="stats-row">
        <view class="stat-item">
          <text class="stat-value">{{ result.score }}</text>
          <text class="stat-label">得分</text>
        </view>
        <view class="stat-divider"></view>
        <view class="stat-item">
          <text class="stat-value">{{ result.totalScore }}</text>
          <text class="stat-label">满分</text>
        </view>
        <view class="stat-divider"></view>
        <view class="stat-item">
          <text class="stat-value">{{ result.correctRate }}%</text>
          <text class="stat-label">正确率</text>
        </view>
        <view class="stat-divider"></view>
        <view class="stat-item">
          <text class="stat-value">{{ result.rank || '-' }}</text>
          <text class="stat-label">排名</text>
        </view>
      </view>

      <!-- Submission Info -->
      <view class="info-card">
        <text class="info-card-title">提交信息</text>
        <view class="info-row">
          <text class="info-label">作业标题</text>
          <text class="info-value">{{ result.homeworkTitle }}</text>
        </view>
        <view class="info-row">
          <text class="info-label">提交时间</text>
          <text class="info-value">{{ result.submitTime }}</text>
        </view>
        <view v-if="result.gradeTime" class="info-row">
          <text class="info-label">批改时间</text>
          <text class="info-value">{{ result.gradeTime }}</text>
        </view>
        <view v-if="result.graderName" class="info-row">
          <text class="info-label">批改人</text>
          <text class="info-value">{{ result.graderName }}</text>
        </view>
      </view>

      <!-- Teacher Feedback -->
      <view v-if="result.feedback" class="feedback-card">
        <text class="feedback-title">教师评语</text>
        <view class="feedback-body">
          <text class="feedback-text">{{ result.feedback }}</text>
        </view>
      </view>

      <!-- My Answer -->
      <view class="answer-card">
        <text class="answer-card-title">我的作答</text>
        <text v-if="result.content" class="answer-content">{{ result.content }}</text>
        <text v-else class="answer-empty">未填写文字答案</text>

        <!-- My Attachments -->
        <view v-if="result.attachments && result.attachments.length > 0" class="my-attachments">
          <text class="attachment-section-title">提交的附件</text>
          <view
            v-for="(att, index) in result.attachments"
            :key="index"
            class="attachment-item"
            @tap="onPreviewFile(att)"
          >
            <text class="attachment-icon">&#xe628;</text>
            <text class="attachment-name">{{ att.name || '附件' + (index + 1) }}</text>
            <text class="attachment-action">查看</text>
          </view>
        </view>
      </view>

      <!-- Correct Answers (if released) -->
      <view v-if="result.answerReleased && result.questions && result.questions.length > 0" class="correct-section">
        <text class="correct-section-title">参考答案</text>
        <view
          v-for="(q, index) in result.questions"
          :key="q.id"
          class="correct-card"
        >
          <view class="correct-header">
            <view class="correct-num-wrap">
              <text class="correct-num">{{ index + 1 }}</text>
            </view>
            <text class="correct-title">{{ q.content }}</text>
          </view>

          <!-- My Answer vs Correct Answer -->
          <view v-if="hasScore" class="answer-compare">
            <view class="compare-row">
              <text class="compare-label">我的答案</text>
              <text :class="['compare-value', q.isCorrect ? 'value-correct' : 'value-wrong']">
                {{ q.userAnswer || '未作答' }}
              </text>
            </view>
            <view class="compare-row">
              <text class="compare-label">参考答案</text>
              <text class="compare-value value-correct">{{ q.correctAnswer }}</text>
            </view>
            <view class="compare-row">
              <text class="compare-label">得分</text>
              <text class="compare-value score-text">{{ q.score }}/{{ q.totalScore }}</text>
            </view>
          </view>

          <!-- Explanation -->
          <view v-if="q.explanation" class="explanation-wrap">
            <text class="explanation-label">解析</text>
            <text class="explanation-text">{{ q.explanation }}</text>
          </view>
        </view>
      </view>
    </template>

    <!-- Bottom Actions -->
    <view class="bottom-bar">
      <view class="action-btn btn-back" @tap="onGoBack">
        <text class="btn-text">返回作业列表</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { getHomeworkResult } from '@/api/class'

interface QuestionResult {
  id: string
  content: string
  userAnswer: string
  correctAnswer: string
  isCorrect: boolean
  score: number
  totalScore: number
  explanation?: string
}

interface HomeworkResult {
  homeworkTitle: string
  score: number
  totalScore: number
  passed: boolean
  correctRate: number
  rank?: number
  submitTime: string
  gradeTime?: string
  graderName?: string
  feedback?: string
  content?: string
  attachments?: { name: string; url: string }[]
  answerReleased: boolean
  questions?: QuestionResult[]
}

const classId = ref('')
const homeworkId = ref('')
const loading = ref(true)
const result = ref<HomeworkResult | null>(null)

const hasScore = computed(() => {
  return result.value && result.value.score !== undefined && result.value.score !== null && result.value.gradeTime
})

// ---------- Data Loading ----------

async function loadResult() {
  loading.value = true
  try {
    const res = await getHomeworkResult(classId.value, homeworkId.value)
    result.value = res.data
  } catch (e) {
    uni.showToast({ title: '加载失败', icon: 'none' })
  } finally {
    loading.value = false
  }
}

// ---------- File Preview ----------

function onPreviewFile(att: { name: string; url: string }) {
  if (att.url) {
    uni.downloadFile({
      url: att.url,
      success: (res) => {
        if (res.statusCode === 200) {
          uni.openDocument({
            filePath: res.tempFilePath,
            showMenu: true
          })
        }
      }
    })
  }
}

// ---------- Navigation ----------

function onGoBack() {
  uni.navigateBack()
}

// ---------- Lifecycle ----------

onMounted(() => {
  const pages = getCurrentPages()
  const currentPage = pages[pages.length - 1] as any
  classId.value = currentPage?.options?.classId || ''
  homeworkId.value = currentPage?.options?.homeworkId || ''

  if (!classId.value || !homeworkId.value) {
    uni.showToast({ title: '参数错误', icon: 'none' })
    setTimeout(() => uni.navigateBack(), 1500)
    return
  }

  loadResult()
})
</script>

<style scoped>
.page {
  min-height: 100vh;
  background-color: #f5f6fa;
  padding-bottom: 140rpx;
}

.loading-wrap {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 400rpx;
}

.loading-text {
  font-size: 28rpx;
  color: #999;
}

/* ---- Score Section ---- */

.score-section {
  padding: 60rpx 40rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.section-pass {
  background: linear-gradient(135deg, #52c41a 0%, #95de64 100%);
}

.section-fail {
  background: linear-gradient(135deg, #ff4d4f 0%, #ff7875 100%);
}

.section-pending {
  background: linear-gradient(135deg, #1677ff 0%, #69b1ff 100%);
}

.score-circle {
  width: 200rpx;
  height: 200rpx;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.2);
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 24rpx;
}

.score-value {
  font-size: 72rpx;
  font-weight: bold;
  color: #fff;
}

.score-unit {
  font-size: 28rpx;
  color: rgba(255, 255, 255, 0.8);
  margin-left: 4rpx;
}

.score-pending-text {
  font-size: 32rpx;
  font-weight: 600;
  color: #fff;
}

.score-info {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8rpx;
}

.pass-text {
  font-size: 36rpx;
  font-weight: bold;
}

.text-pass {
  color: #fff;
}

.text-fail {
  color: #fff;
}

.score-total {
  font-size: 26rpx;
  color: rgba(255, 255, 255, 0.8);
}

.score-hint {
  font-size: 28rpx;
  color: rgba(255, 255, 255, 0.9);
}

/* ---- Stats Row ---- */

.stats-row {
  display: flex;
  align-items: center;
  justify-content: space-around;
  background-color: #fff;
  margin: -30rpx 30rpx 0;
  border-radius: 16rpx;
  padding: 28rpx 0;
  box-shadow: 0 4rpx 20rpx rgba(0, 0, 0, 0.08);
  position: relative;
  z-index: 1;
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8rpx;
}

.stat-value {
  font-size: 36rpx;
  font-weight: 700;
  color: #333;
}

.stat-label {
  font-size: 24rpx;
  color: #999;
}

.stat-divider {
  width: 1rpx;
  height: 48rpx;
  background-color: #e8e8e8;
}

/* ---- Info Card ---- */

.info-card {
  background-color: #fff;
  margin: 20rpx 24rpx;
  border-radius: 16rpx;
  padding: 28rpx;
  box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.04);
}

.info-card-title {
  font-size: 28rpx;
  font-weight: 600;
  color: #333;
  margin-bottom: 16rpx;
  display: block;
}

.info-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12rpx 0;
  border-bottom: 1rpx solid #f5f5f5;
}

.info-row:last-child {
  border-bottom: none;
}

.info-label {
  font-size: 26rpx;
  color: #999;
}

.info-value {
  font-size: 26rpx;
  color: #333;
}

/* ---- Feedback ---- */

.feedback-card {
  background-color: #fff;
  margin: 0 24rpx 20rpx;
  border-radius: 16rpx;
  padding: 28rpx;
  box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.04);
}

.feedback-title {
  font-size: 28rpx;
  font-weight: 600;
  color: #333;
  margin-bottom: 16rpx;
  display: block;
}

.feedback-body {
  background-color: #fffbe6;
  border-radius: 12rpx;
  padding: 20rpx;
  border-left: 6rpx solid #faad14;
}

.feedback-text {
  font-size: 28rpx;
  color: #333;
  line-height: 1.8;
}

/* ---- My Answer ---- */

.answer-card {
  background-color: #fff;
  margin: 0 24rpx 20rpx;
  border-radius: 16rpx;
  padding: 28rpx;
  box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.04);
}

.answer-card-title {
  font-size: 28rpx;
  font-weight: 600;
  color: #333;
  margin-bottom: 16rpx;
  display: block;
}

.answer-content {
  font-size: 28rpx;
  color: #333;
  line-height: 1.8;
  background-color: #f8f9fc;
  border-radius: 12rpx;
  padding: 20rpx;
  display: block;
}

.answer-empty {
  font-size: 26rpx;
  color: #ccc;
  display: block;
}

.my-attachments {
  margin-top: 20rpx;
}

.attachment-section-title {
  font-size: 26rpx;
  font-weight: 600;
  color: #666;
  margin-bottom: 12rpx;
  display: block;
}

.attachment-item {
  display: flex;
  align-items: center;
  padding: 16rpx 0;
  border-bottom: 1rpx solid #f5f5f5;
}

.attachment-item:last-child {
  border-bottom: none;
}

.attachment-icon {
  font-size: 28rpx;
  color: #1677ff;
  margin-right: 12rpx;
}

.attachment-name {
  flex: 1;
  font-size: 26rpx;
  color: #333;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.attachment-action {
  font-size: 24rpx;
  color: #1677ff;
  margin-left: 12rpx;
}

/* ---- Correct Answers ---- */

.correct-section {
  padding: 0 24rpx;
}

.correct-section-title {
  font-size: 30rpx;
  font-weight: 600;
  color: #333;
  margin-bottom: 16rpx;
  display: block;
}

.correct-card {
  background-color: #fff;
  border-radius: 16rpx;
  padding: 24rpx;
  margin-bottom: 16rpx;
  box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.04);
}

.correct-header {
  display: flex;
  align-items: flex-start;
  margin-bottom: 16rpx;
}

.correct-num-wrap {
  width: 40rpx;
  height: 40rpx;
  border-radius: 50%;
  background-color: #1677ff;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 12rpx;
  flex-shrink: 0;
}

.correct-num {
  font-size: 22rpx;
  font-weight: 600;
  color: #fff;
}

.correct-title {
  flex: 1;
  font-size: 28rpx;
  color: #333;
  line-height: 1.6;
}

.answer-compare {
  background-color: #f8f9fc;
  border-radius: 12rpx;
  padding: 20rpx;
  margin-bottom: 16rpx;
}

.compare-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8rpx 0;
}

.compare-label {
  font-size: 26rpx;
  color: #999;
  min-width: 120rpx;
}

.compare-value {
  flex: 1;
  font-size: 26rpx;
  font-weight: 600;
  text-align: right;
}

.value-correct {
  color: #52c41a;
}

.value-wrong {
  color: #ff4d4f;
}

.score-text {
  color: #1677ff;
}

.explanation-wrap {
  background-color: #fffbe6;
  border-radius: 12rpx;
  padding: 16rpx 20rpx;
  border-left: 6rpx solid #faad14;
}

.explanation-label {
  font-size: 24rpx;
  color: #faad14;
  font-weight: 600;
  display: block;
  margin-bottom: 8rpx;
}

.explanation-text {
  font-size: 26rpx;
  color: #666;
  line-height: 1.6;
}

/* ---- Bottom Bar ---- */

.bottom-bar {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  padding: 20rpx 32rpx;
  padding-bottom: calc(20rpx + env(safe-area-inset-bottom));
  background-color: #fff;
  box-shadow: 0 -2rpx 12rpx rgba(0, 0, 0, 0.06);
}

.action-btn {
  height: 88rpx;
  border-radius: 44rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.btn-back {
  background-color: #1677ff;
}

.btn-text {
  font-size: 30rpx;
  font-weight: 600;
  color: #fff;
}
</style>
