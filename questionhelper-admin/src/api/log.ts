import request from "@/utils/request";

export function listOperationLogsApi(params?: any) {
  return request.get("/admin/operation-logs", { params });
}

export function getOperationLogApi(id: number) {
  return request.get(`/admin/operation-logs/${id}`);
}

export function exportOperationLogsApi(params?: any) {
  return request.get("/admin/operation-logs/export", { params, responseType: "blob" });
}

export function cleanOperationLogsApi(days: number) {
  return request.post("/admin/operation-logs/clean", { days });
}

export function listLoginLogsApi(params?: any) {
  return request.get("/admin/login-logs", { params });
}

export function exportLoginLogsApi(params?: any) {
  return request.get("/admin/login-logs/export", { params, responseType: "blob" });
}

export function cleanLoginLogsApi(days: number) {
  return request.post("/admin/login-logs/clean", { days });
}

export function listSystemLogsApi(params?: any) {
  return request.get("/admin/logs/system", { params });
}

export function listErrorLogsApi(params?: any) {
  return request.get("/admin/logs/error", { params });
}

export function getLogStatsApi() {
  return request.get("/admin/logs/stats");
}

export function archiveLogsApi() {
  return request.post("/admin/logs/archive");
}
