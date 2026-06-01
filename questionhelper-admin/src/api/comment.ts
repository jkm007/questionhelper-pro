import request from "@/utils/request";

// 评论管理
export function listCommentsApi(params?: any) {
  return request.get("/admin/comments", { params });
}

export function batchAuditApi(ids: number[], status: number) {
  return request.put("/admin/comments/batch-audit", { ids, status });
}

export function batchDeleteCommentsApi(ids: number[]) {
  return request.post("/admin/comments/batch-delete", { ids });
}

export function pinCommentApi(id: number) {
  return request.post(`/admin/comments/${id}/pin`);
}

export function unpinCommentApi(id: number) {
  return request.delete(`/admin/comments/${id}/pin`);
}

export function featureCommentApi(id: number) {
  return request.post(`/admin/comments/${id}/featured`);
}

export function getCommentStatsApi() {
  return request.get("/admin/comments/stats");
}

export function exportCommentsApi(params?: any) {
  return request.get("/admin/comments/export", { params, responseType: "blob" });
}

// 举报管理
export function listReportsApi(params?: any) {
  return request.get("/admin/comments/reports", { params });
}

export function handleReportApi(id: number, data: any) {
  return request.put(`/admin/comments/reports/${id}`, data);
}

// 黑名单
export function listBlacklistApi(params?: any) {
  return request.get("/admin/comments/blacklist", { params });
}

export function addBlacklistApi(data: any) {
  return request.post("/admin/comments/blacklist", data);
}

export function removeBlacklistApi(id: number) {
  return request.delete(`/admin/comments/blacklist/${id}`);
}

// 审核规则
export function listAuditRulesApi(params?: any) {
  return request.get("/admin/comments/audit-rules", { params });
}

export function createAuditRuleApi(data: any) {
  return request.post("/admin/comments/audit-rules", data);
}

export function updateAuditRuleApi(id: number, data: any) {
  return request.put(`/admin/comments/audit-rules/${id}`, data);
}

export function deleteAuditRuleApi(id: number) {
  return request.delete(`/admin/comments/audit-rules/${id}`);
}
