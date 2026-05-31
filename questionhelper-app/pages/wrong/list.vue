<template>
  <view class="page">
    <!-- Search Bar -->
    <view class="search-bar">
      <view class="search-input-wrap">
        <text class="search-icon">&#xe610;</text>
        <input
          class="search-input"
          v-model="keyword"
          placeholder="搜索错题"
          confirm-type="search"
          @confirm="onSearch"
        />
        <text v-if="keyword" class="clear-btn" @tap="keyword = ''; onSearch()">&#xe621;</text>
      </view>
      <text class="filter-btn" @tap="showFilter = true">&#xe61a;</text>
    </view>

    <!-- Filter Tags -->
    <view v-if="hasActiveFilter" class="active-filters">
      <view v-if="filterForm.category" class="filter-tag">
        <text class="filter-tag-text">{{ filterForm.category }}</text>
        <text class="filter-tag-close" @tap="filterForm.category = ''; loadData(true)">&#xe621;</text>
      </view>
      <view v-if="filterForm.masteryLevel !== ''" class="filter-tag">
        <text class="filter-tag-text">{{ masteryLabels[filterForm.masteryLevel] }}</text>
        <text class="filter-tag-close" @tap="filterForm.masteryLevel = ''; loadData(true)">&#xe621;</text>
      </view>
      <view v-if="filterForm.questionType" class="filter-tag">
        <text class="filter-tag-text">{{ filterForm.questionType }}</text>
        <text class="filter-tag-close" @tap="filterForm.questionType = ''; loadData(true)">&#xe621;</text>
      </view>
      <text class="clear-all" @tap="clearFilters">清除全部</text>
    </view>

    <!-- Batch Operations Bar -->
    <view v-if="batchMode" class="batch-bar">
      <text class="batch-select-all" @tap="toggleSelectAll">
        {{ isAllSelected ? '取消全选' : '全选' }}
      </text>
      <text class="batch-count">已选 {{ selectedIds.length }} 项</text>
      <text class="batch-remove" @tap="onBatchRemove">移除选中</text>
    </view>

    <!-- Wrong Question List -->
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
        class="wrong-card"
        @tap="goDetail(item.id)"
      >
        <view v-if="batchMode" class="card-checkbox" @tap.stop="toggleSelect(item.id)">
          <view :class="['checkbox', selectedIds.includes(item.id) ? 'checked' : '']">
            <text v-if="selectedIds.includes(item.id)" class="check-icon">&#xe623;</text>
          </view>
        </view>
        <view class="card-body">
          <view class="card-header">
            <text class="card-title">{{ item.title }}</text>
            <view class="type-tag" :class="'type-' + item.questionType">
              <text class="type-tag-text">{{ typeLabels[item.questionType] }}</text>
            </view>
          </view>
          <view class="card-meta">
            <text class="wrong-count">错误 {{ item.wrongCount }} 次</text>
            <text class="next-review">下次复习: {{ item.nextReviewTime }}</text>
          </view>
          <view class="mastery-bar-wrap">
            <view class="mastery-bar-bg">
              <view
                class="mastery-bar-fill"
                :style="{ width: item.masteryLevel + '%' }"
                :class="getMasteryColor(item.masteryLevel)"
              ></view>
            </view>
            <text class="mastery-label">掌握度 {{ item.masteryLevel }}%</text>
          </view>
        </view>
      </view>

      <!-- Empty State -->
      <view v-if="!loading && list.length === 0" class="empty-state">
        <image class="empty-img" src="/static/empty/wrong.png" mode="aspectFit" />
        <text class="empty-text">暂无错题记录</text>
      </view>

      <!-- Load More -->
      <view v-if="list.length > 0" class="load-more">
        <text v-if="loading" class="load-more-text">加载中...</text>
        <text v-else-if="noMore" class="load-more-text">没有更多了</text>
      </view>
    </scroll-view>

    <!-- Bottom Action -->
    <view class="bottom-bar">
      <text class="batch-toggle" @tap="toggleBatchMode">
        {{ batchMode ? '退出批量' : '批量操作' }}
      </text>
      <text class="analysis-link" @tap="goAnalysis">错题分析</text>
    </view>

    <!-- Filter Drawer -->
    <view v-if="showFilter" class="drawer-mask" @tap="showFilter = false"></view>
    <view :class="['filter-drawer', showFilter ? 'open' : '']">
      <view class="drawer-header">
        <text class="drawer-title">筛选</text>
        <text class="drawer-close" @tap="showFilter = false">&#xe621;</text>
      </view>
      <scroll-view class="drawer-body" scroll-y>
        <!-- Category -->
        <view class="filter-section">
          <text class="filter-label">分类</text>
          <view class="filter-options">
            <view
              v-for="cat in categoryOptions"
              :key="cat"
              :class="['option-item', filterForm.category === cat ? 'selected' : '']"
              @tap="filterForm.category = filterForm.category === cat ? '' : cat"
            >
              <text class="option-text">{{ cat }}</text>
            </view>
          </view>
        </view>
        <!-- Mastery Level -->
        <view class="filter-section">
          <text class="filter-label">掌握程度</text>
          <view class="filter-options">
            <view
              v-for="(label, key) in masteryLabels"
              :key="key"
              :class="['option-item', filterForm.masteryLevel === key ? 'selected' : '']"
              @tap="filterForm.masteryLevel = filterForm.masteryLevel === key ? '' : key"
            >
              <text class="option-text">{{ label }}</text>
            </view>
          </view>
        </view>
        <!-- Question Type -->
        <view class="filter-section">
          <text class="filter-label">题目类型</text>
          <view class="filter-options">
            <view
              v-for="(label, key) in typeLabels"
              :key="key"
              :class="['option-item', filterForm.questionType === key ? 'selected' : '']"
              @tap="filterForm.questionType = filterForm.questionType === key ? '' : key"
            >
              <text class="option-text">{{ label }}</text>
            </view>
          </view>
        </view>
      </scroll-view>
      <view class="drawer-footer">
        <text class="reset-btn" @tap="resetFilter">重置</text>
        <text class="confirm-btn" @tap="applyFilter">确定</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { getWrongList, removeWrongQuestions } from '@/api/wrong'

interface WrongQuestion {
  id: string
  title: string
  questionType: string
  wrongCount: number
  masteryLevel: number
  nextReviewTime: string
  category: string
}

const keyword = ref('')
const list = ref<WrongQuestion[]>([])
const loading = ref(false)
const refreshing = ref(false)
const noMore = ref(false)
const page = ref(1)
const pageSize = 20
const showFilter = ref(false)
const batchMode = ref(false)
const selectedIds = ref<string[]>([])

const filterForm = ref({
  category: '',
  masteryLevel: '',
  questionType: ''
})

const categoryOptions = ['数学', '英语', '语文', '物理', '化学', '生物', '历史', '地理', '政治']

const masteryLabels: Record<string, string> = {
  'low': '未掌握',
  'medium': '初步掌握',
  'high': '基本掌握'
}

const typeLabels: Record<string, string> = {
  'single': '单选题',
  'multiple': '多选题',
  'judge': '判断题',
  'fill': '填空题',
  'essay': '简答题'
}

const hasActiveFilter = computed(() => {
  return filterForm.value.category || filterForm.value.masteryLevel !== '' || filterForm.value.questionType
})

const isAllSelected = computed(() => {
  return list.value.length > 0 && selectedIds.value.length === list.value.length
})

function getMasteryColor(level: number): string {
  if (level < 30) return 'mastery-low'
  if (level < 70) return 'mastery-medium'
  return 'mastery-high'
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
    const res = await getWrongList({
      keyword: keyword.value,
      page: page.value,
      pageSize,
      ...filterForm.value
    })
    const newData = res.data?.list || []
    if (reset) {
      list.value = newData
    } else {
      list.value = [...list.value, ...newData]
    }
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

function toggleBatchMode() {
  batchMode.value = !batchMode.value
  if (!batchMode.value) {
    selectedIds.value = []
  }
}

function toggleSelect(id: string) {
  const idx = selectedIds.value.indexOf(id)
  if (idx >= 0) {
    selectedIds.value.splice(idx, 1)
  } else {
    selectedIds.value.push(id)
  }
}

function toggleSelectAll() {
  if (isAllSelected.value) {
    selectedIds.value = []
  } else {
    selectedIds.value = list.value.map(item => item.id)
  }
}

async function onBatchRemove() {
  if (selectedIds.value.length === 0) {
    uni.showToast({ title: '请先选择题目', icon: 'none' })
    return
  }
  uni.showModal({
    title: '确认移除',
    content: `确定移除选中的 ${selectedIds.value.length} 道错题吗？`,
    success: async (res) => {
      if (res.confirm) {
        try {
          await removeWrongQuestions(selectedIds.value)
          uni.showToast({ title: '移除成功', icon: 'success' })
          selectedIds.value = []
          batchMode.value = false
          loadData(true)
        } catch (e) {
          uni.showToast({ title: '移除失败', icon: 'none' })
        }
      }
    }
  })
}

function clearFilters() {
  filterForm.value = { category: '', masteryLevel: '', questionType: '' }
  loadData(true)
}

function resetFilter() {
  filterForm.value = { category: '', masteryLevel: '', questionType: '' }
}

function applyFilter() {
  showFilter.value = false
  loadData(true)
}

function goDetail(id: string) {
  if (batchMode.value) {
    toggleSelect(id)
    return
  }
  uni.navigateTo({ url: `/pages/wrong/detail?id=${id}` })
}

function goAnalysis() {
  uni.navigateTo({ url: '/pages/wrong/analysis' })
}

onMounted(() => {
  loadData(true)
})
</script>

<style scoped>
.page {
  min-height: 100vh;
  background-color: #f5f6fa;
  padding-bottom: 120rpx;
}

.search-bar {
  display: flex;
  align-items: center;
  padding: 20rpx 24rpx;
  background-color: #fff;
  gap: 16rpx;
}

.search-input-wrap {
  flex: 1;
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

.filter-btn {
  font-size: 40rpx;
  color: #333;
  padding: 8rpx;
}

.active-filters {
  display: flex;
  align-items: center;
  padding: 16rpx 24rpx;
  background-color: #fff;
  gap: 12rpx;
  flex-wrap: wrap;
  border-bottom: 1rpx solid #eee;
}

.filter-tag {
  display: flex;
  align-items: center;
  background-color: #e8f3ff;
  border-radius: 24rpx;
  padding: 8rpx 20rpx;
  gap: 8rpx;
}

.filter-tag-text {
  font-size: 24rpx;
  color: #1677ff;
}

.filter-tag-close {
  font-size: 24rpx;
  color: #1677ff;
}

.clear-all {
  font-size: 24rpx;
  color: #999;
  margin-left: auto;
}

.batch-bar {
  display: flex;
  align-items: center;
  padding: 16rpx 24rpx;
  background-color: #fff8e6;
  border-bottom: 1rpx solid #f0e6c0;
}

.batch-select-all {
  font-size: 26rpx;
  color: #1677ff;
}

.batch-count {
  flex: 1;
  text-align: center;
  font-size: 26rpx;
  color: #666;
}

.batch-remove {
  font-size: 26rpx;
  color: #ff4d4f;
}

.list-container {
  height: calc(100vh - 280rpx);
  padding: 20rpx 24rpx;
}

.wrong-card {
  display: flex;
  align-items: flex-start;
  background-color: #fff;
  border-radius: 16rpx;
  margin-bottom: 20rpx;
  padding: 24rpx;
  box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.04);
}

.card-checkbox {
  margin-right: 16rpx;
  padding-top: 4rpx;
}

.checkbox {
  width: 40rpx;
  height: 40rpx;
  border: 2rpx solid #d9d9d9;
  border-radius: 50%;
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

.card-body {
  flex: 1;
  min-width: 0;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16rpx;
}

.card-title {
  flex: 1;
  font-size: 30rpx;
  font-weight: 600;
  color: #333;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  margin-right: 16rpx;
}

.type-tag {
  padding: 4rpx 16rpx;
  border-radius: 8rpx;
  flex-shrink: 0;
}

.type-tag-text {
  font-size: 22rpx;
}

.type-single {
  background-color: #e8f3ff;
}

.type-single .type-tag-text {
  color: #1677ff;
}

.type-multiple {
  background-color: #fff0e6;
}

.type-multiple .type-tag-text {
  color: #fa8c16;
}

.type-judge {
  background-color: #f6ffed;
}

.type-judge .type-tag-text {
  color: #52c41a;
}

.type-fill {
  background-color: #f9f0ff;
}

.type-fill .type-tag-text {
  color: #722ed1;
}

.type-essay {
  background-color: #fff1f0;
}

.type-essay .type-tag-text {
  color: #ff4d4f;
}

.card-meta {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 16rpx;
}

.wrong-count {
  font-size: 24rpx;
  color: #ff4d4f;
}

.next-review {
  font-size: 22rpx;
  color: #999;
}

.mastery-bar-wrap {
  display: flex;
  align-items: center;
  gap: 16rpx;
}

.mastery-bar-bg {
  flex: 1;
  height: 12rpx;
  background-color: #f0f1f5;
  border-radius: 6rpx;
  overflow: hidden;
}

.mastery-bar-fill {
  height: 100%;
  border-radius: 6rpx;
  transition: width 0.3s;
}

.mastery-low {
  background-color: #ff4d4f;
}

.mastery-medium {
  background-color: #faad14;
}

.mastery-high {
  background-color: #52c41a;
}

.mastery-label {
  font-size: 22rpx;
  color: #666;
  flex-shrink: 0;
  min-width: 140rpx;
  text-align: right;
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 120rpx 0;
}

.empty-img {
  width: 240rpx;
  height: 240rpx;
  margin-bottom: 24rpx;
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

.bottom-bar {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 20rpx 24rpx;
  padding-bottom: calc(20rpx + env(safe-area-inset-bottom));
  background-color: #fff;
  box-shadow: 0 -2rpx 12rpx rgba(0, 0, 0, 0.06);
}

.batch-toggle {
  font-size: 28rpx;
  color: #1677ff;
  padding: 16rpx 32rpx;
  background-color: #e8f3ff;
  border-radius: 36rpx;
}

.analysis-link {
  font-size: 28rpx;
  color: #fff;
  padding: 16rpx 32rpx;
  background-color: #1677ff;
  border-radius: 36rpx;
}

.drawer-mask {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  z-index: 100;
}

.filter-drawer {
  position: fixed;
  top: 0;
  right: -600rpx;
  bottom: 0;
  width: 600rpx;
  background-color: #fff;
  z-index: 101;
  display: flex;
  flex-direction: column;
  transition: right 0.3s ease;
}

.filter-drawer.open {
  right: 0;
}

.drawer-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 32rpx 24rpx;
  border-bottom: 1rpx solid #eee;
}

.drawer-title {
  font-size: 32rpx;
  font-weight: 600;
  color: #333;
}

.drawer-close {
  font-size: 36rpx;
  color: #999;
  padding: 8rpx;
}

.drawer-body {
  flex: 1;
  padding: 24rpx;
}

.filter-section {
  margin-bottom: 32rpx;
}

.filter-label {
  font-size: 28rpx;
  font-weight: 600;
  color: #333;
  margin-bottom: 16rpx;
  display: block;
}

.filter-options {
  display: flex;
  flex-wrap: wrap;
  gap: 16rpx;
}

.option-item {
  padding: 12rpx 24rpx;
  border: 2rpx solid #e0e0e0;
  border-radius: 32rpx;
  background-color: #fafafa;
}

.option-item.selected {
  border-color: #1677ff;
  background-color: #e8f3ff;
}

.option-text {
  font-size: 26rpx;
  color: #666;
}

.option-item.selected .option-text {
  color: #1677ff;
}

.drawer-footer {
  display: flex;
  padding: 20rpx 24rpx;
  padding-bottom: calc(20rpx + env(safe-area-inset-bottom));
  gap: 20rpx;
  border-top: 1rpx solid #eee;
}

.reset-btn {
  flex: 1;
  text-align: center;
  font-size: 28rpx;
  color: #666;
  padding: 20rpx 0;
  border: 2rpx solid #e0e0e0;
  border-radius: 36rpx;
}

.confirm-btn {
  flex: 1;
  text-align: center;
  font-size: 28rpx;
  color: #fff;
  padding: 20rpx 0;
  background-color: #1677ff;
  border-radius: 36rpx;
}
</style>
