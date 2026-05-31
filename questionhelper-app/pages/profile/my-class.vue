<template>
  <view class="my-class-page">
    <scroll-view class="class-list" scroll-y @scrolltolower="loadMore" refresher-enabled @refresherrefresh="onRefresh" :refresher-triggered="refreshing">
      <view v-if="classList.length === 0 && !loading" class="empty-state">
        <image class="empty-icon" src="/static/icon-empty.png" mode="aspectFit" />
        <text class="empty-text">No classes joined</text>
        <view class="join-btn" @tap="goToClassSquare">
          <text class="join-text">Browse Classes</text>
        </view>
      </view>

      <view
        v-for="cls in classList"
        :key="cls.id"
        class="class-card"
        @tap="goToDetail(cls)"
      >
        <image class="class-cover" :src="cls.cover || '/static/default-class-cover.png'" mode="aspectFill" />
        <view class="class-info">
          <view class="class-header">
            <text class="class-name">{{ cls.name }}</text>
            <view class="role-badge" :class="cls.role">
              <text class="role-text">{{ roleMap[cls.role] || cls.role }}</text>
            </view>
          </view>
          <view class="class-meta">
            <image class="meta-icon" src="/static/icon-members.png" mode="aspectFit" />
            <text class="meta-text">{{ cls.memberCount }} members</text>
          </view>
          <text class="class-desc" v-if="cls.description">{{ cls.description }}</text>
        </view>
      </view>

      <view v-if="loading" class="loading-wrap">
        <text class="loading-text">Loading...</text>
      </view>
      <view v-if="noMore && classList.length > 0" class="no-more-wrap">
        <text class="no-more-text">No more classes</text>
      </view>
    </scroll-view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getMyClasses } from '@/api/class'

interface ClassItem {
  id: string
  name: string
  cover: string
  description: string
  role: 'owner' | 'admin' | 'member'
  memberCount: number
}

const roleMap: Record<string, string> = {
  owner: 'Owner',
  admin: 'Admin',
  member: 'Member'
}

const classList = ref<ClassItem[]>([])
const loading = ref(false)
const noMore = ref(false)
const refreshing = ref(false)
const page = ref(1)
const pageSize = 20

onMounted(() => {
  fetchClasses()
})

async function fetchClasses() {
  if (loading.value) return
  loading.value = true
  try {
    const res = await getMyClasses({
      page: page.value,
      pageSize
    })
    const list = res.data?.list || []
    if (page.value === 1) {
      classList.value = list
    } else {
      classList.value.push(...list)
    }
    if (list.length < pageSize) {
      noMore.value = true
    }
  } catch (e) {
    console.error('Failed to load classes', e)
  } finally {
    loading.value = false
  }
}

function loadMore() {
  if (noMore.value || loading.value) return
  page.value++
  fetchClasses()
}

async function onRefresh() {
  refreshing.value = true
  page.value = 1
  noMore.value = false
  await fetchClasses()
  refreshing.value = false
}

function goToDetail(cls: ClassItem) {
  uni.navigateTo({
    url: `/pages/class/detail?id=${cls.id}`
  })
}

function goToClassSquare() {
  uni.navigateTo({
    url: '/pages/class/list'
  })
}
</script>

<style scoped>
.my-class-page {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background-color: #f5f5f5;
}

.class-list {
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
  margin-bottom: 30rpx;
}

.join-btn {
  padding: 18rpx 48rpx;
  background-color: #4a90d9;
  border-radius: 36rpx;
}

.join-text {
  font-size: 28rpx;
  color: #fff;
}

.class-card {
  display: flex;
  flex-direction: row;
  background-color: #fff;
  border-radius: 16rpx;
  overflow: hidden;
  margin-bottom: 20rpx;
}

.class-cover {
  width: 200rpx;
  height: 200rpx;
  flex-shrink: 0;
}

.class-info {
  flex: 1;
  padding: 20rpx 24rpx;
  display: flex;
  flex-direction: column;
  justify-content: center;
  min-width: 0;
}

.class-header {
  display: flex;
  flex-direction: row;
  align-items: center;
  gap: 12rpx;
}

.class-name {
  font-size: 30rpx;
  color: #333;
  font-weight: 600;
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.role-badge {
  padding: 4rpx 14rpx;
  border-radius: 8rpx;
  flex-shrink: 0;
}

.role-badge.owner {
  background-color: #fff3e0;
}

.role-badge.admin {
  background-color: #e6f0ff;
}

.role-badge.member {
  background-color: #f5f5f5;
}

.role-text {
  font-size: 20rpx;
  color: #666;
}

.class-meta {
  display: flex;
  flex-direction: row;
  align-items: center;
  margin-top: 12rpx;
}

.meta-icon {
  width: 28rpx;
  height: 28rpx;
  margin-right: 8rpx;
}

.meta-text {
  font-size: 24rpx;
  color: #999;
}

.class-desc {
  font-size: 24rpx;
  color: #999;
  margin-top: 10rpx;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
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
