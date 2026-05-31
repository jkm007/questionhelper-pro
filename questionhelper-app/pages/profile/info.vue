<template>
  <view class="info-page">
    <!-- Avatar -->
    <view class="avatar-section" @tap="chooseAvatar">
      <text class="section-label">Avatar</text>
      <view class="avatar-wrap">
        <image class="avatar" :src="form.avatar || '/static/default-avatar.png'" mode="aspectFill" />
        <image class="camera-icon" src="/static/icon-camera.png" mode="aspectFit" />
      </view>
    </view>

    <!-- Form Fields -->
    <view class="form-group">
      <view class="form-item">
        <text class="form-label">Nickname</text>
        <input
          class="form-input"
          v-model="form.nickname"
          placeholder="Enter nickname"
          maxlength="20"
        />
      </view>

      <view class="form-item" @tap="showGenderPicker">
        <text class="form-label">Gender</text>
        <view class="form-value-wrap">
          <text class="form-value">{{ genderText }}</text>
          <image class="form-arrow" src="/static/icon-arrow-right.png" mode="aspectFit" />
        </view>
      </view>

      <view class="form-item" @tap="showBirthdayPicker">
        <text class="form-label">Birthday</text>
        <view class="form-value-wrap">
          <text class="form-value" :class="{ placeholder: !form.birthday }">
            {{ form.birthday || 'Select birthday' }}
          </text>
          <image class="form-arrow" src="/static/icon-arrow-right.png" mode="aspectFit" />
        </view>
      </view>
    </view>

    <view class="form-group">
      <view class="form-item bio-item">
        <text class="form-label">Bio</text>
        <textarea
          class="form-textarea"
          v-model="form.bio"
          placeholder="Write something about yourself..."
          maxlength="200"
          :auto-height="false"
        />
        <text class="char-count">{{ (form.bio || '').length }}/200</text>
      </view>
    </view>

    <!-- Save Button -->
    <view class="save-btn" :class="{ disabled: saving }" @tap="handleSave">
      <text class="save-text">{{ saving ? 'Saving...' : 'Save' }}</text>
    </view>

    <!-- Gender Picker -->
    <uni-popup ref="genderPopup" type="bottom">
      <view class="picker-wrap">
        <view class="picker-header">
          <text class="picker-cancel" @tap="genderPopup?.close()">Cancel</text>
          <text class="picker-title">Select Gender</text>
          <text class="picker-confirm" @tap="confirmGender">Confirm</text>
        </view>
        <view class="picker-options">
          <view
            v-for="option in genderOptions"
            :key="option.value"
            class="picker-option"
            :class="{ selected: tempGender === option.value }"
            @tap="tempGender = option.value"
          >
            <text class="option-text">{{ option.label }}</text>
            <view v-if="tempGender === option.value" class="check-icon">
              <text class="check-text">OK</text>
            </view>
          </view>
        </view>
      </view>
    </uni-popup>

    <!-- Birthday Picker -->
    <uni-popup ref="birthdayPopup" type="bottom">
      <view class="picker-wrap">
        <picker
          mode="date"
          :value="form.birthday || '2000-01-01'"
          :end="today"
          @change="onBirthdayChange"
        >
          <view class="date-picker-inner">
            <text class="date-picker-text">Select Date</text>
          </view>
        </picker>
      </view>
    </uni-popup>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { getUserInfo, updateUserInfo } from '@/api/user'

const genderPopup = ref<any>(null)
const birthdayPopup = ref<any>(null)
const saving = ref(false)
const tempGender = ref(0)

const today = new Date().toISOString().split('T')[0]

const genderOptions = [
  { label: 'Not specified', value: 0 },
  { label: 'Male', value: 1 },
  { label: 'Female', value: 2 }
]

const form = ref({
  avatar: '',
  nickname: '',
  gender: 0,
  birthday: '',
  bio: ''
})

const genderText = computed(() => {
  const found = genderOptions.find((g) => g.value === form.value.gender)
  return found ? found.label : 'Not specified'
})

onMounted(() => {
  fetchUserInfo()
})

async function fetchUserInfo() {
  try {
    const res = await getUserInfo()
    if (res.data) {
      form.value = {
        avatar: res.data.avatar || '',
        nickname: res.data.nickname || '',
        gender: res.data.gender || 0,
        birthday: res.data.birthday || '',
        bio: res.data.bio || ''
      }
    }
  } catch (e) {
    console.error('Failed to load user info', e)
  }
}

function chooseAvatar() {
  uni.chooseImage({
    count: 1,
    sizeType: ['compressed'],
    sourceType: ['album', 'camera'],
    success: (res) => {
      const tempPath = res.tempFilePaths[0]
      uploadAvatar(tempPath)
    }
  })
}

async function uploadAvatar(filePath: string) {
  try {
    uni.showLoading({ title: 'Uploading...' })
    const uploadRes = await uni.uploadFile({
      url: '/api/user/avatar',
      filePath,
      name: 'file',
      header: {
        Authorization: `Bearer ${uni.getStorageSync('token')}`
      }
    })
    const data = JSON.parse(uploadRes.data)
    if (data.data?.url) {
      form.value.avatar = data.data.url
    }
    uni.hideLoading()
  } catch (e) {
    uni.hideLoading()
    console.error('Upload failed', e)
  }
}

function showGenderPicker() {
  tempGender.value = form.value.gender
  genderPopup.value?.open()
}

function confirmGender() {
  form.value.gender = tempGender.value
  genderPopup.value?.close()
}

function showBirthdayPicker() {
  birthdayPopup.value?.open()
}

function onBirthdayChange(e: any) {
  form.value.birthday = e.detail.value
  birthdayPopup.value?.close()
}

async function handleSave() {
  if (saving.value) return
  if (!form.value.nickname.trim()) {
    uni.showToast({ title: 'Please enter nickname', icon: 'none' })
    return
  }

  saving.value = true
  try {
    await updateUserInfo({
      nickname: form.value.nickname.trim(),
      gender: form.value.gender,
      birthday: form.value.birthday,
      bio: form.value.bio,
      avatar: form.value.avatar
    })
    uni.showToast({ title: 'Saved successfully', icon: 'success' })
    setTimeout(() => {
      uni.navigateBack()
    }, 1500)
  } catch (e) {
    console.error('Failed to save', e)
    uni.showToast({ title: 'Save failed', icon: 'none' })
  } finally {
    saving.value = false
  }
}
</script>

<style scoped>
.info-page {
  min-height: 100vh;
  background-color: #f5f5f5;
  padding-bottom: 40rpx;
}

.avatar-section {
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: space-between;
  background-color: #fff;
  padding: 30rpx;
  margin-bottom: 20rpx;
}

.section-label {
  font-size: 30rpx;
  color: #333;
}

.avatar-wrap {
  position: relative;
}

.avatar {
  width: 120rpx;
  height: 120rpx;
  border-radius: 60rpx;
}

.camera-icon {
  position: absolute;
  bottom: 0;
  right: 0;
  width: 36rpx;
  height: 36rpx;
}

.form-group {
  background-color: #fff;
  border-radius: 16rpx;
  margin: 20rpx 24rpx 0;
  overflow: hidden;
}

.form-item {
  display: flex;
  flex-direction: row;
  align-items: center;
  padding: 28rpx 30rpx;
  border-bottom: 1rpx solid #f0f0f0;
}

.form-item:last-child {
  border-bottom: none;
}

.bio-item {
  flex-direction: column;
  align-items: flex-start;
}

.form-label {
  font-size: 30rpx;
  color: #333;
  width: 140rpx;
  flex-shrink: 0;
}

.form-input {
  flex: 1;
  font-size: 28rpx;
  color: #333;
  text-align: right;
}

.form-value-wrap {
  flex: 1;
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: flex-end;
}

.form-value {
  font-size: 28rpx;
  color: #333;
}

.form-value.placeholder {
  color: #999;
}

.form-arrow {
  width: 28rpx;
  height: 28rpx;
  margin-left: 12rpx;
  opacity: 0.4;
}

.form-textarea {
  width: 100%;
  min-height: 160rpx;
  font-size: 28rpx;
  color: #333;
  margin-top: 16rpx;
  line-height: 1.6;
}

.char-count {
  font-size: 22rpx;
  color: #999;
  align-self: flex-end;
  margin-top: 8rpx;
}

.save-btn {
  margin: 48rpx 24rpx 0;
  background-color: #4a90d9;
  border-radius: 44rpx;
  padding: 26rpx 0;
  display: flex;
  align-items: center;
  justify-content: center;
}

.save-btn.disabled {
  opacity: 0.6;
}

.save-text {
  font-size: 32rpx;
  color: #fff;
  font-weight: 600;
}

.picker-wrap {
  background-color: #fff;
  border-radius: 24rpx 24rpx 0 0;
  padding-bottom: env(safe-area-inset-bottom);
}

.picker-header {
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: space-between;
  padding: 24rpx 30rpx;
  border-bottom: 1rpx solid #eee;
}

.picker-cancel {
  font-size: 28rpx;
  color: #999;
}

.picker-title {
  font-size: 30rpx;
  color: #333;
  font-weight: 600;
}

.picker-confirm {
  font-size: 28rpx;
  color: #4a90d9;
}

.picker-options {
  padding: 20rpx 0;
}

.picker-option {
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: space-between;
  padding: 28rpx 30rpx;
}

.picker-option.selected {
  background-color: #f5f8fc;
}

.option-text {
  font-size: 30rpx;
  color: #333;
}

.check-text {
  font-size: 26rpx;
  color: #4a90d9;
}

.date-picker-inner {
  padding: 40rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.date-picker-text {
  font-size: 30rpx;
  color: #4a90d9;
}
</style>
