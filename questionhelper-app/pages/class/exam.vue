<template>
  <view class="page">
    <!-- Status Tabs -->
    <view class="status-tabs">
      <view
        v-for="tab in statusTabs"
        :key="tab.value"
        :class="['status-tab', activeStatus === tab.value ? 'active' : '']"
        @tap="activeStatus = tab.value; onSearch()"
      >
        <text class="status-tab-text">{{ tab.label }}</text>
      </view>
    </view>

    <!-- Exam List -->
    <scroll-view
      class="list-container"
      scroll-y
      refresher-enabled
      :refresher-triggered="refreshing"
      @refresherrefresh="onRefresh"
      @scrolltolower="onLoadMore"
    >
      <view
        v-for="exam in list"
        :key="exam.id"
        class="exam-card"
        @tap="goDetail(exam.id)"
      >
        <view class="card-header">
          <text class="card-title">{{ exam.title }}</text>
          <view :class="['status-badge', 'status-' + exam.status]">
            <text class="status-text">{{ statusLabels[exam.status] }}</text>
          </view>
        </view>
        <view class="card-info">
          <view class="info-row">
            <text class="info-label">考试时间</text>
            <text class="info-value">{{ exam.startTime }} - {{ exam.endTime }}</text>
          </view>
          <view class="info-row">
            <text class="info-label">题目数量</text>
            <text class="info-value">{{ exam.questionCount }}题</text>
          </view>
          <view class="info-row">
            <text class="info-label">考试时长</text>
            <text class="info-value">{{ exam.duration }}分钟</text>
          </view>
          <view v-if="exam.status === 'finished' && exam.score !== undefined" class="info-row">
            <text class="info-label">我的成绩</text>
            <text class="info-value score">{{ exam.score }}分</text>
          </view>
        </view>
        <view class="card-footer">
          <text class="participant-count">{{ exam.participantCount }}人参加</text>
          <view v-if="exam.status === 'ongoing'" class="enter-btn" @tap.stop="onEnterExam(exam)">
            <text class="enter-btn-text">进入考试</text>
          </view>
          <view v-else-if="exam.status === 'finished'" class="enter-btn outline" @tap.stop="goDetail(exam.id)">
            <text class="enter-btn-text">查看详情</text>
          </view>
        </view>
      </view>

      <!-- Empty State -->
      <view v-if="!loading && list.length === 0" class="empty-state">
        <image class="empty-img" src="/static/empty/exam.png" mode="aspectFit" />
        <text class="empty-text">暂无考试</text>
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
import { getClassExamList } from '@/api/exam'

interface ExamItem {
  id: string
  title: string
  status: string
  startTime: string
  endTime: string
  questionCount: number
  duration: number
  score?: number
  participantCount: number
}

const statusTabs = [
  { label: '全部', value: '' },
  { label: '未开始', value: 'pending' },
  { label: '进行中', value: 'ongoing' },
  { label: '已结束', value: 'finished' }
]

const statusLabels: Record<string, string> = {
  'pending': '未开始',
  'ongoing': '进行中',
  'finished': '已结束'
}

const activeStatus = ref('')
const list = ref<ExamItem[]>([])
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
    const res = await getClassExamList(classId.value, {
      status: activeStatus.value,
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

function onEnterExam(exam: ExamItem) {
  uni.showModal({
    title: '进入考试',
    content: `确定进入「${exam.title}」吗？考试时长${exam.duration}分钟。`,
    success: (res) => {
      if (res.confirm) {
        uni.navigateTo({
          url: `/pages/exam/exam?id=${exam.id}`
        })
      }
    }
  })
}

function goDetail(id: string) {
  uni.navigateTo({ url: `/pages/exam/detail?id=${id}` })
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

.status-tabs {
  display: flex;
  background-color: #fff;
  padding: 0 16rpx;
  border-bottom: 1rpx solid #eee;
}

.status-tab {
  flex: 1;
  text-align: center;
  padding: 24rpx 0;
  position: relative;
}

.status-tab.active {
  border-bottom: 4rpx solid #1677ff;
}

.status-tab-text {
  font-size: 28rpx;
  color: #666;
}

.status-tab.active .status-tab-text {
  color: #1677ff;
  font-weight: 600;
}

.list-container {
  height: calc(100vh - 100rpx);
  padding: 20rpx 24rpx;
}

.exam-card {
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
  font-size: 32rpx;
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

.status-pending {
  background-color: #f0f1f5;
}

.status-pending .status-text {
  color: #999;
}

.status-ongoing {
  background-color: #e8f3ff;
}

.status-ongoing .status-text {
  color: #1677ff;
}

.status-finished {
  background-color: #f0f0f0;
}

.status-finished .status-text {
  color: #666;
}

.status-text {
  font-size: 22rpx;
}

.card-info {
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

.info-value.score {
  color: #1677ff;
  font-weight: 600;
  font-size: 30rpx;
}

.card-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-top: 20rpx;
  border-top: 1rpx solid #f5f5f5;
}

.participant-count {
  font-size: 24rpx;
  color: #999;
}

.enter-btn {
  padding: 12rpx 32rpx;
  background-color: #1677ff;
  border-radius: 28rpx;
}

.enter-btn.outline {
  background-color: #fff;
  border: 2rpx solid #1677ff;
}

.enter-btn-text {
  font-size: 26rpx;
  color: #fff;
}

.enter-btn.outline .enter-btn-text {
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
