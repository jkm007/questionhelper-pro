import { request } from './request'

// 开始练习
export const startPractice = (data: { categoryId?: number; type?: number; difficulty?: number; count?: number }) => {
  return request({ url: '/practice/start', method: 'POST', data })
}

// 提交答案
export const submitPracticeAnswer = (data: { practiceId: number; answers: { questionId: number; answer: string; duration: number }[] }) => {
  return request({ url: '/practice/submit', method: 'POST', data })
}

// 练习结果
export const getPracticeResult = (id: number) => {
  return request({ url: `/practice/${id}` })
}

// 练习历史
export const getPracticeHistory = (params?: { page?: number; pageSize?: number }) => {
  return request({ url: '/practice', data: params })
}

// 练习统计
export const getPracticeStatistics = () => {
  return request({ url: '/practice/stats' })
}

// 最近练习记录（别名）
export const getRecentRecords = (params?: { page?: number; pageSize?: number }) => {
  return request({ url: '/practice', data: params })
}

// 今日统计
export const getTodayStats = () => {
  return request({ url: '/practice/stats' })
}

// 获取练习题目
export const getPracticeQuestions = (id: number) => {
  return request({ url: `/practice/${id}` })
}

// 提交练习答案（别名，支持批量）
export const submitPracticeAnswers = (data: { practiceId: number; answers: { questionId: number; answer: string; duration: number }[] }) => {
  return request({ url: '/practice/submit', method: 'POST', data })
}

// 切换题目收藏
export const toggleQuestionBookmark = (practiceId: number, questionId: number) => {
  return request({ url: `/practice/${practiceId}/bookmark`, method: 'POST', data: { questionId } })
}

// 暂停练习
export const pausePractice = (id: number) => {
  return request({ url: `/practice/${id}/pause`, method: 'POST' })
}

// 恢复练习
export const resumePractice = (id: number) => {
  return request({ url: `/practice/${id}/resume`, method: 'POST' })
}

// 结束练习
export const finishPractice = (id: number) => {
  return request({ url: `/practice/${id}/finish`, method: 'POST' })
}

// 获取学习计划列表
export const getPlans = (params?: { page?: number; pageSize?: number }) => {
  return request({ url: '/practice/plans', data: params })
}

// 创建学习计划
export const createPlan = (data: any) => {
  return request({ url: '/practice/plans', method: 'POST', data })
}

// 获取学习计划详情
export const getPlan = (id: number) => {
  return request({ url: `/practice/plans/${id}` })
}

// 更新学习计划
export const updatePlan = (id: number, data: any) => {
  return request({ url: `/practice/plans/${id}`, method: 'PUT', data })
}

// 删除学习计划
export const deletePlan = (id: number) => {
  return request({ url: `/practice/plans/${id}`, method: 'DELETE' })
}

// 执行学习计划
export const executePlan = (id: number) => {
  return request({ url: `/practice/plans/${id}/execute`, method: 'POST' })
}

// 获取今日练习
export const getTodayPractice = () => {
  return request({ url: '/practice/daily/today' })
}

// 完成每日练习
export const completeDailyPractice = (data: any) => {
  return request({ url: '/practice/daily/complete', method: 'POST', data })
}

// 签到
export const checkin = () => {
  return request({ url: '/practice/checkin', method: 'POST' })
}

// 获取签到日历
export const getCheckinCalendar = (params?: { month?: string }) => {
  return request({ url: '/practice/checkin/calendar', data: params })
}

// 获取排行榜
export const getLeaderboard = (params?: { type?: number; page?: number; pageSize?: number }) => {
  return request({ url: '/practice/leaderboard', data: params })
}

// 开始模拟考试
export const startMockExam = (data: any) => {
  return request({ url: '/practice/mock/start', method: 'POST', data })
}

// 获取模拟考试历史
export const getMockExamHistory = (params?: { page?: number; pageSize?: number }) => {
  return request({ url: '/practice/mock/history', data: params })
}

// 提交模拟考试
export const submitMockExam = (id: number, data: any) => {
  return request({ url: `/practice/mock/${id}/submit`, method: 'POST', data })
}
