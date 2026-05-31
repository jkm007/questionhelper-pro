<template>
  <view class="apply-page">
    <!-- Status Display -->
    <view v-if="applyStatus === 'approved'" class="status-card approved">
      <image class="status-icon" src="/static/icon-success.png" mode="aspectFit" />
      <text class="status-title">You are a Creator</text>
      <text class="status-desc">Congratulations! Your creator application has been approved. You can now create and publish content.</text>
    </view>

    <view v-else-if="applyStatus === 'pending'" class="status-card pending">
      <image class="status-icon" src="/static/icon-pending.png" mode="aspectFit" />
      <text class="status-title">Application Under Review</text>
      <text class="status-desc">Your application is being reviewed. We will notify you of the result within 3 business days.</text>
    </view>

    <view v-else-if="applyStatus === 'rejected'" class="status-card rejected">
      <image class="status-icon" src="/static/icon-rejected.png" mode="aspectFit" />
      <text class="status-title">Application Rejected</text>
      <text class="status-desc">Reason: {{ rejectReason || 'Your application did not meet the requirements' }}</text>
      <view class="retry-btn" @tap="resetForm">
        <text class="retry-text">Reapply</text>
      </view>
    </view>

    <!-- Apply Form -->
    <view v-if="applyStatus !== 'approved' && applyStatus !== 'pending'" class="form-section">
      <view class="intro-card">
        <text class="intro-title">Creator Benefits</text>
        <view class="benefit-list">
          <view class="benefit-item">
            <text class="benefit-dot">.</text>
            <text class="benefit-text">Create and publish questions, papers, and exams</text>
          </view>
          <view class="benefit-item">
            <text class="benefit-dot">.</text>
            <text class="benefit-text">Earn income from content usage</text>
          </view>
          <view class="benefit-item">
            <text class="benefit-dot">.</text>
            <text class="benefit-text">Get a verified creator badge</text>
          </view>
          <view class="benefit-item">
            <text class="benefit-dot">.</text>
            <text class="benefit-text">Access advanced creation tools</text>
          </view>
        </view>
      </view>

      <view class="form-group">
        <text class="group-title">Application Reason</text>
        <textarea
          class="form-textarea"
          v-model="form.reason"
          placeholder="Please describe why you want to become a creator, your background and expertise..."
          maxlength="500"
        />
        <text class="char-count">{{ (form.reason || '').length }}/500</text>
      </view>

      <view class="form-group">
        <text class="group-title">Portfolio Description</text>
        <textarea
          class="form-textarea"
          v-model="form.portfolio"
          placeholder="Describe your relevant experience, published works, certifications, or other supporting materials..."
          maxlength="500"
        />
        <text class="char-count">{{ (form.portfolio || '').length }}/500</text>
      </view>

      <!-- Agreement -->
      <view class="agreement-wrap" @tap="toggleAgreement">
        <view class="checkbox" :class="{ checked: agreed }">
          <text v-if="agreed" class="check-mark">OK</text>
        </view>
        <text class="agreement-text">I have read and agree to the </text>
        <text class="agreement-link" @tap.stop="viewTerms">Creator Agreement</text>
      </view>

      <!-- Submit Button -->
      <view class="submit-btn" :class="{ disabled: !canSubmit || submitting }" @tap="handleSubmit">
        <text class="submit-text">{{ submitting ? 'Submitting...' : 'Submit Application' }}</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { getCreatorStatus, submitCreatorApply } from '@/api/user'

const applyStatus = ref<'none' | 'pending' | 'approved' | 'rejected'>('none')
const rejectReason = ref('')
const submitting = ref(false)
const agreed = ref(false)

const form = ref({
  reason: '',
  portfolio: ''
})

const canSubmit = computed(() => {
  return (
    form.value.reason.trim().length >= 10 &&
    form.value.portfolio.trim().length >= 10 &&
    agreed.value
  )
})

onMounted(() => {
  fetchApplyStatus()
})

async function fetchApplyStatus() {
  try {
    const res = await getCreatorStatus()
    if (res.data) {
      applyStatus.value = res.data.status || 'none'
      rejectReason.value = res.data.rejectReason || ''
    }
  } catch (e) {
    console.error('Failed to fetch creator status', e)
  }
}

function resetForm() {
  applyStatus.value = 'none'
  form.value = { reason: '', portfolio: '' }
  agreed.value = false
}

function toggleAgreement() {
  agreed.value = !agreed.value
}

function viewTerms() {
  uni.navigateTo({ url: '/pages/webview/index?url=/creator-agreement' })
}

async function handleSubmit() {
  if (!canSubmit.value || submitting.value) return

  submitting.value = true
  try {
    await submitCreatorApply({
      reason: form.value.reason.trim(),
      portfolio: form.value.portfolio.trim()
    })
    applyStatus.value = 'pending'
    uni.showToast({ title: 'Application submitted', icon: 'success' })
  } catch (e) {
    console.error('Submit failed', e)
    uni.showToast({ title: 'Submission failed', icon: 'none' })
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.apply-page {
  min-height: 100vh;
  background-color: #f5f5f5;
  padding: 24rpx;
}

.status-card {
  background-color: #fff;
  border-radius: 16rpx;
  padding: 40rpx 30rpx;
  align-items: center;
  margin-bottom: 24rpx;
}

.status-icon {
  width: 100rpx;
  height: 100rpx;
  margin-bottom: 20rpx;
}

.status-title {
  font-size: 34rpx;
  color: #333;
  font-weight: 700;
  margin-bottom: 12rpx;
}

.status-desc {
  font-size: 26rpx;
  color: #999;
  text-align: center;
  line-height: 1.5;
}

.status-card.approved {
  border-top: 6rpx solid #4caf50;
}

.status-card.pending {
  border-top: 6rpx solid #ff9800;
}

.status-card.rejected {
  border-top: 6rpx solid #e74c3c;
}

.retry-btn {
  margin-top: 30rpx;
  background-color: #4a90d9;
  border-radius: 40rpx;
  padding: 20rpx 60rpx;
}

.retry-text {
  font-size: 28rpx;
  color: #fff;
}

.intro-card {
  background-color: #fff;
  border-radius: 16rpx;
  padding: 30rpx;
  margin-bottom: 24rpx;
}

.intro-title {
  font-size: 30rpx;
  color: #333;
  font-weight: 600;
  margin-bottom: 16rpx;
}

.benefit-list {
  padding-left: 8rpx;
}

.benefit-item {
  display: flex;
  flex-direction: row;
  align-items: flex-start;
  margin-bottom: 12rpx;
}

.benefit-dot {
  font-size: 28rpx;
  color: #4a90d9;
  margin-right: 12rpx;
  line-height: 1.5;
}

.benefit-text {
  font-size: 26rpx;
  color: #666;
  line-height: 1.5;
}

.form-group {
  background-color: #fff;
  border-radius: 16rpx;
  padding: 24rpx 30rpx;
  margin-bottom: 24rpx;
}

.group-title {
  font-size: 30rpx;
  color: #333;
  font-weight: 600;
  margin-bottom: 16rpx;
}

.form-textarea {
  width: 100%;
  min-height: 200rpx;
  font-size: 28rpx;
  color: #333;
  line-height: 1.6;
}

.char-count {
  font-size: 22rpx;
  color: #999;
  text-align: right;
  margin-top: 8rpx;
}

.agreement-wrap {
  display: flex;
  flex-direction: row;
  align-items: center;
  padding: 20rpx 0;
}

.checkbox {
  width: 36rpx;
  height: 36rpx;
  border: 2rpx solid #ccc;
  border-radius: 6rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 12rpx;
}

.checkbox.checked {
  background-color: #4a90d9;
  border-color: #4a90d9;
}

.check-mark {
  font-size: 22rpx;
  color: #fff;
}

.agreement-text {
  font-size: 26rpx;
  color: #666;
}

.agreement-link {
  font-size: 26rpx;
  color: #4a90d9;
}

.submit-btn {
  margin-top: 40rpx;
  background-color: #4a90d9;
  border-radius: 44rpx;
  padding: 26rpx 0;
  display: flex;
  align-items: center;
  justify-content: center;
}

.submit-btn.disabled {
  opacity: 0.5;
}

.submit-text {
  font-size: 32rpx;
  color: #fff;
  font-weight: 600;
}
</style>
