import { createApp } from "vue";
import App from "./App.vue";
import { setupStore } from "@/stores";
import { setupRouter } from "@/router";
import { setupI18n } from "@/lang";
import { setupPermissionGuard } from "@/router/guards/permission";

// 样式
import "element-plus/dist/index.css";
import "element-plus/theme-chalk/dark/css-vars.css";
import "@/styles/index.scss";
import "virtual:uno.css";

// 全局注册 Element Plus 图标
import * as ElementPlusIconsVue from "@element-plus/icons-vue";

const app = createApp(App);

// 注册所有图标
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component);
}

// 初始化 Pinia
setupStore(app);

// 初始化 i18n
setupI18n(app);

// 初始化 Router
setupRouter(app);

// 初始化路由守卫
setupPermissionGuard();

// 挂载
app.mount("#app");
