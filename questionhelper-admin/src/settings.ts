/**
 * 应用配置
 */

export const appConfig = {
  name: "题小助",
  title: "题小助管理后台",
  tenantEnabled: false,
} as const;

export const defaults = {
  theme: "light",
  themeColor: "#409EFF",
  sidebarColorScheme: "classic-blue",
  layout: "left",
  size: "default",
  language: "zh-cn",
  showTagsView: true,
  showAppLogo: true,
  showWatermark: false,
  pageSwitchingAnimation: "fade-slide",
  showSettings: true,
} as const;

export const themeColorPresets = [
  "#409EFF",
  "#67C23A",
  "#E6A23C",
  "#F56C6C",
  "#909399",
] as const;
