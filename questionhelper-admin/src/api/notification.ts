import request from "@/utils/request";

// 通知模板
export function listTemplatesApi(params?: any) {
  return request.get("/admin/notifications/templates", { params });
}

export function createTemplateApi(data: any) {
  return request.post("/admin/notifications/templates", data);
}

export function updateTemplateApi(id: number, data: any) {
  return request.put(`/admin/notifications/templates/${id}`, data);
}

export function deleteTemplateApi(id: number) {
  return request.delete(`/admin/notifications/templates/${id}`);
}

// 通知渠道
export function listChannelsApi(params?: any) {
  return request.get("/admin/notifications/channels", { params });
}

export function updateChannelApi(id: number, data: any) {
  return request.put(`/admin/notifications/channels/${id}`, data);
}

// 告警规则
export function listAlertRulesApi(params?: any) {
  return request.get("/admin/alerts/rules", { params });
}

export function createAlertRuleApi(data: any) {
  return request.post("/admin/alerts/rules", data);
}

export function updateAlertRuleApi(id: number, data: any) {
  return request.put(`/admin/alerts/rules/${id}`, data);
}

// 告警记录
export function listAlertRecordsApi(params?: any) {
  return request.get("/admin/alerts/records", { params });
}
