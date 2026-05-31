import { request } from './request'

// 登录
export const login = (data: { username: string; password: string }) => {
  return request({ url: '/auth/login', method: 'POST', data })
}

// 注册
export const register = (data: { username: string; password: string; phone?: string; email?: string }) => {
  return request({ url: '/auth/register', method: 'POST', data })
}

// 退出
export const logout = () => {
  return request({ url: '/auth/logout', method: 'POST' })
}

// 刷新 Token
export const refreshToken = () => {
  return request({ url: '/auth/refresh', method: 'POST' })
}

// 获取验证码
export const getCaptcha = () => {
  return request({ url: '/auth/captcha' })
}
