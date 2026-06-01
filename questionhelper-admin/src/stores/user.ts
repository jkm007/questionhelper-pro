import { defineStore } from "pinia";
import { ref } from "vue";
import { getToken, setToken, removeToken, setRefreshToken, removeRefreshToken } from "@/utils/auth";
import { loginApi, getUserInfoApi, logoutApi } from "@/api/auth";

interface UserInfo {
  id: number;
  username: string;
  nickname: string;
  avatar: string;
  roles: string[];
  permissions: string[];
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
    userInfo.value = data;
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
    if (userInfo.value.roles.includes("super_admin")) return true;
    return userInfo.value.permissions.includes(perm);
  }

  return { token, userInfo, login, getUserInfo, logout, hasPermission };
});
