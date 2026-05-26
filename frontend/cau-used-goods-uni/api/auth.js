import { request } from '../utils/request'

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

export const getCurrentUser = () => {
  return request({
    url: '/users/me'
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
