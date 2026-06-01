import { fileURLToPath, URL } from "node:url";
import { defineConfig, loadEnv } from "vite";
import vue from "@vitejs/plugin-vue";
import UnoCSS from "unocss/vite";
import AutoImport from "unplugin-auto-import/vite";
import Components from "unplugin-vue-components/vite";
import { ElementPlusResolver } from "unplugin-vue-components/resolvers";

export default defineConfig(({ mode }) => {
  const env = loadEnv(mode, process.cwd());
  const { VITE_APP_BASE_API, VITE_APP_BASE_URL } = env;

  return {
    base: VITE_APP_BASE_URL || "/",
    plugins: [
      vue(),
      UnoCSS(),
      AutoImport({
        imports: ["vue", "vue-router", "pinia", "@vueuse/core"],
        resolvers: [ElementPlusResolver()],
        dts: "src/auto-imports.d.ts",
      }),
      Components({
        resolvers: [ElementPlusResolver()],
        dts: "src/components.d.ts",
      }),
    ],
    resolve: {
      alias: {
        "@": fileURLToPath(new URL("./src", import.meta.url)),
      },
    },
    server: {
      port: 3000,
      proxy: {
        "/api": {
          target: VITE_APP_BASE_API || "http://localhost:8080",
          changeOrigin: true,
        },
      },
    },
  };
});
