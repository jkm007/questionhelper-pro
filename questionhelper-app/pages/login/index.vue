<template>
  <view class="login-container">
    <!-- Logo -->
    <view class="logo-section">
      <image class="logo" src="/static/images/logo.png" mode="aspectFit"></image>
      <text class="app-name">题小助</text>
      <text class="app-slogan">学习好帮手</text>
    </view>

    <!-- 登录表单 -->
    <view class="form-section">
      <view class="form-item">
        <input
          v-model="form.username"
          class="form-input"
          placeholder="请输入用户名/手机号"
          maxlength="50"
        />
      </view>
      <view class="form-item">
        <input
          v-model="form.password"
          class="form-input"
          placeholder="请输入密码"
          :password="true"
          maxlength="20"
        />
      </view>

      <button class="login-btn" @tap="handleLogin" :loading="loading">
        登 录
      </button>

      <view class="form-footer">
        <text class="link" @tap="goToRegister">注册账号</text>
        <text class="link" @tap="goToForget">忘记密码</text>
      </view>
    </view>

    <!-- 其他登录方式 -->
    <view class="other-login">
      <view class="divider">
        <view class="divider-line"></view>
        <text class="divider-text">其他登录方式</text>
        <view class="divider-line"></view>
      </view>
      <view class="social-login">
        <view class="social-icon" @tap="wechatLogin">
          <text>微信</text>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useUserStore } from '@/store/modules/user'

const userStore = useUserStore()
const loading = ref(false)

const form = reactive({
  username: '',
  password: ''
})

const handleLogin = async () => {
  if (!form.username) {
    uni.showToast({ title: '请输入用户名', icon: 'none' })
    return
  }
  if (!form.password) {
    uni.showToast({ title: '请输入密码', icon: 'none' })
    return
  }

  loading.value = true
  try {
    await userStore.login(form.username, form.password)
    uni.showToast({ title: '登录成功', icon: 'success' })
    setTimeout(() => {
      uni.switchTab({ url: '/pages/index/index' })
    }, 500)
  } catch (e) {
    console.error('登录失败', e)
  } finally {
    loading.value = false
  }
}

const goToRegister = () => {
  uni.navigateTo({ url: '/pages/login/register' })
}

const goToForget = () => {
  uni.navigateTo({ url: '/pages/login/forget' })
}

const wechatLogin = () => {
  uni.showToast({ title: '微信登录开发中', icon: 'none' })
}
</script>

<style lang="scss" scoped>
.login-container {
  min-height: 100vh;
  background-color: #ffffff;
  padding: 0 60rpx;
}

.logo-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding-top: 150rpx;
  margin-bottom: 80rpx;

  .logo {
    width: 160rpx;
    height: 160rpx;
    margin-bottom: 24rpx;
  }

  .app-name {
    font-size: 48rpx;
    font-weight: bold;
    color: #4A90D9;
    margin-bottom: 12rpx;
  }

  .app-slogan {
    font-size: 28rpx;
    color: #909399;
  }
}

.form-section {
  .form-item {
    margin-bottom: 30rpx;
  }

  .form-input {
    height: 100rpx;
    padding: 0 30rpx;
    background-color: #F5F7FA;
    border-radius: 16rpx;
    font-size: 30rpx;
  }

  .login-btn {
    width: 100%;
    height: 100rpx;
    line-height: 100rpx;
    background-color: #4A90D9;
    color: #ffffff;
    font-size: 32rpx;
    border-radius: 16rpx;
    margin-top: 40rpx;
    border: none;

    &:active {
      opacity: 0.8;
    }
  }

  .form-footer {
    display: flex;
    justify-content: space-between;
    margin-top: 30rpx;

    .link {
      font-size: 28rpx;
      color: #4A90D9;
    }
  }
}

.other-login {
  margin-top: 100rpx;

  .divider {
    display: flex;
    align-items: center;
    margin-bottom: 40rpx;

    .divider-line {
      flex: 1;
      height: 1rpx;
      background-color: #DCDFE6;
    }

    .divider-text {
      padding: 0 20rpx;
      font-size: 24rpx;
      color: #909399;
    }
  }

  .social-login {
    display: flex;
    justify-content: center;

    .social-icon {
      width: 100rpx;
      height: 100rpx;
      border-radius: 50%;
      background-color: #07C160;
      display: flex;
      align-items: center;
      justify-content: center;
      color: #ffffff;
      font-size: 24rpx;
    }
  }
}
</style>
