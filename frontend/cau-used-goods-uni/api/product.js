import { request, uploadFile } from '../utils/request'

function buildQuery(params = {}) {
  const parts = []
  Object.keys(params).forEach((key) => {
    const value = params[key]
    if (value === undefined || value === null || value === '') return
    parts.push(`${encodeURIComponent(key)}=${encodeURIComponent(value)}`)
  })
  return parts.length ? `?${parts.join('&')}` : ''
}

export const listCategories = () => request({
  url: '/categories',
  auth: false
})

export const listProducts = (params = {}) => request({
  url: `/products${buildQuery(params)}`,
  auth: false
})

export const getProductById = (id) => request({
  url: `/products/${id}`,
  auth: false
})

export const addFavorite = (productId) => request({
  url: '/favorites',
  method: 'POST',
  data: { productId }
})

export const removeFavorite = (productId) => request({
  url: `/favorites/${productId}`,
  method: 'DELETE'
})

export const checkFavorite = (productId) => request({
  url: `/favorites/check${buildQuery({ productId })}`
})

export const createOrder = (payload) => request({
  url: '/orders',
  method: 'POST',
  data: payload
})

export const createReport = (payload) => request({
  url: '/reports',
  method: 'POST',
  data: payload
})

export const createProduct = (payload) => request({
  url: '/products',
  method: 'POST',
  data: payload
})

export const uploadProductImage = (filePath) => uploadFile({
  url: '/upload/products',
  filePath
})

export const optimizeProductTitle = (payload) => request({
  url: '/ai/optimize-title',
  method: 'POST',
  data: payload
})

export const generateProductDescription = (payload) => request({
  url: '/ai/generate-description',
  method: 'POST',
  data: payload
})
