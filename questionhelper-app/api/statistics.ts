import { request } from './request'

// 用户统计
export const getUserStatistics = () => {
  return request({ url: '/statistics/user' })
}

// 练习统计
export const getPracticeStatistics = () => {
  return request({ url: '/statistics/practice' })
}

// 考试统计
export const getExamStatistics = () => {
  return request({ url: '/statistics/exam' })
}

// 排行榜
export const getRanking = (params?: { type?: number; page?: number; pageSize?: number }) => {
  return request({ url: '/statistics/ranking', data: params })
}

// 个人统计（别名）
export const getPersonalStatistics = getUserStatistics

// 练习趋势
export const getPracticeTrend = (params?: { days?: number }) => {
  return request({ url: '/statistics/practice/trend', data: params })
}

// 分数趋势
export const getScoreTrend = (params?: { days?: number }) => {
  return request({ url: '/statistics/exam/trend', data: params })
}

// 分类正确率
export const getCategoryAccuracy = () => {
  return request({ url: '/statistics/category-accuracy' })
}

// 统计概览
export const getStatisticsSummary = () => {
  return request({ url: '/statistics/mobile/overview' })
}

// 每周练习数据
export const getWeeklyPractice = () => {
  return request({ url: '/statistics/mobile/practice' })
}

// 考试分数趋势
export const getExamScoreTrend = (params?: { days?: number }) => {
  return request({ url: '/statistics/exam/trend', data: params })
}

// 排行榜（别名）
export const getRankings = getRanking

// 移动端正答率趋势
export const getAccuracyTrend = (params?: { days?: number }) => {
  return request({ url: '/statistics/mobile/accuracy-trend', data: params })
}

// 移动端分类分析
export const getMobileCategoryAnalysis = () => {
  return request({ url: '/statistics/mobile/category-analysis' })
}
