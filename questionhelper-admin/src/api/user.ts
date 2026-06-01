import request from "@/utils/request";

// 用户管理
export function listUsersApi(params?: any) {
  return request.get("/admin/users", { params });
}

export function getUserApi(id: number) {
  return request.get(`/admin/users/${id}`);
}

export function createUserApi(data: any) {
  return request.post("/admin/users", data);
}

export function updateUserApi(id: number, data: any) {
  return request.put(`/admin/users/${id}`, data);
}

export function deleteUserApi(id: number) {
  return request.delete(`/admin/users/${id}`);
}

export function updateUserStatusApi(id: number, status: number) {
  return request.put(`/admin/users/${id}/status`, { status });
}

export function resetPasswordApi(id: number) {
  return request.post(`/admin/users/${id}/reset-password`);
}

export function assignRolesApi(id: number, roleIds: number[]) {
  return request.put(`/admin/users/${id}/roles`, { role_ids: roleIds });
}

export function exportUsersApi(params?: any) {
  return request.get("/admin/users/export", { params, responseType: "blob" });
}

// 角色管理
export function listRolesApi(params?: any) {
  return request.get("/admin/roles", { params });
}

export function getRoleApi(id: number) {
  return request.get(`/admin/roles/${id}`);
}

export function createRoleApi(data: any) {
  return request.post("/admin/roles", data);
}

export function updateRoleApi(id: number, data: any) {
  return request.put(`/admin/roles/${id}`, data);
}

export function deleteRoleApi(id: number) {
  return request.delete(`/admin/roles/${id}`);
}

export function assignRoleMenusApi(id: number, menuIds: number[]) {
  return request.put(`/admin/roles/${id}/menus`, { menu_ids: menuIds });
}

export function assignRolePermissionsApi(id: number, permIds: number[]) {
  return request.put(`/admin/roles/${id}/permissions`, { permission_ids: permIds });
}

// 标签管理
export function listTagsApi(params?: any) {
  return request.get("/admin/tags", { params });
}

export function createTagApi(data: any) {
  return request.post("/admin/tags", data);
}

export function updateTagApi(id: number, data: any) {
  return request.put(`/admin/tags/${id}`, data);
}

export function deleteTagApi(id: number) {
  return request.delete(`/admin/tags/${id}`);
}

// 实名审核
export function listRealNamesApi(params?: any) {
  return request.get("/admin/realnames", { params });
}

export function reviewRealNameApi(id: number, data: any) {
  return request.put(`/admin/realnames/${id}/review`, data);
}

// 角色申请审核
export function listApplicationsApi(params?: any) {
  return request.get("/admin/applications", { params });
}

export function getApplicationApi(id: number) {
  return request.get(`/admin/applications/${id}`);
}

export function reviewApplicationApi(id: number, data: any) {
  return request.put(`/admin/applications/${id}/review`, data);
}
