<template>
  <view class="plan-page">
    <!-- Plan List -->
    <view v-if="plans.length" class="plan-list">
      <uni-swipe-action>
        <uni-swipe-action-item v-for="plan in plans" :key="plan.id" :right-options="swipeOptions" @click="onSwipeClick($event, plan.id)">
          <view class="plan-card">
            <view class="plan-header">
              <text class="plan-name">{{ plan.name }}</text>
              <text class="plan-percent">{{ plan.progress }}%</text>
            </view>
            <text class="plan-target">目标：{{ plan.target }}</text>
            <view class="progress-bar">
              <view class="progress-fill" :style="{ width: plan.progress + '%' }" />
            </view>
            <view class="plan-footer">
              <text class="plan-date">{{ plan.startDate }} ~ {{ plan.endDate }}</text>
            </view>
          </view>
        </uni-swipe-action-item>
      </uni-swipe-action>
    </view>

    <!-- Empty State -->
    <view v-else class="empty-state">
      <text class="empty-text">暂无学习计划</text>
      <text class="empty-hint">点击下方按钮创建你的第一个计划</text>
    </view>

    <!-- Create Button -->
    <view class="create-btn" @tap="showCreateDialog">
      <text class="create-text">+ 新建计划</text>
    </view>

    <!-- Create Dialog -->
    <uni-popup ref="createPopup" type="center">
      <view class="dialog">
        <text class="dialog-title">新建学习计划</text>
        <input class="dialog-input" v-model="newPlan.name" placeholder="计划名称" />
        <input class="dialog-input" v-model="newPlan.target" placeholder="学习目标" />
        <view class="dialog-row">
          <picker mode="date" @change="newPlan.startDate = $event.detail.value">
            <text class="date-picker">{{ newPlan.startDate || '开始日期' }}</text>
          </picker>
          <text class="date-sep">~</text>
          <picker mode="date" @change="newPlan.endDate = $event.detail.value">
            <text class="date-picker">{{ newPlan.endDate || '结束日期' }}</text>
          </picker>
        </view>
        <view class="dialog-actions">
          <view class="dialog-btn cancel" @tap="createPopup?.close()">
            <text class="btn-text cancel-text">取消</text>
          </view>
          <view class="dialog-btn confirm" @tap="createPlan">
            <text class="btn-text confirm-text">确定</text>
          </view>
        </view>
      </view>
    </uni-popup>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

const createPopup = ref<any>(null)

interface Plan {
  id: number
  name: string
  target: string
  progress: number
  startDate: string
  endDate: string
}

const plans = ref<Plan[]>([])

const swipeOptions = [{ text: '删除', style: { backgroundColor: '#e74c3c' } }]

const newPlan = ref({ name: '', target: '', startDate: '', endDate: '' })

onMounted(() => { fetchPlans() })

async function fetchPlans() {
  try {
    // TODO: replace with actual API call
    plans.value = [
      { id: 1, name: '算法基础训练', target: '完成200道LeetCode', progress: 65, startDate: '2026-05-01', endDate: '2026-06-30' },
      { id: 2, name: '操作系统复习', target: '通过期末考试', progress: 30, startDate: '2026-05-15', endDate: '2026-06-20' }
    ]
  } catch (e) {
    console.error('Failed to load plans', e)
  }
}

function showCreateDialog() {
  newPlan.value = { name: '', target: '', startDate: '', endDate: '' }
  createPopup.value?.open()
}

async function createPlan() {
  if (!newPlan.value.name.trim()) {
    uni.showToast({ title: '请输入计划名称', icon: 'none' })
    return
  }
  try {
    // TODO: replace with actual API call
    plans.value.push({
      id: Date.now(),
      ...newPlan.value,
      progress: 0
    })
    createPopup.value?.close()
    uni.showToast({ title: '创建成功', icon: 'success' })
  } catch (e) {
    console.error('Failed to create plan', e)
  }
}

function onSwipeClick(e: any, id: number) {
  if (e.index === 0) {
    uni.showModal({
      title: '确认删除',
      content: '确定要删除该学习计划吗？',
      success: async (res) => {
        if (res.confirm) {
          try {
            // TODO: replace with actual API call
            plans.value = plans.value.filter((p) => p.id !== id)
            uni.showToast({ title: '已删除', icon: 'success' })
          } catch (e) {
            console.error('Failed to delete plan', e)
          }
        }
      }
    })
  }
}
</script>

<style scoped>
.plan-page {
  min-height: 100vh;
  background-color: #f5f5f5;
  padding: 20rpx 24rpx;
  padding-bottom: 140rpx;
}

.plan-card {
  background-color: #fff;
  border-radius: 16rpx;
  padding: 28rpx 30rpx;
  margin-bottom: 20rpx;
}

.plan-header {
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10rpx;
}

.plan-name {
  font-size: 32rpx;
  font-weight: 600;
  color: #333;
}

.plan-percent {
  font-size: 32rpx;
  font-weight: 700;
  color: #4a90d9;
}

.plan-target {
  font-size: 26rpx;
  color: #666;
  margin-bottom: 16rpx;
}

.progress-bar {
  height: 14rpx;
  background-color: #f0f0f0;
  border-radius: 7rpx;
  overflow: hidden;
  margin-bottom: 12rpx;
}

.progress-fill {
  height: 100%;
  background-color: #4a90d9;
  border-radius: 7rpx;
  transition: width 0.3s;
}

.plan-footer {
  display: flex;
  justify-content: flex-end;
}

.plan-date {
  font-size: 22rpx;
  color: #bbb;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 120rpx 0;
}

.empty-text {
  font-size: 30rpx;
  color: #999;
  margin-bottom: 12rpx;
}

.empty-hint {
  font-size: 24rpx;
  color: #ccc;
}

.create-btn {
  position: fixed;
  bottom: 40rpx;
  left: 48rpx;
  right: 48rpx;
  background-color: #4a90d9;
  border-radius: 44rpx;
  padding: 26rpx 0;
  display: flex;
  align-items: center;
  justify-content: center;
}

.create-text {
  font-size: 30rpx;
  color: #fff;
  font-weight: 600;
}

.dialog {
  width: 600rpx;
  background-color: #fff;
  border-radius: 24rpx;
  padding: 40rpx 36rpx;
}

.dialog-title {
  font-size: 32rpx;
  font-weight: 600;
  color: #333;
  margin-bottom: 30rpx;
  text-align: center;
}

.dialog-input {
  border: 1rpx solid #e0e0e0;
  border-radius: 12rpx;
  padding: 20rpx 24rpx;
  font-size: 28rpx;
  margin-bottom: 20rpx;
}

.dialog-row {
  display: flex;
  flex-direction: row;
  align-items: center;
  gap: 16rpx;
  margin-bottom: 30rpx;
}

.date-picker {
  flex: 1;
  border: 1rpx solid #e0e0e0;
  border-radius: 12rpx;
  padding: 20rpx 24rpx;
  font-size: 26rpx;
  color: #666;
  text-align: center;
}

.date-sep {
  font-size: 26rpx;
  color: #999;
}

.dialog-actions {
  display: flex;
  flex-direction: row;
  gap: 20rpx;
}

.dialog-btn {
  flex: 1;
  border-radius: 12rpx;
  padding: 22rpx 0;
  display: flex;
  align-items: center;
  justify-content: center;
}

.dialog-btn.cancel {
  background-color: #f0f0f0;
}

.dialog-btn.confirm {
  background-color: #4a90d9;
}

.btn-text {
  font-size: 28rpx;
}

.cancel-text {
  color: #666;
}

.confirm-text {
  color: #fff;
}
</style>
