const BASE_URL = 'http://127.0.0.1:8080'

export function normalizeImage(url) {
  if (!url) return ''
  if (url.indexOf('http://') === 0 || url.indexOf('https://') === 0) return url
  if (url.indexOf('/uploads/') === 0) return BASE_URL + url
  return url
}

export function formatPrice(price) {
  const n = Number(price || 0)
  return Number.isInteger(n) ? String(n) : n.toFixed(2)
}

export function formatProduct(item, categoryMap) {
  const categoryName = categoryMap && categoryMap[item.categoryId] ? categoryMap[item.categoryId] : '未分类'
  const images = (item.images || []).map(normalizeImage)
  const image = images.length ? images[0] : ''
  return {
    ...item,
    images,
    category: categoryName,
    priceText: formatPrice(item.price),
    conditionText: item.conditionLevel || '成色未填写',
    timeText: item.createTime || '',
    coverImage: image
  }
}

export function buildCategoryMap(categories) {
  const map = {}
  ;(categories || []).forEach((item) => {
    map[item.id] = item.name
  })
  return map
}

export function withAllCategory(categories) {
  const list = categories || []
  return list.some((item) => Number(item.id) === 0) ? list : [{ id: 0, name: '全部' }, ...list]
}

export function getStatusText(status) {
  const statusMap = {
    ON_SALE: '在售',
    LOCKED: '已被预约',
    SOLD: '已售出',
    OFF_SHELF: '已下架'
  }
  return statusMap[status] || status || '未知状态'
}
