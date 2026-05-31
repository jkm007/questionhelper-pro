<template>
  <view class="correction-page">
    <!-- Pre-filled question info -->
    <view class="form-section">
      <view class="section-header">
        <text class="section-title">Question Information</text>
        <view class="auto-fill-badge" v-if="questionInfo">
          <text class="auto-fill-text">Auto-filled</text>
        </view>
      </view>

      <view class="question-card" v-if="questionInfo">
        <view class="question-field">
          <text class="field-label">Question ID</text>
          <text class="field-value">{{ questionInfo.id }}</text>
        </view>
        <view class="question-field">
          <text class="field-label">Exam</text>
          <text class="field-value">{{ questionInfo.examName }}</text>
        </view>
        <view class="question-field">
          <text class="field-label">Category</text>
          <text class="field-value">{{ questionInfo.category }}</text>
        </view>
        <view class="question-content">
          <text class="field-label">Question Content</text>
          <text class="question-text">{{ questionInfo.content }}</text>
        </view>
        <view class="question-field" v-if="questionInfo.currentAnswer">
          <text class="field-label">Current Answer</text>
          <text class="field-value current-answer">{{ questionInfo.currentAnswer }}</text>
        </view>
      </view>

      <view class="no-question" v-else>
        <text class="no-question-text">No question information available</text>
      </view>
    </view>

    <!-- Correction form -->
    <view class="form-section">
      <view class="section-header">
        <text class="section-title">Correction Details</text>
      </view>

      <!-- Correct answer -->
      <view class="form-item">
        <text class="form-label">Correct Answer <text class="required">*</text></text>
        <view class="answer-options" v-if="questionInfo?.type === 'choice'">
          <view
            class="answer-option"
            v-for="opt in ['A', 'B', 'C', 'D', 'E']"
            :key="opt"
            :class="{ active: correctAnswer === opt }"
            @tap="correctAnswer = opt"
          >
            <text class="option-text" :class="{ 'active-text': correctAnswer === opt }">{{ opt }}</text>
          </view>
        </view>
        <view class="textarea-wrapper" v-else>
          <textarea
            class="form-textarea"
            v-model="correctAnswer"
            placeholder="Enter the correct answer..."
            :maxlength="200"
          />
        </view>
      </view>

      <!-- Explanation -->
      <view class="form-item">
        <text class="form-label">Explanation <text class="required">*</text></text>
        <view class="textarea-wrapper">
          <textarea
            class="form-textarea explanation-textarea"
            v-model="explanation"
            placeholder="Explain why this is the correct answer..."
            :maxlength="1000"
            :auto-height="false"
          />
          <text class="char-count">{{ explanation.length }}/1000</text>
        </view>
      </view>

      <!-- Reference source -->
      <view class="form-item">
        <text class="form-label">Reference Source</text>
        <view class="input-wrapper">
          <input
            class="form-input"
            v-model="referenceSource"
            placeholder="e.g., Textbook name, page number, or URL"
            :maxlength="200"
          />
        </view>
      </view>

      <!-- Image upload -->
      <view class="form-item">
        <text class="form-label">Supporting Images</text>
        <view class="image-upload-area">
          <view
            class="image-preview"
            v-for="(img, index) in imageList"
            :key="index"
          >
            <image
              class="preview-img"
              :src="img"
              mode="aspectFill"
              @tap="previewImage(index)"
            />
            <view class="remove-btn" @tap.stop="removeImage(index)">
              <text class="remove-icon">×</text>
            </view>
          </view>

          <view
            class="upload-btn"
            v-if="imageList.length < 5"
            @tap="chooseImage"
          >
            <text class="upload-icon">+</text>
            <text class="upload-text">Add Image</text>
          </view>
        </view>
        <text class="upload-hint">Upload screenshots or reference images (max 5)</text>
      </view>
    </view>

    <!-- Submit button -->
    <view class="submit-section">
      <view
        class="submit-btn"
        :class="{ disabled: !canSubmit }"
        @tap="handleSubmit"
      >
        <text class="submit-text">{{ submitting ? 'Submitting...' : 'Submit Correction' }}</text>
      </view>
      <text class="submit-hint">Your correction will be reviewed by our team</text>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { submitQuestionCorrection, type CorrectionParams, type QuestionInfo } from '@/api/feedback'

const questionInfo = ref<QuestionInfo | null>(null)
const correctAnswer = ref('')
const explanation = ref('')
const referenceSource = ref('')
const imageList = ref<string[]>([])
const questionId = ref('')
const submitting = ref(false)

const canSubmit = computed(() => {
  return (
    correctAnswer.value.trim().length > 0 &&
    explanation.value.trim().length > 0 &&
    !submitting.value
  )
})

const chooseImage = () => {
  const remaining = 5 - imageList.value.length
  if (remaining <= 0) return

  uni.chooseImage({
    count: remaining,
    sizeType: ['compressed'],
    sourceType: ['album', 'camera'],
    success: (res) => {
      imageList.value = [...imageList.value, ...res.tempFilePaths]
    }
  })
}

const removeImage = (index: number) => {
  imageList.value.splice(index, 1)
}

const previewImage = (index: number) => {
  uni.previewImage({
    current: index,
    urls: imageList.value
  })
}

const handleSubmit = async () => {
  if (!canSubmit.value) return

  if (!correctAnswer.value.trim()) {
    uni.showToast({ title: 'Please enter the correct answer', icon: 'none' })
    return
  }
  if (!explanation.value.trim()) {
    uni.showToast({ title: 'Please enter explanation', icon: 'none' })
    return
  }

  submitting.value = true

  try {
    const params: CorrectionParams = {
      questionId: questionId.value,
      correctAnswer: correctAnswer.value.trim(),
      explanation: explanation.value.trim(),
      referenceSource: referenceSource.value.trim() || undefined,
      images: imageList.value
    }

    await submitQuestionCorrection(params)

    uni.showModal({
      title: 'Correction Submitted',
      content: 'Thank you for your correction. Our team will review it and update the question if needed.',
      showCancel: false,
      success: () => {
        uni.navigateBack()
      }
    })
  } catch (error: any) {
    uni.showToast({
      title: error.message || 'Submission failed, please try again',
      icon: 'none'
    })
  } finally {
    submitting.value = false
  }
}

const loadQuestionInfo = async (id: string) => {
  try {
    // Simulated question info - replace with actual API call
    questionInfo.value = {
      id: id,
      examName: 'Civil Service Exam',
      category: 'Logical Reasoning',
      content: 'Which of the following conclusions can be logically derived from the given premises?',
      currentAnswer: 'B',
      type: 'choice'
    }
  } catch (error) {
    console.error('Failed to load question info:', error)
  }
}

onLoad((options) => {
  if (options?.questionId) {
    questionId.value = options.questionId
    loadQuestionInfo(options.questionId)
  }
})
</script>

<style lang="scss" scoped>
.correction-page {
  min-height: 100vh;
  background-color: #f5f6fa;
  padding-bottom: 40rpx;
}

.form-section {
  background-color: #ffffff;
  margin: 20rpx;
  border-radius: 16rpx;
  overflow: hidden;
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 24rpx 30rpx;
  border-bottom: 1rpx solid #f5f5f5;
}

.section-title {
  font-size: 30rpx;
  font-weight: 600;
  color: #333333;
}

.auto-fill-badge {
  background-color: #e8f5e9;
  padding: 6rpx 16rpx;
  border-radius: 8rpx;
}

.auto-fill-text {
  font-size: 22rpx;
  color: #4caf50;
}

.question-card {
  padding: 20rpx 30rpx 24rpx;
}

.question-field {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12rpx 0;
}

.field-label {
  font-size: 26rpx;
  color: #999999;
}

.field-value {
  font-size: 26rpx;
  color: #333333;
  font-weight: 500;
}

.current-answer {
  color: #007aff;
  background-color: #e8f0fe;
  padding: 4rpx 16rpx;
  border-radius: 6rpx;
}

.question-content {
  padding: 12rpx 0;
}

.question-text {
  font-size: 28rpx;
  color: #333333;
  line-height: 1.6;
  margin-top: 12rpx;
  display: block;
}

.no-question {
  padding: 40rpx 30rpx;
  text-align: center;
}

.no-question-text {
  font-size: 26rpx;
  color: #999999;
}

.form-item {
  padding: 24rpx 30rpx;
  border-bottom: 1rpx solid #f5f5f5;
}

.form-item:last-child {
  border-bottom: none;
}

.form-label {
  font-size: 28rpx;
  font-weight: 500;
  color: #333333;
  margin-bottom: 16rpx;
  display: block;
}

.required {
  color: #ff4757;
}

.answer-options {
  display: flex;
  gap: 16rpx;
}

.answer-option {
  width: 80rpx;
  height: 80rpx;
  border-radius: 12rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #f8f9ff;
  border: 2rpx solid #e8ecf4;
}

.answer-option.active {
  background-color: #007aff;
  border-color: #007aff;
}

.option-text {
  font-size: 32rpx;
  font-weight: 600;
  color: #666666;
}

.option-text.active-text {
  color: #ffffff;
}

.textarea-wrapper {
  position: relative;
  background-color: #f8f9ff;
  border-radius: 12rpx;
  border: 1rpx solid #e8ecf4;
  padding: 20rpx 24rpx;
}

.form-textarea {
  width: 100%;
  min-height: 150rpx;
  font-size: 28rpx;
  color: #333333;
  line-height: 1.6;
}

.explanation-textarea {
  min-height: 240rpx;
}

.char-count {
  position: absolute;
  bottom: 12rpx;
  right: 16rpx;
  font-size: 22rpx;
  color: #cccccc;
}

.input-wrapper {
  background-color: #f8f9ff;
  border-radius: 12rpx;
  border: 1rpx solid #e8ecf4;
  padding: 20rpx 24rpx;
}

.form-input {
  width: 100%;
  font-size: 28rpx;
  color: #333333;
}

.image-upload-area {
  display: flex;
  flex-wrap: wrap;
  gap: 16rpx;
}

.image-preview {
  position: relative;
  width: 160rpx;
  height: 160rpx;
  border-radius: 12rpx;
  overflow: hidden;
}

.preview-img {
  width: 100%;
  height: 100%;
}

.remove-btn {
  position: absolute;
  top: 0;
  right: 0;
  width: 36rpx;
  height: 36rpx;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  border-bottom-left-radius: 8rpx;
}

.remove-icon {
  font-size: 24rpx;
  color: #ffffff;
  font-weight: bold;
}

.upload-btn {
  width: 160rpx;
  height: 160rpx;
  border-radius: 12rpx;
  border: 2rpx dashed #cccccc;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background-color: #fafafa;
}

.upload-icon {
  font-size: 44rpx;
  color: #cccccc;
  margin-bottom: 4rpx;
}

.upload-text {
  font-size: 20rpx;
  color: #cccccc;
}

.upload-hint {
  font-size: 22rpx;
  color: #999999;
  margin-top: 12rpx;
  display: block;
}

.submit-section {
  padding: 30rpx 20rpx;
  text-align: center;
}

.submit-btn {
  width: 100%;
  height: 96rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #007aff;
  border-radius: 16rpx;
}

.submit-btn.disabled {
  background-color: #cccccc;
}

.submit-text {
  font-size: 32rpx;
  color: #ffffff;
  font-weight: 500;
}

.submit-hint {
  font-size: 24rpx;
  color: #999999;
  margin-top: 16rpx;
  display: block;
}
</style>
