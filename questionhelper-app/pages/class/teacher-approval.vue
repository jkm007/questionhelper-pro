<template>
  <view class="page">
    <!-- Header Stats -->
    <view class="stats-bar">
      <text class="stats-text">待审批 {{ pendingCount }} 条</text>
      <view v-if="!isEditing" class="batch-btn" @tap="isEditing = true">
        <text class="batch-btn-text">批量操作</text>
      </view>
      <view v-else class="batch-btn cancel" @tap="cancelEdit">
        <text class="batch-btn-text">取消</text>
      </view>
    </view>

    <!-- Batch Action Bar -->
    <view v-if="isEditing" class="batch-action-bar">
      <view class="select-all-wrap" @tap="toggleSelectAll">
        <view :class="['checkbox', isAllSelected ? 'checked' : '']">
          <text v-if="isAllSelected" class="check-icon">&#xe622;</text>
        </view>
        <text class="select-all-text">全选</text>
      </view>
      <view class="batch-actions">
        <view class="batch-approve-btn" @tap="onBatchApprove">
          <text class="batch-approve-text">批量通过 ({{ selectedIds.length }})</text>
        </view>
      </view>
    </view>

    <!-- Application List -->
    <scroll-view
      class="list-container"
      scroll-y
      refresher-enabled
      :refresher-triggered="refreshing"
      @refresherrefresh="onRefresh"
      @scrolltolower="onLoadMore"
    >
      <view
        v-for="item in list"
        :key="item.id"
        :class="['application-card', selectedIds.includes(item.id) ? 'selected' : '']"
      >
        <!-- Select Checkbox (batch mode) -->
        <view v-if="isEditing" class="card-checkbox" @tap="toggleSelect(item.id)">
          <view :class="['checkbox', selectedIds.includes(item.id) ? 'checked' : '']">
            <text v-if="selectedIds.includes(item.id)" class="check-icon">&#xe622;</text>
          </view>
        </view>

        <!-- Student Info -->
        <view class="card-header">
          <image class="student-avatar" :src="item.studentAvatar" mode="aspectFill" />
          <view class="student-info">
            <view class="student-name-row">
              <text class="student-name">{{ item.studentName }}</text>
              <view v-if="item.studentNo" class="student-no-badge">
                <text class="student-no-text">{{ item.studentNo }}</text>
              </view>
            </view>
            <text class="apply-time">申请时间: {{ item.applyTime }}</text>
            <text v-if="item.message" class="apply-message">留言: {{ item.message }}</text>
          </view>
        </view>

        <!-- Action Buttons -->
        <view v-if="!isEditing" class="card-actions">
          <view class="reject-btn" @tap="onShowReject(item)">
            <text class="reject-btn-text">拒绝</text>
          </view>
          <view class="approve-btn" @tap="onApprove(item)">
            <text class="approve-btn-text">通过</text>
          </view>
        </view>
      </view>

      <!-- Empty State -->
      <view v-if="!loading && list.length === 0" class="empty-state">
        <text class="empty-icon">&#xe623;</text>
        <text class="empty-text">暂无待审批的申请</text>
      </view>

      <!-- Load More -->
      <view v-if="list.length > 0" class="load-more">
        <text v-if="loading" class="load-more-text">加载中...</text>
        <text v-else-if="noMore" class="load-more-text">没有更多了</text>
      </view>
    </scroll-view>

    <!-- Reject Reason Modal -->
    <view v-if="showRejectModal" class="modal-mask" @tap="showRejectModal = false"></view>
    <view :class="['reject-modal', showRejectModal ? 'open' : '']">
      <view class="modal-header">
        <text class="modal-title">拒绝原因</text>
        <text class="modal-close" @tap="showRejectModal = false">&#xe624;</text>
      </view>
      <view class="modal-body">
        <text class="modal-label">选择常见原因（可选）</text>
        <view class="reason-tags">
          <view
            v-for="tag in reasonTags"
            :key="tag"
            :class="['reason-tag', rejectReason === tag ? 'active' : '']"
            @tap="rejectReason = tag"
          >
            <text class="reason-tag-text">{{ tag }}</text>
          </view>
        </view>
        <text class="modal-label">详细原因</text>
        <textarea
          class="reason-textarea"
          v-model="rejectReason"
          placeholder="请输入拒绝原因（可选）"
          :maxlength="200"
        />
        <text class="char-count">{{ rejectReason.length }}/200</text>
      </view>
      <view class="modal-footer">
        <view class="modal-cancel-btn" @tap="showRejectModal = false">
          <text class="modal-cancel-text">取消</text>
        </view>
        <view class="modal-confirm-btn" @tap="onConfirmReject">
          <text class="modal-confirm-text">确认拒绝</text>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { getClassApplications, approveApplication, rejectApplication, batchApproveApplications } from '@/api/class'

interface Application {
  id: string
  studentId: string
  studentName: string
  studentAvatar: string
  studentNo: string
  applyTime: string
  message: string
}

const classId = ref('')
const list = ref<Application[]>([])
const loading = ref(false)
const refreshing = ref(false)
const noMore = ref(false)
const page = ref(1)
const pageSize = 20
const pendingCount = ref(0)

// Batch mode
const isEditing = ref(false)
const selectedIds = ref<string[]>([])

// Reject modal
const showRejectModal = ref(false)
const rejectTarget = ref<Application | null>(null)
const rejectReason = ref('')
const reasonTags = ['信息不匹配', '非本班学生', '重复申请', '其他原因']

const isAllSelected = computed(() => {
  return list.value.length > 0 && selectedIds.value.length === list.value.length
})

async function loadData(reset = false) {
  if (loading.value) return
  if (reset) {
    page.value = 1
    noMore.value = false
    list.value = []
  }
  loading.value = true
  try {
    const res = await getClassApplications(classId.value, {
      page: page.value,
      pageSize,
      status: 'pending'
    })
    const newData = res.data?.list || []
    if (reset) {
      list.value = newData
    } else {
      list.value = [...list.value, ...newData]
    }
    pendingCount.value = res.data?.total || 0
    noMore.value = newData.length < pageSize
    page.value++
  } catch (e) {
    uni.showToast({ title: '加载失败', icon: 'none' })
  } finally {
    loading.value = false
  }
}

async function onRefresh() {
  refreshing.value = true
  await loadData(true)
  refreshing.value = false
}

function onLoadMore() {
  if (!noMore.value && !loading.value) {
    loadData()
  }
}

// Single approve
async function onApprove(item: Application) {
  uni.showModal({
    title: '通过申请',
    content: `确定通过「${item.studentName}」的加入申请吗？`,
    success: async (res) => {
      if (res.confirm) {
        try {
          await approveApplication(classId.value, item.id)
          uni.showToast({ title: '已通过', icon: 'success' })
          loadData(true)
        } catch (e) {
          uni.showToast({ title: '操作失败', icon: 'none' })
        }
      }
    }
  })
}

// Show reject modal
function onShowReject(item: Application) {
  rejectTarget.value = item
  rejectReason.value = ''
  showRejectModal.value = true
}

// Confirm reject
async function onConfirmReject() {
  if (!rejectTarget.value) return
  try {
    await rejectApplication(classId.value, rejectTarget.value.id, {
      reason: rejectReason.value
    })
    uni.showToast({ title: '已拒绝', icon: 'success' })
    showRejectModal.value = false
    rejectTarget.value = null
    rejectReason.value = ''
    loadData(true)
  } catch (e) {
    uni.showToast({ title: '操作失败', icon: 'none' })
  }
}

// Batch operations
function cancelEdit() {
  isEditing.value = false
  selectedIds.value = []
}

function toggleSelect(id: string) {
  const idx = selectedIds.value.indexOf(id)
  if (idx > -1) {
    selectedIds.value.splice(idx, 1)
  } else {
    selectedIds.value.push(id)
  }
}

function toggleSelectAll() {
  if (isAllSelected.value) {
    selectedIds.value = []
  } else {
    selectedIds.value = list.value.map((item) => item.id)
  }
}

async function onBatchApprove() {
  if (selectedIds.value.length === 0) {
    uni.showToast({ title: '请先选择申请', icon: 'none' })
    return
  }
  uni.showModal({
    title: '批量通过',
    content: `确定通过选中的 ${selectedIds.value.length} 条申请吗？`,
    success: async (res) => {
      if (res.confirm) {
        try {
          await batchApproveApplications(classId.value, {
            ids: selectedIds.value
          })
          uni.showToast({ title: '批量通过成功', icon: 'success' })
          cancelEdit()
          loadData(true)
        } catch (e) {
          uni.showToast({ title: '操作失败', icon: 'none' })
        }
      }
    }
  })
}

onMounted(() => {
  const pages = getCurrentPages()
  const currentPage = pages[pages.length - 1] as any
  classId.value = currentPage?.options?.classId || ''
  if (classId.value) {
    loadData(true)
  } else {
    uni.showToast({ title: '参数错误', icon: 'none' })
    setTimeout(() => uni.navigateBack(), 1500)
  }
})
</script>

<style scoped>
.page {
  min-height: 100vh;
  background-color: #f5f6fa;
}

.stats-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20rpx 24rpx;
  background-color: #fff;
  border-bottom: 1rpx solid #eee;
}

.stats-text {
  font-size: 28rpx;
  color: #666;
}

.batch-btn {
  padding: 8rpx 24rpx;
  border-radius: 28rpx;
  background-color: #e8f3ff;
}

.batch-btn.cancel {
  background-color: #f5f5f5;
}

.batch-btn-text {
  font-size: 24rpx;
  color: #1677ff;
}

.batch-btn.cancel .batch-btn-text {
  color: #999;
}

.batch-action-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16rpx 24rpx;
  background-color: #fffbe6;
  border-bottom: 1rpx solid #ffe58f;
}

.select-all-wrap {
  display: flex;
  align-items: center;
  gap: 12rpx;
}

.select-all-text {
  font-size: 26rpx;
  color: #333;
}

.checkbox {
  width: 40rpx;
  height: 40rpx;
  border: 2rpx solid #d9d9d9;
  border-radius: 8rpx;
  display: flex;
  align-items: center;
  justify-content: center;
}

.checkbox.checked {
  background-color: #1677ff;
  border-color: #1677ff;
}

.check-icon {
  font-size: 24rpx;
  color: #fff;
}

.batch-actions {
  display: flex;
  gap: 16rpx;
}

.batch-approve-btn {
  background-color: #52c41a;
  padding: 12rpx 28rpx;
  border-radius: 28rpx;
}

.batch-approve-text {
  font-size: 24rpx;
  color: #fff;
}

.list-container {
  height: calc(100vh - 120rpx);
  padding: 16rpx 24rpx;
}

.application-card {
  background-color: #fff;
  border-radius: 16rpx;
  padding: 24rpx;
  margin-bottom: 16rpx;
  box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.04);
  display: flex;
  flex-direction: column;
}

.application-card.selected {
  border: 2rpx solid #1677ff;
  background-color: #f0f7ff;
}

.card-checkbox {
  margin-bottom: 16rpx;
}

.card-header {
  display: flex;
  align-items: flex-start;
}

.student-avatar {
  width: 88rpx;
  height: 88rpx;
  border-radius: 50%;
  flex-shrink: 0;
  margin-right: 20rpx;
}

.student-info {
  flex: 1;
  min-width: 0;
}

.student-name-row {
  display: flex;
  align-items: center;
  gap: 12rpx;
  margin-bottom: 8rpx;
}

.student-name {
  font-size: 30rpx;
  font-weight: 600;
  color: #333;
}

.student-no-badge {
  padding: 2rpx 12rpx;
  background-color: #f0f1f5;
  border-radius: 12rpx;
}

.student-no-text {
  font-size: 20rpx;
  color: #999;
}

.apply-time {
  font-size: 24rpx;
  color: #999;
  display: block;
  margin-bottom: 6rpx;
}

.apply-message {
  font-size: 24rpx;
  color: #666;
  line-height: 1.5;
  display: block;
}

.card-actions {
  display: flex;
  justify-content: flex-end;
  gap: 20rpx;
  margin-top: 20rpx;
  padding-top: 20rpx;
  border-top: 1rpx solid #f5f5f5;
}

.reject-btn {
  padding: 12rpx 36rpx;
  border-radius: 28rpx;
  border: 2rpx solid #ff4d4f;
}

.reject-btn-text {
  font-size: 26rpx;
  color: #ff4d4f;
}

.approve-btn {
  padding: 12rpx 36rpx;
  border-radius: 28rpx;
  background-color: #1677ff;
}

.approve-btn-text {
  font-size: 26rpx;
  color: #fff;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 120rpx 0;
}

.empty-icon {
  font-size: 80rpx;
  color: #ddd;
  margin-bottom: 20rpx;
}

.empty-text {
  font-size: 28rpx;
  color: #999;
}

.load-more {
  padding: 24rpx 0;
  text-align: center;
}

.load-more-text {
  font-size: 24rpx;
  color: #999;
}

/* Reject Modal */
.modal-mask {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  z-index: 100;
}

.reject-modal {
  position: fixed;
  left: 0;
  right: 0;
  bottom: -800rpx;
  background-color: #fff;
  border-radius: 24rpx 24rpx 0 0;
  z-index: 101;
  transition: bottom 0.3s ease;
}

.reject-modal.open {
  bottom: 0;
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 28rpx 32rpx;
  border-bottom: 1rpx solid #f0f0f0;
}

.modal-title {
  font-size: 32rpx;
  font-weight: 600;
  color: #333;
}

.modal-close {
  font-size: 32rpx;
  color: #999;
  padding: 8rpx;
}

.modal-body {
  padding: 24rpx 32rpx;
}

.modal-label {
  font-size: 26rpx;
  color: #666;
  margin-bottom: 16rpx;
  display: block;
}

.reason-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 16rpx;
  margin-bottom: 24rpx;
}

.reason-tag {
  padding: 10rpx 24rpx;
  background-color: #f5f5f5;
  border-radius: 28rpx;
  border: 2rpx solid transparent;
}

.reason-tag.active {
  background-color: #fff1f0;
  border-color: #ff4d4f;
}

.reason-tag-text {
  font-size: 24rpx;
  color: #666;
}

.reason-tag.active .reason-tag-text {
  color: #ff4d4f;
}

.reason-textarea {
  width: 100%;
  height: 160rpx;
  font-size: 28rpx;
  padding: 16rpx;
  border: 2rpx solid #e0e0e0;
  border-radius: 12rpx;
  box-sizing: border-box;
}

.char-count {
  display: block;
  text-align: right;
  font-size: 22rpx;
  color: #ccc;
  margin-top: 8rpx;
}

.modal-footer {
  display: flex;
  gap: 20rpx;
  padding: 20rpx 32rpx;
  padding-bottom: calc(20rpx + env(safe-area-inset-bottom));
  border-top: 1rpx solid #f0f0f0;
}

.modal-cancel-btn {
  flex: 1;
  padding: 20rpx 0;
  text-align: center;
  border-radius: 36rpx;
  background-color: #f5f5f5;
}

.modal-cancel-text {
  font-size: 28rpx;
  color: #666;
}

.modal-confirm-btn {
  flex: 1;
  padding: 20rpx 0;
  text-align: center;
  border-radius: 36rpx;
  background-color: #ff4d4f;
}

.modal-confirm-text {
  font-size: 28rpx;
  color: #fff;
}
</style>
