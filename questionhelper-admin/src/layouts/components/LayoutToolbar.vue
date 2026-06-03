<template>
  <div :class="['navbar-actions', navbarActionsClass]">
    <!-- 桌面端工具项 -->
    <template v-if="isDesktop">
      <!-- 全屏 -->
      <div class="navbar-actions__item">
        <Fullscreen />
      </div>

      <!-- 布局大小 -->
      <div class="navbar-actions__item">
        <SizeSelect />
      </div>

      <!-- 语言选择 -->
      <div class="navbar-actions__item">
        <LangSelect />
      </div>
    </template>

    <!-- 用户菜单 -->
    <div class="navbar-actions__item">
      <el-dropdown trigger="click">
        <div class="user-profile">
          <el-avatar :size="28" :src="userStore.userInfo?.avatar || '/images/default-avatar.png'" />
          <span class="user-profile__name">{{ userStore.userInfo?.nickname || userStore.userInfo?.username || "管理员" }}</span>
        </div>
        <template #dropdown>
          <el-dropdown-menu>
            <el-dropdown-item @click="router.push('/profile')">
              个人中心
            </el-dropdown-item>
            <el-dropdown-item divided @click="logout">
              退出登录
            </el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
    </div>

    <!-- 系统设置 -->
    <div v-if="defaults.showSettings" class="navbar-actions__item" @click="settingsStore.settingsVisible = true">
      <el-icon :size="18"><Setting /></el-icon>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRoute, useRouter } from "vue-router";
import { defaults } from "@/settings";
import { useAppStore, useSettingsStore, useUserStore } from "@/stores";
import { Setting } from "@element-plus/icons-vue";

// 导入子组件
import Fullscreen from "@/components/Fullscreen/index.vue";
import SizeSelect from "@/components/SizeSelect/index.vue";
import LangSelect from "@/components/LangSelect/index.vue";

const route = useRoute();
const router = useRouter();
const appStore = useAppStore();
const settingsStore = useSettingsStore();
const userStore = useUserStore();

const isDesktop = computed(() => appStore.device === "desktop");

const navbarActionsClass = computed(() => ({
  "navbar-actions--mobile": !isDesktop.value,
}));

async function logout() {
  await userStore.logout();
  router.push("/login");
}
</script>

<style lang="scss" scoped>
.navbar-actions {
  display: flex;
  align-items: center;
  gap: 4px;

  &__item {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 36px;
    height: 36px;
    border-radius: 6px;
    cursor: pointer;
    transition: background-color 0.2s;

    &:hover {
      background-color: var(--el-fill-color-light);
    }
  }
}

.user-profile {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  padding: 0 8px;

  &__name {
    font-size: 14px;
    color: var(--el-text-color-primary);
  }
}
</style>
