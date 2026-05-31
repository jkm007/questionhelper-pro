import { request } from './request'

// 考试反馈
export const submitExamFeedback = (data: {
  examId: number
  rating: number
  difficulty: number
  suggestion?: string
}) => {
  return request({ url: '/feedback/exam', method: 'POST', data })
}

// 题目纠错
export const submitCorrection = (data: {
  questionId: number
  type: number // 1:答案错误 2:解析错误 3:选项错误 4:其他
  description: string
  images?: string[]
}) => {
  return request({ url: '/feedback/correction', method: 'POST', data })
}

// 用户建议
export const submitSuggestion = (data: {
  content: string
  contact?: string
}) => {
  return request({ url: '/feedback/suggestion', method: 'POST', data })
}

// 题目纠错（别名）
export const submitQuestionCorrection = submitCorrection

// 类型定义
export interface CorrectionParams {
  questionId: number
  type: number
  description: string
  images?: string[]
  correctAnswer?: string
  explanation?: string
  reference?: string
}

export interface QuestionInfo {
  id: number
  title: string
  type: number
  options?: { label: string; content: string }[]
  correctAnswer?: string
}

export interface ExamFeedbackParams {
  examId: number
  type: number
  description: string
  questionId?: number
  images?: string[]
}
