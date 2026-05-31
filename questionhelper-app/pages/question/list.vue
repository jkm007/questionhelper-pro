<template>
  <view class="question-list-container">
    <!-- 搜索栏 -->
    <view class="search-bar" @tap="goToSearch">
      <text class="search-icon">🔍</text>
      <text class="search-placeholder">搜索题目...</text>
    </view>

    <!-- 筛选栏 -->
    <view class="filter-bar">
      <view
        class="filter-item"
        :class="{ active: currentType === 0 }"
        @tap="changeType(0)"
      >
        <text>全部</text>
      </view>
      <view
        class="filter-item"
        :class="{ active: currentType === 1 }"
        @tap="changeType(1)"
      >
        <text>单选题</text>
      </view>
      <view
        class="filter-item"
        :class="{ active: currentType === 2 }"
        @tap="changeType(2)"
      >
        <text>多选题</text>
      </view>
      <view
        class="filter-item"
        :class="{ active: currentType === 3 }"
        @tap="changeType(3)"
      >
        <text>判断题</text>
      </view>
      <view class="filter-item filter-item--category" @tap="showCategoryDrawer = true">
        <text>{{ currentCategoryName || '分类' }}</text>
        <text class="filter-arrow">▼</text>
      </view>
    </view>

    <!-- 题目列表 -->
    <scroll-view
      class="question-scroll"
      scroll-y
      :refresher-enabled="true"
      :refresher-triggered="isRefreshing"
      @refresherrefresh="onRefresh"
      @scrolltolower="onLoadMore"
    >
      <view class="question-content">
        <QuestionCard
          v-for="item in questionList"
          :key="item.id"
          :question="item"
          @tap="goToDetail(item.id)"
        />
        <view v-if="questionList.length === 0 && !loading" class="empty-wrap">
          <Empty text="暂无题目" />
        </view>
        <view v-if="questionList.length > 0" class="load-more">
          <text class="load-more-text">{{ loadMoreText }}</text>
        </view>
      </view>
    </scroll-view>

    <!-- 分类抽屉 -->
    <view v-if="showCategoryDrawer" class="drawer-mask" @tap="showCategoryDrawer = false">
      <view class="drawer" @tap.stop>
        <view class="drawer-header">
          <text class="drawer-title">选择分类</text>
          <text class="drawer-close" @tap="showCategoryDrawer = false">✕</text>
        </view>
        <scroll-view class="drawer-body" scroll-y>
          <view
            class="category-item"
            :class="{ active: currentCategoryId === 0 }"
            @tap="selectCategory(0, '全部分类')"
          >
            <text>全部分类</text>
          </view>
          <view
            v-for="item in categoryTree"
            :key="item.id"
            class="category-group"
          >
            <view class="category-parent" @tap="selectCategory(item.id, item.name)">
              <text :class="{ active: currentCategoryId === item.id }">{{ item.name }}</text>
            </view>
            <view
              v-for="child in item.children"
              :key="child.id"
              class="category-child"
              @tap="selectCategory(child.id, child.name)"
            >
              <text :class="{ active: currentCategoryId === child.id }">{{ child.name }}</text>
            </view>
          </view>
        </scroll-view>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { getQuestions, getCategoryTree } from '@/api/question'
import QuestionCard from '@/components/QuestionCard/index.vue'
import Empty from '@/components/Empty/index.vue'

const questionList = ref<any[]>([])
const categoryTree = ref<any[]>([])
const loading = ref(false)
const isRefreshing = ref(false)
const showCategoryDrawer = ref(false)
const currentType = ref(0)
const currentCategoryId = ref(0)
const currentCategoryName = ref('')

const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const hasMore = ref(true)

const loadMoreText = computed(() => {
  if (loading.value) return '加载中...'
  if (!hasMore.value) return '没有更多了'
  return '上拉加载更多'
})

onMounted(() => {
  loadCategoryTree()
  loadQuestions()
})

const loadCategoryTree = async () => {
  try {
    const res = await getCategoryTree()
    categoryTree.value = res.data || []
  } catch (e) {
    console.error('加载分类失败', e)
  }
}

const loadQuestions = async (reset = false) => {
  if (loading.value) return
  if (reset) {
    page.value = 1
    hasMore.value = true
  }
  if (!hasMore.value) return

  loading.value = true
  try {
    const params: any = {
      page: page.value,
      pageSize: pageSize.value
    }
    if (currentType.value > 0) {
      params.type = currentType.value
    }
    if (currentCategoryId.value > 0) {
      params.categoryId = currentCategoryId.value
    }
    const res = await getQuestions(params)
    const list = res.data?.list || []
    if (reset) {
      questionList.value = list
    } else {
      questionList.value = [...questionList.value, ...list]
    }
    total.value = res.data?.total || 0
    hasMore.value = questionList.value.length < total.value
    page.value++
  } catch (e) {
    console.error('加载题目失败', e)
  } finally {
    loading.value = false
  }
}

const onRefresh = async () => {
  isRefreshing.value = true
  await loadQuestions(true)
  isRefreshing.value = false
}

const onLoadMore = () => {
  if (!hasMore.value || loading.value) return
  loadQuestions()
}

const changeType = (type: number) => {
  currentType.value = type
  loadQuestions(true)
}

const selectCategory = (id: number, name: string) => {
  currentCategoryId.value = id
  currentCategoryName.value = id === 0 ? '' : name
  showCategoryDrawer.value = false
  loadQuestions(true)
}

const goToSearch = () => {
  uni.navigateTo({ url: '/pages/question/search' })
}

const goToDetail = (id: number) => {
  uni.navigateTo({ url: `/pages/question/detail?id=${id}` })
}
</script>

<style lang="scss" scoped>
.question-list-container {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background-color: #F5F7FA;
}

.search-bar {
  display: flex;
  align-items: center;
  margin: 20rpx 30rpx;
  padding: 20rpx 30rpx;
  background-color: #ffffff;
  border-radius: 40rpx;
  box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.05);

  .search-icon {
    margin-right: 16rpx;
  }

  .search-placeholder {
    color: #909399;
    font-size: 28rpx;
  }
}

.filter-bar {
  display: flex;
  align-items: center;
  padding: 16rpx 30rpx;
  background-color: #ffffff;
  border-bottom: 1rpx solid #EBEEF5;
  overflow-x: auto;
  white-space: nowrap;

  .filter-item {
    flex-shrink: 0;
    padding: 12rpx 24rpx;
    margin-right: 16rpx;
    font-size: 26rpx;
    color: #606266;
    border-radius: 30rpx;
    background-color: #F5F7FA;

    &.active {
      background-color: #4A90D9;
      color: #ffffff;
    }

    &--category {
      display: flex;
      align-items: center;
    }

    .filter-arrow {
      font-size: 20rpx;
      margin-left: 8rpx;
    }
  }
}

.question-scroll {
  flex: 1;
  overflow: hidden;
}

.question-content {
  padding: 20rpx 30rpx;
}

.empty-wrap {
  padding-top: 100rpx;
}

.load-more {
  text-align: center;
  padding: 30rpx 0;

  .load-more-text {
    font-size: 26rpx;
    color: #909399;
  }
}

.drawer-mask {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  z-index: 999;
  display: flex;
  justify-content: flex-end;
}

.drawer {
  width: 600rpx;
  height: 100%;
  background-color: #ffffff;
  display: flex;
  flex-direction: column;

  .drawer-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 30rpx;
    border-bottom: 1rpx solid #EBEEF5;

    .drawer-title {
      font-size: 32rpx;
      font-weight: bold;
      color: #303133;
    }

    .drawer-close {
      font-size: 36rpx;
      color: #909399;
      padding: 10rpx;
    }
  }

  .drawer-body {
    flex: 1;
    padding: 20rpx 0;
  }
}

.category-item {
  padding: 24rpx 30rpx;
  font-size: 28rpx;
  color: #606266;

  &.active {
    color: #4A90D9;
    background-color: #F5F7FA;
  }
}

.category-parent {
  padding: 24rpx 30rpx;
  font-size: 30rpx;
  font-weight: bold;
  color: #303133;

  .active {
    color: #4A90D9;
  }
}

.category-child {
  padding: 20rpx 60rpx;
  font-size: 28rpx;
  color: #606266;

  .active {
    color: #4A90D9;
  }
}
</style>
