import { request } from '../utils/request'
import { getToken } from '../utils/auth'

//const BASE_URL = 'http://62.234.163.176:7001'
const BASE_URL = 'http://127.0.0.1:8080'

export const devLogin = (payload = {}) => {
  return request({
    url: '/auth/dev-login',
    method: 'POST',
    data: payload,
    auth: false
  })
}

export const wechatLogin = (code) => {
  return request({
    url: '/auth/wechat-login',
    method: 'POST',
    data: { code },
    auth: false
  })
}

export const reactivateAccount = (reactivationToken) => {
  return request({
    url: '/auth/reactivate',
    method: 'POST',
    data: {
      reactivationToken,
      confirm: true
    },
    auth: false
  })
}

export const getCurrentUser = () => {
  return request({
    url: '/users/me'
  })
}

export const cancelAccount = () => {
  return request({
    url: '/users/cancel',
    method: 'POST',
    data: { confirm: true }
  })
}

export const updateProfile = (payload) => {
  return request({
    url: '/users/profile',
    method: 'PUT',
    data: payload
  })
}

export const submitStudentVerification = (payload) => {
  return request({
    url: '/users/student-verify',
    method: 'POST',
    data: payload
  })
}

export const getStudentVerification = () => {
  return request({
    url: '/users/student-verify'
  })
}

export const uploadAvatar = (filePath) => {
  const token = getToken()
  return new Promise((resolve, reject) => {
    uni.uploadFile({
      url: `${BASE_URL}/users/avatar`,
      filePath,
      name: 'avatar',
      header: {
        Authorization: `Bearer ${token}`
      },
      success: (res) => {
        let body = {}
        try {
          body = JSON.parse(res.data || '{}')
        } catch (error) {
          reject(new Error('头像上传响应解析失败'))
          return
        }
        if (res.statusCode < 200 || res.statusCode >= 300 || body.code !== 0) {
          reject(new Error(body.message || '头像上传失败'))
          return
        }
        resolve(body.data)
      },
      fail: () => reject(new Error('无法连接服务器，请确认后端服务已启动'))
    })
  })
}
