<template>
  <view class="page">
    <view class="profile card">
      <view class="avatar">{{ (user?.nickname || '我').slice(0, 1) }}</view>
      <view>
        <view class="name">{{ user?.nickname || '未登录用户' }}</view>
        <view class="muted">{{ authText }}</view>
      </view>
    </view>
    <view class="card menu">
      <view class="item" @click="goProfile">编辑个人资料 <text>›</text></view>
      <view class="item" @click="goStudentAuth">学生认证 <text>›</text></view>
      <view class="item" @click="goPublish">发布闲置商品 <text>›</text></view>
    </view>
    <button v-if="user" class="logout" @click="logout">退出登录</button>
    <button v-else class="login" @click="goLogin">去登录</button>
  </view>
</template>

<script setup>
import { computed, ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import { clearAuth, getUser } from '../../utils/auth'
const user = ref(null)
const authText = computed(() => user.value?.authStatus === 'VERIFIED' ? '学生认证已通过' : user.value ? '学生认证待完善' : '登录后发布和预约商品')
const goLogin = () => uni.navigateTo({ url: '/pages/login/login' })
const ensureLogin = () => { if (!user.value) { goLogin(); return false } return true }
const goProfile = () => { if (ensureLogin()) uni.navigateTo({ url: '/pages/profile-edit/profile-edit' }) }
const goStudentAuth = () => { if (ensureLogin()) uni.navigateTo({ url: '/pages/student-auth/student-auth' }) }
const goPublish = () => { if (ensureLogin()) uni.switchTab({ url: '/pages/publish/publish' }) }
const logout = () => { clearAuth(); user.value = null; uni.showToast({ title: '已退出登录', icon: 'success' }) }
onShow(() => { user.value = getUser() })
</script>

<style scoped>
.page { min-height: 100vh; padding: 28rpx; box-sizing: border-box; }.card { border-radius: 18rpx; background: #fff; }.profile { display: flex; align-items: center; padding: 32rpx; }.avatar { display: flex; width: 92rpx; height: 92rpx; margin-right: 18rpx; align-items: center; justify-content: center; border-radius: 50%; background: #e7f4ec; color: #23734f; font-size: 34rpx; font-weight: 700; }.name { font-size: 32rpx; font-weight: 700; }.muted { margin-top: 8rpx; color: #929c98; font-size: 23rpx; }
.menu { margin-top: 22rpx; padding: 0 22rpx; }.item { display: flex; justify-content: space-between; padding: 28rpx 0; border-bottom: 2rpx solid #f0f2f1; }.item:last-child { border-bottom: 0; }.item text { color: #23734f; }
.logout,.login { margin-top: 28rpx; border-radius: 999rpx; }.logout { background: #fff; color: #b85d45; }.login { background: #23734f; color: #fff; }
</style>
