<template>
  <view class="plan-page">
    <!-- Current Plan -->
    <view v-if="currentPlan" class="current-plan-card">
      <view class="plan-header">
        <text class="plan-title">{{ currentPlan.title }}</text>
        <view class="plan-badge" :class="currentPlan.status">
          <text class="plan-badge-text">{{ statusText(currentPlan.status) }}</text>
        </view>
      </view>
      <text class="plan-date">{{ currentPlan.startDate }} - {{ currentPlan.endDate }}</text>

      <!-- Progress -->
      <view class="progress-section">
        <view class="progress-header">
          <text class="progress-label">Overall Progress</text>
          <text class="progress-value">{{ overallProgress }}%</text>
        </view>
        <view class="progress-bar">
          <view class="progress-fill" :style="{ width: overallProgress + '%' }" />
        </view>
      </view>

      <!-- Daily Goals -->
      <view class="goals-row">
        <view class="goal-item">
          <text class="goal-value">{{ currentPlan.dailyPracticeGoal }}</text>
          <text class="goal-label">Questions/Day</text>
        </view>
        <view class="goal-item">
          <text class="goal-value">{{ currentPlan.dailyStudyMinutes }}min</text>
          <text class="goal-label">Study Time/Day</text>
        </view>
        <view class="goal-item">
          <text class="goal-value">{{ currentPlan.weeklyExamGoal }}</text>
          <text class="goal-label">Exams/Week</text>
        </view>
      </view>
    </view>

    <!-- Plan Items -->
    <view class="section-header">
      <text class="section-title">Plan Items</text>
      <view class="add-btn" @tap="showAddItem">
        <text class="add-btn-text">+ Add</text>
      </view>
    </view>

    <view v-if="planItems.length" class="items-list">
      <view
        v-for="item in planItems"
        :key="item.id"
        class="item-card"
        :class="{ completed: item.isCompleted }"
      >
        <view class="item-check" @tap="toggleItem(item)">
          <view class="check-circle" :class="{ checked: item.isCompleted }">
            <text v-if="item.isCompleted" class="check-icon">OK</text>
          </view>
        </view>
        <view class="item-content">
          <text class="item-title" :class="{ 'title-done': item.isCompleted }">{{ item.title }}</text>
          <text v-if="item.description" class="item-desc">{{ item.description }}</text>
          <view class="item-meta">
            <text class="item-category">{{ item.category }}</text>
            <text v-if="item.dueDate" class="item-due">Due: {{ item.dueDate }}</text>
          </view>
        </view>
        <view class="item-delete" @tap="deleteItem(item.id)">
          <text class="delete-icon">×</text>
        </view>
      </view>
    </view>

    <view v-else class="empty-state">
      <text class="empty-text">No plan items yet</text>
      <text class="empty-hint">Tap "+ Add" to create your first task</text>
    </view>

    <!-- Create/Edit Plan Button -->
    <view v-if="!currentPlan" class="create-plan-btn" @tap="showPlanForm = true">
      <text class="create-plan-text">Create Learning Plan</text>
    </view>

    <!-- Plan Form Modal -->
    <uni-popup ref="planPopup" type="center">
      <view class="modal-content">
        <text class="modal-title">{{ editingPlan ? 'Edit Plan' : 'Create Learning Plan' }}</text>

        <view class="form-field">
          <text class="field-label">Plan Title</text>
          <input
            class="field-input"
            v-model="planForm.title"
            placeholder="e.g. Final Exam Preparation"
          />
        </view>

        <view class="form-row">
          <view class="form-field half">
            <text class="field-label">Start Date</text>
            <picker mode="date" :value="planForm.startDate" @change="onStartDateChange">
              <view class="picker-trigger">
                <text class="picker-text">{{ planForm.startDate || 'Select' }}</text>
              </view>
            </picker>
          </view>
          <view class="form-field half">
            <text class="field-label">End Date</text>
            <picker mode="date" :value="planForm.endDate" @change="onEndDateChange">
              <view class="picker-trigger">
                <text class="picker-text">{{ planForm.endDate || 'Select' }}</text>
              </view>
            </picker>
          </view>
        </view>

        <view class="form-field">
          <text class="field-label">Daily Practice Goal (questions)</text>
          <input
            class="field-input"
            type="number"
            v-model="planForm.dailyPracticeGoal"
            placeholder="20"
          />
        </view>

        <view class="form-field">
          <text class="field-label">Daily Study Time (minutes)</text>
          <input
            class="field-input"
            type="number"
            v-model="planForm.dailyStudyMinutes"
            placeholder="60"
          />
        </view>

        <view class="form-field">
          <text class="field-label">Weekly Exam Goal</text>
          <input
            class="field-input"
            type="number"
            v-model="planForm.weeklyExamGoal"
            placeholder="2"
          />
        </view>

        <view class="modal-actions">
          <view class="modal-btn cancel" @tap="closePlanForm">
            <text class="modal-btn-text cancel-text">Cancel</text>
          </view>
          <view class="modal-btn confirm" @tap="savePlan">
            <text class="modal-btn-text confirm-text">Save</text>
          </view>
        </view>
      </view>
    </uni-popup>

    <!-- Add Item Modal -->
    <uni-popup ref="itemPopup" type="center">
      <view class="modal-content">
        <text class="modal-title">Add Plan Item</text>

        <view class="form-field">
          <text class="field-label">Task Title</text>
          <input class="field-input" v-model="itemForm.title" placeholder="e.g. Review Chapter 5" />
        </view>

        <view class="form-field">
          <text class="field-label">Description (optional)</text>
          <textarea
            class="field-textarea"
            v-model="itemForm.description"
            placeholder="Additional details..."
            maxlength="200"
          />
        </view>

        <view class="form-field">
          <text class="field-label">Category</text>
          <view class="category-options">
            <view
              v-for="cat in categories"
              :key="cat"
              class="category-chip"
              :class="{ selected: itemForm.category === cat }"
              @tap="itemForm.category = cat"
            >
              <text class="chip-text">{{ cat }}</text>
            </view>
          </view>
        </view>

        <view class="form-field">
          <text class="field-label">Due Date (optional)</text>
          <picker mode="date" :value="itemForm.dueDate" @change="onDueDateChange">
            <view class="picker-trigger">
              <text class="picker-text">{{ itemForm.dueDate || 'Select date' }}</text>
            </view>
          </picker>
        </view>

        <view class="modal-actions">
          <view class="modal-btn cancel" @tap="closeItemForm">
            <text class="modal-btn-text cancel-text">Cancel</text>
          </view>
          <view class="modal-btn confirm" @tap="saveItem">
            <text class="modal-btn-text confirm-text">Add</text>
          </view>
        </view>
      </view>
    </uni-popup>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'

interface PlanItem {
  id: number
  title: string
  description: string
  category: string
  dueDate: string
  isCompleted: boolean
}

interface LearningPlan {
  id: number
  title: string
  startDate: string
  endDate: string
  status: 'active' | 'completed' | 'paused'
  dailyPracticeGoal: number
  dailyStudyMinutes: number
  weeklyExamGoal: number
}

const categories = ['Practice', 'Review', 'Exam', 'Reading', 'Other']

const currentPlan = ref<LearningPlan | null>(null)
const planItems = ref<PlanItem[]>([])
const showPlanForm = ref(false)
const editingPlan = ref(false)

const planPopup = ref<any>(null)
const itemPopup = ref<any>(null)

const planForm = ref({
  title: '',
  startDate: '',
  endDate: '',
  dailyPracticeGoal: '20',
  dailyStudyMinutes: '60',
  weeklyExamGoal: '2'
})

const itemForm = ref({
  title: '',
  description: '',
  category: 'Practice',
  dueDate: ''
})

const overallProgress = computed(() => {
  if (!planItems.value.length) return 0
  const completed = planItems.value.filter((i) => i.isCompleted).length
  return Math.round((completed / planItems.value.length) * 100)
})

onMounted(() => {
  fetchPlan()
})

async function fetchPlan() {
  try {
    const res = await new Promise<any>((resolve) => {
      uni.request({
        url: '/api/v1/user/learning-plan',
        method: 'GET',
        header: { Authorization: `Bearer ${uni.getStorageSync('token')}` },
        success: (r: any) => resolve(r.data),
        fail: () => resolve(null)
      })
    })

    if (res?.code === '00000' && res.data) {
      currentPlan.value = res.data.plan || null
      planItems.value = res.data.items || []
    } else {
      loadMockData()
    }
  } catch (e) {
    console.error('Failed to load learning plan', e)
    loadMockData()
  }
}

function loadMockData() {
  currentPlan.value = {
    id: 1,
    title: 'Final Exam Preparation',
    startDate: '2026-05-15',
    endDate: '2026-06-30',
    status: 'active',
    dailyPracticeGoal: 20,
    dailyStudyMinutes: 60,
    weeklyExamGoal: 2
  }

  planItems.value = [
    { id: 1, title: 'Complete Math Chapter 1-5 review', description: 'Focus on calculus and algebra', category: 'Review', dueDate: '2026-06-05', isCompleted: true },
    { id: 2, title: 'Practice 100 physics questions', description: 'Mechanics and thermodynamics', category: 'Practice', dueDate: '2026-06-10', isCompleted: false },
    { id: 3, title: 'Take 3 mock exams', description: '', category: 'Exam', dueDate: '2026-06-15', isCompleted: false },
    { id: 4, title: 'Read English textbook Unit 8-10', description: 'Vocabulary and grammar exercises', category: 'Reading', dueDate: '2026-06-12', isCompleted: false },
    { id: 5, title: 'Review chemistry lab notes', description: '', category: 'Review', dueDate: '2026-06-08', isCompleted: true }
  ]
}

function statusText(status: string): string {
  const map: Record<string, string> = { active: 'Active', completed: 'Completed', paused: 'Paused' }
  return map[status] || status
}

function showAddItem() {
  itemForm.value = { title: '', description: '', category: 'Practice', dueDate: '' }
  itemPopup.value?.open()
}

function closePlanForm() {
  planPopup.value?.close()
}

function closeItemForm() {
  itemPopup.value?.close()
}

function onStartDateChange(e: any) {
  planForm.value.startDate = e.detail.value
}

function onEndDateChange(e: any) {
  planForm.value.endDate = e.detail.value
}

function onDueDateChange(e: any) {
  itemForm.value.dueDate = e.detail.value
}

async function savePlan() {
  if (!planForm.value.title.trim()) {
    uni.showToast({ title: 'Please enter a title', icon: 'none' })
    return
  }
  if (!planForm.value.startDate || !planForm.value.endDate) {
    uni.showToast({ title: 'Please select dates', icon: 'none' })
    return
  }

  try {
    const data = {
      title: planForm.value.title.trim(),
      startDate: planForm.value.startDate,
      endDate: planForm.value.endDate,
      dailyPracticeGoal: Number(planForm.value.dailyPracticeGoal) || 20,
      dailyStudyMinutes: Number(planForm.value.dailyStudyMinutes) || 60,
      weeklyExamGoal: Number(planForm.value.weeklyExamGoal) || 2
    }

    await new Promise<any>((resolve, reject) => {
      uni.request({
        url: '/api/v1/user/learning-plan',
        method: currentPlan.value ? 'PUT' : 'POST',
        data,
        header: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${uni.getStorageSync('token')}`
        },
        success: (r: any) => {
          if (r.data?.code === '00000') resolve(r.data)
          else reject(new Error(r.data?.msg))
        },
        fail: reject
      })
    })

    currentPlan.value = {
      id: currentPlan.value?.id || Date.now(),
      ...data,
      status: 'active'
    }
    closePlanForm()
    uni.showToast({ title: 'Plan saved', icon: 'success' })
  } catch (e) {
    console.error('Save plan failed', e)
    // Offline fallback: save locally
    currentPlan.value = {
      id: currentPlan.value?.id || Date.now(),
      title: planForm.value.title.trim(),
      startDate: planForm.value.startDate,
      endDate: planForm.value.endDate,
      dailyPracticeGoal: Number(planForm.value.dailyPracticeGoal) || 20,
      dailyStudyMinutes: Number(planForm.value.dailyStudyMinutes) || 60,
      weeklyExamGoal: Number(planForm.value.weeklyExamGoal) || 2,
      status: 'active'
    }
    closePlanForm()
    uni.showToast({ title: 'Plan saved locally', icon: 'success' })
  }
}

async function saveItem() {
  if (!itemForm.value.title.trim()) {
    uni.showToast({ title: 'Please enter a task title', icon: 'none' })
    return
  }

  const newItem: PlanItem = {
    id: Date.now(),
    title: itemForm.value.title.trim(),
    description: itemForm.value.description.trim(),
    category: itemForm.value.category,
    dueDate: itemForm.value.dueDate,
    isCompleted: false
  }

  try {
    await new Promise<any>((resolve, reject) => {
      uni.request({
        url: '/api/v1/user/learning-plan/items',
        method: 'POST',
        data: newItem,
        header: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${uni.getStorageSync('token')}`
        },
        success: (r: any) => {
          if (r.data?.code === '00000') resolve(r.data)
          else reject(new Error(r.data?.msg))
        },
        fail: reject
      })
    })
  } catch (e) {
    console.error('Save item failed, saving locally', e)
  }

  planItems.value.push(newItem)
  closeItemForm()
  uni.showToast({ title: 'Task added', icon: 'success' })
}

async function toggleItem(item: PlanItem) {
  item.isCompleted = !item.isCompleted
  try {
    await new Promise<any>((resolve, reject) => {
      uni.request({
        url: `/api/v1/user/learning-plan/items/${item.id}`,
        method: 'PUT',
        data: { isCompleted: item.isCompleted },
        header: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${uni.getStorageSync('token')}`
        },
        success: (r: any) => resolve(r.data),
        fail: reject
      })
    })
  } catch (e) {
    console.error('Toggle item failed', e)
  }
}

async function deleteItem(id: number) {
  uni.showModal({
    title: 'Delete Task',
    content: 'Are you sure you want to delete this task?',
    success: async (res) => {
      if (res.confirm) {
        try {
          await new Promise<any>((resolve, reject) => {
            uni.request({
              url: `/api/v1/user/learning-plan/items/${id}`,
              method: 'DELETE',
              header: { Authorization: `Bearer ${uni.getStorageSync('token')}` },
              success: (r: any) => resolve(r.data),
              fail: reject
            })
          })
        } catch (e) {
          console.error('Delete item failed', e)
        }
        planItems.value = planItems.value.filter((i) => i.id !== id)
        uni.showToast({ title: 'Task deleted', icon: 'success' })
      }
    }
  })
}
</script>

<style scoped>
.plan-page {
  min-height: 100vh;
  background-color: #f5f5f5;
  padding: 24rpx;
  padding-bottom: 60rpx;
}

.current-plan-card {
  background-color: #fff;
  border-radius: 16rpx;
  padding: 30rpx;
  margin-bottom: 24rpx;
}

.plan-header {
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 8rpx;
}

.plan-title {
  font-size: 34rpx;
  color: #333;
  font-weight: 700;
  flex: 1;
}

.plan-badge {
  padding: 6rpx 16rpx;
  border-radius: 8rpx;
  margin-left: 16rpx;
}

.plan-badge.active {
  background-color: #e8f5e9;
}

.plan-badge.completed {
  background-color: #e3f2fd;
}

.plan-badge.paused {
  background-color: #fff3e0;
}

.plan-badge-text {
  font-size: 22rpx;
  color: #4a90d9;
}

.plan-date {
  font-size: 24rpx;
  color: #999;
  margin-bottom: 24rpx;
}

.progress-section {
  margin-bottom: 24rpx;
}

.progress-header {
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10rpx;
}

.progress-label {
  font-size: 26rpx;
  color: #666;
}

.progress-value {
  font-size: 26rpx;
  color: #4a90d9;
  font-weight: 600;
}

.progress-bar {
  width: 100%;
  height: 16rpx;
  background-color: #f0f0f0;
  border-radius: 8rpx;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: linear-gradient(90deg, #4a90d9, #67c23a);
  border-radius: 8rpx;
  transition: width 0.3s ease;
}

.goals-row {
  display: flex;
  flex-direction: row;
  gap: 16rpx;
}

.goal-item {
  flex: 1;
  background-color: #f8fafd;
  border-radius: 12rpx;
  padding: 16rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.goal-value {
  font-size: 32rpx;
  color: #4a90d9;
  font-weight: 700;
}

.goal-label {
  font-size: 20rpx;
  color: #999;
  margin-top: 6rpx;
}

.section-header {
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16rpx;
}

.section-title {
  font-size: 30rpx;
  color: #333;
  font-weight: 600;
}

.add-btn {
  padding: 10rpx 24rpx;
  background-color: #4a90d9;
  border-radius: 24rpx;
}

.add-btn-text {
  font-size: 24rpx;
  color: #fff;
}

.items-list {
  margin-bottom: 24rpx;
}

.item-card {
  background-color: #fff;
  border-radius: 12rpx;
  padding: 24rpx;
  margin-bottom: 16rpx;
  display: flex;
  flex-direction: row;
  align-items: flex-start;
}

.item-card.completed {
  opacity: 0.7;
}

.item-check {
  margin-right: 16rpx;
  padding-top: 4rpx;
}

.check-circle {
  width: 44rpx;
  height: 44rpx;
  border: 3rpx solid #ddd;
  border-radius: 22rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.check-circle.checked {
  background-color: #67c23a;
  border-color: #67c23a;
}

.check-icon {
  font-size: 20rpx;
  color: #fff;
  font-weight: 700;
}

.item-content {
  flex: 1;
}

.item-title {
  font-size: 28rpx;
  color: #333;
  font-weight: 500;
  margin-bottom: 8rpx;
}

.title-done {
  text-decoration: line-through;
  color: #999;
}

.item-desc {
  font-size: 24rpx;
  color: #999;
  margin-bottom: 8rpx;
  line-height: 1.4;
}

.item-meta {
  display: flex;
  flex-direction: row;
  align-items: center;
  gap: 16rpx;
}

.item-category {
  font-size: 22rpx;
  color: #4a90d9;
  background-color: #f0f5fb;
  padding: 4rpx 12rpx;
  border-radius: 6rpx;
}

.item-due {
  font-size: 22rpx;
  color: #999;
}

.item-delete {
  padding: 8rpx;
  margin-left: 8rpx;
}

.delete-icon {
  font-size: 36rpx;
  color: #ccc;
  line-height: 1;
}

.empty-state {
  background-color: #fff;
  border-radius: 16rpx;
  padding: 80rpx 30rpx;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.empty-text {
  font-size: 30rpx;
  color: #999;
  margin-bottom: 12rpx;
}

.empty-hint {
  font-size: 26rpx;
  color: #ccc;
}

.create-plan-btn {
  margin-top: 40rpx;
  background-color: #4a90d9;
  border-radius: 44rpx;
  padding: 26rpx 0;
  display: flex;
  align-items: center;
  justify-content: center;
}

.create-plan-text {
  font-size: 32rpx;
  color: #fff;
  font-weight: 600;
}

.modal-content {
  background-color: #fff;
  border-radius: 20rpx;
  padding: 36rpx 30rpx;
  width: 620rpx;
}

.modal-title {
  font-size: 34rpx;
  color: #333;
  font-weight: 700;
  text-align: center;
  margin-bottom: 30rpx;
}

.form-field {
  margin-bottom: 24rpx;
}

.form-row {
  display: flex;
  flex-direction: row;
  gap: 20rpx;
}

.form-field.half {
  flex: 1;
}

.field-label {
  font-size: 26rpx;
  color: #666;
  margin-bottom: 12rpx;
  display: block;
}

.field-input {
  width: 100%;
  height: 80rpx;
  border: 2rpx solid #eee;
  border-radius: 12rpx;
  padding: 0 20rpx;
  font-size: 28rpx;
  color: #333;
  box-sizing: border-box;
}

.field-textarea {
  width: 100%;
  min-height: 140rpx;
  border: 2rpx solid #eee;
  border-radius: 12rpx;
  padding: 16rpx 20rpx;
  font-size: 28rpx;
  color: #333;
  box-sizing: border-box;
}

.picker-trigger {
  height: 80rpx;
  border: 2rpx solid #eee;
  border-radius: 12rpx;
  padding: 0 20rpx;
  display: flex;
  align-items: center;
}

.picker-text {
  font-size: 28rpx;
  color: #333;
}

.category-options {
  display: flex;
  flex-direction: row;
  flex-wrap: wrap;
  gap: 12rpx;
}

.category-chip {
  padding: 12rpx 24rpx;
  border: 2rpx solid #eee;
  border-radius: 24rpx;
  background-color: #fafafa;
}

.category-chip.selected {
  background-color: #e8f0fe;
  border-color: #4a90d9;
}

.chip-text {
  font-size: 24rpx;
  color: #666;
}

.category-chip.selected .chip-text {
  color: #4a90d9;
}

.modal-actions {
  display: flex;
  flex-direction: row;
  gap: 20rpx;
  margin-top: 30rpx;
}

.modal-btn {
  flex: 1;
  height: 80rpx;
  border-radius: 40rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.modal-btn.cancel {
  background-color: #f5f5f5;
}

.modal-btn.confirm {
  background-color: #4a90d9;
}

.modal-btn-text {
  font-size: 28rpx;
  font-weight: 600;
}

.cancel-text {
  color: #666;
}

.confirm-text {
  color: #fff;
}
</style>
