<template>
  <view class="page">
    <view class="card">
      <text class="section-title">这次交易体验如何？</text>
      <view class="stars">
        <text v-for="value in 5" :key="value" :class="{ active: value <= form.rating }" @click="form.rating = value">★</text>
      </view>
      <view class="field">
        <text class="field-label">评价内容</text>
        <textarea v-model="form.content" class="textarea" maxlength="300" placeholder="请分享本次线下面交体验" />
      </view>
      <view class="anonymous">
        <switch :checked="form.anonymous" color="#2f6b4f" @change="form.anonymous = $event.detail.value" />
        <text>匿名评价</text>
      </view>
      <button class="btn btn-primary" @click="submit">提交评价</button>
    </view>
  </view>
</template>

<script setup>
import { onLoad } from '@dcloudio/uni-app'
import { reactive, ref } from 'vue'
import { tradeService } from '../../services/trade'
import { showError, showSuccess } from '../../utils/navigation'

const orderId = ref('')
const order = ref(null)
const form = reactive({ rating: 5, content: '', anonymous: false })
onLoad(async (options) => {
  orderId.value = options.orderId
  if (!orderId.value) return
  order.value = await tradeService.getOrder(orderId.value).catch(() => null)
})

function saveSellerReview() {
  const sellerId = order.value?.sellerId || order.value?.seller?.id
  if (!sellerId) return
  const key = `user-reviews-${sellerId}`
  const list = uni.getStorageSync(key) || []
  const next = [{
    id: Date.now(),
    rating: form.rating,
    content: form.content,
    anonymous: form.anonymous,
    productTitle: order.value?.product?.title || order.value?.productTitleSnapshot || '交易商品',
    createTime: new Date().toLocaleString()
  }].concat(Array.isArray(list) ? list : [])
  uni.setStorageSync(key, next.slice(0, 20))
}

async function submit() {
  if (!form.content.trim()) {
    showError(new Error('请填写评价内容'))
    return
  }
  try {
    await tradeService.createReview({ orderId: orderId.value, ...form })
    uni.setStorageSync(`order-reviewed-${orderId.value}`, true)
    saveSellerReview()
    showSuccess('评价成功')
    setTimeout(() => uni.navigateBack(), 500)
  } catch (error) {
    showError(error)
  }
}
</script>

<style scoped lang="scss">
.stars { display: flex; gap: 14rpx; margin: 20rpx 0 32rpx; }
.stars text { color: #d9dfdc; font-size: 64rpx; }
.stars .active { color: #f2a23a; }
.anonymous { display: flex; align-items: center; gap: 14rpx; margin: 12rpx 0 28rpx; color: #738077; font-size: 25rpx; }
</style>
