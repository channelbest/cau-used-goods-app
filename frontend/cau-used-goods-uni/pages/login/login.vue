<template>
  <view class="page">
    <view class="title">CAU Used Goods</view>
    <view class="subtitle">Welcome to campus second-hand market</view>

    <button class="login-button" @click="goStudentAuth">
      WeChat Login
    </button>
  </view>
</template>

<script setup>
const goStudentAuth = () => {
  uni.navigateTo({
    url: '/pages/student-auth/student-auth'
  })
}
</script>

<style scoped>
.page {
  min-height: 100vh;
  padding: 120rpx 48rpx;
  background: #f6f7f9;
  box-sizing: border-box;
}

.title {
  margin-top: 120rpx;
  font-size: 44rpx;
  font-weight: 700;
  color: #1f2933;
  text-align: center;
}

.subtitle {
  margin-top: 24rpx;
  font-size: 28rpx;
  color: #6b7280;
  text-align: center;
}

.login-button {
  margin-top: 120rpx;
  height: 88rpx;
  line-height: 88rpx;
  border-radius: 12rpx;
  background: #1aad19;
  color: #ffffff;
  font-size: 32rpx;
}
</style>
