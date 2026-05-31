<template>
  <view class="class-card" @tap="$emit('tap')">
    <image class="card-cover" :src="classInfo.coverImage || '/static/default-class.png'" mode="aspectFill" />
    <view class="card-body">
      <text class="card-name">{{ classInfo.name }}</text>
      <view class="card-meta">
        <text class="meta-creator">{{ classInfo.creatorName }}</text>
        <text class="meta-members">{{ classInfo.memberCount }}人</text>
      </view>
      <text v-if="classInfo.description" class="card-desc">{{ classInfo.description }}</text>
    </view>
  </view>
</template>

<script setup lang="ts">
defineProps<{
  classInfo: {
    id: number | string
    name: string
    coverImage?: string
    creatorName?: string
    memberCount?: number
    description?: string
  }
}>()

defineEmits(['tap'])
</script>

<style lang="scss" scoped>
.class-card {
  display: flex;
  background: #ffffff;
  border-radius: 16rpx;
  overflow: hidden;
  box-shadow: 0 2rpx 12rpx rgba(0, 0, 0, 0.04);
}

.card-cover {
  width: 200rpx;
  height: 200rpx;
  flex-shrink: 0;
}

.card-body {
  flex: 1;
  padding: 20rpx;
  min-width: 0;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.card-name {
  font-size: 30rpx;
  font-weight: 600;
  color: #303133;
  margin-bottom: 8rpx;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.card-meta {
  display: flex;
  align-items: center;
  gap: 16rpx;
  margin-bottom: 8rpx;
}

.meta-creator {
  font-size: 24rpx;
  color: #909399;
}

.meta-members {
  font-size: 24rpx;
  color: #409eff;
}

.card-desc {
  font-size: 24rpx;
  color: #909399;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
