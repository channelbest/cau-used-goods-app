<template>
  <view class="page">
    <view class="toolbar">
      <view>
        <view class="title">我的申诉</view>
        <view class="subtitle">查看申诉进度和处理结果</view>
      </view>
    </view>

    <view v-if="appeals.length" class="list">
      <view v-for="appeal in appeals" :key="appeal.id" class="card appeal" @click="openDetail(appeal)">
        <view class="dot" :class="{ done: isDone(appeal.status) }" />
        <view class="appeal-body">
          <view class="appeal-head">
            <text class="appeal-title">{{ appeal.reason }}</text>
            <StatusBadge :label="status(appeal.status).label" :tone="status(appeal.status).tone" />
          </view>
          <text class="appeal-meta">{{ appealMeta(appeal) }}</text>
          <view v-if="appeal.result" class="result">处理结果：{{ appeal.result }}</view>
        </view>
      </view>
    </view>

    <EmptyState v-else title="暂无申诉记录" detail="从相关账号、商品或订单问题入口提交申诉后，处理进度会显示在这里" />
  </view>
</template>

<script setup>
import { onShow } from '@dcloudio/uni-app'
import { ref } from 'vue'
import EmptyState from '../../components/EmptyState.vue'
import StatusBadge from '../../components/StatusBadge.vue'
import { tradeService } from '../../services/trade'
import { APPEAL_STATUS } from '../../utils/constants'
import { navigate, showError } from '../../utils/navigation'

const appeals = ref([])

onShow(load)

async function load() {
  try {
    appeals.value = (await tradeService.getAppeals()).items || []
  } catch (error) {
    showError(error)
  }
}

function status(value) {
  return APPEAL_STATUS[value] || { label: value, tone: 'muted' }
}

function isDone(value) {
  return ['APPROVED', 'REJECTED', 'CLOSED'].includes(value)
}

function targetText(value) {
  return {
    PRODUCT: '商品',
    USER: '用户',
    ORDER: '订单',
    REPORT: '举报'
  }[String(value || '').toUpperCase()] || ''
}

function appealMeta(appeal = {}) {
  const parts = []
  const target = targetText(appeal.targetType)
  if (target) parts.push(appeal.targetId ? `${target} #${appeal.targetId}` : target)
  else if (appeal.targetId) parts.push(`#${appeal.targetId}`)
  if (appeal.createdAt) parts.push(appeal.createdAt)
  return parts.join(' · ')
}

function openDetail(appeal) {
  navigate('/pages/interaction/appeal-detail', { id: appeal.id })
}

</script>

<style scoped lang="scss">
.page { min-height: 100vh; padding: 24rpx; background: #f6f7f9; box-sizing: border-box; }
.toolbar { display: flex; align-items: center; justify-content: space-between; gap: 18rpx; margin-bottom: 22rpx; padding: 24rpx; border-radius: 20rpx; background: #fff; }
.title { color: #243129; font-size: 34rpx; font-weight: 800; }
.subtitle { margin-top: 8rpx; color: #8d9892; font-size: 24rpx; }
.list { display: flex; flex-direction: column; gap: 18rpx; }
.appeal { display: flex; gap: 14rpx; padding: 24rpx; border-radius: 18rpx; background: #fff; }
.dot { width: 16rpx; height: 16rpx; margin-top: 10rpx; flex: 0 0 16rpx; border-radius: 50%; background: #f2a23a; }
.dot.done { background: #ccd4d0; }
.appeal-body { display: flex; flex: 1; min-width: 0; flex-direction: column; gap: 12rpx; }
.appeal-head { display: flex; align-items: flex-start; justify-content: space-between; gap: 14rpx; }
.appeal-title { flex: 1; min-width: 0; display: -webkit-box; overflow: hidden; color: #243129; font-size: 29rpx; font-weight: 700; line-height: 1.4; -webkit-box-orient: vertical; -webkit-line-clamp: 2; }
.appeal-meta { color: #98a39d; font-size: 22rpx; line-height: 1.5; }
.result { padding: 16rpx; border-radius: 12rpx; color: #2f6b4f; background: #edf6f1; font-size: 24rpx; line-height: 1.5; }
</style>
