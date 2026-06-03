import request from "@/utils/request";

export function loginApi(data: { username: string; password: string }) {
  return request.post("/auth/login", data);
}

export function getUserInfoApi() {
  return request.get("/users/me");
}

export function logoutApi() {
  return request.post("/user/logout");
}

export function refreshTokenApi(refreshToken: string) {
  return request.post("/auth/refresh", { refresh_token: refreshToken });
}
