<template>
  <view class="page">
    <view class="title">CAU 二手交易平台</view>
    <view class="subtitle">登录后继续使用校园二手交易服务</view>

    <view class="card">
      <button class="login-button" :loading="loading" @click="handleWechatLogin">微信登录</button>

      <view class="dev-divider">开发调试登录</view>
      <text class="label">测试 OpenID</text>
      <input v-model="openid" class="input" placeholder="请输入 openid，或选择测试账号" />

      <view class="quick-title">切换测试账号</view>
      <view class="quick-list">
        <view
          v-for="account in accounts"
          :key="account.openid"
          :class="['quick-chip', openid === account.openid ? 'active' : '']"
          @click="selectAccount(account.openid)"
        >
          {{ account.label }}
        </view>
      </view>

      <button class="login-button secondary" :loading="loading" @click="handleDevLogin">使用测试账号登录</button>
    </view>
  </view>
</template>

<script setup>
import { ref } from 'vue'
import { devLogin, reactivateAccount, wechatLogin } from '../../api/auth'
import { saveLoginResult } from '../../utils/auth'

const loading = ref(false)
const openid = ref(uni.getStorageSync('dev-login-openid') || '')
const accounts = [
  { label: '用户 A', openid: 'frontend_a_dev_user' },
  { label: '用户 B', openid: 'frontend_b_dev_user' },
  { label: '管理员', openid: 'admin_dev_user', role: 'ADMIN' },
  { label: '超级管理员', openid: 'super_admin_dev_user', role: 'SUPER_ADMIN' },
  { label: '待认证', openid: 'pending_dev_user' }
]
const adminOpenids = ['admin_dev_user', 'super_admin_dev_user']

function isAdminLogin(value, user = {}) {
  const role = String(user.role || user.userRole || user.user_role || user.type || user.userType || user.user_type || '').toUpperCase()
  return adminOpenids.includes(value) || role === 'ADMIN' || role === 'ROLE_ADMIN' || role.includes('ADMIN') || user.isAdmin === true || user.admin === true
}

function selectAccount(value) {
  openid.value = value
}

function selectedAccount() {
  return accounts.find((account) => account.openid === openid.value.trim())
}

function goHome() {
  uni.switchTab({ url: '/pages/home/home' })
}

function goAdmin() {
  uni.reLaunch({ url: '/pages/admin/admin' })
}

function loginWithWechatCode() {
  return new Promise((resolve, reject) => {
    uni.login({
      provider: 'weixin',
      success: ({ code }) => {
        if (!code) {
          reject(new Error('微信登录凭证为空'))
          return
        }
        resolve(code)
      },
      fail: () => reject(new Error('微信登录失败，请稍后重试'))
    })
  })
}

function confirmModal(options) {
  return new Promise((resolve) => {
    uni.showModal({
      ...options,
      success: (res) => resolve(Boolean(res.confirm))
    })
  })
}

function routeAfterLogin(result) {
  if (isAdminLogin('', result?.user)) {
    goAdmin()
    return
  }
  goHome()
}

async function handleReactivation(result) {
  const confirmed = await confirmModal({
    title: '恢复账号',
    content: '该账号此前已注销。是否恢复原账号并继续登录？',
    confirmText: '恢复',
    cancelText: '取消'
  })
  if (!confirmed) {
    uni.showToast({ title: '已取消恢复账号', icon: 'none' })
    return
  }
  const restored = await reactivateAccount(result.reactivationToken)
  saveLoginResult(restored)
  uni.showToast({ title: '账号已恢复', icon: 'success' })
  routeAfterLogin(restored)
}

async function finishLogin(result) {
  if (result?.requiresReactivation && result?.reactivationToken) {
    await handleReactivation(result)
    return
  }
  saveLoginResult(result)
  uni.showToast({ title: '登录成功', icon: 'success' })
  routeAfterLogin(result)
}

async function handleWechatLogin() {
  if (loading.value) return
  loading.value = true
  try {
    const code = await loginWithWechatCode()
    const result = await wechatLogin(code)
    await finishLogin(result)
  } catch (error) {
    uni.showToast({ title: error.message || '登录失败', icon: 'none' })
  } finally {
    loading.value = false
  }
}

async function handleDevLogin() {
  const value = openid.value.trim()
  if (!value) {
    uni.showToast({ title: '请输入 openid', icon: 'none' })
    return
  }
  if (loading.value) return
  loading.value = true
  try {
    const account = selectedAccount()
    const payload = account?.role ? { openid: value, role: account.role } : { openid: value }
    const result = await devLogin(payload)
    uni.setStorageSync('dev-login-openid', value)
    if (isAdminLogin(value, result?.user)) {
      result.user = { ...(result.user || {}), role: account?.role || result?.user?.role || 'ADMIN' }
    }
    await finishLogin(result)
  } catch (error) {
    uni.showToast({ title: error.message || '登录失败', icon: 'none' })
  } finally {
    loading.value = false
  }
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
  font-weight: 800;
  color: #1f2933;
  text-align: center;
}

.subtitle {
  margin-top: 14rpx;
  color: #667085;
  font-size: 26rpx;
  text-align: center;
}

.card {
  margin-top: 70rpx;
  padding: 30rpx;
  border-radius: 20rpx;
  background: #fff;
  box-shadow: 0 10rpx 30rpx rgba(31, 41, 51, .06);
}

.label,
.quick-title {
  display: block;
  color: #344054;
  font-size: 26rpx;
  font-weight: 700;
}

.dev-divider {
  display: flex;
  align-items: center;
  gap: 18rpx;
  margin: 36rpx 0 24rpx;
  color: #98a2b3;
  font-size: 24rpx;
}

.dev-divider::before,
.dev-divider::after {
  content: '';
  flex: 1;
  height: 1rpx;
  background: #eef0f3;
}

.input {
  height: 88rpx;
  margin-top: 18rpx;
  padding: 0 24rpx;
  border-radius: 14rpx;
  background: #f2f4f7;
  color: #1f2933;
  font-size: 28rpx;
  box-sizing: border-box;
}

.quick-title {
  margin-top: 34rpx;
}

.quick-list {
  display: flex;
  flex-wrap: wrap;
  gap: 18rpx;
  margin-top: 18rpx;
}

.quick-chip {
  min-width: 150rpx;
  padding: 16rpx 22rpx;
  border: 1rpx solid #d0d5dd;
  border-radius: 999rpx;
  color: #475467;
  font-size: 25rpx;
  text-align: center;
  box-sizing: border-box;
}

.quick-chip.active {
  border-color: #14b86e;
  background: #e9fbf3;
  color: #079455;
  font-weight: 700;
}

.login-button {
  margin-top: 44rpx;
  height: 88rpx;
  line-height: 88rpx;
  border-radius: 999rpx;
  background: #14b86e;
  color: #fff;
  font-size: 30rpx;
  font-weight: 700;
}

.login-button.secondary {
  margin-top: 34rpx;
  background: #eef7f0;
  color: #17a84b;
}
</style>
