import request from "@/utils/request";

// 考试管理
export function listExamsApi(params?: any) {
  return request.get("/admin/exams", { params });
}

export function getExamApi(id: number) {
  return request.get(`/admin/exams/${id}`);
}

export function createExamApi(data: any) {
  return request.post("/admin/exams", data);
}

export function updateExamApi(id: number, data: any) {
  return request.put(`/admin/exams/${id}`, data);
}

export function deleteExamApi(id: number) {
  return request.delete(`/admin/exams/${id}`);
}

export function publishExamApi(id: number) {
  return request.put(`/admin/exams/${id}/publish`);
}

export function closeExamApi(id: number) {
  return request.put(`/admin/exams/${id}/close`);
}

export function extendExamApi(id: number, data: any) {
  return request.post(`/admin/exams/${id}/extend`, data);
}

export function getExamMonitorApi(id: number) {
  return request.get(`/admin/exams/${id}/monitor`);
}

export function getExamAnalysisApi(id: number) {
  return request.get(`/admin/exams/${id}/analysis`);
}

export function getExamScoresApi(id: number, params?: any) {
  return request.get(`/admin/exams/${id}/scores`, { params });
}

export function exportExamScoresApi(id: number) {
  return request.get(`/admin/exams/${id}/scores/export`, { responseType: "blob" });
}

// 试卷管理
export function listPapersApi(params?: any) {
  return request.get("/admin/papers", { params });
}

export function getPaperApi(id: number) {
  return request.get(`/admin/papers/${id}`);
}

export function createPaperApi(data: any) {
  return request.post("/admin/papers", data);
}

export function updatePaperApi(id: number, data: any) {
  return request.put(`/admin/papers/${id}`, data);
}

export function deletePaperApi(id: number) {
  return request.delete(`/admin/papers/${id}`);
}

export function previewPaperApi(id: number) {
  return request.get(`/admin/papers/${id}/preview`);
}

export function copyPaperApi(id: number) {
  return request.post(`/admin/papers/${id}/copy`);
}

export function publishPaperApi(id: number) {
  return request.put(`/admin/papers/${id}/publish`);
}

export function exportPaperApi(id: number) {
  return request.get(`/admin/papers/${id}/export`, { responseType: "blob" });
}

export function getPaperStatsApi(id: number) {
  return request.get(`/admin/papers/${id}/stats`);
}

// 模板管理
export function listTemplatesApi(params?: any) {
  return request.get("/admin/templates", { params });
}

export function createFromTemplateApi(id: number) {
  return request.post(`/admin/templates/${id}/create`);
}

// 成绩管理
export function listScoresApi(params?: any) {
  return request.get("/admin/scores", { params });
}

export function getScoreApi(id: number) {
  return request.get(`/admin/scores/${id}`);
}

export function getScoreAnalysisApi(params?: any) {
  return request.get("/admin/scores/analysis", { params });
}

// 成绩复核
export function listScoreReviewsApi(params?: any) {
  return request.get("/admin/score-reviews", { params });
}

export function handleScoreReviewApi(id: number, data: any) {
  return request.put(`/admin/score-reviews/${id}`, data);
}
