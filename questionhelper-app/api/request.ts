// H5 模式使用相对路径走 Vite 代理；其他平台使用完整地址
const BASE_URL = import.meta.env.VITE_API_BASE_URL || '/api/v1'

interface RequestOptions {
  url: string
  method?: 'GET' | 'POST' | 'PUT' | 'DELETE'
  data?: any
  header?: Record<string, string>
  showLoading?: boolean
}

interface ApiResponse<T = any> {
  code: string
  msg: string
  data: T
}

export const request = <T = any>(options: RequestOptions): Promise<ApiResponse<T>> => {
  return new Promise((resolve, reject) => {
    const token = uni.getStorageSync('token')

    if (options.showLoading !== false) {
      uni.showLoading({ title: '加载中...' })
    }

    uni.request({
      url: BASE_URL + options.url,
      method: options.method || 'GET',
      data: options.data,
      header: {
        'Content-Type': 'application/json',
        'Authorization': token ? `Bearer ${token}` : '',
        ...options.header
      },
      success: (res: any) => {
        if (options.showLoading !== false) {
          uni.hideLoading()
        }

        const data = res.data as ApiResponse<T>

        if (data.code === '00000') {
          resolve(data)
        } else if (data.code === 'A0003' || data.code === 'A0004') {
          // token 过期或无效，跳转登录
          uni.removeStorageSync('token')
          uni.navigateTo({ url: '/pages/login/index' })
          reject(new Error('登录已过期'))
        } else {
          uni.showToast({ title: data.msg || '请求失败', icon: 'none' })
          reject(new Error(data.msg))
        }
      },
      fail: (err: any) => {
        if (options.showLoading !== false) {
          uni.hideLoading()
        }
        uni.showToast({ title: '网络错误', icon: 'none' })
        reject(err)
      }
    })
  })
}
