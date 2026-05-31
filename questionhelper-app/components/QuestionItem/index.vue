<template>
  <view class="question-item" @tap="$emit('tap')">
    <view class="item-header">
      <view class="type-tag" :class="typeClass">{{ typeText }}</view>
      <view class="difficulty-tag" :class="difficultyClass">{{ difficultyText }}</view>
    </view>
    <text class="item-title">{{ question.title }}</text>
    <view class="item-footer">
      <text v-if="question.categoryName" class="footer-category">{{ question.categoryName }}</text>
      <view class="footer-stats">
        <text class="stat">👁 {{ question.viewCount || 0 }}</text>
        <text class="stat">👍 {{ question.likeCount || 0 }}</text>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  question: {
    id: number | string
    title: string
    type?: number
    difficulty?: number
    categoryName?: string
    viewCount?: number
    likeCount?: number
  }
}>()

defineEmits(['tap'])

const typeClass = computed(() => {
  const map: Record<number, string> = { 1: 'type-single', 2: 'type-multi', 3: 'type-judge', 4: 'type-fill', 5: 'type-short' }
  return map[props.question.type || 0] || 'type-single'
})

const typeText = computed(() => {
  const map: Record<number, string> = { 1: '单选', 2: '多选', 3: '判断', 4: '填空', 5: '简答' }
  return map[props.question.type || 0] || '单选'
})

const difficultyClass = computed(() => {
  const map: Record<number, string> = { 1: 'diff-easy', 2: 'diff-medium', 3: 'diff-hard' }
  return map[props.question.difficulty || 0] || 'diff-easy'
})

const difficultyText = computed(() => {
  const map: Record<number, string> = { 1: '简单', 2: '中等', 3: '困难' }
  return map[props.question.difficulty || 0] || '简单'
})
</script>

<style lang="scss" scoped>
.question-item {
  background: #ffffff;
  border-radius: 12rpx;
  padding: 24rpx;
  box-shadow: 0 2rpx 8rpx rgba(0, 0, 0, 0.04);
}

.item-header {
  display: flex;
  gap: 12rpx;
  margin-bottom: 12rpx;
}

.type-tag, .difficulty-tag {
  font-size: 20rpx;
  padding: 2rpx 12rpx;
  border-radius: 6rpx;
}

.type-single { background: #ecf5ff; color: #409eff; }
.type-multi { background: #fdf6ec; color: #e6a23c; }
.type-judge { background: #f0f9eb; color: #67c23a; }
.type-fill { background: #fef0f0; color: #f56c6c; }
.type-short { background: #f4f4f5; color: #909399; }

.diff-easy { background: #f0f9eb; color: #67c23a; }
.diff-medium { background: #fdf6ec; color: #e6a23c; }
.diff-hard { background: #fef0f0; color: #f56c6c; }

.item-title {
  font-size: 28rpx;
  color: #303133;
  line-height: 1.5;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  margin-bottom: 12rpx;
}

.item-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.footer-category {
  font-size: 22rpx;
  color: #909399;
}

.footer-stats {
  display: flex;
  gap: 16rpx;
}

.stat {
  font-size: 22rpx;
  color: #c0c4cc;
}
</style>
