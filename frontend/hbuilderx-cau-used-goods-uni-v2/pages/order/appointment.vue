<template>
  <view class="page">
    <view v-if="product" class="product-card">
      <image v-if="product.image" class="product-image" :src="product.image" mode="aspectFill" />
      <view v-else class="product-image placeholder">商品</view>
      <view class="product-info">
        <text class="product-title">{{ product.title }}</text>
        <text class="product-price">¥{{ product.price }}</text>
        <text class="product-location">{{ product.meetLocation || '预约后协商面交地点' }}</text>
      </view>
    </view>

    <view class="notice">平台不提供线上支付和物流服务，请选择校园内公共区域完成面交。</view>

    <view class="card form-card">
      <view class="field">
        <text class="field-label">期望面交时间</text>
        <view class="datetime-row">
          <picker mode="date" :start="minDate" :value="form.meetDate" @change="changeDate">
            <view class="picker-box">{{ form.meetDate || '选择日期' }}</view>
          </picker>
          <picker mode="time" :value="form.meetClock" @change="changeTime">
            <view class="picker-box">{{ form.meetClock || '选择时间' }}</view>
          </picker>
        </view>
        <text class="field-tip">请选择当前时间之后的面交时间。</text>
      </view>

      <view class="field">
        <text class="field-label">面交地点</text>
        <input v-model="form.meetLocation" class="input appointment-input" placeholder="如：东区图书馆门口" />
      </view>

      <view class="field">
        <text class="field-label">备注（选填）</text>
        <textarea v-model="form.remark" class="textarea" placeholder="可填写时间补充或其他说明" maxlength="200" />
      </view>

      <button class="btn btn-primary submit-btn" :class="{ duplicate: hasActiveOrder }" :disabled="submitting" @click="submit">
        {{ hasActiveOrder ? '已预约，不能重复提交' : (submitting ? '正在提交...' : '确认提交预约') }}
      </button>
    </view>
  </view>
</template>

<script setup>
import { onLoad } from '@dcloudio/uni-app'
import { computed, reactive, ref } from 'vue'
import { tradeService } from '../../services/trade'
import { showError, showSuccess } from '../../utils/navigation'

const product = ref()
const submitting = ref(false)
const hasActiveOrder = ref(false)
const productId = ref('')
const form = reactive({ meetDate: '', meetClock: '', meetLocation: '', remark: '' })
const minDate = computed(() => formatDate(new Date()))

onLoad(async (options) => {
  productId.value = options.productId || 'p-1001'
  initDefaultTime()
  await checkActiveOrder()
  try {
    product.value = await tradeService.getProduct(productId.value)
    form.meetLocation = product.value.meetLocation
  } catch (error) {
    if (!hasActiveOrder.value) showError(error)
  }
})

async function checkActiveOrder() {
  try {
    const orders = await tradeService.getOrders('buyer', { pageSize: 100 })
    hasActiveOrder.value = (orders || []).some((order) => (
      String(order.productId || order.product?.id) === String(productId.value)
      && ['PENDING_CONFIRM', 'WAIT_MEET'].includes(order.status)
    ))
  } catch (error) {
    hasActiveOrder.value = false
  }
  return hasActiveOrder.value
}

function showDuplicateAppointment() {
  uni.showModal({
    title: '不能重复预约',
    content: '你已经预约过该商品，请前往“我的订单”查看当前预约进度。',
    showCancel: false,
    confirmText: '我知道了'
  })
}

function pad(value) {
  return String(value).padStart(2, '0')
}

function formatDate(date) {
  return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())}`
}

function formatClock(date) {
  return `${pad(date.getHours())}:${pad(date.getMinutes())}`
}

function initDefaultTime() {
  const next = new Date(Date.now() + 60 * 60 * 1000)
  form.meetDate = formatDate(next)
  form.meetClock = formatClock(next)
}

function changeDate(event) {
  form.meetDate = event.detail.value
}

function changeTime(event) {
  form.meetClock = event.detail.value
}

function selectedMeetDate() {
  if (!form.meetDate || !form.meetClock) return null
  const date = new Date(`${form.meetDate.replace(/-/g, '/')} ${form.meetClock}:00`)
  return Number.isNaN(date.getTime()) ? null : date
}

function normalizedMeetTime() {
  const date = selectedMeetDate()
  if (!date) return ''
  return `${formatDate(date)} ${formatClock(date)}:00`
}

async function submit() {
  if (hasActiveOrder.value || await checkActiveOrder()) {
    showDuplicateAppointment()
    return
  }
  if (!form.meetDate || !form.meetClock || !form.meetLocation) {
    showError(new Error('请选择面交时间并填写地点'))
    return
  }
  const date = selectedMeetDate()
  if (!date || date.getTime() <= Date.now()) {
    showError(new Error('面交时间必须晚于当前时间'))
    return
  }
  submitting.value = true
  try {
    const order = await tradeService.createAppointment({
      productId: productId.value,
      meetLocation: form.meetLocation,
      remark: form.remark,
      meetTime: normalizedMeetTime()
    })
    showSuccess('预约成功')
    setTimeout(() => uni.redirectTo({ url: `/pages/order/detail?id=${order.id}` }), 500)
  } catch (error) {
    const message = String(error?.message || '')
    if (message.includes('已预约过该商品')) {
      showDuplicateAppointment()
    } else {
      showError(error)
    }
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped lang="scss">
.page { min-height: 100vh; padding: 24rpx; background: #f6f7f9; box-sizing: border-box; }
.product-card { display: flex; gap: 22rpx; padding: 24rpx; border-radius: 26rpx; background: #fff; box-shadow: 0 10rpx 32rpx rgba(26, 61, 43, .06); }
.product-image { width: 220rpx; height: 180rpx; flex: 0 0 220rpx; border-radius: 22rpx; background: #edf2ef; }
.placeholder { display: flex; align-items: center; justify-content: center; color: #9aa5a1; }
.product-info { display: flex; flex: 1; min-width: 0; flex-direction: column; justify-content: center; }
.product-title { display: -webkit-box; overflow: hidden; color: #222; font-size: 34rpx; font-weight: 800; line-height: 1.35; -webkit-box-orient: vertical; -webkit-line-clamp: 2; }
.product-price { margin-top: 16rpx; color: #e36a3e; font-size: 40rpx; font-weight: 800; }
.product-location { margin-top: 10rpx; color: #8a9690; font-size: 24rpx; }
.notice { margin: 22rpx 0; padding: 20rpx 22rpx; border-radius: 18rpx; color: #71501e; background: #fff8e9; font-size: 25rpx; line-height: 1.6; }
.form-card { padding: 28rpx 24rpx; border-radius: 26rpx; background: #fff; }
.datetime-row { display: grid; grid-template-columns: 1fr 1fr; gap: 16rpx; }
.picker-box { min-height: 82rpx; padding: 0 22rpx; border: 2rpx solid #e6ece8; border-radius: 16rpx; background: #fbfcfb; color: #243129; font-size: 28rpx; line-height: 82rpx; box-sizing: border-box; }
.appointment-input { min-height: 82rpx; line-height: 82rpx; }
.field-tip { display: block; margin-top: 10rpx; color: #8a9690; font-size: 22rpx; line-height: 1.5; }
.submit-btn { margin-top: 8rpx; min-height: 84rpx; font-size: 30rpx; }
.submit-btn.duplicate { background: #8a9690; }
</style>
