<template>
  <view class="webview-page">
    <!-- Loading indicator -->
    <view class="loading-bar" v-if="isLoading">
      <view class="loading-progress"></view>
    </view>

    <!-- Error state -->
    <view class="error-state" v-if="hasError">
      <view class="error-content">
        <text class="error-icon">🌐</text>
        <text class="error-title">Unable to Load Page</text>
        <text class="error-message">{{ errorMessage }}</text>
        <view class="error-actions">
          <view class="retry-btn" @tap="handleRetry">
            <text class="retry-text">Retry</text>
          </view>
          <view class="back-btn" @tap="goBack">
            <text class="back-text">Go Back</text>
          </view>
        </view>
      </view>
    </view>

    <!-- WebView -->
    <web-view
      v-if="url && !hasError"
      :src="url"
      @message="onMessage"
      @load="onLoad"
      @error="onError"
    ></web-view>
  </view>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { onLoad as onPageLoad } from '@dcloudio/uni-app'

const url = ref('')
const pageTitle = ref('')
const isLoading = ref(true)
const hasError = ref(false)
const errorMessage = ref('Failed to load the page. Please check your network connection and try again.')

const validateUrl = (urlStr: string): boolean => {
  try {
    const parsed = new URL(urlStr)
    return ['http:', 'https:'].includes(parsed.protocol)
  } catch {
    return false
  }
}

const onMessage = (e: any) => {
  // Handle messages from webview
  const data = e.detail?.data?.[0]
  if (data) {
    if (data.type === 'title' && data.title) {
      pageTitle.value = data.title
      uni.setNavigationBarTitle({ title: data.title })
    }
    if (data.type === 'close') {
      goBack()
    }
  }
}

const onLoad = () => {
  isLoading.value = false
  hasError.value = false
}

const onError = (e: any) => {
  isLoading.value = false
  hasError.value = true

  const detail = e.detail || {}
  if (detail.url) {
    errorMessage.value = `Failed to load: ${detail.url}`
  } else {
    errorMessage.value = 'Page load failed. Please check your network connection and try again.'
  }
}

const handleRetry = () => {
  hasError.value = false
  isLoading.value = true

  // Force reload by temporarily clearing url
  const currentUrl = url.value
  url.value = ''

  setTimeout(() => {
    url.value = currentUrl
  }, 100)
}

const goBack = () => {
  const pages = getCurrentPages()
  if (pages.length > 1) {
    uni.navigateBack()
  } else {
    uni.switchTab({
      url: '/pages/index/index'
    })
  }
}

onPageLoad((options) => {
  if (options?.url) {
    const decodedUrl = decodeURIComponent(options.url)

    if (validateUrl(decodedUrl)) {
      url.value = decodedUrl

      if (options.title) {
        pageTitle.value = decodeURIComponent(options.title)
        uni.setNavigationBarTitle({ title: pageTitle.value })
      }
    } else {
      hasError.value = true
      errorMessage.value = 'Invalid URL provided.'
      isLoading.value = false
    }
  } else {
    hasError.value = true
    errorMessage.value = 'No URL provided.'
    isLoading.value = false
  }
})
</script>

<style lang="scss" scoped>
.webview-page {
  min-height: 100vh;
  background-color: #ffffff;
  position: relative;
}

.loading-bar {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  height: 4rpx;
  z-index: 999;
  background-color: #e8ecf4;
  overflow: hidden;
}

.loading-progress {
  width: 30%;
  height: 100%;
  background-color: #007aff;
  animation: loading-animation 1.5s ease-in-out infinite;
}

@keyframes loading-animation {
  0% {
    transform: translateX(-100%);
  }
  50% {
    transform: translateX(200%);
  }
  100% {
    transform: translateX(600%);
  }
}

.error-state {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #f5f6fa;
  z-index: 1000;
}

.error-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 60rpx;
  text-align: center;
}

.error-icon {
  font-size: 120rpx;
  margin-bottom: 30rpx;
}

.error-title {
  font-size: 36rpx;
  font-weight: 600;
  color: #333333;
  margin-bottom: 16rpx;
}

.error-message {
  font-size: 28rpx;
  color: #666666;
  line-height: 1.6;
  margin-bottom: 40rpx;
  max-width: 500rpx;
}

.error-actions {
  display: flex;
  gap: 24rpx;
}

.retry-btn {
  padding: 20rpx 48rpx;
  background-color: #007aff;
  border-radius: 12rpx;
}

.retry-text {
  font-size: 28rpx;
  color: #ffffff;
  font-weight: 500;
}

.back-btn {
  padding: 20rpx 48rpx;
  background-color: #ffffff;
  border-radius: 12rpx;
  border: 2rpx solid #e8ecf4;
}

.back-text {
  font-size: 28rpx;
  color: #666666;
  font-weight: 500;
}
</style>
