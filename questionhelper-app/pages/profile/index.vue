<template>
  <view class="profile-page">
    <!-- User Card -->
    <view class="user-card" @tap="goToInfo">
      <image class="user-avatar" :src="userInfo.avatar || '/static/default-avatar.png'" mode="aspectFill" />
      <view class="user-info">
        <text class="user-nickname">{{ userInfo.nickname || 'Not logged in' }}</text>
        <view v-if="userInfo.level" class="level-badge">
          <text class="level-text">Lv.{{ userInfo.level }}</text>
        </view>
        <text class="user-id">ID: {{ userInfo.userId || '--' }}</text>
      </view>
      <image class="arrow-icon" src="/static/icon-arrow-right.png" mode="aspectFit" />
    </view>

    <!-- Menu Groups -->
    <view class="menu-group">
      <view class="menu-item" @tap="goTo('/pages/profile/info')">
        <image class="menu-icon" src="/static/icon-profile.png" mode="aspectFit" />
        <text class="menu-label">My Info</text>
        <image class="menu-arrow" src="/static/icon-arrow-right.png" mode="aspectFit" />
      </view>
      <view class="menu-item" @tap="goTo('/pages/profile/exam-history')">
        <image class="menu-icon" src="/static/icon-exam.png" mode="aspectFit" />
        <text class="menu-label">Exam History</text>
        <image class="menu-arrow" src="/static/icon-arrow-right.png" mode="aspectFit" />
      </view>
      <view class="menu-item" @tap="goTo('/pages/profile/statistics')">
        <image class="menu-icon" src="/static/icon-statistics.png" mode="aspectFit" />
        <text class="menu-label">My Statistics</text>
        <image class="menu-arrow" src="/static/icon-arrow-right.png" mode="aspectFit" />
      </view>
      <view class="menu-item" @tap="goTo('/pages/profile/favorites')">
        <image class="menu-icon" src="/static/icon-favorite.png" mode="aspectFit" />
        <text class="menu-label">My Favorites</text>
        <image class="menu-arrow" src="/static/icon-arrow-right.png" mode="aspectFit" />
      </view>
      <view class="menu-item" @tap="goTo('/pages/profile/my-class')">
        <image class="menu-icon" src="/static/icon-class.png" mode="aspectFit" />
        <text class="menu-label">My Classes</text>
        <image class="menu-arrow" src="/static/icon-arrow-right.png" mode="aspectFit" />
      </view>
      <view class="menu-item" @tap="goTo('/pages/profile/my-creation')">
        <image class="menu-icon" src="/static/icon-creation.png" mode="aspectFit" />
        <text class="menu-label">My Creations</text>
        <image class="menu-arrow" src="/static/icon-arrow-right.png" mode="aspectFit" />
      </view>
    </view>

    <view class="menu-group">
      <view class="menu-item" @tap="goTo('/pages/profile/auth')">
        <image class="menu-icon" src="/static/icon-auth.png" mode="aspectFit" />
        <text class="menu-label">Real-name Auth</text>
        <view v-if="userInfo.isVerified" class="status-tag verified">
          <text class="status-tag-text">Verified</text>
        </view>
        <image class="menu-arrow" src="/static/icon-arrow-right.png" mode="aspectFit" />
      </view>
      <view class="menu-item" @tap="goTo('/pages/profile/apply')">
        <image class="menu-icon" src="/static/icon-creator.png" mode="aspectFit" />
        <text class="menu-label">Apply Creator</text>
        <view v-if="userInfo.isCreator" class="status-tag creator">
          <text class="status-tag-text">Creator</text>
        </view>
        <image class="menu-arrow" src="/static/icon-arrow-right.png" mode="aspectFit" />
      </view>
    </view>

    <view class="menu-group">
      <view class="menu-item" @tap="goTo('/pages/profile/settings')">
        <image class="menu-icon" src="/static/icon-settings.png" mode="aspectFit" />
        <text class="menu-label">Settings</text>
        <image class="menu-arrow" src="/static/icon-arrow-right.png" mode="aspectFit" />
      </view>
      <view class="menu-item" @tap="goToAbout">
        <image class="menu-icon" src="/static/icon-about.png" mode="aspectFit" />
        <text class="menu-label">About</text>
        <image class="menu-arrow" src="/static/icon-arrow-right.png" mode="aspectFit" />
      </view>
    </view>

    <!-- Logout -->
    <view class="logout-btn" @tap="handleLogout">
      <text class="logout-text">Log Out</text>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import { getUserInfo } from '@/api/user'
import { useUserStore } from '@/store/modules/user'

interface UserInfo {
  userId: string
  avatar: string
  nickname: string
  level: number
  isVerified: boolean
  isCreator: boolean
}

const userStore = useUserStore()
const userInfo = ref<UserInfo>({
  userId: '',
  avatar: '',
  nickname: '',
  level: 0,
  isVerified: false,
  isCreator: false
})

onShow(() => {
  fetchUserInfo()
})

async function fetchUserInfo() {
  try {
    const res = await getUserInfo()
    if (res.data) {
      userInfo.value = res.data
      userStore.setUserInfo(res.data)
    }
  } catch (e) {
    console.error('Failed to load user info', e)
  }
}

function goToInfo() {
  uni.navigateTo({ url: '/pages/profile/info' })
}

function goTo(url: string) {
  uni.navigateTo({ url })
}

function goToAbout() {
  uni.showModal({
    title: 'About',
    content: 'QuestionHelper v1.0.0\nAn intelligent exam preparation platform.',
    showCancel: false
  })
}

function handleLogout() {
  uni.showModal({
    title: 'Confirm Logout',
    content: 'Are you sure you want to log out?',
    success: (res) => {
      if (res.confirm) {
        userStore.logout()
        uni.reLaunch({ url: '/pages/login/index' })
      }
    }
  })
}
</script>

<style scoped>
.profile-page {
  min-height: 100vh;
  background-color: #f5f5f5;
  padding-bottom: 40rpx;
}

.user-card {
  display: flex;
  flex-direction: row;
  align-items: center;
  background: linear-gradient(135deg, #4a90d9, #357abd);
  padding: 60rpx 40rpx 40rpx;
  margin-bottom: 20rpx;
}

.user-avatar {
  width: 120rpx;
  height: 120rpx;
  border-radius: 60rpx;
  border: 4rpx solid rgba(255, 255, 255, 0.6);
}

.user-info {
  flex: 1;
  margin-left: 24rpx;
}

.user-nickname {
  font-size: 36rpx;
  color: #fff;
  font-weight: 700;
}

.level-badge {
  display: inline-flex;
  background-color: rgba(255, 255, 255, 0.25);
  border-radius: 16rpx;
  padding: 4rpx 16rpx;
  margin-top: 8rpx;
}

.level-text {
  font-size: 22rpx;
  color: #fff;
}

.user-id {
  font-size: 24rpx;
  color: rgba(255, 255, 255, 0.7);
  margin-top: 6rpx;
}

.arrow-icon {
  width: 32rpx;
  height: 32rpx;
  opacity: 0.7;
}

.menu-group {
  background-color: #fff;
  border-radius: 16rpx;
  margin: 20rpx 24rpx 0;
  overflow: hidden;
}

.menu-item {
  display: flex;
  flex-direction: row;
  align-items: center;
  padding: 28rpx 30rpx;
  border-bottom: 1rpx solid #f0f0f0;
}

.menu-item:last-child {
  border-bottom: none;
}

.menu-icon {
  width: 44rpx;
  height: 44rpx;
  flex-shrink: 0;
}

.menu-label {
  flex: 1;
  font-size: 30rpx;
  color: #333;
  margin-left: 20rpx;
}

.menu-arrow {
  width: 28rpx;
  height: 28rpx;
  opacity: 0.4;
}

.status-tag {
  padding: 4rpx 14rpx;
  border-radius: 8rpx;
  margin-right: 12rpx;
}

.status-tag.verified {
  background-color: #e6f7e6;
}

.status-tag.creator {
  background-color: #fff3e0;
}

.status-tag-text {
  font-size: 22rpx;
  color: #4a90d9;
}

.logout-btn {
  margin: 40rpx 24rpx 0;
  background-color: #fff;
  border-radius: 16rpx;
  padding: 28rpx 0;
  display: flex;
  align-items: center;
  justify-content: center;
}

.logout-text {
  font-size: 30rpx;
  color: #e74c3c;
}
</style>
