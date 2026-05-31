<template>
  <view class="auth-page">
    <!-- Status Display (already authenticated) -->
    <view v-if="authStatus === 'approved'" class="status-card approved">
      <image class="status-icon" src="/static/icon-success.png" mode="aspectFit" />
      <text class="status-title">Authentication Approved</text>
      <text class="status-desc">Your real-name identity has been verified.</text>
      <view class="info-display">
        <view class="info-row">
          <text class="info-label">Real Name</text>
          <text class="info-value">{{ maskedName }}</text>
        </view>
        <view class="info-row">
          <text class="info-label">ID Number</text>
          <text class="info-value">{{ maskedIdNumber }}</text>
        </view>
      </view>
    </view>

    <view v-else-if="authStatus === 'pending'" class="status-card pending">
      <image class="status-icon" src="/static/icon-pending.png" mode="aspectFit" />
      <text class="status-title">Under Review</text>
      <text class="status-desc">Your authentication is being reviewed. Please wait patiently.</text>
    </view>

    <view v-else-if="authStatus === 'rejected'" class="status-card rejected">
      <image class="status-icon" src="/static/icon-rejected.png" mode="aspectFit" />
      <text class="status-title">Authentication Rejected</text>
      <text class="status-desc">Reason: {{ rejectReason || 'Information does not meet requirements' }}</text>
      <view class="retry-btn" @tap="resetForm">
        <text class="retry-text">Resubmit</text>
      </view>
    </view>

    <!-- Auth Form -->
    <view v-if="authStatus !== 'approved' && authStatus !== 'pending'" class="form-section">
      <view class="form-group">
        <view class="form-item">
          <text class="form-label">Real Name</text>
          <input
            class="form-input"
            v-model="form.realName"
            placeholder="Enter your real name"
            maxlength="30"
          />
        </view>
        <view class="form-item">
          <text class="form-label">ID Number</text>
          <input
            class="form-input"
            v-model="form.idNumber"
            placeholder="Enter your ID card number"
            maxlength="18"
            type="idcard"
          />
        </view>
      </view>

      <!-- ID Card Upload -->
      <view class="form-group">
        <text class="group-title">ID Card Photos</text>
        <view class="upload-row">
          <view class="upload-item" @tap="chooseImage('front')">
            <image
              v-if="form.idCardFront"
              class="upload-preview"
              :src="form.idCardFront"
              mode="aspectFill"
            />
            <view v-else class="upload-placeholder">
              <image class="upload-icon" src="/static/icon-upload.png" mode="aspectFit" />
              <text class="upload-text">Front Side</text>
            </view>
          </view>
          <view class="upload-item" @tap="chooseImage('back')">
            <image
              v-if="form.idCardBack"
              class="upload-preview"
              :src="form.idCardBack"
              mode="aspectFill"
            />
            <view v-else class="upload-placeholder">
              <image class="upload-icon" src="/static/icon-upload.png" mode="aspectFit" />
              <text class="upload-text">Back Side</text>
            </view>
          </view>
        </view>
        <text class="upload-hint">Please upload clear photos of your ID card (front and back)</text>
      </view>

      <!-- Submit Button -->
      <view class="submit-btn" :class="{ disabled: !canSubmit || submitting }" @tap="handleSubmit">
        <text class="submit-text">{{ submitting ? 'Submitting...' : 'Submit Authentication' }}</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { getAuthStatus, submitAuth } from '@/api/user'

const authStatus = ref<'none' | 'pending' | 'approved' | 'rejected'>('none')
const rejectReason = ref('')
const submitting = ref(false)

const form = ref({
  realName: '',
  idNumber: '',
  idCardFront: '',
  idCardBack: ''
})

const maskedName = computed(() => {
  const name = form.value.realName
  if (!name) return '***'
  if (name.length <= 1) return '*'
  return name[0] + '*'.repeat(name.length - 1)
})

const maskedIdNumber = computed(() => {
  const id = form.value.idNumber
  if (!id || id.length < 8) return '****'
  return id.slice(0, 4) + '****' + id.slice(-4)
})

const canSubmit = computed(() => {
  return (
    form.value.realName.trim() &&
    form.value.idNumber.trim() &&
    form.value.idCardFront &&
    form.value.idCardBack
  )
})

onMounted(() => {
  fetchAuthStatus()
})

async function fetchAuthStatus() {
  try {
    const res = await getAuthStatus()
    if (res.data) {
      authStatus.value = res.data.status || 'none'
      rejectReason.value = res.data.rejectReason || ''
      if (res.data.realName) {
        form.value.realName = res.data.realName
      }
      if (res.data.idNumber) {
        form.value.idNumber = res.data.idNumber
      }
      if (res.data.idCardFront) {
        form.value.idCardFront = res.data.idCardFront
      }
      if (res.data.idCardBack) {
        form.value.idCardBack = res.data.idCardBack
      }
    }
  } catch (e) {
    console.error('Failed to fetch auth status', e)
  }
}

function resetForm() {
  authStatus.value = 'none'
  form.value = {
    realName: '',
    idNumber: '',
    idCardFront: '',
    idCardBack: ''
  }
}

function chooseImage(side: 'front' | 'back') {
  uni.chooseImage({
    count: 1,
    sizeType: ['compressed'],
    sourceType: ['album', 'camera'],
    success: async (res) => {
      const filePath = res.tempFilePaths[0]
      try {
        uni.showLoading({ title: 'Uploading...' })
        const uploadRes = await uni.uploadFile({
          url: '/api/upload/image',
          filePath,
          name: 'file',
          header: {
            Authorization: `Bearer ${uni.getStorageSync('token')}`
          }
        })
        const data = JSON.parse(uploadRes.data)
        if (data.data?.url) {
          if (side === 'front') {
            form.value.idCardFront = data.data.url
          } else {
            form.value.idCardBack = data.data.url
          }
        }
        uni.hideLoading()
      } catch (e) {
        uni.hideLoading()
        console.error('Upload failed', e)
        uni.showToast({ title: 'Upload failed', icon: 'none' })
      }
    }
  })
}

async function handleSubmit() {
  if (!canSubmit.value || submitting.value) return

  // Validate ID number format
  const idRegex = /(^\d{15}$)|(^\d{18}$)|(^\d{17}(\d|X|x)$)/
  if (!idRegex.test(form.value.idNumber)) {
    uni.showToast({ title: 'Invalid ID number format', icon: 'none' })
    return
  }

  submitting.value = true
  try {
    await submitAuth({
      realName: form.value.realName.trim(),
      idNumber: form.value.idNumber.trim(),
      idCardFront: form.value.idCardFront,
      idCardBack: form.value.idCardBack
    })
    authStatus.value = 'pending'
    uni.showToast({ title: 'Submitted successfully', icon: 'success' })
  } catch (e) {
    console.error('Submit failed', e)
    uni.showToast({ title: 'Submission failed', icon: 'none' })
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.auth-page {
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

.info-display {
  width: 100%;
  margin-top: 30rpx;
}

.info-row {
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  padding: 16rpx 0;
  border-bottom: 1rpx solid #f0f0f0;
}

.info-label {
  font-size: 28rpx;
  color: #666;
}

.info-value {
  font-size: 28rpx;
  color: #333;
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

.form-group {
  background-color: #fff;
  border-radius: 16rpx;
  margin-bottom: 24rpx;
  padding: 20rpx 30rpx;
}

.group-title {
  font-size: 30rpx;
  color: #333;
  font-weight: 600;
  padding: 16rpx 0;
}

.form-item {
  display: flex;
  flex-direction: row;
  align-items: center;
  padding: 24rpx 0;
  border-bottom: 1rpx solid #f0f0f0;
}

.form-item:last-child {
  border-bottom: none;
}

.form-label {
  font-size: 30rpx;
  color: #333;
  width: 160rpx;
  flex-shrink: 0;
}

.form-input {
  flex: 1;
  font-size: 28rpx;
  color: #333;
}

.upload-row {
  display: flex;
  flex-direction: row;
  gap: 24rpx;
  padding: 16rpx 0;
}

.upload-item {
  flex: 1;
  aspect-ratio: 1.6;
  border-radius: 12rpx;
  overflow: hidden;
  border: 2rpx dashed #ddd;
}

.upload-preview {
  width: 100%;
  height: 100%;
}

.upload-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background-color: #fafafa;
}

.upload-icon {
  width: 60rpx;
  height: 60rpx;
  margin-bottom: 12rpx;
}

.upload-text {
  font-size: 24rpx;
  color: #999;
}

.upload-hint {
  font-size: 22rpx;
  color: #bbb;
  padding: 8rpx 0 16rpx;
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
