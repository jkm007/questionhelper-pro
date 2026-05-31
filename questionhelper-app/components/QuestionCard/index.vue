<template>
  <view class="question-card" @tap="$emit('tap')">
    <view class="card-header">
      <view class="type-tag" :class="typeClass">{{ typeText }}</view>
      <view class="difficulty" :class="difficultyClass">{{ difficultyText }}</view>
    </view>
    <view class="card-body">
      <text class="question-title">{{ question.title }}</text>
    </view>
    <view class="card-footer">
      <view class="footer-item">
        <text class="footer-icon">📚</text>
        <text class="footer-text">{{ question.category?.name || '未分类' }}</text>
      </view>
      <view class="footer-item">
        <text class="footer-icon">👁️</text>
        <text class="footer-text">{{ question.viewCount || 0 }}</text>
      </view>
      <view class="footer-item">
        <text class="footer-icon">❤️</text>
        <text class="footer-text">{{ question.likeCount || 0 }}</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  question: any
}>()

defineEmits(['tap'])

const typeClass = computed(() => {
  const typeMap: Record<number, string> = {
    1: 'single',
    2: 'multiple',
    3: 'judge',
    4: 'fill',
    5: 'short'
  }
  return typeMap[props.question.type] || ''
})

const typeText = computed(() => {
  const typeMap: Record<number, string> = {
    1: '单选题',
    2: '多选题',
    3: '判断题',
    4: '填空题',
    5: '简答题'
  }
  return typeMap[props.question.type] || '未知'
})

const difficultyClass = computed(() => {
  const map: Record<number, string> = {
    1: 'easy',
    2: 'medium',
    3: 'hard'
  }
  return map[props.question.difficulty] || ''
})

const difficultyText = computed(() => {
  const map: Record<number, string> = {
    1: '简单',
    2: '中等',
    3: '困难'
  }
  return map[props.question.difficulty] || '未知'
})
</script>

<style lang="scss" scoped>
.question-card {
  background-color: #ffffff;
  border-radius: 16rpx;
  padding: 24rpx;
  margin-bottom: 20rpx;
  box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.05);
}

.card-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 16rpx;
}

.type-tag {
  font-size: 22rpx;
  padding: 4rpx 16rpx;
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

.difficulty {
  font-size: 22rpx;

  &.easy {
    color: #67C23A;
  }

  &.medium {
    color: #E6A23C;
  }

  &.hard {
    color: #F56C6C;
  }
}

.card-body {
  margin-bottom: 16rpx;
}

.question-title {
  font-size: 30rpx;
  color: #303133;
  line-height: 1.5;
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
  overflow: hidden;
}

.card-footer {
  display: flex;
  align-items: center;
}

.footer-item {
  display: flex;
  align-items: center;
  margin-right: 30rpx;
}

.footer-icon {
  font-size: 24rpx;
  margin-right: 8rpx;
}

.footer-text {
  font-size: 24rpx;
  color: #909399;
}
</style>
