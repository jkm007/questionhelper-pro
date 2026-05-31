<template>
  <view class="page">
    <!-- Tabs -->
    <view class="tabs">
      <view :class="['tab-item', activeTab === 'my' ? 'active' : '']" @tap="activeTab = 'my'; onSearch()">
        <text class="tab-text">我的班级</text>
      </view>
      <view :class="['tab-item', activeTab === 'discover' ? 'active' : '']" @tap="activeTab = 'discover'; onSearch()">
        <text class="tab-text">发现班级</text>
      </view>
    </view>

    <!-- Search Bar -->
    <view class="search-bar">
      <view class="search-input-wrap">
        <text class="search-icon">&#xe610;</text>
        <input
          class="search-input"
          v-model="keyword"
          placeholder="搜索班级"
          confirm-type="search"
          @confirm="onSearch"
        />
        <text v-if="keyword" class="clear-btn" @tap="keyword = ''; onSearch()">&#xe621;</text>
      </view>
    </view>

    <!-- Class List -->
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
        class="class-card"
        @tap="goDetail(item.id)"
      >
        <image class="card-cover" :src="item.coverImage" mode="aspectFill" />
        <view class="card-body">
          <text class="card-name">{{ item.name }}</text>
          <view class="card-meta">
            <text class="meta-creator">{{ item.creatorName }}</text>
            <text class="meta-members">{{ item.memberCount }}人</text>
          </view>
          <text class="card-desc">{{ item.description }}</text>
          <view v-if="activeTab === 'discover'" class="card-action">
            <text class="join-btn" @tap.stop="onJoin(item)">加入</text>
          </view>
        </view>
      </view>

      <!-- Empty State -->
      <view v-if="!loading && list.length === 0" class="empty-state">
        <image class="empty-img" src="/static/empty/class.png" mode="aspectFit" />
        <text class="empty-text">{{ activeTab === 'my' ? '暂未加入任何班级' : '暂无发现班级' }}</text>
      </view>

      <!-- Load More -->
      <view v-if="list.length > 0" class="load-more">
        <text v-if="loading" class="load-more-text">加载中...</text>
        <text v-else-if="noMore" class="load-more-text">没有更多了</text>
      </view>
    </scroll-view>

    <!-- Create Class Button (my tab) -->
    <view v-if="activeTab === 'my'" class="fab" @tap="goCreate">
      <text class="fab-icon">+</text>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getMyClasses, discoverClasses, joinClass } from '@/api/class'

interface ClassItem {
  id: string
  name: string
  coverImage: string
  memberCount: number
  creatorName: string
  description: string
}

const activeTab = ref<'my' | 'discover'>('my')
const keyword = ref('')
const list = ref<ClassItem[]>([])
const loading = ref(false)
const refreshing = ref(false)
const noMore = ref(false)
const page = ref(1)
const pageSize = 20

async function loadData(reset = false) {
  if (loading.value) return
  if (reset) {
    page.value = 1
    noMore.value = false
    list.value = []
  }
  loading.value = true
  try {
    const api = activeTab.value === 'my' ? getMyClasses : discoverClasses
    const res = await api({
      keyword: keyword.value,
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

function onSearch() {
  loadData(true)
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

async function onJoin(item: ClassItem) {
  uni.showModal({
    title: '加入班级',
    content: `确定加入「${item.name}」吗？`,
    success: async (res) => {
      if (res.confirm) {
        try {
          await joinClass(item.id)
          uni.showToast({ title: '加入成功', icon: 'success' })
          loadData(true)
        } catch (e) {
          uni.showToast({ title: '加入失败', icon: 'none' })
        }
      }
    }
  })
}

function goDetail(id: string) {
  uni.navigateTo({ url: `/pages/class/detail?id=${id}` })
}

function goCreate() {
  uni.navigateTo({ url: '/pages/class/create' })
}

onMounted(() => {
  loadData(true)
})
</script>

<style scoped>
.page {
  min-height: 100vh;
  background-color: #f5f6fa;
}

.tabs {
  display: flex;
  background-color: #fff;
  border-bottom: 1rpx solid #eee;
}

.tab-item {
  flex: 1;
  text-align: center;
  padding: 24rpx 0;
  position: relative;
}

.tab-item.active {
  border-bottom: 4rpx solid #1677ff;
}

.tab-text {
  font-size: 30rpx;
  color: #666;
}

.tab-item.active .tab-text {
  color: #1677ff;
  font-weight: 600;
}

.search-bar {
  padding: 16rpx 24rpx;
  background-color: #fff;
}

.search-input-wrap {
  display: flex;
  align-items: center;
  background-color: #f0f1f5;
  border-radius: 36rpx;
  padding: 0 24rpx;
  height: 72rpx;
}

.search-icon {
  font-size: 32rpx;
  color: #999;
  margin-right: 12rpx;
}

.search-input {
  flex: 1;
  font-size: 28rpx;
  height: 72rpx;
}

.clear-btn {
  font-size: 28rpx;
  color: #999;
  padding: 8rpx;
}

.list-container {
  height: calc(100vh - 200rpx);
  padding: 20rpx 24rpx;
}

.class-card {
  display: flex;
  background-color: #fff;
  border-radius: 16rpx;
  margin-bottom: 20rpx;
  overflow: hidden;
  box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.04);
}

.card-cover {
  width: 200rpx;
  height: 200rpx;
  flex-shrink: 0;
}

.card-body {
  flex: 1;
  padding: 20rpx;
  min-width: 0;
}

.card-name {
  font-size: 30rpx;
  font-weight: 600;
  color: #333;
  margin-bottom: 10rpx;
  display: block;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.card-meta {
  display: flex;
  align-items: center;
  gap: 20rpx;
  margin-bottom: 10rpx;
}

.meta-creator {
  font-size: 24rpx;
  color: #999;
}

.meta-members {
  font-size: 24rpx;
  color: #1677ff;
}

.card-desc {
  font-size: 24rpx;
  color: #666;
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.card-action {
  margin-top: 12rpx;
  display: flex;
  justify-content: flex-end;
}

.join-btn {
  font-size: 24rpx;
  color: #fff;
  background-color: #1677ff;
  padding: 8rpx 28rpx;
  border-radius: 28rpx;
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
