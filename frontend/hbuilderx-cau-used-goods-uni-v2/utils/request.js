import { getToken, clearAuth } from './auth'

//export const BASE_URL = 'http://62.234.163.176:7001'
export const BASE_URL = 'http://127.0.0.1:8080'

export const normalizeErrorMessage = (message, fallback = '请求失败，请稍后重试') => {
  const value = String(message || '').trim()
  if (!value) return fallback
  const lower = value.toLowerCase()
  if (lower.includes('already have an active order')) return '你已预约过该商品，不能重复提交预约'
  if (lower.includes('cannot buy your own product')) return '不能预约自己发布的商品'
  if (lower.includes('you already have an active appeal')) return '该对象已有待处理或处理中的申诉，请勿重复提交'
  if (lower.includes('appeal target not found')) return '申诉对象不存在或已被删除'
  if (lower.includes('appeal cannot be closed')) return '该申诉当前状态不能撤回'
  if (lower.includes('appeal already handled')) return '该申诉已处理，不能重复操作'
  if (lower.includes('appeal must be marked processing')) return '请先将申诉标记为处理中'
  if (lower.includes('report result can only be appealed after handled')) return '举报处理完成后才能申诉处理结果'
  if (lower.includes('user appeal target must be restored before approval')) return '账号解封申诉需先恢复账号，再标记为通过'
  if (lower.includes('not found')) return '未找到相关数据'
  if (lower.includes('already reviewed')) return '该订单已经评价过'
  if (lower.includes('unauthorized') || lower.includes('token')) return '登录已过期，请重新登录'
  if (lower.includes('forbidden') || lower.includes('permission')) return '没有权限进行该操作'
  if (lower.includes('invalid')) return '提交内容格式不正确'
  if (lower.includes('duplicate')) return '数据已存在，请勿重复提交'
  if (lower.includes('failed') || lower.includes('error')) return fallback
  if (/^[\x00-\x7F]+$/.test(value)) return fallback
  return value
}

const getErrorMessage = (body, fallback) => {
  return normalizeErrorMessage(body && body.message, fallback)
}

const buildQuery = (data = {}) => {
  const parts = []
  Object.keys(data || {}).forEach((key) => {
    const value = data[key]
    if (value === undefined || value === null || value === '') return
    parts.push(`${encodeURIComponent(key)}=${encodeURIComponent(value)}`)
  })
  return parts.length ? `?${parts.join('&')}` : ''
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

  const finalURL = method === 'GET' && data && Object.keys(data).length
    ? `${BASE_URL}${url}${url.includes('?') ? '&' + buildQuery(data).slice(1) : buildQuery(data)}`
    : `${BASE_URL}${url}`

  return new Promise((resolve, reject) => {
    uni.request({
      url: finalURL,
      method,
      data: method === 'GET' ? {} : data,
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

export const uploadImage = async (filePath) => {
  const result = await uploadFile({
    url: '/upload/image',
    filePath
  })
  return result?.url || result?.imageUrl || result?.path || result
}
