import type { App } from "vue";
import { createPinia } from "pinia";

export const store = createPinia();

export function setupStore(app: App) {
  app.use(store);
}

// 导出所有 store
export { useAppStore } from "./app";
export { useUserStore } from "./user";
export { usePermissionStore } from "./permission";
export { useSettingsStore } from "./settings";
export { useTagsViewStore } from "./tags-view";
export { useDictStore } from "./dict";
