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

// 创建作业（教师）
export const createHomework = (classId: number | string, data: {
  title: string
  description?: string
  deadline: string
  questionIds: number[]
  attachments?: string[]
}) => {
  return request({ url: `/classes/${classId}/homework`, method: 'POST', data })
}

// 作业详情
export const getHomeworkDetail = (classId: number | string, homeworkId: number | string) => {
  return request({ url: `/classes/${classId}/homework/${homeworkId}` })
}

// 提交作业（学生）
export const submitHomework = (classId: number | string, homeworkId: number | string, data: {
  content: string
  attachments?: string[]
}) => {
  return request({ url: `/classes/${classId}/homework/${homeworkId}/submit`, method: 'POST', data })
}

// 作业提交详情（学生查看自己的提交）
export const getHomeworkSubmit = (classId: number | string, homeworkId: number | string) => {
  return request({ url: `/classes/${classId}/homework/${homeworkId}/submit` })
}

// 作业结果（批改结果）
export const getHomeworkResult = (classId: number | string, homeworkId: number | string) => {
  return request({ url: `/classes/${classId}/homework/${homeworkId}/result` })
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

// ========== 考勤 Attendance ==========

// 考勤会话列表
export const getAttendanceSessions = (classId: number | string, params?: { page?: number; pageSize?: number }) => {
  return request({ url: `/classes/${classId}/attendance/sessions`, data: params })
}

// 考勤会话签到记录
export const getAttendanceRecords = (classId: number | string, sessionId: number | string) => {
  return request({ url: `/classes/${classId}/attendance/sessions/${sessionId}/records` })
}

// 创建考勤会话（教师）
export const createAttendanceSession = (classId: number | string, data: {
  title: string
  deadlineTime?: string
  duration?: number
  remark?: string
}) => {
  return request({ url: `/classes/${classId}/attendance/sessions`, method: 'POST', data })
}

// 学生签到
export const studentCheckIn = (classId: number | string, sessionId: number | string) => {
  return request({ url: `/classes/${classId}/attendance/sessions/${sessionId}/checkin`, method: 'POST' })
}

// 导出考勤记录
export const exportAttendance = (classId: number | string, sessionId: number | string) => {
  return request({ url: `/classes/${classId}/attendance/sessions/${sessionId}/export` })
}

// 考勤日历数据
export const getAttendanceCalendar = (classId: number | string, params: { year: number; month: number }) => {
  return request({ url: `/classes/${classId}/attendance/calendar`, data: params })
}

// ========== 教师审批 Teacher Approval ==========

// 获取加入申请列表
export const getClassApplications = (classId: number | string, params?: {
  page?: number
  pageSize?: number
  status?: string
}) => {
  return request({ url: `/classes/${classId}/applications`, data: params })
}

// 通过加入申请
export const approveApplication = (classId: number | string, applicationId: number | string) => {
  return request({ url: `/classes/${classId}/applications/${applicationId}/approve`, method: 'POST' })
}

// 拒绝加入申请
export const rejectApplication = (classId: number | string, applicationId: number | string, data: {
  reason?: string
}) => {
  return request({ url: `/classes/${classId}/applications/${applicationId}/reject`, method: 'POST', data })
}

// 批量通过加入申请
export const batchApproveApplications = (classId: number | string, data: {
  ids: string[]
}) => {
  return request({ url: `/classes/${classId}/applications/batch-approve`, method: 'POST', data })
}
