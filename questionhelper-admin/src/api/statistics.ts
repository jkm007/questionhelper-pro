import request from "@/utils/request";

export function getOverviewApi() {
  return request.get("/admin/statistics/dashboard");
}

export function getRetentionApi(params?: any) {
  return request.get("/admin/statistics/retention", { params });
}

export function listFunnelsApi(params?: any) {
  return request.get("/admin/statistics/funnels", { params });
}

export function createFunnelApi(data: any) {
  return request.post("/admin/statistics/funnels", data);
}

export function getFunnelStatsApi(id: number) {
  return request.get(`/admin/statistics/funnels/${id}/stats`);
}

export function listSegmentsApi(params?: any) {
  return request.get("/admin/statistics/segments", { params });
}

export function createSegmentApi(data: any) {
  return request.post("/admin/statistics/segments", data);
}

export function updateSegmentApi(id: number, data: any) {
  return request.put(`/admin/statistics/segments/${id}`, data);
}

export function deleteSegmentApi(id: number) {
  return request.delete(`/admin/statistics/segments/${id}`);
}

export function getPathAnalysisApi(params?: any) {
  return request.get("/admin/statistics/paths", { params });
}

export function compareDataApi(params?: any) {
  return request.get("/admin/statistics/compare", { params });
}

export function exportStatisticsApi(data: any) {
  return request.post("/admin/statistics/export", data, { responseType: "blob" });
}
