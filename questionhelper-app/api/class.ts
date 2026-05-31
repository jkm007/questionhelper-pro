import { request } from './request'

// 班级列表
export const getClasses = (params?: { page?: number; pageSize?: number }) => {
  return request({ url: '/classes', data: params })
}

// 班级详情
export const getClassDetail = (id: number) => {
  return request({ url: `/classes/${id}` })
}

// 创建班级
export const createClass = (data: { name: string; description?: string; cover?: string }) => {
  return request({ url: '/classes', method: 'POST', data })
}

// 更新班级
export const updateClass = (id: number, data: { name?: string; description?: string; cover?: string }) => {
  return request({ url: `/classes/${id}`, method: 'PUT', data })
}

// 加入班级
export const joinClass = (id: number, data: { code: string }) => {
  return request({ url: `/classes/${id}/join`, method: 'POST', data })
}

// 退出班级
export const leaveClass = (id: number) => {
  return request({ url: `/classes/${id}/leave`, method: 'POST' })
}

// 成员列表
export const getClassMembers = (id: number, params?: { page?: number; pageSize?: number }) => {
  return request({ url: `/classes/${id}/members`, data: params })
}

// 班级考试
export const getClassExams = (id: number, params?: { page?: number; pageSize?: number }) => {
  return request({ url: `/classes/${id}/exams`, data: params })
}

// 班级作业
export const getClassHomework = (id: number, params?: { page?: number; pageSize?: number }) => {
  return request({ url: `/classes/${id}/homework`, data: params })
}

// 班级公告
export const getClassNotices = (id: number, params?: { page?: number; pageSize?: number }) => {
  return request({ url: `/classes/${id}/notice`, data: params })
}

// 我的班级
export const getMyClasses = (params?: { page?: number; pageSize?: number }) => {
  return request({ url: '/classes', data: { ...params, mine: true } })
}

// 发现班级
export const discoverClasses = (params?: { page?: number; pageSize?: number; keyword?: string }) => {
  return request({ url: '/classes', data: { ...params, discover: true } })
}

// 搜索班级（通过班级码）
export const searchClassByCode = (code: string) => {
  return request({ url: '/classes/search', data: { code } })
}

// 加入历史
export const getJoinHistory = (params?: { page?: number; pageSize?: number }) => {
  return request({ url: '/classes/join-history', data: params })
}

// 班级考试列表
export const getClassExamList = (id: number, params?: { page?: number; pageSize?: number; status?: number }) => {
  return request({ url: `/classes/${id}/exams`, data: params })
}

// 班级作业列表
export const getClassHomeworkList = (id: number, params?: { page?: number; pageSize?: number; status?: number }) => {
  return request({ url: `/classes/${id}/homework`, data: params })
}

// 提交作业
export const submitHomework = (classId: number, homeworkId: number, data: { content: string; attachments?: string[] }) => {
  return request({ url: `/classes/${classId}/homework/${homeworkId}/submit`, method: 'POST', data })
}

// 更新成员角色
export const updateMemberRole = (classId: number, memberId: number, data: { role: string }) => {
  return request({ url: `/classes/${classId}/members/${memberId}/role`, method: 'PUT', data })
}

// 移除成员
export const removeMember = (classId: number, memberId: number) => {
  return request({ url: `/classes/${classId}/members/${memberId}`, method: 'DELETE' })
}

// 班级公告列表
export const getClassNoticeList = (id: number, params?: { page?: number; pageSize?: number }) => {
  return request({ url: `/classes/${id}/notice`, data: params })
}

// 标记公告已读
export const markNoticeRead = (classId: number, noticeId: number) => {
  return request({ url: `/classes/${classId}/notice/${noticeId}/read`, method: 'PUT' })
}
