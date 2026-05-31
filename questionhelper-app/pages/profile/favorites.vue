<template>
  <view class="favorites-page">
    <!-- Tabs -->
    <view class="tab-bar">
      <view
        class="tab-item"
        :class="{ active: currentTab === 'questions' }"
        @tap="switchTab('questions')"
      >
        <text class="tab-text" :class="{ active: currentTab === 'questions' }">Questions</text>
        <view v-if="currentTab === 'questions'" class="tab-line" />
      </view>
      <view
        class="tab-item"
        :class="{ active: currentTab === 'exams' }"
        @tap="switchTab('exams')"
      >
        <text class="tab-text" :class="{ active: currentTab === 'exams' }">Exams</text>
        <view v-if="currentTab === 'exams'" class="tab-line" />
      </view>
    </view>

    <!-- Batch Actions -->
    <view v-if="list.length > 0" class="batch-bar">
      <view class="select-all" @tap="toggleSelectAll">
        <view class="checkbox" :class="{ checked: isAllSelected }">
          <text v-if="isAllSelected" class="check-mark">OK</text>
        </view>
        <text class="select-text">Select All</text>
      </view>
      <view
        class="batch-remove"
        :class="{ disabled: selectedIds.length === 0 }"
        @tap="batchRemove"
      >
        <text class="batch-remove-text">Remove ({{ selectedIds.length }})</text>
      </view>
    </view>

    <!-- List -->
    <scroll-view class="list-wrap" scroll-y @scrolltolower="loadMore" refresher-enabled @refresherrefresh="onRefresh" :refresher-triggered="refreshing">
      <view v-if="list.length === 0 && !loading" class="empty-state">
        <image class="empty-icon" src="/static/icon-empty.png" mode="aspectFit" />
        <text class="empty-text">No favorites yet</text>
      </view>

      <!-- Question Items -->
      <view v-if="currentTab === 'questions'">
        <view
          v-for="item in list"
          :key="item.id"
          class="fav-item"
        >
          <view class="fav-checkbox" @tap="toggleSelect(item.id)">
            <view class="checkbox" :class="{ checked: selectedIds.includes(item.id) }">
              <text v-if="selectedIds.includes(item.id)" class="check-mark">OK</text>
            </view>
          </view>
          <view class="fav-content" @tap="goToDetail(item)">
            <text class="fav-title">{{ item.title || item.content }}</text>
            <view class="fav-meta">
              <text class="fav-type">{{ item.typeName || 'Question' }}</text>
              <text class="fav-time">{{ item.favoriteTime }}</text>
            </view>
          </view>
          <view class="fav-remove" @tap.stop="removeSingle(item.id)">
            <text class="remove-text">Remove</text>
          </view>
        </view>
      </view>

      <!-- Exam Items -->
      <view v-if="currentTab === 'exams'">
        <view
          v-for="item in list"
          :key="item.id"
          class="fav-item"
        >
          <view class="fav-checkbox" @tap="toggleSelect(item.id)">
            <view class="checkbox" :class="{ checked: selectedIds.includes(item.id) }">
              <text v-if="selectedIds.includes(item.id)" class="check-mark">OK</text>
            </view>
          </view>
          <view class="fav-content" @tap="goToDetail(item)">
            <text class="fav-title">{{ item.title }}</text>
            <view class="fav-meta">
              <text class="fav-type">Exam</text>
              <text class="fav-info" v-if="item.questionCount">{{ item.questionCount }} questions</text>
              <text class="fav-time">{{ item.favoriteTime }}</text>
            </view>
          </view>
          <view class="fav-remove" @tap.stop="removeSingle(item.id)">
            <text class="remove-text">Remove</text>
          </view>
        </view>
      </view>

      <view v-if="loading" class="loading-wrap">
        <text class="loading-text">Loading...</text>
      </view>
      <view v-if="noMore && list.length > 0" class="no-more-wrap">
        <text class="no-more-text">No more items</text>
      </view>
    </scroll-view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { getFavorites, removeFavorite } from '@/api/user'

interface FavoriteItem {
  id: string
  title: string
  content?: string
  typeName?: string
  questionCount?: number
  favoriteTime: string
  targetId: string
}

const currentTab = ref<'questions' | 'exams'>('questions')
const list = ref<FavoriteItem[]>([])
const loading = ref(false)
const noMore = ref(false)
const refreshing = ref(false)
const page = ref(1)
const pageSize = 20
const selectedIds = ref<string[]>([])

const isAllSelected = computed(() => {
  return list.value.length > 0 && selectedIds.value.length === list.value.length
})

onMounted(() => {
  fetchList()
})

async function fetchList() {
  if (loading.value) return
  loading.value = true
  try {
    const res = await getFavorites({
      type: currentTab.value,
      page: page.value,
      pageSize
    })
    const items = res.data?.list || []
    if (page.value === 1) {
      list.value = items
    } else {
      list.value.push(...items)
    }
    if (items.length < pageSize) {
      noMore.value = true
    }
  } catch (e) {
    console.error('Failed to load favorites', e)
  } finally {
    loading.value = false
  }
}

function switchTab(tab: 'questions' | 'exams') {
  if (currentTab.value === tab) return
  currentTab.value = tab
  page.value = 1
  noMore.value = false
  list.value = []
  selectedIds.value = []
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
  selectedIds.value = []
  await fetchList()
  refreshing.value = false
}

function toggleSelect(id: string) {
  const index = selectedIds.value.indexOf(id)
  if (index >= 0) {
    selectedIds.value.splice(index, 1)
  } else {
    selectedIds.value.push(id)
  }
}

function toggleSelectAll() {
  if (isAllSelected.value) {
    selectedIds.value = []
  } else {
    selectedIds.value = list.value.map((item) => item.id)
  }
}

async function removeSingle(id: string) {
  uni.showModal({
    title: 'Confirm',
    content: 'Remove this item from favorites?',
    success: async (res) => {
      if (res.confirm) {
        try {
          await removeFavorite({ ids: [id] })
          list.value = list.value.filter((item) => item.id !== id)
          const idx = selectedIds.value.indexOf(id)
          if (idx >= 0) selectedIds.value.splice(idx, 1)
          uni.showToast({ title: 'Removed', icon: 'success' })
        } catch (e) {
          console.error('Failed to remove', e)
          uni.showToast({ title: 'Remove failed', icon: 'none' })
        }
      }
    }
  })
}

async function batchRemove() {
  if (selectedIds.value.length === 0) return
  uni.showModal({
    title: 'Confirm',
    content: `Remove ${selectedIds.value.length} items from favorites?`,
    success: async (res) => {
      if (res.confirm) {
        try {
          await removeFavorite({ ids: [...selectedIds.value] })
          list.value = list.value.filter((item) => !selectedIds.value.includes(item.id))
          selectedIds.value = []
          uni.showToast({ title: 'Removed', icon: 'success' })
        } catch (e) {
          console.error('Failed to batch remove', e)
          uni.showToast({ title: 'Remove failed', icon: 'none' })
        }
      }
    }
  })
}

function goToDetail(item: FavoriteItem) {
  if (currentTab.value === 'questions') {
    uni.navigateTo({ url: `/pages/question/detail?id=${item.targetId}` })
  } else {
    uni.navigateTo({ url: `/pages/exam/detail?id=${item.targetId}` })
  }
}
</script>

<style scoped>
.favorites-page {
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

.batch-bar {
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: space-between;
  padding: 16rpx 24rpx;
  background-color: #fff;
  border-bottom: 1rpx solid #f0f0f0;
}

.select-all {
  display: flex;
  flex-direction: row;
  align-items: center;
}

.checkbox {
  width: 36rpx;
  height: 36rpx;
  border: 2rpx solid #ccc;
  border-radius: 6rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 12rpx;
}

.checkbox.checked {
  background-color: #4a90d9;
  border-color: #4a90d9;
}

.check-mark {
  font-size: 20rpx;
  color: #fff;
}

.select-text {
  font-size: 26rpx;
  color: #666;
}

.batch-remove {
  padding: 10rpx 24rpx;
  background-color: #e74c3c;
  border-radius: 20rpx;
}

.batch-remove.disabled {
  opacity: 0.5;
}

.batch-remove-text {
  font-size: 24rpx;
  color: #fff;
}

.list-wrap {
  flex: 1;
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

.fav-item {
  display: flex;
  flex-direction: row;
  align-items: center;
  background-color: #fff;
  padding: 24rpx;
  margin: 16rpx 24rpx 0;
  border-radius: 12rpx;
}

.fav-checkbox {
  margin-right: 16rpx;
}

.fav-content {
  flex: 1;
  min-width: 0;
}

.fav-title {
  font-size: 28rpx;
  color: #333;
  font-weight: 500;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.fav-meta {
  display: flex;
  flex-direction: row;
  align-items: center;
  margin-top: 10rpx;
  gap: 16rpx;
}

.fav-type {
  font-size: 22rpx;
  color: #4a90d9;
  background-color: #e6f0ff;
  padding: 4rpx 12rpx;
  border-radius: 6rpx;
}

.fav-info {
  font-size: 22rpx;
  color: #999;
}

.fav-time {
  font-size: 22rpx;
  color: #bbb;
}

.fav-remove {
  margin-left: 16rpx;
  padding: 10rpx 16rpx;
}

.remove-text {
  font-size: 24rpx;
  color: #e74c3c;
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
