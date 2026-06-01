import request from "@/utils/request";

export function listBackupsApi(params?: any) {
  return request.get("/admin/backup/list", { params });
}

export function createBackupApi() {
  return request.post("/admin/backup/create");
}

export function restoreBackupApi(id: number) {
  return request.post(`/admin/backup/${id}/restore`);
}

export function deleteBackupApi(id: number) {
  return request.delete(`/admin/backup/${id}`);
}

export function getBackupConfigsApi() {
  return request.get("/admin/backup/config");
}

export function updateBackupConfigApi(data: any) {
  return request.put("/admin/backup/config", data);
}
