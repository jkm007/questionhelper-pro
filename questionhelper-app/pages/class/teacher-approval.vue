<template>
  <view class="approval-page">
    <!-- Pending List -->
    <view v-if="applications.length" class="app-list">
      <view v-for="app in applications" :key="app.id" class="app-card">
        <view class="app-header">
          <text class="app-name">{{ app.applicantName }}</text>
          <text class="app-time">{{ app.applyTime }}</text>
        </view>
        <text class="app-reason">{{ app.reason }}</text>
        <view class="app-actions">
          <view class="action-btn reject" @tap="showRejectDialog(app)">
            <text class="action-text reject-text">拒绝</text>
          </view>
          <view class="action-btn approve" @tap="handleApprove(app.id)">
            <text class="action-text approve-text">通过</text>
          </view>
        </view>
      </view>
    </view>

    <!-- Empty State -->
    <view v-else class="empty-state">
      <text class="empty-text">暂无待审批的申请</text>
    </view>

    <!-- Reject Dialog -->
    <uni-popup ref="rejectPopup" type="center">
      <view class="dialog">
        <text class="dialog-title">拒绝原因</text>
        <textarea
          class="dialog-textarea"
          v-model="rejectReason"
          placeholder="请输入拒绝原因（选填）"
          :maxlength="200"
        />
        <view class="dialog-actions">
          <view class="dialog-btn cancel" @tap="rejectPopup?.close()">
            <text class="btn-text cancel-text">取消</text>
          </view>
          <view class="dialog-btn confirm" @tap="confirmReject">
            <text class="btn-text confirm-text">确认拒绝</text>
          </view>
        </view>
      </view>
    </uni-popup>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'

const rejectPopup = ref<any>(null)
const rejectReason = ref('')
const rejectingId = ref(0)

const classId = ref('')

interface Application {
  id: number
  applicantName: string
  applyTime: string
  reason: string
}

const applications = ref<Application[]>([])

onMounted(() => {
  const pages = getCurrentPages()
  const page = pages[pages.length - 1] as any
  classId.value = page.options?.classId || ''
  fetchApplications()
})

async function fetchApplications() {
  try {
    // TODO: replace with actual API call using classId
    applications.value = [
      { id: 1, applicantName: '张三', applyTime: '2026-05-28 14:30', reason: '我是一名有3年教学经验的高中数学教师，希望加入班级进行教学管理。' },
      { id: 2, applicantName: '李四', applyTime: '2026-05-29 09:15', reason: '我是计算机科学专业的研究生，擅长算法与数据结构，希望协助辅导。' }
    ]
  } catch (e) {
    console.error('Failed to load applications', e)
  }
}

function showRejectDialog(app: Application) {
  rejectingId.value = app.id
  rejectReason.value = ''
  rejectPopup.value?.open()
}

async function handleApprove(id: number) {
  uni.showModal({
    title: '确认通过',
    content: '确定通过该教师申请？',
    success: async (res) => {
      if (res.confirm) {
        try {
          // TODO: replace with actual API call
          applications.value = applications.value.filter((a) => a.id !== id)
          uni.showToast({ title: '已通过', icon: 'success' })
        } catch (e) {
          console.error('Failed to approve', e)
        }
      }
    }
  })
}

async function confirmReject() {
  try {
    // TODO: replace with actual API call with rejectingId and rejectReason
    applications.value = applications.value.filter((a) => a.id !== rejectingId.value)
    rejectPopup.value?.close()
    uni.showToast({ title: '已拒绝', icon: 'success' })
  } catch (e) {
    console.error('Failed to reject', e)
  }
}
</script>

<style scoped>
.approval-page {
  min-height: 100vh;
  background-color: #f5f5f5;
  padding: 20rpx 24rpx;
  padding-bottom: 80rpx;
}

.app-card {
  background-color: #fff;
  border-radius: 16rpx;
  padding: 28rpx 30rpx;
  margin-bottom: 20rpx;
}

.app-header {
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16rpx;
}

.app-name {
  font-size: 32rpx;
  font-weight: 600;
  color: #333;
}

.app-time {
  font-size: 22rpx;
  color: #bbb;
}

.app-reason {
  font-size: 26rpx;
  color: #666;
  line-height: 1.6;
  margin-bottom: 24rpx;
}

.app-actions {
  display: flex;
  flex-direction: row;
  gap: 20rpx;
  justify-content: flex-end;
}

.action-btn {
  padding: 14rpx 40rpx;
  border-radius: 32rpx;
}

.action-btn.reject {
  background-color: #f0f0f0;
}

.action-btn.approve {
  background-color: #4a90d9;
}

.action-text {
  font-size: 26rpx;
}

.reject-text {
  color: #666;
}

.approve-text {
  color: #fff;
}

.empty-state {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 120rpx 0;
}

.empty-text {
  font-size: 28rpx;
  color: #ccc;
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
  margin-bottom: 24rpx;
  text-align: center;
}

.dialog-textarea {
  width: 100%;
  height: 160rpx;
  font-size: 28rpx;
  color: #333;
  padding: 16rpx;
  background-color: #f9f9f9;
  border-radius: 12rpx;
  margin-bottom: 30rpx;
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
  background-color: #e74c3c;
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
