<template>
  <view class="page">
    <view class="search-box" @click="goSearch">
      搜索二手商品
    </view>

    <view class="section-title">商品分类</view>

    <view class="category-list">
      <view class="category-item" @click="goCategory('book')">教材书籍</view>
      <view class="category-item" @click="goCategory('digital')">数码电子</view>
      <view class="category-item" @click="goCategory('daily')">生活用品</view>
      <view class="category-item" @click="goCategory('sport')">运动户外</view>
    </view>

    <view class="section-title">最新商品</view>

    <view class="goods-list">
      <view
        class="goods-card"
        v-for="item in goodsList"
        :key="item.id"
        @click="goDetail(item.id)"
      >
        <view class="goods-image"></view>
        <view class="goods-info">
          <view class="goods-title">{{ item.title }}</view>
          <view class="goods-desc">{{ item.condition }}</view>
          <view class="goods-price">￥{{ item.price }}</view>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
const goodsList = [
  {
    id: 1,
    title: '高等数学教材',
    condition: '八成新',
    price: 20
  },
  {
    id: 2,
    title: '蓝牙耳机',
    condition: '九成新',
    price: 68
  },
  {
    id: 3,
    title: '台灯',
    condition: '七成新',
    price: 15
  }
]

const goSearch = () => {
  uni.navigateTo({
    url: '/pages/search/search'
  })
}

const goCategory = (category) => {
  uni.navigateTo({
    url: `/pages/category/category?category=${category}`
  })
}

const goDetail = (id) => {
  uni.navigateTo({
    url: `/pages/detail/detail?id=${id}`
  })
}
</script>

<style scoped>
.page {
  min-height: 100vh;
  padding: 32rpx;
  background: #f6f7f9;
  box-sizing: border-box;
}

.search-box {
  height: 72rpx;
  line-height: 72rpx;
  padding: 0 28rpx;
  border-radius: 36rpx;
  background: #ffffff;
  color: #9ca3af;
  font-size: 28rpx;
}

.section-title {
  margin-top: 36rpx;
  margin-bottom: 20rpx;
  font-size: 34rpx;
  font-weight: 700;
  color: #1f2933;
}

.category-list {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 20rpx;
}

.category-item {
  height: 96rpx;
  line-height: 96rpx;
  border-radius: 16rpx;
  background: #ffffff;
  text-align: center;
  font-size: 28rpx;
  color: #374151;
}

.goods-list {
  display: flex;
  flex-direction: column;
  gap: 20rpx;
}

.goods-card {
  display: flex;
  padding: 20rpx;
  border-radius: 16rpx;
  background: #ffffff;
}

.goods-image {
  width: 140rpx;
  height: 140rpx;
  border-radius: 12rpx;
  background: #d1d5db;
  flex-shrink: 0;
}

.goods-info {
  margin-left: 24rpx;
  flex: 1;
}

.goods-title {
  font-size: 30rpx;
  font-weight: 600;
  color: #1f2933;
}

.goods-desc {
  margin-top: 16rpx;
  font-size: 24rpx;
  color: #6b7280;
}

.goods-price {
  margin-top: 20rpx;
  font-size: 32rpx;
  font-weight: 700;
  color: #e11d48;
}
</style>