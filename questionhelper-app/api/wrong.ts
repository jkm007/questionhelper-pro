import { request } from './request'

// 错题列表
export const getWrongQuestions = (params?: { page?: number; pageSize?: number; mastered?: boolean }) => {
  return request({ url: '/wrong', data: params })
}

// 错题详情
export const getWrongDetail = (id: number) => {
  return request({ url: `/wrong/${id}` })
}

// 错题复习
export const reviewWrong = (id: number, data: { answer: string }) => {
  return request({ url: `/wrong/${id}/review`, method: 'POST', data })
}

// 移除错题
export const removeWrong = (id: number) => {
  return request({ url: `/wrong/${id}`, method: 'DELETE' })
}

// 错题分析
export const getWrongAnalysis = () => {
  return request({ url: '/wrong/analysis' })
}

// 错题列表（别名）
export const getWrongList = getWrongQuestions

// 批量移除错题
export const removeWrongQuestions = (ids: number[]) => {
  return request({ url: '/wrong/batch/remove', method: 'POST', data: { ids } })
}

// 掌握度趋势
export const getMasteryTrend = (params?: { days?: number }) => {
  return request({ url: '/wrong/mastery-trend', data: params })
}

// 复习推荐
export const getRecommendations = () => {
  return request({ url: '/wrong/recommendations' })
}

// 更新掌握度
export const updateMastery = (id: number, data: { mastered: boolean }) => {
  return request({ url: `/wrong/${id}/mastery`, method: 'PUT', data })
}

// 更新错题笔记
export const updateWrongNotes = (id: number, data: { notes: string }) => {
  return request({ url: `/wrong/${id}/notes`, method: 'PUT', data })
}

// 批量收藏错题
export const batchFavoriteWrong = (ids: number[]) => {
  return request({ url: '/wrong/batch/favorite', method: 'POST', data: { ids } })
}

// 导出错题
export const exportWrongQuestions = (data: { format: string; masteryLevel?: number; categoryId?: number }) => {
  return request({ url: '/wrong/export', method: 'POST', data })
}

// 获取掌握度趋势
export const getWrongTrend = (params?: { days?: number }) => {
  return request({ url: '/wrong/analysis/trend', data: params })
}

// 获取分类分析
export const getWrongCategoryAnalysis = () => {
  return request({ url: '/wrong/analysis/category' })
}

// 获取正确率分析
export const getWrongAccuracyAnalysis = () => {
  return request({ url: '/wrong/analysis/accuracy' })
}

// 获取今日复习题目
export const getTodayReviewQuestions = () => {
  return request({ url: '/wrong/review/today' })
}

// 记录复习结果
export const recordReviewResult = (id: number, data: { result: string; duration?: number }) => {
  return request({ url: `/wrong/${id}/review/record`, method: 'POST', data })
}

// 获取复习历史
export const getReviewHistory = (params?: { page?: number; pageSize?: number }) => {
  return request({ url: '/wrong/review/history', data: params })
}

// 收藏错题
export const favoriteWrong = (id: number) => {
  return request({ url: `/wrong/${id}/favorite`, method: 'POST' })
}

// 取消收藏错题
export const unfavoriteWrong = (id: number) => {
  return request({ url: `/wrong/${id}/favorite`, method: 'DELETE' })
}

// 获取收藏错题列表
export const getWrongFavorites = (params?: { page?: number; pageSize?: number }) => {
  return request({ url: '/wrong/favorites', data: params })
}
