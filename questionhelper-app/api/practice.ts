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
