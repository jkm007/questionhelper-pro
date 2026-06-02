<template>
  <view class="page">
    <!-- App Logo & Info -->
    <view class="logo-section">
      <image class="app-logo" src="/static/logo.png" mode="aspectFit" />
      <text class="app-name">题小助</text>
      <text class="app-slogan">智能题库，轻松备考</text>
      <view class="version-badge">
        <text class="version-text">v{{ appVersion }}</text>
      </view>
    </view>

    <!-- Feature Highlights -->
    <view class="section">
      <text class="section-title">核心功能</text>
      <view class="feature-grid">
        <view class="feature-item">
          <view class="feature-icon feature-icon-blue">
            <text class="feature-icon-text">题</text>
          </view>
          <text class="feature-name">海量题库</text>
          <text class="feature-desc">多学科全覆盖</text>
        </view>
        <view class="feature-item">
          <view class="feature-icon feature-icon-green">
            <text class="feature-icon-text">练</text>
          </view>
          <text class="feature-name">智能练习</text>
          <text class="feature-desc">个性化推荐</text>
        </view>
        <view class="feature-item">
          <view class="feature-icon feature-icon-orange">
            <text class="feature-icon-text">错</text>
          </view>
          <text class="feature-name">错题巩固</text>
          <text class="feature-desc">查漏补缺</text>
        </view>
        <view class="feature-item">
          <view class="feature-icon feature-icon-purple">
            <text class="feature-icon-text">析</text>
          </view>
          <text class="feature-name">数据分析</text>
          <text class="feature-desc">学习报告</text>
        </view>
      </view>
    </view>

    <!-- Links -->
    <view class="section links-section">
      <view class="link-item" @tap="openUserAgreement">
        <text class="link-label">用户协议</text>
        <text class="link-arrow">&#xe61a;</text>
      </view>
      <view class="link-item" @tap="openPrivacyPolicy">
        <text class="link-label">隐私政策</text>
        <text class="link-arrow">&#xe61a;</text>
      </view>
      <view class="link-item" @tap="openThirdParty">
        <text class="link-label">第三方信息共享清单</text>
        <text class="link-arrow">&#xe61a;</text>
      </view>
      <view class="link-item" @tap="openLicenses">
        <text class="link-label">开源许可</text>
        <text class="link-arrow">&#xe61a;</text>
      </view>
    </view>

    <!-- Check Update -->
    <view class="section links-section">
      <view class="link-item" @tap="checkUpdate">
        <text class="link-label">检查更新</text>
        <view class="link-right">
          <text class="link-value">{{ updateStatus }}</text>
          <text class="link-arrow">&#xe61a;</text>
        </view>
      </view>
      <view class="link-item" @tap="clearCache">
        <text class="link-label">清除缓存</text>
        <view class="link-right">
          <text class="link-value">{{ cacheSize }}</text>
          <text class="link-arrow">&#xe61a;</text>
        </view>
      </view>
    </view>

    <!-- Feedback -->
    <view class="section links-section">
      <view class="link-item" @tap="goFeedback">
        <text class="link-label">意见反馈</text>
        <text class="link-arrow">&#xe61a;</text>
      </view>
      <view class="link-item" @tap="goCustomerService">
        <text class="link-label">联系客服</text>
        <text class="link-arrow">&#xe61a;</text>
      </view>
    </view>

    <!-- Copyright -->
    <view class="copyright-section">
      <text class="copyright-text">Copyright &copy; 2024 题小助</text>
      <text class="copyright-text">All Rights Reserved</text>
      <text class="copyright-sub">本应用由题小助团队开发并维护</text>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

const appVersion = ref('1.0.0')
const updateStatus = ref('已是最新版本')
const cacheSize = ref('0 KB')
const checking = ref(false)

function getAppVersion() {
  try {
    // #ifdef MP-WEIXIN
    const accountInfo = uni.getAccountInfoSync()
    appVersion.value = accountInfo.miniProgram?.version || '1.0.0'
    // #endif
    // #ifndef MP-WEIXIN
    appVersion.value = '1.0.0'
    // #endif
  } catch (e) {
    appVersion.value = '1.0.0'
  }
}

function calculateCacheSize() {
  try {
    const res = uni.getStorageInfoSync()
    const sizeKB = res.currentSize || 0
    if (sizeKB > 1024) {
      cacheSize.value = (sizeKB / 1024).toFixed(1) + ' MB'
    } else {
      cacheSize.value = sizeKB + ' KB'
    }
  } catch (e) {
    cacheSize.value = '未知'
  }
}

async function checkUpdate() {
  if (checking.value) return
  checking.value = true
  updateStatus.value = '检查中...'

  try {
    // #ifdef MP-WEIXIN
    const updateManager = uni.getUpdateManager()
    updateManager.onCheckForUpdate((res) => {
      if (res.hasUpdate) {
        updateStatus.value = '发现新版本'
        updateManager.onUpdateReady(() => {
          uni.showModal({
            title: '更新提示',
            content: '新版本已准备好，是否重启应用？',
            success: (modalRes) => {
              if (modalRes.confirm) {
                updateManager.applyUpdate()
              }
            }
          })
        })
        updateManager.onUpdateFailed(() => {
          uni.showToast({ title: '更新失败，请稍后重试', icon: 'none' })
          updateStatus.value = '更新失败'
        })
      } else {
        updateStatus.value = '已是最新版本'
        uni.showToast({ title: '已是最新版本', icon: 'success' })
      }
    })
    // #endif
    // #ifndef MP-WEIXIN
    // H5 / App - simulate check
    await new Promise((r) => setTimeout(r, 1500))
    updateStatus.value = '已是最新版本'
    uni.showToast({ title: '已是最新版本', icon: 'success' })
    // #endif
  } catch (e) {
    updateStatus.value = '检查失败'
    uni.showToast({ title: '检查更新失败', icon: 'none' })
  } finally {
    checking.value = false
  }
}

function clearCache() {
  uni.showModal({
    title: '清除缓存',
    content: `确定清除本地缓存 (${cacheSize.value}) 吗？清除后不会影响账号数据。`,
    success: (res) => {
      if (res.confirm) {
        // Preserve auth token before clearing
        const token = uni.getStorageSync('token')
        try {
          uni.clearStorageSync()
          if (token) {
            uni.setStorageSync('token', token)
          }
          cacheSize.value = '0 KB'
          uni.showToast({ title: '缓存已清除', icon: 'success' })
        } catch (e) {
          uni.showToast({ title: '清除失败', icon: 'none' })
        }
      }
    }
  })
}

function openUserAgreement() {
  uni.navigateTo({ url: '/pages/webview/index?url=/user-agreement' })
}

function openPrivacyPolicy() {
  uni.navigateTo({ url: '/pages/webview/index?url=/privacy-policy' })
}

function openThirdParty() {
  uni.navigateTo({ url: '/pages/webview/index?url=/third-party-sharing' })
}

function openLicenses() {
  uni.navigateTo({ url: '/pages/webview/index?url=/open-source-licenses' })
}

function goFeedback() {
  uni.navigateTo({ url: '/pages/feedback/correction' })
}

function goCustomerService() {
  uni.showModal({
    title: '联系客服',
    content: '客服邮箱: support@questionhelper.com\n工作时间: 周一至周五 9:00-18:00',
    showCancel: false,
    confirmText: '知道了'
  })
}

onMounted(() => {
  getAppVersion()
  calculateCacheSize()
})
</script>

<style scoped>
.page {
  min-height: 100vh;
  background-color: #f5f6fa;
  padding-bottom: 60rpx;
}

/* Logo Section */
.logo-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 60rpx 32rpx 48rpx;
  background-color: #fff;
  margin-bottom: 20rpx;
}

.app-logo {
  width: 140rpx;
  height: 140rpx;
  border-radius: 28rpx;
  margin-bottom: 24rpx;
}

.app-name {
  font-size: 38rpx;
  font-weight: 700;
  color: #333;
  margin-bottom: 8rpx;
}

.app-slogan {
  font-size: 26rpx;
  color: #999;
  margin-bottom: 20rpx;
}

.version-badge {
  padding: 6rpx 24rpx;
  background-color: #e8f3ff;
  border-radius: 20rpx;
}

.version-text {
  font-size: 24rpx;
  color: #1677ff;
}

/* Section */
.section {
  background-color: #fff;
  margin: 0 24rpx 20rpx;
  border-radius: 16rpx;
  padding: 28rpx;
}

.section-title {
  font-size: 30rpx;
  font-weight: 600;
  color: #333;
  margin-bottom: 24rpx;
  display: block;
}

/* Feature Grid */
.feature-grid {
  display: flex;
  justify-content: space-between;
}

.feature-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  flex: 1;
}

.feature-icon {
  width: 88rpx;
  height: 88rpx;
  border-radius: 24rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 16rpx;
}

.feature-icon-blue {
  background-color: #e8f3ff;
}

.feature-icon-green {
  background-color: #f6ffed;
}

.feature-icon-orange {
  background-color: #fff7e6;
}

.feature-icon-purple {
  background-color: #f9f0ff;
}

.feature-icon-text {
  font-size: 36rpx;
  font-weight: 700;
}

.feature-icon-blue .feature-icon-text {
  color: #1677ff;
}

.feature-icon-green .feature-icon-text {
  color: #52c41a;
}

.feature-icon-orange .feature-icon-text {
  color: #fa8c16;
}

.feature-icon-purple .feature-icon-text {
  color: #722ed1;
}

.feature-name {
  font-size: 26rpx;
  font-weight: 600;
  color: #333;
  margin-bottom: 4rpx;
}

.feature-desc {
  font-size: 22rpx;
  color: #999;
}

/* Links */
.links-section {
  padding: 0;
}

.link-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 28rpx 28rpx;
  border-bottom: 1rpx solid #f5f5f5;
}

.link-item:last-child {
  border-bottom: none;
}

.link-label {
  font-size: 30rpx;
  color: #333;
}

.link-right {
  display: flex;
  align-items: center;
  gap: 12rpx;
}

.link-value {
  font-size: 26rpx;
  color: #999;
}

.link-arrow {
  font-size: 28rpx;
  color: #ccc;
}

/* Copyright */
.copyright-section {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 40rpx 32rpx 20rpx;
  gap: 8rpx;
}

.copyright-text {
  font-size: 24rpx;
  color: #c0c4cc;
  text-align: center;
}

.copyright-sub {
  font-size: 22rpx;
  color: #d9d9d9;
  text-align: center;
  margin-top: 8rpx;
}
</style>
