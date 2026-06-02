<template>
  <view class="page">
    <!-- Search Bar -->
    <view class="search-bar">
      <view class="search-input-wrap">
        <input
          class="search-input"
          v-model="keyword"
          placeholder="搜索资源"
          confirm-type="search"
          @confirm="onSearch"
        />
      </view>
    </view>

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
        class="resource-card"
      >
        <view class="resource-icon-wrap">
          <text class="resource-icon">{{ getFileIcon(item.fileType) }}</text>
        </view>
        <view class="resource-info">
          <text class="resource-name">{{ item.fileName }}</text>
          <view class="resource-meta">
            <text class="meta-size">{{ item.fileSize }}</text>
            <text class="meta-time">{{ item.createdAt }}</text>
          </view>
        </view>
        <view class="resource-actions">
          <text class="action-link" @tap="onDownload(item)">下载</text>
          <text v-if="isTeacher" class="action-link danger" @tap="onDelete(item)">删除</text>
        </view>
      </view>

      <view v-if="!loading && list.length === 0" class="empty-state">
        <text class="empty-text">暂无资源</text>
      </view>

      <view v-if="list.length > 0" class="load-more">
        <text v-if="loading" class="load-more-text">加载中...</text>
        <text v-else-if="noMore" class="load-more-text">没有更多了</text>
      </view>
    </scroll-view>

    <view v-if="isTeacher" class="fab" @tap="onUpload">
      <text class="fab-icon">+</text>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

interface ResourceItem {
  id: string
  fileName: string
  fileType: string
  fileSize: string
  fileUrl: string
  createdAt: string
}

const classId = ref('')
const keyword = ref('')
const isTeacher = ref(false)
const list = ref<ResourceItem[]>([])
const loading = ref(false)
const refreshing = ref(false)
const noMore = ref(false)
const page = ref(1)
const pageSize = 20

function getFileIcon(type: string): string {
  const map: Record<string, string> = {
    pdf: 'PDF', doc: 'DOC', docx: 'DOC', xls: 'XLS', xlsx: 'XLS',
    ppt: 'PPT', pptx: 'PPT', jpg: 'IMG', png: 'IMG', mp4: 'VID'
  }
  return map[type?.toLowerCase()] || 'FILE'
}

async function loadData(reset = false) {
  if (loading.value) return
  if (reset) { page.value = 1; noMore.value = false; list.value = [] }
  loading.value = true
  try {
    // TODO: import { getClassResources } from '@/api/class'
    // const res = await getClassResources({ classId: classId.value, keyword: keyword.value, page: page.value, pageSize })
    const newData: ResourceItem[] = []
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

function onSearch() { loadData(true) }

async function onRefresh() {
  refreshing.value = true
  await loadData(true)
  refreshing.value = false
}

function onLoadMore() {
  if (!noMore.value && !loading.value) loadData()
}

function onDownload(item: ResourceItem) {
  uni.downloadFile({
    url: item.fileUrl,
    success: (res) => {
      if (res.statusCode === 200) {
        uni.openDocument({ filePath: res.tempFilePath, showMenu: true })
      }
    }
  })
}

function onDelete(item: ResourceItem) {
  uni.showModal({
    title: '删除资源',
    content: `确定删除「${item.fileName}」吗？`,
    success: async (res) => {
      if (res.confirm) {
        try {
          // TODO: import { deleteClassResource } from '@/api/class'
          // await deleteClassResource(item.id)
          uni.showToast({ title: '删除成功', icon: 'success' })
          loadData(true)
        } catch (e) {
          uni.showToast({ title: '删除失败', icon: 'none' })
        }
      }
    }
  })
}

function onUpload() {
  uni.chooseMessageFile?.({
    count: 1,
    success: async (res) => {
      const file = res.tempFiles[0]
      try {
        // TODO: import { uploadClassResource } from '@/api/class'
        // await uploadClassResource({ classId: classId.value, filePath: file.path, fileName: file.name })
        uni.showToast({ title: '上传成功', icon: 'success' })
        loadData(true)
      } catch (e) {
        uni.showToast({ title: '上传失败', icon: 'none' })
      }
    }
  })
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

.search-bar {
  padding: 16rpx 24rpx;
  background-color: #fff;
}

.search-input-wrap {
  background-color: #f0f1f5;
  border-radius: 36rpx;
  padding: 0 24rpx;
  height: 72rpx;
}

.search-input {
  height: 72rpx;
  font-size: 28rpx;
}

.list-container {
  height: calc(100vh - 104rpx);
  padding: 20rpx 24rpx;
}

.resource-card {
  display: flex;
  align-items: center;
  background-color: #fff;
  border-radius: 16rpx;
  padding: 24rpx;
  margin-bottom: 16rpx;
  box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.04);
}

.resource-icon-wrap {
  width: 80rpx;
  height: 80rpx;
  border-radius: 12rpx;
  background-color: #e8f4ff;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 20rpx;
  flex-shrink: 0;
}

.resource-icon {
  font-size: 24rpx;
  color: #1677ff;
  font-weight: 700;
}

.resource-info {
  flex: 1;
  min-width: 0;
}

.resource-name {
  font-size: 28rpx;
  color: #333;
  font-weight: 500;
  margin-bottom: 8rpx;
  display: block;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.resource-meta {
  display: flex;
  gap: 20rpx;
}

.meta-size, .meta-time {
  font-size: 22rpx;
  color: #999;
}

.resource-actions {
  display: flex;
  gap: 16rpx;
  flex-shrink: 0;
  margin-left: 16rpx;
}

.action-link {
  font-size: 24rpx;
  color: #1677ff;
  padding: 8rpx 16rpx;
}

.action-link.danger {
  color: #ff4d4f;
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
