<template>
  <view v-if="product" class="page">
    <swiper v-if="visibleImages.length" class="gallery" indicator-dots circular>
      <swiper-item v-for="image in visibleImages" :key="image">
        <image class="gallery-image" :src="image" mode="aspectFill" @error="markImageFailed(image)" />
      </swiper-item>
    </swiper>
    <view v-else class="gallery placeholder">图片未找到</view>

    <view class="card">
      <view class="price-line">
        <text class="price">¥{{ product.priceText }}</text>
        <text class="status">{{ statusText }}</text>
      </view>
      <view class="title">{{ product.title }}</view>
      <view class="meta">
        <text>{{ product.conditionText }}</text>
        <text>{{ product.viewCount || 0 }} 浏览</text>
        <text>{{ favoriteCount }} 收藏</text>
        <text>{{ product.timeText }}</text>
      </view>
    </view>

    <view class="card">
      <view class="section-title">商品描述</view>
      <view class="description">{{ product.description || '卖家暂未填写描述' }}</view>
    </view>

    <view class="card seller" @click="openSeller">
      <image v-if="sellerAvatarUrl" class="avatar image-avatar" :src="sellerAvatarUrl" mode="aspectFill" />
      <view v-else class="avatar">{{ sellerAvatarText }}</view>
      <view class="seller-body">
        <view class="seller-label">卖家信息</view>
        <view class="seller-name">{{ sellerName }}</view>
        <view class="meta single">{{ sellerCollege }}</view>
      </view>
      <text class="seller-arrow">›</text>
    </view>

    <view v-if="productReadonly" class="readonly-tip">该商品仅可查看</view>

    <view class="bottom">
      <button
        class="icon-button favorite"
        :class="{ active: isFavorite }"
        :disabled="productReadonly || isOwnProduct"
        @click="toggleFavorite"
      >
        {{ isFavorite ? '★' : '☆' }}
      </button>
      <button
        class="icon-button report"
        :class="{ disabled: isOwnProduct || productReadonly }"
        :disabled="isOwnProduct || productReadonly"
        @click="report"
      >
        !
      </button>
      <button class="chat" :disabled="productReadonly || isOwnProduct" @click="chat">聊一聊</button>
      <button class="primary" :disabled="productReadonly || isOwnProduct || product.status !== 'ON_SALE'" @click="reserve">
        {{ actionText }}
      </button>
    </view>
  </view>

  <view v-else class="page loading-page">
    <view class="load-text">正在加载商品详情...</view>
  </view>
</template>

<script setup>
import { computed, ref } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import {
  addFavorite,
  checkFavorite,
  getProductById,
  listCategories,
  removeFavorite
} from '../../api/product'
import { createOrGetConversation } from '../../api/chat'
import { getPublicProfile } from '../../api/user'
import { buildCategoryMap, formatPrice, formatProduct, getStatusText, normalizeImage } from '../../utils/product-format'
import { getToken, getUser, isVerifiedUser } from '../../utils/auth'
import { navigate } from '../../utils/navigation'
import { accountStatusOf, displayUserName, isBannedUserStatus, isCanceledUserStatus, isDisabledUserStatus } from '../../utils/user-format'
import { addBrowseHistory } from '../../utils/browse-history'

const product = ref(null)
const isFavorite = ref(false)
const failedImages = ref([])
const readonlyMode = ref(false)
const sellerProfile = ref(null)
const sellerAvatarFile = ref('')

const pick = (...values) => values.find((value) => value !== undefined && value !== null && value !== '') || ''
const sellerRestricted = computed(() => {
  const status = accountStatusOf(sellerSource.value)
  return isBannedUserStatus(status) || isDisabledUserStatus(status) || isCanceledUserStatus(status)
})
const productReadonly = computed(() => readonlyMode.value || sellerRestricted.value)
const statusText = computed(() => productReadonly.value ? '仅可查看' : getStatusText(product.value?.status))
const favoriteCount = computed(() => Number(product.value?.favoriteCount || product.value?.favorite_count || 0))
const visibleImages = computed(() => (product.value?.images || []).filter((image) => !failedImages.value.includes(image)))
const sellerId = computed(() => (
  product.value?.sellerId
  || product.value?.seller_id
  || product.value?.seller?.id
  || product.value?.userId
  || product.value?.user_id
  || product.value?.ownerId
  || product.value?.owner_id
  || ''
))
const sellerSource = computed(() => sellerProfile.value || product.value?.seller || {})
const sellerName = computed(() => displayUserName(sellerSource.value, 'CAU 同学'))
const sellerAvatarUrl = computed(() => sellerAvatarFile.value || normalizeImage(pick(
  sellerSource.value?.avatarUrl,
  sellerSource.value?.avatar,
  sellerSource.value?.avatar_url
)))
const sellerCollege = computed(() => sellerSource.value?.college || '中国农业大学')
const isOwnProduct = computed(() => {
  const user = getUser() || {}
  const currentUserId = pick(user.id, user.userId, user.user_id)
  const currentOpenid = pick(user.openid, user.openId, user.open_id)
  const sellerOpenid = pick(
    product.value?.sellerOpenid,
    product.value?.seller_openid,
    product.value?.seller?.openid,
    product.value?.seller?.openId
  )
  return (currentUserId && sellerId.value && String(sellerId.value) === String(currentUserId))
    || (currentOpenid && sellerOpenid && String(currentOpenid) === String(sellerOpenid))
})
const sellerAvatarText = computed(() => {
  const status = sellerSource.value?.accountStatus || sellerSource.value?.account_status || sellerSource.value?.status
  if (isBannedUserStatus(status) || isCanceledUserStatus(status)) return '停'
  return (sellerName.value || '卖').slice(0, 1)
})
const actionText = computed(() => {
  if (productReadonly.value) return '仅可查看'
  if (isOwnProduct.value) return '自己的商品'
  return product.value?.status === 'ON_SALE' ? '提交预约' : statusText.value
})

function toast(title, icon = 'none') {
  uni.showToast({ title, icon })
}

function ensureVerified() {
  if (productReadonly.value) {
    toast('该商品仅可查看')
    return false
  }
  if (!getToken()) {
    uni.navigateTo({ url: '/pages/login/login' })
    return false
  }
  if (!isVerifiedUser()) {
    uni.navigateTo({ url: '/pages/student-auth/student-auth' })
    return false
  }
  return true
}

function adjustFavoriteCount(delta) {
  product.value.favoriteCount = Math.max(0, favoriteCount.value + delta)
}

function markImageFailed(image) {
  if (!failedImages.value.includes(image)) failedImages.value = failedImages.value.concat(image)
}

function localizeHttpImage(url) {
  if (!/^http:\/\//.test(url || '')) return Promise.resolve(url || '')
  return new Promise((resolve) => {
    uni.downloadFile({
      url,
      success: (res) => resolve(res.tempFilePath || url),
      fail: () => resolve(url)
    })
  })
}

async function loadSellerProfile() {
  if (!sellerId.value) return
  try {
    const profile = await getPublicProfile(sellerId.value)
    const avatar = normalizeImage(pick(profile?.avatarUrl, profile?.avatar, profile?.avatar_url))
    sellerAvatarFile.value = await localizeHttpImage(avatar)
    sellerProfile.value = profile
  } catch (error) {
    sellerProfile.value = null
    sellerAvatarFile.value = ''
  }
}

async function loadFavoriteState(id) {
  isFavorite.value = false
  if (!getToken() || !isVerifiedUser() || isOwnProduct.value) return
  try {
    isFavorite.value = Boolean((await checkFavorite(id)).favorited)
  } catch (error) {
    isFavorite.value = false
  }
}

async function toggleFavorite() {
  if (!ensureVerified()) return
  if (isOwnProduct.value) return toast('不能收藏自己的商品')
  try {
    if (isFavorite.value) {
      await removeFavorite(product.value.id)
      isFavorite.value = false
      adjustFavoriteCount(-1)
    } else {
      await addFavorite(product.value.id)
      isFavorite.value = true
      adjustFavoriteCount(1)
    }
    toast(isFavorite.value ? '收藏成功' : '已取消收藏', 'success')
  } catch (error) {
    const message = String(error?.message || '')
    toast(message.toLowerCase().includes('own') ? '不能收藏自己的商品' : '收藏操作失败')
  }
}

function reserve() {
  if (!ensureVerified() || !product.value?.id) return
  if (isOwnProduct.value) return toast('不能预约自己的商品')
  navigate('/pages/order/appointment', { productId: product.value.id })
}

function report() {
  if (!ensureVerified() || !product.value?.id) return
  if (isOwnProduct.value) return toast('不能举报自己的商品')
  navigate('/pages/interaction/report', { targetType: 'PRODUCT', targetId: product.value.id })
}

function openSeller() {
  if (!sellerId.value) return
  navigate('/pages/user-profile/user-profile', {
    id: sellerId.value,
    productId: product.value?.id,
    productTitle: product.value?.title
  })
}

async function chat() {
  if (!ensureVerified() || !product.value?.id) return
  if (isOwnProduct.value) return toast('不能和自己的商品聊天')
  try {
    const conversation = await createOrGetConversation(product.value.id)
    const currentUserId = Number(getUser()?.id || getUser()?.userId || 0)
    const targetUserId = Number(conversation.buyerId) === currentUserId ? conversation.sellerId : conversation.buyerId
    navigate('/pages/chat/chat', {
      conversationId: conversation.id,
      title: product.value.title,
      targetUserId,
      productId: product.value.id
    })
  } catch (error) {
    toast('暂时无法发起私信')
  }
}

function buildSnapshotProduct(id, options = {}) {
  const image = options.snapshotImage ? decodeURIComponent(options.snapshotImage) : ''
  const title = options.snapshotTitle ? decodeURIComponent(options.snapshotTitle) : '订单商品'
  const price = options.snapshotPrice || 0
  const meetLocation = options.snapshotMeetLocation ? decodeURIComponent(options.snapshotMeetLocation) : '订单约定地点'
  const snapshotSellerName = options.snapshotSellerName ? decodeURIComponent(options.snapshotSellerName) : '卖家'
  return {
    id,
    title,
    price,
    priceText: formatPrice(price),
    status: options.snapshotStatus || 'SOLD',
    images: image ? [normalizeImage(image)] : [],
    coverImage: normalizeImage(image),
    description: '该商品来自订单快照，当前为只读详情。',
    conditionText: '订单商品',
    timeText: '',
    meetLocation,
    sellerId: options.snapshotSellerId || '',
    seller: { id: options.snapshotSellerId || '', nickname: snapshotSellerName, college: '中国农业大学' }
  }
}

function getDetailErrorText(error) {
  const message = String(error?.message || '').toLowerCase()
  if (message.includes('not found')) return '商品不存在或已下架'
  if (message.includes('permission') || message.includes('forbidden')) return '暂无权限查看该商品'
  return '商品暂不可查看'
}

onLoad(async (options) => {
  const { id } = options
  readonlyMode.value = options.readonly === '1' || options.readonly === 1
  if (!id) {
    toast('商品不存在')
    return
  }

  try {
    const [detail, categories] = await Promise.all([getProductById(id), listCategories()])
    product.value = formatProduct(detail, buildCategoryMap(categories))
    uni.setStorageSync(`product-detail-cache-${id}`, product.value)
    addBrowseHistory(product.value)
    failedImages.value = []
    await loadSellerProfile()
    await loadFavoriteState(id)
  } catch (error) {
    const cached = uni.getStorageSync(`product-detail-cache-${id}`)
    if (cached) {
      product.value = cached
      failedImages.value = []
      await loadSellerProfile()
      await loadFavoriteState(id)
      if (!readonlyMode.value) toast('商品暂不可查看，显示最近一次详情')
      return
    }
    if (readonlyMode.value) {
      product.value = buildSnapshotProduct(id, options)
      failedImages.value = []
      await loadSellerProfile()
      return
    }
    toast(getDetailErrorText(error))
  }
})
</script>

<style scoped>
.page { min-height: 100vh; padding-bottom: 130rpx; background: #f5f6f8; }
.gallery, .gallery-image { width: 100%; height: 600rpx; background: #e8efeb; }
.placeholder { display: flex; align-items: center; justify-content: center; color: #9aa5a1; font-size: 28rpx; }
.card { margin: 20rpx; padding: 24rpx; border-radius: 18rpx; background: #fff; }
.price-line, .seller { display: flex; align-items: center; justify-content: space-between; }
.price { color: #e36a3e; font-size: 46rpx; font-weight: 700; }
.status { padding: 8rpx 14rpx; border-radius: 999rpx; background: #e7f4ec; color: #23734f; font-size: 23rpx; }
.title { margin-top: 14rpx; font-size: 36rpx; font-weight: 700; color: #1f2933; }
.meta { display: flex; flex-wrap: wrap; gap: 12rpx 18rpx; margin-top: 14rpx; color: #89938f; font-size: 23rpx; }
.meta.single { display: block; margin-top: 6rpx; }
.section-title, .seller-name { font-weight: 700; color: #1f2933; }
.description { margin-top: 16rpx; color: #58645f; line-height: 1.7; white-space: pre-line; word-break: break-word; }
.seller { justify-content: flex-start; gap: 16rpx; }
.seller-label { color: #9aa5a1; font-size: 24rpx; }
.avatar { display: flex; width: 84rpx; height: 84rpx; flex-shrink: 0; align-items: center; justify-content: center; border-radius: 50%; background: #e7f4ec; color: #23734f; font-weight: 700; }
.image-avatar { background: #e8ecef; }
.seller-body { flex: 1; min-width: 0; }
.seller-arrow { color: #98a2b3; font-size: 42rpx; }
.readonly-tip { margin: 20rpx; padding: 20rpx 24rpx; border-radius: 16rpx; background: #fff7e6; color: #a15c00; font-size: 26rpx; }
.bottom { position: fixed; right: 0; bottom: 0; left: 0; display: flex; gap: 12rpx; padding: 16rpx 20rpx calc(16rpx + env(safe-area-inset-bottom)); background: #fff; }
.icon-button, .chat, .primary { height: 74rpx; border-radius: 999rpx; font-size: 26rpx; line-height: 74rpx; }
.icon-button { width: 84rpx; padding: 0; color: #4b5a52; background: #f2f5f3; }
.favorite.active { color: #f59e0b; background: #fff7e6; }
.report.disabled { color: #c3cac6; }
.chat { flex: 1; color: #23734f; background: #e7f4ec; }
.primary { flex: 1.4; color: #fff; background: #23734f; }
button[disabled] { opacity: .48; }
.loading-page { display: flex; align-items: center; justify-content: center; color: #667085; }
</style>
