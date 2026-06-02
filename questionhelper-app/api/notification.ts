import { request } from './request'

// 通知列表
export const getNotifications = (params?: { page?: number; pageSize?: number; type?: number }) => {
  return request({ url: '/notifications', data: params })
}

// 未读通知数量
export const getUnreadCount = () => {
  return request({ url: '/notifications/unread-count' })
}

// 标记已读
export const markAsRead = (id: number) => {
  return request({ url: `/notifications/${id}/read`, method: 'PUT' })
}

// 全部已读
export const markAllAsRead = () => {
  return request({ url: '/notifications/read-all', method: 'PUT' })
}

// 删除通知
export const deleteNotification = (id: number) => {
  return request({ url: `/notifications/${id}`, method: 'DELETE' })
}

// 通知列表（别名）
export const getNotificationList = getNotifications

// 通知详情
export const getNotificationDetail = (id: number) => {
  return request({ url: `/notifications/${id}` })
}

// 获取通知设置
export const getNotificationSettings = () => {
  return request({ url: '/notifications/settings' })
}

// 更新通知设置
export const updateNotificationSettings = (data: {
  systemEnabled?: boolean
  examEnabled?: boolean
  homeworkEnabled?: boolean
  classEnabled?: boolean
  commentEnabled?: boolean
}) => {
  return request({ url: '/notifications/settings', method: 'PUT', data })
}

// 获取通知统计
export const getNotificationStats = () => {
  return request({ url: '/notifications/stats' })
}

// 类型定义
export interface NotificationItem {
  id: number
  type: number
  title: string
  content: string
  isRead: boolean
  createdAt: string
  relatedId?: number
  relatedType?: string
}
