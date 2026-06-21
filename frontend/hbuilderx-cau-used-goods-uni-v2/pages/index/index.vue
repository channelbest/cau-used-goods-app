<template>
  <view class="page">
    <view class="profile-hero">
      <text class="hero-eyebrow">CAU CAMPUS MARKET</text>
      <text class="hero-title">&#20010;&#20154;&#20013;&#24515;</text>
      <text class="hero-copy">&#31649;&#29702;&#20320;&#30340;&#38386;&#32622;&#12289;&#35746;&#21333;&#21644;&#26657;&#22253;&#36523;&#20221;</text>
    </view>

    <view class="profile-card">
      <image v-if="avatarUrl" class="avatar" :src="avatarUrl" mode="aspectFill" />
      <view v-else class="avatar placeholder">&#22836;&#20687;</view>
      <view class="profile-main">
        <view class="nickname">{{ user.nickname || '\u5fae\u4fe1\u7528\u6237' }}</view>
        <view :class="['status', authTone]">{{ identityText }}</view>
      </view>
      <button class="edit-button" size="mini" @click="goProfileEdit">&#20462;&#25913;&#36164;&#26009;</button>
    </view>

    <view v-if="!isAdmin" class="section">
      <view class="section-title">&#25105;&#30340;&#20132;&#26131;</view>
      <view class="quick-list">
        <view :class="['quick-card', !isVerified ? 'locked' : '']" @click="goMyProducts">
          <view class="quick-icon">&#21457;</view>
          <view class="quick-copy">
            <view class="quick-title">&#25105;&#21457;&#24067;&#30340;</view>
            <view class="quick-desc">&#31649;&#29702;&#38386;&#32622;&#21830;&#21697;</view>
          </view>
          <text class="arrow">&#8250;</text>
        </view>
        <view :class="['quick-card', !isVerified ? 'locked' : '']" @click="goSoldOrders">
          <view class="quick-icon warm">&#21334;</view>
          <view class="quick-copy">
            <view class="quick-title">&#25105;&#21334;&#20986;&#30340;</view>
            <view class="quick-desc">&#21334;&#23478;&#35746;&#21333;</view>
          </view>
          <text class="arrow">&#8250;</text>
        </view>
        <view :class="['quick-card', !isVerified ? 'locked' : '']" @click="goBoughtOrders">
          <view class="quick-icon cool">&#20080;</view>
          <view class="quick-copy">
            <view class="quick-title">&#25105;&#20080;&#21040;&#30340;</view>
            <view class="quick-desc">&#20080;&#23478;&#35746;&#21333;</view>
          </view>
          <text class="arrow">&#8250;</text>
        </view>
      </view>
    </view>

    <view v-if="!isAdmin" class="menu-card">
      <view :class="['menu-item', !isVerified ? 'locked' : '']" @click="goFavorites">
        <view>
          <view class="menu-title">&#25105;&#30340;&#25910;&#34255;</view>
          <view class="menu-desc">&#20445;&#23384;&#24863;&#20852;&#36259;&#30340;&#21830;&#21697;</view>
        </view>
        <text>&#8250;</text>
      </view>
      <view :class="['menu-item', !isVerified ? 'locked' : '']" @click="goReportList">
        <view>
          <view class="menu-title">&#25105;&#30340;&#20030;&#25253;</view>
          <view class="menu-desc">&#26597;&#30475;&#22788;&#29702;&#36827;&#24230;</view>
        </view>
        <text>&#8250;</text>
      </view>
      <view class="menu-item" @click="goAppealList">
        <view>
          <view class="menu-title">我的申诉</view>
          <view class="menu-desc">提交复核并查看处理结果</view>
        </view>
        <text>&#8250;</text>
      </view>
      <view :class="['menu-item', !isVerified ? 'locked' : '']" @click="goBrowseHistory">
        <view>
          <view class="menu-title">&#27983;&#35272;&#21382;&#21490;</view>
          <view class="menu-desc">&#26597;&#30475;&#26368;&#36817;&#25171;&#24320;&#30340;&#21830;&#21697;</view>
        </view>
        <text>&#8250;</text>
      </view>
      <view class="menu-item" @click="goStudentAuth">
        <view>
          <view class="menu-title">&#23398;&#29983;&#35748;&#35777;</view>
          <view class="menu-desc">{{ authText }}</view>
        </view>
        <text>&#8250;</text>
      </view>
    </view>

    <view v-if="isAdmin" class="menu-card">
      <view class="menu-item" @click="goAdmin">
        <view>
          <view class="menu-title">&#21518;&#21488;&#31649;&#29702;</view>
          <view class="menu-desc">&#23457;&#26680;&#12289;&#21830;&#21697;&#12289;&#39118;&#25511;&#19982;&#26085;&#24535;</view>
        </view>
        <text>&#8250;</text>
      </view>
    </view>

    <button class="logout-button" @click="logout">&#36864;&#20986;&#30331;&#24405;</button>
    <button v-if="canCancelAccount" class="cancel-account-button" :loading="canceling" @click="cancelMyAccount">注销账号</button>
  </view>
</template>

<script setup>
import { computed, ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import { cancelAccount, getCurrentUser } from '../../api/auth'
import { clearAuth, getUser, setUser } from '../../utils/auth'
import { BASE_URL } from '../../utils/request'

const user = ref(getUser() || {})
const canceling = ref(false)

const authMap = {
  UNVERIFIED: '\u672a\u8ba4\u8bc1',
  PENDING: '\u5ba1\u6838\u4e2d',
  VERIFIED: '\u5df2\u8ba4\u8bc1',
  REJECTED: '\u5df2\u9a73\u56de'
}

const authStatus = computed(() => user.value?.authStatus || user.value?.auth_status || '')
const accountStatus = computed(() => user.value?.accountStatus || user.value?.account_status || '')
const role = computed(() => String(user.value?.role || '').toUpperCase())
const authText = computed(() => authMap[authStatus.value] || '\u672a\u8ba4\u8bc1')
const isVerified = computed(() => authStatus.value === 'VERIFIED')
const authTone = computed(() => {
  if (authStatus.value === 'VERIFIED') return 'verified'
  if (authStatus.value === 'REJECTED') return 'rejected'
  if (authStatus.value === 'PENDING') return 'pending'
  return ''
})
const isAdmin = computed(() => role.value === 'ADMIN' || role.value === 'SUPER_ADMIN')
const canCancelAccount = computed(() => !isAdmin.value && accountStatus.value === 'NORMAL')
const identityText = computed(() => isAdmin.value ? '\u7ba1\u7406\u5458' : `\u5b66\u751f\u8ba4\u8bc1\uff1a${authText.value}`)

const avatarUrl = computed(() => {
  const url = user.value?.avatarUrl || user.value?.avatar_url || ''
  if (!url) return ''
  if (url.startsWith('http://') || url.startsWith('https://')) return url
  if (url.startsWith('/uploads/')) return BASE_URL + url
  return url
})

onShow(async () => {
  try {
    const current = await getCurrentUser()
    user.value = current
    setUser(current)
  } catch (error) {
    if (error.message) {
      uni.showToast({ title: error.message, icon: 'none' })
    }
  }
})

const goProfileEdit = () => uni.navigateTo({ url: '/pages/profile-edit/profile-edit' })
const requireVerified = (next) => {
  if (!isVerified.value) {
    uni.showToast({ title: '\u5b66\u751f\u8ba4\u8bc1\u901a\u8fc7\u540e\u624d\u80fd\u4f7f\u7528', icon: 'none' })
    return
  }
  next()
}
const goMyProducts = () => requireVerified(() => uni.navigateTo({ url: '/pages/my-products/my-products' }))
const goSoldOrders = () => requireVerified(() => uni.navigateTo({ url: '/pages/my-orders/my-orders?role=seller' }))
const goBoughtOrders = () => requireVerified(() => uni.navigateTo({ url: '/pages/my-orders/my-orders?role=buyer' }))
const goStudentAuth = () => uni.navigateTo({ url: '/pages/student-auth/student-auth' })
const goAdmin = () => uni.navigateTo({ url: '/pages/admin/admin' })
const goFavorites = () => requireVerified(() => uni.navigateTo({ url: '/pages/interaction/favorites' }))
const goReportList = () => requireVerified(() => uni.navigateTo({ url: '/pages/interaction/report-list' }))
const goAppealList = () => uni.navigateTo({ url: '/pages/interaction/appeal-list' })
const goBrowseHistory = () => requireVerified(() => uni.navigateTo({ url: '/pages/interaction/browse-history' }))

const confirmModal = (options) => new Promise((resolve) => {
  uni.showModal({
    ...options,
    success: (res) => resolve(Boolean(res.confirm))
  })
})

const cancelMyAccount = async () => {
  if (canceling.value) return
  const confirmed = await confirmModal({
    title: '注销账号',
    content: '注销后账号将退出平台，旧登录状态会失效。在重新登录并确认恢复前，无法继续使用交易功能。',
    confirmText: '确认注销',
    cancelText: '取消'
  })
  if (!confirmed) return
  canceling.value = true
  try {
    await cancelAccount()
    clearAuth()
    uni.showToast({ title: '账号已注销', icon: 'success' })
    setTimeout(() => uni.reLaunch({ url: '/pages/login/login' }), 500)
  } catch (error) {
    uni.showToast({ title: error.message || '注销失败', icon: 'none' })
  } finally {
    canceling.value = false
  }
}

const logout = () => {
  clearAuth()
  uni.reLaunch({ url: '/pages/login/login' })
}
</script>

<style scoped>
.page {
  min-height: 100vh;
  padding: 0 26rpx 42rpx;
  background: #f3f8f5;
  box-sizing: border-box;
}

.profile-hero {
  margin: 0 -26rpx;
  padding: 70rpx 34rpx 86rpx;
  border-radius: 0 0 42rpx 42rpx;
  background: linear-gradient(145deg, #23734f, #2f8b62);
  color: #fff;
}

.hero-eyebrow,
.hero-title,
.hero-copy {
  display: block;
}

.hero-eyebrow {
  color: rgba(255, 255, 255, .72);
  font-size: 20rpx;
  letter-spacing: 3rpx;
}

.hero-title {
  margin-top: 14rpx;
  font-size: 42rpx;
  line-height: 50rpx;
  font-weight: 700;
}

.hero-copy {
  margin-top: 14rpx;
  color: rgba(255, 255, 255, .78);
  font-size: 24rpx;
}

.profile-card {
  display: flex;
  align-items: center;
  min-height: 168rpx;
  padding: 34rpx;
  margin-top: -48rpx;
  border-radius: 28rpx;
  background: #fff;
  box-shadow: 0 14rpx 34rpx rgba(31, 106, 73, .06);
}

.avatar {
  display: flex;
  width: 112rpx;
  height: 112rpx;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  border-radius: 56rpx;
  background: #e8edf2;
  color: #8b98a7;
  font-size: 24rpx;
}

.profile-main {
  flex: 1;
  min-width: 0;
  margin-left: 24rpx;
}

.nickname {
  font-size: 36rpx;
  line-height: 44rpx;
  font-weight: 700;
  color: #1f2933;
}

.status {
  display: inline-flex;
  margin-top: 12rpx;
  padding: 7rpx 16rpx;
  border-radius: 999rpx;
  background: #f2f4f7;
  color: #667085;
  font-size: 24rpx;
}

.status.verified {
  background: #dcfce7;
  color: #16a34a;
}

.status.pending {
  background: #fff7ed;
  color: #c26a18;
}

.status.rejected {
  background: #fee2e2;
  color: #ef4444;
}

.edit-button {
  flex-shrink: 0;
  margin-left: 18rpx;
  padding: 0 26rpx;
  height: 64rpx;
  line-height: 64rpx;
  border-radius: 18rpx;
  background: #eef7f0;
  color: #17a84b;
  font-size: 26rpx;
}

.section {
  margin-top: 34rpx;
}

.section-title {
  margin-bottom: 18rpx;
  color: #20352b;
  font-size: 32rpx;
  font-weight: 700;
}

.quick-list {
  border-radius: 26rpx;
  background: #fff;
  box-shadow: 0 12rpx 32rpx rgba(31, 106, 73, .05);
  overflow: hidden;
}

.quick-card,
.menu-item {
  display: flex;
  align-items: center;
  min-height: 122rpx;
  padding: 22rpx 30rpx;
  border-bottom: 1rpx solid #eef0f3;
  box-sizing: border-box;
}

.quick-card:last-child,
.menu-item:last-child {
  border-bottom: 0;
}

.quick-card.locked,
.menu-item.locked {
  opacity: .62;
}

.quick-icon {
  display: flex;
  width: 68rpx;
  height: 68rpx;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  margin-right: 18rpx;
  border-radius: 22rpx;
  background: #e7f8ef;
  color: #17a84b;
  font-size: 28rpx;
  font-weight: 700;
}

.quick-icon.warm {
  background: #fff2d8;
  color: #a96500;
}

.quick-icon.cool {
  background: #e8f1ff;
  color: #2f6bcc;
}

.quick-copy {
  flex: 1;
  min-width: 0;
}

.quick-title,
.menu-title {
  color: #20352b;
  font-size: 30rpx;
  line-height: 38rpx;
  font-weight: 700;
}

.quick-desc,
.menu-desc {
  margin-top: 8rpx;
  color: #98a2b3;
  font-size: 24rpx;
}

.arrow,
.menu-item text {
  margin-left: auto;
  color: #b2bdca;
  font-size: 38rpx;
}

.menu-card {
  margin-top: 28rpx;
  border-radius: 26rpx;
  background: #fff;
  box-shadow: 0 12rpx 32rpx rgba(31, 106, 73, .05);
  overflow: hidden;
}

.menu-item {
  justify-content: space-between;
}

.logout-button {
  margin-top: 30rpx;
  height: 88rpx;
  line-height: 88rpx;
  border-radius: 24rpx;
  background: #fff;
  color: #ef4444;
  font-size: 30rpx;
  box-shadow: 0 12rpx 32rpx rgba(31, 106, 73, .05);
}

.cancel-account-button {
  margin-top: 30rpx;
  height: 88rpx;
  line-height: 88rpx;
  border-radius: 24rpx;
  background: #fff7ed;
  color: #c26a18;
  font-size: 30rpx;
  box-shadow: 0 12rpx 32rpx rgba(31, 106, 73, .05);
}
</style>
