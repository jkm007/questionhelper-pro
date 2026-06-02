<template>
  <view class="exam-publish-page">
    <view class="form-section">
      <!-- Exam Title -->
      <view class="form-item">
        <text class="form-label required">考试名称</text>
        <view class="input-wrap">
          <input
            class="form-input"
            v-model="form.title"
            placeholder="请输入考试名称"
            :maxlength="100"
          />
        </view>
        <text class="char-count">{{ form.title.length }}/100</text>
      </view>

      <!-- Description -->
      <view class="form-item">
        <text class="form-label">考试说明</text>
        <view class="textarea-wrap">
          <textarea
            class="form-textarea"
            v-model="form.description"
            placeholder="请输入考试说明及注意事项（选填）"
            :maxlength="500"
            auto-height
          />
        </view>
        <text class="char-count">{{ form.description.length }}/500</text>
      </view>

      <!-- Time Limit -->
      <view class="form-item">
        <text class="form-label required">考试时长</text>
        <view class="input-wrap input-row">
          <input
            class="form-input"
            type="number"
            v-model="form.duration"
            placeholder="请输入考试时长"
          />
          <text class="input-suffix">分钟</text>
        </view>
      </view>

      <!-- Start Time -->
      <view class="form-item">
        <text class="form-label required">开始时间</text>
        <picker mode="date" :value="startDate" :start="today" @change="onStartDateChange">
          <view class="picker-row">
            <text :class="['picker-value', startDate ? '' : 'placeholder']">
              {{ startDate || '请选择开始日期' }}
            </text>
            <text class="picker-arrow">></text>
          </view>
        </picker>
        <picker v-if="startDate" mode="time" :value="startTime" @change="onStartTimeChange">
          <view class="picker-row">
            <text :class="['picker-value', startTime ? '' : 'placeholder']">
              {{ startTime || '请选择开始时间' }}
            </text>
            <text class="picker-arrow">></text>
          </view>
        </picker>
      </view>

      <!-- End Time -->
      <view class="form-item">
        <text class="form-label required">结束时间</text>
        <picker mode="date" :value="endDate" :start="startDate || today" @change="onEndDateChange">
          <view class="picker-row">
            <text :class="['picker-value', endDate ? '' : 'placeholder']">
              {{ endDate || '请选择结束日期' }}
            </text>
            <text class="picker-arrow">></text>
          </view>
        </picker>
        <picker v-if="endDate" mode="time" :value="endTime" @change="onEndTimeChange">
          <view class="picker-row">
            <text :class="['picker-value', endTime ? '' : 'placeholder']">
              {{ endTime || '请选择结束时间' }}
            </text>
            <text class="picker-arrow">></text>
          </view>
        </picker>
      </view>
    </view>

    <!-- Select Paper -->
    <view class="section-header">
      <text class="section-title required-label">选择试卷</text>
      <view class="section-action" @tap="onSelectPaper">
        <text class="section-action-text">{{ selectedPaper ? '重新选择' : '从试卷库选择' }}</text>
      </view>
    </view>

    <view v-if="selectedPaper" class="paper-card">
      <view class="paper-info">
        <text class="paper-name">{{ selectedPaper.title }}</text>
        <view class="paper-meta">
          <text class="paper-meta-item">{{ selectedPaper.questionCount }} 题</text>
          <text class="paper-meta-divider">|</text>
          <text class="paper-meta-item">总分 {{ selectedPaper.totalScore }} 分</text>
        </view>
      </view>
      <text class="paper-change" @tap="onSelectPaper">更换</text>
    </view>
    <view v-else class="empty-paper">
      <text class="empty-hint">请选择一份试卷</text>
    </view>

    <!-- Passing Score -->
    <view class="form-section" style="margin-top: 0">
      <view class="form-item">
        <text class="form-label required">及格分数</text>
        <view class="input-wrap input-row">
          <input
            class="form-input"
            type="number"
            v-model="form.passingScore"
            placeholder="请输入及格分数"
          />
          <text class="input-suffix">分</text>
        </view>
        <text v-if="selectedPaper" class="form-hint">
          总分 {{ selectedPaper.totalScore }} 分
        </text>
      </view>
    </view>

    <!-- Select Target Classes -->
    <view class="section-header">
      <text class="section-title required-label">指定班级</text>
      <view class="section-action" @tap="onSelectClasses">
        <text class="section-action-text">
          {{ selectedClasses.length > 0 ? '重新选择' : '选择班级' }}
        </text>
      </view>
    </view>

    <view v-if="selectedClasses.length > 0" class="class-tags">
      <view
        v-for="cls in selectedClasses"
        :key="cls.id"
        class="class-tag"
      >
        <text class="class-tag-name">{{ cls.name }}</text>
        <text class="class-tag-remove" @tap="onRemoveClass(cls.id)">x</text>
      </view>
    </view>
    <view v-else class="empty-classes">
      <text class="empty-hint">请选择参与考试的班级</text>
    </view>

    <!-- Submit Button -->
    <view class="submit-section">
      <view :class="['submit-btn', submitting ? 'disabled' : '']" @tap="onSubmit">
        <text class="submit-btn-text">{{ submitting ? '发布中...' : '发布考试' }}</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { request } from '@/api/request'
import { getMyClasses } from '@/api/class'
import { getPapers } from '@/api/exam'

interface PaperItem {
  id: string
  title: string
  questionCount: number
  totalScore: number
}

interface ClassItem {
  id: string
  name: string
}

const submitting = ref(false)
const today = ref('')

const form = reactive({
  title: '',
  description: '',
  duration: '',
  passingScore: '',
})

const startDate = ref('')
const startTime = ref('')
const endDate = ref('')
const endTime = ref('')

const selectedPaper = ref<PaperItem | null>(null)
const selectedClasses = ref<ClassItem[]>([])

onMounted(() => {
  const now = new Date()
  today.value = `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}-${String(now.getDate()).padStart(2, '0')}`

  // Listen for paper selection result
  uni.$on('examPaperSelected', (paper: PaperItem) => {
    selectedPaper.value = paper
  })

  // Listen for class selection result
  uni.$on('examClassSelected', (classes: ClassItem[]) => {
    selectedClasses.value = classes
  })
})

// ---------- Date / Time Pickers ----------

function onStartDateChange(e: any) {
  startDate.value = e.detail.value
  // Reset end date if it is before start date
  if (endDate.value && endDate.value < startDate.value) {
    endDate.value = ''
    endTime.value = ''
  }
}

function onStartTimeChange(e: any) {
  startTime.value = e.detail.value
}

function onEndDateChange(e: any) {
  endDate.value = e.detail.value
}

function onEndTimeChange(e: any) {
  endTime.value = e.detail.value
}

// ---------- Paper Selection ----------

function onSelectPaper() {
  // Navigate to paper list in select mode
  uni.navigateTo({
    url: '/pages/exam/paper-list?mode=select',
  })
}

// ---------- Class Selection ----------

function onSelectClasses() {
  // Navigate to class list in select mode
  uni.navigateTo({
    url: '/pages/class/list?mode=select&selected=' + selectedClasses.value.map((c) => c.id).join(','),
  })
}

function onRemoveClass(id: string) {
  selectedClasses.value = selectedClasses.value.filter((c) => c.id !== id)
}

// ---------- Validation ----------

function validate(): boolean {
  if (!form.title.trim()) {
    uni.showToast({ title: '请输入考试名称', icon: 'none' })
    return false
  }
  if (!form.duration || Number(form.duration) <= 0) {
    uni.showToast({ title: '请输入有效的考试时长', icon: 'none' })
    return false
  }
  if (!startDate.value || !startTime.value) {
    uni.showToast({ title: '请选择开始时间', icon: 'none' })
    return false
  }
  if (!endDate.value || !endTime.value) {
    uni.showToast({ title: '请选择结束时间', icon: 'none' })
    return false
  }
  const startTs = new Date(`${startDate.value} ${startTime.value}:00`).getTime()
  const endTs = new Date(`${endDate.value} ${endTime.value}:00`).getTime()
  if (endTs <= startTs) {
    uni.showToast({ title: '结束时间必须晚于开始时间', icon: 'none' })
    return false
  }
  if (!selectedPaper.value) {
    uni.showToast({ title: '请选择试卷', icon: 'none' })
    return false
  }
  if (!form.passingScore || Number(form.passingScore) < 0) {
    uni.showToast({ title: '请输入有效的及格分数', icon: 'none' })
    return false
  }
  if (selectedPaper.value && Number(form.passingScore) > selectedPaper.value.totalScore) {
    uni.showToast({ title: '及格分数不能超过总分', icon: 'none' })
    return false
  }
  if (selectedClasses.value.length === 0) {
    uni.showToast({ title: '请至少选择一个班级', icon: 'none' })
    return false
  }
  return true
}

// ---------- Submit ----------

async function onSubmit() {
  if (submitting.value) return
  if (!validate()) return

  submitting.value = true
  try {
    const startTimeStr = `${startDate.value} ${startTime.value}:00`
    const endTimeStr = `${endDate.value} ${endTime.value}:00`

    await request({
      url: '/exam',
      method: 'POST',
      data: {
        title: form.title.trim(),
        description: form.description.trim() || undefined,
        duration: Number(form.duration),
        startTime: startTimeStr,
        endTime: endTimeStr,
        paperId: selectedPaper.value!.id,
        passingScore: Number(form.passingScore),
        classIds: selectedClasses.value.map((c) => Number(c.id)),
      },
    })

    uni.showToast({ title: '发布成功', icon: 'success' })
    setTimeout(() => {
      uni.navigateBack()
    }, 1500)
  } catch (e) {
    uni.showToast({ title: '发布失败', icon: 'none' })
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.exam-publish-page {
  min-height: 100vh;
  background-color: #f5f6fa;
  padding-bottom: 160rpx;
}

/* ---- Form Section ---- */

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

.input-wrap {
  border: 2rpx solid #e0e0e0;
  border-radius: 12rpx;
  padding: 0 20rpx;
  height: 80rpx;
}

.input-row {
  display: flex;
  align-items: center;
}

.input-row .form-input {
  flex: 1;
}

.input-suffix {
  font-size: 28rpx;
  color: #999;
  margin-left: 12rpx;
  flex-shrink: 0;
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

.picker-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16rpx 0;
  border-bottom: 1rpx solid #f5f5f5;
}

.picker-row:last-child {
  border-bottom: none;
}

.picker-value {
  font-size: 28rpx;
  color: #333;
}

.picker-value.placeholder {
  color: #ccc;
}

.picker-arrow {
  font-size: 24rpx;
  color: #ccc;
}

/* ---- Section Header ---- */

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 24rpx 32rpx 12rpx;
}

.section-title {
  font-size: 30rpx;
  font-weight: 600;
  color: #333;
}

.required-label::before {
  content: '*';
  color: #ff4d4f;
  margin-right: 4rpx;
}

.section-action {
  padding: 8rpx 20rpx;
  background-color: #e8f3ff;
  border-radius: 24rpx;
}

.section-action-text {
  font-size: 24rpx;
  color: #1677ff;
}

/* ---- Paper Card ---- */

.paper-card {
  margin: 0 24rpx;
  background: #ffffff;
  border-radius: 16rpx;
  padding: 24rpx;
  display: flex;
  align-items: center;
  justify-content: space-between;
  box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.04);
}

.paper-info {
  flex: 1;
  min-width: 0;
}

.paper-name {
  font-size: 28rpx;
  font-weight: 600;
  color: #333;
  display: block;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.paper-meta {
  display: flex;
  align-items: center;
  margin-top: 8rpx;
}

.paper-meta-item {
  font-size: 24rpx;
  color: #999;
}

.paper-meta-divider {
  font-size: 24rpx;
  color: #e0e0e0;
  margin: 0 12rpx;
}

.paper-change {
  font-size: 24rpx;
  color: #1677ff;
  padding: 8rpx 16rpx;
  flex-shrink: 0;
}

.empty-paper {
  padding: 40rpx 24rpx;
  text-align: center;
}

.empty-hint {
  font-size: 28rpx;
  color: #ccc;
}

/* ---- Class Tags ---- */

.class-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 12rpx;
  padding: 0 24rpx 20rpx;
}

.class-tag {
  display: flex;
  align-items: center;
  background: #e8f3ff;
  border-radius: 24rpx;
  padding: 10rpx 20rpx;
}

.class-tag-name {
  font-size: 24rpx;
  color: #1677ff;
  margin-right: 8rpx;
}

.class-tag-remove {
  font-size: 24rpx;
  color: #999;
  padding: 0 4rpx;
}

.empty-classes {
  padding: 40rpx 24rpx;
  text-align: center;
}

/* ---- Submit ---- */

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
  background-color: #1677ff;
  border-radius: 44rpx;
  padding: 24rpx 0;
  text-align: center;
}

.submit-btn.disabled {
  opacity: 0.6;
}

.submit-btn-text {
  font-size: 30rpx;
  font-weight: 600;
  color: #fff;
}
</style>
