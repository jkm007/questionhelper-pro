<template>
  <view class="page">
    <!-- Homework List -->
    <scroll-view
      class="list-container"
      scroll-y
      refresher-enabled
      :refresher-triggered="refreshing"
      @refresherrefresh="onRefresh"
      @scrolltolower="onLoadMore"
    >
      <view
        v-for="hw in list"
        :key="hw.id"
        class="homework-card"
      >
        <view class="card-header">
          <text class="card-title">{{ hw.title }}</text>
          <view :class="['status-badge', 'status-' + hw.status]">
            <text class="status-text">{{ statusLabels[hw.status] }}</text>
          </view>
        </view>
        <view class="card-body">
          <view class="info-row">
            <text class="info-label">截止时间</text>
            <text :class="['info-value', isDeadlineSoon(hw.deadline) ? 'deadline-warn' : '']">
              {{ hw.deadline }}
            </text>
          </view>
          <view class="info-row">
            <text class="info-label">提交状态</text>
            <view :class="['submit-badge', 'submit-' + hw.submitStatus]">
              <text class="submit-text">{{ submitLabels[hw.submitStatus] }}</text>
            </view>
          </view>
          <view v-if="hw.score !== undefined && hw.score !== null" class="info-row">
            <text class="info-label">得分</text>
            <text class="info-value score">{{ hw.score }}分</text>
          </view>
        </view>
        <view class="card-footer">
          <text class="creator-info">发布者: {{ hw.creatorName }}</text>
          <view class="footer-actions">
            <view
              v-if="hw.submitStatus === 'pending' && hw.status !== 'expired'"
              class="submit-btn"
              @tap="onSubmit(hw)"
            >
              <text class="submit-btn-text">提交作业</text>
            </view>
            <view
              v-if="hw.submitStatus === 'submitted' || hw.score !== undefined"
              class="view-btn"
              @tap="onViewResult(hw)"
            >
              <text class="view-btn-text">查看结果</text>
            </view>
          </view>
        </view>
      </view>

      <!-- Empty State -->
      <view v-if="!loading && list.length === 0" class="empty-state">
        <image class="empty-img" src="/static/empty/homework.png" mode="aspectFit" />
        <text class="empty-text">暂无作业</text>
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
import { getClassHomeworkList, submitHomework } from '@/api/class'

interface HomeworkItem {
  id: string
  title: string
  status: string
  deadline: string
  submitStatus: string
  score?: number
  creatorName: string
}

const statusLabels: Record<string, string> = {
  'ongoing': '进行中',
  'expired': '已截止',
  'closed': '已关闭'
}

const submitLabels: Record<string, string> = {
  'pending': '未提交',
  'submitted': '已提交',
  'graded': '已批改'
}

const list = ref<HomeworkItem[]>([])
const loading = ref(false)
const refreshing = ref(false)
const noMore = ref(false)
const page = ref(1)
const pageSize = 20
const classId = ref('')

function isDeadlineSoon(deadline: string): boolean {
  const deadlineTime = new Date(deadline).getTime()
  const now = Date.now()
  const diff = deadlineTime - now
  return diff > 0 && diff < 24 * 60 * 60 * 1000
}

async function loadData(reset = false) {
  if (loading.value) return
  if (reset) {
    page.value = 1
    noMore.value = false
    list.value = []
  }
  loading.value = true
  try {
    const res = await getClassHomeworkList(classId.value, {
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

function onSubmit(hw: HomeworkItem) {
  uni.navigateTo({
    url: `/pages/class/homework-submit?homeworkId=${hw.id}&classId=${classId.value}`
  })
}

function onViewResult(hw: HomeworkItem) {
  uni.navigateTo({
    url: `/pages/class/homework-result?homeworkId=${hw.id}&classId=${classId.value}`
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

.homework-card {
  background-color: #fff;
  border-radius: 16rpx;
  margin-bottom: 20rpx;
  padding: 28rpx;
  box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.04);
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 20rpx;
}

.card-title {
  flex: 1;
  font-size: 30rpx;
  font-weight: 600;
  color: #333;
  margin-right: 16rpx;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.status-badge {
  padding: 6rpx 16rpx;
  border-radius: 16rpx;
  flex-shrink: 0;
}

.status-ongoing {
  background-color: #e8f3ff;
}

.status-ongoing .status-text {
  color: #1677ff;
}

.status-expired {
  background-color: #fff1f0;
}

.status-expired .status-text {
  color: #ff4d4f;
}

.status-closed {
  background-color: #f0f0f0;
}

.status-closed .status-text {
  color: #999;
}

.status-text {
  font-size: 22rpx;
}

.card-body {
  margin-bottom: 20rpx;
}

.info-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8rpx 0;
}

.info-label {
  font-size: 26rpx;
  color: #999;
}

.info-value {
  font-size: 26rpx;
  color: #333;
}

.deadline-warn {
  color: #ff4d4f;
  font-weight: 600;
}

.info-value.score {
  color: #1677ff;
  font-weight: 600;
  font-size: 30rpx;
}

.submit-badge {
  padding: 4rpx 14rpx;
  border-radius: 12rpx;
}

.submit-pending {
  background-color: #fff0e6;
}

.submit-pending .submit-text {
  color: #fa8c16;
}

.submit-submitted {
  background-color: #e8f3ff;
}

.submit-submitted .submit-text {
  color: #1677ff;
}

.submit-graded {
  background-color: #f6ffed;
}

.submit-graded .submit-text {
  color: #52c41a;
}

.submit-text {
  font-size: 22rpx;
}

.card-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-top: 20rpx;
  border-top: 1rpx solid #f5f5f5;
}

.creator-info {
  font-size: 24rpx;
  color: #999;
}

.footer-actions {
  display: flex;
  gap: 16rpx;
}

.submit-btn {
  padding: 12rpx 28rpx;
  background-color: #1677ff;
  border-radius: 28rpx;
}

.submit-btn-text {
  font-size: 24rpx;
  color: #fff;
}

.view-btn {
  padding: 12rpx 28rpx;
  border: 2rpx solid #1677ff;
  border-radius: 28rpx;
}

.view-btn-text {
  font-size: 24rpx;
  color: #1677ff;
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
