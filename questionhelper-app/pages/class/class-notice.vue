<template>
  <view class="page">
    <scroll-view
      class="list-container"
      scroll-y
      refresher-enabled
      :refresher-triggered="refreshing"
      @refresherrefresh="onRefresh"
      @scrolltolower="onLoadMore"
    >
      <view
        v-for="item in list"
        :key="item.id"
        class="notice-card"
        @tap="onTapNotice(item)"
      >
        <view v-if="!item.isRead" class="unread-dot" />
        <view class="notice-header">
          <text class="notice-title">{{ item.title }}</text>
          <text class="notice-time">{{ item.createdAt }}</text>
        </view>
        <text class="notice-content">{{ item.content }}</text>
        <text class="notice-author">发布人：{{ item.authorName }}</text>
      </view>

      <view v-if="!loading && list.length === 0" class="empty-state">
        <text class="empty-text">暂无通知</text>
      </view>

      <view v-if="list.length > 0" class="load-more">
        <text v-if="loading" class="load-more-text">加载中...</text>
        <text v-else-if="noMore" class="load-more-text">没有更多了</text>
      </view>
    </scroll-view>

    <view v-if="isTeacher" class="fab" @tap="goCreate">
      <text class="fab-icon">+</text>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

interface NoticeItem {
  id: string
  title: string
  content: string
  authorName: string
  createdAt: string
  isRead: boolean
}

const classId = ref('')
const isTeacher = ref(false)
const list = ref<NoticeItem[]>([])
const loading = ref(false)
const refreshing = ref(false)
const noMore = ref(false)
const page = ref(1)
const pageSize = 20

async function loadData(reset = false) {
  if (loading.value) return
  if (reset) { page.value = 1; noMore.value = false; list.value = [] }
  loading.value = true
  try {
    // TODO: import { getClassNotices } from '@/api/class'
    // const res = await getClassNotices({ classId: classId.value, page: page.value, pageSize })
    const newData: NoticeItem[] = []
    if (reset) list.value = newData
    else list.value = [...list.value, ...newData]
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
  if (!noMore.value && !loading.value) loadData()
}

async function onTapNotice(item: NoticeItem) {
  if (!item.isRead) {
    item.isRead = true
    // TODO: import { markNoticeRead } from '@/api/class'
    // await markNoticeRead(item.id)
  }
  uni.showModal({
    title: item.title,
    content: item.content,
    showCancel: false
  })
}

function goCreate() {
  uni.navigateTo({ url: `/pages/class/notice-create?classId=${classId.value}` })
}

onMounted(() => {
  const pages = getCurrentPages()
  const current = pages[pages.length - 1] as any
  classId.value = current.options?.classId || ''
  isTeacher.value = uni.getStorageSync('role') === 'teacher'
  loadData(true)
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
  padding: 24rpx;
  margin-bottom: 20rpx;
  box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.04);
  position: relative;
}

.unread-dot {
  position: absolute;
  top: 24rpx;
  right: 24rpx;
  width: 16rpx;
  height: 16rpx;
  border-radius: 50%;
  background-color: #ff4d4f;
}

.notice-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  margin-bottom: 12rpx;
  padding-right: 32rpx;
}

.notice-title {
  font-size: 30rpx;
  font-weight: 600;
  color: #333;
  flex: 1;
}

.notice-time {
  font-size: 22rpx;
  color: #999;
  flex-shrink: 0;
  margin-left: 16rpx;
}

.notice-content {
  font-size: 26rpx;
  color: #666;
  line-height: 1.6;
  margin-bottom: 12rpx;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.notice-author {
  font-size: 22rpx;
  color: #999;
}

.empty-state {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 120rpx 0;
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

.fab {
  position: fixed;
  right: 40rpx;
  bottom: 80rpx;
  width: 100rpx;
  height: 100rpx;
  border-radius: 50%;
  background-color: #1677ff;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 8rpx 24rpx rgba(22, 119, 255, 0.4);
}

.fab-icon {
  font-size: 48rpx;
  color: #fff;
  font-weight: 300;
}
</style>
