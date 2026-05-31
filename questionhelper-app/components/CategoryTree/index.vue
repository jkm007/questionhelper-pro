<template>
  <view class="category-tree">
    <view
      class="tree-item root-item"
      :class="{ active: selectedId === 0 }"
      @tap="handleSelect(0, '全部')"
    >
      <text class="item-text">全部分类</text>
    </view>
    <view v-for="item in categories" :key="item.id" class="tree-group">
      <view
        class="tree-item parent-item"
        :class="{ active: selectedId === item.id, expanded: expandedIds.has(item.id) }"
        @tap="toggleExpand(item.id)"
      >
        <text class="expand-icon">{{ expandedIds.has(item.id) ? '▼' : '▶' }}</text>
        <text class="item-text">{{ item.name }}</text>
      </view>
      <view v-if="expandedIds.has(item.id) && item.children" class="tree-children">
        <view
          v-for="child in item.children"
          :key="child.id"
          class="tree-item child-item"
          :class="{ active: selectedId === child.id }"
          @tap="handleSelect(child.id, child.name)"
        >
          <text class="item-text">{{ child.name }}</text>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup lang="ts">
import { ref } from 'vue'

defineProps<{
  categories: {
    id: number
    name: string
    children?: { id: number; name: string }[]
  }[]
  selectedId?: number
}>()

const emit = defineEmits(['select'])

const expandedIds = ref(new Set<number>())

function toggleExpand(id: number) {
  if (expandedIds.value.has(id)) {
    expandedIds.value.delete(id)
  } else {
    expandedIds.value.add(id)
  }
}

function handleSelect(id: number, name: string) {
  emit('select', { id, name })
}
</script>

<style lang="scss" scoped>
.category-tree {
  padding: 12rpx 0;
}

.tree-item {
  display: flex;
  align-items: center;
  padding: 20rpx 24rpx;

  &.active {
    background: #ecf5ff;
    .item-text { color: #409eff; font-weight: 600; }
  }
}

.parent-item {
  font-weight: 600;
}

.expand-icon {
  font-size: 20rpx;
  color: #c0c4cc;
  margin-right: 12rpx;
  transition: transform 0.2s;
}

.item-text {
  font-size: 28rpx;
  color: #606266;
}

.tree-children {
  padding-left: 32rpx;
}

.child-item {
  padding-left: 48rpx;
}
</style>
