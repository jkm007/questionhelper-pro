import request from "@/utils/request";

export function listClassesApi(params?: any) {
  return request.get("/admin/class", { params });
}

export function getClassApi(id: number) {
  return request.get(`/admin/class/${id}`);
}

export function createClassApi(data: any) {
  return request.post("/admin/class", data);
}

export function updateClassApi(id: number, data: any) {
  return request.put(`/admin/class/${id}`, data);
}

export function deleteClassApi(id: number) {
  return request.delete(`/admin/class/${id}`);
}

export function listClassMembersApi(id: number, params?: any) {
  return request.get(`/admin/class/${id}/members`, { params });
}

export function removeClassMemberApi(id: number, uid: number) {
  return request.delete(`/admin/class/${id}/members/${uid}`);
}
