<template>
  <view class="page">
    <view v-if="profile" class="profile-card">
      <view class="profile-top">
        <image v-if="avatarUrl" class="avatar image-avatar" :src="avatarUrl" mode="aspectFill" />
        <view v-else class="avatar text-avatar">{{ displayName.slice(0, 1) }}</view>
        <view class="info">
          <text class="name">{{ displayName }}</text>
          <view class="tags">
            <text class="tag">{{ profile.authStatus === 'VERIFIED' ? '学生已认证' : '学生未认证' }}</text>
            <text class="tag" :class="{ danger: !profile.tradeAvailable }">{{ profile.tradeAvailable ? '交易正常' : '交易受限' }}</text>
          </view>
        </view>
      </view>
      <view class="actions">
        <template v-if="adminView">
          <button v-if="accountStatus === 'NORMAL'" class="disable-btn" @click="changeUserStatus('DISABLED')">禁用</button>
          <button v-if="accountStatus === 'NORMAL' || accountStatus === 'DISABLED'" class="ban-btn" @click="changeUserStatus('BANNED')">封禁</button>
          <button v-if="canRecoverUser" class="recover-btn" @click="changeUserStatus('NORMAL')">{{ recoverButtonText }}</button>
        </template>
        <template v-else>
          <button v-if="!isSelf" class="chat-btn" :disabled="isCurrentUserRestricted" @click="chatWithUser">聊一聊</button>
          <button class="report-btn" :disabled="isCurrentUserRestricted" @click="reportUser">举报该用户</button>
        </template>
      </view>
    </view>

    <view class="section">
      <view class="section-head">
        <text class="section-title">TA 在卖</text>
        <text class="section-subtitle">{{ products.length }} 件商品</text>
      </view>
      <view v-if="products.length" class="product-grid">
        <view v-for="item in products" :key="item.id" class="product-card" @click="openProduct(item.id)">
          <image v-if="item.cover" class="product-image" :src="item.cover" mode="aspectFill" />
          <view v-else class="product-image placeholder">商品</view>
          <text class="product-title">{{ item.title }}</text>
          <text class="price">¥{{ formatPrice(item.price) }}</text>
        </view>
      </view>
      <view v-else class="empty">暂无公开在售商品</view>
    </view>

    <view class="section">
      <view class="section-head">
        <text class="section-title">收到的评价</text>
        <text class="section-subtitle">{{ reviews.length }} 条</text>
      </view>
      <view v-if="reviews.length" class="review-list">
        <view v-for="item in reviews" :key="item.id" class="review-card">
          <view class="review-head">
            <text class="stars">{{ starText(item.rating) }}</text>
            <text class="review-time">{{ item.createTime }}</text>
          </view>
          <text class="review-content">{{ item.content }}</text>
          <text class="review-product">{{ item.productTitle }}</text>
        </view>
      </view>
      <view v-else class="empty">暂无公开评价</view>
    </view>
  </view>
</template>

<script setup>
import { computed, ref } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { getPublicProfile } from '../../api/user'
import { getAdminUserDetail, updateAdminUserStatus } from '../../api/admin'
import { listProducts } from '../../api/product'
import { createOrGetConversation } from '../../api/chat'
import { getUser } from '../../utils/auth'
import { BASE_URL } from '../../utils/request'
import { navigate, showError } from '../../utils/navigation'
import { accountStatusOf, isBannedUserStatus, isDisabledUserStatus } from '../../utils/user-format'

const userId = ref('')
const adminView = ref(false)
const relatedType = ref('')
const relatedId = ref('')
const preferredProductId = ref('')
const preferredProductTitle = ref('')
const profile = ref(null)
const adminUser = ref(null)
const products = ref([])
const reviews = ref([])

const currentUserId = computed(() => getUser()?.id || getUser()?.userId || '')
const isSelf = computed(() => String(userId.value) === String(currentUserId.value))
const currentUserAccountStatus = computed(() => accountStatusOf(getUser() || {}))
const isCurrentUserRestricted = computed(() => isBannedUserStatus(currentUserAccountStatus.value) || isDisabledUserStatus(currentUserAccountStatus.value))
const currentUserRestrictionText = computed(() => {
  if (isBannedUserStatus(currentUserAccountStatus.value)) return '你已被封禁，无法进行操作'
  if (isDisabledUserStatus(currentUserAccountStatus.value)) return '你已被禁用，无法进行操作'
  return ''
})
const displayName = computed(() => profile.value?.nickname || 'CAU 同学')
const avatarUrl = computed(() => normalizeImage(profile.value?.avatarUrl))
const accountStatus = computed(() => adminUser.value?.accountStatus || (profile.value?.tradeAvailable ? 'NORMAL' : 'DISABLED'))
const currentUserRole = computed(() => String(getUser()?.role || '').toUpperCase())
const isSuperAdmin = computed(() => currentUserRole.value === 'SUPER_ADMIN')
const hasAppealContext = computed(() => relatedType.value === 'APPEAL' && Number(relatedId.value) > 0)
const canRecoverUser = computed(() => accountStatus.value === 'DISABLED' || (accountStatus.value === 'BANNED' && isSuperAdmin.value && hasAppealContext.value))
const recoverButtonText = computed(() => accountStatus.value === 'BANNED' ? '解封' : '恢复')

function normalizeImage(url) {
  if (!url) return ''
  return /^https?:\/\//.test(url) ? url : `${BASE_URL}${url}`
}

function pickCover(item) {
  const raw = item.coverImage || item.mainImage || item.imageUrl || item.image
    || item.images?.[0]?.url || item.images?.[0]?.imageUrl || item.images?.[0]
  return normalizeImage(raw)
}

function formatPrice(value) {
  const number = Number(value || 0)
  return Number.isInteger(number) ? String(number) : number.toFixed(2)
}

function starText(value) {
  const rating = Math.max(0, Math.min(5, Number(value || 0)))
  return '★★★★★'.slice(0, rating) + '☆☆☆☆☆'.slice(0, 5 - rating)
}

async function loadProducts() {
  const result = await listProducts({ page: 1, pageSize: 80 }).catch(() => ({ list: [], items: [] }))
  const list = result.list || result.items || []
  products.value = list
    .filter((item) => String(item.sellerId || item.userId || item.ownerId) === String(userId.value))
    .map((item) => ({ ...item, cover: pickCover(item) }))
}

function loadLocalReviews() {
  const list = uni.getStorageSync(`user-reviews-${userId.value}`) || []
  reviews.value = Array.isArray(list) ? list : []
}

function profileFromAdminUser(user) {
  if (!user) return null
  return {
    id: user.id,
    nickname: user.nickname,
    avatarUrl: user.avatarUrl,
    authStatus: user.authStatus,
    tradeAvailable: user.accountStatus === 'NORMAL'
  }
}

onLoad(async (options) => {
  userId.value = options.id || ''
  adminView.value = options.adminView === '1' || options.adminView === 1
  relatedType.value = String(options.relatedType || '').toUpperCase()
  relatedId.value = options.relatedId || ''
  preferredProductId.value = options.productId || ''
  preferredProductTitle.value = options.productTitle ? decodeURIComponent(options.productTitle) : ''
  if (!userId.value) {
    showError(new Error('用户不存在'))
    return
  }
  try {
    if (adminView.value) {
      const detail = await getAdminUserDetail(userId.value).catch(() => null)
      adminUser.value = detail?.user || null
      profile.value = await getPublicProfile(userId.value).catch(() => profileFromAdminUser(adminUser.value))
    } else {
      profile.value = await getPublicProfile(userId.value)
    }
    await loadProducts()
    loadLocalReviews()
  } catch (error) {
    showError(error)
  }
})

function openProduct(id) {
  if (!id) return
  navigate('/pages/detail/detail', { id })
}

function reportUser() {
  if (isSelf.value) {
    uni.showToast({ title: '不能举报自己', icon: 'none' })
    return
  }
  if (isCurrentUserRestricted.value) {
    uni.showToast({ title: currentUserRestrictionText.value, icon: 'none' })
    return
  }
  navigate('/pages/interaction/report', { targetType: 'USER', targetId: userId.value })
}

function accountText(status) {
  return { NORMAL: '恢复', DISABLED: '禁用', BANNED: '封禁' }[status] || status
}

function changeUserStatus(status) {
  const action = accountText(status)
  uni.showModal({
    title: `${action}用户`,
    editable: true,
    placeholderText: '请输入处理原因',
    success: async (res) => {
      if (!res.confirm) return
      const reason = (res.content || '').trim()
      if (!reason) {
        uni.showToast({ title: '请填写处理原因', icon: 'none' })
        return
      }
      try {
        const payload = { accountStatus: status, reason }
        if (relatedType.value && relatedId.value) {
          payload.relatedType = relatedType.value
          payload.relatedId = Number(relatedId.value)
        }
        await updateAdminUserStatus(userId.value, payload)
        const detail = await getAdminUserDetail(userId.value).catch(() => null)
        adminUser.value = detail?.user || adminUser.value
        if (status === 'NORMAL') {
          profile.value = { ...profile.value, tradeAvailable: true }
        } else {
          profile.value = { ...profile.value, tradeAvailable: false }
        }
        uni.showToast({ title: '操作成功', icon: 'success' })
      } catch (error) {
        uni.showToast({ title: error.message || '操作失败', icon: 'none' })
      }
    }
  })
}

async function chatWithUser() {
  if (isCurrentUserRestricted.value) {
    uni.showToast({ title: currentUserRestrictionText.value, icon: 'none' })
    return
  }
  const product = preferredProductId.value
    ? { id: preferredProductId.value, title: preferredProductTitle.value || products.value[0]?.title || '商品咨询' }
    : products.value[0]
  if (!product?.id) {
    uni.showToast({ title: 'TA 暂无可咨询商品', icon: 'none' })
    return
  }
  const productId = Number(product.id)
  if (!Number.isFinite(productId) || productId <= 0) {
    uni.showToast({ title: '商品信息异常，暂时无法聊天', icon: 'none' })
    return
  }
  try {
    const conversation = await createOrGetConversation(productId)
    navigate('/pages/chat/chat', {
      conversationId: conversation.id,
      title: product.title,
      targetUserId: userId.value,
      productId
    })
  } catch (error) {
    uni.showToast({ title: error.message || '暂时无法发起聊天', icon: 'none' })
  }
}
</script>

<style scoped>
.page { min-height: 100vh; padding: 24rpx; background: #f6f7f9; box-sizing: border-box; }
.profile-card, .section { border-radius: 28rpx; background: #fff; box-shadow: 0 8rpx 28rpx rgba(23, 33, 43, .04); }
.profile-card { padding: 30rpx; }
.profile-top { display: flex; gap: 22rpx; align-items: center; }
.avatar { display: flex; width: 124rpx; height: 124rpx; flex-shrink: 0; align-items: center; justify-content: center; border-radius: 50%; color: #fff; font-size: 42rpx; font-weight: 800; }
.text-avatar { background: linear-gradient(135deg, #f3b34c, #f47b45); }
.image-avatar { background: #e8ecef; }
.info { flex: 1; min-width: 0; }
.name { display: block; overflow: hidden; color: #222; font-size: 40rpx; font-weight: 800; text-overflow: ellipsis; white-space: nowrap; }
.tags { display: flex; flex-wrap: wrap; gap: 12rpx; margin-top: 18rpx; }
.tag { padding: 8rpx 16rpx; border-radius: 999rpx; background: #edf6f1; color: #23734f; font-size: 23rpx; }
.tag.danger { background: #fff1ef; color: #d85c45; }
.actions { display: flex; gap: 16rpx; margin-top: 28rpx; }
.chat-btn, .report-btn, .disable-btn, .ban-btn, .recover-btn { flex: 1; height: 72rpx; border-radius: 999rpx; font-size: 26rpx; line-height: 72rpx; }
.chat-btn { background: #23734f; color: #fff; }
.report-btn { background: #fff1ef; color: #d85c45; }
.disable-btn { background: #fff7e6; color: #a96500; }
.ban-btn { background: #fff1f2; color: #ef4444; }
.recover-btn { background: #e7f4ec; color: #23734f; }
.section { margin-top: 22rpx; padding: 26rpx; }
.section-head { display: flex; align-items: baseline; justify-content: space-between; margin-bottom: 20rpx; }
.section-title { color: #202124; font-size: 31rpx; font-weight: 800; }
.section-subtitle { color: #9aa2a8; font-size: 23rpx; }
.product-grid { display: grid; grid-template-columns: repeat(2, minmax(0, 1fr)); gap: 18rpx; }
.product-card { overflow: hidden; border-radius: 22rpx; background: #f8faf9; }
.product-image { display: flex; width: 100%; height: 210rpx; align-items: center; justify-content: center; background: #edf2ef; color: #8c9691; font-size: 26rpx; }
.product-title { display: block; overflow: hidden; margin: 16rpx 16rpx 8rpx; color: #26342f; font-size: 27rpx; font-weight: 700; text-overflow: ellipsis; white-space: nowrap; }
.price { display: block; margin: 0 16rpx 18rpx; color: #d8662f; font-size: 30rpx; font-weight: 800; }
.review-list { display: flex; flex-direction: column; gap: 16rpx; }
.review-card { padding: 20rpx; border-radius: 20rpx; background: #f8faf9; }
.review-head { display: flex; align-items: center; justify-content: space-between; gap: 16rpx; }
.stars { color: #f2a23a; font-size: 25rpx; letter-spacing: 2rpx; }
.review-time { color: #a0a6ad; font-size: 22rpx; }
.review-content { display: block; margin-top: 12rpx; color: #26342f; font-size: 26rpx; line-height: 1.55; }
.review-product { display: block; margin-top: 10rpx; color: #8a938e; font-size: 23rpx; }
.empty { padding: 34rpx 0; color: #9aa2a8; font-size: 25rpx; text-align: center; }
</style>
