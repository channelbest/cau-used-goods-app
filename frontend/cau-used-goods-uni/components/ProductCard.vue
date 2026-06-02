<template>
  <view class="card" @click="openDetail">
    <image v-if="product.coverImage" class="cover" :src="product.coverImage" mode="aspectFill" />
    <view v-else class="cover placeholder">暂无图片</view>
    <view class="body">
      <view class="title">{{ product.title }}</view>
      <view class="meta">
        <text>{{ product.conditionText }}</text>
        <text>{{ product.timeText }}</text>
      </view>
      <view class="footer">
        <text class="price">￥{{ product.priceText }}</text>
        <text class="category">{{ product.category }}</text>
      </view>
    </view>
  </view>
</template>

<script setup>
const props = defineProps({
  product: {
    type: Object,
    required: true
  }
})

const openDetail = () => {
  uni.navigateTo({ url: `/pages/detail/detail?id=${props.product.id}` })
}
</script>

<style scoped>
.card { overflow: hidden; border-radius: 18rpx; background: #fff; box-shadow: 0 6rpx 18rpx rgba(28, 68, 52, .06); }
.cover { display: flex; width: 100%; height: 250rpx; align-items: center; justify-content: center; background: #e8efeb; color: #9aa5a1; font-size: 24rpx; }
.body { padding: 18rpx; }
.title { display: -webkit-box; overflow: hidden; min-height: 78rpx; color: #26342f; font-size: 28rpx; font-weight: 600; line-height: 1.4; -webkit-box-orient: vertical; -webkit-line-clamp: 2; }
.meta, .footer { display: flex; align-items: center; justify-content: space-between; margin-top: 14rpx; color: #89938f; font-size: 21rpx; }
.price { color: #e36a3e; font-size: 34rpx; font-weight: 700; }
.category { color: #5a756b; }
</style>
