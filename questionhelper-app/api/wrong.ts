import { request } from './request'

// 错题列表
export const getWrongQuestions = (params?: { page?: number; pageSize?: number; mastered?: boolean }) => {
  return request({ url: '/wrong-questions', data: params })
}

// 错题详情
export const getWrongDetail = (id: number) => {
  return request({ url: `/wrong-questions/${id}` })
}

// 错题复习
export const reviewWrong = (id: number, data: { answer: string }) => {
  return request({ url: `/wrong-questions/${id}/review`, method: 'POST', data })
}

// 移除错题
export const removeWrong = (id: number) => {
  return request({ url: `/wrong-questions/${id}`, method: 'DELETE' })
}

// 错题分析
export const getWrongAnalysis = () => {
  return request({ url: '/wrong-questions/analysis' })
}

// 错题列表（别名）
export const getWrongList = getWrongQuestions

// 批量移除错题
export const removeWrongQuestions = (ids: number[]) => {
  return request({ url: '/wrong-questions/batch/remove', method: 'POST', data: { ids } })
}

// 掌握度趋势
export const getMasteryTrend = (params?: { days?: number }) => {
  return request({ url: '/wrong-questions/mastery-trend', data: params })
}

// 复习推荐
export const getRecommendations = () => {
  return request({ url: '/wrong-questions/recommendations' })
}

// 更新掌握度
export const updateMastery = (id: number, data: { mastered: boolean }) => {
  return request({ url: `/wrong-questions/${id}/mastery`, method: 'PUT', data })
}

// 更新错题笔记
export const updateWrongNotes = (id: number, data: { notes: string }) => {
  return request({ url: `/wrong-questions/${id}/notes`, method: 'PUT', data })
}

// 批量收藏错题
export const batchFavoriteWrong = (ids: number[]) => {
  return request({ url: '/wrong-questions/batch/favorite', method: 'POST', data: { ids } })
}
