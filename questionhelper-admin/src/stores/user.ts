import { store } from "@/stores";
import { loginApi, getUserInfoApi, logoutApi } from "@/api/auth";
import { getToken, setToken, removeToken, setRefreshToken, removeRefreshToken } from "@/utils/auth";
import { usePermissionStoreHook } from "@/stores/permission";

interface UserInfo {
  id: number;
  username: string;
  nickname: string;
  avatar: string;
  roles: string[];
}

export const useUserStore = defineStore("user", () => {
  const userInfo = ref<UserInfo>({} as UserInfo);
  const rememberMe = ref(false);

  async function login(username: string, password: string): Promise<void> {
    const data: any = await loginApi({ username, password });
    const accessToken = data.accessToken;
    const refreshToken = data.refreshToken;
    setToken(accessToken);
    setRefreshToken(refreshToken);
    rememberMe.value = false;
  }

  async function getUserInfo(): Promise<UserInfo> {
    const data: any = await getUserInfoApi();
    const roles = Array.isArray(data.roles)
      ? data.roles.map((r: any) => (typeof r === "string" ? r : r.code))
      : [];
    userInfo.value = { ...data, roles };
    return userInfo.value;
  }

  async function logout(): Promise<void> {
    try {
      await logoutApi();
    } catch {
      // ignore
    }
    resetUserState();
    usePermissionStoreHook().resetRouter();
  }

  function resetUserState(): void {
    removeToken();
    removeRefreshToken();
    userInfo.value = {} as UserInfo;
  }

  function isLoggedIn(): boolean {
    return !!getToken();
  }

  return {
    userInfo,
    rememberMe,
    isLoggedIn,
    login,
    logout,
    getUserInfo,
    resetUserState,
  };
});

export function useUserStoreHook() {
  return useUserStore(store);
}
