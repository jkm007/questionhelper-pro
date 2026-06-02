<template>
  <view class="password-page">
    <view class="form-card">
      <!-- Old Password -->
      <view class="form-item">
        <text class="form-label">原密码</text>
        <view class="input-wrap">
          <input
            class="form-input"
            :password="!showOld"
            v-model="form.oldPassword"
            placeholder="请输入原密码"
          />
          <text class="toggle-eye" @tap="showOld = !showOld">{{ showOld ? '隐藏' : '显示' }}</text>
        </view>
      </view>

      <!-- New Password -->
      <view class="form-item">
        <text class="form-label">新密码</text>
        <view class="input-wrap">
          <input
            class="form-input"
            :password="!showNew"
            v-model="form.newPassword"
            placeholder="请输入新密码"
          />
          <text class="toggle-eye" @tap="showNew = !showNew">{{ showNew ? '隐藏' : '显示' }}</text>
        </view>
        <!-- Strength Indicator -->
        <view class="strength-bar">
          <view class="strength-segment" :class="strengthClass(1)" />
          <view class="strength-segment" :class="strengthClass(2)" />
          <view class="strength-segment" :class="strengthClass(3)" />
          <text class="strength-text">{{ strengthText }}</text>
        </view>
      </view>

      <!-- Confirm Password -->
      <view class="form-item">
        <text class="form-label">确认密码</text>
        <view class="input-wrap">
          <input
            class="form-input"
            :password="!showConfirm"
            v-model="form.confirmPassword"
            placeholder="请再次输入新密码"
          />
          <text class="toggle-eye" @tap="showConfirm = !showConfirm">{{ showConfirm ? '隐藏' : '显示' }}</text>
        </view>
      </view>
    </view>

    <!-- Submit -->
    <view class="submit-btn" :class="{ disabled: submitting }" @tap="handleSubmit">
      <text class="submit-text">{{ submitting ? '提交中...' : '确认修改' }}</text>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'

const showOld = ref(false)
const showNew = ref(false)
const showConfirm = ref(false)
const submitting = ref(false)

const form = ref({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

const passwordStrength = computed(() => {
  const pwd = form.value.newPassword
  if (!pwd) return 0
  let score = 0
  if (pwd.length >= 8) score++
  if (/[A-Z]/.test(pwd) && /[a-z]/.test(pwd)) score++
  if (/[0-9]/.test(pwd) && /[^A-Za-z0-9]/.test(pwd)) score++
  return score
})

const strengthText = computed(() => {
  const map: Record<number, string> = { 0: '', 1: '弱', 2: '中', 3: '强' }
  return map[passwordStrength.value] || ''
})

function strengthClass(level: number) {
  if (passwordStrength.value >= level) {
    return passwordStrength.value === 1 ? 'weak' : passwordStrength.value === 2 ? 'medium' : 'strong'
  }
  return ''
}

async function handleSubmit() {
  if (!form.value.oldPassword) {
    uni.showToast({ title: '请输入原密码', icon: 'none' })
    return
  }
  if (form.value.newPassword.length < 8) {
    uni.showToast({ title: '新密码至少8位', icon: 'none' })
    return
  }
  if (form.value.newPassword !== form.value.confirmPassword) {
    uni.showToast({ title: '两次密码不一致', icon: 'none' })
    return
  }
  submitting.value = true
  try {
    // TODO: call PUT /user/password
    uni.showToast({ title: '密码修改成功', icon: 'success' })
    setTimeout(() => uni.navigateBack(), 1500)
  } catch (e) {
    console.error('Failed to change password', e)
    uni.showToast({ title: '修改失败', icon: 'none' })
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.password-page {
  min-height: 100vh;
  background-color: #f5f5f5;
  padding: 20rpx 24rpx;
  padding-bottom: 80rpx;
}

.form-card {
  background-color: #fff;
  border-radius: 16rpx;
  padding: 12rpx 30rpx;
  margin-bottom: 40rpx;
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

.input-wrap {
  display: flex;
  flex-direction: row;
  align-items: center;
}

.form-input {
  flex: 1;
  font-size: 30rpx;
  color: #333;
  padding: 8rpx 0;
}

.toggle-eye {
  font-size: 24rpx;
  color: #4a90d9;
  padding: 8rpx 0 8rpx 16rpx;
}

.strength-bar {
  display: flex;
  flex-direction: row;
  align-items: center;
  gap: 10rpx;
  margin-top: 16rpx;
}

.strength-segment {
  flex: 1;
  height: 10rpx;
  background-color: #f0f0f0;
  border-radius: 5rpx;
}

.strength-segment.weak {
  background-color: #e74c3c;
}

.strength-segment.medium {
  background-color: #f39c12;
}

.strength-segment.strong {
  background-color: #27ae60;
}

.strength-text {
  font-size: 22rpx;
  color: #999;
  margin-left: 12rpx;
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
  opacity: 0.6;
}

.submit-text {
  font-size: 30rpx;
  color: #fff;
  font-weight: 600;
}
</style>
