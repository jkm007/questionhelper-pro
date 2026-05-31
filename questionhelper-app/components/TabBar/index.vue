<template>
  <view class="custom-tabbar">
    <view
      v-for="(item, index) in tabs"
      :key="index"
      class="tabbar-item"
      :class="{ 'tabbar-active': current === index }"
      @tap="switchTab(index, item.pagePath)"
    >
      <image
        class="tabbar-icon"
        :src="current === index ? item.selectedIconPath : item.iconPath"
        mode="aspectFit"
      />
      <text class="tabbar-text" :class="{ 'text-active': current === index }">{{ item.text }}</text>
    </view>
  </view>
</template>

<script setup lang="ts">
defineProps<{
  tabs: {
    pagePath: string
    text: string
    iconPath: string
    selectedIconPath: string
  }[]
  current: number
}>()

function switchTab(index: number, pagePath: string) {
  uni.switchTab({ url: `/${pagePath}` })
}
</script>

<style lang="scss" scoped>
.custom-tabbar {
  display: flex;
  background: #ffffff;
  box-shadow: 0 -2rpx 10rpx rgba(0, 0, 0, 0.05);
  padding-bottom: env(safe-area-inset-bottom);
}

.tabbar-item {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 8rpx 0 4rpx;
}

.tabbar-icon {
  width: 48rpx;
  height: 48rpx;
  margin-bottom: 4rpx;
}

.tabbar-text {
  font-size: 20rpx;
  color: #909399;

  &.text-active {
    color: #409eff;
    font-weight: 600;
  }
}
</style>
