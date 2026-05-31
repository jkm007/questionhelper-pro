<template>
  <view class="my-creation-page">
    <!-- Tabs -->
    <view class="tab-bar">
      <view
        v-for="tab in tabs"
        :key="tab.value"
        class="tab-item"
        :class="{ active: currentTab === tab.value }"
        @tap="switchTab(tab.value)"
      >
        <text class="tab-text" :class="{ active: currentTab === tab.value }">{{ tab.label }}</text>
        <view v-if="currentTab === tab.value" class="tab-line" />
      </view>
    </view>

    <!-- Creation List -->
    <scroll-view class="list-wrap" scroll-y @scrolltolower="loadMore" refresher-enabled @refresherrefresh="onRefresh" :refresher-triggered="refreshing">
      <view v-if="creationList.length === 0 && !loading" class="empty-state">
        <image class="empty-icon" src="/static/icon-empty.png" mode="aspectFit" />
        <text class="empty-text">No creations yet</text>
      </view>

      <view
        v-for="item in creationList"
        :key="item.id"
        class="creation-card"
        @tap="goToDetail(item)"
      >
        <view class="card-header">
          <text class="card-title">{{ item.title }}</text>
          <view class="status-badge" :class="item.status">
            <text class="status-text">{{ statusMap[item.status] || item.status }}</text>
          </view>
        </view>
        <view class="card-stats">
          <view class="stat-item">
            <image class="stat-icon" src="/static/icon-eye.png" mode="aspectFit" />
            <text class="stat-value">{{ item.viewCount || 0 }}</text>
          </view>
          <view class="stat-item">
            <image class="stat-icon" src="/static/icon-use.png" mode="aspectFit" />
            <text class="stat-value">{{ item.useCount || 0 }}</text>
          </view>
          <text class="card-time">{{ item.createTime }}</text>
        </view>
        <view class="card-actions">
          <view class="action-btn edit" @tap.stop="handleEdit(item)">
            <image class="action-icon" src="/static/icon-edit.png" mode="aspectFit" />
            <text class="action-text">Edit</text>
          </view>
          <view class="action-btn delete" @tap.stop="handleDelete(item)">
            <image class="action-icon" src="/static/icon-delete.png" mode="aspectFit" />
            <text class="action-text">Delete</text>
          </view>
        </view>
      </view>

      <view v-if="loading" class="loading-wrap">
        <text class="loading-text">Loading...</text>
      </view>
      <view v-if="noMore && creationList.length > 0" class="no-more-wrap">
        <text class="no-more-text">No more creations</text>
      </view>
    </scroll-view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getMyCreations, deleteCreation } from '@/api/question'

interface CreationItem {
  id: string
  title: string
  type: 'question' | 'paper' | 'exam'
  status: 'draft' | 'published' | 'rejected' | 'reviewing'
  viewCount: number
  useCount: number
  createTime: string
}

const tabs = [
  { label: 'Questions', value: 'questions' },
  { label: 'Papers', value: 'papers' },
  { label: 'Exams', value: 'exams' }
]

const statusMap: Record<string, string> = {
  draft: 'Draft',
  published: 'Published',
  rejected: 'Rejected',
  reviewing: 'Reviewing'
}

const currentTab = ref('questions')
const creationList = ref<CreationItem[]>([])
const loading = ref(false)
const noMore = ref(false)
const refreshing = ref(false)
const page = ref(1)
const pageSize = 20

onMounted(() => {
  fetchList()
})

async function fetchList() {
  if (loading.value) return
  loading.value = true
  try {
    const res = await getMyCreations({
      type: currentTab.value,
      page: page.value,
      pageSize
    })
    const list = res.data?.list || []
    if (page.value === 1) {
      creationList.value = list
    } else {
      creationList.value.push(...list)
    }
    if (list.length < pageSize) {
      noMore.value = true
    }
  } catch (e) {
    console.error('Failed to load creations', e)
  } finally {
    loading.value = false
  }
}

function switchTab(tab: string) {
  if (currentTab.value === tab) return
  currentTab.value = tab
  page.value = 1
  noMore.value = false
  creationList.value = []
  fetchList()
}

function loadMore() {
  if (noMore.value || loading.value) return
  page.value++
  fetchList()
}

async function onRefresh() {
  refreshing.value = true
  page.value = 1
  noMore.value = false
  await fetchList()
  refreshing.value = false
}

function goToDetail(item: CreationItem) {
  const typeMap: Record<string, string> = {
    questions: 'question',
    papers: 'paper',
    exams: 'exam'
  }
  const type = typeMap[currentTab.value] || 'question'
  uni.navigateTo({
    url: `/pages/${type}/detail?id=${item.id}`
  })
}

function handleEdit(item: CreationItem) {
  const typeMap: Record<string, string> = {
    questions: 'question',
    papers: 'paper',
    exams: 'exam'
  }
  const type = typeMap[currentTab.value] || 'question'
  uni.navigateTo({
    url: `/pages/${type}/edit?id=${item.id}`
  })
}

function handleDelete(item: CreationItem) {
  uni.showModal({
    title: 'Confirm Delete',
    content: `Are you sure you want to delete "${item.title}"? This action cannot be undone.`,
    confirmColor: '#e74c3c',
    success: async (res) => {
      if (res.confirm) {
        try {
          await deleteCreation(item.id)
          creationList.value = creationList.value.filter((c) => c.id !== item.id)
          uni.showToast({ title: 'Deleted', icon: 'success' })
        } catch (e) {
          console.error('Failed to delete', e)
          uni.showToast({ title: 'Delete failed', icon: 'none' })
        }
      }
    }
  })
}
</script>

<style scoped>
.my-creation-page {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background-color: #f5f5f5;
}

.tab-bar {
  display: flex;
  flex-direction: row;
  background-color: #fff;
  border-bottom: 1rpx solid #f0f0f0;
}

.tab-item {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 24rpx 0 16rpx;
  position: relative;
}

.tab-text {
  font-size: 30rpx;
  color: #666;
}

.tab-text.active {
  color: #4a90d9;
  font-weight: 600;
}

.tab-line {
  position: absolute;
  bottom: 0;
  width: 60rpx;
  height: 6rpx;
  background-color: #4a90d9;
  border-radius: 3rpx;
}

.list-wrap {
  flex: 1;
  padding: 20rpx 24rpx;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 200rpx 0;
}

.empty-icon {
  width: 160rpx;
  height: 160rpx;
  margin-bottom: 20rpx;
  opacity: 0.6;
}

.empty-text {
  font-size: 28rpx;
  color: #999;
}

.creation-card {
  background-color: #fff;
  border-radius: 16rpx;
  padding: 24rpx;
  margin-bottom: 20rpx;
}

.card-header {
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: space-between;
}

.card-title {
  font-size: 30rpx;
  color: #333;
  font-weight: 600;
  flex: 1;
  margin-right: 16rpx;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.status-badge {
  padding: 6rpx 16rpx;
  border-radius: 8rpx;
  flex-shrink: 0;
}

.status-badge.draft {
  background-color: #f5f5f5;
}

.status-badge.published {
  background-color: #e6f7e6;
}

.status-badge.rejected {
  background-color: #ffeaea;
}

.status-badge.reviewing {
  background-color: #fff3e0;
}

.status-text {
  font-size: 22rpx;
  color: #666;
}

.card-stats {
  display: flex;
  flex-direction: row;
  align-items: center;
  margin-top: 16rpx;
  gap: 28rpx;
}

.stat-item {
  display: flex;
  flex-direction: row;
  align-items: center;
}

.stat-icon {
  width: 28rpx;
  height: 28rpx;
  margin-right: 6rpx;
}

.stat-value {
  font-size: 24rpx;
  color: #999;
}

.card-time {
  font-size: 22rpx;
  color: #ccc;
  margin-left: auto;
}

.card-actions {
  display: flex;
  flex-direction: row;
  justify-content: flex-end;
  gap: 24rpx;
  margin-top: 20rpx;
  padding-top: 16rpx;
  border-top: 1rpx solid #f5f5f5;
}

.action-btn {
  display: flex;
  flex-direction: row;
  align-items: center;
  padding: 10rpx 20rpx;
  border-radius: 8rpx;
}

.action-btn.edit {
  background-color: #e6f0ff;
}

.action-btn.delete {
  background-color: #ffeaea;
}

.action-icon {
  width: 28rpx;
  height: 28rpx;
  margin-right: 6rpx;
}

.action-text {
  font-size: 24rpx;
  color: #666;
}

.loading-wrap,
.no-more-wrap {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 30rpx 0;
}

.loading-text,
.no-more-text {
  font-size: 24rpx;
  color: #999;
}
</style>
