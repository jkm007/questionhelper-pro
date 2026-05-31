<template>
  <view class="settings-page">
    <!-- Notification Settings -->
    <view class="setting-group">
      <text class="group-title">Notifications</text>
      <view class="setting-item">
        <text class="setting-label">Push Notifications</text>
        <switch
          :checked="settings.pushEnabled"
          color="#4a90d9"
          @change="onToggle('pushEnabled', $event)"
        />
      </view>
      <view class="setting-item">
        <text class="setting-label">Sound</text>
        <switch
          :checked="settings.soundEnabled"
          color="#4a90d9"
          @change="onToggle('soundEnabled', $event)"
        />
      </view>
      <view class="setting-item">
        <text class="setting-label">Vibration</text>
        <switch
          :checked="settings.vibrationEnabled"
          color="#4a90d9"
          @change="onToggle('vibrationEnabled', $event)"
        />
      </view>
    </view>

    <!-- Language -->
    <view class="setting-group">
      <text class="group-title">General</text>
      <view class="setting-item" @tap="showLanguagePicker">
        <text class="setting-label">Language</text>
        <view class="setting-value-wrap">
          <text class="setting-value">{{ currentLanguageText }}</text>
          <image class="setting-arrow" src="/static/icon-arrow-right.png" mode="aspectFit" />
        </view>
      </view>
    </view>

    <!-- Cache & Storage -->
    <view class="setting-group">
      <text class="group-title">Storage</text>
      <view class="setting-item">
        <text class="setting-label">Cache Size</text>
        <view class="setting-value-wrap">
          <text class="setting-value">{{ cacheSize }}</text>
          <view class="clear-btn" @tap="clearCache">
            <text class="clear-text">Clear</text>
          </view>
        </view>
      </view>
    </view>

    <!-- Legal -->
    <view class="setting-group">
      <text class="group-title">Legal</text>
      <view class="setting-item" @tap="openPage('privacy')">
        <text class="setting-label">Privacy Policy</text>
        <image class="setting-arrow" src="/static/icon-arrow-right.png" mode="aspectFit" />
      </view>
      <view class="setting-item" @tap="openPage('terms')">
        <text class="setting-label">Terms of Service</text>
        <image class="setting-arrow" src="/static/icon-arrow-right.png" mode="aspectFit" />
      </view>
    </view>

    <!-- Version -->
    <view class="setting-group">
      <text class="group-title">About</text>
      <view class="setting-item">
        <text class="setting-label">Version</text>
        <text class="setting-value">{{ appVersion }}</text>
      </view>
    </view>

    <!-- Account Deactivation -->
    <view class="deactivate-section">
      <view class="deactivate-btn" @tap="handleDeactivate">
        <text class="deactivate-text">Deactivate Account</text>
      </view>
    </view>

    <!-- Language Picker -->
    <uni-popup ref="langPopup" type="bottom">
      <view class="picker-wrap">
        <view class="picker-header">
          <text class="picker-cancel" @tap="langPopup?.close()">Cancel</text>
          <text class="picker-title">Select Language</text>
          <text class="picker-confirm" @tap="confirmLanguage">Confirm</text>
        </view>
        <view class="picker-options">
          <view
            v-for="lang in languages"
            :key="lang.value"
            class="picker-option"
            :class="{ selected: tempLanguage === lang.value }"
            @tap="tempLanguage = lang.value"
          >
            <text class="option-text">{{ lang.label }}</text>
            <view v-if="tempLanguage === lang.value" class="check-icon">
              <text class="check-text">OK</text>
            </view>
          </view>
        </view>
      </view>
    </uni-popup>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { getSettings, updateSettings, deactivateAccount } from '@/api/user'

const langPopup = ref<any>(null)

const languages = [
  { label: 'English', value: 'en' },
  { label: 'Chinese (Simplified)', value: 'zh-CN' },
  { label: 'Chinese (Traditional)', value: 'zh-TW' },
  { label: 'Japanese', value: 'ja' },
  { label: 'Korean', value: 'ko' }
]

const settings = ref({
  pushEnabled: true,
  soundEnabled: true,
  vibrationEnabled: true,
  language: 'zh-CN'
})

const tempLanguage = ref('zh-CN')
const cacheSize = ref('Calculating...')
const appVersion = ref('1.0.0')

const currentLanguageText = computed(() => {
  const found = languages.find((l) => l.value === settings.value.language)
  return found ? found.label : 'English'
})

onMounted(() => {
  fetchSettings()
  calculateCacheSize()
  getAppVersion()
})

async function fetchSettings() {
  try {
    const res = await getSettings()
    if (res.data) {
      settings.value = {
        pushEnabled: res.data.pushEnabled ?? true,
        soundEnabled: res.data.soundEnabled ?? true,
        vibrationEnabled: res.data.vibrationEnabled ?? true,
        language: res.data.language || 'zh-CN'
      }
    }
  } catch (e) {
    console.error('Failed to load settings', e)
  }
}

async function onToggle(key: string, event: any) {
  const value = event.detail.value
  ;(settings.value as any)[key] = value
  try {
    await updateSettings({ [key]: value })
  } catch (e) {
    console.error('Failed to update setting', e)
    ;(settings.value as any)[key] = !value
    uni.showToast({ title: 'Update failed', icon: 'none' })
  }
}

function showLanguagePicker() {
  tempLanguage.value = settings.value.language
  langPopup.value?.open()
}

async function confirmLanguage() {
  settings.value.language = tempLanguage.value
  langPopup.value?.close()
  try {
    await updateSettings({ language: tempLanguage.value })
    uni.showToast({ title: 'Language updated', icon: 'success' })
  } catch (e) {
    console.error('Failed to update language', e)
  }
}

function calculateCacheSize() {
  try {
    const res = uni.getStorageInfoSync()
    const sizeKB = res.currentSize || 0
    if (sizeKB > 1024) {
      cacheSize.value = (sizeKB / 1024).toFixed(1) + ' MB'
    } else {
      cacheSize.value = sizeKB + ' KB'
    }
  } catch (e) {
    cacheSize.value = 'Unknown'
  }
}

function clearCache() {
  uni.showModal({
    title: 'Clear Cache',
    content: 'Are you sure you want to clear the cache? This will not affect your account data.',
    success: (res) => {
      if (res.confirm) {
        try {
          uni.clearStorageSync()
          cacheSize.value = '0 KB'
          uni.showToast({ title: 'Cache cleared', icon: 'success' })
        } catch (e) {
          console.error('Failed to clear cache', e)
        }
      }
    }
  })
}

function getAppVersion() {
  try {
    const accountInfo = uni.getAccountInfoSync()
    appVersion.value = accountInfo.miniProgram?.version || '1.0.0'
  } catch (e) {
    appVersion.value = '1.0.0'
  }
}

function openPage(type: string) {
  const urlMap: Record<string, string> = {
    privacy: '/pages/webview/index?url=/privacy-policy',
    terms: '/pages/webview/index?url=/terms-of-service'
  }
  uni.navigateTo({ url: urlMap[type] || '' })
}

function handleDeactivate() {
  uni.showModal({
    title: 'Deactivate Account',
    content: 'Warning: This action is irreversible. All your data will be permanently deleted. Are you sure you want to proceed?',
    confirmText: 'Deactivate',
    confirmColor: '#e74c3c',
    success: async (res) => {
      if (res.confirm) {
        uni.showModal({
          title: 'Final Confirmation',
          content: 'Please type "DEACTIVATE" in the next step to confirm.',
          confirmText: 'Confirm',
          confirmColor: '#e74c3c',
          success: async (res2) => {
            if (res2.confirm) {
              try {
                await deactivateAccount()
                uni.showToast({ title: 'Account deactivated', icon: 'none' })
                setTimeout(() => {
                  uni.reLaunch({ url: '/pages/login/index' })
                }, 2000)
              } catch (e) {
                console.error('Failed to deactivate', e)
                uni.showToast({ title: 'Deactivation failed', icon: 'none' })
              }
            }
          }
        })
      }
    }
  })
}
</script>

<style scoped>
.settings-page {
  min-height: 100vh;
  background-color: #f5f5f5;
  padding: 20rpx 24rpx;
  padding-bottom: 80rpx;
}

.setting-group {
  background-color: #fff;
  border-radius: 16rpx;
  margin-bottom: 24rpx;
  overflow: hidden;
}

.group-title {
  font-size: 26rpx;
  color: #999;
  padding: 24rpx 30rpx 12rpx;
}

.setting-item {
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: space-between;
  padding: 28rpx 30rpx;
  border-bottom: 1rpx solid #f5f5f5;
}

.setting-item:last-child {
  border-bottom: none;
}

.setting-label {
  font-size: 30rpx;
  color: #333;
}

.setting-value-wrap {
  display: flex;
  flex-direction: row;
  align-items: center;
}

.setting-value {
  font-size: 28rpx;
  color: #999;
  margin-right: 12rpx;
}

.setting-arrow {
  width: 28rpx;
  height: 28rpx;
  opacity: 0.4;
}

.clear-btn {
  padding: 8rpx 24rpx;
  background-color: #f0f0f0;
  border-radius: 20rpx;
  margin-left: 12rpx;
}

.clear-text {
  font-size: 24rpx;
  color: #666;
}

.deactivate-section {
  margin-top: 40rpx;
}

.deactivate-btn {
  background-color: #fff;
  border-radius: 16rpx;
  padding: 28rpx 0;
  display: flex;
  align-items: center;
  justify-content: center;
}

.deactivate-text {
  font-size: 30rpx;
  color: #e74c3c;
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
</style>
