<template>
  <view v-if="message" class="page">
    <view class="hero" :class="{ 'auth-hero': isStudentAuthResult || isAccountStatusChange, 'notice-hero': isPlainNotice }">
      <template v-if="isStudentAuthResult">
        <text class="auth-title">{{ studentAuthPassed ? '学生认证通过' : '学生认证未通过' }}</text>
      </template>
      <template v-else-if="isAccountStatusChange">
        <text class="auth-title">{{ accountStatusTitle }}</text>
      </template>
      <template v-else-if="isPlainNotice">
        <text class="hero-title">{{ noticeTitle }}</text>
      </template>
      <template v-else>
        <text class="eyebrow">SYSTEM MESSAGE</text>
        <text class="hero-title">消息详情</text>
        <text class="hero-copy">查看订单进度、举报处理和平台通知的完整内容</text>
      </template>
    </view>

    <view v-if="isStudentAuthResult" class="auth-actions">
      <button v-if="studentAuthPassed" class="btn primary" @click="goHome">逛首页</button>
      <button v-else class="btn primary" @click="goAppeal">去申诉</button>
    </view>

    <view v-else-if="isAccountStatusChange" class="account-status-content">
      <view class="status-reason">
        <text class="reason-label">原因</text>
        <text class="reason-value">{{ accountStatusReason }}</text>
        <text class="reason-hint">涉及的所有交易都被关闭</text>
      </view>
      <view v-if="hasRelatedTransactions" class="transactions-hint">
        <text class="hint-text">{{ transactionHint }}</text>
      </view>
      <view class="auth-actions">
        <button class="btn primary" @click="goAppeal">去申诉</button>
      </view>
    </view>

    <template v-else>
    <view class="card message-card">
      <view class="message-top">
        <view class="message-icon">系</view>
        <view class="message-head">
          <text class="title">{{ detailTitle }}</text>
          <text class="time">{{ message.createdAt || message.createTime || '暂无时间' }}</text>
        </view>
        <text :class="['status-pill', statusClass]">{{ statusLabel }}</text>
      </view>
      <text v-if="detailContent" class="content">{{ detailContent }}</text>
    </view>

    <view v-if="relatedCard" class="card related-card" @click="openRelated">
      <image v-if="relatedCard.image" class="cover" :src="relatedCard.image" mode="aspectFill" />
      <view v-else class="cover placeholder">{{ relatedCard.placeholder }}</view>
      <view class="related-body">
        <view class="related-head">
          <text class="related-label">{{ relatedCard.label }}</text>
          <text :class="['mini-status', relatedCard.statusClass]">{{ relatedCard.status }}</text>
        </view>
        <text class="related-title">{{ relatedCard.title }}</text>
        <text class="related-meta">{{ relatedCard.meta }}</text>
      </view>
      <text class="arrow">›</text>
    </view>

    <view class="actions">
      <button class="btn plain" @click="goMessages">返回消息中心</button>
    </view>
    </template>
  </view>

  <view v-else class="page loading-page">
    <view class="load-text">正在加载消息详情...</view>
  </view>
</template>

<script setup>
import { computed, ref } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { tradeService } from '../../services/trade'
import { getProductById, listMyProducts } from '../../api/product'
import { navigate, showError } from '../../utils/navigation'
import { BASE_URL } from '../../utils/request'

const message = ref(null)

const isStudentAuthResult = computed(() => (
  message.value?.title === '学生认证审核结果'
  || String(message.value?.content || '').includes('学生认证已通过')
  || String(message.value?.content || '').includes('学生认证未通过')
))
const studentAuthPassed = computed(() => isStudentAuthResult.value && !String(message.value?.content || '').includes('未通过'))
const isAccountStatusChange = computed(() => (
  message.value?.title === '账号状态变更'
  || String(message.value?.content || '').includes('账号已被')
  || String(message.value?.content || '').includes('已被临时禁用')
  || String(message.value?.content || '').includes('已被永久封禁')
))
const accountStatusTitle = computed(() => {
  const content = String(message.value?.content || '')
  if (content.includes('永久封禁')) return '账号已被永久封禁'
  if (content.includes('临时禁用') || content.includes('已被禁用')) return '账号已被临时禁用'
  return '账号状态已变更'
})
const accountStatusReason = computed(() => {
  const content = String(message.value?.content || '')
  const match = content.match(/原因[：:]\s*(.+)/)
  return match ? match[1] : content
})
const hasRelatedTransactions = computed(() => {
  const content = String(message.value?.content || '')
  if (content.includes('商品') || content.includes('订单') || content.includes('交易')) return true
  // 检查是否有相关商品或订单数据
  const product = message.value?.product || message.value?.relatedProduct
  const order = message.value?.order || message.value?.relatedOrder
  return !!(product || order)
})
const transactionHint = computed(() => {
  const content = String(message.value?.content || '')
  const parts = []
  const productMatch = content.match(/(\d+)\s*件商品已被下架/)
  if (productMatch) parts.push(`你发布的 ${productMatch[1]} 件商品已被下架`)
  const orderMatch = content.match(/(\d+)\s*个进行中的订单已被关闭/)
  if (orderMatch) parts.push(`${orderMatch[1]} 个进行中的订单已被关闭`)
  if (parts.length === 0) return '目前涉及的交易已被关闭'
  return '目前涉及的交易已被关闭，包括：' + parts.join('，')
})
const statusLabel = computed(() => statusText(message.value?.status || message.value?.readStatus || (message.value?.read ? 'READ' : 'UNREAD')))
const statusClass = computed(() => statusClassByValue(message.value?.status || message.value?.readStatus || (message.value?.read ? 'READ' : 'UNREAD')))
const isPlainNotice = computed(() => !hasConcreteRelated(message.value || {}))
const noticeTitle = computed(() => message.value?.title || '平台公告')
const detailTitle = computed(() => {
  if (isPlainNotice.value) return message.value?.content || message.value?.title || '系统消息'
  return message.value?.title || '系统消息'
})
const detailContent = computed(() => {
  if (isPlainNotice.value) return ''
  return message.value?.content || '暂无正文内容'
})

const relatedCard = computed(() => buildRelatedCard(message.value || {}))

onLoad(async (options) => {
  try {
    const data = await tradeService.getMessage(options.id)
    await hydrateRelatedProduct(data)
    message.value = data
    await tradeService.markMessageRead(options.id)
  } catch (error) {
    showError(error)
  }
})

async function hydrateRelatedProduct(item = {}) {
  const targetType = item.targetType || item.relatedType
  const targetId = item.targetId || item.relatedId
  if (targetType !== 'PRODUCT' || !targetId || item.product || item.relatedProduct) return

  try {
    item.product = await getProductById(targetId)
    return
  } catch (error) {
    // 下架商品的公开详情可能不可见，卖家本人再从“我的商品”兜底取展示信息。
  }

  try {
    const result = await listMyProducts()
    const list = Array.isArray(result) ? result : (result?.items || result?.list || [])
    const product = list.find((entry) => String(entry.id) === String(targetId))
    if (product) item.product = product
  } catch (error) {
    // 没有权限或不是卖家本人时保持占位图，不影响消息详情展示。
  }
}

function absoluteImage(url) {
  if (!url) return ''
  if (/^https?:\/\//.test(url)) return url
  return `${BASE_URL}${url}`
}

function statusText(status) {
  return {
    UNREAD: '未读',
    READ: '已读',
    PENDING: '待处理',
    PROCESSING: '处理中',
    APPROVED: '已通过',
    REJECTED: '已驳回',
    CLOSED: '已关闭',
    PENDING_CONFIRM: '待确认',
    WAIT_MEET: '待面交',
    COMPLETED: '已完成',
    CANCELED: '已取消',
    CANCELLED: '已取消',
    EXCEPTION_CLOSED: '异常关闭',
    ON_SALE: '在售',
    OFF_SHELF: '已下架',
    SOLD: '已售出',
    LOCKED: '交易中',
    VERIFIED: '已认证',
    DISABLED: '已禁用',
    BANNED: '已封禁'
  }[status] || status || '待查看'
}

function statusClassByValue(value) {
  const text = String(value || '').toUpperCase()
  if (['APPROVED', 'COMPLETED', 'READ', 'VERIFIED', 'ON_SALE'].includes(text)) return 'success'
  if (['PENDING', 'PROCESSING', 'PENDING_CONFIRM', 'WAIT_MEET', 'LOCKED', 'UNREAD'].includes(text)) return 'warning'
  if (['REJECTED', 'CLOSED', 'CANCELED', 'CANCELLED', 'EXCEPTION_CLOSED', 'OFF_SHELF', 'SOLD', 'DISABLED', 'BANNED'].includes(text)) return 'danger'
  return 'neutral'
}

function pickImage(item = {}) {
  const images = Array.isArray(item.images) ? item.images : []
  return item.image || item.productImage || item.productImageUrl || item.productCover || item.productCoverImage ||
    item.coverImage || item.coverImageUrl || item.imageUrl || item.snapshotImage || item.productImageSnapshot ||
    item.product?.image || item.product?.productImage || item.product?.productImageUrl || item.product?.productCover ||
    item.product?.productCoverImage || item.product?.coverImage || item.product?.coverImageUrl ||
    item.product?.imageUrl || item.product?.images?.[0] || images[0] || ''
}

function buildRelatedCard(item) {
  const targetType = item.targetType || item.relatedType
  const order = item.order || item.relatedOrder
  const product = item.product || item.relatedProduct
  const user = item.user || item.relatedUser
  if (!hasConcreteRelated(item)) return null

  if (targetType === 'ORDER' || order) {
    const data = order || item
    return {
      label: '相关订单',
      title: data.productTitleSnapshot || data.product?.title || item.targetTitle || item.title || '订单详情',
      meta: `${statusText(data.status)} · ¥${data.productPriceSnapshot || data.product?.price || item.productPrice || ''}`,
      image: absoluteImage(pickImage(data)),
      status: statusText(data.status || item.status),
      statusClass: statusClassByValue(data.status || item.status),
      placeholder: '单',
      targetType: 'ORDER',
      targetId: order ? (data.id || item.targetId) : item.targetId
    }
  }

  if (targetType === 'PRODUCT' || product) {
    const data = product || item
    return {
      label: '相关商品',
      title: data.title || item.targetTitle || '商品详情',
      meta: `¥${data.price || item.productPrice || ''} · ${statusText(data.status)}`,
      image: absoluteImage(pickImage(data)),
      status: statusText(data.status || item.status),
      statusClass: statusClassByValue(data.status || item.status),
      placeholder: '物',
      targetType: 'PRODUCT',
      targetId: product ? (data.id || item.targetId) : item.targetId
    }
  }

  if (targetType === 'USER' || user) {
    const data = user || item
    return {
      label: '相关用户',
      title: data.nickname || data.realName || item.targetTitle || '用户主页',
      meta: statusText(data.authStatus || data.accountStatus),
      image: absoluteImage(data.avatarUrl || data.avatar_url),
      status: statusText(data.authStatus || data.accountStatus),
      statusClass: statusClassByValue(data.authStatus || data.accountStatus),
      placeholder: '人',
      targetType: 'USER',
      targetId: user ? (data.id || item.targetId) : item.targetId
    }
  }

  return null
}

function hasConcreteRelated(item = {}) {
  const targetType = item.targetType || item.relatedType
  return ['ORDER', 'PRODUCT', 'USER'].includes(targetType) ||
    Boolean(item.order || item.relatedOrder || item.product || item.relatedProduct || item.user || item.relatedUser)
}

function openRelated() {
  const related = relatedCard.value
  if (!related?.targetId) return
  if (related.targetType === 'ORDER') navigate('/pages/order/detail', { id: related.targetId })
  else if (related.targetType === 'PRODUCT') navigate('/pages/detail/detail', { id: related.targetId })
  else if (related.targetType === 'USER') navigate('/pages/user-profile/user-profile', { id: related.targetId })
}

function goMessages() {
  uni.navigateBack({ delta: 1 })
}

function goHome() {
  uni.switchTab({ url: '/pages/home/home' })
}

function goAppeal() {
  navigate('/pages/interaction/appeal', { targetType: 'USER' })
}
</script>

<style scoped>
.page { min-height: 100vh; padding: 30rpx 28rpx 48rpx; background: #f5f8f6; box-sizing: border-box; }
.loading-page { display: flex; align-items: center; justify-content: center; color: #667085; }
.hero { padding: 34rpx 30rpx; border-radius: 28rpx; background: linear-gradient(135deg, #23734f, #3e9b72); color: #fff; box-shadow: 0 12rpx 32rpx rgba(35,115,79,.18); }
.auth-hero { display: flex; min-height: 172rpx; align-items: center; justify-content: center; text-align: center; }
.notice-hero { display: flex; min-height: 112rpx; align-items: center; justify-content: center; margin-bottom: 24rpx; text-align: center; }
.auth-title { color: #fff; font-size: 42rpx; font-weight: 800; }
.eyebrow, .hero-title, .hero-copy, .title, .time, .content, .related-label, .related-title, .related-meta { display: block; }
.eyebrow { color: rgba(255,255,255,.72); font-size: 20rpx; letter-spacing: 2rpx; }
.hero-title { font-size: 40rpx; font-weight: 800; }
.hero-copy { margin-top: 10rpx; color: rgba(255,255,255,.78); font-size: 24rpx; line-height: 1.5; }
.card { margin-top: 24rpx; padding: 30rpx; border-radius: 24rpx; background: #fff; box-shadow: 0 10rpx 30rpx rgba(28,68,52,.06); box-sizing: border-box; }
.message-top { display: flex; align-items: flex-start; gap: 20rpx; }
.message-icon { display: flex; width: 72rpx; height: 72rpx; flex-shrink: 0; align-items: center; justify-content: center; border-radius: 50%; background: #e7f4ec; color: #23734f; font-size: 28rpx; font-weight: 800; }
.message-head { flex: 1; min-width: 0; }
.title { color: #26342f; font-size: 36rpx; font-weight: 800; line-height: 1.45; }
.time { margin-top: 12rpx; color: #98a39d; font-size: 24rpx; }
.content { margin-top: 34rpx; padding-top: 30rpx; border-top: 1rpx solid #eef2f0; color: #425148; font-size: 29rpx; line-height: 1.9; white-space: pre-line; }
.status-pill, .mini-status { flex-shrink: 0; padding: 8rpx 16rpx; border-radius: 999rpx; font-size: 22rpx; }
.success { background: #e9f7ef; color: #168451; }
.warning { background: #fff4df; color: #9a6b24; }
.danger { background: #ffecea; color: #c23b2d; }
.neutral { background: #eef2f0; color: #526158; }
.related-card { display: flex; align-items: center; gap: 20rpx; }
.cover { width: 120rpx; height: 120rpx; flex-shrink: 0; border-radius: 16rpx; background: #e8efeb; }
.placeholder { display: flex; align-items: center; justify-content: center; color: #7a8780; font-size: 32rpx; font-weight: 800; }
.related-body { flex: 1; min-width: 0; }
.related-head { display: flex; align-items: center; justify-content: space-between; gap: 12rpx; }
.related-label { color: #667085; font-size: 22rpx; }
.related-title { margin-top: 10rpx; color: #26342f; font-size: 29rpx; font-weight: 800; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.related-meta { margin-top: 8rpx; color: #89938f; font-size: 23rpx; }
.arrow { color: #98a39d; font-size: 44rpx; }
.result-row { display: flex; justify-content: space-between; gap: 24rpx; padding: 12rpx 0; }
.result-row.multiline { align-items: flex-start; }
.result-label { flex-shrink: 0; color: #667085; font-size: 25rpx; }
.result-value { color: #26342f; font-size: 26rpx; line-height: 1.6; text-align: right; }
.actions { display: flex; gap: 18rpx; margin-top: 28rpx; }
.auth-actions { margin-top: 34rpx; }
.account-status-content { margin-top: 34rpx; }
.status-reason { padding: 30rpx; border-radius: 24rpx; background: #fff; box-shadow: 0 10rpx 30rpx rgba(28,68,52,.06); }
.reason-label { display: block; color: #667085; font-size: 24rpx; margin-bottom: 12rpx; }
.reason-value { display: block; color: #26342f; font-size: 30rpx; font-weight: 700; line-height: 1.5; }
.reason-hint { display: block; margin-top: 16rpx; color: #89938f; font-size: 24rpx; }
.transactions-hint { margin-top: 20rpx; padding: 20rpx 30rpx; }
.hint-text { color: #89938f; font-size: 24rpx; line-height: 1.6; }
.btn { flex: 1; height: 78rpx; line-height: 78rpx; border-radius: 999rpx; font-size: 26rpx; }
.btn::after { border: 0; }
.primary { background: #23734f; color: #fff; }
.plain { background: #fff; color: #526158; }
</style>
