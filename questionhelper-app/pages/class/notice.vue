<template>
  <view class="page">
    <!-- Notice List -->
    <scroll-view
      class="list-container"
      scroll-y
      refresher-enabled
      :refresher-triggered="refreshing"
      @refresherrefresh="onRefresh"
      @scrolltolower="onLoadMore"
    >
      <view
        v-for="notice in list"
        :key="notice.id"
        :class="['notice-card', notice.isRead ? '' : 'unread']"
        @tap="onTapNotice(notice)"
      >
        <view class="card-header">
          <view class="title-wrap">
            <view v-if="!notice.isRead" class="unread-dot"></view>
            <text class="card-title">{{ notice.title }}</text>
          </view>
          <text class="card-time">{{ notice.publishTime }}</text>
        </view>
        <text class="card-content">{{ notice.content }}</text>
        <view class="card-footer">
          <text class="publisher">发布者: {{ notice.publisherName }}</text>
          <view v-if="!notice.isRead" class="read-badge">
            <text class="read-badge-text">未读</text>
          </view>
        </view>
      </view>

      <!-- Empty State -->
      <view v-if="!loading && list.length === 0" class="empty-state">
        <image class="empty-img" src="/static/empty/notice.png" mode="aspectFit" />
        <text class="empty-text">暂无通知</text>
      </view>

      <!-- Load More -->
      <view v-if="list.length > 0" class="load-more">
        <text v-if="loading" class="load-more-text">加载中...</text>
        <text v-else-if="noMore" class="load-more-text">没有更多了</text>
      </view>
    </scroll-view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getClassNoticeList, markNoticeRead } from '@/api/class'

interface NoticeItem {
  id: string
  title: string
  content: string
  publishTime: string
  publisherName: string
  isRead: boolean
}

const list = ref<NoticeItem[]>([])
const loading = ref(false)
const refreshing = ref(false)
const noMore = ref(false)
const page = ref(1)
const pageSize = 20
const classId = ref('')

async function loadData(reset = false) {
  if (loading.value) return
  if (reset) {
    page.value = 1
    noMore.value = false
    list.value = []
  }
  loading.value = true
  try {
    const res = await getClassNoticeList(classId.value, {
      page: page.value,
      pageSize
    })
    const newData = res.data?.list || []
    if (reset) {
      list.value = newData
    } else {
      list.value = [...list.value, ...newData]
    }
    noMore.value = newData.length < pageSize
    page.value++
  } catch (e) {
    uni.showToast({ title: '加载失败', icon: 'none' })
  } finally {
    loading.value = false
  }
}

async function onRefresh() {
  refreshing.value = true
  await loadData(true)
  refreshing.value = false
}

function onLoadMore() {
  if (!noMore.value && !loading.value) {
    loadData()
  }
}

async function onTapNotice(notice: NoticeItem) {
  // Mark as read if unread
  if (!notice.isRead) {
    try {
      await markNoticeRead(notice.id)
      notice.isRead = true
    } catch (e) {
      // silent
    }
  }
  // Navigate to detail
  uni.navigateTo({
    url: `/pages/class/noticeDetail?noticeId=${notice.id}&classId=${classId.value}`
  })
}

onMounted(() => {
  const pages = getCurrentPages()
  const currentPage = pages[pages.length - 1] as any
  classId.value = currentPage?.options?.classId || ''
  if (classId.value) {
    loadData(true)
  } else {
    uni.showToast({ title: '参数错误', icon: 'none' })
    setTimeout(() => uni.navigateBack(), 1500)
  }
})
</script>

<style scoped>
.page {
  min-height: 100vh;
  background-color: #f5f6fa;
}

.list-container {
  height: 100vh;
  padding: 20rpx 24rpx;
}

.notice-card {
  background-color: #fff;
  border-radius: 16rpx;
  margin-bottom: 20rpx;
  padding: 28rpx;
  box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.04);
}

.notice-card.unread {
  border-left: 6rpx solid #1677ff;
}

.card-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  margin-bottom: 16rpx;
}

.title-wrap {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 10rpx;
  margin-right: 16rpx;
}

.unread-dot {
  width: 14rpx;
  height: 14rpx;
  border-radius: 50%;
  background-color: #ff4d4f;
  flex-shrink: 0;
}

.card-title {
  font-size: 30rpx;
  font-weight: 600;
  color: #333;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.card-time {
  font-size: 22rpx;
  color: #999;
  flex-shrink: 0;
}

.card-content {
  font-size: 26rpx;
  color: #666;
  line-height: 1.6;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
  margin-bottom: 16rpx;
}

.card-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-top: 16rpx;
  border-top: 1rpx solid #f5f5f5;
}

.publisher {
  font-size: 24rpx;
  color: #999;
}

.read-badge {
  padding: 4rpx 14rpx;
  background-color: #fff0e6;
  border-radius: 12rpx;
}

.read-badge-text {
  font-size: 20rpx;
  color: #fa8c16;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 120rpx 0;
}

.empty-img {
  width: 240rpx;
  height: 240rpx;
  margin-bottom: 24rpx;
}

.empty-text {
  font-size: 28rpx;
  color: #999;
}

.load-more {
  padding: 24rpx 0;
  text-align: center;
}

.load-more-text {
  font-size: 24rpx;
  color: #999;
}
</style>
