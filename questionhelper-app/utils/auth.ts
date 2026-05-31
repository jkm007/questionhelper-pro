const TOKEN_KEY = 'token'
const REFRESH_TOKEN_KEY = 'refresh_token'
const USER_INFO_KEY = 'user_info'

// 获取 Token
export const getToken = (): string => {
  return uni.getStorageSync(TOKEN_KEY) || ''
}

// 设置 Token
export const setToken = (token: string): void => {
  uni.setStorageSync(TOKEN_KEY, token)
}

// 移除 Token
export const removeToken = (): void => {
  uni.removeStorageSync(TOKEN_KEY)
}

// 获取刷新 Token
export const getRefreshToken = (): string => {
  return uni.getStorageSync(REFRESH_TOKEN_KEY) || ''
}

// 设置刷新 Token
export const setRefreshToken = (token: string): void => {
  uni.setStorageSync(REFRESH_TOKEN_KEY, token)
}

// 移除刷新 Token
export const removeRefreshToken = (): void => {
  uni.removeStorageSync(REFRESH_TOKEN_KEY)
}

// 获取用户信息
export const getUserInfo = (): any => {
  const info = uni.getStorageSync(USER_INFO_KEY)
  return info ? JSON.parse(info) : null
}

// 设置用户信息
export const setUserInfo = (userInfo: any): void => {
  uni.setStorageSync(USER_INFO_KEY, JSON.stringify(userInfo))
}

// 移除用户信息
export const removeUserInfo = (): void => {
  uni.removeStorageSync(USER_INFO_KEY)
}

// 是否已登录
export const isLoggedIn = (): boolean => {
  return !!getToken()
}

// 清除所有认证信息
export const clearAuth = (): void => {
  removeToken()
  removeRefreshToken()
  removeUserInfo()
}
