<template>
  <div class="tags-view-container">
    <el-scrollbar class="tags-scrollbar">
      <div class="tags-view-wrapper">
        <router-link
          v-for="tag in visitedViews"
          :key="tag.path"
          :to="{ path: tag.path, query: tag.query }"
          :class="['tags-view-item', { active: isActive(tag) }]"
          @click.middle="closeTag(tag)"
        >
          {{ tag.title }}
          <el-icon v-if="!isAffix(tag)" class="tag-close" @click.prevent.stop="closeTag(tag)">
            <Close />
          </el-icon>
        </router-link>
      </div>
    </el-scrollbar>
  </div>
</template>

<script setup lang="ts">
import { computed, watch } from "vue";
import { useRoute, useRouter } from "vue-router";
import { Close } from "@element-plus/icons-vue";

interface TagView {
  path: string;
  title: string;
  query?: any;
  affix?: boolean;
}

const route = useRoute();
const router = useRouter();

const visitedViews = ref<TagView[]>([]);

function isActive(tag: TagView) {
  return tag.path === route.path;
}

function isAffix(tag: TagView) {
  return tag.affix;
}

function addTag() {
  const { path, meta, query } = route;
  if (meta?.hidden) return;
  if (visitedViews.value.some((v) => v.path === path)) return;
  visitedViews.value.push({
    path,
    title: (meta?.title as string) || "未命名",
    query,
  });
}

function closeTag(tag: TagView) {
  const index = visitedViews.value.findIndex((v) => v.path === tag.path);
  if (index === -1) return;
  visitedViews.value.splice(index, 1);
  if (isActive(tag)) {
    const last = visitedViews.value[visitedViews.value.length - 1];
    if (last) {
      router.push(last.path);
    } else {
      router.push("/");
    }
  }
}

watch(
  () => route.path,
  () => addTag(),
  { immediate: true }
);
</script>

<style lang="scss" scoped>
.tags-view-container {
  height: 34px;
  width: 100%;
  background: #fff;
  border-bottom: 1px solid #d8dce5;
  box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.12);
}

.tags-view-wrapper {
  display: flex;
  align-items: center;
  height: 34px;
  padding: 0 8px;
}

.tags-view-item {
  display: inline-flex;
  align-items: center;
  height: 26px;
  padding: 0 10px;
  margin: 0 2px;
  font-size: 12px;
  color: #495060;
  background: #fff;
  border: 1px solid #d8dce5;
  border-radius: 3px;
  text-decoration: none;
  white-space: nowrap;
  cursor: pointer;

  &:hover {
    color: #409eff;
  }

  &.active {
    background-color: #409eff;
    color: #fff;
    border-color: #409eff;

    &::before {
      content: "";
      display: inline-block;
      width: 8px;
      height: 8px;
      border-radius: 50%;
      background: #fff;
      margin-right: 6px;
    }
  }

  .tag-close {
    margin-left: 4px;
    font-size: 12px;
    border-radius: 50%;

    &:hover {
      background-color: rgba(0, 0, 0, 0.15);
    }
  }
}
</style>
