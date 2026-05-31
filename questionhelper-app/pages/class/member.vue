<template>
  <view class="page">
    <!-- Search Bar -->
    <view class="search-bar">
      <view class="search-input-wrap">
        <text class="search-icon">&#xe610;</text>
        <input
          class="search-input"
          v-model="keyword"
          placeholder="搜索成员"
          confirm-type="search"
          @confirm="onSearch"
        />
        <text v-if="keyword" class="clear-btn" @tap="keyword = ''; onSearch()">&#xe621;</text>
      </view>
    </view>

    <!-- Member Count -->
    <view class="count-bar">
      <text class="count-text">共 {{ totalCount }} 名成员</text>
    </view>

    <!-- Member List -->
    <scroll-view
      class="list-container"
      scroll-y
      refresher-enabled
      :refresher-triggered="refreshing"
      @refresherrefresh="onRefresh"
      @scrolltolower="onLoadMore"
    >
      <view
        v-for="member in list"
        :key="member.id"
        class="member-card"
      >
        <image class="member-avatar" :src="member.avatar" mode="aspectFill" />
        <view class="member-body">
          <view class="member-name-row">
            <text class="member-name">{{ member.name }}</text>
            <view :class="['role-badge', 'role-' + member.role]">
              <text class="role-text">{{ roleLabels[member.role] }}</text>
            </view>
          </view>
          <text class="member-join-date">加入时间: {{ member.joinDate }}</text>
        </view>
        <view v-if="isCreator && member.role !== 'creator'" class="member-actions">
          <text class="action-btn" @tap="onShowRoleMenu(member)">&#xe61a;</text>
        </view>
      </view>

      <!-- Empty State -->
      <view v-if="!loading && list.length === 0" class="empty-state">
        <text class="empty-text">暂无成员</text>
      </view>

      <!-- Load More -->
      <view v-if="list.length > 0" class="load-more">
        <text v-if="loading" class="load-more-text">加载中...</text>
        <text v-else-if="noMore" class="load-more-text">没有更多了</text>
      </view>
    </scroll-view>

    <!-- Role Action Sheet -->
    <view v-if="showRoleSheet" class="action-mask" @tap="showRoleSheet = false"></view>
    <view :class="['action-sheet', showRoleSheet ? 'open' : '']">
      <view class="sheet-header">
        <text class="sheet-title">成员管理 - {{ selectedMember?.name }}</text>
      </view>
      <view class="sheet-body">
        <view class="sheet-item" @tap="onChangeRole('admin')">
          <text class="sheet-item-text">设为管理员</text>
        </view>
        <view class="sheet-item" @tap="onChangeRole('member')">
          <text class="sheet-item-text">设为普通成员</text>
        </view>
        <view class="sheet-item danger" @tap="onRemoveMember">
          <text class="sheet-item-text">移出班级</text>
        </view>
      </view>
      <view class="sheet-cancel" @tap="showRoleSheet = false">
        <text class="sheet-cancel-text">取消</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getClassMembers, updateMemberRole, removeMember } from '@/api/class'

interface Member {
  id: string
  name: string
  avatar: string
  role: string
  joinDate: string
}

const keyword = ref('')
const list = ref<Member[]>([])
const loading = ref(false)
const refreshing = ref(false)
const noMore = ref(false)
const page = ref(1)
const pageSize = 20
const totalCount = ref(0)
const isCreator = ref(false)
const classId = ref('')
const showRoleSheet = ref(false)
const selectedMember = ref<Member | null>(null)

const roleLabels: Record<string, string> = {
  'creator': '创建者',
  'admin': '管理员',
  'member': '成员'
}

async function loadData(reset = false) {
  if (loading.value) return
  if (reset) {
    page.value = 1
    noMore.value = false
    list.value = []
  }
  loading.value = true
  try {
    const res = await getClassMembers(classId.value, {
      keyword: keyword.value,
      page: page.value,
      pageSize
    })
    const newData = res.data?.list || []
    if (reset) {
      list.value = newData
    } else {
      list.value = [...list.value, ...newData]
    }
    totalCount.value = res.data?.total || 0
    isCreator.value = res.data?.isCreator || false
    noMore.value = newData.length < pageSize
    page.value++
  } catch (e) {
    uni.showToast({ title: '加载失败', icon: 'none' })
  } finally {
    loading.value = false
  }
}

function onSearch() {
  loadData(true)
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

function onShowRoleMenu(member: Member) {
  selectedMember.value = member
  showRoleSheet.value = true
}

async function onChangeRole(role: string) {
  if (!selectedMember.value) return
  try {
    await updateMemberRole(classId.value, selectedMember.value.id, role)
    uni.showToast({ title: '设置成功', icon: 'success' })
    showRoleSheet.value = false
    loadData(true)
  } catch (e) {
    uni.showToast({ title: '设置失败', icon: 'none' })
  }
}

async function onRemoveMember() {
  if (!selectedMember.value) return
  uni.showModal({
    title: '移出成员',
    content: `确定将「${selectedMember.value.name}」移出班级吗？`,
    success: async (res) => {
      if (res.confirm) {
        try {
          await removeMember(classId.value, selectedMember.value!.id)
          uni.showToast({ title: '已移出', icon: 'success' })
          showRoleSheet.value = false
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

.search-bar {
  padding: 16rpx 24rpx;
  background-color: #fff;
}

.search-input-wrap {
  display: flex;
  align-items: center;
  background-color: #f0f1f5;
  border-radius: 36rpx;
  padding: 0 24rpx;
  height: 72rpx;
}

.search-icon {
  font-size: 32rpx;
  color: #999;
  margin-right: 12rpx;
}

.search-input {
  flex: 1;
  font-size: 28rpx;
  height: 72rpx;
}

.clear-btn {
  font-size: 28rpx;
  color: #999;
  padding: 8rpx;
}

.count-bar {
  padding: 16rpx 24rpx;
  background-color: #fff;
  border-bottom: 1rpx solid #eee;
}

.count-text {
  font-size: 26rpx;
  color: #999;
}

.list-container {
  height: calc(100vh - 200rpx);
  padding: 16rpx 24rpx;
}

.member-card {
  display: flex;
  align-items: center;
  background-color: #fff;
  border-radius: 16rpx;
  padding: 24rpx;
  margin-bottom: 16rpx;
  box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.04);
}

.member-avatar {
  width: 88rpx;
  height: 88rpx;
  border-radius: 50%;
  flex-shrink: 0;
  margin-right: 20rpx;
}

.member-body {
  flex: 1;
  min-width: 0;
}

.member-name-row {
  display: flex;
  align-items: center;
  gap: 12rpx;
  margin-bottom: 8rpx;
}

.member-name {
  font-size: 30rpx;
  font-weight: 600;
  color: #333;
}

.role-badge {
  padding: 4rpx 14rpx;
  border-radius: 16rpx;
}

.role-creator {
  background-color: #fff0e6;
}

.role-creator .role-text {
  color: #fa8c16;
}

.role-admin {
  background-color: #e8f3ff;
}

.role-admin .role-text {
  color: #1677ff;
}

.role-member {
  background-color: #f0f1f5;
}

.role-member .role-text {
  color: #999;
}

.role-text {
  font-size: 20rpx;
}

.member-join-date {
  font-size: 24rpx;
  color: #999;
}

.member-actions {
  flex-shrink: 0;
  margin-left: 16rpx;
}

.action-btn {
  font-size: 36rpx;
  color: #ccc;
  padding: 8rpx;
}

.empty-state {
  padding: 80rpx 0;
  text-align: center;
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

.action-mask {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  z-index: 100;
}

.action-sheet {
  position: fixed;
  left: 0;
  right: 0;
  bottom: -500rpx;
  background-color: #f5f6fa;
  border-radius: 24rpx 24rpx 0 0;
  z-index: 101;
  transition: bottom 0.3s ease;
}

.action-sheet.open {
  bottom: 0;
}

.sheet-header {
  padding: 32rpx;
  text-align: center;
  border-bottom: 1rpx solid #eee;
  background-color: #fff;
  border-radius: 24rpx 24rpx 0 0;
}

.sheet-title {
  font-size: 30rpx;
  font-weight: 600;
  color: #333;
}

.sheet-body {
  background-color: #fff;
  margin-top: 16rpx;
}

.sheet-item {
  padding: 28rpx 32rpx;
  border-bottom: 1rpx solid #f5f5f5;
  text-align: center;
}

.sheet-item:last-child {
  border-bottom: none;
}

.sheet-item-text {
  font-size: 30rpx;
  color: #333;
}

.sheet-item.danger .sheet-item-text {
  color: #ff4d4f;
}

.sheet-cancel {
  background-color: #fff;
  margin-top: 16rpx;
  padding: 28rpx 0;
  text-align: center;
  margin-bottom: env(safe-area-inset-bottom);
}

.sheet-cancel-text {
  font-size: 30rpx;
  color: #999;
}
</style>
