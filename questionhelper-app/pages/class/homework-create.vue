<template>
  <view class="page">
    <view class="form-section">
      <!-- Title -->
      <view class="form-item">
        <text class="form-label required">作业标题</text>
        <view class="input-wrap">
          <input
            class="form-input"
            v-model="form.title"
            placeholder="请输入作业标题"
            :maxlength="50"
          />
        </view>
        <text class="char-count">{{ form.title.length }}/50</text>
      </view>

      <!-- Description -->
      <view class="form-item">
        <text class="form-label">作业描述</text>
        <view class="textarea-wrap">
          <textarea
            class="form-textarea"
            v-model="form.description"
            placeholder="请输入作业要求及说明（选填）"
            :maxlength="500"
            auto-height
          />
        </view>
        <text class="char-count">{{ form.description.length }}/500</text>
      </view>

      <!-- Deadline -->
      <view class="form-item">
        <text class="form-label required">截止时间</text>
        <picker mode="date" :value="deadlineDate" :start="today" @change="onDateChange">
          <view class="picker-row">
            <text :class="['picker-value', deadlineDate ? '' : 'placeholder']">
              {{ deadlineDate || '请选择截止日期' }}
            </text>
            <text class="picker-arrow">&#xe61e;</text>
          </view>
        </picker>
        <picker v-if="deadlineDate" mode="time" :value="deadlineTime" @change="onTimeChange">
          <view class="picker-row">
            <text :class="['picker-value', deadlineTime ? '' : 'placeholder']">
              {{ deadlineTime || '请选择截止时间' }}
            </text>
            <text class="picker-arrow">&#xe61e;</text>
          </view>
        </picker>
      </view>
    </view>

    <!-- Select Questions -->
    <view class="section-header">
      <text class="section-title">选择题目</text>
      <view class="section-action" @tap="onSelectQuestions">
        <text class="section-action-text">{{ selectedQuestions.length > 0 ? '重新选择' : '从题库选择' }}</text>
      </view>
    </view>

    <!-- Selected Questions Preview -->
    <view v-if="selectedQuestions.length > 0" class="questions-section">
      <view
        v-for="(q, index) in selectedQuestions"
        :key="q.id"
        class="question-card"
      >
        <view class="question-header">
          <text class="question-index">{{ index + 1 }}</text>
          <view :class="['question-type-badge', 'type-' + q.type]">
            <text class="question-type-text">{{ typeLabels[q.type] }}</text>
          </view>
        </view>
        <text class="question-content">{{ q.content }}</text>
        <view class="question-footer">
          <text class="question-score">分值: {{ q.score }}分</text>
          <text class="question-remove" @tap="onRemoveQuestion(index)">移除</text>
        </view>
      </view>
      <view class="question-summary">
        <text class="summary-text">共 {{ selectedQuestions.length }} 题，总分 {{ totalScore }} 分</text>
      </view>
    </view>
    <view v-else class="empty-questions">
      <text class="empty-hint">请从题库中选择作业题目</text>
    </view>

    <!-- Attachments -->
    <view class="section-header">
      <text class="section-title">附件</text>
      <view class="section-action" @tap="onChooseFile">
        <text class="section-action-text">添加附件</text>
      </view>
    </view>

    <view v-if="attachments.length > 0" class="attachments-section">
      <view
        v-for="(file, index) in attachments"
        :key="index"
        class="attachment-item"
      >
        <text class="attachment-icon">&#xe628;</text>
        <text class="attachment-name">{{ file.name }}</text>
        <text class="attachment-remove" @tap="onRemoveFile(index)">&#xe621;</text>
      </view>
    </view>

    <!-- Submit Button -->
    <view class="submit-section">
      <view :class="['submit-btn', submitting ? 'disabled' : '']" @tap="onSubmit">
        <text class="submit-btn-text">{{ submitting ? '发布中...' : '发布作业' }}</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { createHomework } from '@/api/class'
import { uploadFile } from '@/api/file'
import { getQuestions } from '@/api/question'

interface QuestionItem {
  id: string
  content: string
  type: number
  score: number
}

interface AttachmentItem {
  name: string
  url: string
  filePath: string
}

const typeLabels: Record<number, string> = {
  1: '单选',
  2: '多选',
  3: '判断',
  4: '填空',
  5: '简答'
}

const classId = ref('')
const submitting = ref(false)

const form = reactive({
  title: '',
  description: ''
})

const deadlineDate = ref('')
const deadlineTime = ref('')
const today = ref('')

const selectedQuestions = ref<QuestionItem[]>([])
const attachments = ref<AttachmentItem[]>([])

const totalScore = computed(() => {
  return selectedQuestions.value.reduce((sum, q) => sum + (q.score || 0), 0)
})

// ---------- Date / Time ----------

function onDateChange(e: any) {
  deadlineDate.value = e.detail.value
}

function onTimeChange(e: any) {
  deadlineTime.value = e.detail.value
}

// ---------- Question Selection ----------

function onSelectQuestions() {
  uni.navigateTo({
    url: `/pages/question/list?mode=select&classId=${classId.value}&callback=homeworkCreate`
  })
}

function onRemoveQuestion(index: number) {
  selectedQuestions.value.splice(index, 1)
}

// ---------- Attachments ----------

function onChooseFile() {
  uni.chooseMessageFile({
    count: 5 - attachments.value.length,
    type: 'file',
    success: (res) => {
      for (const file of res.tempFiles) {
        if (attachments.value.length >= 5) break
        attachments.value.push({
          name: file.name,
          url: '',
          filePath: file.path
        })
      }
    }
  })
}

function onRemoveFile(index: number) {
  attachments.value.splice(index, 1)
}

// ---------- Submit ----------

function validate(): boolean {
  if (!form.title.trim()) {
    uni.showToast({ title: '请输入作业标题', icon: 'none' })
    return false
  }
  if (!deadlineDate.value) {
    uni.showToast({ title: '请选择截止日期', icon: 'none' })
    return false
  }
  if (selectedQuestions.value.length === 0) {
    uni.showToast({ title: '请至少选择一道题目', icon: 'none' })
    return false
  }
  return true
}

async function onSubmit() {
  if (submitting.value) return
  if (!validate()) return

  submitting.value = true
  try {
    // Upload attachments
    const uploadedUrls: string[] = []
    for (const file of attachments.value) {
      if (!file.url && file.filePath) {
        const res: any = await uploadFile(file.filePath)
        file.url = res.data?.url || ''
      }
      if (file.url) {
        uploadedUrls.push(file.url)
      }
    }

    const deadline = deadlineTime.value
      ? `${deadlineDate.value} ${deadlineTime.value}:00`
      : `${deadlineDate.value} 23:59:59`

    await createHomework(classId.value, {
      title: form.title.trim(),
      description: form.description.trim() || undefined,
      deadline,
      questionIds: selectedQuestions.value.map((q) => Number(q.id)),
      attachments: uploadedUrls.length > 0 ? uploadedUrls : undefined
    })

    uni.showToast({ title: '发布成功', icon: 'success' })
    setTimeout(() => {
      uni.navigateBack()
    }, 1500)
  } catch (e) {
    uni.showToast({ title: '发布失败', icon: 'none' })
  } finally {
    submitting.value = false
  }
}

// ---------- Lifecycle ----------

onMounted(() => {
  const pages = getCurrentPages()
  const currentPage = pages[pages.length - 1] as any
  classId.value = currentPage?.options?.classId || ''

  if (!classId.value) {
    uni.showToast({ title: '参数错误', icon: 'none' })
    setTimeout(() => uni.navigateBack(), 1500)
    return
  }

  // Set today for date picker min
  const now = new Date()
  today.value = `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}-${String(now.getDate()).padStart(2, '0')}`

  // Listen for question selection result from question list page
  uni.$on('homeworkQuestionSelected', (questions: QuestionItem[]) => {
    selectedQuestions.value = questions
  })
})
</script>

<style scoped>
.page {
  min-height: 100vh;
  background-color: #f5f6fa;
  padding-bottom: 160rpx;
}

.form-section {
  background-color: #fff;
  margin: 20rpx 24rpx;
  border-radius: 16rpx;
  overflow: hidden;
}

.form-item {
  padding: 24rpx 28rpx;
  border-bottom: 1rpx solid #f5f5f5;
}

.form-item:last-child {
  border-bottom: none;
}

.form-label {
  font-size: 28rpx;
  font-weight: 600;
  color: #333;
  margin-bottom: 16rpx;
  display: block;
}

.form-label.required::before {
  content: '*';
  color: #ff4d4f;
  margin-right: 4rpx;
}

.input-wrap {
  border: 2rpx solid #e0e0e0;
  border-radius: 12rpx;
  padding: 0 20rpx;
  height: 80rpx;
}

.form-input {
  width: 100%;
  height: 80rpx;
  font-size: 28rpx;
}

.textarea-wrap {
  border: 2rpx solid #e0e0e0;
  border-radius: 12rpx;
  padding: 16rpx 20rpx;
  min-height: 160rpx;
}

.form-textarea {
  width: 100%;
  font-size: 28rpx;
  line-height: 1.6;
}

.char-count {
  font-size: 22rpx;
  color: #ccc;
  text-align: right;
  margin-top: 8rpx;
  display: block;
}

.picker-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16rpx 0;
  border-bottom: 1rpx solid #f5f5f5;
}

.picker-row:last-child {
  border-bottom: none;
}

.picker-value {
  font-size: 28rpx;
  color: #333;
}

.picker-value.placeholder {
  color: #ccc;
}

.picker-arrow {
  font-size: 24rpx;
  color: #ccc;
}

/* ---- Section Header ---- */

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 24rpx 32rpx 12rpx;
}

.section-title {
  font-size: 30rpx;
  font-weight: 600;
  color: #333;
}

.section-action {
  padding: 8rpx 20rpx;
  background-color: #e8f3ff;
  border-radius: 24rpx;
}

.section-action-text {
  font-size: 24rpx;
  color: #1677ff;
}

/* ---- Questions ---- */

.questions-section {
  padding: 0 24rpx;
}

.question-card {
  background-color: #fff;
  border-radius: 16rpx;
  padding: 24rpx;
  margin-bottom: 16rpx;
  box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.04);
}

.question-header {
  display: flex;
  align-items: center;
  margin-bottom: 12rpx;
}

.question-index {
  width: 40rpx;
  height: 40rpx;
  border-radius: 50%;
  background-color: #1677ff;
  color: #fff;
  font-size: 22rpx;
  font-weight: 600;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 12rpx;
  flex-shrink: 0;
}

.question-type-badge {
  padding: 4rpx 14rpx;
  border-radius: 12rpx;
}

.type-1 {
  background-color: #e8f3ff;
}

.type-1 .question-type-text {
  color: #1677ff;
}

.type-2 {
  background-color: #f6ffed;
}

.type-2 .question-type-text {
  color: #52c41a;
}

.type-3 {
  background-color: #fff0e6;
}

.type-3 .question-type-text {
  color: #fa8c16;
}

.type-4 {
  background-color: #f9f0ff;
}

.type-4 .question-type-text {
  color: #722ed1;
}

.type-5 {
  background-color: #fff1f0;
}

.type-5 .question-type-text {
  color: #ff4d4f;
}

.question-type-text {
  font-size: 22rpx;
}

.question-content {
  font-size: 28rpx;
  color: #333;
  line-height: 1.6;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  margin-bottom: 12rpx;
}

.question-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.question-score {
  font-size: 24rpx;
  color: #1677ff;
  font-weight: 600;
}

.question-remove {
  font-size: 24rpx;
  color: #ff4d4f;
  padding: 4rpx 12rpx;
}

.question-summary {
  padding: 16rpx 0;
  text-align: center;
}

.summary-text {
  font-size: 26rpx;
  color: #666;
}

.empty-questions {
  padding: 60rpx 24rpx;
  text-align: center;
}

.empty-hint {
  font-size: 28rpx;
  color: #ccc;
}

/* ---- Attachments ---- */

.attachments-section {
  padding: 0 24rpx;
}

.attachment-item {
  display: flex;
  align-items: center;
  background-color: #fff;
  border-radius: 12rpx;
  padding: 20rpx 24rpx;
  margin-bottom: 12rpx;
}

.attachment-icon {
  font-size: 32rpx;
  color: #1677ff;
  margin-right: 16rpx;
}

.attachment-name {
  flex: 1;
  font-size: 26rpx;
  color: #333;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.attachment-remove {
  font-size: 28rpx;
  color: #999;
  padding: 8rpx;
  margin-left: 16rpx;
}

/* ---- Submit ---- */

.submit-section {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  padding: 20rpx 32rpx;
  padding-bottom: calc(20rpx + env(safe-area-inset-bottom));
  background-color: #fff;
  box-shadow: 0 -2rpx 12rpx rgba(0, 0, 0, 0.06);
}

.submit-btn {
  background-color: #1677ff;
  border-radius: 44rpx;
  padding: 24rpx 0;
  text-align: center;
}

.submit-btn.disabled {
  opacity: 0.6;
}

.submit-btn-text {
  font-size: 30rpx;
  font-weight: 600;
  color: #fff;
}
</style>
