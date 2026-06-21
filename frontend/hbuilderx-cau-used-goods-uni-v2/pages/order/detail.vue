<template>
  <view v-if="order" class="page">
    <view :class="['card', 'status-card', 'result-status-card', statusCardClass]">
      <view class="status-copy">
        <text class="status-title">{{ status.label }}</text>
        <text class="status-tip">{{ statusTip }}</text>
      </view>
    </view>

    <view class="card product-card" @click="openProduct">
      <ProductRow v-if="order.product" :product="order.product" />
    </view>

    <view class="card seller-card" @click="openSeller">
      <image v-if="sellerAvatar" class="seller-avatar image-avatar" :src="sellerAvatar" mode="aspectFill" />
      <view v-else class="seller-avatar">{{ sellerName.slice(0, 1) }}</view>
      <view class="seller-body">
        <text class="seller-label">卖家信息</text>
        <text class="seller-name">{{ sellerName }}</text>
      </view>
      <text class="seller-arrow">›</text>
    </view>

    <view class="card info">
      <view class="info-row"><text class="info-label">订单编号</text><text class="info-value">{{ order.id }}</text></view>
      <view class="info-row"><text class="info-label">预约时间</text><text class="info-value">{{ order.meetTime }}</text></view>
      <view class="info-row"><text class="info-label">面交地点</text><text class="info-value">{{ order.meetLocation }}</text></view>
      <view class="info-row"><text class="info-label">备注</text><text class="info-value">{{ order.remark || '无' }}</text></view>
      <view v-if="order.expireTime && order.status === 'PENDING_CONFIRM'" class="info-row">
        <text class="info-label">确认时限</text><text class="info-value">{{ order.expireTime }}</text>
      </view>
    </view>

    <view v-if="order.status === 'WAIT_MEET'" class="notice">
      为保护隐私，联系方式仅在待面交阶段向交易双方展示：{{ order.contact || '后端暂未返回联系方式字段' }}
    </view>

    <view v-if="order.status === 'WAIT_MEET' || order.status === 'COMPLETED'" class="card confirm-card">
      <view class="confirm-title">交易完成确认</view>
      <view class="confirm-row">
        <text>卖家确认</text>
        <text :class="['confirm-state', sellerConfirmed ? 'done' : 'pending']">{{ sellerConfirmed ? '已确认' : '待确认' }}</text>
      </view>
    </view>

    <view v-if="false && adminView" class="notice">
      管理员只读查看订单，如需异常关闭请在风险处理入口继续操作。
    </view>

    <view v-if="adminView" class="notice">
      管理员只读查看订单，可对待确认或待面交订单执行异常关闭。
    </view>

    <view v-if="adminView" class="actions admin-actions">
      <button v-if="canAdminExceptionClose" class="btn btn-danger" :loading="submitting" @click="openCloseModal">异常关闭</button>
      <button v-else class="btn btn-disabled" disabled>{{ status.label || '不可操作' }}</button>
    </view>

    <view v-if="!adminView && !readonlyMode" class="actions">
      <button v-if="isSeller && order.status === 'PENDING_CONFIRM'" class="btn btn-primary" @click="change('confirm')">确认预约</button>
      <button v-if="isSeller && order.status === 'WAIT_MEET'" class="btn btn-primary" @click="change('complete')">确认完成交易</button>
      <button v-if="canCancel" class="btn btn-plain" @click="cancel">取消订单</button>
      <button v-if="order.status === 'COMPLETED' && !isSeller" class="btn btn-primary" :disabled="hasReviewed" @click="review">
        {{ hasReviewed ? '已评价' : '去评价' }}
      </button>
      <button class="btn btn-plain" @click="appeal">申诉订单问题</button>
      <button class="btn btn-plain" @click="report">举报交易问题</button>
    </view>
    <view v-if="fromMessage" class="actions message-center-actions">
      <button class="btn btn-plain" @click="goMessageCenter">返回消息中心</button>
    </view>
    <view v-if="closeModal.visible" class="modal-mask" @click="closeCloseModal">
      <view class="reason-sheet" @click.stop>
        <view class="sheet-title">异常关闭订单</view>
        <view class="sheet-sub">订单 #{{ order.id }}</view>
        <view class="party-row">
          <view class="party-label">责任方</view>
          <view class="party-options">
            <view :class="['party-chip', closeModal.responsibleParty === 'BUYER' ? 'active' : '']" @click="closeModal.responsibleParty = 'BUYER'">买家</view>
            <view :class="['party-chip', closeModal.responsibleParty === 'SELLER' ? 'active' : '']" @click="closeModal.responsibleParty = 'SELLER'">卖家</view>
          </view>
        </view>
        <view class="reason-list">
          <view
            v-for="reason in closeReasons"
            :key="reason"
            :class="['reason-chip', closeModal.reason === reason ? 'active' : '']"
            @click="closeModal.reason = reason"
          >
            {{ reason }}
          </view>
        </view>
        <textarea v-model="closeModal.note" class="reason-input" maxlength="200" placeholder="补充说明，可不填" placeholder-class="reason-placeholder" />
        <view class="sheet-actions">
          <button class="sheet-button cancel" @click="closeCloseModal">取消</button>
          <button class="sheet-button confirm" :loading="submitting" @click="submitExceptionClose">确认关闭</button>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { onLoad } from '@dcloudio/uni-app'
import { computed, reactive, ref } from 'vue'
import ProductRow from '../../components/ProductRow.vue'
import { exceptionCloseAdminOrder } from '../../api/admin'
import { getPublicProfile } from '../../api/user'
import { tradeService } from '../../services/trade'
import { getUser } from '../../utils/auth'
import { ORDER_STATUS } from '../../utils/constants'
import { BASE_URL } from '../../utils/request'
import { navigate, showError, showSuccess } from '../../utils/navigation'

const order = ref()
const sellerProfile = ref(null)
const adminView = ref(false)
const readonlyMode = ref(false)
const relatedType = ref('')
const relatedId = ref('')
const fromMessage = ref(false)
const submitting = ref(false)
const closeModal = reactive({ visible: false, reason: '', note: '', responsibleParty: 'SELLER' })
const closeReasons = ['买卖双方协商取消', '交易存在纠纷', '商品违规或信息异常', '长时间未完成交易', '其他原因']
const currentUserId = computed(() => String(getUser()?.id || ''))
const status = computed(() => ORDER_STATUS[order.value?.status] || { label: '', tone: 'muted' })
const statusCardClass = computed(() => ({
  PENDING_CONFIRM: 'pending-status-card',
  WAIT_MEET: 'active-status-card',
  COMPLETED: 'success-status-card',
  CANCELED: 'danger-status-card',
  CANCELLED: 'danger-status-card',
  EXCEPTION_CLOSED: 'danger-status-card'
}[order.value?.status] || 'neutral-status-card'))
const isSeller = computed(() => String(order.value?.sellerId) === currentUserId.value)
const canCancel = computed(() => ['PENDING_CONFIRM', 'WAIT_MEET'].includes(order.value?.status))
const canAdminExceptionClose = computed(() => ['PENDING_CONFIRM', 'WAIT_MEET'].includes(order.value?.status))
const sellerId = computed(() => order.value?.sellerId || order.value?.seller?.id || '')
const sellerName = computed(() => (
  sellerProfile.value?.nickname
  || order.value?.sellerName
  || order.value?.sellerNickname
  || order.value?.seller?.nickname
  || 'CAU 卖家'
))
const sellerAvatar = computed(() => normalizeImage(
  sellerProfile.value?.avatarUrl
  || order.value?.sellerAvatarUrl
  || order.value?.seller?.avatarUrl
  || order.value?.sellerAvatar
))
const sellerConfirmed = computed(() => order.value?.status === 'COMPLETED')
const hasReviewed = computed(() => order.value?.reviewed === true || uni.getStorageSync(`order-reviewed-${id}`) === true)
const statusTip = computed(() => ({
  PENDING_CONFIRM: '卖家需要在 24 小时内处理预约',
  WAIT_MEET: '请按约定时间在校园公共区域完成面交',
  COMPLETED: '线下交易已完成',
  CANCELED: '订单已关闭，商品将按规则恢复在售',
  CANCELLED: '订单已关闭，商品将按规则恢复在售',
  EXCEPTION_CLOSED: '订单由管理员介入关闭'
}[order.value?.status] || ''))

let id = ''
onLoad((options) => {
  id = options.id
  adminView.value = options.adminView === '1' || options.adminView === 1
  readonlyMode.value = options.readonly === '1' || options.readonly === 1
  relatedType.value = String(options.relatedType || '').toUpperCase()
  relatedId.value = options.relatedId || ''
  fromMessage.value = options.fromMessage === '1' || options.fromMessage === 1
  load()
})

async function load() {
  try {
    const detail = adminView.value ? await tradeService.getAdminOrder(id) : await tradeService.getOrder(id)
    if (detail?.product?.id && !detail.product.image) {
      try {
        const product = await tradeService.getProduct(detail.product.id)
        detail.product = { ...detail.product, ...product, image: product.image || detail.product.image }
      } catch (error) {
        // 商品锁定、售出或下架时，订单详情仍按订单快照展示。
      }
    }
    order.value = detail
    loadSellerProfile()
  } catch (error) {
    showError(error)
  }
}

async function loadSellerProfile() {
  if (!sellerId.value) return
  sellerProfile.value = await getPublicProfile(sellerId.value).catch(() => null)
}

function normalizeImage(url) {
  if (!url) return ''
  return /^https?:\/\//.test(url) ? url : `${BASE_URL}${url}`
}

async function change(action, payload) {
  try {
    order.value = await tradeService.changeOrderStatus(id, action, payload)
    showSuccess('操作成功')
  } catch (error) {
    showError(error)
  }
}

function cancel() {
  uni.showModal({
    title: '取消订单',
    content: '确认取消当前订单吗？取消后商品将按规则恢复在售。',
    success: ({ confirm }) => confirm && change('cancel', { reason: '用户主动取消' })
  })
}

function openCloseModal() {
  closeModal.visible = true
  closeModal.reason = ''
  closeModal.note = ''
  closeModal.responsibleParty = 'SELLER'
}

function closeCloseModal() {
  if (submitting.value) return
  closeModal.visible = false
}

function buildCloseReason() {
  const note = String(closeModal.note || '').trim()
  return note ? `${closeModal.reason}。补充说明：${note}` : closeModal.reason
}

async function submitExceptionClose() {
  if (!canAdminExceptionClose.value) return
  if (!closeModal.reason) {
    showError(new Error('请选择关闭原因'))
    return
  }
  if (submitting.value) return
  submitting.value = true
  try {
    const payload = {
      reason: buildCloseReason(),
      responsibleParty: closeModal.responsibleParty
    }
    if (relatedType.value && relatedId.value) {
      payload.relatedType = relatedType.value
      payload.relatedId = Number(relatedId.value)
    }
    await exceptionCloseAdminOrder(id, payload)
    showSuccess('已异常关闭')
    closeModal.visible = false
    await load()
  } catch (error) {
    showError(error)
  } finally {
    submitting.value = false
  }
}

function review() {
  if (hasReviewed.value) return
  navigate('/pages/interaction/review', { orderId: id })
}

function openProduct() {
  const productId = order.value?.product?.id || order.value?.productId
  if (!productId) {
    showError(new Error('商品信息缺失，暂时无法查看'))
    return
  }
  navigate('/pages/detail/detail', {
    id: productId,
    readonly: adminView.value || readonlyMode.value || order.value?.status === 'COMPLETED' ? 1 : 0,
    adminView: adminView.value ? 1 : '',
    relatedType: relatedType.value,
    relatedId: relatedId.value,
    snapshotTitle: order.value?.product?.title || order.value?.productTitleSnapshot || '',
    snapshotPrice: order.value?.product?.price || order.value?.productPriceSnapshot || '',
    snapshotImage: order.value?.product?.image || order.value?.productImage || '',
    snapshotMeetLocation: order.value?.meetLocation || order.value?.product?.meetLocation || '',
    snapshotSellerId: sellerId.value,
    snapshotSellerName: sellerName.value
  })
}

function openSeller() {
  if (!sellerId.value) return
  const productId = order.value?.product?.id || order.value?.productId
  navigate('/pages/user-profile/user-profile', {
    id: sellerId.value,
    adminView: adminView.value ? 1 : '',
    productId,
    productTitle: order.value?.product?.title || order.value?.productTitleSnapshot || ''
  })
}

function report() {
  navigate('/pages/interaction/report', { targetType: 'ORDER', targetId: id })
}

function appeal() {
  navigate('/pages/interaction/appeal', { targetType: 'ORDER', targetId: id })
}

function goMessageCenter() {
  uni.switchTab({ url: '/pages/messages/messages' })
}
</script>

<style scoped lang="scss">
.page {
  box-sizing: border-box;
  width: 100%;
  max-width: 100vw;
  overflow-x: hidden;
  padding: 20rpx 24rpx 48rpx;
}

.card {
  box-sizing: border-box;
  width: 100%;
  max-width: 100%;
  margin-bottom: 20rpx;
  padding: 24rpx;
  border-radius: 20rpx;
  background: #fff;
}

.status-card {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 18rpx;
}

.result-status-card {
  align-items: center;
  justify-content: center;
  color: #fff;
  text-align: center;
}

.result-status-card .status-copy {
  align-items: center;
}

.result-status-card .status-tip {
  color: rgba(255, 255, 255, .82);
}

.success-status-card {
  background: linear-gradient(135deg, #23734f, #3e9b72);
  box-shadow: 0 12rpx 32rpx rgba(35, 115, 79, .18);
}

.pending-status-card {
  background: linear-gradient(135deg, #d99424, #edb64b);
  box-shadow: 0 12rpx 32rpx rgba(190, 126, 24, .18);
}

.active-status-card {
  background: linear-gradient(135deg, #3478b8, #55a0d8);
  box-shadow: 0 12rpx 32rpx rgba(38, 105, 165, .18);
}

.danger-status-card {
  background: linear-gradient(135deg, #c9443e, #e0645d);
  box-shadow: 0 12rpx 32rpx rgba(180, 48, 43, .18);
}

.neutral-status-card {
  background: linear-gradient(135deg, #65746c, #87958e);
  box-shadow: 0 12rpx 32rpx rgba(68, 83, 75, .16);
}

.status-copy {
  display: flex;
  flex: 1;
  min-width: 0;
  flex-direction: column;
  gap: 10rpx;
}

.status-title {
  font-size: 36rpx;
  font-weight: 700;
}

.status-tip {
  color: #738077;
  font-size: 24rpx;
  line-height: 1.5;
}

.info-row {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 28rpx;
  padding: 16rpx 0;
  border-bottom: 1rpx solid #eef1ef;
  font-size: 25rpx;
  line-height: 1.5;
}

.info-row:last-child {
  border-bottom: 0;
}

.info-label {
  flex: 0 0 150rpx;
  color: #738077;
}

.info-value {
  flex: 1;
  min-width: 0;
  color: #36443c;
  text-align: right;
  overflow-wrap: break-word;
  word-break: break-all;
}

.actions {
  box-sizing: border-box;
  width: 100%;
  max-width: 100%;
}

.actions .btn {
  box-sizing: border-box;
  min-width: 0;
}

.admin-actions {
  margin-bottom: 20rpx;
}

.message-center-actions {
  margin-top: 20rpx;
}

.btn-danger {
  color: #ef4444;
  background: #fee2e2;
}

.btn-disabled {
  color: #8a9690;
  background: #eef2f6;
}

.modal-mask {
  position: fixed;
  right: 0;
  bottom: 0;
  left: 0;
  top: 0;
  z-index: 99;
  display: flex;
  align-items: flex-end;
  background: rgba(15, 23, 42, .42);
}

.reason-sheet {
  box-sizing: border-box;
  width: 100%;
  padding: 30rpx 28rpx 36rpx;
  border-radius: 28rpx 28rpx 0 0;
  background: #fff;
}

.sheet-title {
  color: #1f2933;
  font-size: 34rpx;
  font-weight: 700;
}

.sheet-sub {
  margin-top: 8rpx;
  color: #98a2b3;
  font-size: 24rpx;
}

.party-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 20rpx;
  margin-top: 24rpx;
}

.party-label {
  color: #475467;
  font-size: 26rpx;
  font-weight: 700;
}

.party-options,
.reason-list,
.sheet-actions {
  display: flex;
  gap: 14rpx;
}

.party-chip,
.reason-chip {
  padding: 12rpx 24rpx;
  border: 2rpx solid transparent;
  border-radius: 999rpx;
  background: #f8fafc;
  color: #667085;
  font-size: 24rpx;
}

.party-chip.active,
.reason-chip.active {
  border-color: #17a84b;
  background: #f0fdf4;
  color: #16a34a;
  font-weight: 700;
}

.reason-list {
  flex-wrap: wrap;
  margin-top: 24rpx;
}

.reason-input {
  box-sizing: border-box;
  width: 100%;
  min-height: 150rpx;
  margin-top: 22rpx;
  padding: 18rpx;
  border-radius: 14rpx;
  background: #f8fafc;
  color: #1f2933;
  font-size: 25rpx;
  line-height: 1.5;
}

.sheet-actions {
  margin-top: 24rpx;
}

.sheet-button {
  flex: 1;
  height: 76rpx;
  border-radius: 14rpx;
  font-size: 26rpx;
  line-height: 76rpx;
}

.sheet-button.cancel {
  color: #667085;
  background: #f2f4f7;
}

.sheet-button.confirm {
  color: #fff;
  background: #17a84b;
}

.seller-card {
  display: flex;
  align-items: center;
  gap: 18rpx;
}

.seller-avatar {
  display: flex;
  width: 78rpx;
  height: 78rpx;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  background: #edf6f1;
  color: #23734f;
  font-weight: 700;
}

.seller-body {
  display: flex;
  flex: 1;
  min-width: 0;
  flex-direction: column;
  gap: 8rpx;
}

.seller-label {
  color: #8a9690;
  font-size: 23rpx;
}

.seller-name {
  color: #243129;
  font-size: 29rpx;
  font-weight: 700;
}

.seller-arrow {
  color: #b8c1bd;
  font-size: 44rpx;
}

.product-card {
  cursor: pointer;
}

.image-avatar {
  background: #e8ecef;
}

.confirm-title {
  margin-bottom: 12rpx;
  color: #243129;
  font-size: 29rpx;
  font-weight: 700;
}

.confirm-row {
  display: flex;
  justify-content: space-between;
  padding: 12rpx 0;
  color: #66736b;
  font-size: 25rpx;
}

.confirm-state.done { color: #23734f; }
.confirm-state.pending { color: #b27b1f; }
</style>
