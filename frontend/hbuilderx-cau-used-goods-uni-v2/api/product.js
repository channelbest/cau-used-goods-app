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

export const listMyProducts = () => request({
  url: '/products/my'
})

export const getProductById = (id) => request({
  url: `/products/${id}`
})

export const updateProductStatus = (id, status, reason = '') => request({
  url: `/products/${id}/status`,
  method: 'PUT',
  data: { status, reason }
})

export const batchPutOnSaleProducts = () => request({
  url: '/products/batch-on-sale',
  method: 'POST'
})

export const deleteProduct = (id) => request({
  url: `/products/${id}`,
  method: 'DELETE'
})

export const updateProduct = (id, payload) => request({
  url: `/products/${id}`,
  method: 'PUT',
  data: payload
})

export const addProductImages = (id, images = []) => request({
  url: `/products/${id}/images`,
  method: 'POST',
  data: { images }
})

export const replaceProductImages = (id, images = []) => request({
  url: `/products/${id}/images`,
  method: 'PUT',
  data: { images }
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

export const listMyOrders = (params = {}) => request({
  url: `/orders${buildQuery(params)}`
})

export const confirmOrder = (id) => request({
  url: `/orders/${id}/confirm`,
  method: 'POST'
})

export const completeOrder = (id) => request({
  url: `/orders/${id}/complete`,
  method: 'POST'
})

export const cancelOrder = (id, reason) => request({
  url: `/orders/${id}/cancel`,
  method: 'POST',
  data: { reason }
})

export const createReport = ({ productId, reason }) => request({
  url: '/reports',
  method: 'POST',
  data: {
    targetType: 'PRODUCT',
    targetId: productId,
    reasonType: 'OTHER',
    description: reason
  }
})

export const createProduct = async (payload) => {
  const { images = [], ...product } = payload
  const result = await request({
    url: '/products',
    method: 'POST',
    data: product
  })
  if (images.length) {
    await addProductImages(result.id, images)
  }
  return result
}

export const uploadProductImage = async (filePath) => {
  const result = await uploadFile({
    url: '/upload/image',
    filePath
  })
  return {
    ...result,
    imageUrl: result.imageUrl || result.url
  }
}

function buildAiDescription({ description = '', categoryName = '', conditionLevel = '', meetLocation = '' } = {}) {
  const details = []
  if (description) details.push(`商品描述：${description}`)
  if (categoryName) details.push(`商品分类：${categoryName}`)
  if (conditionLevel) details.push(`商品成色：${conditionLevel}`)
  if (meetLocation) details.push(`建议面交地点：${meetLocation}`)
  return details.join('\n')
}

export const optimizeProductTitle = async (payload = {}) => {
  const { title } = payload
  const result = await request({
    url: '/ai/optimize-product',
    method: 'POST',
    data: {
      title,
      description: buildAiDescription(payload)
    }
  })
  return {
    titles: result.optimizedTitle ? [result.optimizedTitle] : [title]
  }
}

export const generateProductDescription = async (payload = {}) => {
  const { title, description = '' } = payload
  const result = await request({
    url: '/ai/optimize-product',
    method: 'POST',
    data: {
      title,
      description: buildAiDescription(payload)
    }
  })
  return {
    description: result.optimizedDescription || description
  }
}
