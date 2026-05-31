<template>
  <view class="exam-feedback-page">
    <view class="form-section">
      <!-- Exam selector -->
      <view class="form-item">
        <text class="form-label">Select Exam</text>
        <picker
          :range="examOptions"
          :range-key="'name'"
          :value="selectedExamIndex"
          @change="onExamChange"
        >
          <view class="picker-wrapper">
            <text class="picker-text" :class="{ placeholder: selectedExamIndex < 0 }">
              {{ selectedExamIndex >= 0 ? examOptions[selectedExamIndex].name : 'Please select an exam' }}
            </text>
            <text class="picker-arrow">▼</text>
          </view>
        </picker>
      </view>

      <!-- Feedback type -->
      <view class="form-item">
        <text class="form-label">Feedback Type</text>
        <view class="type-options">
          <view
            class="type-option"
            v-for="option in feedbackTypes"
            :key="option.value"
            :class="{ active: selectedType === option.value }"
            @tap="selectedType = option.value"
          >
            <text class="type-option-text" :class="{ 'active-text': selectedType === option.value }">
              {{ option.label }}
            </text>
          </view>
        </view>
      </view>

      <!-- Description -->
      <view class="form-item">
        <text class="form-label">Description</text>
        <view class="textarea-wrapper">
          <textarea
            class="form-textarea"
            v-model="description"
            placeholder="Please describe the issue in detail..."
            :maxlength="500"
            :auto-height="false"
          />
          <text class="char-count">{{ description.length }}/500</text>
        </view>
      </view>

      <!-- Image upload -->
      <view class="form-item">
        <text class="form-label">Upload Images (max 3)</text>
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
            v-if="imageList.length < 3"
            @tap="chooseImage"
          >
            <text class="upload-icon">+</text>
            <text class="upload-text">Add Image</text>
          </view>
        </view>
        <text class="upload-hint">Screenshots help us locate the issue faster</text>
      </view>
    </view>

    <!-- Question info (optional) -->
    <view class="form-section" v-if="questionId">
      <view class="section-title">
        <text class="section-title-text">Related Question</text>
      </view>
      <view class="question-info">
        <text class="question-id">Question ID: {{ questionId }}</text>
      </view>
    </view>

    <!-- Submit button -->
    <view class="submit-section">
      <view
        class="submit-btn"
        :class="{ disabled: !canSubmit }"
        @tap="handleSubmit"
      >
        <text class="submit-text">{{ submitting ? 'Submitting...' : 'Submit Feedback' }}</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { submitExamFeedback, type ExamFeedbackParams } from '@/api/feedback'

interface ExamOption {
  id: string
  name: string
}

const examOptions = ref<ExamOption[]>([
  { id: '1', name: 'Civil Service Exam (行测)' },
  { id: '2', name: 'Civil Service Exam (申论)' },
  { id: '3', name: 'Teacher Certification (综合素质)' },
  { id: '4', name: 'Teacher Certification (教育知识)' },
  { id: '5', name: 'Driving License Exam (Subject 1)' },
  { id: '6', name: 'Driving License Exam (Subject 4)' }
])

const feedbackTypes = [
  { value: 'error_question', label: 'Error in Question' },
  { value: 'wrong_answer', label: 'Wrong Answer' },
  { value: 'other', label: 'Other' }
]

const selectedExamIndex = ref(-1)
const selectedType = ref('')
const description = ref('')
const imageList = ref<string[]>([])
const questionId = ref('')
const submitting = ref(false)

const canSubmit = computed(() => {
  return (
    selectedExamIndex.value >= 0 &&
    selectedType.value !== '' &&
    description.value.trim().length > 0 &&
    !submitting.value
  )
})

const onExamChange = (e: any) => {
  selectedExamIndex.value = e.detail.value
}

const chooseImage = () => {
  const remaining = 3 - imageList.value.length
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

  // Validate
  if (selectedExamIndex.value < 0) {
    uni.showToast({ title: 'Please select an exam', icon: 'none' })
    return
  }
  if (!selectedType.value) {
    uni.showToast({ title: 'Please select feedback type', icon: 'none' })
    return
  }
  if (!description.value.trim()) {
    uni.showToast({ title: 'Please enter description', icon: 'none' })
    return
  }

  submitting.value = true

  try {
    const params: ExamFeedbackParams = {
      examId: examOptions.value[selectedExamIndex.value].id,
      examName: examOptions.value[selectedExamIndex.value].name,
      feedbackType: selectedType.value,
      description: description.value.trim(),
      images: imageList.value,
      questionId: questionId.value || undefined
    }

    await submitExamFeedback(params)

    uni.showModal({
      title: 'Submitted Successfully',
      content: 'Thank you for your feedback. We will review it as soon as possible.',
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

onLoad((options) => {
  if (options?.questionId) {
    questionId.value = options.questionId
  }
  if (options?.examId) {
    const index = examOptions.value.findIndex((e) => e.id === options.examId)
    if (index >= 0) {
      selectedExamIndex.value = index
    }
  }
})
</script>

<style lang="scss" scoped>
.exam-feedback-page {
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

.section-title {
  padding: 24rpx 30rpx 0;
}

.section-title-text {
  font-size: 28rpx;
  font-weight: 600;
  color: #333333;
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

.picker-wrapper {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20rpx 24rpx;
  background-color: #f8f9ff;
  border-radius: 12rpx;
  border: 1rpx solid #e8ecf4;
}

.picker-text {
  font-size: 28rpx;
  color: #333333;
}

.picker-text.placeholder {
  color: #cccccc;
}

.picker-arrow {
  font-size: 20rpx;
  color: #999999;
}

.type-options {
  display: flex;
  flex-wrap: wrap;
  gap: 16rpx;
}

.type-option {
  padding: 16rpx 28rpx;
  background-color: #f8f9ff;
  border-radius: 8rpx;
  border: 2rpx solid #e8ecf4;
}

.type-option.active {
  background-color: #e8f0fe;
  border-color: #007aff;
}

.type-option-text {
  font-size: 26rpx;
  color: #666666;
}

.type-option-text.active-text {
  color: #007aff;
  font-weight: 500;
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
  min-height: 200rpx;
  font-size: 28rpx;
  color: #333333;
  line-height: 1.6;
}

.char-count {
  position: absolute;
  bottom: 12rpx;
  right: 16rpx;
  font-size: 22rpx;
  color: #cccccc;
}

.image-upload-area {
  display: flex;
  flex-wrap: wrap;
  gap: 16rpx;
}

.image-preview {
  position: relative;
  width: 180rpx;
  height: 180rpx;
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
  width: 40rpx;
  height: 40rpx;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  border-bottom-left-radius: 8rpx;
}

.remove-icon {
  font-size: 28rpx;
  color: #ffffff;
  font-weight: bold;
}

.upload-btn {
  width: 180rpx;
  height: 180rpx;
  border-radius: 12rpx;
  border: 2rpx dashed #cccccc;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background-color: #fafafa;
}

.upload-icon {
  font-size: 48rpx;
  color: #cccccc;
  margin-bottom: 8rpx;
}

.upload-text {
  font-size: 22rpx;
  color: #cccccc;
}

.upload-hint {
  font-size: 22rpx;
  color: #999999;
  margin-top: 12rpx;
  display: block;
}

.question-info {
  padding: 20rpx 30rpx 24rpx;
}

.question-id {
  font-size: 24rpx;
  color: #666666;
}

.submit-section {
  padding: 30rpx 20rpx;
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
</style>
