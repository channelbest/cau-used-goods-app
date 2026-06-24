<template>
  <view class="page">
    <view class="page-title">风险处理</view>

    <view class="tab-row">
      <view :class="['tab', activeMode === 'REPORT' ? 'active' : '']" @click="activeMode = 'REPORT'">
        举报处理
        <text>{{ pendingReports }}</text>
      </view>
      <view :class="['tab', activeMode === 'APPEAL' ? 'active' : '']" @click="activeMode = 'APPEAL'">
        申诉处理
        <text>{{ pendingAppeals }}</text>
      </view>
    </view>

    <view class="summary">
      <view>
        <view :class="['summary-number', pendingCount ? 'danger' : '']">{{ pendingCount }}</view>
        <view class="summary-label">{{ activeMode === 'REPORT' ? '待处理举报' : '待处理申诉' }}</view>
      </view>
      <view :class="['summary-status', pendingCount ? 'warning' : 'safe']">
        {{ pendingCount ? '需要处理' : '暂无风险' }}
      </view>
    </view>

    <view class="filter-row">
      <view
        v-for="item in currentFilters"
        :key="item.value"
        :class="['filter-chip', activeTarget === item.value ? 'active' : '']"
        @click="activeTarget = item.value"
      >
        {{ item.label }} {{ countByTarget(item.value) }}
      </view>
    </view>


    <view v-if="filteredItems.length === 0" class="empty">{{ emptyText }}</view>

    <view v-for="item in filteredItems" :key="item.id" class="case-card" @click="goRiskDetail(item)">
      <view class="card-head">
        <view class="card-main">
          <view class="case-title">{{ itemTitle(item) }}</view>
          <view class="case-sub">{{ targetText(item.targetType) }} #{{ item.targetId }}</view>
        </view>
        <view class="badge-row">
          <view class="status-badge" :style="statusBadgeStyle(item.status)">{{ statusText(item.status) }}</view>
          <view v-if="isUserAppeal(item)" class="risk-badge">{{ userAppealRiskText(item) }}</view>
        </view>
      </view>

      <view class="case-desc">{{ itemDescription(item) }}</view>

      <view class="meta-row">
        <text>{{ activeMode === 'REPORT' ? '举报人' : '申诉人' }}：{{ actorName(item) }}</text>
        <text>{{ shortTime(item.createTime) }}</text>
      </view>
    </view>

    <view v-if="reasonModal.visible" class="modal-mask" @click="closeReasonModal">
      <view class="reason-sheet" @click.stop>
        <view class="sheet-head">
          <view>
            <view class="sheet-title">{{ reasonModalTitle }}</view>
            <view class="sheet-subtitle">选择一个原因，也可以补充更具体的说明</view>
          </view>
          <text class="sheet-close" @click="closeReasonModal">×</text>
        </view>

        <view class="reason-list">
          <view
            v-for="reason in currentReasonOptions"
            :key="reason"
            :class="['reason-chip', reasonModal.reason === reason ? 'active' : '']"
            @click="reasonModal.reason = reason"
          >
            {{ reason }}
          </view>
        </view>

        <textarea
          v-model="reasonModal.note"
          class="reason-input"
          maxlength="200"
          placeholder="补充说明，可不填"
          placeholder-class="reason-placeholder"
        />

        <view class="sheet-actions">
          <button class="sheet-button cancel" @click="closeReasonModal">取消</button>
          <button class="sheet-button confirm" :loading="submittingReason" @click="submitReasonAction">确认处理</button>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { computed, reactive, ref, watch } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import {
  getAdminAppeals,
  getAdminReports,
  handleAdminAppeal,
  handleAdminReport,
  markAdminAppealProcessing,
  markAdminReportProcessing,
  getAdminProductById
} from '../../api/admin'
import { normalizeImage } from '../../utils/product-format'
import { displayRelatedUserName } from '../../utils/user-format'

const activeMode = ref('REPORT')
const activeTarget = ref('ALL')
const openedId = ref(null)
const reports = ref([])
const appeals = ref([])
const productMap = ref({})
const submittingReason = ref(false)
const reasonModal = reactive({
  visible: false,
  id: 0,
  status: '',
  mode: 'REPORT',
  reason: '',
  note: ''
})

const reportFilters = [
  { label: '全部', value: 'ALL' },
  { label: '商品', value: 'PRODUCT' },
  { label: '用户', value: 'USER' },
  { label: '订单', value: 'ORDER' }
]

const appealFilters = [
  ...reportFilters,
  { label: '举报', value: 'REPORT' }
]

const reasonOptions = {
  REPORT_REJECTED: [
    '证据不足，无法认定违规',
    '举报内容与对象不符',
    '未发现明显违规行为',
    '重复举报或恶意举报',
    '商品信息已修改，风险已解除'
  ],
  REPORT_CLOSED: [
    '重复举报，合并处理',
    '相关对象已通过其他方式处理',
    '举报人撤回或无法继续核实',
    '不属于平台举报受理范围',
    '已转人工跟进，暂时关闭'
  ],
  APPEAL_REJECTED: [
    '申诉材料不足',
    '原举报处理结论无误',
    '未提供有效证明',
    '申诉理由与处理结果无关',
    '存在重复申诉'
  ],
  APPEAL_CLOSED: [
    '重复申诉，合并处理',
    '申诉人撤回或无法联系',
    '关联举报已关闭',
    '已线下处理，关闭申诉',
    '不属于申诉受理范围'
  ]
}

const currentItems = computed(() => activeMode.value === 'REPORT' ? reports.value : appeals.value)
const currentFilters = computed(() => activeMode.value === 'REPORT' ? reportFilters : appealFilters)

const pendingReports = computed(() => reports.value.filter((item) => ['PENDING', 'PROCESSING'].includes(item.status)).length)
const pendingAppeals = computed(() => appeals.value.filter((item) => ['PENDING', 'PROCESSING'].includes(item.status)).length)
const pendingCount = computed(() => activeMode.value === 'REPORT' ? pendingReports.value : pendingAppeals.value)

const filteredItems = computed(() => {
  if (activeTarget.value === 'ALL') return currentItems.value
  return currentItems.value.filter((item) => item.targetType === activeTarget.value)
})

const emptyText = computed(() => {
  const type = activeTarget.value === 'ALL' ? '全部' : targetText(activeTarget.value)
  const noun = activeMode.value === 'REPORT' ? '举报' : '申诉'
  return `当前没有${type}${noun}记录`
})

const reasonModalTitle = computed(() => {
  const action = reasonModal.status === 'REJECTED' ? '驳回' : '关闭'
  const noun = reasonModal.mode === 'REPORT' ? '举报' : '申诉'
  return `${action}${noun}`
})

const currentReasonOptions = computed(() => {
  const key = `${reasonModal.mode}_${reasonModal.status}`
  return reasonOptions[key] || []
})

watch(activeMode, () => {
  activeTarget.value = 'ALL'
  openedId.value = null
})

const load = async () => {
  try {
    const [reportResult, appealResult] = await Promise.all([
      getAdminReports(),
      getAdminAppeals()
    ])
    reports.value = reportResult?.items || []
    appeals.value = appealResult?.items || []
    loadTargetProducts([...reports.value, ...appeals.value])
  } catch (error) {
    uni.showToast({ title: error.message || '加载失败', icon: 'none' })
  }
}

onShow(load)

const countByTarget = (targetType) => {
  if (targetType === 'ALL') return currentItems.value.length
  return currentItems.value.filter((item) => item.targetType === targetType).length
}

const selectTarget = (targetType) => {
  activeTarget.value = targetType
  openedId.value = null
}

const toggleOpen = (id) => {
  openedId.value = openedId.value === id ? null : id
}

const goRiskDetail = (item) => {
  if (!item?.id) return
  uni.navigateTo({ url: `/pages/admin-risk-detail/admin-risk-detail?mode=${activeMode.value}&id=${item.id}` })
}

const loadTargetProducts = async (items) => {
  const productIds = [...new Set((items || [])
    .filter((item) => item.targetType === 'PRODUCT' && item.targetId)
    .map((item) => item.targetId))]
    .filter((id) => !productMap.value[id])

  if (!productIds.length) return

  const entries = await Promise.all(productIds.map(async (id) => {
    try {
      const product = await getAdminProductById(id)
      return [id, product]
    } catch (error) {
      return [id, null]
    }
  }))

  const next = { ...productMap.value }
  entries.forEach(([id, product]) => {
    if (product) next[id] = product
  })
  productMap.value = next
}

const targetText = (targetType) => {
  const map = {
    PRODUCT: '商品',
    USER: '用户',
    ORDER: '订单',
    REPORT: '举报'
  }
  return map[targetType] || targetType || '对象'
}

const reasonText = (reasonType) => {
  const map = {
    FAKE_PRODUCT: '虚假或违规商品',
    INAPPROPRIATE_CONTENT: '不当内容',
    SCAM: '欺诈风险',
    TRADE_DISPUTE: '交易纠纷',
    FAKE: '虚假信息',
    FRAUD: '疑似诈骗',
    PROHIBITED: '违规商品',
    INAPPROPRIATE: '不当内容',
    HARASSMENT: '骚扰行为',
    OTHER: '其他原因'
  }
  return map[reasonType] || reasonType || '举报'
}

const statusText = (status) => {
  status = normalizeStatus(status)
  const map = {
    PENDING: '待处理',
    PROCESSING: '处理中',
    RESOLVED: '已处理',
    APPROVED: '已通过',
    REJECTED: '已驳回',
    CLOSED: '已关闭'
  }
  return map[status] || status || '未知'
}

const normalizeStatus = (status) => String(status || '').trim().toUpperCase()

const statusTone = (status) => {
  const value = normalizeStatus(status)
  if (value === 'APPROVED' || value === 'RESOLVED') return 'success'
  if (value === 'PENDING') return 'warning'
  if (value === 'PROCESSING') return 'primary'
  if (value === 'REJECTED') return 'danger'
  if (value === 'CLOSED') return 'muted'
  return 'muted'
}

const statusStyleMap = {
  warning: { background: '#fff7e6', color: '#b66a00' },
  primary: { background: '#eff6ff', color: '#2563eb' },
  success: { background: '#dcfce7', color: '#16a34a' },
  danger: { background: '#fee2e2', color: '#ef4444' },
  muted: { background: '#eef2f6', color: '#667085' }
}

const statusBadgeStyle = (status) => statusStyleMap[statusTone(status)] || statusStyleMap.muted

const isUserAppeal = (item) => activeMode.value === 'APPEAL' && item?.targetType === 'USER'

const userAppealRiskText = (item) => {
  const status = normalizeStatus(item?.status)
  if (status === 'PENDING' || status === 'PROCESSING') return '需超管'
  return '账号解封'
}

const itemTitle = (item) => {
  if (activeMode.value === 'REPORT') return reasonText(item.reasonType)
  return item.reason || '申诉'
}

const itemDescription = (item) => {
  return item.description || item.reason || '暂无补充说明'
}

const targetProduct = (item) => {
  if (item?.targetType !== 'PRODUCT') return null
  return productMap.value[item.targetId] || null
}

const productCover = (product) => {
  return normalizeImage(product?.images?.[0] || '')
}

const productStatusText = (status) => {
  const map = { ON_SALE: '在售', OFF_SHELF: '已下架', LOCKED: '交易锁定', SOLD: '已售出', DELETED: '已删除' }
  return map[status] || status || '未知'
}

const goProduct = (id) => {
  if (!id) return
  uni.navigateTo({ url: `/pages/detail/detail?id=${id}&adminView=1&readonly=1` })
}

const actorName = (item) => {
  if (activeMode.value === 'REPORT') return displayRelatedUserName(item, 'reporter', `用户${item.reporterId}`)
  return displayRelatedUserName(item, 'appellant', `用户${item.appellantId}`)
}

const shortTime = (value) => {
  if (!value) return ''
  return String(value).replace('T', ' ').slice(0, 16)
}

const canHandle = (status) => {
  return ['PENDING', 'PROCESSING'].includes(status)
}

const previewImage = (current, urls) => {
  uni.previewImage({ current, urls })
}

const openReasonModal = (id, status) => {
  reasonModal.visible = true
  reasonModal.id = id
  reasonModal.status = status
  reasonModal.mode = activeMode.value
  reasonModal.reason = ''
  reasonModal.note = ''
}

const closeReasonModal = () => {
  if (submittingReason.value) return
  reasonModal.visible = false
}

const buildHandleResult = (reason, note) => {
  const extra = String(note || '').trim()
  return extra ? `${reason}。补充说明：${extra}` : reason
}

const submitReasonAction = async () => {
  if (!reasonModal.reason) {
    uni.showToast({ title: '请选择处理原因', icon: 'none' })
    return
  }
  if (submittingReason.value) return
  submittingReason.value = true
  try {
    const handleResult = buildHandleResult(reasonModal.reason, reasonModal.note)
    if (reasonModal.mode === 'REPORT') {
      await handleAdminReport(reasonModal.id, reasonModal.status, handleResult)
    } else {
      await handleAdminAppeal(reasonModal.id, reasonModal.status, handleResult)
    }
    uni.showToast({ title: '处理成功', icon: 'success' })
    reasonModal.visible = false
    openedId.value = null
    load()
  } catch (error) {
    uni.showToast({ title: error.message || '处理失败', icon: 'none' })
  } finally {
    submittingReason.value = false
  }
}

const handleCurrent = async (id, status) => {
  const action = statusText(status)
  uni.showModal({
    title: '确认处理',
    content: `确定将该${activeMode.value === 'REPORT' ? '举报' : '申诉'}标记为${action}吗？`,
    success: async (res) => {
      if (!res.confirm) return
      try {
        if (activeMode.value === 'REPORT') {
          if (status === 'PROCESSING') {
            await markAdminReportProcessing(id)
          } else {
            await handleAdminReport(id, status, `举报${action}`)
          }
        } else {
          if (status === 'PROCESSING') {
            await markAdminAppealProcessing(id)
          } else {
            await handleAdminAppeal(id, status, `申诉${action}`)
          }
        }
        uni.showToast({ title: '处理成功', icon: 'success' })
        openedId.value = null
        load()
      } catch (error) {
        uni.showToast({ title: error.message || '处理失败', icon: 'none' })
      }
    }
  })
}
</script>

<style scoped>
.page {
  min-height: 100vh;
  padding: 28rpx 24rpx;
  background: #f5f6f8;
  box-sizing: border-box;
}

.page-title {
  margin: 18rpx 0 24rpx;
  font-size: 38rpx;
  font-weight: 700;
  color: #1f2933;
}

.tab-row {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 16rpx;
  margin-bottom: 22rpx;
}

.tab {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 10rpx;
  min-height: 82rpx;
  border-radius: 14rpx;
  background: #fff;
  color: #667085;
  font-size: 28rpx;
  font-weight: 700;
}

.tab.active {
  background: #17a84b;
  color: #fff;
}

.tab text {
  min-width: 34rpx;
  height: 34rpx;
  padding: 0 10rpx;
  border-radius: 999rpx;
  background: rgba(0, 0, 0, 0.08);
  font-size: 22rpx;
  line-height: 34rpx;
  text-align: center;
}

.summary {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  padding: 30rpx;
  border-radius: 18rpx;
  background: #fff;
}

.summary-number {
  font-size: 54rpx;
  line-height: 60rpx;
  font-weight: 700;
  color: #17a84b;
}

.summary-number.danger {
  color: #ef4444;
}

.summary-label {
  margin-top: 10rpx;
  font-size: 26rpx;
  color: #667085;
}

.summary-status {
  padding: 8rpx 16rpx;
  border-radius: 999rpx;
  font-size: 22rpx;
}

.summary-status.warning {
  background: #fee2e2;
  color: #ef4444;
}

.summary-status.safe {
  background: #dcfce7;
  color: #16a34a;
}

.filter-row {
  display: flex;
  gap: 14rpx;
  margin: 24rpx 0;
  overflow-x: auto;
}

.filter-chip {
  flex-shrink: 0;
  padding: 14rpx 22rpx;
  border-radius: 999rpx;
  background: #fff;
  color: #667085;
  font-size: 26rpx;
}

.filter-chip.active {
  background: #17a84b;
  color: #fff;
  font-weight: 700;
}


.case-card,
.empty {
  padding: 28rpx;
  border-radius: 16rpx;
  background: #fff;
  margin-bottom: 18rpx;
}

.card-head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 20rpx;
}

.card-main {
  flex: 1;
  min-width: 0;
}

.case-title {
  font-size: 32rpx;
  font-weight: 700;
  color: #1f2933;
}

.case-sub {
  margin-top: 8rpx;
  font-size: 24rpx;
  color: #8a96a8;
}

.status-badge {
  flex-shrink: 0;
  padding: 8rpx 14rpx;
  border-radius: 999rpx;
  background: #eef2f6;
  color: #667085;
  font-size: 22rpx;
}

.badge-row {
  display: flex;
  flex-shrink: 0;
  align-items: center;
  gap: 8rpx;
  flex-wrap: wrap;
  justify-content: flex-end;
}

.risk-badge {
  flex-shrink: 0;
  padding: 8rpx 14rpx;
  border-radius: 999rpx;
  background: #fff1f2;
  color: #be123c;
  font-size: 22rpx;
}

.case-desc,
.empty {
  margin-top: 18rpx;
  color: #667085;
  font-size: 26rpx;
  line-height: 38rpx;
}

.meta-row {
  display: flex;
  justify-content: space-between;
  gap: 20rpx;
  margin-top: 20rpx;
  font-size: 22rpx;
  color: #98a2b3;
}

.detail-panel {
  margin-top: 22rpx;
  padding-top: 20rpx;
  border-top: 1rpx solid #eef0f3;
}

.target-product {
  display: flex;
  align-items: center;
  gap: 18rpx;
  padding: 18rpx;
  margin-bottom: 18rpx;
  border-radius: 14rpx;
  background: #f8fafc;
}

.target-image {
  width: 116rpx;
  height: 116rpx;
  border-radius: 12rpx;
  background: #eef2f6;
  flex-shrink: 0;
}

.placeholder {
  display: flex;
  align-items: center;
  justify-content: center;
  color: #98a2b3;
  font-size: 22rpx;
}

.target-main {
  flex: 1;
  min-width: 0;
}

.target-title {
  font-size: 28rpx;
  font-weight: 700;
  color: #1f2933;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.target-meta {
  margin-top: 8rpx;
  font-size: 24rpx;
  color: #667085;
}

.target-arrow {
  color: #b2bdca;
  font-size: 40rpx;
}

.detail-row {
  display: flex;
  justify-content: space-between;
  gap: 24rpx;
  padding: 12rpx 0;
  font-size: 24rpx;
  color: #667085;
}

.detail-row text:last-child {
  flex: 1;
  color: #1f2933;
  text-align: right;
  word-break: break-all;
}

.image-grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 12rpx;
  margin-top: 16rpx;
}

.case-image {
  width: 100%;
  height: 150rpx;
  border-radius: 12rpx;
  background: #eef2f6;
}

.actions {
  margin-top: 22rpx;
  display: flex;
  gap: 14rpx;
  flex-wrap: wrap;
}

.actions button {
  margin: 0;
}

.pass {
  background: #17a84b;
  color: #fff;
}

.process {
  background: #eff6ff;
  color: #2563eb;
}

.reject {
  background: #fff1f2;
  color: #ef4444;
}

.close {
  background: #f8fafc;
  color: #667085;
}

.handle-result {
  margin-top: 18rpx;
  padding: 16rpx;
  border-radius: 12rpx;
  background: #f8fafc;
  color: #667085;
  font-size: 24rpx;
}

.modal-mask {
  position: fixed;
  left: 0;
  right: 0;
  top: 0;
  bottom: 0;
  z-index: 99;
  display: flex;
  align-items: flex-end;
  background: rgba(15, 23, 42, 0.42);
}

.reason-sheet {
  width: 100%;
  padding: 30rpx 28rpx 36rpx;
  border-radius: 28rpx 28rpx 0 0;
  background: #fff;
  box-sizing: border-box;
}

.sheet-head {
  display: flex;
  justify-content: space-between;
  gap: 24rpx;
}

.sheet-title {
  font-size: 34rpx;
  font-weight: 700;
  color: #1f2933;
}

.sheet-subtitle {
  margin-top: 8rpx;
  font-size: 24rpx;
  color: #8a96a8;
}

.sheet-close {
  flex-shrink: 0;
  width: 54rpx;
  height: 54rpx;
  border-radius: 50%;
  background: #f1f5f9;
  color: #667085;
  font-size: 42rpx;
  line-height: 50rpx;
  text-align: center;
}

.reason-list {
  display: flex;
  flex-wrap: wrap;
  gap: 14rpx;
  margin-top: 26rpx;
}

.reason-chip {
  padding: 14rpx 18rpx;
  border-radius: 999rpx;
  background: #f8fafc;
  color: #475467;
  font-size: 24rpx;
  border: 2rpx solid transparent;
}

.reason-chip.active {
  border-color: #17a84b;
  background: #f0fdf4;
  color: #16a34a;
  font-weight: 700;
}

.reason-input {
  width: 100%;
  min-height: 150rpx;
  margin-top: 24rpx;
  padding: 20rpx;
  border-radius: 16rpx;
  background: #f8fafc;
  font-size: 26rpx;
  color: #1f2933;
  box-sizing: border-box;
}

.reason-placeholder {
  color: #98a2b3;
}

.sheet-actions {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 18rpx;
  margin-top: 24rpx;
}

.sheet-button {
  height: 84rpx;
  line-height: 84rpx;
  border-radius: 14rpx;
  font-size: 28rpx;
}

.sheet-button.cancel {
  background: #f8fafc;
  color: #667085;
}

.sheet-button.confirm {
  background: #17a84b;
  color: #fff;
}
</style>
