import { request, BASE_URL } from './request'

// 上传文件
export const uploadFile = (filePath: string) => {
  return new Promise((resolve, reject) => {
    const token = uni.getStorageSync('token')
    uni.uploadFile({
      url: `${BASE_URL}/files/upload`,
      filePath,
      name: 'file',
      header: {
        'Authorization': `Bearer ${token}`
      },
      success: (res) => {
        const data = JSON.parse(res.data)
        if (data.code === '00000') {
          resolve(data)
        } else {
          reject(new Error(data.msg))
        }
      },
      fail: reject
    })
  })
}

// 上传图片
export const uploadImage = (filePath: string) => {
  return new Promise((resolve, reject) => {
    const token = uni.getStorageSync('token')
    uni.uploadFile({
      url: `${BASE_URL}/files/upload/image`,
      filePath,
      name: 'file',
      header: {
        'Authorization': `Bearer ${token}`
      },
      success: (res) => {
        const data = JSON.parse(res.data)
        if (data.code === '00000') {
          resolve(data)
        } else {
          reject(new Error(data.msg))
        }
      },
      fail: reject
    })
  })
}

// 获取文件信息
export const getFileInfo = (id: number) => {
  return request({ url: `/files/${id}` })
}

// 删除文件
export const deleteFile = (id: number) => {
  return request({ url: `/files/${id}`, method: 'DELETE' })
}
