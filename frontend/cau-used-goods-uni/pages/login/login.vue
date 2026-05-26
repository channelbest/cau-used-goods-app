<template>
  <view class="page">
    <view class="title">CAU Used Goods</view>
    <view class="subtitle">Campus second-hand market</view>

    <button class="login-button" :loading="loading" @click="handleDevLogin">
      Dev Login
    </button>

    <button class="secondary-button" @click="previewStudentAuth">
      Preview Student Verification
    </button>

    <view class="tip">
      Dev Login calls /auth/dev-login and stores the returned token.
    </view>
  </view>
</template>

<script setup>
import { ref } from 'vue'
import { devLogin } from '../../api/auth'
import { saveLoginResult } from '../../utils/auth'

const loading = ref(false)

const goNext = (user) => {
  if (user?.authStatus === 'VERIFIED') {
    uni.navigateTo({
      url: '/pages/home/home'
    })
    return
  }

  uni.navigateTo({
    url: '/pages/student-auth/student-auth'
  })
}

const handleDevLogin = async () => {
  if (loading.value) return

  loading.value = true
  try {
    const result = await devLogin({
      openid: 'frontend_a_dev_user',
      role: 'USER'
    })

    saveLoginResult(result)
    uni.showToast({
      title: 'Login success',
      icon: 'success'
    })
    goNext(result.user)
  } catch (error) {
    uni.showToast({
      title: error.message || 'Login failed',
      icon: 'none'
    })
  } finally {
    loading.value = false
  }
}

const previewStudentAuth = () => {
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

.secondary-button {
  margin-top: 24rpx;
  height: 88rpx;
  line-height: 88rpx;
  border-radius: 12rpx;
  background: #ffffff;
  color: #374151;
  font-size: 30rpx;
}

.tip {
  margin-top: 32rpx;
  line-height: 1.6;
  font-size: 24rpx;
  color: #9ca3af;
  text-align: center;
}
</style>
