<template>
  <view class="settings-page">
    <view class="section">
      <view class="section-title">通知类型开关</view>
      <view class="setting-item" v-for="item in notifyTypes" :key="item.type">
        <view class="item-left">
          <text class="item-icon">{{ item.icon }}</text>
          <text class="item-label">{{ item.label }}</text>
        </view>
        <switch :checked="item.enabled" @change="onToggle(item)" color="#4A90D9" />
      </view>
    </view>

    <view class="section">
      <view class="section-title">通知统计</view>
      <view class="stats-grid">
        <view class="stat-card">
          <text class="stat-num">{{ stats.total || 0 }}</text>
          <text class="stat-label">全部通知</text>
        </view>
        <view class="stat-card">
          <text class="stat-num">{{ stats.unread || 0 }}</text>
          <text class="stat-label">未读通知</text>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { getNotificationSettings, updateNotificationSettings, getNotificationStats } from '@/api/notification'

const notifyTypes = reactive([
  { type: 1, label: '系统通知', icon: '🔔', enabled: true },
  { type: 2, label: '考试通知', icon: '📝', enabled: true },
  { type: 3, label: '作业通知', icon: '📚', enabled: true },
  { type: 4, label: '班级通知', icon: '🏫', enabled: true },
  { type: 5, label: '评论通知', icon: '💬', enabled: true }
])

const stats = reactive({ total: 0, unread: 0 })

const loadSettings = async () => {
  try {
    const res = await getNotificationSettings()
    if (res.data) {
      const settings = res.data
      notifyTypes.forEach(item => {
        const key = `${['system', 'exam', 'homework', 'class', 'comment'][item.type - 1]}Enabled`
        if (settings[key] !== undefined) {
          item.enabled = settings[key]
        }
      })
    }
  } catch (e) {
    console.error('加载通知设置失败', e)
  }
}

const loadStats = async () => {
  try {
    const res = await getNotificationStats()
    if (res.data) {
      stats.total = res.data.total || 0
      stats.unread = res.data.unread || 0
    }
  } catch (e) {
    console.error('加载通知统计失败', e)
  }
}

const onToggle = async (item: { type: number; enabled: boolean }) => {
  item.enabled = !item.enabled
  try {
    const keys = ['system', 'exam', 'homework', 'class', 'comment']
    const data: any = {}
    data[`${keys[item.type - 1]}Enabled`] = item.enabled
    await updateNotificationSettings(data)
    uni.showToast({ title: '设置已更新', icon: 'success' })
  } catch (e) {
    item.enabled = !item.enabled
    uni.showToast({ title: '更新失败', icon: 'none' })
  }
}

onMounted(() => {
  loadSettings()
  loadStats()
})
</script>

<style scoped>
.settings-page {
  padding: 20rpx;
  background-color: #f5f5f5;
  min-height: 100vh;
}

.section {
  background-color: #ffffff;
  border-radius: 16rpx;
  margin-bottom: 24rpx;
  overflow: hidden;
}

.section-title {
  font-size: 28rpx;
  color: #999;
  padding: 24rpx 32rpx 12rpx;
}

.setting-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 28rpx 32rpx;
  border-bottom: 1rpx solid #f0f0f0;
}

.setting-item:last-child {
  border-bottom: none;
}

.item-left {
  display: flex;
  align-items: center;
  gap: 20rpx;
}

.item-icon {
  font-size: 40rpx;
}

.item-label {
  font-size: 30rpx;
  color: #333;
}

.stats-grid {
  display: flex;
  gap: 20rpx;
  padding: 20rpx 32rpx 32rpx;
}

.stat-card {
  flex: 1;
  background-color: #f8f9fa;
  border-radius: 12rpx;
  padding: 24rpx;
  text-align: center;
}

.stat-num {
  display: block;
  font-size: 40rpx;
  font-weight: bold;
  color: #4A90D9;
}

.stat-label {
  display: block;
  font-size: 24rpx;
  color: #999;
  margin-top: 8rpx;
}
</style>
