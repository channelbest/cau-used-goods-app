<template>
  <view class="page">
    <view class="title">学生认证</view>

    <view class="form-item">
      <text class="label">姓名</text>
      <input class="input" placeholder="请输入姓名" />
    </view>

    <view class="form-item">
      <text class="label">学号</text>
      <input class="input" placeholder="请输入学号" />
    </view>

    <view class="form-item">
      <text class="label">学院</text>
      <input class="input" placeholder="请输入学院" />
    </view>

    <button class="submit-button" @click="goHome">
      提交认证
    </button>
  </view>
</template>

<script setup>
const goHome = () => {
  uni.navigateTo({
    url: '/pages/home/home'
  })
}
</script>

<style scoped>
.page {
  min-height: 100vh;
  padding: 48rpx;
  background: #f6f7f9;
  box-sizing: border-box;
}

.title {
  margin-bottom: 48rpx;
  font-size: 40rpx;
  font-weight: 700;
  color: #1f2933;
}

.form-item {
  margin-bottom: 32rpx;
}

.label {
  display: block;
  margin-bottom: 12rpx;
  font-size: 28rpx;
  color: #374151;
}

.input {
  height: 88rpx;
  padding: 0 24rpx;
  border-radius: 12rpx;
  background: #ffffff;
  font-size: 28rpx;
}

.submit-button {
  margin-top: 56rpx;
  height: 88rpx;
  line-height: 88rpx;
  border-radius: 12rpx;
  background: #1aad19;
  color: #ffffff;
  font-size: 32rpx;
}
</style>