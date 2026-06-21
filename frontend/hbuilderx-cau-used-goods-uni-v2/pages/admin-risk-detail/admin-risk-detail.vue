<template>
  <view class="page">
    <view v-if="loading" class="empty">加载中</view>
    <view v-else-if="!item" class="empty">未找到对应记录</view>

    <view v-else class="detail-card">
      <view class="card-head">
        <view>
          <view class="title">{{ itemTitle(item) }}</view>
          <view class="sub">{{ targetText(item.targetType) }} #{{ item.targetId }}</view>
        </view>
        <view class="badge-row">
          <view class="status-badge" :style="statusBadgeStyle(item.status)">{{ statusText(item.status) }}</view>
          <view v-if="isUserAppeal" class="risk-badge">{{ userAppealRiskText }}</view>
        </view>
      </view>

      <view class="desc">{{ itemDescription(item) }}</view>
      <view class="meta-row">
        <text>{{ mode === 'REPORT' ? '举报人' : '申诉人' }}：{{ actorName(item) }}</text>
        <text>{{ shortTime(item.createTime) }}</text>
      </view>

      <view v-if="targetProduct" class="target-product" @click="goProduct(targetProduct.id)">
        <image v-if="productCover(targetProduct)" class="target-image" :src="productCover(targetProduct)" mode="aspectFill" />
        <view v-else class="target-image placeholder">商品</view>
        <view class="target-main">
          <view class="target-title">{{ targetProduct.title || '商品' }}</view>
          <view class="target-meta">￥{{ targetProduct.price || 0 }} · {{ productStatusText(targetProduct.status) }}</view>
        </view>
        <text class="target-arrow">›</text>
      </view>

      <view v-else-if="targetEntryVisible" class="target-product" @click="goTargetManage">
        <image v-if="targetEntryImage" class="target-image" :src="targetEntryImage" mode="aspectFill" />
        <view v-else class="target-image placeholder">{{ targetEntryPlaceholder }}</view>
        <view class="target-main">
          <view class="target-title">{{ targetEntryTitle }}</view>
          <view class="target-meta">{{ targetEntryMeta }}</view>
        </view>
        <text class="target-arrow">›</text>
      </view>

      <view class="section-title">具体信息</view>
      <template v-if="item.targetType === 'PRODUCT'">
        <view class="detail-row"><text>商品名称</text><text>{{ productInfoTitle }}</text></view>
        <view class="detail-row"><text>商品状态</text><text>{{ productInfoStatus }}</text></view>
        <view class="detail-row"><text>商品价格</text><text>￥{{ productInfoPrice }}</text></view>
      </template>
      <template v-else>
        <view class="detail-row"><text>对象类型</text><text>{{ targetText(item.targetType) }}</text></view>
        <view class="detail-row"><text>对象 ID</text><text>#{{ item.targetId }}</text></view>
      </template>
      <view class="detail-row"><text>当前状态</text><text>{{ statusText(item.status) }}</text></view>
      <view class="detail-row"><text>{{ mode === 'REPORT' ? '举报原因' : '申诉理由' }}</text><text>{{ itemTitle(item) }}</text></view>
      <view v-if="item.handleTime" class="detail-row"><text>处理时间</text><text>{{ shortTime(item.handleTime) }}</text></view>

      <view v-if="item.images?.length" class="image-grid">
        <image v-for="url in item.images" :key="url" class="case-image" :src="normalizeImage(url)" mode="aspectFill" @click="previewImage(url, item.images)" />
      </view>

      <view v-if="item.handleResult" class="handle-result">{{ item.handleResult }}</view>

      <view v-if="canHandle(item.status)" class="actions">
        <button v-if="item.status === 'PENDING'" size="mini" class="process" @click="handleCurrent('PROCESSING')">开始处理</button>
        <button v-if="item.status === 'PROCESSING'" size="mini" class="pass" @click="handleCurrent('APPROVED')">{{ mode === 'REPORT' ? '处理完成' : '通过申诉' }}</button>
        <button v-if="item.status === 'PROCESSING'" size="mini" class="reject" @click="openReasonModal">{{ mode === 'REPORT' ? '驳回' : '驳回申诉' }}</button>
      </view>
    </view>

    <view v-if="reasonModal.visible" class="modal-mask" @click="closeReasonModal">
      <view class="reason-sheet" @click.stop>
        <view class="sheet-title">{{ mode === 'REPORT' ? '驳回举报' : '驳回申诉' }}</view>
        <view class="reason-list">
          <view v-for="reason in reasonOptions" :key="reason" :class="['reason-chip', reasonModal.reason === reason ? 'active' : '']" @click="reasonModal.reason = reason">
            {{ reason }}
          </view>
        </view>
        <textarea v-model="reasonModal.note" class="reason-input" maxlength="200" placeholder="补充说明，可不填" placeholder-class="reason-placeholder" />
        <view class="sheet-actions">
          <button class="sheet-button cancel" @click="closeReasonModal">取消</button>
          <button class="sheet-button confirm" :loading="submitting" @click="submitReject">确认处理</button>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { computed, reactive, ref } from 'vue'
import { onLoad, onShow } from '@dcloudio/uni-app'
import {
  getAdminAppeals,
  getAdminOrders,
  getAdminReports,
  handleAdminAppeal,
  handleAdminReport,
  markAdminAppealProcessing,
  markAdminReportProcessing,
  getAdminProductById
} from '../../api/admin'
import { getPublicProfile } from '../../api/user'
import { normalizeImage } from '../../utils/product-format'
import { displayRelatedUserName } from '../../utils/user-format'

const mode = ref('REPORT')
const id = ref(0)
const item = ref(null)
const product = ref(null)
const order = ref(null)
const targetUser = ref(null)
const loading = ref(false)
const submitting = ref(false)
const reasonModal = reactive({ visible: false, reason: '', note: '' })

const reasonOptions = computed(() => mode.value === 'REPORT'
  ? ['证据不足，无法认定违规', '举报内容与对象不符', '未发现明显违规行为', '重复举报或恶意举报', '商品信息已修正，风险已解除']
  : ['申诉材料不足', '原举报处理结论无误', '未提供有效证明', '申诉理由与处理结果无关', '存在重复申诉'])

const targetProduct = computed(() => item.value?.targetType === 'PRODUCT' ? product.value : null)
const targetOrder = computed(() => item.value?.targetType === 'ORDER' ? order.value : null)
const isUserAppeal = computed(() => mode.value === 'APPEAL' && item.value?.targetType === 'USER')
const userAppealRiskText = computed(() => {
  const status = normalizeStatus(item.value?.status)
  if (status === 'PENDING' || status === 'PROCESSING') return '需超管'
  return '账号解封'
})
const productInfoTitle = computed(() => targetProduct.value?.title || `商品 #${item.value?.targetId || ''}`)
const productInfoPrice = computed(() => targetProduct.value?.price ?? 0)
const productInfoStatus = computed(() => productStatusText(targetProduct.value?.status))
const targetEntryVisible = computed(() => ['USER', 'PRODUCT', 'ORDER'].includes(item.value?.targetType) || Boolean(targetManageUrl.value))
const targetUserAvatar = computed(() => item.value?.targetType === 'USER' ? normalizeImage(targetUser.value?.avatarUrl || targetUser.value?.avatar || targetUser.value?.avatar_url || '') : '')
const targetUserName = computed(() => targetUser.value?.nickname || targetUser.value?.realName || targetUser.value?.real_name || '')
const targetEntryImage = computed(() => {
  if (item.value?.targetType === 'USER') return targetUserAvatar.value
  if (item.value?.targetType === 'ORDER') return orderProductImage(targetOrder.value)
  return ''
})
const targetEntryPlaceholder = computed(() => {
  if (item.value?.targetType === 'USER') return targetUserName.value ? targetUserName.value.slice(0, 1) : '用户'
  if (item.value?.targetType === 'ORDER') return '订单'
  return targetText(item.value?.targetType)
})
const targetEntryTitle = computed(() => {
  if (item.value?.targetType === 'USER') return '查看用户主页'
  if (item.value?.targetType === 'PRODUCT') return '查看商品详情'
  if (item.value?.targetType === 'ORDER') return '查看订单详情'
  return `查看${targetText(item.value?.targetType)}管理`
})
const targetEntryMeta = computed(() => {
  if (item.value?.targetType === 'USER') {
    return targetUserName.value || '用户主页'
  }
  if (item.value?.targetType === 'ORDER') {
    const buyerName = displayRelatedUserName(targetOrder.value || {}, 'buyer', '买家')
    const productTitle = targetOrder.value?.productTitleSnapshot || targetOrder.value?.product?.title || '订单商品'
    return `${buyerName} · ${productTitle}`
  }
  return `${targetText(item.value?.targetType)} #${item.value?.targetId}`
})
const targetManageUrl = computed(() => {
  if (!item.value) return ''
  const map = {
  }
  return map[item.value.targetType] || ''
})

onLoad((query) => {
  mode.value = String(query.mode || 'REPORT').toUpperCase()
  id.value = Number(query.id || 0)
})

onShow(() => {
  load()
})

const load = async () => {
  if (!id.value) return
  loading.value = true
  try {
    const result = mode.value === 'APPEAL' ? await getAdminAppeals() : await getAdminReports()
    const list = result?.items || []
    item.value = list.find((record) => Number(record.id) === id.value) || null
    product.value = null
    order.value = null
    targetUser.value = null
    if (item.value?.targetType === 'PRODUCT' && item.value.targetId) {
      try {
        product.value = await getAdminProductById(item.value.targetId)
      } catch (error) {
        product.value = null
      }
    }
    if (item.value?.targetType === 'USER' && item.value.targetId) {
      targetUser.value = await getPublicProfile(item.value.targetId).catch(() => null)
    }
    if (item.value?.targetType === 'ORDER' && item.value.targetId) {
      order.value = await findAdminOrder(item.value.targetId).catch(() => null)
    }
  } catch (error) {
    uni.showToast({ title: error.message || '加载失败', icon: 'none' })
  } finally {
    loading.value = false
  }
}

const targetText = (targetType) => ({ PRODUCT: '商品', USER: '用户', ORDER: '订单', REPORT: '举报' }[targetType] || targetType || '对象')
const normalizeStatus = (status) => String(status || '').trim().toUpperCase()
const statusText = (status) => ({ PENDING: '待处理', PROCESSING: '处理中', RESOLVED: '已处理', APPROVED: '已通过', REJECTED: '已驳回', CLOSED: '已关闭' }[normalizeStatus(status)] || status || '未知')
const statusTone = (status) => {
  const value = normalizeStatus(status)
  if (value === 'APPROVED' || value === 'RESOLVED') return 'success'
  if (value === 'PENDING') return 'warning'
  if (value === 'PROCESSING') return 'primary'
  if (value === 'REJECTED') return 'danger'
  if (value === 'CLOSED') return 'muted'
  return 'muted'
}
const statusStyleMap = {
  warning: { background: '#fff7e6', color: '#b66a00' },
  primary: { background: '#eff6ff', color: '#2563eb' },
  success: { background: '#dcfce7', color: '#16a34a' },
  danger: { background: '#fee2e2', color: '#ef4444' },
  muted: { background: '#eef2f6', color: '#667085' }
}
const statusBadgeStyle = (status) => statusStyleMap[statusTone(status)] || statusStyleMap.muted
const reasonText = (reasonType) => ({
  FAKE_PRODUCT: '虚假或违规商品',
  INAPPROPRIATE_CONTENT: '不当内容',
  SCAM: '欺诈风险',
  TRADE_DISPUTE: '交易纠纷',
  FAKE: '虚假信息',
  FRAUD: '疑似诈骗',
  PROHIBITED: '违规商品',
  INAPPROPRIATE: '不当内容',
  HARASSMENT: '骚扰行为',
  OTHER: '其他原因'
}[reasonType] || reasonType || '举报')
const itemTitle = (record) => mode.value === 'REPORT' ? reasonText(record.reasonType) : (record.reason || '申诉')
const itemDescription = (record) => record.description || record.reason || '暂无补充说明'
const actorName = (record) => mode.value === 'REPORT'
  ? displayRelatedUserName(record, 'reporter', `用户${record.reporterId}`)
  : displayRelatedUserName(record, 'appellant', `用户${record.appellantId}`)
const shortTime = (value) => value ? String(value).replace('T', ' ').slice(0, 16) : ''
const canHandle = (status) => ['PENDING', 'PROCESSING'].includes(status)
const productCover = (record) => normalizeImage(
  record?.imageUrl
  || record?.image_url
  || record?.coverImage
  || record?.cover_image
  || record?.image
  || record?.productImage
  || record?.product_image
  || record?.images?.[0]
  || ''
)
const productStatusText = (status) => ({ ON_SALE: '在售', OFF_SHELF: '已下架', LOCKED: '交易锁定', SOLD: '已售出', DELETED: '已删除' }[status] || status || '未知')
const findAdminOrder = async (orderId) => {
  const result = await getAdminOrders('ALL', { pageSize: 200 })
  const list = Array.isArray(result) ? result : (result?.items || result?.list || result?.records || [])
  return list.find((record) => Number(record.id) === Number(orderId)) || null
}
const orderProductImage = (record) => normalizeImage(
  record?.productImage
  || record?.productImageUrl
  || record?.productImageSnapshot
  || record?.productCover
  || record?.productCoverImage
  || record?.coverImage
  || record?.imageUrl
  || record?.product?.image
  || record?.product?.imageUrl
  || record?.product?.coverImage
  || record?.product?.images?.[0]
  || ''
)

const previewImage = (current, urls) => uni.previewImage({ current, urls: urls.map((url) => normalizeImage(url)) })
const relatedQuery = () => `relatedType=${mode.value}&relatedId=${id.value}`
const productDetailUrl = (productId, record = null) => {
  if (!productId) return ''
  const params = [`id=${productId}`, 'readonly=1', 'adminView=1', relatedQuery()]
  if (record?.title) params.push(`snapshotTitle=${encodeURIComponent(record.title)}`)
  if (record?.price !== undefined && record?.price !== null) params.push(`snapshotPrice=${record.price}`)
  if (record?.status) params.push(`snapshotStatus=${record.status}`)
  const image = productCover(record)
  if (image) params.push(`snapshotImage=${encodeURIComponent(image)}`)
  const sellerId = record?.sellerId || record?.seller?.id
  const sellerName = record?.sellerName || record?.seller?.nickname
  if (sellerId) params.push(`snapshotSellerId=${sellerId}`)
  if (sellerName) params.push(`snapshotSellerName=${encodeURIComponent(sellerName)}`)
  return `/pages/detail/detail?${params.join('&')}`
}
const goProduct = (productId) => {
  const url = productDetailUrl(productId, targetProduct.value)
  if (url) uni.navigateTo({ url })
}
const orderDetailUrl = (orderId) => {
  if (!orderId) return ''
  return `/pages/order/detail?id=${orderId}&adminView=1&readonly=1&${relatedQuery()}`
}
const goTargetManage = () => {
  if (item.value?.targetType === 'USER' && item.value.targetId) {
    uni.navigateTo({ url: `/pages/user-profile/user-profile?id=${item.value.targetId}&adminView=1&${relatedQuery()}` })
    return
  }
  if (item.value?.targetType === 'PRODUCT' && item.value.targetId) {
    const url = productDetailUrl(item.value.targetId)
    if (url) uni.navigateTo({ url })
    return
  }
  if (item.value?.targetType === 'ORDER' && item.value.targetId) {
    const url = orderDetailUrl(item.value.targetId)
    if (url) uni.navigateTo({ url })
    return
  }
  if (!targetManageUrl.value) return
  uni.navigateTo({ url: `${targetManageUrl.value}?${relatedQuery()}` })
}
const openReasonModal = () => { reasonModal.visible = true; reasonModal.reason = ''; reasonModal.note = '' }
const closeReasonModal = () => { if (!submitting.value) reasonModal.visible = false }
const buildHandleResult = (reason, note) => {
  const extra = String(note || '').trim()
  return extra ? `${reason}。补充说明：${extra}` : reason
}

const handleCurrent = (status) => {
  const action = statusText(status)
  uni.showModal({
    title: '确认处理',
    content: `确定将该${mode.value === 'REPORT' ? '举报' : '申诉'}标记为${action}吗？`,
    success: async (res) => {
      if (!res.confirm || submitting.value) return
      submitting.value = true
      try {
        if (mode.value === 'REPORT') {
          if (status === 'PROCESSING') await markAdminReportProcessing(id.value)
          else await handleAdminReport(id.value, status, `举报${action}`)
        } else {
          if (status === 'PROCESSING') await markAdminAppealProcessing(id.value)
          else await handleAdminAppeal(id.value, status, `申诉${action}`)
        }
        uni.showToast({ title: '处理成功', icon: 'success' })
        await load()
      } catch (error) {
        uni.showToast({ title: error.message || '处理失败', icon: 'none' })
      } finally {
        submitting.value = false
      }
    }
  })
}

const submitReject = async () => {
  if (!reasonModal.reason) {
    uni.showToast({ title: '请选择处理原因', icon: 'none' })
    return
  }
  if (submitting.value) return
  submitting.value = true
  try {
    const handleResult = buildHandleResult(reasonModal.reason, reasonModal.note)
    if (mode.value === 'REPORT') await handleAdminReport(id.value, 'REJECTED', handleResult)
    else await handleAdminAppeal(id.value, 'REJECTED', handleResult)
    uni.showToast({ title: '处理成功', icon: 'success' })
    reasonModal.visible = false
    await load()
  } catch (error) {
    uni.showToast({ title: error.message || '处理失败', icon: 'none' })
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.page { min-height: 100vh; padding: 24rpx; background: #f5f6f8; box-sizing: border-box; }
.detail-card, .empty { padding: 28rpx; border-radius: 16rpx; background: #fff; margin-bottom: 18rpx; }
.card-head { display: flex; align-items: flex-start; justify-content: space-between; gap: 20rpx; }
.title { font-size: 34rpx; font-weight: 700; color: #1f2933; }
.sub { margin-top: 8rpx; font-size: 24rpx; color: #8a96a8; }
.badge-row { display: flex; flex-shrink: 0; align-items: center; gap: 8rpx; flex-wrap: wrap; justify-content: flex-end; }
.status-badge { flex-shrink: 0; padding: 8rpx 14rpx; border-radius: 999rpx; background: #eef2f6; color: #667085; font-size: 22rpx; }
.risk-badge { flex-shrink: 0; padding: 8rpx 14rpx; border-radius: 999rpx; background: #fff1f2; color: #be123c; font-size: 22rpx; }
.desc { margin-top: 18rpx; color: #667085; font-size: 26rpx; line-height: 38rpx; }
.meta-row { display: flex; justify-content: space-between; gap: 20rpx; margin-top: 20rpx; font-size: 22rpx; color: #98a2b3; }
.target-product { display: flex; align-items: center; gap: 18rpx; padding: 18rpx; margin-top: 24rpx; border-radius: 14rpx; background: #f8fafc; }
.target-image { width: 116rpx; height: 116rpx; border-radius: 12rpx; background: #eef2f6; flex-shrink: 0; }
.placeholder { display: flex; align-items: center; justify-content: center; color: #98a2b3; font-size: 22rpx; }
.target-main { flex: 1; min-width: 0; }
.target-title { font-size: 28rpx; font-weight: 700; color: #1f2933; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.target-meta { margin-top: 8rpx; font-size: 24rpx; color: #667085; }
.target-arrow { color: #b2bdca; font-size: 40rpx; }
.section-title { margin: 26rpx 0 12rpx; font-size: 28rpx; font-weight: 700; color: #1f2933; }
.detail-row { display: flex; justify-content: space-between; gap: 24rpx; padding: 14rpx 0; font-size: 24rpx; color: #667085; }
.detail-row text:last-child { flex: 1; color: #1f2933; text-align: right; word-break: break-all; }
.image-grid { display: grid; grid-template-columns: repeat(3, minmax(0, 1fr)); gap: 12rpx; margin-top: 16rpx; }
.case-image { width: 100%; height: 150rpx; border-radius: 12rpx; background: #eef2f6; }
.handle-result { margin-top: 18rpx; padding: 16rpx; border-radius: 12rpx; background: #f8fafc; color: #667085; font-size: 24rpx; }
.actions { margin-top: 22rpx; display: flex; gap: 14rpx; flex-wrap: wrap; }
.actions button { margin: 0; }
.pass { background: #17a84b; color: #fff; }
.process { background: #eff6ff; color: #2563eb; }
.reject { background: #fff1f2; color: #ef4444; }
.modal-mask { position: fixed; left: 0; right: 0; top: 0; bottom: 0; z-index: 99; display: flex; align-items: flex-end; background: rgba(15, 23, 42, 0.42); }
.reason-sheet { width: 100%; padding: 30rpx 28rpx 36rpx; border-radius: 28rpx 28rpx 0 0; background: #fff; box-sizing: border-box; }
.sheet-title { font-size: 34rpx; font-weight: 700; color: #1f2933; }
.reason-list { display: flex; flex-wrap: wrap; gap: 14rpx; margin-top: 26rpx; }
.reason-chip { padding: 14rpx 18rpx; border-radius: 999rpx; background: #f8fafc; color: #475467; font-size: 24rpx; border: 2rpx solid transparent; }
.reason-chip.active { border-color: #17a84b; background: #f0fdf4; color: #16a34a; font-weight: 700; }
.reason-input { width: 100%; min-height: 150rpx; margin-top: 24rpx; padding: 20rpx; border-radius: 16rpx; background: #f8fafc; font-size: 26rpx; color: #1f2933; box-sizing: border-box; }
.reason-placeholder { color: #98a2b3; }
.sheet-actions { display: grid; grid-template-columns: 1fr 1fr; gap: 18rpx; margin-top: 24rpx; }
.sheet-button { height: 84rpx; line-height: 84rpx; border-radius: 14rpx; font-size: 28rpx; }
.sheet-button.cancel { background: #f8fafc; color: #667085; }
.sheet-button.confirm { background: #17a84b; color: #fff; }
</style>
