<template>
  <view class="forget-container">
    <!-- 步骤指示器 -->
    <view class="steps">
      <view class="step" :class="{ active: step >= 1, done: step > 1 }">
        <view class="step-dot">1</view>
        <text class="step-text">验证手机</text>
      </view>
      <view class="step-line" :class="{ active: step > 1 }"></view>
      <view class="step" :class="{ active: step >= 2, done: step > 2 }">
        <view class="step-dot">2</view>
        <text class="step-text">重置密码</text>
      </view>
      <view class="step-line" :class="{ active: step > 2 }"></view>
      <view class="step" :class="{ active: step >= 3 }">
        <view class="step-dot">3</view>
        <text class="step-text">完成</text>
      </view>
    </view>

    <!-- 步骤1: 验证手机 -->
    <view v-if="step === 1" class="form-section">
      <view class="section-title">
        <text class="title-text">验证手机号</text>
        <text class="title-desc">请输入注册时使用的手机号</text>
      </view>
      <view class="form-item">
        <input
          v-model="form.phone"
          class="form-input"
          placeholder="请输入手机号"
          type="number"
          maxlength="11"
        />
      </view>
      <view class="form-item form-item--row">
        <input
          v-model="form.code"
          class="form-input form-input--flex"
          placeholder="请输入验证码"
          type="number"
          maxlength="6"
        />
        <button
          class="sms-btn"
          :disabled="smsDisabled"
          @tap="handleSendCode"
        >
          {{ smsText }}
        </button>
      </view>
      <button class="submit-btn" :loading="loading" @tap="handleVerifyCode">
        下一步
      </button>
    </view>

    <!-- 步骤2: 重置密码 -->
    <view v-if="step === 2" class="form-section">
      <view class="section-title">
        <text class="title-text">重置密码</text>
        <text class="title-desc">请设置新密码</text>
      </view>
      <view class="form-item form-item--password">
        <input
          v-model="form.password"
          class="form-input form-input--flex"
          :password="!showPassword"
          placeholder="请输入新密码"
          maxlength="20"
        />
        <view class="toggle-eye" @tap="showPassword = !showPassword">
          <text>{{ showPassword ? '🙈' : '👁️' }}</text>
        </view>
      </view>
      <view class="form-item form-item--password">
        <input
          v-model="form.confirmPassword"
          class="form-input form-input--flex"
          :password="!showConfirmPassword"
          placeholder="请确认新密码"
          maxlength="20"
        />
        <view class="toggle-eye" @tap="showConfirmPassword = !showConfirmPassword">
          <text>{{ showConfirmPassword ? '🙈' : '👁️' }}</text>
        </view>
      </view>
      <button class="submit-btn" :loading="loading" @tap="handleResetPassword">
        确认重置
      </button>
    </view>

    <!-- 步骤3: 完成 -->
    <view v-if="step === 3" class="success-section">
      <view class="success-icon">
        <text class="icon-text">✓</text>
      </view>
      <text class="success-title">密码重置成功</text>
      <text class="success-desc">请使用新密码登录</text>
      <button class="submit-btn" @tap="goToLogin">
        去登录
      </button>
    </view>

    <!-- 返回登录 -->
    <view v-if="step < 3" class="form-footer">
      <text class="link" @tap="goToLogin">返回登录</text>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
// import { sendSmsCode, verifySmsCode, resetPassword } from '@/api/auth'

const step = ref(1)
const loading = ref(false)
const showPassword = ref(false)
const showConfirmPassword = ref(false)
const countdown = ref(0)
let timer: ReturnType<typeof setInterval> | null = null

const form = reactive({
  phone: '',
  code: '',
  password: '',
  confirmPassword: ''
})

const smsText = computed(() => {
  return countdown.value > 0 ? `${countdown.value}s后重新获取` : '获取验证码'
})

const smsDisabled = computed(() => {
  return countdown.value > 0 || !form.phone || form.phone.length !== 11
})

const startCountdown = () => {
  countdown.value = 60
  timer = setInterval(() => {
    countdown.value--
    if (countdown.value <= 0) {
      if (timer) {
        clearInterval(timer)
        timer = null
      }
    }
  }, 1000)
}

const handleSendCode = async () => {
  if (!form.phone || form.phone.length !== 11) {
    uni.showToast({ title: '请输入正确的手机号', icon: 'none' })
    return
  }
  try {
    uni.showLoading({ title: '发送中...' })
    // await sendSmsCode({ phone: form.phone, type: 'forget' })
    uni.hideLoading()
    uni.showToast({ title: '验证码已发送', icon: 'success' })
    startCountdown()
  } catch (e) {
    uni.hideLoading()
    console.error('发送验证码失败', e)
  }
}

const handleVerifyCode = async () => {
  if (!form.phone || form.phone.length !== 11) {
    uni.showToast({ title: '请输入正确的手机号', icon: 'none' })
    return
  }
  if (!form.code) {
    uni.showToast({ title: '请输入验证码', icon: 'none' })
    return
  }

  loading.value = true
  try {
    // await verifySmsCode({ phone: form.phone, code: form.code })
    step.value = 2
  } catch (e) {
    console.error('验证码验证失败', e)
  } finally {
    loading.value = false
  }
}

const handleResetPassword = async () => {
  if (!form.password || form.password.length < 6) {
    uni.showToast({ title: '密码不能少于6位', icon: 'none' })
    return
  }
  if (form.password !== form.confirmPassword) {
    uni.showToast({ title: '两次密码不一致', icon: 'none' })
    return
  }

  loading.value = true
  try {
    // await resetPassword({
    //   phone: form.phone,
    //   code: form.code,
    //   newPassword: form.password
    // })
    step.value = 3
  } catch (e) {
    console.error('重置密码失败', e)
  } finally {
    loading.value = false
  }
}

const goToLogin = () => {
  uni.navigateTo({ url: '/pages/login/index' })
}
</script>

<style lang="scss" scoped>
.forget-container {
  min-height: 100vh;
  background-color: #ffffff;
  padding: 0 60rpx;
}

.steps {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 60rpx 0 40rpx;
}

.step {
  display: flex;
  flex-direction: column;
  align-items: center;

  .step-dot {
    width: 48rpx;
    height: 48rpx;
    border-radius: 50%;
    background-color: #DCDFE6;
    color: #ffffff;
    font-size: 24rpx;
    display: flex;
    align-items: center;
    justify-content: center;
    margin-bottom: 12rpx;
  }

  .step-text {
    font-size: 24rpx;
    color: #909399;
  }

  &.active {
    .step-dot {
      background-color: #4A90D9;
    }

    .step-text {
      color: #4A90D9;
    }
  }

  &.done {
    .step-dot {
      background-color: #67C23A;
    }
  }
}

.step-line {
  width: 100rpx;
  height: 4rpx;
  background-color: #DCDFE6;
  margin: 0 16rpx;
  margin-bottom: 30rpx;

  &.active {
    background-color: #4A90D9;
  }
}

.form-section {
  .section-title {
    margin-bottom: 40rpx;

    .title-text {
      display: block;
      font-size: 40rpx;
      font-weight: bold;
      color: #303133;
      margin-bottom: 12rpx;
    }

    .title-desc {
      font-size: 28rpx;
      color: #909399;
    }
  }

  .form-item {
    margin-bottom: 30rpx;

    &--row {
      display: flex;
      align-items: center;
      gap: 20rpx;
    }

    &--password {
      display: flex;
      align-items: center;
    }
  }

  .form-input {
    height: 100rpx;
    padding: 0 30rpx;
    background-color: #F5F7FA;
    border-radius: 16rpx;
    font-size: 30rpx;

    &--flex {
      flex: 1;
    }
  }

  .sms-btn {
    flex-shrink: 0;
    height: 100rpx;
    padding: 0 24rpx;
    background-color: #4A90D9;
    color: #ffffff;
    font-size: 26rpx;
    border-radius: 16rpx;
    border: none;
    white-space: nowrap;

    &[disabled] {
      background-color: #A0CFFF;
    }
  }

  .toggle-eye {
    flex-shrink: 0;
    width: 80rpx;
    height: 100rpx;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 36rpx;
  }

  .submit-btn {
    width: 100%;
    height: 100rpx;
    line-height: 100rpx;
    background-color: #4A90D9;
    color: #ffffff;
    font-size: 32rpx;
    border-radius: 16rpx;
    margin-top: 20rpx;
    border: none;

    &:active {
      opacity: 0.8;
    }
  }
}

.success-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding-top: 100rpx;

  .success-icon {
    width: 120rpx;
    height: 120rpx;
    border-radius: 50%;
    background-color: #67C23A;
    display: flex;
    align-items: center;
    justify-content: center;
    margin-bottom: 30rpx;

    .icon-text {
      font-size: 60rpx;
      color: #ffffff;
    }
  }

  .success-title {
    font-size: 36rpx;
    font-weight: bold;
    color: #303133;
    margin-bottom: 16rpx;
  }

  .success-desc {
    font-size: 28rpx;
    color: #909399;
    margin-bottom: 60rpx;
  }

  .submit-btn {
    width: 100%;
    height: 100rpx;
    line-height: 100rpx;
    background-color: #4A90D9;
    color: #ffffff;
    font-size: 32rpx;
    border-radius: 16rpx;
    border: none;

    &:active {
      opacity: 0.8;
    }
  }
}

.form-footer {
  display: flex;
  justify-content: center;
  margin-top: 40rpx;

  .link {
    font-size: 28rpx;
    color: #4A90D9;
  }
}
</style>
