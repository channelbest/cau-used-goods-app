<template>
  <view class="page">
    <scroll-view class="filter-row" scroll-x enhanced :show-scrollbar="false">
      <view class="filter-inner">
        <view
          v-for="item in logFilters"
          :key="item.value"
          :class="['filter-chip', activeFilter === item.value ? 'active' : '']"
          @tap.stop="selectFilter(item.value)"
        >
          {{ item.label }} {{ countByFilter(item.value) }}
        </view>
      </view>
    </scroll-view>

    <view v-if="filteredLogs.length === 0" class="empty">暂无日志</view>
    <view v-for="log in filteredLogs" :key="log.id" class="log-item" @click="showLogDetail(log)">
      <view class="log-main">
        <view>
          <view class="name">{{ operationLabel(log.operationType, log) }}</view>
          <view class="desc">{{ formatDateTime(log.createTime) }}</view>
        </view>
        <text class="arrow">›</text>
      </view>
      <view v-if="formatDescription(log)" class="detail">{{ formatDescription(log) }}</view>
      <view v-if="sourceText(log)" class="source">{{ sourceText(log) }}</view>
    </view>

    <!-- 日志详情弹窗 -->
    <view v-if="detailVisible" class="detail-modal" @click="closeDetail">
      <view class="detail-content" @click.stop>
        <view class="detail-header">
          <text class="detail-title">操作详情</text>
          <text class="detail-close" @click="closeDetail">×</text>
        </view>
        <view v-if="currentLog" class="detail-body">
          <view class="detail-row">
            <text class="detail-label">操作类型</text>
            <text class="detail-value">{{ operationLabel(currentLog.operationType, currentLog) }}</text>
          </view>
          <view class="detail-row">
            <text class="detail-label">操作人</text>
            <view class="detail-value">
              <image v-if="currentLog.adminAvatar" class="admin-avatar" :src="currentLog.adminAvatar" mode="aspectFill" />
              <text>{{ currentLog.adminName || `管理员 #${currentLog.adminId}` }}</text>
            </view>
          </view>
          <view class="detail-row">
            <text class="detail-label">操作对象</text>
            <text class="detail-value">{{ targetDetail || `${targetLabelMap[currentLog.targetType] || currentLog.targetType} #${currentLog.targetId}` }}</text>
          </view>
          <view v-if="currentLog.targetType === 'USER' && currentLog.targetName" class="detail-row">
            <text class="detail-label">被操作者</text>
            <text class="detail-value">{{ currentLog.targetName }}（ID: {{ currentLog.targetId }}，{{ currentLog.targetCollege || '未知学院' }}）</text>
          </view>
          <view v-if="currentLog.targetType === 'PRODUCT' && currentLog.targetName" class="detail-row">
            <text class="detail-label">商品名称</text>
            <text class="detail-value">{{ currentLog.targetName }}</text>
          </view>
          <view v-if="currentLog.description" class="detail-row">
            <text class="detail-label">操作说明</text>
            <text class="detail-value">{{ formatDescription(currentLog) }}</text>
          </view>
          <view v-if="currentLog.relatedType && currentLog.relatedId" class="detail-row">
            <text class="detail-label">关联来源</text>
            <text class="detail-value">{{ targetLabelMap[currentLog.relatedType] || currentLog.relatedType }} #{{ currentLog.relatedId }}</text>
          </view>
          <view class="detail-row">
            <text class="detail-label">操作时间</text>
            <text class="detail-value">{{ formatDateTime(currentLog.createTime) }}</text>
          </view>
        </view>
      </view>
    </view>
  </view>
</template>

<script setup>
import { computed, ref } from 'vue'
import { onShow } from '@dcloudio/uni-app'
import { getAdminLogs, getAdminLogDetail, getAdminUserDetail, getAdminReportDetail, getAdminAppealDetail } from '../../api/admin'
import { getProductById } from '../../api/product'

const logs = ref([])
const activeFilter = ref('ALL')
const detailVisible = ref(false)
const currentLog = ref(null)
const targetDetail = ref('')

const logFilters = [
  { label: '全部', value: 'ALL' },
  { label: '商品', value: 'PRODUCT' },
  { label: '用户', value: 'USER' },
  { label: '订单', value: 'ORDER' },
  { label: '举报', value: 'REPORT' },
  { label: '申诉', value: 'APPEAL' },
  { label: '公告', value: 'NOTICE' },
  { label: '敏感词', value: 'WORD' },
  { label: '标签', value: 'CATEGORY' }
]

const operationMap = {
  USER_DISABLE: '禁用用户',
  USER_ENABLE: '启用用户',
  USER_BAN: '永久封禁用户',
  USER_UNBAN: '解除封禁用户',
  USER_ROLE_CHANGE: '修改用户角色',
  PRODUCT_OFF_SHELF: '下架商品',
  PRODUCT_ON_SALE: '上架商品',
  REPORT_RESOLVE: '处理举报',
  REPORT_APPROVE: '通过举报',
  APPROVE_REPORT: '通过举报',
  REPORT_HANDLE: '处理举报',
  MARK_REPORT_PROCESSING: '举报标记为处理中',
  REPORT_REJECT: '驳回举报',
  REJECT_REPORT: '驳回举报',
  REPORT_CLOSE: '关闭举报',
  HANDLE_APPEAL: '处理申诉',
  MARK_APPEAL_PROCESSING: '申诉标记为处理中',
  APPROVE_APPEAL: '通过申诉',
  REJECT_APPEAL: '驳回申诉',
  NOTICE_PUBLISH: '发布公告',
  NOTICE_OFFLINE: '下线公告',
  STATUS_NOTICE: '更新公告状态',
  CREATE_NOTICE: '新增公告',
  UPDATE_NOTICE: '编辑公告',
  DELETE_NOTICE: '删除公告',
  WORD_CREATE: '新增敏感词',
  WORD_DISABLE: '禁用敏感词',
  CREATE_WORD: '新增敏感词',
  UPDATE_WORD: '编辑敏感词',
  DELETE_WORD: '禁用敏感词',
  CATEGORY_CREATE: '新增标签',
  CATEGORY_UPDATE: '编辑标签',
  CATEGORY_ENABLE: '启用标签',
  CATEGORY_DISABLE: '停用标签',
  ORDER_EXCEPTION_CLOSE: '异常关闭订单',
  UPDATE_ORDER_STATUS: '修改订单状态',
  STUDENT_VERIFY_APPROVE: '通过学生认证',
  STUDENT_VERIFY_REJECT: '驳回学生认证'
}

const statusLabelMap = {
  DRAFT: '草稿',
  PUBLISHED: '已发布',
  OFFLINE: '已下线',
  PROCESSING: '处理中',
  APPROVED: '已通过',
  RESOLVED: '已处理',
  REJECTED: '已驳回',
  CLOSED: '已关闭',
  ENABLED: '已启用',
  DISABLED: '已停用',
  NORMAL: '正常',
  HIDDEN: '已隐藏',
  DELETED: '已删除'
}

const targetLabelMap = {
  USER: '用户',
  PRODUCT: '商品',
  REPORT: '举报',
  APPEAL: '申诉',
  NOTICE: '公告',
  WORD: '敏感词',
  CATEGORY: '标签',
  ORDER: '订单'
}

const descriptionMap = [
  [/^update announcement status:\s*(.+)$/i, (_, status) => announcementStatusDescription(status)],
  [/^create announcement:\s*(.+)$/i, (_, title) => `新增公告：${title}`],
  [/^disable sensitive word by delete operation$/i, () => '禁用敏感词'],
  [/^create sensitive word:\s*(.+)$/i, (_, word) => `新增敏感词：${word}`],
  [/^update sensitive word:\s*(.+)$/i, (_, word) => `编辑敏感词：${word}`],
  [/^mark report as processing$/i, () => '举报标记为处理中'],
  [/^handle report.*APPROVED$/i, () => '举报已通过'],
  [/^handle report.*REJECTED$/i, () => '举报已驳回'],
  [/^mark appeal #?(\d+)? as PROCESSING$/i, (_, id) => `申诉${id ? ` #${id}` : ''}标记为处理中`],
  [/^handle appeal #?(\d+)?:\s*APPROVED$/i, (_, id) => `申诉 #${id}已通过`],
  [/^handle appeal #?(\d+)?:\s*REJECTED$/i, (_, id) => `申诉 #${id}已驳回`]
]

const logCategory = (log) => {
  if (!log) return ''
  if (['STUDENT_VERIFY_APPROVE', 'STUDENT_VERIFY_REJECT', 'USER_DISABLE', 'USER_ENABLE', 'USER_BAN', 'USER_UNBAN', 'USER_ROLE_CHANGE'].includes(log.operationType)) return 'USER'
  if (log.operationType?.includes('NOTICE')) return 'NOTICE'
  if (log.operationType?.includes('WORD')) return 'WORD'
  if (log.operationType?.includes('CATEGORY')) return 'CATEGORY'
  if (log.operationType?.includes('ORDER')) return 'ORDER'
  return log.targetType || ''
}

const filteredLogs = computed(() => {
  if (activeFilter.value === 'ALL') return logs.value
  return logs.value.filter((log) => logCategory(log) === activeFilter.value)
})

const countByFilter = (value) => {
  if (value === 'ALL') return logs.value.length
  return logs.value.filter((log) => logCategory(log) === value).length
}

const selectFilter = (value) => {
  activeFilter.value = value
}

const load = async () => {
  try {
    const result = await getAdminLogs()
    logs.value = result?.items || []
  } catch (error) {
    uni.showToast({ title: error.message || '日志加载失败', icon: 'none' })
  }
}

const showLogDetail = async (log) => {
  if (!log) return
  currentLog.value = log
  targetDetail.value = ''

  // 根据操作对象类型获取详细信息
  if (log.targetType === 'USER' && log.targetId) {
    try {
      const result = await getAdminUserDetail(log.targetId)
      const data = result?.data || result
      const user = data?.user || data
      if (user && user.nickname) {
        targetDetail.value = `${user.nickname}（ID: ${log.targetId}，${user.college || '未知学院'}）`
      } else {
        targetDetail.value = `用户 #${log.targetId}（已删除或不存在）`
      }
    } catch (e) {
      console.error('获取用户详情失败', e)
      targetDetail.value = `用户 #${log.targetId}（已删除或不存在）`
    }
  } else if (log.targetType === 'PRODUCT' && log.targetId) {
    try {
      const result = await getProductById(log.targetId)
      const data = result?.data || result
      if (data && data.title) {
        targetDetail.value = `${data.title}（ID: ${log.targetId}）`
      } else {
        targetDetail.value = `商品 #${log.targetId}（已删除或不存在）`
      }
    } catch (e) {
      console.error('获取商品详情失败', e)
      targetDetail.value = `商品 #${log.targetId}（已删除或不存在）`
    }
  } else if (log.targetType === 'APPEAL' && log.targetId) {
    try {
      const result = await getAdminAppealDetail(log.targetId)
      const data = result?.data || result
      const appeal = data?.appeal || data
      if (appeal) {
        targetDetail.value = `申诉 #${log.targetId} - ${appeal.reason || '查看详情'}`
      }
    } catch (e) {
      console.error('获取申诉详情失败', e)
    }
  } else if (log.targetType === 'REPORT' && log.targetId) {
    try {
      const result = await getAdminReportDetail(log.targetId)
      const data = result?.data || result
      const report = data?.report || data
      if (report) {
        targetDetail.value = `举报 #${log.targetId} - ${report.reason || '查看详情'}`
      }
    } catch (e) {
      console.error('获取举报详情失败', e)
    }
  } else if (log.targetType === 'ORDER' && log.targetId) {
    targetDetail.value = `订单 #${log.targetId}`
  } else if (log.targetType === 'NOTICE' && log.targetId) {
    targetDetail.value = `公告 #${log.targetId}`
  }

  detailVisible.value = true
}

const closeDetail = () => {
  detailVisible.value = false
  currentLog.value = null
}

const operationLabel = (value, log = null) => {
  if (value === 'STATUS_NOTICE') {
    const raw = String(log?.description || '')
    if (raw.includes('PUBLISHED') || raw.includes('已发布')) return '发布公告'
    if (raw.includes('OFFLINE') || raw.includes('已下线')) return '下线公告'
    if (raw.includes('DRAFT') || raw.includes('草稿')) return '转为草稿'
  }
  return operationMap[value] || translateOperation(value)
}

const translateOperation = (value) => {
  if (!value) return '操作'
  return String(value)
    .replaceAll('PROCESSING', '处理中')
    .replaceAll('APPROVED', '已通过')
    .replaceAll('APPROVE', '通过')
    .replaceAll('REJECTED', '已驳回')
    .replaceAll('REJECT', '驳回')
    .replaceAll('MARK', '标记')
    .replaceAll('HANDLE', '处理')
    .replaceAll('STATUS', '更新状态')
    .replaceAll('CREATE', '新增')
    .replaceAll('UPDATE', '编辑')
    .replaceAll('DELETE', '删除')
    .replaceAll('NOTICE', '公告')
    .replaceAll('WORD', '敏感词')
    .replaceAll('REPORT', '举报')
    .replaceAll('APPEAL', '申诉')
    .replaceAll('PRODUCT', '商品')
    .replaceAll('USER', '用户')
    .replaceAll('ORDER', '订单')
    .replaceAll('_', '')
}

const translateStatus = (status) => statusLabelMap[String(status || '').trim()] || status || ''

const replaceKnownWords = (value) => {
  let text = String(value || '')
  text = text.replaceAll('auto exception close by account status change', '因用户账号状态变更，系统自动异常关闭订单')
  text = text.replaceAll('responsibleParty', '责任方')
  text = text.replaceAll('reason', '原因')
  text = text.replaceAll('BUYER', '买家')
  text = text.replaceAll('SELLER', '卖家')
  Object.keys(statusLabelMap).forEach((key) => {
    text = text.replaceAll(key, statusLabelMap[key])
  })
  Object.keys(targetLabelMap).forEach((key) => {
    text = text.replaceAll(key, targetLabelMap[key])
  })
  return text
}

const formatDescription = (log) => {
  const raw = String(log?.description || '').trim()
  const target = targetLabelMap[log?.targetType] || '对象'
  const idText = log?.targetId ? ` #${log.targetId}` : ''

  for (const [pattern, formatter] of descriptionMap) {
    const matched = raw.match(pattern)
    if (matched) return formatter(...matched)
  }

  const replaced = replaceKnownWords(raw)
  if (replaced && !/[A-Za-z_]/.test(replaced)) return replaced

  switch (log?.operationType) {
    case 'STATUS_NOTICE':
      return replaced || `公告${idText}状态已更新`
    case 'REPORT_HANDLE':
    case 'MARK_REPORT_PROCESSING':
      return `举报${idText}标记为处理中`
    case 'REPORT_APPROVE':
    case 'APPROVE_REPORT':
      return `举报${idText}已通过`
    case 'REPORT_REJECT':
    case 'REJECT_REPORT':
      return `举报${idText}已驳回`
    case 'REPORT_CLOSE':
      return `举报${idText}已关闭`
    case 'ORDER_EXCEPTION_CLOSE':
      return replaced || `订单${idText}已异常关闭`
    case 'HANDLE_APPEAL':
    case 'MARK_APPEAL_PROCESSING':
      return replaced || `申诉${idText}已处理`
    case 'APPROVE_APPEAL':
      return `申诉${idText}已通过`
    case 'REJECT_APPEAL':
      return `申诉${idText}已驳回`
    case 'USER_DISABLE':
      return `用户${idText}已被禁用`
    case 'USER_ENABLE':
      return `用户${idText}已恢复启用`
    case 'USER_BAN':
      return `用户${idText}已被永久封禁`
    case 'USER_UNBAN':
      return `用户${idText}已解除封禁`
    case 'STUDENT_VERIFY_APPROVE':
      return `学生认证${idText}已通过`
    case 'STUDENT_VERIFY_REJECT':
      return `学生认证${idText}已驳回`
    case 'UPDATE_PRODUCT_STATUS':
      return `商品${idText}状态已更新`
    case 'UPDATE_ORDER_STATUS':
      return `订单${idText}状态已更新`
    default:
      return replaced || `${target}${idText}已处理`
  }
}

const sourceText = (log) => {
  if (!log?.relatedType || !log?.relatedId) return ''
  const typeText = targetLabelMap[log.relatedType] || log.relatedType
  return `来源：${typeText} #${log.relatedId}`
}

const pad = (value) => String(value).padStart(2, '0')

const formatDateTime = (value) => {
  if (!value) return ''
  const date = new Date(value)
  if (!Number.isNaN(date.getTime())) {
    return `${date.getFullYear()}-${pad(date.getMonth() + 1)}-${pad(date.getDate())} ${pad(date.getHours())}:${pad(date.getMinutes())}:${pad(date.getSeconds())}`
  }
  return String(value).replace('T', ' ').replace(/\+\d{2}:\d{2}$/, '')
}

onShow(load)
</script>

<style scoped>
.page {
  min-height: 100vh;
  padding: 24rpx;
  background: #f5f6f8;
  box-sizing: border-box;
}

.filter-row {
  width: 100%;
  margin: 4rpx 0 22rpx;
  white-space: nowrap;
}

.filter-inner {
  display: inline-flex;
  gap: 14rpx;
  min-width: 100%;
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

.log-item,
.empty {
  padding: 28rpx;
  border-radius: 16rpx;
  background: #fff;
  margin-bottom: 18rpx;
}

.log-main {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 24rpx;
}

.name {
  font-size: 30rpx;
  font-weight: 700;
  color: #1f2933;
}

.desc,
.empty,
.detail {
  margin-top: 8rpx;
  color: #667085;
  font-size: 26rpx;
}

.detail {
  line-height: 38rpx;
  color: #8a96a8;
}

.source {
  margin-top: 8rpx;
  font-size: 24rpx;
  color: #17a84b;
}

.arrow {
  color: #b2bdca;
  font-size: 42rpx;
  line-height: 42rpx;
}

/* 详情弹窗 */
.detail-modal {
  position: fixed;
  top: 0;
  right: 0;
  bottom: 0;
  left: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0, 0, 0, 0.5);
  z-index: 1000;
}

.detail-content {
  width: 80%;
  max-width: 600rpx;
  border-radius: 24rpx;
  background: #fff;
  overflow: hidden;
}

.detail-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 28rpx 32rpx;
  border-bottom: 1rpx solid #f0f0f0;
}

.detail-title {
  font-size: 32rpx;
  font-weight: 700;
  color: #1f2933;
}

.detail-close {
  font-size: 44rpx;
  color: #999;
  line-height: 1;
}

.detail-body {
  padding: 24rpx 32rpx;
}

.detail-row {
  display: flex;
  padding: 16rpx 0;
  border-bottom: 1rpx solid #f5f5f5;
}

.detail-row:last-child {
  border-bottom: none;
}

.detail-label {
  width: 160rpx;
  flex-shrink: 0;
  color: #667085;
  font-size: 26rpx;
}

.detail-value {
  flex: 1;
  color: #1f2933;
  font-size: 26rpx;
  word-break: break-word;
  display: flex;
  align-items: center;
  gap: 12rpx;
}

.admin-avatar {
  width: 48rpx;
  height: 48rpx;
  border-radius: 50%;
  background: #e8ecef;
}
</style>
