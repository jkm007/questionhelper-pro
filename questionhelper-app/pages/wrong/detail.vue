<template>
  <view class="page">
    <!-- Loading Skeleton -->
    <view v-if="loading" class="skeleton">
      <view class="skeleton-title"></view>
      <view class="skeleton-line"></view>
      <view class="skeleton-line short"></view>
    </view>

    <template v-else-if="detail">
      <!-- Question Content -->
      <view class="section question-section">
        <view class="section-header">
          <text class="section-title">题目内容</text>
          <view class="type-tag" :class="'type-' + detail.questionType">
            <text class="type-tag-text">{{ typeLabels[detail.questionType] }}</text>
          </view>
        </view>
        <text class="question-content">{{ detail.content }}</text>
        <!-- Options if choice question -->
        <view v-if="detail.options && detail.options.length" class="options-list">
          <view
            v-for="(opt, idx) in detail.options"
            :key="idx"
            :class="['option-item', getOptionClass(opt.key)]"
          >
            <text class="option-key">{{ opt.key }}</text>
            <text class="option-value">{{ opt.value }}</text>
          </view>
        </view>
      </view>

      <!-- Answer Comparison -->
      <view class="section comparison-section">
        <text class="section-title">答案对比</text>
        <view class="compare-cards">
          <view class="compare-card wrong-card">
            <text class="compare-label">我的答案</text>
            <text class="compare-answer wrong-answer">{{ detail.myAnswer || '未作答' }}</text>
          </view>
          <view class="compare-arrow">
            <text class="arrow-icon">&#xe61e;</text>
          </view>
          <view class="compare-card correct-card">
            <text class="compare-label">正确答案</text>
            <text class="compare-answer correct-answer">{{ detail.correctAnswer }}</text>
          </view>
        </view>
      </view>

      <!-- Analysis -->
      <view class="section analysis-section">
        <text class="section-title">题目解析</text>
        <text class="analysis-content">{{ detail.analysis }}</text>
        <view v-if="detail.knowledgePoints && detail.knowledgePoints.length" class="knowledge-tags">
          <text class="knowledge-label">知识点:</text>
          <view v-for="(kp, idx) in detail.knowledgePoints" :key="idx" class="knowledge-tag">
            <text class="knowledge-tag-text">{{ kp }}</text>
          </view>
        </view>
      </view>

      <!-- Wrong History -->
      <view class="section history-section">
        <text class="section-title">错误历史</text>
        <view v-for="(record, idx) in detail.wrongHistory" :key="idx" class="history-item">
          <view class="history-dot"></view>
          <view class="history-content">
            <text class="history-date">{{ record.date }}</text>
            <text class="history-answer">我的答案: {{ record.myAnswer }}</text>
          </view>
        </view>
        <view v-if="!detail.wrongHistory || detail.wrongHistory.length === 0" class="no-history">
          <text class="no-history-text">暂无历史记录</text>
        </view>
      </view>

      <!-- Mastery Indicator -->
      <view class="section mastery-section">
        <text class="section-title">掌握程度</text>
        <view class="mastery-indicator">
          <view class="mastery-visual">
            <view class="mastery-ring">
              <text class="mastery-percent">{{ detail.masteryLevel }}%</text>
            </view>
          </view>
          <view class="mastery-info">
            <text class="mastery-status" :class="getMasteryStatusClass(detail.masteryLevel)">
              {{ getMasteryStatusText(detail.masteryLevel) }}
            </text>
            <text class="mastery-hint">{{ getMasteryHint(detail.masteryLevel) }}</text>
          </view>
        </view>
        <view class="mastery-bar-wrap">
          <view class="mastery-bar-bg">
            <view
              class="mastery-bar-fill"
              :style="{ width: detail.masteryLevel + '%' }"
              :class="getMasteryColor(detail.masteryLevel)"
            ></view>
          </view>
        </view>
      </view>

      <!-- Notes -->
      <view class="section notes-section">
        <text class="section-title">学习笔记</text>
        <textarea
          class="notes-input"
          v-model="notes"
          placeholder="记录你的学习心得..."
          :maxlength="500"
          auto-height
        />
        <text class="notes-count">{{ notes.length }}/500</text>
      </view>
    </template>

    <!-- Bottom Actions -->
    <view v-if="detail" class="bottom-bar">
      <view class="action-btn review-btn" @tap="onReviewAgain">
        <text class="btn-icon">&#xe61d;</text>
        <text class="btn-text">再做一次</text>
      </view>
      <view class="action-btn master-btn" @tap="onMarkMastered">
        <text class="btn-icon">&#xe623;</text>
        <text class="btn-text">标记已掌握</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getWrongDetail, updateMastery, updateWrongNotes } from '@/api/wrong'

interface WrongRecord {
  date: string
  myAnswer: string
}

interface WrongDetail {
  id: string
  content: string
  questionType: string
  options: { key: string; value: string }[]
  myAnswer: string
  correctAnswer: string
  analysis: string
  knowledgePoints: string[]
  wrongHistory: WrongRecord[]
  masteryLevel: number
}

const loading = ref(true)
const detail = ref<WrongDetail | null>(null)
const notes = ref('')
const questionId = ref('')

const typeLabels: Record<string, string> = {
  'single': '单选题',
  'multiple': '多选题',
  'judge': '判断题',
  'fill': '填空题',
  'essay': '简答题'
}

function getOptionClass(key: string): string {
  if (!detail.value) return ''
  const classes = []
  if (detail.value.myAnswer === key) classes.push('selected-wrong')
  if (detail.value.correctAnswer === key) classes.push('selected-correct')
  return classes.join(' ')
}

function getMasteryColor(level: number): string {
  if (level < 30) return 'mastery-low'
  if (level < 70) return 'mastery-medium'
  return 'mastery-high'
}

function getMasteryStatusClass(level: number): string {
  if (level < 30) return 'status-low'
  if (level < 70) return 'status-medium'
  return 'status-high'
}

function getMasteryStatusText(level: number): string {
  if (level < 30) return '未掌握'
  if (level < 70) return '初步掌握'
  if (level < 100) return '基本掌握'
  return '已完全掌握'
}

function getMasteryHint(level: number): string {
  if (level < 30) return '建议加强复习，重点理解解题思路'
  if (level < 70) return '继续练习，巩固知识点'
  if (level < 100) return '快要掌握了，再巩固一下'
  return '恭喜！已经完全掌握这道题'
}

async function loadDetail() {
  loading.value = true
  try {
    const res = await getWrongDetail(questionId.value)
    detail.value = res.data
    notes.value = res.data?.notes || ''
  } catch (e) {
    uni.showToast({ title: '加载失败', icon: 'none' })
  } finally {
    loading.value = false
  }
}

async function onMarkMastered() {
  if (!detail.value) return
  uni.showModal({
    title: '确认',
    content: '确定将此题标记为已掌握吗？',
    success: async (res) => {
      if (res.confirm) {
        try {
          await updateMastery(detail.value!.id, 100)
          detail.value!.masteryLevel = 100
          uni.showToast({ title: '已标记为已掌握', icon: 'success' })
          saveNotes()
        } catch (e) {
          uni.showToast({ title: '操作失败', icon: 'none' })
        }
      }
    }
  })
}

function onReviewAgain() {
  if (!detail.value) return
  saveNotes()
  uni.navigateTo({
    url: `/pages/exam/practice?wrongId=${detail.value.id}`
  })
}

async function saveNotes() {
  if (!detail.value || !notes.value.trim()) return
  try {
    await updateWrongNotes(detail.value.id, notes.value)
  } catch (e) {
    // silent fail for notes save
  }
}

onMounted(() => {
  const pages = getCurrentPages()
  const currentPage = pages[pages.length - 1] as any
  questionId.value = currentPage?.options?.id || ''
  if (questionId.value) {
    loadDetail()
  } else {
    uni.showToast({ title: '参数错误', icon: 'none' })
    setTimeout(() => uni.navigateBack(), 1500)
  }
})
</script>

<style scoped>
.page {
  min-height: 100vh;
  background-color: #f5f6fa;
  padding-bottom: 140rpx;
}

.skeleton {
  padding: 32rpx 24rpx;
}

.skeleton-title {
  height: 40rpx;
  width: 60%;
  background-color: #e8e8e8;
  border-radius: 8rpx;
  margin-bottom: 24rpx;
}

.skeleton-line {
  height: 28rpx;
  width: 100%;
  background-color: #e8e8e8;
  border-radius: 8rpx;
  margin-bottom: 16rpx;
}

.skeleton-line.short {
  width: 40%;
}

.section {
  background-color: #fff;
  margin: 20rpx 24rpx;
  border-radius: 16rpx;
  padding: 28rpx;
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 20rpx;
}

.section-title {
  font-size: 30rpx;
  font-weight: 600;
  color: #333;
  margin-bottom: 20rpx;
  display: block;
}

.section-header .section-title {
  margin-bottom: 0;
}

.type-tag {
  padding: 6rpx 16rpx;
  border-radius: 8rpx;
}

.type-tag-text {
  font-size: 22rpx;
}

.type-single { background-color: #e8f3ff; }
.type-single .type-tag-text { color: #1677ff; }
.type-multiple { background-color: #fff0e6; }
.type-multiple .type-tag-text { color: #fa8c16; }
.type-judge { background-color: #f6ffed; }
.type-judge .type-tag-text { color: #52c41a; }
.type-fill { background-color: #f9f0ff; }
.type-fill .type-tag-text { color: #722ed1; }
.type-essay { background-color: #fff1f0; }
.type-essay .type-tag-text { color: #ff4d4f; }

.question-content {
  font-size: 30rpx;
  color: #333;
  line-height: 1.8;
}

.options-list {
  margin-top: 24rpx;
  display: flex;
  flex-direction: column;
  gap: 16rpx;
}

.option-item {
  display: flex;
  align-items: flex-start;
  padding: 20rpx;
  border: 2rpx solid #e8e8e8;
  border-radius: 12rpx;
  gap: 16rpx;
}

.option-item.selected-wrong {
  border-color: #ff4d4f;
  background-color: #fff1f0;
}

.option-item.selected-correct {
  border-color: #52c41a;
  background-color: #f6ffed;
}

.option-key {
  font-size: 28rpx;
  font-weight: 600;
  color: #666;
  flex-shrink: 0;
  width: 40rpx;
}

.option-item.selected-wrong .option-key { color: #ff4d4f; }
.option-item.selected-correct .option-key { color: #52c41a; }

.option-value {
  font-size: 28rpx;
  color: #333;
  line-height: 1.6;
}

.compare-cards {
  display: flex;
  align-items: center;
  gap: 16rpx;
}

.compare-card {
  flex: 1;
  padding: 24rpx;
  border-radius: 12rpx;
  text-align: center;
}

.wrong-card {
  background-color: #fff1f0;
  border: 2rpx solid #ffccc7;
}

.correct-card {
  background-color: #f6ffed;
  border: 2rpx solid #b7eb8f;
}

.compare-label {
  font-size: 22rpx;
  color: #999;
  display: block;
  margin-bottom: 12rpx;
}

.compare-answer {
  font-size: 32rpx;
  font-weight: 600;
}

.wrong-answer { color: #ff4d4f; }
.correct-answer { color: #52c41a; }

.compare-arrow {
  flex-shrink: 0;
}

.arrow-icon {
  font-size: 32rpx;
  color: #999;
}

.analysis-content {
  font-size: 28rpx;
  color: #555;
  line-height: 1.8;
}

.knowledge-tags {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 12rpx;
  margin-top: 20rpx;
  padding-top: 20rpx;
  border-top: 1rpx solid #f0f0f0;
}

.knowledge-label {
  font-size: 24rpx;
  color: #999;
}

.knowledge-tag {
  padding: 6rpx 16rpx;
  background-color: #e8f3ff;
  border-radius: 20rpx;
}

.knowledge-tag-text {
  font-size: 22rpx;
  color: #1677ff;
}

.history-item {
  display: flex;
  align-items: flex-start;
  gap: 16rpx;
  padding: 16rpx 0;
  border-bottom: 1rpx solid #f5f5f5;
}

.history-item:last-child {
  border-bottom: none;
}

.history-dot {
  width: 16rpx;
  height: 16rpx;
  border-radius: 50%;
  background-color: #1677ff;
  margin-top: 8rpx;
  flex-shrink: 0;
}

.history-content {
  flex: 1;
}

.history-date {
  font-size: 24rpx;
  color: #999;
  display: block;
  margin-bottom: 8rpx;
}

.history-answer {
  font-size: 26rpx;
  color: #666;
}

.no-history {
  padding: 32rpx 0;
  text-align: center;
}

.no-history-text {
  font-size: 26rpx;
  color: #ccc;
}

.mastery-indicator {
  display: flex;
  align-items: center;
  gap: 32rpx;
  margin-bottom: 24rpx;
}

.mastery-visual {
  flex-shrink: 0;
}

.mastery-ring {
  width: 120rpx;
  height: 120rpx;
  border-radius: 50%;
  border: 8rpx solid #1677ff;
  display: flex;
  align-items: center;
  justify-content: center;
}

.mastery-percent {
  font-size: 32rpx;
  font-weight: 700;
  color: #1677ff;
}

.mastery-info {
  flex: 1;
}

.mastery-status {
  font-size: 30rpx;
  font-weight: 600;
  display: block;
  margin-bottom: 8rpx;
}

.status-low { color: #ff4d4f; }
.status-medium { color: #faad14; }
.status-high { color: #52c41a; }

.mastery-hint {
  font-size: 24rpx;
  color: #999;
  line-height: 1.6;
}

.mastery-bar-wrap {
  margin-top: 8rpx;
}

.mastery-bar-bg {
  height: 16rpx;
  background-color: #f0f1f5;
  border-radius: 8rpx;
  overflow: hidden;
}

.mastery-bar-fill {
  height: 100%;
  border-radius: 8rpx;
  transition: width 0.5s ease;
}

.mastery-low { background-color: #ff4d4f; }
.mastery-medium { background-color: #faad14; }
.mastery-high { background-color: #52c41a; }

.notes-input {
  width: 100%;
  min-height: 160rpx;
  font-size: 28rpx;
  color: #333;
  line-height: 1.8;
  padding: 16rpx;
  background-color: #fafafa;
  border-radius: 12rpx;
  border: 2rpx solid #e8e8e8;
  box-sizing: border-box;
}

.notes-count {
  font-size: 22rpx;
  color: #ccc;
  text-align: right;
  display: block;
  margin-top: 8rpx;
}

.bottom-bar {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 24rpx;
  padding: 20rpx 32rpx;
  padding-bottom: calc(20rpx + env(safe-area-inset-bottom));
  background-color: #fff;
  box-shadow: 0 -2rpx 12rpx rgba(0, 0, 0, 0.06);
}

.action-btn {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 12rpx;
  padding: 24rpx 0;
  border-radius: 44rpx;
}

.review-btn {
  background-color: #fff;
  border: 2rpx solid #1677ff;
}

.review-btn .btn-icon,
.review-btn .btn-text {
  color: #1677ff;
}

.master-btn {
  background-color: #1677ff;
}

.master-btn .btn-icon,
.master-btn .btn-text {
  color: #fff;
}

.btn-icon {
  font-size: 32rpx;
}

.btn-text {
  font-size: 28rpx;
  font-weight: 600;
}
</style>
