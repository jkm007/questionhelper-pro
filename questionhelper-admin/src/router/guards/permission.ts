import router from "@/router";
import { useUserStore } from "@/stores/user";
import { usePermissionStore } from "@/stores/permission";
import { getToken } from "@/utils/auth";

const whiteList = ["/login", "/401", "/404"];

router.beforeEach(async (to, _from, next) => {
  const token = getToken();

  if (token) {
    if (to.path === "/login") {
      next({ path: "/" });
      return;
    }

    const userStore = useUserStore();
    const permissionStore = usePermissionStore();

    if (!userStore.userInfo) {
      try {
        await userStore.getUserInfo();
        const routes = await permissionStore.generateRoutes();
        routes.forEach((route) => router.addRoute(route));
        next({ ...to, replace: true });
      } catch {
        await userStore.logout();
        next(`/login?redirect=${to.path}`);
      }
      return;
    }

    next();
  } else {
    if (whiteList.includes(to.path)) {
      next();
    } else {
      next(`/login?redirect=${to.path}`);
    }
  }
});
