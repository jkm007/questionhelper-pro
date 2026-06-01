import { defineStore } from "pinia";
import { ref } from "vue";
import type { RouteRecordRaw } from "vue-router";
import { getRoutesApi, getButtonsApi } from "@/api/menu";

export const usePermissionStore = defineStore("permission", () => {
  const routes = ref<RouteRecordRaw[]>([]);
  const buttons = ref<string[]>([]);
  const isRouteGenerated = ref(false);

  async function generateRoutes() {
    const routeData: any = await getRoutesApi();
    const buttonData: any = await getButtonsApi();

    routes.value = transformRoutes(routeData);
    buttons.value = buttonData || [];
    isRouteGenerated.value = true;
    return routes.value;
  }

  function transformRoutes(routeList: any[]): RouteRecordRaw[] {
    return routeList.map((item) => {
      const route: RouteRecordRaw = {
        path: item.path,
        name: item.name,
        component: resolveComponent(item.component),
        redirect: item.redirect,
        meta: item.meta,
        children: item.children ? transformRoutes(item.children) : [],
      };
      return route;
    });
  }

  function resolveComponent(component: string) {
    if (!component) return undefined;
    if (component === "Layout") return () => import("@/layouts/BaseLayout.vue");
    return () => import(`@/views/${component}.vue`);
  }

  function resetRoutes() {
    routes.value = [];
    buttons.value = [];
    isRouteGenerated.value = false;
  }

  return { routes, buttons, isRouteGenerated, generateRoutes, resetRoutes };
});
