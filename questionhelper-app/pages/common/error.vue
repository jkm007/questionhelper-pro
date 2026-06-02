<template>
  <view class="page">
    <view class="error-container">
      <!-- Error Icon -->
      <view class="error-icon-wrap">
        <view class="error-icon" :class="'icon-' + errorType">
          <text class="icon-symbol">{{ errorInfo.icon }}</text>
        </view>
      </view>

      <!-- Error Code -->
      <text class="error-code">{{ errorInfo.code }}</text>

      <!-- Error Message -->
      <text class="error-message">{{ displayMessage }}</text>

      <!-- Error Description -->
      <text class="error-desc">{{ errorInfo.description }}</text>

      <!-- Actions -->
      <view class="action-group">
        <view class="action-btn primary-btn" @tap="onRetry">
          <text class="action-btn-text primary-text">{{ errorInfo.retryText }}</text>
        </view>
        <view class="action-btn secondary-btn" @tap="onBackHome">
          <text class="action-btn-text secondary-text">返回首页</text>
        </view>
      </view>

      <!-- Help Text -->
      <view v-if="errorType === 'network'" class="help-section">
        <text class="help-title">网络问题排查建议</text>
        <view class="help-list">
          <text class="help-item">1. 检查手机是否已连接网络</text>
          <text class="help-item">2. 尝试切换 Wi-Fi 或移动数据</text>
          <text class="help-item">3. 检查是否开启了 VPN 或代理</text>
          <text class="help-item">4. 稍后再试，可能是服务器维护中</text>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'

interface ErrorInfo {
  icon: string
  code: string
  title: string
  description: string
  retryText: string
}

const props = withDefaults(
  defineProps<{
    errorCode?: string | number
    errorMessage?: string
  }>(),
  {
    errorCode: '',
    errorMessage: ''
  }
)

const pageErrorCode = ref('')
const pageErrorMessage = ref('')

const errorType = computed(() => {
  const code = String(pageErrorCode.value || props.errorCode)
  if (code === '401' || code === '403') return 'auth'
  if (code === '404') return 'notfound'
  if (code === 'network' || code === '0' || code === 'timeout') return 'network'
  return 'unknown'
})

const displayMessage = computed(() => {
  const msg = pageErrorMessage.value || props.errorMessage
  return msg || errorInfo.value.title
})

const errorInfoMap: Record<string, ErrorInfo> = {
  auth: {
    icon: '!',
    code: '401',
    title: '未授权访问',
    description: '您的登录状态已过期，请重新登录后再试',
    retryText: '重新登录'
  },
  notfound: {
    icon: '?',
    code: '404',
    title: '页面不存在',
    description: '您访问的页面可能已被删除或地址有误',
    retryText: '刷新页面'
  },
  network: {
    icon: 'x',
    code: '--',
    title: '网络连接失败',
    description: '无法连接到服务器，请检查您的网络设置',
    retryText: '重新加载'
  },
  unknown: {
    icon: '!',
    code: '500',
    title: '服务器错误',
    description: '服务器开小差了，请稍后再试',
    retryText: '重试'
  }
}

const errorInfo = computed<ErrorInfo>(() => {
  return errorInfoMap[errorType.value] || errorInfoMap.unknown
})

function onRetry() {
  const type = errorType.value
  if (type === 'auth') {
    uni.reLaunch({ url: '/pages/login/index' })
    return
  }
  if (type === 'notfound') {
    uni.navigateBack({
      fail: () => {
        uni.reLaunch({ url: '/pages/index/index' })
      }
    })
    return
  }
  // For network/unknown errors, reload the previous page
  const pages = getCurrentPages()
  if (pages.length > 1) {
    uni.navigateBack()
  } else {
    uni.reLaunch({ url: '/pages/index/index' })
  }
}

function onBackHome() {
  uni.reLaunch({ url: '/pages/index/index' })
}

onMounted(() => {
  // Try to read error code/message from page options (for navigation params)
  const pages = getCurrentPages()
  const currentPage = pages[pages.length - 1] as any
  if (currentPage?.options) {
    pageErrorCode.value = currentPage.options.errorCode || ''
    pageErrorMessage.value = currentPage.options.errorMessage || ''
  }
})
</script>

<style scoped>
.page {
  min-height: 100vh;
  background-color: #f5f6fa;
  display: flex;
  align-items: center;
  justify-content: center;
}

.error-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 60rpx 48rpx;
  width: 100%;
  box-sizing: border-box;
}

.error-icon-wrap {
  margin-bottom: 40rpx;
}

.error-icon {
  width: 160rpx;
  height: 160rpx;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.icon-auth {
  background-color: #fff7e6;
  border: 4rpx solid #ffd591;
}

.icon-notfound {
  background-color: #f0f5ff;
  border: 4rpx solid #adc6ff;
}

.icon-network {
  background-color: #fff1f0;
  border: 4rpx solid #ffa39e;
}

.icon-unknown {
  background-color: #f6ffed;
  border: 4rpx solid #b7eb8f;
}

.icon-symbol {
  font-size: 72rpx;
  font-weight: 700;
}

.icon-auth .icon-symbol {
  color: #fa8c16;
}

.icon-notfound .icon-symbol {
  color: #1677ff;
}

.icon-network .icon-symbol {
  color: #ff4d4f;
}

.icon-unknown .icon-symbol {
  color: #52c41a;
}

.error-code {
  font-size: 80rpx;
  font-weight: 800;
  color: #d9d9d9;
  margin-bottom: 16rpx;
  letter-spacing: 8rpx;
}

.error-message {
  font-size: 34rpx;
  font-weight: 600;
  color: #333;
  margin-bottom: 16rpx;
}

.error-desc {
  font-size: 26rpx;
  color: #999;
  text-align: center;
  line-height: 1.6;
  margin-bottom: 60rpx;
  padding: 0 20rpx;
}

.action-group {
  width: 100%;
  display: flex;
  flex-direction: column;
  gap: 24rpx;
  margin-bottom: 40rpx;
}

.action-btn {
  width: 100%;
  padding: 24rpx 0;
  border-radius: 44rpx;
  text-align: center;
}

.primary-btn {
  background-color: #1677ff;
}

.secondary-btn {
  background-color: #fff;
  border: 2rpx solid #e0e0e0;
}

.action-btn-text {
  font-size: 30rpx;
  font-weight: 600;
}

.primary-text {
  color: #fff;
}

.secondary-text {
  color: #666;
}

.help-section {
  width: 100%;
  background-color: #fff;
  border-radius: 16rpx;
  padding: 28rpx;
  margin-top: 20rpx;
}

.help-title {
  font-size: 28rpx;
  font-weight: 600;
  color: #333;
  margin-bottom: 20rpx;
  display: block;
}

.help-list {
  display: flex;
  flex-direction: column;
  gap: 12rpx;
}

.help-item {
  font-size: 26rpx;
  color: #666;
  line-height: 1.6;
}
</style>
