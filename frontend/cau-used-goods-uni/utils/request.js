import { getToken, clearAuth } from './auth'

export const BASE_URL = 'http://127.0.0.1:8080'

const getErrorMessage = (body, fallback) => {
  return body && body.message ? body.message : fallback
}

export const request = ({
  url,
  method = 'GET',
  data = {},
  header = {},
  auth = true
}) => {
  const token = getToken()
  const requestHeader = {
    'Content-Type': 'application/json',
    ...header
  }

  if (auth && token) {
    requestHeader.Authorization = `Bearer ${token}`
  }

  return new Promise((resolve, reject) => {
    uni.request({
      url: `${BASE_URL}${url}`,
      method,
      data,
      header: requestHeader,
      success: (res) => {
        const body = res.data || {}

        if (res.statusCode === 401) {
          clearAuth()
          reject(new Error(getErrorMessage(body, '登录已过期，请重新登录')))
          return
        }

        if (res.statusCode < 200 || res.statusCode >= 300 || body.code !== 0) {
          reject(new Error(getErrorMessage(body, '请求失败，请稍后重试')))
          return
        }

        resolve(body.data)
      },
      fail: () => {
        reject(new Error('无法连接服务器，请确认后端服务已启动'))
      }
    })
  })
}

export const uploadFile = ({
  url,
  filePath,
  name = 'file',
  formData = {},
  auth = true
}) => {
  const token = getToken()
  const header = {}
  if (auth && token) header.Authorization = `Bearer ${token}`

  return new Promise((resolve, reject) => {
    uni.uploadFile({
      url: `${BASE_URL}${url}`,
      filePath,
      name,
      formData,
      header,
      success: (res) => {
        let body = {}
        try {
          body = JSON.parse(res.data || '{}')
        } catch (error) {
          reject(new Error('上传响应解析失败'))
          return
        }
        if (res.statusCode === 401) clearAuth()
        if (res.statusCode < 200 || res.statusCode >= 300 || body.code !== 0) {
          reject(new Error(getErrorMessage(body, '图片上传失败，请重新上传')))
          return
        }
        resolve(body.data)
      },
      fail: () => reject(new Error('图片上传失败，请检查网络后重试'))
    })
  })
}
