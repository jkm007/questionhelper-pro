import type { RouteRecordRaw } from "vue-router";
import router from "@/router";
import { usePermissionStore, useUserStore } from "@/stores";

/**
 * 路由权限守卫
 */
export function setupPermissionGuard() {
  const whiteList = ["/login"];

  router.beforeEach(async (to, _from) => {
    try {
      const userStore = useUserStore();
      const isLoggedIn = userStore.isLoggedIn();

      // 未登录处理
      if (!isLoggedIn) {
        if (whiteList.includes(to.path)) {
          return;
        }
        return `/login?redirect=${encodeURIComponent(to.fullPath)}`;
      }

      // 已登录访问登录页，重定向到首页
      if (to.path === "/login") {
        return { path: "/" };
      }

      const permissionStore = usePermissionStore();

      // 动态路由生成
      if (!permissionStore.isRouteGenerated) {
        if (!userStore.userInfo?.roles?.length) {
          await userStore.getUserInfo();
        }

        const dynamicRoutes = await permissionStore.generateRoutes();
        dynamicRoutes.forEach((route: RouteRecordRaw) => {
          router.addRoute(route);
        });

        return { ...to, replace: true };
      }

      // 路由 404 检查
      if (to.matched.length === 0) {
        if (_from.path === "/login") {
          return { path: "/", replace: true };
        }
        return "/404";
      }
    } catch (error) {
      console.error("Route guard error:", error);
      userStore.resetUserState();
      return "/login";
    }
  });
}
