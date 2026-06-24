<template>
  <view class="page">
    <view class="section-title">商品状态</view>

    <image v-if="mainImage" class="product-image" :src="mainImage" mode="aspectFill" @click="previewProductImage" />
    <view v-else class="product-image placeholder">暂无图片</view>

    <view class="card">
      <view class="name">{{ product.title || '商品' }}</view>
      <view class="price">￥{{ product.price || 0 }}</view>
      <view class="desc">商品ID：{{ product.id || '-' }}</view>
      <view class="desc">卖家ID：{{ sellerId || '-' }}</view>
      <view class="desc">当前状态：{{ statusText(product.status) }}</view>
      <view class="desc">成色：{{ conditionText(product.conditionLevel) }}</view>
      <view class="desc">交易地点：{{ product.meetLocation || '线下面交' }}</view>
      <button class="profile-button" :disabled="!sellerId" @click="goSellerProfile">查看用户主页</button>
    </view>

    <view class="actions">
      <button v-if="product.status === 'OFF_SHELF'" class="pass" @click="changeStatus('ON_SALE')">上架商品</button>
      <button v-if="product.status === 'ON_SALE'" class="reject" @click="changeStatus('OFF_SHELF')">下架商品</button>
    </view>
  </view>
</template>

<script setup>
import { computed, ref } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { getAdminProductById, updateAdminProductStatus } from '../../api/admin'
import { normalizeImage } from '../../utils/product-format'

const product = ref({})
const productId = ref('')
const relatedType = ref('')
const relatedId = ref(0)
const mainImage = computed(() => normalizeImage(product.value?.images?.[0] || ''))
const sellerId = computed(() => product.value?.sellerId || product.value?.seller_id || product.value?.seller?.id || '')

onLoad((query = {}) => {
  productId.value = query.id || ''
  relatedType.value = String(query.relatedType || '').toUpperCase()
  relatedId.value = Number(query.relatedId || 0)
  load()
})

const load = async () => {
  if (!productId.value) {
    uni.showToast({ title: '商品不存在', icon: 'none' })
    return
  }
  try {
    product.value = await getAdminProductById(productId.value)
  } catch (error) {
    uni.showToast({ title: error.message || '商品加载失败', icon: 'none' })
  }
}

const changeStatus = async (status) => {
  try {
    const extra = {}
    if (relatedType.value && relatedId.value) {
      extra.relatedType = relatedType.value
      extra.relatedId = relatedId.value
    }
    await updateAdminProductStatus(productId.value, status, extra)
    uni.showToast({ title: status === 'ON_SALE' ? '已上架' : '已下架', icon: 'success' })
    product.value = { ...product.value, status }
  } catch (error) {
    uni.showToast({ title: error.message || '商品状态更新失败', icon: 'none' })
  }
}

const statusText = (status) => {
  const map = { ON_SALE: '在售', OFF_SHELF: '已下架', LOCKED: '交易锁定', SOLD: '已售出', DELETED: '已删除' }
  return map[status] || status || '未知'
}

const conditionText = (level) => {
  const map = { NEW: '全新', LIKE_NEW: '九成新', GOOD: '八成新', FAIR: '七成新', OLD: '旧物' }
  return map[level] || level || '成色未填写'
}

const previewProductImage = () => {
  const urls = (product.value?.images || []).map((url) => normalizeImage(url)).filter(Boolean)
  if (!urls.length) return
  uni.previewImage({ current: urls[0], urls })
}

const goSellerProfile = () => {
  if (!sellerId.value) {
    uni.showToast({ title: '暂无卖家信息', icon: 'none' })
    return
  }
  const query = [`id=${sellerId.value}`, 'adminView=1']
  if (relatedType.value && relatedId.value) {
    query.push(`relatedType=${relatedType.value}`, `relatedId=${relatedId.value}`)
  }
  uni.navigateTo({ url: `/pages/user-profile/user-profile?${query.join('&')}` })
}
</script>

<style scoped>
.page { min-height: 100vh; padding: 24rpx; background: #f5f6f8; box-sizing: border-box; }
.section-title { margin: 24rpx 0 18rpx; font-size: 32rpx; font-weight: 700; color: #1f2933; }
.product-image { width: 100%; height: 360rpx; margin-bottom: 18rpx; border-radius: 16rpx; background: #eef2f6; }
.placeholder { display: flex; align-items: center; justify-content: center; color: #98a2b3; font-size: 28rpx; }
.card { padding: 28rpx; border-radius: 16rpx; background: #fff; margin-bottom: 18rpx; }
.name { font-size: 34rpx; font-weight: 700; color: #1f2933; }
.price { margin-top: 18rpx; font-size: 40rpx; font-weight: 700; color: #e11d48; }
.desc { margin-top: 12rpx; color: #667085; font-size: 26rpx; line-height: 1.6; }
.profile-button { height: 72rpx; margin-top: 22rpx; border-radius: 12rpx; background: #23734f; color: #fff; font-size: 27rpx; line-height: 72rpx; }
.profile-button[disabled] { background: #d0d5dd; color: #fff; }
.actions { margin-top: 24rpx; }
.pass, .reject { height: 88rpx; line-height: 88rpx; border-radius: 12rpx; font-size: 30rpx; }
.pass { background: #17a84b; color: #fff; }
.reject { background: #fff1f2; color: #ef4444; }
</style>
