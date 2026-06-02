<template>
  <view class="publish-page">
    <view class="form-card">
      <!-- Title -->
      <view class="form-item">
        <text class="form-label">考试标题</text>
        <input class="form-input" v-model="form.title" placeholder="请输入考试标题" />
      </view>

      <!-- Description -->
      <view class="form-item">
        <text class="form-label">考试说明</text>
        <textarea class="form-textarea" v-model="form.description" placeholder="请输入考试说明（选填）" />
      </view>

      <!-- Time Limit -->
      <view class="form-item">
        <text class="form-label">考试时长（分钟）</text>
        <input class="form-input" v-model="form.timeLimit" type="number" placeholder="如：120" />
      </view>

      <!-- Passing Score -->
      <view class="form-item">
        <text class="form-label">及格分数</text>
        <input class="form-input" v-model="form.passingScore" type="number" placeholder="如：60" />
      </view>

      <!-- Start Time -->
      <view class="form-item">
        <text class="form-label">开始时间</text>
        <picker mode="multiSelector" :range="dateTimeRange" :value="startDateTimeIndex" @change="onStartChange">
          <text class="picker-display">{{ form.startTime || '请选择开始时间' }}</text>
        </picker>
      </view>

      <!-- End Time -->
      <view class="form-item">
        <text class="form-label">结束时间</text>
        <picker mode="multiSelector" :range="dateTimeRange" :value="endDateTimeIndex" @change="onEndChange">
          <text class="picker-display">{{ form.endTime || '请选择结束时间' }}</text>
        </picker>
      </view>

      <!-- Class Selection -->
      <view class="form-item">
        <text class="form-label">发布班级</text>
        <view class="class-list">
          <view
            v-for="cls in classes"
            :key="cls.id"
            class="class-item"
            :class="{ selected: selectedClasses.includes(cls.id) }"
            @tap="toggleClass(cls.id)"
          >
            <view class="checkbox" :class="{ checked: selectedClasses.includes(cls.id) }">
              <text v-if="selectedClasses.includes(cls.id)" class="check-mark">✓</text>
            </view>
            <text class="class-name">{{ cls.name }}</text>
            <text class="class-count">{{ cls.studentCount }}人</text>
          </view>
        </view>
      </view>
    </view>

    <!-- Submit -->
    <view class="submit-btn" :class="{ disabled: submitting }" @tap="handlePublish">
      <text class="submit-text">{{ submitting ? '发布中...' : '发布考试' }}</text>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'

const submitting = ref(false)
const selectedClasses = ref<number[]>([])

const form = ref({
  title: '',
  description: '',
  timeLimit: '',
  passingScore: '',
  startTime: '',
  endTime: ''
})

const startDateTimeIndex = ref([0, 0])
const endDateTimeIndex = ref([0, 0])

interface ClassItem { id: number; name: string; studentCount: number }
const classes = ref<ClassItem[]>([])

const dateTimeRange = computed(() => {
  const dates: string[] = []
  const now = new Date()
  for (let i = 0; i < 30; i++) {
    const d = new Date(now)
    d.setDate(d.getDate() + i)
    dates.push(`${d.getMonth() + 1}月${d.getDate()}日`)
  }
  const times: string[] = []
  for (let h = 0; h < 24; h++) {
    for (let m = 0; m < 60; m += 30) {
      times.push(`${String(h).padStart(2, '0')}:${String(m).padStart(2, '0')}`)
    }
  }
  return [dates, times]
})

onMounted(() => { fetchClasses() })

async function fetchClasses() {
  try {
    // TODO: replace with actual API call
    classes.value = [
      { id: 1, name: '数据结构2026春季班', studentCount: 45 },
      { id: 2, name: '算法设计2026春季班', studentCount: 38 },
      { id: 3, name: '操作系统2026春季班', studentCount: 52 }
    ]
  } catch (e) {
    console.error('Failed to load classes', e)
  }
}

function onStartChange(e: any) {
  startDateTimeIndex.value = e.detail.value
  const dates = dateTimeRange.value[0]
  const times = dateTimeRange.value[1]
  form.value.startTime = `${dates[e.detail.value[0]]} ${times[e.detail.value[1]]}`
}

function onEndChange(e: any) {
  endDateTimeIndex.value = e.detail.value
  const dates = dateTimeRange.value[0]
  const times = dateTimeRange.value[1]
  form.value.endTime = `${dates[e.detail.value[0]]} ${times[e.detail.value[1]]}`
}

function toggleClass(id: number) {
  const idx = selectedClasses.value.indexOf(id)
  if (idx > -1) {
    selectedClasses.value.splice(idx, 1)
  } else {
    selectedClasses.value.push(id)
  }
}

async function handlePublish() {
  if (!form.value.title.trim()) {
    uni.showToast({ title: '请输入考试标题', icon: 'none' })
    return
  }
  if (!form.value.timeLimit || Number(form.value.timeLimit) <= 0) {
    uni.showToast({ title: '请输入有效的考试时长', icon: 'none' })
    return
  }
  if (!form.value.passingScore) {
    uni.showToast({ title: '请输入及格分数', icon: 'none' })
    return
  }
  if (!form.value.startTime || !form.value.endTime) {
    uni.showToast({ title: '请选择考试时间', icon: 'none' })
    return
  }
  if (selectedClasses.value.length === 0) {
    uni.showToast({ title: '请选择发布班级', icon: 'none' })
    return
  }
  submitting.value = true
  try {
    // TODO: replace with actual API call
    uni.showToast({ title: '发布成功', icon: 'success' })
    setTimeout(() => uni.navigateBack(), 1500)
  } catch (e) {
    console.error('Failed to publish exam', e)
    uni.showToast({ title: '发布失败', icon: 'none' })
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.publish-page {
  min-height: 100vh;
  background-color: #f5f5f5;
  padding: 20rpx 24rpx;
  padding-bottom: 80rpx;
}

.form-card {
  background-color: #fff;
  border-radius: 16rpx;
  padding: 12rpx 30rpx;
  margin-bottom: 40rpx;
}

.form-item {
  padding: 28rpx 0;
  border-bottom: 1rpx solid #f5f5f5;
}

.form-item:last-child {
  border-bottom: none;
}

.form-label {
  font-size: 26rpx;
  color: #999;
  margin-bottom: 16rpx;
}

.form-input {
  font-size: 30rpx;
  color: #333;
  padding: 8rpx 0;
}

.form-textarea {
  width: 100%;
  height: 160rpx;
  font-size: 28rpx;
  color: #333;
  padding: 16rpx;
  background-color: #f9f9f9;
  border-radius: 12rpx;
}

.picker-display {
  font-size: 28rpx;
  color: #666;
  padding: 8rpx 0;
}

.class-list {
  display: flex;
  flex-direction: column;
  gap: 16rpx;
  margin-top: 8rpx;
}

.class-item {
  display: flex;
  flex-direction: row;
  align-items: center;
  padding: 20rpx 24rpx;
  background-color: #f9f9f9;
  border-radius: 12rpx;
  border: 2rpx solid transparent;
}

.class-item.selected {
  border-color: #4a90d9;
  background-color: #f0f7ff;
}

.checkbox {
  width: 36rpx;
  height: 36rpx;
  border: 2rpx solid #ddd;
  border-radius: 8rpx;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-right: 16rpx;
}

.checkbox.checked {
  background-color: #4a90d9;
  border-color: #4a90d9;
}

.check-mark {
  font-size: 22rpx;
  color: #fff;
  font-weight: 700;
}

.class-name {
  flex: 1;
  font-size: 28rpx;
  color: #333;
}

.class-count {
  font-size: 24rpx;
  color: #bbb;
}

.submit-btn {
  background-color: #4a90d9;
  border-radius: 44rpx;
  padding: 26rpx 0;
  display: flex;
  align-items: center;
  justify-content: center;
}

.submit-btn.disabled {
  opacity: 0.6;
}

.submit-text {
  font-size: 30rpx;
  color: #fff;
  font-weight: 600;
}
</style>
