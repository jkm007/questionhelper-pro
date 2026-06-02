<template>
  <view class="apply-page">
    <!-- Status Display -->
    <view v-if="application.status !== 'none'" class="status-card" :class="application.status">
      <text class="status-icon">{{ statusIcon }}</text>
      <text class="status-text">{{ statusText }}</text>
      <text v-if="application.rejectReason" class="status-reason">原因：{{ application.rejectReason }}</text>
    </view>

    <!-- Application Form -->
    <view v-if="application.status === 'none' || application.status === 'rejected'" class="form-card">
      <view class="form-item">
        <text class="form-label">申请理由</text>
        <textarea
          class="form-textarea"
          v-model="form.reason"
          placeholder="请说明你申请成为教师的理由"
          :maxlength="500"
        />
        <text class="char-count">{{ form.reason.length }}/500</text>
      </view>

      <view class="form-item">
        <text class="form-label">资质说明</text>
        <input class="form-input" v-model="form.qualifications" placeholder="如：计算机科学硕士、5年教学经验" />
      </view>

      <view class="form-item">
        <text class="form-label">证书上传</text>
        <view class="upload-area" @tap="chooseCertificate">
          <image v-if="form.certificate" class="cert-preview" :src="form.certificate" mode="aspectFill" />
          <view v-else class="upload-placeholder">
            <text class="upload-icon">+</text>
            <text class="upload-hint">点击上传证书照片</text>
          </view>
        </view>
      </view>

      <view class="submit-btn" :class="{ disabled: submitting }" @tap="handleSubmit">
        <text class="submit-text">{{ submitting ? '提交中...' : '提交申请' }}</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'

const submitting = ref(false)

const form = ref({
  reason: '',
  qualifications: '',
  certificate: ''
})

const application = ref({
  status: 'none' as 'none' | 'pending' | 'approved' | 'rejected',
  rejectReason: ''
})

const statusIcon = computed(() => {
  const map: Record<string, string> = { pending: '⏳', approved: '✅', rejected: '❌' }
  return map[application.value.status] || ''
})

const statusText = computed(() => {
  const map: Record<string, string> = { pending: '审核中，请耐心等待', approved: '已通过', rejected: '未通过，可重新申请' }
  return map[application.value.status] || ''
})

onMounted(() => { fetchStatus() })

async function fetchStatus() {
  try {
    // TODO: replace with actual API call
    // application.value = { status: 'pending', rejectReason: '' }
  } catch (e) {
    console.error('Failed to fetch application status', e)
  }
}

function chooseCertificate() {
  uni.chooseImage({
    count: 1,
    sizeType: ['compressed'],
    sourceType: ['album', 'camera'],
    success: (res) => {
      form.value.certificate = res.tempFilePaths[0]
    }
  })
}

async function handleSubmit() {
  if (!form.value.reason.trim()) {
    uni.showToast({ title: '请输入申请理由', icon: 'none' })
    return
  }
  if (!form.value.qualifications.trim()) {
    uni.showToast({ title: '请输入资质说明', icon: 'none' })
    return
  }
  submitting.value = true
  try {
    // TODO: replace with actual API call
    application.value = { status: 'pending', rejectReason: '' }
    uni.showToast({ title: '申请已提交', icon: 'success' })
  } catch (e) {
    console.error('Failed to submit application', e)
    uni.showToast({ title: '提交失败', icon: 'none' })
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.apply-page {
  min-height: 100vh;
  background-color: #f5f5f5;
  padding: 20rpx 24rpx;
  padding-bottom: 80rpx;
}

.status-card {
  border-radius: 16rpx;
  padding: 36rpx 30rpx;
  margin-bottom: 24rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.status-card.pending {
  background-color: #fff8e1;
}

.status-card.approved {
  background-color: #e8f5e9;
}

.status-card.rejected {
  background-color: #ffebee;
}

.status-icon {
  font-size: 60rpx;
  margin-bottom: 12rpx;
}

.status-text {
  font-size: 30rpx;
  font-weight: 600;
  color: #333;
  margin-bottom: 8rpx;
}

.status-reason {
  font-size: 24rpx;
  color: #999;
}

.form-card {
  background-color: #fff;
  border-radius: 16rpx;
  padding: 12rpx 30rpx;
}

.form-item {
  padding: 28rpx 0;
  border-bottom: 1rpx solid #f5f5f5;
}

.form-item:last-child {
  border-bottom: none;
}

.form-label {
  font-size: 26rpx;
  color: #999;
  margin-bottom: 16rpx;
}

.form-input {
  font-size: 30rpx;
  color: #333;
  padding: 8rpx 0;
}

.form-textarea {
  width: 100%;
  height: 200rpx;
  font-size: 28rpx;
  color: #333;
  padding: 16rpx;
  background-color: #f9f9f9;
  border-radius: 12rpx;
}

.char-count {
  display: block;
  text-align: right;
  font-size: 22rpx;
  color: #ccc;
  margin-top: 8rpx;
}

.upload-area {
  width: 280rpx;
  height: 200rpx;
  border: 2rpx dashed #ddd;
  border-radius: 12rpx;
  overflow: hidden;
}

.cert-preview {
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
}

.upload-icon {
  font-size: 60rpx;
  color: #ccc;
}

.upload-hint {
  font-size: 22rpx;
  color: #ccc;
  margin-top: 8rpx;
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
  opacity: 0.6;
}

.submit-text {
  font-size: 30rpx;
  color: #fff;
  font-weight: 600;
}
</style>
