import request from "@/utils/request";

// 基础设置
export function getSettingsApi() {
  return request.get("/admin/settings");
}

export function updateSettingsApi(data: any) {
  return request.put("/admin/settings", data);
}

export function getClassSettingsApi() {
  return request.get("/admin/settings/class");
}

export function updateClassSettingsApi(data: any) {
  return request.put("/admin/settings/class", data);
}

export function getResourceSettingsApi() {
  return request.get("/admin/settings/resource");
}

export function updateResourceSettingsApi(data: any) {
  return request.put("/admin/settings/resource", data);
}

// 安全设置
export function getSecurityConfigsApi() {
  return request.get("/admin/security");
}

export function updateSecurityConfigsApi(data: any) {
  return request.put("/admin/security", data);
}

// 功能开关
export function listFeatureFlagsApi() {
  return request.get("/admin/features");
}

export function updateFeatureFlagApi(key: string, data: any) {
  return request.put(`/admin/features/${key}`, data);
}

// 主题设置
export function getThemeConfigApi() {
  return request.get("/admin/theme");
}

export function updateThemeConfigApi(data: any) {
  return request.put("/admin/theme", data);
}

// 存储配置
export function listStorageConfigsApi() {
  return request.get("/admin/storage");
}

export function createStorageConfigApi(data: any) {
  return request.post("/admin/storage", data);
}

export function updateStorageConfigApi(id: number, data: any) {
  return request.put(`/admin/storage/${id}`, data);
}

// 邮件配置
export function getEmailConfigApi() {
  return request.get("/admin/email/config");
}

export function updateEmailConfigApi(data: any) {
  return request.put("/admin/email/config", data);
}

export function listEmailTemplatesApi() {
  return request.get("/admin/email/templates");
}

export function createEmailTemplateApi(data: any) {
  return request.post("/admin/email/templates", data);
}

// 短信配置
export function getSmsConfigApi() {
  return request.get("/admin/sms/config");
}

export function updateSmsConfigApi(data: any) {
  return request.put("/admin/sms/config", data);
}

export function listSmsTemplatesApi() {
  return request.get("/admin/sms/templates");
}

export function createSmsTemplateApi(data: any) {
  return request.post("/admin/sms/templates", data);
}

// 缓存管理
export function getCacheStatsApi() {
  return request.get("/admin/cache/stats");
}

export function clearCacheApi() {
  return request.post("/admin/cache/clear");
}
