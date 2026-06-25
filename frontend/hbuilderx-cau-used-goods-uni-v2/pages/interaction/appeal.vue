<template>
  <view class="page">
    <view class="notice">申诉用于对账号、商品、订单或举报处理结果申请管理员复核。请填写真实原因，必要时上传凭证图片。</view>

    <view class="card">
      <view class="field">
        <text class="field-label">申诉对象</text>
        <picker :range="targetLabels" :value="targetIndex" :disabled="targetLocked" @change="selectTargetType">
          <view class="picker-value" :class="{ locked: targetLocked }">{{ targetTypeLabel }}</view>
        </picker>
      </view>

      <view class="field">
        <text class="field-label">对象 ID</text>
        <input v-model="form.targetId" class="input" :class="{ locked: targetLocked }" :disabled="targetLocked" type="number" placeholder="请输入商品、订单、用户或举报 ID" />
      </view>

      <view class="field">
        <text class="field-label">申诉理由</text>
        <textarea v-model="form.reason" class="textarea" maxlength="500" placeholder="请说明为什么需要复核，最多 500 字" />
      </view>

      <view class="field">
        <text class="field-label">凭证图片（选填，最多 3 张）</text>
        <view class="images">
          <view v-for="src in form.images" :key="src" class="image-wrap">
            <image :src="src" mode="aspectFill" />
            <text class="remove" @click.stop="removeImage(src)">×</text>
          </view>
          <view v-if="form.images.length < 3" class="image-add" @click="chooseImage">+</view>
        </view>
      </view>

      <button class="btn btn-primary" :disabled="submitting" @click="submit">
        {{ submitting ? '提交中...' : '提交申诉' }}
      </button>
    </view>
  </view>
</template>

<script setup>
import { computed, reactive, ref } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import { tradeService } from '../../services/trade'
import { getUser } from '../../utils/auth'
import { showError, showSuccess } from '../../utils/navigation'

const targetTypes = [
  { value: 'PRODUCT', label: '商品' },
  { value: 'ORDER', label: '订单' },
  { value: 'USER', label: '账号' },
  { value: 'REPORT', label: '举报' }
]
const targetLabels = targetTypes.map((item) => item.label)
const form = reactive({ targetType: 'PRODUCT', targetId: '', reason: '', images: [] })
const submitting = ref(false)
const targetLocked = ref(false)

const targetIndex = computed(() => Math.max(0, targetTypes.findIndex((item) => item.value === form.targetType)))
const targetTypeLabel = computed(() => targetTypes[targetIndex.value]?.label || '商品')

onLoad((options) => {
  if (options.targetType) form.targetType = String(options.targetType).toUpperCase()
  if (options.targetId) form.targetId = options.targetId
  targetLocked.value = options.lockTarget === '1' || options.lockTarget === 1
  if (form.targetType === 'USER' && !form.targetId) fillCurrentUserId()
})

function selectTargetType(event) {
  if (targetLocked.value) return
  form.targetType = targetTypes[Number(event.detail.value)]?.value || 'PRODUCT'
  if (form.targetType === 'USER' && !form.targetId) fillCurrentUserId()
}

function fillCurrentUserId() {
  const user = getUser() || {}
  form.targetId = user.id || user.userId || user.user_id || ''
}

function chooseImage() {
  uni.chooseImage({
    count: 3 - form.images.length,
    success: ({ tempFilePaths }) => form.images.push(...tempFilePaths)
  })
}

function removeImage(src) {
  form.images = form.images.filter((item) => item !== src)
}

async function submit() {
  if (!form.targetId) {
    showError(new Error('请填写申诉对象 ID'))
    return
  }
  if (!form.reason.trim()) {
    showError(new Error('请填写申诉理由'))
    return
  }
  submitting.value = true
  try {
    await tradeService.createAppeal({
      targetType: form.targetType,
      targetId: form.targetId,
      reason: form.reason.trim(),
      images: form.images
    })
    showSuccess('申诉已提交')
    setTimeout(() => uni.redirectTo({ url: '/pages/interaction/appeal-list' }), 500)
  } catch (error) {
    showError(error)
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped lang="scss">
.page { min-height: 100vh; padding: 24rpx; background: #f6f8f5; box-sizing: border-box; }
.notice { margin-bottom: 20rpx; padding: 20rpx 24rpx; border-radius: 18rpx; background: #fff7e6; color: #8a5a12; font-size: 25rpx; line-height: 1.6; }
.card { padding: 26rpx 24rpx; border-radius: 20rpx; background: #fff; }
.field { margin-bottom: 28rpx; }
.field-label { display: block; margin-bottom: 14rpx; color: #425148; font-size: 27rpx; line-height: 1.5; }
.picker-value, .input { box-sizing: border-box; min-height: 78rpx; padding: 18rpx 20rpx; border-radius: 14rpx; background: #f7faf8; color: #27352f; font-size: 27rpx; line-height: 1.5; }
.picker-value.locked, .input.locked { color: #667085; background: #eef2f0; }
.textarea { box-sizing: border-box; width: 100%; min-height: 220rpx; padding: 18rpx 20rpx; border-radius: 14rpx; background: #f7faf8; color: #27352f; font-size: 27rpx; line-height: 1.6; }
.images { display: flex; gap: 16rpx; flex-wrap: wrap; }
.image-wrap { position: relative; }
.image-wrap image, .image-add { width: 144rpx; height: 144rpx; border-radius: 14rpx; }
.image-add { display: flex; align-items: center; justify-content: center; border: 1rpx dashed #b8c3bd; color: #91a098; background: #fbfcfb; font-size: 54rpx; }
.remove { position: absolute; top: -10rpx; right: -10rpx; display: flex; width: 36rpx; height: 36rpx; align-items: center; justify-content: center; border-radius: 50%; background: rgba(0, 0, 0, .62); color: #fff; font-size: 28rpx; line-height: 36rpx; }
.btn-primary { background: #23734f; color: #fff; }
button[disabled] { opacity: .56; }
</style>
