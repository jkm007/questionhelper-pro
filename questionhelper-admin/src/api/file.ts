import request from "@/utils/request";

export function getStorageStatisticsApi() {
  return request.get("/admin/file/storage/statistics");
}

export function listHotlinkRulesApi(params?: any) {
  return request.get("/admin/file/hotlink-rules", { params });
}

export function createHotlinkRuleApi(data: any) {
  return request.post("/admin/file/hotlink-rules", data);
}

export function updateHotlinkRuleApi(id: number, data: any) {
  return request.put(`/admin/file/hotlink-rules/${id}`, data);
}

export function deleteHotlinkRuleApi(id: number) {
  return request.delete(`/admin/file/hotlink-rules/${id}`);
}

export function listWatermarkConfigsApi(params?: any) {
  return request.get("/admin/file/watermark-configs", { params });
}

export function createWatermarkConfigApi(data: any) {
  return request.post("/admin/file/watermark-configs", data);
}

export function updateWatermarkConfigApi(id: number, data: any) {
  return request.put(`/admin/file/watermark-configs/${id}`, data);
}

export function cleanupOrphanFilesApi() {
  return request.post("/admin/file/cleanup/orphan");
}

export function getCleanupLogsApi(params?: any) {
  return request.get("/admin/file/cleanup/logs", { params });
}
