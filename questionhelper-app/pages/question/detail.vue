<template>
  <view class="detail-container">
    <!-- 加载状态 -->
    <view v-if="loading" class="loading-wrap">
      <text class="loading-text">加载中...</text>
    </view>

    <template v-else-if="question">
      <!-- 题目信息 -->
      <view class="question-section">
        <view class="question-header">
          <view class="type-tag" :class="typeClass">{{ typeText }}</view>
          <view class="difficulty-badge" :class="difficultyClass">{{ difficultyText }}</view>
        </view>
        <text class="question-title">{{ question.title }}</text>
        <view class="question-meta">
          <text class="meta-item">{{ question.category?.name || '未分类' }}</text>
          <text class="meta-divider">|</text>
          <text class="meta-item">浏览 {{ question.viewCount || 0 }}</text>
          <text class="meta-divider">|</text>
          <text class="meta-item">点赞 {{ question.likeCount || 0 }}</text>
        </view>
      </view>

      <!-- 题目内容 -->
      <view class="content-section">
        <text class="section-label">题目内容</text>
        <view class="rich-content">
          <text class="content-text">{{ question.content }}</text>
        </view>
      </view>

      <!-- 选项 (单选/多选/判断题) -->
      <view v-if="showOptions" class="options-section">
        <text class="section-label">选项</text>
        <view class="option-list">
          <view
            v-for="(option, index) in question.options"
            :key="index"
            class="option-item"
            :class="{
              'option-item--selected': selectedOptions.includes(index),
              'option-item--correct': showAnswer && isCorrectOption(index),
              'option-item--wrong': showAnswer && selectedOptions.includes(index) && !isCorrectOption(index)
            }"
            @tap="selectOption(index)"
          >
            <view class="option-label">{{ optionLabels[index] }}</view>
            <text class="option-text">{{ option.content }}</text>
          </view>
        </view>
      </view>

      <!-- 答案区域 -->
      <view class="answer-section">
        <view class="section-toggle" @tap="showAnswer = !showAnswer">
          <text class="section-label">参考答案</text>
          <text class="toggle-icon">{{ showAnswer ? '▲' : '▼' }}</text>
        </view>
        <view v-if="showAnswer" class="answer-content">
          <text class="answer-text">{{ question.answer || '暂无答案' }}</text>
        </view>
      </view>

      <!-- 解析区域 -->
      <view class="analysis-section">
        <view class="section-toggle" @tap="showAnalysis = !showAnalysis">
          <text class="section-label">题目解析</text>
          <text class="toggle-icon">{{ showAnalysis ? '▲' : '▼' }}</text>
        </view>
        <view v-if="showAnalysis" class="analysis-content">
          <text class="analysis-text">{{ question.analysis || '暂无解析' }}</text>
        </view>
      </view>

      <!-- 知识点标签 -->
      <view v-if="question.knowledgePoints?.length" class="knowledge-section">
        <text class="section-label">知识点</text>
        <view class="knowledge-tags">
          <view
            v-for="(point, index) in question.knowledgePoints"
            :key="index"
            class="knowledge-tag"
          >
            <text>{{ point.name }}</text>
          </view>
        </view>
      </view>

      <!-- 创建者信息 -->
      <view v-if="question.creator" class="creator-section">
        <text class="section-label">出题人</text>
        <view class="creator-info">
          <image
            class="creator-avatar"
            :src="question.creator.avatar || '/static/images/default-avatar.png'"
            mode="aspectFill"
          ></image>
          <text class="creator-name">{{ question.creator.nickname || '匿名用户' }}</text>
        </view>
      </view>
    </template>

    <!-- 空状态 -->
    <view v-else class="empty-wrap">
      <Empty text="题目不存在" />
    </view>

    <!-- 底部操作栏 -->
    <view v-if="question" class="bottom-bar">
      <view class="bar-item" @tap="handleFavorite">
        <text class="bar-icon">{{ isFavorited ? '⭐' : '☆' }}</text>
        <text class="bar-text">{{ isFavorited ? '已收藏' : '收藏' }}</text>
      </view>
      <view class="bar-item" @tap="handleLike">
        <text class="bar-icon">{{ isLiked ? '❤️' : '🤍' }}</text>
        <text class="bar-text">{{ isLiked ? '已点赞' : '点赞' }}</text>
      </view>
      <view class="bar-item" @tap="handleShare">
        <text class="bar-icon">📤</text>
        <text class="bar-text">分享</text>
      </view>
      <view class="bar-item" @tap="handleReport">
        <text class="bar-icon">⚠️</text>
        <text class="bar-text">纠错</text>
      </view>
      <view class="bar-item" @tap="handleAddWrong">
        <text class="bar-icon">📝</text>
        <text class="bar-text">错题本</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { getQuestionDetail, favoriteQuestion, likeQuestion } from '@/api/question'
import Empty from '@/components/Empty/index.vue'

const question = ref<any>(null)
const loading = ref(true)
const showAnswer = ref(false)
const showAnalysis = ref(false)
const isFavorited = ref(false)
const isLiked = ref(false)
const selectedOptions = ref<number[]>([])

const optionLabels = ['A', 'B', 'C', 'D', 'E', 'F', 'G', 'H']

const showOptions = computed(() => {
  return question.value?.type && question.value.type <= 3 && question.value.options?.length
})

const typeClass = computed(() => {
  const typeMap: Record<number, string> = {
    1: 'single',
    2: 'multiple',
    3: 'judge',
    4: 'fill',
    5: 'short'
  }
  return typeMap[question.value?.type] || ''
})

const typeText = computed(() => {
  const typeMap: Record<number, string> = {
    1: '单选题',
    2: '多选题',
    3: '判断题',
    4: '填空题',
    5: '简答题'
  }
  return typeMap[question.value?.type] || '未知'
})

const difficultyClass = computed(() => {
  const map: Record<number, string> = {
    1: 'easy',
    2: 'medium',
    3: 'hard'
  }
  return map[question.value?.difficulty] || ''
})

const difficultyText = computed(() => {
  const map: Record<number, string> = {
    1: '简单',
    2: '中等',
    3: '困难'
  }
  return map[question.value?.difficulty] || '未知'
})

onMounted(() => {
  const pages = getCurrentPages()
  const currentPage = pages[pages.length - 1] as any
  const id = currentPage?.options?.id
  if (id) {
    loadDetail(Number(id))
  } else {
    loading.value = false
  }
})

const loadDetail = async (id: number) => {
  loading.value = true
  try {
    const res = await getQuestionDetail(id)
    question.value = res.data
    isFavorited.value = res.data?.isFavorited || false
    isLiked.value = res.data?.isLiked || false
  } catch (e) {
    console.error('加载题目详情失败', e)
  } finally {
    loading.value = false
  }
}

const isCorrectOption = (index: number) => {
  const correctAnswer = question.value?.correctOptions || []
  return correctAnswer.includes(index)
}

const selectOption = (index: number) => {
  if (showAnswer.value) return
  if (question.value?.type === 1 || question.value?.type === 3) {
    selectedOptions.value = [index]
  } else if (question.value?.type === 2) {
    const idx = selectedOptions.value.indexOf(index)
    if (idx > -1) {
      selectedOptions.value.splice(idx, 1)
    } else {
      selectedOptions.value.push(index)
    }
  }
}

const handleFavorite = async () => {
  if (!question.value?.id) return
  try {
    await favoriteQuestion(question.value.id)
    isFavorited.value = !isFavorited.value
    uni.showToast({
      title: isFavorited.value ? '收藏成功' : '已取消收藏',
      icon: 'success'
    })
  } catch (e) {
    console.error('收藏操作失败', e)
  }
}

const handleLike = async () => {
  if (!question.value?.id) return
  try {
    await likeQuestion(question.value.id)
    isLiked.value = !isLiked.value
    if (isLiked.value) {
      question.value.likeCount = (question.value.likeCount || 0) + 1
    } else {
      question.value.likeCount = Math.max(0, (question.value.likeCount || 0) - 1)
    }
    uni.showToast({
      title: isLiked.value ? '点赞成功' : '已取消点赞',
      icon: 'success'
    })
  } catch (e) {
    console.error('点赞操作失败', e)
  }
}

const handleShare = () => {
  uni.showToast({ title: '分享功能开发中', icon: 'none' })
}

const handleReport = () => {
  if (question.value?.id) {
    uni.navigateTo({ url: `/pages/feedback/correction?id=${question.value.id}` })
  }
}

const handleAddWrong = () => {
  uni.showToast({ title: '已加入错题本', icon: 'success' })
}
</script>

<style lang="scss" scoped>
.detail-container {
  min-height: 100vh;
  background-color: #F5F7FA;
  padding-bottom: 140rpx;
}

.loading-wrap {
  display: flex;
  align-items: center;
  justify-content: center;
  padding-top: 200rpx;

  .loading-text {
    font-size: 28rpx;
    color: #909399;
  }
}

.empty-wrap {
  padding-top: 200rpx;
}

.question-section {
  background-color: #ffffff;
  padding: 30rpx;
  margin-bottom: 20rpx;

  .question-header {
    display: flex;
    align-items: center;
    margin-bottom: 20rpx;
    gap: 16rpx;
  }

  .type-tag {
    font-size: 22rpx;
    padding: 6rpx 16rpx;
    border-radius: 8rpx;

    &.single {
      background-color: #E6F7FF;
      color: #1890FF;
    }

    &.multiple {
      background-color: #FFF7E6;
      color: #FA8C16;
    }

    &.judge {
      background-color: #F6FFED;
      color: #52C41A;
    }

    &.fill {
      background-color: #F9F0FF;
      color: #722ED1;
    }

    &.short {
      background-color: #FFF1F0;
      color: #F5222D;
    }
  }

  .difficulty-badge {
    font-size: 22rpx;
    padding: 6rpx 16rpx;
    border-radius: 8rpx;
    border: 1rpx solid;

    &.easy {
      color: #67C23A;
      border-color: #67C23A;
    }

    &.medium {
      color: #E6A23C;
      border-color: #E6A23C;
    }

    &.hard {
      color: #F56C6C;
      border-color: #F56C6C;
    }
  }

  .question-title {
    font-size: 34rpx;
    font-weight: bold;
    color: #303133;
    line-height: 1.6;
    margin-bottom: 20rpx;
  }

  .question-meta {
    display: flex;
    align-items: center;

    .meta-item {
      font-size: 24rpx;
      color: #909399;
    }

    .meta-divider {
      margin: 0 12rpx;
      color: #DCDFE6;
    }
  }
}

.section-label {
  display: block;
  font-size: 30rpx;
  font-weight: bold;
  color: #303133;
  margin-bottom: 20rpx;
}

.section-toggle {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 30rpx;

  .toggle-icon {
    font-size: 24rpx;
    color: #909399;
  }
}

.content-section {
  background-color: #ffffff;
  padding: 30rpx;
  margin-bottom: 20rpx;

  .content-text {
    font-size: 30rpx;
    color: #303133;
    line-height: 1.8;
  }
}

.options-section {
  background-color: #ffffff;
  padding: 30rpx;
  margin-bottom: 20rpx;

  .option-list {
    .option-item {
      display: flex;
      align-items: flex-start;
      padding: 20rpx 24rpx;
      margin-bottom: 16rpx;
      border: 2rpx solid #EBEEF5;
      border-radius: 12rpx;

      &.option-item--selected {
        border-color: #4A90D9;
        background-color: #E6F7FF;
      }

      &.option-item--correct {
        border-color: #67C23A;
        background-color: #F6FFED;
      }

      &.option-item--wrong {
        border-color: #F56C6C;
        background-color: #FFF1F0;
      }

      .option-label {
        flex-shrink: 0;
        width: 48rpx;
        height: 48rpx;
        border-radius: 50%;
        background-color: #F5F7FA;
        display: flex;
        align-items: center;
        justify-content: center;
        font-size: 26rpx;
        color: #606266;
        margin-right: 16rpx;
      }

      .option-text {
        flex: 1;
        font-size: 28rpx;
        color: #303133;
        line-height: 1.5;
        padding-top: 6rpx;
      }
    }
  }
}

.answer-section {
  background-color: #ffffff;
  margin-bottom: 20rpx;

  .answer-content {
    padding: 0 30rpx 30rpx;

    .answer-text {
      font-size: 28rpx;
      color: #67C23A;
      line-height: 1.6;
    }
  }
}

.analysis-section {
  background-color: #ffffff;
  margin-bottom: 20rpx;

  .analysis-content {
    padding: 0 30rpx 30rpx;

    .analysis-text {
      font-size: 28rpx;
      color: #606266;
      line-height: 1.8;
    }
  }
}

.knowledge-section {
  background-color: #ffffff;
  padding: 30rpx;
  margin-bottom: 20rpx;

  .knowledge-tags {
    display: flex;
    flex-wrap: wrap;
    gap: 16rpx;
  }

  .knowledge-tag {
    padding: 8rpx 24rpx;
    background-color: #F5F7FA;
    border-radius: 30rpx;
    font-size: 24rpx;
    color: #606266;
  }
}

.creator-section {
  background-color: #ffffff;
  padding: 30rpx;
  margin-bottom: 20rpx;

  .creator-info {
    display: flex;
    align-items: center;
  }

  .creator-avatar {
    width: 72rpx;
    height: 72rpx;
    border-radius: 50%;
    margin-right: 20rpx;
  }

  .creator-name {
    font-size: 28rpx;
    color: #303133;
  }
}

.bottom-bar {
  position: fixed;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  align-items: center;
  justify-content: space-around;
  height: 120rpx;
  background-color: #ffffff;
  border-top: 1rpx solid #EBEEF5;
  padding-bottom: env(safe-area-inset-bottom);
  z-index: 100;

  .bar-item {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 10rpx 20rpx;
  }

  .bar-icon {
    font-size: 40rpx;
    margin-bottom: 4rpx;
  }

  .bar-text {
    font-size: 22rpx;
    color: #606266;
  }
}
</style>
