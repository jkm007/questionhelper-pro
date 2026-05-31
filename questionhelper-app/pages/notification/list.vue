<template>
  <view class="notification-page">
    <!-- Header actions -->
    <view class="header-actions">
      <view class="unread-badge" v-if="unreadCount > 0">
        <text class="unread-text">{{ unreadCount }} unread</text>
      </view>
      <view class="mark-all-btn" @tap="handleMarkAllRead">
        <text class="mark-all-text">Mark all as read</text>
      </view>
    </view>

    <!-- Notification list -->
    <scroll-view
      class="notification-list"
      scroll-y
      :refresher-enabled="true"
      :refresher-triggered="isRefreshing"
      @refresherrefresh="onPullDownRefresh"
      @scrolltolower="onLoadMore"
    >
      <view
        class="notification-item"
        v-for="item in notificationList"
        :key="item.id"
        @tap="goToDetail(item)"
      >
        <view class="item-left">
          <view class="icon-wrapper" :class="'icon-' + item.type">
            <text class="icon-text">{{ getTypeIcon(item.type) }}</text>
          </view>
          <view class="unread-dot" v-if="!item.isRead"></view>
        </view>

        <view class="item-content">
          <view class="item-header">
            <text class="item-title" :class="{ 'title-read': item.isRead }">{{ item.title }}</text>
            <text class="item-time">{{ formatTime(item.createdAt) }}</text>
          </view>
          <text class="item-preview" :class="{ 'preview-read': item.isRead }">{{ item.content }}</text>
        </view>

        <view class="delete-btn" @tap.stop="handleDelete(item)">
          <text class="delete-text">Delete</text>
        </view>
      </view>

      <!-- Empty state -->
      <view class="empty-state" v-if="notificationList.length === 0 && !loading">
        <text class="empty-icon">📭</text>
        <text class="empty-text">No notifications yet</text>
      </view>

      <!-- Load more -->
      <view class="load-more" v-if="notificationList.length > 0">
        <text class="load-more-text" v-if="loading">Loading...</text>
        <text class="load-more-text" v-else-if="noMore">No more notifications</text>
      </view>
    </scroll-view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { onPullDownRefresh as onPullDownRefreshHook, onReachBottom } from '@dcloudio/uni-app'
import {
  getNotificationList,
  markAsRead,
  markAllAsRead,
  deleteNotification,
  type NotificationItem
} from '@/api/notification'

const notificationList = ref<NotificationItem[]>([])
const loading = ref(false)
const isRefreshing = ref(false)
const noMore = ref(false)
const page = ref(1)
const pageSize = 20

const unreadCount = computed(() => {
  return notificationList.value.filter((item) => !item.isRead).length
})

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

const formatTime = (time: string): string => {
  const date = new Date(time)
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  const minutes = Math.floor(diff / 60000)
  const hours = Math.floor(diff / 3600000)
  const days = Math.floor(diff / 86400000)

  if (minutes < 1) return 'Just now'
  if (minutes < 60) return `${minutes}m ago`
  if (hours < 24) return `${hours}h ago`
  if (days < 7) return `${days}d ago`

  const month = (date.getMonth() + 1).toString().padStart(2, '0')
  const day = date.getDate().toString().padStart(2, '0')
  return `${month}-${day}`
}

const fetchNotifications = async (isRefresh = false) => {
  if (loading.value) return

  loading.value = true
  try {
    if (isRefresh) {
      page.value = 1
      noMore.value = false
    }

    const res = await getNotificationList({
      page: page.value,
      pageSize
    })

    const newList = res.data?.list || []
    if (isRefresh) {
      notificationList.value = newList
    } else {
      notificationList.value = [...notificationList.value, ...newList]
    }

    if (newList.length < pageSize) {
      noMore.value = true
    } else {
      page.value++
    }
  } catch (error) {
    uni.showToast({
      title: 'Failed to load notifications',
      icon: 'none'
    })
  } finally {
    loading.value = false
    isRefreshing.value = false
  }
}

const onPullDownRefresh = async () => {
  isRefreshing.value = true
  await fetchNotifications(true)
}

const onLoadMore = () => {
  if (!noMore.value && !loading.value) {
    fetchNotifications()
  }
}

const goToDetail = async (item: NotificationItem) => {
  if (!item.isRead) {
    try {
      await markAsRead(item.id)
      item.isRead = true
    } catch (error) {
      // Continue navigation even if mark as read fails
    }
  }
  uni.navigateTo({
    url: `/pages/notification/detail?id=${item.id}`
  })
}

const handleDelete = (item: NotificationItem) => {
  uni.showModal({
    title: 'Confirm Delete',
    content: 'Are you sure you want to delete this notification?',
    success: async (res) => {
      if (res.confirm) {
        try {
          await deleteNotification(item.id)
          notificationList.value = notificationList.value.filter((n) => n.id !== item.id)
          uni.showToast({
            title: 'Deleted',
            icon: 'success'
          })
        } catch (error) {
          uni.showToast({
            title: 'Failed to delete',
            icon: 'none'
          })
        }
      }
    }
  })
}

const handleMarkAllRead = () => {
  if (unreadCount.value === 0) return

  uni.showModal({
    title: 'Mark All as Read',
    content: 'Mark all notifications as read?',
    success: async (res) => {
      if (res.confirm) {
        try {
          await markAllAsRead()
          notificationList.value.forEach((item) => {
            item.isRead = true
          })
          uni.showToast({
            title: 'All marked as read',
            icon: 'success'
          })
        } catch (error) {
          uni.showToast({
            title: 'Operation failed',
            icon: 'none'
          })
        }
      }
    }
  })
}

onMounted(() => {
  fetchNotifications(true)
})

onPullDownRefreshHook(() => {
  onPullDownRefresh()
})

onReachBottom(() => {
  onLoadMore()
})
</script>

<style lang="scss" scoped>
.notification-page {
  min-height: 100vh;
  background-color: #f5f6fa;
}

.header-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20rpx 30rpx;
  background-color: #ffffff;
  border-bottom: 1rpx solid #eee;
}

.unread-badge {
  background-color: #ff4757;
  border-radius: 20rpx;
  padding: 6rpx 16rpx;
}

.unread-text {
  font-size: 24rpx;
  color: #ffffff;
}

.mark-all-btn {
  padding: 10rpx 20rpx;
}

.mark-all-text {
  font-size: 26rpx;
  color: #007aff;
}

.notification-list {
  height: calc(100vh - 90rpx);
}

.notification-item {
  display: flex;
  align-items: center;
  background-color: #ffffff;
  padding: 24rpx 30rpx;
  margin-bottom: 2rpx;
  position: relative;
}

.item-left {
  position: relative;
  margin-right: 20rpx;
  flex-shrink: 0;
}

.icon-wrapper {
  width: 80rpx;
  height: 80rpx;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: #e8f0fe;
}

.icon-system {
  background-color: #fff3e0;
}

.icon-exam {
  background-color: #e8f5e9;
}

.icon-feedback {
  background-color: #e3f2fd;
}

.icon-update {
  background-color: #f3e5f5;
}

.icon-promotion {
  background-color: #fce4ec;
}

.icon-text {
  font-size: 36rpx;
}

.unread-dot {
  position: absolute;
  top: 0;
  right: 0;
  width: 16rpx;
  height: 16rpx;
  border-radius: 50%;
  background-color: #ff4757;
}

.item-content {
  flex: 1;
  overflow: hidden;
}

.item-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8rpx;
}

.item-title {
  font-size: 28rpx;
  font-weight: 600;
  color: #333333;
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.title-read {
  font-weight: 400;
  color: #999999;
}

.item-time {
  font-size: 22rpx;
  color: #bbbbbb;
  margin-left: 16rpx;
  flex-shrink: 0;
}

.item-preview {
  font-size: 24rpx;
  color: #666666;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.preview-read {
  color: #aaaaaa;
}

.delete-btn {
  padding: 16rpx;
  margin-left: 16rpx;
  flex-shrink: 0;
}

.delete-text {
  font-size: 24rpx;
  color: #ff4757;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 120rpx 0;
}

.empty-icon {
  font-size: 80rpx;
  margin-bottom: 20rpx;
}

.empty-text {
  font-size: 28rpx;
  color: #999999;
}

.load-more {
  padding: 30rpx 0;
  text-align: center;
}

.load-more-text {
  font-size: 24rpx;
  color: #999999;
}
</style>
