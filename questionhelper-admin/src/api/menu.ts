import request from "@/utils/request";

export function getRoutesApi() {
  return request.get("/menus/routes");
}

export function getButtonsApi() {
  return request.get("/menus/buttons");
}

export function listMenusApi() {
  return request.get("/admin/menus");
}

export function getMenuTreeApi() {
  return request.get("/admin/menus/tree");
}

export function getMenuApi(id: number) {
  return request.get(`/admin/menus/${id}`);
}

export function createMenuApi(data: any) {
  return request.post("/admin/menus", data);
}

export function updateMenuApi(id: number, data: any) {
  return request.put(`/admin/menus/${id}`, data);
}

export function deleteMenuApi(id: number) {
  return request.delete(`/admin/menus/${id}`);
}
