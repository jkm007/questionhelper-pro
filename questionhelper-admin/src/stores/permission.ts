import type { RouteRecordRaw } from "vue-router";
import { constantRoutes } from "@/router";
import { store } from "@/stores";
import router from "@/router";
import { getRoutesApi, getButtonsApi } from "@/api/menu";

const modules = import.meta.glob("../views/**/**.vue");
const Layout = () => import("../layouts/index.vue");

function resolveViewComponent(componentPath: string) {
  const normalized = componentPath
    .trim()
    .replace(/^\/+/, "")
    .replace(/\.vue$/i, "");
  return (
    modules[`../views/${normalized}.vue`] ||
    modules[`../views/${normalized}/index.vue`] ||
    modules[`../views/error/404.vue`]
  );
}

export const usePermissionStore = defineStore("permission", () => {
  const routes = ref<RouteRecordRaw[]>([]);
  const mixLayoutSideMenus = ref<RouteRecordRaw[]>([]);
  const isRouteGenerated = ref(false);
  const buttons = ref<string[]>([]);

  async function generateRoutes(): Promise<RouteRecordRaw[]> {
    try {
      const routeData: any = await getRoutesApi();
      const buttonData: any = await getButtonsApi();

      const dynamicRoutes = transformRoutes(routeData);
      buttons.value = buttonData || [];
      routes.value = [...constantRoutes, ...dynamicRoutes];
      isRouteGenerated.value = true;

      return dynamicRoutes;
    } catch (error) {
      isRouteGenerated.value = false;
      throw error;
    }
  }

  const setMixLayoutSideMenus = (parentPath: string) => {
    const parentMenu = routes.value.find((item: RouteRecordRaw) => item.path === parentPath);
    mixLayoutSideMenus.value = parentMenu?.children || [];
  };

  const resetRouter = () => {
    const constantRouteNames = new Set(constantRoutes.map((route) => route.name).filter(Boolean));
    routes.value.forEach((route: RouteRecordRaw) => {
      if (route.name && !constantRouteNames.has(route.name)) {
        router.removeRoute(route.name);
      }
    });

    routes.value = [...constantRoutes];
    mixLayoutSideMenus.value = [];
    isRouteGenerated.value = false;
    buttons.value = [];
  };

  return {
    routes,
    mixLayoutSideMenus,
    isRouteGenerated,
    buttons,
    generateRoutes,
    setMixLayoutSideMenus,
    resetRouter,
  };
});

const transformRoutes = (routes: any[], isTopLevel: boolean = true): RouteRecordRaw[] => {
  return routes.map((route) => {
    const { component, children, ...args } = route;

    const processedComponent = isTopLevel || component !== "Layout" ? component : undefined;

    const normalizedRoute = { ...args } as RouteRecordRaw;

    if (!processedComponent) {
      normalizedRoute.component = undefined;
    } else {
      normalizedRoute.component =
        processedComponent === "Layout" ? Layout : resolveViewComponent(processedComponent);
    }

    if (children && children.length > 0) {
      normalizedRoute.children = transformRoutes(children, false);
    }

    return normalizedRoute;
  });
};

export function usePermissionStoreHook() {
  return usePermissionStore(store);
}
