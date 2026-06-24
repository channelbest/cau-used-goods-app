const BASE_URL = 'http://62.234.163.176:7001'
// const BASE_URL = 'http://127.0.0.1:8080'
const INVALID_TEST_IMAGES = [
  '1780413573623001600.png',
  '1781224931041565400.png',
  '1781226398817582900.png',
  'integration-demo.png'
]

export function normalizeImage(url) {
  if (!url) return ''
  const value = String(url).trim()
  if (!value) return ''
  if (INVALID_TEST_IMAGES.some((name) => value.includes(name))) return ''
  if (value.indexOf('http://') === 0 || value.indexOf('https://') === 0) return value
  if (value.indexOf('/uploads/') === 0) return BASE_URL + value
  if (value.indexOf('uploads/') === 0) return `${BASE_URL}/${value}`
  return value
}

export function formatPrice(price) {
  const n = Number(price || 0)
  return Number.isInteger(n) ? String(n) : n.toFixed(2)
}

export function formatProduct(item, categoryMap) {
  const categoryId = item.categoryId || item.category_id
  const categoryName = categoryMap && categoryMap[categoryId] ? categoryMap[categoryId] : '\u672a\u5206\u7c7b'
  const rawImages = item.images || item.imageUrls || item.image_urls || item.productImages || []
  const imageList = Array.isArray(rawImages) ? rawImages : String(rawImages || '').split(',').filter(Boolean)
  const normalizedImages = imageList.map(normalizeImage).filter(Boolean)
  const coverImage = normalizeImage(item.coverImage || item.cover_image || item.image || item.productImage || normalizedImages[0] || '')
  const images = coverImage && !normalizedImages.length ? [coverImage] : normalizedImages

  return {
    ...item,
    categoryId,
    images,
    category: categoryName,
    priceText: formatPrice(item.price),
    conditionText: item.conditionLevel || item.condition_level || '\u6210\u8272\u672a\u586b\u5199',
    timeText: item.createTime || item.create_time || '',
    coverImage
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
  return list.some((item) => Number(item.id) === 0) ? list : [{ id: 0, name: '\u5168\u90e8' }, ...list]
}

export function getStatusText(status) {
  const statusMap = {
    ON_SALE: '\u5728\u552e',
    LOCKED: '\u5df2\u88ab\u9884\u7ea6',
    SOLD: '\u5df2\u552e\u51fa',
    OFF_SHELF: '\u5df2\u4e0b\u67b6'
  }
  return statusMap[status] || status || '\u672a\u77e5\u72b6\u6001'
}
