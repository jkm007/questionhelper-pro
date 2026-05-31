<template>
  <view class="page">
    <!-- Loading -->
    <view v-if="loading" class="loading-wrap">
      <text class="loading-text">加载中...</text>
    </view>

    <template v-else-if="detail">
      <!-- Cover Header -->
      <view class="header">
        <image class="header-cover" :src="detail.coverImage" mode="aspectFill" />
        <view class="header-overlay"></view>
        <view class="header-content">
          <text class="header-name">{{ detail.name }}</text>
          <view class="header-meta">
            <text class="meta-creator">创建者: {{ detail.creatorName }}</text>
            <text class="meta-members">{{ detail.memberCount }}名成员</text>
          </view>
        </view>
      </view>

      <!-- Class Info -->
      <view class="info-section">
        <text class="info-desc">{{ detail.description }}</text>
        <view class="info-stats">
          <view class="stat-item">
            <text class="stat-num">{{ detail.resourceCount || 0 }}</text>
            <text class="stat-label">资源</text>
          </view>
          <view class="stat-divider"></view>
          <view class="stat-item">
            <text class="stat-num">{{ detail.examCount || 0 }}</text>
            <text class="stat-label">考试</text>
          </view>
          <view class="stat-divider"></view>
          <view class="stat-item">
            <text class="stat-num">{{ detail.homeworkCount || 0 }}</text>
            <text class="stat-label">作业</text>
          </view>
          <view class="stat-divider"></view>
          <view class="stat-item">
            <text class="stat-num">{{ detail.noticeCount || 0 }}</text>
            <text class="stat-label">通知</text>
          </view>
        </view>
      </view>

      <!-- Action Buttons -->
      <view class="action-section">
        <view v-if="!detail.isMember" class="join-btn" @tap="onJoin">
          <text class="join-btn-text">加入班级</text>
        </view>
        <template v-else>
          <view class="member-btn" @tap="goMember">
            <text class="member-btn-text">成员列表</text>
          </view>
          <view v-if="!detail.isCreator" class="leave-btn" @tap="onLeave">
            <text class="leave-btn-text">退出班级</text>
          </view>
        </template>
      </view>

      <!-- Tab Sections -->
      <view class="tab-sections">
        <!-- Resources -->
        <view class="section-card" @tap="goSection('resource')">
          <view class="section-icon-wrap resource-icon">
            <text class="section-icon">&#xe618;</text>
          </view>
          <view class="section-body">
            <text class="section-name">班级资源</text>
            <text class="section-count">{{ detail.resourceCount || 0 }}个资源</text>
          </view>
          <text class="section-arrow">&#xe61e;</text>
        </view>

        <!-- Exams -->
        <view class="section-card" @tap="goSection('exam')">
          <view class="section-icon-wrap exam-icon">
            <text class="section-icon">&#xe619;</text>
          </view>
          <view class="section-body">
            <text class="section-name">班级考试</text>
            <text class="section-count">{{ detail.examCount || 0 }}场考试</text>
          </view>
          <text class="section-arrow">&#xe61e;</text>
        </view>

        <!-- Homework -->
        <view class="section-card" @tap="goSection('homework')">
          <view class="section-icon-wrap homework-icon">
            <text class="section-icon">&#xe61a;</text>
          </view>
          <view class="section-body">
            <text class="section-name">班级作业</text>
            <text class="section-count">{{ detail.homeworkCount || 0 }}个作业</text>
          </view>
          <text class="section-arrow">&#xe61e;</text>
        </view>

        <!-- Notices -->
        <view class="section-card" @tap="goSection('notice')">
          <view class="section-icon-wrap notice-icon">
            <text class="section-icon">&#xe61b;</text>
          </view>
          <view class="section-body">
            <text class="section-name">班级通知</text>
            <text class="section-count">{{ detail.noticeCount || 0 }}条通知</text>
          </view>
          <text class="section-arrow">&#xe61e;</text>
        </view>

        <!-- Discussions -->
        <view class="section-card" @tap="goSection('discussion')">
          <view class="section-icon-wrap discuss-icon">
            <text class="section-icon">&#xe61c;</text>
          </view>
          <view class="section-body">
            <text class="section-name">班级讨论</text>
            <text class="section-count">{{ detail.discussionCount || 0 }}条讨论</text>
          </view>
          <text class="section-arrow">&#xe61e;</text>
        </view>
      </view>
    </template>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getClassDetail, joinClass, leaveClass } from '@/api/class'

interface ClassDetail {
  id: string
  name: string
  coverImage: string
  description: string
  creatorName: string
  memberCount: number
  resourceCount: number
  examCount: number
  homeworkCount: number
  noticeCount: number
  discussionCount: number
  isMember: boolean
  isCreator: boolean
}

const loading = ref(true)
const detail = ref<ClassDetail | null>(null)
const classId = ref('')

async function loadDetail() {
  loading.value = true
  try {
    const res = await getClassDetail(classId.value)
    detail.value = res.data
  } catch (e) {
    uni.showToast({ title: '加载失败', icon: 'none' })
  } finally {
    loading.value = false
  }
}

async function onJoin() {
  if (!detail.value) return
  try {
    await joinClass(detail.value.id)
    uni.showToast({ title: '加入成功', icon: 'success' })
    loadDetail()
  } catch (e) {
    uni.showToast({ title: '加入失败', icon: 'none' })
  }
}

async function onLeave() {
  if (!detail.value) return
  uni.showModal({
    title: '退出班级',
    content: `确定退出「${detail.value.name}」吗？`,
    success: async (res) => {
      if (res.confirm) {
        try {
          await leaveClass(detail.value!.id)
          uni.showToast({ title: '已退出', icon: 'success' })
          setTimeout(() => uni.navigateBack(), 1500)
        } catch (e) {
          uni.showToast({ title: '操作失败', icon: 'none' })
        }
      }
    }
  })
}

function goMember() {
  uni.navigateTo({ url: `/pages/class/member?classId=${classId.value}` })
}

function goSection(type: string) {
  const urlMap: Record<string, string> = {
    resource: `/pages/class/resource?classId=${classId.value}`,
    exam: `/pages/class/exam?classId=${classId.value}`,
    homework: `/pages/class/homework?classId=${classId.value}`,
    notice: `/pages/class/notice?classId=${classId.value}`,
    discussion: `/pages/class/discussion?classId=${classId.value}`
  }
  const url = urlMap[type]
  if (url) {
    uni.navigateTo({ url })
  }
}

onMounted(() => {
  const pages = getCurrentPages()
  const currentPage = pages[pages.length - 1] as any
  classId.value = currentPage?.options?.id || ''
  if (classId.value) {
    loadDetail()
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

.loading-wrap {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 400rpx;
}

.loading-text {
  font-size: 28rpx;
  color: #999;
}

.header {
  position: relative;
  height: 360rpx;
  overflow: hidden;
}

.header-cover {
  width: 100%;
  height: 100%;
}

.header-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(transparent 40%, rgba(0, 0, 0, 0.6));
}

.header-content {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  padding: 32rpx;
}

.header-name {
  font-size: 40rpx;
  font-weight: 700;
  color: #fff;
  margin-bottom: 12rpx;
  display: block;
}

.header-meta {
  display: flex;
  align-items: center;
  gap: 24rpx;
}

.meta-creator,
.meta-members {
  font-size: 24rpx;
  color: rgba(255, 255, 255, 0.85);
}

.info-section {
  background-color: #fff;
  padding: 28rpx;
  margin-bottom: 20rpx;
}

.info-desc {
  font-size: 28rpx;
  color: #666;
  line-height: 1.6;
  margin-bottom: 24rpx;
  display: block;
}

.info-stats {
  display: flex;
  align-items: center;
  justify-content: space-around;
  background-color: #f8f9fc;
  border-radius: 16rpx;
  padding: 24rpx 0;
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8rpx;
}

.stat-num {
  font-size: 36rpx;
  font-weight: 700;
  color: #1677ff;
}

.stat-label {
  font-size: 24rpx;
  color: #999;
}

.stat-divider {
  width: 1rpx;
  height: 48rpx;
  background-color: #e8e8e8;
}

.action-section {
  display: flex;
  gap: 20rpx;
  padding: 0 24rpx 20rpx;
}

.join-btn {
  flex: 1;
  background-color: #1677ff;
  border-radius: 44rpx;
  padding: 24rpx 0;
  text-align: center;
}

.join-btn-text {
  font-size: 30rpx;
  font-weight: 600;
  color: #fff;
}

.member-btn {
  flex: 1;
  background-color: #1677ff;
  border-radius: 44rpx;
  padding: 24rpx 0;
  text-align: center;
}

.member-btn-text {
  font-size: 30rpx;
  font-weight: 600;
  color: #fff;
}

.leave-btn {
  flex: 1;
  background-color: #fff;
  border: 2rpx solid #ff4d4f;
  border-radius: 44rpx;
  padding: 24rpx 0;
  text-align: center;
}

.leave-btn-text {
  font-size: 30rpx;
  font-weight: 600;
  color: #ff4d4f;
}

.tab-sections {
  background-color: #fff;
  border-radius: 16rpx;
  margin: 0 24rpx;
  overflow: hidden;
}

.section-card {
  display: flex;
  align-items: center;
  padding: 28rpx;
  border-bottom: 1rpx solid #f5f5f5;
}

.section-card:last-child {
  border-bottom: none;
}

.section-icon-wrap {
  width: 72rpx;
  height: 72rpx;
  border-radius: 16rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 20rpx;
  flex-shrink: 0;
}

.resource-icon { background-color: #e8f3ff; }
.exam-icon { background-color: #fff0e6; }
.homework-icon { background-color: #f6ffed; }
.notice-icon { background-color: #f9f0ff; }
.discuss-icon { background-color: #fff1f0; }

.section-icon {
  font-size: 36rpx;
  color: #1677ff;
}

.section-body {
  flex: 1;
  min-width: 0;
}

.section-name {
  font-size: 30rpx;
  font-weight: 600;
  color: #333;
  display: block;
  margin-bottom: 4rpx;
}

.section-count {
  font-size: 24rpx;
  color: #999;
}

.section-arrow {
  font-size: 28rpx;
  color: #ccc;
  flex-shrink: 0;
  margin-left: 16rpx;
}
</style>
