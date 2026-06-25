import { request, uploadFile } from '../utils/request'

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
  return uploadFile({
    url: '/users/avatar',
    filePath,
    name: 'avatar'
  })
}
