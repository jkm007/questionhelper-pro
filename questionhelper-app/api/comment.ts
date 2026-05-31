import { request } from './request'

// 评论列表
export const getComments = (params: { targetType: number; targetId: number; page?: number; pageSize?: number }) => {
  return request({ url: '/comments', data: params })
}

// 发表评论
export const addComment = (data: { targetType: number; targetId: number; content: string; parentId?: number }) => {
  return request({ url: '/comments', method: 'POST', data })
}

// 删除评论
export const deleteComment = (id: number) => {
  return request({ url: `/comments/${id}`, method: 'DELETE' })
}

// 点赞评论
export const likeComment = (id: number) => {
  return request({ url: `/comments/${id}/like`, method: 'POST' })
}

// 举报评论
export const reportComment = (id: number, data: { reason: string }) => {
  return request({ url: `/comments/${id}/report`, method: 'POST', data })
}

// 评论列表（别名）
export const getCommentList = getComments

// 发表评论（别名）
export const createComment = addComment

// 获取回复列表
export const getReplyList = (commentId: number, params?: { page?: number; pageSize?: number }) => {
  return request({ url: `/comments/${commentId}/replies`, data: params })
}
