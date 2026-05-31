import { request } from './request'

// 题目列表
export const getQuestions = (params?: {
  page?: number
  pageSize?: number
  categoryId?: number
  type?: number
  difficulty?: number
  keyword?: string
}) => {
  return request({ url: '/questions', data: params })
}

// 题目详情
export const getQuestionDetail = (id: number) => {
  return request({ url: `/questions/${id}` })
}

// 搜索题目
export const searchQuestions = (keyword: string, params?: { page?: number; pageSize?: number }) => {
  return request({ url: '/questions/search', data: { keyword, ...params } })
}

// 收藏题目
export const favoriteQuestion = (id: number) => {
  return request({ url: `/questions/${id}/favorite`, method: 'POST' })
}

// 点赞题目
export const likeQuestion = (id: number) => {
  return request({ url: `/questions/${id}/like`, method: 'POST' })
}

// 分类列表
export const getCategories = () => {
  return request({ url: '/categories' })
}

// 分类树
export const getCategoryTree = () => {
  return request({ url: '/categories/tree' })
}

// 知识点列表
export const getKnowledgePoints = (params?: { categoryId?: number }) => {
  return request({ url: '/knowledge-points', data: params })
}

// 我的创作列表
export const getMyCreations = (params?: { page?: number; pageSize?: number; type?: number }) => {
  return request({ url: '/questions/my-creations', data: params })
}

// 删除创作
export const deleteCreation = (id: number) => {
  return request({ url: `/questions/${id}`, method: 'DELETE' })
}
