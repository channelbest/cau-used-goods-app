<template>
  <view class="page">
    <view class="goods-image"></view>

    <view class="info-card">
      <view class="title">高等数学教材</view>
      <view class="price">￥20</view>
      <view class="meta">分类：教材书籍</view>
      <view class="meta">成色：八成新</view>
      <view class="meta">发布时间：2026-05-27</view>
    </view>

    <view class="info-card">
      <view class="section-title">商品描述</view>
      <view class="desc">
        课程使用教材，书页完整，有少量笔记，适合自用复习。
      </view>
    </view>

    <view class="info-card">
      <view class="section-title">卖家信息</view>
      <view class="meta">卖家：CAU学生</view>
      <view class="meta">联系方式需预约后在订单详情中查看</view>
    </view>

    <view class="button-row">
      <button class="plain-button">收藏</button>
      <button class="main-button">提交预约</button>
      <button class="danger-button">举报</button>
    </view>
  </view>
</template>

<script setup>
</script>

<style scoped>
.page {
  min-height: 100vh;
  padding: 32rpx;
  padding-bottom: 140rpx;
  background: #f6f7f9;
  box-sizing: border-box;
}

.goods-image {
  height: 420rpx;
  border-radius: 20rpx;
  background: #d1d5db;
}

.info-card {
  margin-top: 24rpx;
  padding: 28rpx;
  border-radius: 16rpx;
  background: #ffffff;
}

.title {
  font-size: 38rpx;
  font-weight: 700;
  color: #1f2933;
}

.price {
  margin-top: 20rpx;
  font-size: 42rpx;
  font-weight: 700;
  color: #e11d48;
}

.meta {
  margin-top: 16rpx;
  font-size: 26rpx;
  color: #6b7280;
}

.section-title {
  font-size: 32rpx;
  font-weight: 700;
  color: #1f2933;
}

.desc {
  margin-top: 16rpx;
  line-height: 1.7;
  font-size: 28rpx;
  color: #374151;
}

.button-row {
  position: fixed;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  gap: 16rpx;
  padding: 20rpx 32rpx;
  background: #ffffff;
  box-sizing: border-box;
}

.plain-button,
.main-button,
.danger-button {
  flex: 1;
  height: 76rpx;
  line-height: 76rpx;
  border-radius: 12rpx;
  font-size: 26rpx;
}

.plain-button {
  background: #f3f4f6;
  color: #374151;
}

.main-button {
  background: #1aad19;
  color: #ffffff;
}

.danger-button {
  background: #fee2e2;
  color: #dc2626;
}
</style>