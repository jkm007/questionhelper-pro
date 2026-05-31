<template>
  <view class="detail-page">
    <!-- Loading state -->
    <view class="loading-state" v-if="loading">
      <text class="loading-text">Loading...</text>
    </view>

    <!-- Content -->
    <view class="detail-content" v-else-if="notification">
      <!-- Type icon and header -->
      <view class="detail-header">
        <view class="type-icon-wrapper" :class="'type-' + notification.type">
          <text class="type-icon-text">{{ getTypeIcon(notification.type) }}</text>
        </view>
        <view class="header-info">
          <text class="detail-title">{{ notification.title }}</text>
          <text class="detail-time">{{ formatTime(notification.createdAt) }}</text>
        </view>
      </view>

      <!-- Divider -->
      <view class="divider"></view>

      <!-- Rich text content -->
      <view class="detail-body">
        <rich-text :nodes="notification.content" class="rich-content"></rich-text>
      </view>

      <!-- Related resource link -->
      <view
        class="related-link"
        v-if="notification.relatedUrl"
        @tap="openRelatedLink"
      >
        <view class="link-icon-wrapper">
          <text class="link-icon">🔗</text>
        </view>
        <view class="link-info">
          <text class="link-label">Related Resource</text>
          <text class="link-url">{{ notification.relatedUrl }}</text>
        </view>
        <text class="link-arrow">></text>
      </view>

      <!-- Notification metadata -->
      <view class="meta-section">
        <view class="meta-item">
          <text class="meta-label">Type</text>
          <text class="meta-value">{{ getTypeName(notification.type) }}</text>
        </view>
        <view class="meta-item">
          <text class="meta-label">ID</text>
          <text class="meta-value">{{ notification.id }}</text>
        </view>
      </view>
    </view>

    <!-- Error state -->
    <view class="error-state" v-else-if="error">
      <text class="error-icon">⚠️</text>
      <text class="error-text">{{ error }}</text>
      <view class="retry-btn" @tap="fetchDetail">
        <text class="retry-text">Retry</text>
      </view>
    </view>

    <!-- Bottom action bar -->
    <view class="bottom-bar" v-if="notification">
      <view class="delete-action-btn" @tap="handleDelete">
        <text class="delete-action-text">Delete This Notification</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import {
  getNotificationDetail,
  deleteNotification,
  type NotificationItem
} from '@/api/notification'

const notification = ref<NotificationItem | null>(null)
const loading = ref(true)
const error = ref<string | null>(null)
const notificationId = ref('')

const getTypeIcon = (type: string): string => {
  const iconMap: Record<string, string> = {
    system: '🔔',
    exam: '📝',
    feedback: '💬',
    update: '🔄',
    promotion: '🎁'
  }
  return iconMap[type] || '📢'
}

const getTypeName = (type: string): string => {
  const nameMap: Record<string, string> = {
    system: 'System Notification',
    exam: 'Exam Notification',
    feedback: 'Feedback Reply',
    update: 'App Update',
    promotion: 'Promotion'
  }
  return nameMap[type] || 'Notification'
}

const formatTime = (time: string): string => {
  const date = new Date(time)
  const year = date.getFullYear()
  const month = (date.getMonth() + 1).toString().padStart(2, '0')
  const day = date.getDate().toString().padStart(2, '0')
  const hours = date.getHours().toString().padStart(2, '0')
  const minutes = date.getMinutes().toString().padStart(2, '0')
  return `${year}-${month}-${day} ${hours}:${minutes}`
}

const fetchDetail = async () => {
  loading.value = true
  error.value = null

  try {
    const res = await getNotificationDetail(notificationId.value)
    notification.value = res.data || null
  } catch (err: any) {
    error.value = err.message || 'Failed to load notification details'
  } finally {
    loading.value = false
  }
}

const openRelatedLink = () => {
  if (!notification.value?.relatedUrl) return

  const url = notification.value.relatedUrl

  // Check if it's an internal page
  if (url.startsWith('/pages/')) {
    uni.navigateTo({ url })
  } else {
    // Open in webview
    uni.navigateTo({
      url: `/pages/webview/index?url=${encodeURIComponent(url)}`
    })
  }
}

const handleDelete = () => {
  uni.showModal({
    title: 'Confirm Delete',
    content: 'Are you sure you want to delete this notification? This action cannot be undone.',
    confirmColor: '#ff4757',
    success: async (res) => {
      if (res.confirm) {
        try {
          await deleteNotification(notificationId.value)
          uni.showToast({
            title: 'Deleted',
            icon: 'success'
          })
          setTimeout(() => {
            uni.navigateBack()
          }, 1500)
        } catch (err) {
          uni.showToast({
            title: 'Failed to delete',
            icon: 'none'
          })
        }
      }
    }
  })
}

onLoad((options) => {
  if (options?.id) {
    notificationId.value = options.id
    fetchDetail()
  } else {
    loading.value = false
    error.value = 'Notification ID is missing'
  }
})
</script>

<style lang="scss" scoped>
.detail-page {
  min-height: 100vh;
  background-color: #f5f6fa;
  padding-bottom: 120rpx;
}

.loading-state {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 120rpx 0;
}

.loading-text {
  font-size: 28rpx;
  color: #999999;
}

.detail-content {
  background-color: #ffffff;
  margin: 20rpx;
  border-radius: 16rpx;
  overflow: hidden;
}

.detail-header {
  display: flex;
  align-items: center;
  padding: 30rpx;
}

.type-icon-wrapper {
  width: 90rpx;
  height: 90rpx;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 24rpx;
  background-color: #e8f0fe;
}

.type-system {
  background-color: #fff3e0;
}

.type-exam {
  background-color: #e8f5e9;
}

.type-feedback {
  background-color: #e3f2fd;
}

.type-update {
  background-color: #f3e5f5;
}

.type-promotion {
  background-color: #fce4ec;
}

.type-icon-text {
  font-size: 44rpx;
}

.header-info {
  flex: 1;
}

.detail-title {
  font-size: 32rpx;
  font-weight: 600;
  color: #333333;
  margin-bottom: 8rpx;
}

.detail-time {
  font-size: 24rpx;
  color: #999999;
}

.divider {
  height: 1rpx;
  background-color: #f0f0f0;
  margin: 0 30rpx;
}

.detail-body {
  padding: 30rpx;
}

.rich-content {
  font-size: 28rpx;
  color: #333333;
  line-height: 1.8;
}

.related-link {
  display: flex;
  align-items: center;
  padding: 24rpx 30rpx;
  background-color: #f8f9ff;
  margin: 0 30rpx 30rpx;
  border-radius: 12rpx;
  border: 1rpx solid #e8ecf4;
}

.link-icon-wrapper {
  width: 60rpx;
  height: 60rpx;
  border-radius: 12rpx;
  background-color: #e8f0fe;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 20rpx;
}

.link-icon {
  font-size: 28rpx;
}

.link-info {
  flex: 1;
  overflow: hidden;
}

.link-label {
  font-size: 26rpx;
  color: #666666;
  margin-bottom: 4rpx;
}

.link-url {
  font-size: 24rpx;
  color: #007aff;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.link-arrow {
  font-size: 28rpx;
  color: #cccccc;
  margin-left: 16rpx;
}

.meta-section {
  padding: 20rpx 30rpx 30rpx;
}

.meta-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12rpx 0;
}

.meta-label {
  font-size: 24rpx;
  color: #999999;
}

.meta-value {
  font-size: 24rpx;
  color: #666666;
}

.error-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 120rpx 0;
}

.error-icon {
  font-size: 80rpx;
  margin-bottom: 20rpx;
}

.error-text {
  font-size: 28rpx;
  color: #666666;
  margin-bottom: 30rpx;
}

.retry-btn {
  padding: 16rpx 40rpx;
  background-color: #007aff;
  border-radius: 8rpx;
}

.retry-text {
  font-size: 28rpx;
  color: #ffffff;
}

.bottom-bar {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  background-color: #ffffff;
  padding: 20rpx 30rpx;
  padding-bottom: calc(20rpx + env(safe-area-inset-bottom));
  box-shadow: 0 -2rpx 10rpx rgba(0, 0, 0, 0.05);
}

.delete-action-btn {
  width: 100%;
  height: 88rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #fff0f0;
  border-radius: 12rpx;
  border: 1rpx solid #ffd0d0;
}

.delete-action-text {
  font-size: 28rpx;
  color: #ff4757;
  font-weight: 500;
}
</style>
