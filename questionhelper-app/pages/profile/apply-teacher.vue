<template>
  <view class="apply-teacher-page">
    <!-- Status Display -->
    <view v-if="applyStatus === 'approved'" class="status-card approved">
      <image class="status-icon" src="/static/icon-success.png" mode="aspectFit" />
      <text class="status-title">Application Approved</text>
      <text class="status-desc">Congratulations! You are now a certified teacher. You can create classes and assign homework.</text>
    </view>

    <view v-else-if="applyStatus === 'pending'" class="status-card pending">
      <image class="status-icon" src="/static/icon-pending.png" mode="aspectFit" />
      <text class="status-title">Under Review</text>
      <text class="status-desc">Your teacher application is being reviewed. We will notify you within 3-5 business days.</text>
    </view>

    <view v-else-if="applyStatus === 'rejected'" class="status-card rejected">
      <image class="status-icon" src="/static/icon-rejected.png" mode="aspectFit" />
      <text class="status-title">Application Rejected</text>
      <text class="status-desc">Reason: {{ rejectReason || 'Did not meet requirements' }}</text>
      <view class="retry-btn" @tap="resetForm">
        <text class="retry-text">Reapply</text>
      </view>
    </view>

    <!-- Application Form -->
    <view v-if="applyStatus !== 'approved' && applyStatus !== 'pending'" class="form-section">
      <!-- Benefits -->
      <view class="benefits-card">
        <text class="benefits-title">Teacher Benefits</text>
        <view class="benefit-list">
          <view class="benefit-item">
            <text class="benefit-dot">.</text>
            <text class="benefit-text">Create and manage classes</text>
          </view>
          <view class="benefit-item">
            <text class="benefit-dot">.</text>
            <text class="benefit-text">Assign homework and exams to students</text>
          </view>
          <view class="benefit-item">
            <text class="benefit-dot">.</text>
            <text class="benefit-text">Track student progress and analytics</text>
          </view>
          <view class="benefit-item">
            <text class="benefit-dot">.</text>
            <text class="benefit-text">Verified teacher badge on your profile</text>
          </view>
        </view>
      </view>

      <!-- Reason -->
      <view class="form-group">
        <text class="group-title">Application Reason</text>
        <textarea
          class="form-textarea"
          v-model="form.reason"
          placeholder="Describe why you want to become a teacher, your teaching experience, and goals..."
          maxlength="500"
        />
        <text class="char-count">{{ (form.reason || '').length }}/500</text>
      </view>

      <!-- Qualifications -->
      <view class="form-group">
        <text class="group-title">Qualifications</text>
        <textarea
          class="form-textarea"
          v-model="form.qualifications"
          placeholder="List your educational background, degrees, certifications, teaching years, subject expertise..."
          maxlength="500"
        />
        <text class="char-count">{{ (form.qualifications || '').length }}/500</text>
      </view>

      <!-- School / Institution -->
      <view class="form-group">
        <text class="group-title">School / Institution</text>
        <input
          class="form-input"
          v-model="form.school"
          placeholder="Enter your school or institution name"
        />
      </view>

      <!-- Subject -->
      <view class="form-group">
        <text class="group-title">Teaching Subject</text>
        <view class="subject-options">
          <view
            v-for="subject in subjects"
            :key="subject"
            class="subject-chip"
            :class="{ selected: form.subject === subject }"
            @tap="form.subject = subject"
          >
            <text class="chip-text">{{ subject }}</text>
          </view>
        </view>
      </view>

      <!-- Certificate Upload -->
      <view class="form-group">
        <text class="group-title">Certificates (optional)</text>
        <text class="group-hint">Upload teaching certificates, degree photos, or other supporting documents</text>
        <view class="upload-grid">
          <view
            v-for="(cert, index) in certificates"
            :key="index"
            class="upload-item"
          >
            <image class="upload-preview" :src="cert" mode="aspectFill" @tap="previewImage(index)" />
            <view class="upload-remove" @tap="removeCert(index)">
              <text class="remove-icon">×</text>
            </view>
          </view>
          <view v-if="certificates.length < 5" class="upload-add" @tap="chooseCert">
            <text class="add-icon">+</text>
            <text class="add-text">Add Photo</text>
          </view>
        </view>
      </view>

      <!-- Agreement -->
      <view class="agreement-wrap" @tap="toggleAgreement">
        <view class="checkbox" :class="{ checked: agreed }">
          <text v-if="agreed" class="check-mark">OK</text>
        </view>
        <text class="agreement-text">I confirm that the information provided is true and accurate</text>
      </view>

      <!-- Submit -->
      <view class="submit-btn" :class="{ disabled: !canSubmit || submitting }" @tap="handleSubmit">
        <text class="submit-text">{{ submitting ? 'Submitting...' : 'Submit Application' }}</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { BASE_URL } from '@/api/request'

const subjects = ['Math', 'English', 'Physics', 'Chemistry', 'Biology', 'Chinese', 'History', 'Geography', 'Other']

const applyStatus = ref<'none' | 'pending' | 'approved' | 'rejected'>('none')
const rejectReason = ref('')
const submitting = ref(false)
const agreed = ref(false)
const certificates = ref<string[]>([])

const form = ref({
  reason: '',
  qualifications: '',
  school: '',
  subject: ''
})

const canSubmit = computed(() => {
  return (
    form.value.reason.trim().length >= 10 &&
    form.value.qualifications.trim().length >= 10 &&
    form.value.school.trim().length > 0 &&
    form.value.subject.length > 0 &&
    agreed.value
  )
})

onMounted(() => {
  fetchApplyStatus()
})

async function fetchApplyStatus() {
  try {
    const res = await new Promise<any>((resolve) => {
      uni.request({
        url: '/api/v1/user/teacher/status',
        method: 'GET',
        header: { Authorization: `Bearer ${uni.getStorageSync('token')}` },
        success: (r: any) => resolve(r.data),
        fail: () => resolve(null)
      })
    })
    if (res?.code === '00000' && res.data) {
      applyStatus.value = res.data.status || 'none'
      rejectReason.value = res.data.rejectReason || ''
    }
  } catch (e) {
    console.error('Failed to fetch teacher status', e)
  }
}

function resetForm() {
  applyStatus.value = 'none'
  form.value = { reason: '', qualifications: '', school: '', subject: '' }
  certificates.value = []
  agreed.value = false
}

function toggleAgreement() {
  agreed.value = !agreed.value
}

function chooseCert() {
  uni.chooseImage({
    count: 5 - certificates.value.length,
    sizeType: ['compressed'],
    sourceType: ['album', 'camera'],
    success: (res) => {
      certificates.value.push(...res.tempFilePaths)
    }
  })
}

function removeCert(index: number) {
  certificates.value.splice(index, 1)
}

function previewImage(index: number) {
  uni.previewImage({
    current: index,
    urls: certificates.value
  })
}

async function uploadCertificates(): Promise<string[]> {
  if (!certificates.value.length) return []

  const urls: string[] = []
  const token = uni.getStorageSync('token')

  for (const filePath of certificates.value) {
    try {
      const res = await new Promise<any>((resolve, reject) => {
        uni.uploadFile({
          url: `${BASE_URL}/upload/image`,
          filePath,
          name: 'file',
          header: { Authorization: `Bearer ${token}` },
          success: (r) => {
            const data = JSON.parse(r.data)
            if (data.code === '00000') resolve(data)
            else reject(new Error(data.msg))
          },
          fail: reject
        })
      })
      if (res.data?.url) {
        urls.push(res.data.url)
      }
    } catch (e) {
      console.error('Upload failed for', filePath, e)
    }
  }

  return urls
}

async function handleSubmit() {
  if (!canSubmit.value || submitting.value) return

  submitting.value = true
  try {
    uni.showLoading({ title: 'Submitting...' })

    const certUrls = await uploadCertificates()

    await new Promise<any>((resolve, reject) => {
      uni.request({
        url: '/api/v1/user/teacher/apply',
        method: 'POST',
        data: {
          reason: form.value.reason.trim(),
          qualifications: form.value.qualifications.trim(),
          school: form.value.school.trim(),
          subject: form.value.subject,
          certificates: certUrls
        },
        header: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${uni.getStorageSync('token')}`
        },
        success: (r: any) => {
          if (r.data?.code === '00000') resolve(r.data)
          else reject(new Error(r.data?.msg))
        },
        fail: reject
      })
    })

    uni.hideLoading()
    applyStatus.value = 'pending'
    uni.showToast({ title: 'Application submitted', icon: 'success' })
  } catch (e: any) {
    uni.hideLoading()
    console.error('Submit failed', e)
    uni.showToast({ title: e?.message || 'Submission failed', icon: 'none' })
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.apply-teacher-page {
  min-height: 100vh;
  background-color: #f5f5f5;
  padding: 24rpx;
  padding-bottom: 60rpx;
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

.form-section {
  margin-bottom: 24rpx;
}

.benefits-card {
  background-color: #fff;
  border-radius: 16rpx;
  padding: 30rpx;
  margin-bottom: 24rpx;
}

.benefits-title {
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

.group-hint {
  font-size: 24rpx;
  color: #999;
  margin-bottom: 16rpx;
  display: block;
}

.form-textarea {
  width: 100%;
  min-height: 180rpx;
  font-size: 28rpx;
  color: #333;
  line-height: 1.6;
}

.form-input {
  width: 100%;
  height: 80rpx;
  font-size: 28rpx;
  color: #333;
}

.char-count {
  font-size: 22rpx;
  color: #999;
  text-align: right;
  margin-top: 8rpx;
}

.subject-options {
  display: flex;
  flex-direction: row;
  flex-wrap: wrap;
  gap: 16rpx;
}

.subject-chip {
  padding: 14rpx 28rpx;
  border: 2rpx solid #eee;
  border-radius: 30rpx;
  background-color: #fafafa;
}

.subject-chip.selected {
  background-color: #e8f0fe;
  border-color: #4a90d9;
}

.chip-text {
  font-size: 26rpx;
  color: #666;
}

.subject-chip.selected .chip-text {
  color: #4a90d9;
  font-weight: 600;
}

.upload-grid {
  display: flex;
  flex-direction: row;
  flex-wrap: wrap;
  gap: 16rpx;
}

.upload-item {
  width: 180rpx;
  height: 180rpx;
  border-radius: 12rpx;
  overflow: hidden;
  position: relative;
}

.upload-preview {
  width: 100%;
  height: 100%;
}

.upload-remove {
  position: absolute;
  top: 0;
  right: 0;
  width: 40rpx;
  height: 40rpx;
  background-color: rgba(0, 0, 0, 0.5);
  border-radius: 0 0 0 12rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.remove-icon {
  font-size: 28rpx;
  color: #fff;
  line-height: 1;
}

.upload-add {
  width: 180rpx;
  height: 180rpx;
  border: 2rpx dashed #ccc;
  border-radius: 12rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background-color: #fafafa;
}

.add-icon {
  font-size: 48rpx;
  color: #ccc;
  line-height: 1;
}

.add-text {
  font-size: 22rpx;
  color: #ccc;
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
  flex: 1;
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
