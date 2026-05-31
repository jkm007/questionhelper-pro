<template>
  <view class="page">
    <!-- Join Methods -->
    <view class="methods-section">
      <view class="method-card" @tap="showCodeInput = true">
        <view class="method-icon-wrap">
          <text class="method-icon">&#xe620;</text>
        </view>
        <text class="method-title">输入班级码</text>
        <text class="method-desc">输入老师提供的班级码加入</text>
      </view>
      <view class="method-card" @tap="onScanQR">
        <view class="method-icon-wrap">
          <text class="method-icon">&#xe621;</text>
        </view>
        <text class="method-title">扫码加入</text>
        <text class="method-desc">扫描班级二维码加入</text>
      </view>
    </view>

    <!-- Code Input Section -->
    <view v-if="showCodeInput" class="code-section">
      <text class="code-label">请输入班级码</text>
      <view class="code-input-wrap">
        <input
          class="code-input"
          v-model="classCode"
          placeholder="请输入6位班级码"
          :maxlength="10"
          @confirm="onSearchByCode"
        />
      </view>
      <text class="search-btn" @tap="onSearchByCode">查询班级</text>
    </view>

    <!-- Search Result -->
    <view v-if="searchResult" class="result-section">
      <text class="result-title">查询结果</text>
      <view class="result-card">
        <image class="result-cover" :src="searchResult.coverImage" mode="aspectFill" />
        <view class="result-body">
          <text class="result-name">{{ searchResult.name }}</text>
          <text class="result-creator">创建者: {{ searchResult.creatorName }}</text>
          <text class="result-members">成员: {{ searchResult.memberCount }}人</text>
          <text class="result-desc">{{ searchResult.description }}</text>
        </view>
      </view>
      <text class="confirm-join-btn" @tap="onConfirmJoin">确认加入</text>
    </view>

    <!-- Recent History -->
    <view v-if="historyList.length > 0" class="history-section">
      <text class="history-title">最近加入</text>
      <view v-for="item in historyList" :key="item.id" class="history-item">
        <image class="history-avatar" :src="item.coverImage" mode="aspectFill" />
        <view class="history-body">
          <text class="history-name">{{ item.name }}</text>
          <text class="history-time">{{ item.joinTime }}</text>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { searchClassByCode, joinClass, getJoinHistory } from '@/api/class'

interface SearchResult {
  id: string
  name: string
  coverImage: string
  creatorName: string
  memberCount: number
  description: string
  needApproval: boolean
}

interface HistoryItem {
  id: string
  name: string
  coverImage: string
  joinTime: string
}

const showCodeInput = ref(false)
const classCode = ref('')
const searchResult = ref<SearchResult | null>(null)
const historyList = ref<HistoryItem[]>([])

async function onSearchByCode() {
  if (!classCode.value.trim()) {
    uni.showToast({ title: '请输入班级码', icon: 'none' })
    return
  }
  uni.showLoading({ title: '查询中' })
  try {
    const res = await searchClassByCode(classCode.value.trim())
    searchResult.value = res.data
    if (!res.data) {
      uni.showToast({ title: '未找到该班级', icon: 'none' })
    }
  } catch (e) {
    uni.showToast({ title: '查询失败', icon: 'none' })
  } finally {
    uni.hideLoading()
  }
}

function onScanQR() {
  uni.scanCode({
    onlyFromCamera: false,
    scanType: ['qrCode'],
    success: (res) => {
      try {
        const data = JSON.parse(res.result)
        if (data.classId) {
          classCode.value = data.classId
          onSearchByCode()
        } else {
          uni.showToast({ title: '无效的二维码', icon: 'none' })
        }
      } catch (e) {
        uni.showToast({ title: '无效的二维码', icon: 'none' })
      }
    },
    fail: () => {
      uni.showToast({ title: '扫码取消', icon: 'none' })
    }
  })
}

async function onConfirmJoin() {
  if (!searchResult.value) return

  const doJoin = async () => {
    try {
      await joinClass(searchResult.value!.id)
      uni.showToast({ title: '加入成功', icon: 'success' })
      setTimeout(() => {
        uni.navigateTo({
          url: `/pages/class/detail?id=${searchResult.value!.id}`
        })
      }, 1500)
    } catch (e) {
      uni.showToast({ title: '加入失败', icon: 'none' })
    }
  }

  if (searchResult.value.needApproval) {
    uni.showModal({
      title: '加入班级',
      content: '该班级需要管理员审批，确定申请加入吗？',
      success: (res) => {
        if (res.confirm) doJoin()
      }
    })
  } else {
    uni.showModal({
      title: '加入班级',
      content: `确定加入「${searchResult.value.name}」吗？`,
      success: (res) => {
        if (res.confirm) doJoin()
      }
    })
  }
}

async function loadHistory() {
  try {
    const res = await getJoinHistory()
    historyList.value = res.data || []
  } catch (e) {
    // silent
  }
}

onMounted(() => {
  loadHistory()
})
</script>

<style scoped>
.page {
  min-height: 100vh;
  background-color: #f5f6fa;
  padding: 24rpx;
}

.methods-section {
  display: flex;
  gap: 20rpx;
  margin-bottom: 32rpx;
}

.method-card {
  flex: 1;
  background-color: #fff;
  border-radius: 16rpx;
  padding: 32rpx 20rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
  box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.04);
}

.method-icon-wrap {
  width: 88rpx;
  height: 88rpx;
  border-radius: 50%;
  background-color: #e8f3ff;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 16rpx;
}

.method-icon {
  font-size: 40rpx;
  color: #1677ff;
}

.method-title {
  font-size: 28rpx;
  font-weight: 600;
  color: #333;
  margin-bottom: 8rpx;
}

.method-desc {
  font-size: 22rpx;
  color: #999;
  text-align: center;
}

.code-section {
  background-color: #fff;
  border-radius: 16rpx;
  padding: 32rpx;
  margin-bottom: 24rpx;
  box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.04);
}

.code-label {
  font-size: 28rpx;
  font-weight: 600;
  color: #333;
  margin-bottom: 16rpx;
  display: block;
}

.code-input-wrap {
  border: 2rpx solid #e0e0e0;
  border-radius: 12rpx;
  padding: 0 20rpx;
  height: 80rpx;
  margin-bottom: 24rpx;
}

.code-input {
  width: 100%;
  height: 80rpx;
  font-size: 30rpx;
}

.search-btn {
  display: block;
  text-align: center;
  font-size: 28rpx;
  color: #fff;
  background-color: #1677ff;
  padding: 20rpx 0;
  border-radius: 36rpx;
}

.result-section {
  background-color: #fff;
  border-radius: 16rpx;
  padding: 32rpx;
  margin-bottom: 24rpx;
  box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.04);
}

.result-title {
  font-size: 28rpx;
  font-weight: 600;
  color: #333;
  margin-bottom: 20rpx;
  display: block;
}

.result-card {
  display: flex;
  gap: 20rpx;
  margin-bottom: 24rpx;
}

.result-cover {
  width: 160rpx;
  height: 160rpx;
  border-radius: 12rpx;
  flex-shrink: 0;
}

.result-body {
  flex: 1;
  min-width: 0;
}

.result-name {
  font-size: 30rpx;
  font-weight: 600;
  color: #333;
  margin-bottom: 8rpx;
  display: block;
}

.result-creator {
  font-size: 24rpx;
  color: #999;
  display: block;
  margin-bottom: 4rpx;
}

.result-members {
  font-size: 24rpx;
  color: #1677ff;
  display: block;
  margin-bottom: 8rpx;
}

.result-desc {
  font-size: 24rpx;
  color: #666;
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.confirm-join-btn {
  display: block;
  text-align: center;
  font-size: 28rpx;
  color: #fff;
  background-color: #52c41a;
  padding: 20rpx 0;
  border-radius: 36rpx;
}

.history-section {
  background-color: #fff;
  border-radius: 16rpx;
  padding: 32rpx;
  box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.04);
}

.history-title {
  font-size: 28rpx;
  font-weight: 600;
  color: #333;
  margin-bottom: 20rpx;
  display: block;
}

.history-item {
  display: flex;
  align-items: center;
  gap: 16rpx;
  padding: 16rpx 0;
  border-bottom: 1rpx solid #f5f5f5;
}

.history-item:last-child {
  border-bottom: none;
}

.history-avatar {
  width: 72rpx;
  height: 72rpx;
  border-radius: 50%;
  flex-shrink: 0;
}

.history-body {
  flex: 1;
}

.history-name {
  font-size: 28rpx;
  color: #333;
  display: block;
  margin-bottom: 4rpx;
}

.history-time {
  font-size: 22rpx;
  color: #999;
}
</style>
