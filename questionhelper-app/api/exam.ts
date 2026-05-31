import { request } from './request'

// 考试列表
export const getExams = (params?: { page?: number; pageSize?: number; status?: number }) => {
  return request({ url: '/exam', data: params })
}

// 考试详情
export const getExamDetail = (id: number) => {
  return request({ url: `/exam/${id}` })
}

// 开始考试
export const startExam = (id: number) => {
  return request({ url: `/exam/${id}/start`, method: 'POST' })
}

// 提交答案
export const submitExam = (id: number, data: { answers: any[] }) => {
  return request({ url: `/exam/${id}/submit`, method: 'POST', data })
}

// 考试结果
export const getExamResult = (id: number) => {
  return request({ url: `/exam/${id}/result` })
}

// 考试历史
export const getExamHistory = (params?: { page?: number; pageSize?: number }) => {
  return request({ url: '/exam/history', data: params })
}

// 试卷列表
export const getPapers = (params?: { page?: number; pageSize?: number }) => {
  return request({ url: '/paper', data: params })
}

// 试卷详情
export const getPaperDetail = (id: number) => {
  return request({ url: `/paper/${id}` })
}

// 考试列表（别名）
export const getExamList = getExams

// 考试历史列表（别名）
export const getExamHistoryList = getExamHistory

// 考试参与者
export const getExamParticipants = (id: number, params?: { page?: number; pageSize?: number }) => {
  return request({ url: `/exam/${id}/participants`, data: params })
}

// 获取考试题目
export const getExamQuestions = (id: number) => {
  return request({ url: `/exam/${id}/questions` })
}

// 提交考试答案（别名）
export const submitExamAnswers = submitExam

// 保存考试进度
export const saveExamProgress = (id: number, data: { answers: any[]; currentQuestion?: number }) => {
  return request({ url: `/exam/${id}/save-answer`, method: 'POST', data })
}
