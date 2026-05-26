<template>
  <view class="page">
    <view class="title">Student Verification</view>

    <view class="form-item">
      <text class="label">Real Name</text>
      <input class="input" v-model="form.realName" placeholder="Enter real name" />
    </view>

    <view class="form-item">
      <text class="label">Student ID</text>
      <input class="input" v-model="form.studentId" placeholder="Enter student ID" />
    </view>

    <view class="form-item">
      <text class="label">College</text>
      <input class="input" v-model="form.college" placeholder="Enter college" />
    </view>

    <button class="submit-button" :loading="loading" @click="handleSubmit">
      Submit Verification
    </button>

    <button class="secondary-button" @click="goHome">
      Preview Home
    </button>
  </view>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { submitStudentVerification } from '../../api/auth'

const loading = ref(false)
const form = reactive({
  realName: '',
  studentId: '',
  college: ''
})

const validateForm = () => {
  if (!form.realName.trim()) return 'Enter real name'
  if (!form.studentId.trim()) return 'Enter student ID'
  if (!form.college.trim()) return 'Enter college'
  return ''
}

const goHome = () => {
  uni.navigateTo({
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
      title: 'Submit success',
      icon: 'success'
    })
    goHome()
  } catch (error) {
    uni.showToast({
      title: error.message || 'Submit failed',
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

.secondary-button {
  margin-top: 24rpx;
  height: 88rpx;
  line-height: 88rpx;
  border-radius: 12rpx;
  background: #ffffff;
  color: #374151;
  font-size: 30rpx;
}
</style>
