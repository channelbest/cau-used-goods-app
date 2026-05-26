<template>
  <view class="page">
    <view class="search-row">
      <input class="search-input" v-model="keyword" placeholder="请输入关键词" />
      <button class="search-button">搜索</button>
    </view>

    <view class="filter-row">
      <view class="filter-item">分类</view>
      <view class="filter-item">价格</view>
      <view class="filter-item">成色</view>
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
import { ref } from 'vue'

const keyword = ref('')

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

.search-row {
  display: flex;
  gap: 16rpx;
}

.search-input {
  flex: 1;
  height: 76rpx;
  padding: 0 24rpx;
  border-radius: 12rpx;
  background: #ffffff;
  font-size: 28rpx;
}

.search-button {
  width: 140rpx;
  height: 76rpx;
  line-height: 76rpx;
  border-radius: 12rpx;
  background: #1aad19;
  color: #ffffff;
  font-size: 28rpx;
}

.filter-row {
  display: flex;
  gap: 16rpx;
  margin-top: 24rpx;
  margin-bottom: 24rpx;
}

.filter-item {
  flex: 1;
  height: 64rpx;
  line-height: 64rpx;
  border-radius: 12rpx;
  background: #ffffff;
  text-align: center;
  font-size: 26rpx;
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