<template>
  <view class="register-container">
    <!-- 标题 -->
    <view class="header">
      <text class="title">注册账号</text>
      <text class="subtitle">创建您的题小助账号</text>
    </view>

    <!-- 注册表单 -->
    <view class="form-section">
      <view class="form-item">
        <input
          v-model="form.nickname"
          class="form-input"
          placeholder="请输入昵称"
          maxlength="20"
        />
      </view>
      <view class="form-item form-item--row">
        <input
          v-model="form.phone"
          class="form-input form-input--flex"
          placeholder="请输入手机号"
          type="number"
          maxlength="11"
        />
        <button
          class="sms-btn"
          :disabled="smsDisabled"
          @tap="handleSendCode"
        >
          {{ smsText }}
        </button>
      </view>
      <view class="form-item">
        <input
          v-model="form.code"
          class="form-input"
          placeholder="请输入验证码"
          type="number"
          maxlength="6"
        />
      </view>
      <view class="form-item form-item--password">
        <input
          v-model="form.password"
          class="form-input form-input--flex"
          :password="!showPassword"
          placeholder="请输入密码"
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
          placeholder="请确认密码"
          maxlength="20"
        />
        <view class="toggle-eye" @tap="showConfirmPassword = !showConfirmPassword">
          <text>{{ showConfirmPassword ? '🙈' : '👁️' }}</text>
        </view>
      </view>

      <!-- 协议勾选 -->
      <view class="agreement">
        <view class="checkbox" :class="{ checked: agreed }" @tap="agreed = !agreed">
          <text v-if="agreed" class="check-icon">✓</text>
        </view>
        <text class="agreement-text">
          我已阅读并同意
          <text class="link" @tap.stop="goToUserAgreement">《用户协议》</text>
          和
          <text class="link" @tap.stop="goToPrivacy">《隐私政策》</text>
        </text>
      </view>

      <button class="register-btn" :loading="loading" @tap="handleRegister">
        注 册
      </button>

      <view class="form-footer">
        <text class="link" @tap="goToLogin">已有账号？去登录</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import { register } from '@/api/auth'

const loading = ref(false)
const showPassword = ref(false)
const showConfirmPassword = ref(false)
const agreed = ref(false)
const countdown = ref(0)
let timer: ReturnType<typeof setInterval> | null = null

const form = reactive({
  nickname: '',
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
    // 假设发送验证码的接口
    // await sendSmsCode({ phone: form.phone, type: 'register' })
    uni.hideLoading()
    uni.showToast({ title: '验证码已发送', icon: 'success' })
    startCountdown()
  } catch (e) {
    uni.hideLoading()
    console.error('发送验证码失败', e)
  }
}

const handleRegister = async () => {
  if (!form.nickname) {
    uni.showToast({ title: '请输入昵称', icon: 'none' })
    return
  }
  if (!form.phone || form.phone.length !== 11) {
    uni.showToast({ title: '请输入正确的手机号', icon: 'none' })
    return
  }
  if (!form.code) {
    uni.showToast({ title: '请输入验证码', icon: 'none' })
    return
  }
  if (!form.password || form.password.length < 6) {
    uni.showToast({ title: '密码不能少于6位', icon: 'none' })
    return
  }
  if (form.password !== form.confirmPassword) {
    uni.showToast({ title: '两次密码不一致', icon: 'none' })
    return
  }
  if (!agreed.value) {
    uni.showToast({ title: '请同意用户协议', icon: 'none' })
    return
  }

  loading.value = true
  try {
    await register({
      username: form.nickname,
      password: form.password,
      phone: form.phone
    })
    uni.showToast({ title: '注册成功', icon: 'success' })
    setTimeout(() => {
      uni.navigateTo({ url: '/pages/login/index' })
    }, 500)
  } catch (e) {
    console.error('注册失败', e)
  } finally {
    loading.value = false
  }
}

const goToLogin = () => {
  uni.navigateTo({ url: '/pages/login/index' })
}

const goToUserAgreement = () => {
  uni.navigateTo({ url: '/pages/webview/index?url=https://example.com/agreement' })
}

const goToPrivacy = () => {
  uni.navigateTo({ url: '/pages/webview/index?url=https://example.com/privacy' })
}
</script>

<style lang="scss" scoped>
.register-container {
  min-height: 100vh;
  background-color: #ffffff;
  padding: 0 60rpx;
}

.header {
  padding-top: 80rpx;
  margin-bottom: 60rpx;

  .title {
    display: block;
    font-size: 48rpx;
    font-weight: bold;
    color: #303133;
    margin-bottom: 16rpx;
  }

  .subtitle {
    font-size: 28rpx;
    color: #909399;
  }
}

.form-section {
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

  .agreement {
    display: flex;
    align-items: flex-start;
    margin-bottom: 40rpx;
    margin-top: 10rpx;
  }

  .checkbox {
    flex-shrink: 0;
    width: 36rpx;
    height: 36rpx;
    border: 2rpx solid #DCDFE6;
    border-radius: 8rpx;
    margin-right: 12rpx;
    margin-top: 4rpx;
    display: flex;
    align-items: center;
    justify-content: center;

    &.checked {
      background-color: #4A90D9;
      border-color: #4A90D9;
    }

    .check-icon {
      font-size: 24rpx;
      color: #ffffff;
    }
  }

  .agreement-text {
    font-size: 24rpx;
    color: #909399;
    line-height: 1.5;

    .link {
      color: #4A90D9;
    }
  }

  .register-btn {
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

  .form-footer {
    display: flex;
    justify-content: center;
    margin-top: 30rpx;

    .link {
      font-size: 28rpx;
      color: #4A90D9;
    }
  }
}
</style>
