import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { getToken, setToken, removeToken, setRefreshToken, removeRefreshToken, getUserInfo, setUserInfo, removeUserInfo } from '@/utils/auth'
import { login as loginApi, logout as logoutApi } from '@/api/auth'
import { getProfile } from '@/api/user'

export const useUserStore = defineStore('user', () => {
  const token = ref<string>(getToken())
  const userInfo = ref<any>(getUserInfo())

  const isLoggedIn = computed(() => !!token.value)
  const userId = computed(() => userInfo.value?.id)
  const nickname = computed(() => userInfo.value?.nickname || userInfo.value?.username || '')
  const avatar = computed(() => userInfo.value?.avatar || '')

  // 登录
  const login = async (username: string, password: string) => {
    const res = await loginApi({ username, password })
    token.value = res.data.accessToken
    setToken(res.data.accessToken)
    if (res.data.refreshToken) {
      setRefreshToken(res.data.refreshToken)
    }
    if (res.data.user) {
      userInfo.value = res.data.user
      setUserInfo(res.data.user)
    }
    return res
  }

  // 获取用户信息
  const fetchUserInfo = async () => {
    const res = await getProfile()
    userInfo.value = res.data
    setUserInfo(res.data)
    return res.data
  }

  // 退出登录
  const logout = async () => {
    try {
      await logoutApi()
    } catch (e) {
      // 忽略错误
    }
    token.value = ''
    userInfo.value = null
    removeToken()
    removeRefreshToken()
    removeUserInfo()
  }

  return {
    token,
    userInfo,
    isLoggedIn,
    userId,
    nickname,
    avatar,
    login,
    fetchUserInfo,
    logout
  }
})
