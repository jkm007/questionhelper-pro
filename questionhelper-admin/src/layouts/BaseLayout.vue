<template>
  <div class="app-wrapper">
    <el-container>
      <el-aside width="210px" class="sidebar">
        <div class="logo">
          <h3>题小助</h3>
        </div>
        <el-menu
          :default-active="activeMenu"
          :collapse="appStore.sidebarCollapsed"
          router
          background-color="#304156"
          text-color="#bfcbd9"
          active-text-color="#409eff"
        >
          <template v-for="route in menuRoutes" :key="route.path">
            <el-sub-menu v-if="route.children?.length" :index="route.path">
              <template #title>
                <el-icon v-if="route.meta?.icon"><component :is="route.meta.icon" /></el-icon>
                <span>{{ route.meta?.title }}</span>
              </template>
              <el-menu-item
                v-for="child in route.children"
                :key="child.path"
                :index="`${route.path}/${child.path}`"
              >
                <span>{{ child.meta?.title }}</span>
              </el-menu-item>
            </el-sub-menu>
            <el-menu-item v-else :index="route.path">
              <el-icon v-if="route.meta?.icon"><component :is="route.meta.icon" /></el-icon>
              <span>{{ route.meta?.title }}</span>
            </el-menu-item>
          </template>
        </el-menu>
      </el-aside>

      <el-container>
        <el-header class="navbar">
          <div class="navbar-left">
            <el-icon class="collapse-btn" @click="appStore.toggleSidebar">
              <Fold v-if="!appStore.sidebarCollapsed" />
              <Expand v-else />
            </el-icon>
            <Breadcrumb />
          </div>
          <div class="navbar-right">
            <el-dropdown>
              <span class="user-info">
                {{ userStore.userInfo?.nickname || "管理员" }}
                <el-icon><ArrowDown /></el-icon>
              </span>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item @click="$router.push('/profile')">个人中心</el-dropdown-item>
                  <el-dropdown-item divided @click="handleLogout">退出登录</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
        </el-header>

        <el-main>
          <router-view />
        </el-main>
      </el-container>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import { computed } from "vue";
import { useRoute, useRouter } from "vue-router";
import { useAppStore } from "@/stores/app";
import { useUserStore } from "@/stores/user";
import { usePermissionStore } from "@/stores/permission";
import Breadcrumb from "@/components/Breadcrumb/index.vue";
import { Fold, Expand, ArrowDown } from "@element-plus/icons-vue";

const route = useRoute();
const router = useRouter();
const appStore = useAppStore();
const userStore = useUserStore();
const permissionStore = usePermissionStore();

const activeMenu = computed(() => route.path);
const menuRoutes = computed(() =>
  permissionStore.routes.filter((r) => !r.meta?.hidden)
);

async function handleLogout() {
  await userStore.logout();
  router.push("/login");
}
</script>

<style scoped>
.app-wrapper {
  height: 100vh;
}
.sidebar {
  background: #304156;
  overflow-y: auto;
}
.logo {
  height: 50px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
}
.navbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1px solid #eee;
  background: #fff;
}
.navbar-left {
  display: flex;
  align-items: center;
  gap: 12px;
}
.collapse-btn {
  cursor: pointer;
  font-size: 20px;
}
.user-info {
  display: flex;
  align-items: center;
  cursor: pointer;
  gap: 4px;
}
</style>
