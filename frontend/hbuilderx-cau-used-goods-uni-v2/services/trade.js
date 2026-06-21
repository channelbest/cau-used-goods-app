import { BASE_URL } from '../utils/request'
import * as api from '../api/trade'
import { getAdminOrders } from '../api/admin'
import { REPORT_REASON, TARGET_TYPE } from '../utils/constants'
import { displayRelatedUserName } from '../utils/user-format'

function absoluteImage(url) {
  if (!url || /^https?:\/\//.test(url)) return url
  return `${BASE_URL}${url}`
}

function pad(value) {
  return String(value).padStart(2, '0')
}

function formatDateTime(value) {
  if (!value) return ''
  if (typeof value === 'string' && /^\d{4}-\d{2}-\d{2} \d{2}:\d{2}/.test(value)) {
    return value.slice(0, 16)
  }
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) return value
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())} ${pad(date.getHours())}:${pad(date.getMinutes())}`
}

function normalizeProduct(item = {}) {
  return {
    ...item,
    image: absoluteImage(pickProductImage(item)),
    meetLocation: item.meetLocation || '预约后协商'
  }
}

function pickProductImage(item = {}) {
  const images = Array.isArray(item.images) ? item.images : []
  return item.image
    || item.productImage
    || item.productImageUrl
    || item.productCover
    || item.productCoverImage
    || item.coverImage
    || item.coverImageUrl
    || item.imageUrl
    || item.snapshotImage
    || item.productImageSnapshot
    || item.product?.image
    || item.product?.productImage
    || item.product?.coverImage
    || item.product?.imageUrl
    || item.product?.images?.[0]
    || images[0]
    || ''
}

function normalizeOrder(item = {}) {
  const productImage = pickProductImage(item)
  return {
    ...item,
    id: String(item.id),
    sellerName: displayRelatedUserName(item, 'seller', '卖家'),
    buyerName: displayRelatedUserName(item, 'buyer', '买家'),
    createdAt: formatDateTime(item.createdAt || item.createTime),
    meetTime: formatDateTime(item.meetTime),
    expireTime: formatDateTime(item.expireTime),
    confirmTime: formatDateTime(item.confirmTime),
    finishTime: formatDateTime(item.finishTime),
    closeTime: formatDateTime(item.closeTime),
    product: normalizeProduct({
      ...(item.product || {}),
      id: item.product?.id || item.productId,
      title: item.product?.title || item.productTitleSnapshot,
      price: item.product?.price || item.productPriceSnapshot,
      image: productImage,
      meetLocation: item.product?.meetLocation || item.meetLocation
    })
  }
}

function normalizeMessage(item = {}) {
  return {
    ...item,
    id: String(item.id),
    type: item.type || item.messageType,
    targetType: item.targetType || item.relatedType,
    targetId: item.targetId || item.relatedId,
    createdAt: formatDateTime(item.createdAt || item.createTime),
    read: item.read ?? item.readStatus === 'READ'
  }
}

function normalizeAppeal(item = {}) {
  return {
    ...item,
    id: String(item.id),
    targetTypeLabel: TARGET_TYPE[item.targetType] || item.targetType,
    result: item.handleResult,
    createdAt: formatDateTime(item.createTime || item.createdAt),
    updatedAt: formatDateTime(item.updateTime || item.updatedAt),
    handleTime: formatDateTime(item.handleTime)
  }
}

function normalizeReport(item = {}) {
  return {
    ...item,
    id: String(item.id),
    reason: item.reasonType,
    reasonLabel: REPORT_REASON[item.reasonType] || item.reasonType,
    detail: item.description,
    result: item.handleResult,
    targetTypeLabel: TARGET_TYPE[item.targetType] || item.targetType,
    createdAt: formatDateTime(item.createTime || item.createdAt),
    handleTime: formatDateTime(item.handleTime)
  }
}

async function getOrderWithImage(id) {
  const detail = normalizeOrder(await api.getOrder(id))
  if (detail.product?.image) return detail

  for (const role of ['buyer', 'seller']) {
    try {
      const result = await api.getOrders(role, { pageSize: 100 })
      const matched = (result?.items || []).find((item) => String(item.id) === String(id))
      if (!matched) continue

      const orderFromList = normalizeOrder(matched)
      if (!orderFromList.product?.image) continue
      return {
        ...detail,
        productImage: orderFromList.product.image,
        product: {
          ...detail.product,
          ...orderFromList.product,
          image: orderFromList.product.image
        }
      }
    } catch (error) {
      // 当前用户可能只拥有买家或卖家其中一种订单列表权限，继续尝试另一种角色。
    }
  }

  return detail
}

export const tradeService = {
  getProduct: async (id) => normalizeProduct(await api.getProduct(id)),
  createAppointment: async (data) => normalizeOrder(await api.createAppointment(data)),
  getOrders: async (role) => (await api.getOrders(role)).items.map(normalizeOrder),
  getOrder: getOrderWithImage,
  getAdminOrder: async (id) => {
    const result = await getAdminOrders('ALL', { pageSize: 200 })
    const list = result?.items || []
    const order = list.find((item) => Number(item.id) === Number(id))
    if (!order) throw new Error('订单不存在或不在当前管理员列表中')
    return normalizeOrder(order)
  },
  changeOrderStatus: async (id, action, data) => normalizeOrder(await api.changeOrderStatus(id, action, data)),
  getFavorites: async () => (await api.getFavorites()).items.map((item) => normalizeProduct({
    id: item.productId,
    title: item.productTitle,
    price: item.productPrice,
    status: item.productStatus,
    image: item.productImage,
    sellerName: item.sellerNickname
  })),
  addFavorite: api.addFavorite,
  removeFavorite: api.removeFavorite,
  getMessages: async () => (await api.getMessages()).items.map(normalizeMessage),
  getMessage: async (id) => normalizeMessage(await api.getMessage(id)),
  markMessageRead: api.markMessageRead,
  deleteMessage: api.deleteMessage,
  deleteAllMessages: api.deleteAllMessages,
  createReview: api.createReview,
  createReport: api.createReport,
  getReports: async () => (await api.getReports()).items.map(normalizeReport),
  getReport: async (id) => normalizeReport(await api.getReport(id)),
  createAppeal: api.createAppeal,
  getAppeals: async (params) => {
    const result = await api.getAppeals(params)
    return {
      ...result,
      items: (result?.items || []).map(normalizeAppeal)
    }
  },
  getAppeal: async (id) => normalizeAppeal(await api.getAppeal(id)),
  closeAppeal: api.closeAppeal
}
