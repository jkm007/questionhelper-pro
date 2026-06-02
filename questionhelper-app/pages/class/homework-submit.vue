<template>
  <view class="page">
    <!-- Loading -->
    <view v-if="loading" class="loading-wrap">
      <text class="loading-text">加载中...</text>
    </view>

    <template v-else-if="homework">
      <!-- Homework Info Card -->
      <view class="info-card">
        <view class="info-header">
          <text class="info-title">{{ homework.title }}</text>
          <view :class="['status-badge', 'status-' + homework.status]">
            <text class="status-text">{{ statusLabels[homework.status] }}</text>
          </view>
        </view>
        <text v-if="homework.description" class="info-desc">{{ homework.description }}</text>
        <view class="info-meta">
          <view class="meta-row">
            <text class="meta-label">发布者</text>
            <text class="meta-value">{{ homework.creatorName }}</text>
          </view>
          <view class="meta-row">
            <text class="meta-label">截止时间</text>
            <text :class="['meta-value', isDeadlineSoon ? 'deadline-warn' : '']">
              {{ homework.deadline }}
            </text>
          </view>
          <view class="meta-row">
            <text class="meta-label">题目数量</text>
            <text class="meta-value">{{ homework.questionCount }}题</text>
          </view>
          <view v-if="homework.totalScore" class="meta-row">
            <text class="meta-label">总分</text>
            <text class="meta-value highlight">{{ homework.totalScore }}分</text>
          </view>
        </view>
      </view>

      <!-- Questions Section -->
      <view v-if="homework.questions && homework.questions.length > 0" class="questions-section">
        <text class="section-title">作业题目</text>
        <view
          v-for="(q, index) in homework.questions"
          :key="q.id"
          class="question-card"
        >
          <view class="question-header">
            <view class="question-num-wrap">
              <text class="question-num">{{ index + 1 }}</text>
            </view>
            <view :class="['question-type-badge', 'type-' + q.type]">
              <text class="question-type-text">{{ typeLabels[q.type] || '未知' }}</text>
            </view>
            <text class="question-score">{{ q.score }}分</text>
          </view>
          <text class="question-content">{{ q.content }}</text>

          <!-- Options for choice questions -->
          <view v-if="q.options && q.options.length > 0" class="options-list">
            <view
              v-for="(opt, oIndex) in q.options"
              :key="oIndex"
              class="option-item"
            >
              <text class="option-label">{{ optionLabels[oIndex] }}.</text>
              <text class="option-text">{{ opt.content }}</text>
            </view>
          </view>
        </view>
      </view>

      <!-- Attachments -->
      <view v-if="homework.attachments && homework.attachments.length > 0" class="attachments-card">
        <text class="section-title">作业附件</text>
        <view
          v-for="(att, index) in homework.attachments"
          :key="index"
          class="attachment-item"
          @tap="onPreviewAttachment(att)"
        >
          <text class="attachment-icon">&#xe628;</text>
          <text class="attachment-name">{{ att.name || '附件' + (index + 1) }}</text>
          <text class="attachment-action">查看</text>
        </view>
      </view>

      <!-- Answer Section -->
      <view class="answer-section">
        <text class="section-title">提交作业</text>

        <!-- Text Answer -->
        <view class="answer-card">
          <text class="answer-label">文字作答</text>
          <view class="textarea-wrap">
            <textarea
              class="answer-textarea"
              v-model="answerContent"
              placeholder="请输入你的答案..."
              :maxlength="2000"
              auto-height
            />
          </view>
          <text class="char-count">{{ answerContent.length }}/2000</text>
        </view>

        <!-- File Upload -->
        <view class="answer-card">
          <view class="answer-label-row">
            <text class="answer-label">附件上传</text>
            <text class="answer-hint">最多5个文件</text>
          </view>
          <view class="upload-area">
            <view
              v-for="(file, index) in uploadFiles"
              :key="index"
              class="upload-file-item"
            >
              <text class="file-icon">&#xe628;</text>
              <text class="file-name">{{ file.name }}</text>
              <view v-if="file.uploading" class="file-progress">
                <text class="file-progress-text">上传中...</text>
              </view>
              <text v-else class="file-remove" @tap="onRemoveFile(index)">&#xe621;</text>
            </view>
            <view
              v-if="uploadFiles.length < 5"
              class="upload-trigger"
              @tap="onChooseFile"
            >
              <text class="upload-icon">+</text>
              <text class="upload-text">添加文件</text>
            </view>
          </view>
        </view>
      </view>
    </template>

    <!-- Submit Button -->
    <view v-if="!loading && homework && homework.status !== 'expired'" class="submit-section">
      <view :class="['submit-btn', submitting ? 'disabled' : '']" @tap="onSubmit">
        <text class="submit-btn-text">{{ submitting ? '提交中...' : '提交作业' }}</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { getHomeworkDetail, submitHomework } from '@/api/class'
import { uploadFile } from '@/api/file'

interface HomeworkQuestion {
  id: string
  content: string
  type: number
  score: number
  options?: { content: string }[]
}

interface HomeworkDetail {
  id: string
  title: string
  description: string
  status: string
  deadline: string
  creatorName: string
  questionCount: number
  totalScore: number
  questions: HomeworkQuestion[]
  attachments?: { name: string; url: string }[]
}

interface UploadFile {
  name: string
  url: string
  filePath: string
  uploading: boolean
}

const statusLabels: Record<string, string> = {
  ongoing: '进行中',
  expired: '已截止',
  closed: '已关闭'
}

const typeLabels: Record<number, string> = {
  1: '单选题',
  2: '多选题',
  3: '判断题',
  4: '填空题',
  5: '简答题'
}

const optionLabels = ['A', 'B', 'C', 'D', 'E', 'F', 'G', 'H']

const classId = ref('')
const homeworkId = ref('')
const loading = ref(true)
const submitting = ref(false)
const homework = ref<HomeworkDetail | null>(null)
const answerContent = ref('')
const uploadFiles = ref<UploadFile[]>([])

const isDeadlineSoon = computed(() => {
  if (!homework.value) return false
  const deadlineTime = new Date(homework.value.deadline).getTime()
  const now = Date.now()
  const diff = deadlineTime - now
  return diff > 0 && diff < 24 * 60 * 60 * 1000
})

// ---------- Data Loading ----------

async function loadHomework() {
  loading.value = true
  try {
    const res = await getHomeworkDetail(classId.value, homeworkId.value)
    homework.value = res.data
  } catch (e) {
    uni.showToast({ title: '加载失败', icon: 'none' })
  } finally {
    loading.value = false
  }
}

// ---------- File Handling ----------

function onChooseFile() {
  uni.chooseMessageFile({
    count: 5 - uploadFiles.value.length,
    type: 'file',
    success: (res) => {
      for (const file of res.tempFiles) {
        if (uploadFiles.value.length >= 5) break
        uploadFiles.value.push({
          name: file.name,
          url: '',
          filePath: file.path,
          uploading: false
        })
      }
    }
  })
}

function onRemoveFile(index: number) {
  uploadFiles.value.splice(index, 1)
}

function onPreviewAttachment(att: { name: string; url: string }) {
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

// ---------- Submit ----------

function validate(): boolean {
  if (!answerContent.value.trim() && uploadFiles.value.length === 0) {
    uni.showToast({ title: '请输入答案或上传附件', icon: 'none' })
    return false
  }
  return true
}

async function onSubmit() {
  if (submitting.value) return
  if (!validate()) return

  uni.showModal({
    title: '确认提交',
    content: '提交后将无法修改，确定提交作业吗？',
    success: async (res) => {
      if (!res.confirm) return

      submitting.value = true
      try {
        // Upload files
        const attachmentUrls: string[] = []
        for (const file of uploadFiles.value) {
          if (!file.url) {
            file.uploading = true
            try {
              const uploadRes: any = await uploadFile(file.filePath)
              file.url = uploadRes.data?.url || ''
            } catch (e) {
              file.uploading = false
              uni.showToast({ title: `文件「${file.name}」上传失败`, icon: 'none' })
              submitting.value = false
              return
            }
            file.uploading = false
          }
          if (file.url) {
            attachmentUrls.push(file.url)
          }
        }

        await submitHomework(classId.value, homeworkId.value, {
          content: answerContent.value.trim(),
          attachments: attachmentUrls.length > 0 ? attachmentUrls : undefined
        })

        uni.showToast({ title: '提交成功', icon: 'success' })
        setTimeout(() => {
          uni.navigateBack()
        }, 1500)
      } catch (e) {
        uni.showToast({ title: '提交失败', icon: 'none' })
      } finally {
        submitting.value = false
      }
    }
  })
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

  loadHomework()
})
</script>

<style scoped>
.page {
  min-height: 100vh;
  background-color: #f5f6fa;
  padding-bottom: 160rpx;
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

/* ---- Info Card ---- */

.info-card {
  background-color: #fff;
  margin: 20rpx 24rpx;
  border-radius: 16rpx;
  padding: 28rpx;
  box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.04);
}

.info-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16rpx;
}

.info-title {
  flex: 1;
  font-size: 32rpx;
  font-weight: 600;
  color: #333;
  margin-right: 16rpx;
}

.status-badge {
  padding: 6rpx 16rpx;
  border-radius: 16rpx;
  flex-shrink: 0;
}

.status-ongoing {
  background-color: #e8f3ff;
}

.status-ongoing .status-text {
  color: #1677ff;
}

.status-expired {
  background-color: #fff1f0;
}

.status-expired .status-text {
  color: #ff4d4f;
}

.status-closed {
  background-color: #f0f0f0;
}

.status-closed .status-text {
  color: #999;
}

.status-text {
  font-size: 22rpx;
}

.info-desc {
  font-size: 26rpx;
  color: #666;
  line-height: 1.6;
  margin-bottom: 20rpx;
  display: block;
}

.info-meta {
  background-color: #f8f9fc;
  border-radius: 12rpx;
  padding: 20rpx;
}

.meta-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8rpx 0;
}

.meta-label {
  font-size: 26rpx;
  color: #999;
}

.meta-value {
  font-size: 26rpx;
  color: #333;
}

.meta-value.deadline-warn {
  color: #ff4d4f;
  font-weight: 600;
}

.meta-value.highlight {
  color: #1677ff;
  font-weight: 600;
}

/* ---- Questions ---- */

.questions-section {
  padding: 0 24rpx;
  margin-bottom: 20rpx;
}

.section-title {
  font-size: 30rpx;
  font-weight: 600;
  color: #333;
  margin-bottom: 16rpx;
  display: block;
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
  margin-bottom: 16rpx;
}

.question-num-wrap {
  width: 44rpx;
  height: 44rpx;
  border-radius: 50%;
  background-color: #1677ff;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 12rpx;
  flex-shrink: 0;
}

.question-num {
  font-size: 24rpx;
  font-weight: 600;
  color: #fff;
}

.question-type-badge {
  padding: 4rpx 14rpx;
  border-radius: 12rpx;
  margin-right: 12rpx;
}

.type-1 { background-color: #e8f3ff; }
.type-1 .question-type-text { color: #1677ff; }

.type-2 { background-color: #f6ffed; }
.type-2 .question-type-text { color: #52c41a; }

.type-3 { background-color: #fff0e6; }
.type-3 .question-type-text { color: #fa8c16; }

.type-4 { background-color: #f9f0ff; }
.type-4 .question-type-text { color: #722ed1; }

.type-5 { background-color: #fff1f0; }
.type-5 .question-type-text { color: #ff4d4f; }

.question-type-text {
  font-size: 22rpx;
}

.question-score {
  font-size: 24rpx;
  color: #1677ff;
  font-weight: 600;
  margin-left: auto;
}

.question-content {
  font-size: 28rpx;
  color: #333;
  line-height: 1.8;
  display: block;
  margin-bottom: 16rpx;
}

.options-list {
  padding: 8rpx 0;
}

.option-item {
  display: flex;
  padding: 10rpx 0;
}

.option-label {
  font-size: 26rpx;
  color: #666;
  margin-right: 12rpx;
  min-width: 40rpx;
}

.option-text {
  font-size: 26rpx;
  color: #333;
  flex: 1;
}

/* ---- Attachments ---- */

.attachments-card {
  background-color: #fff;
  margin: 0 24rpx 20rpx;
  border-radius: 16rpx;
  padding: 28rpx;
  box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.04);
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

.attachment-action {
  font-size: 24rpx;
  color: #1677ff;
  margin-left: 16rpx;
}

/* ---- Answer Section ---- */

.answer-section {
  padding: 0 24rpx;
}

.answer-card {
  background-color: #fff;
  border-radius: 16rpx;
  padding: 24rpx;
  margin-bottom: 16rpx;
  box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.04);
}

.answer-label {
  font-size: 28rpx;
  font-weight: 600;
  color: #333;
  margin-bottom: 16rpx;
  display: block;
}

.answer-label-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16rpx;
}

.answer-hint {
  font-size: 22rpx;
  color: #999;
}

.textarea-wrap {
  border: 2rpx solid #e0e0e0;
  border-radius: 12rpx;
  padding: 16rpx 20rpx;
  min-height: 200rpx;
}

.answer-textarea {
  width: 100%;
  font-size: 28rpx;
  line-height: 1.8;
}

.char-count {
  font-size: 22rpx;
  color: #ccc;
  text-align: right;
  margin-top: 8rpx;
  display: block;
}

/* ---- Upload ---- */

.upload-area {
  display: flex;
  flex-wrap: wrap;
  gap: 16rpx;
}

.upload-file-item {
  display: flex;
  align-items: center;
  background-color: #f8f9fc;
  border-radius: 12rpx;
  padding: 16rpx 20rpx;
  width: 100%;
}

.file-icon {
  font-size: 28rpx;
  color: #1677ff;
  margin-right: 12rpx;
}

.file-name {
  flex: 1;
  font-size: 26rpx;
  color: #333;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.file-progress {
  margin-left: 12rpx;
}

.file-progress-text {
  font-size: 22rpx;
  color: #fa8c16;
}

.file-remove {
  font-size: 28rpx;
  color: #999;
  padding: 4rpx;
  margin-left: 12rpx;
}

.upload-trigger {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 120rpx;
  border: 2rpx dashed #d9d9d9;
  border-radius: 12rpx;
  background-color: #fafafa;
}

.upload-icon {
  font-size: 48rpx;
  color: #d9d9d9;
  line-height: 1;
  margin-bottom: 4rpx;
}

.upload-text {
  font-size: 24rpx;
  color: #999;
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
