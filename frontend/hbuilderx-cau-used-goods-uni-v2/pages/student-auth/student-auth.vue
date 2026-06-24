<template>
  <view class="page">
    <view class="hero">
      <text class="eyebrow">CAU CAMPUS MARKET</text>
      <text class="hero-title">学生认证</text>
      <text class="hero-copy">{{ heroCopy }}</text>
    </view>

    <view v-if="authStatus === 'PENDING'" class="state-card">
      <view class="state-icon pending">审</view>
      <view class="state-title">认证已提交，等待管理员审核</view>
      <view class="state-desc">你的认证资料已经进入审核队列，审核通过后就可以发布商品、收藏、举报和查看交易记录。</view>
      <view class="info-list">
        <view class="info-row"><text>姓名</text><text>{{ form.realName || '-' }}</text></view>
        <view class="info-row"><text>学号</text><text>{{ form.studentId || '-' }}</text></view>
        <view class="info-row"><text>学院</text><text>{{ form.college || '-' }}</text></view>
        <view class="info-row"><text>当前状态</text><text class="pending-text">审核中</text></view>
      </view>
      <button class="secondary-button" @click="goHome">返回首页</button>
    </view>

    <view v-else-if="authStatus === 'VERIFIED'" class="state-card">
      <view class="state-icon success">✓</view>
      <view class="state-title">学生认证已通过</view>
      <view class="state-desc">你现在可以正常发布闲置、预约交易、收藏商品和提交举报。</view>
      <view class="info-list">
        <view class="info-row"><text>姓名</text><text>{{ form.realName || '-' }}</text></view>
        <view class="info-row"><text>学号</text><text>{{ form.studentId || '-' }}</text></view>
        <view class="info-row"><text>学院</text><text>{{ form.college || '-' }}</text></view>
      </view>
      <button class="submit-button" @click="goHome">返回首页</button>
    </view>

    <view v-else class="form-card">
      <view class="status-pill" :class="{ rejected: authStatus === 'REJECTED' }">{{ statusText }}</view>
      <view v-if="authStatus === 'REJECTED'" class="reject-note">上次认证未通过，请核对信息后重新提交。</view>

      <view class="form-item">
        <text class="label">姓名</text>
        <input class="input" v-model="form.realName" placeholder="请输入姓名" />
      </view>

      <view class="form-item">
        <text class="label">学号</text>
        <input class="input" v-model="form.studentId" type="number" placeholder="请输入学号" />
      </view>

      <view class="form-item">
        <text class="label">学院</text>
        <input class="input" v-model="form.college" placeholder="请输入学院" />
      </view>

      <button class="submit-button" :loading="loading" @click="handleSubmit">
        提交认证
      </button>

      <button class="secondary-button" @click="goHome">
        返回首页
      </button>
    </view>
  </view>
</template>

<script setup>
import { computed, reactive, ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import { getCurrentUser, getStudentVerification, submitStudentVerification } from '../../api/auth'
import { setUser } from '../../utils/auth'

const loading = ref(false)
const authStatus = ref('UNVERIFIED')
const form = reactive({
  realName: '',
  studentId: '',
  college: ''
})

const statusMap = {
  UNVERIFIED: '未认证',
  PENDING: '审核中',
  VERIFIED: '已认证',
  REJECTED: '已驳回'
}

const statusText = computed(() => statusMap[authStatus.value] || '未认证')
const heroCopy = computed(() => {
  if (authStatus.value === 'PENDING') return '资料已提交，请等待管理员审核。'
  if (authStatus.value === 'VERIFIED') return '你已完成校园身份认证。'
  if (authStatus.value === 'REJECTED') return '认证未通过，请重新核对并提交。'
  return '完成认证后即可发布、收藏和参与交易。'
})

onShow(async () => {
  try {
    const user = await getCurrentUser()
    setUser(user)
    authStatus.value = user.authStatus || 'UNVERIFIED'
    const verification = await getStudentVerification()
    form.realName = verification.realName || ''
    form.studentId = verification.studentId || ''
    form.college = verification.college || ''
  } catch (error) {
    if (error.message) {
      uni.showToast({ title: error.message, icon: 'none' })
    }
  }
})

const validateForm = () => {
  if (!/^[\u4e00-\u9fa5]{2,20}$/.test(form.realName.trim())) return '姓名需填写2到20个汉字'
  if (!/^\d{6,20}$/.test(form.studentId.trim())) return '学号需填写6到20位数字'
  if (!/^[\u4e00-\u9fa5]{2,30}$/.test(form.college.trim())) return '学院需填写2到30个汉字'
  return ''
}

const goHome = () => {
  uni.switchTab({
    url: '/pages/home/home'
  })
}

const handleSubmit = async () => {
  const message = validateForm()
  if (message) {
    uni.showToast({
      title: message,
      icon: 'none'
    })
    return
  }

  if (loading.value) return

  loading.value = true
  try {
    await submitStudentVerification({
      studentId: form.studentId.trim(),
      realName: form.realName.trim(),
      college: form.college.trim()
    })

    uni.showToast({
      title: '提交成功',
      icon: 'success'
    })
    const user = await getCurrentUser()
    setUser(user)
    goHome()
  } catch (error) {
    uni.showToast({
      title: error.message || '提交失败',
      icon: 'none'
    })
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.page {
  min-height: 100vh;
  padding-bottom: 42rpx;
  background: #f5f8f6;
  box-sizing: border-box;
}

.hero {
  padding: 72rpx 40rpx 78rpx;
  border-radius: 0 0 42rpx 42rpx;
  background: linear-gradient(145deg, #23734f, #2f8b62);
  color: #fff;
}

.eyebrow,
.hero-title,
.hero-copy {
  display: block;
}

.eyebrow {
  color: rgba(255, 255, 255, .72);
  font-size: 20rpx;
  letter-spacing: 3rpx;
}

.hero-title {
  margin-top: 16rpx;
  font-size: 42rpx;
  font-weight: 700;
}

.hero-copy {
  margin-top: 14rpx;
  color: rgba(255, 255, 255, .78);
  font-size: 26rpx;
}

.form-card,
.state-card {
  margin: -42rpx 30rpx 0;
  padding: 34rpx;
  border-radius: 26rpx;
  background: #fff;
  box-shadow: 0 16rpx 38rpx rgba(35, 115, 79, .1);
}

.state-card {
  text-align: center;
}

.state-icon {
  display: flex;
  width: 96rpx;
  height: 96rpx;
  margin: 4rpx auto 24rpx;
  align-items: center;
  justify-content: center;
  border-radius: 32rpx;
  font-size: 36rpx;
  font-weight: 700;
}

.state-icon.pending {
  background: #fff7ed;
  color: #c26a18;
}

.state-icon.success {
  background: #e8f8ef;
  color: #17a84b;
}

.state-title {
  color: #1f2933;
  font-size: 34rpx;
  font-weight: 700;
}

.state-desc {
  margin-top: 16rpx;
  color: #667085;
  font-size: 26rpx;
  line-height: 1.65;
}

.info-list {
  margin-top: 30rpx;
  overflow: hidden;
  border-radius: 18rpx;
  background: #f6f8f7;
  text-align: left;
}

.info-row {
  display: flex;
  justify-content: space-between;
  gap: 22rpx;
  padding: 22rpx 24rpx;
  border-bottom: 1rpx solid #edf1ef;
  color: #667085;
  font-size: 26rpx;
}

.info-row:last-child {
  border-bottom: 0;
}

.info-row text:last-child {
  color: #1f2933;
  font-weight: 600;
  text-align: right;
}

.info-row .pending-text {
  color: #c26a18;
}

.status-pill {
  display: inline-flex;
  margin-bottom: 26rpx;
  padding: 8rpx 18rpx;
  border-radius: 999rpx;
  background: #eef7f0;
  color: #17a84b;
  font-size: 24rpx;
}

.status-pill.rejected {
  background: #fee2e2;
  color: #ef4444;
}

.reject-note {
  margin: -8rpx 0 26rpx;
  color: #ef4444;
  font-size: 24rpx;
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
  border-radius: 18rpx;
  background: #f6f8f7;
  font-size: 28rpx;
}

.submit-button {
  margin-top: 42rpx;
  height: 88rpx;
  line-height: 88rpx;
  border-radius: 999rpx;
  background: #23734f;
  color: #ffffff;
  font-size: 32rpx;
}

.secondary-button {
  margin-top: 24rpx;
  height: 88rpx;
  line-height: 88rpx;
  border-radius: 999rpx;
  background: #f6f8f7;
  color: #374151;
  font-size: 30rpx;
}
</style>
