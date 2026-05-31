<template>
  <view class="container">
    <!-- 自定义导航栏 -->
    <view class="navbar">
      <view class="navbar-content">
        <text class="navbar-title">题小助</text>
        <view class="navbar-right">
          <view class="notification-icon" @tap="goToNotification">
            <text class="iconfont icon-notification"></text>
            <view v-if="unreadCount > 0" class="badge">{{ unreadCount > 99 ? '99+' : unreadCount }}</view>
          </view>
        </view>
      </view>
    </view>

    <!-- 搜索栏 -->
    <view class="search-bar" @tap="goToSearch">
      <text class="search-icon">🔍</text>
      <text class="search-placeholder">搜索题目、试卷、用户...</text>
    </view>

    <!-- 快捷入口 -->
    <view class="quick-entry">
      <view class="entry-item" @tap="goToPractice">
        <view class="entry-icon">📝</view>
        <text class="entry-text">刷题练习</text>
      </view>
      <view class="entry-item" @tap="goToExam">
        <view class="entry-icon">📋</view>
        <text class="entry-text">考试中心</text>
      </view>
      <view class="entry-item" @tap="goToClass">
        <view class="entry-icon">👥</view>
        <text class="entry-text">班级学习</text>
      </view>
      <view class="entry-item" @tap="goToWrong">
        <view class="entry-icon">❌</view>
        <text class="entry-text">我的错题</text>
      </view>
    </view>

    <!-- 推荐题目 -->
    <view class="section">
      <view class="section-header">
        <text class="section-title">推荐题目</text>
        <text class="section-more" @tap="goToQuestionList">更多 ></text>
      </view>
      <view class="question-list">
        <QuestionCard
          v-for="item in recommendQuestions"
          :key="item.id"
          :question="item"
          @tap="goToQuestionDetail(item.id)"
        />
      </view>
    </view>

    <!-- 热门考试 -->
    <view class="section">
      <view class="section-header">
        <text class="section-title">热门考试</text>
        <text class="section-more" @tap="goToExam">更多 ></text>
      </view>
      <view class="exam-list">
        <ExamCard
          v-for="item in hotExams"
          :key="item.id"
          :exam="item"
          @tap="goToExamDetail(item.id)"
        />
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getQuestions } from '@/api/question'
import { getExams } from '@/api/exam'
import { getUnreadCount } from '@/api/notification'
import QuestionCard from '@/components/QuestionCard/index.vue'
import ExamCard from '@/components/ExamCard/index.vue'

const unreadCount = ref(0)
const recommendQuestions = ref([])
const hotExams = ref([])

onMounted(async () => {
  loadData()
})

const loadData = async () => {
  try {
    const [questionRes, examRes, unreadRes] = await Promise.all([
      getQuestions({ page: 1, pageSize: 5 }),
      getExams({ page: 1, pageSize: 3 }),
      getUnreadCount()
    ])
    recommendQuestions.value = questionRes.data.list || []
    hotExams.value = examRes.data.list || []
    unreadCount.value = unreadRes.data?.count || 0
  } catch (e) {
    console.error('加载数据失败', e)
  }
}

const goToNotification = () => {
  uni.navigateTo({ url: '/pages/notification/list' })
}

const goToSearch = () => {
  uni.navigateTo({ url: '/pages/question/search' })
}

const goToPractice = () => {
  uni.switchTab({ url: '/pages/practice/index' })
}

const goToExam = () => {
  uni.switchTab({ url: '/pages/exam/list' })
}

const goToClass = () => {
  uni.switchTab({ url: '/pages/class/list' })
}

const goToWrong = () => {
  uni.navigateTo({ url: '/pages/wrong/list' })
}

const goToQuestionList = () => {
  uni.switchTab({ url: '/pages/question/list' })
}

const goToQuestionDetail = (id: number) => {
  uni.navigateTo({ url: `/pages/question/detail?id=${id}` })
}

const goToExamDetail = (id: number) => {
  uni.navigateTo({ url: `/pages/exam/detail?id=${id}` })
}
</script>

<style lang="scss" scoped>
.container {
  padding-bottom: 20rpx;
}

.navbar {
  background-color: #4A90D9;
  padding-top: var(--status-bar-height);

  .navbar-content {
    display: flex;
    align-items: center;
    justify-content: space-between;
    height: 88rpx;
    padding: 0 30rpx;
  }

  .navbar-title {
    color: #ffffff;
    font-size: 36rpx;
    font-weight: bold;
  }

  .navbar-right {
    display: flex;
    align-items: center;
  }

  .notification-icon {
    position: relative;
    padding: 10rpx;
  }

  .icon-notification {
    font-size: 40rpx;
    color: #ffffff;
  }

  .badge {
    position: absolute;
    top: 0;
    right: 0;
    background-color: #F56C6C;
    color: #ffffff;
    font-size: 20rpx;
    padding: 2rpx 10rpx;
    border-radius: 20rpx;
  }
}

.search-bar {
  display: flex;
  align-items: center;
  margin: 20rpx 30rpx;
  padding: 20rpx 30rpx;
  background-color: #ffffff;
  border-radius: 40rpx;
  box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.05);

  .search-icon {
    margin-right: 16rpx;
  }

  .search-placeholder {
    color: #909399;
    font-size: 28rpx;
  }
}

.quick-entry {
  display: flex;
  justify-content: space-around;
  padding: 30rpx;
  margin: 0 30rpx;
  background-color: #ffffff;
  border-radius: 16rpx;
  box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.05);

  .entry-item {
    display: flex;
    flex-direction: column;
    align-items: center;
  }

  .entry-icon {
    font-size: 48rpx;
    margin-bottom: 12rpx;
  }

  .entry-text {
    font-size: 24rpx;
    color: #606266;
  }
}

.section {
  margin: 30rpx;

  .section-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20rpx;
  }

  .section-title {
    font-size: 32rpx;
    font-weight: bold;
    color: #303133;
  }

  .section-more {
    font-size: 24rpx;
    color: #909399;
  }
}
</style>
