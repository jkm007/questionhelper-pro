<template>
  <view class="password-page">
    <view class="form-card">
      <view class="form-field">
        <text class="field-label">Current Password</text>
        <view class="input-wrap">
          <input
            class="field-input"
            :password="!showOld"
            v-model="form.oldPassword"
            placeholder="Enter current password"
          />
          <view class="eye-btn" @tap="showOld = !showOld">
            <text class="eye-text">{{ showOld ? 'Hide' : 'Show' }}</text>
          </view>
        </view>
      </view>

      <view class="form-field">
        <text class="field-label">New Password</text>
        <view class="input-wrap">
          <input
            class="field-input"
            :password="!showNew"
            v-model="form.newPassword"
            placeholder="Enter new password (6-20 characters)"
          />
          <view class="eye-btn" @tap="showNew = !showNew">
            <text class="eye-text">{{ showNew ? 'Hide' : 'Show' }}</text>
          </view>
        </view>
        <!-- Password Strength -->
        <view v-if="form.newPassword" class="strength-bar">
          <view
            class="strength-segment"
            :class="passwordStrength >= 1 ? strengthColor : 'empty'"
          />
          <view
            class="strength-segment"
            :class="passwordStrength >= 2 ? strengthColor : 'empty'"
          />
          <view
            class="strength-segment"
            :class="passwordStrength >= 3 ? strengthColor : 'empty'"
          />
          <text class="strength-text">{{ strengthText }}</text>
        </view>
      </view>

      <view class="form-field">
        <text class="field-label">Confirm New Password</text>
        <view class="input-wrap">
          <input
            class="field-input"
            :password="!showConfirm"
            v-model="form.confirmPassword"
            placeholder="Re-enter new password"
          />
          <view class="eye-btn" @tap="showConfirm = !showConfirm">
            <text class="eye-text">{{ showConfirm ? 'Hide' : 'Show' }}</text>
          </view>
        </view>
        <text v-if="confirmError" class="error-text">Passwords do not match</text>
      </view>
    </view>

    <!-- Password Rules -->
    <view class="rules-card">
      <text class="rules-title">Password Requirements</text>
      <view class="rule-item">
        <view class="rule-dot" :class="{ valid: hasMinLength }" />
        <text class="rule-text" :class="{ valid: hasMinLength }">6-20 characters long</text>
      </view>
      <view class="rule-item">
        <view class="rule-dot" :class="{ valid: hasLetter }" />
        <text class="rule-text" :class="{ valid: hasLetter }">Contains at least one letter</text>
      </view>
      <view class="rule-item">
        <view class="rule-dot" :class="{ valid: hasNumber }" />
        <text class="rule-text" :class="{ valid: hasNumber }">Contains at least one number</text>
      </view>
      <view class="rule-item">
        <view class="rule-dot" :class="{ valid: noRepeats }" />
        <text class="rule-text" :class="{ valid: noRepeats }">No 3+ repeated characters</text>
      </view>
    </view>

    <!-- Submit -->
    <view class="submit-btn" :class="{ disabled: !canSubmit || submitting }" @tap="handleSubmit">
      <text class="submit-text">{{ submitting ? 'Updating...' : 'Update Password' }}</text>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { updatePassword } from '@/api/user'

const form = ref({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const showOld = ref(false)
const showNew = ref(false)
const showConfirm = ref(false)
const submitting = ref(false)

const hasMinLength = computed(() => {
  const len = form.value.newPassword.length
  return len >= 6 && len <= 20
})

const hasLetter = computed(() => /[a-zA-Z]/.test(form.value.newPassword))

const hasNumber = computed(() => /[0-9]/.test(form.value.newPassword))

const noRepeats = computed(() => !/(.)\1{2,}/.test(form.value.newPassword))

const passwordStrength = computed(() => {
  const pwd = form.value.newPassword
  if (!pwd) return 0
  let score = 0
  if (pwd.length >= 8) score++
  if (/[a-z]/.test(pwd) && /[A-Z]/.test(pwd)) score++
  if (/[0-9]/.test(pwd) && /[^a-zA-Z0-9]/.test(pwd)) score++
  return score
})

const strengthColor = computed(() => {
  const map: Record<number, string> = { 1: 'weak', 2: 'medium', 3: 'strong' }
  return map[passwordStrength.value] || 'weak'
})

const strengthText = computed(() => {
  const map: Record<number, string> = { 1: 'Weak', 2: 'Fair', 3: 'Strong' }
  return map[passwordStrength.value] || ''
})

const confirmError = computed(() => {
  return (
    form.value.confirmPassword.length > 0 &&
    form.value.newPassword !== form.value.confirmPassword
  )
})

const canSubmit = computed(() => {
  return (
    form.value.oldPassword.length > 0 &&
    hasMinLength.value &&
    hasLetter.value &&
    hasNumber.value &&
    noRepeats.value &&
    form.value.newPassword === form.value.confirmPassword &&
    form.value.confirmPassword.length > 0
  )
})

async function handleSubmit() {
  if (!canSubmit.value || submitting.value) return

  if (form.value.oldPassword === form.value.newPassword) {
    uni.showToast({ title: 'New password must differ from old', icon: 'none' })
    return
  }

  submitting.value = true
  try {
    await updatePassword({
      oldPassword: form.value.oldPassword,
      newPassword: form.value.newPassword
    })
    uni.showToast({ title: 'Password updated', icon: 'success' })
    setTimeout(() => {
      uni.navigateBack()
    }, 1500)
  } catch (e: any) {
    console.error('Failed to update password', e)
    uni.showToast({ title: e?.message || 'Update failed', icon: 'none' })
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.password-page {
  min-height: 100vh;
  background-color: #f5f5f5;
  padding: 24rpx;
  padding-bottom: 60rpx;
}

.form-card {
  background-color: #fff;
  border-radius: 16rpx;
  padding: 30rpx;
  margin-bottom: 24rpx;
}

.form-field {
  margin-bottom: 32rpx;
}

.form-field:last-child {
  margin-bottom: 0;
}

.field-label {
  font-size: 28rpx;
  color: #333;
  font-weight: 600;
  margin-bottom: 16rpx;
  display: block;
}

.input-wrap {
  display: flex;
  flex-direction: row;
  align-items: center;
  border: 2rpx solid #eee;
  border-radius: 12rpx;
  overflow: hidden;
}

.field-input {
  flex: 1;
  height: 88rpx;
  padding: 0 20rpx;
  font-size: 28rpx;
  color: #333;
}

.eye-btn {
  padding: 0 20rpx;
  height: 88rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.eye-text {
  font-size: 24rpx;
  color: #4a90d9;
}

.strength-bar {
  display: flex;
  flex-direction: row;
  align-items: center;
  margin-top: 12rpx;
  gap: 8rpx;
}

.strength-segment {
  flex: 1;
  height: 8rpx;
  border-radius: 4rpx;
  background-color: #eee;
}

.strength-segment.weak {
  background-color: #e74c3c;
}

.strength-segment.medium {
  background-color: #ff9800;
}

.strength-segment.strong {
  background-color: #67c23a;
}

.strength-segment.empty {
  background-color: #eee;
}

.strength-text {
  font-size: 22rpx;
  color: #999;
  margin-left: 8rpx;
}

.error-text {
  font-size: 24rpx;
  color: #e74c3c;
  margin-top: 8rpx;
  display: block;
}

.rules-card {
  background-color: #fff;
  border-radius: 16rpx;
  padding: 30rpx;
  margin-bottom: 40rpx;
}

.rules-title {
  font-size: 28rpx;
  color: #333;
  font-weight: 600;
  margin-bottom: 20rpx;
  display: block;
}

.rule-item {
  display: flex;
  flex-direction: row;
  align-items: center;
  margin-bottom: 14rpx;
}

.rule-item:last-child {
  margin-bottom: 0;
}

.rule-dot {
  width: 16rpx;
  height: 16rpx;
  border-radius: 8rpx;
  background-color: #ddd;
  margin-right: 14rpx;
  flex-shrink: 0;
}

.rule-dot.valid {
  background-color: #67c23a;
}

.rule-text {
  font-size: 26rpx;
  color: #999;
}

.rule-text.valid {
  color: #67c23a;
}

.submit-btn {
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
