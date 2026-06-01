import { defineStore } from "pinia";
import { ref } from "vue";

export const useSettingsStore = defineStore("settings", () => {
  const theme = ref<string>("light");
  const primaryColor = ref<string>("#409eff");
  const showTagsView = ref(true);
  const showBreadcrumb = ref(true);

  function setTheme(val: string) {
    theme.value = val;
    document.documentElement.setAttribute("data-theme", val);
  }

  function setPrimaryColor(val: string) {
    primaryColor.value = val;
    document.documentElement.style.setProperty("--el-color-primary", val);
  }

  return { theme, primaryColor, showTagsView, showBreadcrumb, setTheme, setPrimaryColor };
});
