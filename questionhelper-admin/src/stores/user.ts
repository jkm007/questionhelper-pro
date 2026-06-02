import { defineStore } from "pinia";
import { ref } from "vue";
import { getToken, setToken, removeToken, setRefreshToken, removeRefreshToken } from "@/utils/auth";
import { loginApi, getUserInfoApi, logoutApi } from "@/api/auth";
import { usePermissionStore } from "@/stores/permission";

interface UserInfo {
  id: number;
  username: string;
  nickname: string;
  avatar: string;
  roles: string[];
}

export const useUserStore = defineStore("user", () => {
  const token = ref<string>(getToken() || "");
  const userInfo = ref<UserInfo | null>(null);

  async function login(username: string, password: string) {
    const data: any = await loginApi({ username, password });
    token.value = data.access_token;
    setToken(data.access_token);
    setRefreshToken(data.refresh_token);
    return data;
  }

  async function getUserInfo() {
    const data: any = await getUserInfoApi();
    // 后端返回 roles 为 RoleInfo[] 对象数组，需转换为 string[] 角色编码
    const roles = Array.isArray(data.roles)
      ? data.roles.map((r: any) => (typeof r === "string" ? r : r.code))
      : [];
    userInfo.value = { ...data, roles };
    return data;
  }

  async function logout() {
    try {
      await logoutApi();
    } catch {
      // ignore
    }
    token.value = "";
    userInfo.value = null;
    removeToken();
    removeRefreshToken();
  }

  function hasPermission(perm: string): boolean {
    if (!userInfo.value) return false;
    // 超级管理员拥有全部权限（需求文档 11.4）
    if (userInfo.value.roles.includes("super_admin")) return true;
    // 检查按钮权限列表（来自 GET /menus/buttons，menus.type=3）
    const permissionStore = usePermissionStore();
    return permissionStore.buttons.includes(perm);
  }

  return { token, userInfo, login, getUserInfo, logout, hasPermission };
});
