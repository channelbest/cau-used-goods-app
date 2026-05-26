<template>
  <view class="page">
    <view class="title">分类商品</view>

    <view class="sort-row">
      <view class="sort-item active">最新发布</view>
      <view class="sort-item">价格最低</view>
      <view class="sort-item">价格最高</view>
    </view>

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
    title: '教材书籍示例商品',
    condition: '八成新',
    price: 20
  },
  {
    id: 2,
    title: '分类下的商品',
    condition: '九成新',
    price: 35
  }
]

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

.title {
  font-size: 40rpx;
  font-weight: 700;
  color: #1f2933;
}

.sort-row {
  display: flex;
  gap: 16rpx;
  margin-top: 28rpx;
  margin-bottom: 28rpx;
}

.sort-item {
  flex: 1;
  height: 64rpx;
  line-height: 64rpx;
  border-radius: 12rpx;
  background: #ffffff;
  text-align: center;
  font-size: 26rpx;
  color: #374151;
}

.sort-item.active {
  background: #1aad19;
  color: #ffffff;
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