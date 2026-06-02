// H5 模式使用相对路径走 Vite 代理；其他平台使用完整地址
import { getToken, setToken, getRefreshToken, setRefreshToken, clearAuth } from '@/utils/auth'
import { refreshToken as refreshTokenApi } from './auth'

export const BASE_URL = import.meta.env.VITE_API_BASE_URL || '/api/v1'

interface RequestOptions {
  url: string
  method?: 'GET' | 'POST' | 'PUT' | 'DELETE'
  data?: any
  header?: Record<string, string>
}

interface ApiResponse<T = any> {
  code: string
  msg: string
  data: T
}

// Token 刷新状态
let isRefreshing = false
let refreshSubscribers: ((token: string) => void)[] = []

function subscribeTokenRefresh(cb: (token: string) => void) {
  refreshSubscribers.push(cb)
}

function onTokenRefreshed(newToken: string) {
  refreshSubscribers.forEach(cb => cb(newToken))
  refreshSubscribers = []
}

// 重试请求
function retryRequest(options: RequestOptions, newToken: string): Promise<ApiResponse> {
  return new Promise((resolve, reject) => {
    uni.request({
      url: BASE_URL + options.url,
      method: options.method || 'GET',
      data: options.data,
      header: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${newToken}`,
        ...options.header
      },
      success: (res: any) => {
        const data = res.data
        if (data.code === '00000') {
          resolve(data)
        } else {
          reject(new Error(data.msg))
        }
      },
      fail: reject
    })
  })
}

export const request = <T = any>(options: RequestOptions): Promise<ApiResponse<T>> => {
  return new Promise((resolve, reject) => {
    const token = getToken()

    uni.request({
      url: BASE_URL + options.url,
      method: options.method || 'GET',
      data: options.data,
      header: {
        'Content-Type': 'application/json',
        'Authorization': token ? `Bearer ${token}` : '',
        ...options.header
      },
      success: async (res: any) => {
        const data = res.data as ApiResponse<T>

        if (data.code === '00000') {
          resolve(data)
        } else if (data.code === 'A0003') {
          // Access Token 过期，尝试刷新
          const refreshToken = getRefreshToken()
          if (!refreshToken) {
            clearAuth()
            uni.reLaunch({ url: '/pages/login/index' })
            reject(new Error('登录已过期'))
            return
          }

          if (!isRefreshing) {
            isRefreshing = true
            try {
              const refreshRes = await refreshTokenApi()
              const newToken = refreshRes.data.accessToken
              const newRefreshToken = refreshRes.data.refreshToken
              setToken(newToken)
              setRefreshToken(newRefreshToken)
              onTokenRefreshed(newToken)

              // 重试当前请求
              const retryRes = await retryRequest(options, newToken)
              resolve(retryRes)
            } catch (e) {
              clearAuth()
              uni.reLaunch({ url: '/pages/login/index' })
              reject(new Error('刷新 Token 失败'))
            } finally {
              isRefreshing = false
              refreshSubscribers = []
            }
          } else {
            // 正在刷新，将请求加入队列
            subscribeTokenRefresh(async (newToken: string) => {
              const retryRes = await retryRequest(options, newToken)
              resolve(retryRes)
            })
          }
        } else if (data.code === 'A0004') {
          // Refresh Token 也无效
          clearAuth()
          uni.reLaunch({ url: '/pages/login/index' })
          reject(new Error('登录已过期'))
        } else {
          uni.showToast({ title: data.msg || '请求失败', icon: 'none' })
          reject(new Error(data.msg))
        }
      },
      fail: (err: any) => {
        uni.showToast({ title: '网络错误', icon: 'none' })
        reject(err)
      }
    })
  })
}
