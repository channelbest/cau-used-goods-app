<template>
  <view class="page">
    <view v-if="report" class="card">
      <view class="head">
        <text class="title">{{ report.reasonLabel || report.reason }}</text>
        <StatusBadge :label="status(report.status).label" :tone="status(report.status).tone" />
      </view>
      <view class="row"><text>举报对象</text><text>{{ report.targetTypeLabel || report.targetType }} · {{ report.targetId }}</text></view>
      <view class="row"><text>提交时间</text><text>{{ report.createdAt }}</text></view>
      <view class="section">
        <text class="section-title">举报说明</text>
        <text class="desc">{{ report.detail || '未填写补充说明' }}</text>
      </view>
      <view v-if="report.result" class="result">处理结果：{{ report.result }}</view>
      <button v-if="canAppealReport" class="appeal-button" @click="appealReport">申诉举报处理结果</button>
      <view v-else-if="reportAppealHint" class="appeal-hint">{{ reportAppealHint }}</view>
    </view>

    <view v-if="target" class="card">
      <text class="section-title">{{ targetTitle }}</text>

      <view v-if="report?.targetType === 'ORDER'" class="target-info">
        <view class="row"><text>订单编号</text><text>{{ target.orderNo || target.id }}</text></view>
        <view class="row"><text>商品名称</text><text>{{ target.product?.title || target.productTitleSnapshot }}</text></view>
        <view class="row"><text>预约时间</text><text>{{ target.meetTime || '-' }}</text></view>
        <view class="row"><text>面交地点</text><text>{{ target.meetLocation || '-' }}</text></view>
      </view>

      <view v-else-if="report?.targetType === 'PRODUCT'" class="product" @click="openProduct">
        <image v-if="target.image" class="cover" :src="target.image" mode="aspectFill" />
        <view v-else class="cover placeholder">商品</view>
        <view class="product-body">
          <text class="product-title">{{ target.title }}</text>
          <text class="price">¥{{ target.price }}</text>
          <text class="meta">{{ target.meetLocation || '预约后协商' }}</text>
        </view>
      </view>

      <view v-else-if="report?.targetType === 'USER'" class="user" @click="openUser">
        <image v-if="target.avatarUrl" class="avatar" :src="target.avatarUrl" mode="aspectFill" />
        <view v-else class="avatar text-avatar">{{ (target.nickname || '同').slice(0, 1) }}</view>
        <view>
          <view class="user-name">{{ target.nickname || 'CAU 同学' }}</view>
          <view class="meta">{{ target.authStatus === 'VERIFIED' ? '已认证' : '未认证' }}</view>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { computed, ref } from 'vue'
import { onLoad } from '@dcloudio/uni-app'
import StatusBadge from '../../components/StatusBadge.vue'
import { getPublicProfile } from '../../api/user'
import { tradeService } from '../../services/trade'
import { REPORT_STATUS } from '../../utils/constants'
import { BASE_URL } from '../../utils/request'
import { navigate, showError } from '../../utils/navigation'

const report = ref(null)
const target = ref(null)
const targetTitle = ref('关联详情')
const canAppealReport = computed(() => {
  const value = String(report.value?.status || '').toUpperCase()
  return value === 'APPROVED' || value === 'REJECTED'
})
// const reportAppealHint = computed(() => {
//   const value = String(report.value?.status || '').toUpperCase()
//   if (value === 'PENDING' || value === 'PROCESSING') return '举报处理完成后可申诉处理结果'
//   if (value === 'CLOSED') return '举报已关闭，暂无处理结果可申诉'
//   return ''
// })

function status(value) {
  return REPORT_STATUS[value] || { label: value, tone: 'muted' }
}

function normalizeImage(url) {
  if (!url) return ''
  return /^https?:\/\//.test(url) ? url : `${BASE_URL}${url}`
}

onLoad(async (options) => {
  try {
    const list = await tradeService.getReports()
    report.value = list.find((item) => String(item.id) === String(options.id))
    if (!report.value) {
      showError(new Error('举报记录不存在'))
      return
    }
    await loadTarget()
  } catch (error) {
    showError(error)
  }
})

async function loadTarget() {
  try {
    if (report.value.targetType === 'ORDER') {
      targetTitle.value = '订单详情'
      target.value = await tradeService.getOrder(report.value.targetId)
    } else if (report.value.targetType === 'PRODUCT') {
      targetTitle.value = '商品详情'
      target.value = await tradeService.getProduct(report.value.targetId)
    } else if (report.value.targetType === 'USER') {
      targetTitle.value = '用户详情'
      const user = await getPublicProfile(report.value.targetId)
      target.value = { ...user, avatarUrl: normalizeImage(user.avatarUrl) }
    }
  } catch (error) {
    target.value = null
  }
}

function openProduct() {
  navigate('/pages/detail/detail', { id: report.value.targetId })
}

function openUser() {
  navigate('/pages/user-profile/user-profile', { id: report.value.targetId })
}

function appealReport() {
  navigate('/pages/interaction/appeal', { targetType: 'REPORT', targetId: report.value.id })
}
</script>

<style scoped>
.page { min-height: 100vh; padding: 24rpx; background: #f6f7f9; box-sizing: border-box; }
.card { margin-bottom: 18rpx; padding: 26rpx; border-radius: 22rpx; background: #fff; box-shadow: 0 8rpx 28rpx rgba(23, 33, 43, .04); }
.head { display: flex; justify-content: space-between; gap: 18rpx; align-items: flex-start; }
.title { flex: 1; color: #243129; font-size: 34rpx; font-weight: 800; }
.row { display: flex; justify-content: space-between; gap: 24rpx; padding: 16rpx 0; border-bottom: 1rpx solid #edf1ee; color: #7a8580; font-size: 25rpx; }
.row text:last-child { flex: 1; color: #33413a; text-align: right; word-break: break-all; }
.section { margin-top: 18rpx; }
.section-title { display: block; margin-bottom: 14rpx; color: #243129; font-size: 29rpx; font-weight: 700; }
.desc { color: #59675f; font-size: 26rpx; line-height: 1.7; }
.result { margin-top: 20rpx; padding: 18rpx; border-radius: 14rpx; color: #2f6b4f; background: #edf6f1; font-size: 25rpx; line-height: 1.5; }
.appeal-button { margin-top: 22rpx; height: 76rpx; border-radius: 999rpx; background: #fff7e6; color: #a96500; font-size: 27rpx; line-height: 76rpx; }
.appeal-hint { margin-top: 22rpx; padding: 18rpx 22rpx; border-radius: 18rpx; background: #f8faf9; color: #8a948f; font-size: 25rpx; line-height: 1.5; text-align: center; }
.product, .user { display: flex; gap: 18rpx; align-items: center; }
.cover { width: 150rpx; height: 122rpx; flex: 0 0 150rpx; border-radius: 16rpx; background: #edf2ef; }
.placeholder { display: flex; align-items: center; justify-content: center; color: #9aa5a1; font-size: 23rpx; }
.product-body { flex: 1; min-width: 0; }
.product-title { display: block; overflow: hidden; color: #243129; font-size: 30rpx; font-weight: 700; text-overflow: ellipsis; white-space: nowrap; }
.price { display: block; margin-top: 8rpx; color: #e36a3e; font-size: 32rpx; font-weight: 800; }
.meta { margin-top: 8rpx; color: #87918c; font-size: 24rpx; }
.avatar { display: flex; width: 92rpx; height: 92rpx; align-items: center; justify-content: center; border-radius: 50%; background: #e8ecef; color: #fff; font-size: 34rpx; font-weight: 800; }
.text-avatar { background: linear-gradient(135deg, #f3b34c, #f47b45); }
.user-name { color: #243129; font-size: 31rpx; font-weight: 800; }
</style>
