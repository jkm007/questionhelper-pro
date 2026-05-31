<template>
  <view class="page">
    <view class="form-section">
      <!-- Cover Image -->
      <view class="form-item">
        <text class="form-label">班级封面</text>
        <view class="cover-upload" @tap="onChooseCover">
          <image v-if="form.coverImage" class="cover-preview" :src="form.coverImage" mode="aspectFill" />
          <view v-else class="cover-placeholder">
            <text class="upload-icon">+</text>
            <text class="upload-text">上传封面</text>
          </view>
        </view>
      </view>

      <!-- Name -->
      <view class="form-item">
        <text class="form-label required">班级名称</text>
        <view class="input-wrap">
          <input
            class="form-input"
            v-model="form.name"
            placeholder="请输入班级名称"
            :maxlength="30"
          />
        </view>
        <text class="char-count">{{ form.name.length }}/30</text>
      </view>

      <!-- Description -->
      <view class="form-item">
        <text class="form-label">班级描述</text>
        <view class="textarea-wrap">
          <textarea
            class="form-textarea"
            v-model="form.description"
            placeholder="请输入班级描述（选填）"
            :maxlength="200"
            auto-height
          />
        </view>
        <text class="char-count">{{ form.description.length }}/200</text>
      </view>

      <!-- Max Members -->
      <view class="form-item">
        <text class="form-label">最大成员数</text>
        <view class="input-wrap">
          <input
            class="form-input"
            v-model="form.maxMembers"
            type="number"
            placeholder="请输入最大成员数量"
          />
        </view>
        <text class="form-hint">范围 10-1000，不填默认200</text>
      </view>

      <!-- Settings -->
      <view class="form-item settings-item">
        <text class="form-label">班级设置</text>

        <view class="setting-row">
          <view class="setting-info">
            <text class="setting-title">允许搜索加入</text>
            <text class="setting-desc">开启后其他用户可通过搜索发现并加入班级</text>
          </view>
          <switch
            :checked="form.allowJoin"
            @change="form.allowJoin = $event.detail.value"
            color="#1677ff"
          />
        </view>

        <view class="setting-row">
          <view class="setting-info">
            <text class="setting-title">需要审批</text>
            <text class="setting-desc">开启后新成员加入需要管理员审批</text>
          </view>
          <switch
            :checked="form.requireApproval"
            @change="form.requireApproval = $event.detail.value"
            color="#1677ff"
          />
        </view>
      </view>
    </view>

    <!-- Submit Button -->
    <view class="submit-section">
      <text :class="['submit-btn', submitting ? 'disabled' : '']" @tap="onSubmit">
        {{ submitting ? '创建中...' : '创建班级' }}
      </text>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { createClass } from '@/api/class'

const form = reactive({
  name: '',
  description: '',
  coverImage: '',
  maxMembers: '',
  allowJoin: true,
  requireApproval: false
})

const submitting = ref(false)

function onChooseCover() {
  uni.chooseImage({
    count: 1,
    sizeType: ['compressed'],
    sourceType: ['album', 'camera'],
    success: (res) => {
      form.coverImage = res.tempFilePaths[0]
    }
  })
}

function validate(): boolean {
  if (!form.name.trim()) {
    uni.showToast({ title: '请输入班级名称', icon: 'none' })
    return false
  }
  if (form.name.trim().length < 2) {
    uni.showToast({ title: '班级名称至少2个字符', icon: 'none' })
    return false
  }
  if (form.maxMembers) {
    const num = parseInt(form.maxMembers)
    if (isNaN(num) || num < 10 || num > 1000) {
      uni.showToast({ title: '最大成员数范围10-1000', icon: 'none' })
      return false
    }
  }
  return true
}

async function onSubmit() {
  if (submitting.value) return
  if (!validate()) return

  submitting.value = true
  try {
    // Upload cover if local file
    let coverUrl = form.coverImage
    if (form.coverImage && !form.coverImage.startsWith('http')) {
      const uploadRes = await uni.uploadFile({
        url: '/api/upload/image',
        filePath: form.coverImage,
        name: 'file'
      })
      const data = JSON.parse(uploadRes.data)
      coverUrl = data.data?.url || ''
    }

    await createClass({
      name: form.name.trim(),
      description: form.description.trim(),
      coverImage: coverUrl,
      maxMembers: form.maxMembers ? parseInt(form.maxMembers) : 200,
      allowJoin: form.allowJoin,
      requireApproval: form.requireApproval
    })

    uni.showToast({ title: '创建成功', icon: 'success' })
    setTimeout(() => {
      uni.navigateBack()
    }, 1500)
  } catch (e) {
    uni.showToast({ title: '创建失败', icon: 'none' })
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.page {
  min-height: 100vh;
  background-color: #f5f6fa;
  padding-bottom: 140rpx;
}

.form-section {
  background-color: #fff;
  margin: 20rpx 24rpx;
  border-radius: 16rpx;
  overflow: hidden;
}

.form-item {
  padding: 24rpx 28rpx;
  border-bottom: 1rpx solid #f5f5f5;
}

.form-item:last-child {
  border-bottom: none;
}

.form-label {
  font-size: 28rpx;
  font-weight: 600;
  color: #333;
  margin-bottom: 16rpx;
  display: block;
}

.form-label.required::before {
  content: '*';
  color: #ff4d4f;
  margin-right: 4rpx;
}

.cover-upload {
  width: 320rpx;
  height: 200rpx;
  border-radius: 12rpx;
  overflow: hidden;
  border: 2rpx dashed #d9d9d9;
}

.cover-preview {
  width: 100%;
  height: 100%;
}

.cover-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background-color: #fafafa;
}

.upload-icon {
  font-size: 56rpx;
  color: #d9d9d9;
  line-height: 1;
  margin-bottom: 8rpx;
}

.upload-text {
  font-size: 24rpx;
  color: #999;
}

.input-wrap {
  border: 2rpx solid #e0e0e0;
  border-radius: 12rpx;
  padding: 0 20rpx;
  height: 80rpx;
}

.form-input {
  width: 100%;
  height: 80rpx;
  font-size: 28rpx;
}

.textarea-wrap {
  border: 2rpx solid #e0e0e0;
  border-radius: 12rpx;
  padding: 16rpx 20rpx;
  min-height: 160rpx;
}

.form-textarea {
  width: 100%;
  font-size: 28rpx;
  line-height: 1.6;
}

.char-count {
  font-size: 22rpx;
  color: #ccc;
  text-align: right;
  margin-top: 8rpx;
  display: block;
}

.form-hint {
  font-size: 22rpx;
  color: #999;
  margin-top: 8rpx;
  display: block;
}

.settings-item {
  padding-bottom: 12rpx;
}

.setting-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16rpx 0;
}

.setting-info {
  flex: 1;
  margin-right: 20rpx;
}

.setting-title {
  font-size: 28rpx;
  color: #333;
  display: block;
  margin-bottom: 4rpx;
}

.setting-desc {
  font-size: 22rpx;
  color: #999;
  line-height: 1.4;
}

.submit-section {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  padding: 20rpx 32rpx;
  padding-bottom: calc(20rpx + env(safe-area-inset-bottom));
  background-color: #fff;
  box-shadow: 0 -2rpx 12rpx rgba(0, 0, 0, 0.06);
}

.submit-btn {
  display: block;
  text-align: center;
  font-size: 30rpx;
  font-weight: 600;
  color: #fff;
  background-color: #1677ff;
  padding: 24rpx 0;
  border-radius: 44rpx;
}

.submit-btn.disabled {
  opacity: 0.6;
}
</style>
