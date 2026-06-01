import request from "@/utils/request";

// 题目管理
export function listQuestionsApi(params?: any) {
  return request.get("/admin/questions", { params });
}

export function getQuestionApi(id: number) {
  return request.get(`/admin/questions/${id}`);
}

export function createQuestionApi(data: any) {
  return request.post("/admin/questions", data);
}

export function updateQuestionApi(id: number, data: any) {
  return request.put(`/admin/questions/${id}`, data);
}

export function deleteQuestionApi(id: number) {
  return request.delete(`/admin/questions/${id}`);
}

export function updateQuestionStatusApi(id: number, status: number) {
  return request.put(`/admin/questions/${id}/status`, { status });
}

export function importQuestionsApi(file: File) {
  const formData = new FormData();
  formData.append("file", file);
  return request.post("/admin/questions/import", formData);
}

export function exportQuestionsApi(params?: any) {
  return request.get("/admin/questions/export", { params, responseType: "blob" });
}

// 批量操作
export function batchPublishApi(ids: number[]) {
  return request.post("/admin/questions/batch/publish", { ids });
}

export function batchArchiveApi(ids: number[]) {
  return request.post("/admin/questions/batch/archive", { ids });
}

export function batchDeleteApi(ids: number[]) {
  return request.post("/admin/questions/batch/delete", { ids });
}

// 版本管理
export function listVersionsApi(questionId: number) {
  return request.get(`/admin/questions/${questionId}/versions`);
}

export function rollbackVersionApi(questionId: number, version: number) {
  return request.post(`/admin/questions/${questionId}/versions/${version}/rollback`);
}

// 分类管理
export function listCategoriesApi(params?: any) {
  return request.get("/admin/categories", { params });
}

export function createCategoryApi(data: any) {
  return request.post("/admin/categories", data);
}

export function updateCategoryApi(id: number, data: any) {
  return request.put(`/admin/categories/${id}`, data);
}

export function deleteCategoryApi(id: number) {
  return request.delete(`/admin/categories/${id}`);
}

// 知识点管理
export function listKnowledgePointsApi(params?: any) {
  return request.get("/admin/knowledge-points", { params });
}

export function createKnowledgePointApi(data: any) {
  return request.post("/admin/knowledge-points", data);
}

export function updateKnowledgePointApi(id: number, data: any) {
  return request.put(`/admin/knowledge-points/${id}`, data);
}

export function deleteKnowledgePointApi(id: number) {
  return request.delete(`/admin/knowledge-points/${id}`);
}

// 题目审核
export function listReviewsApi(params?: any) {
  return request.get("/admin/reviews", { params });
}

export function approveReviewApi(id: number) {
  return request.post(`/admin/reviews/${id}/approve`);
}

export function rejectReviewApi(id: number, reason: string) {
  return request.post(`/admin/reviews/${id}/reject`, { reason });
}

// 分享管理
export function listSharesApi(params?: any) {
  return request.get("/admin/shares", { params });
}

export function revokeShareApi(id: number) {
  return request.delete(`/admin/shares/${id}`);
}

// 敏感词管理
export function listSensitiveWordsApi(params?: any) {
  return request.get("/admin/sensitive-words", { params });
}

export function createSensitiveWordApi(data: any) {
  return request.post("/admin/sensitive-words", data);
}

export function deleteSensitiveWordApi(id: number) {
  return request.delete(`/admin/sensitive-words/${id}`);
}

export function importSensitiveWordsApi(file: File) {
  const formData = new FormData();
  formData.append("file", file);
  return request.post("/admin/sensitive-words/import", formData);
}

// 题目统计
export function getQuestionStatsApi() {
  return request.get("/admin/questions/stats");
}
