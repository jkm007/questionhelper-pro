<template>
  <view class="search-container">
    <!-- 搜索头部 -->
    <view class="search-header">
      <view class="search-input-wrap">
        <text class="search-icon">🔍</text>
        <input
          v-model="keyword"
          class="search-input"
          placeholder="搜索题目、知识点..."
          focus
          confirm-type="search"
          @confirm="handleSearch"
        />
        <text v-if="keyword" class="clear-icon" @tap="clearKeyword">✕</text>
      </view>
      <text class="cancel-btn" @tap="goBack">取消</text>
    </view>

    <!-- 搜索历史 & 热搜 (未搜索时显示) -->
    <view v-if="!hasSearched" class="suggestion-section">
      <!-- 搜索历史 -->
      <view v-if="searchHistory.length > 0" class="history-section">
        <view class="section-header">
          <text class="section-title">搜索历史</text>
          <text class="section-action" @tap="clearHistory">清空</text>
        </view>
        <view class="tag-list">
          <view
            v-for="(item, index) in searchHistory"
            :key="index"
            class="history-tag"
            @tap="searchByKeyword(item)"
          >
            <text>{{ item }}</text>
          </view>
        </view>
      </view>

      <!-- 热门搜索 -->
      <view class="hot-section">
        <view class="section-header">
          <text class="section-title">热门搜索</text>
        </view>
        <view class="hot-list">
          <view
            v-for="(item, index) in hotKeywords"
            :key="index"
            class="hot-item"
            @tap="searchByKeyword(item.keyword)"
          >
            <text class="hot-rank" :class="{ 'hot-rank--top': index < 3 }">{{ index + 1 }}</text>
            <text class="hot-text">{{ item.keyword }}</text>
            <text v-if="item.hot" class="hot-icon">🔥</text>
          </view>
        </view>
      </view>
    </view>

    <!-- 搜索结果 -->
    <view v-else class="result-section">
      <view v-if="loading" class="loading-wrap">
        <text class="loading-text">搜索中...</text>
      </view>
      <template v-else>
        <view v-if="resultList.length > 0" class="result-list">
          <QuestionCard
            v-for="item in resultList"
            :key="item.id"
            :question="item"
            @tap="goToDetail(item.id)"
          />
          <view class="load-more">
            <text class="load-more-text">{{ loadMoreText }}</text>
          </view>
        </view>
        <view v-else class="empty-wrap">
          <Empty text="未找到相关题目" />
        </view>
      </template>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { searchQuestions } from '@/api/question'
import QuestionCard from '@/components/QuestionCard/index.vue'
import Empty from '@/components/Empty/index.vue'

const keyword = ref('')
const hasSearched = ref(false)
const loading = ref(false)
const resultList = ref<any[]>([])
const searchHistory = ref<string[]>([])
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const hasMore = ref(true)

const hotKeywords = ref([
  { keyword: '数据结构', hot: true },
  { keyword: '操作系统', hot: true },
  { keyword: '计算机网络', hot: true },
  { keyword: '算法设计', hot: false },
  { keyword: '数据库原理', hot: false },
  { keyword: '软件工程', hot: false },
  { keyword: '编译原理', hot: false },
  { keyword: '人工智能', hot: false }
])

const loadMoreText = computed(() => {
  if (loading.value) return '加载中...'
  if (!hasMore.value) return '没有更多了'
  return '上拉加载更多'
})

onMounted(() => {
  loadHistory()
})

const loadHistory = () => {
  try {
    const history = uni.getStorageSync('search_history')
    if (history) {
      searchHistory.value = JSON.parse(history)
    }
  } catch (e) {
    console.error('读取搜索历史失败', e)
  }
}

const saveHistory = (kw: string) => {
  const index = searchHistory.value.indexOf(kw)
  if (index > -1) {
    searchHistory.value.splice(index, 1)
  }
  searchHistory.value.unshift(kw)
  if (searchHistory.value.length > 20) {
    searchHistory.value = searchHistory.value.slice(0, 20)
  }
  try {
    uni.setStorageSync('search_history', JSON.stringify(searchHistory.value))
  } catch (e) {
    console.error('保存搜索历史失败', e)
  }
}

const clearHistory = () => {
  searchHistory.value = []
  try {
    uni.removeStorageSync('search_history')
  } catch (e) {
    console.error('清空搜索历史失败', e)
  }
}

const clearKeyword = () => {
  keyword.value = ''
  hasSearched.value = false
  resultList.value = []
}

const searchByKeyword = (kw: string) => {
  keyword.value = kw
  handleSearch()
}

const handleSearch = async () => {
  const kw = keyword.value.trim()
  if (!kw) {
    uni.showToast({ title: '请输入搜索内容', icon: 'none' })
    return
  }

  saveHistory(kw)
  hasSearched.value = true
  page.value = 1
  hasMore.value = true
  resultList.value = []

  await loadResults()
}

const loadResults = async () => {
  if (loading.value || !hasMore.value) return

  loading.value = true
  try {
    const res = await searchQuestions(keyword.value, {
      page: page.value,
      pageSize: pageSize.value
    })
    const list = res.data?.list || []
    resultList.value = [...resultList.value, ...list]
    total.value = res.data?.total || 0
    hasMore.value = resultList.value.length < total.value
    page.value++
  } catch (e) {
    console.error('搜索失败', e)
  } finally {
    loading.value = false
  }
}

const goToDetail = (id: number) => {
  uni.navigateTo({ url: `/pages/question/detail?id=${id}` })
}

const goBack = () => {
  uni.navigateBack()
}
</script>

<style lang="scss" scoped>
.search-container {
  min-height: 100vh;
  background-color: #F5F7FA;
}

.search-header {
  display: flex;
  align-items: center;
  padding: 16rpx 30rpx;
  background-color: #ffffff;
  border-bottom: 1rpx solid #EBEEF5;
}

.search-input-wrap {
  flex: 1;
  display: flex;
  align-items: center;
  height: 72rpx;
  padding: 0 24rpx;
  background-color: #F5F7FA;
  border-radius: 36rpx;

  .search-icon {
    font-size: 28rpx;
    margin-right: 12rpx;
  }

  .search-input {
    flex: 1;
    font-size: 28rpx;
    color: #303133;
  }

  .clear-icon {
    font-size: 28rpx;
    color: #909399;
    padding: 10rpx;
  }
}

.cancel-btn {
  flex-shrink: 0;
  margin-left: 20rpx;
  font-size: 28rpx;
  color: #4A90D9;
}

.suggestion-section {
  padding: 30rpx;
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 24rpx;

  .section-title {
    font-size: 30rpx;
    font-weight: bold;
    color: #303133;
  }

  .section-action {
    font-size: 26rpx;
    color: #909399;
  }
}

.history-section {
  margin-bottom: 40rpx;

  .tag-list {
    display: flex;
    flex-wrap: wrap;
    gap: 16rpx;
  }

  .history-tag {
    padding: 12rpx 28rpx;
    background-color: #ffffff;
    border-radius: 30rpx;
    font-size: 26rpx;
    color: #606266;
    box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.04);
  }
}

.hot-section {
  .hot-list {
    background-color: #ffffff;
    border-radius: 16rpx;
    overflow: hidden;
    box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.05);
  }

  .hot-item {
    display: flex;
    align-items: center;
    padding: 24rpx 30rpx;
    border-bottom: 1rpx solid #F5F7FA;

    &:last-child {
      border-bottom: none;
    }

    .hot-rank {
      flex-shrink: 0;
      width: 40rpx;
      font-size: 28rpx;
      font-weight: bold;
      color: #909399;
      margin-right: 20rpx;

      &--top {
        color: #F56C6C;
      }
    }

    .hot-text {
      flex: 1;
      font-size: 28rpx;
      color: #303133;
    }

    .hot-icon {
      font-size: 28rpx;
      margin-left: 8rpx;
    }
  }
}

.result-section {
  padding: 20rpx 30rpx;
}

.loading-wrap {
  display: flex;
  align-items: center;
  justify-content: center;
  padding-top: 100rpx;

  .loading-text {
    font-size: 28rpx;
    color: #909399;
  }
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
</style>
