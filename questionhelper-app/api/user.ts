import { request, BASE_URL } from './request'

// 获取个人信息
export const getProfile = () => {
  return request({ url: '/user/profile' })
}

// 更新个人信息
export const updateProfile = (data: any) => {
  return request({ url: '/user/profile', method: 'PUT', data })
}

// 修改密码
export const updatePassword = (data: { oldPassword: string; newPassword: string }) => {
  return request({ url: '/user/password', method: 'PUT', data })
}

// 上传头像
export const uploadAvatar = (filePath: string) => {
  return new Promise((resolve, reject) => {
    const token = uni.getStorageSync('token')
    uni.uploadFile({
      url: `${BASE_URL}/user/avatar`,
      filePath,
      name: 'file',
      header: {
        'Authorization': `Bearer ${token}`
      },
      success: (res) => {
        const data = JSON.parse(res.data)
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

// 实名认证
export const realNameAuth = (data: { realName: string; idCard: string }) => {
  return request({ url: '/user/realname', method: 'POST', data })
}

// 收藏列表
export const getFavorites = (params?: { page?: number; pageSize?: number }) => {
  return request({ url: '/user/favorites', data: params })
}

// 添加收藏
export const addFavorite = (data: { targetType: number; targetId: number }) => {
  return request({ url: '/user/favorites', method: 'POST', data })
}

// 取消收藏
export const removeFavorite = (id: number) => {
  return request({ url: `/user/favorites/${id}`, method: 'DELETE' })
}

// 获取个人信息（别名）
export const getUserInfo = getProfile

// 更新个人信息（别名）
export const updateUserInfo = updateProfile

// 获取实名认证状态
export const getAuthStatus = () => {
  return request({ url: '/user/realname/status' })
}

// 提交实名认证
export const submitAuth = (data: { realName: string; idCard: string; frontImage: string; backImage: string }) => {
  return request({ url: '/user/realname/submit', method: 'POST', data })
}

// 获取创作者状态
export const getCreatorStatus = () => {
  return request({ url: '/user/creator/status' })
}

// 申请创作者
export const submitCreatorApply = (data: { reason: string; portfolio: string }) => {
  return request({ url: '/user/apply/creator', method: 'POST', data })
}

// 获取设置
export const getSettings = () => {
  return request({ url: '/user/settings' })
}

// 更新设置
export const updateSettings = (data: { pushEnabled?: boolean; soundEnabled?: boolean; vibrationEnabled?: boolean; language?: string }) => {
  return request({ url: '/user/settings', method: 'PUT', data })
}

// 账号注销
export const deactivateAccount = (data: { password: string; reason?: string }) => {
  return request({ url: '/user/account/deactivate', method: 'POST', data })
}

// 绑定邮箱
export const bindEmail = (data: { email: string; code: string }) => {
  return request({ url: '/user/bindemail', method: 'POST', data })
}

// 绑定手机
export const bindPhone = (data: { phone: string; code: string }) => {
  return request({ url: '/user/bindphone', method: 'POST', data })
}

// 隐私设置
export const getPrivacy = () => {
  return request({ url: '/user/privacy' })
}

// 更新隐私设置
export const updatePrivacy = (data: any) => {
  return request({ url: '/user/privacy', method: 'PUT', data })
}

// 获取教师申请状态
export const getTeacherStatus = () => {
  return request({ url: '/user/teacher/status' })
}

// 提交教师申请
export const submitTeacherApply = (data: {
  reason: string
  qualifications: string
  school: string
  subject: string
  certificates?: string[]
}) => {
  return request({ url: '/user/teacher/apply', method: 'POST', data })
}

// 获取学习报告
export const getLearningReport = (params?: { period?: number }) => {
  return request({ url: '/statistics/learning-report', data: params })
}

// 获取学习计划
export const getLearningPlan = () => {
  return request({ url: '/user/learning-plan' })
}

// 创建学习计划
export const createLearningPlan = (data: any) => {
  return request({ url: '/user/learning-plan', method: 'POST', data })
}

// 更新学习计划
export const updateLearningPlan = (data: any) => {
  return request({ url: '/user/learning-plan', method: 'PUT', data })
}

// 添加学习计划项
export const addPlanItem = (data: any) => {
  return request({ url: '/user/learning-plan/items', method: 'POST', data })
}

// 更新学习计划项
export const updatePlanItem = (id: number, data: any) => {
  return request({ url: `/user/learning-plan/items/${id}`, method: 'PUT', data })
}

// 删除学习计划项
export const deletePlanItem = (id: number) => {
  return request({ url: `/user/learning-plan/items/${id}`, method: 'DELETE' })
}
