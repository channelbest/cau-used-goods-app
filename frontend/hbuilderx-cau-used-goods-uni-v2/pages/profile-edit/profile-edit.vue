<template>
  <view class="page">
    <view class="form-card">
      <button class="profile-row avatar-picker" :open-type="wechatAvatarOpenType" @chooseavatar="useWechatAvatar" @click="chooseAvatarFile">
        <text class="row-label">头像</text>
        <view class="row-value">
          <image v-if="avatarUrl" class="avatar" :src="avatarUrl" mode="aspectFill" />
          <view v-else class="avatar placeholder">头像</view>
          <text class="arrow">›</text>
        </view>
      </button>

      <view class="profile-row">
        <text class="row-label">昵称</text>
        <input class="row-input" v-model="nickname" type="nickname" placeholder="请输入昵称" />
      </view>

      <view class="profile-row">
        <text class="row-label">手机号</text>
        <input class="row-input" v-model="phone" type="number" maxlength="11" placeholder="请输入手机号" />
      </view>
    </view>

    <button class="primary-button" :loading="loading" @click="saveProfile">保存资料</button>
  </view>
</template>

<script setup>
import { ref } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { getUser, setUser } from '../../utils/auth'
import { updateProfile, uploadAvatar } from '../../api/auth'
import { BASE_URL } from '../../utils/request'

const nickname = ref('')
const phone = ref('')
const avatarUrl = ref('')
const loading = ref(false)
const selectedAvatarPath = ref('')
const canUseWechatAvatar = typeof uni.canIUse === 'function' && uni.canIUse('button.open-type.chooseAvatar')
const wechatAvatarOpenType = canUseWechatAvatar ? 'chooseAvatar' : ''

const normalizeAvatar = (url) => {
  if (!url) return ''
  if (url.startsWith('http://') || url.startsWith('https://')) return url
  if (url.startsWith('/uploads/')) return BASE_URL + url
  return url
}

onLoad(() => {
  const user = getUser() || {}
  nickname.value = user.nickname || ''
  phone.value = user.phone || ''
  avatarUrl.value = normalizeAvatar(user.avatarUrl)
})

const useWechatAvatar = (event) => {
  const filePath = event?.detail?.avatarUrl || ''
  if (!filePath) return
  selectedAvatarPath.value = filePath
  avatarUrl.value = filePath
}

const chooseAvatarFile = () => {
  if (canUseWechatAvatar) return
  uni.chooseImage({
    count: 1,
    sizeType: ['compressed'],
    sourceType: ['album', 'camera'],
    success: (res) => {
      const filePath = res.tempFilePaths?.[0] || ''
      if (!filePath) return
      selectedAvatarPath.value = filePath
      avatarUrl.value = filePath
    }
  })
}

const saveProfile = async () => {
  if (loading.value) return

  const nextNickname = nickname.value.trim()
  const nextPhone = phone.value.trim()
  if (!nextNickname && !nextPhone && !selectedAvatarPath.value) {
    uni.showToast({ title: '请填写或选择要保存的资料', icon: 'none' })
    return
  }
  if (nextPhone && !/^1[3-9]\d{9}$/.test(nextPhone)) {
    uni.showToast({ title: '手机号格式不正确', icon: 'none' })
    return
  }

  const payload = {}
  if (nextNickname) payload.nickname = nextNickname
  if (nextPhone) payload.phone = nextPhone

  loading.value = true
  try {
    let user = null
    if (selectedAvatarPath.value) {
      const data = await uploadAvatar(selectedAvatarPath.value)
      user = data.user || null
      selectedAvatarPath.value = ''
    }
    if (Object.keys(payload).length) {
      user = await updateProfile(payload)
    }
    if (user) {
      setUser(user)
      avatarUrl.value = normalizeAvatar(user.avatarUrl)
      nickname.value = user.nickname || nextNickname
      phone.value = user.phone || nextPhone
    }
    uni.showToast({ title: '资料已保存', icon: 'success' })
    setTimeout(() => uni.navigateBack(), 600)
  } catch (error) {
    uni.showToast({ title: error.message || '保存失败', icon: 'none' })
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.page {
  min-height: 100vh;
  padding: 24rpx 24rpx 48rpx;
  background: #f5f6f8;
  box-sizing: border-box;
}

.form-card {
  overflow: hidden;
  border-radius: 18rpx;
  background: #ffffff;
}

.profile-row {
  display: flex;
  width: 100%;
  min-height: 112rpx;
  align-items: center;
  justify-content: space-between;
  margin: 0;
  padding: 0 28rpx;
  border-radius: 0;
  background: #ffffff;
  box-sizing: border-box;
}

.profile-row + .profile-row {
  border-top: 1rpx solid #eef0f3;
}

.avatar-picker {
  height: 144rpx;
  line-height: normal;
  text-align: left;
}

.avatar-picker::after {
  border: none;
}

.row-label {
  flex-shrink: 0;
  color: #1f2933;
  font-size: 30rpx;
  font-weight: 600;
}

.row-value {
  display: flex;
  align-items: center;
  gap: 18rpx;
}

.avatar {
  display: flex;
  width: 92rpx;
  height: 92rpx;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  background: #dce3ea;
  color: #98a2b3;
  font-size: 24rpx;
}

.arrow {
  color: #c0c7d0;
  font-size: 42rpx;
  line-height: 1;
}

.row-input {
  flex: 1;
  height: 112rpx;
  padding-left: 32rpx;
  color: #1f2933;
  font-size: 30rpx;
  text-align: right;
  box-sizing: border-box;
}

.primary-button {
  margin-top: 44rpx;
  height: 88rpx;
  line-height: 88rpx;
  border-radius: 18rpx;
  background: #17a84b;
  color: #ffffff;
  font-size: 30rpx;
  font-weight: 700;
}
</style>
